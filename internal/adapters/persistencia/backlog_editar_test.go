package persistencia

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir los puertos de HU-06 (se verifica al compilar).
var (
	_ app.RepositorioEdicionBacklog     = (*RepositorioBacklog)(nil)
	_ app.RepositorioEliminacionBacklog = (*RepositorioBacklog)(nil)
)

func guardarElemento(t *testing.T, repo *RepositorioBacklog, proyectoID int64, titulo string) domain.ElementoBacklog {
	t.Helper()
	cinco := 5
	e := elementoDePrueba(t, proyectoID, titulo, &cinco, "Uno", "Dos")
	if err := repo.GuardarElemento(context.Background(), &e); err != nil {
		t.Fatal(err)
	}
	return e
}

// Traza: HU-06 CA-06.2
func TestRepositorioBacklog_ActualizarElemento_TodosLosCampos(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")
	criterios := []string{"Nuevo uno", "Nuevo dos", "Nuevo tres"}
	hecho := domain.EstadoHecho
	c := domain.CambiosElemento{Titulo: &criterios[0], CambiarStoryPoints: true, Criterios: &criterios, Estado: &hecho}
	editado, err := e.Editar(c)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.ActualizarElemento(ctx, editado, c); err != nil {
		t.Fatalf("actualizar: %v", err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, editado) {
		t.Errorf("se recuperó %+v, se esperaba %+v", got, editado)
	}
}

// Traza: HU-06 sección 5: dos ediciones de campos distintos no se pisan.
func TestRepositorioBacklog_ActualizarElemento_SoloLosCamposEnviados(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")

	// Las dos personas parten del mismo elemento: una cambia la prioridad y la otra el estado.
	baja, enProgreso := domain.PrioridadBaja, domain.EstadoEnProgreso
	c1 := domain.CambiosElemento{Prioridad: &baja}
	c2 := domain.CambiosElemento{Estado: &enProgreso}
	e1, _ := e.Editar(c1)
	e2, _ := e.Editar(c2)
	if err := repo.ActualizarElemento(ctx, e1, c1); err != nil {
		t.Fatal(err)
	}
	if err := repo.ActualizarElemento(ctx, e2, c2); err != nil {
		t.Fatal(err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Prioridad != domain.PrioridadBaja || got.Estado != domain.EstadoEnProgreso || !reflect.DeepEqual(got.Criterios, e.Criterios) {
		t.Errorf("se perdió un cambio: %+v", got)
	}
}

// Traza: HU-06 CA-06.5, CA-06.7
func TestRepositorioBacklog_ActualizarElemento_Errores(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	login := guardarElemento(t, repo, ids[0], "Login")
	guardarElemento(t, repo, ids[0], "Registro")

	nuevo := "registro"
	c := domain.CambiosElemento{Titulo: &nuevo}
	editado, _ := login.Editar(c)
	if err := repo.ActualizarElemento(ctx, editado, c); !errors.Is(err, domain.ErrTituloDuplicado) {
		t.Errorf("título duplicado: error = %v, se esperaba %v", err, domain.ErrTituloDuplicado)
	}

	editado.ID = 999
	otro := "Otro"
	if err := repo.ActualizarElemento(ctx, editado, domain.CambiosElemento{Titulo: &otro}); !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("inexistente: error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-06 CA-06.8
func TestRepositorioBacklog_EliminarElemento_BorraElElementoYSusCriterios(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")

	if err := repo.EliminarElemento(ctx, e.ID); err != nil {
		t.Fatalf("eliminar: %v", err)
	}

	if _, err := repo.BuscarElementoPorID(ctx, e.ID); !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("buscar eliminado: error = %v", err)
	}
	var criterios int
	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM backlog_criterios WHERE elemento_id = ?`, e.ID).Scan(&criterios); err != nil || criterios != 0 {
		t.Errorf("quedaron %d criterios (err %v)", criterios, err)
	}
	if err := repo.EliminarElemento(ctx, e.ID); !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("eliminar dos veces: error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-06 CA-06.9 (RN-06.7, garantizada también por la clave foránea)
func TestRepositorioBacklog_ElementoConEsfuerzo(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")
	if tiene, err := repo.TieneEsfuerzo(ctx, e.ID); err != nil || tiene {
		t.Fatalf("TieneEsfuerzo antes de registrar = %v, %v", tiene, err)
	}

	r, err := domain.NuevoRegistroEsfuerzo(ids[0], e.ID, "Ana", time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		"Desarrollo", 2, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if err := NuevoRepositorioEsfuerzo(repo.db).GuardarRegistro(ctx, &r); err != nil {
		t.Fatal(err)
	}

	if tiene, err := repo.TieneEsfuerzo(ctx, e.ID); err != nil || !tiene {
		t.Errorf("TieneEsfuerzo = %v, %v; se esperaba true", tiene, err)
	}
	if err := repo.EliminarElemento(ctx, e.ID); !errors.Is(err, domain.ErrElementoConEsfuerzo) {
		t.Errorf("eliminar con esfuerzo: error = %v, se esperaba %v", err, domain.ErrElementoConEsfuerzo)
	}
	if _, err := repo.BuscarElementoPorID(ctx, e.ID); err != nil {
		t.Errorf("el elemento debería seguir existiendo: %v", err)
	}
}
