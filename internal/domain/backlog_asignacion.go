package domain

import "errors"

// Errores de HU-09 (spec, sección 7).
var (
	ErrSprintDeOtroProyecto = errors.New("el sprint pertenece a otro proyecto")
	ErrElementoYaAsignado   = errors.New("el elemento ya está asignado a otro sprint activo")
	ErrSprintCerrado        = errors.New("el sprint ya está cerrado")
	ErrElementoSinSprint    = errors.New("el elemento no tiene un sprint asignado")
)

// AsignarASprint asigna el elemento a un sprint. No depende de domain.Sprint (HU-08 todavía no
// está mergeada): recibe el proyecto y el estado del sprint destino, y si el sprint que el
// elemento ya tenía asignado (si tiene) sigue activo, todo calculado por el caso de uso.
// Implementa HU-09 (RN-09.1 a RN-09.3, RN-09.5, RN-09.6).
func (e *ElementoBacklog) AsignarASprint(sprintID, sprintProyectoID int64, sprintDestinoAbierto, sprintActualActivo bool) error {
	switch {
	case sprintProyectoID != e.ProyectoID:
		return ErrSprintDeOtroProyecto
	case e.SprintID != nil && *e.SprintID != sprintID && sprintActualActivo:
		return ErrElementoYaAsignado
	case !sprintDestinoAbierto:
		return ErrSprintCerrado
	}
	e.SprintID = &sprintID
	return nil
}

// DesasignarDeSprint quita al elemento de su sprint actual, mientras ese sprint siga activo.
// Implementa HU-09 (RN-09.4).
func (e *ElementoBacklog) DesasignarDeSprint(sprintActualActivo bool) error {
	switch {
	case e.SprintID == nil:
		return ErrElementoSinSprint
	case !sprintActualActivo:
		return ErrSprintCerrado
	}
	e.SprintID = nil
	return nil
}
