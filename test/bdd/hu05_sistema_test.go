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

// erroresHU05 traduce la descripción del escenario al error definido en la spec.
func erroresHU05() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"título vacío":           domain.ErrTituloVacio,
		"título demasiado largo": domain.ErrTituloDemasiadoLargo,
		"descripción vacía":      domain.ErrDescripcionVacia,
		"prioridad inválida":     domain.ErrPrioridadInvalida,
		"story points inválidos": domain.ErrStoryPointsInvalidos,
		"sin criterios":          domain.ErrSinCriterios,
		"criterio vacío":         domain.ErrCriterioVacio,
		"título duplicado":       domain.ErrTituloDuplicado,
	}
}

// sistemaBacklogReal crea los proyectos con CrearProyecto (HU-01) y los elementos con
// CrearElementoBacklog (HU-05), sobre SQLite real, en una base nueva por escenario.
type sistemaBacklogReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	repo          *persistencia.RepositorioBacklog
	db            *sql.DB
	dir           string
}

func nuevoSistemaBacklog() (sistemaBacklog, error) {
	dir, err := os.MkdirTemp("", "bdd-hu05-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repo := persistencia.NuevoRepositorioBacklog(db)
	return sistemaBacklogReal{
		crearProyecto: app.NuevoCrearProyecto(persistencia.NuevoRepositorioProyectos(db)),
		crearElemento: app.NuevoCrearElementoBacklog(repo),
		repo:          repo,
		db:            db,
		dir:           dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaBacklogReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaBacklogReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

func (s sistemaBacklogReal) CrearElemento(proyectoID int64, d datosElemento) (elementoVisto, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      d.Titulo,
		Descripcion: d.Descripcion,
		Prioridad:   domain.Prioridad(d.Prioridad),
		StoryPoints: d.StoryPoints,
		Criterios:   d.Criterios,
	})
	if err != nil {
		return elementoVisto{}, err
	}
	return verElemento(e), nil
}

func (s sistemaBacklogReal) BuscarElemento(id int64) (elementoVisto, error) {
	e, err := s.repo.BuscarElementoPorID(context.Background(), id)
	if err != nil {
		return elementoVisto{}, err
	}
	return verElemento(e), nil
}

func (s sistemaBacklogReal) ContarElementos(proyectoID int64) (int, error) {
	return s.repo.ContarElementos(context.Background(), proyectoID)
}

func verElemento(e domain.ElementoBacklog) elementoVisto {
	v := elementoVisto{ID: e.ID, Titulo: e.Titulo, Estado: string(e.Estado), Criterios: e.Criterios}
	if e.StoryPoints != nil {
		v.Estimado, v.StoryPoints = true, *e.StoryPoints
	}
	return v
}
