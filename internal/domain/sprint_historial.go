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

// PorcentajeCompletado calcula qué porcentaje de los PuntosPlanificados de un Sprint se completó
// (RN-24.1): PuntosCompletados / PuntosPlanificados × 100, redondeado a un decimal (redondearUnDecimal,
// igual criterio que la desviación de HU-19, RN-19.6), para que la API y la pantalla siempre
// coincidan. Si PuntosPlanificados es 0, hayPorcentaje es false y porcentaje no significa nada
// (RN-24.2): 0% sería un dato real (se planificó algo y no se completó nada), distinto de "no se
// planificó nada". Implementa HU-24 (RN-24.1, RN-24.2).
func PorcentajeCompletado(s Sprint) (porcentaje float64, hayPorcentaje bool) {
	if s.PuntosPlanificados == 0 {
		return 0, false
	}
	return redondearUnDecimal(float64(s.PuntosCompletados) / float64(s.PuntosPlanificados) * 100), true
}

// ContarDefectosDeSprint cuenta, entre los defectos recibidos, cuántos tienen a sprintID como
// Sprint de detección (RN-24.3) y cuántos lo tienen como Sprint de resolución (RN-24.4), sin
// importar el estado actual de cada defecto. Un defecto sin Sprint asignado (SprintDeteccionID o
// SprintResolucionID nil) no cuenta en ningún Sprint. Implementa HU-24 (RN-24.3, RN-24.4).
func ContarDefectosDeSprint(defectos []Defecto, sprintID int64) (detectados, resueltos int) {
	for _, d := range defectos {
		if d.SprintDeteccionID != nil && *d.SprintDeteccionID == sprintID {
			detectados++
		}
		if d.SprintResolucionID != nil && *d.SprintResolucionID == sprintID {
			resueltos++
		}
	}
	return detectados, resueltos
}
