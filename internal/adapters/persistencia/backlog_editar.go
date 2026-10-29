package persistencia

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ActualizarElemento guarda de e solo los campos que c indica como enviados, así dos ediciones
// simultáneas de campos distintos no se pisan. Si se envían criterios, reemplaza la lista en la
// misma transacción. Si otro elemento del proyecto tomó el título, devuelve
// domain.ErrTituloDuplicado; si el elemento ya no existe, domain.ErrElementoNoEncontrado.
// Implementa HU-06 (app.RepositorioEdicionBacklog, RN-06.1, sección 5).
func (r *RepositorioBacklog) ActualizarElemento(ctx context.Context, e domain.ElementoBacklog, c domain.CambiosElemento) error {
	var (
		sets []string
		args []any
	)
	if c.Titulo != nil {
		sets = append(sets, "titulo = ?", "titulo_clave = ?")
		args = append(args, e.Titulo, domain.ClaveNombre(e.Titulo))
	}
	if c.Descripcion != nil {
		sets = append(sets, "descripcion = ?")
		args = append(args, e.Descripcion)
	}
	if c.Prioridad != nil {
		sets = append(sets, "prioridad = ?")
		args = append(args, string(e.Prioridad))
	}
	if c.CambiarStoryPoints {
		var sp sql.NullInt64
		if e.StoryPoints != nil {
			sp = sql.NullInt64{Int64: int64(*e.StoryPoints), Valid: true}
		}
		sets = append(sets, "story_points = ?")
		args = append(args, sp)
	}
	if c.Estado != nil {
		sets = append(sets, "estado = ?")
		args = append(args, string(e.Estado))
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar edición del elemento %d: %w", e.ID, err)
	}
	defer tx.Rollback() // no hace nada si ya se hizo Commit

	// Si solo cambian los criterios, el UPDATE sin columnas igual sirve para saber si el elemento existe.
	set := "id = id"
	if len(sets) > 0 {
		set = strings.Join(sets, ", ")
	}
	res, err := tx.ExecContext(ctx, "UPDATE backlog_elementos SET "+set+" WHERE id = ?", append(args, e.ID)...)
	if esUnique(err) { // la única restricción UNIQUE de la tabla es (proyecto_id, titulo_clave)
		return domain.ErrTituloDuplicado
	}
	if err != nil {
		return fmt.Errorf("actualizar elemento %d: %w", e.ID, err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("verificar edición del elemento %d: %w", e.ID, err)
	} else if n == 0 {
		return domain.ErrElementoNoEncontrado
	}

	if c.Criterios != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM backlog_criterios WHERE elemento_id = ?`, e.ID); err != nil {
			return fmt.Errorf("borrar criterios del elemento %d: %w", e.ID, err)
		}
		for i, criterio := range e.Criterios {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO backlog_criterios (elemento_id, orden, texto) VALUES (?, ?, ?)`, e.ID, i, criterio); err != nil {
				return fmt.Errorf("insertar criterio %d del elemento %d: %w", i+1, e.ID, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar edición del elemento %d: %w", e.ID, err)
	}
	return nil
}

// TieneEsfuerzo indica si el elemento tiene registros de esfuerzo (HU-18).
// Implementa HU-06 (app.RepositorioEliminacionBacklog, RN-06.7).
func (r *RepositorioBacklog) TieneEsfuerzo(ctx context.Context, elementoID int64) (bool, error) {
	var tiene bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM registros_esfuerzo WHERE elemento_id = ?)`, elementoID).Scan(&tiene)
	if err != nil {
		return false, fmt.Errorf("consultar esfuerzo del elemento %d: %w", elementoID, err)
	}
	return tiene, nil
}

// TieneDefectos indica si el elemento tiene defectos registrados (HU-20).
// Implementa HU-06 (app.RepositorioEliminacionBacklog, RN-06.10, DEF-02).
func (r *RepositorioBacklog) TieneDefectos(ctx context.Context, elementoID int64) (bool, error) {
	var tiene bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM defectos WHERE elemento_id = ?)`, elementoID).Scan(&tiene)
	if err != nil {
		return false, fmt.Errorf("consultar defectos del elemento %d: %w", elementoID, err)
	}
	return tiene, nil
}

// EliminarElemento borra el elemento; sus criterios se borran con él (ON DELETE CASCADE). Si
// tiene esfuerzo o defectos registrados, las claves foráneas de registros_esfuerzo y defectos lo
// impiden y errorDeEliminacion informa cuál; si no existe, domain.ErrElementoNoEncontrado.
// Implementa HU-06 (app.RepositorioEliminacionBacklog, RN-06.7, RN-06.10).
func (r *RepositorioBacklog) EliminarElemento(ctx context.Context, elementoID int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM backlog_elementos WHERE id = ?`, elementoID)
	if esForeignKey(err) {
		return r.errorDeEliminacion(ctx, elementoID)
	}
	if err != nil {
		return fmt.Errorf("eliminar elemento %d: %w", elementoID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar eliminación del elemento %d: %w", elementoID, err)
	}
	if n == 0 {
		return domain.ErrElementoNoEncontrado
	}
	return nil
}

// errorDeEliminacion determina, tras un DELETE rechazado por una clave foránea, qué lo bloqueó: el
// esfuerzo (se informa primero, RN-06.7) o los defectos (RN-06.10). Antes de DEF-02 se informaba
// siempre el esfuerzo.
func (r *RepositorioBacklog) errorDeEliminacion(ctx context.Context, elementoID int64) error {
	conEsfuerzo, err := r.TieneEsfuerzo(ctx, elementoID)
	if err != nil {
		return err
	}
	if conEsfuerzo {
		return domain.ErrElementoConEsfuerzo
	}
	return domain.ErrElementoConDefectos
}
