package domain

// FilaMetrica es la calidad de un Sprint dentro de un reporte (HU-26, RN-26.5, RN-26.6): mismo
// contenido que app.CalidadSprint de HU-24, pero en el dominio, porque el dominio no puede depender
// de internal/app.
type FilaMetrica struct {
	Sprint               Sprint
	PorcentajeCompletado float64
	HayPorcentaje        bool
	DefectosDetectados   int
	DefectosResueltos    int
}

// Reporte es el reporte de un proyecto o de un Sprint, con sus datos ya resueltos: Sprint es nil
// para el reporte de un proyecto. No tiene ningún dato ya formateado como texto (eso es tarea de la
// pantalla, igual que en el resto del proyecto): solo junta lo que ya calcularon otras historias.
// Lo arman ArmarReporteProyecto y ArmarReporteSprint. Implementa HU-26.
type Reporte struct {
	Proyecto               Proyecto
	Sprint                 *Sprint // nil: reporte de un proyecto
	Historias              []ElementoBacklog
	HistoriasCompletadas   int
	StoryPointsCompletados int
	StoryPointsTotal       int
	Esfuerzo               TotalHoras
	HayVelocidad           bool
	Velocidad              float64
	Metricas               []FilaMetrica
	Defectos               []Defecto
}

// ArmarReporteProyecto arma el reporte de un proyecto (RN-26.1, RN-26.5, RN-26.7): historias,
// resumen (RN-26.3), esfuerzo, velocidad y métricas de los Sprints cerrados (ya calculados por el
// caso de uso) y defectos, todos tal cual los recibe. Implementa HU-26.
func ArmarReporteProyecto(p Proyecto, historias []ElementoBacklog, esfuerzo TotalHoras,
	velocidad float64, hayVelocidad bool, metricas []FilaMetrica, defectos []Defecto,
) Reporte {
	completadas, spCompletados, spTotal := resumenHistorias(historias)
	return Reporte{
		Proyecto:               p,
		Historias:              historias,
		HistoriasCompletadas:   completadas,
		StoryPointsCompletados: spCompletados,
		StoryPointsTotal:       spTotal,
		Esfuerzo:               esfuerzo,
		HayVelocidad:           hayVelocidad,
		Velocidad:              velocidad,
		Metricas:               metricas,
		Defectos:               defectos,
	}
}

// ArmarReporteSprint arma el reporte de un Sprint (RN-26.2, RN-26.6, RN-26.8): mismas secciones que
// ArmarReporteProyecto, acotadas a ese Sprint y sin velocidad (que es del equipo, no de un Sprint).
// Implementa HU-26.
func ArmarReporteSprint(p Proyecto, s Sprint, historias []ElementoBacklog, esfuerzo TotalHoras,
	metrica FilaMetrica, defectos []Defecto,
) Reporte {
	completadas, spCompletados, spTotal := resumenHistorias(historias)
	return Reporte{
		Proyecto:               p,
		Sprint:                 &s,
		Historias:              historias,
		HistoriasCompletadas:   completadas,
		StoryPointsCompletados: spCompletados,
		StoryPointsTotal:       spTotal,
		Esfuerzo:               esfuerzo,
		Metricas:               []FilaMetrica{metrica},
		Defectos:               defectos,
	}
}

// resumenHistorias cuenta las historias completadas (RN-26.3) y sus story points, sobre el total de
// historias y story points del reporte. Una historia sin estimar suma al total de historias, no al
// de story points (ninguno de los dos).
func resumenHistorias(historias []ElementoBacklog) (completadas, spCompletados, spTotal int) {
	for _, h := range historias {
		if h.Estado == EstadoHecho {
			completadas++
		}
		if h.StoryPoints == nil {
			continue
		}
		spTotal += *h.StoryPoints
		if h.Estado == EstadoHecho {
			spCompletados += *h.StoryPoints
		}
	}
	return completadas, spCompletados, spTotal
}

// HistoriasDeSprint filtra las historias asignadas a un Sprint (RN-26.2), mismo criterio que ya usa
// la pantalla de asignación de HU-09. Implementa HU-26.
func HistoriasDeSprint(historias []ElementoBacklog, sprintID int64) []ElementoBacklog {
	del := []ElementoBacklog{}
	for _, h := range historias {
		if h.SprintID != nil && *h.SprintID == sprintID {
			del = append(del, h)
		}
	}
	return del
}

// DefectosDeSprint filtra los defectos que tienen ese Sprint como Sprint de detección o como Sprint
// de resolución (RN-26.8), mismo criterio que ContarDefectosDeSprint (HU-24) pero devolviendo la
// lista, no solo los conteos. Implementa HU-26.
func DefectosDeSprint(defectos []Defecto, sprintID int64) []Defecto {
	del := []Defecto{}
	for _, d := range defectos {
		detectado := d.SprintDeteccionID != nil && *d.SprintDeteccionID == sprintID
		resuelto := d.SprintResolucionID != nil && *d.SprintResolucionID == sprintID
		if detectado || resuelto {
			del = append(del, d)
		}
	}
	return del
}
