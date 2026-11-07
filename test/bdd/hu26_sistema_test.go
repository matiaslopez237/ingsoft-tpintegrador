package bdd

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sistemaReporteReal prepara el proyecto con el caso de uso de HU-01, y los sprints, las historias,
// las horas reales y los defectos directamente en la base (igual que HU-22/HU-24 con los sprints),
// y consulta el reporte con los casos de uso de HU-26, sobre SQLite real, en una base nueva por
// escenario.
type sistemaReporteReal struct {
	crearProyecto   *app.CrearProyecto
	reporteProyecto *app.ArmarReporteProyecto
	reporteSprint   *app.ArmarReporteSprint
	db              *sql.DB
	dir             string
	siguienteSprint map[int64]int
	elementoDummy   map[int64]int64 // elemento de prueba por proyecto, para los defectos
}

func nuevoSistemaReporte() (sistemaReporte, error) {
	dir, err := os.MkdirTemp("", "bdd-hu26-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	sprints := persistencia.NuevoRepositorioSprints(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	defectos := persistencia.NuevoRepositorioDefectos(db)
	return &sistemaReporteReal{
		crearProyecto:   app.NuevoCrearProyecto(proyectos),
		reporteProyecto: app.NuevoArmarReporteProyecto(proyectos, backlog, sprints, esfuerzo, defectos),
		reporteSprint:   app.NuevoArmarReporteSprint(proyectos, sprints, backlog, esfuerzo, defectos),
		db:              db,
		dir:             dir,
		siguienteSprint: map[int64]int{},
		elementoDummy:   map[int64]int64{},
	}, nil
}

func (s *sistemaReporteReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

// CrearSprint inserta el sprint con el número siguiente del proyecto, igual criterio que HU-22/HU-24:
// no pasa por HU-08 ni por HU-11. Devuelve su ID.
func (s *sistemaReporteReal) CrearSprint(proyectoID int64, estado domain.EstadoSprint, planificados, completados int) (int64, error) {
	numero := s.siguienteSprint[proyectoID] + 1
	s.siguienteSprint[proyectoID] = numero
	inicio := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, numero-1, 0)
	fin := inicio.AddDate(0, 0, 10)
	r, err := s.db.Exec(
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado, puntos_planificados, puntos_completados)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, numero, inicio.Format(time.DateOnly), fin.Format(time.DateOnly), "Objetivo de prueba",
		string(estado), planificados, completados)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

// CrearHistoria inserta una historia del backlog directo en la base (no pasa por HU-05/HU-06/HU-09/
// HU-10): el escenario prueba que HU-26 junta bien los datos, no cómo se llega a ellos.
// sprintID en 0 es "sin asignar".
func (s *sistemaReporteReal) CrearHistoria(proyectoID int64, titulo string, estado domain.EstadoElemento, storyPoints *int, horasEstimadas *float64, sprintID int64) (int64, error) {
	var sp, sid sql.NullInt64
	if storyPoints != nil {
		sp = sql.NullInt64{Int64: int64(*storyPoints), Valid: true}
	}
	if sprintID != 0 {
		sid = sql.NullInt64{Int64: sprintID, Valid: true}
	}
	r, err := s.db.Exec(
		`INSERT INTO backlog_elementos (proyecto_id, titulo, titulo_clave, descripcion, prioridad, estado, story_points, sprint_id, medias_horas_estimadas)
		 VALUES (?, ?, ?, 'Historia de prueba', 'Alta', ?, ?, ?, ?)`,
		proyectoID, titulo, domain.ClaveNombre(titulo), string(estado), sp, sid, aMediasHorasReporte(horasEstimadas))
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func aMediasHorasReporte(horas *float64) sql.NullInt64 {
	if horas == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(math.Round(*horas * 2)), Valid: true}
}

// RegistrarHoras inserta un registro de esfuerzo real directo en la base (no pasa por HU-18).
func (s *sistemaReporteReal) RegistrarHoras(proyectoID, elementoID int64, horas float64) error {
	_, err := s.db.Exec(
		`INSERT INTO registros_esfuerzo (proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas)
		 VALUES (?, ?, 'Ana', '2026-01-02', 'Desarrollo', ?)`,
		proyectoID, elementoID, int(math.Round(horas*2)))
	return err
}

// CrearDefecto inserta un defecto directo en la base (no pasa por HU-20/HU-21), con el Sprint de
// detección y el de resolución que indique el escenario (0 = sin Sprint). Lo asocia a un elemento de
// prueba del proyecto, creado la primera vez que hace falta.
func (s *sistemaReporteReal) CrearDefecto(proyectoID int64, sprintDeteccionID, sprintResolucionID int64) error {
	eid, err := s.idElementoDummy(proyectoID)
	if err != nil {
		return err
	}
	var detID, resID sql.NullInt64
	if sprintDeteccionID != 0 {
		detID = sql.NullInt64{Int64: sprintDeteccionID, Valid: true}
	}
	if sprintResolucionID != 0 {
		resID = sql.NullInt64{Int64: sprintResolucionID, Valid: true}
	}
	_, err = s.db.Exec(
		`INSERT INTO defectos (proyecto_id, elemento_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id)
		 VALUES (?, ?, 'Defecto de prueba', 'Alta', 'Abierto', ?, ?)`,
		proyectoID, eid, detID, resID)
	return err
}

func (s *sistemaReporteReal) idElementoDummy(proyectoID int64) (int64, error) {
	if id, ok := s.elementoDummy[proyectoID]; ok {
		return id, nil
	}
	id, err := s.CrearHistoria(proyectoID, "Elemento para defectos de prueba", domain.EstadoPendiente, nil, nil, 0)
	if err != nil {
		return 0, err
	}
	s.elementoDummy[proyectoID] = id
	return id, nil
}

func (s *sistemaReporteReal) ConsultarReporteProyecto(proyectoID int64) (domain.Reporte, error) {
	return s.reporteProyecto.Ejecutar(context.Background(), proyectoID)
}

func (s *sistemaReporteReal) ConsultarReporteSprint(proyectoID, sprintID int64) (domain.Reporte, error) {
	return s.reporteSprint.Ejecutar(context.Background(), proyectoID, sprintID)
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaReporteReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
