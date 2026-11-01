package http

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// respuestaSprintAnterior es cada elemento de GET /proyectos/{id}/sprints/anteriores (HU-12,
// sección 3). Fechas AAAA-MM-DD.
type respuestaSprintAnterior struct {
	ID                 int64  `json:"id"`
	Numero             int    `json:"numero"`
	FechaInicio        string `json:"fecha_inicio"`
	FechaFin           string `json:"fecha_fin"`
	Objetivo           string `json:"objetivo"`
	PuntosPlanificados int    `json:"puntos_planificados"`
	PuntosCompletados  int    `json:"puntos_completados"`
}

// RegistrarSprintsAnteriores agrega al mux la ruta de HU-12. El caso de uso tiene la misma forma
// que el de listar sprints de HU-08, así que se reutiliza la interfaz listadorSprints. La ruta
// termina en el literal "anteriores", que tiene prioridad sobre las rutas con {sprintID}.
func RegistrarSprintsAnteriores(mux *http.ServeMux, uc listadorSprints) {
	mux.HandleFunc("GET /proyectos/{id}/sprints/anteriores", handlerSprintsAnteriores(uc))
}

// handlerSprintsAnteriores atiende GET /proyectos/{id}/sprints/anteriores: 200 con los sprints
// cerrados (siempre un array, [] si no hay), 400 si el ID no es un número, 404 si el proyecto no
// existe y 500 ante una falla interna, sin el detalle. Implementa HU-12 (CA-12.6).
func handlerSprintsAnteriores(uc listadorSprints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		sprints, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("listar sprints anteriores del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al listar los sprints anteriores"})
		default:
			respuesta := make([]respuestaSprintAnterior, 0, len(sprints))
			for _, s := range sprints {
				respuesta = append(respuesta, respuestaSprintAnterior{
					ID: s.ID, Numero: s.Numero,
					FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
					Objetivo: s.Objetivo, PuntosPlanificados: s.PuntosPlanificados, PuntosCompletados: s.PuntosCompletados,
				})
			}
			escribirJSON(w, http.StatusOK, respuesta)
		}
	}
}
