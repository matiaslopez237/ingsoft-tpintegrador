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

// creadorFalso deja probar el handler sin base: usa el dominio real y un ID fijo.
type creadorFalso struct {
	err      error
	recibido app.DatosNuevoProyecto
}

func (c *creadorFalso) Ejecutar(_ context.Context, d app.DatosNuevoProyecto) (domain.Proyecto, error) {
	c.recibido = d
	if c.err != nil {
		return domain.Proyecto{}, c.err
	}
	p, err := domain.NuevoProyecto(d.Nombre, d.FechaInicio, d.FechaFin)
	if err != nil {
		return domain.Proyecto{}, err
	}
	p.ID = 1
	return p, nil
}

func postProyectos(t *testing.T, uc creadorProyectos, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/proyectos", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarCrearProyecto(mux, uc)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-01 CA-01.1
func TestPostProyectos_DatosValidos_Devuelve201ConProyecto(t *testing.T) {
	rec := postProyectos(t, &creadorFalso{},
		`{"nombre":"TP Integrador","fecha_inicio":"2026-03-01","fecha_fin":"2026-11-30"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	var got respuestaProyecto
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	if got.ID != 1 || got.Nombre != "TP Integrador" || got.Estado != "Activo" ||
		got.FechaInicio != "2026-03-01" || got.FechaFin == nil || *got.FechaFin != "2026-11-30" {
		t.Errorf("respuesta inesperada: %+v", got)
	}
}

// Traza: HU-01 CA-01.7
func TestPostProyectos_SinFechaFin_DevuelveFechaFinNula(t *testing.T) {
	rec := postProyectos(t, &creadorFalso{}, `{"nombre":"Alfa","fecha_inicio":"2026-03-01"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), `"fecha_fin":null`) {
		t.Errorf("se esperaba fecha_fin null, cuerpo: %s", rec.Body)
	}
}

// Traza: HU-01 CA-01.2, CA-01.4, CA-01.5, CA-01.6
func TestPostProyectos_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		cuerpo  string
		creador *creadorFalso
		wantMsg string
	}{
		{"nombre vacío", `{"nombre":"","fecha_inicio":"2026-03-01"}`, &creadorFalso{}, domain.ErrNombreVacio.Error()},
		{"sin fecha de inicio", `{"nombre":"Alfa"}`, &creadorFalso{}, domain.ErrFechaInicioObligatoria.Error()},
		{"fin anterior", `{"nombre":"Alfa","fecha_inicio":"2026-03-02","fecha_fin":"2026-03-01"}`, &creadorFalso{}, domain.ErrFechaFinAnterior.Error()},
		{"nombre duplicado", `{"nombre":"Alfa","fecha_inicio":"2026-03-01"}`, &creadorFalso{err: domain.ErrNombreDuplicado}, domain.ErrNombreDuplicado.Error()},
		{"fecha con formato inválido", `{"nombre":"Alfa","fecha_inicio":"01/03/2026"}`, &creadorFalso{}, "fecha_inicio debe tener el formato AAAA-MM-DD"},
		{"fecha de fin con formato inválido", `{"nombre":"Alfa","fecha_inicio":"2026-03-01","fecha_fin":"x"}`, &creadorFalso{}, "fecha_fin debe tener el formato AAAA-MM-DD"},
		{"JSON inválido", `{nombre`, &creadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", `{"nombre":"Alfa","fecha_inicio":"2026-03-01","estado":"Cancelado"}`, &creadorFalso{}, "el cuerpo debe ser un JSON válido"},
		{"cuerpo demasiado grande", `{"nombre":"` + strings.Repeat("a", tamanoMaximoCuerpo) + `"}`, &creadorFalso{}, "el cuerpo debe ser un JSON válido"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postProyectos(t, tt.creador, tt.cuerpo)

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

func TestPostProyectos_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postProyectos(t, &creadorFalso{err: errors.New("disco lleno")},
		`{"nombre":"Alfa","fecha_inicio":"2026-03-01"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
