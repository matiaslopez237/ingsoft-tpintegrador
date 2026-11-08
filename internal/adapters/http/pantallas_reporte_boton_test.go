package http

import (
	"html"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// botonDescargarPDF es el enlace de CA-27.21, tal cual: sin aria-label, sin target, sin download.
func botonDescargarPDF(ruta string) string {
	return `<a class="boton" href="` + ruta + `">Descargar PDF</a>`
}

// reBotones toma el contenido de <div class="botones"> de la cabecera de la pantalla.
var reBotones = regexp.MustCompile(`(?s)<div class="botones">(.*?)</div>`)

// Traza: HU-27 CA-27.21: la pantalla del reporte de un proyecto y la de un Sprint tienen el enlace
// "Descargar PDF" a su ruta de descarga, dentro de los botones de la cabecera y después de "Volver…",
// sin formulario ni JavaScript en la página.
func TestPantallaReporte_BotonDescargarPDF(t *testing.T) {
	sprint := reporteDeSprintDeEjemplo() // Sprint con ID 9 y número 2: el enlace lleva el ID
	tests := []struct {
		nombre, ruta, href, volver string
		reporte                    domain.Reporte
	}{
		{"proyecto", "/web/proyectos/7/reporte", "/web/proyectos/7/reporte/pdf", "Volver al proyecto", reporteParaDocumento()},
		{"sprint", "/web/proyectos/7/sprints/9/reporte", "/web/proyectos/7/sprints/9/reporte/pdf", "Volver a sprints", sprint},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConReporte(tt.reporte, nil), http.MethodGet, tt.ruta, nil)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200", rec.Code)
			}
			cuerpo := rec.Body.String()

			botones := reBotones.FindStringSubmatch(cuerpo)
			if botones == nil {
				t.Fatalf("la pantalla no tiene <div class=\"botones\">")
			}
			boton := botonDescargarPDF(tt.href)
			volver := strings.Index(botones[1], tt.volver)
			enlace := strings.Index(botones[1], boton)
			if volver < 0 || enlace < 0 || enlace < volver {
				t.Errorf("los botones tienen que tener %q después de %q; tienen %q", boton, tt.volver, botones[1])
			}
			if n := strings.Count(cuerpo, "Descargar PDF"); n != 1 {
				t.Errorf("\"Descargar PDF\" aparece %d veces, se esperaba 1", n)
			}
			for _, prohibido := range []string{"<form", "<script"} {
				if strings.Contains(cuerpo, prohibido) {
					t.Errorf("la pantalla del reporte no tiene que tener %s", prohibido)
				}
			}
		})
	}
}

// reBotonPDF toma el href del enlace "Descargar PDF".
var reBotonPDF = regexp.MustCompile(`<a class="boton" href="([^"]*)">Descargar PDF</a>`)

// Traza: HU-27 CA-27.21, CA-27.14, CA-27.15: en NuevoRouter, el enlace de la pantalla del reporte lleva a
// la descarga del PDF de ese mismo proyecto o Sprint.
func TestPantallaReporte_BotonDescargarPDF_LlevaALaDescarga(t *testing.T) {
	tests := []struct {
		nombre, ruta         string
		reporte              domain.Reporte
		proyectoID, sprintID int64
	}{
		{"proyecto", "/web/proyectos/7/reporte", reporteParaDocumento(), 7, 0},
		{"sprint", "/web/proyectos/7/sprints/9/reporte", reporteDeSprintDeEjemplo(), 7, 9},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			archivo := app.ArchivoPDF{Nombre: "reporte.pdf", Contenido: contenidoPDFFalso}
			proyecto := &exportadorReporteProyectoFalso{archivo: archivo}
			sprint := &exportadorReporteSprintFalso{archivo: archivo}
			cu := casosConReporte(tt.reporte, nil)
			cu.ExportarReporteProyecto, cu.ExportarReporteSprint = proyecto, sprint
			mux := NuevoRouter(cu)

			pantalla := httptest.NewRecorder()
			mux.ServeHTTP(pantalla, httptest.NewRequest(http.MethodGet, tt.ruta, nil))
			m := reBotonPDF.FindStringSubmatch(pantalla.Body.String())
			if m == nil {
				t.Fatalf("la pantalla (status %d) no tiene el enlace Descargar PDF", pantalla.Code)
			}
			descarga := httptest.NewRecorder()
			mux.ServeHTTP(descarga, httptest.NewRequest(http.MethodGet, html.UnescapeString(m[1]), nil))

			if descarga.Code != http.StatusOK || descarga.Header().Get("Content-Type") != "application/pdf" {
				t.Fatalf("GET %s: status %d, Content-Type %q; se esperaba 200 con application/pdf",
					m[1], descarga.Code, descarga.Header().Get("Content-Type"))
			}
			if tt.sprintID == 0 && (!proyecto.llamado || proyecto.proyectoID != tt.proyectoID || sprint.llamado) {
				t.Errorf("se esperaba la descarga del proyecto %d; proyecto llamado %v con %d, sprint llamado %v",
					tt.proyectoID, proyecto.llamado, proyecto.proyectoID, sprint.llamado)
			}
			if tt.sprintID != 0 && (!sprint.llamado || sprint.proyectoID != tt.proyectoID || sprint.sprintID != tt.sprintID || proyecto.llamado) {
				t.Errorf("se esperaba la descarga del sprint %d del proyecto %d; sprint llamado %v con %d y %d, proyecto llamado %v",
					tt.sprintID, tt.proyectoID, sprint.llamado, sprint.proyectoID, sprint.sprintID, proyecto.llamado)
			}
		})
	}
}
