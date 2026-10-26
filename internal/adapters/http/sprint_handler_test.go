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

// creadorSprintsFalso deja probar el handler sin base: devuelve el sprint número 1 con los datos recibidos.
type creadorSprintsFalso struct {
	err      error
	recibido app.DatosNuevoSprint
}

func (f *creadorSprintsFalso) Ejecutar(_ context.Context, d app.DatosNuevoSprint) (domain.Sprint, error) {
	f.recibido = d
	if f.err != nil {
		return domain.Sprint{}, f.err
	}
	return domain.Sprint{ID: 5, ProyectoID: d.ProyectoID, Numero: 1, FechaInicio: d.FechaInicio, FechaFin: d.FechaFin,
		Objetivo: strings.TrimSpace(d.Objetivo), Estado: domain.SprintAbierto}, nil
}

// listadorSprintsFalso devuelve una lista fija de sprints, o un error.
type listadorSprintsFalso struct {
	sprints  []domain.Sprint
	err      error
	recibido int64
}

func (f *listadorSprintsFalso) Ejecutar(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
	f.recibido = proyectoID
	return f.sprints, f.err
}

func pedirSprints(t *testing.T, crear creadorSprints, listar listadorSprints, metodo, id, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, "/proyectos/"+id+"/sprints", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarSprints(mux, crear, listar)
	mux.ServeHTTP(rec, req)
	return rec
}

// Traza: HU-08 CA-08.1, sección 3
func TestPostSprints_DatosValidos_Devuelve201ConElSprint(t *testing.T) {
	uc := &creadorSprintsFalso{}

	rec := pedirSprints(t, uc, &listadorSprintsFalso{}, http.MethodPost, "1",
		`{"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"Gestionar el backlog"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	d := uc.recibido
	if d.ProyectoID != 1 || d.FechaInicio.Format(time.DateOnly) != "2026-10-26" ||
		d.FechaFin.Format(time.DateOnly) != "2026-10-27" || d.Objetivo != "Gestionar el backlog" {
		t.Errorf("datos recibidos por el caso de uso: %+v", d)
	}
	mismoJSON(t, rec.Body.String(), `{"id":5,"proyecto_id":1,"numero":1,"fecha_inicio":"2026-10-26",
		"fecha_fin":"2026-10-27","objetivo":"Gestionar el backlog","estado":"Abierto"}`)
}

// Traza: HU-08 CA-08.5: una fecha que no se envía llega al caso de uso como fecha cero (sin fecha).
func TestPostSprints_SinFechas_LlegaComoFechaCero(t *testing.T) {
	uc := &creadorSprintsFalso{err: domain.ErrFechaInicioObligatoria}

	pedirSprints(t, uc, &listadorSprintsFalso{}, http.MethodPost, "1", `{"objetivo":"x"}`)

	if !uc.recibido.FechaInicio.IsZero() || !uc.recibido.FechaFin.IsZero() {
		t.Errorf("fechas recibidas = %v y %v, se esperaban cero", uc.recibido.FechaInicio, uc.recibido.FechaFin)
	}
}

// Traza: HU-08 CA-08.4 a CA-08.8, sección 7
func TestPostSprints_DatosInvalidos_Devuelve400(t *testing.T) {
	valido := `{"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"x"}`
	tests := []struct {
		nombre  string
		id      string
		cuerpo  string
		err     error
		wantMsg string
	}{
		{"proyecto no activo", "1", valido, domain.ErrProyectoNoActivo, domain.ErrProyectoNoActivo.Error()},
		{"sin fecha de inicio", "1", valido, domain.ErrFechaInicioObligatoria, domain.ErrFechaInicioObligatoria.Error()},
		{"sin fecha de fin", "1", valido, domain.ErrFechaFinObligatoria, domain.ErrFechaFinObligatoria.Error()},
		{"fin anterior", "1", valido, domain.ErrFechaFinAnterior, domain.ErrFechaFinAnterior.Error()},
		{"Sprint Goal vacío", "1", valido, domain.ErrObjetivoVacio, domain.ErrObjetivoVacio.Error()},
		{"Sprint Goal largo", "1", valido, domain.ErrObjetivoDemasiadoLargo, domain.ErrObjetivoDemasiadoLargo.Error()},
		{"fuera del proyecto", "1", valido, domain.ErrSprintFueraDelProyecto, domain.ErrSprintFueraDelProyecto.Error()},
		{"superpuesto", "1", valido, domain.ErrSprintSuperpuesto, domain.ErrSprintSuperpuesto.Error()},
		{"fecha con formato inválido", "1", `{"fecha_inicio":"26/10/2026","fecha_fin":"2026-10-27","objetivo":"x"}`, nil, "fecha_inicio debe tener el formato AAAA-MM-DD"},
		{"JSON inválido", "1", `{objetivo`, nil, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", "1", `{"objetivo":"x","extra":true}`, nil, "el cuerpo debe ser un JSON válido"},
		{"id no numérico", "abc", valido, nil, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirSprints(t, &creadorSprintsFalso{err: tt.err}, &listadorSprintsFalso{}, http.MethodPost, tt.id, tt.cuerpo)

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

// Traza: HU-08 CA-08.3
func TestPostSprints_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := pedirSprints(t, &creadorSprintsFalso{err: domain.ErrProyectoNoEncontrado}, &listadorSprintsFalso{}, http.MethodPost, "99",
		`{"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"x"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}

func TestPostSprints_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := pedirSprints(t, &creadorSprintsFalso{err: errors.New("disco lleno")}, &listadorSprintsFalso{}, http.MethodPost, "1",
		`{"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"x"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}

// Traza: HU-08 CA-08.9
func TestGetSprints_DevuelveLosSprintsDelProyecto(t *testing.T) {
	uc := &listadorSprintsFalso{sprints: []domain.Sprint{
		{ID: 1, ProyectoID: 3, Numero: 1, FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
			FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC), Objetivo: "A", Estado: domain.SprintAbierto},
		{ID: 2, ProyectoID: 3, Numero: 2, FechaInicio: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC),
			FechaFin: time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC), Objetivo: "B", Estado: domain.SprintAbierto},
	}}

	rec := pedirSprints(t, &creadorSprintsFalso{}, uc, http.MethodGet, "3", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 3 {
		t.Errorf("proyecto recibido = %d, se esperaba 3", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `[
		{"id":1,"proyecto_id":3,"numero":1,"fecha_inicio":"2026-10-26","fecha_fin":"2026-10-27","objetivo":"A","estado":"Abierto"},
		{"id":2,"proyecto_id":3,"numero":2,"fecha_inicio":"2026-10-28","fecha_fin":"2026-10-29","objetivo":"B","estado":"Abierto"}]`)
}

// Traza: HU-08 CA-08.9: un proyecto sin sprints devuelve [] y no null.
func TestGetSprints_SinSprints_DevuelveArrayVacio(t *testing.T) {
	rec := pedirSprints(t, &creadorSprintsFalso{}, &listadorSprintsFalso{}, http.MethodGet, "3", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `[]`)
}

func TestGetSprints_Errores(t *testing.T) {
	tests := []struct {
		nombre string
		id     string
		err    error
		want   int
	}{
		{"proyecto inexistente", "99", domain.ErrProyectoNoEncontrado, http.StatusNotFound},
		{"id inválido", "abc", nil, http.StatusBadRequest},
		{"falla interna", "1", errors.New("disco lleno"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirSprints(t, &creadorSprintsFalso{}, &listadorSprintsFalso{err: tt.err}, http.MethodGet, tt.id, "")
			if rec.Code != tt.want {
				t.Errorf("status = %d, se esperaba %d", rec.Code, tt.want)
			}
		})
	}
}
