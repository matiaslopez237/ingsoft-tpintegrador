package bdd

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sistemaCalidadSprintReal prepara el proyecto con el caso de uso de HU-01, el elemento con el de
// HU-05, y los sprints y los defectos directamente en la base (igual que HU-22 con los sprints), y
// consulta la calidad con el caso de uso de HU-24, sobre SQLite real, en una base nueva por
// escenario.
type sistemaCalidadSprintReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	consultar     *app.ConsultarCalidadSprints
	db            *sql.DB
	dir           string
	siguiente     map[int64]int // próximo número de sprint de cada proyecto
}

func nuevoSistemaCalidadSprint() (sistemaCalidadSprint, error) {
	dir, err := os.MkdirTemp("", "bdd-hu24-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	elementos := persistencia.NuevoRepositorioBacklog(db)
	defectos := persistencia.NuevoRepositorioDefectos(db)
	return &sistemaCalidadSprintReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(elementos),
		consultar:     app.NuevoConsultarCalidadSprints(proyectos, persistencia.NuevoRepositorioSprints(db), defectos),
		db:            db,
		dir:           dir,
		siguiente:     map[int64]int{},
	}, nil
}

func (s *sistemaCalidadSprintReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-24 no miran: lo
// necesita un defecto para su referencia obligatoria (RN-20.1).
func (s *sistemaCalidadSprintReal) CrearElemento(proyectoID int64) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      "Elemento para la calidad del sprint",
		Descripcion: "Elemento de prueba",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Sirve para adjuntar defectos de prueba"},
	})
	return e.ID, err
}

// CrearSprint inserta el sprint con el número siguiente del proyecto, fechas no superpuestas (un
// mes de separación por número) y sus puntos planificados/completados, igual criterio que HU-22:
// no pasa por HU-08 ni por HU-11. Devuelve su ID, para que los defectos lo referencien.
func (s *sistemaCalidadSprintReal) CrearSprint(proyectoID int64, estado domain.EstadoSprint, planificados, completados int) (int64, error) {
	numero := s.siguiente[proyectoID] + 1
	s.siguiente[proyectoID] = numero
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

// CrearDefecto inserta el defecto directo en la base, con el estado, el Sprint de detección y el de
// resolución que indique el escenario (0 = sin Sprint): no pasa por HU-20 ni por HU-21, que asignan
// el Sprint automáticamente según la fecha de hoy, algo que estos escenarios no necesitan probar de
// nuevo (ya lo cubren las features de esas historias).
func (s *sistemaCalidadSprintReal) CrearDefecto(proyectoID, elementoID int64, estado domain.EstadoDefecto, sprintDeteccionID, sprintResolucionID int64) error {
	_, err := s.db.Exec(
		`INSERT INTO defectos (proyecto_id, elemento_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id)
		 VALUES (?, ?, 'Defecto de prueba', 'Alta', ?, ?, ?)`,
		proyectoID, elementoID, string(estado), nuloSiCero(sprintDeteccionID), nuloSiCero(sprintResolucionID))
	return err
}

func (s *sistemaCalidadSprintReal) ConsultarCalidad(proyectoID int64) ([]calidadVista, error) {
	resumen, err := s.consultar.Ejecutar(context.Background(), proyectoID)
	if err != nil {
		return nil, err
	}
	vistas := make([]calidadVista, len(resumen))
	for i, r := range resumen {
		vistas[i] = calidadVista{
			Numero: r.Sprint.Numero, PorcentajeCompletado: r.PorcentajeCompletado, HayPorcentaje: r.HayPorcentaje,
			DefectosDetectados: r.DefectosDetectados, DefectosResueltos: r.DefectosResueltos,
		}
	}
	return vistas, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaCalidadSprintReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

// nuloSiCero traduce 0 (sin Sprint, en los steps) a NULL en la base: 0 nunca es un ID real de
// Sprint (SQLite autoincremental arranca en 1).
func nuloSiCero(id int64) any {
	if id == 0 {
		return nil
	}
	return id
}
