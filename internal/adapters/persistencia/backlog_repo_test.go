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
var _ app.RepositorioBacklog = (*RepositorioBacklog)(nil)

// nuevoBacklogDePrueba abre una base SQLite real y crea los proyectos indicados; devuelve sus IDs.
func nuevoBacklogDePrueba(t *testing.T, proyectos ...string) (*RepositorioBacklog, []int64) {
	t.Helper()
	db, err := AbrirDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	repoProyectos := NuevoRepositorioProyectos(db)
	var ids []int64
	for _, nombre := range proyectos {
		p, err := domain.NuevoProyecto(nombre, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := repoProyectos.Guardar(context.Background(), &p); err != nil {
			t.Fatalf("guardar proyecto: %v", err)
		}
		ids = append(ids, p.ID)
	}
	return NuevoRepositorioBacklog(db), ids
}

func elementoDePrueba(t *testing.T, proyectoID int64, titulo string, sp *int, criterios ...string) domain.ElementoBacklog {
	t.Helper()
	if len(criterios) == 0 {
		criterios = []string{"Criterio"}
	}
	e, err := domain.NuevoElementoBacklog(proyectoID, titulo, "Descripción", domain.PrioridadMedia, sp, criterios)
	if err != nil {
		t.Fatalf("crear elemento de prueba: %v", err)
	}
	return e
}

// Traza: HU-05 CA-05.10
func TestRepositorioBacklog_GuardarYBuscar_DevuelveLoMismo(t *testing.T) {
	ctx := context.Background()
	cinco := 5
	tests := []struct {
		nombre string
		sp     *int
	}{
		{"con story points", &cinco},
		{"sin estimar", nil},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			repo, ids := nuevoBacklogDePrueba(t, "Alfa")
			e := elementoDePrueba(t, ids[0], "Gestión de usuarios", tt.sp, "Tercero", "Primero", "Segundo")

			if err := repo.GuardarElemento(ctx, &e); err != nil {
				t.Fatalf("guardar: %v", err)
			}
			if e.ID <= 0 {
				t.Fatalf("ID = %d, se esperaba un ID asignado", e.ID)
			}
			got, err := repo.BuscarElementoPorID(ctx, e.ID)
			if err != nil {
				t.Fatalf("buscar: %v", err)
			}
			if !reflect.DeepEqual(got, e) {
				t.Errorf("se recuperó %+v, se esperaba %+v", got, e)
			}
		})
	}
}

// Traza: HU-05 CA-05.9 (RN-05.3 garantizada también por la base)
func TestRepositorioBacklog_TituloPorProyecto(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa", "Beta")
	login := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &login); err != nil {
		t.Fatal(err)
	}

	if existe, err := repo.ExisteTitulo(ctx, ids[0], domain.ClaveNombre("LOGIN")); err != nil || !existe {
		t.Errorf("ExisteTitulo en Alfa = %v, %v; se esperaba true", existe, err)
	}
	if existe, err := repo.ExisteTitulo(ctx, ids[1], domain.ClaveNombre("Login")); err != nil || existe {
		t.Errorf("ExisteTitulo en Beta = %v, %v; se esperaba false", existe, err)
	}

	otro := elementoDePrueba(t, ids[0], "login", nil)
	if err := repo.GuardarElemento(ctx, &otro); !errors.Is(err, domain.ErrTituloDuplicado) {
		t.Errorf("guardar duplicado: error = %v, se esperaba %v", err, domain.ErrTituloDuplicado)
	}
	enBeta := elementoDePrueba(t, ids[1], "Login", nil)
	if err := repo.GuardarElemento(ctx, &enBeta); err != nil {
		t.Errorf("mismo título en otro proyecto: %v", err)
	}
	if n, err := repo.ContarElementos(ctx, ids[0]); err != nil || n != 1 {
		t.Errorf("ContarElementos(Alfa) = %d, %v; se esperaba 1", n, err)
	}
}

// Traza: HU-05 CA-05.3
func TestRepositorioBacklog_ProyectoInexistente(t *testing.T) {
	ctx := context.Background()
	repo, _ := nuevoBacklogDePrueba(t)

	if existe, err := repo.ExisteProyecto(ctx, 42); err != nil || existe {
		t.Errorf("ExisteProyecto(42) = %v, %v; se esperaba false", existe, err)
	}
	e := elementoDePrueba(t, 42, "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("guardar en proyecto inexistente: error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}

func TestRepositorioBacklog_BuscarInexistente(t *testing.T) {
	repo, _ := nuevoBacklogDePrueba(t)
	if _, err := repo.BuscarElementoPorID(context.Background(), 999); !errors.Is(err, ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrElementoNoEncontrado)
	}
}
