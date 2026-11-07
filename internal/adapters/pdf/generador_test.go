package pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"

	lectorpdf "github.com/ledongthuc/pdf"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El generador implementa el puerto del caso de uso (spec HU-27, sección 5).
var _ app.GeneradorPDF = (*Generador)(nil)

var emisionDePrueba = time.Date(2026, 11, 7, 14, 5, 0, 0, time.FixedZone("UTC-3", -3*60*60))

// generar genera el PDF y falla el test si hay error.
func generar(t *testing.T, doc domain.DocumentoPDF) []byte {
	t.Helper()
	b, err := NuevoGenerador().Generar(context.Background(), doc)
	if err != nil {
		t.Fatalf("Generar: %v", err)
	}
	return b
}

var espacios = regexp.MustCompile(`\s+`)

// paginas lee el PDF de vuelta y devuelve el texto de cada página, con los espacios normalizados: el
// extractor corta donde el PDF cambia de línea (ADR-0006). Leerlo prueba que el PDF es válido; la
// cantidad de páginas es len del resultado.
func paginas(t *testing.T, b []byte) []string {
	t.Helper()
	r, err := lectorpdf.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatalf("el PDF no se puede leer: %v", err)
	}
	var textos []string
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			t.Fatalf("la página %d no existe", i)
		}
		s, err := p.GetPlainText(nil)
		if err != nil {
			t.Fatalf("texto de la página %d: %v", i, err)
		}
		textos = append(textos, strings.TrimSpace(espacios.ReplaceAllString(s, " ")))
	}
	return textos
}

// contieneEnOrden falla si los fragmentos no aparecen en el texto en ese orden.
func contieneEnOrden(t *testing.T, texto string, fragmentos ...string) {
	t.Helper()
	desde := 0
	for _, f := range fragmentos {
		i := strings.Index(texto[desde:], f)
		if i < 0 {
			t.Fatalf("no aparece %q (en orden) en el texto:\n%s", f, texto)
		}
		desde += i + len(f)
	}
}

func documentoNormal() domain.DocumentoPDF {
	return domain.DocumentoPDF{
		Titulo:    "Proyecto Ñandú: planificación ¿cumplida? ¡Sí!",
		Subtitulo: "Del 26/10/2026 al 08/11/2026",
		Emision:   emisionDePrueba,
		Aviso:     domain.TextoAvisoReemplazos(1),
		Secciones: []domain.SeccionPDF{
			{Titulo: "Resumen", Bloques: []domain.BloquePDF{{Parrafo: "El equipo terminó la migración; sigue la revisión."}}},
			{Titulo: "Historias y estimaciones", Bloques: []domain.BloquePDF{{Tabla: &domain.TablaPDF{
				Columnas: []string{"Historia", "Estado", "Puntos"},
				Filas: [][]string{
					{"Diseño de la pantalla", "Hecho", "5"},
					{"Reportes", "En progreso", "3"},
					{"Señales: áéíóúüñ€°ªº–—“”…→≤≥%", "Pendiente", "8"},
				},
			}}}},
		},
	}
}

// CA-27.1, CA-27.2: el documento sale completo y en orden, con acentos y ñ (RN-27.1, RN-27.2, RN-27.14).
func TestGenerar_DocumentoNormal_TextoEnOrdenEnUnaPagina(t *testing.T) {
	textos := paginas(t, generar(t, documentoNormal()))
	if len(textos) != 1 {
		t.Fatalf("páginas = %d, se esperaba 1", len(textos))
	}
	contieneEnOrden(t, textos[0],
		"Proyecto Ñandú: planificación ¿cumplida? ¡Sí!",
		"Del 26/10/2026 al 08/11/2026",
		"Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3).",
		"Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?».",
		"Resumen",
		"El equipo terminó la migración; sigue la revisión.",
		"Historias y estimaciones",
		"Historia", "Estado", "Puntos",
		"Diseño de la pantalla", "Hecho", "5",
		"Reportes", "En progreso", "3",
		"Señales: áéíóúüñ€°ªº–—“”…→≤≥%", "Pendiente", "8",
	)
}

// CA-27.7: sin subtítulo ni aviso, esas líneas no aparecen.
func TestGenerar_SinSubtituloNiAviso_NoLosMuestra(t *testing.T) {
	doc := documentoNormal()
	doc.Subtitulo, doc.Aviso = "", ""
	texto := paginas(t, generar(t, doc))[0]
	contieneEnOrden(t, texto, "Proyecto Ñandú", "Emitido el 07/11/2026", "Resumen")
	if strings.Contains(texto, "Del 26/10") || strings.Contains(texto, "Aviso:") {
		t.Errorf("aparece un subtítulo o un aviso que el documento no tiene:\n%s", texto)
	}
}

// CA-27.24: el resultado es un PDF completo: empieza con %PDF- y termina con %%EOF (como mucho seguido
// de un salto de línea), y se puede volver a leer (RN-27.15).
func TestGenerar_EsUnPDFCompleto(t *testing.T) {
	b := generar(t, documentoNormal())
	if !bytes.HasPrefix(b, []byte("%PDF-")) {
		t.Errorf("no empieza con %%PDF-: %q", b[:min(len(b), 8)])
	}
	fin := bytes.TrimRight(b, "\r\n")
	if !bytes.HasSuffix(fin, []byte("%%EOF")) || len(b)-len(fin) > 2 {
		t.Errorf("no termina con %%%%EOF: %q", b[max(0, len(b)-12):])
	}
	paginas(t, b)
}

// Sección 7: si fpdf falla, el generador devuelve un error, no entra en pánico ni devuelve bytes. Un
// emoji no llega normalizado desde el caso de uso, pero es la forma conocida de hacer fallar a fpdf
// (ADR-0006).
func TestGenerar_FallaDeFpdf_DevuelveError(t *testing.T) {
	doc := documentoNormal()
	doc.Secciones[0].Bloques[0].Parrafo = "Texto con 😀"
	b, err := NuevoGenerador().Generar(context.Background(), doc)
	if err == nil {
		t.Fatal("se esperaba un error de fpdf")
	}
	if b != nil {
		t.Errorf("devolvió %d bytes junto con el error", len(b))
	}
}

// Un pedido cancelado no genera nada.
func TestGenerar_ContextoCancelado_DevuelveError(t *testing.T) {
	ctx, cancelar := context.WithCancel(context.Background())
	cancelar()
	b, err := NuevoGenerador().Generar(ctx, documentoNormal())
	if !errors.Is(err, context.Canceled) || b != nil {
		t.Errorf("Generar = (%d bytes, %v), se esperaba context.Canceled sin bytes", len(b), err)
	}
}

// compacto quita todos los espacios: el extractor pega las líneas de un párrafo sin separador.
func compacto(s string) string { return espacios.ReplaceAllString(s, "") }

// documentoConTabla arma un documento con una sola sección que tiene la tabla.
func documentoConTabla(t *domain.TablaPDF) domain.DocumentoPDF {
	return domain.DocumentoPDF{
		Titulo:    "Reporte con tabla",
		Emision:   emisionDePrueba,
		Secciones: []domain.SeccionPDF{{Titulo: "Datos", Bloques: []domain.BloquePDF{{Tabla: t}}}},
	}
}

// filaDePrueba es la fila i de una tabla de prueba: cada celda dice su columna y su número de fila.
func filaDePrueba(i int) []string {
	return []string{fmt.Sprintf("A-%03d", i), fmt.Sprintf("B-%03d", i), fmt.Sprintf("C-%03d", i)}
}

// paginaDe devuelve en qué página aparece el texto y cuántas veces aparece en todo el documento.
func paginaDe(textos []string, s string) (pagina, veces int) {
	pagina = -1
	for i, texto := range textos {
		if n := strings.Count(texto, s); n > 0 {
			veces += n
			if pagina < 0 {
				pagina = i
			}
		}
	}
	return pagina, veces
}

// CA-27.3: una tabla de 200 filas cruza páginas; cada fila aparece completa, una sola vez, con sus tres
// celdas en la misma página y en orden (RN-27.4, RN-27.5). El texto solo no alcanza (ADR-0006): también
// se comprueba la cantidad de páginas: las 201 filas (con la de encabezados) de una línea, de
// altoLineaTabla cada una, ocupan lo que dan, repartido en hojas A4 con los márgenes del adaptador; como
// las filas no se parten y la primera hoja también tiene el título, se acepta una página más.
func TestGenerar_TablaLarga_CruzaPaginasSinPerderNiRepetirFilas(t *testing.T) {
	tabla := &domain.TablaPDF{Columnas: []string{"Columna A", "Columna B", "Columna C"}}
	for i := 1; i <= 200; i++ {
		tabla.Filas = append(tabla.Filas, filaDePrueba(i))
	}
	textos := paginas(t, generar(t, documentoConTabla(tabla)))
	minimo := int(math.Ceil(201 * altoLineaTabla / (297 - margenSuperior - margenInferior)))
	if len(textos) < minimo || len(textos) > minimo+1 {
		t.Errorf("páginas = %d, se esperaban %d o %d", len(textos), minimo, minimo+1)
	}
	anterior := 0
	for i := 1; i <= 200; i++ {
		celdas := filaDePrueba(i)
		pagina, veces := paginaDe(textos, celdas[0])
		if veces != 1 {
			t.Fatalf("la fila %d aparece %d veces", i, veces)
		}
		for _, c := range celdas[1:] {
			if p, v := paginaDe(textos, c); v != 1 || p != pagina {
				t.Fatalf("la celda %s está en la página %d (%d veces); su fila, en la %d", c, p+1, v, pagina+1)
			}
		}
		if pagina < anterior {
			t.Fatalf("la fila %d está en la página %d, antes que la fila anterior", i, pagina+1)
		}
		anterior = pagina
	}
	if anterior != len(textos)-1 {
		t.Errorf("la última fila está en la página %d de %d: sobran páginas", anterior+1, len(textos))
	}
}

// CA-27.3, RN-27.5: el texto de una celda se ajusta al ancho de su columna en varias líneas. 30 filas
// con una celda de 300 caracteres no entran en una hoja si cada celda ocupa varias líneas: con tres
// columnas iguales (unos 60 mm, unos 35 caracteres por línea) cada fila tiene unas 9 líneas, unos
// 40 mm, así que entran unas 6 filas por hoja y hacen falta 5 o más páginas.
func TestGenerar_CeldaLarga_SeAjustaEnVariasLineas(t *testing.T) {
	largo := strings.TrimSpace(strings.Repeat("palabra ", 37) + "fin") // 299 caracteres
	tabla := &domain.TablaPDF{Columnas: []string{"Texto", "B", "C"}}
	for i := 1; i <= 30; i++ {
		tabla.Filas = append(tabla.Filas, []string{largo, fmt.Sprintf("B-%03d", i), "c"})
	}
	textos := paginas(t, generar(t, documentoConTabla(tabla)))
	if len(textos) < 5 {
		t.Errorf("páginas = %d: las celdas no se ajustaron en varias líneas", len(textos))
	}
	if _, veces := paginaDe(textos, "B-030"); veces != 1 {
		t.Errorf("la última fila aparece %d veces", veces)
	}
}

// CA-27.4: un párrafo de 10 000 caracteres continúa en las páginas siguientes y se extrae completo.
func TestGenerar_ParrafoLargo_ContinuaEnOtrasPaginas(t *testing.T) {
	var sb strings.Builder
	for i := 1; sb.Len() < 10000; i++ {
		fmt.Fprintf(&sb, "frase%04d con acentos y eñes. ", i)
	}
	parrafo := strings.TrimSpace(sb.String())
	doc := domain.DocumentoPDF{Titulo: "Párrafo largo", Emision: emisionDePrueba, Secciones: []domain.SeccionPDF{
		{Titulo: "Texto", Bloques: []domain.BloquePDF{{Parrafo: parrafo}}},
	}}
	textos := paginas(t, generar(t, doc))
	if len(textos) < 2 {
		t.Fatalf("páginas = %d, se esperaban 2 o más", len(textos))
	}
	if !strings.Contains(compacto(strings.Join(textos, "")), compacto(parrafo)) {
		t.Errorf("el párrafo no se extrae completo y en orden")
	}
}

// CA-27.7: un documento sin subtítulo ni secciones es una página con título y emisión, y nada más.
func TestGenerar_DocumentoMinimo(t *testing.T) {
	doc := domain.DocumentoPDF{Titulo: "Reporte vacío", Emision: emisionDePrueba}
	textos := paginas(t, generar(t, doc))
	if len(textos) != 1 {
		t.Fatalf("páginas = %d, se esperaba 1", len(textos))
	}
	if got, want := compacto(textos[0]), compacto("Reporte vacío Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3)."); got != want {
		t.Errorf("texto = %q, se esperaba %q", got, want)
	}
}

// CA-27.7, CA-27.22: las secciones salen en orden aunque estén vacías; una tabla sin filas muestra sus
// encabezados y "Sin datos."; un párrafo vacío no muestra nada (RN-27.7, RN-27.11).
func TestGenerar_SeccionesVacias(t *testing.T) {
	doc := domain.DocumentoPDF{Titulo: "Vacíos", Emision: emisionDePrueba, Secciones: []domain.SeccionPDF{
		{Titulo: "Notas", Bloques: []domain.BloquePDF{{}}},
		{Titulo: "Esfuerzo"},
		{Titulo: "Defectos", Bloques: []domain.BloquePDF{{Tabla: &domain.TablaPDF{Columnas: []string{"Defecto", "Severidad"}}}}},
		{Titulo: "Métricas", Bloques: []domain.BloquePDF{{Parrafo: "Velocidad: 8 puntos."}}},
	}}
	textos := paginas(t, generar(t, doc))
	if len(textos) != 1 {
		t.Fatalf("páginas = %d, se esperaba 1", len(textos))
	}
	want := compacto("Vacíos Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3). " +
		"Notas Esfuerzo Defectos Defecto Severidad Sin datos. Métricas Velocidad: 8 puntos.")
	if got := compacto(textos[0]); got != want {
		t.Errorf("texto = %q\nse esperaba %q", got, want)
	}
}
