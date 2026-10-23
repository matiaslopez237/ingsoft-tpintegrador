package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ErrElementoNoEncontrado se devuelve cuando no existe un elemento del backlog con el ID buscado.
var ErrElementoNoEncontrado = errors.New("elemento de backlog no encontrado")

// RepositorioBacklog implementa app.RepositorioBacklog sobre SQLite.
// Implementa HU-05 (RN-05.1, RN-05.3, RN-05.9).
type RepositorioBacklog struct {
	db *sql.DB
}

func NuevoRepositorioBacklog(db *sql.DB) *RepositorioBacklog {
	return &RepositorioBacklog{db: db}
}

func (r *RepositorioBacklog) ExisteProyecto(ctx context.Context, id int64) (bool, error) {
	var existe bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM proyectos WHERE id = ?)`, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("consultar proyecto %d: %w", id, err)
	}
	return existe, nil
}

func (r *RepositorioBacklog) ExisteTitulo(ctx context.Context, proyectoID int64, clave string) (bool, error) {
	var existe bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM backlog_elementos WHERE proyecto_id = ? AND titulo_clave = ?)`,
		proyectoID, clave).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("consultar título del backlog: %w", err)
	}
	return existe, nil
}

// GuardarElemento inserta el elemento y sus criterios en una transacción y le asigna su ID.
// Si otro pedido guardó el mismo título en el medio, la restricción UNIQUE lo rechaza y se
// devuelve domain.ErrTituloDuplicado; si el proyecto ya no existe, domain.ErrProyectoNoEncontrado.
func (r *RepositorioBacklog) GuardarElemento(ctx context.Context, e *domain.ElementoBacklog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar alta del elemento %q: %w", e.Titulo, err)
	}
	defer tx.Rollback() // no hace nada si ya se hizo Commit

	var sp sql.NullInt64
	if e.StoryPoints != nil {
		sp = sql.NullInt64{Int64: int64(*e.StoryPoints), Valid: true}
	}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO backlog_elementos (proyecto_id, titulo, titulo_clave, descripcion, prioridad, estado, story_points)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ProyectoID, e.Titulo, domain.ClaveNombre(e.Titulo), e.Descripcion, string(e.Prioridad), string(e.Estado), sp)
	switch {
	case esUnique(err): // la única restricción UNIQUE de la tabla es (proyecto_id, titulo_clave)
		return domain.ErrTituloDuplicado
	case esForeignKey(err):
		return domain.ErrProyectoNoEncontrado
	case err != nil:
		return fmt.Errorf("insertar elemento %q: %w", e.Titulo, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del elemento %q: %w", e.Titulo, err)
	}

	for i, c := range e.Criterios {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO backlog_criterios (elemento_id, orden, texto) VALUES (?, ?, ?)`, id, i, c); err != nil {
			return fmt.Errorf("insertar criterio %d del elemento %q: %w", i+1, e.Titulo, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar alta del elemento %q: %w", e.Titulo, err)
	}
	e.ID = id
	return nil
}

// BuscarElementoPorID devuelve el elemento con sus criterios en el orden en que se guardaron.
// No es parte del puerto de HU-05: lo usan los tests y el BDD, y lo van a usar HU-06 y HU-07.
func (r *RepositorioBacklog) BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error) {
	var (
		e                 domain.ElementoBacklog
		prioridad, estado string
		sp                sql.NullInt64
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, titulo, descripcion, prioridad, estado, story_points
		   FROM backlog_elementos WHERE id = ?`, id).
		Scan(&e.ID, &e.ProyectoID, &e.Titulo, &e.Descripcion, &prioridad, &estado, &sp)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ElementoBacklog{}, ErrElementoNoEncontrado
	}
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("buscar elemento %d: %w", id, err)
	}
	e.Prioridad, e.Estado = domain.Prioridad(prioridad), domain.EstadoElemento(estado)
	if sp.Valid {
		n := int(sp.Int64)
		e.StoryPoints = &n
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT texto FROM backlog_criterios WHERE elemento_id = ? ORDER BY orden`, id)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("buscar criterios del elemento %d: %w", id, err)
	}
	defer rows.Close()
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return domain.ElementoBacklog{}, fmt.Errorf("leer criterio del elemento %d: %w", id, err)
		}
		e.Criterios = append(e.Criterios, c)
	}
	if err := rows.Err(); err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("recorrer criterios del elemento %d: %w", id, err)
	}
	return e, nil
}

// ContarElementos devuelve la cantidad de elementos del backlog del proyecto (lo usan los tests y el BDD).
func (r *RepositorioBacklog) ContarElementos(ctx context.Context, proyectoID int64) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM backlog_elementos WHERE proyecto_id = ?`, proyectoID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("contar elementos del proyecto %d: %w", proyectoID, err)
	}
	return n, nil
}
