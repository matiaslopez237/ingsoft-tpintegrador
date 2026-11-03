package persistencia

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// GuardarEstimacion guarda en una transacción la estimación que r tiene y la copia a los story points
// del elemento (RN-17.1). Las reglas se verifican en los mismos UPDATE, que SQLite ejecuta de forma
// atómica: la ronda tiene que seguir siendo la última, Cerrada y sin estimación (RN-17.3, RN-17.4), y
// el elemento no estar Hecho. Así, si dos pedidos registran a la vez, se guarda uno solo. Del elemento
// solo se actualiza story_points, para no pisar otros campos. Si algo falla, no se guarda nada.
// Implementa HU-17 (app.RepositorioEstimacion).
func (r *RepositorioRondas) GuardarEstimacion(ctx context.Context, ronda domain.Ronda) error {
	if ronda.EstimacionAcordada == nil { // el caso de uso siempre la registra antes de guardarla
		return fmt.Errorf("guardar estimación de la ronda %d: la ronda no la tiene", ronda.ID)
	}
	valor := *ronda.EstimacionAcordada

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback() // sin efecto si ya se hizo Commit

	res, err := tx.ExecContext(ctx,
		`UPDATE rondas SET estimacion_acordada = ?
		 WHERE id = ? AND estado = ? AND estimacion_acordada IS NULL
		   AND numero = (SELECT MAX(numero) FROM rondas AS otra WHERE otra.elemento_id = rondas.elemento_id)`,
		valor, ronda.ID, string(domain.RondaCerrada))
	if err != nil {
		return fmt.Errorf("guardar estimación de la ronda %d: %w", ronda.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar estimación de la ronda %d: %w", ronda.ID, err)
	}
	if n == 0 {
		tx.Rollback() // libera la conexión antes de releer
		return r.errorDeEstimacion(ctx, ronda.ID)
	}

	res, err = tx.ExecContext(ctx,
		`UPDATE backlog_elementos SET story_points = ?
		 WHERE estado <> ? AND id = (SELECT elemento_id FROM rondas WHERE id = ?)`,
		valor, string(domain.EstadoHecho), ronda.ID)
	if err != nil {
		return fmt.Errorf("guardar story points de la ronda %d: %w", ronda.ID, err)
	}
	if n, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("verificar story points de la ronda %d: %w", ronda.ID, err)
	}
	if n == 0 { // la ronda existe, y el elemento se borra con ella: lo único que falla es RN-17.3
		return domain.ErrElementoTerminado
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar estimación de la ronda %d: %w", ronda.ID, err)
	}
	return nil
}

// errorDeEstimacion determina, tras un UPDATE que no afectó ninguna fila, qué regla falló: relee la
// ronda (puede no existir) y sigue el orden de la tabla de la spec.
func (r *RepositorioRondas) errorDeEstimacion(ctx context.Context, id int64) error {
	actual, err := r.BuscarRonda(ctx, id)
	switch {
	case err != nil:
		return err
	case actual.Estado != domain.RondaCerrada:
		return domain.ErrEstimacionRondaAbierta
	case actual.EstimacionAcordada != nil:
		return domain.ErrEstimacionYaRegistrada
	}
	return domain.ErrRondaNoEsLaUltima
}
