package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// comparadorEsfuerzo es lo que el handler necesita del caso de uso CompararEsfuerzo.
type comparadorEsfuerzo interface {
	Ejecutar(ctx context.Context, proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error)
}

// respuestaComparacion es el JSON de la comparación (spec HU-19, sección 3). Los punteros salen como
// null cuando no hay valor: sin estimar, o sin porcentaje con estimado 0.
type respuestaComparacion struct {
	ProyectoID           int64    `json:"proyecto_id"`
	ElementoID           int64    `json:"elemento_id"`
	HorasEstimadas       *float64 `json:"horas_estimadas"`
	HorasReales          float64  `json:"horas_reales"`
	DesviacionHoras      *float64 `json:"desviacion_horas"`
	DesviacionPorcentaje *float64 `json:"desviacion_porcentaje"`
}

// erroresComparacion son los errores de HU-19 que se informan al cliente (spec, sección 7). Se
// responde con el mensaje del dominio y no con el que devolvió el caso de uso, que puede venir envuelto.
var erroresComparacion = []error{domain.ErrProyectoNoEncontrado, domain.ErrElementoNoEncontrado}

// RegistrarCompararEsfuerzo agrega al mux la ruta de HU-19: comparar el esfuerzo de un elemento.
func RegistrarCompararEsfuerzo(mux *http.ServeMux, uc comparadorEsfuerzo) {
	mux.HandleFunc("GET /proyectos/{id}/backlog/{elementoID}/esfuerzo/comparacion", handlerCompararEsfuerzo(uc))
}

// handlerCompararEsfuerzo atiende GET /proyectos/{id}/backlog/{elementoID}/esfuerzo/comparacion: 200
// con la comparación, 400 si los IDs no son válidos, 404 si el proyecto o el elemento no existen y
// 500 ante una falla interna. Implementa HU-19 (CA-19.17).
func handlerCompararEsfuerzo(uc comparadorEsfuerzo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, elementoID, err := idsElementoDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		c, err := uc.Ejecutar(r.Context(), proyectoID, elementoID)
		if err != nil {
			for _, e := range erroresComparacion {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusNotFound, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("comparar esfuerzo del elemento %d del proyecto %d: %v", elementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al comparar el esfuerzo"})
			return
		}

		escribirJSON(w, http.StatusOK, respuestaComparacion{
			ProyectoID:           proyectoID,
			ElementoID:           elementoID,
			HorasEstimadas:       c.HorasEstimadas,
			HorasReales:          c.HorasReales,
			DesviacionHoras:      c.DesviacionHoras,
			DesviacionPorcentaje: c.DesviacionPorcentaje,
		})
	}
}
