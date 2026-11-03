package http

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const (
	rutaProyectoWeb      = "/web/proyectos/7"
	textoSprintCerrado   = "Solo cuentan los elementos que siguen en este sprint."
	textoSinElementos    = "Todavía no hay elementos."
	textoRealesSinEstim  = "sus 5 h están en las horas reales, pero no en la desviación."
	mensajeFallaInternaW = "Ocurrió un error interno"
)

// pantallasConHoras son las dos pantallas que muestran el bloque, con la sección que tiene que venir
// después de él (spec, sección 3).
var pantallasConHoras = []struct{ nombre, ruta, seccionSiguiente string }{
	{"sprint", rutaSprintDetalle, `id="titulo-asignados"`},
	{"proyecto", rutaProyectoWeb, `id="titulo-integrantes"`},
}

// casosConHoras arma las dos pantallas con el total real de esos elementos.
func casosConHoras(elementos ...domain.HorasElemento) CasosDeUso {
	total := domain.TotalizarHoras(elementos)
	cu := casosConAsignacion(nil, nil)
	cu.ConsultarHorasProyecto = &consultorHorasProyectoFalso{total: total}
	cu.ConsultarHorasSprint = &consultorHorasSprintFalso{total: total}
	return cu
}

// casosConFallaDeHoras arma las dos pantallas con un error al consultar las horas.
func casosConFallaDeHoras(err error) CasosDeUso {
	cu := casosConAsignacion(nil, nil)
	cu.ConsultarHorasProyecto = &consultorHorasProyectoFalso{err: err}
	cu.ConsultarHorasSprint = &consultorHorasSprintFalso{err: err}
	return cu
}

// seccionHoras devuelve el HTML del bloque "Horas y desviación".
func seccionHoras(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	_, resto, ok := strings.Cut(rec.Body.String(), `id="total-horas"`)
	if !ok {
		t.Fatalf("la página no tiene el bloque de horas; status = %d, cuerpo: %.400s", rec.Code, rec.Body)
	}
	seccion, _, ok := strings.Cut(resto, "</section>")
	if !ok {
		t.Fatalf("el bloque de horas no cierra; cuerpo: %.400s", resto)
	}
	return seccion
}

// Los mismos valores en las dos pantallas: coma decimal, un decimal en el porcentaje y signo en las
// positivas, como HU-19. Sin ningún elemento estimado, un texto y no un 0.
// Traza: HU-23 CA-23.17, CA-23.18 (RN-23.4, RN-23.5, RN-23.6; sección 3)
func TestPantallasHoras_MuestranLosValores(t *testing.T) {
	const sinEstimarTexto = "Sin estimar"
	tests := []struct {
		nombre                             string
		elementos                          []domain.HorasElemento
		wantElem                           string
		wantEst, wantReal, wantHs, wantPor string
		wantAviso                          string
	}{
		{"de más", []domain.HorasElemento{estimado(8, 10), estimado(4, 3)}, "2", "12 h", "13 h", "+1 h", "+8,3 %", ""},
		{"de menos", []domain.HorasElemento{estimado(8, 6), estimado(4, 3)}, "2", "12 h", "9 h", "-3 h", "-25,0 %", ""},
		{"sin desviación: 0 y sin signo", []domain.HorasElemento{estimado(8, 10), estimado(8, 6)}, "2", "16 h", "16 h", "0 h", "0,0 %", ""},
		{"con media hora", []domain.HorasElemento{estimado(20, 20.5), estimado(20, 20)}, "2", "40 h", "40,5 h", "+0,5 h", "+1,3 %", ""},
		{"estimado sin esfuerzo", []domain.HorasElemento{estimado(8, 0)}, "1", "8 h", "0 h", "-8 h", "-100,0 %", ""},
		{"con un elemento sin estimar", []domain.HorasElemento{estimado(8, 10), sinEstimar(5)}, "2", "8 h", "15 h", "+2 h", "+25,0 %",
			"1 de 2 elementos sin estimar: " + textoRealesSinEstim},
		{"ningún elemento estimado", []domain.HorasElemento{sinEstimar(5), sinEstimar(3)}, "2", "0 h", "8 h", sinEstimarTexto, sinEstimarTexto,
			"2 de 2 elementos sin estimar: sus 8 h están en las horas reales, pero no en la desviación."},
		{"un solo elemento, sin estimar", []domain.HorasElemento{sinEstimar(5)}, "1", "0 h", "5 h", sinEstimarTexto, sinEstimarTexto,
			"1 de 1 elemento sin estimar: " + textoRealesSinEstim},
		{"sin estimar y sin esfuerzo", []domain.HorasElemento{sinEstimar(0)}, "1", "0 h", "0 h", sinEstimarTexto, sinEstimarTexto,
			"1 de 1 elemento sin estimar: sus 0 h están en las horas reales, pero no en la desviación."},
		{"sin elementos", nil, "0", "0 h", "0 h", sinEstimarTexto, sinEstimarTexto, textoSinElementos},
	}
	for _, p := range pantallasConHoras {
		for _, tt := range tests {
			t.Run(p.nombre+": "+tt.nombre, func(t *testing.T) {
				rec := pedirPantalla(t, casosConHoras(tt.elementos...), http.MethodGet, p.ruta, nil)

				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
				}
				s := seccionHoras(t, rec)
				seccionContiene(t, s, `<h2 id="titulo-total-horas">Horas y desviación</h2>`, valorDeComparacion("Elementos", tt.wantElem))
				bloqueMuestra(t, s, tt.wantEst, tt.wantReal, tt.wantHs, tt.wantPor)
				if tt.wantAviso != "" {
					seccionContiene(t, s, tt.wantAviso)
				} else {
					seccionNoContiene(t, s, "elementos sin estimar", "elemento sin estimar", textoSinElementos)
				}
			})
		}
	}
}

// Traza: HU-23 CA-23.17, CA-23.18 (sección 3): después de los datos y antes de la primera sección, y
// lo que la pantalla ya mostraba sigue ahí.
func TestPantallasHoras_ElBloqueVaAntesDeLaPrimeraSeccion(t *testing.T) {
	for _, p := range pantallasConHoras {
		t.Run(p.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConHoras(estimado(8, 10)), http.MethodGet, p.ruta, nil)

			cuerpo := rec.Body.String()
			datos, horas, siguiente := strings.Index(cuerpo, `class="datos"`), strings.Index(cuerpo, `id="total-horas"`), strings.Index(cuerpo, p.seccionSiguiente)
			if datos < 0 || horas < datos || siguiente < horas {
				t.Errorf("posiciones: datos %d, horas %d, sección siguiente %d; cuerpo: %.400s", datos, horas, siguiente, cuerpo)
			}
			contiene(t, rec, `aria-labelledby="titulo-total-horas"`)
		})
	}
}

// Traza: HU-23 CA-23.17, CA-23.18 (sección 3): lo que ya mostraban las pantallas no cambia.
func TestPantallasHoras_NoRompenLoQueYaMostraban(t *testing.T) {
	rec := pedirPantalla(t, casosConHoras(estimado(8, 10)), http.MethodGet, rutaSprintDetalle, nil)
	contiene(t, rec, "Gestionar el backlog", "Elementos del Sprint Backlog", "Login de usuarios", "Logout",
		`action="/web/proyectos/7/sprints/1/asignar"`, `action="/web/proyectos/7/sprints/1/quitar"`)

	rec = pedirPantalla(t, casosConHoras(estimado(8, 10)), http.MethodGet, rutaProyectoWeb, nil)
	contiene(t, rec, "Alfa", "Integrantes", "Ana", "Agregar integrante", `action="/web/proyectos/7/integrantes"`,
		`href="/web/proyectos/7/sprints"`)
}

// Traza: HU-23 CA-23.17, CA-23.18 (sección 3; P-6): el bloque es de solo lectura.
func TestPantallasHoras_ElBloqueEsDeSoloLectura(t *testing.T) {
	for _, p := range pantallasConHoras {
		t.Run(p.nombre, func(t *testing.T) {
			s := seccionHoras(t, pedirPantalla(t, casosConHoras(estimado(8, 10)), http.MethodGet, p.ruta, nil))

			seccionNoContiene(t, s, "<form", "<input", "<button", "<script", "<a ")
		})
	}
}

// Traza: HU-23 CA-23.17, CA-23.18 (P-6): solo en los dos detalles: ni en las listas de sprints ni en
// la pantalla del elemento, que tiene su propia comparación (HU-19).
func TestPantallasHoras_NoApareceEnOtrasPantallas(t *testing.T) {
	cu := casosConHoras(estimado(8, 10))
	cu.ListarSprintsAnteriores = &listadorSprintsFalso{}
	comparacion := casosConComparacion(comparacionFija(pHoras(8), 10))
	tests := []struct {
		nombre, ruta string
		cu           CasosDeUso
	}{
		{"lista de sprints", "/web/proyectos/7/sprints", cu},
		{"sprints anteriores", "/web/proyectos/7/sprints/anteriores", cu},
		{"elemento", rutaElementoWeb, comparacion},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, tt.cu, http.MethodGet, tt.ruta, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			if strings.Contains(rec.Body.String(), `id="total-horas"`) {
				t.Errorf("la página no debería tener el bloque de horas; cuerpo: %.400s", rec.Body)
			}
		})
	}
	bloqueMuestra(t, seccionComparacion(t, pedirPantalla(t, comparacion, http.MethodGet, rutaElementoWeb, nil)), "8 h", "10 h", "+2 h", "+25,0 %")
}

// La aclaración es solo para un sprint Cerrado: en uno Abierto y en el proyecto no hay limitación.
// Traza: HU-23 CA-23.17, CA-23.18 (RN-23.8)
func TestPantallasHoras_AclaracionDelSprintCerrado(t *testing.T) {
	cerrado := casosConHoras(estimado(8, 10))
	cerrado.ListarSprints = listarSprintsFunc(func(context.Context, int64) ([]domain.Sprint, error) {
		s := sprintDeAlfa
		s.Estado = domain.SprintCerrado
		return []domain.Sprint{s}, nil
	})
	seccionContiene(t, seccionHoras(t, pedirPantalla(t, cerrado, http.MethodGet, rutaSprintDetalle, nil)),
		textoSprintCerrado, "Los que volvieron al backlog al cerrarlo no suman acá, pero sí en el total del proyecto.")

	abierto := casosConHoras(estimado(8, 10))
	seccionNoContiene(t, seccionHoras(t, pedirPantalla(t, abierto, http.MethodGet, rutaSprintDetalle, nil)), textoSprintCerrado)
	seccionNoContiene(t, seccionHoras(t, pedirPantalla(t, cerrado, http.MethodGet, rutaProyectoWeb, nil)), textoSprintCerrado)
}

// Traza: HU-23 CA-23.17, CA-23.18 (sección 5: rutas y pantallas): cada pantalla consulta con sus IDs.
func TestPantallasHoras_ConsultanConLosIDsDeLaRuta(t *testing.T) {
	proyecto, sprint := &consultorHorasProyectoFalso{}, &consultorHorasSprintFalso{}
	cu := casosConAsignacion(nil, nil)
	cu.ConsultarHorasProyecto, cu.ConsultarHorasSprint = proyecto, sprint

	pedirPantalla(t, cu, http.MethodGet, rutaSprintDetalle, nil)
	if sprint.proyectoID != 7 || sprint.sprintID != 1 || proyecto.llamado {
		t.Errorf("detalle del sprint: sprint %+v, proyecto llamado: %t", sprint, proyecto.llamado)
	}

	pedirPantalla(t, cu, http.MethodGet, rutaProyectoWeb, nil)
	if proyecto.proyectoID != 7 {
		t.Errorf("detalle del proyecto: el caso de uso recibió el proyecto %d", proyecto.proyectoID)
	}
}

// Al volver con el error de un formulario, el bloque se sigue mostrando, y una falla interna al
// consultar las horas tiene prioridad sobre el mensaje de validación, como en HU-19.
// Traza: HU-23 CA-23.17, CA-23.18 (sección 3)
func TestPantallasHoras_AlVolverConError_SigueElBloque(t *testing.T) {
	tests := []struct {
		nombre, ruta string
		form         url.Values
		errAsignar   error
		wantMsg      string
	}{
		{"sprint", rutaSprintDetalle + "/asignar", url.Values{"elemento_id": {"11"}}, domain.ErrElementoYaAsignado, domain.ErrElementoYaAsignado.Error()},
		{"proyecto", rutaProyectoWeb + "/integrantes", url.Values{"nombre": {"Ana"}, "rol": {"Tester"}}, nil, domain.ErrIntegranteDuplicado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu := casosConHoras(estimado(8, 10), estimado(4, 3))
			if tt.errAsignar != nil {
				cu.AsignarHistoria = &asignarHistoriaFalso{err: tt.errAsignar}
			}

			rec := pedirPantalla(t, cu, http.MethodPost, tt.ruta, tt.form)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, tt.wantMsg)
			bloqueMuestra(t, seccionHoras(t, rec), "12 h", "13 h", "+1 h", "+8,3 %")
		})
		t.Run(tt.nombre+": falla interna", func(t *testing.T) {
			cu := casosConFallaDeHoras(errors.New("database is locked"))
			if tt.errAsignar != nil {
				cu.AsignarHistoria = &asignarHistoriaFalso{err: tt.errAsignar}
			}

			rec := pedirPantalla(t, cu, http.MethodPost, tt.ruta, tt.form)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, se esperaba 500; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, mensajeFallaInternaW)
		})
	}
}

// El mensaje es el del dominio, aunque el error llegue envuelto.
// Traza: HU-23 CA-23.14, CA-23.15, CA-23.17, CA-23.18 (RN-23.11; sección 7)
func TestPantallasHoras_NoEncontrados_Devuelven404(t *testing.T) {
	tests := []struct {
		nombre, ruta string
		err, want    error
	}{
		{"sprint: proyecto", rutaSprintDetalle, domain.ErrProyectoNoEncontrado, domain.ErrProyectoNoEncontrado},
		{"sprint: proyecto, envuelto", rutaSprintDetalle, fmt.Errorf("buscar registro: %w", domain.ErrProyectoNoEncontrado), domain.ErrProyectoNoEncontrado},
		{"sprint: sprint", rutaSprintDetalle, domain.ErrSprintNoEncontrado, domain.ErrSprintNoEncontrado},
		{"sprint: sprint, envuelto", rutaSprintDetalle, fmt.Errorf("buscar registro: %w", domain.ErrSprintNoEncontrado), domain.ErrSprintNoEncontrado},
		{"proyecto: proyecto", rutaProyectoWeb, domain.ErrProyectoNoEncontrado, domain.ErrProyectoNoEncontrado},
		{"proyecto: proyecto, envuelto", rutaProyectoWeb, fmt.Errorf("buscar registro: %w", domain.ErrProyectoNoEncontrado), domain.ErrProyectoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConFallaDeHoras(tt.err), http.MethodGet, tt.ruta, nil)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, tt.want.Error())
			if strings.Contains(rec.Body.String(), "buscar registro") {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", rec.Body)
			}
		})
	}
}

// Los que ya resuelve la pantalla no llegan a consultar las horas.
// Traza: HU-23 CA-23.14, CA-23.15, CA-23.17, CA-23.18 (RN-23.11; sección 7)
func TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras(t *testing.T) {
	const (
		idProyecto = "el id del proyecto debe ser un número entero positivo"
		idSprint   = "el id del sprint debe ser un número entero positivo"
	)
	tests := []struct {
		nombre, ruta string
		wantCode     int
		wantMsg      string
	}{
		{"proyecto: id no numérico", "/web/proyectos/abc", http.StatusBadRequest, idProyecto},
		{"proyecto: id en cero", "/web/proyectos/0", http.StatusBadRequest, idProyecto},
		{"proyecto: inexistente", "/web/proyectos/99", http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"sprint: id de proyecto no numérico", "/web/proyectos/abc/sprints/1", http.StatusBadRequest, idProyecto},
		{"sprint: id no numérico", "/web/proyectos/7/sprints/abc", http.StatusBadRequest, idSprint},
		{"sprint: id en cero", "/web/proyectos/7/sprints/0", http.StatusBadRequest, idSprint},
		{"sprint: proyecto inexistente", "/web/proyectos/99/sprints/1", http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"sprint: inexistente o de otro proyecto", "/web/proyectos/7/sprints/99", http.StatusNotFound, domain.ErrSprintNoEncontrado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			proyecto, sprint := &consultorHorasProyectoFalso{}, &consultorHorasSprintFalso{}
			cu := casosConAsignacion(nil, nil)
			cu.ConsultarHorasProyecto, cu.ConsultarHorasSprint = proyecto, sprint

			rec := pedirPantalla(t, cu, http.MethodGet, tt.ruta, nil)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantCode, rec.Body)
			}
			esPaginaDeError(t, rec, tt.wantMsg)
			if proyecto.llamado || sprint.llamado {
				t.Error("no se esperaba consultar las horas")
			}
		})
	}
}

// Traza: HU-23 CA-23.17, CA-23.18 (sección 7): cualquier otro error es 500 genérico, sin el detalle.
func TestPantallasHoras_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	for _, p := range pantallasConHoras {
		t.Run(p.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConFallaDeHoras(fmt.Errorf("horas: %w", errors.New("database is locked"))), http.MethodGet, p.ruta, nil)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, se esperaba 500; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, mensajeFallaInternaW)
			if strings.Contains(rec.Body.String(), "database") {
				t.Errorf("la página muestra el detalle interno; cuerpo: %.400s", rec.Body)
			}
		})
	}
}

// html/template escribe el + como &#43;. Se mira el HTML sin decodificar: si el valor saliera como
// template.HTML, o si la plantilla dejara de escapar, el + saldría tal cual.
// Traza: HU-23 CA-23.17, CA-23.18, ADR-0003
func TestPantallasHoras_LosValoresSeEscapan(t *testing.T) {
	for _, p := range pantallasConHoras {
		t.Run(p.nombre, func(t *testing.T) {
			s := seccionHoras(t, pedirPantalla(t, casosConHoras(estimado(8, 10), estimado(4, 3)), http.MethodGet, p.ruta, nil))

			seccionContiene(t, s, "<dd>&#43;1 h</dd>", "<dd>&#43;8,3 %</dd>")
			if html.UnescapeString(s) == s {
				t.Errorf("el bloque no tiene nada escapado: %s", s)
			}
		})
	}
}
