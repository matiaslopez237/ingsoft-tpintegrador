package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ListarSprintsAnteriores es el caso de uso de HU-12. Reutiliza LectorProyectos y el ListarSprints
// del puerto RepositorioSprints de HU-08: no necesita persistencia nueva, porque ListarSprints ya
// lee los puntos que guarda HU-11.
type ListarSprintsAnteriores struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
}

func NuevoListarSprintsAnteriores(proyectos LectorProyectos, sprints RepositorioSprints) *ListarSprintsAnteriores {
	return &ListarSprintsAnteriores{proyectos: proyectos, sprints: sprints}
}

// Ejecutar devuelve los sprints cerrados del proyecto, del más reciente al más antiguo
// (domain.SprintsAnteriores), o domain.ErrProyectoNoEncontrado si el proyecto no existe, para no
// confundirlo con un proyecto sin sprints cerrados. Implementa HU-12 (RN-12.1 a RN-12.3).
func (uc *ListarSprintsAnteriores) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return nil, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	return domain.SprintsAnteriores(sprints), nil
}
