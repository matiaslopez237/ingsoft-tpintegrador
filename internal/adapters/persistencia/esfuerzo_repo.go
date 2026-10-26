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

// GuardarRegistro inserta el registro y le asigna su ID. Si el elemento ya no existe (otro pedido
// lo borró después de que el caso de uso lo verificara), la clave foránea lo rechaza y se devuelve
// domain.ErrElementoNoEncontrado. El proyecto no puede faltar si el elemento existe, porque el
// elemento también tiene clave foránea hacia su proyecto.
func (r *RepositorioEsfuerzo) GuardarRegistro(ctx context.Context, reg *domain.RegistroEsfuerzo) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO registros_esfuerzo (proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		reg.ProyectoID, reg.ElementoID, reg.Integrante, reg.Fecha.Format(time.DateOnly), reg.Actividad,
		aMediasHoras(reg.Horas))
	switch {
	case esForeignKey(err):
		return fmt.Errorf("elemento %d: %w", reg.ElementoID, domain.ErrElementoNoEncontrado)
	case err != nil:
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
	reg.Horas = desdeMediasHoras(medias)
	return reg, nil
}

// aMediasHoras pasa las horas a un entero de medias horas (1,5 h → 3). El dominio ya garantiza
// que son múltiplo de 0,5 (RN-18.7); Round solo absorbe la representación binaria del float64.
func aMediasHoras(horas float64) int {
	return int(math.Round(horas * 2))
}

func desdeMediasHoras(medias int) float64 {
	return float64(medias) / 2
}
