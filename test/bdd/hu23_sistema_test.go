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

// hu23SistemaReal conecta los escenarios de HU-23 a los casos de uso ConsultarHorasProyecto y
// ConsultarHorasSprint con SQLite real, en una base nueva por escenario. Los datos se preparan con los
// casos de uso de HU-01, HU-02, HU-03, HU-05, HU-06, HU-08, HU-09, HU-10, HU-11 y HU-18, y el elemento y
// el sprint se leen con los mismos repositorios que usan los casos de uso.
type hu23SistemaReal struct {
	crearProyecto   *app.CrearProyecto
	modificar       *app.ModificarProyecto
	agregar         *app.AgregarIntegrante
	quitar          *app.QuitarIntegrante
	crearSprint     *app.CrearSprint
	crearElemento   *app.CrearElementoBacklog
	editar          *app.EditarElementoBacklog
	registrar       *app.RegistrarEsfuerzo
	asignar         *app.AsignarHistoria
	completar       *app.CompletarHistoria
	cerrar          *app.CerrarSprint
	horasDeProyecto *app.ConsultarHorasProyecto
	horasDeSprint   *app.ConsultarHorasSprint
	repoBacklog     *persistencia.RepositorioBacklog
	repoSprints     *persistencia.RepositorioSprints
	db              *sql.DB
	dir             string
}

func hu23NuevoSistema() (hu23Sistema, error) {
	dir, err := os.MkdirTemp("", "bdd-hu23-")
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
	sprints := persistencia.NuevoRepositorioSprints(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	reloj := func() time.Time { return hu23Hoy }
	return &hu23SistemaReal{
		crearProyecto:   app.NuevoCrearProyecto(proyectos),
		modificar:       app.NuevoModificarProyecto(proyectos),
		agregar:         app.NuevoAgregarIntegrante(proyectos),
		quitar:          app.NuevoQuitarIntegrante(proyectos),
		crearSprint:     app.NuevoCrearSprint(proyectos, sprints),
		crearElemento:   app.NuevoCrearElementoBacklog(backlog),
		editar:          app.NuevoEditarElementoBacklog(backlog),
		registrar:       app.NuevoRegistrarEsfuerzo(proyectos, backlog, esfuerzo, reloj),
		asignar:         app.NuevoAsignarHistoria(backlog),
		completar:       app.NuevoCompletarHistoria(proyectos, backlog),
		cerrar:          app.NuevoCerrarSprint(proyectos, sprints),
		horasDeProyecto: app.NuevoConsultarHorasProyecto(proyectos, esfuerzo),
		horasDeSprint:   app.NuevoConsultarHorasSprint(proyectos, sprints, esfuerzo),
		repoBacklog:     backlog,
		repoSprints:     sprints,
		db:              db,
		dir:             dir,
	}, nil
}

var (
	// hu23Hoy es la fecha de los registros: cae dentro del sprint 1 (ver hu23FechasDeSprint).
	hu23Hoy = time.Date(2026, 11, 3, 0, 0, 0, 0, time.Local)
	// hu23FechaFueraDelSprint es anterior a todos los sprints y no es futura (RN-18.4, RN-23.9).
	hu23FechaFueraDelSprint = time.Date(2026, 9, 1, 0, 0, 0, 0, time.Local)
)

// hu23FechasDeSprint da dos semanas a cada sprint, sin superposición: el 1 va del 26/10 al 8/11.
func hu23FechasDeSprint(numero int) (inicio, fin time.Time) {
	inicio = time.Date(2026, 10, 26, 0, 0, 0, 0, time.Local).AddDate(0, 0, 14*(numero-1))
	return inicio, inicio.AddDate(0, 0, 13)
}

// CrearProyecto arranca en la fecha fuera del sprint, para que todos los sprints queden dentro del
// proyecto y esa fecha sea anterior a ellos.
func (s *hu23SistemaReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: hu23FechaFueraDelSprint,
	})
	return p.ID, err
}

// AgregarIntegrante usa un rol fijo: los escenarios de HU-23 no miran los roles.
func (s *hu23SistemaReal) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrolladora",
	})
}

func (s *hu23SistemaReal) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s *hu23SistemaReal) CambiarEstadoProyecto(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

// CrearSprint usa el caso de uso de HU-08: numero solo elige las fechas, y el número del sprint lo
// asigna la persistencia.
func (s *hu23SistemaReal) CrearSprint(proyectoID int64, numero int) (int64, error) {
	inicio, fin := hu23FechasDeSprint(numero)
	sp, err := s.crearSprint.Ejecutar(context.Background(), app.DatosNuevoSprint{
		ProyectoID: proyectoID, FechaInicio: inicio, FechaFin: fin,
		Objetivo: fmt.Sprintf("Objetivo del sprint %d", numero),
	})
	return sp.ID, err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-23 no miran.
func (s *hu23SistemaReal) CrearElemento(proyectoID int64, titulo string, horasEstimadas *float64) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Elemento para totalizar horas",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Se pueden totalizar sus horas"},
	})
	if err != nil {
		return 0, err
	}
	if horasEstimadas != nil {
		if err := s.EditarHorasEstimadas(proyectoID, e.ID, *horasEstimadas); err != nil {
			return 0, err
		}
	}
	return e.ID, nil
}

func (s *hu23SistemaReal) RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64, fecha time.Time) error {
	_, err := s.registrar.Ejecutar(context.Background(), app.DatosRegistroEsfuerzo{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Integrante: integrante,
		Fecha:      fecha,
		Actividad:  "Desarrollo",
		Horas:      horas,
	})
	return err
}

func (s *hu23SistemaReal) EditarHorasEstimadas(proyectoID, elementoID int64, horas float64) error {
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: &horas})
}

func (s *hu23SistemaReal) EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error {
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &storyPoints})
}

func (s *hu23SistemaReal) PasarAEnProgreso(proyectoID, elementoID int64) error {
	estado := domain.EstadoEnProgreso
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{Estado: &estado})
}

func (s *hu23SistemaReal) editarElemento(proyectoID, elementoID int64, c domain.CambiosElemento) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID, ElementoID: elementoID, Cambios: c,
	})
	return err
}

func (s *hu23SistemaReal) Asignar(proyectoID, elementoID, sprintID int64) error {
	return s.asignar.Ejecutar(context.Background(), app.DatosAsignacion{
		ProyectoID: proyectoID, ElementoID: elementoID, SprintID: sprintID,
	})
}

func (s *hu23SistemaReal) Completar(proyectoID, elementoID int64) error {
	return s.completar.Ejecutar(context.Background(), app.DatosCompletar{ProyectoID: proyectoID, ElementoID: elementoID})
}

func (s *hu23SistemaReal) CerrarSprint(proyectoID, sprintID int64) error {
	return s.cerrar.Ejecutar(context.Background(), app.DatosCerrarSprint{ProyectoID: proyectoID, SprintID: sprintID})
}

func (s *hu23SistemaReal) HorasDeProyecto(proyectoID int64) (domain.TotalHoras, error) {
	return s.horasDeProyecto.Ejecutar(context.Background(), proyectoID)
}

func (s *hu23SistemaReal) HorasDeSprint(proyectoID, sprintID int64) (domain.TotalHoras, error) {
	return s.horasDeSprint.Ejecutar(context.Background(), proyectoID, sprintID)
}

func (s *hu23SistemaReal) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return s.repoBacklog.BuscarElementoPorID(context.Background(), elementoID)
}

func (s *hu23SistemaReal) Sprint(sprintID int64) (domain.Sprint, error) {
	return s.repoSprints.BuscarSprintPorID(context.Background(), sprintID)
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *hu23SistemaReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
