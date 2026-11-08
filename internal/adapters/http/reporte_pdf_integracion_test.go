package http

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	lectorpdf "github.com/ledongthuc/pdf"

	pdfadapter "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/pdf"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// letraTituloSeccion es el tamaño de letra de los títulos de sección en el PDF (adaptador de PDF): ningún
// otro texto del documento usa ese tamaño, así que identifica los títulos sin confundirlos con una
// columna que diga lo mismo (por ejemplo, "Defectos detectados").
const letraTituloSeccion = 12

// lectoresDelReporte son el lector de proyectos y el de sprints del caso de uso de HU-27: alfa (ID 7) y
// su Sprint 9.
type lectoresDelReporte struct{}

func (lectoresDelReporte) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if id != alfa.ID {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return alfa, nil
}

func (lectoresDelReporte) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	if id != 9 {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return domain.Sprint{ID: 9, ProyectoID: alfa.ID, Numero: 2}, nil
}

// titulosDeSeccionDelPDF lee el PDF y devuelve, en orden, los textos escritos con la letra de los
// títulos de sección. Un título que sigue en la línea siguiente se junta con un espacio.
func titulosDeSeccionDelPDF(t *testing.T, b []byte) []string {
	t.Helper()
	r, err := lectorpdf.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatalf("el PDF no se puede leer: %v", err)
	}
	var titulos []string
	var actual strings.Builder
	cerrar := func() {
		if s := strings.Join(strings.Fields(actual.String()), " "); s != "" {
			titulos = append(titulos, s)
		}
		actual.Reset()
	}
	for i := 1; i <= r.NumPage(); i++ {
		yAnterior := -1.0
		for _, letra := range r.Page(i).Content().Text {
			if letra.FontSize != letraTituloSeccion {
				cerrar()
				yAnterior = -1
				continue
			}
			if yAnterior >= 0 && letra.Y != yAnterior {
				actual.WriteString(" ")
			}
			actual.WriteString(letra.S)
			yAnterior = letra.Y
		}
		cerrar()
	}
	return titulos
}

// Traza: HU-27 CA-27.23 (RN-27.1, RN-27.2, RN-27.16; P-11): de punta a punta con la fuente, el caso de
// uso y el generador reales, los títulos de sección del PDF son los <h2> de la pantalla del reporte,
// exactos y en el mismo orden, para el reporte de un proyecto y el de un Sprint. El del proyecto tiene
// 50 defectos largos, así que las secciones siguen en varias páginas.
func TestExportarReporte_TitulosDelPDFSonLosDeLaPantalla(t *testing.T) {
	reloj := func() time.Time { return time.Date(2026, 11, 7, 17, 5, 0, 0, time.UTC) }
	tests := []struct {
		nombre, ruta  string
		reporte       domain.Reporte
		variasPaginas bool
		exportar      func(fuente *fuenteDocumentoReporte) (app.ArchivoPDF, error)
	}{
		{"proyecto", "/web/proyectos/7/reporte", reporteConDefectosLargos(), true,
			func(fuente *fuenteDocumentoReporte) (app.ArchivoPDF, error) {
				uc := app.NuevoExportarReporteProyecto(lectoresDelReporte{}, fuente, pdfadapter.NuevoGenerador(), reloj)
				return uc.Ejecutar(context.Background(), 7)
			}},
		{"sprint", "/web/proyectos/7/sprints/9/reporte", reporteDeSprintDeEjemplo(), false,
			func(fuente *fuenteDocumentoReporte) (app.ArchivoPDF, error) {
				uc := app.NuevoExportarReporteSprint(lectoresDelReporte{}, lectoresDelReporte{}, fuente, pdfadapter.NuevoGenerador(), reloj)
				return uc.Ejecutar(context.Background(), 7, 9)
			}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, proyecto, sprint := armadoresQueAnotan(tt.reporte, nil)
			fuente := NuevaFuenteDocumentoReporte(proyecto, sprint)
			var h2 []string
			for _, m := range reH2.FindAllStringSubmatch(pantallaDelReporte(t, casoReporteDocumento{tt.nombre, tt.reporte, tt.ruta}), -1) {
				h2 = append(h2, textoVisible(m[1]))
			}

			archivo, err := tt.exportar(fuente)

			if err != nil {
				t.Fatalf("exportar: %v", err)
			}
			if r, err := lectorpdf.NewReader(bytes.NewReader(archivo.Contenido), int64(len(archivo.Contenido))); err != nil || (r.NumPage() > 1) != tt.variasPaginas {
				t.Fatalf("el PDF no se puede leer o no tiene las páginas esperadas (varias: %v): %v", tt.variasPaginas, err)
			}
			if got := titulosDeSeccionDelPDF(t, archivo.Contenido); !reflect.DeepEqual(got, h2) {
				t.Errorf("títulos de sección del PDF = %q, se esperaban los <h2> de la pantalla %q", got, h2)
			}
		})
	}
}
