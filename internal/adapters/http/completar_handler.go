package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// completadorHistorias es lo que el handler necesita del caso de uso CompletarHistoria.
type completadorHistorias interface {
	Ejecutar(ctx context.Context, d app.DatosCompletar) error
}

// erroresNoEncontradoCompletar son los errores de HU-10 que se responden como 404.
var erroresNoEncontradoCompletar = []error{
	domain.ErrProyectoNoEncontrado,
	domain.ErrElementoNoEncontrado,
}

// erroresDeValidacionCompletar son los errores de HU-10 que se responden como 400.
var erroresDeValidacionCompletar = []error{
	domain.ErrElementoSinSprint,
	domain.ErrSprintCerrado,
}

// RegistrarCompletar agrega al mux la ruta de HU-10: marcar un elemento del sprint como Hecho.
func RegistrarCompletar(mux *http.ServeMux, uc completadorHistorias) {
	mux.HandleFunc("POST /proyectos/{id}/backlog/{elementoID}/completar", handlerCompletarHistoria(uc))
}

// handlerCompletarHistoria atiende POST /proyectos/{id}/backlog/{elementoID}/completar (sin
// cuerpo): 204 si se completó, 404 si el proyecto o el elemento no existen, 400 si el elemento
// no tiene sprint o el sprint está cerrado, 500 ante una falla interna. Implementa HU-10 (CA-10.6).
func handlerCompletarHistoria(uc completadorHistorias) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		elementoID, err := elementoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		err = uc.Ejecutar(r.Context(), app.DatosCompletar{ProyectoID: proyectoID, ElementoID: elementoID})
		if err != nil {
			for _, e := range erroresNoEncontradoCompletar {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusNotFound, respuestaError{e.Error()})
					return
				}
			}
			for _, e := range erroresDeValidacionCompletar {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("completar elemento %d: %v", elementoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al completar el elemento"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
