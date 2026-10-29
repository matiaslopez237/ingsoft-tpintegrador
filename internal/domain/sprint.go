package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// LargoMaximoObjetivo es la cantidad máxima de caracteres del Sprint Goal (RN-08.6).
const LargoMaximoObjetivo = 200

// EstadoSprint es el estado de un sprint. Implementa HU-08 (RN-08.8).
type EstadoSprint string

const (
	SprintAbierto EstadoSprint = "Abierto"
	SprintCerrado EstadoSprint = "Cerrado" // se llega con Cerrar (HU-11)
)

// Errores de HU-08 (spec, sección 7). Las fechas obligatorias y en orden reutilizan los de HU-01.
var (
	ErrProyectoNoActivo       = errors.New("solo se pueden crear sprints en proyectos activos")
	ErrFechaFinObligatoria    = errors.New("la fecha de fin es obligatoria")
	ErrObjetivoVacio          = errors.New("el Sprint Goal es obligatorio")
	ErrObjetivoDemasiadoLargo = errors.New("el Sprint Goal no puede superar los 200 caracteres")
	ErrSprintFueraDelProyecto = errors.New("las fechas del sprint tienen que estar dentro de las del proyecto")
	ErrSprintSuperpuesto      = errors.New("las fechas se superponen con otro sprint del proyecto")
	// ErrSprintNoEncontrado: no existe un sprint con el ID buscado. No es una regla de HU-08, pero
	// vive acá por ser del mismo agregado; la usa HU-09 al buscar el sprint destino de una asignación.
	ErrSprintNoEncontrado = errors.New("sprint no encontrado")
)

// ErrSprintYaCerrado es el error de HU-11 al cerrar un sprint que ya estaba Cerrado (RN-11.1).
var ErrSprintYaCerrado = errors.New("el sprint ya fue cerrado")

// Sprint es un período de trabajo de un proyecto, con su objetivo (Sprint Goal). Implementa HU-08.
type Sprint struct {
	ID          int64 // lo asigna la persistencia al guardar (RN-08.9)
	ProyectoID  int64
	Numero      int // lo asigna la persistencia al guardar: 1, 2, 3… dentro del proyecto (RN-08.7)
	FechaInicio time.Time
	FechaFin    time.Time
	Objetivo    string // el Sprint Goal
	Estado      EstadoSprint
	// Resumen que registra el cierre (HU-11, RN-11.6); mientras el sprint está Abierto valen 0.
	PuntosPlanificados int
	PuntosCompletados  int
}

// NuevoSprint valida los datos contra el proyecto y sus sprints existentes, y devuelve un sprint
// Abierto, todavía sin ID ni número. Los errores se informan en el orden de la tabla de la spec.
// Implementa HU-08 (RN-08.2 a RN-08.6, RN-08.8).
func NuevoSprint(p Proyecto, inicio, fin time.Time, objetivo string, existentes []Sprint) (Sprint, error) {
	objetivo = strings.TrimSpace(objetivo)
	switch {
	case p.Estado != EstadoActivo:
		return Sprint{}, ErrProyectoNoActivo
	case inicio.IsZero():
		return Sprint{}, ErrFechaInicioObligatoria
	case fin.IsZero():
		return Sprint{}, ErrFechaFinObligatoria
	case diaCalendario(fin).Before(diaCalendario(inicio)):
		return Sprint{}, ErrFechaFinAnterior
	case objetivo == "":
		return Sprint{}, ErrObjetivoVacio
	case utf8.RuneCountInString(objetivo) > LargoMaximoObjetivo:
		return Sprint{}, ErrObjetivoDemasiadoLargo
	case !dentroDelProyecto(p, inicio, fin):
		return Sprint{}, ErrSprintFueraDelProyecto
	}
	for _, s := range existentes {
		if s.ProyectoID == p.ID && seSuperponen(inicio, fin, s.FechaInicio, s.FechaFin) {
			return Sprint{}, ErrSprintSuperpuesto
		}
	}
	return Sprint{
		ProyectoID:  p.ID,
		FechaInicio: inicio,
		FechaFin:    fin,
		Objetivo:    objetivo,
		Estado:      SprintAbierto,
	}, nil
}

// PuedeCerrarse indica si el sprint se puede cerrar: solo uno Abierto (RN-11.1). Lo usa la pantalla
// de sprints para mostrar el botón "Cerrar sprint" sin repetir la regla. Implementa HU-11 (CA-11.5).
func (s Sprint) PuedeCerrarse() bool {
	return s.Estado == SprintAbierto
}

// Cerrar pasa el sprint de Abierto a Cerrado sin tocar sus demás datos. Si ya estaba Cerrado
// devuelve ErrSprintYaCerrado y no lo modifica: un sprint cerrado no se reabre ni se vuelve a cerrar.
// Implementa HU-11 (RN-11.1, RN-11.5).
func (s *Sprint) Cerrar() error {
	if s.Estado == SprintCerrado {
		return ErrSprintYaCerrado
	}
	s.Estado = SprintCerrado
	return nil
}

// dentroDelProyecto indica si el sprint cae dentro de las fechas del proyecto (RN-08.4).
func dentroDelProyecto(p Proyecto, inicio, fin time.Time) bool {
	if diaCalendario(inicio).Before(diaCalendario(p.FechaInicio)) {
		return false
	}
	return p.FechaFin == nil || !diaCalendario(fin).After(diaCalendario(*p.FechaFin))
}

// seSuperponen indica si dos períodos de días completos comparten al menos un día (RN-08.5).
func seSuperponen(inicioA, finA, inicioB, finB time.Time) bool {
	return !diaCalendario(inicioA).After(diaCalendario(finB)) && !diaCalendario(inicioB).After(diaCalendario(finA))
}
