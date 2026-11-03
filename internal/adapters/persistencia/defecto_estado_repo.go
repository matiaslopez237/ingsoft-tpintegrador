package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const columnasDefecto = "id, proyecto_id, elemento_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id"

// filaDefecto es lo que tienen en común *sql.Row y *sql.Rows para escanear una fila.
type filaDefecto interface {
	Scan(dest ...any) error
}

// escanearDefecto arma un domain.Defecto desde una fila con columnasDefecto, en ese orden.
func escanearDefecto(fila filaDefecto) (domain.Defecto, error) {
	var (
		d                 domain.Defecto
		severidad, estado string
		sprintDeteccion   sql.NullInt64
		sprintResolucion  sql.NullInt64
	)
	err := fila.Scan(&d.ID, &d.ProyectoID, &d.ElementoID, &d.Descripcion, &severidad, &estado, &sprintDeteccion, &sprintResolucion)
	if err != nil {
		return domain.Defecto{}, err
	}
	d.Severidad = domain.Severidad(severidad)
	d.Estado = domain.EstadoDefecto(estado)
	if sprintDeteccion.Valid {
		d.SprintDeteccionID = &sprintDeteccion.Int64
	}
	if sprintResolucion.Valid {
		d.SprintResolucionID = &sprintResolucion.Int64
	}
	return d, nil
}

// BuscarDefectoPorID devuelve el defecto con ese ID, o domain.ErrDefectoNoEncontrado si no existe.
// Implementa HU-21 (app.RepositorioEstadoDefecto).
func (r *RepositorioDefectos) BuscarDefectoPorID(ctx context.Context, id int64) (domain.Defecto, error) {
	d, err := escanearDefecto(r.db.QueryRowContext(ctx, `SELECT `+columnasDefecto+` FROM defectos WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Defecto{}, domain.ErrDefectoNoEncontrado
	}
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("buscar defecto %d: %w", id, err)
	}
	return d, nil
}

// ActualizarEstadoDefecto guarda el estado y el SprintResolucionID de d, solo si el defecto sigue
// en anterior. Si otro pedido lo cambió en el medio, no afecta ninguna fila y se relee para
// devolver el error de dominio que corresponda a anterior (ErrDefectoNoAbierto, ErrDefectoNoEnProgreso
// o ErrDefectoNoResuelto), igual que CerrarSprint en HU-11 (spec HU-21, sección 5: pedidos simultáneos).
// Implementa HU-21 (app.RepositorioEstadoDefecto).
func (r *RepositorioDefectos) ActualizarEstadoDefecto(ctx context.Context, d *domain.Defecto, anterior domain.EstadoDefecto) error {
	var sprintResolucion sql.NullInt64
	if d.SprintResolucionID != nil {
		sprintResolucion = sql.NullInt64{Int64: *d.SprintResolucionID, Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE defectos SET estado = ?, sprint_resolucion_id = ? WHERE id = ? AND estado = ?`,
		string(d.Estado), sprintResolucion, d.ID, string(anterior))
	if err != nil {
		return fmt.Errorf("actualizar estado del defecto %d: %w", d.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar actualización del defecto %d: %w", d.ID, err)
	}
	if n == 0 {
		return r.errorDeActualizacion(ctx, d.ID, anterior)
	}
	return nil
}

// errorDeActualizacion determina, tras un UPDATE que no afectó ninguna fila, por qué: el defecto ya
// no está en el estado esperado por quien llamó (otro pedido lo cambió en el medio).
func (r *RepositorioDefectos) errorDeActualizacion(ctx context.Context, id int64, anterior domain.EstadoDefecto) error {
	actual, err := r.BuscarDefectoPorID(ctx, id)
	if err != nil {
		return err
	}
	switch anterior {
	case domain.DefectoAbierto:
		return domain.ErrDefectoNoAbierto
	case domain.DefectoEnProgreso:
		return domain.ErrDefectoNoEnProgreso
	case domain.DefectoResuelto:
		return domain.ErrDefectoNoResuelto
	default:
		return fmt.Errorf("actualizar estado del defecto %d: estado guardado inesperado %q", id, actual.Estado)
	}
}

// ListarDefectos devuelve los defectos del proyecto ordenados por id (orden estable, no depende de
// cómo los devuelva SQLite); lista vacía si no tiene ninguno. Implementa HU-21 (app.LectorDefectos).
func (r *RepositorioDefectos) ListarDefectos(ctx context.Context, proyectoID int64) ([]domain.Defecto, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+columnasDefecto+` FROM defectos WHERE proyecto_id = ? ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}
	defer rows.Close()

	var defectos []domain.Defecto
	for rows.Next() {
		d, err := escanearDefecto(rows)
		if err != nil {
			return nil, fmt.Errorf("leer defecto del proyecto %d: %w", proyectoID, err)
		}
		defectos = append(defectos, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer defectos del proyecto %d: %w", proyectoID, err)
	}
	return defectos, nil
}
