package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioEdicionBacklog es el puerto de persistencia que necesita EditarElementoBacklog.
type RepositorioEdicionBacklog interface {
	// BuscarElementoPorID devuelve domain.ErrElementoNoEncontrado si no existe.
	BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error)
	// ExisteTitulo indica si el proyecto ya tiene un elemento cuya domain.ClaveNombre(titulo) es clave.
	ExisteTitulo(ctx context.Context, proyectoID int64, clave string) (bool, error)
	// ActualizarElemento guarda de e solo los campos que c indica como enviados, para que dos
	// ediciones simultáneas de campos distintos no se pisen (HU-06, sección 5).
	ActualizarElemento(ctx context.Context, e domain.ElementoBacklog, c domain.CambiosElemento) error
}

// RepositorioEliminacionBacklog es el puerto de persistencia que necesita EliminarElementoBacklog.
type RepositorioEliminacionBacklog interface {
	BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error)
	// BuscarSprint es el de HU-09 (RepositorioAsignacion): indica si el sprint está Abierto.
	BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error)
	// TieneEsfuerzo indica si el elemento tiene registros de esfuerzo (HU-18).
	TieneEsfuerzo(ctx context.Context, elementoID int64) (bool, error)
	// EliminarElemento borra el elemento y sus criterios. Si otro pedido le registró esfuerzo en
	// el medio, devuelve domain.ErrElementoConEsfuerzo.
	EliminarElemento(ctx context.Context, elementoID int64) error
}

// DatosEditarElemento son las entradas de la edición de HU-06.
type DatosEditarElemento struct {
	ProyectoID int64
	ElementoID int64
	Cambios    domain.CambiosElemento
}

// DatosEliminarElemento son las entradas de la eliminación de HU-06.
type DatosEliminarElemento struct {
	ProyectoID int64
	ElementoID int64
}

// EditarElementoBacklog es el caso de uso de HU-06 para editar un elemento.
// Implementa HU-06 (RN-06.2, RN-06.6).
type EditarElementoBacklog struct {
	repo RepositorioEdicionBacklog
}

func NuevoEditarElementoBacklog(repo RepositorioEdicionBacklog) *EditarElementoBacklog {
	return &EditarElementoBacklog{repo: repo}
}

// Ejecutar busca el elemento en el proyecto, aplica los cambios en el dominio, verifica que el
// título no sea de otro elemento del proyecto y guarda. Devuelve el elemento como quedó guardado.
// Si algo falla, no se guarda nada.
func (uc *EditarElementoBacklog) Ejecutar(ctx context.Context, d DatosEditarElemento) (domain.ElementoBacklog, error) {
	actual, err := buscarElementoDelProyecto(ctx, uc.repo, d.ProyectoID, d.ElementoID)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}
	editado, err := actual.Editar(d.Cambios)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}

	// Si la clave no cambia, el título sigue siendo de este mismo elemento (RN-06.2).
	clave := domain.ClaveNombre(editado.Titulo)
	if clave != domain.ClaveNombre(actual.Titulo) {
		existe, err := uc.repo.ExisteTitulo(ctx, d.ProyectoID, clave)
		if err != nil {
			return domain.ElementoBacklog{}, fmt.Errorf("verificar título duplicado: %w", err)
		}
		if existe {
			return domain.ElementoBacklog{}, domain.ErrTituloDuplicado
		}
	}

	if err := uc.repo.ActualizarElemento(ctx, editado, d.Cambios); err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("actualizar elemento %d: %w", d.ElementoID, err)
	}
	// Se relee porque otro pedido pudo cambiar en el medio los campos que este no envió.
	guardado, err := uc.repo.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("releer elemento %d: %w", d.ElementoID, err)
	}
	return guardado, nil
}

// EliminarElementoBacklog es el caso de uso de HU-06 para eliminar un elemento.
// Implementa HU-06 (RN-06.6 a RN-06.8).
type EliminarElementoBacklog struct {
	repo RepositorioEliminacionBacklog
}

func NuevoEliminarElementoBacklog(repo RepositorioEliminacionBacklog) *EliminarElementoBacklog {
	return &EliminarElementoBacklog{repo: repo}
}

// Ejecutar busca el elemento en el proyecto, verifica que no esté en un Sprint activo ni tenga
// esfuerzo registrado y lo elimina.
func (uc *EliminarElementoBacklog) Ejecutar(ctx context.Context, d DatosEliminarElemento) error {
	e, err := buscarElementoDelProyecto(ctx, uc.repo, d.ProyectoID, d.ElementoID)
	if err != nil {
		return err
	}
	var sprintActivo bool
	if e.SprintID != nil {
		if _, sprintActivo, err = uc.repo.BuscarSprint(ctx, *e.SprintID); err != nil {
			return fmt.Errorf("buscar sprint %d del elemento %d: %w", *e.SprintID, d.ElementoID, err)
		}
	}
	if err := e.ValidarEliminacion(sprintActivo); err != nil {
		return err
	}
	conEsfuerzo, err := uc.repo.TieneEsfuerzo(ctx, d.ElementoID)
	if err != nil {
		return fmt.Errorf("verificar esfuerzo del elemento %d: %w", d.ElementoID, err)
	}
	if conEsfuerzo {
		return domain.ErrElementoConEsfuerzo
	}
	if err := uc.repo.EliminarElemento(ctx, d.ElementoID); err != nil {
		return fmt.Errorf("eliminar elemento %d: %w", d.ElementoID, err)
	}
	return nil
}

// lectorElemento es lo que necesita buscarElementoDelProyecto.
type lectorElemento interface {
	BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error)
}

// buscarElementoDelProyecto devuelve el elemento si pertenece al proyecto. Un elemento de otro
// proyecto se trata como inexistente (RN-06.6, como RN-18.1).
func buscarElementoDelProyecto(ctx context.Context, repo lectorElemento, proyectoID, elementoID int64) (domain.ElementoBacklog, error) {
	e, err := repo.BuscarElementoPorID(ctx, elementoID)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("buscar elemento %d: %w", elementoID, err)
	}
	if e.ProyectoID != proyectoID {
		return domain.ElementoBacklog{}, fmt.Errorf("elemento %d en el proyecto %d: %w", elementoID, proyectoID, domain.ErrElementoNoEncontrado)
	}
	return e, nil
}
