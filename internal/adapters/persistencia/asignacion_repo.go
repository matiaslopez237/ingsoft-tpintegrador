package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// BuscarSprint consulta el sprint por ID en la tabla de HU-08, sin depender de domain.Sprint
// (spec HU-09, sección 5). Implementa HU-09 (app.RepositorioAsignacion).
func (r *RepositorioBacklog) BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error) {
	var estado string
	err = r.db.QueryRowContext(ctx, `SELECT proyecto_id, estado FROM sprints WHERE id = ?`, id).Scan(&proyectoID, &estado)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, app.ErrSprintNoEncontrado
	}
	if err != nil {
		return 0, false, fmt.Errorf("buscar sprint %d: %w", id, err)
	}
	return proyectoID, estado == string(domain.SprintAbierto), nil
}

// GuardarAsignacion persiste el SprintID (o su ausencia, si es nil) de un elemento ya existente.
// Implementa HU-09 (app.RepositorioAsignacion).
func (r *RepositorioBacklog) GuardarAsignacion(ctx context.Context, e *domain.ElementoBacklog) error {
	var sprintID sql.NullInt64
	if e.SprintID != nil {
		sprintID = sql.NullInt64{Int64: *e.SprintID, Valid: true}
	}
	if _, err := r.db.ExecContext(ctx, `UPDATE backlog_elementos SET sprint_id = ? WHERE id = ?`, sprintID, e.ID); err != nil {
		return fmt.Errorf("guardar asignación del elemento %d: %w", e.ID, err)
	}
	return nil
}
