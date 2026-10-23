package app

import (
	"context"

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
// Pendiente: se implementa en la etapa GREEN de la Tarea 4.4.
func (uc *ConsultarProyecto) Ejecutar(ctx context.Context, id int64) (domain.Proyecto, error) {
	return domain.Proyecto{}, nil
}
