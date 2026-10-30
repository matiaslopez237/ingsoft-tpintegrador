package domain

import "slices"

// SprintsAnteriores devuelve, de los sprints recibidos, los Cerrados (RN-12.1), del más reciente al
// más antiguo por fecha de inicio (RN-12.3), con todos sus datos (RN-12.2). Como los sprints de un
// proyecto no se superponen (RN-08.5), no hay empates. No modifica la lista recibida.
// Implementa HU-12.
func SprintsAnteriores(sprints []Sprint) []Sprint {
	anteriores := []Sprint{}
	for _, s := range sprints {
		if s.Estado == SprintCerrado {
			anteriores = append(anteriores, s)
		}
	}
	slices.SortFunc(anteriores, func(a, b Sprint) int {
		return b.FechaInicio.Compare(a.FechaInicio) // b antes que a: de la más reciente a la más antigua
	})
	return anteriores
}
