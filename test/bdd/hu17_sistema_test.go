package bdd

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RegistrarEstimacion queda pendiente hasta que exista el caso de uso (tarea 17.5): mientras tanto,
// los escenarios que lo usan quedan en rojo.
func (s sistemaPlanningPokerReal) RegistrarEstimacion(proyectoID, rondaID int64, valor int) error {
	return godog.ErrPending
}

// EstimacionDeRonda también queda pendiente: la ronda todavía no guarda la estimación (tarea 17.4).
func (s sistemaPlanningPokerReal) EstimacionDeRonda(proyectoID, rondaID int64) (*int, error) {
	return nil, godog.ErrPending
}

// Elemento lee el elemento guardado, para ver sus story points y sus demás datos.
func (s sistemaPlanningPokerReal) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return persistencia.NuevoRepositorioBacklog(s.db).BuscarElementoPorID(context.Background(), elementoID)
}

// EditarStoryPoints edita los story points con el caso de uso de HU-06 (RN-17.5).
func (s sistemaPlanningPokerReal) EditarStoryPoints(proyectoID, elementoID int64, valor int) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID, ElementoID: elementoID,
		Cambios: domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &valor},
	})
	return err
}

// FijarEstadoElemento cambia el estado con el caso de uso de HU-06.
func (s sistemaPlanningPokerReal) FijarEstadoElemento(proyectoID, elementoID int64, estado string) error {
	nuevo := domain.EstadoElemento(estado)
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID, ElementoID: elementoID, Cambios: domain.CambiosElemento{Estado: &nuevo},
	})
	return err
}
