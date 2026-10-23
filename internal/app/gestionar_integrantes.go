package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioIntegrantes es el puerto de persistencia que necesitan AgregarIntegrante y QuitarIntegrante.
type RepositorioIntegrantes interface {
	// BuscarPorID devuelve domain.ErrProyectoNoEncontrado si no existe.
	BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error)
	// GuardarIntegrantes persiste la lista de integrantes de un proyecto ya existente.
	GuardarIntegrantes(ctx context.Context, p *domain.Proyecto) error
}

// DatosIntegrante son las entradas para agregar un integrante (HU-03).
type DatosIntegrante struct {
	ProyectoID int64
	Nombre     string
	Rol        string
}

// AgregarIntegrante es el caso de uso de HU-03 para sumar un integrante a un proyecto.
type AgregarIntegrante struct {
	repo RepositorioIntegrantes
}

func NuevoAgregarIntegrante(repo RepositorioIntegrantes) *AgregarIntegrante {
	return &AgregarIntegrante{repo: repo}
}

// Ejecutar busca el proyecto, aplica la regla de dominio y persiste el resultado.
// Si el proyecto no existe o los datos son inválidos, no se guarda nada.
// Implementa HU-03 (CA-03.1, CA-03.2, CA-03.5, CA-03.6).
func (uc *AgregarIntegrante) Ejecutar(ctx context.Context, d DatosIntegrante) error {
	p, err := uc.repo.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return err
	}
	if err := p.AgregarIntegrante(d.Nombre, d.Rol); err != nil {
		return err
	}
	if err := uc.repo.GuardarIntegrantes(ctx, &p); err != nil {
		return fmt.Errorf("guardar integrantes: %w", err)
	}
	return nil
}

// DatosQuitarIntegrante son las entradas para quitar un integrante (HU-03).
type DatosQuitarIntegrante struct {
	ProyectoID int64
	Nombre     string
}

// QuitarIntegrante es el caso de uso de HU-03 para sacar un integrante de un proyecto.
type QuitarIntegrante struct {
	repo RepositorioIntegrantes
}

func NuevoQuitarIntegrante(repo RepositorioIntegrantes) *QuitarIntegrante {
	return &QuitarIntegrante{repo: repo}
}

// Ejecutar busca el proyecto, aplica la regla de dominio y persiste el resultado.
// Si el proyecto no existe o el integrante no pertenece a él, no se guarda nada.
// Implementa HU-03 (CA-03.3, CA-03.4).
func (uc *QuitarIntegrante) Ejecutar(ctx context.Context, d DatosQuitarIntegrante) error {
	p, err := uc.repo.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return err
	}
	if err := p.QuitarIntegrante(d.Nombre); err != nil {
		return err
	}
	if err := uc.repo.GuardarIntegrantes(ctx, &p); err != nil {
		return fmt.Errorf("guardar integrantes: %w", err)
	}
	return nil
}
