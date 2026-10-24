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

// creadorElementosFalso usa el dominio real y un ID fijo; el proyecto 404 no existe.
type creadorElementosFalso struct {
	err error
}

func (c *creadorElementosFalso) Ejecutar(_ context.Context, d app.DatosNuevoElemento) (domain.ElementoBacklog, error) {
	if c.err != nil {
		return domain.ElementoBacklog{}, c.err
	}
	if d.ProyectoID == 404 {
		return domain.ElementoBacklog{}, domain.ErrProyectoNoEncontrado
	}
	e, err := domain.NuevoElementoBacklog(d.ProyectoID, d.Titulo, d.Descripcion, d.Prioridad, d.StoryPoints, d.Criterios)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}
	e.ID = 1
	return e, nil
}

func postBacklog(t *testing.T, uc creadorElementos, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarCrearElementoBacklog(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return rec
}

const cuerpoValido = `{"titulo":"Login","descripcion":"Ingresar","prioridad":"Alta","story_points":5,"criterios":["Con email","Con clave"]}`

// Traza: HU-05 CA-05.1
func TestPostBacklog_DatosValidos_Devuelve201ConElemento(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{}, "/proyectos/7/backlog", cuerpoValido)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	var got respuestaElemento
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	if got.ID != 1 || got.ProyectoID != 7 || got.Titulo != "Login" || got.Estado != "Pendiente" ||
		got.Prioridad != "Alta" || got.StoryPoints == nil || *got.StoryPoints != 5 || len(got.Criterios) != 2 {
		t.Errorf("respuesta inesperada: %+v", got)
	}
}

// Traza: HU-05 CA-05.2
func TestPostBacklog_SinStoryPoints_DevuelveNull(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{}, "/proyectos/7/backlog",
		`{"titulo":"Login","descripcion":"Ingresar","prioridad":"Baja","criterios":["Uno"]}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), `"story_points":null`) {
		t.Errorf("se esperaba story_points null, cuerpo: %s", rec.Body)
	}
}

// Traza: HU-05 CA-05.3
func TestPostBacklog_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{}, "/proyectos/404/backlog", cuerpoValido)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}

// Traza: HU-05 CA-05.4 a CA-05.9
func TestPostBacklog_DatosInvalidos_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		cuerpo  string
		creador *creadorElementosFalso
		wantMsg string
	}{
		{"título vacío", "/proyectos/7/backlog", `{"titulo":" ","descripcion":"D","prioridad":"Alta","criterios":["C"]}`, &creadorElementosFalso{}, domain.ErrTituloVacio.Error()},
		{"descripción vacía", "/proyectos/7/backlog", `{"titulo":"T","prioridad":"Alta","criterios":["C"]}`, &creadorElementosFalso{}, domain.ErrDescripcionVacia.Error()},
		{"prioridad inválida", "/proyectos/7/backlog", `{"titulo":"T","descripcion":"D","prioridad":"alta","criterios":["C"]}`, &creadorElementosFalso{}, domain.ErrPrioridadInvalida.Error()},
		{"story points en cero", "/proyectos/7/backlog", `{"titulo":"T","descripcion":"D","prioridad":"Alta","story_points":0,"criterios":["C"]}`, &creadorElementosFalso{}, domain.ErrStoryPointsInvalidos.Error()},
		{"sin criterios", "/proyectos/7/backlog", `{"titulo":"T","descripcion":"D","prioridad":"Alta","criterios":[]}`, &creadorElementosFalso{}, domain.ErrSinCriterios.Error()},
		{"criterio vacío", "/proyectos/7/backlog", `{"titulo":"T","descripcion":"D","prioridad":"Alta","criterios":["C",""]}`, &creadorElementosFalso{}, domain.ErrCriterioVacio.Error()},
		{"título duplicado", "/proyectos/7/backlog", cuerpoValido, &creadorElementosFalso{err: domain.ErrTituloDuplicado}, domain.ErrTituloDuplicado.Error()},
		{"story points no enteros", "/proyectos/7/backlog", `{"titulo":"T","descripcion":"D","prioridad":"Alta","story_points":2.5,"criterios":["C"]}`, &creadorElementosFalso{}, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", "/proyectos/7/backlog", `{"titulo":"T","estado":"Hecho"}`, &creadorElementosFalso{}, "el cuerpo debe ser un JSON válido"},
		{"id de proyecto inválido", "/proyectos/abc/backlog", cuerpoValido, &creadorElementosFalso{}, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postBacklog(t, tt.creador, tt.ruta, tt.cuerpo)

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

func TestPostBacklog_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postBacklog(t, &creadorElementosFalso{err: errors.New("disco lleno")}, "/proyectos/7/backlog", cuerpoValido)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
