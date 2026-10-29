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

// registradorDefectoFalso usa el dominio real. El proyecto 404 y el elemento 404 no existen; el
// elemento 5 tiene un Sprint de detección (5) vigente. Guarda lo que recibió y cuántas veces.
type registradorDefectoFalso struct {
	err      error
	recibido app.DatosDefecto
	llamadas int
}

func (f *registradorDefectoFalso) Ejecutar(_ context.Context, d app.DatosDefecto) (domain.Defecto, error) {
	f.llamadas++
	f.recibido = d
	switch {
	case f.err != nil:
		return domain.Defecto{}, f.err
	case d.ProyectoID == 404:
		return domain.Defecto{}, domain.ErrProyectoNoEncontrado
	case d.ElementoID == 404:
		return domain.Defecto{}, domain.ErrElementoNoEncontrado
	}
	def, err := domain.NuevoDefecto(d.ProyectoID, d.ElementoID, d.Descripcion, d.Severidad)
	if err != nil {
		return domain.Defecto{}, err
	}
	def.ID = 1
	if d.ElementoID == 5 {
		sprintID := int64(5)
		def.SprintDeteccionID = &sprintID
	}
	return def, nil
}

func postDefecto(t *testing.T, uc registradorDefecto, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarRegistrarDefecto(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return rec
}

const (
	rutaDefecto         = "/proyectos/7/defectos"
	cuerpoDefectoValido = `{"elemento_id":3,"descripcion":"No carga la pantalla de login","severidad":"Alta"}`
)

// Traza: HU-20 CA-20.1, CA-20.6 (spec, secciones 2 y 3: JSON en snake_case)
func TestPostDefecto_DatosValidos_Devuelve201ConElDefecto(t *testing.T) {
	uc := &registradorDefectoFalso{}

	rec := postDefecto(t, uc, rutaDefecto, cuerpoDefectoValido)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	var got map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	want := map[string]any{
		"id": 1.0, "proyecto_id": 7.0, "elemento_id": 3.0,
		"descripcion": "No carga la pantalla de login", "severidad": "Alta", "estado": "Abierto",
		"sprint_deteccion_id": nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("respuesta = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-20 CA-20.1 (RN-20.5: el Sprint de detección se informa cuando hay uno)
func TestPostDefecto_ConSprintDeteccion_LoInformaEnLaRespuesta(t *testing.T) {
	uc := &registradorDefectoFalso{}

	rec := postDefecto(t, uc, "/proyectos/7/defectos", `{"elemento_id":5,"descripcion":"Falla el login","severidad":"Baja"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201; cuerpo: %s", rec.Code, rec.Body)
	}
	var got map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	if got["sprint_deteccion_id"] != 5.0 {
		t.Errorf("sprint_deteccion_id = %v, se esperaba 5", got["sprint_deteccion_id"])
	}
}

// Traza: HU-20 CA-20.1
func TestPostDefecto_DatosValidos_PasaLosDatosAlCasoDeUso(t *testing.T) {
	uc := &registradorDefectoFalso{}

	postDefecto(t, uc, rutaDefecto, cuerpoDefectoValido)

	want := app.DatosDefecto{
		ProyectoID: 7, ElementoID: 3, Descripcion: "No carga la pantalla de login", Severidad: domain.SeveridadAlta,
	}
	if uc.recibido != want {
		t.Errorf("el caso de uso recibió %+v, se esperaba %+v", uc.recibido, want)
	}
}

// Traza: HU-20 CA-20.6, CA-20.5 (el elemento va en el cuerpo, no en la ruta, a diferencia de HU-18)
func TestPostDefecto_ProyectoOElementoInexistente_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		cuerpo  string
		wantMsg string
	}{
		{"proyecto inexistente", "/proyectos/404/defectos", cuerpoDefectoValido, domain.ErrProyectoNoEncontrado.Error()},
		{"elemento inexistente", rutaDefecto, `{"elemento_id":404,"descripcion":"Falla el login","severidad":"Alta"}`, domain.ErrElementoNoEncontrado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postDefecto(t, &registradorDefectoFalso{}, tt.ruta, tt.cuerpo)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("error = %q, se esperaba %q", got, tt.wantMsg)
			}
		})
	}
}

// Traza: HU-20 CA-20.3, CA-20.4
func TestPostDefecto_DatosInvalidos_Devuelve400ConElError(t *testing.T) {
	tests := []struct {
		nombre string
		cuerpo string
		want   error
	}{
		{"descripción vacía", `{"elemento_id":3,"descripcion":"  ","severidad":"Alta"}`, domain.ErrDescripcionDefectoVacia},
		{"severidad inválida", `{"elemento_id":3,"descripcion":"Falla el login","severidad":"urgente"}`, domain.ErrSeveridadInvalida},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postDefecto(t, &registradorDefectoFalso{}, rutaDefecto, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.want.Error() {
				t.Errorf("error = %q, se esperaba %q", got, tt.want.Error())
			}
		})
	}
}

// Traza: HU-20 sección 7 (cuerpo mal formado o id de proyecto inválido → 400 sin llamar al caso de uso)
func TestPostDefecto_PedidoMalFormado_Devuelve400SinLlamarAlCasoDeUso(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		cuerpo  string
		wantMsg string
	}{
		{"JSON mal formado", rutaDefecto, `{"descripcion":`, "el cuerpo debe ser un JSON válido"},
		{"campo desconocido", rutaDefecto, `{"descripcion":"x","comentario":"y"}`, "el cuerpo debe ser un JSON válido"},
		{"elemento_id como texto", rutaDefecto, `{"elemento_id":"3","descripcion":"x","severidad":"Alta"}`, "el cuerpo debe ser un JSON válido"},
		{"id de proyecto no numérico", "/proyectos/abc/defectos", cuerpoDefectoValido, "el id del proyecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorDefectoFalso{}

			rec := postDefecto(t, uc, tt.ruta, tt.cuerpo)

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

func TestPostDefecto_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	rec := postDefecto(t, &registradorDefectoFalso{err: errors.New("disco lleno")}, rutaDefecto, cuerpoDefectoValido)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", rec.Body)
	}
}
