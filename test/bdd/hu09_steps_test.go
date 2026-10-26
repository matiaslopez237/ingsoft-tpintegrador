package bdd

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
)

// sistemaAsignacion es la puerta de entrada de los escenarios al núcleo.
// La implementación real está en hu09_sistema_test.go.
type sistemaAsignacion interface {
	CrearProyecto(nombre string) (int64, error)
	CrearElemento(proyectoID int64, titulo string) (int64, error)
	CrearSprint(proyectoID int64, abierto bool) int64
	AsignarElemento(elementoID, sprintID int64) error
	DesasignarElemento(elementoID int64) error
	SprintDeElemento(elementoID int64) (sprintID int64, asignado bool, err error)
	Cerrar() error
}

type asignacionWorld struct {
	sistema           sistemaAsignacion
	proyectoID        int64
	elementoPorNombre map[string]int64
	sprintPorNombre   map[string]int64
	ultimoErr         error
}

func registrarStepsHU09(ctx *godog.ScenarioContext) {
	w := &asignacionWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaAsignacion()
		if err != nil {
			return c, err
		}
		*w = asignacionWorld{sistema: sistema, elementoPorNombre: map[string]int64{}, sprintPorNombre: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^un proyecto "([^"]*)" con un elemento de backlog "([^"]*)"$`, w.unProyectoConElemento)
	ctx.Step(`^un sprint abierto "([^"]*)" del proyecto$`, w.unSprintAbiertoDelProyecto)
	ctx.Step(`^un sprint cerrado "([^"]*)" del proyecto$`, w.unSprintCerradoDelProyecto)
	ctx.Step(`^un sprint abierto "([^"]*)" de otro proyecto$`, w.unSprintAbiertoDeOtroProyecto)
	ctx.Step(`^"([^"]*)" ya está asignado al sprint "([^"]*)"$`, w.yaEstaAsignado)

	ctx.Step(`^se asigna "([^"]*)" al sprint "([^"]*)"$`, w.seAsigna)
	ctx.Step(`^se desasigna "([^"]*)"$`, w.seDesasigna)

	ctx.Step(`^"([^"]*)" queda asignado al sprint "([^"]*)"$`, w.estaAsignadoAlSprint)
	ctx.Step(`^"([^"]*)" sigue asignado al sprint "([^"]*)"$`, w.estaAsignadoAlSprint)
	ctx.Step(`^"([^"]*)" no tiene sprint asignado$`, w.noTieneSprintAsignado)
	ctx.Step(`^se rechaza la asignación con el error "([^"]*)"$`, w.seRechazaConError)
}

// --- Given / When ---

func (w *asignacionWorld) unProyectoConElemento(nombreProyecto, nombreElemento string) error {
	pid, err := w.sistema.CrearProyecto(nombreProyecto)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombreProyecto, err)
	}
	w.proyectoID = pid

	eid, err := w.sistema.CrearElemento(pid, nombreElemento)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el elemento %q: %w", nombreElemento, err)
	}
	w.elementoPorNombre[nombreElemento] = eid
	return nil
}

func (w *asignacionWorld) unSprintAbiertoDelProyecto(nombre string) error {
	w.sprintPorNombre[nombre] = w.sistema.CrearSprint(w.proyectoID, true)
	return nil
}

func (w *asignacionWorld) unSprintCerradoDelProyecto(nombre string) error {
	w.sprintPorNombre[nombre] = w.sistema.CrearSprint(w.proyectoID, false)
	return nil
}

func (w *asignacionWorld) unSprintAbiertoDeOtroProyecto(nombre string) error {
	otroProyectoID, err := w.sistema.CrearProyecto("Otro proyecto")
	if err != nil {
		return fmt.Errorf("no se pudo preparar el otro proyecto: %w", err)
	}
	w.sprintPorNombre[nombre] = w.sistema.CrearSprint(otroProyectoID, true)
	return nil
}

func (w *asignacionWorld) yaEstaAsignado(nombreElemento, nombreSprint string) error {
	err := w.sistema.AsignarElemento(w.elementoPorNombre[nombreElemento], w.sprintPorNombre[nombreSprint])
	if err != nil {
		return fmt.Errorf("no se pudo preparar la asignación de %q al sprint %q: %w", nombreElemento, nombreSprint, err)
	}
	return nil
}

func (w *asignacionWorld) seAsigna(nombreElemento, nombreSprint string) error {
	w.ultimoErr = w.sistema.AsignarElemento(w.elementoPorNombre[nombreElemento], w.sprintPorNombre[nombreSprint])
	return nil
}

func (w *asignacionWorld) seDesasigna(nombreElemento string) error {
	w.ultimoErr = w.sistema.DesasignarElemento(w.elementoPorNombre[nombreElemento])
	return nil
}

// --- Then ---

func (w *asignacionWorld) estaAsignadoAlSprint(nombreElemento, nombreSprint string) error {
	sprintID, asignado, err := w.sistema.SprintDeElemento(w.elementoPorNombre[nombreElemento])
	if err != nil {
		return fmt.Errorf("consultar sprint de %q: %w", nombreElemento, err)
	}
	if !asignado {
		return fmt.Errorf("%q no tiene sprint asignado, se esperaba %q", nombreElemento, nombreSprint)
	}
	if esperado := w.sprintPorNombre[nombreSprint]; sprintID != esperado {
		return fmt.Errorf("%q está asignado al sprint %d, se esperaba %d (%q)", nombreElemento, sprintID, esperado, nombreSprint)
	}
	return nil
}

func (w *asignacionWorld) noTieneSprintAsignado(nombreElemento string) error {
	_, asignado, err := w.sistema.SprintDeElemento(w.elementoPorNombre[nombreElemento])
	if err != nil {
		return fmt.Errorf("consultar sprint de %q: %w", nombreElemento, err)
	}
	if asignado {
		return fmt.Errorf("%q tiene un sprint asignado, se esperaba que no tuviera ninguno", nombreElemento)
	}
	return nil
}

func (w *asignacionWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU09()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la operación se realizó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}
