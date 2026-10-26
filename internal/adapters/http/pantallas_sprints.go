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
		Numero                int
		FechaInicio, FechaFin string
		Objetivo, Estado      string
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
		vistas[i] = vistaSprint{
			Numero:      s.Numero,
			FechaInicio: s.FechaInicio.Format(time.DateOnly),
			FechaFin:    s.FechaFin.Format(time.DateOnly),
			Objetivo:    s.Objetivo,
			Estado:      string(s.Estado),
		}
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
