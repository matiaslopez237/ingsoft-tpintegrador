// Package app contiene los casos de uso: orquestan el dominio y definen los puertos
// que implementan los adaptadores.
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioProyectos es el puerto de persistencia que necesita CrearProyecto.
type RepositorioProyectos interface {
	// ExisteNombre indica si ya hay un proyecto cuya domain.ClaveNombre es clave.
	ExisteNombre(ctx context.Context, clave string) (bool, error)
	// Guardar persiste el proyecto y le asigna su ID (RN-01.8).
	Guardar(ctx context.Context, p *domain.Proyecto) error
}

// DatosNuevoProyecto son las entradas de HU-01.
type DatosNuevoProyecto struct {
	Nombre      string
	FechaInicio time.Time
	FechaFin    *time.Time
}

// CrearProyecto es el caso de uso de HU-01.
type CrearProyecto struct {
	repo RepositorioProyectos
}

func NuevoCrearProyecto(repo RepositorioProyectos) *CrearProyecto {
	return &CrearProyecto{repo: repo}
}

// Ejecutar valida los datos, verifica que el nombre no esté repetido y guarda el proyecto.
func (uc *CrearProyecto) Ejecutar(ctx context.Context, d DatosNuevoProyecto) (domain.Proyecto, error) {
	p, err := domain.NuevoProyecto(d.Nombre, d.FechaInicio, d.FechaFin)
	if err != nil {
		return domain.Proyecto{}, err
	}

	existe, err := uc.repo.ExisteNombre(ctx, domain.ClaveNombre(p.Nombre))
	if err != nil {
		return domain.Proyecto{}, fmt.Errorf("verificar nombre duplicado: %w", err)
	}
	if existe {
		return domain.Proyecto{}, domain.ErrNombreDuplicado
	}

	if err := uc.repo.Guardar(ctx, &p); err != nil {
		return domain.Proyecto{}, fmt.Errorf("guardar proyecto: %w", err)
	}
	return p, nil
}
