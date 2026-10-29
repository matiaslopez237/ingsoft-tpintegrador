package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioCierre es el puerto de persistencia de HU-11. Es angosto a propósito: no se ensancha
// RepositorioSprints de HU-08 (spec HU-11, sección 5).
type RepositorioCierre interface {
	// BuscarSprintPorID devuelve domain.ErrSprintNoEncontrado si no existe.
	BuscarSprintPorID(ctx context.Context, id int64) (domain.Sprint, error)
	// CerrarSprint deja el sprint guardado en Cerrado solo si sigue Abierto; si otro pedido lo
	// cerró en el medio, devuelve domain.ErrSprintYaCerrado (spec HU-11, sección 5, P-8). Persiste
	// siempre el estado Cerrado (RN-11.1) y del sprint usa solo el ID: no copia s.Estado ni otros
	// campos, así no puede guardar otro estado aunque se lo llame sin pasar por domain.Sprint.Cerrar.
	// En la misma operación aplica domain.ResumirCierre a las historias asignadas: guarda los puntos
	// planificados y completados, los deja en s, y devuelve al backlog las sin terminar (RN-11.3,
	// RN-11.6). Si falla, no guarda nada.
	CerrarSprint(ctx context.Context, s *domain.Sprint) error
}

// DatosCerrarSprint son las entradas de HU-11.
type DatosCerrarSprint struct {
	ProyectoID int64 // de la ruta /proyectos/{id}/sprints/{sprintID}/cerrar
	SprintID   int64
}

// CerrarSprint es el caso de uso de HU-11. Implementa HU-11 (RN-11.1, RN-11.2, RN-11.4).
type CerrarSprint struct {
	proyectos LectorProyectos
	sprints   RepositorioCierre
}

func NuevoCerrarSprint(proyectos LectorProyectos, sprints RepositorioCierre) *CerrarSprint {
	return &CerrarSprint{proyectos: proyectos, sprints: sprints}
}

// Ejecutar verifica que el proyecto exista y que el sprint sea suyo, lo cierra en el dominio y lo
// guarda. Los errores se informan en el orden de la tabla de la spec; si algo falla, no se guarda
// nada. El estado del proyecto no importa (RN-11.2, P-10).
func (uc *CerrarSprint) Ejecutar(ctx context.Context, d DatosCerrarSprint) error {
	if _, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID); err != nil {
		return err
	}
	s, err := uc.sprints.BuscarSprintPorID(ctx, d.SprintID)
	if err != nil {
		return err
	}
	if s.ProyectoID != d.ProyectoID { // RN-11.4: de otro proyecto es "inexistente"
		return domain.ErrSprintNoEncontrado
	}

	if err := s.Cerrar(); err != nil {
		return err
	}
	if err := uc.sprints.CerrarSprint(ctx, &s); err != nil {
		return fmt.Errorf("cerrar sprint: %w", err)
	}
	return nil
}
