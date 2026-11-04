package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// actualizadorEstadoDefecto es lo que el handler necesita del caso de uso ActualizarEstadoDefecto.
type actualizadorEstadoDefecto interface {
	Iniciar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
	Resolver(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
	Cerrar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
}

// erroresEstadoDefecto son los errores de HU-21 que se informan al cliente, con su código (spec,
// sección 7). Se responde con el mensaje del error del dominio y no con el que devolvió el caso de
// uso, igual que HU-11.
var erroresEstadoDefecto = []struct {
	err    error
	status int
}{
	{domain.ErrProyectoNoEncontrado, http.StatusNotFound},
	{domain.ErrDefectoNoEncontrado, http.StatusNotFound},
	{domain.ErrDefectoNoAbierto, http.StatusBadRequest},
	{domain.ErrDefectoNoEnProgreso, http.StatusBadRequest},
	{domain.ErrDefectoNoResuelto, http.StatusBadRequest},
}

// RegistrarActualizarEstadoDefecto agrega al mux las tres rutas de acción de HU-21.
//
// La llamada a uc.Iniciar/uc.Resolver/uc.Cerrar se hace recién al atender cada pedido (nunca acá,
// al registrar), igual que todos los demás handlers: extraer un method value de una interfaz nil
// (por ejemplo cuando un escenario BDD arma el router con un CasosDeUso parcial para otra historia)
// hace panic de inmediato, aunque esa ruta nunca se use.
func RegistrarActualizarEstadoDefecto(mux *http.ServeMux, uc actualizadorEstadoDefecto) {
	mux.HandleFunc("POST /proyectos/{id}/defectos/{defectoID}/iniciar", handlerAccionEstadoDefecto(uc, "iniciar"))
	mux.HandleFunc("POST /proyectos/{id}/defectos/{defectoID}/resolver", handlerAccionEstadoDefecto(uc, "resolver"))
	mux.HandleFunc("POST /proyectos/{id}/defectos/{defectoID}/cerrar", handlerAccionEstadoDefecto(uc, "cerrar"))
}

// handlerAccionEstadoDefecto atiende iniciar, resolver o cerrar: 204 si se actualizó, 400 si los
// IDs no son válidos o el defecto no estaba en el estado de origen requerido, 404 si el proyecto o
// el defecto no existen, 500 ante una falla interna. Implementa HU-21 (CA-21.9).
func handlerAccionEstadoDefecto(uc actualizadorEstadoDefecto, accion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		defectoID, err := defectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		datos := app.DatosActualizarEstadoDefecto{ProyectoID: proyectoID, DefectoID: defectoID}
		switch accion {
		case "iniciar":
			_, err = uc.Iniciar(r.Context(), datos)
		case "resolver":
			_, err = uc.Resolver(r.Context(), datos)
		case "cerrar":
			_, err = uc.Cerrar(r.Context(), datos)
		}
		if err != nil {
			for _, e := range erroresEstadoDefecto {
				if errors.Is(err, e.err) {
					escribirJSON(w, e.status, respuestaError{e.err.Error()})
					return
				}
			}
			log.Printf("%s defecto %d del proyecto %d: %v", accion, defectoID, proyectoID, err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al actualizar el defecto"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// defectoIDDesdeRuta lee {defectoID} de la ruta, igual que proyectoIDDesdeRuta lee {id}.
func defectoIDDesdeRuta(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("defectoID"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("el id del defecto debe ser un número entero positivo")
	}
	return id, nil
}
