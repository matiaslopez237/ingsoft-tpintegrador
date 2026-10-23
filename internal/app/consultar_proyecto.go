package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorProyectos es el puerto de solo lectura que necesita ConsultarProyecto.
// No tiene métodos de escritura: así la consulta no puede modificar nada (HU-04, sección 5).
type LectorProyectos interface {
	// BuscarPorID devuelve domain.ErrProyectoNoEncontrado si no existe.
	BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error)
}

// ConsultarProyecto es el caso de uso de HU-04.
type ConsultarProyecto struct {
	repo LectorProyectos
}

func NuevoConsultarProyecto(repo LectorProyectos) *ConsultarProyecto {
	return &ConsultarProyecto{repo: repo}
}

// Ejecutar devuelve el proyecto con sus integrantes y roles.
// Los errores se envuelven con %w: errors.Is sigue distinguiendo domain.ErrProyectoNoEncontrado
// de una falla interna (HU-04, sección 7).
// Implementa HU-04 (CA-04.1, CA-04.2, CA-04.5).
func (uc *ConsultarProyecto) Ejecutar(ctx context.Context, id int64) (domain.Proyecto, error) {
	p, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return domain.Proyecto{}, fmt.Errorf("consultar proyecto %d: %w", id, err)
	}
	return p, nil
}
