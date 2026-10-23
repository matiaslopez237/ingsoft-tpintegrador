// Package domain contiene las entidades y reglas de negocio, sin dependencias externas.
package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// LargoMaximoNombre es la cantidad máxima de caracteres del nombre de un proyecto (RN-01.2).
const LargoMaximoNombre = 100

// Estado es el estado de un proyecto (RN-01.7).
type Estado string

const (
	EstadoActivo     Estado = "Activo"
	EstadoEnPausa    Estado = "En pausa"
	EstadoFinalizado Estado = "Finalizado"
	EstadoPostergado Estado = "Postergado"
	EstadoCancelado  Estado = "Cancelado"
)

// Errores de HU-01 (spec, sección 7).
var (
	ErrNombreVacio            = errors.New("el nombre del proyecto es obligatorio")
	ErrNombreDemasiadoLargo   = errors.New("el nombre no puede superar los 100 caracteres")
	ErrNombreDuplicado        = errors.New("ya existe un proyecto con ese nombre")
	ErrFechaInicioObligatoria = errors.New("la fecha de inicio es obligatoria")
	ErrFechaFinAnterior       = errors.New("la fecha de fin no puede ser anterior a la de inicio")

	// ErrProyectoNoEncontrado: no existe un proyecto con el ID buscado (lo usarán HU-02 y HU-04).
	ErrProyectoNoEncontrado = errors.New("proyecto no encontrado")
)

// Proyecto es la entidad raíz: backlog, sprints, esfuerzo y métricas pertenecen a un proyecto.
type Proyecto struct {
	ID          int64 // lo asigna la persistencia al guardar (RN-01.8)
	Nombre      string
	FechaInicio time.Time
	FechaFin    *time.Time // nil: el proyecto no tiene fecha de fin
	Estado      Estado
}

// NuevoProyecto valida los datos y devuelve un proyecto en estado Activo, todavía sin ID.
// Los errores se informan en el orden de la tabla de la spec.
func NuevoProyecto(nombre string, fechaInicio time.Time, fechaFin *time.Time) (Proyecto, error) {
	nombre = strings.TrimSpace(nombre)
	switch {
	case nombre == "":
		return Proyecto{}, ErrNombreVacio
	case utf8.RuneCountInString(nombre) > LargoMaximoNombre:
		return Proyecto{}, ErrNombreDemasiadoLargo
	case fechaInicio.IsZero():
		return Proyecto{}, ErrFechaInicioObligatoria
	case fechaFin != nil && fechaFin.Before(fechaInicio):
		return Proyecto{}, ErrFechaFinAnterior
	}
	return Proyecto{
		Nombre:      nombre,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      EstadoActivo,
	}, nil
}

// ClaveNombre normaliza un nombre para comparar duplicados: sin espacios alrededor
// y sin distinguir mayúsculas (RN-01.3).
func ClaveNombre(nombre string) string {
	return strings.ToLower(strings.TrimSpace(nombre))
}
