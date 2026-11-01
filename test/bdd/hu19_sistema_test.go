package bdd

import (
	"context"
	"slices"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RED de HU-19: no compila hasta que existan app.CompararEsfuerzo, app.LectorHorasReales,
// domain.ComparacionEsfuerzo, ElementoBacklog.HorasEstimadas y los campos nuevos de
// domain.CambiosElemento (lista completa en hu19_steps_test.go).

// hu19SistemaMemoria conecta los escenarios de HU-19 a los casos de uso reales, sobre
// hu19RepoMemoria. Los datos se preparan con los casos de uso de HU-01, HU-02, HU-03, HU-05, HU-06 y
// HU-18. Más adelante pasa a SQLite, como en HU-11: la interfaz hu19Sistema no cambia.
type hu19SistemaMemoria struct {
	crearProyecto *app.CrearProyecto
	modificar     *app.ModificarProyecto
	agregar       *app.AgregarIntegrante
	quitar        *app.QuitarIntegrante
	crearElemento *app.CrearElementoBacklog
	editar        *app.EditarElementoBacklog
	registrar     *app.RegistrarEsfuerzo
	comparar      *app.CompararEsfuerzo
	repo          *hu19RepoMemoria
}

func hu19NuevoSistema() (hu19Sistema, error) {
	repo := hu19NuevoRepoMemoria()
	reloj := func() time.Time { return hu19Hoy }
	return &hu19SistemaMemoria{
		crearProyecto: app.NuevoCrearProyecto(repo),
		modificar:     app.NuevoModificarProyecto(repo),
		agregar:       app.NuevoAgregarIntegrante(repo),
		quitar:        app.NuevoQuitarIntegrante(repo),
		crearElemento: app.NuevoCrearElementoBacklog(repo),
		editar:        app.NuevoEditarElementoBacklog(repo),
		registrar:     app.NuevoRegistrarEsfuerzo(repo, repo, repo, reloj),
		comparar:      app.NuevoCompararEsfuerzo(repo, repo, repo),
		repo:          repo,
	}, nil
}

// hu19Hoy es la fecha de los registros: los escenarios de HU-19 no miran las fechas.
var hu19Hoy = time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local)

func (s *hu19SistemaMemoria) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: hu19Hoy,
	})
	return p.ID, err
}

// AgregarIntegrante usa un rol fijo: los escenarios de HU-19 no miran los roles.
func (s *hu19SistemaMemoria) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrolladora",
	})
}

func (s *hu19SistemaMemoria) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s *hu19SistemaMemoria) CambiarEstadoProyecto(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-19 no miran.
func (s *hu19SistemaMemoria) CrearElemento(proyectoID int64, titulo string, storyPoints *int, horasEstimadas *float64) (int64, error) {
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

func (s *hu19SistemaMemoria) RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64) error {
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

func (s *hu19SistemaMemoria) EditarHorasEstimadas(proyectoID, elementoID int64, horas *float64) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Cambios:    domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: horas},
	})
	return err
}

func (s *hu19SistemaMemoria) EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error {
	_, err := s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{
		ProyectoID: proyectoID,
		ElementoID: elementoID,
		Cambios:    domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &storyPoints},
	})
	return err
}

func (s *hu19SistemaMemoria) Comparar(proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error) {
	return s.comparar.Ejecutar(context.Background(), proyectoID, elementoID)
}

func (s *hu19SistemaMemoria) Elemento(elementoID int64) (domain.ElementoBacklog, error) {
	return s.repo.BuscarElementoPorID(context.Background(), elementoID)
}

// Cerrar no tiene nada que liberar en memoria; existe porque el sistema sobre SQLite lo va a usar.
func (s *hu19SistemaMemoria) Cerrar() error { return nil }

// hu19RepoMemoria guarda todo en memoria y cumple los puertos de los casos de uso que usa el sistema.
type hu19RepoMemoria struct {
	proyectos map[int64]domain.Proyecto
	elementos map[int64]domain.ElementoBacklog
	registros []domain.RegistroEsfuerzo
	ultimoID  int64
}

// Puertos de HU-19 y de los casos de uso con que se preparan los datos.
var (
	_ app.LectorProyectos           = (*hu19RepoMemoria)(nil)
	_ app.LectorElementos           = (*hu19RepoMemoria)(nil)
	_ app.LectorHorasReales         = (*hu19RepoMemoria)(nil)
	_ app.RepositorioProyectos      = (*hu19RepoMemoria)(nil)
	_ app.RepositorioModificacion   = (*hu19RepoMemoria)(nil)
	_ app.RepositorioIntegrantes    = (*hu19RepoMemoria)(nil)
	_ app.RepositorioBacklog        = (*hu19RepoMemoria)(nil)
	_ app.RepositorioEdicionBacklog = (*hu19RepoMemoria)(nil)
	_ app.RepositorioEsfuerzo       = (*hu19RepoMemoria)(nil)
)

func hu19NuevoRepoMemoria() *hu19RepoMemoria {
	return &hu19RepoMemoria{
		proyectos: map[int64]domain.Proyecto{},
		elementos: map[int64]domain.ElementoBacklog{},
	}
}

func (r *hu19RepoMemoria) nuevoID() int64 {
	r.ultimoID++
	return r.ultimoID
}

// --- Proyectos e integrantes ---

// BuscarPorID devuelve una copia de los integrantes: los casos de uso de HU-03 modifican la lista
// que reciben, y lo guardado solo tiene que cambiar por AgregarIntegrante y QuitarIntegrante.
func (r *hu19RepoMemoria) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	p, ok := r.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	p.Integrantes = slices.Clone(p.Integrantes)
	return p, nil
}

func (r *hu19RepoMemoria) ExisteNombre(_ context.Context, clave string) (bool, error) {
	for _, p := range r.proyectos {
		if domain.ClaveNombre(p.Nombre) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *hu19RepoMemoria) Guardar(_ context.Context, p *domain.Proyecto) error {
	p.ID = r.nuevoID()
	r.proyectos[p.ID] = *p
	return nil
}

// Actualizar solo aplica el estado: es el único cambio de proyecto que usan los escenarios de HU-19.
func (r *hu19RepoMemoria) Actualizar(_ context.Context, p domain.Proyecto, c domain.CambiosProyecto) error {
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

func (r *hu19RepoMemoria) AgregarIntegrante(_ context.Context, proyectoID int64, i domain.Integrante) error {
	p, ok := r.proyectos[proyectoID]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	p.Integrantes = append(slices.Clone(p.Integrantes), i)
	r.proyectos[proyectoID] = p
	return nil
}

func (r *hu19RepoMemoria) QuitarIntegrante(_ context.Context, proyectoID int64, nombre string) error {
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

// AgregarRol completa el puerto de HU-03: los escenarios de HU-19 no miran los roles.
func (r *hu19RepoMemoria) AgregarRol(context.Context, int64, string, string) error {
	return nil
}

// --- Elementos del backlog ---

func (r *hu19RepoMemoria) ExisteProyecto(_ context.Context, id int64) (bool, error) {
	_, ok := r.proyectos[id]
	return ok, nil
}

func (r *hu19RepoMemoria) ExisteTitulo(_ context.Context, proyectoID int64, clave string) (bool, error) {
	for _, e := range r.elementos {
		if e.ProyectoID == proyectoID && domain.ClaveNombre(e.Titulo) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *hu19RepoMemoria) GuardarElemento(_ context.Context, e *domain.ElementoBacklog) error {
	e.ID = r.nuevoID()
	r.elementos[e.ID] = *e
	return nil
}

func (r *hu19RepoMemoria) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	e, ok := r.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

// ActualizarElemento guarda solo los campos que c indica como enviados, como el repositorio SQLite
// (HU-06, sección 5). Los escenarios de HU-19 solo editan los story points y las horas estimadas.
func (r *hu19RepoMemoria) ActualizarElemento(_ context.Context, e domain.ElementoBacklog, c domain.CambiosElemento) error {
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
	r.elementos[e.ID] = guardado
	return nil
}

// --- Esfuerzo ---

func (r *hu19RepoMemoria) GuardarRegistro(_ context.Context, reg *domain.RegistroEsfuerzo) error {
	reg.ID = r.nuevoID()
	r.registros = append(r.registros, *reg)
	return nil
}

// HorasRealesDeElemento suma los registros del elemento; 0 si no tiene ninguno (RN-19.4).
func (r *hu19RepoMemoria) HorasRealesDeElemento(_ context.Context, elementoID int64) (float64, error) {
	var total float64
	for _, reg := range r.registros {
		if reg.ElementoID == elementoID {
			total += reg.Horas
		}
	}
	return total, nil
}
