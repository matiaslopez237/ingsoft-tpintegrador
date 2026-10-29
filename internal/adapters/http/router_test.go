package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Cada historia registra sus rutas con su propia función RegistrarX; este test verifica que
// NuevoRouter (lo que usa main.go) las registra todas. Si una historia agrega una ruta y se
// olvida de sumarla a NuevoRouter, el pedido da 404 o 405 y el test falla.
func TestNuevoRouter_RegistraTodasLasRutas(t *testing.T) {
	mux := NuevoRouter(CasosDeUso{
		CrearProyecto:            &creadorFalso{},
		ModificarProyecto:        &modificadorFalso{},
		AgregarIntegrante:        &agregadorFalso{},
		QuitarIntegrante:         &quitadorFalso{},
		AgregarRol:               &agregadorRolFalso{},
		ConsultarProyecto:        &consultorFalso{},
		CrearElementoBacklog:     &creadorElementosFalso{},
		RegistrarEsfuerzo:        &registradorEsfuerzoFalso{},
		CrearSprint:              &creadorSprintsFalso{},
		ListarSprints:            &listadorSprintsFalso{sprints: []domain.Sprint{{ID: 1, ProyectoID: 1}}}, // para que /web/.../sprints/1 no dé 404
		ListarBacklog:            &listadorBacklogFalso{},
		AsignarHistoria:          &asignarHistoriaFalso{},
		DesasignarHistoria:       &desasignarHistoriaFalso{},
		CerrarSprint:             &cerrarSprintFalso{},
		CompletarHistoria:        &completarHistoriaFalso{},
		ListarProyectos:          casosFalsos().ListarProyectos,
		EditarElementoBacklog:    &editorElementosFalso{},
		EliminarElementoBacklog:  &eliminadorElementosFalso{},
		ConsultarElementoBacklog: casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{}).ConsultarElementoBacklog,
	})
	rutas := []struct{ metodo, ruta, cuerpo string }{
		{http.MethodGet, "/salud", ""},
		{http.MethodPost, "/proyectos", `{"nombre":"Alfa","fecha_inicio":"2026-03-01"}`},                                                     // HU-01
		{http.MethodPatch, "/proyectos/1", `{"nombre":"Beta"}`},                                                                              // HU-02
		{http.MethodPost, "/proyectos/1/integrantes", `{"nombre":"Ana","rol":"Tester"}`},                                                     // HU-03
		{http.MethodDelete, "/proyectos/1/integrantes/Ana", ""},                                                                              // HU-03
		{http.MethodPost, "/proyectos/1/integrantes/Ana/roles", `{"rol":"Tester"}`},                                                          // HU-03
		{http.MethodGet, "/proyectos/1", ""},                                                                                                 // HU-04
		{http.MethodPost, "/proyectos/1/backlog", `{"titulo":"Login","descripcion":"x","prioridad":"Alta","criterios":["a"]}`},               // HU-05
		{http.MethodPatch, "/proyectos/1/backlog/7", `{"prioridad":"Baja"}`},                                                                 // HU-06
		{http.MethodDelete, "/proyectos/1/backlog/7", ""},                                                                                    // HU-06
		{http.MethodPost, "/proyectos/1/backlog/1/esfuerzo", `{"integrante":"Ana","fecha":"2026-10-25","actividad":"Desarrollo","horas":2}`}, // HU-18
		{http.MethodPost, "/proyectos/1/sprints", `{"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"x"}`},                   // HU-08
		{http.MethodGet, "/proyectos/1/sprints", ""},                                                                                         // HU-08
		{http.MethodGet, "/proyectos/1/backlog", ""},                                                                                         // HU-07
		{http.MethodPost, "/proyectos/1/backlog/1/sprint", `{"sprint_id":1}`},                                                                // HU-09
		{http.MethodDelete, "/proyectos/1/backlog/1/sprint", ""},                                                                             // HU-09
		{http.MethodPost, "/proyectos/1/backlog/1/completar", ""},                                                                            // HU-10
		{http.MethodPost, "/proyectos/1/sprints/1/cerrar", ""},                                                                               // HU-11
		{http.MethodGet, "/", ""},                                     // HU-30
		{http.MethodGet, "/web/proyectos", ""},                        // HU-30
		{http.MethodGet, "/static/estilos.css", ""},                   // HU-30
		{http.MethodGet, "/web/proyectos/nuevo", ""},                  // HU-30
		{http.MethodGet, "/web/proyectos/1/sprints", ""},              // HU-08
		{http.MethodPost, "/web/proyectos/1/sprints/1/cerrar", ""},    // HU-11
		{http.MethodGet, "/web/proyectos/1/backlog", ""},              // HU-07
		{http.MethodGet, "/web/proyectos/7/backlog/7", ""},            // HU-06
		{http.MethodPost, "/web/proyectos/1/backlog/7/editar", ""},    // HU-06 (proyecto 1: el editor falso responde 404 para otros)
		{http.MethodPost, "/web/proyectos/7/backlog/7/eliminar", ""},  // HU-06
		{http.MethodPost, "/web/proyectos/1/sprints/1/completar", ""}, // HU-10 (sin elemento_id: 400 en el detalle del sprint)
	}
	for _, r := range rutas {
		t.Run(r.metodo+" "+r.ruta, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(r.metodo, r.ruta, strings.NewReader(r.cuerpo)))
			if rec.Code == http.StatusNotFound || rec.Code == http.StatusMethodNotAllowed {
				t.Errorf("status = %d: la ruta no está registrada en NuevoRouter", rec.Code)
			}
		})
	}
}
