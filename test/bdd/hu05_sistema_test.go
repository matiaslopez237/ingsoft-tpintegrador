package bdd

import "errors"

// errBacklogNoImplementado hace fallar los escenarios mientras CrearElementoBacklog no existe (rojo esperado).
var errBacklogNoImplementado = errors.New("pendiente: CrearElementoBacklog todavía no está implementado")

// erroresHU05 traduce la descripción del escenario al error definido en la spec.
// Se completa cuando existan los errores del dominio (Tarea 5.4).
func erroresHU05() map[string]error {
	return map[string]error{}
}

type backlogSinImplementar struct{}

func nuevoSistemaBacklog() (sistemaBacklog, error) {
	return backlogSinImplementar{}, nil
}

func (backlogSinImplementar) CrearProyecto(string) (int64, error) {
	return 0, errBacklogNoImplementado
}

func (backlogSinImplementar) CrearElemento(int64, datosElemento) (elementoVisto, error) {
	return elementoVisto{}, errBacklogNoImplementado
}

func (backlogSinImplementar) BuscarElemento(int64) (elementoVisto, error) {
	return elementoVisto{}, errBacklogNoImplementado
}

func (backlogSinImplementar) ContarElementos(int64) (int, error) {
	return 0, errBacklogNoImplementado
}

func (backlogSinImplementar) Cerrar() error { return nil }
