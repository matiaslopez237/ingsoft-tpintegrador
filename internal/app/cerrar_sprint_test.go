package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// cierreFalso guarda proyectos y sprints en memoria y permite simular fallas. Hace de
// LectorProyectos y de RepositorioCierre.
type cierreFalso struct {
	proyectos       map[int64]domain.Proyecto
	sprints         map[int64]domain.Sprint
	errBuscarSprint error
	errCerrar       error
	vecesCerro      int
}

// nuevoCierreFalso arranca con el proyecto 1 (Activo) y su sprint 10 Abierto.
func nuevoCierreFalso() *cierreFalso {
	return &cierreFalso{
		proyectos: map[int64]domain.Proyecto{1: {ID: 1, Nombre: "Alfa", Estado: domain.EstadoActivo}},
		sprints:   map[int64]domain.Sprint{10: sprintAbiertoFalso(10, 1)},
	}
}

func sprintAbiertoFalso(id, proyectoID int64) domain.Sprint {
	return domain.Sprint{
		ID: id, ProyectoID: proyectoID, Numero: 1,
		FechaInicio: time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC),
		Objetivo:    "Login de usuarios", Estado: domain.SprintAbierto,
	}
}

func (f *cierreFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *cierreFalso) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	if f.errBuscarSprint != nil {
		return domain.Sprint{}, f.errBuscarSprint
	}
	s, ok := f.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return s, nil
}

func (f *cierreFalso) CerrarSprint(_ context.Context, s *domain.Sprint) error {
	f.vecesCerro++
	if f.errCerrar != nil {
		return f.errCerrar
	}
	f.sprints[s.ID] = *s
	return nil
}

func cerrar(f *cierreFalso, proyectoID, sprintID int64) error {
	return NuevoCerrarSprint(f, f).Ejecutar(context.Background(), DatosCerrarSprint{ProyectoID: proyectoID, SprintID: sprintID})
}

// Traza: HU-11 CA-11.1
func TestCerrarSprint_SprintAbierto_QuedaCerradoConSusDatos(t *testing.T) {
	f := nuevoCierreFalso()
	antes := f.sprints[10]

	err := cerrar(f, 1, 10)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	despues := f.sprints[10]
	if despues.Estado != domain.SprintCerrado {
		t.Errorf("Estado = %q, se esperaba %q", despues.Estado, domain.SprintCerrado)
	}
	despues.Estado, antes.Estado = "", ""
	if despues != antes {
		t.Errorf("el sprint cambió: era %+v, es %+v", antes, despues)
	}
	if f.vecesCerro != 1 {
		t.Errorf("se llamó a CerrarSprint %d veces, se esperaba 1", f.vecesCerro)
	}
}

// Traza: HU-11 CA-11.3
func TestCerrarSprint_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoCierreFalso()

	err := cerrar(f, 99, 10)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}

// Traza: HU-11 CA-11.3
func TestCerrarSprint_SprintInexistente_RetornaError(t *testing.T) {
	f := nuevoCierreFalso()

	err := cerrar(f, 1, 99)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
}

// Traza: HU-11 CA-11.3, RN-11.4: un sprint de otro proyecto se trata como inexistente.
func TestCerrarSprint_SprintDeOtroProyecto_RetornaErrorYSigueAbierto(t *testing.T) {
	f := nuevoCierreFalso()
	f.proyectos[2] = domain.Proyecto{ID: 2, Nombre: "Beta", Estado: domain.EstadoActivo}
	f.sprints[20] = sprintAbiertoFalso(20, 2)

	err := cerrar(f, 1, 20)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if got := f.sprints[20].Estado; got != domain.SprintAbierto {
		t.Errorf("Estado = %q, se esperaba %q", got, domain.SprintAbierto)
	}
}

// Traza: HU-11 CA-11.2
func TestCerrarSprint_SprintYaCerrado_RetornaError(t *testing.T) {
	f := nuevoCierreFalso()
	s := f.sprints[10]
	s.Estado = domain.SprintCerrado
	f.sprints[10] = s

	err := cerrar(f, 1, 10)

	if !errors.Is(err, domain.ErrSprintYaCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintYaCerrado)
	}
	if got := f.sprints[10].Estado; got != domain.SprintCerrado {
		t.Errorf("Estado = %q, se esperaba %q", got, domain.SprintCerrado)
	}
}

// Traza: HU-11, sección 3: si hay un error, no se modifica nada.
func TestCerrarSprint_ValidacionFallida_NoLlamaACerrarSprint(t *testing.T) {
	tests := []struct {
		nombre     string
		proyectoID int64
		sprintID   int64
		preparar   func(f *cierreFalso)
	}{
		{"proyecto inexistente", 99, 10, nil},
		{"sprint inexistente", 1, 99, nil},
		{"sprint de otro proyecto", 1, 20, func(f *cierreFalso) { f.sprints[20] = sprintAbiertoFalso(20, 2) }},
		{"sprint ya cerrado", 1, 10, func(f *cierreFalso) {
			s := f.sprints[10]
			s.Estado = domain.SprintCerrado
			f.sprints[10] = s
		}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoCierreFalso()
			if tt.preparar != nil {
				tt.preparar(f)
			}

			if err := cerrar(f, tt.proyectoID, tt.sprintID); err == nil {
				t.Fatal("se esperaba un error")
			}
			if f.vecesCerro != 0 {
				t.Errorf("se llamó a CerrarSprint %d veces, no se esperaba ninguna", f.vecesCerro)
			}
		})
	}
}

// Las fallas del repositorio se propagan envueltas. ErrSprintYaCerrado al guardar es el caso de dos
// pedidos simultáneos (spec HU-11, sección 5, P-8): lo detecta el repositorio, no el dominio.
func TestCerrarSprint_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre   string
		preparar func(f *cierreFalso)
		esperado error
	}{
		{"al buscar el sprint", func(f *cierreFalso) { f.errBuscarSprint = errBase }, errBase},
		{"al cerrar", func(f *cierreFalso) { f.errCerrar = errBase }, errBase},
		{"al cerrar, otro pedido lo cerró antes", func(f *cierreFalso) { f.errCerrar = domain.ErrSprintYaCerrado }, domain.ErrSprintYaCerrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoCierreFalso()
			tt.preparar(f)

			err := cerrar(f, 1, 10)

			if !errors.Is(err, tt.esperado) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, tt.esperado)
			}
		})
	}
}
