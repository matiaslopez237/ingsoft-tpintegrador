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

// errHU20SinImplementar es el rojo de la tarea 20.3: los escenarios compilan y corren, pero el
// caso de uso todavía no existe. Se va cuando las tareas 20.4 a 20.6 conecten RegistrarDefecto.
var errHU20SinImplementar = errors.New("HU-20 todavía no está implementada: falta conectar el caso de uso RegistrarDefecto (tareas 20.4 a 20.6)")

// defectoErrorEsperado es un error de la spec de HU-20 (sección 7). Los reutilizados de otras
// historias ya existen en el dominio y se comparan con errors.Is. Los nuevos todavía no existen,
// así que se comparan por el mensaje de la spec, igual que hizo HU-18 en su tarea 18.3.
// TODO(HU-20, tarea 20.4): cuando el dominio defina los errores nuevos, completar err y quitar mensaje.
type defectoErrorEsperado struct {
	err     error
	mensaje string
}

// coincide recorre la cadena de errores envueltos (%w) buscando el error esperado.
func (e defectoErrorEsperado) coincide(got error) bool {
	if e.err != nil {
		return errors.Is(got, e.err)
	}
	for ; got != nil; got = errors.Unwrap(got) {
		if got.Error() == e.mensaje {
			return true
		}
	}
	return false
}

// erroresHU20 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU20() map[string]defectoErrorEsperado {
	return map[string]defectoErrorEsperado{
		"proyecto no encontrado": {err: domain.ErrProyectoNoEncontrado},
		"elemento no encontrado": {err: domain.ErrElementoNoEncontrado},
		"descripción vacía":      {mensaje: "la descripción es obligatoria"},                     // ErrDescripcionDefectoVacia
		"severidad inválida":     {mensaje: "la severidad debe ser Crítica, Alta, Media o Baja"}, // ErrSeveridadInvalida
	}
}

// sistemaDefectoReal prepara los proyectos y elementos con los casos de uso de HU-01 y HU-05, y
// los sprints con los de HU-08, sobre SQLite real, en una base nueva por escenario. Registrar un
// defecto todavía devuelve errHU20SinImplementar; reloj queda listo para cuando las tareas 20.4 a
// 20.6 conecten RegistrarDefecto, igual que el reloj falso de HU-18.
type sistemaDefectoReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	crearSprint   *app.CrearSprint
	reloj         func() time.Time
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
	return sistemaDefectoReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(persistencia.NuevoRepositorioBacklog(db)),
		crearSprint:   app.NuevoCrearSprint(proyectos, persistencia.NuevoRepositorioSprints(db)),
		reloj:         reloj,
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

func (s sistemaDefectoReal) RegistrarDefecto(datosDefecto) (defectoVisto, error) {
	return defectoVisto{}, errHU20SinImplementar
}
