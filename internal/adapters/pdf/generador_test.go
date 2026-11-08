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

// caracterMasAncho es el carácter más ancho de la lista de RN-27.8 en DejaVu Sans Condensed (unos 3 mm
// a 9 puntos). Un texto sin espacios de este carácter es el peor caso para el alto de una fila (P-10).
const caracterMasAncho = "Œ"

// Traza: HU-27 CA-27.6 (RN-27.5, RN-27.6, P-10, sección 5): justo en los topes, 6 columnas con
// encabezados de 100 caracteres y celdas de 300 del carácter más ancho, cada página tiene la fila de
// encabezados y exactamente una fila entera. Se cuentan los Œ de cada página: 6×100 + 6×300.
func TestGenerar_TablaEnLosTopes_UnaFilaPorPagina(t *testing.T) {
	const filas = 3
	tabla := &domain.TablaPDF{}
	for range 6 {
		tabla.Columnas = append(tabla.Columnas, strings.Repeat(caracterMasAncho, 100))
	}
	for range filas {
		var fila []string
		for range 6 {
			fila = append(fila, strings.Repeat(caracterMasAncho, 300))
		}
		tabla.Filas = append(tabla.Filas, fila)
	}
	doc := documentoConTabla(tabla)
	if err := doc.Validar(); err != nil {
		t.Fatalf("la tabla en los topes tiene que ser válida: %v", err)
	}

	textos := paginas(t, generar(t, doc))

	if len(textos) != filas {
		t.Fatalf("páginas = %d, se esperaba una por fila (%d)", len(textos), filas)
	}
	for i, texto := range textos {
		if n := strings.Count(texto, caracterMasAncho); n != 6*100+6*300 {
			t.Errorf("página %d: %d caracteres %s, se esperaban %d (encabezados y una fila entera)", i+1, n, caracterMasAncho, 6*100+6*300)
		}
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
	if !strings.Contains(compacto(strings.Join(cuerpos(textos, doc.Titulo), "")), compacto(parrafo)) {
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
	if got, want := compacto(cuerpos(textos, doc.Titulo)[0]), compacto("Reporte vacío Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3)."); got != want {
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
	if got := compacto(cuerpos(textos, doc.Titulo)[0]); got != want {
		t.Errorf("texto = %q\nse esperaba %q", got, want)
	}
}

// tablaLarga es el documento de la tabla de 200 filas de CA-27.3.
func tablaLarga() domain.DocumentoPDF {
	tabla := &domain.TablaPDF{Columnas: []string{"Columna A", "Columna B", "Columna C"}}
	for i := 1; i <= 200; i++ {
		tabla.Filas = append(tabla.Filas, filaDePrueba(i))
	}
	return documentoConTabla(tabla)
}

// CA-27.3, RN-27.3: todas las páginas tienen el pie "Página k de N", con N igual a la cantidad de
// páginas que lee el extractor; desde la segunda, empiezan con el título del documento.
func TestGenerar_EncabezadoYPie(t *testing.T) {
	textos := paginas(t, generar(t, tablaLarga()))
	if len(textos) < 2 {
		t.Fatalf("páginas = %d, se esperaban 2 o más", len(textos))
	}
	for i, texto := range textos {
		pie := fmt.Sprintf("Página %d de %d", i+1, len(textos))
		if strings.Count(texto, pie) != 1 {
			t.Errorf("la página %d no tiene el pie %q:\n%s", i+1, pie, texto)
		}
		titulo := strings.Count(texto, "Reporte con tabla")
		if i == 0 && titulo != 1 {
			t.Errorf("la primera página tiene el título %d veces, se esperaba 1 (sin encabezado)", titulo)
		}
		if i > 0 && !strings.HasPrefix(texto, "Reporte con tabla") {
			t.Errorf("la página %d no empieza con el encabezado:\n%s", i+1, texto)
		}
	}
}

// CA-27.3, RN-27.5: la fila de encabezados se repite en cada página donde la tabla continúa.
func TestGenerar_TablaLarga_RepiteLosEncabezadosEnCadaPagina(t *testing.T) {
	textos := paginas(t, generar(t, tablaLarga()))
	for i, texto := range textos {
		if strings.Count(texto, "Columna A Columna B Columna C") != 1 {
			t.Errorf("la página %d no tiene la fila de encabezados una vez:\n%s", i+1, texto)
		}
	}
}

// CA-27.6, RN-27.5: una fila de varias líneas nunca se parte entre dos páginas: el texto entero de cada
// celda larga está en una sola página.
func TestGenerar_FilaDeVariasLineas_NoSeParte(t *testing.T) {
	tabla := &domain.TablaPDF{Columnas: []string{"Texto", "Fila"}}
	for i := 1; i <= 40; i++ {
		celda := fmt.Sprintf("inicio%02d ", i) + strings.Repeat("palabra ", 34) + fmt.Sprintf("final%02d", i)
		tabla.Filas = append(tabla.Filas, []string{celda, fmt.Sprintf("F-%02d", i)})
	}
	textos := paginas(t, generar(t, documentoConTabla(tabla)))
	if len(textos) < 3 {
		t.Fatalf("páginas = %d: la prueba necesita que la tabla cruce varias páginas", len(textos))
	}
	for i, fila := range tabla.Filas {
		enteras := 0
		for _, texto := range textos {
			if strings.Contains(compacto(texto), compacto(fila[0])) {
				enteras++
			}
		}
		if enteras != 1 {
			t.Errorf("la celda de la fila %d está entera en %d páginas, se esperaba 1", i+1, enteras)
		}
	}
}

// RN-27.4: el título de una sección no queda solo al pie. Se prueba con tablas previas de distinto
// largo, para que el título caiga en todas las alturas posibles cerca del final de la primera página: el
// título siempre está en la misma página que la primera fila de su tabla, o que la primera línea de su
// párrafo.
func TestGenerar_TituloDeSeccion_NoQuedaSoloAlPie(t *testing.T) {
	for _, bloque := range []domain.BloquePDF{
		{Tabla: &domain.TablaPDF{Columnas: []string{"Encabezado final"}, Filas: [][]string{{"primera-fila"}}}},
		{Parrafo: "primera-linea del párrafo final"},
	} {
		for filas := 45; filas <= 60; filas++ {
			previa := &domain.TablaPDF{Columnas: []string{"A", "B", "C"}}
			for i := 1; i <= filas; i++ {
				previa.Filas = append(previa.Filas, filaDePrueba(i))
			}
			doc := documentoConTabla(previa)
			doc.Secciones = append(doc.Secciones, domain.SeccionPDF{Titulo: "Sección final", Bloques: []domain.BloquePDF{bloque}})
			textos := paginas(t, generar(t, doc))
			titulo, _ := paginaDe(textos, "Sección final")
			primera, _ := paginaDe(textos, "primera-")
			if titulo != primera {
				t.Errorf("con %d filas previas, el título está en la página %d y su primer contenido en la %d", filas, titulo+1, primera+1)
			}
		}
	}
}

// CA-27.12: con la misma emisión, dos generaciones dan los mismos bytes, también en paralelo; con otra
// emisión, otros bytes, porque la fecha de creación del PDF es la emisión (RN-27.10).
func TestGenerar_Determinista(t *testing.T) {
	primero := generar(t, documentoNormal())
	resultados := make(chan []byte, 8)
	for range 8 {
		go func() {
			b, _ := NuevoGenerador().Generar(context.Background(), documentoNormal())
			resultados <- b
		}()
	}
	for range 8 {
		if b := <-resultados; !bytes.Equal(b, primero) {
			t.Fatal("dos generaciones del mismo documento dieron bytes distintos")
		}
	}
	otro := documentoNormal()
	otro.Emision = otro.Emision.Add(time.Minute)
	if bytes.Equal(generar(t, otro), primero) {
		t.Error("con otra emisión dio los mismos bytes: la fecha del PDF no es la emisión")
	}
}

var piePagina = regexp.MustCompile(`\s*Página \d+ de \d+\s*$`)

// cuerpos devuelve el texto de cada página sin el encabezado (el título, desde la segunda) ni el pie,
// para comparar solo el contenido del documento (RN-27.3).
func cuerpos(textos []string, titulo string) []string {
	var r []string
	for i, texto := range textos {
		texto = piePagina.ReplaceAllString(texto, "")
		if i > 0 {
			texto = strings.TrimSpace(strings.TrimPrefix(texto, titulo))
		}
		r = append(r, texto)
	}
	return r
}

// CA-27.12, RN-27.10: la fecha de creación y la de modificación del PDF son la emisión, y no la hora en
// que se generó. Sin esto, dos exportaciones en segundos distintos darían bytes distintos.
func TestGenerar_LaFechaDelPDFEsLaEmision(t *testing.T) {
	b := generar(t, documentoNormal())
	for _, campo := range []string{"/CreationDate (D:20261107140500)", "/ModDate (D:20261107140500)"} {
		if !bytes.Contains(b, []byte(campo)) {
			t.Errorf("el PDF no tiene %s", campo)
		}
	}
}

// RN-27.4, RN-27.5: una tabla que no es el primer bloque de su sección tampoco deja su fila de
// encabezados sola al pie: los encabezados están en la misma página que la primera fila.
func TestGenerar_TablaDespuesDeUnParrafo_EncabezadosJuntoALaPrimeraFila(t *testing.T) {
	for filas := 45; filas <= 60; filas++ {
		previa := &domain.TablaPDF{Columnas: []string{"A", "B", "C"}}
		for i := 1; i <= filas; i++ {
			previa.Filas = append(previa.Filas, filaDePrueba(i))
		}
		doc := documentoConTabla(previa)
		doc.Secciones = append(doc.Secciones, domain.SeccionPDF{Titulo: "Sección final", Bloques: []domain.BloquePDF{
			{Parrafo: "Texto antes de la tabla."},
			{Tabla: &domain.TablaPDF{Columnas: []string{"Encabezado final"}, Filas: [][]string{{"primera-fila"}}}},
		}})
		textos := paginas(t, generar(t, doc))
		encabezados, veces := paginaDe(textos, "Encabezado final")
		primera, _ := paginaDe(textos, "primera-fila")
		if veces != 1 || encabezados != primera {
			t.Errorf("con %d filas previas, los encabezados están en la página %d (%d veces) y la primera fila en la %d",
				filas, encabezados+1, veces, primera+1)
		}
	}
}
