package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// comparadorEsfuerzoFalso devuelve la comparación que se le indica y guarda lo que recibió.
type comparadorEsfuerzoFalso struct {
	comparacion domain.ComparacionEsfuerzo
	err         error
	proyectoID  int64
	elementoID  int64
	llamado     bool
}

func (f *comparadorEsfuerzoFalso) Ejecutar(_ context.Context, proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error) {
	f.llamado = true
	f.proyectoID, f.elementoID = proyectoID, elementoID
	return f.comparacion, f.err
}

func pHoras(h float64) *float64 { return &h }

const rutaCompararEsfuerzo = "/proyectos/7/backlog/3/esfuerzo/comparacion"

func getCompararEsfuerzo(t *testing.T, uc comparadorEsfuerzo, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarCompararEsfuerzo(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

// camposJSON lee el cuerpo sin interpretar los valores, para comparar el JSON literal: un null, un
// número y un texto se distinguen, y un -0 se vería tal cual.
func camposJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var crudo map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &crudo); err != nil {
		t.Fatalf("el cuerpo no es un objeto JSON: %v; cuerpo: %s", err, rec.Body)
	}
	campos := make(map[string]string, len(crudo))
	for k, v := range crudo {
		campos[k] = string(v)
	}
	return campos
}

// Cada caso usa la función real del dominio, para probar el JSON de lo que el caso de uso devuelve de
// verdad. Los valores esperados son el JSON literal: 25,0 sale como el número 25.
// Traza: HU-19 CA-19.17 (RN-19.5, RN-19.6, RN-19.7, RN-19.8, RN-19.9; spec, tabla de la sección 3)
func TestGetCompararEsfuerzo_Devuelve200ConLaComparacion(t *testing.T) {
	tests := []struct {
		nombre     string
		estimadas  *float64
		reales     float64
		wantEst    string
		wantReal   string
		wantHoras  string
		wantPorcen string
	}{
		{"se trabajó más de lo estimado", pHoras(8), 10, "8", "10", "2", "25"},
		{"se trabajó menos de lo estimado", pHoras(8), 6, "8", "6", "-2", "-25"},
		{"coincide", pHoras(8), 8, "8", "8", "0", "0"},
		{"sin esfuerzo registrado", pHoras(8), 0, "8", "0", "-8", "-100"},
		{"sin estimar, con esfuerzo", nil, 10, "null", "10", "null", "null"},
		{"sin estimar y sin esfuerzo", nil, 0, "null", "0", "null", "null"},
		{"estimado 0", pHoras(0), 10, "0", "10", "10", "null"},
		{"porcentaje con decimal", pHoras(3), 4, "3", "4", "1", "33.3"},
		{"horas y porcentaje con decimal", pHoras(40), 40.5, "40", "40.5", "0.5", "1.3"},
		{"negativos con decimal", pHoras(40), 39.5, "40", "39.5", "-0.5", "-1.3"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &comparadorEsfuerzoFalso{comparacion: domain.CompararEsfuerzo(tt.estimadas, tt.reales)}

			rec := getCompararEsfuerzo(t, uc, rutaCompararEsfuerzo)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
			}
			want := map[string]string{
				"proyecto_id": "7", "elemento_id": "3",
				"horas_estimadas": tt.wantEst, "horas_reales": tt.wantReal,
				"desviacion_horas": tt.wantHoras, "desviacion_porcentaje": tt.wantPorcen,
			}
			got := camposJSON(t, rec)
			if len(got) != len(want) {
				t.Errorf("campos = %v, se esperaban %v", got, want)
			}
			for campo, valor := range want {
				if got[campo] != valor {
					t.Errorf("%s = %s, se esperaba %s", campo, got[campo], valor)
				}
			}
			if uc.proyectoID != 7 || uc.elementoID != 3 {
				t.Errorf("el caso de uso recibió el proyecto %d y el elemento %d", uc.proyectoID, uc.elementoID)
			}
		})
	}
}

// Traza: HU-19 CA-19.17 (RN-19.6; sección 3): el porcentaje es un número en JSON; el formato con coma es de la pantalla.
func TestGetCompararEsfuerzo_ElPorcentajeEsUnNumero(t *testing.T) {
	uc := &comparadorEsfuerzoFalso{comparacion: domain.CompararEsfuerzo(pHoras(3), 4)}

	rec := getCompararEsfuerzo(t, uc, rutaCompararEsfuerzo)

	var got struct {
		Porcentaje float64 `json:"desviacion_porcentaje"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("no se pudo leer el porcentaje como número: %v; cuerpo: %s", err, rec.Body)
	}
	if got.Porcentaje != 33.3 {
		t.Errorf("porcentaje = %v, se esperaba 33.3", got.Porcentaje)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type = %q, se esperaba JSON", ct)
	}
}

// encoding/json escribe un -0 de Go como -0. Con 2000 estimadas y media hora de menos, el porcentaje
// redondeado sería -0: el dominio lo evita (RN-19.6), y acá se verifica lo que le llega al cliente.
// Traza: HU-19 CA-19.17 (RN-19.6)
func TestGetCompararEsfuerzo_NuncaEscribeCeroNegativo(t *testing.T) {
	tests := []struct {
		nombre            string
		estimadas, reales float64
		wantHoras         string
	}{
		{"el redondeo de un negativo da cero", 2000, 1999.5, "-0.5"},
		{"el redondeo de un positivo da cero", 2000, 2000.5, "0.5"},
		{"sin desviación", 8, 8, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &comparadorEsfuerzoFalso{comparacion: domain.CompararEsfuerzo(&tt.estimadas, tt.reales)}

			got := camposJSON(t, getCompararEsfuerzo(t, uc, rutaCompararEsfuerzo))

			if got["desviacion_porcentaje"] != "0" {
				t.Errorf("desviacion_porcentaje = %s, se esperaba 0", got["desviacion_porcentaje"])
			}
			if got["desviacion_horas"] != tt.wantHoras {
				t.Errorf("desviacion_horas = %s, se esperaba %s", got["desviacion_horas"], tt.wantHoras)
			}
		})
	}
}

// El mensaje es el del dominio, aunque el error llegue envuelto: se reconoce con errors.Is.
// Traza: HU-19 CA-19.9, CA-19.10, CA-19.17 (RN-19.10; sección 7)
func TestGetCompararEsfuerzo_NoEncontrado_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre  string
		err     error
		wantMsg string
	}{
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, "proyecto no encontrado"},
		{"elemento inexistente", domain.ErrElementoNoEncontrado, "elemento de backlog no encontrado"},
		{"proyecto, envuelto", fmt.Errorf("buscar proyecto 7: %w", domain.ErrProyectoNoEncontrado), "proyecto no encontrado"},
		{"elemento, envuelto", fmt.Errorf("buscar elemento 3: %w", domain.ErrElementoNoEncontrado), "elemento de backlog no encontrado"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := getCompararEsfuerzo(t, &comparadorEsfuerzoFalso{err: tt.err}, rutaCompararEsfuerzo)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
		})
	}
}

// IDs inválidos en la ruta: 400, como en los demás handlers, sin llamar al caso de uso.
// Traza: HU-19 CA-19.17 (sección 7)
func TestGetCompararEsfuerzo_IDInvalido_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"id de proyecto no numérico", "/proyectos/abc/backlog/3/esfuerzo/comparacion", "el id del proyecto debe ser un número entero positivo"},
		{"id de proyecto en cero", "/proyectos/0/backlog/3/esfuerzo/comparacion", "el id del proyecto debe ser un número entero positivo"},
		{"id de elemento no numérico", "/proyectos/7/backlog/abc/esfuerzo/comparacion", "el id del elemento debe ser un número entero positivo"},
		{"id de elemento negativo", "/proyectos/7/backlog/-1/esfuerzo/comparacion", "el id del elemento debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &comparadorEsfuerzoFalso{}

			rec := getCompararEsfuerzo(t, uc, tt.ruta)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
			if uc.llamado {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}

// Cualquier otro error es una falla interna: 500, sin mostrar el detalle. El texto no se mira: uno
// que solo se parece a "proyecto no encontrado" sigue siendo una falla interna.
// Traza: HU-19 CA-19.17 (sección 7)
func TestGetCompararEsfuerzo_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	for _, err := range []error{
		errors.New("consultar horas reales: database is locked"),
		errors.New("proyecto no encontrado"),
	} {
		rec := getCompararEsfuerzo(t, &comparadorEsfuerzoFalso{err: err}, rutaCompararEsfuerzo)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("%v: status = %d, se esperaba 500; cuerpo: %s", err, rec.Code, rec.Body)
		}
		if got := mensajeDeError(t, rec); got != "error interno al comparar el esfuerzo" {
			t.Errorf("%v: mensaje = %q: no tiene que mostrar el detalle interno", err, got)
		}
	}
}
