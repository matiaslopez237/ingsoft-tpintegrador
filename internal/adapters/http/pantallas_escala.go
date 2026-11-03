package http

import (
	"slices"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// opcionStoryPoints es una opción de la lista de story points de las pantallas de alta y edición.
type opcionStoryPoints struct {
	Valor, Etiqueta string
	Elegida         bool
}

// opcionesStoryPoints arma la lista de story points: "Sin estimar" y la escala del dominio, con
// actual (el valor del formulario, "" si no tiene) elegido. Si actual no es de la escala (un valor
// cargado antes de HU-13), se agrega marcado, así guardar sin tocarlo no lo cambia.
// Implementa HU-13 (CA-13.7, RN-13.3, RN-13.4).
func opcionesStoryPoints(actual string) []opcionStoryPoints {
	escala := domain.EscalaStoryPoints()
	opciones := []opcionStoryPoints{{Valor: "", Etiqueta: "Sin estimar", Elegida: actual == ""}}
	n, err := strconv.Atoi(actual)
	fuera := actual != "" && (err != nil || !slices.Contains(escala, n))
	if fuera {
		opciones = append(opciones, opcionStoryPoints{Valor: actual, Etiqueta: actual + " (fuera de la escala)", Elegida: true})
	}
	for _, v := range escala {
		s := strconv.Itoa(v)
		opciones = append(opciones, opcionStoryPoints{Valor: s, Etiqueta: s, Elegida: !fuera && s == actual})
	}
	return opciones
}
