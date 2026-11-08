package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorDashboardFalso devuelve lo que se le cargue, o err si está seteado.
type consultorDashboardFalso struct {
	dashboard app.Dashboard
	err       error
	recibido  int64
	llamadas  int
}

func (f *consultorDashboardFalso) Ejecutar(_ context.Context, proyectoID int64) (app.Dashboard, error) {
	f.recibido = proyectoID
	f.llamadas++
	if f.err != nil {
		return app.Dashboard{}, f.err
	}
	return f.dashboard, nil
}

func pedirDashboard(t *testing.T, uc consultorDashboard, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/proyectos/"+id+"/dashboard", nil)
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarDashboard(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

func horasDashboard(h float64) *float64 { return &h }

// dashboardCompleto es un panel con datos en todas las secciones.
func dashboardCompleto() app.Dashboard {
	enCurso := domain.Sprint{
		ID: 30, Numero: 3, Objetivo: "Cerrar el login", Estado: domain.SprintAbierto,
		FechaInicio: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC),
	}
	return app.Dashboard{
		Backlog:          domain.ConteoBacklog{Pendientes: 2, EnProgreso: 1, Hechos: 3, Total: 6},
		SprintEnCurso:    domain.SprintEnCurso{Sprint: enCurso, PuntosPlanificados: 8, PuntosCompletados: 5},
		HaySprintEnCurso: true,
		SprintsCerrados: []app.CalidadSprint{
			{
				Sprint:               domain.Sprint{ID: 10, Numero: 1, PuntosPlanificados: 8, PuntosCompletados: 6},
				PorcentajeCompletado: 75, HayPorcentaje: true, DefectosDetectados: 2, DefectosResueltos: 0,
			},
			{
				Sprint:             domain.Sprint{ID: 20, Numero: 2},
				DefectosDetectados: 0, DefectosResueltos: 1,
			},
		},
		Velocidad: app.VelocidadProyecto{Velocidad: 3, HayDatos: true, SprintsCerrados: 2},
		Horas: domain.TotalHoras{
			HorasEstimadas: 10, HorasReales: 12, DesviacionHoras: horasDashboard(2), DesviacionPorcentaje: horasDashboard(20),
			Elementos: 1,
		},
	}
}

// Traza: HU-25 CA-25.8, sección 3: 200 con el backlog, el sprint en curso, los sprints cerrados,
// la velocidad y las horas.
func TestGetDashboard_ConDatos_DevuelveElPanel(t *testing.T) {
	uc := &consultorDashboardFalso{dashboard: dashboardCompleto()}

	rec := pedirDashboard(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 3 {
		t.Errorf("proyecto recibido = %d, se esperaba 3", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `{
		"backlog": {"pendientes":2,"en_progreso":1,"hechos":3,"total":6},
		"sprint_en_curso": {"sprint_id":30,"numero":3,"objetivo":"Cerrar el login","fecha_inicio":"2026-11-01","fecha_fin":"2026-11-14",
			"puntos_planificados":8,"puntos_completados":5},
		"sprints_cerrados": [
			{"sprint_id":10,"numero":1,"puntos_planificados":8,"puntos_completados":6,"porcentaje_completado":75,
				"defectos_detectados":2,"defectos_resueltos":0},
			{"sprint_id":20,"numero":2,"puntos_planificados":0,"puntos_completados":0,"porcentaje_completado":null,
				"defectos_detectados":0,"defectos_resueltos":1}
		],
		"velocidad": 3,
		"horas": {"horas_estimadas":10,"horas_reales":12,"horas_reales_sin_estimar":0,"desviacion_horas":2,
			"desviacion_porcentaje":20,"elementos":1,"elementos_sin_estimar":0}
	}`)
}

// Traza: HU-25 CA-25.8, CA-25.11: sin datos, sprint_en_curso y velocidad en null y un array vacío.
func TestGetDashboard_SinDatos_DevuelveNullYArrayVacio(t *testing.T) {
	uc := &consultorDashboardFalso{dashboard: app.Dashboard{SprintsCerrados: []app.CalidadSprint{}}}

	rec := pedirDashboard(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	mismoJSON(t, rec.Body.String(), `{
		"backlog": {"pendientes":0,"en_progreso":0,"hechos":0,"total":0},
		"sprint_en_curso": null,
		"sprints_cerrados": [],
		"velocidad": null,
		"horas": {"horas_estimadas":0,"horas_reales":0,"horas_reales_sin_estimar":0,"desviacion_horas":null,
			"desviacion_porcentaje":null,"elementos":0,"elementos_sin_estimar":0}
	}`)
}

// Traza: HU-25 CA-25.8: los sprints cerrados en nil (caso de uso sin sprints) salen como [], no null.
func TestGetDashboard_SprintsCerradosNil_DevuelveArrayVacio(t *testing.T) {
	rec := pedirDashboard(t, &consultorDashboardFalso{}, "3")

	mismoJSON(t, rec.Body.String(), `{
		"backlog": {"pendientes":0,"en_progreso":0,"hechos":0,"total":0},
		"sprint_en_curso": null, "sprints_cerrados": [], "velocidad": null,
		"horas": {"horas_estimadas":0,"horas_reales":0,"horas_reales_sin_estimar":0,"desviacion_horas":null,
			"desviacion_porcentaje":null,"elementos":0,"elementos_sin_estimar":0}
	}`)
}

// Traza: HU-25 CA-25.8, sección 7: ID inválido da 400 sin llamar al caso de uso.
func TestGetDashboard_IDInvalido_Devuelve400(t *testing.T) {
	for _, id := range []string{"abc", "0", "-1"} {
		t.Run(id, func(t *testing.T) {
			uc := &consultorDashboardFalso{}

			rec := pedirDashboard(t, uc, id)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, se esperaba 400", rec.Code)
			}
			if uc.llamadas != 0 {
				t.Errorf("se llamó al caso de uso %d veces con un ID inválido", uc.llamadas)
			}
		})
	}
}

// Traza: HU-25 CA-25.7, CA-25.8: proyecto inexistente da 404.
func TestGetDashboard_ProyectoInexistente_Devuelve404(t *testing.T) {
	uc := &consultorDashboardFalso{err: fmt.Errorf("buscar proyecto 9: %w", domain.ErrProyectoNoEncontrado)}

	rec := pedirDashboard(t, uc, "9")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `{"error":"proyecto no encontrado"}`)
}

// Traza: HU-25, sección 7: una falla interna da 500 sin el detalle.
func TestGetDashboard_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	uc := &consultorDashboardFalso{err: errors.New("disco lleno")}

	rec := pedirDashboard(t, uc, "3")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `{"error":"error interno al consultar el dashboard"}`)
}
