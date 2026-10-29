package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// BuscarSprint consulta el sprint por ID en la tabla de HU-08, sin depender de domain.Sprint
// (spec HU-09, sección 5). Implementa HU-09 (app.RepositorioAsignacion).
func (r *RepositorioBacklog) BuscarSprint(ctx context.Context, id int64) (proyectoID int64, abierto bool, err error) {
	var estado string
	err = r.db.QueryRowContext(ctx, `SELECT proyecto_id, estado FROM sprints WHERE id = ?`, id).Scan(&proyectoID, &estado)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, domain.ErrSprintNoEncontrado
	}
	if err != nil {
		return 0, false, fmt.Errorf("buscar sprint %d: %w", id, err)
	}
	return proyectoID, estado == string(domain.SprintAbierto), nil
}

// GuardarAsignacion persiste el SprintID (o su ausencia, si es nil) de un elemento ya existente.
// RN-09.2 a RN-09.4 se verifican en el propio UPDATE, que SQLite ejecuta de forma atómica: dos
// pedidos simultáneos no pueden asignar el mismo elemento a dos sprints Abiertos distintos, ni
// desasignarlo (o reasignarlo) si su sprint se cerró en el medio. Si el UPDATE no afecta ninguna
// fila, errorDeAsignacion relee el estado actual para informar cuál regla falló.
// Implementa HU-09 (app.RepositorioAsignacion, RN-09.2 a RN-09.4).
func (r *RepositorioBacklog) GuardarAsignacion(ctx context.Context, e *domain.ElementoBacklog) error {
	abierto := string(domain.SprintAbierto)
	var res sql.Result
	var err error
	if e.SprintID != nil {
		// El elemento se puede mover a *e.SprintID si hoy no tiene sprint, ya es ese mismo sprint,
		// o el que tiene quedó Cerrado (RN-09.2); el destino tiene que estar Abierto (RN-09.3).
		res, err = r.db.ExecContext(ctx, `
			UPDATE backlog_elementos
			SET sprint_id = ?
			WHERE id = ?
			  AND (
			    sprint_id IS NULL
			    OR sprint_id = ?
			    OR NOT EXISTS (SELECT 1 FROM sprints WHERE id = backlog_elementos.sprint_id AND estado = ?)
			  )
			  AND EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = ?)`,
			*e.SprintID, e.ID, *e.SprintID, abierto, *e.SprintID, abierto)
	} else {
		// Solo se puede desasignar si el elemento tiene un sprint y ese sprint sigue Abierto (RN-09.4).
		res, err = r.db.ExecContext(ctx, `
			UPDATE backlog_elementos
			SET sprint_id = NULL
			WHERE id = ?
			  AND sprint_id IS NOT NULL
			  AND EXISTS (SELECT 1 FROM sprints WHERE id = backlog_elementos.sprint_id AND estado = ?)`,
			e.ID, abierto)
	}
	if err != nil {
		return fmt.Errorf("guardar asignación del elemento %d: %w", e.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar guardado de la asignación del elemento %d: %w", e.ID, err)
	}
	if n == 0 {
		return r.errorDeAsignacion(ctx, e)
	}
	return nil
}

// errorDeAsignacion determina, tras un UPDATE que no afectó ninguna fila, cuál regla de RN-09.2 a
// RN-09.4 falló: relee el elemento (puede ya no existir) y, si se estaba asignando, el estado del
// sprint destino. Cubre tanto una validación que ya había fallado como la carrera: otro pedido
// cambió el sprint del elemento, o lo cerró, entre la validación del caso de uso y este guardado.
func (r *RepositorioBacklog) errorDeAsignacion(ctx context.Context, e *domain.ElementoBacklog) error {
	actual, err := r.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		return err // el elemento ya no existe: lo borraron en el medio
	}
	if e.SprintID == nil { // se estaba desasignando (RN-09.4)
		if actual.SprintID == nil {
			return domain.ErrElementoSinSprint
		}
		return domain.ErrSprintCerrado
	}
	// se estaba asignando a *e.SprintID (RN-09.2, RN-09.3)
	_, destinoAbierto, err := r.BuscarSprint(ctx, *e.SprintID)
	if err != nil {
		return err
	}
	if !destinoAbierto {
		return domain.ErrSprintCerrado
	}
	// el destino está Abierto: por eliminación, lo que bloqueó el UPDATE es que el elemento ya
	// estaba en otro sprint Abierto distinto (la otra condición de la cláusula WHERE).
	return domain.ErrElementoYaAsignado
}
