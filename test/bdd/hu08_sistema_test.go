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

// erroresHU08 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU08() map[string]error {
	return map[string]error{
		"proyecto no encontrado":      domain.ErrProyectoNoEncontrado,
		"proyecto no activo":          domain.ErrProyectoNoActivo,
		"fecha de inicio obligatoria": domain.ErrFechaInicioObligatoria,
		"fecha de fin obligatoria":    domain.ErrFechaFinObligatoria,
		"fecha de fin anterior":       domain.ErrFechaFinAnterior,
		"Sprint Goal vacío":           domain.ErrObjetivoVacio,
		"Sprint Goal demasiado largo": domain.ErrObjetivoDemasiadoLargo,
		"sprint fuera del proyecto":   domain.ErrSprintFueraDelProyecto,
		"sprint superpuesto":          domain.ErrSprintSuperpuesto,
	}
}

// sistemaSprintsReal prepara los proyectos con los casos de uso de HU-01 y HU-02, y crea y lista
// sprints con los de HU-08, sobre SQLite real, en una base nueva por escenario.
type sistemaSprintsReal struct {
	crearProyecto     *app.CrearProyecto
	modificarProyecto *app.ModificarProyecto
	crearSprint       *app.CrearSprint
	listarSprints     *app.ListarSprints
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
	sprints := persistencia.NuevoRepositorioSprints(db)
	return sistemaSprintsReal{
		crearProyecto:     app.NuevoCrearProyecto(proyectos),
		modificarProyecto: app.NuevoModificarProyecto(proyectos),
		crearSprint:       app.NuevoCrearSprint(proyectos, sprints),
		listarSprints:     app.NuevoListarSprints(proyectos, sprints),
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

func (s sistemaSprintsReal) CrearSprint(proyectoID int64, inicio, fin *time.Time, objetivo string) (sprintVisto, error) {
	datos := app.DatosNuevoSprint{ProyectoID: proyectoID, Objetivo: objetivo}
	if inicio != nil {
		datos.FechaInicio = *inicio
	}
	if fin != nil {
		datos.FechaFin = *fin
	}
	sp, err := s.crearSprint.Ejecutar(context.Background(), datos)
	if err != nil {
		return sprintVisto{}, err
	}
	return verSprint(sp), nil
}

func (s sistemaSprintsReal) ListarSprints(proyectoID int64) ([]sprintVisto, error) {
	sprints, err := s.listarSprints.Ejecutar(context.Background(), proyectoID)
	if err != nil {
		return nil, err
	}
	vistos := make([]sprintVisto, len(sprints))
	for i, sp := range sprints {
		vistos[i] = verSprint(sp)
	}
	return vistos, nil
}

func verSprint(sp domain.Sprint) sprintVisto {
	return sprintVisto{
		Numero:   sp.Numero,
		Inicio:   sp.FechaInicio.Format(time.DateOnly),
		Fin:      sp.FechaFin.Format(time.DateOnly),
		Objetivo: sp.Objetivo,
		Estado:   string(sp.Estado),
	}
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaSprintsReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
