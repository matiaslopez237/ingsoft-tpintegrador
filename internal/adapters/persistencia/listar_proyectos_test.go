package persistencia

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto de HU-30 (se verifica al compilar).
var _ app.ListadorProyectos = (*RepositorioProyectos)(nil)

// Traza: HU-30 CA-30.2
func TestRepositorioProyectos_ListarProyectos_DevuelveTodosConSusDatos(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	alfa := proyectoDePrueba(t, "Alfa", &fin)
	beta := proyectoDePrueba(t, "Beta", nil)
	for _, p := range []*domain.Proyecto{&alfa, &beta} {
		if err := repo.Guardar(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	got, err := repo.ListarProyectos(ctx)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	porID := map[int64]domain.Proyecto{}
	for _, p := range got {
		porID[p.ID] = p
	}
	if len(got) != 2 || !reflect.DeepEqual(porID[alfa.ID], alfa) || !reflect.DeepEqual(porID[beta.ID], beta) {
		t.Errorf("se obtuvo %+v, se esperaban %+v y %+v", got, alfa, beta)
	}
}

// Traza: HU-30 CA-30.2 (lista vacía)
func TestRepositorioProyectos_ListarProyectos_SinProyectos(t *testing.T) {
	got, err := nuevoRepoDePrueba(t).ListarProyectos(context.Background())
	if err != nil || len(got) != 0 {
		t.Errorf("se obtuvo %v, %v; se esperaba una lista vacía sin error", got, err)
	}
}
