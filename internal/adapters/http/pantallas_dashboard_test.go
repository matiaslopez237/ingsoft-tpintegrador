package http

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const rutaDashboardDeAlfa = "/web/proyectos/7/dashboard"

// casosConDashboard arma los casos de uso de las pantallas con un Dashboard fijo para el proyecto Alfa.
func casosConDashboard(d app.Dashboard, err error) (CasosDeUso, *consultorDashboardFalso) {
	cu := casosFalsos()
	f := &consultorDashboardFalso{dashboard: d, err: err}
	cu.ConsultarDashboard = f
	return cu, f
}

// Traza: HU-25 CA-25.9 (ADR-0005): al menos 3 gráficos SVG, cada uno con título accesible y una tabla
// con sus datos, sin scripts ni estilos en línea (la CSP los bloquearía).
func TestPantallaDashboard_MuestraLosGraficos(t *testing.T) {
	cu, f := casosConDashboard(dashboardCompleto(), nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaDashboardDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if f.recibido != 7 {
		t.Errorf("proyecto recibido = %d, se esperaba 7", f.recibido)
	}
	cuerpo := rec.Body.String()
	if n := strings.Count(cuerpo, `<svg`); n != 4 {
		t.Errorf("hay %d gráficos SVG, se esperaban 4 (puntos, porcentaje, defectos y horas)", n)
	}
	if n := strings.Count(cuerpo, `role="img"`); n != 4 {
		t.Errorf("hay %d gráficos con role=\"img\", se esperaban 4", n)
	}
	contiene(t, rec,
		`aria-labelledby="grafico-puntos-titulo"`, `<title id="grafico-puntos-titulo">Puntos por sprint</title>`,
		`<title id="grafico-porcentaje-titulo">`, `<title id="grafico-defectos-titulo">`, `<title id="grafico-horas-titulo">`,
		`<caption>Puntos por sprint</caption>`, `<th scope="row">Sprint 1</th>`,
		"Velocidad: 3", "75%", "sin datos",
		"Desviación: 20,0 %",
	)
	noContiene(t, cuerpo, "<script")
	noContiene(t, cuerpo, `style="`)
	if csp := rec.Header().Get("Content-Security-Policy"); csp != "default-src 'self'" {
		t.Errorf("Content-Security-Policy = %q, se esperaba default-src 'self'", csp)
	}
}

// Traza: HU-25 CA-25.1, CA-25.2: el estado del proyecto muestra el backlog por estado y el sprint en curso.
func TestPantallaDashboard_MuestraElEstadoDelProyecto(t *testing.T) {
	cu, _ := casosConDashboard(dashboardCompleto(), nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaDashboardDeAlfa, nil)

	contiene(t, rec,
		"Dashboard de Alfa",
		"<dt>Pendientes</dt><dd>2</dd>", "<dt>En progreso</dt><dd>1</dd>", "<dt>Hechos</dt><dd>3</dd>", "<dt>Total</dt><dd>6</dd>",
		"Sprint 3", "Cerrar el login", "2026-11-01", "2026-11-14",
		`<progress id="avance-sprint" max="8" value="5">`, "5 de 8 puntos completados",
		`href="/web/proyectos/7/sprints/30"`,
	)
}

// Traza: HU-25 CA-25.3: sin sprint en curso lo dice, sin barra de avance.
func TestPantallaDashboard_SinSprintEnCurso(t *testing.T) {
	d := dashboardCompleto()
	d.HaySprintEnCurso, d.SprintEnCurso = false, domain.SprintEnCurso{}
	cu, _ := casosConDashboard(d, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaDashboardDeAlfa, nil)

	contiene(t, rec, "No hay sprint en curso")
	noContiene(t, rec.Body.String(), "<progress")
}

// Traza: HU-25 CA-25.11: un proyecto sin datos responde 200, con avisos en lugar de los gráficos.
func TestPantallaDashboard_ProyectoSinDatos(t *testing.T) {
	cu, _ := casosConDashboard(app.Dashboard{}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaDashboardDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "<dt>Total</dt><dd>0</dd>", "No hay sprint en curso",
		"Todavía no hay sprints cerrados", "Todavía no hay horas estimadas")
	noContiene(t, rec.Body.String(), "<svg")
}

// Traza: HU-25 CA-25.7: un proyecto inexistente muestra la página de error con 404.
func TestPantallaDashboard_ProyectoInexistente_Devuelve404(t *testing.T) {
	cu, _ := casosConDashboard(app.Dashboard{}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/99/dashboard", nil)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404", rec.Code)
	}
	contiene(t, rec, "proyecto no encontrado")
}

// Traza: HU-25, sección 7: una falla del caso de uso da 500 sin el detalle.
func TestPantallaDashboard_FallaInterna_Devuelve500(t *testing.T) {
	cu, _ := casosConDashboard(app.Dashboard{}, errors.New("disco lleno"))

	rec := pedirPantalla(t, cu, http.MethodGet, rutaDashboardDeAlfa, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	noContiene(t, rec.Body.String(), "disco lleno")
}

// Traza: HU-25 CA-25.10 (RN-25.6): la pantalla del proyecto tiene el enlace al Dashboard, el primero.
func TestPantallaProyecto_EnlaceAlDashboard(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/7", nil)

	cuerpo := rec.Body.String()
	dashboard := strings.Index(cuerpo, `href="/web/proyectos/7/dashboard"`)
	if dashboard < 0 {
		t.Fatalf("la pantalla del proyecto no tiene el enlace al Dashboard; cuerpo: %.400s", cuerpo)
	}
	if backlog := strings.Index(cuerpo, `href="/web/proyectos/7/backlog"`); backlog < dashboard {
		t.Errorf("el enlace al Dashboard tiene que ser el primero de la fila, antes que el del backlog")
	}
}
