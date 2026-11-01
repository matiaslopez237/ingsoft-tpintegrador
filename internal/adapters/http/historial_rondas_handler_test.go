package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// historialRondasFalso devuelve el historial o el error indicados y guarda lo que recibió.
type historialRondasFalso struct {
	historial              []domain.ResumenRonda
	err                    error
	proyectoID, elementoID int64
}

func (f *historialRondasFalso) Ejecutar(_ context.Context, proyectoID, elementoID int64) ([]domain.ResumenRonda, error) {
	f.proyectoID, f.elementoID = proyectoID, elementoID
	return f.historial, f.err
}

// pedirHistorialRondas registra solo la ruta de HU-16 y hace el pedido.
func pedirHistorialRondas(t *testing.T, uc *historialRondasFalso, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarHistorialRondas(mux, uc)
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

// Traza: HU-16 CA-16.7, sección 3: 200 con las rondas en el orden del caso de uso, con la forma del
// resumen de HU-14 y HU-15; la ronda abierta no trae revelación ni valores.
func TestGetHistorialRondas_DevuelveLasRondas(t *testing.T) {
	uc := &historialRondasFalso{historial: []domain.ResumenRonda{
		{ID: 12, ProyectoID: 1, ElementoID: 7, Numero: 2, Estado: domain.RondaAbierta,
			Votantes: []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto"}}},
		{ID: 11, ProyectoID: 1, ElementoID: 7, Numero: 1, Estado: domain.RondaCerrada,
			Votantes: []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto", YaVoto: true}},
			Revelacion: &domain.Revelacion{Votos: []domain.VotoRevelado{{Nombre: "Ana", Valor: 3}, {Nombre: "Beto", Valor: 13}},
				Minimo: 3, Maximo: 13, DiferenciaSignificativa: true}},
	}}

	rec := pedirHistorialRondas(t, uc, "/proyectos/1/backlog/7/rondas")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.proyectoID != 1 || uc.elementoID != 7 {
		t.Errorf("recibió proyecto %d y elemento %d, se esperaba 1 y 7", uc.proyectoID, uc.elementoID)
	}
	mismoJSON(t, rec.Body.String(), `[
		{"id":12,"proyecto_id":1,"elemento_id":7,"numero":2,"estado":"Abierta",
		 "votantes":[{"nombre":"Ana","voto":true},{"nombre":"Beto","voto":false}]},
		{"id":11,"proyecto_id":1,"elemento_id":7,"numero":1,"estado":"Cerrada",
		 "votantes":[{"nombre":"Ana","voto":true},{"nombre":"Beto","voto":true}],
		 "revelacion":{"votos":[{"nombre":"Ana","valor":3},{"nombre":"Beto","valor":13}],"minimo":3,"maximo":13,"diferencia_significativa":true}}
	]`)
}

// Traza: HU-16 CA-16.7 (CA-16.5): sin rondas es un array vacío, no null.
func TestGetHistorialRondas_SinRondas_DevuelveArrayVacio(t *testing.T) {
	rec := pedirHistorialRondas(t, &historialRondasFalso{}, "/proyectos/1/backlog/7/rondas")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	mismoJSON(t, rec.Body.String(), `[]`)
}

// Traza: HU-16 CA-16.7, sección 7: 404 con el mensaje del dominio aunque llegue envuelto; 400 con un
// ID inválido, sin llamar al caso de uso; 500 sin el detalle.
func TestGetHistorialRondas_Errores(t *testing.T) {
	tests := []struct {
		nombre string
		ruta   string
		err    error
		want   int
		msg    string
	}{
		{"proyecto inexistente", "/proyectos/99/backlog/7/rondas", fmt.Errorf("buscar proyecto 99: %w", domain.ErrProyectoNoEncontrado), http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"elemento inexistente o de otro proyecto", "/proyectos/1/backlog/99/rondas", domain.ErrElementoNoEncontrado, http.StatusNotFound, domain.ErrElementoNoEncontrado.Error()},
		{"proyecto inválido", "/proyectos/abc/backlog/7/rondas", nil, http.StatusBadRequest, "proyecto"},
		{"elemento inválido", "/proyectos/1/backlog/abc/rondas", nil, http.StatusBadRequest, "elemento"},
		{"falla interna", "/proyectos/1/backlog/7/rondas", errors.New("disco lleno"), http.StatusInternalServerError, "error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &historialRondasFalso{err: tt.err}

			rec := pedirHistorialRondas(t, uc, tt.ruta)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %s", rec.Code, tt.want, rec.Body)
			}
			cuerpo := rec.Body.String()
			if !strings.Contains(cuerpo, tt.msg) || strings.Contains(cuerpo, "disco lleno") || strings.Contains(cuerpo, "buscar proyecto") {
				t.Errorf("cuerpo = %s: se esperaba %q, sin el detalle interno", cuerpo, tt.msg)
			}
			if tt.want == http.StatusBadRequest && uc.proyectoID != 0 {
				t.Error("con un ID inválido no se tiene que llamar al caso de uso")
			}
		})
	}
}
