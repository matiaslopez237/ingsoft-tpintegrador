package app

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ListadorProyectos es el puerto de solo lectura que necesita ListarProyectos (HU-30).
type ListadorProyectos interface {
	// ListarProyectos devuelve todos los proyectos, en cualquier orden y sin sus integrantes.
	ListarProyectos(ctx context.Context) ([]domain.Proyecto, error)
}

// ListarProyectos es el caso de uso de la lista de proyectos de la interfaz web.
// Implementa HU-30 (RN-30.1).
type ListarProyectos struct {
	repo ListadorProyectos
}

func NuevoListarProyectos(repo ListadorProyectos) *ListarProyectos {
	return &ListarProyectos{repo: repo}
}

// Ejecutar devuelve los proyectos ordenados por nombre sin distinguir mayúsculas
// (domain.ClaveNombre) y, a igual nombre, por ID. No modifica la lista del repositorio.
func (uc *ListarProyectos) Ejecutar(ctx context.Context) ([]domain.Proyecto, error) {
	proyectos, err := uc.repo.ListarProyectos(ctx)
	if err != nil {
		return nil, fmt.Errorf("listar proyectos: %w", err)
	}
	ordenados := slices.Clone(proyectos)
	slices.SortFunc(ordenados, func(a, b domain.Proyecto) int {
		return cmp.Or(
			cmp.Compare(domain.ClaveNombre(a.Nombre), domain.ClaveNombre(b.Nombre)),
			cmp.Compare(a.ID, b.ID),
		)
	})
	return ordenados, nil
}
