package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Cada historia registra sus rutas con su propia función RegistrarX; este test verifica que
// NuevoRouter (lo que usa main.go) las registra todas. Si una historia agrega una ruta y se
// olvida de sumarla a NuevoRouter, el pedido da 404 o 405 y el test falla.
func TestNuevoRouter_RegistraTodasLasRutas(t *testing.T) {
	mux := NuevoRouter(CasosDeUso{
		CrearProyecto:     &creadorFalso{},
		ModificarProyecto: &modificadorFalso{},
		AgregarIntegrante: &agregadorFalso{},
		QuitarIntegrante:  &quitadorFalso{},
		AgregarRol:        &agregadorRolFalso{},
		ConsultarProyecto: &consultorFalso{},
	})
	rutas := []struct{ metodo, ruta, cuerpo string }{
		{http.MethodGet, "/salud", ""},
		{http.MethodPost, "/proyectos", `{"nombre":"Alfa","fecha_inicio":"2026-03-01"}`}, // HU-01
		{http.MethodPatch, "/proyectos/1", `{"nombre":"Beta"}`},                          // HU-02
		{http.MethodPost, "/proyectos/1/integrantes", `{"nombre":"Ana","rol":"Tester"}`}, // HU-03
		{http.MethodDelete, "/proyectos/1/integrantes/Ana", ""},                          // HU-03
		{http.MethodPost, "/proyectos/1/integrantes/Ana/roles", `{"rol":"Tester"}`},      // HU-03
		{http.MethodGet, "/proyectos/1", ""},                                             // HU-04
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
