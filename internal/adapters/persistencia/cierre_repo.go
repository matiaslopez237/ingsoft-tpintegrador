package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// BuscarSprintPorID devuelve el sprint con ese ID, o domain.ErrSprintNoEncontrado si no existe.
// Implementa HU-11 (app.RepositorioCierre).
func (r *RepositorioSprints) BuscarSprintPorID(ctx context.Context, id int64) (domain.Sprint, error) {
	s, err := escanearSprint(r.db.QueryRowContext(ctx, `SELECT `+columnasSprint+` FROM sprints WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("buscar sprint %d: %w", id, err)
	}
	return s, nil
}

// CerrarSprint pasa el sprint a Cerrado solo si sigue Abierto, guarda su resumen de puntos y
// devuelve al backlog las historias sin terminar (RN-11.3, RN-11.6, con domain.ResumirCierre).
// Todo va en una transacción que empieza por el UPDATE del estado: la condición va en ese UPDATE,
// que SQLite ejecuta de forma atómica, así que si dos pedidos cierran el mismo sprint a la vez uno
// lo cierra y el otro no afecta ninguna fila (errorDeCierre relee el sprint para informar por qué).
// Desde ese UPDATE la transacción tiene el lock de escritura: una asignación (HU-09) o un completado
// (HU-10) simultáneos esperan, así que no pueden caer entre la lectura de las historias y la
// devolución. Escribe siempre el estado Cerrado (RN-11.1), de s usa solo el ID, como define el
// puerto, y al terminar le deja los puntos registrados. Implementa HU-11 (app.RepositorioCierre,
// sección 5, P-8).
func (r *RepositorioSprints) CerrarSprint(ctx context.Context, s *domain.Sprint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cerrar sprint %d: %w", s.ID, err)
	}
	defer tx.Rollback() // sin efecto si ya se hizo Commit

	res, err := tx.ExecContext(ctx, `UPDATE sprints SET estado = ? WHERE id = ? AND estado = ?`,
		string(domain.SprintCerrado), s.ID, string(domain.SprintAbierto))
	if err != nil {
		return fmt.Errorf("cerrar sprint %d: %w", s.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar cierre del sprint %d: %w", s.ID, err)
	}
	if n == 0 {
		tx.Rollback() // libera la conexión antes de releer
		return r.errorDeCierre(ctx, s.ID)
	}

	historias, err := historiasDelSprint(ctx, tx, s.ID)
	if err != nil {
		return err
	}
	resumen := domain.ResumirCierre(historias)
	if _, err := tx.ExecContext(ctx, `UPDATE sprints SET puntos_planificados = ?, puntos_completados = ? WHERE id = ?`,
		resumen.PuntosPlanificados, resumen.PuntosCompletados, s.ID); err != nil {
		return fmt.Errorf("guardar resumen del sprint %d: %w", s.ID, err)
	}
	for _, id := range resumen.DevueltasAlBacklog {
		if _, err := tx.ExecContext(ctx, `UPDATE backlog_elementos SET sprint_id = NULL WHERE id = ?`, id); err != nil {
			return fmt.Errorf("devolver al backlog el elemento %d: %w", id, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar cierre del sprint %d: %w", s.ID, err)
	}
	s.PuntosPlanificados, s.PuntosCompletados = resumen.PuntosPlanificados, resumen.PuntosCompletados
	return nil
}

// historiasDelSprint lee, dentro de la transacción del cierre, lo que domain.ResumirCierre necesita
// de cada historia asignada al sprint: ID, estado y story points.
func historiasDelSprint(ctx context.Context, tx *sql.Tx, sprintID int64) ([]domain.ElementoBacklog, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, estado, story_points FROM backlog_elementos WHERE sprint_id = ? ORDER BY id`, sprintID)
	if err != nil {
		return nil, fmt.Errorf("leer historias del sprint %d: %w", sprintID, err)
	}
	defer rows.Close()

	var historias []domain.ElementoBacklog
	for rows.Next() {
		var (
			e      domain.ElementoBacklog
			estado string
			sp     sql.NullInt64
		)
		if err := rows.Scan(&e.ID, &estado, &sp); err != nil {
			return nil, fmt.Errorf("leer historia del sprint %d: %w", sprintID, err)
		}
		e.Estado = domain.EstadoElemento(estado)
		if sp.Valid {
			n := int(sp.Int64)
			e.StoryPoints = &n
		}
		historias = append(historias, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer historias del sprint %d: %w", sprintID, err)
	}
	return historias, nil
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
