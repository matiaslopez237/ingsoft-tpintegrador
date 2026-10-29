package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// editorElementos es lo que el handler necesita del caso de uso EditarElementoBacklog.
type editorElementos interface {
	Ejecutar(ctx context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error)
}

// eliminadorElementos es lo que el handler necesita del caso de uso EliminarElementoBacklog.
type eliminadorElementos interface {
	Ejecutar(ctx context.Context, d app.DatosEliminarElemento) error
}

// pedidoEditarElemento es el JSON de PATCH /proyectos/{id}/backlog/{elementoID}. Un campo ausente
// no cambia. story_points es json.RawMessage para distinguir "ausente" (nil) de null (quitar la estimación).
type pedidoEditarElemento struct {
	Titulo      *string         `json:"titulo"`
	Descripcion *string         `json:"descripcion"`
	Prioridad   *string         `json:"prioridad"`
	StoryPoints json.RawMessage `json:"story_points"`
	Criterios   *[]string       `json:"criterios"`
	Estado      *string         `json:"estado"`
}

// erroresDeValidacionEdicion son los errores de HU-06 que se responden como 400.
var erroresDeValidacionEdicion = append([]error{domain.ErrSinCambios, domain.ErrEstadoElementoInvalido}, erroresDeValidacionBacklog...)

// RegistrarEditarEliminarElementoBacklog agrega las rutas de HU-06 al mux.
func RegistrarEditarEliminarElementoBacklog(mux *http.ServeMux, editar editorElementos, eliminar eliminadorElementos) {
	mux.HandleFunc("PATCH /proyectos/{id}/backlog/{elementoID}", handlerEditarElemento(editar))
	mux.HandleFunc("DELETE /proyectos/{id}/backlog/{elementoID}", handlerEliminarElemento(eliminar))
}

// idsElementoDesdeRuta lee el proyecto y el elemento de la ruta.
func idsElementoDesdeRuta(r *http.Request) (int64, int64, error) {
	proyectoID, err := proyectoIDDesdeRuta(r)
	if err != nil {
		return 0, 0, err
	}
	elementoID, err := elementoIDDesdeRuta(r)
	if err != nil {
		return 0, 0, err
	}
	return proyectoID, elementoID, nil
}

// handlerEditarElemento atiende PATCH /proyectos/{id}/backlog/{elementoID}: 200 con el elemento
// como quedó guardado, 400 si los datos no son válidos, 404 si el elemento no existe en el
// proyecto y 500 ante una falla interna. Implementa HU-06 (expone EditarElementoBacklog por HTTP).
func handlerEditarElemento(uc editorElementos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, elementoID, err := idsElementoDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoEditarElemento
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		cambios, err := pedido.aCambios()
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		e, err := uc.Ejecutar(r.Context(), app.DatosEditarElemento{ProyectoID: proyectoID, ElementoID: elementoID, Cambios: cambios})
		if err != nil {
			if errors.Is(err, domain.ErrElementoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrElementoNoEncontrado.Error()})
				return
			}
			for _, v := range erroresDeValidacionEdicion {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{v.Error()})
					return
				}
			}
			log.Printf("editar elemento %d del proyecto %d: %v", elementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al editar el elemento"})
			return
		}
		escribirJSON(w, http.StatusOK, aRespuestaElemento(e))
	}
}

// aCambios traduce el JSON a los cambios del dominio. story_points en null quita la estimación.
func (p pedidoEditarElemento) aCambios() (domain.CambiosElemento, error) {
	c := domain.CambiosElemento{Titulo: p.Titulo, Descripcion: p.Descripcion, Criterios: p.Criterios}
	if p.Prioridad != nil {
		pr := domain.Prioridad(*p.Prioridad)
		c.Prioridad = &pr
	}
	if p.Estado != nil {
		e := domain.EstadoElemento(*p.Estado)
		c.Estado = &e
	}
	if p.StoryPoints != nil {
		c.CambiarStoryPoints = true
		if !bytes.Equal(p.StoryPoints, []byte("null")) {
			var sp int
			if err := json.Unmarshal(p.StoryPoints, &sp); err != nil {
				return domain.CambiosElemento{}, errors.New("story_points debe ser un entero o null")
			}
			c.StoryPoints = &sp
		}
	}
	return c, nil
}

// handlerEliminarElemento atiende DELETE /proyectos/{id}/backlog/{elementoID}: 204 si se eliminó,
// 404 si no existe en el proyecto, 409 si tiene esfuerzo registrado y 500 ante una falla interna.
// Implementa HU-06 (expone EliminarElementoBacklog por HTTP).
func handlerEliminarElemento(uc eliminadorElementos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, elementoID, err := idsElementoDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		err = uc.Ejecutar(r.Context(), app.DatosEliminarElemento{ProyectoID: proyectoID, ElementoID: elementoID})
		switch {
		case err == nil:
			w.WriteHeader(http.StatusNoContent)
		case errors.Is(err, domain.ErrElementoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrElementoNoEncontrado.Error()})
		case errors.Is(err, domain.ErrElementoConEsfuerzo):
			escribirJSON(w, http.StatusConflict, respuestaError{domain.ErrElementoConEsfuerzo.Error()})
		case errors.Is(err, domain.ErrElementoConDefectos): // DEF-02
			escribirJSON(w, http.StatusConflict, respuestaError{domain.ErrElementoConDefectos.Error()})
		default:
			log.Printf("eliminar elemento %d del proyecto %d: %v", elementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al eliminar el elemento"})
		}
	}
}
