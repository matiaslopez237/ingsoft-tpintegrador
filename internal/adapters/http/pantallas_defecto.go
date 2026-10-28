package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// severidadesDefecto son las opciones del <select> de severidad (RN-20.3). Validarlas es tarea
// del dominio: acá solo se listan.
var severidadesDefecto = []domain.Severidad{domain.SeveridadCritica, domain.SeveridadAlta, domain.SeveridadMedia, domain.SeveridadBaja}

// erroresDeValidacionDefectoPantalla son los errores de HU-20 que la pantalla muestra en el mismo
// formulario con HTTP 400. Incluye ErrElementoNoEncontrado: el elemento se elige en un <select>,
// pero un pedido manual podría traer un ID que ya no existe.
var erroresDeValidacionDefectoPantalla = []error{
	domain.ErrDescripcionDefectoVacia,
	domain.ErrSeveridadInvalida,
	domain.ErrElementoNoEncontrado,
}

// Datos de la pantalla de registrar un defecto (HU-20): solo texto ya formateado.
type (
	opcionElementoDefecto struct {
		ID     int64
		Titulo string
	}
	formDefecto struct {
		ElementoID, Descripcion, Severidad string
	}
	datosFormDefecto struct {
		Proyecto    vistaProyecto
		Elementos   []opcionElementoDefecto
		Form        formDefecto
		Error       string
		Severidades []domain.Severidad
		// Creado y SprintID solo se usan cuando se acaba de registrar el defecto (confirmación).
		Creado   bool
		SprintID *int64
	}
)

// registrarPantallasDefecto agrega la pantalla para registrar un defecto. Implementa HU-20 (CA-20.8).
func (p *pantallas) registrarPantallasDefecto(mux *http.ServeMux) {
	mux.HandleFunc("GET /web/proyectos/{id}/defectos/nuevo", p.formNuevoDefecto)
	mux.HandleFunc("POST /web/proyectos/{id}/defectos/nuevo", p.formulario(p.crearDefecto))
}

func (p *pantallas) formNuevoDefecto(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	p.mostrarFormDefecto(w, r, id, http.StatusOK, formDefecto{}, "", false, nil)
}

// crearDefecto registra el defecto con el caso de uso de HU-20 y muestra la confirmación con el
// Sprint de detección asignado (o que no se asignó ninguno), o vuelve a mostrar el formulario con
// el error si los datos no son válidos. Implementa HU-20 (CA-20.1 a CA-20.5, CA-20.8).
func (p *pantallas) crearDefecto(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	form := formDefecto{
		ElementoID:  r.PostFormValue("elemento_id"),
		Descripcion: r.PostFormValue("descripcion"),
		Severidad:   r.PostFormValue("severidad"),
	}
	datos, err := datosNuevoDefecto(id, form)
	var def domain.Defecto
	if err == nil {
		def, err = p.cu.RegistrarDefecto.Ejecutar(r.Context(), datos)
	}
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if mensaje, ok := mensajeDeValidacionDefecto(err); ok {
			p.mostrarFormDefecto(w, r, id, http.StatusBadRequest, form, mensaje, false, nil)
			return
		}
		p.fallaInterna(w, "registrar defecto", err)
		return
	}
	p.mostrarFormDefecto(w, r, id, http.StatusOK, formDefecto{}, "", true, def.SprintDeteccionID)
}

// datosNuevoDefecto convierte el formulario: el elemento tiene que ser un entero (viene de un
// <select>, pero un pedido manual podría traer cualquier texto).
func datosNuevoDefecto(proyectoID int64, f formDefecto) (app.DatosDefecto, error) {
	elementoID, err := strconv.ParseInt(strings.TrimSpace(f.ElementoID), 10, 64)
	if err != nil {
		return app.DatosDefecto{}, errFormulario{"elegí un elemento del backlog"}
	}
	return app.DatosDefecto{
		ProyectoID: proyectoID, ElementoID: elementoID,
		Descripcion: f.Descripcion, Severidad: domain.Severidad(f.Severidad),
	}, nil
}

// mostrarFormDefecto renderiza el formulario (o la confirmación, si creado es true) con la lista
// de elementos del backlog del proyecto para el <select>.
func (p *pantallas) mostrarFormDefecto(w http.ResponseWriter, r *http.Request, id int64, status int,
	form formDefecto, mensaje string, creado bool, sprintID *int64) {
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	elementos, err := p.cu.ListarBacklog.Ejecutar(r.Context(), app.DatosListarBacklog{ProyectoID: id})
	if err != nil {
		p.fallaInterna(w, "listar el backlog", err)
		return
	}
	opciones := make([]opcionElementoDefecto, len(elementos))
	for i, e := range elementos {
		opciones[i] = opcionElementoDefecto{ID: e.ID, Titulo: e.Titulo}
	}
	p.renderizar(w, status, "defecto_form.html", datosFormDefecto{
		Proyecto: aVistaProyecto(proyecto), Elementos: opciones, Form: form, Error: mensaje,
		Severidades: severidadesDefecto, Creado: creado, SprintID: sprintID,
	})
}

func mensajeDeValidacionDefecto(err error) (string, bool) {
	var ef errFormulario
	if errors.As(err, &ef) {
		return ef.mensaje, true
	}
	for _, e := range erroresDeValidacionDefectoPantalla {
		if errors.Is(err, e) {
			return e.Error(), true
		}
	}
	return "", false
}
