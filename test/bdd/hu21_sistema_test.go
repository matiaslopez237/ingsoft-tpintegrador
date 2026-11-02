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

// erroresHU21 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU21() map[string]defectoErrorEsperado {
	return map[string]defectoErrorEsperado{
		"proyecto no encontrado":      {err: domain.ErrProyectoNoEncontrado},
		"defecto no encontrado":       {err: domain.ErrDefectoNoEncontrado},
		"el defecto no está Abierto":  {err: domain.ErrDefectoNoAbierto},
		"el defecto no está Resuelto": {err: domain.ErrDefectoNoResuelto},
	}
}

// defectoErrorEsperado es un error de la spec de HU-21 (sección 7), comparado con errors.Is.
type defectoErrorEsperado struct {
	err error
}

func (e defectoErrorEsperado) coincide(got error) bool {
	return errors.Is(got, e.err)
}

// sistemaEstadoDefectoReal prepara proyectos, elementos, sprints y defectos con los casos de uso
// ya mergeados (HU-01, HU-05, HU-08, HU-20) y actualiza su estado con el de HU-21
// (ActualizarEstadoDefecto), sobre SQLite real, en una base nueva por escenario.
type sistemaEstadoDefectoReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	crearSprint   *app.CrearSprint
	registrar     *app.RegistrarDefecto
	actualizar    *app.ActualizarEstadoDefecto
	db            *sql.DB
	dir           string
}

func nuevoSistemaEstadoDefecto(reloj func() time.Time) (sistemaEstadoDefecto, error) {
	dir, err := os.MkdirTemp("", "bdd-hu21-")
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
	defectos := persistencia.NuevoRepositorioDefectos(db)
	return sistemaEstadoDefectoReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(elementos),
		crearSprint:   app.NuevoCrearSprint(proyectos, sprints),
		registrar:     app.NuevoRegistrarDefecto(proyectos, elementos, sprints, defectos, reloj),
		actualizar:    app.NuevoActualizarEstadoDefecto(proyectos, sprints, defectos, reloj),
		db:            db,
		dir:           dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaEstadoDefectoReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaEstadoDefectoReal) CrearProyecto(nombre string, inicio time.Time) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: inicio,
	})
	return p.ID, err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-21 no miran.
func (s sistemaEstadoDefectoReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Elemento para actualizar el estado de un defecto",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Se puede actualizar el estado de un defecto"},
	})
	return e.ID, err
}

// CrearSprint completa el Sprint Goal con un valor fijo que los escenarios de HU-21 no miran.
func (s sistemaEstadoDefectoReal) CrearSprint(proyectoID int64, inicio, fin time.Time) error {
	_, err := s.crearSprint.Ejecutar(context.Background(), app.DatosNuevoSprint{
		ProyectoID: proyectoID, FechaInicio: inicio, FechaFin: fin, Objetivo: "Sprint para HU-21",
	})
	return err
}

// CrearDefecto registra el defecto de partida de cada escenario con el caso de uso ya mergeado de
// HU-20: esta historia parte siempre de un defecto Abierto existente.
func (s sistemaEstadoDefectoReal) CrearDefecto(d datosDefecto) (defectoVisto, error) {
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

func (s sistemaEstadoDefectoReal) ResolverDefecto(d datosActualizarDefecto) (defectoVisto, error) {
	r, err := s.actualizar.Resolver(context.Background(), app.DatosActualizarEstadoDefecto{
		ProyectoID: d.ProyectoID, DefectoID: d.DefectoID,
	})
	return aDefectoVisto(r), err
}

func (s sistemaEstadoDefectoReal) CerrarDefecto(d datosActualizarDefecto) (defectoVisto, error) {
	r, err := s.actualizar.Cerrar(context.Background(), app.DatosActualizarEstadoDefecto{
		ProyectoID: d.ProyectoID, DefectoID: d.DefectoID,
	})
	return aDefectoVisto(r), err
}

func aDefectoVisto(r domain.Defecto) defectoVisto {
	return defectoVisto{
		ID: r.ID, ProyectoID: r.ProyectoID, ElementoID: r.ElementoID,
		Descripcion: r.Descripcion, Severidad: string(r.Severidad), Estado: string(r.Estado),
		SprintDeteccionID: r.SprintDeteccionID, SprintResolucionID: r.SprintResolucionID,
	}
}
