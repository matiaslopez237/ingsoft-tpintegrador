// Package pdf genera el PDF del reporte (HU-27) con codeberg.org/go-pdf/fpdf (ADR-0006). Es el único
// paquete que importa fpdf: el dominio y los casos de uso no lo conocen.
package pdf

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"

	"codeberg.org/go-pdf/fpdf"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// DejaVu Sans Condensed 2.37, regular y negrita (ADR-0006). La licencia está en
// docs/third-party-licenses/dejavu-fonts-LICENSE.txt.
var (
	//go:embed fuentes/DejaVuSansCondensed.ttf
	fuenteRegular []byte
	//go:embed fuentes/DejaVuSansCondensed-Bold.ttf
	fuenteNegrita []byte
)

const familia = "DejaVu"

// Generador implementa app.GeneradorPDF. No guarda estado: cada llamada arma su propio documento.
type Generador struct{}

func NuevoGenerador() *Generador { return &Generador{} }

// Generar dibuja el documento en A4 vertical (RN-27.14). Recibe el documento ya normalizado y validado
// por el caso de uso, así que no lo vuelve a normalizar ni a validar. Si fpdf falla, devuelve su error
// envuelto y ningún byte.
func (g *Generador) Generar(ctx context.Context, doc domain.DocumentoPDF) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f := fpdf.New("P", "mm", "A4", "")
	f.AddUTF8FontFromBytes(familia, "", fuenteRegular)
	f.AddUTF8FontFromBytes(familia, "B", fuenteNegrita)
	f.AddPage()

	f.SetFont(familia, "B", 16)
	f.MultiCell(0, 8, doc.Titulo, "", "L", false)
	if doc.Subtitulo != "" {
		f.SetFont(familia, "", 11)
		f.MultiCell(0, 6, doc.Subtitulo, "", "L", false)
	}
	f.SetFont(familia, "", 9)
	f.MultiCell(0, 5, textoEmision(doc), "", "L", false)
	if doc.Aviso != "" {
		f.MultiCell(0, 5, doc.Aviso, "", "L", false)
	}
	for _, s := range doc.Secciones {
		f.Ln(4)
		f.SetFont(familia, "B", 12)
		f.MultiCell(0, 7, s.Titulo, "", "L", false)
		for _, b := range s.Bloques {
			if b.Tabla != nil {
				dibujarTabla(f, b.Tabla)
				continue
			}
			if b.Parrafo != "" {
				f.SetFont(familia, "", 10)
				f.MultiCell(0, 5, b.Parrafo, "", "L", false)
			}
		}
	}

	var buf bytes.Buffer
	if err := f.Output(&buf); err != nil {
		return nil, fmt.Errorf("generar el PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// textoEmision es la línea de emisión de la sección 3 de la spec (RN-27.10). La emisión ya viene en
// UTC-3 desde el caso de uso.
func textoEmision(doc domain.DocumentoPDF) string {
	return "Emitido el " + doc.Emision.Format("02/01/2006") + " a las " + doc.Emision.Format("15:04") +
		" (hora de Argentina, UTC-3)."
}

func dibujarTabla(f *fpdf.Fpdf, t *domain.TablaPDF) {
	ancho, _ := f.GetPageSize()
	izq, _, der, _ := f.GetMargins()
	w := (ancho - izq - der) / float64(len(t.Columnas))
	f.SetFont(familia, "B", 9)
	for _, c := range t.Columnas {
		f.CellFormat(w, 6, c, "1", 0, "L", false, 0, "")
	}
	f.Ln(-1)
	f.SetFont(familia, "", 9)
	for _, fila := range t.Filas {
		for _, celda := range fila {
			f.CellFormat(w, 6, celda, "1", 0, "L", false, 0, "")
		}
		f.Ln(-1)
	}
}
