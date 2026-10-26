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
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// registradorEsfuerzoFalso usa el dominio real con un "hoy" fijo y un ID fijo. El proyecto 404 y
// el elemento 404 no existen, y "Caro" no es integrante. Guarda lo que recibió y cuántas veces.
type registradorEsfuerzoFalso struct {
	err      error
	recibido app.DatosRegistroEsfuerzo
	llamadas int
}

var hoyHandler = time.Date(2026, time.October, 25, 12, 0, 0, 0, time.Local)

func (f *registradorEsfuerzoFalso) Ejecutar(_ context.Context, d app.DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error) {
	f.llamadas++
	f.recibido = d
	switch {
	case f.err != nil:
		return domain.RegistroEsfuerzo{}, f.err
	case d.ProyectoID == 404:
		return domain.RegistroEsfuerzo{}, domain.ErrProyectoNoEncontrado
	case d.ElementoID == 404:
		return domain.RegistroEsfuerzo{}, domain.ErrElementoNoEncontrado
	}
	r, err := domain.NuevoRegistroEsfuerzo(d.ProyectoID, d.ElementoID, d.Integrante, d.Fecha, d.Actividad, d.Horas, hoyHandler)
	if err != nil {
		return domain.RegistroEsfuerzo{}, err
	}
	if r.Integrante == "Caro" {
		return domain.RegistroEsfuerzo{}, domain.ErrIntegranteNoEncontrado
	}
	r.ID = 1
	return r, nil
}

func postEsfuerzo(t *testing.T, uc registradorEsfuerzo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarRegistrarEsfuerzo(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return rec
}

const (
	rutaEsfuerzo         = "/proyectos/7/backlog/3/esfuerzo"
	cuerpoEsfuerzoValido = `{"integrante":"Ana","fecha":"2026-10-25","actividad":"Desarrollo","horas":1.5}`
)

// Traza: HU-18 CA-18.15, CA-18.1 (spec, secciones 2 y 3: JSON en snake_case, fecha AAAA-MM-DD)
func TestPostEsfuerzo_DatosValidos_Devuelve201ConElRegistro(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}

	rec := postEsfuerzo(t, uc, rutaEsfuerzo, cuerpoEsfuerzoValido)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	var got map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	want := map[string]any{
		"id": 1.0, "proyecto_id": 7.0, "elemento_id": 3.0, "integrante": "Ana",
		"fecha": "2026-10-25", "actividad": "Desarrollo", "horas": 1.5,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("respuesta = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-18 RN-18.4 (la fecha del cuerpo se interpreta como un día de la zona del servidor)
func TestPostEsfuerzo_DatosValidos_PasaLosDatosAlCasoDeUso(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}

	postEsfuerzo(t, uc, rutaEsfuerzo, cuerpoEsfuerzoValido)

	want := app.DatosRegistroEsfuerzo{
		ProyectoID: 7, ElementoID: 3, Integrante: "Ana",
		Fecha:     time.Date(2026, time.October, 25, 0, 0, 0, 0, time.Local),
		Actividad: "Desarrollo", Horas: 1.5,
	}
	if uc.recibido != want {
		t.Errorf("el caso de uso recibió %+v, se esperaba %+v", uc.recibido, want)
	}
}

// Traza: HU-18 CA-18.15, CA-18.7
func TestPostEsfuerzo_ProyectoOElementoInexistente_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"proyecto inexistente", "/proyectos/404/backlog/3/esfuerzo", domain.ErrProyectoNoEncontrado.Error()},
		{"elemento inexistente", "/proyectos/7/backlog/404/esfuerzo", domain.ErrElementoNoEncontrado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postEsfuerzo(t, &registradorEsfuerzoFalso{}, tt.ruta, cuerpoEsfuerzoValido)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got, tt.wantMsg)
			}
		})
	}
}

// Traza: HU-18 CA-18.15, CA-18.3 a CA-18.6, CA-18.8 a CA-18.10 (sección 7: errores de datos → 400,
// incluido el integrante que no pertenece, porque viene en el cuerpo y no en la ruta)
func TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError(t *testing.T) {
	tests := []struct {
		nombre string
		cuerpo string
		want   error
	}{
		{"integrante vacío", `{"integrante":"  ","fecha":"2026-10-25","actividad":"Desarrollo","horas":2}`, domain.ErrNombreIntegranteVacio},
		{"sin fecha", `{"integrante":"Ana","actividad":"Desarrollo","horas":2}`, domain.ErrFechaEsfuerzoObligatoria},
		{"fecha futura", `{"integrante":"Ana","fecha":"2026-10-26","actividad":"Desarrollo","horas":2}`, domain.ErrFechaEsfuerzoFutura},
		{"actividad vacía", `{"integrante":"Ana","fecha":"2026-10-25","actividad":"","horas":2}`, domain.ErrActividadVacia},
		{"horas en cero", `{"integrante":"Ana","fecha":"2026-10-25","actividad":"Desarrollo","horas":0}`, domain.ErrHorasInvalidas},
		{"horas no múltiplo de media hora", `{"integrante":"Ana","fecha":"2026-10-25","actividad":"Desarrollo","horas":1.25}`, domain.ErrHorasNoMultiploMediaHora},
		{"más de 24 horas", `{"integrante":"Ana","fecha":"2026-10-25","actividad":"Desarrollo","horas":24.5}`, domain.ErrHorasExcedenMaximo},
		{"integrante que no pertenece", `{"integrante":"Caro","fecha":"2026-10-25","actividad":"Desarrollo","horas":2}`, domain.ErrIntegranteNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postEsfuerzo(t, &registradorEsfuerzoFalso{}, rutaEsfuerzo, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.want.Error() {
				t.Errorf("error = %q, se esperaba %q", got, tt.want.Error())
			}
		})
	}
}

// Traza: HU-18 sección 7 (cuerpo mal formado o ids inválidos → 400 sin llamar al caso de uso)
func TestPostEsfuerzo_PedidoMalFormado_Devuelve400SinLlamarAlCasoDeUso(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		cuerpo  string
		wantMsg string
	}{
		{"JSON mal formado", rutaEsfuerzo, `{"integrante":`, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", rutaEsfuerzo, `{"integrante":"Ana","minutos":30}`, "el cuerpo debe ser un JSON válido"},
		{"horas como texto", rutaEsfuerzo, `{"integrante":"Ana","fecha":"2026-10-25","actividad":"D","horas":"2"}`, "el cuerpo debe ser un JSON válido"},
		{"fecha con otro formato", rutaEsfuerzo, `{"integrante":"Ana","fecha":"25/10/2026","actividad":"D","horas":2}`, "la fecha debe tener el formato AAAA-MM-DD"},
		{"fecha con hora", rutaEsfuerzo, `{"integrante":"Ana","fecha":"2026-10-25T10:00:00Z","actividad":"D","horas":2}`, "la fecha debe tener el formato AAAA-MM-DD"},
		{"id de proyecto no numérico", "/proyectos/abc/backlog/3/esfuerzo", cuerpoEsfuerzoValido, "el id del proyecto debe ser un número entero positivo"},
		{"id de elemento no numérico", "/proyectos/7/backlog/x/esfuerzo", cuerpoEsfuerzoValido, "el id del elemento debe ser un número entero positivo"},
		{"id de elemento en cero", "/proyectos/7/backlog/0/esfuerzo", cuerpoEsfuerzoValido, "el id del elemento debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorEsfuerzoFalso{}

			rec := postEsfuerzo(t, uc, tt.ruta, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got, tt.wantMsg)
			}
			if uc.llamadas != 0 {
				t.Errorf("no se esperaba llamar al caso de uso, se llamó %d veces", uc.llamadas)
			}
		})
	}
}

func TestPostEsfuerzo_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postEsfuerzo(t, &registradorEsfuerzoFalso{err: errors.New("disco lleno")}, rutaEsfuerzo, cuerpoEsfuerzoValido)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}

func mensajeDeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var got respuestaError
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	return got.Error
}
