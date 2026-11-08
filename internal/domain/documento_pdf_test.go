package domain

import (
	"errors"
	"strings"
	"testing"
)

// documentoValido arma un documento que cumple RN-27.11, para que cada caso cambie una sola cosa.
func documentoValido() DocumentoPDF {
	return DocumentoPDF{
		Titulo: "Reporte del proyecto Alfa",
		Secciones: []SeccionPDF{
			{Titulo: "Resumen", Bloques: []BloquePDF{{Parrafo: "El proyecto está Activo."}}},
			{Titulo: "Historias", Bloques: []BloquePDF{{Tabla: &TablaPDF{
				Columnas: []string{"Historia", "Puntos"},
				Filas:    [][]string{{"Login", "5"}},
			}}}},
		},
	}
}

// tablaDeColumnas arma una tabla de n columnas con una fila, todas con el mismo texto.
func tablaDeColumnas(n int, texto string) *TablaPDF {
	t := &TablaPDF{Filas: [][]string{{}}}
	for range n {
		t.Columnas = append(t.Columnas, texto)
		t.Filas[0] = append(t.Filas[0], texto)
	}
	return t
}

// Traza: HU-27 CA-27.11 (RN-27.11, P-10): los errores de quien arma el documento se rechazan con
// ErrDocumentoInvalido.
func TestDocumentoPDFValidar_DocumentoMalArmado_RetornaErrDocumentoInvalido(t *testing.T) {
	casos := []struct {
		nombre  string
		cambiar func(d *DocumentoPDF)
	}{
		{"título vacío", func(d *DocumentoPDF) { d.Titulo = "" }},
		{"título solo de espacios", func(d *DocumentoPDF) { d.Titulo = "   " }},
		{"sección sin título", func(d *DocumentoPDF) { d.Secciones[0].Titulo = "" }},
		{"sección con título solo de espacios", func(d *DocumentoPDF) { d.Secciones[0].Titulo = "  " }},
		{"bloque con párrafo y tabla", func(d *DocumentoPDF) {
			d.Secciones[0].Bloques[0].Tabla = &TablaPDF{Columnas: []string{"A"}}
		}},
		{"tabla sin columnas", func(d *DocumentoPDF) { d.Secciones[1].Bloques[0].Tabla.Columnas = nil }},
		{"fila con menos celdas que columnas", func(d *DocumentoPDF) {
			d.Secciones[1].Bloques[0].Tabla.Filas = [][]string{{"Login"}}
		}},
		{"fila con más celdas que columnas", func(d *DocumentoPDF) {
			d.Secciones[1].Bloques[0].Tabla.Filas = [][]string{{"Login", "5", "extra"}}
		}},
		// RN-27.6, P-10: los topes de columnas y de encabezados.
		{"tabla de 7 columnas", func(d *DocumentoPDF) { d.Secciones[1].Bloques[0].Tabla = tablaDeColumnas(7, "A") }},
		{"encabezado de 101 caracteres", func(d *DocumentoPDF) {
			d.Secciones[1].Bloques[0].Tabla.Columnas[0] = strings.Repeat("Œ", 101)
		}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			d := documentoValido()
			c.cambiar(&d)
			if err := d.Validar(); !errors.Is(err, ErrDocumentoInvalido) {
				t.Errorf("Validar() = %v, se esperaba ErrDocumentoInvalido", err)
			}
		})
	}
}

// Traza: HU-27 CA-27.7, CA-27.11 (RN-27.6, RN-27.11, P-10): lo vacío que no es un error de armado se
// acepta, y también una tabla justo en los topes de columnas y de encabezados.
func TestDocumentoPDFValidar_ContenidoVacioPermitido_NoEsError(t *testing.T) {
	casos := []struct {
		nombre  string
		cambiar func(d *DocumentoPDF)
	}{
		{"documento completo", func(d *DocumentoPDF) {}},
		{"sin subtítulo ni secciones", func(d *DocumentoPDF) { d.Secciones = nil }},
		{"sección sin bloques", func(d *DocumentoPDF) { d.Secciones[0].Bloques = nil }},
		{"tabla sin filas", func(d *DocumentoPDF) { d.Secciones[1].Bloques[0].Tabla.Filas = nil }},
		{"bloque sin párrafo ni tabla (párrafo vacío)", func(d *DocumentoPDF) {
			d.Secciones[0].Bloques = []BloquePDF{{}}
		}},
		// RN-27.6, P-10: justo en los topes. Œ ocupa 2 bytes: el encabezado se cuenta en caracteres.
		{"tabla de 6 columnas", func(d *DocumentoPDF) { d.Secciones[1].Bloques[0].Tabla = tablaDeColumnas(6, "A") }},
		{"encabezado de 100 caracteres", func(d *DocumentoPDF) {
			d.Secciones[1].Bloques[0].Tabla.Columnas[0] = strings.Repeat("Œ", 100)
		}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			d := documentoValido()
			c.cambiar(&d)
			if err := d.Validar(); err != nil {
				t.Errorf("Validar() = %v, se esperaba nil", err)
			}
		})
	}
}
