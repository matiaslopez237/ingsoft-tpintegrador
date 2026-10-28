package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioCompletar es el puerto de persistencia que necesita CompletarHistoria (HU-10).
type RepositorioCompletar interface {
	LectorElementos // BuscarElementoPorID (de HU-18): devuelve domain.ErrElementoNoEncontrado si no existe.
	// BuscarSprint es el mismo que en RepositorioAsignacion (HU-09): el proyecto del sprint y si
	// está Abierto. Devuelve domain.ErrSprintNoEncontrado si no existe.
	BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error)
	// GuardarCompletado persiste el estado del elemento ya existente.
	GuardarCompletado(ctx context.Context, e *domain.ElementoBacklog) error
}

// DatosCompletar son las entradas para marcar un elemento como completado (HU-10).
type DatosCompletar struct {
	ProyectoID int64 // de la ruta /proyectos/{id}/backlog/{elementoID}/completar
	ElementoID int64
}

// CompletarHistoria es el caso de uso de HU-10 para marcar un elemento del sprint como Hecho.
type CompletarHistoria struct {
	proyectos LectorProyectos
	repo      RepositorioCompletar
}

func NuevoCompletarHistoria(proyectos LectorProyectos, repo RepositorioCompletar) *CompletarHistoria {
	return &CompletarHistoria{proyectos: proyectos, repo: repo}
}

// Ejecutar busca el proyecto y el elemento, consulta si el sprint del elemento (si tiene) está
// abierto, aplica la regla de dominio y persiste el resultado. Los errores salen en el orden de
// la sección 7 de la spec; si algo falla, no se guarda nada.
// Implementa HU-10 (CA-10.1 a CA-10.5, RN-10.4).
func (uc *CompletarHistoria) Ejecutar(ctx context.Context, d DatosCompletar) error {
	if _, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID); err != nil {
		return fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}

	e, err := uc.repo.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return fmt.Errorf("buscar elemento %d: %w", d.ElementoID, err)
	}
	if e.ProyectoID != d.ProyectoID { // RN-10.4, igual que RN-18.1: de otro proyecto es "inexistente"
		return fmt.Errorf("elemento %d en el proyecto %d: %w", d.ElementoID, d.ProyectoID, domain.ErrElementoNoEncontrado)
	}

	sprintAbierto, err := sprintDelElementoAbierto(ctx, uc.repo, e)
	if err != nil {
		return fmt.Errorf("buscar sprint %d: %w", *e.SprintID, err)
	}

	if err := e.Completar(sprintAbierto); err != nil {
		return err
	}

	if err := uc.repo.GuardarCompletado(ctx, &e); err != nil {
		return fmt.Errorf("guardar completado: %w", err)
	}
	return nil
}
