package persistencia

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ListarRondas devuelve todas las rondas del elemento con sus votantes y votos, leídas con
// BuscarRonda (HU-14). Un elemento tiene pocas rondas, así que leerlas de a una no pesa. El orden
// del historial lo decide el dominio (domain.HistorialRondas). Implementa HU-16
// (app.LectorHistorialRondas).
func (r *RepositorioRondas) ListarRondas(ctx context.Context, elementoID int64) ([]domain.Ronda, error) {
	filas, err := r.db.QueryContext(ctx, `SELECT id FROM rondas WHERE elemento_id = ? ORDER BY numero`, elementoID)
	if err != nil {
		return nil, fmt.Errorf("listar rondas del elemento %d: %w", elementoID, err)
	}
	var ids []int64
	for filas.Next() {
		var id int64
		if err := filas.Scan(&id); err != nil {
			filas.Close()
			return nil, fmt.Errorf("leer ronda del elemento %d: %w", elementoID, err)
		}
		ids = append(ids, id)
	}
	filas.Close() // antes de las consultas de BuscarRonda, así no retiene una conexión
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("recorrer rondas del elemento %d: %w", elementoID, err)
	}

	rondas := make([]domain.Ronda, 0, len(ids))
	for _, id := range ids {
		ronda, err := r.BuscarRonda(ctx, id)
		if err != nil {
			return nil, err
		}
		rondas = append(rondas, ronda)
	}
	return rondas, nil
}
