package domain

import (
	"errors"
	"strings"
	"time"
)

// Severidad de un defecto. Implementa HU-20 (RN-20.3).
type Severidad string

const (
	SeveridadCritica Severidad = "Crítica"
	SeveridadAlta    Severidad = "Alta"
	SeveridadMedia   Severidad = "Media"
	SeveridadBaja    Severidad = "Baja"
)

// EstadoDefecto es el estado de un defecto. Implementa HU-20 (RN-20.4) y HU-21 (RN-21.1): un
// defecto nuevo siempre arranca Abierto; Iniciar, Resolver y Cerrar lo avanzan hasta Cerrado, un
// paso por vez, sin retroceder ni saltear ninguno (RN-21.5).
type EstadoDefecto string

const (
	DefectoAbierto    EstadoDefecto = "Abierto"
	DefectoEnProgreso EstadoDefecto = "En progreso" // se llega con Iniciar (HU-21)
	DefectoResuelto   EstadoDefecto = "Resuelto"    // se llega con Resolver (HU-21)
	DefectoCerrado    EstadoDefecto = "Cerrado"     // se llega con Cerrar (HU-21)
)

// Errores de HU-20 (spec, sección 7).
var (
	ErrDescripcionDefectoVacia = errors.New("la descripción es obligatoria")
	ErrSeveridadInvalida       = errors.New("la severidad debe ser Crítica, Alta, Media o Baja")
)

// Errores de HU-21 (spec, sección 7).
var (
	ErrDefectoNoEncontrado = errors.New("defecto no encontrado")
	ErrDefectoNoAbierto    = errors.New("el defecto no está Abierto")
	ErrDefectoNoEnProgreso = errors.New("el defecto no está En progreso")
	ErrDefectoNoResuelto   = errors.New("el defecto no está Resuelto")
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
	// SprintResolucionID queda en nil hasta que el caso de uso lo asigna al resolver el defecto
	// (RN-21.5): el dominio tampoco conoce los sprints del proyecto para esta asignación.
	SprintResolucionID *int64
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

// PuedeIniciarse indica si el defecto se puede iniciar: solo uno Abierto (RN-21.2). Lo usa la
// pantalla de defectos para mostrar el botón "Marcar en progreso" sin repetir la regla, igual que
// Sprint.PuedeCerrarse en HU-11. Implementa HU-21 (CA-21.11).
func (d Defecto) PuedeIniciarse() bool {
	return d.Estado == DefectoAbierto
}

// Iniciar pasa el defecto de Abierto a En progreso sin tocar sus demás datos. Si no está Abierto
// (ya En progreso, Resuelto o Cerrado) devuelve ErrDefectoNoAbierto y no lo modifica: no se puede
// iniciar dos veces ni volver a iniciar uno que ya avanzó. Implementa HU-21 (RN-21.2, RN-21.5).
func (d *Defecto) Iniciar() error {
	if !d.PuedeIniciarse() {
		return ErrDefectoNoAbierto
	}
	d.Estado = DefectoEnProgreso
	return nil
}

// PuedeResolverse indica si el defecto se puede resolver: solo uno En progreso (RN-21.3). Lo usa
// la pantalla de defectos para mostrar el botón "Marcar resuelto" sin repetir la regla, igual que
// Sprint.PuedeCerrarse en HU-11. Implementa HU-21 (CA-21.11).
func (d Defecto) PuedeResolverse() bool {
	return d.Estado == DefectoEnProgreso
}

// Resolver pasa el defecto de En progreso a Resuelto sin tocar sus demás datos. Si no está En
// progreso (todavía Abierto, ya Resuelto o Cerrado) devuelve ErrDefectoNoEnProgreso y no lo
// modifica: no se puede saltear En progreso, resolver dos veces, ni volver a resolver uno ya
// Cerrado. Asignar el Sprint de resolución es tarea del caso de uso (RN-21.6): el dominio no
// conoce los sprints del proyecto. Implementa HU-21 (RN-21.3, RN-21.5).
func (d *Defecto) Resolver() error {
	if !d.PuedeResolverse() {
		return ErrDefectoNoEnProgreso
	}
	d.Estado = DefectoResuelto
	return nil
}

// PuedeCerrarse indica si el defecto se puede cerrar: solo uno Resuelto (RN-21.4). Lo usa la
// pantalla de defectos para mostrar el botón "Cerrar" sin repetir la regla. Implementa HU-21
// (CA-21.11).
func (d Defecto) PuedeCerrarse() bool {
	return d.Estado == DefectoResuelto
}

// Cerrar pasa el defecto de Resuelto a Cerrado sin tocar sus demás datos (el SprintResolucionID ya
// quedó asignado al resolver). Si no está Resuelto (todavía Abierto, En progreso, o ya Cerrado)
// devuelve ErrDefectoNoResuelto y no lo modifica: no se puede saltear Resuelto ni cerrar dos veces.
// Implementa HU-21 (RN-21.4, RN-21.5).
func (d *Defecto) Cerrar() error {
	if !d.PuedeCerrarse() {
		return ErrDefectoNoResuelto
	}
	d.Estado = DefectoCerrado
	return nil
}

func severidadValida(s Severidad) bool {
	switch s {
	case SeveridadCritica, SeveridadAlta, SeveridadMedia, SeveridadBaja:
		return true
	}
	return false
}

// SprintVigente busca, entre los sprints dados, el que está Abierto y cuyas fechas (día completo,
// igual que RN-08.5) incluyen hoy. Como RN-08.5 impide que los sprints de un proyecto se
// superpongan en fechas, a lo sumo uno califica. Devuelve su ID y true, o 0 y false si ninguno
// califica. Implementa HU-20 (RN-20.5).
func SprintVigente(sprints []Sprint, hoy time.Time) (int64, bool) {
	h := diaCalendario(hoy)
	for _, s := range sprints {
		if s.Estado != SprintAbierto {
			continue
		}
		if !diaCalendario(s.FechaInicio).After(h) && !diaCalendario(s.FechaFin).Before(h) {
			return s.ID, true
		}
	}
	return 0, false
}
