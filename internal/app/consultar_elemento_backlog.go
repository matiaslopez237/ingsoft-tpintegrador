package app

import (
	"context"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ConsultarElementoBacklog devuelve un elemento del backlog de un proyecto, para la pantalla de
// HU-06. Usa el mismo puerto de solo lectura que HU-18 (LectorElementos).
// Implementa HU-06 (RN-06.6).
type ConsultarElementoBacklog struct {
	repo LectorElementos
}

func NuevoConsultarElementoBacklog(repo LectorElementos) *ConsultarElementoBacklog {
	return &ConsultarElementoBacklog{repo: repo}
}

// Ejecutar devuelve el elemento si pertenece al proyecto; si no, domain.ErrElementoNoEncontrado.
func (uc *ConsultarElementoBacklog) Ejecutar(ctx context.Context, proyectoID, elementoID int64) (domain.ElementoBacklog, error) {
	return buscarElementoDelProyecto(ctx, uc.repo, proyectoID, elementoID)
}
