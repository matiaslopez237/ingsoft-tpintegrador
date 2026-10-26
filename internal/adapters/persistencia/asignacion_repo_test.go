package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioAsignacion = (*RepositorioBacklog)(nil)

// cerrarSprintDePrueba marca un sprint como Cerrado directamente en la base: HU-11 (cerrar sprint)
// todavía no existe, así que no hay otra forma de dejar un sprint cerrado para probar.
func cerrarSprintDePrueba(t *testing.T, db *sql.DB, sprintID int64) {
	t.Helper()
	if _, err := db.Exec(`UPDATE sprints SET estado = ? WHERE id = ?`, string(domain.SprintCerrado), sprintID); err != nil {
		t.Fatalf("cerrar sprint de prueba: %v", err)
	}
}

// Traza: HU-09, sección 7 (sprint encontrado, abierto)
func TestRepositorioBacklog_BuscarSprint_Abierto(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	repoSprints := NuevoRepositorioSprints(repo.db)
	s := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))

	proyectoID, abierto, err := repo.BuscarSprint(ctx, s.ID)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if proyectoID != ids[0] || !abierto {
		t.Errorf("proyectoID = %d, abierto = %v; se esperaba %d, true", proyectoID, abierto, ids[0])
	}
}

// Traza: HU-09 CA-09.8
func TestRepositorioBacklog_BuscarSprint_Cerrado(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	repoSprints := NuevoRepositorioSprints(repo.db)
	s := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))
	cerrarSprintDePrueba(t, repo.db, s.ID)

	_, abierto, err := repo.BuscarSprint(ctx, s.ID)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if abierto {
		t.Error("abierto = true, se esperaba false (sprint cerrado)")
	}
}

// Traza: HU-09, sección 7: sprint inexistente
func TestRepositorioBacklog_BuscarSprint_Inexistente_RetornaError(t *testing.T) {
	repo, _ := nuevoBacklogDePrueba(t, "Alfa")

	_, _, err := repo.BuscarSprint(context.Background(), 999)

	if !errors.Is(err, app.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, app.ErrSprintNoEncontrado)
	}
}

// Traza: HU-09 CA-09.1
func TestRepositorioBacklog_GuardarAsignacion_Asigna(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	repoSprints := NuevoRepositorioSprints(repo.db)
	s := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))

	e.SprintID = &s.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("guardar asignación: %v", err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.SprintID == nil || *got.SprintID != s.ID {
		t.Errorf("SprintID = %v, se esperaba %d", got.SprintID, s.ID)
	}
}

// Traza: HU-09 CA-09.4
func TestRepositorioBacklog_GuardarAsignacion_Desasigna(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	repoSprints := NuevoRepositorioSprints(repo.db)
	s := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))
	e.SprintID = &s.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	e.SprintID = nil
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("guardar desasignación: %v", err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.SprintID != nil {
		t.Errorf("SprintID = %v, se esperaba nil", got.SprintID)
	}
}
