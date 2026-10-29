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

type cerrarSprintFalso struct {
	err      error
	recibido app.DatosCerrarSprint
	llamado  bool
}

func (f *cerrarSprintFalso) Ejecutar(_ context.Context, d app.DatosCerrarSprint) error {
	f.llamado = true
	f.recibido = d
	return f.err
}

func postCerrarSprint(t *testing.T, uc cerradorSprints, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarCerrarSprint(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, nil))
	return rec
}

// Los cuerpos de error se leen con mensajeDeError (esfuerzo_handler_test.go).
const rutaCerrarSprint = "/proyectos/7/sprints/3/cerrar"

// Traza: HU-11 CA-11.4: pedido válido, 204 sin cuerpo.
func TestPostCerrarSprint_PedidoValido_Devuelve204SinCuerpo(t *testing.T) {
	uc := &cerrarSprintFalso{}

	rec := postCerrarSprint(t, uc, rutaCerrarSprint)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("cuerpo = %q, se esperaba vacío", rec.Body)
	}
	if uc.recibido != (app.DatosCerrarSprint{ProyectoID: 7, SprintID: 3}) {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-11 CA-11.4: proyecto o sprint inexistente, 404.
func TestPostCerrarSprint_NoEncontrado_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
	}{
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado},
		{"sprint inexistente", domain.ErrSprintNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postCerrarSprint(t, &cerrarSprintFalso{err: tt.err}, rutaCerrarSprint)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.err.Error() {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.err.Error())
			}
		})
	}
}

// Traza: HU-11 CA-11.4: sprint ya cerrado, 400.
func TestPostCerrarSprint_YaCerrado_Devuelve400(t *testing.T) {
	rec := postCerrarSprint(t, &cerrarSprintFalso{err: domain.ErrSprintYaCerrado}, rutaCerrarSprint)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := mensajeDeError(t, rec); got != domain.ErrSprintYaCerrado.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", got, domain.ErrSprintYaCerrado.Error())
	}
}

// Traza: HU-11, sección 5 (P-8): si dos pedidos cierran a la vez, el caso de uso devuelve
// ErrSprintYaCerrado envuelto por el guardado. El 400 muestra el mensaje limpio del dominio.
func TestPostCerrarSprint_YaCerradoEnvuelto_Devuelve400ConMensajeLimpio(t *testing.T) {
	envuelto := fmt.Errorf("cerrar sprint: %w", domain.ErrSprintYaCerrado)

	rec := postCerrarSprint(t, &cerrarSprintFalso{err: envuelto}, rutaCerrarSprint)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := mensajeDeError(t, rec); got != domain.ErrSprintYaCerrado.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", got, domain.ErrSprintYaCerrado.Error())
	}
}

// IDs inválidos en la ruta: 400, como en los demás handlers, sin llamar al caso de uso.
func TestPostCerrarSprint_IDInvalido_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"id de proyecto no numérico", "/proyectos/abc/sprints/3/cerrar", "el id del proyecto debe ser un número entero positivo"},
		{"id de proyecto en cero", "/proyectos/0/sprints/3/cerrar", "el id del proyecto debe ser un número entero positivo"},
		{"id de sprint no numérico", "/proyectos/7/sprints/abc/cerrar", "el id del sprint debe ser un número entero positivo"},
		{"id de sprint negativo", "/proyectos/7/sprints/-1/cerrar", "el id del sprint debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &cerrarSprintFalso{}

			rec := postCerrarSprint(t, uc, tt.ruta)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
			if uc.llamado {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}

// Cualquier otro error es una falla interna: 500, sin mostrar el detalle.
func TestPostCerrarSprint_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	rec := postCerrarSprint(t, &cerrarSprintFalso{err: errors.New("cerrar sprint: database is locked")}, rutaCerrarSprint)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := mensajeDeError(t, rec); strings.Contains(got, "database") {
		t.Errorf("mensaje = %q: no tiene que mostrar el detalle interno", got)
	}
}
