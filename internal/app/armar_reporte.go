package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ArmarReporteProyecto es el caso de uso de HU-26 para un proyecto. Reutiliza LectorProyectos
// (HU-04), ListadorElementos (HU-07), RepositorioSprints (HU-08), LectorHoras (HU-23) y
// LectorDefectos (HU-21): no agrega ningún puerto de persistencia nuevo.
type ArmarReporteProyecto struct {
	proyectos LectorProyectos
	elementos ListadorElementos
	sprints   RepositorioSprints
	horas     LectorHoras
	defectos  LectorDefectos
}

func NuevoArmarReporteProyecto(proyectos LectorProyectos, elementos ListadorElementos, sprints RepositorioSprints, horas LectorHoras, defectos LectorDefectos) *ArmarReporteProyecto {
	return &ArmarReporteProyecto{proyectos: proyectos, elementos: elementos, sprints: sprints, horas: horas, defectos: defectos}
}

// Ejecutar busca el proyecto y arma su reporte: todas sus historias (RN-26.1), su esfuerzo total
// (HU-23), la velocidad del equipo y la calidad de sus Sprints cerrados (HU-22, HU-24, RN-26.5) y
// todos sus defectos (RN-26.7). Implementa HU-26.
func (uc *ArmarReporteProyecto) Ejecutar(ctx context.Context, proyectoID int64) (domain.Reporte, error) {
	proyecto, err := uc.proyectos.BuscarPorID(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	historias, err := uc.elementos.ListarElementos(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar historias del proyecto %d: %w", proyectoID, err)
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	horasElementos, err := uc.horas.HorasDeProyecto(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("leer horas del proyecto %d: %w", proyectoID, err)
	}
	defectos, err := uc.defectos.ListarDefectos(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}

	cerrados := domain.SprintsAnteriores(sprints)
	velocidad, _, hayVelocidad := domain.Velocidad(cerrados)
	metricas := make([]domain.FilaMetrica, len(cerrados))
	for i, s := range cerrados {
		porcentaje, hayPorcentaje := domain.PorcentajeCompletado(s)
		detectados, resueltos := domain.ContarDefectosDeSprint(defectos, s.ID)
		metricas[i] = domain.FilaMetrica{
			Sprint: s, PorcentajeCompletado: porcentaje, HayPorcentaje: hayPorcentaje,
			DefectosDetectados: detectados, DefectosResueltos: resueltos,
		}
	}
	esfuerzo := domain.TotalizarHoras(horasElementos)
	reporte := domain.ArmarReporteProyecto(proyecto, historias, esfuerzo, velocidad, hayVelocidad, metricas, defectos)
	reporte.NumerosSprint = domain.NumerosDeSprint(sprints)                          // DEF-05
	reporte.SprintsAbiertos, reporte.SprintsCerrados = domain.ContarSprints(sprints) // DEF-05
	return reporte, nil
}

// LectorSprintsReporte es lo que el reporte de un sprint necesita de los sprints: buscar el del reporte
// y listar los del proyecto, para los números de los sprints de sus defectos (DEF-05).
type LectorSprintsReporte interface {
	LectorSprintPorID
	ListarSprints(ctx context.Context, proyectoID int64) ([]domain.Sprint, error)
}

// ArmarReporteSprint es el caso de uso de HU-26 para un Sprint. Reutiliza LectorProyectos (HU-04),
// LectorSprintsReporte (BuscarSprintPorID de HU-23 y ListarSprints de HU-08, este último por DEF-05),
// ListadorElementos (HU-07), LectorHoras (HU-23) y LectorDefectos (HU-21): no agrega ningún puerto
// de persistencia nuevo.
type ArmarReporteSprint struct {
	proyectos LectorProyectos
	sprints   LectorSprintsReporte
	elementos ListadorElementos
	horas     LectorHoras
	defectos  LectorDefectos
}

func NuevoArmarReporteSprint(proyectos LectorProyectos, sprints LectorSprintsReporte, elementos ListadorElementos, horas LectorHoras, defectos LectorDefectos) *ArmarReporteSprint {
	return &ArmarReporteSprint{proyectos: proyectos, sprints: sprints, elementos: elementos, horas: horas, defectos: defectos}
}

// Ejecutar busca el proyecto y el Sprint (verifica que sea suyo, RN-26.10) y arma su reporte: sus
// historias asignadas (RN-26.2), su esfuerzo (HU-23), su propia calidad sin exigir que esté Cerrado
// (RN-26.6) y sus defectos detectados o resueltos ahí (RN-26.8). Implementa HU-26.
func (uc *ArmarReporteSprint) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.Reporte, error) {
	proyecto, err := uc.proyectos.BuscarPorID(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	sprint, err := uc.sprints.BuscarSprintPorID(ctx, sprintID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("buscar sprint %d: %w", sprintID, err)
	}
	if sprint.ProyectoID != proyectoID { // RN-26.10: de otro proyecto es "inexistente"
		return domain.Reporte{}, domain.ErrSprintNoEncontrado
	}
	todasHistorias, err := uc.elementos.ListarElementos(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar historias del proyecto %d: %w", proyectoID, err)
	}
	horasElementos, err := uc.horas.HorasDeSprint(ctx, sprintID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("leer horas del sprint %d: %w", sprintID, err)
	}
	// Todos los sprints del proyecto, para mostrar con su número los sprints de los defectos, que
	// pueden ser otros que el del reporte (RN-26.8, DEF-05).
	sprintsProyecto, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	todosDefectos, err := uc.defectos.ListarDefectos(ctx, proyectoID)
	if err != nil {
		return domain.Reporte{}, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}

	porcentaje, hayPorcentaje := domain.PorcentajeCompletado(sprint)
	detectados, resueltos := domain.ContarDefectosDeSprint(todosDefectos, sprintID)
	metrica := domain.FilaMetrica{
		Sprint: sprint, PorcentajeCompletado: porcentaje, HayPorcentaje: hayPorcentaje,
		DefectosDetectados: detectados, DefectosResueltos: resueltos,
	}
	esfuerzo := domain.TotalizarHoras(horasElementos)
	historias := domain.HistoriasDeSprint(todasHistorias, sprintID)
	defectos := domain.DefectosDeSprint(todosDefectos, sprintID)
	reporte := domain.ArmarReporteSprint(proyecto, sprint, historias, esfuerzo, metrica, defectos)
	reporte.NumerosSprint = domain.NumerosDeSprint(sprintsProyecto) // DEF-05
	return reporte, nil
}
