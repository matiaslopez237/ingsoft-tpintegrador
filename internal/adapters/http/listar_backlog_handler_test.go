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

// listadorBacklogFalso devuelve una lista fija, o un error, y guarda lo que recibió.
type listadorBacklogFalso struct {
	elementos []domain.ElementoBacklog
	err       error
	recibido  app.DatosListarBacklog
}

func (f *listadorBacklogFalso) Ejecutar(_ context.Context, d app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
	f.recibido = d
	return f.elementos, f.err
}

func getBacklog(t *testing.T, uc listadorBacklog, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarListarBacklog(mux, uc)
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

// Traza: HU-07 CA-07.1, sección 3: devuelve los elementos completos en el orden del caso de uso.
func TestGetBacklog_DevuelveLosElementos(t *testing.T) {
	cinco := 5
	uc := &listadorBacklogFalso{elementos: []domain.ElementoBacklog{
		{ID: 2, ProyectoID: 1, Titulo: "Login", Descripcion: "Ingresar", Prioridad: domain.PrioridadAlta,
			Estado: domain.EstadoPendiente, StoryPoints: &cinco, Criterios: []string{"Con clave"}},
		{ID: 1, ProyectoID: 1, Titulo: "Reportes", Descripcion: "Ver", Prioridad: domain.PrioridadBaja,
			Estado: domain.EstadoHecho, Criterios: []string{"Exporta PDF"}},
	}}

	rec := getBacklog(t, uc, "/proyectos/1/backlog")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != (app.DatosListarBacklog{ProyectoID: 1}) {
		t.Errorf("datos recibidos = %+v, se esperaba solo el proyecto 1", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `[
		{"id":2,"proyecto_id":1,"titulo":"Login","descripcion":"Ingresar","prioridad":"Alta","estado":"Pendiente","story_points":5,"criterios":["Con clave"]},
		{"id":1,"proyecto_id":1,"titulo":"Reportes","descripcion":"Ver","prioridad":"Baja","estado":"Hecho","story_points":null,"criterios":["Exporta PDF"]}]`)
}

// Traza: HU-07 CA-07.2, CA-07.4: el orden y el estado llegan desde la URL.
func TestGetBacklog_PasaOrdenYEstado(t *testing.T) {
	uc := &listadorBacklogFalso{}

	getBacklog(t, uc, "/proyectos/1/backlog?orden=asc&estado=En+progreso")

	if uc.recibido != (app.DatosListarBacklog{ProyectoID: 1, Orden: "asc", Estado: "En progreso"}) {
		t.Errorf("datos recibidos = %+v", uc.recibido)
	}
}

// Traza: HU-07 CA-07.5: sin elementos responde [] y no null.
func TestGetBacklog_SinElementos_DevuelveArrayVacio(t *testing.T) {
	rec := getBacklog(t, &listadorBacklogFalso{}, "/proyectos/1/backlog")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `[]`)
}

// Traza: HU-07 CA-07.7, CA-07.8, sección 7
func TestGetBacklog_Errores(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		err     error
		status  int
		wantMsg string
	}{
		{"proyecto inexistente", "/proyectos/99/backlog", domain.ErrProyectoNoEncontrado, http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"orden inválido", "/proyectos/1/backlog?orden=ASC", domain.ErrOrdenInvalido, http.StatusBadRequest, domain.ErrOrdenInvalido.Error()},
		{"estado inválido", "/proyectos/1/backlog?estado=pendiente", domain.ErrEstadoElementoInvalido, http.StatusBadRequest, domain.ErrEstadoElementoInvalido.Error()},
		{"id inválido", "/proyectos/abc/backlog", nil, http.StatusBadRequest, "el id del proyecto debe ser un número entero positivo"},
		{"falla interna", "/proyectos/1/backlog", errors.New("disco lleno"), http.StatusInternalServerError, "error interno al listar el backlog"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := getBacklog(t, &listadorBacklogFalso{err: tt.err}, tt.ruta)

			if rec.Code != tt.status {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %s", rec.Code, tt.status, rec.Body)
			}
			var got respuestaError
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decodificar respuesta: %v", err)
			}
			if got.Error != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got.Error, tt.wantMsg)
			}
			if strings.Contains(rec.Body.String(), "disco lleno") {
				t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
			}
		})
	}
}
