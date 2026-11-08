package http

import (
	"html"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// columnasDatoValor son los encabezados de las tablas de dato/valor del documento (Resumen y Esfuerzo).
// En la pantalla esas secciones son una lista de definiciones sin encabezados, así que son los únicos
// textos del documento que no aparecen en ella.
var columnasDatoValor = []string{"Dato", "Valor"}

// reporteDeSprintDeEjemplo es un reporte de un Sprint, con su métrica sin porcentaje (Sprint abierto).
func reporteDeSprintDeEjemplo() domain.Reporte {
	r := reporteDeEjemplo()
	r.Sprint = &domain.Sprint{
		ID: 9, ProyectoID: alfa.ID, Numero: 2, Objetivo: "Cerrar el login",
		FechaInicio: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC),
		Estado:      domain.SprintAbierto,
	}
	r.HayVelocidad, r.Velocidad = false, 0
	r.Metricas = []domain.FilaMetrica{{Sprint: *r.Sprint, DefectosDetectados: 1}}
	return r
}

// reporteConDefectosLargos es el peor caso para el PDF: la tabla más ancha del reporte (Defectos, 6
// columnas) con textos largos sin espacios del carácter más ancho de RN-27.8.
func reporteConDefectosLargos() domain.Reporte {
	r := reporteDeEjemplo()
	largo := strings.Repeat("Œ", 1000)
	for range 50 {
		r.Defectos = append(r.Defectos, domain.Defecto{Descripcion: largo, Severidad: domain.SeveridadAlta, Estado: domain.DefectoAbierto})
	}
	return r
}

// casoReporteDocumento es un reporte y la ruta de su pantalla.
type casoReporteDocumento struct {
	nombre  string
	reporte domain.Reporte
	ruta    string
}

func casosReporteDocumento() []casoReporteDocumento {
	return []casoReporteDocumento{
		{"proyecto", reporteDeEjemplo(), "/web/proyectos/7/reporte"},
		{"sprint", reporteDeSprintDeEjemplo(), "/web/proyectos/7/sprints/9/reporte"},
		{"proyecto vacío", domain.Reporte{Proyecto: alfa}, "/web/proyectos/7/reporte"},
		{"sprint vacío", domain.Reporte{Proyecto: alfa, Sprint: reporteDeSprintDeEjemplo().Sprint}, "/web/proyectos/7/sprints/9/reporte"},
		{"defectos largos", reporteConDefectosLargos(), "/web/proyectos/7/reporte"},
	}
}

var (
	reH1       = regexp.MustCompile(`(?s)<h1>(.*?)</h1>`)
	reH2       = regexp.MustCompile(`(?s)<h2[^>]*>(.*?)</h2>`)
	reSeccion  = regexp.MustCompile(`(?s)<section[^>]*>(.*?)</section>`)
	reEtiqueta = regexp.MustCompile(`<[^>]*>`)
)

// textoVisible saca las etiquetas de un fragmento de HTML, des-escapa las entidades y deja un solo
// espacio entre palabras: lo que lee la persona en la pantalla.
func textoVisible(fragmento string) string {
	return strings.Join(strings.Fields(html.UnescapeString(reEtiqueta.ReplaceAllString(fragmento, " "))), " ")
}

// pantallaDelReporte pide la pantalla del reporte y devuelve su HTML.
func pantallaDelReporte(t *testing.T, c casoReporteDocumento) string {
	t.Helper()
	rec := pedirPantalla(t, casosConReporte(c.reporte, nil), http.MethodGet, c.ruta, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("pantalla del reporte: status = %d, se esperaba 200", rec.Code)
	}
	return rec.Body.String()
}

// textosDeSeccion junta los textos de los bloques de una sección del documento, en orden.
func textosDeSeccion(s domain.SeccionPDF) []string {
	var textos []string
	for _, b := range s.Bloques {
		if b.Parrafo != "" {
			textos = append(textos, b.Parrafo)
		}
		if b.Tabla != nil {
			textos = append(textos, b.Tabla.Columnas...)
			for _, fila := range b.Tabla.Filas {
				textos = append(textos, fila...)
			}
		}
	}
	return textos
}

// Traza: HU-26 CA-26.14 (RN-26.11): las secciones del documento son los <h2> de la pantalla, con los
// mismos títulos y en el mismo orden, y el título del documento es el <h1>.
func TestADocumentoPDF_SeccionesSonLasDeLaPantalla(t *testing.T) {
	for _, c := range casosReporteDocumento() {
		t.Run(c.nombre, func(t *testing.T) {
			pantalla := pantallaDelReporte(t, c)

			doc := aDocumentoPDF(aDatosReporte(c.reporte))

			var h2 []string
			for _, m := range reH2.FindAllStringSubmatch(pantalla, -1) {
				h2 = append(h2, textoVisible(m[1]))
			}
			if len(h2) != 5 {
				t.Fatalf("la pantalla tiene %d <h2>, se esperaban 5: %q", len(h2), h2)
			}
			var titulos []string
			for _, s := range doc.Secciones {
				titulos = append(titulos, s.Titulo)
			}
			if !reflect.DeepEqual(titulos, h2) {
				t.Errorf("secciones del documento = %q, se esperaban los <h2> de la pantalla %q", titulos, h2)
			}
			h1 := reH1.FindStringSubmatch(pantalla)
			if h1 == nil || doc.Titulo != textoVisible(h1[1]) {
				t.Errorf("título del documento = %q, se esperaba el <h1> de la pantalla %v", doc.Titulo, h1)
			}
			if doc.Subtitulo != "" || doc.Aviso != "" || !doc.Emision.IsZero() {
				t.Errorf("subtítulo, aviso y emisión tienen que quedar vacíos: %q, %q, %v", doc.Subtitulo, doc.Aviso, doc.Emision)
			}
		})
	}
}

// Traza: HU-26 CA-26.14 (RN-26.11; HU-27 P-11): cada texto de una sección del documento aparece en la
// misma sección de la pantalla, salvo los encabezados de las tablas de dato/valor.
func TestADocumentoPDF_TextosSonLosDeLaPantalla(t *testing.T) {
	for _, c := range casosReporteDocumento() {
		t.Run(c.nombre, func(t *testing.T) {
			pantalla := pantallaDelReporte(t, c)
			secciones := reSeccion.FindAllStringSubmatch(pantalla, -1)

			doc := aDocumentoPDF(aDatosReporte(c.reporte))

			if len(secciones) != len(doc.Secciones) {
				t.Fatalf("la pantalla tiene %d <section> y el documento %d secciones", len(secciones), len(doc.Secciones))
			}
			for i, s := range doc.Secciones {
				visible := textoVisible(secciones[i][1])
				if len(s.Bloques) == 0 {
					t.Errorf("la sección %q no tiene bloques; en la pantalla dice %q", s.Titulo, visible)
				}
				for _, texto := range textosDeSeccion(s) {
					if esColumnaDatoValor(texto) {
						continue
					}
					if !strings.Contains(visible, texto) {
						t.Errorf("sección %q: el texto %q no está en la pantalla, que dice %.300q", s.Titulo, texto, visible)
					}
				}
			}
		})
	}
}

func esColumnaDatoValor(texto string) bool {
	for _, c := range columnasDatoValor {
		if texto == c {
			return true
		}
	}
	return false
}

// tabla arma una tabla del documento, para escribir los esperados más cortos.
func tabla(columnas []string, filas ...[]string) domain.BloquePDF {
	return domain.BloquePDF{Tabla: &domain.TablaPDF{Columnas: columnas, Filas: filas}}
}

func parrafo(texto string) domain.BloquePDF { return domain.BloquePDF{Parrafo: texto} }

// Traza: HU-26 CA-26.14 (RN-26.11; HU-27 P-11): el reporte de un proyecto, sección por sección, con los
// textos ya formateados de datosReporte: Resumen y Esfuerzo como tablas de dato/valor, el resumen de
// historias y la velocidad como párrafos antes de su tabla.
func TestADocumentoPDF_ReporteDeProyecto(t *testing.T) {
	d := aDatosReporte(reporteDeEjemplo())

	doc := aDocumentoPDF(d)

	want := domain.DocumentoPDF{
		Titulo: "Reporte de Alfa",
		Secciones: []domain.SeccionPDF{
			{Titulo: "Resumen", Bloques: []domain.BloquePDF{tabla(columnasDatoValor,
				[]string{"Nombre", "Alfa"},
				[]string{"Estado", d.Proyecto.Estado},
				[]string{"Inicio", "2026-03-01"},
				[]string{"Fin", "Sin fecha de fin"},
			)}},
			{Titulo: "Historias y estimaciones", Bloques: []domain.BloquePDF{
				parrafo("1 de 2 historias completadas (5 de 13 story points)"),
				tabla([]string{"Título", "Prioridad", "Estado", "Story points"},
					[]string{"Login", d.Historias[0].Prioridad, d.Historias[0].Estado, "5"},
					[]string{"Logout", d.Historias[1].Prioridad, d.Historias[1].Estado, "8"},
				),
			}},
			{Titulo: "Esfuerzo", Bloques: []domain.BloquePDF{tabla(columnasDatoValor,
				[]string{"Horas estimadas", d.EsfuerzoEstimadas},
				[]string{"Horas reales", d.EsfuerzoReales},
				[]string{"Desviación en horas", d.EsfuerzoDesviacionHoras},
				[]string{"Desviación en porcentaje", d.EsfuerzoDesviacionPorcentaje},
			)}},
			{Titulo: "Métricas", Bloques: []domain.BloquePDF{
				parrafo("Velocidad del equipo: 6 puntos por Sprint"),
				tabla([]string{"Sprint", "% completado", "Defectos detectados", "Defectos resueltos"},
					[]string{"1", "75,0 %", "2", "1"},
				),
			}},
			{Titulo: "Defectos", Bloques: []domain.BloquePDF{
				tabla([]string{"Elemento", "Descripción", "Severidad", "Estado", "Sprint detección", "Sprint resolución"},
					// reporteDeEjemplo no les pone ID a las historias ni al defecto: el defecto queda con
					// ElementoID 0, que es el de la última historia, y la pantalla también dice "Logout".
					[]string{"Logout", "Falla el login", d.Defectos[0].Severidad, d.Defectos[0].Estado, "#1", "—"},
				),
			}},
		},
	}
	if !reflect.DeepEqual(doc, want) {
		t.Errorf("aDocumentoPDF() =\n%+v\nse esperaba\n%+v", doc, want)
	}
}

// Traza: HU-26 CA-26.14 (RN-26.5, RN-26.6, RN-26.11): el reporte de un Sprint lleva su resumen propio y
// no tiene velocidad; un porcentaje que no hay dice "Sin datos", como la pantalla.
func TestADocumentoPDF_ReporteDeSprint(t *testing.T) {
	d := aDatosReporte(reporteDeSprintDeEjemplo())

	doc := aDocumentoPDF(d)

	if doc.Titulo != "Reporte del Sprint 2 de Alfa" {
		t.Errorf("título = %q, se esperaba %q", doc.Titulo, "Reporte del Sprint 2 de Alfa")
	}
	wantResumen := []domain.BloquePDF{tabla(columnasDatoValor,
		[]string{"Número", "2"},
		[]string{"Sprint Goal", "Cerrar el login"},
		[]string{"Inicio", "2026-04-01"},
		[]string{"Fin", "2026-04-14"},
		[]string{"Estado", "Abierto"},
	)}
	if !reflect.DeepEqual(doc.Secciones[0].Bloques, wantResumen) {
		t.Errorf("Resumen = %+v, se esperaba %+v", doc.Secciones[0].Bloques, wantResumen)
	}
	wantMetricas := []domain.BloquePDF{
		tabla([]string{"Sprint", "% completado", "Defectos detectados", "Defectos resueltos"},
			[]string{"2", "Sin datos", "1", "0"},
		),
	}
	if !reflect.DeepEqual(doc.Secciones[3].Bloques, wantMetricas) {
		t.Errorf("Métricas = %+v, se esperaba solo la tabla, sin velocidad: %+v", doc.Secciones[3].Bloques, wantMetricas)
	}
}

// Traza: HU-26 CA-26.14, CA-26.15 (RN-26.9, RN-26.11): un reporte sin datos lleva en cada sección el
// mismo mensaje que la pantalla, como párrafo y no como tabla vacía.
func TestADocumentoPDF_ReporteVacio_MensajesDeLaPantalla(t *testing.T) {
	tests := []struct {
		nombre  string
		reporte domain.Reporte
		want    [][]domain.BloquePDF // Historias, Métricas, Defectos
	}{
		{"proyecto", domain.Reporte{Proyecto: alfa}, [][]domain.BloquePDF{
			{parrafo("El reporte no tiene historias.")},
			{parrafo("Todavía no hay datos de velocidad (sin sprints cerrados)."), parrafo("El reporte no tiene sprints cerrados.")},
			{parrafo("El reporte no tiene defectos.")},
		}},
		{"sprint", domain.Reporte{Proyecto: alfa, Sprint: reporteDeSprintDeEjemplo().Sprint}, [][]domain.BloquePDF{
			{parrafo("El reporte no tiene historias.")},
			{parrafo("El reporte no tiene sprints cerrados.")},
			{parrafo("El reporte no tiene defectos.")},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			doc := aDocumentoPDF(aDatosReporte(tt.reporte))

			got := [][]domain.BloquePDF{doc.Secciones[1].Bloques, doc.Secciones[3].Bloques, doc.Secciones[4].Bloques}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Historias, Métricas y Defectos =\n%+v\nse esperaba\n%+v", got, tt.want)
			}
			if len(doc.Secciones[0].Bloques) != 1 || len(doc.Secciones[2].Bloques) != 1 {
				t.Errorf("Resumen y Esfuerzo tienen que seguir con su tabla: %+v, %+v", doc.Secciones[0].Bloques, doc.Secciones[2].Bloques)
			}
		})
	}
}

// Traza: HU-26 CA-26.14 (RN-26.11; HU-27 RN-27.11, P-10): el documento que sale de cualquier reporte es
// válido, también con la tabla de 6 columnas llena de textos largos, antes y después de normalizarlo
// como hace el caso de uso de HU-27.
func TestADocumentoPDF_SiempreValido(t *testing.T) {
	for _, c := range casosReporteDocumento() {
		t.Run(c.nombre, func(t *testing.T) {
			doc := aDocumentoPDF(aDatosReporte(c.reporte))

			if err := doc.Validar(); err != nil {
				t.Errorf("Validar() = %v, se esperaba nil", err)
			}
			normalizado, _ := domain.NormalizarDocumento(doc)
			if err := normalizado.Validar(); err != nil {
				t.Errorf("Validar() del documento normalizado = %v, se esperaba nil", err)
			}
			for _, s := range doc.Secciones {
				for _, b := range s.Bloques {
					if b.Tabla != nil && len(b.Tabla.Columnas) > domain.MaximoColumnas {
						t.Errorf("sección %q: tabla de %d columnas", s.Titulo, len(b.Tabla.Columnas))
					}
				}
			}
		})
	}
}
