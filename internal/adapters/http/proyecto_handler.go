package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// creadorProyectos es lo que el handler necesita del caso de uso CrearProyecto.
type creadorProyectos interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoProyecto) (domain.Proyecto, error)
}

// pedidoCrearProyecto es el JSON de entrada de POST /proyectos. Fechas en formato AAAA-MM-DD.
type pedidoCrearProyecto struct {
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

type respuestaProyecto struct {
	ID          int64   `json:"id"`
	Nombre      string  `json:"nombre"`
	FechaInicio string  `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
	Estado      string  `json:"estado"`
}

type respuestaError struct {
	Error string `json:"error"`
}

// tamanoMaximoCuerpo limita el JSON de entrada (1 MiB) para no leer cuerpos enormes en memoria.
const tamanoMaximoCuerpo = 1 << 20

// erroresDeValidacion son los errores de HU-01 que se responden como 400.
var erroresDeValidacion = []error{
	domain.ErrNombreVacio,
	domain.ErrNombreDemasiadoLargo,
	domain.ErrNombreDuplicado,
	domain.ErrFechaInicioObligatoria,
	domain.ErrFechaFinAnterior,
}

// handlerCrearProyecto atiende POST /proyectos: 201 con el proyecto creado, 400 si los datos
// no son válidos y 500 ante una falla interna.
func handlerCrearProyecto(uc creadorProyectos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoCrearProyecto
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		datos, err := pedido.aDatos()
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		p, err := uc.Ejecutar(r.Context(), datos)
		if err != nil {
			for _, e := range erroresDeValidacion {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("crear proyecto: %v", err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al crear el proyecto"})
			return
		}

		escribirJSON(w, http.StatusCreated, aRespuestaProyecto(p))
	}
}

// aDatos convierte las fechas del JSON. Una fecha vacía se toma como no informada.
func (p pedidoCrearProyecto) aDatos() (app.DatosNuevoProyecto, error) {
	datos := app.DatosNuevoProyecto{Nombre: p.Nombre}
	if p.FechaInicio != "" {
		f, err := time.Parse(time.DateOnly, p.FechaInicio)
		if err != nil {
			return app.DatosNuevoProyecto{}, errors.New("fecha_inicio debe tener el formato AAAA-MM-DD")
		}
		datos.FechaInicio = f
	}
	if p.FechaFin != "" {
		f, err := time.Parse(time.DateOnly, p.FechaFin)
		if err != nil {
			return app.DatosNuevoProyecto{}, errors.New("fecha_fin debe tener el formato AAAA-MM-DD")
		}
		datos.FechaFin = &f
	}
	return datos, nil
}

func aRespuestaProyecto(p domain.Proyecto) respuestaProyecto {
	resp := respuestaProyecto{
		ID:          p.ID,
		Nombre:      p.Nombre,
		FechaInicio: p.FechaInicio.Format(time.DateOnly),
		Estado:      string(p.Estado),
	}
	if p.FechaFin != nil {
		fin := p.FechaFin.Format(time.DateOnly)
		resp.FechaFin = &fin
	}
	return resp
}

func escribirJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("escribir respuesta JSON: %v", err)
	}
}
