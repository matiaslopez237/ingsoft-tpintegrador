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
	if err := migrarSprintID(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrar columna sprint_id: %w", err)
	}
	if err := migrarPuntosDeSprint(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrar columnas de puntos de sprints: %w", err)
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
);

-- HU-03. nombre_clave = domain.ClaveNombre(nombre): garantiza RN-03.3 también en la base.
CREATE TABLE IF NOT EXISTS integrantes (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
	nombre       TEXT    NOT NULL,
	nombre_clave TEXT    NOT NULL,
	UNIQUE (proyecto_id, nombre_clave)
);

-- HU-03. rol_clave = domain.ClaveNombre(rol): garantiza RN-03.6 también en la base.
-- Al quitar un integrante se borran sus roles (ON DELETE CASCADE).
CREATE TABLE IF NOT EXISTS integrante_roles (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	integrante_id INTEGER NOT NULL REFERENCES integrantes (id) ON DELETE CASCADE,
	rol           TEXT    NOT NULL,
	rol_clave     TEXT    NOT NULL,
	UNIQUE (integrante_id, rol_clave)
);

-- HU-05. titulo_clave = domain.ClaveNombre(titulo): garantiza RN-05.3 también en la base.
-- story_points NULL: elemento sin estimar (RN-05.6).
CREATE TABLE IF NOT EXISTS backlog_elementos (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
	titulo       TEXT    NOT NULL,
	titulo_clave TEXT    NOT NULL,
	descripcion  TEXT    NOT NULL,
	prioridad    TEXT    NOT NULL,
	estado       TEXT    NOT NULL,
	story_points INTEGER,
	UNIQUE (proyecto_id, titulo_clave)
);

-- HU-05. Criterios de aceptación de cada elemento; orden conserva el orden en que se enviaron.
CREATE TABLE IF NOT EXISTS backlog_criterios (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	elemento_id INTEGER NOT NULL REFERENCES backlog_elementos (id) ON DELETE CASCADE,
	orden       INTEGER NOT NULL,
	texto       TEXT    NOT NULL,
	UNIQUE (elemento_id, orden)
);

-- HU-18. integrante guarda el nombre y no referencia a integrantes: quitar a un integrante no
-- borra ni bloquea sus registros (RN-18.3). medias_horas = horas * 2, entero para no redondear.
CREATE TABLE IF NOT EXISTS registros_esfuerzo (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
	elemento_id  INTEGER NOT NULL REFERENCES backlog_elementos (id),
	integrante   TEXT    NOT NULL,
	fecha        TEXT    NOT NULL,
	actividad    TEXT    NOT NULL,
	medias_horas INTEGER NOT NULL
);

-- HU-08. numero: 1, 2, 3… dentro de cada proyecto (RN-08.7); lo calcula el INSERT de GuardarSprint.
-- Fechas AAAA-MM-DD: compararlas como texto respeta el orden de las fechas (RN-08.5).
-- HU-11. puntos_*: el resumen que registra el cierre (RN-11.6); 0 mientras el sprint está Abierto.
CREATE TABLE IF NOT EXISTS sprints (
	id                  INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id         INTEGER NOT NULL REFERENCES proyectos (id),
	numero              INTEGER NOT NULL,
	fecha_inicio        TEXT    NOT NULL,
	fecha_fin           TEXT    NOT NULL,
	objetivo            TEXT    NOT NULL,
	estado              TEXT    NOT NULL,
	puntos_planificados INTEGER NOT NULL DEFAULT 0,
	puntos_completados  INTEGER NOT NULL DEFAULT 0,
	UNIQUE (proyecto_id, numero)
);`)
	return err
}

// migrarSprintID agrega la columna sprint_id a backlog_elementos si todavía no existe (HU-09:
// NULL = sin asignar a un sprint). CREATE TABLE IF NOT EXISTS no modifica una tabla que ya existe,
// así que una base creada antes de HU-09 se queda sin la columna; esto la agrega sin tocar los
// datos. Es seguro llamarla en cada arranque: si la columna ya está, no hace nada.
func migrarSprintID(db *sql.DB) error {
	return agregarColumnas(db, "backlog_elementos", map[string]string{
		"sprint_id": "INTEGER REFERENCES sprints (id)",
	})
}

// migrarPuntosDeSprint agrega a sprints las columnas del resumen de cierre de HU-11 (RN-11.6) si
// todavía no existen, por el mismo motivo que migrarSprintID. Los sprints que ya estaban quedan con
// 0 puntos. Es seguro llamarla en cada arranque.
func migrarPuntosDeSprint(db *sql.DB) error {
	return agregarColumnas(db, "sprints", map[string]string{
		"puntos_planificados": "INTEGER NOT NULL DEFAULT 0",
		"puntos_completados":  "INTEGER NOT NULL DEFAULT 0",
	})
}

// agregarColumnas agrega a la tabla las columnas que le falten, con su definición, sin tocar los
// datos; las que ya están se dejan como están. tabla y las columnas son siempre constantes del
// código, nunca datos externos: por eso es seguro concatenarlas (PRAGMA y ALTER TABLE no admiten "?").
func agregarColumnas(db *sql.DB, tabla string, columnas map[string]string) error {
	rows, err := db.Query(`PRAGMA table_info(` + tabla + `)`)
	if err != nil {
		return fmt.Errorf("consultar columnas de %s: %w", tabla, err)
	}
	existentes := map[string]bool{}
	for rows.Next() {
		var (
			cid, notNull, pk int
			nombre, tipo     string
			valorPorOmision  sql.NullString
		)
		if err := rows.Scan(&cid, &nombre, &tipo, &notNull, &valorPorOmision, &pk); err != nil {
			rows.Close()
			return fmt.Errorf("leer columna de %s: %w", tabla, err)
		}
		existentes[nombre] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("recorrer columnas de %s: %w", tabla, err)
	}

	for columna, definicion := range columnas {
		if existentes[columna] {
			continue
		}
		if _, err := db.Exec(`ALTER TABLE ` + tabla + ` ADD COLUMN ` + columna + ` ` + definicion); err != nil {
			return fmt.Errorf("agregar columna %s a %s: %w", columna, tabla, err)
		}
	}
	return nil
}
