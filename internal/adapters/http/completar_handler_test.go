package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type completarHistoriaFalso struct {
	err      error
	llamado  bool
	recibido app.DatosCompletar
}

func (f *completarHistoriaFalso) Ejecutar(_ context.Context, d app.DatosCompletar) error {
	f.llamado = true
	f.recibido = d
	return f.err
}

func postCompletar(t *testing.T, uc completadorHistorias, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarCompletar(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, nil))
	return rec
}

const rutaCompletar = "/proyectos/7/backlog/3/completar"

// Traza: HU-10 CA-10.6
func TestPostCompletar_PedidoValido_Devuelve204(t *testing.T) {
	uc := &completarHistoriaFalso{}

	rec := postCompletar(t, uc, rutaCompletar)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("la respuesta 204 no debería tener cuerpo: %s", rec.Body)
	}
	if uc.recibido.ProyectoID != 7 || uc.recibido.ElementoID != 3 {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-10 CA-10.6, sección 7
func TestPostCompletar_ErroresDelCasoDeUso(t *testing.T) {
	tests := []struct {
		nombre     string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"elemento inexistente", domain.ErrElementoNoEncontrado, http.StatusNotFound, domain.ErrElementoNoEncontrado.Error()},
		{"proyecto inexistente envuelto", fmt.Errorf("buscar proyecto 7: %w", domain.ErrProyectoNoEncontrado), http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"elemento sin sprint", domain.ErrElementoSinSprint, http.StatusBadRequest, "el elemento no tiene un sprint asignado"},
		{"sprint cerrado", domain.ErrSprintCerrado, http.StatusBadRequest, "el sprint ya está cerrado"},
		{"sprint cerrado envuelto", fmt.Errorf("guardar completado: %w", domain.ErrSprintCerrado), http.StatusBadRequest, "el sprint ya está cerrado"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postCompletar(t, &completarHistoriaFalso{err: tt.err}, rutaCompletar)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %s", rec.Code, tt.wantStatus, rec.Body)
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

// Traza: HU-10, sección 2: IDs de la ruta
func TestPostCompletar_IDsInvalidos_Devuelve400SinLlamarAlCasoDeUso(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"id de proyecto no numérico", "/proyectos/abc/backlog/3/completar", "el id del proyecto debe ser un número entero positivo"},
		{"id de proyecto en cero", "/proyectos/0/backlog/3/completar", "el id del proyecto debe ser un número entero positivo"},
		{"id de proyecto negativo", "/proyectos/-1/backlog/3/completar", "el id del proyecto debe ser un número entero positivo"},
		{"id de elemento no numérico", "/proyectos/7/backlog/abc/completar", "el id del elemento debe ser un número entero positivo"},
		{"id de elemento en cero", "/proyectos/7/backlog/0/completar", "el id del elemento debe ser un número entero positivo"},
		{"id de elemento negativo", "/proyectos/7/backlog/-1/completar", "el id del elemento debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &completarHistoriaFalso{}
			rec := postCompletar(t, uc, tt.ruta)

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
			if uc.llamado {
				t.Error("con un ID inválido no se debería llamar al caso de uso")
			}
		})
	}
}

func TestPostCompletar_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postCompletar(t, &completarHistoriaFalso{err: errors.New("disco lleno")}, rutaCompletar)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
