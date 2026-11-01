package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

func patchElemento(t *testing.T, editor *editorElementosFalso, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	return pedirBacklogEdicion(t, editor, &eliminadorElementosFalso{}, http.MethodPatch, "/proyectos/1/backlog/7", cuerpo)
}

// Ausente no cambia nada, null quita la estimación y un número la carga. La respuesta del PATCH
// devuelve el elemento como quedó (sección 3).
// Traza: HU-19 CA-19.13 (RN-19.3)
func TestPatchElemento_HorasEstimadas(t *testing.T) {
	tests := []struct {
		nombre        string
		cuerpo        string
		wantCambiar   bool
		wantHoras     *float64
		wantRespuesta string
	}{
		{"un número las carga", `{"horas_estimadas":8}`, true, pHoras(8), "8"},
		{"con media hora", `{"horas_estimadas":0.5}`, true, pHoras(0.5), "0.5"},
		{"en el tope", `{"horas_estimadas":1000}`, true, pHoras(1000), "1000"},
		{"null las quita", `{"horas_estimadas":null}`, true, nil, "null"},
		{"ausente no las toca", `{"titulo":"Otro"}`, false, nil, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			editor := &editorElementosFalso{}

			rec := patchElemento(t, editor, tt.cuerpo)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
			}
			c := editor.recibido.Cambios
			if c.CambiarHorasEstimadas != tt.wantCambiar {
				t.Errorf("CambiarHorasEstimadas = %t, se esperaba %t", c.CambiarHorasEstimadas, tt.wantCambiar)
			}
			switch {
			case tt.wantHoras == nil && c.HorasEstimadas != nil:
				t.Errorf("HorasEstimadas = %v, se esperaba nil", *c.HorasEstimadas)
			case tt.wantHoras != nil && (c.HorasEstimadas == nil || *c.HorasEstimadas != *tt.wantHoras):
				t.Errorf("HorasEstimadas = %v, se esperaba %v", c.HorasEstimadas, *tt.wantHoras)
			}
			if got := camposJSON(t, rec)["horas_estimadas"]; got != tt.wantRespuesta {
				t.Errorf("horas_estimadas en la respuesta = %s, se esperaba %s", got, tt.wantRespuesta)
			}
		})
	}
}

// Las horas no se mezclan con los story points: cada uno tiene su marca (RN-19.1).
// Traza: HU-19 CA-19.15
func TestPatchElemento_HorasEstimadasYStoryPoints_SonIndependientes(t *testing.T) {
	editor := &editorElementosFalso{}

	rec := patchElemento(t, editor, `{"story_points":13}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo: %s", rec.Code, rec.Body)
	}
	if c := editor.recibido.Cambios; c.CambiarHorasEstimadas || c.HorasEstimadas != nil || !c.CambiarStoryPoints {
		t.Errorf("cambios enviados al caso de uso: %+v", c)
	}
}

// Un valor que no es un número da 400 sin llamar al caso de uso.
// Traza: HU-19 sección 7
func TestPatchElemento_HorasEstimadasQueNoSonNumero_Devuelve400SinLlamarAlCasoDeUso(t *testing.T) {
	for _, valor := range []string{`"ocho"`, `"8"`, `true`, `[8]`, `{"horas":8}`, `1e999`} {
		t.Run(valor, func(t *testing.T) {
			editor := &editorElementosFalso{}

			rec := patchElemento(t, editor, `{"horas_estimadas":`+valor+`}`)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != "horas_estimadas debe ser un número o null" {
				t.Errorf("mensaje = %q", got)
			}
			if editor.recibido.ProyectoID != 0 {
				t.Errorf("no se esperaba que se llamara al caso de uso: %+v", editor.recibido)
			}
		})
	}
}

// El falso usa el dominio real, así que estos errores salen de la validación verdadera (RN-19.2).
// Traza: HU-19 CA-19.14
func TestPatchElemento_HorasEstimadasInvalidas_Devuelve400ConElMensajeDelDominio(t *testing.T) {
	tests := []struct {
		cuerpo  string
		wantErr error
	}{
		{`{"horas_estimadas":0}`, domain.ErrHorasInvalidas},
		{`{"horas_estimadas":-1}`, domain.ErrHorasInvalidas},
		{`{"horas_estimadas":1.25}`, domain.ErrHorasNoMultiploMediaHora},
		{`{"horas_estimadas":1000.5}`, domain.ErrHorasEstimadasExcedenMaximo},
	}
	for _, tt := range tests {
		t.Run(tt.cuerpo, func(t *testing.T) {
			rec := patchElemento(t, &editorElementosFalso{}, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantErr.Error() {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantErr.Error())
			}
		})
	}
}

// Sin estos errores en erroresDeValidacionEdicion responderían 500. Si llegan envueltos, el mensaje
// sigue siendo el del dominio.
// Traza: HU-19 CA-19.14 (sección 7)
func TestPatchElemento_ErrorDeHorasEnvuelto_Devuelve400ConMensajeLimpio(t *testing.T) {
	for _, err := range []error{domain.ErrHorasInvalidas, domain.ErrHorasNoMultiploMediaHora, domain.ErrHorasEstimadasExcedenMaximo} {
		t.Run(err.Error(), func(t *testing.T) {
			envuelto := fmt.Errorf("editar elemento 7: %w", err)

			rec := patchElemento(t, &editorElementosFalso{err: envuelto}, `{"horas_estimadas":8}`)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != err.Error() {
				t.Errorf("mensaje = %q, se esperaba %q", got, err.Error())
			}
		})
	}
}

// Traza: HU-19 sección 3 (RN-19.3): la respuesta de crear suma el campo, siempre null porque un
// elemento nace sin estimar.
func TestPostBacklog_RespuestaConHorasEstimadasEnNull(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{}, "/proyectos/1/backlog",
		`{"titulo":"Login","descripcion":"Ingresar","prioridad":"Alta","criterios":["Con email"]}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	got, ok := camposJSON(t, rec)["horas_estimadas"]
	if !ok || got != "null" {
		t.Errorf("horas_estimadas = %q (presente: %t), se esperaba null", got, ok)
	}
}

// Traza: HU-19 sección 5 (fuera de alcance): las horas no se cargan al crear, sino editando.
func TestPostBacklog_ConHorasEstimadas_Devuelve400(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{}, "/proyectos/1/backlog",
		`{"titulo":"Login","descripcion":"Ingresar","prioridad":"Alta","criterios":["Con email"],"horas_estimadas":8}`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
	}
}

// Traza: HU-19 sección 3: el listado trae el campo en cada elemento, con valor o null.
func TestGetBacklog_RespuestaConHorasEstimadas(t *testing.T) {
	uc := &listadorBacklogFalso{elementos: []domain.ElementoBacklog{
		{ID: 1, ProyectoID: 1, Titulo: "Login", HorasEstimadas: pHoras(7.5)},
		{ID: 2, ProyectoID: 1, Titulo: "Reportes"},
	}}

	rec := getBacklog(t, uc, "/proyectos/1/backlog")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	var got []map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("el cuerpo no es un array JSON: %v; cuerpo: %s", err, rec.Body)
	}
	if len(got) != 2 {
		t.Fatalf("elementos = %d, se esperaban 2", len(got))
	}
	for i, want := range []string{"7.5", "null"} {
		if v, ok := got[i]["horas_estimadas"]; !ok || string(v) != want {
			t.Errorf("elemento %d: horas_estimadas = %q (presente: %t), se esperaba %s", i, v, ok, want)
		}
	}
}
