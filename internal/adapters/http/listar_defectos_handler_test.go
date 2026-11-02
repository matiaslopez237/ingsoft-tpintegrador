package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type listadorDefectosFalso struct {
	defectos []domain.Defecto
	err      error
}

func (f *listadorDefectosFalso) Ejecutar(_ context.Context, _ int64) ([]domain.Defecto, error) {
	return f.defectos, f.err
}

func getListarDefectos(t *testing.T, uc listadorDefectos, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarListarDefectos(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

// Traza: HU-21 CA-21.8
func TestGetListarDefectos_ConDefectos_Devuelve200ConJSON(t *testing.T) {
	uc := &listadorDefectosFalso{defectos: []domain.Defecto{
		{ID: 1, ProyectoID: 7, Estado: domain.DefectoAbierto},
		{ID: 2, ProyectoID: 7, Estado: domain.DefectoResuelto},
	}}

	rec := getListarDefectos(t, uc, "/proyectos/7/defectos")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	var got []respuestaDefecto
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decodificar cuerpo: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("se esperaban 2 defectos, hay %d", len(got))
	}
}

// Traza: HU-21 CA-21.8 (lista vacía si no tiene ninguno)
func TestGetListarDefectos_SinDefectos_Devuelve200ConListaVacia(t *testing.T) {
	uc := &listadorDefectosFalso{}

	rec := getListarDefectos(t, uc, "/proyectos/7/defectos")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := rec.Body.String(); got != "[]\n" && got != "[]" {
		t.Errorf("cuerpo = %q, se esperaba una lista vacía", got)
	}
}

// Traza: HU-21 CA-21.8: proyecto inexistente, 404.
func TestGetListarDefectos_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := getListarDefectos(t, &listadorDefectosFalso{err: domain.ErrProyectoNoEncontrado}, "/proyectos/99/defectos")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
}
