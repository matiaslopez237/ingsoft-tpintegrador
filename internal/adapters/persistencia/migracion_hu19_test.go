package persistencia

import (
	"context"
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// esquemaAnteriorAHU19 es, sin cambios, el DDL de migrar() en origin/main antes de HU-19 (6767f6d,
// internal/adapters/persistencia/db.go). No se arma a partir de migrar(): tiene que seguir siendo el
// esquema viejo aunque el actual cambie.
const esquemaAnteriorAHU19 = `
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
);

-- HU-20. sprint_deteccion_id queda NULL cuando no hay Sprint vigente al registrarlo (RN-20.5).
CREATE TABLE IF NOT EXISTS defectos (
	id                  INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id         INTEGER NOT NULL REFERENCES proyectos (id),
	elemento_id         INTEGER NOT NULL REFERENCES backlog_elementos (id),
	descripcion         TEXT    NOT NULL,
	severidad           TEXT    NOT NULL,
	estado              TEXT    NOT NULL,
	sprint_deteccion_id INTEGER REFERENCES sprints (id)
);

-- HU-14. Rondas de Planning Poker. Se borran con el elemento (ON DELETE CASCADE): una ronda no
-- impide eliminarlo (HU-14, sección 5).
CREATE TABLE IF NOT EXISTS rondas (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id INTEGER NOT NULL REFERENCES proyectos (id),
	elemento_id INTEGER NOT NULL REFERENCES backlog_elementos (id) ON DELETE CASCADE,
	numero      INTEGER NOT NULL,
	estado      TEXT    NOT NULL,
	UNIQUE (elemento_id, numero)
);

-- HU-14. RN-14.1 también en la base: a lo sumo una ronda Abierta por elemento.
CREATE UNIQUE INDEX IF NOT EXISTS rondas_una_abierta_por_elemento ON rondas (elemento_id) WHERE estado = 'Abierta';

-- HU-14. Votantes de cada ronda (RN-14.2), en orden; voto NULL mientras no votó.
-- nombre_clave = domain.ClaveNombre(nombre), como en integrantes.
CREATE TABLE IF NOT EXISTS ronda_votantes (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	ronda_id     INTEGER NOT NULL REFERENCES rondas (id) ON DELETE CASCADE,
	orden        INTEGER NOT NULL,
	nombre       TEXT    NOT NULL,
	nombre_clave TEXT    NOT NULL,
	voto         INTEGER,
	UNIQUE (ronda_id, nombre_clave)
);
`

// columnaSprintAnteriorAHU19 es lo que migrarSprintID le agregaba a esa base al arrancar. Las columnas
// de migrarPuntosDeSprint ya están en el CREATE TABLE de sprints, así que no agregaba nada.
const columnaSprintAnteriorAHU19 = `ALTER TABLE backlog_elementos ADD COLUMN sprint_id INTEGER REFERENCES sprints (id);`

// datosAnterioresAHU19 son datos de una base en uso: "Login" está en un sprint y tiene 9,5 h, una
// parte de "Beto", que ya no está en el proyecto; "Reportes" no tiene sprint y tiene 2 h.
const datosAnterioresAHU19 = `
INSERT INTO proyectos (id, nombre, nombre_clave, fecha_inicio, estado) VALUES (1, 'Alfa', 'alfa', '2026-10-01', 'Activo');
INSERT INTO integrantes (id, proyecto_id, nombre, nombre_clave) VALUES (1, 1, 'Ana', 'ana');
INSERT INTO sprints (id, proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado)
VALUES (7, 1, 1, '2026-10-12', '2026-10-25', 'Primer sprint', 'Abierto');
INSERT INTO backlog_elementos (id, proyecto_id, titulo, titulo_clave, descripcion, prioridad, estado, story_points, sprint_id)
VALUES (1, 1, 'Login', 'login', 'Entrar', 'Alta', 'En progreso', 5, 7),
       (2, 1, 'Reportes', 'reportes', 'Ver', 'Media', 'Pendiente', NULL, NULL);
INSERT INTO backlog_criterios (elemento_id, orden, texto) VALUES (1, 0, 'Con email'), (2, 0, 'En PDF');
INSERT INTO registros_esfuerzo (proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas)
VALUES (1, 1, 'Ana', '2026-10-20', 'Desarrollo', 16),
       (1, 1, 'Beto', '2026-10-21', 'Testing', 3),
       (1, 2, 'Ana', '2026-10-22', 'Desarrollo', 4);`

// crearBaseAnteriorAHU19 crea el archivo con el esquema y los datos viejos, sin pasar por AbrirDB, y
// devuelve la suma de medias horas de cada elemento, leída con SQL antes de migrar.
func crearBaseAnteriorAHU19(t *testing.T, ruta string) map[int64]int {
	t.Helper()
	db, err := sql.Open("sqlite", ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("abrir base vieja: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(esquemaAnteriorAHU19 + columnaSprintAnteriorAHU19); err != nil {
		t.Fatalf("crear esquema viejo: %v", err)
	}
	if _, err := db.Exec(datosAnterioresAHU19); err != nil {
		t.Fatalf("cargar datos viejos: %v", err)
	}
	if tieneColumna(t, db, "backlog_elementos", "medias_horas_estimadas") {
		t.Fatal("el esquema viejo no debería tener medias_horas_estimadas")
	}
	return sumasDeMediasHoras(t, db)
}

func sumasDeMediasHoras(t *testing.T, db *sql.DB) map[int64]int {
	t.Helper()
	rows, err := db.Query(`SELECT elemento_id, SUM(medias_horas) FROM registros_esfuerzo GROUP BY elemento_id`)
	if err != nil {
		t.Fatalf("sumar medias horas: %v", err)
	}
	defer rows.Close()
	sumas := map[int64]int{}
	for rows.Next() {
		var (
			id   int64
			suma int
		)
		if err := rows.Scan(&id, &suma); err != nil {
			t.Fatalf("leer suma: %v", err)
		}
		sumas[id] = suma
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("recorrer sumas: %v", err)
	}
	return sumas
}

// abrirBaseAnteriorAHU19 abre con el código actual una base recién creada con el esquema viejo.
func abrirBaseAnteriorAHU19(t *testing.T) (db *sql.DB, ruta string, sumasAntes map[int64]int) {
	t.Helper()
	ruta = filepath.Join(t.TempDir(), "anterior-hu19.db")
	sumasAntes = crearBaseAnteriorAHU19(t, ruta)
	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir la base vieja con la migración actual: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db, ruta, sumasAntes
}

// Traza: HU-19 CA-19.16 (RN-19.1): los elementos se leen sin estimar y conservan todo lo demás, también en el listado.
func TestMigracionHU19_ElementosSinEstimarYConSusDatos(t *testing.T) {
	ctx := context.Background()
	db, _, _ := abrirBaseAnteriorAHU19(t)
	repo := NuevoRepositorioBacklog(db)
	cinco, sprint := 5, int64(7)
	want := []domain.ElementoBacklog{
		{ID: 1, ProyectoID: 1, Titulo: "Login", Descripcion: "Entrar", Prioridad: domain.PrioridadAlta,
			Estado: domain.EstadoEnProgreso, StoryPoints: &cinco, Criterios: []string{"Con email"}, SprintID: &sprint},
		{ID: 2, ProyectoID: 1, Titulo: "Reportes", Descripcion: "Ver", Prioridad: domain.PrioridadMedia,
			Estado: domain.EstadoPendiente, Criterios: []string{"En PDF"}},
	}

	for _, w := range want {
		got, err := repo.BuscarElementoPorID(ctx, w.ID)
		if err != nil {
			t.Fatalf("leer el elemento %d: %v", w.ID, err)
		}
		if !reflect.DeepEqual(got, w) {
			t.Errorf("elemento %d: se leyó %+v, se esperaba %+v", w.ID, got, w)
		}
	}
	lista, err := repo.ListarElementos(ctx, 1)
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if !reflect.DeepEqual(lista, want) {
		t.Errorf("listado = %+v, se esperaba %+v", lista, want)
	}
}

// Con los casos de uso reales, se cargan horas sobre la base migrada, sin perder el sprint, y la
// comparación da el resultado completo.
// Traza: HU-19 CA-19.16 (RN-19.3)
func TestMigracionHU19_SeCarganHorasYSeComparan(t *testing.T) {
	ctx := context.Background()
	db, _, _ := abrirBaseAnteriorAHU19(t)
	backlog := NuevoRepositorioBacklog(db)

	e, err := app.NuevoEditarElementoBacklog(backlog).Ejecutar(ctx, app.DatosEditarElemento{
		ProyectoID: 1, ElementoID: 1, Cambios: cambiarHoras(pHoras(8)),
	})

	if err != nil {
		t.Fatalf("cargar horas: %v", err)
	}
	if !mismasHoras(e.HorasEstimadas, pHoras(8)) {
		t.Errorf("HorasEstimadas = %v, se esperaba 8", verHoras(e.HorasEstimadas))
	}
	if e.SprintID == nil || *e.SprintID != 7 {
		t.Errorf("SprintID = %v, se esperaba 7", e.SprintID)
	}
	comparar := app.NuevoCompararEsfuerzo(NuevoRepositorioProyectos(db), backlog, NuevoRepositorioEsfuerzo(db))
	got, err := comparar.Ejecutar(ctx, 1, 1)
	if err != nil {
		t.Fatalf("comparar: %v", err)
	}
	want := domain.ComparacionEsfuerzo{HorasEstimadas: pHoras(8), HorasReales: 9.5, DesviacionHoras: pHoras(1.5), DesviacionPorcentaje: pHoras(18.8)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("comparación = %+v, se esperaba %+v", got, want)
	}
}

// La suma de horas reales es la misma antes de migrar, después de migrar, después de cargar horas
// estimadas y después de volver a abrir la base.
// Traza: HU-19 CA-19.16 (RN-19.4)
func TestMigracionHU19_LaSumaDeHorasRealesNoCambia(t *testing.T) {
	ctx := context.Background()
	db, ruta, sumasAntes := abrirBaseAnteriorAHU19(t)
	if !reflect.DeepEqual(sumasAntes, map[int64]int{1: 19, 2: 4}) {
		t.Fatalf("sumas antes de migrar = %v, se esperaban 19 y 4 medias horas", sumasAntes)
	}
	verificar := func(momento string, db *sql.DB) {
		t.Helper()
		if got := sumasDeMediasHoras(t, db); !reflect.DeepEqual(got, sumasAntes) {
			t.Errorf("%s: medias horas = %v, se esperaba %v", momento, got, sumasAntes)
		}
		repo := NuevoRepositorioEsfuerzo(db)
		for id, want := range map[int64]float64{1: 9.5, 2: 2} {
			got, err := repo.HorasRealesDeElemento(ctx, id)
			if err != nil {
				t.Fatalf("%s: horas reales del elemento %d: %v", momento, id, err)
			}
			if got != want {
				t.Errorf("%s: horas reales del elemento %d = %v, se esperaba %v", momento, id, got, want)
			}
		}
	}

	verificar("después de migrar", db)

	backlog := NuevoRepositorioBacklog(db)
	e, err := backlog.BuscarElementoPorID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	editarYGuardar(t, backlog, e, cambiarHoras(pHoras(8)))
	verificar("después de cargar horas", db)

	db.Close()
	otra, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("volver a abrir la base migrada: %v", err)
	}
	defer otra.Close()
	verificar("después de volver a abrir", otra)
	got, err := NuevoRepositorioBacklog(otra).BuscarElementoPorID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !mismasHoras(got.HorasEstimadas, pHoras(8)) || got.SprintID == nil || *got.SprintID != 7 {
		t.Errorf("después de volver a abrir: HorasEstimadas = %v, SprintID = %v", verHoras(got.HorasEstimadas), got.SprintID)
	}
}
