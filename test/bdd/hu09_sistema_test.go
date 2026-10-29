package bdd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU09 traduce la descripción del escenario al error definido en la spec.
func erroresHU09() map[string]error {
	return map[string]error{
		"elemento ya asignado":    domain.ErrElementoYaAsignado,
		"sprint cerrado":          domain.ErrSprintCerrado,
		"elemento sin sprint":     domain.ErrElementoSinSprint,
		"sprint de otro proyecto": domain.ErrSprintDeOtroProyecto,
	}
}

// sistemaAsignacionReal conecta los escenarios a los casos de uso de HU-09 con SQLite real, sobre
// los repositorios reales de HU-01, HU-05 y HU-08.
type sistemaAsignacionReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	crearSprint   *app.CrearSprint
	asignar       *app.AsignarHistoria
	desasignar    *app.DesasignarHistoria
	repoBacklog   *persistencia.RepositorioBacklog
	db            *sql.DB
	dir           string
	siguienteDia  int // para que los sprints de un mismo proyecto no se superpongan (RN-08.5)
}

func nuevoSistemaAsignacion() (sistemaAsignacion, error) {
	dir, err := os.MkdirTemp("", "bdd-hu09-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repoProyectos := persistencia.NuevoRepositorioProyectos(db)
	repoBacklog := persistencia.NuevoRepositorioBacklog(db)
	repoSprints := persistencia.NuevoRepositorioSprints(db)
	return &sistemaAsignacionReal{
		crearProyecto: app.NuevoCrearProyecto(repoProyectos),
		crearElemento: app.NuevoCrearElementoBacklog(repoBacklog),
		crearSprint:   app.NuevoCrearSprint(repoProyectos, repoSprints),
		asignar:       app.NuevoAsignarHistoria(repoBacklog),
		desasignar:    app.NuevoDesasignarHistoria(repoBacklog),
		repoBacklog:   repoBacklog,
		db:            db,
		dir:           dir,
	}, nil
}

func (s *sistemaAsignacionReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s *sistemaAsignacionReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Now(),
	})
	return p.ID, err
}

func (s *sistemaAsignacionReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Descripción de prueba",
		Prioridad:   domain.PrioridadMedia,
		Criterios:   []string{"Criterio de prueba"},
	})
	return e.ID, err
}

// CrearSprint crea un sprint real (HU-08) de un día. Los sprints se van espaciando dos días entre
// sí para que nunca se superpongan (RN-08.5), aunque sean de proyectos distintos. Si abierto es
// false, lo cierra directamente en la base: HU-11 (cerrar sprint) todavía no existe, así que no
// hay otra forma de dejarlo así.
func (s *sistemaAsignacionReal) CrearSprint(proyectoID int64, abierto bool) int64 {
	dia := time.Now().AddDate(0, 0, s.siguienteDia)
	s.siguienteDia += 2
	sp, err := s.crearSprint.Ejecutar(context.Background(), app.DatosNuevoSprint{
		ProyectoID:  proyectoID,
		FechaInicio: dia,
		FechaFin:    dia,
		Objetivo:    "Objetivo de prueba",
	})
	if err != nil {
		panic(fmt.Sprintf("crear sprint de prueba: %v", err)) // no hay forma de devolver el error (interfaz de CrearSprint)
	}
	if !abierto {
		if _, err := s.db.Exec(`UPDATE sprints SET estado = ? WHERE id = ?`, string(domain.SprintCerrado), sp.ID); err != nil {
			panic(fmt.Sprintf("cerrar sprint de prueba: %v", err))
		}
	}
	return sp.ID
}

// AsignarDirectamente fuerza la asignación en la base, sin pasar por el caso de uso: lo usan los
// escenarios para dejar preparado un elemento "ya asignado" a un sprint que ahora está cerrado
// (en la vida real, la asignación pasó antes de que el sprint se cerrara).
func (s *sistemaAsignacionReal) AsignarDirectamente(elementoID, sprintID int64) error {
	_, err := s.db.Exec(`UPDATE backlog_elementos SET sprint_id = ? WHERE id = ?`, sprintID, elementoID)
	return err
}

func (s *sistemaAsignacionReal) AsignarElemento(proyectoID, elementoID, sprintID int64) error {
	return s.asignar.Ejecutar(context.Background(), app.DatosAsignacion{
		ProyectoID: proyectoID, ElementoID: elementoID, SprintID: sprintID,
	})
}

func (s *sistemaAsignacionReal) DesasignarElemento(proyectoID, elementoID int64) error {
	return s.desasignar.Ejecutar(context.Background(), app.DatosDesasignacion{
		ProyectoID: proyectoID, ElementoID: elementoID,
	})
}

func (s *sistemaAsignacionReal) SprintDeElemento(elementoID int64) (sprintID int64, asignado bool, err error) {
	e, err := s.repoBacklog.BuscarElementoPorID(context.Background(), elementoID)
	if err != nil {
		return 0, false, err
	}
	if e.SprintID == nil {
		return 0, false, nil
	}
	return *e.SprintID, true, nil
}
