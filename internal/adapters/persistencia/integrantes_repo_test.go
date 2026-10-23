package persistencia

import (
	"context"
	"errors"
	"slices"
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

func agregarIntegranteDePrueba(t *testing.T, repo *RepositorioProyectos, proyectoID int64, nombre, rol string) {
	t.Helper()
	if err := repo.AgregarIntegrante(context.Background(), proyectoID, domain.Integrante{Nombre: nombre, Roles: []string{rol}}); err != nil {
		t.Fatalf("agregar integrante %q: %v", nombre, err)
	}
}

// Traza: HU-03 CA-03.1: agregar un integrante persiste y se puede recuperar.
func TestRepositorioProyectos_AgregarIntegrante_Persiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if len(got.Integrantes) != 1 || got.Integrantes[0].Nombre != "Ana" || !slices.Equal(got.Integrantes[0].Roles, []string{"Desarrolladora"}) {
		t.Errorf("Integrantes = %+v, se esperaba 1 integrante Ana/Desarrolladora", got.Integrantes)
	}
}

// Traza: HU-03 CA-03.3: quitar un integrante persiste la lista actualizada.
func TestRepositorioProyectos_QuitarIntegrante_Persiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")
	agregarIntegranteDePrueba(t, repo, p.ID, "Beto", "Tester")

	if err := repo.QuitarIntegrante(ctx, p.ID, "Ana"); err != nil {
		t.Fatalf("quitar: %v", err)
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
func TestRepositorioProyectos_QuitarIntegrante_UnicoRestante_QuedaVacio(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")

	if err := repo.QuitarIntegrante(ctx, p.ID, "Ana"); err != nil {
		t.Fatalf("quitar: %v", err)
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
func TestRepositorioProyectos_DuplicadoTrasRecargar_RetornaError(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")

	recargado, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if err := recargado.AgregarIntegrante("ana", "Tester"); err != domain.ErrIntegranteDuplicado {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteDuplicado)
	}
}

func rolesDe(t *testing.T, repo *RepositorioProyectos, proyectoID int64, nombre string) []string {
	t.Helper()
	p, err := repo.BuscarPorID(context.Background(), proyectoID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	for _, i := range p.Integrantes {
		if i.Nombre == nombre {
			return i.Roles
		}
	}
	t.Fatalf("%q no está en el proyecto: %+v", nombre, p.Integrantes)
	return nil
}

// Traza: HU-03 CA-03.7, RN-03.5: los roles sumados se persisten en el orden en que se agregaron.
func TestRepositorioProyectos_AgregarRol_Persiste(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")

	if err := repo.AgregarRol(ctx, p.ID, "ana", "Tester"); err != nil {
		t.Fatalf("agregar rol: %v", err)
	}

	if got, want := rolesDe(t, repo, p.ID, "Ana"), []string{"Desarrolladora", "Tester"}; !slices.Equal(got, want) {
		t.Errorf("roles = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-03 sección 5 y RN-03.6: si otro pedido ya sumó el mismo rol, la base lo rechaza.
func TestRepositorioProyectos_AgregarRol_Duplicado_RetornaError(t *testing.T) {
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")

	err := repo.AgregarRol(context.Background(), p.ID, "Ana", "desarrolladora")

	if !errors.Is(err, domain.ErrRolDuplicado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrRolDuplicado)
	}
}

// Traza: HU-03 sección 5 y CA-03.9: si otro pedido quitó al integrante en el medio, no hay a quién sumarle el rol.
func TestRepositorioProyectos_AgregarRol_IntegranteInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)

	err := repo.AgregarRol(context.Background(), p.ID, "Caro", "Tester")

	if !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteNoEncontrado)
	}
}

// Traza: HU-03 CA-03.3 (caso límite): quitar a un integrante borra todos sus roles; si vuelve a
// entrar, empieza solo con el rol nuevo.
func TestRepositorioProyectos_QuitarIntegrante_BorraSusRoles(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")
	if err := repo.AgregarRol(ctx, p.ID, "Ana", "Tester"); err != nil {
		t.Fatalf("agregar rol: %v", err)
	}

	if err := repo.QuitarIntegrante(ctx, p.ID, "Ana"); err != nil {
		t.Fatalf("quitar: %v", err)
	}
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Scrum Master")

	if got, want := rolesDe(t, repo, p.ID, "Ana"), []string{"Scrum Master"}; !slices.Equal(got, want) {
		t.Errorf("roles = %v, se esperaba %v", got, want)
	}
}
