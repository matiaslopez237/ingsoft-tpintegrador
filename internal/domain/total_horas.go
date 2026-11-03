package domain

// HorasElemento son las horas de un elemento que entran en un total. Implementa HU-23 (sección 5).
type HorasElemento struct {
	HorasEstimadas *float64 // nil: sin estimar (RN-19.1)
	HorasReales    float64
}

// TotalHoras es el total de horas de un sprint o de un proyecto. Implementa HU-23 (RN-23.1).
type TotalHoras struct {
	HorasEstimadas        float64  // suma de los elementos estimados; 0 si no hay ninguno (RN-23.2, RN-23.6)
	HorasReales           float64  // suma de todos los elementos, estimados o no (RN-23.3)
	HorasRealesSinEstimar float64  // la parte de HorasReales que viene de elementos sin estimar (RN-23.5)
	DesviacionHoras       *float64 // nil: ningún elemento estimado (RN-23.6)
	DesviacionPorcentaje  *float64 // con un decimal; nil: ningún elemento estimado (RN-23.6)
	Elementos             int
	ElementosSinEstimar   int
}

// SinEstimar indica que no hay ningún elemento estimado, o ningún elemento, para que los adaptadores
// no repitan la regla (RN-23.6).
func (t TotalHoras) SinEstimar() bool {
	return t.Elementos == t.ElementosSinEstimar
}

// TotalizarHoras suma las horas de los elementos y calcula la desviación una sola vez, sobre las sumas
// de los elementos estimados (RN-23.4, RN-23.5). Sin ninguno estimado le pasa nil a CompararEsfuerzo, no
// 0, para no mostrar todo lo trabajado como desviación (RN-23.6). Implementa HU-23 (RN-23.2 a RN-23.6).
func TotalizarHoras(elementos []HorasElemento) TotalHoras {
	t := TotalHoras{Elementos: len(elementos)}
	for _, e := range elementos {
		t.HorasReales += e.HorasReales
		if e.HorasEstimadas == nil {
			t.ElementosSinEstimar++
			t.HorasRealesSinEstimar += e.HorasReales
			continue
		}
		t.HorasEstimadas += *e.HorasEstimadas
	}
	var estimadas *float64
	if !t.SinEstimar() {
		estimadas = &t.HorasEstimadas
	}
	c := CompararEsfuerzo(estimadas, t.HorasReales-t.HorasRealesSinEstimar)
	t.DesviacionHoras, t.DesviacionPorcentaje = c.DesviacionHoras, c.DesviacionPorcentaje
	return t
}
