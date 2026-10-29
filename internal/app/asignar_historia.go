package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioAsignacion es el puerto de persistencia que necesitan AsignarHistoria y DesasignarHistoria.
type RepositorioAsignacion interface {
	LectorElementos // BuscarElementoPorID (de HU-18): devuelve domain.ErrElementoNoEncontrado si no existe.
	// BuscarSprint devuelve el proyecto al que pertenece el sprint y si está Abierto, sin depender
	// de domain.Sprint (de HU-08). Devuelve domain.ErrSprintNoEncontrado si no existe.
	BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error)
	// GuardarAsignacion persiste el SprintID del elemento ya existente.
	GuardarAsignacion(ctx context.Context, e *domain.ElementoBacklog) error
}

// DatosAsignacion son las entradas para asignar un elemento a un sprint (HU-09).
type DatosAsignacion struct {
	ProyectoID int64 // de la ruta /proyectos/{id}/backlog/{elementoID}/sprint
	ElementoID int64
	SprintID   int64
}

// AsignarHistoria es el caso de uso de HU-09 para asignar un elemento del backlog a un sprint.
type AsignarHistoria struct {
	repo RepositorioAsignacion
}

func NuevoAsignarHistoria(repo RepositorioAsignacion) *AsignarHistoria {
	return &AsignarHistoria{repo: repo}
}

// Ejecutar busca el elemento y el sprint destino, calcula si el sprint que el elemento ya tenía
// (si tenía) sigue activo, aplica la regla de dominio y persiste el resultado. Si el elemento o
// el sprint no existen, o los datos son inválidos, no se guarda nada.
// Implementa HU-09 (CA-09.1, CA-09.2, CA-09.3, CA-09.7, CA-09.8, RN-09.7).
func (uc *AsignarHistoria) Ejecutar(ctx context.Context, d DatosAsignacion) error {
	e, err := uc.repo.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return err
	}
	if e.ProyectoID != d.ProyectoID { // RN-09.7, igual que RN-18.1: de otro proyecto es "inexistente"
		return domain.ErrElementoNoEncontrado
	}

	destinoProyectoID, destinoAbierto, err := uc.repo.BuscarSprint(ctx, d.SprintID)
	if err != nil {
		return err
	}

	sprintActualActivo, err := uc.sprintActualActivo(ctx, e, d.SprintID)
	if err != nil {
		return err
	}

	if err := e.AsignarASprint(d.SprintID, destinoProyectoID, destinoAbierto, sprintActualActivo); err != nil {
		return err
	}

	if err := uc.repo.GuardarAsignacion(ctx, &e); err != nil {
		return fmt.Errorf("guardar asignación: %w", err)
	}
	return nil
}

// sprintActualActivo consulta si el sprint que el elemento ya tiene asignado (si tiene, y es
// distinto del destino) sigue activo. Si no tiene sprint asignado, no hace falta consultar nada.
func (uc *AsignarHistoria) sprintActualActivo(ctx context.Context, e domain.ElementoBacklog, sprintDestinoID int64) (bool, error) {
	if e.SprintID == nil || *e.SprintID == sprintDestinoID {
		return false, nil
	}
	_, abierto, err := uc.repo.BuscarSprint(ctx, *e.SprintID)
	if err != nil {
		return false, err
	}
	return abierto, nil
}

// DatosDesasignacion son las entradas para desasignar un elemento de su sprint (HU-09).
type DatosDesasignacion struct {
	ProyectoID int64 // de la ruta /proyectos/{id}/backlog/{elementoID}/sprint
	ElementoID int64
}

// DesasignarHistoria es el caso de uso de HU-09 para sacar un elemento del backlog de su sprint.
type DesasignarHistoria struct {
	repo RepositorioAsignacion
}

func NuevoDesasignarHistoria(repo RepositorioAsignacion) *DesasignarHistoria {
	return &DesasignarHistoria{repo: repo}
}

// Ejecutar busca el elemento, calcula si su sprint actual sigue activo, aplica la regla de
// dominio y persiste el resultado. Si el elemento no existe o no se puede desasignar, no se
// guarda nada. Implementa HU-09 (CA-09.4, CA-09.5, CA-09.6, RN-09.7).
func (uc *DesasignarHistoria) Ejecutar(ctx context.Context, d DatosDesasignacion) error {
	e, err := uc.repo.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return err
	}
	if e.ProyectoID != d.ProyectoID { // RN-09.7, igual que RN-18.1: de otro proyecto es "inexistente"
		return domain.ErrElementoNoEncontrado
	}

	var sprintActualActivo bool
	if e.SprintID != nil {
		_, abierto, err := uc.repo.BuscarSprint(ctx, *e.SprintID)
		if err != nil {
			return err
		}
		sprintActualActivo = abierto
	}

	if err := e.DesasignarDeSprint(sprintActualActivo); err != nil {
		return err
	}

	if err := uc.repo.GuardarAsignacion(ctx, &e); err != nil {
		return fmt.Errorf("guardar asignación: %w", err)
	}
	return nil
}
