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

// Títulos de sección y párrafos.
const (
	letraSeccion     = 12
	altoLineaSeccion = 7.0
	espacioSeccion   = 4.0 // antes del título de cada sección
	letraParrafo     = 10
	altoLineaParrafo = 5.0
)

// aliasPaginas es el texto que fpdf reemplaza por la cantidad total de páginas al cerrar el documento.
const aliasPaginas = "{nb}"

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
	// CA-27.12: la única fecha es la emisión, y el catálogo sale ordenado: mismos datos, mismos bytes.
	// No se usan las funciones SetDefault*, que cambian estado global de fpdf (ADR-0006).
	f.SetCreationDate(doc.Emision)
	f.SetModificationDate(doc.Emision)
	f.SetCatalogSort(true)
	f.AddUTF8FontFromBytes(familia, "", fuenteRegular)
	f.AddUTF8FontFromBytes(familia, "B", fuenteNegrita)
	f.AliasNbPages(aliasPaginas)
	f.SetHeaderFunc(func() { encabezado(f, doc.Titulo) })
	f.SetFooterFunc(func() { pie(f) })
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
		dibujarSeccion(f, s)
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

// encabezado escribe el título del documento arriba de cada página, desde la segunda (RN-27.3).
func encabezado(f *fpdf.Fpdf, titulo string) {
	if f.PageNo() == 1 {
		return
	}
	f.SetFont(familia, "", 8)
	f.SetTextColor(90, 90, 90)
	f.CellFormat(0, 5, titulo, "B", 1, "L", false, 0, "")
	f.SetTextColor(0, 0, 0)
	f.Ln(3)
}

// pie escribe "Página k de N" abajo de cada página (RN-27.3).
func pie(f *fpdf.Fpdf) {
	f.SetY(-margenInferior + 5)
	f.SetFont(familia, "", 8)
	f.SetTextColor(90, 90, 90)
	f.CellFormat(0, 5, fmt.Sprintf("Página %d de %s", f.PageNo(), aliasPaginas), "", 0, "C", false, 0, "")
	f.SetTextColor(0, 0, 0)
}

// dibujarSeccion escribe el título y los bloques en orden. El título no queda solo al pie: si no entra
// junto con el comienzo de su primer bloque visible, pasa a la página siguiente (RN-27.4).
func dibujarSeccion(f *fpdf.Fpdf, s domain.SeccionPDF) {
	ancho := anchoUtil(f)
	f.SetFont(familia, "B", letraSeccion)
	alto := espacioSeccion + float64(len(f.SplitText(s.Titulo, ancho)))*altoLineaSeccion + altoComienzo(f, s.Bloques)
	if f.GetY()+alto > limiteInferior(f) {
		f.AddPage()
	} else {
		f.Ln(espacioSeccion)
	}
	f.SetFont(familia, "B", letraSeccion)
	f.MultiCell(0, altoLineaSeccion, s.Titulo, "", "L", false)
	for _, b := range s.Bloques {
		switch {
		case b.Tabla != nil:
			dibujarTabla(f, b.Tabla)
		case b.Parrafo != "": // un párrafo vacío no se muestra (RN-27.11)
			f.SetFont(familia, "", letraParrafo)
			f.MultiCell(0, altoLineaParrafo, b.Parrafo, "", "L", false)
		}
	}
}

// altoComienzo es lo que tiene que entrar junto con el título de la sección: la primera línea del
// primer párrafo con texto, o los encabezados y la primera fila de la primera tabla.
func altoComienzo(f *fpdf.Fpdf, bloques []domain.BloquePDF) float64 {
	for _, b := range bloques {
		switch {
		case b.Tabla != nil:
			return altoComienzoTabla(f, b.Tabla)
		case b.Parrafo != "":
			return altoLineaParrafo
		}
	}
	return 0
}

// dibujarTabla dibuja la fila de encabezados y las filas, todas las columnas del mismo ancho (el ancho
// no es regla de negocio; spec, sección 5). El texto de cada celda se ajusta en varias líneas; una fila
// que no entra en lo que queda de la página pasa entera a la siguiente, donde se repiten los encabezados
// (RN-27.5). Sin filas, muestra "Sin datos." (RN-27.7).
func dibujarTabla(f *fpdf.Fpdf, t *domain.TablaPDF) {
	w := anchoUtil(f) / float64(len(t.Columnas))
	if f.GetY()+altoComienzoTabla(f, t) > limiteInferior(f) {
		f.AddPage()
	}
	dibujarFila(f, t.Columnas, w, "B")
	if len(t.Filas) == 0 {
		f.SetFont(familia, "", letraTabla)
		f.MultiCell(0, altoLineaTabla+1, "Sin datos.", "", "L", false)
		return
	}
	for _, fila := range t.Filas {
		f.SetFont(familia, "", letraTabla)
		if f.GetY()+altoFila(f, fila, w) > limiteInferior(f) {
			f.AddPage()
			dibujarFila(f, t.Columnas, w, "B")
		}
		dibujarFila(f, fila, w, "")
	}
}

// altoComienzoTabla es el alto de la fila de encabezados más el de la primera fila (o de "Sin datos.").
func altoComienzoTabla(f *fpdf.Fpdf, t *domain.TablaPDF) float64 {
	w := anchoUtil(f) / float64(len(t.Columnas))
	f.SetFont(familia, "B", letraTabla)
	alto := altoFila(f, t.Columnas, w)
	f.SetFont(familia, "", letraTabla)
	if len(t.Filas) == 0 {
		return alto + altoLineaTabla + 1
	}
	return alto + altoFila(f, t.Filas[0], w)
}

// dibujarFila dibuja una fila con la altura de su celda más alta, desde la posición actual.
func dibujarFila(f *fpdf.Fpdf, celdas []string, w float64, estilo string) {
	f.SetFont(familia, estilo, letraTabla)
	alto := altoFila(f, celdas, w)
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

func anchoUtil(f *fpdf.Fpdf) float64 {
	ancho, _ := f.GetPageSize()
	izq, _, der, _ := f.GetMargins()
	return ancho - izq - der
}

// limiteInferior es la altura desde la que fpdf salta de página solo.
func limiteInferior(f *fpdf.Fpdf) float64 {
	_, alto := f.GetPageSize()
	_, margen := f.GetAutoPageBreak()
	return alto - margen
}
