package domain

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// NormalizarDocumento aplica RN-27.8 a todos los textos del documento y devuelve cuántos "?" insertó
// (RN-27.9). Es pura: no modifica el documento que recibe. Emision y Aviso no se tocan.
func NormalizarDocumento(d DocumentoPDF) (DocumentoPDF, int) {
	total := 0
	norm := func(s string) string {
		r, n := normalizarTextoPDF(s)
		total += n
		return r
	}
	c := d
	c.Titulo = norm(d.Titulo)
	c.Subtitulo = norm(d.Subtitulo)
	c.Secciones = nil
	for _, s := range d.Secciones {
		cs := SeccionPDF{Titulo: norm(s.Titulo)}
		for _, b := range s.Bloques {
			cb := BloquePDF{Parrafo: norm(b.Parrafo)}
			if b.Tabla != nil {
				t := TablaPDF{}
				for _, col := range b.Tabla.Columnas {
					t.Columnas = append(t.Columnas, norm(col))
				}
				for _, fila := range b.Tabla.Filas {
					var celdas []string
					for _, celda := range fila {
						celdas = append(celdas, norm(celda))
					}
					t.Filas = append(t.Filas, celdas)
				}
				cb.Tabla = &t
			}
			cs.Bloques = append(cs.Bloques, cb)
		}
		c.Secciones = append(c.Secciones, cs)
	}
	return c, total
}

// TextoAvisoReemplazos es el aviso de RN-27.9 para n "?" insertados; vacío si no hubo reemplazos.
func TextoAvisoReemplazos(n int) string {
	switch {
	case n <= 0:
		return ""
	case n == 1:
		return "Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?»."
	default:
		return fmt.Sprintf("Aviso: %d símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?».", n)
	}
}

// normalizarTextoPDF implementa los cuatro pasos de RN-27.8 sobre un texto y devuelve cuántos "?"
// insertó. Un byte que no es UTF-8 válido se trata como un carácter no representable.
func normalizarTextoPDF(s string) (string, int) {
	var b strings.Builder
	reemplazos := 0
	enReemplazo := false // el último carácter escrito fue un "?" de reemplazo
	for i := 0; i < len(s); {
		r, tam := utf8.DecodeRuneInString(s[i:])
		i += tam
		switch {
		case r == utf8.RuneError && tam == 1: // byte inválido: paso 4
		case unicode.IsSpace(r): // paso 1
			b.WriteByte(' ')
			enReemplazo = false
			continue
		case invisibleDeFormato(r): // paso 2: no corta una racha de "?"
			continue
		case representableEnPDF(r): // paso 3
			b.WriteRune(r)
			enReemplazo = false
			continue
		}
		if !enReemplazo { // paso 4: varios seguidos dan un solo "?"
			b.WriteByte('?')
			reemplazos++
			enReemplazo = true
		}
	}
	return b.String(), reemplazos
}

// invisibleDeFormato: categoría Unicode Cf y selectores de variante (RN-27.8, paso 2).
func invisibleDeFormato(r rune) bool {
	return unicode.Is(unicode.Cf, r) ||
		(r >= 0xFE00 && r <= 0xFE0F) || // selectores de variante
		(r >= 0xE0100 && r <= 0xE01EF) // selectores de variante suplementarios
}

// representableEnPDF es la lista corta de RN-27.8, paso 3, probada contra DejaVu Sans Condensed
// (ADR-0006).
func representableEnPDF(r rune) bool {
	switch {
	case r >= 0x21 && r <= 0x7E: // ASCII imprimible
		return true
	case r >= 0xA1 && r <= 0x17F && r != 0xAD: // Latin-1 y Latín extendido A, sin el guion opcional
		return true
	}
	switch r {
	case '–', '—', '‘', '’', '“', '”', '•', '…', '€', '←', '↑', '→', '↓', '≤', '≥':
		return true
	}
	return false
}
