package bdd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// velocidadVista es lo que los escenarios de HU-22 observan de la velocidad de un proyecto.
type velocidadVista struct {
	Velocidad       float64
	HayDatos        bool
	SprintsCerrados int
}

// sistemaVelocidad es la puerta de entrada de los escenarios de HU-22 al núcleo. La implementa
// sistemaVelocidadReal, sobre SQLite (hu22_sistema_test.go).
type sistemaVelocidad interface {
	CrearProyecto(nombre string) (int64, error)
	// CrearSprint prepara un sprint cerrado (con sus puntos completados) o abierto. Es preparación
	// del escenario: no pasa por HU-08 ni por HU-11, igual que en HU-12.
	CrearSprint(proyectoID int64, estado domain.EstadoSprint, completados int) error
	ConsultarVelocidad(proyectoID int64) (velocidadVista, error)
	Cerrar() error
}

// erroresHU22 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU22() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

type velocidadWorld struct {
	sistema   sistemaVelocidad
	proyectos map[string]int64
	ultima    velocidadVista
	ultimoErr error
}

func registrarStepsHU22(ctx *godog.ScenarioContext) {
	w := &velocidadWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = velocidadWorld{}
		sistema, err := nuevoSistemaVelocidad()
		if err != nil {
			return c, err
		}
		*w = velocidadWorld{sistema: sistema, proyectos: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para la velocidad$`, w.proyectoParaVelocidad)
	ctx.Step(`^"([^"]*)" tiene un sprint cerrado que completó (\d+) puntos$`, w.tieneSprintCerrado)
	ctx.Step(`^"([^"]*)" tiene un sprint abierto$`, w.tieneSprintAbierto)

	// When
	ctx.Step(`^consulto la velocidad de "([^"]*)"$`, w.consultoVelocidad)
	ctx.Step(`^consulto la velocidad de un proyecto inexistente$`, w.consultoVelocidadDeInexistente)

	// Then
	ctx.Step(`^la velocidad es (-?\d+(?:,\d+)?) sobre (\d+) sprints cerrados$`, w.velocidadEs)
	ctx.Step(`^no hay velocidad todavía$`, w.noHayVelocidad)
	ctx.Step(`^se rechaza la consulta de velocidad con el error "([^"]*)"$`, w.seRechazaConsulta)
}

// --- Given ---

func (w *velocidadWorld) proyectoParaVelocidad(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *velocidadWorld) tieneSprintCerrado(proyecto string, completados int) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearSprint(id, domain.SprintCerrado, completados); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint cerrado del proyecto %q: %w", proyecto, err)
	}
	return nil
}

func (w *velocidadWorld) tieneSprintAbierto(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearSprint(id, domain.SprintAbierto, 0); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint abierto del proyecto %q: %w", proyecto, err)
	}
	return nil
}

// --- When ---

func (w *velocidadWorld) consultoVelocidad(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultima, w.ultimoErr = w.sistema.ConsultarVelocidad(id)
	return nil
}

func (w *velocidadWorld) consultoVelocidadDeInexistente() error {
	w.ultima, w.ultimoErr = w.sistema.ConsultarVelocidad(idInexistente)
	return nil
}

// --- Then ---

func (w *velocidadWorld) velocidadEs(velocidad string, sprintsCerrados int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba calcular la velocidad, pero hubo error: %w", w.ultimoErr)
	}
	if !w.ultima.HayDatos {
		return fmt.Errorf("se esperaba velocidad con datos, pero no hay")
	}
	want, err := strconv.ParseFloat(strings.ReplaceAll(velocidad, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("velocidad inválida en el escenario %q: %w", velocidad, err)
	}
	if w.ultima.Velocidad != want || w.ultima.SprintsCerrados != sprintsCerrados {
		return fmt.Errorf("velocidad = %v sobre %d sprints cerrados, se esperaba %v sobre %d",
			w.ultima.Velocidad, w.ultima.SprintsCerrados, want, sprintsCerrados)
	}
	return nil
}

func (w *velocidadWorld) noHayVelocidad() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la velocidad, pero hubo error: %w", w.ultimoErr)
	}
	if w.ultima.HayDatos {
		return fmt.Errorf("se esperaba que no hubiera velocidad, pero hay %v sobre %d sprints cerrados", w.ultima.Velocidad, w.ultima.SprintsCerrados)
	}
	return nil
}

func (w *velocidadWorld) seRechazaConsulta(descripcion string) error {
	esperado, ok := erroresHU22()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la consulta funcionó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

// --- Auxiliares ---

func (w *velocidadWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}
