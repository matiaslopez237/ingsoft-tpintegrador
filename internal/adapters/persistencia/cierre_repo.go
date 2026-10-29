package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// BuscarSprintPorID devuelve el sprint con ese ID, o domain.ErrSprintNoEncontrado si no existe.
// Implementa HU-11 (app.RepositorioCierre).
func (r *RepositorioSprints) BuscarSprintPorID(ctx context.Context, id int64) (domain.Sprint, error) {
	var (
		s                   = domain.Sprint{ID: id}
		inicio, fin, estado string
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado FROM sprints WHERE id = ?`, id).
		Scan(&s.ProyectoID, &s.Numero, &inicio, &fin, &s.Objetivo, &estado)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("buscar sprint %d: %w", id, err)
	}
	if s.FechaInicio, err = time.Parse(time.DateOnly, inicio); err != nil {
		return domain.Sprint{}, fmt.Errorf("fecha de inicio guardada inválida: %w", err)
	}
	if s.FechaFin, err = time.Parse(time.DateOnly, fin); err != nil {
		return domain.Sprint{}, fmt.Errorf("fecha de fin guardada inválida: %w", err)
	}
	s.Estado = domain.EstadoSprint(estado)
	return s, nil
}

// CerrarSprint pasa el sprint a Cerrado solo si sigue Abierto. La condición va en el mismo UPDATE,
// que SQLite ejecuta de forma atómica: si dos pedidos cierran el mismo sprint a la vez, uno lo
// cierra y el otro no afecta ninguna fila. En ese caso errorDeCierre relee el sprint para informar
// por qué. Implementa HU-11 (app.RepositorioCierre, sección 5, P-8).
func (r *RepositorioSprints) CerrarSprint(ctx context.Context, s *domain.Sprint) error {
	res, err := r.db.ExecContext(ctx, `UPDATE sprints SET estado = ? WHERE id = ? AND estado = ?`,
		string(domain.SprintCerrado), s.ID, string(domain.SprintAbierto))
	if err != nil {
		return fmt.Errorf("cerrar sprint %d: %w", s.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar cierre del sprint %d: %w", s.ID, err)
	}
	if n == 0 {
		return r.errorDeCierre(ctx, s.ID)
	}
	return nil
}

// errorDeCierre determina, tras un UPDATE que no afectó ninguna fila, por qué no se cerró el
// sprint: no existe (domain.ErrSprintNoEncontrado) o ya estaba Cerrado (domain.ErrSprintYaCerrado).
func (r *RepositorioSprints) errorDeCierre(ctx context.Context, id int64) error {
	actual, err := r.BuscarSprintPorID(ctx, id)
	if err != nil {
		return err
	}
	if actual.Estado == domain.SprintCerrado {
		return domain.ErrSprintYaCerrado
	}
	return fmt.Errorf("cerrar sprint %d: estado guardado inesperado %q", id, actual.Estado)
}
