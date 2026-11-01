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

// errHU22SinImplementar es el rojo de la tarea 22.3: los escenarios compilan y corren, pero el
// caso de uso todavía no existe. Se va cuando las tareas 22.4 y 22.5 conecten ConsultarVelocidad.
var errHU22SinImplementar = errors.New("HU-22 todavía no está implementada: falta conectar el caso de uso ConsultarVelocidad (tareas 22.4 y 22.5)")

// sistemaVelocidadReal prepara los proyectos con el caso de uso de HU-01 y los sprints
// directamente en la base (igual que HU-12), sobre SQLite real, en una base nueva por escenario.
// Consultar la velocidad todavía devuelve errHU22SinImplementar.
type sistemaVelocidadReal struct {
	crearProyecto *app.CrearProyecto
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

func (s *sistemaVelocidadReal) ConsultarVelocidad(int64) (velocidadVista, error) {
	return velocidadVista{}, errHU22SinImplementar
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaVelocidadReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
