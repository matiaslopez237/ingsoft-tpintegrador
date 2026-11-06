package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorCalidadSprintsFalso devuelve lo que se le cargue, o err si está seteado.
type consultorCalidadSprintsFalso struct {
	calidad  []app.CalidadSprint
	err      error
	recibido int64
}

func (f *consultorCalidadSprintsFalso) Ejecutar(_ context.Context, proyectoID int64) ([]app.CalidadSprint, error) {
	f.recibido = proyectoID
	if f.err != nil {
		return nil, f.err
	}
	return f.calidad, nil
}

func pedirCalidad(t *testing.T, uc consultorCalidadSprints, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/proyectos/"+id+"/calidad", nil)
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarCalidadSprints(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-24 CA-24.8, sección 3: 200 con el porcentaje (o null) y los defectos de cada sprint.
func TestGetCalidad_ConSprints_DevuelveLaCalidad(t *testing.T) {
	uc := &consultorCalidadSprintsFalso{calidad: []app.CalidadSprint{
		{
			Sprint: domain.Sprint{ID: 5, Numero: 2}, PorcentajeCompletado: 75, HayPorcentaje: true,
			DefectosDetectados: 2, DefectosResueltos: 1,
		},
		{
			Sprint: domain.Sprint{ID: 4, Numero: 1}, HayPorcentaje: false,
			DefectosDetectados: 0, DefectosResueltos: 0,
		},
	}}

	rec := pedirCalidad(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 3 {
		t.Errorf("proyecto recibido = %d, se esperaba 3", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `[
		{"sprint_id":5,"numero":2,"porcentaje_completado":75,"defectos_detectados":2,"defectos_resueltos":1},
		{"sprint_id":4,"numero":1,"porcentaje_completado":null,"defectos_detectados":0,"defectos_resueltos":0}
	]`)
}

// Traza: HU-24 CA-24.8: sin sprints cerrados, un array vacío, no null.
func TestGetCalidad_SinSprints_DevuelveArrayVacio(t *testing.T) {
	uc := &consultorCalidadSprintsFalso{calidad: []app.CalidadSprint{}}

	rec := pedirCalidad(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `[]`)
}

// Traza: HU-24 CA-24.8, sección 7. El 404 muestra el mensaje del dominio aunque el caso de uso lo
// devuelva envuelto, y el 500 no muestra el detalle de la falla.
func TestGetCalidad_Errores(t *testing.T) {
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
			uc := &consultorCalidadSprintsFalso{err: tt.err}

			rec := pedirCalidad(t, uc, tt.id)

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
