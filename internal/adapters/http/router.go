// Package http contiene los adaptadores HTTP: ruteo y handlers que exponen
// los casos de uso de internal/app sobre net/http (ver docs/adr/0001-framework-web.md).
package http

import "net/http"

// CasosDeUso son los casos de uso que exponen las rutas HTTP. Cada historia de usuario suma
// los suyos acá y registra sus rutas con su propia función RegistrarX, en el archivo de su handler.
type CasosDeUso struct {
	CrearProyecto            creadorProyectos     // HU-01
	ModificarProyecto        modificadorProyectos // HU-02
	AgregarIntegrante        agregadorIntegrantes // HU-03
	QuitarIntegrante         quitadorIntegrantes  // HU-03
	AgregarRol               agregadorRoles       // HU-03
	ConsultarProyecto        consultorProyectos   // HU-04
	CrearElementoBacklog     creadorElementos     // HU-05
	RegistrarEsfuerzo        registradorEsfuerzo  // HU-18
	CrearSprint              creadorSprints       // HU-08
	ListarSprints            listadorSprints      // HU-08
	ListarBacklog            listadorBacklog      // HU-07
	AsignarHistoria          agregadorAsignacion  // HU-09
	DesasignarHistoria       quitadorAsignacion   // HU-09
	ListarProyectos          listadorProyectos    // HU-30
	EditarElementoBacklog    editorElementos      // HU-06
	EliminarElementoBacklog  eliminadorElementos  // HU-06
	ConsultarElementoBacklog consultorElementos   // HU-06
}

// NuevoRouter arma el mux con todas las rutas de la aplicación. Lo usan main.go y el test
// de rutas (router_test.go), así ninguna ruta queda sin registrar.
func NuevoRouter(cu CasosDeUso) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /salud", handlerSalud)
	RegistrarCrearProyecto(mux, cu.CrearProyecto)
	RegistrarModificarProyecto(mux, cu.ModificarProyecto)
	RegistrarIntegrantes(mux, cu.AgregarIntegrante, cu.QuitarIntegrante, cu.AgregarRol)
	RegistrarConsultarProyecto(mux, cu.ConsultarProyecto)
	RegistrarCrearElementoBacklog(mux, cu.CrearElementoBacklog)
	RegistrarEditarEliminarElementoBacklog(mux, cu.EditarElementoBacklog, cu.EliminarElementoBacklog)
	RegistrarRegistrarEsfuerzo(mux, cu.RegistrarEsfuerzo)
	RegistrarSprints(mux, cu.CrearSprint, cu.ListarSprints)
	RegistrarListarBacklog(mux, cu.ListarBacklog)
	RegistrarAsignacion(mux, cu.AsignarHistoria, cu.DesasignarHistoria)
	RegistrarPantallas(mux, cu)
	return mux
}

// handlerSalud es un endpoint mínimo para verificar que el servidor está vivo.
func handlerSalud(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
}
