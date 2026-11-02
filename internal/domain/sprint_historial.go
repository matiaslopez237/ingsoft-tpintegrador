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

// VentanaVelocidad es la cantidad de Sprints Cerrados más recientes que entran en el promedio de
// la velocidad (RN-22.2): si el proyecto tiene más, los más antiguos no cuentan.
const VentanaVelocidad = 3

// Velocidad calcula el promedio de PuntosCompletados de los últimos VentanaVelocidad sprints
// recibidos: quien llama le pasa los ya filtrados por SprintsAnteriores, del más reciente al más
// antiguo (Velocidad no vuelve a filtrar por estado ni a ordenar, solo se queda con los primeros
// VentanaVelocidad). cantidad es cuántos entraron en el promedio (como mucho VentanaVelocidad), no
// el total de sprints recibidos. Si la lista está vacía, hayDatos es false y velocidad no significa
// nada (RN-22.3): no hay Sprints cerrados de los que promediar, que es distinto de un promedio real
// de 0 puntos. Implementa HU-22 (RN-22.2 a RN-22.4).
func Velocidad(cerrados []Sprint) (velocidad float64, cantidad int, hayDatos bool) {
	if len(cerrados) == 0 {
		return 0, 0, false
	}
	ventana := cerrados
	if len(ventana) > VentanaVelocidad {
		ventana = ventana[:VentanaVelocidad]
	}
	total := 0
	for _, s := range ventana {
		total += s.PuntosCompletados
	}
	return float64(total) / float64(len(ventana)), len(ventana), true
}
