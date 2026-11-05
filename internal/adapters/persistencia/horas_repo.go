package persistencia

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Las dos consultas de horas solo difieren en el WHERE. El LEFT JOIN deja a los elementos sin registros,
// con 0 horas reales (spec HU-23, sección 5), y no se mira la fecha de los registros (RN-23.9).
const (
	seleccionHoras = `
		SELECT e.medias_horas_estimadas, COALESCE(SUM(reg.medias_horas), 0)
		FROM backlog_elementos e
		LEFT JOIN registros_esfuerzo reg ON reg.elemento_id = e.id `
	agrupacionHoras = `
		GROUP BY e.id
		ORDER BY e.id`
)

// HorasDeProyecto devuelve una fila por elemento del proyecto, esté o no en un sprint (RN-23.10), con una
// sola consulta. Implementa HU-23 (app.LectorHoras).
func (r *RepositorioEsfuerzo) HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error) {
	filas, err := r.leerHoras(ctx, seleccionHoras+`WHERE e.proyecto_id = ?`+agrupacionHoras, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("horas del proyecto %d: %w", proyectoID, err)
	}
	return filas, nil
}

// HorasDeSprint devuelve una fila por elemento que tiene hoy ese sprint (RN-23.7), con una sola consulta.
// Implementa HU-23 (app.LectorHoras).
func (r *RepositorioEsfuerzo) HorasDeSprint(ctx context.Context, sprintID int64) ([]domain.HorasElemento, error) {
	filas, err := r.leerHoras(ctx, seleccionHoras+`WHERE e.sprint_id = ?`+agrupacionHoras, sprintID)
	if err != nil {
		return nil, fmt.Errorf("horas del sprint %d: %w", sprintID, err)
	}
	return filas, nil
}

// leerHoras corre la consulta y convierte cada fila: de medias horas a horas, y NULL a sin estimar.
func (r *RepositorioEsfuerzo) leerHoras(ctx context.Context, consulta string, id int64) ([]domain.HorasElemento, error) {
	rows, err := r.db.QueryContext(ctx, consulta, id)
	if err != nil {
		return nil, fmt.Errorf("consultar: %w", err)
	}
	defer rows.Close()

	var filas []domain.HorasElemento
	for rows.Next() {
		var (
			estimadas sql.NullInt64
			reales    int
		)
		if err := rows.Scan(&estimadas, &reales); err != nil {
			return nil, fmt.Errorf("leer fila: %w", err)
		}
		filas = append(filas, domain.HorasElemento{
			HorasEstimadas: desdeMediasHorasEstimadas(estimadas),
			HorasReales:    desdeMediasHoras(reales),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer filas: %w", err)
	}
	return filas, nil
}
