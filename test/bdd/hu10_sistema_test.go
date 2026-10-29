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

// sistemaCompletarReal conecta los escenarios al caso de uso CompletarHistoria (HU-10) con SQLite
// real, sobre los repositorios reales de HU-01, HU-05 y HU-08, en una base nueva por escenario.
type sistemaCompletarReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	crearSprint   *app.CrearSprint
	completar     *app.CompletarHistoria
	repoBacklog   *persistencia.RepositorioBacklog
	db            *sql.DB
	dir           string
	siguienteDia  int // para que los sprints no se superpongan (RN-08.5)
}

func nuevoSistemaCompletar() (sistemaCompletar, error) {
	dir, err := os.MkdirTemp("", "bdd-hu10-")
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
	return &sistemaCompletarReal{
		crearProyecto: app.NuevoCrearProyecto(repoProyectos),
		crearElemento: app.NuevoCrearElementoBacklog(repoBacklog),
		crearSprint:   app.NuevoCrearSprint(repoProyectos, repoSprints),
		completar:     app.NuevoCompletarHistoria(repoProyectos, repoBacklog),
		repoBacklog:   repoBacklog,
		db:            db,
		dir:           dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaCompletarReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s *sistemaCompletarReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Now(),
	})
	return p.ID, err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-10 no miran.
func (s *sistemaCompletarReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Descripción de prueba",
		Prioridad:   domain.PrioridadMedia,
		Criterios:   []string{"Criterio de prueba"},
	})
	return e.ID, err
}

// CrearSprint crea un sprint real (HU-08) de un día, espaciados dos días entre sí para que nunca
// se superpongan (RN-08.5). Si abierto es false, lo cierra directamente en la base: HU-11 (cerrar
// sprint) todavía no está en esta rama, así que no hay otra forma de dejarlo así.
func (s *sistemaCompletarReal) CrearSprint(proyectoID int64, abierto bool) (int64, error) {
	dia := time.Now().AddDate(0, 0, s.siguienteDia)
	s.siguienteDia += 2
	sp, err := s.crearSprint.Ejecutar(context.Background(), app.DatosNuevoSprint{
		ProyectoID:  proyectoID,
		FechaInicio: dia,
		FechaFin:    dia,
		Objetivo:    "Objetivo de prueba",
	})
	if err != nil {
		return 0, err
	}
	if !abierto {
		if _, err := s.db.Exec(`UPDATE sprints SET estado = ? WHERE id = ?`, string(domain.SprintCerrado), sp.ID); err != nil {
			return 0, fmt.Errorf("cerrar sprint de prueba: %w", err)
		}
	}
	return sp.ID, nil
}

// PrepararElemento deja el elemento en el sprint (nil: sin sprint) y el estado indicados,
// directamente en la base, sin pasar por los casos de uso: como AsignarDirectamente en HU-09, los
// escenarios necesitan, por ejemplo, un elemento en un sprint que después se cerró.
func (s *sistemaCompletarReal) PrepararElemento(elementoID int64, sprintID *int64, estado domain.EstadoElemento) error {
	res, err := s.db.Exec(`UPDATE backlog_elementos SET sprint_id = ?, estado = ? WHERE id = ?`,
		sprintID, string(estado), elementoID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrElementoNoEncontrado
	}
	return nil
}

func (s *sistemaCompletarReal) CompletarHistoria(proyectoID, elementoID int64) error {
	return s.completar.Ejecutar(context.Background(), app.DatosCompletar{ProyectoID: proyectoID, ElementoID: elementoID})
}

func (s *sistemaCompletarReal) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return s.repoBacklog.BuscarElementoPorID(context.Background(), elementoID)
}
