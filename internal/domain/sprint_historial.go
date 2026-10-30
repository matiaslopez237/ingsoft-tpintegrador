package domain

import "sort"

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
	sort.Slice(anteriores, func(i, j int) bool {
		return anteriores[i].FechaInicio.After(anteriores[j].FechaInicio)
	})
	return anteriores
}
