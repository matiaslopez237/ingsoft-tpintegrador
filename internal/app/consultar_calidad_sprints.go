package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// CalidadSprint es la calidad de un Sprint Cerrado (HU-24): su porcentaje de story points
// completados y sus defectos detectados y resueltos. HayPorcentaje es false si el Sprint se cerró
// sin puntos planificados (RN-24.2): en ese caso PorcentajeCompletado no significa nada.
type CalidadSprint struct {
	Sprint               domain.Sprint
	PorcentajeCompletado float64
	HayPorcentaje        bool
	DefectosDetectados   int
	DefectosResueltos    int
}

// ConsultarCalidadSprints es el caso de uso de HU-24. Reutiliza LectorProyectos, el ListarSprints
// del puerto RepositorioSprints de HU-08 y el ListarDefectos del puerto LectorDefectos de HU-21: no
// necesita persistencia nueva.
type ConsultarCalidadSprints struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
	defectos  LectorDefectos
}

func NuevoConsultarCalidadSprints(proyectos LectorProyectos, sprints RepositorioSprints, defectos LectorDefectos) *ConsultarCalidadSprints {
	return &ConsultarCalidadSprints{proyectos: proyectos, sprints: sprints, defectos: defectos}
}

// Ejecutar lista los sprints y los defectos del proyecto, se queda con los sprints Cerrados
// (domain.SprintsAnteriores, de HU-12) y arma, para cada uno, su porcentaje completado
// (domain.PorcentajeCompletado) y sus defectos detectados/resueltos (domain.ContarDefectosDeSprint).
// Implementa HU-24 (RN-24.1 a RN-24.5).
func (uc *ConsultarCalidadSprints) Ejecutar(ctx context.Context, proyectoID int64) ([]CalidadSprint, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return nil, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	defectos, err := uc.defectos.ListarDefectos(ctx, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}

	cerrados := domain.SprintsAnteriores(sprints)
	resumen := make([]CalidadSprint, len(cerrados))
	for i, s := range cerrados {
		resumen[i] = calidadDe(s, defectos)
	}
	return resumen, nil
}

// calidadDe arma la calidad de un Sprint Cerrado: su porcentaje completado (domain.PorcentajeCompletado)
// y sus defectos detectados/resueltos (domain.ContarDefectosDeSprint). La usan HU-24 y el Dashboard
// de HU-25.
func calidadDe(s domain.Sprint, defectos []domain.Defecto) CalidadSprint {
	porcentaje, hayPorcentaje := domain.PorcentajeCompletado(s)
	detectados, resueltos := domain.ContarDefectosDeSprint(defectos, s.ID)
	return CalidadSprint{
		Sprint: s, PorcentajeCompletado: porcentaje, HayPorcentaje: hayPorcentaje,
		DefectosDetectados: detectados, DefectosResueltos: resueltos,
	}
}
