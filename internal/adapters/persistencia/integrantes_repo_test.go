package persistencia

import (
	"context"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioIntegrantes = (*RepositorioProyectos)(nil)

func proyectoGuardadoDePrueba(t *testing.T, repo *RepositorioProyectos) domain.Proyecto {
	t.Helper()
	p := proyectoDePrueba(t, "TP Integrador", nil)
	if err := repo.Guardar(context.Background(), &p); err != nil {
		t.Fatalf("guardar proyecto de prueba: %v", err)
	}
	return p
}

// Traza: HU-03 CA-03.1: agregar un integrante persiste y se puede recuperar.
func TestRepositorioProyectos_GuardarIntegrantes_AgregaYPersiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	if err := p.AgregarIntegrante("Ana", "Desarrolladora"); err != nil {
		t.Fatalf("agregar en memoria: %v", err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatalf("guardar integrantes: %v", err)
	}

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if len(got.Integrantes) != 1 || got.Integrantes[0] != (domain.Integrante{Nombre: "Ana", Rol: "Desarrolladora"}) {
		t.Errorf("Integrantes = %+v, se esperaba 1 integrante Ana/Desarrolladora", got.Integrantes)
	}
}

// Traza: HU-03 CA-03.3: quitar un integrante persiste la lista actualizada.
func TestRepositorioProyectos_GuardarIntegrantes_QuitaYPersiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	if err := p.AgregarIntegrante("Ana", "Desarrolladora"); err != nil {
		t.Fatal(err)
	}
	if err := p.AgregarIntegrante("Beto", "Tester"); err != nil {
		t.Fatal(err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatalf("guardar con 2 integrantes: %v", err)
	}

	if err := p.QuitarIntegrante("Ana"); err != nil {
		t.Fatal(err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatalf("guardar tras quitar: %v", err)
	}

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if len(got.Integrantes) != 1 || got.Integrantes[0].Nombre != "Beto" {
		t.Errorf("Integrantes = %+v, se esperaba solo Beto", got.Integrantes)
	}
}

// Traza: HU-03 caso límite: quitar al único integrante deja la lista vacía en la base.
func TestRepositorioProyectos_GuardarIntegrantes_UnicoRestante_QuedaVacio(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	if err := p.AgregarIntegrante("Ana", "Desarrolladora"); err != nil {
		t.Fatal(err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatal(err)
	}
	if err := p.QuitarIntegrante("Ana"); err != nil {
		t.Fatal(err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatal(err)
	}

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if len(got.Integrantes) != 0 {
		t.Errorf("Integrantes = %+v, se esperaba lista vacía", got.Integrantes)
	}
}

// Traza: RN-03.3 también se respeta al recargar desde la base (vía domain.Proyecto.AgregarIntegrante).
func TestRepositorioProyectos_GuardarIntegrantes_DuplicadoTrasRecargar_RetornaError(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	if err := p.AgregarIntegrante("Ana", "Desarrolladora"); err != nil {
		t.Fatal(err)
	}
	if err := repo.GuardarIntegrantes(ctx, &p); err != nil {
		t.Fatal(err)
	}

	recargado, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if err := recargado.AgregarIntegrante("ana", "Tester"); err != domain.ErrIntegranteDuplicado {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteDuplicado)
	}
}
