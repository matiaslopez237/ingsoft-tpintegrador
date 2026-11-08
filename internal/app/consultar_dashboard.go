package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorHorasDeProyecto es la parte de LectorHoras (HU-23) que necesita el Dashboard.
type LectorHorasDeProyecto interface {
	HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error)
}

// Dashboard es el panel de un proyecto (HU-25): su estado y los datos de los gráficos.
type Dashboard struct {
	Backlog          domain.ConteoBacklog
	SprintEnCurso    domain.SprintEnCurso
	HaySprintEnCurso bool // false: ningún sprint Abierto incluye hoy (RN-25.3)
	// SprintsCerrados va del más antiguo al más reciente (RN-25.5), con el porcentaje y los defectos
	// de HU-24.
	SprintsCerrados []CalidadSprint
	Velocidad       VelocidadProyecto
	Horas           domain.TotalHoras
}

// ConsultarDashboard es el caso de uso de HU-25. Reutiliza los puertos de HU-04, HU-07, HU-08, HU-21
// y HU-23: no necesita persistencia nueva. El reloj se inyecta para fijar "hoy" en los tests, igual
// que en HU-18 y HU-20.
type ConsultarDashboard struct {
	proyectos LectorProyectos
	sprints   RepositorioSprints
	elementos ListadorElementos
	defectos  LectorDefectos
	horas     LectorHorasDeProyecto
	reloj     func() time.Time
}

func NuevoConsultarDashboard(proyectos LectorProyectos, sprints RepositorioSprints, elementos ListadorElementos,
	defectos LectorDefectos, horas LectorHorasDeProyecto, reloj func() time.Time) *ConsultarDashboard {
	return &ConsultarDashboard{proyectos: proyectos, sprints: sprints, elementos: elementos, defectos: defectos, horas: horas, reloj: reloj}
}

// Ejecutar lee de nuevo todos los datos del proyecto en cada consulta (RN-25.1) y arma el panel con
// las reglas que ya existen: ContarBacklog y ResumirSprintEnCurso (HU-25), PorcentajeCompletado y
// ContarDefectosDeSprint (HU-24), Velocidad (HU-22) y TotalizarHoras (HU-23). Implementa HU-25
// (RN-25.1 a RN-25.5).
func (uc *ConsultarDashboard) Ejecutar(ctx context.Context, proyectoID int64) (Dashboard, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return Dashboard{}, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	sprints, err := uc.sprints.ListarSprints(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	elementos, err := uc.elementos.ListarElementos(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, fmt.Errorf("listar elementos del proyecto %d: %w", proyectoID, err)
	}
	defectos, err := uc.defectos.ListarDefectos(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}
	horas, err := uc.horas.HorasDeProyecto(ctx, proyectoID)
	if err != nil {
		return Dashboard{}, fmt.Errorf("consultar horas del proyecto %d: %w", proyectoID, err)
	}

	d := Dashboard{Backlog: domain.ContarBacklog(elementos), Horas: domain.TotalizarHoras(horas)}
	d.SprintEnCurso, d.HaySprintEnCurso = domain.ResumirSprintEnCurso(sprints, elementos, uc.reloj())

	cerrados := domain.SprintsCerradosCronologicos(sprints)
	d.SprintsCerrados = make([]CalidadSprint, len(cerrados))
	for i, s := range cerrados {
		d.SprintsCerrados[i] = calidadDe(s, defectos)
	}

	// Velocidad promedia los más recientes primero (orden de SprintsAnteriores, RN-22.2).
	v, cantidad, hayDatos := domain.Velocidad(domain.SprintsAnteriores(sprints))
	d.Velocidad = VelocidadProyecto{Velocidad: v, HayDatos: hayDatos, SprintsCerrados: cantidad}
	return d, nil
}
