package pdf

import (
	"bytes"
	"context"
	"errors"
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
