package bdd

import (
	"errors"
	"net/url"
)

// errWebNoImplementado hace fallar los escenarios mientras las pantallas no existen (rojo esperado).
var errWebNoImplementado = errors.New("pendiente: las pantallas de HU-30 todavía no están implementadas")

type webSinImplementar struct{}

func nuevoSistemaWeb() (sistemaWeb, error) {
	return webSinImplementar{}, nil
}

func (webSinImplementar) GuardarProyecto(string) (int64, error) { return 0, errWebNoImplementado }

func (webSinImplementar) AgregarIntegrante(int64, string, string) error { return errWebNoImplementado }

func (webSinImplementar) Pedir(string, string, url.Values) (respuestaWeb, error) {
	return respuestaWeb{}, errWebNoImplementado
}

func (webSinImplementar) ContarProyectos() (int, error) { return 0, errWebNoImplementado }

func (webSinImplementar) BuscarProyecto(int64) (proyectoGuardado, error) {
	return proyectoGuardado{}, errWebNoImplementado
}

func (webSinImplementar) Cerrar() error { return nil }
