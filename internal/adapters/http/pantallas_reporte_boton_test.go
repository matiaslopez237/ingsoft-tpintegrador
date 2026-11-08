package http

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

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
