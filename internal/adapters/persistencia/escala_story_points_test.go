package persistencia

import (
	"context"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: HU-13 CA-13.5 (RN-13.3): un elemento guardado antes de la escala con 4 story points se
// lee igual y, al actualizar otro campo, la base conserva los 4. No hizo falta persistencia nueva
// ni migrar datos: la columna sigue siendo un entero.
func TestRepositorioBacklog_ValorAnteriorALaEscala(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")
	if _, err := repo.db.Exec(`UPDATE backlog_elementos SET story_points = 4 WHERE id = ?`, e.ID); err != nil {
		t.Fatal(err)
	}

	viejo, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil || viejo.StoryPoints == nil || *viejo.StoryPoints != 4 {
		t.Fatalf("leer el elemento: %v, %v; se esperaba 4 story points", viejo.StoryPoints, err)
	}
	titulo := "Ingreso"
	c := domain.CambiosElemento{Titulo: &titulo}
	editado, err := viejo.Editar(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.ActualizarElemento(ctx, editado, c); err != nil {
		t.Fatal(err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil || got.Titulo != "Ingreso" || got.StoryPoints == nil || *got.StoryPoints != 4 {
		t.Errorf("después de editar el título: %q, %v, %v; se esperaba conservar los 4", got.Titulo, got.StoryPoints, err)
	}
}
