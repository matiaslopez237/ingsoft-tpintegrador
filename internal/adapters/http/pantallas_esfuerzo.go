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

// Datos de la sección "Registrar esfuerzo" de la pantalla del elemento (HU-18): solo texto ya formateado.
type (
	formEsfuerzo struct {
		Integrante, Fecha, Actividad, Horas string
	}
	datosSeccionEsfuerzo struct {
		Integrantes []domain.Integrante
		Form        formEsfuerzo
		Error       string
		Registrado  bool // se acaba de registrar un esfuerzo (confirmación)
	}
)

// registrarPantallasEsfuerzo agrega la acción de registrar esfuerzo desde la pantalla del elemento.
// Implementa HU-18 (CA-18.18 a CA-18.23).
func (p *pantallas) registrarPantallasEsfuerzo(mux *http.ServeMux) {
	mux.HandleFunc("POST /web/proyectos/{id}/backlog/{elementoID}/esfuerzo", p.formulario(p.registrarEsfuerzoPantalla))
}

// registrarEsfuerzoPantalla registra el esfuerzo con el mismo caso de uso que el endpoint JSON y
// vuelve al elemento (POST-redirect-GET), o muestra la pantalla con el error conservando lo tipeado.
func (p *pantallas) registrarEsfuerzoPantalla(w http.ResponseWriter, r *http.Request) {
	pid, eid, ok := p.idsElemento(w, r)
	if !ok {
		return
	}
	form := formEsfuerzo{
		Integrante: r.PostFormValue("integrante"),
		Fecha:      r.PostFormValue("fecha"),
		Actividad:  r.PostFormValue("actividad"),
		Horas:      r.PostFormValue("horas"),
	}
	datos, err := datosRegistroDeForm(pid, eid, form)
	if err == nil {
		_, err = p.cu.RegistrarEsfuerzo.Ejecutar(r.Context(), datos)
	}
	if err != nil {
		if mensaje, ok := mensajeDeValidacionEsfuerzo(err); ok {
			p.mostrarErrorDeEsfuerzo(w, r, pid, eid, form, mensaje)
			return
		}
		for _, noEncontrado := range erroresNoEncontradoEsfuerzo {
			if errors.Is(err, noEncontrado) {
				p.errorDeElemento(w, err)
				return
			}
		}
		p.fallaInterna(w, "registrar esfuerzo", err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/backlog/%d?esfuerzo=registrado#esfuerzo", pid, eid), http.StatusSeeOther)
}

// mostrarErrorDeEsfuerzo vuelve a la pantalla del elemento con el mensaje en la sección de esfuerzo.
// El formulario de edición se rearma con los valores guardados (cada formulario se envía por separado).
// El proyecto se busca primero (sección 7): un error de formato no pasó por el caso de uso, y
// ConsultarElementoBacklog no distingue un proyecto inexistente.
func (p *pantallas) mostrarErrorDeEsfuerzo(w http.ResponseWriter, r *http.Request, pid, eid int64, form formEsfuerzo, mensaje string) {
	if _, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), pid); err != nil {
		p.errorDeElemento(w, err)
		return
	}
	e, err := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, eid)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	f := aFormElemento(e)
	p.mostrarElemento(w, r, http.StatusBadRequest, pid, eid, f, f, "", datosSeccionEsfuerzo{Form: form, Error: mensaje})
}

// seccionEsfuerzoInicial son los valores con los que se abre el formulario: la fecha de hoy en la
// zona del servidor y lo demás vacío. Validar la fecha es tarea del dominio: acá solo se propone.
func (p *pantallas) seccionEsfuerzoInicial() datosSeccionEsfuerzo {
	return datosSeccionEsfuerzo{Form: formEsfuerzo{Fecha: p.reloj().Format(time.DateOnly)}}
}

// datosRegistroDeForm convierte el formulario. Los errores de formato no llegan al caso de uso,
// como el cuerpo mal formado del endpoint JSON (spec, sección 7).
func datosRegistroDeForm(pid, eid int64, f formEsfuerzo) (app.DatosRegistroEsfuerzo, error) {
	fecha, err := fechaDelPedido(f.Fecha)
	if err != nil {
		return app.DatosRegistroEsfuerzo{}, errFormulario{err.Error()}
	}
	horas, err := horasDeForm(f.Horas)
	if err != nil {
		return app.DatosRegistroEsfuerzo{}, err
	}
	return app.DatosRegistroEsfuerzo{
		ProyectoID: pid, ElementoID: eid, Integrante: f.Integrante,
		Fecha: fecha, Actividad: f.Actividad, Horas: horas,
	}, nil
}

// horasDeForm convierte las horas. Vacías valen 0: el dominio responde ErrHorasInvalidas, igual que
// cuando faltan en el JSON.
func horasDeForm(valor string) (float64, error) {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return 0, nil
	}
	horas, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		return 0, errFormulario{"las horas deben ser un número"}
	}
	return horas, nil
}

func mensajeDeValidacionEsfuerzo(err error) (string, bool) {
	return formMensajeDeErrorConocido(err, erroresDeValidacionEsfuerzo)
}
