package domain

import (
	"errors"
	"slices"
)

// Errores de HU-17 (spec, sección 7).
var (
	ErrEstimacionRondaAbierta = errors.New("la estimación se registra cuando votaron todos")
	ErrRondaNoEsLaUltima      = errors.New("solo se puede registrar la estimación en la última ronda")
	ErrEstimacionYaRegistrada = errors.New("la ronda ya tiene una estimación acordada")
	ErrEstimacionInvalida     = errors.New("la estimación debe ser 1, 2, 3, 5, 8, 13 o 21")
)

// RegistrarEstimacion guarda en la ronda el valor acordado. ultima es la última ronda del elemento,
// nil si no se conoce. El valor no tiene que coincidir con un voto (RN-17.2). Si hay un error, la
// ronda no cambia; los errores se informan en el orden de la tabla de la spec. Copiar el valor a los
// story points del elemento (RN-17.1) es parte de guardarlo. Implementa HU-17 (RN-17.2 a RN-17.4).
func (r *Ronda) RegistrarEstimacion(valor int, ultima *Ronda, e ElementoBacklog) error {
	var ultimaID *int64
	if ultima != nil {
		ultimaID = &ultima.ID
	}
	if err := errorParaRegistrar(r.ID, r.Estado, r.EstimacionAcordada, ultimaID, e); err != nil {
		return err
	}
	if !slices.Contains(valoresVoto, valor) {
		return ErrEstimacionInvalida
	}
	r.EstimacionAcordada = &valor
	return nil
}

// PuedeRegistrarEstimacion indica si la pantalla de la ronda ofrece registrar la estimación: la misma
// regla que RegistrarEstimacion, sin mirar el valor. Implementa HU-17 (RN-17.3, RN-17.4).
func PuedeRegistrarEstimacion(r ResumenRonda, ultima *ResumenRonda, e ElementoBacklog) bool {
	var ultimaID *int64
	if ultima != nil {
		ultimaID = &ultima.ID
	}
	return errorParaRegistrar(r.ID, r.Estado, r.EstimacionAcordada, ultimaID, e) == nil
}

// errorParaRegistrar aplica RN-17.3 y RN-17.4, en el orden de la tabla de la spec.
func errorParaRegistrar(id int64, estado EstadoRonda, estimacion *int, ultimaID *int64, e ElementoBacklog) error {
	switch {
	case estado != RondaCerrada:
		return ErrEstimacionRondaAbierta
	case ultimaID == nil || *ultimaID != id:
		return ErrRondaNoEsLaUltima
	case e.Estado == EstadoHecho:
		return ErrElementoTerminado
	case estimacion != nil:
		return ErrEstimacionYaRegistrada
	}
	return nil
}
