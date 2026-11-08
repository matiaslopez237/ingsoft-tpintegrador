package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// ErrDocumentoInvalido indica que quien armó el documento del PDF (la conversión del reporte, spec HU-27,
// P-11) no cumplió RN-27.11. Es un error de programación, no del usuario: el adaptador HTTP lo responde
// como falla interna. Implementa HU-27 (RN-27.11).
var ErrDocumentoInvalido = errors.New("documento del PDF inválido")

// Topes de una tabla (RN-27.6, P-10): con ellos, la fila de encabezados y una fila con celdas de
// LargoMaximoCelda caracteres entran juntas en una página aun en el peor caso.
const (
	MaximoColumnas        = 6
	LargoMaximoEncabezado = 100 // en caracteres (runes), contados después de RN-27.8
)

// DocumentoPDF es el modelo neutral del que se genera el PDF (spec HU-27, sección 2). La conversión del
// reporte de HU-26 (RN-27.16, P-11) completa todos los campos menos Emision y Aviso, que pone el caso de
// uso de exportación (RN-27.9, RN-27.10). Implementa HU-27 (RN-27.1).
type DocumentoPDF struct {
	Titulo    string    // obligatorio
	Subtitulo string    // opcional
	Emision   time.Time // fecha y hora de emisión, con su huso horario (RN-27.10)
	Aviso     string    // aviso de reemplazos (RN-27.9); HU-26 nunca lo completa ni lo lee
	Secciones []SeccionPDF
}

// SeccionPDF tiene un título obligatorio y sus bloques, en el orden en que se muestran.
type SeccionPDF struct {
	Titulo  string
	Bloques []BloquePDF
}

// BloquePDF es un párrafo o una tabla, no los dos (RN-27.11). Sin ninguno de los dos es un párrafo
// vacío: no es un error y no se muestra.
type BloquePDF struct {
	Parrafo string
	Tabla   *TablaPDF
}

// TablaPDF tiene entre 1 y MaximoColumnas columnas, y cada fila tantas celdas como columnas (RN-27.11).
type TablaPDF struct {
	Columnas []string
	Filas    [][]string
}

// Validar implementa RN-27.11. Los títulos se comprueban sin contar espacios; se llama después de
// normalizar (RN-27.8), así que un título hecho solo de invisibles llega vacío. Implementa HU-27
// (RN-27.6, RN-27.11).
func (d DocumentoPDF) Validar() error {
	if strings.TrimSpace(d.Titulo) == "" {
		return ErrDocumentoInvalido
	}
	for _, s := range d.Secciones {
		if strings.TrimSpace(s.Titulo) == "" {
			return ErrDocumentoInvalido
		}
		for _, b := range s.Bloques {
			if b.Tabla == nil {
				continue // párrafo, con texto o vacío
			}
			if b.Parrafo != "" || len(b.Tabla.Columnas) == 0 || len(b.Tabla.Columnas) > MaximoColumnas {
				return ErrDocumentoInvalido
			}
			for _, col := range b.Tabla.Columnas {
				if utf8.RuneCountInString(col) > LargoMaximoEncabezado {
					return ErrDocumentoInvalido
				}
			}
			for _, fila := range b.Tabla.Filas {
				if len(fila) != len(b.Tabla.Columnas) {
					return ErrDocumentoInvalido
				}
			}
		}
	}
	return nil
}
