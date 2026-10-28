package persistencia

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// GuardarCompletado persiste el estado de un elemento ya existente, sin tocar sus otros campos.
// RN-10.2 se verifica en el propio UPDATE, que SQLite ejecuta de forma atómica: el elemento tiene
// que seguir en el sprint que se leyó (e.SprintID) y ese sprint tiene que seguir Abierto, así un
// pedido simultáneo que lo quitó del sprint (HU-09) o lo cerró (HU-11) no se pisa. Si el elemento
// ya estaba Hecho, SQLite cuenta la fila como afectada aunque el valor no cambie (RN-10.3). Si el
// UPDATE no afecta ninguna fila, errorDeCompletado relee el elemento para informar qué falló.
// Implementa HU-10 (app.RepositorioCompletar, RN-10.2, RN-10.3, spec sección 5).
func (r *RepositorioBacklog) GuardarCompletado(ctx context.Context, e *domain.ElementoBacklog) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE backlog_elementos
		SET estado = ?
		WHERE id = ?
		  AND sprint_id = ?
		  AND EXISTS (SELECT 1 FROM sprints WHERE id = backlog_elementos.sprint_id AND estado = ?)`,
		string(e.Estado), e.ID, e.SprintID, string(domain.SprintAbierto))
	if err != nil {
		return fmt.Errorf("guardar completado del elemento %d: %w", e.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar guardado del completado del elemento %d: %w", e.ID, err)
	}
	if n == 0 {
		return r.errorDeCompletado(ctx, e)
	}
	return nil
}

// errorDeCompletado determina, tras un UPDATE que no afectó ninguna fila, qué parte de RN-10.2
// falló: relee el elemento (puede ya no existir) y compara su sprint con el que se leyó antes.
// Si ya no está en ese sprint (lo quitaron, o lo movieron a otro), es ErrElementoSinSprint; si
// sigue en él, lo único que pudo bloquear el UPDATE es que el sprint ya no esté Abierto.
func (r *RepositorioBacklog) errorDeCompletado(ctx context.Context, e *domain.ElementoBacklog) error {
	actual, err := r.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		return err // el elemento ya no existe: lo borraron en el medio
	}
	if actual.SprintID == nil || e.SprintID == nil || *actual.SprintID != *e.SprintID {
		return domain.ErrElementoSinSprint
	}
	return domain.ErrSprintCerrado
}
