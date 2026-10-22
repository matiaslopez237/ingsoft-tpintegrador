// Package persistencia contiene los adaptadores que implementan la persistencia
// de las entidades del dominio. La base elegida es SQLite (ver docs/adr/0002-persistencia.md).
package persistencia

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// AbrirDB abre (y crea si no existe) el archivo de base de datos SQLite indicado,
// y aplica la migración inicial del esquema.
func AbrirDB(ruta string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		return nil, fmt.Errorf("abrir base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("conectar a la base de datos: %w", err)
	}

	if err := migrar(db); err != nil {
		return nil, fmt.Errorf("migrar esquema: %w", err)
	}

	return db, nil
}

// migrar aplica el esquema inicial. Por ahora está vacío: cada historia de usuario
// agrega sus propias tablas a medida que se implementa (ver Tarea 1.6, 2.6, 3.6, etc.
// de docs/scrum/sprint-1/tareas-sprint-1.md).
func migrar(db *sql.DB) error {
	_, err := db.Exec(`PRAGMA foreign_keys = ON;`)
	return err
}
