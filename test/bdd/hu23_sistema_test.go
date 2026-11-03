package bdd

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RED de HU-23: no compila hasta que existan domain.TotalHoras, domain.HorasElemento, app.LectorHoras,
// app.LectorSprintPorID, app.ConsultarHorasProyecto y app.ConsultarHorasSprint (lista completa en
// hu23_steps_test.go).

// hu23SistemaMemoria conecta los escenarios de HU-23 a los casos de uso reales, sobre hu23RepoMemoria.
// Los datos se preparan con los casos de uso de HU-01, HU-02, HU-03, HU-05, HU-06, HU-09, HU-10, HU-11
// y HU-18. El sprint se guarda directamente porque sus fechas no son lo que se prueba. Más adelante
// pasa a SQLite, como en HU-19: la interfaz hu23Sistema no cambia.
type hu23SistemaMemoria struct {
	crearProyecto   *app.CrearProyecto
	modificar       *app.ModificarProyecto
	agregar         *app.AgregarIntegrante
	quitar          *app.QuitarIntegrante
	crearElemento   *app.CrearElementoBacklog
	editar          *app.EditarElementoBacklog
	registrar       *app.RegistrarEsfuerzo
	asignar         *app.AsignarHistoria
	completar       *app.CompletarHistoria
	cerrar          *app.CerrarSprint
	horasDeProyecto *app.ConsultarHorasProyecto
	horasDeSprint   *app.ConsultarHorasSprint
	repo            *hu23RepoMemoria
}

func hu23NuevoSistema() (hu23Sistema, error) {
	repo := hu23NuevoRepoMemoria()
	reloj := func() time.Time { return hu23Hoy }
	return &hu23SistemaMemoria{
		crearProyecto:   app.NuevoCrearProyecto(repo),
		modificar:       app.NuevoModificarProyecto(repo),
		agregar:         app.NuevoAgregarIntegrante(repo),
		quitar:          app.NuevoQuitarIntegrante(repo),
		crearElemento:   app.NuevoCrearElementoBacklog(repo),
		editar:          app.NuevoEditarElementoBacklog(repo),
		registrar:       app.NuevoRegistrarEsfuerzo(repo, repo, repo, reloj),
		asignar:         app.NuevoAsignarHistoria(repo),
		completar:       app.NuevoCompletarHistoria(repo, repo),
		cerrar:          app.NuevoCerrarSprint(repo, repo),
		horasDeProyecto: app.NuevoConsultarHorasProyecto(repo, repo),
		horasDeSprint:   app.NuevoConsultarHorasSprint(repo, repo, repo),
		repo:            repo,
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

func (s *hu23SistemaMemoria) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: hu23FechaFueraDelSprint,
	})
	return p.ID, err
}

// AgregarIntegrante usa un rol fijo: los escenarios de HU-23 no miran los roles.
func (s *hu23SistemaMemoria) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrolladora",
	})
}

func (s *hu23SistemaMemoria) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s *hu23SistemaMemoria) CambiarEstadoProyecto(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

func (s *hu23SistemaMemoria) CrearSprint(proyectoID int64, numero int) (int64, error) {
	inicio, fin := hu23FechasDeSprint(numero)
	return s.repo.guardarSprint(domain.Sprint{
		ProyectoID: proyectoID, Numero: numero, FechaInicio: inicio, FechaFin: fin,
		Objetivo: fmt.Sprintf("Objetivo del sprint %d", numero), Estado: domain.SprintAbierto,
	}), nil
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-23 no miran.
func (s *hu23SistemaMemoria) CrearElemento(proyectoID int64, titulo string, horasEstimadas *float64) (int64, error) {
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

func (s *hu23SistemaMemoria) RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64, fecha time.Time) error {
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

func (s *hu23SistemaMemoria) EditarHorasEstimadas(proyectoID, elementoID int64, horas float64) error {
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: &horas})
}

func (s *hu23SistemaMemoria) EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error {
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &storyPoints})
}

func (s *hu23SistemaMemoria) PasarAEnProgreso(proyectoID, elementoID int64) error {
	estado := domain.EstadoEnProgreso
	return s.editarElemento(proyectoID, elementoID, domain.CambiosElemento{Estado: &estado})
}

func (s *hu23SistemaMemoria) editarElemento(proyectoID, elementoID int64, c domain.CambiosElemento) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID, ElementoID: elementoID, Cambios: c,
	})
	return err
}

func (s *hu23SistemaMemoria) Asignar(proyectoID, elementoID, sprintID int64) error {
	return s.asignar.Ejecutar(context.Background(), app.DatosAsignacion{
		ProyectoID: proyectoID, ElementoID: elementoID, SprintID: sprintID,
	})
}

func (s *hu23SistemaMemoria) Completar(proyectoID, elementoID int64) error {
	return s.completar.Ejecutar(context.Background(), app.DatosCompletar{ProyectoID: proyectoID, ElementoID: elementoID})
}

func (s *hu23SistemaMemoria) CerrarSprint(proyectoID, sprintID int64) error {
	return s.cerrar.Ejecutar(context.Background(), app.DatosCerrarSprint{ProyectoID: proyectoID, SprintID: sprintID})
}

func (s *hu23SistemaMemoria) HorasDeProyecto(proyectoID int64) (domain.TotalHoras, error) {
	return s.horasDeProyecto.Ejecutar(context.Background(), proyectoID)
}

func (s *hu23SistemaMemoria) HorasDeSprint(proyectoID, sprintID int64) (domain.TotalHoras, error) {
	return s.horasDeSprint.Ejecutar(context.Background(), proyectoID, sprintID)
}

func (s *hu23SistemaMemoria) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return s.repo.BuscarElementoPorID(context.Background(), elementoID)
}

func (s *hu23SistemaMemoria) Sprint(sprintID int64) (domain.Sprint, error) {
	return s.repo.BuscarSprintPorID(context.Background(), sprintID)
}

// Cerrar no tiene nada que liberar en memoria; existe porque el sistema sobre SQLite lo va a usar.
func (s *hu23SistemaMemoria) Cerrar() error { return nil }

// hu23RepoMemoria guarda todo en memoria y cumple los puertos de los casos de uso que usa el sistema.
type hu23RepoMemoria struct {
	proyectos map[int64]domain.Proyecto
	elementos map[int64]domain.ElementoBacklog
	sprints   map[int64]domain.Sprint
	registros []domain.RegistroEsfuerzo
	ultimoID  int64
}

// Puertos de HU-23 y de los casos de uso con que se preparan los datos.
var (
	_ app.LectorProyectos           = (*hu23RepoMemoria)(nil)
	_ app.LectorHoras               = (*hu23RepoMemoria)(nil)
	_ app.LectorSprintPorID         = (*hu23RepoMemoria)(nil)
	_ app.RepositorioProyectos      = (*hu23RepoMemoria)(nil)
	_ app.RepositorioModificacion   = (*hu23RepoMemoria)(nil)
	_ app.RepositorioIntegrantes    = (*hu23RepoMemoria)(nil)
	_ app.RepositorioBacklog        = (*hu23RepoMemoria)(nil)
	_ app.RepositorioEdicionBacklog = (*hu23RepoMemoria)(nil)
	_ app.RepositorioEsfuerzo       = (*hu23RepoMemoria)(nil)
	_ app.RepositorioAsignacion     = (*hu23RepoMemoria)(nil)
	_ app.RepositorioCompletar      = (*hu23RepoMemoria)(nil)
	_ app.RepositorioCierre         = (*hu23RepoMemoria)(nil)
)

func hu23NuevoRepoMemoria() *hu23RepoMemoria {
	return &hu23RepoMemoria{
		proyectos: map[int64]domain.Proyecto{},
		elementos: map[int64]domain.ElementoBacklog{},
		sprints:   map[int64]domain.Sprint{},
	}
}

func (r *hu23RepoMemoria) nuevoID() int64 {
	r.ultimoID++
	return r.ultimoID
}

// --- Proyectos e integrantes ---

// BuscarPorID devuelve una copia de los integrantes: los casos de uso de HU-03 modifican la lista
// que reciben, y lo guardado solo tiene que cambiar por AgregarIntegrante y QuitarIntegrante.
func (r *hu23RepoMemoria) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	p, ok := r.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	p.Integrantes = slices.Clone(p.Integrantes)
	return p, nil
}

func (r *hu23RepoMemoria) ExisteNombre(_ context.Context, clave string) (bool, error) {
	for _, p := range r.proyectos {
		if domain.ClaveNombre(p.Nombre) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *hu23RepoMemoria) Guardar(_ context.Context, p *domain.Proyecto) error {
	p.ID = r.nuevoID()
	r.proyectos[p.ID] = *p
	return nil
}

// Actualizar solo aplica el estado: es el único cambio de proyecto que usan los escenarios de HU-23.
func (r *hu23RepoMemoria) Actualizar(_ context.Context, p domain.Proyecto, c domain.CambiosProyecto) error {
	guardado, ok := r.proyectos[p.ID]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	if c.Estado != nil {
		guardado.Estado = *c.Estado
	}
	r.proyectos[p.ID] = guardado
	return nil
}

func (r *hu23RepoMemoria) AgregarIntegrante(_ context.Context, proyectoID int64, i domain.Integrante) error {
	p, ok := r.proyectos[proyectoID]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	p.Integrantes = append(slices.Clone(p.Integrantes), i)
	r.proyectos[proyectoID] = p
	return nil
}

func (r *hu23RepoMemoria) QuitarIntegrante(_ context.Context, proyectoID int64, nombre string) error {
	p, ok := r.proyectos[proyectoID]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	clave := domain.ClaveNombre(nombre)
	p.Integrantes = slices.DeleteFunc(slices.Clone(p.Integrantes), func(i domain.Integrante) bool {
		return domain.ClaveNombre(i.Nombre) == clave
	})
	r.proyectos[proyectoID] = p
	return nil
}

// AgregarRol completa el puerto de HU-03: los escenarios de HU-23 no miran los roles.
func (r *hu23RepoMemoria) AgregarRol(context.Context, int64, string, string) error {
	return nil
}

// --- Elementos del backlog ---

func (r *hu23RepoMemoria) ExisteProyecto(_ context.Context, id int64) (bool, error) {
	_, ok := r.proyectos[id]
	return ok, nil
}

func (r *hu23RepoMemoria) ExisteTitulo(_ context.Context, proyectoID int64, clave string) (bool, error) {
	for _, e := range r.elementos {
		if e.ProyectoID == proyectoID && domain.ClaveNombre(e.Titulo) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *hu23RepoMemoria) GuardarElemento(_ context.Context, e *domain.ElementoBacklog) error {
	e.ID = r.nuevoID()
	r.elementos[e.ID] = *e
	return nil
}

func (r *hu23RepoMemoria) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	e, ok := r.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

// ActualizarElemento guarda solo los campos que c indica como enviados, como el repositorio SQLite
// (HU-06, sección 5). Los escenarios de HU-23 solo editan estos tres.
func (r *hu23RepoMemoria) ActualizarElemento(_ context.Context, e domain.ElementoBacklog, c domain.CambiosElemento) error {
	guardado, ok := r.elementos[e.ID]
	if !ok {
		return domain.ErrElementoNoEncontrado
	}
	if c.CambiarStoryPoints {
		guardado.StoryPoints = e.StoryPoints
	}
	if c.CambiarHorasEstimadas {
		guardado.HorasEstimadas = e.HorasEstimadas
	}
	if c.Estado != nil {
		guardado.Estado = e.Estado
	}
	r.elementos[e.ID] = guardado
	return nil
}

// --- Sprints: asignación (HU-09), completado (HU-10) y cierre (HU-11) ---

func (r *hu23RepoMemoria) guardarSprint(s domain.Sprint) int64 {
	s.ID = r.nuevoID()
	r.sprints[s.ID] = s
	return s.ID
}

func (r *hu23RepoMemoria) BuscarSprint(_ context.Context, id int64) (proyectoID int64, abierto bool, err error) {
	s, ok := r.sprints[id]
	if !ok {
		return 0, false, domain.ErrSprintNoEncontrado
	}
	return s.ProyectoID, s.Estado == domain.SprintAbierto, nil
}

func (r *hu23RepoMemoria) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	s, ok := r.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return s, nil
}

func (r *hu23RepoMemoria) GuardarAsignacion(_ context.Context, e *domain.ElementoBacklog) error {
	guardado, ok := r.elementos[e.ID]
	if !ok {
		return domain.ErrElementoNoEncontrado
	}
	guardado.SprintID = e.SprintID
	r.elementos[e.ID] = guardado
	return nil
}

func (r *hu23RepoMemoria) GuardarCompletado(_ context.Context, e *domain.ElementoBacklog) error {
	guardado, ok := r.elementos[e.ID]
	if !ok {
		return domain.ErrElementoNoEncontrado
	}
	guardado.Estado = e.Estado
	r.elementos[e.ID] = guardado
	return nil
}

// CerrarSprint aplica domain.ResumirCierre, como el repositorio SQLite: guarda los puntos y devuelve
// al backlog las historias sin terminar (RN-11.3, RN-11.6).
func (r *hu23RepoMemoria) CerrarSprint(_ context.Context, s *domain.Sprint) error {
	guardado, ok := r.sprints[s.ID]
	if !ok {
		return domain.ErrSprintNoEncontrado
	}
	if guardado.Estado != domain.SprintAbierto {
		return domain.ErrSprintYaCerrado
	}
	var historias []domain.ElementoBacklog
	for _, e := range r.elementos {
		if e.SprintID != nil && *e.SprintID == s.ID {
			historias = append(historias, e)
		}
	}
	resumen := domain.ResumirCierre(historias)
	for _, id := range resumen.DevueltasAlBacklog {
		e := r.elementos[id]
		e.SprintID = nil
		r.elementos[id] = e
	}
	guardado.Estado = domain.SprintCerrado
	guardado.PuntosPlanificados, guardado.PuntosCompletados = resumen.PuntosPlanificados, resumen.PuntosCompletados
	r.sprints[s.ID] = guardado
	s.PuntosPlanificados, s.PuntosCompletados = resumen.PuntosPlanificados, resumen.PuntosCompletados
	return nil
}

// --- Esfuerzo ---

func (r *hu23RepoMemoria) GuardarRegistro(_ context.Context, reg *domain.RegistroEsfuerzo) error {
	reg.ID = r.nuevoID()
	r.registros = append(r.registros, *reg)
	return nil
}

// --- Horas (HU-23) ---

func (r *hu23RepoMemoria) HorasDeProyecto(_ context.Context, proyectoID int64) ([]domain.HorasElemento, error) {
	return r.horasDe(func(e domain.ElementoBacklog) bool { return e.ProyectoID == proyectoID }), nil
}

func (r *hu23RepoMemoria) HorasDeSprint(_ context.Context, sprintID int64) ([]domain.HorasElemento, error) {
	return r.horasDe(func(e domain.ElementoBacklog) bool { return e.SprintID != nil && *e.SprintID == sprintID }), nil
}

// horasDe devuelve una fila por cada elemento elegido, también si no tiene registros (spec, sección 5).
func (r *hu23RepoMemoria) horasDe(elegido func(domain.ElementoBacklog) bool) []domain.HorasElemento {
	var filas []domain.HorasElemento
	for _, e := range r.elementos {
		if !elegido(e) {
			continue
		}
		var reales float64
		for _, reg := range r.registros {
			if reg.ElementoID == e.ID {
				reales += reg.Horas
			}
		}
		filas = append(filas, domain.HorasElemento{HorasEstimadas: e.HorasEstimadas, HorasReales: reales})
	}
	return filas
}
