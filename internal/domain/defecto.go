package domain

import (
	"errors"
	"strings"
)

// Severidad de un defecto. Implementa HU-20 (RN-20.3).
type Severidad string

const (
	SeveridadCritica Severidad = "Crítica"
	SeveridadAlta    Severidad = "Alta"
	SeveridadMedia   Severidad = "Media"
	SeveridadBaja    Severidad = "Baja"
)

// EstadoDefecto es el estado de un defecto. Implementa HU-20 (RN-20.4). Por ahora el único valor
// posible es Abierto: cambiarlo de estado queda para una historia futura de la Épica 6.
type EstadoDefecto string

const (
	DefectoAbierto EstadoDefecto = "Abierto"
)

// Errores de HU-20 (spec, sección 7).
var (
	ErrDescripcionDefectoVacia = errors.New("la descripción es obligatoria")
	ErrSeveridadInvalida       = errors.New("la severidad debe ser Crítica, Alta, Media o Baja")
)

// Defecto es un problema encontrado sobre un elemento del backlog. Implementa HU-20.
type Defecto struct {
	ID          int64 // lo asigna la persistencia al guardar (RN-20.6)
	ProyectoID  int64
	ElementoID  int64
	Descripcion string
	Severidad   Severidad
	Estado      EstadoDefecto
	// SprintDeteccionID queda en nil hasta que el caso de uso lo asigna (RN-20.5): el dominio no
	// conoce los sprints del proyecto.
	SprintDeteccionID *int64
}

// NuevoDefecto valida los datos y devuelve un defecto en estado Abierto, todavía sin ID ni Sprint
// de detección. Los errores se informan en el orden de la tabla de la spec.
// Implementa HU-20 (RN-20.2 a RN-20.4).
func NuevoDefecto(proyectoID, elementoID int64, descripcion string, severidad Severidad) (Defecto, error) {
	descripcion = strings.TrimSpace(descripcion)
	switch {
	case descripcion == "":
		return Defecto{}, ErrDescripcionDefectoVacia
	case !severidadValida(severidad):
		return Defecto{}, ErrSeveridadInvalida
	}
	return Defecto{
		ProyectoID:  proyectoID,
		ElementoID:  elementoID,
		Descripcion: descripcion,
		Severidad:   severidad,
		Estado:      DefectoAbierto,
	}, nil
}

func severidadValida(s Severidad) bool {
	switch s {
	case SeveridadCritica, SeveridadAlta, SeveridadMedia, SeveridadBaja:
		return true
	}
	return false
}
