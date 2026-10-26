package bdd

import (
	"errors"
	"net/url"
)

// errEdicionNoImplementada hace fallar los escenarios mientras HU-06 no existe (rojo esperado).
var errEdicionNoImplementada = errors.New("pendiente: editar y eliminar elementos todavía no está implementado")

// erroresHU06 traduce la descripción del escenario al error definido en la spec.
// Se completa cuando existan los errores del dominio (Tarea 6.4).
func erroresHU06() map[string]error {
	return map[string]error{}
}

type edicionSinImplementar struct{}

func nuevoSistemaEdicionBacklog() (sistemaEdicionBacklog, error) {
	return edicionSinImplementar{}, nil
}

func (edicionSinImplementar) CrearProyecto(string) (int64, error) { return 0, errEdicionNoImplementada }

func (edicionSinImplementar) CrearElemento(int64, string, string, int) (int64, error) {
	return 0, errEdicionNoImplementada
}

func (edicionSinImplementar) RegistrarEsfuerzo(int64, int64, string, float64) error {
	return errEdicionNoImplementada
}

func (edicionSinImplementar) EditarElemento(int64, int64, map[string]string) error {
	return errEdicionNoImplementada
}

func (edicionSinImplementar) EliminarElemento(int64, int64) error { return errEdicionNoImplementada }

func (edicionSinImplementar) BuscarElemento(int64) (elementoEditable, error) {
	return elementoEditable{}, errEdicionNoImplementada
}

func (edicionSinImplementar) ContarElementos(int64) (int, error) { return 0, errEdicionNoImplementada }

func (edicionSinImplementar) Pedir(string, string, url.Values) (respuestaWeb, error) {
	return respuestaWeb{}, errEdicionNoImplementada
}

func (edicionSinImplementar) Cerrar() error { return nil }
