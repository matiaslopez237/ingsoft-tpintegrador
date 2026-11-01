package bdd

import (
	"context"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
)

// sistemaHistorialRondas es lo que los escenarios de HU-16 suman al sistema de HU-14. Lo cumple
// sistemaPlanningPokerReal con el método de abajo (los métodos de un tipo pueden estar en otro
// archivo del mismo paquete), así no se toca la interfaz de HU-14.
type sistemaHistorialRondas interface {
	HistorialRondas(proyectoID, elementoID int64) ([]rondaVista, error)
}

// HistorialRondas consulta el historial con el caso de uso de HU-16, sobre la misma base SQLite del
// escenario. El caso de uso se arma acá para no sumar campos al sistema de HU-14.
func (s sistemaPlanningPokerReal) HistorialRondas(proyectoID, elementoID int64) ([]rondaVista, error) {
	uc := app.NuevoListarRondas(persistencia.NuevoRepositorioProyectos(s.db),
		persistencia.NuevoRepositorioBacklog(s.db), persistencia.NuevoRepositorioRondas(s.db))
	historial, err := uc.Ejecutar(context.Background(), proyectoID, elementoID)
	if err != nil {
		return nil, err
	}
	vistas := make([]rondaVista, len(historial))
	for i, r := range historial {
		vistas[i] = aVista(r)
	}
	return vistas, nil
}
