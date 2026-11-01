package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorVelocidad es lo que el handler necesita del caso de uso ConsultarVelocidad.
type consultorVelocidad interface {
	Ejecutar(ctx context.Context, proyectoID int64) (app.VelocidadProyecto, error)
}

// respuestaVelocidad es el JSON de GET /proyectos/{id}/velocidad (HU-22, sección 3). Velocidad es
// un puntero: queda null si el proyecto todavía no tiene Sprints cerrados (RN-22.3), para no
// confundirlo con un promedio real de 0.
type respuestaVelocidad struct {
	Velocidad       *float64 `json:"velocidad"`
	SprintsCerrados int      `json:"sprints_cerrados"`
}

// RegistrarVelocidad agrega la ruta de HU-22 al mux.
func RegistrarVelocidad(mux *http.ServeMux, uc consultorVelocidad) {
	mux.HandleFunc("GET /proyectos/{id}/velocidad", handlerVelocidad(uc))
}

// handlerVelocidad atiende GET /proyectos/{id}/velocidad: 200 con la velocidad (o null si todavía
// no hay Sprints cerrados), 400 si el ID no es un número, 404 si el proyecto no existe y 500 ante
// una falla interna, sin el detalle. Implementa HU-22 (CA-22.7).
func handlerVelocidad(uc consultorVelocidad) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		v, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("consultar velocidad del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al consultar la velocidad"})
		default:
			escribirJSON(w, http.StatusOK, aRespuestaVelocidad(v))
		}
	}
}

func aRespuestaVelocidad(v app.VelocidadProyecto) respuestaVelocidad {
	r := respuestaVelocidad{SprintsCerrados: v.SprintsCerrados}
	if v.HayDatos {
		r.Velocidad = &v.Velocidad
	}
	return r
}
