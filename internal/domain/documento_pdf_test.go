package domain

import (
	"errors"
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

// CA-27.11: los errores de quien arma el documento se rechazan con ErrDocumentoInvalido (RN-27.11).
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

// CA-27.7 y RN-27.11: lo vacío que no es un error de armado se acepta.
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
