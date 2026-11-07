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
	f.SetMargins(margenLateral, margenSuperior, margenLateral)
	f.SetAutoPageBreak(true, margenInferior)
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

// Márgenes de la hoja, en mm. El inferior deja lugar para el pie.
const (
	margenLateral  = 15.0
	margenSuperior = 15.0
	margenInferior = 20.0
)

// Medidas de la tabla, en mm.
const (
	altoLineaTabla = 4.5 // con letra de 9 puntos
	letraTabla     = 9
)

// dibujarTabla dibuja la fila de encabezados y las filas, todas las columnas del mismo ancho (el ancho
// no es regla de negocio; spec, sección 5). El texto de cada celda se ajusta en varias líneas, y una fila
// que no entra en lo que queda de la página pasa entera a la siguiente (RN-27.5). Sin filas, muestra
// "Sin datos." (RN-27.7).
func dibujarTabla(f *fpdf.Fpdf, t *domain.TablaPDF) {
	ancho, _ := f.GetPageSize()
	izq, _, der, _ := f.GetMargins()
	w := (ancho - izq - der) / float64(len(t.Columnas))
	dibujarFila(f, t.Columnas, w, "B")
	if len(t.Filas) == 0 {
		f.SetFont(familia, "", letraTabla)
		f.MultiCell(0, altoLineaTabla+1, "Sin datos.", "", "L", false)
		return
	}
	for _, fila := range t.Filas {
		dibujarFila(f, fila, w, "")
	}
}

// dibujarFila dibuja una fila con la altura de su celda más alta. Si no entra, salta de página antes.
func dibujarFila(f *fpdf.Fpdf, celdas []string, w float64, estilo string) {
	f.SetFont(familia, estilo, letraTabla)
	alto := altoFila(f, celdas, w)
	if f.GetY()+alto > limiteInferior(f) {
		f.AddPage()
	}
	x0, y0 := f.GetXY()
	for i, c := range celdas {
		x := x0 + float64(i)*w
		f.Rect(x, y0, w, alto, "D")
		f.SetXY(x, y0)
		f.MultiCell(w, altoLineaTabla, c, "", "L", false)
	}
	f.SetXY(x0, y0+alto)
}

// altoFila es la cantidad de líneas de la celda más alta por el alto de una línea, con la letra actual.
func altoFila(f *fpdf.Fpdf, celdas []string, w float64) float64 {
	lineas := 1
	for _, c := range celdas {
		lineas = max(lineas, len(f.SplitText(c, w)))
	}
	return float64(lineas) * altoLineaTabla
}

// limiteInferior es la altura desde la que fpdf salta de página solo.
func limiteInferior(f *fpdf.Fpdf) float64 {
	_, alto := f.GetPageSize()
	_, margen := f.GetAutoPageBreak()
	return alto - margen
}
