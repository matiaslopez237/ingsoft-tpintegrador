package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorHoras es el puerto para leer las horas de los elementos que entran en un total (HU-23). Es
// angosto a propósito: no se ensanchan RepositorioEsfuerzo ni LectorHorasReales (spec HU-23, sección 5).
type LectorHoras interface {
	// HorasDeProyecto devuelve una fila por elemento del proyecto, esté o no en un sprint (RN-23.10);
	// un elemento sin registros viene con 0 horas reales.
	HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error)
	// HorasDeSprint devuelve una fila por elemento que tiene hoy ese sprint (RN-23.7); un elemento sin
	// registros viene con 0 horas reales.
	HorasDeSprint(ctx context.Context, sprintID int64) ([]domain.HorasElemento, error)
}

// LectorSprintPorID es el puerto para buscar un sprint por su ID (HU-23).
type LectorSprintPorID interface {
	// BuscarSprintPorID devuelve domain.ErrSprintNoEncontrado si no existe.
	BuscarSprintPorID(ctx context.Context, id int64) (domain.Sprint, error)
}

// ConsultarHorasProyecto es el caso de uso de HU-23 para un proyecto. Reutiliza el puerto de lectura de
// HU-04. Implementa HU-23 (RN-23.1, RN-23.10, RN-23.11, RN-23.12).
type ConsultarHorasProyecto struct {
	proyectos LectorProyectos
	horas     LectorHoras
}

func NuevoConsultarHorasProyecto(proyectos LectorProyectos, horas LectorHoras) *ConsultarHorasProyecto {
	return &ConsultarHorasProyecto{proyectos: proyectos, horas: horas}
}

// Ejecutar verifica que el proyecto exista y totaliza las horas de sus elementos. El estado del proyecto
// no importa (RN-23.12).
func (uc *ConsultarHorasProyecto) Ejecutar(ctx context.Context, proyectoID int64) (domain.TotalHoras, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return domain.TotalHoras{}, err
	}
	elementos, err := uc.horas.HorasDeProyecto(ctx, proyectoID)
	if err != nil {
		return domain.TotalHoras{}, fmt.Errorf("consultar horas del proyecto %d: %w", proyectoID, err)
	}
	return domain.TotalizarHoras(elementos), nil
}

// ConsultarHorasSprint es el caso de uso de HU-23 para un sprint. Reutiliza el puerto de lectura de
// HU-04. Implementa HU-23 (RN-23.1, RN-23.7, RN-23.11, RN-23.12).
type ConsultarHorasSprint struct {
	proyectos LectorProyectos
	sprints   LectorSprintPorID
	horas     LectorHoras
}

func NuevoConsultarHorasSprint(proyectos LectorProyectos, sprints LectorSprintPorID, horas LectorHoras) *ConsultarHorasSprint {
	return &ConsultarHorasSprint{proyectos: proyectos, sprints: sprints, horas: horas}
}

// Ejecutar verifica que el proyecto exista y que el sprint sea suyo, y totaliza las horas de los
// elementos del sprint. Los errores se informan en el orden de la tabla de la spec. El estado del
// sprint y el del proyecto no importan (RN-23.12).
func (uc *ConsultarHorasSprint) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.TotalHoras, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return domain.TotalHoras{}, err
	}
	s, err := uc.sprints.BuscarSprintPorID(ctx, sprintID)
	if err != nil {
		return domain.TotalHoras{}, err
	}
	if s.ProyectoID != proyectoID { // RN-23.11: de otro proyecto es "inexistente"
		return domain.TotalHoras{}, domain.ErrSprintNoEncontrado
	}
	elementos, err := uc.horas.HorasDeSprint(ctx, sprintID)
	if err != nil {
		return domain.TotalHoras{}, fmt.Errorf("consultar horas del sprint %d: %w", sprintID, err)
	}
	return domain.TotalizarHoras(elementos), nil
}
