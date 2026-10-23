package app

import (
	"context"
	"fmt"
	"slices"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioIntegrantes es el puerto de persistencia que necesitan AgregarIntegrante, QuitarIntegrante y AgregarRol.
type RepositorioIntegrantes interface {
	// BuscarPorID devuelve domain.ErrProyectoNoEncontrado si no existe.
	BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error)
	// AgregarIntegrante inserta un integrante en el proyecto sin tocar a los demás,
	// para que pedidos simultáneos no se pisen (HU-03, sección 5).
	AgregarIntegrante(ctx context.Context, proyectoID int64, i domain.Integrante) error
	// QuitarIntegrante borra un integrante del proyecto sin tocar a los demás.
	QuitarIntegrante(ctx context.Context, proyectoID int64, nombre string) error
	// AgregarRol inserta un rol del integrante sin tocar sus otros roles.
	AgregarRol(ctx context.Context, proyectoID int64, nombre, rol string) error
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
	nuevo := p.Integrantes[len(p.Integrantes)-1] // el que acaba de validar y normalizar el dominio
	if err := uc.repo.AgregarIntegrante(ctx, p.ID, nuevo); err != nil {
		return fmt.Errorf("guardar integrante: %w", err)
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
	if err := uc.repo.QuitarIntegrante(ctx, p.ID, d.Nombre); err != nil {
		return fmt.Errorf("quitar integrante: %w", err)
	}
	return nil
}

// DatosRol son las entradas para sumarle un rol a un integrante (HU-03).
type DatosRol struct {
	ProyectoID int64
	Nombre     string
	Rol        string
}

// AgregarRol es el caso de uso de HU-03 para sumarle un rol a un integrante del proyecto.
type AgregarRol struct {
	repo RepositorioIntegrantes
}

func NuevoAgregarRol(repo RepositorioIntegrantes) *AgregarRol {
	return &AgregarRol{repo: repo}
}

// Ejecutar busca el proyecto, aplica la regla de dominio y persiste solo el rol nuevo.
// Si el proyecto no existe o el dominio rechaza el pedido, no se guarda nada.
// Implementa HU-03 (CA-03.7, CA-03.8, CA-03.9, CA-03.10).
func (uc *AgregarRol) Ejecutar(ctx context.Context, d DatosRol) error {
	p, err := uc.repo.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return err
	}
	if err := p.AgregarRol(d.Nombre, d.Rol); err != nil {
		return err
	}
	integ := p.Integrantes[slices.IndexFunc(p.Integrantes, func(i domain.Integrante) bool {
		return domain.ClaveNombre(i.Nombre) == domain.ClaveNombre(d.Nombre)
	})]
	nuevo := integ.Roles[len(integ.Roles)-1] // el que acaba de validar y normalizar el dominio
	if err := uc.repo.AgregarRol(ctx, p.ID, integ.Nombre, nuevo); err != nil {
		return fmt.Errorf("guardar rol: %w", err)
	}
	return nil
}
