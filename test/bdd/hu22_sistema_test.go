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

// sistemaVelocidadReal prepara los proyectos con el caso de uso de HU-01 y los sprints
// directamente en la base (igual que HU-12), y consulta la velocidad con el caso de uso de HU-22,
// sobre SQLite real, en una base nueva por escenario.
type sistemaVelocidadReal struct {
	crearProyecto *app.CrearProyecto
	consultar     *app.ConsultarVelocidad
	db            *sql.DB
	dir           string
	siguiente     map[int64]int // próximo número de sprint de cada proyecto
}

func nuevoSistemaVelocidad() (sistemaVelocidad, error) {
	dir, err := os.MkdirTemp("", "bdd-hu22-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	return &sistemaVelocidadReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		consultar:     app.NuevoConsultarVelocidad(proyectos, persistencia.NuevoRepositorioSprints(db)),
		db:            db,
		dir:           dir,
		siguiente:     map[int64]int{},
	}, nil
}

func (s *sistemaVelocidadReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

// CrearSprint inserta el sprint con el número siguiente del proyecto, fechas no superpuestas
// (un mes de separación por número) y, si está cerrado, los puntos completados del escenario: los
// planificados no le importan a la velocidad, así que van iguales a los completados.
func (s *sistemaVelocidadReal) CrearSprint(proyectoID int64, estado domain.EstadoSprint, completados int) error {
	numero := s.siguiente[proyectoID] + 1
	s.siguiente[proyectoID] = numero
	inicio := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, numero-1, 0)
	fin := inicio.AddDate(0, 0, 10)
	_, err := s.db.Exec(
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado, puntos_planificados, puntos_completados)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, numero, inicio.Format(time.DateOnly), fin.Format(time.DateOnly), "Objetivo de prueba",
		string(estado), completados, completados)
	return err
}

func (s *sistemaVelocidadReal) ConsultarVelocidad(proyectoID int64) (velocidadVista, error) {
	v, err := s.consultar.Ejecutar(context.Background(), proyectoID)
	if err != nil {
		return velocidadVista{}, err
	}
	return velocidadVista{Velocidad: v.Velocidad, HayDatos: v.HayDatos, SprintsCerrados: v.SprintsCerrados}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaVelocidadReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
