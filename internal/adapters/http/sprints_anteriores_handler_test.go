package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// pedirSprintsAnteriores registra solo la ruta de HU-12 y hace el pedido.
func pedirSprintsAnteriores(t *testing.T, uc listadorSprints, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/proyectos/"+id+"/sprints/anteriores", nil)
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarSprintsAnteriores(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-12 CA-12.6, sección 3: 200 con los sprints en el orden del caso de uso y con sus puntos.
func TestGetSprintsAnteriores_DevuelveLosSprintsConSusPuntos(t *testing.T) {
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.UTC) }
	uc := &listadorSprintsFalso{sprints: []domain.Sprint{
		{ID: 2, ProyectoID: 3, Numero: 2, FechaInicio: dia(13), FechaFin: dia(25), Objetivo: "Interfaz web",
			Estado: domain.SprintCerrado, PuntosPlanificados: 8, PuntosCompletados: 5},
		{ID: 1, ProyectoID: 3, Numero: 1, FechaInicio: dia(1), FechaFin: dia(12), Objetivo: "Armar el MVP",
			Estado: domain.SprintCerrado, PuntosPlanificados: 13, PuntosCompletados: 10},
	}}

	rec := pedirSprintsAnteriores(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 3 {
		t.Errorf("proyecto recibido = %d, se esperaba 3", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `[
		{"id":2,"numero":2,"fecha_inicio":"2026-10-13","fecha_fin":"2026-10-25","objetivo":"Interfaz web","puntos_planificados":8,"puntos_completados":5},
		{"id":1,"numero":1,"fecha_inicio":"2026-10-01","fecha_fin":"2026-10-12","objetivo":"Armar el MVP","puntos_planificados":13,"puntos_completados":10}
	]`)
}

// Traza: HU-12 CA-12.6: sin sprints cerrados es un array vacío, no null.
func TestGetSprintsAnteriores_SinSprints_DevuelveArrayVacio(t *testing.T) {
	rec := pedirSprintsAnteriores(t, &listadorSprintsFalso{}, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `[]`)
}

// Traza: HU-12 CA-12.6, sección 7. El 404 muestra el mensaje del dominio aunque el caso de uso lo
// devuelva envuelto, y el 500 no muestra el detalle de la falla.
func TestGetSprintsAnteriores_Errores(t *testing.T) {
	tests := []struct {
		nombre  string
		id      string
		err     error
		want    int
		wantMsg string
	}{
		{"proyecto inexistente", "99", fmt.Errorf("buscar proyecto 99: %w", domain.ErrProyectoNoEncontrado), http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"id inválido", "abc", nil, http.StatusBadRequest, ""},
		{"falla interna", "1", errors.New("disco lleno"), http.StatusInternalServerError, "error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &listadorSprintsFalso{err: tt.err}

			rec := pedirSprintsAnteriores(t, uc, tt.id)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d", rec.Code, tt.want)
			}
			if !strings.Contains(rec.Body.String(), tt.wantMsg) {
				t.Errorf("cuerpo = %s, se esperaba que contuviera %q", rec.Body, tt.wantMsg)
			}
			if strings.Contains(rec.Body.String(), "disco lleno") || strings.Contains(rec.Body.String(), "buscar proyecto") {
				t.Errorf("cuerpo = %s: no tiene que mostrar el detalle interno", rec.Body)
			}
			if tt.id == "abc" && uc.recibido != 0 {
				t.Errorf("con un ID inválido no se tiene que llamar al caso de uso")
			}
		})
	}
}
