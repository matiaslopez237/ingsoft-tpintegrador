package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioBacklog es el puerto de persistencia que necesita CrearElementoBacklog.
type RepositorioBacklog interface {
	ExisteProyecto(ctx context.Context, id int64) (bool, error)
	// ExisteTitulo indica si el proyecto ya tiene un elemento cuya domain.ClaveNombre(titulo) es clave.
	ExisteTitulo(ctx context.Context, proyectoID int64, clave string) (bool, error)
	// GuardarElemento persiste el elemento con sus criterios y le asigna su ID (RN-05.9).
	GuardarElemento(ctx context.Context, e *domain.ElementoBacklog) error
}

// DatosNuevoElemento son las entradas de HU-05.
type DatosNuevoElemento struct {
	ProyectoID  int64
	Titulo      string
	Descripcion string
	Prioridad   domain.Prioridad
	StoryPoints *int
	Criterios   []string
}

// CrearElementoBacklog es el caso de uso de HU-05. Implementa HU-05 (RN-05.1, RN-05.3, RN-05.9).
type CrearElementoBacklog struct {
	repo RepositorioBacklog
}

func NuevoCrearElementoBacklog(repo RepositorioBacklog) *CrearElementoBacklog {
	return &CrearElementoBacklog{repo: repo}
}

// Ejecutar verifica que el proyecto exista, valida el elemento en el dominio, verifica que el
// título no esté repetido en el proyecto y lo guarda. Si algo falla, no se guarda nada.
func (uc *CrearElementoBacklog) Ejecutar(ctx context.Context, d DatosNuevoElemento) (domain.ElementoBacklog, error) {
	existe, err := uc.repo.ExisteProyecto(ctx, d.ProyectoID)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("buscar proyecto: %w", err)
	}
	if !existe {
		return domain.ElementoBacklog{}, domain.ErrProyectoNoEncontrado
	}

	e, err := domain.NuevoElementoBacklog(d.ProyectoID, d.Titulo, d.Descripcion, d.Prioridad, d.StoryPoints, d.Criterios)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}

	repetido, err := uc.repo.ExisteTitulo(ctx, d.ProyectoID, domain.ClaveNombre(e.Titulo))
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("verificar título duplicado: %w", err)
	}
	if repetido {
		return domain.ElementoBacklog{}, domain.ErrTituloDuplicado
	}

	if err := uc.repo.GuardarElemento(ctx, &e); err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("guardar elemento: %w", err)
	}
	return e, nil
}
