package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioRondas implementa app.RepositorioRondas sobre SQLite. Implementa HU-14.
type RepositorioRondas struct {
	db *sql.DB
}

func NuevoRepositorioRondas(db *sql.DB) *RepositorioRondas {
	return &RepositorioRondas{db: db}
}

// UltimaRonda devuelve la ronda de mayor número del elemento, o nil si no tuvo ninguna.
func (r *RepositorioRondas) UltimaRonda(ctx context.Context, elementoID int64) (*domain.Ronda, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM rondas WHERE elemento_id = ? ORDER BY numero DESC LIMIT 1`, elementoID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("buscar la última ronda del elemento %d: %w", elementoID, err)
	}
	ronda, err := r.BuscarRonda(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ronda, nil
}

// GuardarRonda inserta la ronda y sus votantes en una transacción y le asigna su ID. El índice
// único de rondas abiertas garantiza RN-14.1 aunque dos pedidos abran una ronda a la vez.
func (r *RepositorioRondas) GuardarRonda(ctx context.Context, ronda *domain.Ronda) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO rondas (proyecto_id, elemento_id, numero, estado) VALUES (?, ?, ?, ?)`,
		ronda.ProyectoID, ronda.ElementoID, ronda.Numero, string(ronda.Estado))
	switch {
	// Choca con el índice de rondas abiertas, o con UNIQUE (elemento_id, numero) cuando otro pedido
	// abrió la misma ronda siguiente en el medio: en los dos casos, otra ronda le ganó.
	case esUnique(err):
		return domain.ErrRondaYaAbierta
	case esForeignKey(err):
		return fmt.Errorf("elemento %d: %w", ronda.ElementoID, domain.ErrElementoNoEncontrado)
	case err != nil:
		return fmt.Errorf("insertar ronda del elemento %d: %w", ronda.ElementoID, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID de la ronda: %w", err)
	}
	for i, v := range ronda.Votantes {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO ronda_votantes (ronda_id, orden, nombre, nombre_clave, voto) VALUES (?, ?, ?, ?, ?)`,
			id, i, v.Nombre, domain.ClaveNombre(v.Nombre), v.Voto); err != nil {
			return fmt.Errorf("insertar votante %q: %w", v.Nombre, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar ronda: %w", err)
	}
	ronda.ID = id
	return nil
}

// BuscarRonda devuelve la ronda con sus votantes en orden, o domain.ErrRondaNoEncontrada.
func (r *RepositorioRondas) BuscarRonda(ctx context.Context, id int64) (domain.Ronda, error) {
	ronda := domain.Ronda{ID: id}
	var estado string
	err := r.db.QueryRowContext(ctx,
		`SELECT proyecto_id, elemento_id, numero, estado FROM rondas WHERE id = ?`, id).
		Scan(&ronda.ProyectoID, &ronda.ElementoID, &ronda.Numero, &estado)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Ronda{}, domain.ErrRondaNoEncontrada
	}
	if err != nil {
		return domain.Ronda{}, fmt.Errorf("buscar ronda %d: %w", id, err)
	}
	ronda.Estado = domain.EstadoRonda(estado)

	filas, err := r.db.QueryContext(ctx,
		`SELECT nombre, voto FROM ronda_votantes WHERE ronda_id = ? ORDER BY orden`, id)
	if err != nil {
		return domain.Ronda{}, fmt.Errorf("buscar votantes de la ronda %d: %w", id, err)
	}
	defer filas.Close()
	for filas.Next() {
		var v domain.Votante
		var voto sql.NullInt64
		if err := filas.Scan(&v.Nombre, &voto); err != nil {
			return domain.Ronda{}, fmt.Errorf("leer votante de la ronda %d: %w", id, err)
		}
		if voto.Valid {
			valor := int(voto.Int64)
			v.Voto = &valor
		}
		ronda.Votantes = append(ronda.Votantes, v)
	}
	if err := filas.Err(); err != nil {
		return domain.Ronda{}, fmt.Errorf("leer votantes de la ronda %d: %w", id, err)
	}
	return ronda, nil
}

// GuardarVoto guarda el voto que el votante tiene en ronda y, si ya votaron todos, cierra la ronda,
// en una transacción. El voto solo se guarda si la ronda sigue Abierta en la base: así un voto que
// llega después de que otro pedido la cerró recibe domain.ErrRondaCerrada (RN-14.6). El cierre lo
// decide la base y no la ronda recibida, que pudo leerse antes de un voto simultáneo.
func (r *RepositorioRondas) GuardarVoto(ctx context.Context, ronda domain.Ronda, votante string) error {
	clave := domain.ClaveNombre(votante)
	i := slices.IndexFunc(ronda.Votantes, func(v domain.Votante) bool { return domain.ClaveNombre(v.Nombre) == clave })
	switch {
	case i < 0:
		return fmt.Errorf("guardar voto de %q: %w", votante, domain.ErrNoEsVotante)
	case ronda.Votantes[i].Voto == nil: // el caso de uso siempre aplica el voto antes de guardarlo
		return fmt.Errorf("guardar voto de %q en la ronda %d: la ronda no tiene su voto", votante, ronda.ID)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`UPDATE ronda_votantes SET voto = ?
		 WHERE ronda_id = ? AND nombre_clave = ?
		   AND EXISTS (SELECT 1 FROM rondas WHERE id = ? AND estado = ?)`,
		*ronda.Votantes[i].Voto, ronda.ID, clave, ronda.ID, string(domain.RondaAbierta))
	if err != nil {
		return fmt.Errorf("guardar voto en la ronda %d: %w", ronda.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar voto en la ronda %d: %w", ronda.ID, err)
	}
	if n == 0 {
		return domain.ErrRondaCerrada
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE rondas SET estado = ?
		 WHERE id = ? AND NOT EXISTS (SELECT 1 FROM ronda_votantes WHERE ronda_id = ? AND voto IS NULL)`,
		string(domain.RondaCerrada), ronda.ID, ronda.ID); err != nil {
		return fmt.Errorf("cerrar la ronda %d: %w", ronda.ID, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar voto: %w", err)
	}
	return nil
}
