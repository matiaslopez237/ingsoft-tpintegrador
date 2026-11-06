package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorCalidadSprints es lo que el handler necesita del caso de uso ConsultarCalidadSprints.
type consultorCalidadSprints interface {
	Ejecutar(ctx context.Context, proyectoID int64) ([]app.CalidadSprint, error)
}

// itemCalidad es un Sprint en el JSON de GET /proyectos/{id}/calidad (HU-24, sección 3).
// PorcentajeCompletado es un puntero: queda null si el Sprint se cerró sin puntos planificados
// (RN-24.2), para no confundirlo con un 0% real.
type itemCalidad struct {
	SprintID             int64    `json:"sprint_id"`
	Numero               int      `json:"numero"`
	PorcentajeCompletado *float64 `json:"porcentaje_completado"`
	DefectosDetectados   int      `json:"defectos_detectados"`
	DefectosResueltos    int      `json:"defectos_resueltos"`
}

// RegistrarCalidadSprints agrega la ruta de HU-24 al mux.
func RegistrarCalidadSprints(mux *http.ServeMux, uc consultorCalidadSprints) {
	mux.HandleFunc("GET /proyectos/{id}/calidad", handlerCalidadSprints(uc))
}

// handlerCalidadSprints atiende GET /proyectos/{id}/calidad: 200 con un array (uno por Sprint
// Cerrado, vacío si no tiene ninguno), 400 si el ID no es un número, 404 si el proyecto no existe y
// 500 ante una falla interna, sin el detalle. Implementa HU-24 (CA-24.8).
func handlerCalidadSprints(uc consultorCalidadSprints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		c, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("consultar calidad de los sprints del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al consultar la calidad de los sprints"})
		default:
			escribirJSON(w, http.StatusOK, aRespuestaCalidad(c))
		}
	}
}

func aRespuestaCalidad(c []app.CalidadSprint) []itemCalidad {
	r := make([]itemCalidad, len(c))
	for i, s := range c {
		item := itemCalidad{
			SprintID: s.Sprint.ID, Numero: s.Sprint.Numero,
			DefectosDetectados: s.DefectosDetectados, DefectosResueltos: s.DefectosResueltos,
		}
		if s.HayPorcentaje {
			item.PorcentajeCompletado = &s.PorcentajeCompletado
		}
		r[i] = item
	}
	return r
}
