package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// agregadorAsignacion es lo que el handler necesita del caso de uso AsignarHistoria.
type agregadorAsignacion interface {
	Ejecutar(ctx context.Context, d app.DatosAsignacion) error
}

// quitadorAsignacion es lo que el handler necesita del caso de uso DesasignarHistoria.
type quitadorAsignacion interface {
	Ejecutar(ctx context.Context, d app.DatosDesasignacion) error
}

// pedidoAsignacion es el JSON de entrada de POST /proyectos/{id}/backlog/{elementoID}/sprint.
type pedidoAsignacion struct {
	SprintID int64 `json:"sprint_id"`
}

// erroresDeValidacionAsignacion son los errores de HU-09 que se responden como 400.
var erroresDeValidacionAsignacion = []error{
	domain.ErrSprintDeOtroProyecto,
	domain.ErrElementoYaAsignado,
	domain.ErrSprintCerrado,
	domain.ErrElementoSinSprint,
}

// RegistrarAsignacion agrega al mux las rutas de HU-09: asignar y desasignar un elemento del
// backlog a/de un sprint.
func RegistrarAsignacion(mux *http.ServeMux, asignar agregadorAsignacion, desasignar quitadorAsignacion) {
	mux.HandleFunc("POST /proyectos/{id}/backlog/{elementoID}/sprint", handlerAsignarHistoria(asignar))
	mux.HandleFunc("DELETE /proyectos/{id}/backlog/{elementoID}/sprint", handlerDesasignarHistoria(desasignar))
}

// handlerAsignarHistoria atiende POST /proyectos/{id}/backlog/{elementoID}/sprint: 204 si se
// asignó, 400 si los datos no son válidos, 404 si el elemento o el sprint no existen, 500 ante
// una falla interna. Implementa HU-09 (expone AsignarHistoria por HTTP).
func handlerAsignarHistoria(uc agregadorAsignacion) http.HandlerFunc {
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

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoAsignacion
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		if pedido.SprintID <= 0 {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el sprint_id debe ser un número entero positivo"})
			return
		}

		err = uc.Ejecutar(r.Context(), app.DatosAsignacion{ProyectoID: proyectoID, ElementoID: elementoID, SprintID: pedido.SprintID})
		if err != nil {
			if errors.Is(err, domain.ErrElementoNoEncontrado) || errors.Is(err, domain.ErrSprintNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
				return
			}
			for _, e := range erroresDeValidacionAsignacion {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("asignar elemento %d a un sprint: %v", elementoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al asignar el elemento"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handlerDesasignarHistoria atiende DELETE /proyectos/{id}/backlog/{elementoID}/sprint: 204 si se
// desasignó, 400 si no se puede, 404 si el elemento no existe, 500 ante una falla interna.
// Implementa HU-09 (expone DesasignarHistoria por HTTP).
func handlerDesasignarHistoria(uc quitadorAsignacion) http.HandlerFunc {
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

		err = uc.Ejecutar(r.Context(), app.DatosDesasignacion{ProyectoID: proyectoID, ElementoID: elementoID})
		if err != nil {
			if errors.Is(err, domain.ErrElementoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
				return
			}
			for _, e := range erroresDeValidacionAsignacion {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("desasignar elemento %d: %v", elementoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al desasignar el elemento"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
