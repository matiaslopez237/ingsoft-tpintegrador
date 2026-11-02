package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorElementos es lo que la pantalla del elemento necesita de ConsultarElementoBacklog.
type consultorElementos interface {
	Ejecutar(ctx context.Context, proyectoID, elementoID int64) (domain.ElementoBacklog, error)
}

// Las opciones de los <select> (prioridadesElemento y estadosElemento) son las de la pantalla
// del backlog (HU-07, pantallas_backlog.go).

// formEdicionElemento son los campos del formulario del elemento, como texto. Los criterios van uno por línea.
type formEdicionElemento struct {
	Titulo, Descripcion, Prioridad, StoryPoints, Criterios, Estado, HorasEstimadas string
	// HorasEnviadas indica que el pedido trae el campo horas_estimadas: ausente no cambia las horas.
	HorasEnviadas bool
}

type datosEdicionElemento struct {
	ProyectoID, ElementoID int64
	Form                   formEdicionElemento
	Original               formEdicionElemento // valores con los que se abrió la pantalla
	Error                  string
	Prioridades            []domain.Prioridad
	Estados                []domain.EstadoElemento
	Poker                  *datosPokerElemento  // nil si la aplicación no tiene Planning Poker (HU-14)
	Esfuerzo               datosSeccionEsfuerzo // HU-18: tiene su propio formulario y su propio error
	Comparacion            datosComparacion     // HU-19: siempre con los valores guardados
}

// registrarPantallasElemento agrega la pantalla para ver, editar y eliminar un elemento.
// Implementa HU-06 (CA-06.10).
func (p *pantallas) registrarPantallasElemento(mux *http.ServeMux) {
	mux.HandleFunc("GET /web/proyectos/{id}/backlog/{elementoID}", p.elemento)
	mux.HandleFunc("POST /web/proyectos/{id}/backlog/{elementoID}/editar", p.formulario(p.editarElemento))
	mux.HandleFunc("POST /web/proyectos/{id}/backlog/{elementoID}/eliminar", p.formulario(p.eliminarElemento))
}

func (p *pantallas) elemento(w http.ResponseWriter, r *http.Request) {
	pid, eid, ok := p.idsElemento(w, r)
	if !ok {
		return
	}
	e, err := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, eid)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	f := aFormElemento(e)
	esfuerzo := p.seccionEsfuerzoInicial()
	esfuerzo.Registrado = r.URL.Query().Get("esfuerzo") == "registrado" // HU-18: confirmación tras la redirección
	p.mostrarElemento(w, r, http.StatusOK, pid, eid, f, f, "", esfuerzo)
}

func (p *pantallas) editarElemento(w http.ResponseWriter, r *http.Request) {
	pid, eid, ok := p.idsElemento(w, r)
	if !ok {
		return
	}
	f, original, conOriginal := leerFormElemento(r)
	cambios, err := cambiosDeFormElemento(f, original, conOriginal)
	if err == nil && cambios == (domain.CambiosElemento{}) {
		// Si no cambió ningún campo, se vuelve al elemento sin mensaje: no es un error del
		// usuario (igual que el formulario de proyecto, DEF-01).
		http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid), http.StatusSeeOther)
		return
	}
	if err == nil {
		_, err = p.cu.EditarElementoBacklog.Ejecutar(r.Context(), app.DatosEditarElemento{ProyectoID: pid, ElementoID: eid, Cambios: cambios})
	}
	if err != nil {
		if errors.Is(err, domain.ErrElementoNoEncontrado) {
			p.errorDeElemento(w, err)
			return
		}
		if mensaje, ok := mensajeDeValidacionElemento(err); ok {
			p.mostrarElemento(w, r, http.StatusBadRequest, pid, eid, f, original, mensaje, p.seccionEsfuerzoInicial())
			return
		}
		p.fallaInterna(w, "editar elemento", err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid), http.StatusSeeOther)
}

func (p *pantallas) eliminarElemento(w http.ResponseWriter, r *http.Request) {
	pid, eid, ok := p.idsElemento(w, r)
	if !ok {
		return
	}
	err := p.cu.EliminarElementoBacklog.Ejecutar(r.Context(), app.DatosEliminarElemento{ProyectoID: pid, ElementoID: eid})
	switch {
	case err == nil:
		// Se vuelve a la lista del backlog (HU-07), desde donde se elimina.
		http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/backlog", pid), http.StatusSeeOther)
	case errors.Is(err, domain.ErrElementoConEsfuerzo), errors.Is(err, domain.ErrElementoConDefectos): // DEF-02
		e, errConsulta := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, eid)
		if errConsulta != nil {
			p.errorDeElemento(w, errConsulta)
			return
		}
		f := aFormElemento(e)
		p.mostrarElemento(w, r, http.StatusConflict, pid, eid, f, f, mensajeDeEliminacion(err), p.seccionEsfuerzoInicial())
	default:
		p.errorDeElemento(w, err)
	}
}

// --- Ayudantes ---

// mostrarElemento renderiza la pantalla. Consulta el proyecto por los integrantes del desplegable
// de esfuerzo (HU-18), también cuando vuelve con un error de editar o eliminar, la última ronda
// para la sección de Planning Poker (HU-14) y la comparación de esfuerzo (HU-19), que se lee de lo
// guardado y no de lo que se acaba de tipear.
func (p *pantallas) mostrarElemento(w http.ResponseWriter, r *http.Request, status int, pid, eid int64,
	f, original formEdicionElemento, mensaje string, esfuerzo datosSeccionEsfuerzo) {
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), pid)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	poker, err := p.pokerDelElemento(r, pid, eid, original.Estado)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	comparacion, err := p.cu.CompararEsfuerzo.Ejecutar(r.Context(), pid, eid)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	esfuerzo.Integrantes = proyecto.Integrantes
	p.renderizar(w, status, "elemento.html", datosEdicionElemento{
		ProyectoID: pid, ElementoID: eid, Form: f, Original: original, Error: mensaje,
		Prioridades: prioridadesElemento, Estados: estadosElemento, Poker: poker, Esfuerzo: esfuerzo,
		Comparacion: vistaComparacion(comparacion),
	})
}

func (p *pantallas) idsElemento(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	pid, eid, err := idsElementoDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return 0, 0, false
	}
	return pid, eid, true
}

// errorDeElemento muestra la página 404 si no existe el proyecto o el elemento en el proyecto, o la
// de falla interna.
func (p *pantallas) errorDeElemento(w http.ResponseWriter, err error) {
	for _, noEncontrado := range []error{domain.ErrProyectoNoEncontrado, domain.ErrElementoNoEncontrado} {
		if errors.Is(err, noEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{noEncontrado.Error()})
			return
		}
	}
	p.fallaInterna(w, "elemento del backlog", err)
}

func mensajeDeValidacionElemento(err error) (string, bool) {
	return formMensajeDeErrorConocido(err, erroresDeValidacionEdicion)
}

func aFormElemento(e domain.ElementoBacklog) formEdicionElemento {
	f := formEdicionElemento{
		Titulo: e.Titulo, Descripcion: e.Descripcion, Prioridad: string(e.Prioridad),
		Criterios: strings.Join(e.Criterios, "\n"), Estado: string(e.Estado),
	}
	if e.StoryPoints != nil {
		f.StoryPoints = strconv.Itoa(*e.StoryPoints)
	}
	if e.HorasEstimadas != nil {
		f.HorasEstimadas = strconv.FormatFloat(*e.HorasEstimadas, 'f', -1, 64) // con punto: lo exige el input numérico
	}
	return f
}

// leerFormElemento lee el formulario y los valores originales (campos ocultos original_*).
// conOriginal es false si el pedido no trae esos campos (no viene de la pantalla).
func leerFormElemento(r *http.Request) (f, original formEdicionElemento, conOriginal bool) {
	leer := func(prefijo string) formEdicionElemento {
		_, horasEnviadas := r.PostForm[prefijo+"horas_estimadas"]
		return formEdicionElemento{
			Titulo:         r.PostFormValue(prefijo + "titulo"),
			Descripcion:    r.PostFormValue(prefijo + "descripcion"),
			Prioridad:      r.PostFormValue(prefijo + "prioridad"),
			StoryPoints:    r.PostFormValue(prefijo + "story_points"),
			Criterios:      r.PostFormValue(prefijo + "criterios"),
			Estado:         r.PostFormValue(prefijo + "estado"),
			HorasEstimadas: r.PostFormValue(prefijo + "horas_estimadas"),
			HorasEnviadas:  horasEnviadas,
		}
	}
	_, conOriginal = r.PostForm["original_titulo"]
	return leer(""), leer("original_"), conOriginal
}

// cambiosDeFormElemento arma la edición con los campos que cambiaron respecto de los originales,
// para no pisar ediciones simultáneas (HU-06, sección 5; lección de DEF-01). Si el pedido no trae
// los originales (conOriginal en false), se envían todos los campos, como cambiosDeForm en DEF-01.
// Story points u horas estimadas vacíos quitan la estimación; los criterios van uno por línea y se
// ignoran las vacías. Las horas estimadas son la excepción: si el campo no vino en el pedido, no se
// envían aunque falten los originales, para que un POST incompleto no las borre.
func cambiosDeFormElemento(f, original formEdicionElemento, conOriginal bool) (domain.CambiosElemento, error) {
	cambio := func(nuevo, viejo string) bool { return !conOriginal || nuevo != viejo }
	var c domain.CambiosElemento
	// Título y descripción se comparan recortados: el dominio los recorta igual al guardar.
	if cambio(strings.TrimSpace(f.Titulo), strings.TrimSpace(original.Titulo)) {
		c.Titulo = &f.Titulo
	}
	if cambio(strings.TrimSpace(f.Descripcion), strings.TrimSpace(original.Descripcion)) {
		c.Descripcion = &f.Descripcion
	}
	if cambio(f.Prioridad, original.Prioridad) {
		pr := domain.Prioridad(f.Prioridad)
		c.Prioridad = &pr
	}
	if cambio(strings.TrimSpace(f.StoryPoints), strings.TrimSpace(original.StoryPoints)) {
		c.CambiarStoryPoints = true
		if s := strings.TrimSpace(f.StoryPoints); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return domain.CambiosElemento{}, errFormulario{"story_points debe ser un número entero"}
			}
			c.StoryPoints = &n
		}
	}
	if f.HorasEnviadas && cambio(strings.TrimSpace(f.HorasEstimadas), strings.TrimSpace(original.HorasEstimadas)) {
		c.CambiarHorasEstimadas = true
		if s := strings.TrimSpace(f.HorasEstimadas); s != "" {
			horas, err := strconv.ParseFloat(s, 64) // NaN e Inf pasan: los rechaza el dominio (RN-19.2)
			if err != nil {
				return domain.CambiosElemento{}, errFormulario{"las horas estimadas deben ser un número"}
			}
			c.HorasEstimadas = &horas
		}
	}
	// Los criterios se comparan como lista: el navegador envía el textarea con saltos CRLF.
	if criterios := lineas(f.Criterios); !conOriginal || !slices.Equal(criterios, lineas(original.Criterios)) {
		c.Criterios = &criterios
	}
	if cambio(f.Estado, original.Estado) {
		e := domain.EstadoElemento(f.Estado)
		c.Estado = &e
	}
	return c, nil
}

// lineas separa un texto en líneas recortadas, sin las vacías.
func lineas(texto string) []string {
	resultado := []string{}
	for _, l := range strings.Split(strings.ReplaceAll(texto, "\r\n", "\n"), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			resultado = append(resultado, l)
		}
	}
	return resultado
}

// mensajeDeEliminacion devuelve el mensaje del dominio del error que impidió eliminar el elemento,
// sin el texto con que lo envuelve el caso de uso.
func mensajeDeEliminacion(err error) string {
	if errors.Is(err, domain.ErrElementoConDefectos) {
		return domain.ErrElementoConDefectos.Error()
	}
	return domain.ErrElementoConEsfuerzo.Error()
}
