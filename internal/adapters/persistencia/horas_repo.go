package persistencia

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// HorasDeProyecto devuelve una fila por elemento del proyecto, esté o no en un sprint (RN-23.10), con una
// sola consulta. El LEFT JOIN deja a los elementos sin registros, con 0 horas reales; no se mira la fecha
// de los registros (RN-23.9). Implementa HU-23 (app.LectorHoras).
func (r *RepositorioEsfuerzo) HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.medias_horas_estimadas, COALESCE(SUM(reg.medias_horas), 0)
		FROM backlog_elementos e
		LEFT JOIN registros_esfuerzo reg ON reg.elemento_id = e.id
		WHERE e.proyecto_id = ?
		GROUP BY e.id
		ORDER BY e.id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("consultar horas del proyecto %d: %w", proyectoID, err)
	}
	defer rows.Close()

	var filas []domain.HorasElemento
	for rows.Next() {
		var (
			estimadas sql.NullInt64
			reales    int
		)
		if err := rows.Scan(&estimadas, &reales); err != nil {
			return nil, fmt.Errorf("leer horas del proyecto %d: %w", proyectoID, err)
		}
		filas = append(filas, domain.HorasElemento{
			HorasEstimadas: desdeMediasHorasEstimadas(estimadas),
			HorasReales:    desdeMediasHoras(reales),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer horas del proyecto %d: %w", proyectoID, err)
	}
	return filas, nil
}

// HorasDeSprint devuelve una fila por elemento que tiene hoy ese sprint (RN-23.7), con una sola consulta,
// igual que HorasDeProyecto. Implementa HU-23 (app.LectorHoras).
func (r *RepositorioEsfuerzo) HorasDeSprint(ctx context.Context, sprintID int64) ([]domain.HorasElemento, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.medias_horas_estimadas, COALESCE(SUM(reg.medias_horas), 0)
		FROM backlog_elementos e
		LEFT JOIN registros_esfuerzo reg ON reg.elemento_id = e.id
		WHERE e.sprint_id = ?
		GROUP BY e.id
		ORDER BY e.id`, sprintID)
	if err != nil {
		return nil, fmt.Errorf("consultar horas del sprint %d: %w", sprintID, err)
	}
	defer rows.Close()

	var filas []domain.HorasElemento
	for rows.Next() {
		var (
			estimadas sql.NullInt64
			reales    int
		)
		if err := rows.Scan(&estimadas, &reales); err != nil {
			return nil, fmt.Errorf("leer horas del sprint %d: %w", sprintID, err)
		}
		filas = append(filas, domain.HorasElemento{
			HorasEstimadas: desdeMediasHorasEstimadas(estimadas),
			HorasReales:    desdeMediasHoras(reales),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer horas del sprint %d: %w", sprintID, err)
	}
	return filas, nil
}
