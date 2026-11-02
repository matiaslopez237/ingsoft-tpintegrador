package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const accionEstimacion = rutaRonda + "/estimacion"

// Traza: HU-17 CA-17.10 (RN-17.3): la última ronda, cerrada y sin estimación, ofrece el formulario
// con la escala.
func TestPantallaRonda_CerradaSinEstimacion_OfreceRegistrar(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(conRonda(t, 3, 13)), http.MethodGet, rutaRonda, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, `action="`+accionEstimacion+`"`, `name="valor"`, `<option value="21">`, "Registrar estimación")
}

// Traza: HU-17 CA-17.10 (RN-17.3, RN-17.4): en los demás casos no se ofrece.
func TestPantallaRonda_NoOfreceRegistrar(t *testing.T) {
	tests := []struct {
		name  string
		casos func(*testing.T) CasosDeUso
	}{
		{"abierta", func(t *testing.T) CasosDeUso { return casosPokerPantalla(conRonda(t, 5)) }},
		{"no es la última", func(t *testing.T) CasosDeUso {
			f := conRonda(t, 3, 13)
			cu := casosPokerPantalla(f)
			cu.ConsultarRonda = consultorConOtraUltima{consultorPokerFalso{f}}
			return cu
		}},
		{"elemento Hecho", func(t *testing.T) CasosDeUso {
			cu := casosPokerPantalla(conRonda(t, 3, 13))
			cu.ConsultarElementoBacklog = consultarElementoFunc(func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
				e := elementoLogin()
				e.ProyectoID, e.Estado = 7, domain.EstadoHecho
				return e, nil
			})
			return cu
		}},
		{"ya tiene estimación", func(t *testing.T) CasosDeUso {
			f := conRonda(t, 3, 13)
			v := 5
			f.ronda.EstimacionAcordada = &v
			return casosPokerPantalla(f)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := pedirPantalla(t, tt.casos(t), http.MethodGet, rutaRonda, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			noContiene(t, rec.Body.String(), accionEstimacion)
		})
	}
}

// Traza: HU-17 sección 3: la ronda con estimación la muestra.
func TestPantallaRonda_ConEstimacion_LaMuestra(t *testing.T) {
	f := conRonda(t, 3, 13)
	v := 8
	f.ronda.EstimacionAcordada = &v

	rec := pedirPantalla(t, casosPokerPantalla(f), http.MethodGet, rutaRonda, nil)

	contiene(t, rec, "Estimación acordada: 8")
}

// Traza: HU-17 CA-17.10: al registrar se vuelve a la ronda (POST-redirect-GET).
func TestPantallaEstimacion_Valida_RedirigeALaRonda(t *testing.T) {
	cu := casosPokerPantalla(conRonda(t, 3, 13))
	uc := &registradorEstimacionFalso{}
	cu.RegistrarEstimacion = uc

	rec := pedirPantalla(t, cu, http.MethodPost, accionEstimacion, url.Values{"valor": {"8"}})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaRonda {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a %s", rec.Code, rec.Header().Get("Location"), rutaRonda)
	}
	if uc.recibido == nil || *uc.recibido != (app.DatosEstimacion{ProyectoID: 7, RondaID: 1, Valor: 8}) {
		t.Errorf("recibido = %+v", uc.recibido)
	}
}

// Traza: HU-17 CA-17.10, sección 7: los errores de la regla se muestran en la misma pantalla, con el
// mensaje del dominio aunque llegue envuelto.
func TestPantallaEstimacion_ErroresDeValidacion(t *testing.T) {
	for _, e := range []error{domain.ErrEstimacionRondaAbierta, domain.ErrRondaNoEsLaUltima, domain.ErrElementoTerminado,
		domain.ErrEstimacionYaRegistrada, domain.ErrEstimacionInvalida} {
		t.Run(e.Error(), func(t *testing.T) {
			cu := casosPokerPantalla(conRonda(t, 3, 13))
			cu.RegistrarEstimacion = &registradorEstimacionFalso{err: fmt.Errorf("guardar estimación de la ronda 1: %w", e)}

			rec := pedirPantalla(t, cu, http.MethodPost, accionEstimacion, url.Values{"valor": {"8"}})

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.300s", rec.Code, rec.Body)
			}
			contiene(t, rec, "Planning Poker", e.Error())
			noContiene(t, rec.Body.String(), "guardar estimación")
		})
	}
}

// Traza: HU-17 sección 7: sin un valor numérico se avisa sin llamar al caso de uso.
func TestPantallaEstimacion_ValorNoNumerico_NoLlamaAlCasoDeUso(t *testing.T) {
	cu := casosPokerPantalla(conRonda(t, 3, 13))
	uc := &registradorEstimacionFalso{}
	cu.RegistrarEstimacion = uc

	rec := pedirPantalla(t, cu, http.MethodPost, accionEstimacion, url.Values{"valor": {""}})

	if rec.Code != http.StatusBadRequest || uc.recibido != nil {
		t.Errorf("status = %d, recibido = %+v; se esperaba 400 sin llamar al caso de uso", rec.Code, uc.recibido)
	}
	contiene(t, rec, "elegí un valor de la escala")
}

// Traza: HU-17 sección 7: lo inexistente da 404 y una falla interna, 500 sin el detalle.
func TestPantallaEstimacion_OtrosErrores(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
		want   int
	}{
		{"ronda inexistente", domain.ErrRondaNoEncontrada, http.StatusNotFound},
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, http.StatusNotFound},
		{"falla interna", errors.New("disco lleno"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu := casosPokerPantalla(conRonda(t, 3, 13))
			cu.RegistrarEstimacion = &registradorEstimacionFalso{err: tt.err}

			rec := pedirPantalla(t, cu, http.MethodPost, accionEstimacion, url.Values{"valor": {"8"}})

			if rec.Code != tt.want {
				t.Errorf("status = %d, se esperaba %d", rec.Code, tt.want)
			}
			noContiene(t, rec.Body.String(), "disco lleno")
		})
	}
}

// Traza: HU-17 sección 5: un envío desde otro sitio se rechaza, sin llamar al caso de uso.
func TestPantallaEstimacion_DesdeOtroSitio_SeRechaza(t *testing.T) {
	cu := casosPokerPantalla(conRonda(t, 3, 13))
	uc := &registradorEstimacionFalso{}
	cu.RegistrarEstimacion = uc
	mux := http.NewServeMux()
	RegistrarPantallas(mux, cu)
	req := httptest.NewRequest(http.MethodPost, accionEstimacion, strings.NewReader("valor=8"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden || uc.recibido != nil {
		t.Errorf("status = %d, recibido = %+v; se esperaba 403 sin llamar al caso de uso", rec.Code, uc.recibido)
	}
}

// Traza: HU-17 CA-17.10: el historial del elemento muestra lo acordado en cada ronda.
func TestPantallaElemento_HistorialMuestraLaEstimacion(t *testing.T) {
	v := 8
	con := domain.ResumenRonda{ID: 11, ProyectoID: 7, ElementoID: 7, Numero: 1, Estado: domain.RondaCerrada, EstimacionAcordada: &v}
	sin := domain.ResumenRonda{ID: 12, ProyectoID: 7, ElementoID: 7, Numero: 2, Estado: domain.RondaCerrada}
	cu := casosPokerPantalla(nuevoPokerFalso())
	cu.ListarRondas = &historialRondasFalso{historial: []domain.ResumenRonda{sin, con}}

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	contiene(t, rec, "acordada: 8")
	if n := strings.Count(rec.Body.String(), "acordada:"); n != 1 {
		t.Errorf("aparece %d veces \"acordada:\", se esperaba solo en la ronda que la tiene", n)
	}
}
