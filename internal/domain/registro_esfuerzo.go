package domain

import "time"

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
	return RegistroEsfuerzo{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Integrante: integrante,
		Fecha:      fecha,
		Actividad:  actividad,
		Horas:      horas,
	}, nil
}
