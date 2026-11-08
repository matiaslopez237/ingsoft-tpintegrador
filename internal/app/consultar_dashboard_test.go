package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// dashboardFalso implementa en memoria los cinco puertos que necesita ConsultarDashboard. Reutiliza
// calidadSprintsFalso (HU-24) para proyectos, sprints y defectos, y le suma elementos y horas.
type dashboardFalso struct {
	*calidadSprintsFalso
	elementos    []domain.ElementoBacklog
	horas        []domain.HorasElemento
	errElementos error
	errHoras     error
	pedidosHoras int
}

func nuevoDashboardFalso(proyectos ...domain.Proyecto) *dashboardFalso {
	return &dashboardFalso{calidadSprintsFalso: nuevoCalidadSprintsFalso(proyectos...)}
}

func (f *dashboardFalso) ListarElementos(_ context.Context, proyectoID int64) ([]domain.ElementoBacklog, error) {
	if f.errElementos != nil {
		return nil, f.errElementos
	}
	var del []domain.ElementoBacklog
	for _, e := range f.elementos {
		if e.ProyectoID == proyectoID {
			del = append(del, e)
		}
	}
	return del, nil
}

func (f *dashboardFalso) HorasDeProyecto(_ context.Context, _ int64) ([]domain.HorasElemento, error) {
	f.pedidosHoras++
	if f.errHoras != nil {
		return nil, f.errHoras
	}
	return f.horas, nil
}

// hoyDashboard es el "hoy" fijo de estos tests.
var hoyDashboard = time.Date(2026, 11, 10, 15, 0, 0, 0, time.UTC)

func dashboardCon(f *dashboardFalso) *ConsultarDashboard {
	return NuevoConsultarDashboard(f, f, f, f, f, func() time.Time { return hoyDashboard })
}

func sprintConFechas(id int64, numero int, inicio, fin time.Time, estado domain.EstadoSprint, planificados, completados int) domain.Sprint {
	return domain.Sprint{ID: id, ProyectoID: 1, Numero: numero, FechaInicio: inicio, FechaFin: fin, Objetivo: "Objetivo",
		Estado: estado, PuntosPlanificados: planificados, PuntosCompletados: completados}
}

func dia(mes time.Month, d int) time.Time { return time.Date(2026, mes, d, 0, 0, 0, 0, time.UTC) }

func elementoDelDashboard(id int64, estado domain.EstadoElemento, puntos int, sprintID *int64) domain.ElementoBacklog {
	return domain.ElementoBacklog{ID: id, ProyectoID: 1, Titulo: "Elemento", Estado: estado, StoryPoints: &puntos, SprintID: sprintID}
}

func ptrID(id int64) *int64 { return &id }

func ptrHoras(h float64) *float64 { return &h }

// Traza: HU-25 CA-25.1 a CA-25.5: junta el backlog, el sprint en curso, los sprints cerrados (en
// orden cronológico, con porcentaje y defectos), la velocidad y las horas del proyecto.
func TestConsultarDashboard_JuntaTodasLasMetricas(t *testing.T) {
	f := nuevoDashboardFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{
		sprintConFechas(20, 2, dia(time.October, 15), dia(time.October, 28), domain.SprintCerrado, 4, 4),
		sprintConFechas(10, 1, dia(time.October, 1), dia(time.October, 14), domain.SprintCerrado, 8, 6),
		sprintConFechas(30, 3, dia(time.November, 1), dia(time.November, 14), domain.SprintAbierto, 0, 0),
	}
	f.elementos = []domain.ElementoBacklog{
		elementoDelDashboard(1, domain.EstadoPendiente, 3, ptrID(30)),
		elementoDelDashboard(2, domain.EstadoHecho, 5, ptrID(30)),
		elementoDelDashboard(3, domain.EstadoEnProgreso, 2, nil),
	}
	f.defectos = []domain.Defecto{
		{ID: 1, ProyectoID: 1, SprintDeteccionID: ptrID(10), SprintResolucionID: ptrID(20)},
		{ID: 2, ProyectoID: 1, SprintDeteccionID: ptrID(10)},
	}
	f.horas = []domain.HorasElemento{{HorasEstimadas: ptrHoras(10), HorasReales: 12}}

	d, err := dashboardCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := (domain.ConteoBacklog{Pendientes: 1, EnProgreso: 1, Hechos: 1, Total: 3}); d.Backlog != want {
		t.Errorf("Backlog = %+v, se esperaba %+v", d.Backlog, want)
	}
	if !d.HaySprintEnCurso || d.SprintEnCurso.Sprint.ID != 30 ||
		d.SprintEnCurso.PuntosPlanificados != 8 || d.SprintEnCurso.PuntosCompletados != 5 {
		t.Errorf("sprint en curso = %+v (hay: %v), se esperaba el 30 con 8 planificados y 5 completados", d.SprintEnCurso, d.HaySprintEnCurso)
	}
	if len(d.SprintsCerrados) != 2 {
		t.Fatalf("SprintsCerrados = %d, se esperaban 2", len(d.SprintsCerrados))
	}
	primero, segundo := d.SprintsCerrados[0], d.SprintsCerrados[1]
	if primero.Sprint.ID != 10 || segundo.Sprint.ID != 20 {
		t.Errorf("orden de los sprints cerrados = %d, %d; se esperaba 10, 20 (del más antiguo al más reciente)", primero.Sprint.ID, segundo.Sprint.ID)
	}
	if !primero.HayPorcentaje || primero.PorcentajeCompletado != 75 || primero.DefectosDetectados != 2 || primero.DefectosResueltos != 0 {
		t.Errorf("sprint 1 = %+v, se esperaba 75%%, 2 detectados y 0 resueltos", primero)
	}
	if !segundo.HayPorcentaje || segundo.PorcentajeCompletado != 100 || segundo.DefectosDetectados != 0 || segundo.DefectosResueltos != 1 {
		t.Errorf("sprint 2 = %+v, se esperaba 100%%, 0 detectados y 1 resuelto", segundo)
	}
	if want := (VelocidadProyecto{Velocidad: 5, HayDatos: true, SprintsCerrados: 2}); d.Velocidad != want {
		t.Errorf("Velocidad = %+v, se esperaba %+v", d.Velocidad, want)
	}
	if d.Horas.HorasEstimadas != 10 || d.Horas.HorasReales != 12 || d.Horas.DesviacionPorcentaje == nil || *d.Horas.DesviacionPorcentaje != 20 {
		t.Errorf("Horas = %+v, se esperaban 10 estimadas, 12 reales y 20%% de desviación", d.Horas)
	}
}

// Traza: HU-25 CA-25.11: un proyecto sin datos devuelve el dashboard vacío, sin error.
func TestConsultarDashboard_ProyectoSinDatos(t *testing.T) {
	f := nuevoDashboardFalso(proyectoActivo(1))

	d, err := dashboardCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Backlog != (domain.ConteoBacklog{}) || d.HaySprintEnCurso || len(d.SprintsCerrados) != 0 ||
		d.Velocidad.HayDatos || !d.Horas.SinEstimar() {
		t.Errorf("dashboard = %+v, se esperaba vacío", d)
	}
}

// Traza: HU-25 CA-25.6 (RN-25.1): cada consulta lee los datos de nuevo, sin guardar nada aparte.
func TestConsultarDashboard_LeeLosDatosVigentesEnCadaConsulta(t *testing.T) {
	f := nuevoDashboardFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{sprintConFechas(30, 1, dia(time.November, 1), dia(time.November, 14), domain.SprintAbierto, 0, 0)}
	f.elementos = []domain.ElementoBacklog{elementoDelDashboard(1, domain.EstadoEnProgreso, 3, ptrID(30))}
	uc := dashboardCon(f)
	if _, err := uc.Ejecutar(context.Background(), 1); err != nil {
		t.Fatalf("primera consulta: %v", err)
	}

	f.elementos[0].Estado = domain.EstadoHecho
	d, err := uc.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("segunda consulta: %v", err)
	}
	if d.SprintEnCurso.PuntosCompletados != 3 || d.Backlog.Hechos != 1 {
		t.Errorf("después de completar la historia: %d puntos completados y %d hechos; se esperaba 3 y 1",
			d.SprintEnCurso.PuntosCompletados, d.Backlog.Hechos)
	}
}

// Traza: HU-25 CA-25.3 (RN-25.3): usa el reloj para decidir el sprint en curso.
func TestConsultarDashboard_SprintAbiertoVencido_NoEstaEnCurso(t *testing.T) {
	f := nuevoDashboardFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{sprintConFechas(30, 1, dia(time.October, 1), dia(time.October, 14), domain.SprintAbierto, 0, 0)}

	d, err := dashboardCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.HaySprintEnCurso {
		t.Errorf("se esperaba sin sprint en curso, pero está el %d", d.SprintEnCurso.Sprint.ID)
	}
}

// Traza: HU-25 CA-25.7: un proyecto que no existe devuelve ErrProyectoNoEncontrado sin leer el resto.
func TestConsultarDashboard_ProyectoInexistente(t *testing.T) {
	f := nuevoDashboardFalso()

	_, err := dashboardCon(f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba ErrProyectoNoEncontrado", err)
	}
	if f.pedidosHoras != 0 {
		t.Errorf("se leyeron las horas %d veces de un proyecto inexistente", f.pedidosHoras)
	}
}

// Traza: HU-25, sección 7: una falla de cualquier puerto se informa envuelta, sin dashboard a medias.
func TestConsultarDashboard_FallaDeUnPuerto_SeInforma(t *testing.T) {
	falla := errors.New("falla de la base")
	tests := []struct {
		name   string
		romper func(f *dashboardFalso)
	}{
		{"al buscar el proyecto", func(f *dashboardFalso) { f.errBuscar = falla }},
		{"al listar los sprints", func(f *dashboardFalso) { f.errSprints = falla }},
		{"al listar los elementos", func(f *dashboardFalso) { f.errElementos = falla }},
		{"al listar los defectos", func(f *dashboardFalso) { f.errDefectos = falla }},
		{"al leer las horas", func(f *dashboardFalso) { f.errHoras = falla }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoDashboardFalso(proyectoActivo(1))
			tt.romper(f)

			_, err := dashboardCon(f).Ejecutar(context.Background(), 1)

			if !errors.Is(err, falla) {
				t.Errorf("error = %v, se esperaba que envolviera la falla", err)
			}
		})
	}
}
