package domain

import (
	"errors"
	"math"
	"strings"
	"time"
)

// HorasMaximasPorRegistro es el tope de horas de un registro de esfuerzo (RN-18.8).
const HorasMaximasPorRegistro = 24

// Errores de HU-18 (spec, sección 7). El integrante vacío reutiliza ErrNombreIntegranteVacio (HU-03).
var (
	ErrFechaEsfuerzoObligatoria = errors.New("la fecha es obligatoria")
	ErrFechaEsfuerzoFutura      = errors.New("la fecha no puede ser futura")
	ErrActividadVacia           = errors.New("la actividad es obligatoria")
	ErrHorasInvalidas           = errors.New("las horas deben ser mayores a 0")
	ErrHorasNoMultiploMediaHora = errors.New("las horas se registran en múltiplos de media hora")
	ErrHorasExcedenMaximo       = errors.New("un registro no puede superar las 24 horas")
	// ErrRegistroEsfuerzoNoEncontrado: no existe un registro con el ID buscado. Lo devuelve la
	// persistencia, igual que ErrElementoNoEncontrado (sección 3: el registro se recupera por su ID).
	ErrRegistroEsfuerzoNoEncontrado = errors.New("registro de esfuerzo no encontrado")
)

// RegistroEsfuerzo deja constancia de las horas que un integrante trabajó en un elemento del backlog.
// Implementa HU-18.
type RegistroEsfuerzo struct {
	ID         int64 // lo asigna la persistencia al guardar (RN-18.10)
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      time.Time // solo el día: la hora queda en 00:00
	Actividad  string
	Horas      float64
}

// NuevoRegistroEsfuerzo valida los datos y devuelve un registro todavía sin ID.
// hoy lo pasa el caso de uso: el dominio no lee el reloj (spec, sección 5).
// Los errores se informan en el orden de la tabla de la spec.
// Implementa HU-18 (RN-18.2, RN-18.4 a RN-18.8).
func NuevoRegistroEsfuerzo(proyectoID, elementoID int64, integrante string, fecha time.Time,
	actividad string, horas float64, hoy time.Time) (RegistroEsfuerzo, error) {
	integrante = strings.TrimSpace(integrante)
	actividad = strings.TrimSpace(actividad)
	switch {
	case integrante == "":
		return RegistroEsfuerzo{}, ErrNombreIntegranteVacio
	case fecha.IsZero():
		return RegistroEsfuerzo{}, ErrFechaEsfuerzoObligatoria
	case diaCalendario(fecha).After(diaCalendario(hoy)):
		return RegistroEsfuerzo{}, ErrFechaEsfuerzoFutura
	case actividad == "":
		return RegistroEsfuerzo{}, ErrActividadVacia
	}
	if err := validarHoras(horas); err != nil {
		return RegistroEsfuerzo{}, err
	}
	// Recién acá: soloDia rearma la fecha y una fecha cero en otra zona dejaría de ser cero.
	fecha = soloDia(fecha)

	return RegistroEsfuerzo{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Integrante: integrante,
		Fecha:      fecha,
		Actividad:  actividad,
		Horas:      horas,
	}, nil
}

// validarHoras aplica RN-18.6 a RN-18.8 en este orden: > 0, múltiplo de 0,5 y máximo 24.
func validarHoras(horas float64) error {
	switch {
	case !(horas > 0): // así también NaN, que no es mayor ni menor que 0, cae acá
		return ErrHorasInvalidas
	case horas*2 != math.Trunc(horas*2): // en medias horas tiene que dar un número entero
		return ErrHorasNoMultiploMediaHora
	case horas > HorasMaximasPorRegistro:
		return ErrHorasExcedenMaximo
	}
	return nil
}

// soloDia descarta la hora de t y conserva su día y su zona (RN-18.4: se compara por día).
func soloDia(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// diaCalendario lleva el día de t, tal como se lee en su zona, a UTC. Así dos fechas se comparan
// por año, mes y día aunque vengan en zonas distintas (After compara instantes, no días).
func diaCalendario(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
