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

// Actualizar guarda el nombre, las fechas y el estado de un proyecto existente. Solo cambia la
// fila de proyectos, así que no toca sus integrantes (RN-02.10). Si otro proyecto tomó el nombre
// en el medio, devuelve domain.ErrNombreDuplicado; si el proyecto ya no existe,
// domain.ErrProyectoNoEncontrado.
// Implementa HU-02 (app.RepositorioModificacion, RN-02.3, CA-02.10).
func (r *RepositorioProyectos) Actualizar(ctx context.Context, p domain.Proyecto) error {
	var fin sql.NullString
	if p.FechaFin != nil {
		fin = sql.NullString{String: p.FechaFin.Format(time.DateOnly), Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE proyectos SET nombre = ?, nombre_clave = ?, fecha_inicio = ?, fecha_fin = ?, estado = ? WHERE id = ?`,
		p.Nombre, domain.ClaveNombre(p.Nombre), p.FechaInicio.Format(time.DateOnly), fin, string(p.Estado), p.ID)
	if esUnique(err) {
		return domain.ErrNombreDuplicado
	}
	if err != nil {
		return fmt.Errorf("actualizar proyecto %d: %w", p.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar actualización del proyecto %d: %w", p.ID, err)
	}
	if n == 0 {
		return domain.ErrProyectoNoEncontrado
	}
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

	p.Integrantes, err = r.buscarIntegrantes(ctx, p.ID)
	if err != nil {
		return domain.Proyecto{}, err
	}
	return p, nil
}

// buscarIntegrantes devuelve los integrantes del proyecto con sus roles, ambos en el orden en que se agregaron.
func (r *RepositorioProyectos) buscarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Integrante, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT i.id, i.nombre, ir.rol
		   FROM integrantes i JOIN integrante_roles ir ON ir.integrante_id = i.id
		  WHERE i.proyecto_id = ?
		  ORDER BY i.id, ir.id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("buscar integrantes del proyecto %d: %w", proyectoID, err)
	}
	defer rows.Close()

	var (
		integrantes []domain.Integrante
		ultimoID    int64
	)
	for rows.Next() {
		var (
			id          int64
			nombre, rol string
		)
		if err := rows.Scan(&id, &nombre, &rol); err != nil {
			return nil, fmt.Errorf("leer integrante del proyecto %d: %w", proyectoID, err)
		}
		// Las filas vienen agrupadas por integrante: una por cada rol.
		if len(integrantes) == 0 || id != ultimoID {
			integrantes = append(integrantes, domain.Integrante{Nombre: nombre})
			ultimoID = id
		}
		ultimo := &integrantes[len(integrantes)-1]
		ultimo.Roles = append(ultimo.Roles, rol)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer integrantes del proyecto %d: %w", proyectoID, err)
	}
	return integrantes, nil
}

// AgregarIntegrante inserta un integrante y sus roles sin tocar a los demás del proyecto, así dos
// altas simultáneas no se pisan. Si otro pedido agregó a la misma persona en el medio, la
// restricción UNIQUE lo rechaza y se devuelve domain.ErrIntegranteDuplicado.
// Implementa HU-03 (app.RepositorioIntegrantes, RN-03.3, sección 5).
func (r *RepositorioProyectos) AgregarIntegrante(ctx context.Context, proyectoID int64, i domain.Integrante) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar alta del integrante %q: %w", i.Nombre, err)
	}
	defer tx.Rollback() // no hace nada si ya se hizo Commit

	res, err := tx.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, nombre, nombre_clave) VALUES (?, ?, ?)`,
		proyectoID, i.Nombre, domain.ClaveNombre(i.Nombre))
	if esUnique(err) {
		return domain.ErrIntegranteDuplicado
	}
	if err != nil {
		return fmt.Errorf("insertar integrante %q del proyecto %d: %w", i.Nombre, proyectoID, err)
	}
	integranteID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener ID del integrante %q: %w", i.Nombre, err)
	}
	for _, rol := range i.Roles {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO integrante_roles (integrante_id, rol, rol_clave) VALUES (?, ?, ?)`,
			integranteID, rol, domain.ClaveNombre(rol)); err != nil {
			return fmt.Errorf("insertar rol %q del integrante %q: %w", rol, i.Nombre, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar alta del integrante %q: %w", i.Nombre, err)
	}
	return nil
}

// AgregarRol inserta un rol del integrante sin tocar sus otros roles, así dos pedidos simultáneos
// no se pisan. Si otro pedido sumó el mismo rol en el medio, la restricción UNIQUE lo rechaza
// (domain.ErrRolDuplicado); si otro pedido quitó al integrante, no se inserta nada
// (domain.ErrIntegranteNoEncontrado).
// Implementa HU-03 (app.RepositorioIntegrantes, RN-03.5, RN-03.6, sección 5).
func (r *RepositorioProyectos) AgregarRol(ctx context.Context, proyectoID int64, nombre, rol string) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO integrante_roles (integrante_id, rol, rol_clave)
		 SELECT id, ?, ? FROM integrantes WHERE proyecto_id = ? AND nombre_clave = ?`,
		rol, domain.ClaveNombre(rol), proyectoID, domain.ClaveNombre(nombre))
	if esUnique(err) {
		return domain.ErrRolDuplicado
	}
	if err != nil {
		return fmt.Errorf("insertar rol %q del integrante %q: %w", rol, nombre, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar alta del rol %q: %w", rol, err)
	}
	if n == 0 {
		return domain.ErrIntegranteNoEncontrado
	}
	return nil
}

// QuitarIntegrante borra un integrante sin tocar a los demás del proyecto. Si otro pedido ya lo
// quitó, no borra ninguna fila y devuelve domain.ErrIntegranteNoEncontrado.
// Implementa HU-03 (app.RepositorioIntegrantes, CA-03.4, sección 5).
func (r *RepositorioProyectos) QuitarIntegrante(ctx context.Context, proyectoID int64, nombre string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM integrantes WHERE proyecto_id = ? AND nombre_clave = ?`,
		proyectoID, domain.ClaveNombre(nombre))
	if err != nil {
		return fmt.Errorf("borrar integrante %q del proyecto %d: %w", nombre, proyectoID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar borrado del integrante %q: %w", nombre, err)
	}
	if n == 0 {
		return domain.ErrIntegranteNoEncontrado
	}
	return nil
}

// esUnique indica si err es la violación de una restricción UNIQUE de SQLite.
func esUnique(err error) bool {
	var se *sqlite.Error
	return errors.As(err, &se) && se.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}

// Contar devuelve la cantidad de proyectos guardados.
func (r *RepositorioProyectos) Contar(ctx context.Context) (int, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM proyectos`).Scan(&n); err != nil {
		return 0, fmt.Errorf("contar proyectos: %w", err)
	}
	return n, nil
}
