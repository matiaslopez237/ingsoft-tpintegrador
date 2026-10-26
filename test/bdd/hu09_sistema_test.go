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

// errAsignacionNoImplementada hace fallar los escenarios de HU-09 mientras asignar/desasignar
// un elemento a un sprint no existe (rojo esperado, Tareas 9.4 a 9.8).
var errAsignacionNoImplementada = errors.New("pendiente: asignar/desasignar historias a un sprint todavía no está implementado")

// erroresHU09 traduce la descripción del escenario al error definido en la spec.
// Se completa cuando existan los errores del dominio (Tarea 9.4).
func erroresHU09() map[string]error {
	return map[string]error{}
}

// sprintFalso representa un sprint solo para los escenarios de HU-09: todavía no existe un
// domain.Sprint real (HU-08 no está mergeada), así que esta es la única noción de "sprint" que
// usamos acá, sin depender de su tipo.
type sprintFalso struct {
	proyectoID int64
	abierto    bool
}

// sistemaAsignacionReal crea proyectos (HU-01) y elementos de backlog (HU-05) reales, y lleva
// los sprints en memoria (ver sprintFalso) porque HU-08 todavía no está mergeada.
type sistemaAsignacionReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	db            *sql.DB
	dir           string
	sprints       map[int64]sprintFalso
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
	return sistemaAsignacionReal{
		crearProyecto: app.NuevoCrearProyecto(persistencia.NuevoRepositorioProyectos(db)),
		crearElemento: app.NuevoCrearElementoBacklog(persistencia.NuevoRepositorioBacklog(db)),
		db:            db,
		dir:           dir,
		sprints:       map[int64]sprintFalso{},
	}, nil
}

func (s sistemaAsignacionReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaAsignacionReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Now(),
	})
	return p.ID, err
}

func (s sistemaAsignacionReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Descripción de prueba",
		Prioridad:   domain.PrioridadMedia,
		Criterios:   []string{"Criterio de prueba"},
	})
	return e.ID, err
}

func (s sistemaAsignacionReal) CrearSprint(proyectoID int64, abierto bool) int64 {
	id := int64(len(s.sprints) + 1)
	s.sprints[id] = sprintFalso{proyectoID: proyectoID, abierto: abierto}
	return id
}

func (s sistemaAsignacionReal) AsignarElemento(elementoID, sprintID int64) error {
	return errAsignacionNoImplementada
}

func (s sistemaAsignacionReal) DesasignarElemento(elementoID int64) error {
	return errAsignacionNoImplementada
}

func (s sistemaAsignacionReal) SprintDeElemento(elementoID int64) (sprintID int64, asignado bool, err error) {
	return 0, false, errAsignacionNoImplementada
}
