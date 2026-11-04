package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorConOtraUltima responde la ronda del falso de HU-14, pero informa como última del
// elemento otra ronda (ID 2): así la ronda 1 deja de ser la última.
type consultorConOtraUltima struct{ consultorPokerFalso }

func (f consultorConOtraUltima) UltimaDelElemento(context.Context, int64, int64) (*domain.ResumenRonda, error) {
	return &domain.ResumenRonda{ID: 2, ProyectoID: 7, ElementoID: 7, Numero: 2, Estado: domain.RondaAbierta}, nil
}

const accionNuevaRonda = `action="/web/proyectos/7/backlog/7/rondas"`

// Traza: HU-16 CA-16.8, RN-16.4: la última ronda, cerrada, ofrece "Nueva ronda", aunque no haya habido
// diferencias; la abierta no.
func TestPantallaRonda_NuevaRonda(t *testing.T) {
	tests := []struct {
		name   string
		votos  []int
		ofrece bool
	}{
		{"cerrada, con diferencia", []int{3, 13}, true},
		{"cerrada, sin diferencia", []int{5, 5}, true},
		{"abierta", []int{5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := pedirPantalla(t, casosPokerPantalla(conRonda(t, tt.votos...)), http.MethodGet, rutaRonda, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			if got := strings.Contains(rec.Body.String(), accionNuevaRonda) && strings.Contains(rec.Body.String(), "Nueva ronda"); got != tt.ofrece {
				t.Errorf("ofrece \"Nueva ronda\" = %v, se esperaba %v; cuerpo: %.600s", got, tt.ofrece, rec.Body)
			}
		})
	}
}

// Traza: HU-16 RN-16.4: una ronda cerrada que no es la última del elemento no ofrece "Nueva ronda".
func TestPantallaRonda_CerradaPeroNoEsLaUltima_SinNuevaRonda(t *testing.T) {
	f := conRonda(t, 3, 13)
	cu := casosPokerPantalla(f)
	cu.ConsultarRonda = consultorConOtraUltima{consultorPokerFalso{f}}

	rec := pedirPantalla(t, cu, http.MethodGet, rutaRonda, nil)

	noContiene(t, rec.Body.String(), "Nueva ronda")
}

// Traza: HU-16 RN-16.4: con el elemento Hecho no se ofrece "Nueva ronda".
func TestPantallaRonda_ElementoHecho_SinNuevaRonda(t *testing.T) {
	cu := casosPokerPantalla(conRonda(t, 3, 13))
	cu.ConsultarElementoBacklog = consultarElementoFunc(func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
		e := elementoLogin()
		e.ProyectoID, e.Estado = 7, domain.EstadoHecho
		return e, nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, rutaRonda, nil)

	noContiene(t, rec.Body.String(), "Nueva ronda")
}

// historialDePrueba son la ronda 2, abierta, y la 1, cerrada de 3 a 13, del elemento 7 del proyecto 7.
var historialDePrueba = []domain.ResumenRonda{
	{ID: 12, ProyectoID: 7, ElementoID: 7, Numero: 2, Estado: domain.RondaAbierta,
		Votantes: []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto"}}},
	{ID: 11, ProyectoID: 7, ElementoID: 7, Numero: 1, Estado: domain.RondaCerrada,
		Votantes: []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto", YaVoto: true}},
		Revelacion: &domain.Revelacion{Votos: []domain.VotoRevelado{{Nombre: "Ana", Valor: 3}, {Nombre: "Beto", Valor: 13}},
			Minimo: 3, Maximo: 13, DiferenciaSignificativa: true}},
}

// Traza: HU-16 CA-16.8, RN-16.3: la pantalla del elemento lista las rondas en el orden del caso de
// uso, con un enlace a cada una; de la cerrada, el mínimo, el máximo y la diferencia.
func TestPantallaElemento_HistorialDeRondas(t *testing.T) {
	cu := casosPokerPantalla(nuevoPokerFalso())
	cu.ListarRondas = &historialRondasFalso{historial: historialDePrueba}

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Historial de rondas", "Ronda 2", "Ronda 1", "de 3 a 13", "diferencia significativa")
	cuerpo := rec.Body.String()
	i2, i1 := strings.Index(cuerpo, `href="/web/proyectos/7/rondas/12"`), strings.Index(cuerpo, `href="/web/proyectos/7/rondas/11"`)
	if i2 < 0 || i1 < 0 || i2 > i1 {
		t.Errorf("se esperaban los enlaces a la ronda 2 y después a la 1; posiciones %d y %d", i2, i1)
	}
}

// Traza: HU-16 CA-16.5: sin rondas, la pantalla del elemento no muestra el historial.
func TestPantallaElemento_SinRondas_SinHistorial(t *testing.T) {
	cu := casosPokerPantalla(nuevoPokerFalso())
	cu.ListarRondas = &historialRondasFalso{}

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	noContiene(t, rec.Body.String(), "Historial de rondas")
}

// Traza: HU-16 sección 7: una falla al leer el historial es una falla interna, sin el detalle.
func TestPantallaElemento_FallaDelHistorial(t *testing.T) {
	cu := casosPokerPantalla(nuevoPokerFalso())
	cu.ListarRondas = &historialRondasFalso{err: errors.New("disco lleno")}

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, se esperaba 500", rec.Code)
	}
	noContiene(t, rec.Body.String(), "disco lleno")
}
