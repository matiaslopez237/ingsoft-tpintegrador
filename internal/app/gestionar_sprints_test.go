package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sprintsFalso guarda proyectos y sprints en memoria, numera los sprints como la persistencia
// y permite simular fallas. Cumple LectorProyectos y RepositorioSprints.
type sprintsFalso struct {
	proyectos    map[int64]domain.Proyecto
	sprints      []domain.Sprint
	errBuscar    error
	errListar    error
	errGuardar   error
	vecesGuardo  int
	listoProyect int64 // último proyecto cuyos sprints se listaron
}

func nuevoSprintsFalso(proyectos ...domain.Proyecto) *sprintsFalso {
	f := &sprintsFalso{proyectos: map[int64]domain.Proyecto{}}
	for _, p := range proyectos {
		f.proyectos[p.ID] = p
	}
	return f
}

func (f *sprintsFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errBuscar != nil {
		return domain.Proyecto{}, f.errBuscar
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *sprintsFalso) ListarSprints(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
	f.listoProyect = proyectoID
	if f.errListar != nil {
		return nil, f.errListar
	}
	var del []domain.Sprint
	for _, s := range f.sprints {
		if s.ProyectoID == proyectoID {
			del = append(del, s)
		}
	}
	return del, nil
}

func (f *sprintsFalso) GuardarSprint(_ context.Context, s *domain.Sprint) error {
	if f.errGuardar != nil {
		return f.errGuardar
	}
	f.vecesGuardo++
	numero := 0
	for _, e := range f.sprints {
		if e.ProyectoID == s.ProyectoID && e.Numero > numero {
			numero = e.Numero
		}
	}
	s.ID, s.Numero = int64(len(f.sprints)+1), numero+1
	f.sprints = append(f.sprints, *s)
	return nil
}

func proyectoActivo(id int64) domain.Proyecto {
	fin := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	return domain.Proyecto{ID: id, Nombre: "Alfa", FechaInicio: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), FechaFin: &fin, Estado: domain.EstadoActivo}
}

func datosSprint(proyectoID int64, desde, hasta int) DatosNuevoSprint {
	return DatosNuevoSprint{
		ProyectoID:  proyectoID,
		FechaInicio: time.Date(2026, 10, desde, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 10, hasta, 0, 0, 0, 0, time.UTC),
		Objetivo:    "Gestionar el backlog",
	}
}

// Traza: HU-08 CA-08.1, CA-08.2
func TestCrearSprint_Valido_GuardaYDevuelveConNumero(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	uc := NuevoCrearSprint(f, f)

	primero, err := uc.Ejecutar(context.Background(), datosSprint(1, 26, 27))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	segundo, err := uc.Ejecutar(context.Background(), datosSprint(1, 28, 29))
	if err != nil {
		t.Fatalf("error inesperado en el segundo sprint: %v", err)
	}

	if primero.Numero != 1 || segundo.Numero != 2 || primero.ID == 0 {
		t.Errorf("sprints = %+v y %+v, se esperaban los números 1 y 2 con ID asignado", primero, segundo)
	}
	if primero.Estado != domain.SprintAbierto || primero.Objetivo != "Gestionar el backlog" {
		t.Errorf("primer sprint = %+v", primero)
	}
	if f.vecesGuardo != 2 {
		t.Errorf("se guardó %d veces, se esperaban 2", f.vecesGuardo)
	}
}

// Traza: HU-08 CA-08.3, CA-08.4, CA-08.8: si el proyecto no existe o el dominio rechaza el sprint
// (incluida la superposición con los sprints que ya tiene), no se guarda nada.
func TestCrearSprint_Invalido_RetornaErrorYNoGuarda(t *testing.T) {
	cancelado := proyectoActivo(2)
	cancelado.Estado = domain.EstadoCancelado
	tests := []struct {
		nombre  string
		datos   DatosNuevoSprint
		wantErr error
	}{
		{"proyecto inexistente", datosSprint(99, 26, 27), domain.ErrProyectoNoEncontrado},
		{"proyecto no activo", datosSprint(2, 26, 27), domain.ErrProyectoNoActivo},
		{"se superpone con un sprint existente", datosSprint(1, 27, 28), domain.ErrSprintSuperpuesto},
		{"Sprint Goal vacío", DatosNuevoSprint{ProyectoID: 1, FechaInicio: datosSprint(1, 28, 29).FechaInicio, FechaFin: datosSprint(1, 28, 29).FechaFin}, domain.ErrObjetivoVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoSprintsFalso(proyectoActivo(1), cancelado)
			f.sprints = []domain.Sprint{{ID: 1, ProyectoID: 1, Numero: 1,
				FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)}}

			_, err := NuevoCrearSprint(f, f).Ejecutar(context.Background(), tt.datos)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if f.vecesGuardo != 0 {
				t.Errorf("no se esperaba ningún guardado, hubo %d", f.vecesGuardo)
			}
		})
	}
}

func TestCrearSprint_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		ajuste func(f *sprintsFalso)
	}{
		{"al buscar el proyecto", func(f *sprintsFalso) { f.errBuscar = errBase }},
		{"al listar los sprints", func(f *sprintsFalso) { f.errListar = errBase }},
		{"al guardar", func(f *sprintsFalso) { f.errGuardar = errBase }},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoSprintsFalso(proyectoActivo(1))
			tt.ajuste(f)
			_, err := NuevoCrearSprint(f, f).Ejecutar(context.Background(), datosSprint(1, 26, 27))
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}

// Traza: HU-08 CA-08.9
func TestListarSprints_DevuelveLosDelProyecto(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{{ID: 1, ProyectoID: 1, Numero: 1}, {ID: 2, ProyectoID: 5, Numero: 1}, {ID: 3, ProyectoID: 1, Numero: 2}}

	got, err := NuevoListarSprints(f, f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 3 {
		t.Errorf("sprints = %+v, se esperaban los 2 del proyecto 1", got)
	}
}

// Traza: HU-08 sección 7: listar los sprints de un proyecto inexistente es un error, no una lista vacía.
func TestListarSprints_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoSprintsFalso()

	_, err := NuevoListarSprints(f, f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}
