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

// agregadorFalso deja probar el handler sin base.
type agregadorFalso struct {
	err      error
	recibido app.DatosIntegrante
}

func (f *agregadorFalso) Ejecutar(_ context.Context, d app.DatosIntegrante) error {
	f.recibido = d
	return f.err
}

// quitadorFalso deja probar el handler sin base.
type quitadorFalso struct {
	err      error
	recibido app.DatosQuitarIntegrante
}

func (f *quitadorFalso) Ejecutar(_ context.Context, d app.DatosQuitarIntegrante) error {
	f.recibido = d
	return f.err
}

// agregadorRolFalso deja probar el handler sin base.
type agregadorRolFalso struct {
	err      error
	recibido app.DatosRol
}

func (f *agregadorRolFalso) Ejecutar(_ context.Context, d app.DatosRol) error {
	f.recibido = d
	return f.err
}

func postRol(t *testing.T, uc agregadorRoles, id, nombre, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/proyectos/"+id+"/integrantes/"+nombre+"/roles", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	NuevoRouter(&creadorFalso{}, &agregadorFalso{}, &quitadorFalso{}, uc, &modificadorFalso{}).ServeHTTP(rec, req)
	return rec
}

func postIntegrante(t *testing.T, uc agregadorIntegrantes, id, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/proyectos/"+id+"/integrantes", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	NuevoRouter(&creadorFalso{}, uc, &quitadorFalso{}, &agregadorRolFalso{}, &modificadorFalso{}).ServeHTTP(rec, req)
	return rec
}

func deleteIntegrante(t *testing.T, uc quitadorIntegrantes, id, nombre string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodDelete, "/proyectos/"+id+"/integrantes/"+nombre, nil)
	rec := httptest.NewRecorder()
	NuevoRouter(&creadorFalso{}, &agregadorFalso{}, uc, &agregadorRolFalso{}, &modificadorFalso{}).ServeHTTP(rec, req)
	return rec
}

// Traza: HU-03 CA-03.1
func TestPostIntegrantes_DatosValidos_Devuelve204(t *testing.T) {
	uc := &agregadorFalso{}

	rec := postIntegrante(t, uc, "1", `{"nombre":"Ana","rol":"Desarrolladora"}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido.ProyectoID != 1 || uc.recibido.Nombre != "Ana" || uc.recibido.Rol != "Desarrolladora" {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-03 CA-03.2, CA-03.5, CA-03.6
func TestPostIntegrantes_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		id      string
		cuerpo  string
		uc      *agregadorFalso
		wantMsg string
	}{
		{"nombre vacío", "1", `{"nombre":"","rol":"Desarrolladora"}`, &agregadorFalso{err: domain.ErrNombreIntegranteVacio}, domain.ErrNombreIntegranteVacio.Error()},
		{"rol vacío", "1", `{"nombre":"Ana","rol":""}`, &agregadorFalso{err: domain.ErrRolIntegranteVacio}, domain.ErrRolIntegranteVacio.Error()},
		{"nombre duplicado", "1", `{"nombre":"Ana","rol":"Tester"}`, &agregadorFalso{err: domain.ErrIntegranteDuplicado}, domain.ErrIntegranteDuplicado.Error()},
		{"JSON inválido", "1", `{nombre`, &agregadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", "1", `{"nombre":"Ana","rol":"Tester","extra":true}`, &agregadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"id no numérico", "abc", `{"nombre":"Ana","rol":"Tester"}`, &agregadorFalso{}, "el id del proyecto debe ser un número entero positivo"},
		{"id negativo", "-1", `{"nombre":"Ana","rol":"Tester"}`, &agregadorFalso{}, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postIntegrante(t, tt.uc, tt.id, tt.cuerpo)

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

// Traza: HU-03, sección 7: proyecto inexistente
func TestPostIntegrantes_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := postIntegrante(t, &agregadorFalso{err: domain.ErrProyectoNoEncontrado}, "99", `{"nombre":"Ana","rol":"Tester"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}

func TestPostIntegrantes_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postIntegrante(t, &agregadorFalso{err: errors.New("disco lleno")}, "1", `{"nombre":"Ana","rol":"Tester"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}

// Traza: HU-03 CA-03.3
func TestDeleteIntegrantes_Existente_Devuelve204(t *testing.T) {
	uc := &quitadorFalso{}

	rec := deleteIntegrante(t, uc, "1", "Ana")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido.ProyectoID != 1 || uc.recibido.Nombre != "Ana" {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-03 CA-03.4, sección 7: proyecto inexistente
func TestDeleteIntegrantes_Inexistente_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre string
		uc     *quitadorFalso
	}{
		{"integrante no pertenece al proyecto", &quitadorFalso{err: domain.ErrIntegranteNoEncontrado}},
		{"proyecto inexistente", &quitadorFalso{err: domain.ErrProyectoNoEncontrado}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := deleteIntegrante(t, tt.uc, "1", "Caro")
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestDeleteIntegrantes_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := deleteIntegrante(t, &quitadorFalso{err: errors.New("disco lleno")}, "1", "Ana")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}

// Traza: HU-03 CA-03.7
func TestPostRoles_DatosValidos_Devuelve204(t *testing.T) {
	uc := &agregadorRolFalso{}

	rec := postRol(t, uc, "1", "Ana", `{"rol":"Tester"}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != (app.DatosRol{ProyectoID: 1, Nombre: "Ana", Rol: "Tester"}) {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-03 CA-03.8, CA-03.10
func TestPostRoles_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		id      string
		cuerpo  string
		uc      *agregadorRolFalso
		wantMsg string
	}{
		{"rol vacío", "1", `{"rol":""}`, &agregadorRolFalso{err: domain.ErrRolIntegranteVacio}, domain.ErrRolIntegranteVacio.Error()},
		{"rol duplicado", "1", `{"rol":"desarrolladora"}`, &agregadorRolFalso{err: domain.ErrRolDuplicado}, domain.ErrRolDuplicado.Error()},
		{"JSON inválido", "1", `{rol`, &agregadorRolFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", "1", `{"rol":"Tester","extra":true}`, &agregadorRolFalso{}, "el cuerpo debe ser un JSON válido"},
		{"id no numérico", "abc", `{"rol":"Tester"}`, &agregadorRolFalso{}, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postRol(t, tt.uc, tt.id, "Ana", tt.cuerpo)

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

// Traza: HU-03 CA-03.9, sección 7: proyecto inexistente
func TestPostRoles_Inexistente_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre string
		uc     *agregadorRolFalso
	}{
		{"integrante no pertenece al proyecto", &agregadorRolFalso{err: domain.ErrIntegranteNoEncontrado}},
		{"proyecto inexistente", &agregadorRolFalso{err: domain.ErrProyectoNoEncontrado}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postRol(t, tt.uc, "1", "Caro", `{"rol":"Tester"}`)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestPostRoles_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postRol(t, &agregadorRolFalso{err: errors.New("disco lleno")}, "1", "Ana", `{"rol":"Tester"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
