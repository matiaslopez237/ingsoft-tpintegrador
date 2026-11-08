package http

import (
	"strconv"
	"strings"
)

// Medidas de los gráficos de barras del Dashboard (HU-25), en unidades del viewBox del SVG. El SVG se
// estira al ancho de la pantalla con CSS, así que con muchos sprints las barras se ven más finas
// (P-2), pero nunca se superponen.
const (
	altoGrafico       = 220
	margenArriba      = 24 // lugar para el valor escrito sobre la barra más alta
	margenAbajo       = 28 // lugar para la etiqueta de cada grupo
	margenLateral     = 12
	anchoBarra        = 28
	separacionBarras  = 4
	separacionGrupos  = 24
	anchoMinimo       = 320
	distanciaTexto    = 6  // del valor a la punta de su barra
	distanciaEtiqueta = 18 // de la etiqueta del grupo al eje
)

// especGrafico son los datos de un gráfico de barras, tal como los arma la pantalla a partir del caso
// de uso. La geometría la calcula armarGraficoBarras.
type especGrafico struct {
	ID          string // prefijo de los id del SVG, único en la página
	Titulo      string
	Descripcion string
	// EncabezadoGrupos encabeza la columna de los grupos en la tabla de datos ("Sprint").
	EncabezadoGrupos string
	Series           []serieGrafico
	Grupos           []grupoGrafico
	Sufijo           string // se agrega al valor escrito sobre cada barra ("%")
	// MinimoEscala es el valor que como mínimo ocupa todo el alto: 100 en un porcentaje, para que 50%
	// se vea como la mitad aunque ningún sprint llegue a 100.
	MinimoEscala       float64
	Referencia         *float64 // línea horizontal opcional (la velocidad del equipo)
	EtiquetaReferencia string
}

// serieGrafico es una serie de barras: su nombre para la leyenda y la tabla, y la clase CSS de su color.
type serieGrafico struct {
	Nombre string
	Clase  string
}

// grupoGrafico es un grupo de barras (un sprint): un valor por serie, en el mismo orden. nil: sin datos.
type grupoGrafico struct {
	Etiqueta string
	Valores  []*float64
}

// graficoBarras es el gráfico listo para la plantilla: todo en coordenadas del viewBox.
type graficoBarras struct {
	especGrafico
	Ancho, Alto int
	BaseY       float64 // y del eje horizontal
	Grupos      []grupoDibujado
	Referencia  *lineaReferencia
}

type grupoDibujado struct {
	Etiqueta             string
	XEtiqueta, YEtiqueta float64
	Barras               []barraDibujada
}

type barraDibujada struct {
	X, Y, Ancho, Alto float64
	Clase             string
	Serie             string // para el título accesible de la barra
	Texto             string // el valor ya formateado, o "sin datos"
	XTexto, YTexto    float64
	SinDatos          bool
}

type lineaReferencia struct {
	X1, X2, Y float64
	Etiqueta  string
}

// armarGraficoBarras calcula la posición y el alto de cada barra. El mayor valor (o MinimoEscala, si es
// más grande) ocupa todo el alto útil; con todo en 0 la escala queda en 1, así no se divide por cero.
func armarGraficoBarras(e especGrafico) graficoBarras {
	maximo := e.MinimoEscala
	for _, gr := range e.Grupos {
		for _, v := range gr.Valores {
			if v != nil && *v > maximo {
				maximo = *v
			}
		}
	}
	if e.Referencia != nil && *e.Referencia > maximo {
		maximo = *e.Referencia
	}
	if maximo <= 0 {
		maximo = 1
	}

	baseY := float64(altoGrafico - margenAbajo)
	util := float64(altoGrafico - margenArriba - margenAbajo)
	alto := func(v float64) float64 { return v / maximo * util }

	anchoGrupo := float64(len(e.Series)*anchoBarra + (len(e.Series)-1)*separacionBarras)
	paso := anchoGrupo + separacionGrupos
	ancho := int(2*margenLateral + float64(len(e.Grupos))*paso)
	if ancho < anchoMinimo {
		ancho = anchoMinimo
	}

	g := graficoBarras{especGrafico: e, Ancho: ancho, Alto: altoGrafico, BaseY: baseY,
		Grupos: make([]grupoDibujado, len(e.Grupos))}
	for i, gr := range e.Grupos {
		x := margenLateral + separacionGrupos/2 + float64(i)*paso
		d := grupoDibujado{Etiqueta: gr.Etiqueta, XEtiqueta: x + anchoGrupo/2, YEtiqueta: baseY + distanciaEtiqueta}
		for j, v := range gr.Valores {
			b := barraDibujada{X: x + float64(j*(anchoBarra+separacionBarras)), Ancho: anchoBarra, Y: baseY}
			if j < len(e.Series) {
				b.Clase, b.Serie = e.Series[j].Clase, e.Series[j].Nombre
			}
			if v == nil {
				b.SinDatos, b.Texto = true, "sin datos"
			} else {
				b.Alto = alto(*v)
				b.Y = baseY - b.Alto
				b.Texto = formatearValorGrafico(*v) + e.Sufijo
			}
			b.XTexto, b.YTexto = b.X+anchoBarra/2, b.Y-distanciaTexto
			d.Barras = append(d.Barras, b)
		}
		g.Grupos[i] = d
	}
	if e.Referencia != nil {
		g.Referencia = &lineaReferencia{
			X1: margenLateral, X2: float64(ancho - margenLateral), Y: baseY - alto(*e.Referencia),
			Etiqueta: e.EtiquetaReferencia + ": " + formatearValorGrafico(*e.Referencia),
		}
	}
	return g
}

// formatearValorGrafico escribe el número con un decimal y coma, sin ",0" de más: 66,7 y 100.
func formatearValorGrafico(v float64) string {
	s := strconv.FormatFloat(v, 'f', 1, 64)
	s = strings.TrimSuffix(s, ".0")
	return strings.Replace(s, ".", ",", 1)
}
