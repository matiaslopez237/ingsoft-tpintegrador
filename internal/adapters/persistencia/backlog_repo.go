package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

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
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
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

// ListarElementos devuelve los elementos del proyecto con sus criterios en orden. Usa dos
// consultas (elementos y criterios) en lugar de una por elemento. El orden de la lista lo decide
// el dominio (HU-07, RN-07.2): acá solo se ordena por ID para que el resultado sea estable.
// Implementa HU-07 (app.ListadorElementos, RN-07.1).
func (r *RepositorioBacklog) ListarElementos(ctx context.Context, proyectoID int64) ([]domain.ElementoBacklog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, titulo, descripcion, prioridad, estado, story_points
		   FROM backlog_elementos WHERE proyecto_id = ? ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar elementos del proyecto %d: %w", proyectoID, err)
	}
	defer rows.Close()

	var (
		elementos []domain.ElementoBacklog
		posicion  = map[int64]int{} // ID del elemento -> índice en elementos
	)
	for rows.Next() {
		var (
			e                 = domain.ElementoBacklog{ProyectoID: proyectoID}
			prioridad, estado string
			sp                sql.NullInt64
		)
		if err := rows.Scan(&e.ID, &e.Titulo, &e.Descripcion, &prioridad, &estado, &sp); err != nil {
			return nil, fmt.Errorf("leer elemento del proyecto %d: %w", proyectoID, err)
		}
		e.Prioridad, e.Estado = domain.Prioridad(prioridad), domain.EstadoElemento(estado)
		if sp.Valid {
			n := int(sp.Int64)
			e.StoryPoints = &n
		}
		posicion[e.ID] = len(elementos)
		elementos = append(elementos, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer elementos del proyecto %d: %w", proyectoID, err)
	}

	criterios, err := r.db.QueryContext(ctx,
		`SELECT c.elemento_id, c.texto
		   FROM backlog_criterios c JOIN backlog_elementos e ON e.id = c.elemento_id
		  WHERE e.proyecto_id = ? ORDER BY c.elemento_id, c.orden`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar criterios del proyecto %d: %w", proyectoID, err)
	}
	defer criterios.Close()
	for criterios.Next() {
		var (
			elementoID int64
			texto      string
		)
		if err := criterios.Scan(&elementoID, &texto); err != nil {
			return nil, fmt.Errorf("leer criterio del proyecto %d: %w", proyectoID, err)
		}
		if i, ok := posicion[elementoID]; ok {
			elementos[i].Criterios = append(elementos[i].Criterios, texto)
		}
	}
	if err := criterios.Err(); err != nil {
		return nil, fmt.Errorf("recorrer criterios del proyecto %d: %w", proyectoID, err)
	}
	return elementos, nil
}
