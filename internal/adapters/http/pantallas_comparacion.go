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
	v := datosComparacion{HorasEstimadas: textoSinEstimar, HorasReales: textoHoras(c.HorasReales)}
	if !c.SinEstimar() {
		v.HorasEstimadas = textoHoras(*c.HorasEstimadas)
	}
	v.DesviacionHoras, v.DesviacionPorcentaje = textosDesviacion(c.DesviacionHoras, c.DesviacionPorcentaje)
	return v
}

// textosDesviacion da formato a las dos desviaciones, para la comparación de HU-19 y el total de HU-23.
// Sin desviación en horas no hay estimación (RN-19.7, RN-23.6); con ella y sin porcentaje, el estimado
// es 0 (RN-19.8).
func textosDesviacion(horas, porcentaje *float64) (enHoras, enPorcentaje string) {
	if horas == nil {
		return textoSinEstimar, textoSinEstimar
	}
	enHoras = conSigno(*horas, textoHoras(*horas))
	if porcentaje == nil {
		return enHoras, textoSinPorcentaje
	}
	return enHoras, conSigno(*porcentaje, textoPorcentaje(*porcentaje))
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
