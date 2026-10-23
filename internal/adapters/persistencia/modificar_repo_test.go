package persistencia

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioModificacion = (*RepositorioProyectos)(nil)

func guardarDePrueba(t *testing.T, repo *RepositorioProyectos, nombre string) domain.Proyecto {
	t.Helper()
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	p := proyectoDePrueba(t, nombre, &fin)
	if err := repo.Guardar(context.Background(), &p); err != nil {
		t.Fatalf("guardar %q: %v", nombre, err)
	}
	return p
}

// Traza: HU-02 CA-02.10, CA-02.5: la modificación queda persistida al recuperarla, con el mismo ID,
// incluso quitando la fecha de fin.
func TestRepositorioProyectos_Actualizar_Persiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, repo, "Alfa")

	p.Nombre = "Alfa Renovado"
	p.FechaInicio = time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	p.FechaFin = nil
	p.Estado = domain.EstadoFinalizado
	if err := repo.Actualizar(ctx, p); err != nil {
		t.Fatalf("actualizar: %v", err)
	}

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.ID != p.ID || got.Nombre != "Alfa Renovado" || !got.FechaInicio.Equal(p.FechaInicio) ||
		got.FechaFin != nil || got.Estado != domain.EstadoFinalizado {
		t.Errorf("recuperado = %+v, se esperaba %+v", got, p)
	}
	if existe, _ := repo.ExisteNombre(ctx, domain.ClaveNombre("Alfa Renovado")); !existe {
		t.Error("la clave del nombre nuevo no quedó guardada")
	}
}

// Traza: HU-02 CA-02.12, RN-02.10: actualizar no borra los integrantes ni sus roles.
func TestRepositorioProyectos_Actualizar_ConservaLosIntegrantes(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, repo, "Alfa")
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")
	if err := repo.AgregarRol(ctx, p.ID, "Ana", "Tester"); err != nil {
		t.Fatalf("agregar rol: %v", err)
	}

	p.Nombre = "Alfa Renovado"
	if err := repo.Actualizar(ctx, p); err != nil {
		t.Fatalf("actualizar: %v", err)
	}

	if got, want := rolesDe(t, repo, p.ID, "Ana"), []string{"Desarrolladora", "Tester"}; !slices.Equal(got, want) {
		t.Errorf("roles de Ana = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-02 RN-02.3: renombrar el mismo proyecto con otras mayúsculas no choca consigo mismo.
func TestRepositorioProyectos_Actualizar_MismoNombreConOtrasMayusculas(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, repo, "Alfa")

	p.Nombre = "ALFA"
	if err := repo.Actualizar(ctx, p); err != nil {
		t.Fatalf("actualizar: %v", err)
	}
}

// Traza: HU-02 RN-02.3: si otro pedido tomó el nombre en el medio, la base lo rechaza.
func TestRepositorioProyectos_Actualizar_NombreDeOtro_RetornaDuplicado(t *testing.T) {
	repo := nuevoRepoDePrueba(t)
	guardarDePrueba(t, repo, "Alfa")
	beta := guardarDePrueba(t, repo, "Beta")

	beta.Nombre = "alfa"
	err := repo.Actualizar(context.Background(), beta)

	if !errors.Is(err, domain.ErrNombreDuplicado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrNombreDuplicado)
	}
}

// Traza: HU-02 CA-02.9: si el proyecto se borró en el medio, no hay nada que actualizar.
func TestRepositorioProyectos_Actualizar_Inexistente_RetornaNoEncontrado(t *testing.T) {
	repo := nuevoRepoDePrueba(t)
	p := proyectoDePrueba(t, "Fantasma", nil)
	p.ID = 99

	err := repo.Actualizar(context.Background(), p)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}
