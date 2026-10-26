package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// registradorEsfuerzo es lo que el handler necesita del caso de uso RegistrarEsfuerzo.
type registradorEsfuerzo interface {
	Ejecutar(ctx context.Context, d app.DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error)
}

// pedidoRegistrarEsfuerzo es el JSON de entrada de POST /proyectos/{id}/backlog/{elementoID}/esfuerzo.
// fecha ausente o vacía: el dominio responde que es obligatoria.
type pedidoRegistrarEsfuerzo struct {
	Integrante string  `json:"integrante"`
	Fecha      string  `json:"fecha"`
	Actividad  string  `json:"actividad"`
	Horas      float64 `json:"horas"`
}

type respuestaRegistroEsfuerzo struct {
	ID         int64   `json:"id"`
	ProyectoID int64   `json:"proyecto_id"`
	ElementoID int64   `json:"elemento_id"`
	Integrante string  `json:"integrante"`
	Fecha      string  `json:"fecha"`
	Actividad  string  `json:"actividad"`
	Horas      float64 `json:"horas"`
}

// erroresNoEncontradoEsfuerzo son los errores de HU-18 que se responden como 404: los datos de la ruta.
var erroresNoEncontradoEsfuerzo = []error{
	domain.ErrProyectoNoEncontrado,
	domain.ErrElementoNoEncontrado,
}

// erroresDeValidacionEsfuerzo son los errores de HU-18 que se responden como 400. Incluye
// ErrIntegranteNoEncontrado: el integrante viene en el cuerpo, no en la ruta (spec, sección 7).
var erroresDeValidacionEsfuerzo = []error{
	domain.ErrNombreIntegranteVacio,
	domain.ErrFechaEsfuerzoObligatoria,
	domain.ErrFechaEsfuerzoFutura,
	domain.ErrActividadVacia,
	domain.ErrHorasInvalidas,
	domain.ErrHorasNoMultiploMediaHora,
	domain.ErrHorasExcedenMaximo,
	domain.ErrIntegranteNoEncontrado,
}

// RegistrarRegistrarEsfuerzo agrega la ruta de HU-18 al mux.
func RegistrarRegistrarEsfuerzo(mux *http.ServeMux, uc registradorEsfuerzo) {
	mux.HandleFunc("POST /proyectos/{id}/backlog/{elementoID}/esfuerzo", handlerRegistrarEsfuerzo(uc))
}

// handlerRegistrarEsfuerzo atiende POST /proyectos/{id}/backlog/{elementoID}/esfuerzo: 201 con el
// registro creado, 404 si no existe el proyecto o el elemento, 400 si los datos no son válidos
// (sin llamar al caso de uso si el pedido está mal formado) y 500 ante una falla interna.
// Implementa HU-18 (CA-18.15).
func handlerRegistrarEsfuerzo(uc registradorEsfuerzo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		elementoID, err := strconv.ParseInt(r.PathValue("elementoID"), 10, 64)
		if err != nil || elementoID <= 0 {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el id del elemento debe ser un número entero positivo"})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoRegistrarEsfuerzo
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		var fecha time.Time
		if pedido.Fecha != "" {
			// RN-18.4: la fecha es un día de la zona del servidor, igual que "hoy".
			if fecha, err = time.ParseInLocation(time.DateOnly, pedido.Fecha, time.Local); err != nil {
				escribirJSON(w, http.StatusBadRequest, respuestaError{"la fecha debe tener el formato AAAA-MM-DD"})
				return
			}
		}

		reg, err := uc.Ejecutar(r.Context(), app.DatosRegistroEsfuerzo{
			ProyectoID: proyectoID,
			ElementoID: elementoID,
			Integrante: pedido.Integrante,
			Fecha:      fecha,
			Actividad:  pedido.Actividad,
			Horas:      pedido.Horas,
		})
		if err != nil {
			for _, v := range erroresNoEncontradoEsfuerzo {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusNotFound, respuestaError{v.Error()})
					return
				}
			}
			for _, v := range erroresDeValidacionEsfuerzo {
				if errors.Is(err, v) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{v.Error()})
					return
				}
			}
			log.Printf("registrar esfuerzo en el elemento %d del proyecto %d: %v", elementoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al registrar el esfuerzo"})
			return
		}

		escribirJSON(w, http.StatusCreated, respuestaRegistroEsfuerzo{
			ID:         reg.ID,
			ProyectoID: reg.ProyectoID,
			ElementoID: reg.ElementoID,
			Integrante: reg.Integrante,
			Fecha:      reg.Fecha.Format(time.DateOnly),
			Actividad:  reg.Actividad,
			Horas:      reg.Horas,
		})
	}
}
