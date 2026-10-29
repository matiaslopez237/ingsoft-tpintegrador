package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// LargoMaximoTitulo es la cantidad máxima de caracteres del título de un elemento (RN-05.2).
const LargoMaximoTitulo = 100

// Prioridad de un elemento del backlog. Implementa HU-05 (RN-05.5).
type Prioridad string

const (
	PrioridadAlta  Prioridad = "Alta"
	PrioridadMedia Prioridad = "Media"
	PrioridadBaja  Prioridad = "Baja"
)

// EstadoElemento es el estado de un elemento del backlog. Implementa HU-05 (RN-05.8).
type EstadoElemento string

const (
	EstadoPendiente  EstadoElemento = "Pendiente"
	EstadoEnProgreso EstadoElemento = "En progreso"
	EstadoHecho      EstadoElemento = "Hecho"
)

// Errores de HU-05 (spec, sección 7).
var (
	ErrTituloVacio          = errors.New("el título es obligatorio")
	ErrTituloDemasiadoLargo = errors.New("el título no puede superar los 100 caracteres")
	ErrDescripcionVacia     = errors.New("la descripción es obligatoria")
	ErrPrioridadInvalida    = errors.New("la prioridad debe ser Alta, Media o Baja")
	ErrStoryPointsInvalidos = errors.New("los story points deben ser un entero positivo")
	ErrSinCriterios         = errors.New("debe haber al menos un criterio de aceptación")
	ErrCriterioVacio        = errors.New("los criterios de aceptación no pueden estar vacíos")
	ErrTituloDuplicado      = errors.New("ya existe un elemento con ese título en el proyecto")
	// ErrElementoNoEncontrado: no existe un elemento con el ID buscado. Lo devuelve la persistencia
	// y lo van a usar los casos de uso de HU-06 y HU-07, igual que ErrProyectoNoEncontrado.
	ErrElementoNoEncontrado = errors.New("elemento de backlog no encontrado")
)

// ElementoBacklog es una necesidad del producto registrada en el backlog de un proyecto.
type ElementoBacklog struct {
	ID          int64 // lo asigna la persistencia al guardar (RN-05.9)
	ProyectoID  int64
	Titulo      string
	Descripcion string
	Prioridad   Prioridad
	Estado      EstadoElemento
	StoryPoints *int // nil: todavía sin estimar (RN-05.6)
	Criterios   []string
	SprintID    *int64 // nil: sin asignar a un sprint todavía (HU-09)
}

// NuevoElementoBacklog valida los datos y devuelve un elemento en estado Pendiente, todavía sin ID.
// Los errores se informan en el orden de la tabla de la spec.
// Implementa HU-05 (RN-05.2, RN-05.4 a RN-05.8).
func NuevoElementoBacklog(proyectoID int64, titulo, descripcion string, prioridad Prioridad,
	storyPoints *int, criterios []string) (ElementoBacklog, error) {
	titulo = strings.TrimSpace(titulo)
	descripcion = strings.TrimSpace(descripcion)
	switch {
	case titulo == "":
		return ElementoBacklog{}, ErrTituloVacio
	case utf8.RuneCountInString(titulo) > LargoMaximoTitulo:
		return ElementoBacklog{}, ErrTituloDemasiadoLargo
	case descripcion == "":
		return ElementoBacklog{}, ErrDescripcionVacia
	case !prioridadValida(prioridad):
		return ElementoBacklog{}, ErrPrioridadInvalida
	case storyPoints != nil && *storyPoints <= 0:
		return ElementoBacklog{}, ErrStoryPointsInvalidos
	case len(criterios) == 0:
		return ElementoBacklog{}, ErrSinCriterios
	}

	recortados := make([]string, len(criterios))
	for i, c := range criterios {
		recortados[i] = strings.TrimSpace(c)
		if recortados[i] == "" {
			return ElementoBacklog{}, ErrCriterioVacio
		}
	}

	e := ElementoBacklog{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: descripcion,
		Prioridad:   prioridad,
		Estado:      EstadoPendiente,
		Criterios:   recortados,
	}
	if storyPoints != nil {
		sp := *storyPoints // copia: el elemento no comparte el puntero recibido
		e.StoryPoints = &sp
	}
	return e, nil
}

func prioridadValida(p Prioridad) bool {
	switch p {
	case PrioridadAlta, PrioridadMedia, PrioridadBaja:
		return true
	}
	return false
}
