package bdd

import (
	"errors"
	"time"
)

// errNoImplementado hace fallar los escenarios mientras CrearProyecto no existe (rojo esperado).
var errNoImplementado = errors.New("pendiente: CrearProyecto todavía no está implementado")

// erroresHU01 traduce la descripción del escenario al error definido en la spec.
// Se completa cuando existan los errores del dominio (Tarea 1.4).
func erroresHU01() map[string]error {
	return map[string]error{}
}

type sistemaSinImplementar struct{}

func nuevoSistemaProyectos() (sistemaProyectos, error) {
	return sistemaSinImplementar{}, nil
}

func (sistemaSinImplementar) Crear(string, *time.Time, *time.Time) (proyectoVisto, error) {
	return proyectoVisto{}, errNoImplementado
}

func (sistemaSinImplementar) BuscarPorID(int64) (proyectoVisto, error) {
	return proyectoVisto{}, errNoImplementado
}

func (sistemaSinImplementar) Cantidad() (int, error) {
	return 0, errNoImplementado
}
