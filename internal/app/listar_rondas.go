package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorHistorialRondas es el puerto angosto de HU-16: no se ensancha RepositorioRondas de HU-14,
// para no tocar sus fakes (spec HU-16, sección 5).
type LectorHistorialRondas interface {
	// ListarRondas devuelve todas las rondas del elemento, con sus votantes y votos.
	ListarRondas(ctx context.Context, elementoID int64) ([]domain.Ronda, error)
}

// ListarRondas es el caso de uso de HU-16 que devuelve el historial de rondas de un elemento.
type ListarRondas struct {
	proyectos LectorProyectos
	elementos LectorElementos
	rondas    LectorHistorialRondas
}

func NuevoListarRondas(proyectos LectorProyectos, elementos LectorElementos, rondas LectorHistorialRondas) *ListarRondas {
	return &ListarRondas{proyectos: proyectos, elementos: elementos, rondas: rondas}
}

// Ejecutar verifica el proyecto y el elemento (uno de otro proyecto se trata como inexistente) y
// devuelve sus rondas de la más reciente a la más antigua (domain.HistorialRondas), sin los votos
// de la ronda abierta. Implementa HU-16 (RN-16.3, RN-16.5).
func (uc *ListarRondas) Ejecutar(ctx context.Context, proyectoID, elementoID int64) ([]domain.ResumenRonda, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return nil, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	e, err := buscarElementoDelProyecto(ctx, uc.elementos, proyectoID, elementoID) // RN-16.5, como RN-14.7
	if err != nil {
		return nil, err
	}
	rondas, err := uc.rondas.ListarRondas(ctx, e.ID)
	if err != nil {
		return nil, fmt.Errorf("listar rondas del elemento %d: %w", e.ID, err)
	}
	return domain.HistorialRondas(rondas), nil
}
