package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// listadorBacklog es lo que el handler necesita del caso de uso ListarBacklog.
type listadorBacklog interface {
	Ejecutar(ctx context.Context, d app.DatosListarBacklog) ([]domain.ElementoBacklog, error)
}

// erroresDeValidacionListado son los errores de HU-07 que se responden como 400.
var erroresDeValidacionListado = []error{domain.ErrOrdenInvalido, domain.ErrEstadoElementoInvalido}

// RegistrarListarBacklog agrega GET /proyectos/{id}/backlog al mux. Implementa HU-07.
func RegistrarListarBacklog(mux *http.ServeMux, uc listadorBacklog) {
	mux.HandleFunc("GET /proyectos/{id}/backlog", handlerListarBacklog(uc))
}

// handlerListarBacklog atiende GET /proyectos/{id}/backlog?orden=…&estado=…: 200 con los elementos
// (siempre un array), 400 si el orden o el estado no son válidos, 404 si el proyecto no existe y
// 500 ante una falla interna. Implementa HU-07 (expone ListarBacklog por HTTP).
func handlerListarBacklog(uc listadorBacklog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		q := r.URL.Query()
		elementos, err := uc.Ejecutar(r.Context(), app.DatosListarBacklog{
			ProyectoID: id, Orden: q.Get("orden"), Estado: q.Get("estado"),
		})
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
				return
			}
			for _, e := range erroresDeValidacionListado {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("listar el backlog del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al listar el backlog"})
			return
		}

		respuesta := make([]respuestaElemento, 0, len(elementos))
		for _, e := range elementos {
			respuesta = append(respuesta, aRespuestaElemento(e))
		}
		escribirJSON(w, http.StatusOK, respuesta)
	}
}

// aRespuestaElemento arma el JSON de un elemento con el formato de HU-05. Los criterios salen
// siempre como array.
func aRespuestaElemento(e domain.ElementoBacklog) respuestaElemento {
	criterios := e.Criterios
	if criterios == nil {
		criterios = []string{}
	}
	return respuestaElemento{
		ID:          e.ID,
		ProyectoID:  e.ProyectoID,
		Titulo:      e.Titulo,
		Descripcion: e.Descripcion,
		Prioridad:   string(e.Prioridad),
		Estado:      string(e.Estado),
		StoryPoints: e.StoryPoints,
		Criterios:   criterios,
	}
}
