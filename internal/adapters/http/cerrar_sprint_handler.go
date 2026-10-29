package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// cerradorSprints es lo que el handler necesita del caso de uso CerrarSprint.
type cerradorSprints interface {
	Ejecutar(ctx context.Context, d app.DatosCerrarSprint) error
}

// erroresCierreSprint son los errores de HU-11 que se informan al cliente, con su código (spec,
// sección 7). Se responde con el mensaje del error del dominio y no con el que devolvió el caso de
// uso: si dos pedidos cierran a la vez, ErrSprintYaCerrado llega envuelto ("cerrar sprint: ...").
var erroresCierreSprint = []struct {
	err    error
	status int
}{
	{domain.ErrProyectoNoEncontrado, http.StatusNotFound},
	{domain.ErrSprintNoEncontrado, http.StatusNotFound},
	{domain.ErrSprintYaCerrado, http.StatusBadRequest},
}

// RegistrarCerrarSprint agrega al mux la ruta de HU-11: cerrar un sprint.
func RegistrarCerrarSprint(mux *http.ServeMux, uc cerradorSprints) {
	mux.HandleFunc("POST /proyectos/{id}/sprints/{sprintID}/cerrar", handlerCerrarSprint(uc))
}

// handlerCerrarSprint atiende POST /proyectos/{id}/sprints/{sprintID}/cerrar: 204 si se cerró, 400
// si los IDs no son válidos o el sprint ya estaba cerrado, 404 si el proyecto o el sprint no
// existen, 500 ante una falla interna. Implementa HU-11 (CA-11.4).
func handlerCerrarSprint(uc cerradorSprints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		sprintID, err := sprintIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		err = uc.Ejecutar(r.Context(), app.DatosCerrarSprint{ProyectoID: proyectoID, SprintID: sprintID})
		if err != nil {
			for _, e := range erroresCierreSprint {
				if errors.Is(err, e.err) {
					escribirJSON(w, e.status, respuestaError{e.err.Error()})
					return
				}
			}
			log.Printf("cerrar sprint %d del proyecto %d: %v", sprintID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al cerrar el sprint"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
