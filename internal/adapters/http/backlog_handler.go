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

// creadorElementos es lo que el handler necesita del caso de uso CrearElementoBacklog.
type creadorElementos interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoElemento) (domain.ElementoBacklog, error)
}

// pedidoCrearElemento es el JSON de entrada de POST /proyectos/{id}/backlog.
// story_points ausente o null: el elemento queda sin estimar.
type pedidoCrearElemento struct {
	Titulo      string   `json:"titulo"`
	Descripcion string   `json:"descripcion"`
	Prioridad   string   `json:"prioridad"`
	StoryPoints *int     `json:"story_points"`
	Criterios   []string `json:"criterios"`
}

type respuestaElemento struct {
	ID          int64    `json:"id"`
	ProyectoID  int64    `json:"proyecto_id"`
	Titulo      string   `json:"titulo"`
	Descripcion string   `json:"descripcion"`
	Prioridad   string   `json:"prioridad"`
	Estado      string   `json:"estado"`
	StoryPoints *int     `json:"story_points"`
	Criterios   []string `json:"criterios"`
}

// erroresDeValidacionBacklog son los errores de HU-05 que se responden como 400.
var erroresDeValidacionBacklog = []error{
	domain.ErrTituloVacio,
	domain.ErrTituloDemasiadoLargo,
	domain.ErrDescripcionVacia,
	domain.ErrPrioridadInvalida,
	domain.ErrStoryPointsInvalidos,
	domain.ErrSinCriterios,
	domain.ErrCriterioVacio,
	domain.ErrTituloDuplicado,
}

// RegistrarCrearElementoBacklog agrega la ruta de HU-05 al mux.
func RegistrarCrearElementoBacklog(mux *http.ServeMux, uc creadorElementos) {
	mux.HandleFunc("POST /proyectos/{id}/backlog", handlerCrearElementoBacklog(uc))
}

// handlerCrearElementoBacklog atiende POST /proyectos/{id}/backlog: 201 con el elemento creado,
// 400 si los datos no son válidos, 404 si el proyecto no existe y 500 ante una falla interna.
// Implementa HU-05 (expone CrearElementoBacklog por HTTP).
func handlerCrearElementoBacklog(uc creadorElementos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoCrearElemento
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		e, err := uc.Ejecutar(r.Context(), app.DatosNuevoElemento{
			ProyectoID:  id,
			Titulo:      pedido.Titulo,
			Descripcion: pedido.Descripcion,
			Prioridad:   domain.Prioridad(pedido.Prioridad),
			StoryPoints: pedido.StoryPoints,
			Criterios:   pedido.Criterios,
		})
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
				return
			}
			for _, v := range erroresDeValidacionBacklog {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{v.Error()})
					return
				}
			}
			log.Printf("crear elemento de backlog en el proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al crear el elemento"})
			return
		}

		escribirJSON(w, http.StatusCreated, aRespuestaElemento(e))
	}
}
