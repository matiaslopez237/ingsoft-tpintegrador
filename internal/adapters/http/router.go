// Package http contiene los adaptadores HTTP: ruteo y handlers que exponen
// los casos de uso de internal/app sobre net/http (ver docs/adr/0001-framework-web.md).
package http

import "net/http"

// NuevoRouter arma el mux con todas las rutas de la aplicación.
// Cada historia de usuario agrega sus propias rutas acá a medida que se implementa.
func NuevoRouter(
	crearProyecto creadorProyectos,
	agregarIntegrante agregadorIntegrantes,
	quitarIntegrante quitadorIntegrantes,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /salud", handlerSalud)
	mux.HandleFunc("POST /proyectos", handlerCrearProyecto(crearProyecto))                                   // HU-01
	mux.HandleFunc("POST /proyectos/{id}/integrantes", handlerAgregarIntegrante(agregarIntegrante))          // HU-03
	mux.HandleFunc("DELETE /proyectos/{id}/integrantes/{nombre}", handlerQuitarIntegrante(quitarIntegrante)) // HU-03

	return mux
}

// handlerSalud es un endpoint mínimo para verificar que el servidor está vivo.
func handlerSalud(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
}
