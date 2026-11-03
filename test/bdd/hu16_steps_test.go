package bdd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU16 traduce la descripción del escenario al error de la spec (sección 7).
func erroresHU16() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"elemento no encontrado": domain.ErrElementoNoEncontrado,
	}
}

// historialWorld16 es el estado propio de los escenarios de HU-16, sobre el mundo de HU-14 (las
// rondas se arman y se votan con sus pasos).
type historialWorld16 struct {
	*planningPokerWorld
	historial []rondaVista
	err       error
}

// registrarStepsHU16 agrega los pasos de HU-16 sobre el mismo mundo de HU-14 (registrarStepsHU14
// llama a esta función, como a la de HU-15).
func registrarStepsHU16(ctx *godog.ScenarioContext, pw *planningPokerWorld) {
	w := &historialWorld16{planningPokerWorld: pw}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		w.historial, w.err = nil, nil
		return c, nil
	})

	// Given
	ctx.Step(`^"([^"]*)" ya tuvo una ronda cerrada con los votos "([^"]*)"$`, w.yaTuvoRondaCerrada)

	// When
	ctx.Step(`^consulto el historial de rondas de "([^"]*)"$`, w.consultoHistorial)
	ctx.Step(`^consulto el historial de rondas de "([^"]*)" desde el proyecto "([^"]*)"$`, w.consultoHistorialDesde)
	ctx.Step(`^consulto el historial de rondas de "([^"]*)" desde un proyecto inexistente$`, w.consultoHistorialDesdeInexistente)
	ctx.Step(`^desde la pantalla de la ronda (\d+) de "([^"]*)" uso "Nueva ronda"$`, w.usoNuevaRonda)

	// Then
	ctx.Step(`^el historial de rondas es "([^"]*)"$`, w.historialEs)
	ctx.Step(`^el historial de rondas de "([^"]*)" es "([^"]*)"$`, w.historialDeEs)
	ctx.Step(`^el historial de rondas está vacío$`, w.historialVacio)
	ctx.Step(`^en el historial, la ronda (\d+) revela los votos "([^"]*)" con mínimo (\d+) y máximo (\d+)$`, w.historialRevela)
	ctx.Step(`^en el historial de "([^"]*)", la ronda (\d+) revela los votos "([^"]*)" con mínimo (\d+) y máximo (\d+)$`, w.historialDeRevela)
	ctx.Step(`^en el historial, la ronda (\d+) no revela ningún voto$`, w.historialNoRevela)
	ctx.Step(`^se rechaza el historial con el error "([^"]*)"$`, w.seRechazaHistorial)
	ctx.Step(`^la pantalla de la ronda (\d+) de "([^"]*)" no ofrece "Nueva ronda"$`, w.noOfreceNuevaRonda)
	ctx.Step(`^la pantalla del elemento "([^"]*)" lista las rondas "([^"]*)"$`, w.elementoListaRondas)
}

// --- Given ---

// yaTuvoRondaCerrada abre una ronda sobre el elemento y la cierra con los votos indicados, que
// tienen que ser de todos los votantes.
func (w *historialWorld16) yaTuvoRondaCerrada(titulo, lista string) error {
	if err := w.hayRondaAbierta(titulo); err != nil {
		return err
	}
	votos, err := votosDe(lista)
	if err != nil {
		return err
	}
	for _, v := range votos {
		if err := w.yaVoto(v.Nombre, v.Valor); err != nil {
			return err
		}
	}
	if w.ronda.Estado != string(domain.RondaCerrada) {
		return fmt.Errorf("la ronda %d no se cerró con los votos %q: faltan votantes", w.ronda.Numero, lista)
	}
	return nil
}

// --- When ---

func (w *historialWorld16) consultar(pid int64, titulo string) error {
	_, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	w.historial, w.err = w.sistema.(sistemaHistorialRondas).HistorialRondas(pid, eid)
	if errors.Is(w.err, godog.ErrPending) {
		return godog.ErrPending
	}
	return nil
}

func (w *historialWorld16) consultoHistorial(titulo string) error {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	return w.consultar(pid, titulo)
}

func (w *historialWorld16) consultoHistorialDesde(titulo, proyecto string) error {
	pid, ok := w.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	return w.consultar(pid, titulo)
}

func (w *historialWorld16) consultoHistorialDesdeInexistente(titulo string) error {
	return w.consultar(idProyectoInexistente, titulo)
}

// usoNuevaRonda abre la pantalla de la ronda, verifica que ofrezca "Nueva ronda" y envía su formulario.
func (w *historialWorld16) usoNuevaRonda(numero int, titulo string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	r, err := w.pantallaDeRonda(numero, titulo)
	if err != nil {
		return err
	}
	accion := fmt.Sprintf("/web/proyectos/%d/backlog/%d/rondas", pid, eid)
	if !ofreceNuevaRonda(r.Cuerpo, accion) {
		return fmt.Errorf("la pantalla de la ronda %d no ofrece \"Nueva ronda\" (estado %d)", numero, r.Estado)
	}
	w.pantalla, err = w.sistema.Pedir("POST", accion, url.Values{})
	return err
}

// --- Then ---

func (w *historialWorld16) consultado() error {
	if w.err != nil {
		return fmt.Errorf("se esperaba consultar el historial, pero hubo error: %w", w.err)
	}
	return nil
}

// historialEs compara el historial con una lista como "3: Abierta, 2: Cerrada", en orden.
func (w *historialWorld16) historialEs(esperado string) error {
	if err := w.consultado(); err != nil {
		return err
	}
	var got []string
	for _, r := range w.historial {
		got = append(got, fmt.Sprintf("%d: %s", r.Numero, r.Estado))
	}
	if strings.Join(got, ", ") != esperado {
		return fmt.Errorf("historial = %q, se esperaba %q", strings.Join(got, ", "), esperado)
	}
	return nil
}

func (w *historialWorld16) historialDeEs(titulo, esperado string) error {
	if err := w.consultoHistorial(titulo); err != nil {
		return err
	}
	return w.historialEs(esperado)
}

func (w *historialWorld16) historialVacio() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if len(w.historial) != 0 {
		return fmt.Errorf("se esperaba un historial vacío, tiene %d rondas", len(w.historial))
	}
	return nil
}

func (w *historialWorld16) rondaDelHistorial(numero int) (rondaVista, error) {
	if err := w.consultado(); err != nil {
		return rondaVista{}, err
	}
	i := slices.IndexFunc(w.historial, func(r rondaVista) bool { return r.Numero == numero })
	if i < 0 {
		return rondaVista{}, fmt.Errorf("la ronda %d no está en el historial", numero)
	}
	return w.historial[i], nil
}

func (w *historialWorld16) historialRevela(numero int, lista string, minimo, maximo int) error {
	r, err := w.rondaDelHistorial(numero)
	if err != nil {
		return err
	}
	if r.Revelacion == nil {
		return fmt.Errorf("la ronda %d del historial no revela los votos (estado %q)", numero, r.Estado)
	}
	want, err := votosDe(lista)
	if err != nil {
		return err
	}
	if !slices.Equal(r.Revelacion.Votos, want) || r.Revelacion.Minimo != minimo || r.Revelacion.Maximo != maximo {
		return fmt.Errorf("ronda %d: votos %v, mínimo %d y máximo %d; se esperaba %v, %d y %d",
			numero, r.Revelacion.Votos, r.Revelacion.Minimo, r.Revelacion.Maximo, want, minimo, maximo)
	}
	return nil
}

func (w *historialWorld16) historialDeRevela(titulo string, numero int, lista string, minimo, maximo int) error {
	if err := w.consultoHistorial(titulo); err != nil {
		return err
	}
	return w.historialRevela(numero, lista, minimo, maximo)
}

func (w *historialWorld16) historialNoRevela(numero int) error {
	r, err := w.rondaDelHistorial(numero)
	if err != nil {
		return err
	}
	if r.Revelacion != nil {
		return fmt.Errorf("la ronda %d (%s) no debería revelar votos: %+v", numero, r.Estado, r.Revelacion)
	}
	return nil
}

func (w *historialWorld16) seRechazaHistorial(descripcion string) error {
	esperado, ok := erroresHU16()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.err == nil {
		return fmt.Errorf("se esperaba el error %q, pero la consulta funcionó", descripcion)
	}
	if !errors.Is(w.err, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.err)
	}
	return nil
}

func (w *historialWorld16) noOfreceNuevaRonda(numero int, titulo string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	r, err := w.pantallaDeRonda(numero, titulo)
	if err != nil {
		return err
	}
	if ofreceNuevaRonda(r.Cuerpo, fmt.Sprintf("/web/proyectos/%d/backlog/%d/rondas", pid, eid)) {
		return fmt.Errorf("la pantalla de la ronda %d no debería ofrecer \"Nueva ronda\"", numero)
	}
	return nil
}

// elementoListaRondas verifica que la pantalla del elemento enlace a las rondas en el orden dado,
// por ejemplo "2, 1".
func (w *historialWorld16) elementoListaRondas(titulo, numeros string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid), nil)
	if err != nil {
		return err
	}
	if err := w.consultar(pid, titulo); err != nil {
		return err
	}
	desde := 0
	for _, n := range nombres(numeros) {
		var numero int
		if _, err := fmt.Sscan(n, &numero); err != nil {
			return fmt.Errorf("número de ronda mal escrito en el escenario: %q", n)
		}
		ronda, err := w.rondaDelHistorial(numero)
		if err != nil {
			return err
		}
		enlace := fmt.Sprintf(`href="/web/proyectos/%d/rondas/%d"`, pid, ronda.ID)
		i := strings.Index(r.Cuerpo[desde:], enlace)
		if i < 0 {
			return fmt.Errorf("la pantalla del elemento no enlaza a la ronda %d en el orden %q; cuerpo: %.400s", numero, numeros, r.Cuerpo)
		}
		desde += i + len(enlace)
	}
	return nil
}

// --- Auxiliares ---

// pantallaDeRonda abre la pantalla de la ronda con ese número, buscando su ID en el historial.
func (w *historialWorld16) pantallaDeRonda(numero int, titulo string) (respuestaWeb, error) {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return respuestaWeb{}, err
	}
	if err := w.consultar(pid, titulo); err != nil {
		return respuestaWeb{}, err
	}
	ronda, err := w.rondaDelHistorial(numero)
	if err != nil {
		return respuestaWeb{}, err
	}
	return w.sistema.Pedir("GET", fmt.Sprintf("/web/proyectos/%d/rondas/%d", pid, ronda.ID), nil)
}

// ofreceNuevaRonda indica si la pantalla tiene el botón "Nueva ronda" con la acción de abrir una ronda.
func ofreceNuevaRonda(cuerpo, accion string) bool {
	return strings.Contains(cuerpo, `action="`+accion+`"`) && strings.Contains(cuerpo, "Nueva ronda")
}
