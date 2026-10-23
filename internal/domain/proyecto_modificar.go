package domain

import (
	"errors"
	"slices"
	"time"
)

// Errores nuevos de HU-02 (spec, sección 7). Los de nombre y fechas son los de HU-01.
var (
	ErrSinCambios     = errors.New("no se indicó ningún dato para modificar")
	ErrEstadoInvalido = errors.New("el estado indicado no es válido")
)

// CambiosProyecto son los datos de una modificación parcial (HU-02, RN-02.1).
// Un campo en nil no se envía y conserva su valor actual.
type CambiosProyecto struct {
	Nombre *string
	// FechaInicio: una fecha cero equivale a quitarla, y eso no está permitido (RN-02.4).
	FechaInicio *time.Time
	// CambiarFechaFin indica que se envía la fecha de fin. Con FechaFin en nil, se quita (RN-02.5).
	CambiarFechaFin bool
	FechaFin        *time.Time
	Estado          *Estado
}

func (c CambiosProyecto) vacios() bool {
	return c.Nombre == nil && c.FechaInicio == nil && !c.CambiarFechaFin && c.Estado == nil
}

// Modificar devuelve una copia del proyecto con los cambios aplicados; el original no cambia.
// Valida el resultado con las mismas reglas que al crear (NuevoProyecto) y luego el estado.
// La unicidad del nombre se verifica en el caso de uso, como en HU-01.
// Implementa HU-02 (RN-02.1, RN-02.2, RN-02.4 a RN-02.10).
func (p Proyecto) Modificar(c CambiosProyecto) (Proyecto, error) {
	if c.vacios() {
		return Proyecto{}, ErrSinCambios
	}

	nombre, inicio, fin, estado := p.Nombre, p.FechaInicio, p.FechaFin, p.Estado
	if c.Nombre != nil {
		nombre = *c.Nombre
	}
	if c.FechaInicio != nil {
		inicio = *c.FechaInicio
	}
	if c.CambiarFechaFin {
		fin = c.FechaFin
	}
	if fin != nil {
		copia := *fin // la copia no comparte la fecha de fin con el original ni con c
		fin = &copia
	}
	if c.Estado != nil {
		estado = *c.Estado
	}

	modificado, err := NuevoProyecto(nombre, inicio, fin)
	if err != nil {
		return Proyecto{}, err
	}
	if !EstadoValido(estado) {
		return Proyecto{}, ErrEstadoInvalido
	}
	modificado.ID = p.ID
	modificado.Estado = estado
	modificado.Integrantes = clonarIntegrantes(p.Integrantes) // RN-02.10
	return modificado, nil
}

// clonarIntegrantes copia la lista y los roles de cada integrante, para que la copia
// modificada no comparta memoria con el original.
func clonarIntegrantes(integrantes []Integrante) []Integrante {
	copia := slices.Clone(integrantes)
	for i := range copia {
		copia[i].Roles = slices.Clone(copia[i].Roles)
	}
	return copia
}

// EstadoValido indica si e es uno de los estados de RN-01.7 (distingue mayúsculas).
func EstadoValido(e Estado) bool {
	switch e {
	case EstadoActivo, EstadoEnPausa, EstadoFinalizado, EstadoPostergado, EstadoCancelado:
		return true
	}
	return false
}
