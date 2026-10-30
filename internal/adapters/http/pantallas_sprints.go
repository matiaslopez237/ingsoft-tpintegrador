package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Datos de la pantalla de sprints (HU-08): solo texto ya formateado, sin lógica de negocio.
type (
	vistaSprint struct {
		ID                    int64
		Numero                int
		FechaInicio, FechaFin string
		Objetivo, Estado      string
		PuedeCerrarse         bool // muestra el botón "Cerrar sprint" (HU-11, CA-11.5)
	}
	formularioSprint struct {
		FechaInicio, FechaFin, Objetivo string
	}
	datosSprints struct {
		Proyecto vistaProyecto
		Sprints  []vistaSprint
		Form     formularioSprint
		Error    string
	}
)

// sprints muestra los sprints del proyecto y el formulario para crear uno. Implementa HU-08 (CA-08.9).
func (p *pantallas) sprints(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	p.mostrarSprints(w, r, id, http.StatusOK, formularioSprint{}, "")
}

// crearSprint crea el sprint con el caso de uso de HU-08 y vuelve a la pantalla de sprints, o la
// muestra con el error y lo que se había escrito. Implementa HU-08 (CA-08.1 a CA-08.8).
func (p *pantallas) crearSprint(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	form := formularioSprint{
		FechaInicio: r.PostFormValue("fecha_inicio"),
		FechaFin:    r.PostFormValue("fecha_fin"),
		Objetivo:    r.PostFormValue("objetivo"),
	}
	datos, err := datosNuevoSprint(id, form)
	if err == nil {
		_, err = p.cu.CrearSprint.Ejecutar(r.Context(), datos)
	}
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if mensaje, ok := mensajeDeValidacion(err); ok {
			p.mostrarSprints(w, r, id, http.StatusBadRequest, form, mensaje)
			return
		}
		p.fallaInterna(w, "crear sprint", err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/sprints", id), http.StatusSeeOther)
}

// cerrarSprintPantalla cierra el sprint con el caso de uso de HU-11 y vuelve a la pantalla de
// sprints. Si ya estaba cerrado, la muestra con el error; si el proyecto o el sprint no existen,
// muestra la página 404. Implementa HU-11 (CA-11.5).
func (p *pantallas) cerrarSprintPantalla(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	sprintID, err := sprintIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return
	}

	err = p.cu.CerrarSprint.Ejecutar(r.Context(), app.DatosCerrarSprint{ProyectoID: id, SprintID: sprintID})
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if errors.Is(err, domain.ErrSprintNoEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrSprintNoEncontrado.Error()})
			return
		}
		if mensaje, ok := mensajeDeValidacion(err); ok {
			p.mostrarSprints(w, r, id, http.StatusBadRequest, formularioSprint{}, mensaje)
			return
		}
		p.fallaInterna(w, "cerrar sprint", err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/sprints", id), http.StatusSeeOther)
}

// mostrarSprints renderiza la pantalla con el proyecto y sus sprints, el formulario y el error si lo hay.
func (p *pantallas) mostrarSprints(w http.ResponseWriter, r *http.Request, id int64, status int, form formularioSprint, mensaje string) {
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
	vistas := make([]vistaSprint, len(sprints))
	for i, s := range sprints {
		vistas[i] = aVistaSprint(s)
	}
	p.renderizar(w, status, "sprints.html", datosSprints{
		Proyecto: aVistaProyecto(proyecto), Sprints: vistas, Form: form, Error: mensaje,
	})
}

// datosNuevoSprint convierte el formulario. Una fecha vacía queda en cero (el dominio la rechaza);
// una con formato inválido se informa como error del formulario.
func datosNuevoSprint(proyectoID int64, f formularioSprint) (app.DatosNuevoSprint, error) {
	d := app.DatosNuevoSprint{ProyectoID: proyectoID, Objetivo: f.Objetivo}
	inicio, err := fechaDeForm(f.FechaInicio, "fecha_inicio")
	if err != nil {
		return app.DatosNuevoSprint{}, err
	}
	fin, err := fechaDeForm(f.FechaFin, "fecha_fin")
	if err != nil {
		return app.DatosNuevoSprint{}, err
	}
	if inicio != nil {
		d.FechaInicio = *inicio
	}
	if fin != nil {
		d.FechaFin = *fin
	}
	return d, nil
}

// vistaSprintAnterior es cada fila de la pantalla de sprints anteriores (HU-12).
type vistaSprintAnterior struct {
	ID                              int64
	Numero                          int
	FechaInicio, FechaFin, Objetivo string
	Planificados, Completados       int
}

type datosSprintsAnteriores struct {
	Proyecto vistaProyecto
	Sprints  []vistaSprintAnterior
}

// sprintsAnteriores muestra el historial de sprints cerrados del proyecto, en el orden que da el
// caso de uso. Implementa HU-12 (CA-12.7).
func (p *pantallas) sprintsAnteriores(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	sprints, err := p.cu.ListarSprintsAnteriores.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	vistas := make([]vistaSprintAnterior, len(sprints))
	for i, s := range sprints {
		vistas[i] = vistaSprintAnterior{
			ID: s.ID, Numero: s.Numero,
			FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
			Objetivo: s.Objetivo, Planificados: s.PuntosPlanificados, Completados: s.PuntosCompletados,
		}
	}
	p.renderizar(w, http.StatusOK, "sprints_anteriores.html", datosSprintsAnteriores{Proyecto: aVistaProyecto(proyecto), Sprints: vistas})
}
