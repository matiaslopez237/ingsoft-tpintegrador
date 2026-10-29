package persistencia

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// Traza: HU-09 (comentario de revisión: migrar no agregaba sprint_id a una base ya existente,
// porque CREATE TABLE IF NOT EXISTS no toca una tabla que ya está creada).
func TestAbrirDB_BaseExistenteSinSprintID_AgregaLaColumna(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "vieja.db")
	crearEsquemaSinSprintID(t, ruta)

	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir base vieja con la migración actual: %v", err)
	}
	defer db.Close()

	if !tieneColumna(t, db, "backlog_elementos", "sprint_id") {
		t.Error("se esperaba que backlog_elementos tuviera la columna sprint_id después de migrar")
	}
}

// Traza: HU-11 (spec, sección 5): una base con la tabla sprints de HU-08 recibe las columnas del
// resumen de cierre, sin perder los sprints que ya tenía (quedan en 0 puntos).
func TestAbrirDB_BaseExistenteSinPuntosDeSprint_AgregaLasColumnas(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "vieja.db")
	crearEsquemaSinSprintID(t, ruta)
	vieja, err := sql.Open("sqlite", ruta)
	if err != nil {
		t.Fatalf("abrir base vieja: %v", err)
	}
	_, err = vieja.Exec(`
CREATE TABLE sprints (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
	numero       INTEGER NOT NULL,
	fecha_inicio TEXT    NOT NULL,
	fecha_fin    TEXT    NOT NULL,
	objetivo     TEXT    NOT NULL,
	estado       TEXT    NOT NULL,
	UNIQUE (proyecto_id, numero)
);
INSERT INTO proyectos (nombre, nombre_clave, fecha_inicio, estado) VALUES ('Alfa', 'alfa', '2026-10-01', 'Activo');
INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado)
VALUES (1, 1, '2026-10-12', '2026-10-25', 'Viejo', 'Cerrado');`)
	vieja.Close()
	if err != nil {
		t.Fatalf("crear sprints viejos: %v", err)
	}

	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir base vieja con la migración actual: %v", err)
	}
	defer db.Close()

	for _, col := range []string{"puntos_planificados", "puntos_completados"} {
		if !tieneColumna(t, db, "sprints", col) {
			t.Errorf("se esperaba que sprints tuviera la columna %s después de migrar", col)
		}
	}
	s, err := NuevoRepositorioSprints(db).BuscarSprintPorID(context.Background(), 1)
	if err != nil {
		t.Fatalf("leer el sprint viejo: %v", err)
	}
	if s.Objetivo != "Viejo" || s.PuntosPlanificados != 0 || s.PuntosCompletados != 0 {
		t.Errorf("sprint viejo = %+v, se esperaba el mismo con 0 puntos", s)
	}
}

// Traza: HU-09 (la migración tiene que poder correr dos veces sin fallar, por ejemplo si el
// servidor arranca de nuevo sobre una base ya migrada).
func TestAbrirDB_DosVeces_NoFalla(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "nueva.db")

	for i := 0; i < 2; i++ {
		db, err := AbrirDB(ruta)
		if err != nil {
			t.Fatalf("abrir base (vuelta %d): %v", i+1, err)
		}
		db.Close()
	}
}

// crearEsquemaSinSprintID simula una base creada antes de HU-09, con las tablas que ya existían
// (abreviadas: solo proyectos, sprints y backlog_elementos sin sprint_id, que es lo que migrar
// necesita para decidir si falta la columna).
func crearEsquemaSinSprintID(t *testing.T, ruta string) {
	t.Helper()
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		t.Fatalf("abrir base vieja: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
CREATE TABLE proyectos (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	nombre       TEXT    NOT NULL,
	nombre_clave TEXT    NOT NULL UNIQUE,
	fecha_inicio TEXT    NOT NULL,
	fecha_fin    TEXT,
	estado       TEXT    NOT NULL
);
CREATE TABLE backlog_elementos (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
	titulo       TEXT    NOT NULL,
	titulo_clave TEXT    NOT NULL,
	descripcion  TEXT    NOT NULL,
	prioridad    TEXT    NOT NULL,
	estado       TEXT    NOT NULL,
	story_points INTEGER,
	UNIQUE (proyecto_id, titulo_clave)
);`)
	if err != nil {
		t.Fatalf("crear esquema viejo: %v", err)
	}
}

// tieneColumna consulta PRAGMA table_info: tabla es siempre una constante del test, nunca un
// dato externo, así que concatenarla en el texto del PRAGMA es seguro (no admite "?").
func tieneColumna(t *testing.T, db *sql.DB, tabla, columna string) bool {
	t.Helper()
	rows, err := db.Query(`PRAGMA table_info(` + tabla + `)`)
	if err != nil {
		t.Fatalf("consultar columnas de %s: %v", tabla, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid             int
			nombre          string
			tipo            string
			notNull         int
			valorPorOmision sql.NullString
			pk              int
		)
		if err := rows.Scan(&cid, &nombre, &tipo, &notNull, &valorPorOmision, &pk); err != nil {
			t.Fatalf("leer columna de %s: %v", tabla, err)
		}
		if nombre == columna {
			return true
		}
	}
	return false
}
