package domain

import (
	"errors"
	"math"
	"strings"
	"time"
)

// Errores de HU-18 (spec, sección 7). El integrante vacío reutiliza ErrNombreIntegranteVacio (HU-03).
var (
	ErrFechaEsfuerzoObligatoria = errors.New("la fecha es obligatoria")
	ErrFechaEsfuerzoFutura      = errors.New("la fecha no puede ser futura")
	ErrActividadVacia           = errors.New("la actividad es obligatoria")
	ErrHorasInvalidas           = errors.New("las horas deben ser mayores a 0")
	ErrHorasNoMultiploMediaHora = errors.New("las horas se registran en múltiplos de media hora")
	ErrHorasExcedenMaximo       = errors.New("un registro no puede superar las 24 horas")
)

// RegistroEsfuerzo deja constancia de las horas que un integrante trabajó en un elemento del backlog.
// Implementa HU-18.
type RegistroEsfuerzo struct {
	ID         int64 // lo asigna la persistencia al guardar (RN-18.10)
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
}

// NuevoRegistroEsfuerzo valida los datos y devuelve un registro todavía sin ID.
// hoy lo pasa el caso de uso: el dominio no lee el reloj (spec, sección 5).
func NuevoRegistroEsfuerzo(proyectoID, elementoID int64, integrante string, fecha time.Time,
	actividad string, horas float64, hoy time.Time) (RegistroEsfuerzo, error) {
	integrante = strings.TrimSpace(integrante)
	actividad = strings.TrimSpace(actividad)
	fecha = time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, hoy.Location())
	switch {
	case integrante == "":
		return RegistroEsfuerzo{}, ErrNombreIntegranteVacio
	case fecha.IsZero():
		return RegistroEsfuerzo{}, ErrFechaEsfuerzoObligatoria
	case fecha.After(hoy):
		return RegistroEsfuerzo{}, ErrFechaEsfuerzoFutura
	case actividad == "":
		return RegistroEsfuerzo{}, ErrActividadVacia
	case horas <= 0:
		return RegistroEsfuerzo{}, ErrHorasInvalidas
	case horas*2 != math.Trunc(horas*2):
		return RegistroEsfuerzo{}, ErrHorasNoMultiploMediaHora
	case horas > 24:
		return RegistroEsfuerzo{}, ErrHorasExcedenMaximo
	}

	return RegistroEsfuerzo{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Integrante: integrante,
		Fecha:      fecha,
		Actividad:  actividad,
		Horas:      horas,
	}, nil
}
