package domain

import (
	"cmp"
	"slices"
)

// HistorialRondas devuelve el resumen de cada ronda de un elemento, de la más reciente (número
// mayor) a la más antigua. Usa Ronda.Resumen, así la ronda abierta no muestra valores (RN-14.3) y
// las cerradas traen la revelación de HU-15. No modifica la lista recibida.
// Implementa HU-16 (RN-16.3).
func HistorialRondas(rondas []Ronda) []ResumenRonda {
	historial := make([]ResumenRonda, 0, len(rondas))
	for _, r := range rondas {
		historial = append(historial, r.Resumen())
	}
	slices.SortFunc(historial, func(a, b ResumenRonda) int { return cmp.Compare(b.Numero, a.Numero) })
	return historial
}

// PuedeRepetirse indica si la pantalla de la ronda ofrece "Nueva ronda": la ronda está Cerrada, es
// la última del elemento (ultima, nil si no se conoce) y el elemento no está Hecho. No mira si hubo
// diferencia significativa: repetir es decisión del equipo. Implementa HU-16 (RN-16.4).
func PuedeRepetirse(ronda ResumenRonda, ultima *ResumenRonda, e ElementoBacklog) bool {
	return ronda.Estado == RondaCerrada && ultima != nil && ultima.ID == ronda.ID && e.Estado != EstadoHecho
}
