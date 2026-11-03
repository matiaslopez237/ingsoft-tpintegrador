package bdd

import (
	"context"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RegistrarEstimacion registra con el caso de uso de HU-17, sobre la misma base SQLite del escenario.
// El caso de uso se arma acá para no sumar campos al sistema de HU-14.
func (s sistemaPlanningPokerReal) RegistrarEstimacion(proyectoID, rondaID int64, valor int) error {
	uc := app.NuevoRegistrarEstimacion(persistencia.NuevoRepositorioProyectos(s.db),
		persistencia.NuevoRepositorioBacklog(s.db), persistencia.NuevoRepositorioRondas(s.db))
	_, err := uc.Ejecutar(context.Background(), app.DatosEstimacion{ProyectoID: proyectoID, RondaID: rondaID, Valor: valor})
	return err
}

// EstimacionDeRonda lee la estimación de la ronda con el resumen de HU-14, que es lo que se muestra.
func (s sistemaPlanningPokerReal) EstimacionDeRonda(proyectoID, rondaID int64) (*int, error) {
	r, err := s.consultar.Ejecutar(context.Background(), proyectoID, rondaID)
	if err != nil {
		return nil, err
	}
	return r.EstimacionAcordada, nil
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
