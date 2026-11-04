package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// listadorDefectos es lo que el handler necesita del caso de uso ListarDefectos.
type listadorDefectos interface {
	Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Defecto, error)
}

// RegistrarListarDefectos agrega al mux la ruta de apoyo de HU-21: listar los defectos de un
// proyecto (HU-20 no expone ninguno; esta historia lo necesita para su pantalla).
func RegistrarListarDefectos(mux *http.ServeMux, uc listadorDefectos) {
	mux.HandleFunc("GET /proyectos/{id}/defectos", handlerListarDefectos(uc))
}

// handlerListarDefectos atiende GET /proyectos/{id}/defectos: 200 con los defectos del proyecto en
// JSON (lista vacía si no tiene ninguno), 404 si el proyecto no existe, 500 ante una falla interna.
// Implementa HU-21 (CA-21.8).
func handlerListarDefectos(uc listadorDefectos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		defectos, err := uc.Ejecutar(r.Context(), proyectoID)
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
				return
			}
			log.Printf("listar defectos del proyecto %d: %v", proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al listar los defectos"})
			return
		}

		respuesta := make([]respuestaDefecto, len(defectos))
		for i, d := range defectos {
			respuesta[i] = aRespuestaDefecto(d)
		}
		escribirJSON(w, http.StatusOK, respuesta)
	}
}
