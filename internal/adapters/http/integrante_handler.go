package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// agregadorIntegrantes es lo que el handler necesita del caso de uso AgregarIntegrante.
type agregadorIntegrantes interface {
	Ejecutar(ctx context.Context, d app.DatosIntegrante) error
}

// quitadorIntegrantes es lo que el handler necesita del caso de uso QuitarIntegrante.
type quitadorIntegrantes interface {
	Ejecutar(ctx context.Context, d app.DatosQuitarIntegrante) error
}

// agregadorRoles es lo que el handler necesita del caso de uso AgregarRol.
type agregadorRoles interface {
	Ejecutar(ctx context.Context, d app.DatosRol) error
}

// pedidoIntegrante es el JSON de entrada de POST /proyectos/{id}/integrantes.
type pedidoIntegrante struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// erroresDeValidacionIntegrante son los errores de HU-03 que se responden como 400.
var erroresDeValidacionIntegrante = []error{
	domain.ErrNombreIntegranteVacio,
	domain.ErrRolIntegranteVacio,
	domain.ErrIntegranteDuplicado,
}

// handlerAgregarIntegrante atiende POST /proyectos/{id}/integrantes: 204 si se agregó, 400 si
// los datos no son válidos, 404 si el proyecto no existe y 500 ante una falla interna.
// Implementa HU-03 (expone AgregarIntegrante por HTTP).
func handlerAgregarIntegrante(uc agregadorIntegrantes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoIntegrante
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		err = uc.Ejecutar(r.Context(), app.DatosIntegrante{ProyectoID: id, Nombre: pedido.Nombre, Rol: pedido.Rol})
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
				return
			}
			for _, e := range erroresDeValidacionIntegrante {
				if errors.Is(err, e) {
					escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
					return
				}
			}
			log.Printf("agregar integrante: %v", err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al agregar el integrante"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handlerQuitarIntegrante atiende DELETE /proyectos/{id}/integrantes/{nombre}: 204 si se quitó,
// 404 si el proyecto o el integrante no existen, 500 ante una falla interna. Implementa HU-03
// (expone QuitarIntegrante por HTTP).
func handlerQuitarIntegrante(uc quitadorIntegrantes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		nombre := r.PathValue("nombre")

		err = uc.Ejecutar(r.Context(), app.DatosQuitarIntegrante{ProyectoID: id, Nombre: nombre})
		if err != nil {
			if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrIntegranteNoEncontrado) {
				escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
				return
			}
			log.Printf("quitar integrante: %v", err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al quitar el integrante"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// pedidoRol es el JSON de entrada de POST /proyectos/{id}/integrantes/{nombre}/roles.
type pedidoRol struct {
	Rol string `json:"rol"`
}

// handlerAgregarRol atiende POST /proyectos/{id}/integrantes/{nombre}/roles: 204 si se agregó,
// 400 si el rol es vacío o repetido, 404 si el proyecto o el integrante no existen y 500 ante una
// falla interna. Implementa HU-03 (expone AgregarRol por HTTP).
func handlerAgregarRol(uc agregadorRoles) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoRol
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}

		err = uc.Ejecutar(r.Context(), app.DatosRol{ProyectoID: id, Nombre: r.PathValue("nombre"), Rol: pedido.Rol})
		switch {
		case err == nil:
			w.WriteHeader(http.StatusNoContent)
		case errors.Is(err, domain.ErrProyectoNoEncontrado), errors.Is(err, domain.ErrIntegranteNoEncontrado):
			escribirJSON(w, http.StatusNotFound, respuestaError{err.Error()})
		case errors.Is(err, domain.ErrRolIntegranteVacio):
			escribirJSON(w, http.StatusBadRequest, respuestaError{domain.ErrRolIntegranteVacio.Error()})
		case errors.Is(err, domain.ErrRolDuplicado):
			escribirJSON(w, http.StatusBadRequest, respuestaError{domain.ErrRolDuplicado.Error()})
		default:
			log.Printf("agregar rol: %v", err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al agregar el rol"})
		}
	}
}

// proyectoIDDesdeRuta parsea el {id} de la ruta como entero positivo.
func proyectoIDDesdeRuta(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("el id del proyecto debe ser un número entero positivo")
	}
	return id, nil
}
