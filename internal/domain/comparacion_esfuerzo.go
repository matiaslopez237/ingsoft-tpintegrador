package domain

import "math"

// ComparacionEsfuerzo es el resultado de comparar las horas estimadas con las reales.
// Implementa HU-19 (RN-19.5 a RN-19.9).
type ComparacionEsfuerzo struct {
	HorasEstimadas       *float64 // nil: sin estimar (RN-19.7)
	HorasReales          float64
	DesviacionHoras      *float64 // reales - estimadas; nil: sin estimar
	DesviacionPorcentaje *float64 // con un decimal; nil: sin estimar o con estimado 0 (RN-19.8)
}

// SinEstimar indica que no hay estimación contra la que comparar, para que los adaptadores no
// repitan la regla (RN-19.7). Un estimado 0 no es "sin estimar" (RN-19.8).
func (c ComparacionEsfuerzo) SinEstimar() bool {
	return c.HorasEstimadas == nil
}

// CompararEsfuerzo calcula la desviación entre las horas estimadas y las reales. Recibe números y no
// entidades para que HU-23 la use con totales (spec, sección 5). Los valores llegan ya validados, así
// que no devuelve error. Implementa HU-19 (RN-19.5 a RN-19.9).
func CompararEsfuerzo(estimadas *float64, reales float64) ComparacionEsfuerzo {
	c := ComparacionEsfuerzo{HorasReales: reales}
	if estimadas == nil {
		return c
	}
	estimado := *estimadas // copia: el resultado no comparte el puntero recibido
	desviacion := reales - estimado
	c.HorasEstimadas = &estimado
	c.DesviacionHoras = &desviacion
	if estimado != 0 { // RN-19.8: no se divide por cero
		porcentaje := redondearUnDecimal(desviacion * 100 / estimado) // se multiplica antes de dividir (spec, sección 5)
		c.DesviacionPorcentaje = &porcentaje
	}
	return c
}

// redondearUnDecimal redondea con los medios alejándose de cero, que es lo que hace math.Round.
// Implementa RN-19.6.
func redondearUnDecimal(v float64) float64 {
	r := math.Round(v*10) / 10
	if r == 0 { // también es cierto para -0, y RN-19.6 pide 0, nunca -0
		return 0
	}
	return r
}
