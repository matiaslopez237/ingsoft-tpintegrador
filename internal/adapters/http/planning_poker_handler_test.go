package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// pokerFalso implementa los tres casos de uso de HU-14 con el dominio real sobre una sola ronda en
// memoria (ID 1, elemento 7, votantes Ana y Beto). El proyecto 404, el elemento 404 y la ronda 404
// no existen. Con err, todas las operaciones fallan con ese error.
type pokerFalso struct {
	ronda     *domain.Ronda
	err       error
	votos     int
	sinUltima bool // UltimaDelElemento no encuentra rondas
	iniciadas int
}

func nuevoPokerFalso() *pokerFalso { return &pokerFalso{} }

func (f *pokerFalso) fallo(proyectoID int64) error {
	switch {
	case f.err != nil:
		return f.err
	case proyectoID == 404:
		return domain.ErrProyectoNoEncontrado
	}
	return nil
}

func (f *pokerFalso) Ejecutar(_ context.Context, d app.DatosIniciarRonda) (domain.ResumenRonda, error) {
	if err := f.fallo(d.ProyectoID); err != nil {
		return domain.ResumenRonda{}, err
	}
	if d.ElementoID == 404 {
		return domain.ResumenRonda{}, domain.ErrElementoNoEncontrado
	}
	e := domain.ElementoBacklog{ID: d.ElementoID, ProyectoID: d.ProyectoID, Estado: domain.EstadoPendiente}
	r, err := domain.NuevaRonda(e, []domain.Integrante{{Nombre: "Ana"}, {Nombre: "Beto"}}, f.ronda)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	r.ID = int64(r.Numero)
	f.ronda = &r
	f.iniciadas++
	return r.Resumen(), nil
}

// votadorPokerFalso y consultorPokerFalso adaptan el mismo falso a las otras dos interfaces.
type (
	votadorPokerFalso   struct{ *pokerFalso }
	consultorPokerFalso struct{ *pokerFalso }
)

func (f *pokerFalso) buscar(proyectoID, rondaID int64) (*domain.Ronda, error) {
	if err := f.fallo(proyectoID); err != nil {
		return nil, err
	}
	if f.ronda == nil || rondaID != f.ronda.ID || proyectoID != f.ronda.ProyectoID {
		return nil, domain.ErrRondaNoEncontrada
	}
	return f.ronda, nil
}

func (f votadorPokerFalso) Ejecutar(_ context.Context, d app.DatosVoto) (domain.ResumenRonda, error) {
	r, err := f.buscar(d.ProyectoID, d.RondaID)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	if err := r.Votar(d.Integrante, d.Valor); err != nil {
		return domain.ResumenRonda{}, err
	}
	f.votos++
	return r.Resumen(), nil
}

func (f consultorPokerFalso) Ejecutar(_ context.Context, proyectoID, rondaID int64) (domain.ResumenRonda, error) {
	r, err := f.buscar(proyectoID, rondaID)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	return r.Resumen(), nil
}

func (f consultorPokerFalso) UltimaDelElemento(_ context.Context, proyectoID, elementoID int64) (*domain.ResumenRonda, error) {
	if err := f.fallo(proyectoID); err != nil {
		return nil, err
	}
	if elementoID == 404 {
		return nil, domain.ErrElementoNoEncontrado
	}
	if f.sinUltima || f.ronda == nil {
		return nil, nil
	}
	res := f.ronda.Resumen()
	return &res, nil
}

func pedirPoker(t *testing.T, f *pokerFalso, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarPlanningPoker(mux, f, votadorPokerFalso{f}, consultorPokerFalso{f})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo)))
	return rec
}

func decodificar(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var got map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	return got
}

// Traza: HU-14 CA-14.1, CA-14.8 (spec, sección 3: JSON del resumen, sin valores)
func TestPostRonda_Devuelve201ConElResumen(t *testing.T) {
	f := nuevoPokerFalso()

	rec := pedirPoker(t, f, http.MethodPost, "/proyectos/1/backlog/7/rondas", "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	want := map[string]any{
		"id": 1.0, "proyecto_id": 1.0, "elemento_id": 7.0, "numero": 1.0, "estado": "Abierta",
		"votantes": []any{map[string]any{"nombre": "Ana", "voto": false}, map[string]any{"nombre": "Beto", "voto": false}},
	}
	if got := decodificar(t, rec); !reflect.DeepEqual(got, want) {
		t.Errorf("respuesta = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-14 CA-14.2, CA-14.8 (RN-14.3): el voto aparece como "votó", nunca su valor.
func TestPostVoto_Devuelve200SinElValor(t *testing.T) {
	f := nuevoPokerFalso()
	pedirPoker(t, f, http.MethodPost, "/proyectos/1/backlog/7/rondas", "")

	rec := pedirPoker(t, f, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":13}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "13") {
		t.Errorf("la respuesta muestra el valor del voto: %s", rec.Body)
	}
	votantes := decodificar(t, rec)["votantes"].([]any)
	if ana := votantes[0].(map[string]any); ana["voto"] != true || len(ana) != 2 {
		t.Errorf("Ana = %v, se esperaba {nombre, voto: true}", ana)
	}
}

// Traza: HU-14 CA-14.8
func TestGetRonda_Devuelve200ConElResumen(t *testing.T) {
	f := nuevoPokerFalso()
	pedirPoker(t, f, http.MethodPost, "/proyectos/1/backlog/7/rondas", "")
	pedirPoker(t, f, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":8}`)

	rec := pedirPoker(t, f, http.MethodGet, "/proyectos/1/rondas/1", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := decodificar(t, rec); got["estado"] != "Abierta" || strings.Contains(rec.Body.String(), `"valor"`) {
		t.Errorf("respuesta = %v", got)
	}
}

// Traza: HU-14 CA-14.5 a CA-14.8 (sección 7)
func TestPlanningPoker_Errores(t *testing.T) {
	tests := []struct {
		nombre        string
		preparar      bool // abre la ronda 1 antes del pedido
		metodo, ruta  string
		cuerpo        string
		err           error
		wantCode      int
		wantMsg       string
		noLlegaAlCaso bool
	}{
		{"iniciar en proyecto inexistente", false, http.MethodPost, "/proyectos/404/backlog/7/rondas", "", nil, http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error(), false},
		{"iniciar con elemento inexistente", false, http.MethodPost, "/proyectos/1/backlog/404/rondas", "", nil, http.StatusNotFound, domain.ErrElementoNoEncontrado.Error(), false},
		{"iniciar con una ronda abierta", true, http.MethodPost, "/proyectos/1/backlog/7/rondas", "", nil, http.StatusBadRequest, domain.ErrRondaYaAbierta.Error(), false},
		{"iniciar con elemento terminado", false, http.MethodPost, "/proyectos/1/backlog/7/rondas", "", domain.ErrElementoTerminado, http.StatusBadRequest, domain.ErrElementoTerminado.Error(), false},
		{"iniciar sin integrantes", false, http.MethodPost, "/proyectos/1/backlog/7/rondas", "", domain.ErrRondaSinVotantes, http.StatusBadRequest, domain.ErrRondaSinVotantes.Error(), false},
		{"iniciar con id inválido", false, http.MethodPost, "/proyectos/x/backlog/7/rondas", "", nil, http.StatusBadRequest, "", true},
		{"votar en ronda inexistente", true, http.MethodPost, "/proyectos/1/rondas/404/votos", `{"integrante":"Ana","valor":5}`, nil, http.StatusNotFound, domain.ErrRondaNoEncontrada.Error(), false},
		{"votar un valor inválido", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":4}`, nil, http.StatusBadRequest, domain.ErrVotoInvalido.Error(), false},
		{"votar sin ser votante", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Dani","valor":5}`, nil, http.StatusBadRequest, domain.ErrNoEsVotante.Error(), false},
		{"votar en ronda cerrada", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":5}`, domain.ErrRondaCerrada, http.StatusBadRequest, domain.ErrRondaCerrada.Error(), false},
		{"votar con JSON mal formado", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":`, nil, http.StatusBadRequest, "", true},
		{"votar con un campo desconocido", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":5,"otro":1}`, nil, http.StatusBadRequest, "", true},
		{"votar con un valor que no es número", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":"5"}`, nil, http.StatusBadRequest, "", true},
		{"votar con ronda inválida", true, http.MethodPost, "/proyectos/1/rondas/x/votos", `{"integrante":"Ana","valor":5}`, nil, http.StatusBadRequest, "", true},
		{"consultar ronda de otro proyecto", true, http.MethodGet, "/proyectos/2/rondas/1", "", nil, http.StatusNotFound, domain.ErrRondaNoEncontrada.Error(), false},
		{"falla interna al votar", true, http.MethodPost, "/proyectos/1/rondas/1/votos", `{"integrante":"Ana","valor":5}`, errors.New("disco lleno"), http.StatusInternalServerError, "error interno", false},
		{"falla interna al consultar", true, http.MethodGet, "/proyectos/1/rondas/1", "", errors.New("disco lleno"), http.StatusInternalServerError, "error interno", false},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoPokerFalso()
			if tt.preparar {
				pedirPoker(t, f, http.MethodPost, "/proyectos/1/backlog/7/rondas", "")
			}
			f.err = tt.err
			antes := f.votos + f.iniciadas

			rec := pedirPoker(t, f, tt.metodo, tt.ruta, tt.cuerpo)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %s", rec.Code, tt.wantCode, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tt.wantMsg) {
				t.Errorf("cuerpo = %s, se esperaba que contenga %q", rec.Body, tt.wantMsg)
			}
			if tt.noLlegaAlCaso && f.votos+f.iniciadas != antes {
				t.Error("un pedido mal formado no tendría que llegar al caso de uso")
			}
			if tt.wantCode == http.StatusInternalServerError && strings.Contains(rec.Body.String(), "disco lleno") {
				t.Error("la respuesta 500 no tiene que mostrar el detalle interno")
			}
		})
	}
}
