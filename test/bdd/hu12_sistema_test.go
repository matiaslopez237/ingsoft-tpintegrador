package bdd

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sistemaHistorialReal prepara los proyectos con el caso de uso de HU-01 y los sprints
// directamente en la base, sobre SQLite real, en una base nueva por escenario.
type sistemaHistorialReal struct {
	crearProyecto *app.CrearProyecto
	db            *sql.DB
	dir           string
}

func nuevoSistemaHistorial() (sistemaHistorial, error) {
	dir, err := os.MkdirTemp("", "bdd-hu12-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return sistemaHistorialReal{
		crearProyecto: app.NuevoCrearProyecto(persistencia.NuevoRepositorioProyectos(db)),
		db:            db,
		dir:           dir,
	}, nil
}

func (s sistemaHistorialReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

// CrearSprint inserta el sprint con el número, las fechas, el estado y los puntos del escenario:
// HU-08 calcula el número solo y crea sprints Abiertos, y los puntos los registraría el cierre de
// HU-11 a partir de las historias, que acá no importan.
func (s sistemaHistorialReal) CrearSprint(proyectoID int64, numero int, inicio, fin time.Time, objetivo string,
	estado domain.EstadoSprint, planificados, completados int) error {
	_, err := s.db.Exec(
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado, puntos_planificados, puntos_completados)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, numero, inicio.Format(time.DateOnly), fin.Format(time.DateOnly), objetivo, string(estado),
		planificados, completados)
	return err
}

// ConsultarAnteriores queda pendiente hasta que exista el caso de uso ListarSprintsAnteriores
// (tarea 12.5): mientras tanto, los escenarios quedan en rojo.
func (s sistemaHistorialReal) ConsultarAnteriores(proyectoID int64) ([]sprintAnteriorVisto, error) {
	return nil, godog.ErrPending
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaHistorialReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
