package persistencia

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
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

// Traza: HU-19 CA-19.16: una base con elementos, creada antes de HU-19, recibe la columna al abrirla.
// Sus elementos se leen sin estimar y se les pueden cargar horas, que siguen ahí al volver a abrirla.
func TestAbrirDB_BaseExistenteSinHorasEstimadas_AgregaLaColumna(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "vieja.db")
	crearEsquemaSinSprintID(t, ruta)
	vieja, err := sql.Open("sqlite", ruta)
	if err != nil {
		t.Fatalf("abrir base vieja: %v", err)
	}
	_, err = vieja.Exec(`
CREATE TABLE backlog_criterios (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	elemento_id INTEGER NOT NULL REFERENCES backlog_elementos (id) ON DELETE CASCADE,
	orden       INTEGER NOT NULL,
	texto       TEXT    NOT NULL,
	UNIQUE (elemento_id, orden)
);
INSERT INTO proyectos (nombre, nombre_clave, fecha_inicio, estado) VALUES ('Alfa', 'alfa', '2026-10-01', 'Activo');
INSERT INTO backlog_elementos (proyecto_id, titulo, titulo_clave, descripcion, prioridad, estado, story_points)
VALUES (1, 'Login', 'login', 'Entrar', 'Alta', 'En progreso', 5);
INSERT INTO backlog_criterios (elemento_id, orden, texto) VALUES (1, 0, 'Con email');`)
	vieja.Close()
	if err != nil {
		t.Fatalf("crear datos viejos: %v", err)
	}

	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir base vieja con la migración actual: %v", err)
	}
	defer func() { db.Close() }()

	if !tieneColumna(t, db, "backlog_elementos", "medias_horas_estimadas") {
		t.Fatal("se esperaba que backlog_elementos tuviera la columna medias_horas_estimadas después de migrar")
	}
	repo := NuevoRepositorioBacklog(db)
	e, err := repo.BuscarElementoPorID(ctx, 1)
	if err != nil {
		t.Fatalf("leer el elemento viejo: %v", err)
	}
	if e.Titulo != "Login" || e.Estado != domain.EstadoEnProgreso || e.StoryPoints == nil || *e.StoryPoints != 5 || e.HorasEstimadas != nil {
		t.Errorf("elemento viejo = %+v, se esperaba el mismo, sin estimar", e)
	}
	lista, err := repo.ListarElementos(ctx, 1)
	if err != nil {
		t.Fatalf("listar los elementos viejos: %v", err)
	}
	if len(lista) != 1 || lista[0].HorasEstimadas != nil {
		t.Errorf("listado = %+v, se esperaba un elemento sin estimar", lista)
	}

	got := editarYGuardar(t, repo, e, cambiarHoras(pHoras(8)))
	if !mismasHoras(got.HorasEstimadas, pHoras(8)) {
		t.Fatalf("HorasEstimadas = %v, se esperaba 8", verHoras(got.HorasEstimadas))
	}

	db.Close()
	if db, err = AbrirDB(ruta); err != nil {
		t.Fatalf("volver a abrir la base migrada: %v", err)
	}
	got, err = NuevoRepositorioBacklog(db).BuscarElementoPorID(ctx, 1)
	if err != nil {
		t.Fatalf("releer después de volver a abrir: %v", err)
	}
	if !mismasHoras(got.HorasEstimadas, pHoras(8)) {
		t.Errorf("después de volver a abrir: HorasEstimadas = %v, se esperaba 8", verHoras(got.HorasEstimadas))
	}
}

// Traza: HU-19 CA-19.16: una base nueva nace con la columna, sin pasar por la migración.
func TestAbrirDB_BaseNueva_TieneLaColumnaDeHorasEstimadas(t *testing.T) {
	db, err := AbrirDB(filepath.Join(t.TempDir(), "nueva.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	defer db.Close()

	if !tieneColumna(t, db, "backlog_elementos", "medias_horas_estimadas") {
		t.Error("se esperaba que backlog_elementos tuviera la columna medias_horas_estimadas")
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
