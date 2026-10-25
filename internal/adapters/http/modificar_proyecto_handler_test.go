package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// modificadorFalso deja probar el handler sin base: aplica los cambios con el dominio real
// sobre un proyecto fijo con ID 1.
type modificadorFalso struct {
	err      error
	recibido app.DatosModificarProyecto
}

func (m *modificadorFalso) Ejecutar(_ context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error) {
	m.recibido = d
	if m.err != nil {
		return domain.Proyecto{}, m.err
	}
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	base := domain.Proyecto{ID: 1, Nombre: "Alfa", FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), FechaFin: &fin, Estado: domain.EstadoActivo}
	return base.Modificar(d.Cambios)
}

func patchProyecto(t *testing.T, uc modificadorProyectos, id, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPatch, "/proyectos/"+id, strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarModificarProyecto(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-02 CA-02.1, RN-02.1: solo se envía el nombre; el resto queda como "no enviado".
func TestPatchProyecto_SoloNombre_Devuelve200ConElProyecto(t *testing.T) {
	uc := &modificadorFalso{}

	rec := patchProyecto(t, uc, "1", `{"nombre":"Alfa Renovado"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	c := uc.recibido.Cambios
	if uc.recibido.ProyectoID != 1 || c.Nombre == nil || *c.Nombre != "Alfa Renovado" ||
		c.FechaInicio != nil || c.CambiarFechaFin || c.Estado != nil {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
	var got respuestaProyecto
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	if got.ID != 1 || got.Nombre != "Alfa Renovado" || got.FechaInicio != "2026-03-01" ||
		got.FechaFin == nil || *got.FechaFin != "2026-11-30" || got.Estado != "Activo" {
		t.Errorf("respuesta inesperada: %+v", got)
	}
}

// Traza: HU-02 sección 2 (distinguir "no enviado" de "enviado"), CA-02.5, CA-02.6, CA-02.11
func TestPatchProyecto_TraduceLosCampos(t *testing.T) {
	tests := []struct {
		nombre    string
		cuerpo    string
		verificar func(t *testing.T, c domain.CambiosProyecto)
	}{
		{"fecha_fin en null la quita", `{"fecha_fin":null}`, func(t *testing.T, c domain.CambiosProyecto) {
			if !c.CambiarFechaFin || c.FechaFin != nil {
				t.Errorf("CambiarFechaFin = %v, FechaFin = %v; se esperaba quitarla", c.CambiarFechaFin, c.FechaFin)
			}
		}},
		{"fecha_fin con fecha la reemplaza", `{"fecha_fin":"2026-12-31"}`, func(t *testing.T, c domain.CambiosProyecto) {
			if !c.CambiarFechaFin || c.FechaFin == nil || c.FechaFin.Format(time.DateOnly) != "2026-12-31" {
				t.Errorf("CambiarFechaFin = %v, FechaFin = %v", c.CambiarFechaFin, c.FechaFin)
			}
		}},
		{"fecha_inicio con fecha la reemplaza", `{"fecha_inicio":"2026-04-01"}`, func(t *testing.T, c domain.CambiosProyecto) {
			if c.FechaInicio == nil || c.FechaInicio.Format(time.DateOnly) != "2026-04-01" {
				t.Errorf("FechaInicio = %v", c.FechaInicio)
			}
		}},
		{"estado", `{"estado":"Finalizado"}`, func(t *testing.T, c domain.CambiosProyecto) {
			if c.Estado == nil || *c.Estado != domain.EstadoFinalizado {
				t.Errorf("Estado = %v", c.Estado)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &modificadorFalso{}

			rec := patchProyecto(t, uc, "1", tt.cuerpo)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
			}
			tt.verificar(t, uc.recibido.Cambios)
		})
	}
}

// Traza: HU-02 CA-02.2, CA-02.3, CA-02.4, CA-02.7, CA-02.8, CA-02.11
func TestPatchProyecto_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		id      string
		cuerpo  string
		uc      *modificadorFalso
		wantMsg string
	}{
		{"sin cambios", "1", `{}`, &modificadorFalso{}, domain.ErrSinCambios.Error()},
		{"nombre vacío", "1", `{"nombre":""}`, &modificadorFalso{}, domain.ErrNombreVacio.Error()},
		{"nombre duplicado", "1", `{"nombre":"Beta"}`, &modificadorFalso{err: domain.ErrNombreDuplicado}, domain.ErrNombreDuplicado.Error()},
		{"quitar fecha de inicio", "1", `{"fecha_inicio":null}`, &modificadorFalso{}, domain.ErrFechaInicioObligatoria.Error()},
		{"fecha de fin anterior", "1", `{"fecha_fin":"2026-01-01"}`, &modificadorFalso{}, domain.ErrFechaFinAnterior.Error()},
		{"estado inválido", "1", `{"estado":"Cerrado"}`, &modificadorFalso{}, domain.ErrEstadoInvalido.Error()},
		{"fecha con formato inválido", "1", `{"fecha_fin":"30/11/2026"}`, &modificadorFalso{}, "fecha_fin debe tener el formato AAAA-MM-DD"},
		{"JSON inválido", "1", `{nombre`, &modificadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", "1", `{"nombre":"Alfa","extra":true}`, &modificadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"id no numérico", "abc", `{"nombre":"Alfa"}`, &modificadorFalso{}, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := patchProyecto(t, tt.uc, tt.id, tt.cuerpo)

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

// Traza: HU-02 CA-02.9
func TestPatchProyecto_Inexistente_Devuelve404(t *testing.T) {
	rec := patchProyecto(t, &modificadorFalso{err: domain.ErrProyectoNoEncontrado}, "99", `{"nombre":"Gamma"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}

func TestPatchProyecto_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := patchProyecto(t, &modificadorFalso{err: errors.New("disco lleno")}, "1", `{"nombre":"Gamma"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
