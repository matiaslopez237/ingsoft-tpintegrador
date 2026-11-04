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

// hu19SistemaReal conecta los escenarios de HU-19 al caso de uso CompararEsfuerzo con SQLite real, en
// una base nueva por escenario. Los datos se preparan con los casos de uso de HU-01, HU-02, HU-03,
// HU-05, HU-06 y HU-18, y el elemento se lee con el mismo repositorio que usan los casos de uso.
type hu19SistemaReal struct {
	crearProyecto *app.CrearProyecto
	modificar     *app.ModificarProyecto
	agregar       *app.AgregarIntegrante
	quitar        *app.QuitarIntegrante
	crearElemento *app.CrearElementoBacklog
	editar        *app.EditarElementoBacklog
	registrar     *app.RegistrarEsfuerzo
	comparar      *app.CompararEsfuerzo
	repoBacklog   *persistencia.RepositorioBacklog
	db            *sql.DB
	dir           string
}

func hu19NuevoSistema() (hu19Sistema, error) {
	dir, err := os.MkdirTemp("", "bdd-hu19-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	reloj := func() time.Time { return hu19Hoy }
	return &hu19SistemaReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		modificar:     app.NuevoModificarProyecto(proyectos),
		agregar:       app.NuevoAgregarIntegrante(proyectos),
		quitar:        app.NuevoQuitarIntegrante(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(backlog),
		editar:        app.NuevoEditarElementoBacklog(backlog),
		registrar:     app.NuevoRegistrarEsfuerzo(proyectos, backlog, esfuerzo, reloj),
		comparar:      app.NuevoCompararEsfuerzo(proyectos, backlog, esfuerzo),
		repoBacklog:   backlog,
		db:            db,
		dir:           dir,
	}, nil
}

// hu19Hoy es la fecha de los registros: los escenarios de HU-19 no miran las fechas.
var hu19Hoy = time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local)

func (s *hu19SistemaReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: hu19Hoy,
	})
	return p.ID, err
}

// AgregarIntegrante usa un rol fijo: los escenarios de HU-19 no miran los roles.
func (s *hu19SistemaReal) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrolladora",
	})
}

func (s *hu19SistemaReal) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s *hu19SistemaReal) CambiarEstadoProyecto(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-19 no miran.
func (s *hu19SistemaReal) CrearElemento(proyectoID int64, titulo string, storyPoints *int, horasEstimadas *float64) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Elemento para comparar esfuerzo",
		Prioridad:   domain.PrioridadAlta,
		StoryPoints: storyPoints,
		Criterios:   []string{"Se puede comparar el esfuerzo"},
	})
	if err != nil {
		return 0, err
	}
	if horasEstimadas != nil {
		if err := s.EditarHorasEstimadas(proyectoID, e.ID, horasEstimadas); err != nil {
			return 0, err
		}
	}
	return e.ID, nil
}

func (s *hu19SistemaReal) RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64) error {
	_, err := s.registrar.Ejecutar(context.Background(), app.DatosRegistroEsfuerzo{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Integrante: integrante,
		Fecha:      hu19Hoy,
		Actividad:  "Desarrollo",
		Horas:      horas,
	})
	return err
}

func (s *hu19SistemaReal) EditarHorasEstimadas(proyectoID, elementoID int64, horas *float64) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Cambios:    domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: horas},
	})
	return err
}

func (s *hu19SistemaReal) EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Cambios:    domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &storyPoints},
	})
	return err
}

func (s *hu19SistemaReal) Comparar(proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error) {
	return s.comparar.Ejecutar(context.Background(), proyectoID, elementoID)
}

func (s *hu19SistemaReal) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return s.repoBacklog.BuscarElementoPorID(context.Background(), elementoID)
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *hu19SistemaReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
