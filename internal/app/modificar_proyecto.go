package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioModificacion es el puerto de persistencia que necesita ModificarProyecto.
type RepositorioModificacion interface {
	// BuscarPorID devuelve domain.ErrProyectoNoEncontrado si no existe.
	BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error)
	// ExisteNombre indica si ya hay un proyecto cuya domain.ClaveNombre es clave.
	ExisteNombre(ctx context.Context, clave string) (bool, error)
	// Actualizar guarda de p solo los campos que c indica como enviados, para que modificaciones
	// simultáneas no se pisen, y vuelve a verificar RN-02.5 al guardar (HU-02, sección 5).
	// No toca los integrantes (RN-02.10).
	Actualizar(ctx context.Context, p domain.Proyecto, c domain.CambiosProyecto) error
}

// DatosModificarProyecto son las entradas de HU-02.
type DatosModificarProyecto struct {
	ProyectoID int64
	Cambios    domain.CambiosProyecto
}

// ModificarProyecto es el caso de uso de HU-02.
type ModificarProyecto struct {
	repo RepositorioModificacion
}

func NuevoModificarProyecto(repo RepositorioModificacion) *ModificarProyecto {
	return &ModificarProyecto{repo: repo}
}

// Ejecutar busca el proyecto, aplica los cambios en el dominio, verifica que el nombre no sea
// de otro proyecto y guarda el resultado. Si algo falla, no se guarda nada.
// Implementa HU-02 (RN-02.1 a RN-02.3, RN-02.9, CA-02.9, CA-02.10, sección 7).
func (uc *ModificarProyecto) Ejecutar(ctx context.Context, d DatosModificarProyecto) (domain.Proyecto, error) {
	actual, err := uc.repo.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.Proyecto{}, err
	}
	modificado, err := actual.Modificar(d.Cambios)
	if err != nil {
		return domain.Proyecto{}, err
	}

	// Si la clave no cambia, el nombre sigue siendo de este mismo proyecto (RN-02.3).
	clave := domain.ClaveNombre(modificado.Nombre)
	if clave != domain.ClaveNombre(actual.Nombre) {
		existe, err := uc.repo.ExisteNombre(ctx, clave)
		if err != nil {
			return domain.Proyecto{}, fmt.Errorf("verificar nombre duplicado: %w", err)
		}
		if existe {
			return domain.Proyecto{}, domain.ErrNombreDuplicado
		}
	}

	if err := uc.repo.Actualizar(ctx, modificado, d.Cambios); err != nil {
		return domain.Proyecto{}, fmt.Errorf("actualizar proyecto: %w", err)
	}
	// Se relee porque otro pedido pudo cambiar en el medio los campos que este no envió (sección 5).
	guardado, err := uc.repo.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.Proyecto{}, fmt.Errorf("releer proyecto modificado: %w", err)
	}
	return guardado, nil
}
