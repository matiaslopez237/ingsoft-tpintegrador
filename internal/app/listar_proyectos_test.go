package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type listadorFalso struct {
	proyectos []domain.Proyecto
	err       error
}

func (l *listadorFalso) ListarProyectos(context.Context) ([]domain.Proyecto, error) {
	return l.proyectos, l.err
}

func nombres(ps []domain.Proyecto) []string {
	var n []string
	for _, p := range ps {
		n = append(n, p.Nombre)
	}
	return n
}

// Traza: HU-30 CA-30.2 (RN-30.1)
func TestListarProyectos_OrdenaPorNombreSinDistinguirMayusculas(t *testing.T) {
	repo := &listadorFalso{proyectos: []domain.Proyecto{
		{ID: 1, Nombre: "beta"},
		{ID: 2, Nombre: "Gamma"},
		{ID: 3, Nombre: "Alfa"},
		{ID: 4, Nombre: "ñandú"},
		{ID: 5, Nombre: "Ñu"},
	}}

	got, err := NuevoListarProyectos(repo).Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := []string{"Alfa", "beta", "Gamma", "ñandú", "Ñu"}
	if !reflect.DeepEqual(nombres(got), want) {
		t.Errorf("orden = %v, se esperaba %v", nombres(got), want)
	}
}

// Traza: HU-30 RN-30.1: a igual clave de nombre, decide el ID.
func TestListarProyectos_AIgualNombreOrdenaPorID(t *testing.T) {
	repo := &listadorFalso{proyectos: []domain.Proyecto{
		{ID: 9, Nombre: "alfa"},
		{ID: 2, Nombre: "Alfa"},
	}}

	got, err := NuevoListarProyectos(repo).Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got[0].ID != 2 || got[1].ID != 9 {
		t.Errorf("IDs = %d, %d; se esperaba 2, 9", got[0].ID, got[1].ID)
	}
}

// Traza: HU-30 CA-30.2 (lista vacía)
func TestListarProyectos_SinProyectos_DevuelveListaVacia(t *testing.T) {
	got, err := NuevoListarProyectos(&listadorFalso{}).Ejecutar(context.Background())
	if err != nil || len(got) != 0 {
		t.Errorf("se obtuvo %v, %v; se esperaba una lista vacía sin error", got, err)
	}
}

// El caso de uso no reordena la lista que le devolvió el repositorio.
func TestListarProyectos_NoModificaLaListaDelRepositorio(t *testing.T) {
	original := []domain.Proyecto{{ID: 1, Nombre: "Beta"}, {ID: 2, Nombre: "Alfa"}}
	repo := &listadorFalso{proyectos: original}

	if _, err := NuevoListarProyectos(repo).Ejecutar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if original[0].Nombre != "Beta" {
		t.Errorf("la lista del repositorio quedó reordenada: %v", nombres(original))
	}
}

func TestListarProyectos_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	_, err := NuevoListarProyectos(&listadorFalso{err: errBase}).Ejecutar(context.Background())
	if !errors.Is(err, errBase) {
		t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
	}
}
