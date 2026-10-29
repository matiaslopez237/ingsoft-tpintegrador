package persistencia

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioSprints implementa app.RepositorioSprints sobre SQLite. Implementa HU-08.
type RepositorioSprints struct {
	db *sql.DB
}

func NuevoRepositorioSprints(db *sql.DB) *RepositorioSprints {
	return &RepositorioSprints{db: db}
}

// GuardarSprint inserta el sprint y le asigna su ID y su número dentro del proyecto. El número
// (RN-08.7) y el control de superposición (RN-08.5) se resuelven en el mismo INSERT, que SQLite
// ejecuta de forma atómica: dos pedidos simultáneos no pueden tomar el mismo número ni crear
// sprints superpuestos. Si el sprint se superpone con otro devuelve domain.ErrSprintSuperpuesto;
// si el proyecto no existe, domain.ErrProyectoNoEncontrado.
// Implementa HU-08 (app.RepositorioSprints, sección 5).
func (r *RepositorioSprints) GuardarSprint(ctx context.Context, s *domain.Sprint) error {
	inicio, fin := s.FechaInicio.Format(time.DateOnly), s.FechaFin.Format(time.DateOnly)
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado)
		 SELECT ?, (SELECT COALESCE(MAX(numero), 0) + 1 FROM sprints WHERE proyecto_id = ?), ?, ?, ?, ?
		 WHERE NOT EXISTS (
		     SELECT 1 FROM sprints WHERE proyecto_id = ? AND fecha_inicio <= ? AND ? <= fecha_fin)`,
		s.ProyectoID, s.ProyectoID, inicio, fin, s.Objetivo, string(s.Estado),
		s.ProyectoID, fin, inicio)
	if esForeignKey(err) {
		return domain.ErrProyectoNoEncontrado
	}
	if err != nil {
		return fmt.Errorf("insertar sprint del proyecto %d: %w", s.ProyectoID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar alta del sprint: %w", err)
	}
	if n == 0 {
		return domain.ErrSprintSuperpuesto
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del sprint: %w", err)
	}
	var numero int
	if err := r.db.QueryRowContext(ctx, `SELECT numero FROM sprints WHERE id = ?`, id).Scan(&numero); err != nil {
		return fmt.Errorf("leer número del sprint %d: %w", id, err)
	}
	s.ID, s.Numero = id, numero
	return nil
}

// ListarSprints devuelve los sprints del proyecto ordenados por número (CA-08.9).
func (r *RepositorioSprints) ListarSprints(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+columnasSprint+` FROM sprints WHERE proyecto_id = ? ORDER BY numero`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	defer rows.Close()

	var sprints []domain.Sprint
	for rows.Next() {
		s, err := escanearSprint(rows)
		if err != nil {
			return nil, fmt.Errorf("leer sprint del proyecto %d: %w", proyectoID, err)
		}
		sprints = append(sprints, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer sprints del proyecto %d: %w", proyectoID, err)
	}
	return sprints, nil
}

// columnasSprint son las columnas que lee escanearSprint, en su orden.
const columnasSprint = `id, proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado`

// filaSprint es lo que tienen en común *sql.Row y *sql.Rows para escanear una fila.
type filaSprint interface {
	Scan(dest ...any) error
}

// escanearSprint arma un domain.Sprint a partir de una fila con columnasSprint. Devuelve el error
// de Scan sin envolver, así quien llama puede reconocer sql.ErrNoRows.
func escanearSprint(f filaSprint) (domain.Sprint, error) {
	var (
		s                   domain.Sprint
		inicio, fin, estado string
	)
	if err := f.Scan(&s.ID, &s.ProyectoID, &s.Numero, &inicio, &fin, &s.Objetivo, &estado); err != nil {
		return domain.Sprint{}, err
	}
	var err error
	if s.FechaInicio, err = time.Parse(time.DateOnly, inicio); err != nil {
		return domain.Sprint{}, fmt.Errorf("fecha de inicio guardada inválida: %w", err)
	}
	if s.FechaFin, err = time.Parse(time.DateOnly, fin); err != nil {
		return domain.Sprint{}, fmt.Errorf("fecha de fin guardada inválida: %w", err)
	}
	s.Estado = domain.EstadoSprint(estado)
	return s, nil
}
