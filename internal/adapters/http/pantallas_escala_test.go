package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// opcionesEscala son las opciones de la lista de story points (HU-13, sección 3).
var opcionesEscala = []string{
	`<select id="story_points" name="story_points">`, `<option value="">Sin estimar</option>`,
	`<option value="1">1</option>`, `<option value="2">2</option>`, `<option value="3">3</option>`,
	`<option value="8">8</option>`, `<option value="13">13</option>`, `<option value="21">21</option>`,
}

// Traza: HU-13 CA-13.7: en el alta, los story points se eligen de la escala.
func TestPantallaBacklog_AltaConLaEscala(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet, "/web/proyectos/7/backlog", nil)

	contiene(t, rec, opcionesEscala...)
	contiene(t, rec, `<option value="5">5</option>`)
	noContiene(t, rec.Body.String(), `name="story_points" type="number"`)
}

// Traza: HU-13 CA-13.7: en la edición, la lista marca el valor actual.
func TestPantallaElemento_EdicionConLaEscala(t *testing.T) {
	rec := pedirPantalla(t, casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{}), http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	contiene(t, rec, opcionesEscala...)
	contiene(t, rec, `<option value="5" selected>5</option>`)
	noContiene(t, rec.Body.String(), `type="number"`)
}

// Traza: HU-13 CA-13.7 (RN-13.3): un valor anterior a la escala aparece seleccionado y marcado,
// así guardar sin tocarlo no lo cambia.
func TestPantallaElemento_ValorAnteriorALaEscala(t *testing.T) {
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})
	cu.ConsultarElementoBacklog = consultarElementoFunc(func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
		e := elementoLogin()
		cuatro := 4
		e.ProyectoID, e.StoryPoints = 7, &cuatro
		return e, nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	contiene(t, rec, `<option value="4" selected>4 (fuera de la escala)</option>`, `name="original_story_points" value="4"`)
}

// Traza: HU-13 CA-13.7: un elemento sin estimar tiene "Sin estimar" seleccionado.
func TestOpcionesStoryPoints(t *testing.T) {
	tests := []struct {
		actual   string
		elegida  string
		cantidad int
	}{
		{"", "Sin estimar", 8},
		{"13", "13", 8},
		{"4", "4 (fuera de la escala)", 9},
	}
	for _, tt := range tests {
		ops := opcionesStoryPoints(tt.actual)
		if len(ops) != tt.cantidad {
			t.Errorf("%q: %d opciones, se esperaban %d", tt.actual, len(ops), tt.cantidad)
		}
		elegidas := 0
		for _, o := range ops {
			if o.Elegida {
				elegidas++
				if o.Etiqueta != tt.elegida {
					t.Errorf("%q: elegida %q, se esperaba %q", tt.actual, o.Etiqueta, tt.elegida)
				}
			}
		}
		if elegidas != 1 {
			t.Errorf("%q: %d opciones elegidas, se esperaba 1", tt.actual, elegidas)
		}
	}
}
