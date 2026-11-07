package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorDashboard es lo que el handler y la pantalla necesitan del caso de uso ConsultarDashboard.
type consultorDashboard interface {
	Ejecutar(ctx context.Context, proyectoID int64) (app.Dashboard, error)
}

// respuestaDashboard es el JSON de GET /proyectos/{id}/dashboard (HU-25, sección 3). SprintEnCurso y
// Velocidad son punteros: quedan null si no hay sprint en curso (RN-25.3) o sprints cerrados (RN-22.3).
type respuestaDashboard struct {
	Backlog         respuestaConteoBacklog   `json:"backlog"`
	SprintEnCurso   *respuestaSprintEnCurso  `json:"sprint_en_curso"`
	SprintsCerrados []respuestaSprintCerrado `json:"sprints_cerrados"`
	Velocidad       *float64                 `json:"velocidad"`
	Horas           respuestaTotalHoras      `json:"horas"` // el mismo formato que GET /proyectos/{id}/horas (HU-23)
}

type respuestaConteoBacklog struct {
	Pendientes int `json:"pendientes"`
	EnProgreso int `json:"en_progreso"`
	Hechos     int `json:"hechos"`
	Total      int `json:"total"`
}

type respuestaSprintEnCurso struct {
	SprintID           int64  `json:"sprint_id"`
	Numero             int    `json:"numero"`
	Objetivo           string `json:"objetivo"`
	FechaInicio        string `json:"fecha_inicio"`
	FechaFin           string `json:"fecha_fin"`
	PuntosPlanificados int    `json:"puntos_planificados"`
	PuntosCompletados  int    `json:"puntos_completados"`
}

// respuestaSprintCerrado suma los puntos de HU-11 a los datos de calidad de HU-24. PorcentajeCompletado
// queda null si el sprint se cerró sin puntos planificados (RN-24.2).
type respuestaSprintCerrado struct {
	SprintID             int64    `json:"sprint_id"`
	Numero               int      `json:"numero"`
	PuntosPlanificados   int      `json:"puntos_planificados"`
	PuntosCompletados    int      `json:"puntos_completados"`
	PorcentajeCompletado *float64 `json:"porcentaje_completado"`
	DefectosDetectados   int      `json:"defectos_detectados"`
	DefectosResueltos    int      `json:"defectos_resueltos"`
}

// RegistrarDashboard agrega la ruta JSON de HU-25 al mux.
func RegistrarDashboard(mux *http.ServeMux, uc consultorDashboard) {
	mux.HandleFunc("GET /proyectos/{id}/dashboard", handlerDashboard(uc))
}

// handlerDashboard atiende GET /proyectos/{id}/dashboard: 200 con el panel, 400 si el ID no es un
// número positivo, 404 si el proyecto no existe y 500 ante una falla interna, sin el detalle.
// Implementa HU-25 (CA-25.8).
func handlerDashboard(uc consultorDashboard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		d, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("consultar el dashboard del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al consultar el dashboard"})
		default:
			escribirJSON(w, http.StatusOK, aRespuestaDashboard(d))
		}
	}
}

func aRespuestaDashboard(d app.Dashboard) respuestaDashboard {
	r := respuestaDashboard{
		Backlog: respuestaConteoBacklog{
			Pendientes: d.Backlog.Pendientes, EnProgreso: d.Backlog.EnProgreso, Hechos: d.Backlog.Hechos, Total: d.Backlog.Total,
		},
		SprintsCerrados: make([]respuestaSprintCerrado, len(d.SprintsCerrados)),
		Horas:           aRespuestaTotalHoras(d.Horas),
	}
	if d.HaySprintEnCurso {
		s := d.SprintEnCurso.Sprint
		r.SprintEnCurso = &respuestaSprintEnCurso{
			SprintID: s.ID, Numero: s.Numero, Objetivo: s.Objetivo,
			FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
			PuntosPlanificados: d.SprintEnCurso.PuntosPlanificados, PuntosCompletados: d.SprintEnCurso.PuntosCompletados,
		}
	}
	for i, c := range d.SprintsCerrados {
		item := respuestaSprintCerrado{
			SprintID: c.Sprint.ID, Numero: c.Sprint.Numero,
			PuntosPlanificados: c.Sprint.PuntosPlanificados, PuntosCompletados: c.Sprint.PuntosCompletados,
			DefectosDetectados: c.DefectosDetectados, DefectosResueltos: c.DefectosResueltos,
		}
		if c.HayPorcentaje {
			item.PorcentajeCompletado = &c.PorcentajeCompletado
		}
		r.SprintsCerrados[i] = item
	}
	if d.Velocidad.HayDatos {
		r.Velocidad = &d.Velocidad.Velocidad
	}
	return r
}
