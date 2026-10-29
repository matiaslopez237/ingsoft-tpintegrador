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
	Titulo, Descripcion, Prioridad, StoryPoints, Criterios, Estado string
}

type datosEdicionElemento struct {
	ProyectoID, ElementoID int64
	Form                   formEdicionElemento
	Original               formEdicionElemento // valores con los que se abrió la pantalla
	Error                  string
	Prioridades            []domain.Prioridad
	Estados                []domain.EstadoElemento
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
	p.mostrarElemento(w, http.StatusOK, pid, eid, f, f, "")
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
			p.mostrarElemento(w, http.StatusBadRequest, pid, eid, f, original, mensaje)
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
	conflicto, esConflicto := conflictoAlEliminar(err)
	switch {
	case err == nil:
		// Se vuelve a la lista del backlog (HU-07), desde donde se elimina.
		http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/backlog", pid), http.StatusSeeOther)
	case esConflicto: // RN-06.7 o RN-06.8: se muestra el elemento con el motivo
		e, errConsulta := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, eid)
		if errConsulta != nil {
			p.errorDeElemento(w, errConsulta)
			return
		}
		f := aFormElemento(e)
		p.mostrarElemento(w, http.StatusConflict, pid, eid, f, f, conflicto.Error())
	default:
		p.errorDeElemento(w, err)
	}
}

// --- Ayudantes ---

func (p *pantallas) mostrarElemento(w http.ResponseWriter, status int, pid, eid int64, f, original formEdicionElemento, mensaje string) {
	p.renderizar(w, status, "elemento.html", datosEdicionElemento{
		ProyectoID: pid, ElementoID: eid, Form: f, Original: original, Error: mensaje,
		Prioridades: prioridadesElemento, Estados: estadosElemento,
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

// errorDeElemento muestra la página 404 si el elemento no existe en el proyecto, o la de falla interna.
func (p *pantallas) errorDeElemento(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrElementoNoEncontrado) {
		p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrElementoNoEncontrado.Error()})
		return
	}
	p.fallaInterna(w, "elemento del backlog", err)
}

func mensajeDeValidacionElemento(err error) (string, bool) {
	var ef errFormulario
	if errors.As(err, &ef) {
		return ef.mensaje, true
	}
	for _, e := range erroresDeValidacionEdicion {
		if errors.Is(err, e) {
			return e.Error(), true
		}
	}
	return "", false
}

func aFormElemento(e domain.ElementoBacklog) formEdicionElemento {
	f := formEdicionElemento{
		Titulo: e.Titulo, Descripcion: e.Descripcion, Prioridad: string(e.Prioridad),
		Criterios: strings.Join(e.Criterios, "\n"), Estado: string(e.Estado),
	}
	if e.StoryPoints != nil {
		f.StoryPoints = strconv.Itoa(*e.StoryPoints)
	}
	return f
}

// leerFormElemento lee el formulario y los valores originales (campos ocultos original_*).
// conOriginal es false si el pedido no trae esos campos (no viene de la pantalla).
func leerFormElemento(r *http.Request) (f, original formEdicionElemento, conOriginal bool) {
	leer := func(prefijo string) formEdicionElemento {
		return formEdicionElemento{
			Titulo:      r.PostFormValue(prefijo + "titulo"),
			Descripcion: r.PostFormValue(prefijo + "descripcion"),
			Prioridad:   r.PostFormValue(prefijo + "prioridad"),
			StoryPoints: r.PostFormValue(prefijo + "story_points"),
			Criterios:   r.PostFormValue(prefijo + "criterios"),
			Estado:      r.PostFormValue(prefijo + "estado"),
		}
	}
	_, conOriginal = r.PostForm["original_titulo"]
	return leer(""), leer("original_"), conOriginal
}

// cambiosDeFormElemento arma la edición con los campos que cambiaron respecto de los originales,
// para no pisar ediciones simultáneas (HU-06, sección 5; lección de DEF-01). Si el pedido no trae
// los originales (conOriginal en false), se envían todos los campos, como cambiosDeForm en DEF-01.
// Story points vacío quita la estimación; los criterios van uno por línea y se ignoran las vacías.
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
