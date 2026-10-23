package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// modificadorProyectos es lo que el handler necesita del caso de uso ModificarProyecto.
type modificadorProyectos interface {
	Ejecutar(ctx context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error)
}

// pedidoModificarProyecto es el JSON de entrada de PATCH /proyectos/{id}. Un campo ausente no
// cambia. Las fechas son json.RawMessage para distinguir "ausente" (nil) de "null" (quitarla).
type pedidoModificarProyecto struct {
	Nombre      *string         `json:"nombre"`
	FechaInicio json.RawMessage `json:"fecha_inicio"`
	FechaFin    json.RawMessage `json:"fecha_fin"`
	Estado      *string         `json:"estado"`
}

// erroresDeValidacionModificar son los errores de HU-02 que se responden como 400.
var erroresDeValidacionModificar = append([]error{domain.ErrSinCambios, domain.ErrEstadoInvalido}, erroresDeValidacion...)

// handlerModificarProyecto atiende PATCH /proyectos/{id}: 200 con el proyecto modificado, 400 si
// los datos no son válidos, 404 si el proyecto no existe y 500 ante una falla interna.
// Implementa HU-02 (expone ModificarProyecto por HTTP).
func handlerModificarProyecto(uc modificadorProyectos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoModificarProyecto
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		cambios, err := pedido.aCambios()
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		p, err := uc.Ejecutar(r.Context(), app.DatosModificarProyecto{ProyectoID: id, Cambios: cambios})
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
				return
			}
			for _, e := range erroresDeValidacionModificar {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("modificar proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al modificar el proyecto"})
			return
		}

		escribirJSON(w, http.StatusOK, aRespuestaProyecto(p))
	}
}

// aCambios traduce el JSON a los cambios del dominio. Una fecha de inicio en null queda como
// fecha cero, que el dominio rechaza (RN-02.4); una fecha de fin en null la quita (RN-02.5).
func (p pedidoModificarProyecto) aCambios() (domain.CambiosProyecto, error) {
	c := domain.CambiosProyecto{Nombre: p.Nombre}
	if p.FechaInicio != nil {
		f, err := fechaOpcional(p.FechaInicio, "fecha_inicio")
		if err != nil {
			return domain.CambiosProyecto{}, err
		}
		if f == nil {
			f = &time.Time{}
		}
		c.FechaInicio = f
	}
	if p.FechaFin != nil {
		f, err := fechaOpcional(p.FechaFin, "fecha_fin")
		if err != nil {
			return domain.CambiosProyecto{}, err
		}
		c.CambiarFechaFin = true
		c.FechaFin = f
	}
	if p.Estado != nil {
		e := domain.Estado(*p.Estado)
		c.Estado = &e
	}
	return c, nil
}

// fechaOpcional convierte una fecha JSON AAAA-MM-DD; null devuelve nil.
func fechaOpcional(raw json.RawMessage, campo string) (*time.Time, error) {
	if bytes.Equal(raw, []byte("null")) {
		return nil, nil
	}
	var texto string
	if err := json.Unmarshal(raw, &texto); err != nil {
		return nil, fmt.Errorf("%s debe tener el formato AAAA-MM-DD", campo)
	}
	f, err := time.Parse(time.DateOnly, texto)
	if err != nil {
		return nil, fmt.Errorf("%s debe tener el formato AAAA-MM-DD", campo)
	}
	return &f, nil
}
