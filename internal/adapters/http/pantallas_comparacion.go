package http

import (
	"strconv"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// textoSinEstimar reemplaza a la estimación y a las dos desviaciones de un elemento sin estimar.
const textoSinEstimar = "Sin estimar"

// textoSinPorcentaje es lo que se muestra si no hay porcentaje pero sí estimación (estimado 0, RN-19.8).
const textoSinPorcentaje = "No aplica"

// datosComparacion es el bloque "Esfuerzo estimado y real" de la pantalla del elemento (HU-19): solo
// texto ya formateado. Los números los calcula el dominio.
type datosComparacion struct {
	HorasEstimadas, HorasReales, DesviacionHoras, DesviacionPorcentaje string
}

// vistaComparacion da formato a la comparación. Implementa HU-19 (CA-19.18, sección 3).
func vistaComparacion(c domain.ComparacionEsfuerzo) datosComparacion {
	v := datosComparacion{
		HorasEstimadas: textoSinEstimar, HorasReales: textoHoras(c.HorasReales),
		DesviacionHoras: textoSinEstimar, DesviacionPorcentaje: textoSinEstimar,
	}
	if c.SinEstimar() {
		return v
	}
	v.HorasEstimadas = textoHoras(*c.HorasEstimadas)
	v.DesviacionHoras = conSigno(*c.DesviacionHoras, textoHoras(*c.DesviacionHoras))
	v.DesviacionPorcentaje = textoSinPorcentaje
	if c.DesviacionPorcentaje != nil {
		v.DesviacionPorcentaje = conSigno(*c.DesviacionPorcentaje, textoPorcentaje(*c.DesviacionPorcentaje))
	}
	return v
}

// textoHoras muestra las horas con coma decimal: 7,5 h.
func textoHoras(horas float64) string {
	return conComa(strconv.FormatFloat(horas, 'f', -1, 64)) + " h"
}

// textoPorcentaje muestra el porcentaje con coma y un decimal: 25,0 %.
func textoPorcentaje(porcentaje float64) string {
	return conComa(strconv.FormatFloat(porcentaje, 'f', 1, 64)) + " %"
}

func conComa(numero string) string {
	return strings.ReplaceAll(numero, ".", ",")
}

// conSigno agrega el + a una desviación positiva; las negativas ya traen el signo y el cero no lleva.
func conSigno(valor float64, texto string) string {
	if valor > 0 {
		return "+" + texto
	}
	return texto
}
