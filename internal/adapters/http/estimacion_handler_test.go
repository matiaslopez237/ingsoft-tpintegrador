package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// registradorEstimacionFalso devuelve el resumen o el error indicados y guarda lo que recibió.
type registradorEstimacionFalso struct {
	res      domain.ResumenRonda
	err      error
	recibido *app.DatosEstimacion
}

func (f *registradorEstimacionFalso) Ejecutar(_ context.Context, d app.DatosEstimacion) (domain.ResumenRonda, error) {
	f.recibido = &d
	return f.res, f.err
}

func pedirEstimacion(t *testing.T, uc *registradorEstimacionFalso, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	RegistrarEstimacionAcordada(mux, uc)
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return rec
}

// Traza: HU-17 CA-17.9, sección 3: 200 con el resumen de la ronda, que suma estimacion_acordada.
func TestPostEstimacion_Valida_Devuelve200ConLaEstimacion(t *testing.T) {
	valor := 8
	uc := &registradorEstimacionFalso{res: domain.ResumenRonda{ID: 3, ProyectoID: 1, ElementoID: 7, Numero: 1,
		Estado: domain.RondaCerrada, Votantes: []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}}, EstimacionAcordada: &valor}}

	rec := pedirEstimacion(t, uc, "/proyectos/1/rondas/3/estimacion", `{"valor":8}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	if uc.recibido == nil || *uc.recibido != (app.DatosEstimacion{ProyectoID: 1, RondaID: 3, Valor: 8}) {
		t.Errorf("recibido = %+v", uc.recibido)
	}
	mismoJSON(t, rec.Body.String(), `{"id":3,"proyecto_id":1,"elemento_id":7,"numero":1,"estado":"Cerrada",
		"votantes":[{"nombre":"Ana","voto":true}],"estimacion_acordada":8}`)
}

// Traza: HU-17 CA-17.9: una ronda sin estimación no trae el campo, así no cambia el JSON de HU-14.
func TestRespuestaRonda_SinEstimacion_NoTraeElCampo(t *testing.T) {
	rec := pedirEstimacion(t, &registradorEstimacionFalso{res: domain.ResumenRonda{ID: 3}}, "/proyectos/1/rondas/3/estimacion", `{"valor":8}`)

	if strings.Contains(rec.Body.String(), "estimacion_acordada") {
		t.Errorf("cuerpo = %s: sin estimación no debería traer el campo", rec.Body)
	}
}

// Traza: HU-17 CA-17.9, sección 7: 404 y 400 con el mensaje del dominio, aunque llegue envuelto, y 500
// sin el detalle.
func TestPostEstimacion_ErroresDelCasoDeUso(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{domain.ErrProyectoNoEncontrado, http.StatusNotFound},
		{domain.ErrRondaNoEncontrada, http.StatusNotFound},
		{domain.ErrEstimacionRondaAbierta, http.StatusBadRequest},
		{domain.ErrRondaNoEsLaUltima, http.StatusBadRequest},
		{domain.ErrElementoTerminado, http.StatusBadRequest},
		{domain.ErrEstimacionYaRegistrada, http.StatusBadRequest},
		{domain.ErrEstimacionInvalida, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			uc := &registradorEstimacionFalso{err: fmt.Errorf("guardar estimación de la ronda 3: %w", tt.err)}

			rec := pedirEstimacion(t, uc, "/proyectos/1/rondas/3/estimacion", `{"valor":8}`)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d", rec.Code, tt.want)
			}
			if cuerpo := rec.Body.String(); !strings.Contains(cuerpo, tt.err.Error()) || strings.Contains(cuerpo, "guardar estimación") {
				t.Errorf("cuerpo = %s: se esperaba solo el mensaje del dominio", cuerpo)
			}
		})
	}

	t.Run("falla interna", func(t *testing.T) {
		rec := pedirEstimacion(t, &registradorEstimacionFalso{err: errors.New("disco lleno")}, "/proyectos/1/rondas/3/estimacion", `{"valor":8}`)

		if rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), "disco lleno") {
			t.Errorf("status = %d, cuerpo = %s: se esperaba 500 sin el detalle", rec.Code, rec.Body)
		}
	})
}

// Traza: HU-17 sección 7: un pedido mal formado da 400 sin llamar al caso de uso.
func TestPostEstimacion_PedidoInvalido_NoLlamaAlCasoDeUso(t *testing.T) {
	tests := []struct{ nombre, ruta, cuerpo string }{
		{"JSON mal formado", "/proyectos/1/rondas/3/estimacion", `{"valor":`},
		{"campo desconocido", "/proyectos/1/rondas/3/estimacion", `{"valor":8,"otro":1}`},
		{"valor que no es un número", "/proyectos/1/rondas/3/estimacion", `{"valor":"ocho"}`},
		{"proyecto inválido", "/proyectos/abc/rondas/3/estimacion", `{"valor":8}`},
		{"ronda inválida", "/proyectos/1/rondas/abc/estimacion", `{"valor":8}`},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorEstimacionFalso{}

			rec := pedirEstimacion(t, uc, tt.ruta, tt.cuerpo)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, se esperaba 400", rec.Code)
			}
			if uc.recibido != nil {
				t.Error("no se tendría que llamar al caso de uso")
			}
		})
	}
}
