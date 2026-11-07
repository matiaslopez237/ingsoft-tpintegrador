package domain

import (
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// CA-27.2, CA-27.8, CA-27.9, CA-27.10: cada texto se normaliza carácter por carácter (RN-27.8) y se
// cuentan los "?" insertados (RN-27.9).
func TestNormalizarTextoPDF(t *testing.T) {
	casos := []struct {
		nombre      string
		texto       string
		normalizado string
		reemplazos  int
	}{
		// Paso 3: representables, quedan iguales (CA-27.2).
		{"acentos y signos del español", "Proyecto Ñandú: planificación ¿cumplida? ¡Sí!", "Proyecto Ñandú: planificación ¿cumplida? ¡Sí!", 0},
		{"lista de la spec", "áéíóúüñÁÉÍÓÚÜÑ€°ªº–—‘’“”•…←↑→↓≤≥%", "áéíóúüñÁÉÍÓÚÜÑ€°ªº–—‘’“”•…←↑→↓≤≥%", 0},
		{"latín extendido A", "Łódź Ŝ ő", "Łódź Ŝ ő", 0},
		{"signos de pregunta propios", "¿cumplida?", "¿cumplida?", 0},
		{"texto vacío", "", "", 0},
		// Paso 1: espacios de Unicode a espacio común (CA-27.10).
		{"tabulación", "Sprint\tfinal", "Sprint final", 0},
		{"salto de línea", "Sprint\nfinal", "Sprint final", 0},
		{"retorno de carro y salto", "Sprint\r\nfinal", "Sprint  final", 0},
		{"espacio duro", "Sprint final", "Sprint final", 0},
		{"espacio ideográfico", "Sprint　final", "Sprint final", 0},
		// Paso 2: invisibles de formato y selectores de variante, se quitan (CA-27.10).
		{"espacio de ancho cero", "Sprint​final", "Sprintfinal", 0},
		{"guion opcional", "pla­nificación", "planificación", 0},
		{"selector de variante solo", "a️b", "ab", 0},
		// Paso 4: el resto pasa a "?", y varios seguidos dan uno solo (CA-27.8, CA-27.9, CA-27.10).
		{"un emoji", "Sprint 😀 final", "Sprint ? final", 1},
		{"emojis seguidos e ideogramas", "Equipo 😀😀 y 中文", "Equipo ? y ?", 2},
		{"emojis unidos con ZWJ", "Familia 👨‍👩‍👧", "Familia ?", 1},
		{"emoji con selector de variante", "Amor ❤️", "Amor ?", 1},
		{"emojis separados por espacio", "😀 😀", "? ?", 2},
		{"carácter de control", "Sprint\u0007final", "Sprint?final", 1},
		{"byte que no es UTF-8 válido", "Sprint\xfffinal", "Sprint?final", 1},
		{"griego, aunque la fuente lo tenga", "αβγ", "?", 1},
		{"reemplazo al final", "Fin 🚀", "Fin ?", 1},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, n := normalizarTextoPDF(c.texto)
			if got != c.normalizado || n != c.reemplazos {
				t.Errorf("normalizarTextoPDF(%q) = (%q, %d), se esperaba (%q, %d)", c.texto, got, n, c.normalizado, c.reemplazos)
			}
		})
	}
}

// RN-27.8 se aplica a todos los textos: título, subtítulo, títulos de sección, párrafos, encabezados y
// celdas (CA-27.8). La cuenta es la suma de todos.
func TestNormalizarDocumento_TodosLosTextos(t *testing.T) {
	d := DocumentoPDF{
		Titulo:    "Reporte 😀",
		Subtitulo: "Avance 🚀",
		Secciones: []SeccionPDF{
			{Titulo: "Equipo 😀", Bloques: []BloquePDF{{Parrafo: "Listo ✅"}, {}}},
			{Titulo: "Estados", Bloques: []BloquePDF{{Tabla: &TablaPDF{
				Columnas: []string{"Estado 🔥"},
				Filas:    [][]string{{"Hecho 👍"}},
			}}}},
		},
	}
	got, n := NormalizarDocumento(d)
	want := DocumentoPDF{
		Titulo:    "Reporte ?",
		Subtitulo: "Avance ?",
		Secciones: []SeccionPDF{
			{Titulo: "Equipo ?", Bloques: []BloquePDF{{Parrafo: "Listo ?"}, {}}},
			{Titulo: "Estados", Bloques: []BloquePDF{{Tabla: &TablaPDF{
				Columnas: []string{"Estado ?"},
				Filas:    [][]string{{"Hecho ?"}},
			}}}},
		},
	}
	if n != 6 {
		t.Errorf("reemplazos = %d, se esperaban 6", n)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NormalizarDocumento() = %+v, se esperaba %+v", got, want)
	}
}

// Sección 5: NormalizarDocumento es pura, no modifica el documento que recibe.
func TestNormalizarDocumento_NoModificaElOriginal(t *testing.T) {
	d := DocumentoPDF{
		Titulo: "Reporte 😀",
		Secciones: []SeccionPDF{{Titulo: "S 😀", Bloques: []BloquePDF{{Tabla: &TablaPDF{
			Columnas: []string{"C 😀"}, Filas: [][]string{{"F 😀"}},
		}}}}},
	}
	NormalizarDocumento(d)
	tabla := d.Secciones[0].Bloques[0].Tabla
	if d.Titulo != "Reporte 😀" || d.Secciones[0].Titulo != "S 😀" || tabla.Columnas[0] != "C 😀" || tabla.Filas[0][0] != "F 😀" {
		t.Errorf("NormalizarDocumento modificó el original: %+v %+v", d, *tabla)
	}
}

// Sección 5: Emision y Aviso no son texto que normalizar; se conservan.
func TestNormalizarDocumento_ConservaEmisionYAviso(t *testing.T) {
	d := documentoValido()
	d.Emision = time.Date(2026, 11, 7, 14, 5, 0, 0, time.FixedZone("UTC-3", -3*60*60))
	d.Aviso = "Aviso previo"
	got, _ := NormalizarDocumento(d)
	if !got.Emision.Equal(d.Emision) || got.Aviso != d.Aviso {
		t.Errorf("NormalizarDocumento cambió la emisión o el aviso: %v %q", got.Emision, got.Aviso)
	}
}

// CA-27.8, CA-27.9: el texto del aviso, en singular y en plural; sin reemplazos, no hay aviso (RN-27.9).
func TestTextoAvisoReemplazos(t *testing.T) {
	casos := []struct {
		reemplazos int
		aviso      string
	}{
		{0, ""},
		{1, "Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?»."},
		{2, "Aviso: 2 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?»."},
		{5, "Aviso: 5 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?»."},
	}
	for _, c := range casos {
		if got := TextoAvisoReemplazos(c.reemplazos); got != c.aviso {
			t.Errorf("TextoAvisoReemplazos(%d) = %q, se esperaba %q", c.reemplazos, got, c.aviso)
		}
	}
}

// celdaNormalizada normaliza un documento con una sola celda y devuelve esa celda y los reemplazos.
func celdaNormalizada(t *testing.T, celda string) (string, int) {
	t.Helper()
	d := DocumentoPDF{Titulo: "T", Secciones: []SeccionPDF{{Titulo: "S", Bloques: []BloquePDF{{Tabla: &TablaPDF{
		Columnas: []string{"C"}, Filas: [][]string{{celda}},
	}}}}}}
	got, n := NormalizarDocumento(d)
	return got.Secciones[0].Bloques[0].Tabla.Filas[0][0], n
}

// CA-27.6: una celda tiene como máximo 300 caracteres (runes), contados después de RN-27.8; si tiene
// más, quedan los primeros 299 y "…" (RN-27.6). El recorte no es un reemplazo: no suma al aviso (RN-27.9).
func TestNormalizarDocumento_TopeDeCaracteresPorCelda(t *testing.T) {
	a := strings.Repeat
	casos := []struct {
		nombre     string
		celda      string
		want       string
		reemplazos int
	}{
		{"vacía", "", "", 0},
		{"299 caracteres", a("x", 299), a("x", 299), 0},
		{"300 caracteres, el límite", a("x", 300), a("x", 300), 0},
		{"301 caracteres", a("x", 301), a("x", 299) + "…", 0},
		{"400 caracteres", a("x", 400), a("x", 299) + "…", 0},
		{"300 acentos, que ocupan 600 bytes", a("á", 300), a("á", 300), 0},
		{"301 acentos", a("á", 301), a("á", 299) + "…", 0},
		{"301 eñes y espacios mezclados", a("ñ ", 150) + "ñ", a("ñ ", 149) + "ñ…", 0},
		{"un emoji y 299 letras: se reemplaza antes de contar", "😀" + a("x", 299), "?" + a("x", 299), 1},
		{"un emoji y 300 letras", "😀" + a("x", 300), "?" + a("x", 298) + "…", 1},
		{"350 emojis seguidos quedan en un solo ?", a("😀", 350), "?", 1},
		{"301 invisibles más 2 letras: se quitan antes de contar", a("​", 301) + "xy", "xy", 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, n := celdaNormalizada(t, c.celda)
			if got != c.want || n != c.reemplazos {
				t.Errorf("celda de %d runes = (%d runes %q…, %d), se esperaba (%d runes %q…, %d)",
					utf8.RuneCountInString(c.celda), utf8.RuneCountInString(got), primeros(got), n,
					utf8.RuneCountInString(c.want), primeros(c.want), c.reemplazos)
			}
			if utf8.RuneCountInString(got) > 300 {
				t.Errorf("la celda quedó con %d runes", utf8.RuneCountInString(got))
			}
		})
	}
}

// RN-27.6 es solo para las celdas: los encabezados, los párrafos y los títulos no tienen tope.
func TestNormalizarDocumento_TopeSoloEnLasCeldas(t *testing.T) {
	largo := strings.Repeat("x", 301)
	d := DocumentoPDF{Titulo: largo, Subtitulo: largo, Secciones: []SeccionPDF{
		{Titulo: largo, Bloques: []BloquePDF{{Parrafo: largo}, {Tabla: &TablaPDF{
			Columnas: []string{largo}, Filas: [][]string{{"corta"}},
		}}}},
	}}
	got, _ := NormalizarDocumento(d)
	s := got.Secciones[0]
	for nombre, texto := range map[string]string{
		"título": got.Titulo, "subtítulo": got.Subtitulo, "título de sección": s.Titulo,
		"párrafo": s.Bloques[0].Parrafo, "encabezado": s.Bloques[1].Tabla.Columnas[0],
	} {
		if texto != largo {
			t.Errorf("el %s se cortó: quedó con %d runes", nombre, utf8.RuneCountInString(texto))
		}
	}
}

// Sección 5: el recorte no modifica el documento original.
func TestNormalizarDocumento_TopeNoModificaElOriginal(t *testing.T) {
	largo := strings.Repeat("x", 301)
	filas := [][]string{{largo}}
	d := DocumentoPDF{Titulo: "T", Secciones: []SeccionPDF{{Titulo: "S", Bloques: []BloquePDF{{Tabla: &TablaPDF{
		Columnas: []string{"C"}, Filas: filas,
	}}}}}}
	NormalizarDocumento(d)
	if filas[0][0] != largo {
		t.Errorf("la celda original quedó con %d runes", utf8.RuneCountInString(filas[0][0]))
	}
}

// primeros acorta un texto largo para los mensajes de error.
func primeros(s string) string {
	r := []rune(s)
	if len(r) > 12 {
		return string(r[:12])
	}
	return s
}
