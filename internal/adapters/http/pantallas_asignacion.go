package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Datos de la pantalla de detalle del sprint (HU-09): solo texto ya formateado, sin lógica de negocio.
type (
	vistaElementoBacklog struct {
		ID     int64
		Titulo string
		Estado string
		// PuedeCompletarse decide si se muestra "Marcar como hecha" (HU-10): sprint Abierto y
		// elemento no Hecho. Es solo para la vista; la regla la vuelve a validar el dominio.
		PuedeCompletarse bool
	}
	datosSprintDetalle struct {
		Proyecto    vistaProyecto
		Sprint      vistaSprint
		Asignados   []vistaElementoBacklog
		Disponibles []vistaElementoBacklog
		Error       string
	}
)

// sprintDetalle muestra el sprint, sus elementos asignados y los disponibles para sumar.
// Implementa HU-09 (Tarea 9.8).
func (p *pantallas) sprintDetalle(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	sprintID, err := sprintIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return
	}
	p.mostrarSprintDetalle(w, r, id, sprintID, http.StatusOK, "")
}

// asignarElementoPantalla asigna el elemento del formulario a este sprint, con el caso de uso de
// HU-09, y vuelve al detalle del sprint. Implementa HU-09 (CA-09.1 a CA-09.3, CA-09.7, CA-09.8).
func (p *pantallas) asignarElementoPantalla(w http.ResponseWriter, r *http.Request) {
	id, sprintID, elementoID, ok := p.leerAccionSprint(w, r)
	if !ok {
		return
	}
	err := p.cu.AsignarHistoria.Ejecutar(r.Context(),
		app.DatosAsignacion{ProyectoID: id, ElementoID: elementoID, SprintID: sprintID})
	p.resultadoAsignacion(w, r, id, sprintID, err)
}

// quitarElementoPantalla desasigna el elemento del formulario de este sprint, con el caso de uso
// de HU-09, y vuelve al detalle del sprint. Implementa HU-09 (CA-09.4 a CA-09.6).
func (p *pantallas) quitarElementoPantalla(w http.ResponseWriter, r *http.Request) {
	id, sprintID, elementoID, ok := p.leerAccionSprint(w, r)
	if !ok {
		return
	}
	err := p.cu.DesasignarHistoria.Ejecutar(r.Context(), app.DatosDesasignacion{ProyectoID: id, ElementoID: elementoID})
	p.resultadoAsignacion(w, r, id, sprintID, err)
}

// completarElementoPantalla marca como Hecho el elemento del formulario, con el mismo caso de uso
// que el endpoint JSON de HU-10, y vuelve al detalle del sprint. Implementa HU-10 (CA-10.7).
func (p *pantallas) completarElementoPantalla(w http.ResponseWriter, r *http.Request) {
	id, sprintID, elementoID, ok := p.leerAccionSprint(w, r)
	if !ok {
		return
	}
	err := p.cu.CompletarHistoria.Ejecutar(r.Context(), app.DatosCompletar{ProyectoID: id, ElementoID: elementoID})
	if err == nil {
		redirigirAlDetalleDelSprint(w, r, id, sprintID)
		return
	}
	// Se muestra el texto del error de dominio, no el envuelto por el caso de uso.
	for _, noEncontrado := range []error{domain.ErrProyectoNoEncontrado, domain.ErrElementoNoEncontrado} {
		if errors.Is(err, noEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{noEncontrado.Error()})
			return
		}
	}
	if mensaje, ok := mensajeDeValidacion(err); ok {
		p.mostrarSprintDetalle(w, r, id, sprintID, http.StatusBadRequest, mensaje)
		return
	}
	p.fallaInterna(w, "completar un elemento del sprint", err)
}

// leerAccionSprint lee lo común a las acciones del detalle del sprint (asignar, quitar,
// completar): {id} y {sprintID} de la ruta y elemento_id del formulario. Si algo es inválido,
// ya respondió (página de error, o el detalle con el mensaje) y devuelve ok en false.
func (p *pantallas) leerAccionSprint(w http.ResponseWriter, r *http.Request) (id, sprintID, elementoID int64, ok bool) {
	id, ok = p.idDeRuta(w, r)
	if !ok {
		return 0, 0, 0, false
	}
	sprintID, err := sprintIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return 0, 0, 0, false
	}
	elementoID, err = elementoIDDeForm(r)
	if err != nil {
		p.mostrarSprintDetalle(w, r, id, sprintID, http.StatusBadRequest, err.Error())
		return 0, 0, 0, false
	}
	return id, sprintID, elementoID, true
}

// resultadoAsignacion redirige al detalle del sprint si la acción salió bien, o vuelve a
// mostrarlo con el error.
func (p *pantallas) resultadoAsignacion(w http.ResponseWriter, r *http.Request, id, sprintID int64, err error) {
	if err == nil {
		redirigirAlDetalleDelSprint(w, r, id, sprintID)
		return
	}
	if errors.Is(err, domain.ErrElementoNoEncontrado) || errors.Is(err, domain.ErrSprintNoEncontrado) {
		p.renderizar(w, http.StatusNotFound, "error.html", datosError{err.Error()})
		return
	}
	if mensaje, ok := mensajeDeValidacion(err); ok {
		p.mostrarSprintDetalle(w, r, id, sprintID, http.StatusBadRequest, mensaje)
		return
	}
	p.fallaInterna(w, "asignar o quitar un elemento del sprint", err)
}

// mostrarSprintDetalle busca el proyecto, el sprint (dentro de la lista de HU-08) y el backlog del
// proyecto, y renderiza la página separando lo ya asignado a este sprint de lo disponible.
func (p *pantallas) mostrarSprintDetalle(w http.ResponseWriter, r *http.Request, id, sprintID int64, status int, mensaje string) {
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	sprints, err := p.cu.ListarSprints.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	sprint, ok := buscarSprint(sprints, sprintID)
	if !ok {
		p.renderizar(w, http.StatusNotFound, "error.html", datosError{"sprint no encontrado"})
		return
	}
	elementos, err := p.cu.ListarBacklog.Ejecutar(r.Context(), app.DatosListarBacklog{ProyectoID: id})
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}

	abierto := sprint.Estado == domain.SprintAbierto
	var asignados, disponibles []vistaElementoBacklog
	for _, e := range elementos {
		v := vistaElementoBacklog{ID: e.ID, Titulo: e.Titulo, Estado: string(e.Estado),
			PuedeCompletarse: abierto && e.Estado != domain.EstadoHecho}
		switch {
		case e.SprintID != nil && *e.SprintID == sprintID:
			asignados = append(asignados, v)
		case e.SprintID == nil:
			disponibles = append(disponibles, v)
		}
	}

	p.renderizar(w, status, "sprint_detalle.html", datosSprintDetalle{
		Proyecto: aVistaProyecto(proyecto), Sprint: aVistaSprint(sprint),
		Asignados: asignados, Disponibles: disponibles, Error: mensaje,
	})
}

func buscarSprint(sprints []domain.Sprint, id int64) (domain.Sprint, bool) {
	for _, s := range sprints {
		if s.ID == id {
			return s, true
		}
	}
	return domain.Sprint{}, false
}

func aVistaSprint(s domain.Sprint) vistaSprint {
	return vistaSprint{
		ID: s.ID, Numero: s.Numero,
		FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
		Objetivo: s.Objetivo, Estado: string(s.Estado),
	}
}

// redirigirAlDetalleDelSprint usa 303 (Ver Otro) para que el navegador reconsulte con GET y no
// reenvíe el formulario si el usuario recarga la página, igual que redirigirAlDetalle.
func redirigirAlDetalleDelSprint(w http.ResponseWriter, r *http.Request, id, sprintID int64) {
	http.Redirect(w, r, rutaSprintDetalleWeb(id, sprintID), http.StatusSeeOther)
}

func rutaSprintDetalleWeb(id, sprintID int64) string {
	return fmt.Sprintf("/web/proyectos/%d/sprints/%d", id, sprintID)
}

// sprintIDDesdeRuta lee {sprintID} de la ruta, igual que proyectoIDDesdeRuta lee {id}.
func sprintIDDesdeRuta(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("sprintID"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("el id del sprint debe ser un número entero positivo")
	}
	return id, nil
}

// elementoIDDeForm lee el campo elemento_id del formulario.
func elementoIDDeForm(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PostFormValue("elemento_id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errFormulario{"elegí un elemento del backlog"}
	}
	return id, nil
}
