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

// erroresHU20 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU20() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"elemento no encontrado": domain.ErrElementoNoEncontrado,
		"descripción vacía":      domain.ErrDescripcionDefectoVacia,
		"severidad inválida":     domain.ErrSeveridadInvalida,
	}
}

// sistemaDefectoReal prepara los proyectos y elementos con los casos de uso de HU-01 y HU-05, los
// sprints con los de HU-08, y registra defectos con los de HU-20 (RegistrarDefecto y
// RepositorioDefectos), sobre SQLite real, en una base nueva por escenario.
type sistemaDefectoReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	crearSprint   *app.CrearSprint
	registrar     *app.RegistrarDefecto
	db            *sql.DB
	dir           string
}

func nuevoSistemaDefecto(reloj func() time.Time) (sistemaDefecto, error) {
	dir, err := os.MkdirTemp("", "bdd-hu20-")
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
	sprints := persistencia.NuevoRepositorioSprints(db)
	return sistemaDefectoReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(elementos),
		crearSprint:   app.NuevoCrearSprint(proyectos, sprints),
		registrar:     app.NuevoRegistrarDefecto(proyectos, elementos, sprints, persistencia.NuevoRepositorioDefectos(db), reloj),
		db:            db,
		dir:           dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaDefectoReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaDefectoReal) CrearProyecto(nombre string, inicio time.Time) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: inicio,
	})
	return p.ID, err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-20 no miran.
func (s sistemaDefectoReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Elemento para registrar un defecto",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Se puede registrar un defecto"},
	})
	return e.ID, err
}

// CrearSprint completa el Sprint Goal con un valor fijo que los escenarios de HU-20 no miran.
func (s sistemaDefectoReal) CrearSprint(proyectoID int64, inicio, fin time.Time) error {
	_, err := s.crearSprint.Ejecutar(context.Background(), app.DatosNuevoSprint{
		ProyectoID: proyectoID, FechaInicio: inicio, FechaFin: fin, Objetivo: "Sprint para HU-20",
	})
	return err
}

func (s sistemaDefectoReal) RegistrarDefecto(d datosDefecto) (defectoVisto, error) {
	r, err := s.registrar.Ejecutar(context.Background(), app.DatosDefecto{
		ProyectoID: d.ProyectoID, ElementoID: d.ElementoID, Descripcion: d.Descripcion, Severidad: domain.Severidad(d.Severidad),
	})
	if err != nil {
		return defectoVisto{}, err
	}
	return defectoVisto{
		ID: r.ID, ProyectoID: r.ProyectoID, ElementoID: r.ElementoID,
		Descripcion: r.Descripcion, Severidad: string(r.Severidad), Estado: string(r.Estado),
		SprintDeteccionID: r.SprintDeteccionID,
	}, nil
}
