package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// datosPokerElemento es la sección "Planning Poker" de la pantalla del elemento (HU-14): un enlace
// a la ronda abierta, si tiene, o el botón para iniciar una. Historial son sus rondas, de la más
// reciente a la más antigua (HU-16).
type datosPokerElemento struct {
	RondaAbierta  *domain.ResumenRonda
	PuedeIniciar  bool
	ElementoHecho bool
	Historial     []domain.ResumenRonda
}

// datosRonda es la pantalla de una ronda: quién votó y quién falta, nunca los valores (RN-14.3).
type datosRonda struct {
	ProyectoID int64
	Titulo     string // del elemento
	Ronda      domain.ResumenRonda
	Abierta    bool
	// PuedeRepetirse decide si se ofrece "Nueva ronda" (HU-16, RN-16.4): la regla es del dominio.
	PuedeRepetirse bool
	Valores        []int // la escala, para el formulario (RN-14.4)
	Integrante     string
	Error          string
}

// erroresDeVotoPantalla son los errores de HU-14 que la pantalla de la ronda muestra en la misma
// página con 400.
var erroresDeVotoPantalla = []error{domain.ErrRondaCerrada, domain.ErrNoEsVotante, domain.ErrVotoInvalido}

// erroresDeInicioPantalla son los errores de HU-14 que la pantalla del elemento muestra con 400.
var erroresDeInicioPantalla = []error{domain.ErrElementoTerminado, domain.ErrRondaSinVotantes, domain.ErrRondaYaAbierta}

// registrarPantallasPlanningPoker agrega las pantallas para iniciar una ronda y votar.
// Implementa HU-14 (CA-14.9).
func (p *pantallas) registrarPantallasPlanningPoker(mux *http.ServeMux) {
	mux.HandleFunc("POST /web/proyectos/{id}/backlog/{elementoID}/rondas", p.formulario(p.iniciarRondaPantalla))
	mux.HandleFunc("GET /web/proyectos/{id}/rondas/{rondaID}", p.rondaPantalla)
	mux.HandleFunc("POST /web/proyectos/{id}/rondas/{rondaID}/votar", p.formulario(p.votarPantalla))
}

// pokerDelElemento arma la sección "Planning Poker" de la pantalla del elemento. Devuelve nil si
// la aplicación no tiene los casos de uso de HU-14 (por ejemplo, en los tests de otras historias).
func (p *pantallas) pokerDelElemento(r *http.Request, pid, eid int64, estado string) (*datosPokerElemento, error) {
	if p.cu.ConsultarRonda == nil {
		return nil, nil
	}
	ultima, err := p.cu.ConsultarRonda.UltimaDelElemento(r.Context(), pid, eid)
	if err != nil {
		return nil, err
	}
	d := &datosPokerElemento{ElementoHecho: estado == string(domain.EstadoHecho)}
	if ultima != nil && ultima.Estado == domain.RondaAbierta {
		d.RondaAbierta = ultima
	}
	d.PuedeIniciar = d.RondaAbierta == nil && !d.ElementoHecho
	if p.cu.ListarRondas != nil { // HU-16; nil en los tests de otras historias
		if d.Historial, err = p.cu.ListarRondas.Ejecutar(r.Context(), pid, eid); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// iniciarRondaPantalla abre una ronda sobre el elemento y va a su pantalla; si no se puede, vuelve
// a mostrar el elemento con el error. Implementa HU-14 (CA-14.9).
func (p *pantallas) iniciarRondaPantalla(w http.ResponseWriter, r *http.Request) {
	pid, eid, ok := p.idsElemento(w, r)
	if !ok {
		return
	}
	res, err := p.cu.IniciarRonda.Ejecutar(r.Context(), app.DatosIniciarRonda{ProyectoID: pid, ElementoID: eid})
	if err == nil {
		http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/rondas/%d", pid, res.ID), http.StatusSeeOther)
		return
	}
	if mensaje, ok := mensajeDeErrorConocido(err, erroresDeInicioPantalla); ok {
		e, errConsulta := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, eid)
		if errConsulta != nil {
			p.errorDeElemento(w, errConsulta)
			return
		}
		f := aFormElemento(e)
		p.mostrarElemento(w, r, http.StatusBadRequest, pid, eid, f, f, mensaje, p.seccionEsfuerzoInicial())
		return
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		p.errorDeConsulta(w, err)
		return
	}
	p.errorDeElemento(w, err)
}

func (p *pantallas) rondaPantalla(w http.ResponseWriter, r *http.Request) {
	pid, rid, ok := p.idsRonda(w, r)
	if !ok {
		return
	}
	p.mostrarRonda(w, r, http.StatusOK, pid, rid, "", "")
}

// votarPantalla registra el voto y vuelve a la ronda; con un error, la muestra con el mensaje. El
// valor elegido no se vuelve a mostrar: la pantalla puede estar a la vista de otros (RN-14.3).
// Implementa HU-14 (CA-14.9).
func (p *pantallas) votarPantalla(w http.ResponseWriter, r *http.Request) {
	pid, rid, ok := p.idsRonda(w, r)
	if !ok {
		return
	}
	integrante := r.PostFormValue("integrante")
	valor, err := strconv.Atoi(r.PostFormValue("valor"))
	if err != nil {
		p.mostrarRonda(w, r, http.StatusBadRequest, pid, rid, integrante, "elegí un valor de la escala")
		return
	}
	_, err = p.cu.Votar.Ejecutar(r.Context(), app.DatosVoto{ProyectoID: pid, RondaID: rid, Integrante: integrante, Valor: valor})
	if err == nil {
		http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d/rondas/%d", pid, rid), http.StatusSeeOther)
		return
	}
	if mensaje, ok := mensajeDeErrorConocido(err, erroresDeVotoPantalla); ok {
		p.mostrarRonda(w, r, http.StatusBadRequest, pid, rid, integrante, mensaje)
		return
	}
	p.errorDeRonda(w, err)
}

// mostrarRonda renderiza la pantalla de la ronda con el resumen y el título del elemento.
func (p *pantallas) mostrarRonda(w http.ResponseWriter, r *http.Request, status int, pid, rid int64, integrante, mensaje string) {
	res, err := p.cu.ConsultarRonda.Ejecutar(r.Context(), pid, rid)
	if err != nil {
		p.errorDeRonda(w, err)
		return
	}
	e, err := p.cu.ConsultarElementoBacklog.Ejecutar(r.Context(), pid, res.ElementoID)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	ultima, err := p.cu.ConsultarRonda.UltimaDelElemento(r.Context(), pid, res.ElementoID)
	if err != nil {
		p.errorDeElemento(w, err)
		return
	}
	p.renderizar(w, status, "ronda.html", datosRonda{
		ProyectoID: pid, Titulo: e.Titulo, Ronda: res, Abierta: res.Estado == domain.RondaAbierta,
		PuedeRepetirse: domain.PuedeRepetirse(res, ultima, e),
		Valores:        domain.ValoresVoto(), Integrante: integrante, Error: mensaje,
	})
}

func (p *pantallas) idsRonda(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	pid, rid, err := idsRondaDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return 0, 0, false
	}
	return pid, rid, true
}

// errorDeRonda muestra la página 404 si el proyecto o la ronda no existen, o la de falla interna.
func (p *pantallas) errorDeRonda(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrRondaNoEncontrada):
		p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrRondaNoEncontrada.Error()})
	case errors.Is(err, domain.ErrProyectoNoEncontrado):
		p.errorDeConsulta(w, err)
	default:
		p.fallaInterna(w, "ronda de Planning Poker", err)
	}
}

// mensajeDeErrorConocido devuelve el mensaje del primer error de la lista que envuelve err.
func mensajeDeErrorConocido(err error, conocidos []error) (string, bool) {
	for _, e := range conocidos {
		if errors.Is(err, e) {
			return e.Error(), true
		}
	}
	return "", false
}
