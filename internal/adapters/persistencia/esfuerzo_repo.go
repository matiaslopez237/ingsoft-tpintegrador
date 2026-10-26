package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioEsfuerzo implementa app.RepositorioEsfuerzo sobre SQLite.
// Implementa HU-18 (RN-18.3, RN-18.10).
type RepositorioEsfuerzo struct {
	db *sql.DB
}

func NuevoRepositorioEsfuerzo(db *sql.DB) *RepositorioEsfuerzo {
	return &RepositorioEsfuerzo{db: db}
}

// GuardarRegistro inserta el registro y le asigna su ID.
func (r *RepositorioEsfuerzo) GuardarRegistro(ctx context.Context, reg *domain.RegistroEsfuerzo) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO registros_esfuerzo (proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		reg.ProyectoID, reg.ElementoID, reg.Integrante, reg.Fecha.Format(time.DateOnly), reg.Actividad,
		int(math.Round(reg.Horas*2)))
	if err != nil {
		return fmt.Errorf("insertar registro de esfuerzo de %q: %w", reg.Integrante, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del registro de esfuerzo: %w", err)
	}
	reg.ID = id
	return nil
}

// BuscarRegistroPorID devuelve el registro con la fecha como día de la zona del servidor (RN-18.4).
func (r *RepositorioEsfuerzo) BuscarRegistroPorID(ctx context.Context, id int64) (domain.RegistroEsfuerzo, error) {
	var (
		reg    domain.RegistroEsfuerzo
		fecha  string
		medias int
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas
		   FROM registros_esfuerzo WHERE id = ?`, id).
		Scan(&reg.ID, &reg.ProyectoID, &reg.ElementoID, &reg.Integrante, &fecha, &reg.Actividad, &medias)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RegistroEsfuerzo{}, domain.ErrRegistroEsfuerzoNoEncontrado
	}
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("buscar registro de esfuerzo %d: %w", id, err)
	}
	if reg.Fecha, err = time.ParseInLocation(time.DateOnly, fecha, time.Local); err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("leer fecha del registro de esfuerzo %d: %w", id, err)
	}
	reg.Horas = float64(medias) / 2
	return reg, nil
}

// ContarRegistros devuelve la cantidad de registros del elemento (lo usan los tests y el BDD).
func (r *RepositorioEsfuerzo) ContarRegistros(ctx context.Context, elementoID int64) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM registros_esfuerzo WHERE elemento_id = ?`, elementoID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("contar registros de esfuerzo del elemento %d: %w", elementoID, err)
	}
	return n, nil
}
