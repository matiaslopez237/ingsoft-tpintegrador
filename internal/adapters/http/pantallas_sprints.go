package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	HayPorcentaje                   bool   // false: el sprint se cerró sin puntos planificados (HU-24, RN-24.2)
	PorcentajeCompletado            string // ya formateado, solo si HayPorcentaje
	DefectosDetectados              int
	DefectosResueltos               int
}

type datosSprintsAnteriores struct {
	Proyecto     vistaProyecto
	Sprints      []vistaSprintAnterior
	HayVelocidad bool   // false: todavía no hay Sprints cerrados (HU-22, RN-22.3)
	Velocidad    string // ya formateada con coma decimal, solo si HayVelocidad
}

// sprintsAnteriores muestra el historial de sprints cerrados del proyecto, en el orden que da el
// caso de uso, con la velocidad del equipo y la calidad de cada sprint. Implementa HU-12 (CA-12.7),
// HU-22 (CA-22.8) y HU-24 (CA-24.9).
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
	velocidad, err := p.cu.ConsultarVelocidad.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	calidad, err := p.cu.ConsultarCalidadSprints.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	calidadPorSprint := make(map[int64]app.CalidadSprint, len(calidad))
	for _, c := range calidad {
		calidadPorSprint[c.Sprint.ID] = c
	}
	vistas := make([]vistaSprintAnterior, len(sprints))
	for i, s := range sprints {
		v := vistaSprintAnterior{
			ID: s.ID, Numero: s.Numero,
			FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
			Objetivo: s.Objetivo, Planificados: s.PuntosPlanificados, Completados: s.PuntosCompletados,
		}
		if c, ok := calidadPorSprint[s.ID]; ok {
			v.DefectosDetectados, v.DefectosResueltos = c.DefectosDetectados, c.DefectosResueltos
			if c.HayPorcentaje {
				v.HayPorcentaje = true
				v.PorcentajeCompletado = formatearPorcentaje(c.PorcentajeCompletado)
			}
		}
		vistas[i] = v
	}
	datos := datosSprintsAnteriores{Proyecto: aVistaProyecto(proyecto), Sprints: vistas}
	if velocidad.HayDatos {
		datos.HayVelocidad = true
		datos.Velocidad = formatearVelocidad(velocidad.Velocidad)
	}
	p.renderizar(w, http.StatusOK, "sprints_anteriores.html", datos)
}

// formatearVelocidad muestra la velocidad con un decimal y coma, como el resto de los números de
// la interfaz (spec HU-22, P-2: es un detalle de la plantilla, no una regla de negocio).
func formatearVelocidad(v float64) string {
	return strings.Replace(strconv.FormatFloat(v, 'f', 1, 64), ".", ",", 1)
}

// formatearPorcentaje muestra el porcentaje completado con coma y un decimal (75,0%), igual
// criterio que textoPorcentaje de HU-19/HU-23: v ya viene redondeado a un decimal (RN-24.1,
// domain.PorcentajeCompletado), así que la API y la pantalla siempre coinciden.
func formatearPorcentaje(v float64) string {
	return strings.Replace(strconv.FormatFloat(v, 'f', 1, 64), ".", ",", 1) + "%"
}
