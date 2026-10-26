package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioSprints es el puerto de persistencia de los sprints (HU-08).
type RepositorioSprints interface {
	// ListarSprints devuelve los sprints del proyecto ordenados por número.
	ListarSprints(ctx context.Context, proyectoID int64) ([]domain.Sprint, error)
	// GuardarSprint persiste el sprint y le asigna su ID y su número dentro del proyecto. Vuelve a
	// verificar la superposición al guardar, por si otro pedido creó un sprint en el medio
	// (spec HU-08, sección 5): en ese caso devuelve domain.ErrSprintSuperpuesto.
	GuardarSprint(ctx context.Context, s *domain.Sprint) error
}

// DatosNuevoSprint son las entradas de HU-08. Una fecha cero es "sin fecha".
type DatosNuevoSprint struct {
	ProyectoID  int64
	FechaInicio time.Time
	FechaFin    time.Time
	Objetivo    string
}

// CrearSprint es el caso de uso de HU-08. Implementa HU-08 (RN-08.1, RN-08.7, RN-08.9).
type CrearSprint struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
}

func NuevoCrearSprint(proyectos LectorProyectos, sprints RepositorioSprints) *CrearSprint {
	return &CrearSprint{proyectos: proyectos, sprints: sprints}
}

// Ejecutar busca el proyecto y sus sprints, valida el sprint nuevo en el dominio y lo guarda.
// Si algo falla, no se guarda nada.
func (uc *CrearSprint) Ejecutar(ctx context.Context, d DatosNuevoSprint) (domain.Sprint, error) {
	p, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	existentes, err := uc.sprints.ListarSprints(ctx, p.ID)
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("listar sprints del proyecto: %w", err)
	}

	s, err := domain.NuevoSprint(p, d.FechaInicio, d.FechaFin, d.Objetivo, existentes)
	if err != nil {
		return domain.Sprint{}, err
	}
	if err := uc.sprints.GuardarSprint(ctx, &s); err != nil {
		return domain.Sprint{}, fmt.Errorf("guardar sprint: %w", err)
	}
	return s, nil
}

// ListarSprints devuelve los sprints de un proyecto. Implementa HU-08 (CA-08.9).
type ListarSprints struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
}

func NuevoListarSprints(proyectos LectorProyectos, sprints RepositorioSprints) *ListarSprints {
	return &ListarSprints{proyectos: proyectos, sprints: sprints}
}

// Ejecutar devuelve los sprints del proyecto ordenados por número, o domain.ErrProyectoNoEncontrado
// si el proyecto no existe (así no se confunde con un proyecto sin sprints).
func (uc *ListarSprints) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return nil, err
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto: %w", err)
	}
	return sprints, nil
}
