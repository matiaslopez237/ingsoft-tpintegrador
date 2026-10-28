package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type asignarHistoriaFalso struct {
	err      error
	recibido app.DatosAsignacion
}

func (f *asignarHistoriaFalso) Ejecutar(_ context.Context, d app.DatosAsignacion) error {
	f.recibido = d
	return f.err
}

type desasignarHistoriaFalso struct {
	err      error
	recibido app.DatosDesasignacion
}

func (f *desasignarHistoriaFalso) Ejecutar(_ context.Context, d app.DatosDesasignacion) error {
	f.recibido = d
	return f.err
}

func postAsignacion(t *testing.T, uc agregadorAsignacion, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarAsignacion(mux, uc, &desasignarHistoriaFalso{})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return rec
}

func deleteAsignacion(t *testing.T, uc quitadorAsignacion, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarAsignacion(mux, &asignarHistoriaFalso{}, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, ruta, nil))
	return rec
}

const rutaAsignacion = "/proyectos/7/backlog/3/sprint"

// --- POST (asignar) ---

// Traza: HU-09 CA-09.1
func TestPostAsignacion_DatosValidos_Devuelve204(t *testing.T) {
	uc := &asignarHistoriaFalso{}

	rec := postAsignacion(t, uc, rutaAsignacion, `{"sprint_id":10}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido.ProyectoID != 7 || uc.recibido.ElementoID != 3 || uc.recibido.SprintID != 10 {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-09, sección 7: elemento o sprint inexistentes
func TestPostAsignacion_NoEncontrado_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
	}{
		{"elemento inexistente", domain.ErrElementoNoEncontrado},
		{"sprint inexistente", domain.ErrSprintNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postAsignacion(t, &asignarHistoriaFalso{err: tt.err}, rutaAsignacion, `{"sprint_id":10}`)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
		})
	}
}

// Traza: HU-09 CA-09.2, CA-09.7, CA-09.8
func TestPostAsignacion_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		cuerpo  string
		uc      *asignarHistoriaFalso
		wantMsg string
	}{
		{"sprint de otro proyecto", rutaAsignacion, `{"sprint_id":10}`, &asignarHistoriaFalso{err: domain.ErrSprintDeOtroProyecto}, domain.ErrSprintDeOtroProyecto.Error()},
		{"elemento ya asignado", rutaAsignacion, `{"sprint_id":10}`, &asignarHistoriaFalso{err: domain.ErrElementoYaAsignado}, domain.ErrElementoYaAsignado.Error()},
		{"sprint cerrado", rutaAsignacion, `{"sprint_id":10}`, &asignarHistoriaFalso{err: domain.ErrSprintCerrado}, domain.ErrSprintCerrado.Error()},
		{"JSON inválido", rutaAsignacion, `{sprint`, &asignarHistoriaFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", rutaAsignacion, `{"sprint_id":10,"extra":true}`, &asignarHistoriaFalso{}, "el cuerpo debe ser un JSON válido"},
		{"id de proyecto no numérico", "/proyectos/abc/backlog/3/sprint", `{"sprint_id":10}`, &asignarHistoriaFalso{}, "el id del proyecto debe ser un número entero positivo"},
		{"id de elemento no numérico", "/proyectos/7/backlog/abc/sprint", `{"sprint_id":10}`, &asignarHistoriaFalso{}, "el id del elemento debe ser un número entero positivo"},
		{"sprint_id en cero", rutaAsignacion, `{"sprint_id":0}`, &asignarHistoriaFalso{}, "el sprint_id debe ser un número entero positivo"},
		{"sprint_id negativo", rutaAsignacion, `{"sprint_id":-1}`, &asignarHistoriaFalso{}, "el sprint_id debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postAsignacion(t, tt.uc, tt.ruta, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			var got respuestaError
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decodificar respuesta: %v", err)
			}
			if got.Error != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got.Error, tt.wantMsg)
			}
		})
	}
}

func TestPostAsignacion_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postAsignacion(t, &asignarHistoriaFalso{err: errors.New("disco lleno")}, rutaAsignacion, `{"sprint_id":10}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}

// --- DELETE (desasignar) ---

// Traza: HU-09 CA-09.4
func TestDeleteAsignacion_DatosValidos_Devuelve204(t *testing.T) {
	uc := &desasignarHistoriaFalso{}

	rec := deleteAsignacion(t, uc, rutaAsignacion)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido.ProyectoID != 7 || uc.recibido.ElementoID != 3 {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-09 CA-09.5, CA-09.6
func TestDeleteAsignacion_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		uc      *desasignarHistoriaFalso
		wantMsg string
	}{
		{"sin sprint asignado", &desasignarHistoriaFalso{err: domain.ErrElementoSinSprint}, domain.ErrElementoSinSprint.Error()},
		{"sprint cerrado", &desasignarHistoriaFalso{err: domain.ErrSprintCerrado}, domain.ErrSprintCerrado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := deleteAsignacion(t, tt.uc, rutaAsignacion)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			var got respuestaError
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decodificar respuesta: %v", err)
			}
			if got.Error != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got.Error, tt.wantMsg)
			}
		})
	}
}

// Traza: HU-09, sección 7: elemento inexistente
func TestDeleteAsignacion_ElementoInexistente_Devuelve404(t *testing.T) {
	rec := deleteAsignacion(t, &desasignarHistoriaFalso{err: domain.ErrElementoNoEncontrado}, rutaAsignacion)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}

func TestDeleteAsignacion_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := deleteAsignacion(t, &desasignarHistoriaFalso{err: errors.New("disco lleno")}, rutaAsignacion)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
