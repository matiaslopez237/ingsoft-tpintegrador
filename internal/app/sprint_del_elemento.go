package app

import (
	"context"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// buscadorSprint es la parte de RepositorioAsignacion y RepositorioCompletar que consulta un sprint.
type buscadorSprint interface {
	BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error)
}

// sprintDelElementoAbierto consulta si el sprint asignado al elemento está abierto. Si el
// elemento no tiene sprint, no consulta nada y devuelve false: la regla de dominio informa
// ErrElementoSinSprint. Lo usan DesasignarHistoria (HU-09) y CompletarHistoria (HU-10).
func sprintDelElementoAbierto(ctx context.Context, repo buscadorSprint, e domain.ElementoBacklog) (bool, error) {
	if e.SprintID == nil {
		return false, nil
	}
	_, abierto, err := repo.BuscarSprint(ctx, *e.SprintID)
	if err != nil {
		return false, err
	}
	return abierto, nil
}
