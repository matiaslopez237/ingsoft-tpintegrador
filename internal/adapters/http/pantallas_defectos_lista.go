package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresEstadoDefectoPantalla son los errores de HU-21 que la pantalla muestra en la misma lista
// con HTTP 400 (los 404 de proyecto/defecto van a la página de error, como las demás historias).
var erroresEstadoDefectoPantalla = []error{
	domain.ErrDefectoNoAbierto,
	domain.ErrDefectoNoEnProgreso,
	domain.ErrDefectoNoResuelto,
}

// accionEstadoDefecto es la firma común de ActualizarEstadoDefecto.Iniciar, .Resolver y .Cerrar.
// Se forma recién dentro de iniciarDefectoPantalla/resolverDefectoPantalla/cerrarDefectoPantalla
// (por pedido), nunca al registrar las rutas: ver el comentario de RegistrarActualizarEstadoDefecto.
type accionEstadoDefecto func(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)

// Datos de la pantalla de defectos (HU-21): solo texto ya formateado, sin lógica de negocio.
type (
	vistaDefecto struct {
		ID                    int64
		Elemento, Descripcion string
		Severidad, Estado     string
		SprintDeteccion       *int64
		SprintResolucion      *int64
		PuedeIniciarse        bool
		PuedeResolverse       bool
		PuedeCerrarse         bool
	}
	datosDefectos struct {
		Proyecto vistaProyecto
		Defectos []vistaDefecto
		Error    string
	}
)

// registrarPantallasDefectosLista agrega las pantallas de HU-21: listar los defectos de un
// proyecto y actualizar el estado de uno. Implementa HU-21 (CA-21.11).
func (p *pantallas) registrarPantallasDefectosLista(mux *http.ServeMux) {
	mux.HandleFunc("GET /web/proyectos/{id}/defectos", p.defectos)
	mux.HandleFunc("POST /web/proyectos/{id}/defectos/{defectoID}/iniciar", p.formulario(p.iniciarDefectoPantalla))
	mux.HandleFunc("POST /web/proyectos/{id}/defectos/{defectoID}/resolver", p.formulario(p.resolverDefectoPantalla))
	mux.HandleFunc("POST /web/proyectos/{id}/defectos/{defectoID}/cerrar", p.formulario(p.cerrarDefectoPantalla))
}

func (p *pantallas) defectos(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	p.mostrarDefectos(w, r, id, http.StatusOK, "")
}

// iniciarDefectoPantalla inicia el defecto con el caso de uso de HU-21 y vuelve a la lista. Si no
// estaba Abierto, la muestra con el error; si el proyecto o el defecto no existen, la página 404.
func (p *pantallas) iniciarDefectoPantalla(w http.ResponseWriter, r *http.Request) {
	p.accionEstadoDefectoPantalla(w, r, p.cu.ActualizarEstadoDefecto.Iniciar, "iniciar defecto")
}

// resolverDefectoPantalla resuelve el defecto con el caso de uso de HU-21 y vuelve a la lista. Si
// no estaba En progreso, la muestra con el error; si el proyecto o el defecto no existen, la
// página 404.
func (p *pantallas) resolverDefectoPantalla(w http.ResponseWriter, r *http.Request) {
	p.accionEstadoDefectoPantalla(w, r, p.cu.ActualizarEstadoDefecto.Resolver, "resolver defecto")
}

// cerrarDefectoPantalla cierra el defecto con el caso de uso de HU-21 y vuelve a la lista. Si no
// estaba Resuelto, la muestra con el error; si el proyecto o el defecto no existen, la página 404.
func (p *pantallas) cerrarDefectoPantalla(w http.ResponseWriter, r *http.Request) {
	p.accionEstadoDefectoPantalla(w, r, p.cu.ActualizarEstadoDefecto.Cerrar, "cerrar defecto")
}

func (p *pantallas) accionEstadoDefectoPantalla(w http.ResponseWriter, r *http.Request, accion accionEstadoDefecto, nombre string) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	defectoID, err := defectoIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return
	}

	_, err = accion(r.Context(), app.DatosActualizarEstadoDefecto{ProyectoID: id, DefectoID: defectoID})
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if errors.Is(err, domain.ErrDefectoNoEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrDefectoNoEncontrado.Error()})
			return
		}
		for _, e := range erroresEstadoDefectoPantalla {
			if errors.Is(err, e) {
				p.mostrarDefectos(w, r, id, http.StatusBadRequest, e.Error())
				return
			}
		}
		p.fallaInterna(w, nombre, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/defectos", id), http.StatusSeeOther)
}

// mostrarDefectos renderiza la pantalla con el proyecto y sus defectos, enriquecidos con el título
// del elemento relacionado (el dominio solo guarda el ElementoID), y el error si lo hay.
func (p *pantallas) mostrarDefectos(w http.ResponseWriter, r *http.Request, id int64, status int, mensaje string) {
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	defectos, err := p.cu.ListarDefectos.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	elementos, err := p.cu.ListarBacklog.Ejecutar(r.Context(), app.DatosListarBacklog{ProyectoID: id})
	if err != nil {
		p.fallaInterna(w, "listar el backlog", err)
		return
	}
	titulos := make(map[int64]string, len(elementos))
	for _, e := range elementos {
		titulos[e.ID] = e.Titulo
	}

	vistas := make([]vistaDefecto, len(defectos))
	for i, d := range defectos {
		vistas[i] = aVistaDefecto(d, titulos[d.ElementoID])
	}
	p.renderizar(w, status, "defectos.html", datosDefectos{Proyecto: aVistaProyecto(proyecto), Defectos: vistas, Error: mensaje})
}

func aVistaDefecto(d domain.Defecto, elemento string) vistaDefecto {
	return vistaDefecto{
		ID: d.ID, Elemento: elemento, Descripcion: d.Descripcion,
		Severidad: string(d.Severidad), Estado: string(d.Estado),
		SprintDeteccion: d.SprintDeteccionID, SprintResolucion: d.SprintResolucionID,
		PuedeIniciarse: d.PuedeIniciarse(), PuedeResolverse: d.PuedeResolverse(), PuedeCerrarse: d.PuedeCerrarse(),
	}
}
