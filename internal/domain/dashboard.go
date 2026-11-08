package domain

import (
	"slices"
	"time"
)

// ConteoBacklog es cuántos elementos del backlog hay en cada estado. Implementa HU-25 (RN-25.2).
type ConteoBacklog struct {
	Pendientes int
	EnProgreso int
	Hechos     int
	Total      int
}

// ContarBacklog cuenta los elementos por su estado actual, estén o no asignados a un sprint. El
// total es la suma de los tres estados. Implementa HU-25 (RN-25.2).
func ContarBacklog(elementos []ElementoBacklog) ConteoBacklog {
	var c ConteoBacklog
	for _, e := range elementos {
		switch e.Estado {
		case EstadoPendiente:
			c.Pendientes++
		case EstadoEnProgreso:
			c.EnProgreso++
		case EstadoHecho:
			c.Hechos++
		}
	}
	c.Total = c.Pendientes + c.EnProgreso + c.Hechos
	return c
}

// SprintEnCurso es el sprint que se está trabajando hoy, con sus puntos hasta el momento.
type SprintEnCurso struct {
	Sprint             Sprint
	PuntosPlanificados int
	PuntosCompletados  int
}

// ResumirSprintEnCurso busca el sprint en curso (el Abierto cuyas fechas incluyen hoy, la misma regla
// de SprintVigente de HU-20) y le calcula los puntos con la regla del cierre (ResumirCierre, RN-11.6)
// sobre los elementos que tiene asignados. Si no hay ninguno, devuelve false. Implementa HU-25
// (RN-25.3, RN-25.4).
func ResumirSprintEnCurso(sprints []Sprint, elementos []ElementoBacklog, hoy time.Time) (SprintEnCurso, bool) {
	id, ok := SprintVigente(sprints, hoy)
	if !ok {
		return SprintEnCurso{}, false
	}
	i := slices.IndexFunc(sprints, func(s Sprint) bool { return s.ID == id })

	var delSprint []ElementoBacklog
	for _, e := range elementos {
		if e.SprintID != nil && *e.SprintID == id {
			delSprint = append(delSprint, e)
		}
	}
	r := ResumirCierre(delSprint)
	return SprintEnCurso{Sprint: sprints[i], PuntosPlanificados: r.PuntosPlanificados, PuntosCompletados: r.PuntosCompletados}, true
}

// SprintsCerradosCronologicos devuelve los sprints Cerrados del más antiguo al más reciente, todos,
// sin límite: el orden en que se leen los gráficos, de izquierda a derecha. Es el orden inverso de
// SprintsAnteriores (HU-12). No modifica la lista recibida. Implementa HU-25 (RN-25.5).
func SprintsCerradosCronologicos(sprints []Sprint) []Sprint {
	cerrados := SprintsAnteriores(sprints)
	slices.Reverse(cerrados)
	return cerrados
}
