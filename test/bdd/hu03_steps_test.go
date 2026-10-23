package bdd

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cucumber/godog"
)

// sistemaIntegrantes es la puerta de entrada de los escenarios al núcleo.
// La implementación real está en hu03_sistema_test.go.
type sistemaIntegrantes interface {
	CrearProyecto(nombre string) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre, rol string) error
	QuitarIntegrante(proyectoID int64, nombre string) error
	Integrantes(proyectoID int64) ([]string, error)
	Cerrar() error
}

type integrantesWorld struct {
	sistema    sistemaIntegrantes
	proyectoID int64
	ultimoErr  error
}

func registrarStepsHU03(ctx *godog.ScenarioContext) {
	w := &integrantesWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaIntegrantes()
		if err != nil {
			return c, err
		}
		*w = integrantesWorld{sistema: sistema}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^un proyecto "([^"]*)" sin integrantes$`, w.unProyectoSinIntegrantes)
	ctx.Step(`^el integrante "([^"]*)" con rol "([^"]*)" ya está en el proyecto$`, w.yaEstaEnElProyecto)
	ctx.Step(`^los integrantes "([^"]*)" y "([^"]*)" ya están en el proyecto$`, w.yaEstanEnElProyecto)

	ctx.Step(`^se agrega el integrante "([^"]*)" con rol "([^"]*)"$`, w.seAgrega)
	ctx.Step(`^se quita al integrante "([^"]*)"$`, w.seQuita)

	ctx.Step(`^el proyecto tiene como integrantes a "([^"]*)"$`, w.tieneComoIntegrantesA)
	ctx.Step(`^el proyecto tiene como integrantes a "([^"]*)" y "([^"]*)"$`, w.tieneComoIntegrantesAAmbos)
	ctx.Step(`^el proyecto no tiene integrantes$`, w.noTieneIntegrantes)
	ctx.Step(`^se rechaza la operación con el error "([^"]*)"$`, w.seRechazaConError)
}

// --- Given / When ---

func (w *integrantesWorld) unProyectoSinIntegrantes(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectoID = id
	return nil
}

func (w *integrantesWorld) yaEstaEnElProyecto(nombre, rol string) error {
	if err := w.sistema.AgregarIntegrante(w.proyectoID, nombre, rol); err != nil {
		return fmt.Errorf("no se pudo preparar el integrante %q: %w", nombre, err)
	}
	return nil
}

func (w *integrantesWorld) yaEstanEnElProyecto(a, b string) error {
	if err := w.yaEstaEnElProyecto(a, "Integrante"); err != nil {
		return err
	}
	return w.yaEstaEnElProyecto(b, "Integrante")
}

func (w *integrantesWorld) seAgrega(nombre, rol string) error {
	w.ultimoErr = w.sistema.AgregarIntegrante(w.proyectoID, nombre, rol)
	return nil
}

func (w *integrantesWorld) seQuita(nombre string) error {
	w.ultimoErr = w.sistema.QuitarIntegrante(w.proyectoID, nombre)
	return nil
}

// --- Then ---

func (w *integrantesWorld) tieneComoIntegrantesA(nombre string) error {
	return w.integrantesSonExactamente(nombre)
}

func (w *integrantesWorld) tieneComoIntegrantesAAmbos(a, b string) error {
	return w.integrantesSonExactamente(a, b)
}

func (w *integrantesWorld) integrantesSonExactamente(esperados ...string) error {
	actuales, err := w.sistema.Integrantes(w.proyectoID)
	if err != nil {
		return fmt.Errorf("consultar integrantes: %w", err)
	}
	got := slices.Clone(actuales)
	want := slices.Clone(esperados)
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		return fmt.Errorf("integrantes: se esperaba %v, se obtuvo %v", want, got)
	}
	return nil
}

func (w *integrantesWorld) noTieneIntegrantes() error {
	return w.integrantesSonExactamente()
}

func (w *integrantesWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU03()[descripcion]
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
