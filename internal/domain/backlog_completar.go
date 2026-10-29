package domain

// Completar pasa el elemento a Hecho. No depende de domain.Sprint: recibe si el sprint asignado
// al elemento está abierto, calculado por el caso de uso. El sprint cerrado se valida antes que
// "ya está Hecho" (P-8). Si hay error, el elemento no cambia.
// Implementa HU-10 (RN-10.1 a RN-10.3).
func (e *ElementoBacklog) Completar(sprintAbierto bool) error {
	switch {
	case e.SprintID == nil:
		return ErrElementoSinSprint
	case !sprintAbierto:
		return ErrSprintCerrado
	}
	e.Estado = EstadoHecho // ya Hecho: queda igual, sin error (RN-10.3)
	return nil
}
