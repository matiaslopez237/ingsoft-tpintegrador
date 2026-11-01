package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// casosPokerPantalla son los casos de uso de la pantalla del elemento Login (proyecto 7,
// elemento 7) con el Planning Poker de HU-14 sobre el falso f.
func casosPokerPantalla(f *pokerFalso) CasosDeUso {
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})
	cu.IniciarRonda, cu.Votar, cu.ConsultarRonda = f, votadorPokerFalso{f}, consultorPokerFalso{f}
	return cu
}

// conRonda devuelve el falso con la ronda 1 abierta sobre Login y los votos dados (Ana, Beto…).
func conRonda(t *testing.T, votos ...int) *pokerFalso {
	t.Helper()
	f := nuevoPokerFalso()
	if _, err := f.Ejecutar(context.Background(), app.DatosIniciarRonda{ProyectoID: 7, ElementoID: 7}); err != nil {
		t.Fatal(err)
	}
	for i, v := range votos {
		nombre := []string{"Ana", "Beto"}[i]
		if _, err := (votadorPokerFalso{f}).Ejecutar(context.Background(), app.DatosVoto{ProyectoID: 7, RondaID: 1, Integrante: nombre, Valor: v}); err != nil {
			t.Fatal(err)
		}
	}
	f.votos = 0
	return f
}

const rutaRonda = "/web/proyectos/7/rondas/1"

// Traza: HU-14 CA-14.9: sin ronda abierta, la pantalla del elemento ofrece iniciar una.
func TestPantallaElemento_SinRonda_OfreceIniciar(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(nuevoPokerFalso()), http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Planning Poker", `action="/web/proyectos/7/backlog/7/rondas"`, "Iniciar ronda")
}

// Traza: HU-14 CA-14.9 (RN-14.1): con una ronda abierta, se enlaza a ella en lugar de ofrecer otra.
func TestPantallaElemento_ConRondaAbierta_EnlazaALaRonda(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(conRonda(t)), http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	contiene(t, rec, `href="`+rutaRonda+`"`)
	noContiene(t, rec.Body.String(), "Iniciar ronda")
}

// Traza: HU-14 CA-14.9 (RN-14.1): con la última ronda cerrada se puede iniciar otra.
func TestPantallaElemento_ConRondaCerrada_OfreceIniciarOtra(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(conRonda(t, 5, 8)), http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	contiene(t, rec, "Iniciar ronda")
}

// Traza: HU-14 CA-14.9 (RN-14.1): un elemento Hecho no se estima.
func TestPantallaElemento_Hecho_NoOfreceIniciar(t *testing.T) {
	cu := casosPokerPantalla(nuevoPokerFalso())
	cu.ConsultarElementoBacklog = consultarElementoFunc(func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
		e := elementoLogin()
		e.ProyectoID, e.Estado = 7, domain.EstadoHecho
		return e, nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	noContiene(t, rec.Body.String(), "Iniciar ronda")
}

// Traza: HU-14 CA-14.9: al iniciar se va a la pantalla de la ronda (POST-redirect-GET).
func TestPantallaIniciarRonda_RedirigeALaRonda(t *testing.T) {
	f := nuevoPokerFalso()

	rec := pedirPantalla(t, casosPokerPantalla(f), http.MethodPost, "/web/proyectos/7/backlog/7/rondas", url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaRonda {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a %s", rec.Code, rec.Header().Get("Location"), rutaRonda)
	}
	if f.iniciadas != 1 {
		t.Errorf("rondas iniciadas = %d", f.iniciadas)
	}
}

// Traza: HU-14 CA-14.9 (CA-14.6): el error se muestra en la pantalla del elemento.
func TestPantallaIniciarRonda_Error_SeMuestraEnElElemento(t *testing.T) {
	f := nuevoPokerFalso()
	f.err = domain.ErrRondaSinVotantes

	rec := pedirPantalla(t, casosPokerPantalla(f), http.MethodPost, "/web/proyectos/7/backlog/7/rondas", url.Values{})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, domain.ErrRondaSinVotantes.Error(), `name="titulo"`)
}

// Traza: HU-14 CA-14.9 (RN-14.3): se ve quién votó y quién falta, nunca el valor.
func TestPantallaRonda_MuestraQuienVotoSinLosValores(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(conRonda(t, 13)), http.MethodGet, rutaRonda, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Login", "Ronda 1", "<td>Ana</td>", "<td>Votó</td>", "<td>Beto</td>", "<td>Falta votar</td>",
		`action="`+rutaRonda+`/votar"`, `name="integrante"`, `name="valor"`)
	noContiene(t, rec.Body.String(), "selected>13")
}

// Traza: HU-14 CA-14.9 (RN-14.6): con la ronda cerrada no hay formulario ni valores.
func TestPantallaRonda_Cerrada_SinFormulario(t *testing.T) {
	rec := pedirPantalla(t, casosPokerPantalla(conRonda(t, 13, 8)), http.MethodGet, rutaRonda, nil)

	contiene(t, rec, "Votaron todos")
	cuerpo := rec.Body.String()
	for _, texto := range []string{`name="valor"`, ">13<", ">8<"} {
		noContiene(t, cuerpo, texto)
	}
}

// Traza: HU-14 CA-14.9: votar vuelve a la pantalla de la ronda.
func TestPantallaVotar_GuardaYVuelveALaRonda(t *testing.T) {
	f := conRonda(t)

	rec := pedirPantalla(t, casosPokerPantalla(f), http.MethodPost, rutaRonda+"/votar", url.Values{"integrante": {"Ana"}, "valor": {"5"}})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaRonda {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a %s", rec.Code, rec.Header().Get("Location"), rutaRonda)
	}
	if f.votos != 1 {
		t.Errorf("votos guardados = %d", f.votos)
	}
}

// Traza: HU-14 CA-14.9 (CA-14.5): los errores se muestran en la misma pantalla de la ronda.
func TestPantallaVotar_Errores(t *testing.T) {
	tests := []struct {
		nombre   string
		form     url.Values
		wantCode int
		wantMsg  string
	}{
		{"sin valor", url.Values{"integrante": {"Ana"}}, http.StatusBadRequest, "elegí un valor"},
		{"valor que no es número", url.Values{"integrante": {"Ana"}, "valor": {"mucho"}}, http.StatusBadRequest, "elegí un valor"},
		{"valor fuera de la escala", url.Values{"integrante": {"Ana"}, "valor": {"4"}}, http.StatusBadRequest, domain.ErrVotoInvalido.Error()},
		{"no es votante", url.Values{"integrante": {"Dani"}, "valor": {"5"}}, http.StatusBadRequest, domain.ErrNoEsVotante.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := conRonda(t)
			rec := pedirPantalla(t, casosPokerPantalla(f), http.MethodPost, rutaRonda+"/votar", tt.form)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantCode, rec.Body)
			}
			contiene(t, rec, tt.wantMsg, "Ronda 1")
			if f.votos != 0 {
				t.Errorf("no se esperaba guardar votos: %d", f.votos)
			}
		})
	}
}

// Traza: HU-14 CA-14.7: una ronda inexistente o de otro proyecto da 404.
func TestPantallaRonda_Inexistente(t *testing.T) {
	cu := casosPokerPantalla(conRonda(t))
	for _, ruta := range []string{"/web/proyectos/7/rondas/9", "/web/proyectos/8/rondas/1"} {
		if rec := pedirPantalla(t, cu, http.MethodGet, ruta, nil); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, se esperaba 404", ruta, rec.Code)
		}
	}
	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/rondas/9/votar", url.Values{"integrante": {"Ana"}, "valor": {"5"}})
	if rec.Code != http.StatusNotFound {
		t.Errorf("votar en una ronda inexistente: status = %d, se esperaba 404", rec.Code)
	}
}
