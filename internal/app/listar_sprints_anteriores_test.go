package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sprintGuardado arma un sprint del proyecto, con su número, el día del mes de inicio (de octubre
// de 2026) y su estado.
func sprintGuardado(id, proyectoID int64, numero, diaInicio int, estado domain.EstadoSprint) domain.Sprint {
	inicio := time.Date(2026, 10, diaInicio, 0, 0, 0, 0, time.UTC)
	return domain.Sprint{ID: id, ProyectoID: proyectoID, Numero: numero, FechaInicio: inicio,
		FechaFin: inicio.AddDate(0, 0, 11), Objetivo: "Objetivo", Estado: estado}
}

// Traza: HU-12 CA-12.1, CA-12.4: los cerrados del proyecto, del más reciente al más antiguo, sin
// los abiertos ni los de otro proyecto.
func TestListarSprintsAnteriores_DevuelveLosCerradosDelProyecto(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{
		sprintGuardado(1, 1, 1, 1, domain.SprintCerrado),
		sprintGuardado(2, 1, 2, 13, domain.SprintCerrado),
		sprintGuardado(3, 1, 3, 26, domain.SprintAbierto),
		sprintGuardado(4, 5, 1, 13, domain.SprintCerrado), // de otro proyecto
	}

	got, err := NuevoListarSprintsAnteriores(f, f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 2 || got[0].ID != 2 || got[1].ID != 1 {
		t.Errorf("sprints = %+v, se esperaban los IDs 2 y 1", got)
	}
}

// Traza: HU-12 CA-12.2: sin sprints cerrados es una lista vacía, no un error.
func TestListarSprintsAnteriores_SinCerrados_ListaVacia(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{sprintGuardado(1, 1, 1, 1, domain.SprintAbierto)}

	got, err := NuevoListarSprintsAnteriores(f, f).Ejecutar(context.Background(), 1)

	if err != nil || len(got) != 0 {
		t.Errorf("sprints = %+v, error = %v; se esperaba una lista vacía sin error", got, err)
	}
}

// Traza: HU-12 CA-12.5: con el proyecto inexistente no se leen sus sprints.
func TestListarSprintsAnteriores_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoSprintsFalso()

	_, err := NuevoListarSprintsAnteriores(f, f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if f.listoProyect != 0 {
		t.Errorf("se listaron los sprints del proyecto %d; no se esperaba listar nada", f.listoProyect)
	}
}

// Traza: HU-12 sección 7: una falla al leer se propaga, sin perder la causa.
func TestListarSprintsAnteriores_FallaDelRepositorio_SePropaga(t *testing.T) {
	disco := errors.New("disco lleno")
	tests := []struct {
		name  string
		falla func(*sprintsFalso)
	}{
		{"al buscar el proyecto", func(f *sprintsFalso) { f.errBuscar = disco }},
		{"al listar los sprints", func(f *sprintsFalso) { f.errListar = disco }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoSprintsFalso(proyectoActivo(1))
			tt.falla(f)

			_, err := NuevoListarSprintsAnteriores(f, f).Ejecutar(context.Background(), 1)

			if !errors.Is(err, disco) {
				t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
			}
		})
	}
}
