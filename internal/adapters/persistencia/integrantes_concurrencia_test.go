package persistencia

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// repoConLecturaVieja simula pedidos simultáneos: todos leyeron el proyecto antes de que
// cualquiera guardara, así que BuscarPorID devuelve siempre la misma foto tomada al principio.
// Las escrituras van al repositorio real.
type repoConLecturaVieja struct {
	*RepositorioProyectos
	foto domain.Proyecto
}

func (r repoConLecturaVieja) BuscarPorID(context.Context, int64) (domain.Proyecto, error) {
	p := r.foto
	p.Integrantes = slices.Clone(r.foto.Integrantes)
	return p, nil
}

func nombresIntegrantes(t *testing.T, repo *RepositorioProyectos, id int64) []string {
	t.Helper()
	p, err := repo.BuscarPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	var nombres []string
	for _, i := range p.Integrantes {
		nombres = append(nombres, i.Nombre)
	}
	slices.Sort(nombres)
	return nombres
}

// Traza: HU-03 sección 5 (pedidos simultáneos): dos altas a la vez no se pisan.
func TestIntegrantes_DosAltasSimultaneas_QuedanAmbas(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	vieja := repoConLecturaVieja{RepositorioProyectos: repo, foto: p}

	if err := app.NuevoAgregarIntegrante(vieja).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: "Ana", Rol: "Desarrolladora"}); err != nil {
		t.Fatalf("agregar Ana: %v", err)
	}
	if err := app.NuevoAgregarIntegrante(vieja).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: "Beto", Rol: "Tester"}); err != nil {
		t.Fatalf("agregar Beto: %v", err)
	}

	if got, want := nombresIntegrantes(t, repo, p.ID), []string{"Ana", "Beto"}; !slices.Equal(got, want) {
		t.Errorf("integrantes = %v, se esperaba %v (un alta pisó a la otra)", got, want)
	}
}

// Traza: HU-03 sección 5 (pedidos simultáneos): una baja y un alta a la vez no se pisan.
func TestIntegrantes_BajaYAltaSimultaneas_NoReviveElQuitado(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	for _, nombre := range []string{"Ana", "Beto"} {
		if err := app.NuevoAgregarIntegrante(repo).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: nombre, Rol: "Integrante"}); err != nil {
			t.Fatalf("preparar %s: %v", nombre, err)
		}
	}
	foto, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	vieja := repoConLecturaVieja{RepositorioProyectos: repo, foto: foto}

	if err := app.NuevoQuitarIntegrante(vieja).Ejecutar(ctx, app.DatosQuitarIntegrante{ProyectoID: p.ID, Nombre: "Ana"}); err != nil {
		t.Fatalf("quitar Ana: %v", err)
	}
	if err := app.NuevoAgregarIntegrante(vieja).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: "Caro", Rol: "Tester"}); err != nil {
		t.Fatalf("agregar Caro: %v", err)
	}

	if got, want := nombresIntegrantes(t, repo, p.ID), []string{"Beto", "Caro"}; !slices.Equal(got, want) {
		t.Errorf("integrantes = %v, se esperaba %v (el alta revivió a quien se había quitado)", got, want)
	}
}

// Traza: HU-03 sección 5 y RN-03.3: si dos pedidos agregan a la misma persona a la vez, los dos
// pasan la validación del dominio; la base rechaza al segundo y se informa como duplicado.
func TestIntegrantes_MismaAltaSimultanea_SegundaEsDuplicado(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	vieja := repoConLecturaVieja{RepositorioProyectos: repo, foto: p}

	if err := app.NuevoAgregarIntegrante(vieja).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: "Ana", Rol: "Desarrolladora"}); err != nil {
		t.Fatalf("primera alta: %v", err)
	}
	err := app.NuevoAgregarIntegrante(vieja).Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: "ana", Rol: "Tester"})

	if !errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteDuplicado)
	}
	if got, want := nombresIntegrantes(t, repo, p.ID), []string{"Ana"}; !slices.Equal(got, want) {
		t.Errorf("integrantes = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-03 sección 5 y CA-03.4: si dos pedidos quitan a la misma persona a la vez,
// el segundo no encuentra nada que borrar y se informa como no encontrado.
func TestIntegrantes_MismaBajaSimultanea_SegundaEsNoEncontrado(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := proyectoGuardadoDePrueba(t, repo)
	agregarIntegranteDePrueba(t, repo, p.ID, "Ana", "Desarrolladora")
	foto, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	vieja := repoConLecturaVieja{RepositorioProyectos: repo, foto: foto}

	if err := app.NuevoQuitarIntegrante(vieja).Ejecutar(ctx, app.DatosQuitarIntegrante{ProyectoID: p.ID, Nombre: "Ana"}); err != nil {
		t.Fatalf("primera baja: %v", err)
	}
	err = app.NuevoQuitarIntegrante(vieja).Ejecutar(ctx, app.DatosQuitarIntegrante{ProyectoID: p.ID, Nombre: "Ana"})

	if !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteNoEncontrado)
	}
}
