package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioEstadoDefecto es el puerto de persistencia de HU-21. Es angosto a propósito: no se
// ensancha RepositorioDefectos de HU-20 (spec HU-21, sección 5).
type RepositorioEstadoDefecto interface {
	// BuscarDefectoPorID devuelve domain.ErrDefectoNoEncontrado si no existe.
	BuscarDefectoPorID(ctx context.Context, id int64) (domain.Defecto, error)
	// ActualizarEstadoDefecto guarda el estado nuevo de d (y, al resolver, el SprintResolucionID)
	// solo si el defecto sigue en anterior; si otro pedido lo cambió en el medio, devuelve el error
	// de dominio que corresponda al estado con el que choca (ErrDefectoNoAbierto/ErrDefectoNoEnProgreso/
	// ErrDefectoNoResuelto), igual que CerrarSprint en HU-11 (spec, sección 5: pedidos simultáneos).
	ActualizarEstadoDefecto(ctx context.Context, d *domain.Defecto, anterior domain.EstadoDefecto) error
}

// DatosActualizarEstadoDefecto son las entradas de HU-21 (spec, sección 2).
type DatosActualizarEstadoDefecto struct {
	ProyectoID int64 // de la ruta /proyectos/{id}/defectos/{defectoID}/iniciar|resolver|cerrar
	DefectoID  int64
}

// ActualizarEstadoDefecto es el caso de uso de HU-21, con un método por acción (Iniciar/Resolver/
// Cerrar) que comparten las mismas dependencias. Implementa HU-21 (RN-21.2 a RN-21.7).
type ActualizarEstadoDefecto struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
	defectos  RepositorioEstadoDefecto
	reloj     func() time.Time
}

// NuevoActualizarEstadoDefecto recibe el reloj para que "hoy" se pueda fijar en los tests, igual
// que HU-18 y HU-20.
func NuevoActualizarEstadoDefecto(proyectos LectorProyectos, sprints RepositorioSprints,
	defectos RepositorioEstadoDefecto, reloj func() time.Time) *ActualizarEstadoDefecto {
	return &ActualizarEstadoDefecto{proyectos: proyectos, sprints: sprints, defectos: defectos, reloj: reloj}
}

// Iniciar busca el defecto, lo pasa a En progreso en el dominio y lo guarda. Si algo falla, no se
// guarda nada. Implementa HU-21 (RN-21.2, RN-21.5, RN-21.7).
func (uc *ActualizarEstadoDefecto) Iniciar(ctx context.Context, d DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	defecto, err := uc.buscar(ctx, d)
	if err != nil {
		return domain.Defecto{}, err
	}
	anterior := defecto.Estado
	if err := defecto.Iniciar(); err != nil {
		return domain.Defecto{}, err
	}
	if err := uc.defectos.ActualizarEstadoDefecto(ctx, &defecto, anterior); err != nil {
		return domain.Defecto{}, fmt.Errorf("iniciar defecto: %w", err)
	}
	return defecto, nil
}

// Resolver busca el defecto, lo pasa a Resuelto en el dominio, le asigna el Sprint de resolución
// vigente (si hay uno) y lo guarda. Si algo falla, no se guarda nada.
// Implementa HU-21 (RN-21.3, RN-21.6, RN-21.7).
func (uc *ActualizarEstadoDefecto) Resolver(ctx context.Context, d DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	defecto, err := uc.buscar(ctx, d)
	if err != nil {
		return domain.Defecto{}, err
	}
	anterior := defecto.Estado
	if err := defecto.Resolver(); err != nil {
		return domain.Defecto{}, err
	}

	sprints, err := uc.sprints.ListarSprints(ctx, defecto.ProyectoID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("listar sprints del proyecto %d: %w", defecto.ProyectoID, err)
	}
	if id, ok := domain.SprintVigente(sprints, uc.reloj()); ok {
		defecto.SprintResolucionID = &id
	}

	if err := uc.defectos.ActualizarEstadoDefecto(ctx, &defecto, anterior); err != nil {
		return domain.Defecto{}, fmt.Errorf("resolver defecto: %w", err)
	}
	return defecto, nil
}

// Cerrar busca el defecto, lo pasa a Cerrado en el dominio (sin tocar el Sprint de resolución, ya
// asignado al resolver) y lo guarda. Si algo falla, no se guarda nada.
// Implementa HU-21 (RN-21.4, RN-21.7).
func (uc *ActualizarEstadoDefecto) Cerrar(ctx context.Context, d DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	defecto, err := uc.buscar(ctx, d)
	if err != nil {
		return domain.Defecto{}, err
	}
	anterior := defecto.Estado
	if err := defecto.Cerrar(); err != nil {
		return domain.Defecto{}, err
	}
	if err := uc.defectos.ActualizarEstadoDefecto(ctx, &defecto, anterior); err != nil {
		return domain.Defecto{}, fmt.Errorf("cerrar defecto: %w", err)
	}
	return defecto, nil
}

// buscar resuelve el proyecto y el defecto de la ruta (RN-21.7): un defecto de otro proyecto se
// trata como inexistente, igual que el elemento en HU-20 (RN-20.1).
func (uc *ActualizarEstadoDefecto) buscar(ctx context.Context, d DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID); err != nil {
		return domain.Defecto{}, fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}
	defecto, err := uc.defectos.BuscarDefectoPorID(ctx, d.DefectoID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("buscar defecto %d: %w", d.DefectoID, err)
	}
	if defecto.ProyectoID != d.ProyectoID {
		return domain.Defecto{}, fmt.Errorf("defecto %d en el proyecto %d: %w", d.DefectoID, d.ProyectoID, domain.ErrDefectoNoEncontrado)
	}
	return defecto, nil
}
