package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// repoFalso guarda los proyectos en memoria y permite simular fallas.
type repoFalso struct {
	guardados   []domain.Proyecto
	siguienteID int64
	errExiste   error
	errGuardar  error
}

func (r *repoFalso) ExisteNombre(_ context.Context, clave string) (bool, error) {
	if r.errExiste != nil {
		return false, r.errExiste
	}
	for _, p := range r.guardados {
		if domain.ClaveNombre(p.Nombre) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *repoFalso) Guardar(_ context.Context, p *domain.Proyecto) error {
	if r.errGuardar != nil {
		return r.errGuardar
	}
	r.siguienteID++
	p.ID = r.siguienteID
	r.guardados = append(r.guardados, *p)
	return nil
}

func datosValidos(nombre string) DatosNuevoProyecto {
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	return DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    &fin,
	}
}

// Traza: HU-01 CA-01.1, CA-01.9
func TestCrearProyecto_DatosValidos_GuardaYDevuelveConID(t *testing.T) {
	repo := &repoFalso{}
	uc := NuevoCrearProyecto(repo)

	p, err := uc.Ejecutar(context.Background(), datosValidos("TP Integrador"))

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.ID == 0 {
		t.Error("se esperaba un ID asignado")
	}
	if p.Estado != domain.EstadoActivo {
		t.Errorf("Estado = %q, se esperaba %q", p.Estado, domain.EstadoActivo)
	}
	if len(repo.guardados) != 1 || repo.guardados[0].ID != p.ID {
		t.Errorf("se esperaba 1 proyecto guardado con ID %d, hay %+v", p.ID, repo.guardados)
	}
}

// Traza: HU-01 CA-01.9
func TestCrearProyecto_DosProyectos_TienenIDsDistintos(t *testing.T) {
	uc := NuevoCrearProyecto(&repoFalso{})

	a, errA := uc.Ejecutar(context.Background(), datosValidos("Alfa"))
	b, errB := uc.Ejecutar(context.Background(), datosValidos("Beta"))

	if errA != nil || errB != nil {
		t.Fatalf("errores inesperados: %v, %v", errA, errB)
	}
	if a.ID == b.ID {
		t.Errorf("ambos proyectos tienen el ID %d", a.ID)
	}
}

// Traza: HU-01 CA-01.2: si el dominio rechaza los datos, no se consulta ni se guarda nada.
func TestCrearProyecto_DatosInvalidos_NoGuarda(t *testing.T) {
	repo := &repoFalso{errExiste: errors.New("no debería consultarse el repositorio")}
	uc := NuevoCrearProyecto(repo)

	_, err := uc.Ejecutar(context.Background(), datosValidos("   "))

	if !errors.Is(err, domain.ErrNombreVacio) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrNombreVacio)
	}
	if len(repo.guardados) != 0 {
		t.Errorf("no se esperaban proyectos guardados, hay %d", len(repo.guardados))
	}
}

// Traza: HU-01 CA-01.4
func TestCrearProyecto_NombreDuplicado_RetornaErrorYNoGuarda(t *testing.T) {
	for _, nombre := range []string{"Alfa", "alfa", "  ALFA "} {
		t.Run(nombre, func(t *testing.T) {
			repo := &repoFalso{}
			uc := NuevoCrearProyecto(repo)
			if _, err := uc.Ejecutar(context.Background(), datosValidos("Alfa")); err != nil {
				t.Fatalf("preparar: %v", err)
			}

			_, err := uc.Ejecutar(context.Background(), datosValidos(nombre))

			if !errors.Is(err, domain.ErrNombreDuplicado) {
				t.Errorf("error = %v, se esperaba %v", err, domain.ErrNombreDuplicado)
			}
			if len(repo.guardados) != 1 {
				t.Errorf("se esperaba 1 proyecto guardado, hay %d", len(repo.guardados))
			}
		})
	}
}

func TestCrearProyecto_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		repo   *repoFalso
	}{
		{"al verificar duplicados", &repoFalso{errExiste: errBase}},
		{"al guardar", &repoFalso{errGuardar: errBase}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoCrearProyecto(tt.repo).Ejecutar(context.Background(), datosValidos("Alfa"))
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}
