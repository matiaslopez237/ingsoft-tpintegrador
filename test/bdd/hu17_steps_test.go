package bdd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU17 traduce la descripción del escenario al error de la spec (sección 7).
func erroresHU17() map[string]error {
	return map[string]error{
		"proyecto no encontrado":   domain.ErrProyectoNoEncontrado,
		"ronda no encontrada":      domain.ErrRondaNoEncontrada,
		"elemento terminado":       domain.ErrElementoTerminado,
		"ronda abierta":            domain.ErrEstimacionRondaAbierta,
		"no es la última":          domain.ErrRondaNoEsLaUltima,
		"estimación ya registrada": domain.ErrEstimacionYaRegistrada,
		"estimación inválida":      domain.ErrEstimacionInvalida,
	}
}

// sistemaEstimacion es lo que los escenarios de HU-17 suman al sistema de HU-14. Lo cumple
// sistemaPlanningPokerReal con los métodos de hu17_sistema_test.go, sin tocar la interfaz de HU-14.
type sistemaEstimacion interface {
	RegistrarEstimacion(proyectoID, rondaID int64, valor int) error
	EstimacionDeRonda(proyectoID, rondaID int64) (*int, error)
	Elemento(elementoID int64) (domain.ElementoBacklog, error)
	EditarStoryPoints(proyectoID, elementoID int64, valor int) error
	FijarEstadoElemento(proyectoID, elementoID int64, estado string) error
}

// estimacionWorld es el estado propio de los escenarios de HU-17, sobre el mundo de HU-14.
type estimacionWorld struct {
	*planningPokerWorld
	err      error
	anterior *domain.ElementoBacklog // el elemento antes de registrar, para ver qué no cambió
}

// registrarStepsHU17 agrega los pasos de HU-17 sobre el mismo mundo de HU-14 (registrarStepsHU14
// llama a esta función, como a las de HU-15 y HU-16).
func registrarStepsHU17(ctx *godog.ScenarioContext, pw *planningPokerWorld) {
	w := &estimacionWorld{planningPokerWorld: pw}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		w.err, w.anterior = nil, nil
		return c, nil
	})

	// Given
	ctx.Step(`^"([^"]*)" tiene (\d+) story points cargados a mano$`, w.tieneStoryPoints)
	ctx.Step(`^"([^"]*)" pasó a estado "([^"]*)"$`, w.pasoAEstado)
	ctx.Step(`^registré la estimación (\d+) en la ronda (\d+) de "([^"]*)"$`, w.registreEstimacion)

	// When
	ctx.Step(`^registro la estimación (-?\d+) en la ronda (\d+) de "([^"]*)"$`, w.registroEstimacion)
	ctx.Step(`^registro la estimación (\d+) en una ronda inexistente de "([^"]*)"$`, w.registroEnRondaInexistente)
	ctx.Step(`^registro la estimación (\d+) en la ronda (\d+) de "([^"]*)" desde un proyecto inexistente$`, w.registroDesdeInexistente)
	ctx.Step(`^registro la estimación (\d+) en la ronda (\d+) de "([^"]*)" desde el proyecto "([^"]*)"$`, w.registroDesdeProyecto)
	ctx.Step(`^edito a mano los story points de "([^"]*)" a (\d+)$`, w.editoStoryPoints)
	ctx.Step(`^desde la pantalla de la ronda (\d+) de "([^"]*)" registro la estimación (\d+)$`, w.desdePantallaRegistro)

	// Then
	ctx.Step(`^la estimación se registra sin errores$`, w.sinErrores)
	ctx.Step(`^se rechaza la estimación con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^la ronda (\d+) de "([^"]*)" tiene la estimación acordada (\d+)$`, w.rondaTieneEstimacion)
	ctx.Step(`^la ronda (\d+) de "([^"]*)" no tiene estimación acordada$`, w.rondaSinEstimacion)
	ctx.Step(`^"([^"]*)" tiene (\d+) story points$`, w.elementoTienePuntos)
	ctx.Step(`^"([^"]*)" sigue sin estimar$`, w.elementoSinEstimar)
	ctx.Step(`^los demás datos de "([^"]*)" no cambiaron$`, w.demasDatosIguales)
	ctx.Step(`^la pantalla de la ronda (\d+) de "([^"]*)" muestra la estimación acordada (\d+)$`, w.pantallaMuestraEstimacion)
	ctx.Step(`^la pantalla de la ronda (\d+) de "([^"]*)" no ofrece registrar la estimación$`, w.pantallaNoOfrece)
	ctx.Step(`^el historial de la pantalla de "([^"]*)" muestra acordada (\d+)$`, w.historialPantallaMuestra)
}

func (w *estimacionWorld) estimacion() sistemaEstimacion { return w.sistema.(sistemaEstimacion) }

// rondaID busca el ID de la ronda con ese número en el historial del elemento (HU-16).
func (w *estimacionWorld) rondaID(numero int, titulo string) (int64, error) {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return 0, err
	}
	historial, err := w.sistema.(sistemaHistorialRondas).HistorialRondas(pid, eid)
	if err != nil {
		return 0, err
	}
	i := slices.IndexFunc(historial, func(r rondaVista) bool { return r.Numero == numero })
	if i < 0 {
		return 0, fmt.Errorf("%q no tiene la ronda %d", titulo, numero)
	}
	return historial[i].ID, nil
}

// --- Given ---

func (w *estimacionWorld) tieneStoryPoints(titulo string, valor int) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if err := w.estimacion().EditarStoryPoints(pid, eid, valor); err != nil {
		return fmt.Errorf("preparar los story points de %q: %w", titulo, err)
	}
	e, err := w.estimacion().Elemento(eid)
	if err != nil {
		return err
	}
	w.anterior = &e
	return nil
}

func (w *estimacionWorld) pasoAEstado(titulo, estado string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	return w.estimacion().FijarEstadoElemento(pid, eid, estado)
}

func (w *estimacionWorld) registreEstimacion(valor, numero int, titulo string) error {
	if err := w.registroEstimacion(valor, numero, titulo); err != nil {
		return err
	}
	if w.err != nil {
		return fmt.Errorf("preparar la estimación de la ronda %d: %w", numero, w.err)
	}
	return nil
}

// --- When ---

func (w *estimacionWorld) registrar(pid, rid int64, valor int) error {
	w.err = w.estimacion().RegistrarEstimacion(pid, rid, valor)
	if errors.Is(w.err, godog.ErrPending) {
		return godog.ErrPending
	}
	return nil
}

func (w *estimacionWorld) registroEstimacion(valor, numero int, titulo string) error {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	rid, err := w.rondaID(numero, titulo)
	if err != nil {
		return err
	}
	return w.registrar(pid, rid, valor)
}

func (w *estimacionWorld) registroEnRondaInexistente(valor int, proyecto string) error {
	pid, ok := w.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	return w.registrar(pid, 999999, valor)
}

func (w *estimacionWorld) registroDesdeInexistente(valor, numero int, titulo string) error {
	rid, err := w.rondaID(numero, titulo)
	if err != nil {
		return err
	}
	return w.registrar(idProyectoInexistente, rid, valor)
}

func (w *estimacionWorld) registroDesdeProyecto(valor, numero int, titulo, proyecto string) error {
	rid, err := w.rondaID(numero, titulo)
	if err != nil {
		return err
	}
	pid, ok := w.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	return w.registrar(pid, rid, valor)
}

func (w *estimacionWorld) editoStoryPoints(titulo string, valor int) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	return w.estimacion().EditarStoryPoints(pid, eid, valor)
}

// desdePantallaRegistro abre la pantalla de la ronda, verifica que ofrezca el formulario y lo envía.
func (w *estimacionWorld) desdePantallaRegistro(numero int, titulo string, valor int) error {
	ruta, accion, err := w.rutasRonda(numero, titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if !strings.Contains(r.Cuerpo, `action="`+accion+`"`) {
		return fmt.Errorf("la pantalla de la ronda %d no ofrece registrar la estimación (estado %d)", numero, r.Estado)
	}
	w.pantalla, err = w.sistema.Pedir("POST", accion, url.Values{"estimacion": {strconv.Itoa(valor)}})
	if err != nil {
		return err
	}
	if w.pantalla.Estado != 303 {
		return fmt.Errorf("se esperaba volver a la ronda con 303, se obtuvo %d: %.300s", w.pantalla.Estado, w.pantalla.Cuerpo)
	}
	return nil
}

// --- Then ---

func (w *estimacionWorld) sinErrores() error {
	if w.err != nil {
		return fmt.Errorf("no se esperaba un error: %w", w.err)
	}
	return nil
}

func (w *estimacionWorld) seRechaza(descripcion string) error {
	esperado, ok := erroresHU17()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.err == nil {
		return fmt.Errorf("se esperaba el error %q, pero la estimación se registró", descripcion)
	}
	if !errors.Is(w.err, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.err)
	}
	return nil
}

func (w *estimacionWorld) estimacionDe(numero int, titulo string) (*int, error) {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return nil, err
	}
	rid, err := w.rondaID(numero, titulo)
	if err != nil {
		return nil, err
	}
	return w.estimacion().EstimacionDeRonda(pid, rid)
}

func (w *estimacionWorld) rondaTieneEstimacion(numero int, titulo string, valor int) error {
	got, err := w.estimacionDe(numero, titulo)
	if err != nil {
		return err
	}
	if got == nil || *got != valor {
		return fmt.Errorf("estimación de la ronda %d = %v, se esperaba %d", numero, got, valor)
	}
	return nil
}

func (w *estimacionWorld) rondaSinEstimacion(numero int, titulo string) error {
	got, err := w.estimacionDe(numero, titulo)
	if err != nil {
		return err
	}
	if got != nil {
		return fmt.Errorf("la ronda %d tiene la estimación %d, se esperaba ninguna", numero, *got)
	}
	return nil
}

func (w *estimacionWorld) elementoGuardado(titulo string) (domain.ElementoBacklog, error) {
	_, eid, err := w.elemento(titulo)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}
	return w.estimacion().Elemento(eid)
}

func (w *estimacionWorld) elementoTienePuntos(titulo string, valor int) error {
	e, err := w.elementoGuardado(titulo)
	if err != nil {
		return err
	}
	if e.StoryPoints == nil || *e.StoryPoints != valor {
		return fmt.Errorf("story points de %q = %v, se esperaba %d", titulo, e.StoryPoints, valor)
	}
	return nil
}

func (w *estimacionWorld) elementoSinEstimar(titulo string) error {
	e, err := w.elementoGuardado(titulo)
	if err != nil {
		return err
	}
	if e.StoryPoints != nil {
		return fmt.Errorf("%q tiene %d story points, se esperaba sin estimar", titulo, *e.StoryPoints)
	}
	return nil
}

func (w *estimacionWorld) demasDatosIguales(titulo string) error {
	e, err := w.elementoGuardado(titulo)
	if err != nil {
		return err
	}
	if w.anterior == nil {
		return errors.New("el escenario no guardó el elemento antes de registrar")
	}
	a := *w.anterior
	if e.Titulo != a.Titulo || e.Descripcion != a.Descripcion || e.Prioridad != a.Prioridad || e.Estado != a.Estado ||
		!slices.Equal(e.Criterios, a.Criterios) {
		return fmt.Errorf("el elemento cambió: %+v, se esperaba %+v salvo los story points", e, a)
	}
	return nil
}

func (w *estimacionWorld) rutasRonda(numero int, titulo string) (ruta, accion string, err error) {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return "", "", err
	}
	rid, err := w.rondaID(numero, titulo)
	if err != nil {
		return "", "", err
	}
	ruta = fmt.Sprintf("/web/proyectos/%d/rondas/%d", pid, rid)
	return ruta, ruta + "/estimacion", nil
}

func (w *estimacionWorld) pantallaMuestraEstimacion(numero int, titulo string, valor int) error {
	ruta, _, err := w.rutasRonda(numero, titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if want := fmt.Sprintf("Estimación acordada: %d", valor); !strings.Contains(r.Cuerpo, want) {
		return fmt.Errorf("la pantalla de la ronda no muestra %q: %.500s", want, r.Cuerpo)
	}
	return nil
}

func (w *estimacionWorld) pantallaNoOfrece(numero int, titulo string) error {
	ruta, accion, err := w.rutasRonda(numero, titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if strings.Contains(r.Cuerpo, `action="`+accion+`"`) {
		return fmt.Errorf("la pantalla de la ronda %d no debería ofrecer registrar la estimación", numero)
	}
	return nil
}

func (w *estimacionWorld) historialPantallaMuestra(titulo string, valor int) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid), nil)
	if err != nil {
		return err
	}
	if want := fmt.Sprintf("acordada: %d", valor); !strings.Contains(r.Cuerpo, want) {
		return fmt.Errorf("el historial del elemento no muestra %q", want)
	}
	return nil
}
