// Package http contiene los adaptadores HTTP: ruteo y handlers que exponen
// los casos de uso de internal/app sobre net/http (ver docs/adr/0001-framework-web.md).
package http

import "net/http"

// NuevoRouter arma el mux con todas las rutas de la aplicación.
// Cada historia de usuario agrega sus propias rutas acá a medida que se implementa.
func NuevoRouter(crearProyecto creadorProyectos) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /salud", handlerSalud)
	mux.HandleFunc("POST /proyectos", handlerCrearProyecto(crearProyecto)) // HU-01

	return mux
}

// handlerSalud es un endpoint mínimo para verificar que el servidor está vivo.
func handlerSalud(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"estado":"ok"}`))
}
