package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// RepositorioProyectos implementa app.RepositorioProyectos sobre SQLite.
// Implementa HU-01 (RN-01.3, RN-01.8).
type RepositorioProyectos struct {
	db *sql.DB
}

func NuevoRepositorioProyectos(db *sql.DB) *RepositorioProyectos {
	return &RepositorioProyectos{db: db}
}

func (r *RepositorioProyectos) ExisteNombre(ctx context.Context, clave string) (bool, error) {
	var existe bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM proyectos WHERE nombre_clave = ?)`, clave).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("consultar nombre de proyecto: %w", err)
	}
	return existe, nil
}

// Guardar inserta el proyecto y le asigna el ID autoincremental generado por SQLite.
// Si otro proyecto con el mismo nombre se guardó en el medio, devuelve domain.ErrNombreDuplicado.
func (r *RepositorioProyectos) Guardar(ctx context.Context, p *domain.Proyecto) error {
	var fin sql.NullString
	if p.FechaFin != nil {
		fin = sql.NullString{String: p.FechaFin.Format(time.DateOnly), Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO proyectos (nombre, nombre_clave, fecha_inicio, fecha_fin, estado) VALUES (?, ?, ?, ?, ?)`,
		p.Nombre, domain.ClaveNombre(p.Nombre), p.FechaInicio.Format(time.DateOnly), fin, string(p.Estado))
	if err != nil {
		// La única restricción UNIQUE de la tabla es nombre_clave.
		var se *sqlite.Error
		if errors.As(err, &se) && se.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return domain.ErrNombreDuplicado
		}
		return fmt.Errorf("insertar proyecto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del proyecto: %w", err)
	}
	p.ID = id
	return nil
}

func (r *RepositorioProyectos) BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error) {
	var (
		p           domain.Proyecto
		inicio, est string
		fin         sql.NullString
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, fecha_inicio, fecha_fin, estado FROM proyectos WHERE id = ?`, id).
		Scan(&p.ID, &p.Nombre, &inicio, &fin, &est)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	if err != nil {
		return domain.Proyecto{}, fmt.Errorf("buscar proyecto %d: %w", id, err)
	}

	if p.FechaInicio, err = time.Parse(time.DateOnly, inicio); err != nil {
		return domain.Proyecto{}, fmt.Errorf("fecha de inicio guardada inválida: %w", err)
	}
	if fin.Valid {
		f, err := time.Parse(time.DateOnly, fin.String)
		if err != nil {
			return domain.Proyecto{}, fmt.Errorf("fecha de fin guardada inválida: %w", err)
		}
		p.FechaFin = &f
	}
	p.Estado = domain.Estado(est)
	return p, nil
}

// Contar devuelve la cantidad de proyectos guardados.
func (r *RepositorioProyectos) Contar(ctx context.Context) (int, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM proyectos`).Scan(&n); err != nil {
		return 0, fmt.Errorf("contar proyectos: %w", err)
	}
	return n, nil
}
