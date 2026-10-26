package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorElementos es el puerto para buscar un elemento del backlog por su ID (HU-18).
type LectorElementos interface {
	// BuscarElementoPorID devuelve domain.ErrElementoNoEncontrado si no existe.
	BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error)
}

// RepositorioEsfuerzo es el puerto para guardar los registros de esfuerzo (HU-18).
type RepositorioEsfuerzo interface {
	// GuardarRegistro persiste el registro y le asigna su ID (RN-18.10).
	GuardarRegistro(ctx context.Context, r *domain.RegistroEsfuerzo) error
}

// DatosRegistroEsfuerzo son las entradas de HU-18.
type DatosRegistroEsfuerzo struct {
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
}

// RegistrarEsfuerzo es el caso de uso de HU-18.
type RegistrarEsfuerzo struct {
	proyectos LectorProyectos
	elementos LectorElementos
	registros RepositorioEsfuerzo
	reloj     func() time.Time
}

// NuevoRegistrarEsfuerzo recibe el reloj para que "hoy" se pueda fijar en los tests (spec, sección 5).
func NuevoRegistrarEsfuerzo(proyectos LectorProyectos, elementos LectorElementos,
	registros RepositorioEsfuerzo, reloj func() time.Time) *RegistrarEsfuerzo {
	return &RegistrarEsfuerzo{proyectos: proyectos, elementos: elementos, registros: registros, reloj: reloj}
}

// Ejecutar busca el proyecto y el elemento, valida el registro en el dominio y lo guarda.
// Si algo falla, no se guarda nada.
func (uc *RegistrarEsfuerzo) Ejecutar(ctx context.Context, d DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error) {
	p, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}

	e, err := uc.elementos.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("buscar elemento %d: %w", d.ElementoID, err)
	}
	if e.ProyectoID != p.ID { // RN-18.1: un elemento de otro proyecto se trata como inexistente
		return domain.RegistroEsfuerzo{}, fmt.Errorf("elemento %d en el proyecto %d: %w", d.ElementoID, p.ID, domain.ErrElementoNoEncontrado)
	}

	r, err := domain.NuevoRegistroEsfuerzo(p.ID, e.ID, d.Integrante, d.Fecha, d.Actividad, d.Horas, uc.reloj())
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("registrar esfuerzo: %w", err)
	}

	if err := uc.registros.GuardarRegistro(ctx, &r); err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("guardar registro de esfuerzo: %w", err)
	}
	return r, nil
}
