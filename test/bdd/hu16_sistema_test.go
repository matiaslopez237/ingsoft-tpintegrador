package bdd

import "github.com/cucumber/godog"

// sistemaHistorialRondas es lo que los escenarios de HU-16 suman al sistema de HU-14. Lo cumple
// sistemaPlanningPokerReal con el método de abajo (los métodos de un tipo pueden estar en otro
// archivo del mismo paquete), así no se toca la interfaz de HU-14.
type sistemaHistorialRondas interface {
	HistorialRondas(proyectoID, elementoID int64) ([]rondaVista, error)
}

// HistorialRondas queda pendiente hasta que exista el caso de uso ListarRondas (tarea 16.5):
// mientras tanto, los escenarios que lo usan quedan en rojo.
func (s sistemaPlanningPokerReal) HistorialRondas(proyectoID, elementoID int64) ([]rondaVista, error) {
	return nil, godog.ErrPending
}
