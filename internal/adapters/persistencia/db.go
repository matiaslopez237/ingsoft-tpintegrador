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
// Los PRAGMA van en la cadena de conexión para que apliquen a todas las conexiones del pool.
func AbrirDB(ruta string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ruta+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("abrir base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("conectar a la base de datos: %w", err)
	}

	if err := migrar(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrar esquema: %w", err)
	}

	return db, nil
}

// migrar aplica el esquema. Cada historia de usuario agrega sus propias tablas
// a medida que se implementa (ver Tarea 1.6, 2.6, 3.6, etc.
// de docs/scrum/sprint-1/tareas-sprint-1.md).
func migrar(db *sql.DB) error {
	_, err := db.Exec(`
-- HU-01. nombre_clave = domain.ClaveNombre(nombre): garantiza RN-01.3 también en la base.
CREATE TABLE IF NOT EXISTS proyectos (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	nombre       TEXT    NOT NULL,
	nombre_clave TEXT    NOT NULL UNIQUE,
	fecha_inicio TEXT    NOT NULL,
	fecha_fin    TEXT,
	estado       TEXT    NOT NULL
);`)
	return err
}
