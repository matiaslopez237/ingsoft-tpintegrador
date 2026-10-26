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

// errSprintsNoImplementado hace fallar los escenarios mientras CrearSprint no existe (rojo esperado).
var errSprintsNoImplementado = errors.New("pendiente: CrearSprint y ListarSprints todavía no están implementados")

// erroresHU08 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU08() map[string]error {
	return map[string]error{
		"proyecto no encontrado":      domain.ErrProyectoNoEncontrado,
		"fecha de inicio obligatoria": domain.ErrFechaInicioObligatoria,
		"fecha de fin anterior":       domain.ErrFechaFinAnterior,
	}
}

// sistemaSprintsReal prepara los proyectos con los casos de uso de HU-01 y HU-02 sobre SQLite,
// en una base nueva por escenario. Crear y listar sprints todavía no está implementado.
type sistemaSprintsReal struct {
	crearProyecto     *app.CrearProyecto
	modificarProyecto *app.ModificarProyecto
	db                *sql.DB
	dir               string
}

func nuevoSistemaSprints() (sistemaSprints, error) {
	dir, err := os.MkdirTemp("", "bdd-hu08-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	return sistemaSprintsReal{
		crearProyecto:     app.NuevoCrearProyecto(proyectos),
		modificarProyecto: app.NuevoModificarProyecto(proyectos),
		db:                db,
		dir:               dir,
	}, nil
}

func (s sistemaSprintsReal) CrearProyecto(nombre string, inicio time.Time, fin *time.Time) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{Nombre: nombre, FechaInicio: inicio, FechaFin: fin})
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (s sistemaSprintsReal) FijarEstadoProyecto(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificarProyecto.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

func (s sistemaSprintsReal) CrearSprint(int64, *time.Time, *time.Time, string) (sprintVisto, error) {
	return sprintVisto{}, errSprintsNoImplementado
}

func (s sistemaSprintsReal) ListarSprints(int64) ([]sprintVisto, error) {
	return nil, errSprintsNoImplementado
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaSprintsReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
