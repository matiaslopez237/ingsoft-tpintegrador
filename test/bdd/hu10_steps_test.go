package bdd

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sistemaCompletar es la puerta de entrada de los escenarios de HU-10 al núcleo.
// La implementa sistemaCompletarReal, con SQLite real (hu10_sistema_test.go).
type sistemaCompletar interface {
	CrearProyecto(nombre string) (int64, error)
	CrearElemento(proyectoID int64, titulo string) (int64, error)
	CrearSprint(proyectoID int64, abierto bool) (int64, error)
	// PrepararElemento deja el elemento en un sprint (nil: sin sprint) y un estado, sin pasar por
	// los casos de uso: es preparación del escenario, como AsignarDirectamente en HU-09.
	PrepararElemento(elementoID int64, sprintID *int64, estado domain.EstadoElemento) error
	CompletarHistoria(proyectoID, elementoID int64) error
	Elemento(elementoID int64) (domain.ElementoBacklog, error)
	Cerrar() error
}

// erroresHU10 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU10() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"elemento no encontrado": domain.ErrElementoNoEncontrado,
		"elemento sin sprint":    domain.ErrElementoSinSprint,
		"sprint cerrado":         domain.ErrSprintCerrado,
	}
}

type completarHistoriaWorld struct {
	sistema    sistemaCompletar
	proyectos  map[string]int64
	elementos  map[string]int64
	sprints    map[string]int64
	proyectoID int64 // el proyecto del Background: ahí se marca si no se indica otro
	ultimoErr  error
}

func registrarStepsHU10(ctx *godog.ScenarioContext) {
	w := &completarHistoriaWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaCompletar()
		if err != nil {
			return c, err
		}
		*w = completarHistoriaWorld{
			sistema:   sistema,
			proyectos: map[string]int64{},
			elementos: map[string]int64{},
			sprints:   map[string]int64{},
		}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^el proyecto "([^"]*)" tiene en su backlog el elemento "([^"]*)"$`, w.proyectoTieneElemento)
	ctx.Step(`^el proyecto "([^"]*)" tiene un sprint (abierto|cerrado) "([^"]*)"$`, w.proyectoTieneSprint)
	ctx.Step(`^"([^"]*)" está en el sprint "([^"]*)" con estado "([^"]*)"$`, w.estaEnSprintConEstado)
	ctx.Step(`^"([^"]*)" no está en ningún sprint y tiene estado "([^"]*)"$`, w.sinSprintConEstado)

	ctx.Step(`^se marca "([^"]*)" como completada$`, w.seMarcaCompletada)
	ctx.Step(`^se marca "([^"]*)" como completada en un proyecto inexistente$`, w.seMarcaEnProyectoInexistente)
	ctx.Step(`^se marca como completada una historia inexistente$`, w.seMarcaHistoriaInexistente)

	ctx.Step(`^la historia se marca como completada$`, w.seMarcaCorrectamente)
	ctx.Step(`^"([^"]*)" queda en estado "([^"]*)"$`, w.tieneEstado)
	ctx.Step(`^"([^"]*)" sigue en estado "([^"]*)"$`, w.tieneEstado)
	ctx.Step(`^"([^"]*)" sigue en el sprint "([^"]*)"$`, w.sigueEnSprint)
	ctx.Step(`^se rechaza completar la historia con el error "([^"]*)"$`, w.seRechazaConError)
}

// --- Given ---

func (w *completarHistoriaWorld) proyectoTieneElemento(proyecto, elemento string) error {
	pid, ok := w.proyectos[proyecto]
	if !ok {
		var err error
		if pid, err = w.sistema.CrearProyecto(proyecto); err != nil {
			return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
		}
		w.proyectos[proyecto] = pid
	}
	if w.proyectoID == 0 {
		w.proyectoID = pid
	}
	eid, err := w.sistema.CrearElemento(pid, elemento)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el elemento %q: %w", elemento, err)
	}
	w.elementos[elemento] = eid
	return nil
}

func (w *completarHistoriaWorld) proyectoTieneSprint(proyecto, estado, sprint string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.sistema.CrearSprint(pid, estado == "abierto")
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint %q: %w", sprint, err)
	}
	w.sprints[sprint] = sid
	return nil
}

func (w *completarHistoriaWorld) estaEnSprintConEstado(elemento, sprint, estado string) error {
	eid, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(sprint)
	if err != nil {
		return err
	}
	if err := w.sistema.PrepararElemento(eid, &sid, domain.EstadoElemento(estado)); err != nil {
		return fmt.Errorf("no se pudo preparar %q en el sprint %q: %w", elemento, sprint, err)
	}
	return nil
}

func (w *completarHistoriaWorld) sinSprintConEstado(elemento, estado string) error {
	eid, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	if err := w.sistema.PrepararElemento(eid, nil, domain.EstadoElemento(estado)); err != nil {
		return fmt.Errorf("no se pudo preparar %q sin sprint: %w", elemento, err)
	}
	return nil
}

// --- When ---

func (w *completarHistoriaWorld) seMarcaCompletada(elemento string) error {
	eid, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.CompletarHistoria(w.proyectoID, eid)
	return nil
}

func (w *completarHistoriaWorld) seMarcaEnProyectoInexistente(elemento string) error {
	eid, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.CompletarHistoria(idInexistente, eid)
	return nil
}

func (w *completarHistoriaWorld) seMarcaHistoriaInexistente() error {
	w.ultimoErr = w.sistema.CompletarHistoria(w.proyectoID, idInexistente)
	return nil
}

// --- Then ---

func (w *completarHistoriaWorld) seMarcaCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba marcar la historia como completada, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *completarHistoriaWorld) tieneEstado(elemento, estado string) error {
	e, err := w.elemento(elemento)
	if err != nil {
		return err
	}
	if e.Estado != domain.EstadoElemento(estado) {
		return fmt.Errorf("estado de %q: se esperaba %q, es %q", elemento, estado, e.Estado)
	}
	return nil
}

func (w *completarHistoriaWorld) sigueEnSprint(elemento, sprint string) error {
	e, err := w.elemento(elemento)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(sprint)
	if err != nil {
		return err
	}
	if e.SprintID == nil || *e.SprintID != sid {
		return fmt.Errorf("%q: se esperaba en el sprint %q (%d), tiene %v", elemento, sprint, sid, e.SprintID)
	}
	return nil
}

func (w *completarHistoriaWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU10()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la historia se marcó como completada", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

// --- Auxiliares ---

func (w *completarHistoriaWorld) elemento(titulo string) (domain.ElementoBacklog, error) {
	eid, err := w.idElemento(titulo)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}
	e, err := w.sistema.Elemento(eid)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("no se pudo leer %q: %w", titulo, err)
	}
	return e, nil
}

func (w *completarHistoriaWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *completarHistoriaWorld) idElemento(titulo string) (int64, error) {
	id, ok := w.elementos[titulo]
	if !ok {
		return 0, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return id, nil
}

func (w *completarHistoriaWorld) idSprint(nombre string) (int64, error) {
	id, ok := w.sprints[nombre]
	if !ok {
		return 0, fmt.Errorf("el sprint %q no se preparó en el escenario", nombre)
	}
	return id, nil
}
