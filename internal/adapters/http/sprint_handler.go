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

// creadorSprints es lo que el handler necesita del caso de uso CrearSprint.
type creadorSprints interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoSprint) (domain.Sprint, error)
}

// listadorSprints es lo que el handler necesita del caso de uso ListarSprints.
type listadorSprints interface {
	Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Sprint, error)
}

// pedidoCrearSprint es el JSON de entrada de POST /proyectos/{id}/sprints. Fechas AAAA-MM-DD;
// una fecha vacía o ausente se toma como no informada.
type pedidoCrearSprint struct {
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Objetivo    string `json:"objetivo"`
}

type respuestaSprint struct {
	ID          int64  `json:"id"`
	ProyectoID  int64  `json:"proyecto_id"`
	Numero      int    `json:"numero"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Objetivo    string `json:"objetivo"`
	Estado      string `json:"estado"`
}

// erroresDeValidacionSprint son los errores de HU-08 que se responden como 400.
var erroresDeValidacionSprint = []error{
	domain.ErrProyectoNoActivo,
	domain.ErrFechaInicioObligatoria,
	domain.ErrFechaFinObligatoria,
	domain.ErrFechaFinAnterior,
	domain.ErrObjetivoVacio,
	domain.ErrObjetivoDemasiadoLargo,
	domain.ErrSprintFueraDelProyecto,
	domain.ErrSprintSuperpuesto,
}

// RegistrarSprints agrega al mux las rutas de HU-08: crear y listar los sprints de un proyecto.
func RegistrarSprints(mux *http.ServeMux, crear creadorSprints, listar listadorSprints) {
	mux.HandleFunc("POST /proyectos/{id}/sprints", handlerCrearSprint(crear))
	mux.HandleFunc("GET /proyectos/{id}/sprints", handlerListarSprints(listar))
}

// handlerCrearSprint atiende POST /proyectos/{id}/sprints: 201 con el sprint creado, 400 si los
// datos no son válidos, 404 si el proyecto no existe y 500 ante una falla interna.
// Implementa HU-08 (expone CrearSprint por HTTP).
func handlerCrearSprint(uc creadorSprints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoCrearSprint
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		datos, err := pedido.aDatos(id)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		s, err := uc.Ejecutar(r.Context(), datos)
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
				return
			}
			for _, e := range erroresDeValidacionSprint {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("crear sprint en el proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al crear el sprint"})
			return
		}
		escribirJSON(w, http.StatusCreated, aRespuestaSprint(s))
	}
}

// handlerListarSprints atiende GET /proyectos/{id}/sprints: 200 con los sprints ordenados por
// número (siempre un array, [] si no hay), 404 si el proyecto no existe y 500 ante una falla
// interna. Implementa HU-08 (CA-08.9).
func handlerListarSprints(uc listadorSprints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		sprints, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("listar sprints del proyecto %d: %v", id, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al listar los sprints"})
		default:
			respuesta := make([]respuestaSprint, 0, len(sprints))
			for _, s := range sprints {
				respuesta = append(respuesta, aRespuestaSprint(s))
			}
			escribirJSON(w, http.StatusOK, respuesta)
		}
	}
}

// aDatos convierte las fechas del JSON. Una fecha vacía queda en cero: el dominio la rechaza.
func (p pedidoCrearSprint) aDatos(proyectoID int64) (app.DatosNuevoSprint, error) {
	datos := app.DatosNuevoSprint{ProyectoID: proyectoID, Objetivo: p.Objetivo}
	var err error
	if datos.FechaInicio, err = fechaDePedido(p.FechaInicio, "fecha_inicio"); err != nil {
		return app.DatosNuevoSprint{}, err
	}
	if datos.FechaFin, err = fechaDePedido(p.FechaFin, "fecha_fin"); err != nil {
		return app.DatosNuevoSprint{}, err
	}
	return datos, nil
}

// fechaDePedido convierte una fecha AAAA-MM-DD; vacía devuelve la fecha cero.
func fechaDePedido(valor, campo string) (time.Time, error) {
	if valor == "" {
		return time.Time{}, nil
	}
	f, err := time.Parse(time.DateOnly, valor)
	if err != nil {
		return time.Time{}, errors.New(campo + " debe tener el formato AAAA-MM-DD")
	}
	return f, nil
}

func aRespuestaSprint(s domain.Sprint) respuestaSprint {
	return respuestaSprint{
		ID:          s.ID,
		ProyectoID:  s.ProyectoID,
		Numero:      s.Numero,
		FechaInicio: s.FechaInicio.Format(time.DateOnly),
		FechaFin:    s.FechaFin.Format(time.DateOnly),
		Objetivo:    s.Objetivo,
		Estado:      string(s.Estado),
	}
}
