package domain

import (
	"cmp"
	"errors"
	"slices"
)

// Órdenes del listado del backlog (HU-07, RN-07.2).
const (
	OrdenDescendente = "desc" // Alta primero: el orden por defecto
	OrdenAscendente  = "asc"  // Baja primero
)

// Errores de HU-07 (spec, sección 7).
var (
	ErrOrdenInvalido          = errors.New("el orden debe ser asc o desc")
	ErrEstadoElementoInvalido = errors.New("el estado debe ser Pendiente, En progreso o Hecho")
)

// rangoPrioridad ubica cada prioridad de mayor a menor: Alta va antes que Media y Media antes que Baja.
var rangoPrioridad = map[Prioridad]int{PrioridadAlta: 0, PrioridadMedia: 1, PrioridadBaja: 2}

// OrdenarBacklog devuelve los elementos ordenados por prioridad (descendente por defecto) y, a
// igual prioridad, por orden de creación (ID). Si estado no está vacío, deja solo los de ese
// estado. No modifica la lista recibida. Orden y estado vacíos toman los valores por defecto.
// Implementa HU-07 (RN-07.2 a RN-07.5).
func OrdenarBacklog(elementos []ElementoBacklog, orden, estado string) ([]ElementoBacklog, error) {
	switch orden {
	case "", OrdenDescendente, OrdenAscendente:
	default:
		return nil, ErrOrdenInvalido
	}
	if estado != "" && !estadoElementoValido(EstadoElemento(estado)) {
		return nil, ErrEstadoElementoInvalido
	}

	lista := make([]ElementoBacklog, 0, len(elementos))
	for _, e := range elementos {
		if estado == "" || e.Estado == EstadoElemento(estado) {
			lista = append(lista, e)
		}
	}
	slices.SortFunc(lista, func(a, b ElementoBacklog) int {
		porPrioridad := cmp.Compare(rangoPrioridad[a.Prioridad], rangoPrioridad[b.Prioridad])
		if orden == OrdenAscendente {
			porPrioridad = -porPrioridad
		}
		return cmp.Or(porPrioridad, cmp.Compare(a.ID, b.ID)) // RN-07.3: el desempate no se invierte
	})
	return lista, nil
}

func estadoElementoValido(e EstadoElemento) bool {
	switch e {
	case EstadoPendiente, EstadoEnProgreso, EstadoHecho:
		return true
	}
	return false
}
