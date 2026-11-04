package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorDefectos es el puerto de solo lectura que necesita ListarDefectos (HU-21).
type LectorDefectos interface {
	// ListarDefectos devuelve los defectos del proyecto, en cualquier orden (lista vacía si no
	// tiene ninguno).
	ListarDefectos(ctx context.Context, proyectoID int64) ([]domain.Defecto, error)
}

// ListarDefectos es el caso de uso de apoyo de HU-21 (spec, sección 3): HU-20 no expone ningún
// listado, y la pantalla de defectos de esta historia lo necesita para elegir cuál actualizar.
type ListarDefectos struct {
	proyectos LectorProyectos
	defectos  LectorDefectos
}

func NuevoListarDefectos(proyectos LectorProyectos, defectos LectorDefectos) *ListarDefectos {
	return &ListarDefectos{proyectos: proyectos, defectos: defectos}
}

// Ejecutar verifica que el proyecto exista y trae sus defectos.
func (uc *ListarDefectos) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Defecto, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return nil, err
	}
	defectos, err := uc.defectos.ListarDefectos(ctx, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}
	return defectos, nil
}
