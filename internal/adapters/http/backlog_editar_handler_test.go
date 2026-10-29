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

// editorElementosFalso usa el dominio real sobre un elemento 7 del proyecto 1.
type editorElementosFalso struct {
	err      error
	recibido app.DatosEditarElemento
}

func elementoLogin() domain.ElementoBacklog {
	cinco := 5
	e, _ := domain.NuevoElementoBacklog(1, "Login", "Ingresar", domain.PrioridadAlta, &cinco, []string{"Con email"})
	e.ID = 7
	return e
}

func (f *editorElementosFalso) Ejecutar(_ context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error) {
	f.recibido = d
	if f.err != nil {
		return domain.ElementoBacklog{}, f.err
	}
	if d.ProyectoID != 1 || d.ElementoID != 7 {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return elementoLogin().Editar(d.Cambios)
}

type eliminadorElementosFalso struct {
	err      error
	recibido app.DatosEliminarElemento
}

func (f *eliminadorElementosFalso) Ejecutar(_ context.Context, d app.DatosEliminarElemento) error {
	f.recibido = d
	return f.err
}

func pedirBacklogEdicion(t *testing.T, editar editorElementos, eliminar eliminadorElementos, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarEditarEliminarElementoBacklog(mux, editar, eliminar)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo)))
	return rec
}

// Traza: HU-06 CA-06.1, CA-06.2
func TestPatchElemento_Valido_Devuelve200YSoloEnviaLoRecibido(t *testing.T) {
	editor := &editorElementosFalso{}
	rec := pedirBacklogEdicion(t, editor, &eliminadorElementosFalso{}, http.MethodPatch, "/proyectos/1/backlog/7",
		`{"prioridad":"Baja","estado":"En progreso"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %s", rec.Code, rec.Body)
	}
	var got respuestaElemento
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 7 || got.Prioridad != "Baja" || got.Estado != "En progreso" || got.Titulo != "Login" {
		t.Errorf("respuesta: %+v", got)
	}
	c := editor.recibido.Cambios
	if c.Titulo != nil || c.Descripcion != nil || c.CambiarStoryPoints || c.Criterios != nil || c.Prioridad == nil || c.Estado == nil {
		t.Errorf("cambios enviados al caso de uso: %+v", c)
	}
}

// Traza: HU-06 CA-06.3: story_points en null quita la estimación; ausente no la toca.
func TestPatchElemento_StoryPoints(t *testing.T) {
	tests := []struct {
		nombre      string
		cuerpo      string
		wantCambiar bool
		wantNil     bool
	}{
		{"null quita la estimación", `{"story_points":null}`, true, true},
		{"número la cambia", `{"story_points":8}`, true, false},
		{"ausente no la toca", `{"titulo":"Otro"}`, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			editor := &editorElementosFalso{}
			rec := pedirBacklogEdicion(t, editor, &eliminadorElementosFalso{}, http.MethodPatch, "/proyectos/1/backlog/7", tt.cuerpo)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo %s", rec.Code, rec.Body)
			}
			c := editor.recibido.Cambios
			if c.CambiarStoryPoints != tt.wantCambiar || (c.StoryPoints == nil) != tt.wantNil {
				t.Errorf("CambiarStoryPoints=%v StoryPoints=%v", c.CambiarStoryPoints, c.StoryPoints)
			}
		})
	}
}

// Traza: HU-06 CA-06.4 a CA-06.7
func TestPatchElemento_Errores(t *testing.T) {
	tests := []struct {
		nombre   string
		ruta     string
		cuerpo   string
		err      error
		wantCode int
		wantMsg  string
	}{
		{"sin cambios", "/proyectos/1/backlog/7", `{}`, nil, http.StatusBadRequest, domain.ErrSinCambios.Error()},
		{"título vacío", "/proyectos/1/backlog/7", `{"titulo":""}`, nil, http.StatusBadRequest, domain.ErrTituloVacio.Error()},
		{"estado inválido", "/proyectos/1/backlog/7", `{"estado":"Cerrado"}`, nil, http.StatusBadRequest, domain.ErrEstadoElementoInvalido.Error()},
		{"criterios vacíos", "/proyectos/1/backlog/7", `{"criterios":[]}`, nil, http.StatusBadRequest, domain.ErrSinCriterios.Error()},
		{"story points no enteros", "/proyectos/1/backlog/7", `{"story_points":2.5}`, nil, http.StatusBadRequest, "story_points debe ser un entero o null"},
		{"título duplicado", "/proyectos/1/backlog/7", `{"titulo":"Registro"}`, domain.ErrTituloDuplicado, http.StatusBadRequest, domain.ErrTituloDuplicado.Error()},
		{"de otro proyecto", "/proyectos/2/backlog/7", `{"titulo":"Otro"}`, nil, http.StatusNotFound, domain.ErrElementoNoEncontrado.Error()},
		{"campo desconocido", "/proyectos/1/backlog/7", `{"id":9}`, nil, http.StatusBadRequest, "el cuerpo debe ser un JSON válido"},
		{"id de elemento inválido", "/proyectos/1/backlog/x", `{"titulo":"Otro"}`, nil, http.StatusBadRequest, "el id del elemento debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirBacklogEdicion(t, &editorElementosFalso{err: tt.err}, &eliminadorElementosFalso{}, http.MethodPatch, tt.ruta, tt.cuerpo)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo %s", rec.Code, tt.wantCode, rec.Body)
			}
			var got respuestaError
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if got.Error != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got.Error, tt.wantMsg)
			}
		})
	}
}

// Traza: HU-06 CA-06.7 a CA-06.9
func TestDeleteElemento(t *testing.T) {
	tests := []struct {
		nombre   string
		err      error
		wantCode int
	}{
		{"sin esfuerzo", nil, http.StatusNoContent},
		{"con esfuerzo", domain.ErrElementoConEsfuerzo, http.StatusConflict},
		{"en un sprint activo", domain.ErrElementoEnSprintActivo, http.StatusConflict}, // RN-06.8
		{"inexistente", domain.ErrElementoNoEncontrado, http.StatusNotFound},
		{"falla interna", errors.New("disco lleno"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			eliminador := &eliminadorElementosFalso{err: tt.err}
			rec := pedirBacklogEdicion(t, &editorElementosFalso{}, eliminador, http.MethodDelete, "/proyectos/1/backlog/7", "")
			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, se esperaba %d; cuerpo %s", rec.Code, tt.wantCode, rec.Body)
			}
			if eliminador.recibido != (app.DatosEliminarElemento{ProyectoID: 1, ElementoID: 7}) {
				t.Errorf("datos recibidos: %+v", eliminador.recibido)
			}
			if strings.Contains(rec.Body.String(), "disco lleno") {
				t.Error("la respuesta no debería exponer el error interno")
			}
		})
	}
}
