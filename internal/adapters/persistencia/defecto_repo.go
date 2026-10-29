package persistencia

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioDefectos implementa app.RepositorioDefectos sobre SQLite. Implementa HU-20 (RN-20.6).
type RepositorioDefectos struct {
	db *sql.DB
}

func NuevoRepositorioDefectos(db *sql.DB) *RepositorioDefectos {
	return &RepositorioDefectos{db: db}
}

// GuardarDefecto inserta el defecto y le asigna su ID. Si el elemento ya no existe (otro pedido lo
// borró después de que el caso de uso lo verificara), la clave foránea lo rechaza y se devuelve
// domain.ErrElementoNoEncontrado, igual que en HU-18 (esfuerzo_repo.go).
func (r *RepositorioDefectos) GuardarDefecto(ctx context.Context, d *domain.Defecto) error {
	var sprintID sql.NullInt64
	if d.SprintDeteccionID != nil {
		sprintID = sql.NullInt64{Int64: *d.SprintDeteccionID, Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO defectos (proyecto_id, elemento_id, descripcion, severidad, estado, sprint_deteccion_id)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		d.ProyectoID, d.ElementoID, d.Descripcion, string(d.Severidad), string(d.Estado), sprintID)
	switch {
	case esForeignKey(err):
		return fmt.Errorf("elemento %d: %w", d.ElementoID, domain.ErrElementoNoEncontrado)
	case err != nil:
		return fmt.Errorf("insertar defecto del elemento %d: %w", d.ElementoID, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del defecto: %w", err)
	}
	d.ID = id
	return nil
}
