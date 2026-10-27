package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioDefectos es el puerto para guardar los defectos (HU-20).
type RepositorioDefectos interface {
	// GuardarDefecto persiste el defecto y le asigna su ID (RN-20.6).
	GuardarDefecto(ctx context.Context, d *domain.Defecto) error
}

// DatosDefecto son las entradas de HU-20 (spec, sección 2).
type DatosDefecto struct {
	ProyectoID  int64
	ElementoID  int64
	Descripcion string
	Severidad   domain.Severidad
}

// RegistrarDefecto es el caso de uso de HU-20. Reutiliza LectorProyectos (HU-04), LectorElementos
// (HU-18) y RepositorioSprints (HU-08).
type RegistrarDefecto struct {
	proyectos LectorProyectos
	elementos LectorElementos
	sprints   RepositorioSprints
	defectos  RepositorioDefectos
	reloj     func() time.Time
}

// NuevoRegistrarDefecto recibe el reloj para que "hoy" se pueda fijar en los tests, igual que HU-18.
func NuevoRegistrarDefecto(proyectos LectorProyectos, elementos LectorElementos, sprints RepositorioSprints,
	defectos RepositorioDefectos, reloj func() time.Time) *RegistrarDefecto {
	return &RegistrarDefecto{proyectos: proyectos, elementos: elementos, sprints: sprints, defectos: defectos, reloj: reloj}
}

// Ejecutar busca el proyecto y el elemento, valida el defecto en el dominio, le asigna el Sprint
// de detección vigente (si hay uno) y lo guarda. Si algo falla, no se guarda nada.
// Implementa HU-20 (RN-20.1, RN-20.5, RN-20.6).
func (uc *RegistrarDefecto) Ejecutar(ctx context.Context, d DatosDefecto) (domain.Defecto, error) {
	p, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}

	e, err := uc.elementos.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("buscar elemento %d: %w", d.ElementoID, err)
	}
	if e.ProyectoID != p.ID { // RN-20.1: un elemento de otro proyecto se trata como inexistente
		return domain.Defecto{}, fmt.Errorf("elemento %d en el proyecto %d: %w", d.ElementoID, p.ID, domain.ErrElementoNoEncontrado)
	}

	defecto, err := domain.NuevoDefecto(p.ID, e.ID, d.Descripcion, d.Severidad)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("registrar defecto: %w", err)
	}

	sprints, err := uc.sprints.ListarSprints(ctx, p.ID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("listar sprints del proyecto %d: %w", p.ID, err)
	}
	if id, ok := domain.SprintVigente(sprints, uc.reloj()); ok {
		defecto.SprintDeteccionID = &id
	}

	if err := uc.defectos.GuardarDefecto(ctx, &defecto); err != nil {
		return domain.Defecto{}, fmt.Errorf("guardar defecto: %w", err)
	}
	return defecto, nil
}
