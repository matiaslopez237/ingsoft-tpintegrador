package persistencia

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioProyectos = (*RepositorioProyectos)(nil)

// nuevoRepoDePrueba abre una base SQLite real en un archivo temporal del test.
func nuevoRepoDePrueba(t *testing.T) *RepositorioProyectos {
	t.Helper()
	db, err := AbrirDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NuevoRepositorioProyectos(db)
}

func proyectoDePrueba(t *testing.T, nombre string, fin *time.Time) domain.Proyecto {
	t.Helper()
	p, err := domain.NuevoProyecto(nombre, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), fin)
	if err != nil {
		t.Fatalf("crear proyecto de prueba: %v", err)
	}
	return p
}

// Traza: HU-01 CA-01.9: lo que se guarda se recupera igual, con ID asignado.
func TestRepositorioProyectos_GuardarYBuscar_DevuelveLoMismo(t *testing.T) {
	ctx := context.Background()
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		nombre string
		fin    *time.Time
	}{
		{"con fecha de fin", &fin},
		{"sin fecha de fin", nil},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			repo := nuevoRepoDePrueba(t)
			p := proyectoDePrueba(t, "Gestión Año 2026", tt.fin)

			if err := repo.Guardar(ctx, &p); err != nil {
				t.Fatalf("guardar: %v", err)
			}
			if p.ID <= 0 {
				t.Fatalf("ID = %d, se esperaba un ID asignado", p.ID)
			}

			got, err := repo.BuscarPorID(ctx, p.ID)
			if err != nil {
				t.Fatalf("buscar: %v", err)
			}
			if !reflect.DeepEqual(got, p) {
				t.Errorf("se recuperó %+v, se esperaba %+v", got, p)
			}
		})
	}
}

func TestRepositorioProyectos_IDsAutoincrementalesDistintos(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	a, b := proyectoDePrueba(t, "Alfa", nil), proyectoDePrueba(t, "Beta", nil)

	if err := repo.Guardar(ctx, &a); err != nil {
		t.Fatal(err)
	}
	if err := repo.Guardar(ctx, &b); err != nil {
		t.Fatal(err)
	}

	if a.ID == b.ID {
		t.Errorf("ambos proyectos tienen el ID %d", a.ID)
	}
	if n, err := repo.Contar(ctx); err != nil || n != 2 {
		t.Errorf("Contar() = %d, %v; se esperaba 2", n, err)
	}
}

// Traza: HU-01 RN-01.3: ExisteNombre usa la clave normalizada y la base rechaza duplicados.
func TestRepositorioProyectos_NombreDuplicado(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	alfa := proyectoDePrueba(t, "Alfa", nil)
	if err := repo.Guardar(ctx, &alfa); err != nil {
		t.Fatal(err)
	}

	existe, err := repo.ExisteNombre(ctx, domain.ClaveNombre("ALFA"))
	if err != nil || !existe {
		t.Errorf("ExisteNombre(alfa) = %v, %v; se esperaba true", existe, err)
	}
	existe, err = repo.ExisteNombre(ctx, domain.ClaveNombre("Beta"))
	if err != nil || existe {
		t.Errorf("ExisteNombre(beta) = %v, %v; se esperaba false", existe, err)
	}

	otro := proyectoDePrueba(t, "alfa", nil)
	if err := repo.Guardar(ctx, &otro); !errors.Is(err, domain.ErrNombreDuplicado) {
		t.Errorf("Guardar duplicado: error = %v, se esperaba %v", err, domain.ErrNombreDuplicado)
	}
}

func TestRepositorioProyectos_BuscarInexistente(t *testing.T) {
	_, err := nuevoRepoDePrueba(t).BuscarPorID(context.Background(), 999)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}
