package domain

import "testing"

func historia(titulo string, estado EstadoElemento, sp *int) ElementoBacklog {
	return ElementoBacklog{Titulo: titulo, Estado: estado, StoryPoints: sp}
}

func historiaDeSprint(titulo string, sprintID int64) ElementoBacklog {
	return ElementoBacklog{Titulo: titulo, SprintID: &sprintID}
}

// Traza: HU-26 CA-26.1, CA-26.3
func TestArmarReporteProyecto_HistoriasYResumen(t *testing.T) {
	historias := []ElementoBacklog{
		historia("Login", EstadoHecho, puntos(5)),
		historia("Logout", EstadoHecho, puntos(8)),
		historia("Perfil", EstadoPendiente, puntos(2)),
	}
	p := Proyecto{ID: 1, Nombre: "Alfa"}

	r := ArmarReporteProyecto(p, historias, TotalHoras{}, 0, false, nil, nil)

	if len(r.Historias) != 3 {
		t.Fatalf("Historias = %d, se esperaban 3", len(r.Historias))
	}
	if r.HistoriasCompletadas != 2 {
		t.Errorf("HistoriasCompletadas = %d, se esperaban 2", r.HistoriasCompletadas)
	}
	if r.StoryPointsCompletados != 13 {
		t.Errorf("StoryPointsCompletados = %d, se esperaban 13", r.StoryPointsCompletados)
	}
	if r.StoryPointsTotal != 15 {
		t.Errorf("StoryPointsTotal = %d, se esperaban 15", r.StoryPointsTotal)
	}
}

// Traza: HU-26 CA-26.3: una historia sin estimar no suma a ningún total de story points.
func TestArmarReporteProyecto_HistoriaSinEstimar_NoSumaStoryPoints(t *testing.T) {
	historias := []ElementoBacklog{
		historia("Login", EstadoHecho, puntos(5)),
		historia("Sin estimar", EstadoHecho, nil),
	}

	r := ArmarReporteProyecto(Proyecto{Nombre: "Alfa"}, historias, TotalHoras{}, 0, false, nil, nil)

	if r.HistoriasCompletadas != 2 {
		t.Errorf("HistoriasCompletadas = %d, se esperaban 2", r.HistoriasCompletadas)
	}
	if r.StoryPointsCompletados != 5 || r.StoryPointsTotal != 5 {
		t.Errorf("story points = %d/%d, se esperaban 5/5", r.StoryPointsCompletados, r.StoryPointsTotal)
	}
}

// Traza: HU-26 CA-26.4, CA-26.5
func TestArmarReporteProyecto_EsfuerzoYMetricas_SePropaganTalCual(t *testing.T) {
	esfuerzo := TotalHoras{HorasEstimadas: 2, HorasReales: 5, Elementos: 1}
	metricas := []FilaMetrica{{Sprint: Sprint{Numero: 1}, PorcentajeCompletado: 75, HayPorcentaje: true}}

	r := ArmarReporteProyecto(Proyecto{Nombre: "Alfa"}, nil, esfuerzo, 7, true, metricas, nil)

	if r.Esfuerzo != esfuerzo {
		t.Errorf("Esfuerzo = %+v, se esperaba %+v", r.Esfuerzo, esfuerzo)
	}
	if !r.HayVelocidad || r.Velocidad != 7 {
		t.Errorf("velocidad = %v/%v, se esperaba 7 con datos", r.Velocidad, r.HayVelocidad)
	}
	if len(r.Metricas) != 1 || r.Metricas[0] != metricas[0] {
		t.Errorf("Metricas = %+v, se esperaba %+v", r.Metricas, metricas)
	}
}

// Traza: HU-26 CA-26.6: sin sprints cerrados, sin velocidad y sin filas de métricas.
func TestArmarReporteProyecto_SinMetricas(t *testing.T) {
	r := ArmarReporteProyecto(Proyecto{Nombre: "Alfa"}, nil, TotalHoras{}, 0, false, nil, nil)

	if r.HayVelocidad {
		t.Errorf("se esperaba sin velocidad")
	}
	if len(r.Metricas) != 0 {
		t.Errorf("Metricas = %+v, se esperaba vacío", r.Metricas)
	}
}

// Traza: HU-26 CA-26.9
func TestArmarReporteProyecto_Defectos_SePropaganTalCual(t *testing.T) {
	defectos := []Defecto{{ID: 1, Descripcion: "Falla X"}, {ID: 2, Descripcion: "Falla Y"}}

	r := ArmarReporteProyecto(Proyecto{Nombre: "Alfa"}, nil, TotalHoras{}, 0, false, nil, defectos)

	if len(r.Defectos) != 2 {
		t.Errorf("Defectos = %+v, se esperaban 2", r.Defectos)
	}
}

// Traza: HU-26 CA-26.7, CA-26.8
func TestArmarReporteSprint_Metrica_SePropagaTalCual(t *testing.T) {
	sprint := Sprint{ID: 1, Numero: 1}
	metrica := FilaMetrica{Sprint: sprint, HayPorcentaje: false}

	r := ArmarReporteSprint(Proyecto{Nombre: "Alfa"}, sprint, nil, TotalHoras{}, metrica, nil)

	if r.Metricas[0] != metrica {
		t.Errorf("Metricas[0] = %+v, se esperaba %+v", r.Metricas[0], metrica)
	}
}

// Traza: HU-26 CA-26.2: el filtro de historias de un Sprint, mismo criterio que HU-09
// (pantallas_asignacion.go).
func TestHistoriasDeSprint(t *testing.T) {
	historias := []ElementoBacklog{
		historiaDeSprint("Login", 1),
		historiaDeSprint("Logout", 2),
		{Titulo: "Sin sprint"},
	}

	got := HistoriasDeSprint(historias, 1)

	if len(got) != 1 || got[0].Titulo != "Login" {
		t.Errorf("HistoriasDeSprint() = %+v, se esperaba solo \"Login\"", got)
	}
}

// Traza: HU-26 CA-26.9, CA-26.10: los defectos de un Sprint son los detectados o resueltos ahí.
func TestDefectosDeSprint(t *testing.T) {
	uno, dos := int64(1), int64(2)
	defectos := []Defecto{
		{ID: 1, Descripcion: "Detectado en 1", SprintDeteccionID: &uno},
		{ID: 2, Descripcion: "Resuelto en 1", SprintResolucionID: &uno},
		{ID: 3, Descripcion: "Detectado en 1, resuelto en 2", SprintDeteccionID: &uno, SprintResolucionID: &dos},
		{ID: 4, Descripcion: "De otro sprint", SprintDeteccionID: &dos},
		{ID: 5, Descripcion: "Sin sprint"},
	}

	got := DefectosDeSprint(defectos, 1)

	if len(got) != 3 {
		t.Fatalf("DefectosDeSprint() = %d defectos, se esperaban 3", len(got))
	}
	for _, id := range []int64{1, 2, 3} {
		encontrado := false
		for _, d := range got {
			if d.ID == id {
				encontrado = true
			}
		}
		if !encontrado {
			t.Errorf("falta el defecto %d en el resultado", id)
		}
	}
}
