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

// consultorVelocidadFalso devuelve lo que se le cargue, o err si está seteado.
type consultorVelocidadFalso struct {
	velocidad app.VelocidadProyecto
	err       error
	recibido  int64
}

func (f *consultorVelocidadFalso) Ejecutar(_ context.Context, proyectoID int64) (app.VelocidadProyecto, error) {
	f.recibido = proyectoID
	if f.err != nil {
		return app.VelocidadProyecto{}, f.err
	}
	return f.velocidad, nil
}

func pedirVelocidad(t *testing.T, uc consultorVelocidad, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/proyectos/"+id+"/velocidad", nil)
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarVelocidad(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-22 CA-22.7, sección 3: 200 con la velocidad y la cantidad de sprints cerrados.
func TestGetVelocidad_ConDatos_DevuelveLaVelocidad(t *testing.T) {
	uc := &consultorVelocidadFalso{velocidad: app.VelocidadProyecto{Velocidad: 4.5, HayDatos: true, SprintsCerrados: 2}}

	rec := pedirVelocidad(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 3 {
		t.Errorf("proyecto recibido = %d, se esperaba 3", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `{"velocidad":4.5,"sprints_cerrados":2}`)
}

// Traza: HU-22 CA-22.7, RN-22.3: sin sprints cerrados, velocidad es null, no 0.
func TestGetVelocidad_SinDatos_DevuelveVelocidadNull(t *testing.T) {
	uc := &consultorVelocidadFalso{velocidad: app.VelocidadProyecto{HayDatos: false}}

	rec := pedirVelocidad(t, uc, "3")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `{"velocidad":null,"sprints_cerrados":0}`)
}

// Traza: HU-22 CA-22.7, sección 7. El 404 muestra el mensaje del dominio aunque el caso de uso lo
// devuelva envuelto, y el 500 no muestra el detalle de la falla.
func TestGetVelocidad_Errores(t *testing.T) {
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
			uc := &consultorVelocidadFalso{err: tt.err}

			rec := pedirVelocidad(t, uc, tt.id)

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
