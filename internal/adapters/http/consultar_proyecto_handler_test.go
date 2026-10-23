package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorFalso deja probar el handler sin base.
type consultorFalso struct {
	proyecto domain.Proyecto
	err      error
	recibido int64
	llamado  bool
}

func (f *consultorFalso) Ejecutar(_ context.Context, id int64) (domain.Proyecto, error) {
	f.llamado = true
	f.recibido = id
	if f.err != nil {
		return domain.Proyecto{}, f.err
	}
	return f.proyecto, nil
}

// getProyecto registra la ruta con RegistrarConsultarProyecto, igual que main.go, así los tests
// cubren también el patrón de la ruta y r.PathValue("id") funciona como en la aplicación.
func getProyecto(t *testing.T, uc consultorProyectos, id string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarConsultarProyecto(mux, uc)
	req := httptest.NewRequest(http.MethodGet, "/proyectos/"+id, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// mismoJSON compara el cuerpo con el JSON esperado ya decodificados: no importan los espacios ni
// el orden de las claves, pero sí distingue [] de null.
func mismoJSON(t *testing.T, cuerpo, esperado string) {
	t.Helper()
	var got, want any
	if err := json.Unmarshal([]byte(cuerpo), &got); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v; cuerpo: %q", err, cuerpo)
	}
	if err := json.Unmarshal([]byte(esperado), &want); err != nil {
		t.Fatalf("el JSON esperado del test es inválido: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JSON = %s\nse esperaba %s", cuerpo, esperado)
	}
}

func fecha(s string) time.Time {
	f, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return f
}

// Traza: HU-04 CA-04.1, CA-04.2
func TestGetProyecto_Existente_Devuelve200ConElJSONDeLaSpec(t *testing.T) {
	fin := fecha("2026-11-30")
	uc := &consultorFalso{proyecto: domain.Proyecto{
		ID:          7,
		Nombre:      "TP Integrador",
		FechaInicio: fecha("2026-03-01"),
		FechaFin:    &fin,
		Estado:      domain.EstadoActivo,
		Integrantes: []domain.Integrante{
			{Nombre: "Beto", Roles: []string{"Tester"}},
			{Nombre: "Ana", Roles: []string{"Tester", "Desarrolladora"}},
		},
	}}

	rec := getProyecto(t, uc, "7")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido != 7 {
		t.Errorf("id recibido por el caso de uso = %d, se esperaba 7", uc.recibido)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, se esperaba application/json", ct)
	}
	mismoJSON(t, rec.Body.String(), `{
		"id": 7,
		"nombre": "TP Integrador",
		"fecha_inicio": "2026-03-01",
		"fecha_fin": "2026-11-30",
		"estado": "Activo",
		"integrantes": [
			{"nombre": "Beto", "roles": ["Tester"]},
			{"nombre": "Ana", "roles": ["Tester", "Desarrolladora"]}
		]
	}`)
}

// Traza: HU-04 CA-04.3, CA-04.4
func TestGetProyecto_SinIntegrantesNiFechaFin_DevuelveArrayVacioYFechaFinNull(t *testing.T) {
	// Integrantes nil, como lo devuelve la persistencia cuando no hay filas.
	uc := &consultorFalso{proyecto: domain.Proyecto{
		ID:          7,
		Nombre:      "TP Integrador",
		FechaInicio: fecha("2026-03-01"),
		Estado:      domain.EstadoActivo,
	}}

	rec := getProyecto(t, uc, "7")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	mismoJSON(t, rec.Body.String(), `{
		"id": 7,
		"nombre": "TP Integrador",
		"fecha_inicio": "2026-03-01",
		"fecha_fin": null,
		"estado": "Activo",
		"integrantes": []
	}`)
}

// Traza: HU-04 CA-04.5
func TestGetProyecto_Inexistente_Devuelve404ConElMensajeDelError(t *testing.T) {
	// Envuelto, como lo devuelve el caso de uso ConsultarProyecto.
	uc := &consultorFalso{err: fmt.Errorf("consultar proyecto 9999: %w", domain.ErrProyectoNoEncontrado)}

	rec := getProyecto(t, uc, "9999")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
	}
	var got respuestaError
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar respuesta: %v", err)
	}
	if want := domain.ErrProyectoNoEncontrado.Error(); got.Error != want {
		t.Errorf("error = %q, se esperaba %q", got.Error, want)
	}
}

// Traza: HU-04 secciones 6 y 7 (id inválido); no tiene un CA propio en la spec.
func TestGetProyecto_IDInvalido_Devuelve400SinLlamarAlCasoDeUso(t *testing.T) {
	for _, id := range []string{"abc", "0", "-1", "99999999999999999999"} {
		t.Run(id, func(t *testing.T) {
			uc := &consultorFalso{}

			rec := getProyecto(t, uc, id)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			var got respuestaError
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decodificar respuesta: %v", err)
			}
			if want := "el id del proyecto debe ser un número entero positivo"; got.Error != want {
				t.Errorf("error = %q, se esperaba %q", got.Error, want)
			}
			if uc.llamado {
				t.Errorf("no se esperaba llamar al caso de uso con un id inválido")
			}
		})
	}
}

// Traza: HU-04 sección 7 (falla interna); no tiene un CA propio en la spec.
func TestGetProyecto_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	uc := &consultorFalso{err: fmt.Errorf("consultar proyecto 7: %w", errors.New("disco lleno"))}

	rec := getProyecto(t, uc, "7")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500; cuerpo: %s", rec.Code, rec.Body)
	}
	cuerpo := rec.Body.String() // se lee una vez: decodificar consume el buffer
	var got respuestaError
	if err := json.Unmarshal([]byte(cuerpo), &got); err != nil {
		t.Fatalf("decodificar respuesta: %v; cuerpo: %q", err, cuerpo)
	}
	if want := "error interno al consultar el proyecto"; got.Error != want {
		t.Errorf("error = %q, se esperaba %q", got.Error, want)
	}
	if strings.Contains(cuerpo, "disco lleno") || strings.Contains(cuerpo, "consultar proyecto 7") {
		t.Errorf("la respuesta no debería exponer el error interno: %s", cuerpo)
	}
}
