package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// VelocidadProyecto es la velocidad del equipo de un proyecto (HU-22). HayDatos es false si el
// proyecto todavía no tiene ningún Sprint Cerrado (RN-22.3): en ese caso Velocidad no significa
// nada y no hay que mostrarlo como si fuera un promedio real de 0 puntos.
type VelocidadProyecto struct {
	Velocidad float64
	HayDatos  bool
	// SprintsCerrados es cuántos Sprints Cerrados entraron en el promedio (como mucho
	// domain.VentanaVelocidad), no el total de Sprints Cerrados que tiene el proyecto si son más
	// (RN-22.2): con 5 Sprints cerrados, SprintsCerrados vale 3, no 5.
	SprintsCerrados int
}

// ConsultarVelocidad es el caso de uso de HU-22. Reutiliza LectorProyectos y el ListarSprints del
// puerto RepositorioSprints de HU-08: no necesita persistencia nueva, porque ListarSprints ya lee
// los puntos que guarda HU-11 (RN-11.6).
type ConsultarVelocidad struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
}

func NuevoConsultarVelocidad(proyectos LectorProyectos, sprints RepositorioSprints) *ConsultarVelocidad {
	return &ConsultarVelocidad{proyectos: proyectos, sprints: sprints}
}

// Ejecutar lista los sprints del proyecto, se queda con los Cerrados (domain.SprintsAnteriores,
// de HU-12) y calcula el promedio de sus puntos completados (domain.Velocidad).
// Implementa HU-22 (RN-22.1 a RN-22.4).
func (uc *ConsultarVelocidad) Ejecutar(ctx context.Context, proyectoID int64) (VelocidadProyecto, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return VelocidadProyecto{}, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return VelocidadProyecto{}, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	cerrados := domain.SprintsAnteriores(sprints)
	velocidad, cantidad, hayDatos := domain.Velocidad(cerrados)
	return VelocidadProyecto{Velocidad: velocidad, HayDatos: hayDatos, SprintsCerrados: cantidad}, nil
}
