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

// registradorDefecto es lo que el handler necesita del caso de uso RegistrarDefecto.
type registradorDefecto interface {
	Ejecutar(ctx context.Context, d app.DatosDefecto) (domain.Defecto, error)
}

// pedidoRegistrarDefecto es el JSON de entrada de POST /proyectos/{id}/defectos.
type pedidoRegistrarDefecto struct {
	ElementoID  int64  `json:"elemento_id"`
	Descripcion string `json:"descripcion"`
	Severidad   string `json:"severidad"`
}

type respuestaDefecto struct {
	ID                int64  `json:"id"`
	ProyectoID        int64  `json:"proyecto_id"`
	ElementoID        int64  `json:"elemento_id"`
	Descripcion       string `json:"descripcion"`
	Severidad         string `json:"severidad"`
	Estado            string `json:"estado"`
	SprintDeteccionID *int64 `json:"sprint_deteccion_id"`
}

// erroresNoEncontradoDefecto son los errores de HU-20 que se responden como 404: los datos que
// identifican entidades existentes (el proyecto de la ruta y el elemento del cuerpo).
var erroresNoEncontradoDefecto = []error{
	domain.ErrProyectoNoEncontrado,
	domain.ErrElementoNoEncontrado,
}

// erroresDeValidacionDefecto son los errores de HU-20 que se responden como 400.
var erroresDeValidacionDefecto = []error{
	domain.ErrDescripcionDefectoVacia,
	domain.ErrSeveridadInvalida,
}

// RegistrarRegistrarDefecto agrega la ruta de HU-20 al mux.
func RegistrarRegistrarDefecto(mux *http.ServeMux, uc registradorDefecto) {
	mux.HandleFunc("POST /proyectos/{id}/defectos", handlerRegistrarDefecto(uc))
}

// handlerRegistrarDefecto atiende POST /proyectos/{id}/defectos: 201 con el defecto creado, 404 si
// no existe el proyecto o el elemento (elemento_id va en el cuerpo, no en la ruta), 400 si los
// datos no son válidos (sin llamar al caso de uso si el pedido está mal formado) y 500 ante una
// falla interna. Implementa HU-20 (CA-20.6).
func handlerRegistrarDefecto(uc registradorDefecto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoRegistrarDefecto
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		def, err := uc.Ejecutar(r.Context(), app.DatosDefecto{
			ProyectoID:  proyectoID,
			ElementoID:  pedido.ElementoID,
			Descripcion: pedido.Descripcion,
			Severidad:   domain.Severidad(pedido.Severidad),
		})
		if err != nil {
			for _, v := range erroresNoEncontradoDefecto {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusNotFound, respuestaError{v.Error()})
					return
				}
			}
			for _, v := range erroresDeValidacionDefecto {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{v.Error()})
					return
				}
			}
			log.Printf("registrar defecto en el elemento %d del proyecto %d: %v", pedido.ElementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al registrar el defecto"})
			return
		}

		escribirJSON(w, http.StatusCreated, aRespuestaDefecto(def))
	}
}

func aRespuestaDefecto(d domain.Defecto) respuestaDefecto {
	return respuestaDefecto{
		ID:                d.ID,
		ProyectoID:        d.ProyectoID,
		ElementoID:        d.ElementoID,
		Descripcion:       d.Descripcion,
		Severidad:         string(d.Severidad),
		Estado:            string(d.Estado),
		SprintDeteccionID: d.SprintDeteccionID,
	}
}
