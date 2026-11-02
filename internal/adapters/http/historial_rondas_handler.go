package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// historiadorRondas es lo que el handler y la pantalla del elemento necesitan del caso de uso
// ListarRondas (HU-16).
type historiadorRondas interface {
	Ejecutar(ctx context.Context, proyectoID, elementoID int64) ([]domain.ResumenRonda, error)
}

// RegistrarHistorialRondas agrega al mux la ruta de HU-16. Comparte la ruta con el POST de HU-14
// (abrir una ronda), con otro método.
func RegistrarHistorialRondas(mux *http.ServeMux, uc historiadorRondas) {
	mux.HandleFunc("GET /proyectos/{id}/backlog/{elementoID}/rondas", handlerHistorialRondas(uc))
}

// handlerHistorialRondas atiende GET /proyectos/{id}/backlog/{elementoID}/rondas: 200 con las rondas
// del elemento, de la más reciente a la más antigua, cada una con la forma del resumen de HU-14 y
// HU-15 (siempre un array, [] si no hay); 400 si un ID no es un número; 404 si no existe el
// proyecto o el elemento; 500 ante una falla interna, sin el detalle. Implementa HU-16 (CA-16.7).
func handlerHistorialRondas(uc historiadorRondas) http.HandlerFunc {
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
		historial, err := uc.Ejecutar(r.Context(), proyectoID, elementoID)
		if err != nil {
			for _, e := range []error{domain.ErrProyectoNoEncontrado, domain.ErrElementoNoEncontrado} {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusNotFound, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("historial de rondas del elemento %d del proyecto %d: %v", elementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al listar las rondas"})
			return
		}
		respuesta := make([]respuestaRonda, 0, len(historial))
		for _, ronda := range historial {
			respuesta = append(respuesta, aRespuestaRonda(ronda))
		}
		escribirJSON(w, http.StatusOK, respuesta)
	}
}
