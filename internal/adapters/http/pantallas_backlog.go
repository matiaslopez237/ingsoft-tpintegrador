package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Opciones de los <select> de la pantalla del backlog. Validarlas es tarea del dominio: acá solo
// se listan.
var (
	prioridadesElemento = []domain.Prioridad{domain.PrioridadAlta, domain.PrioridadMedia, domain.PrioridadBaja}
	estadosElemento     = []domain.EstadoElemento{domain.EstadoPendiente, domain.EstadoEnProgreso, domain.EstadoHecho}
)

// Datos de la pantalla del backlog (HU-07): solo texto ya formateado, sin lógica de negocio.
type (
	vistaElemento struct {
		Titulo, Prioridad, Estado string
		StoryPoints               int // 0: sin estimar
		Criterios                 []string
	}
	formularioElemento struct {
		Titulo, Descripcion, Prioridad, StoryPoints, Criterios string
	}
	datosBacklog struct {
		Proyecto      vistaProyecto
		Elementos     []vistaElemento
		Orden, Estado string // el orden y el filtro elegidos
		Form          formularioElemento
		Error         string
		Prioridades   []domain.Prioridad
		Estados       []domain.EstadoElemento
	}
)

// backlog muestra el backlog del proyecto con el orden y el filtro de la URL. Implementa HU-07.
func (p *pantallas) backlog(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	p.mostrarBacklog(w, r, id, http.StatusOK, formularioElemento{Prioridad: string(domain.PrioridadMedia)}, "")
}

// crearElemento agrega un elemento con el caso de uso de HU-05 y vuelve al backlog con el mismo
// orden y filtro, o lo muestra con el error y lo que se había escrito. Implementa HU-07 (sección 5).
func (p *pantallas) crearElemento(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	form := formularioElemento{
		Titulo:      r.PostFormValue("titulo"),
		Descripcion: r.PostFormValue("descripcion"),
		Prioridad:   r.PostFormValue("prioridad"),
		StoryPoints: r.PostFormValue("story_points"),
		Criterios:   r.PostFormValue("criterios"),
	}
	datos, err := datosNuevoElemento(id, form)
	if err == nil {
		_, err = p.cu.CrearElementoBacklog.Ejecutar(r.Context(), datos)
	}
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if mensaje, ok := mensajeDeValidacion(err); ok {
			p.mostrarBacklog(w, r, id, http.StatusBadRequest, form, mensaje)
			return
		}
		p.fallaInterna(w, "crear elemento del backlog", err)
		return
	}
	http.Redirect(w, r, rutaBacklog(id, r.URL.Query()), http.StatusSeeOther)
}

// rutaBacklog arma la ruta del backlog conservando el orden y el filtro elegidos, si los hay.
func rutaBacklog(id int64, q url.Values) string {
	conservar := url.Values{}
	for _, clave := range []string{"orden", "estado"} {
		if v := q.Get(clave); v != "" {
			conservar.Set(clave, v)
		}
	}
	ruta := fmt.Sprintf("/web/proyectos/%d/backlog", id)
	if len(conservar) > 0 {
		ruta += "?" + conservar.Encode()
	}
	return ruta
}

// mostrarBacklog renderiza la pantalla con el proyecto, su backlog ordenado y filtrado según la
// URL, el formulario y el error si lo hay. Un orden o estado inválidos muestran la página de error.
func (p *pantallas) mostrarBacklog(w http.ResponseWriter, r *http.Request, id int64, status int, form formularioElemento, mensaje string) {
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	q := r.URL.Query()
	orden, estado := q.Get("orden"), q.Get("estado")
	elementos, err := p.cu.ListarBacklog.Ejecutar(r.Context(), app.DatosListarBacklog{ProyectoID: id, Orden: orden, Estado: estado})
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if m, ok := mensajeDeValidacion(err); ok {
			p.renderizar(w, http.StatusBadRequest, "error.html", datosError{m})
			return
		}
		p.fallaInterna(w, "listar el backlog", err)
		return
	}

	vistas := make([]vistaElemento, len(elementos))
	for i, e := range elementos {
		vistas[i] = vistaElemento{Titulo: e.Titulo, Prioridad: string(e.Prioridad), Estado: string(e.Estado), Criterios: e.Criterios}
		if e.StoryPoints != nil {
			vistas[i].StoryPoints = *e.StoryPoints
		}
	}
	p.renderizar(w, status, "backlog.html", datosBacklog{
		Proyecto: aVistaProyecto(proyecto), Elementos: vistas, Orden: orden, Estado: estado,
		Form: form, Error: mensaje, Prioridades: prioridadesElemento, Estados: estadosElemento,
	})
}

// datosNuevoElemento convierte el formulario: story points vacíos quedan sin estimar, y los
// criterios van uno por línea, ignorando las líneas vacías (spec HU-07, sección 5).
func datosNuevoElemento(proyectoID int64, f formularioElemento) (app.DatosNuevoElemento, error) {
	d := app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      f.Titulo,
		Descripcion: f.Descripcion,
		Prioridad:   domain.Prioridad(f.Prioridad),
	}
	if sp := strings.TrimSpace(f.StoryPoints); sp != "" {
		n, err := strconv.Atoi(sp)
		if err != nil {
			return app.DatosNuevoElemento{}, errFormulario{"story_points debe ser un número entero"}
		}
		d.StoryPoints = &n
	}
	for _, linea := range strings.Split(f.Criterios, "\n") {
		if c := strings.TrimSpace(linea); c != "" {
			d.Criterios = append(d.Criterios, c)
		}
	}
	return d, nil
}
