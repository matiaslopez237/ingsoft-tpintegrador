package bdd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// sistemaListadoBacklog es la puerta de entrada de los escenarios de HU-07 al núcleo.
// La implementación real está en hu07_sistema_test.go.
type sistemaListadoBacklog interface {
	CrearProyecto(nombre string) (int64, error)
	// CrearElemento carga un elemento con el caso de uso de HU-05 (queda Pendiente).
	CrearElemento(proyectoID int64, titulo, prioridad string) error
	// FijarEstadoElemento prepara un escenario con un elemento en otro estado: todavía no hay
	// una historia que cambie el estado de un elemento (lo hará HU-10).
	FijarEstadoElemento(proyectoID int64, titulo, estado string) error
	// ListarBacklog devuelve los títulos en el orden en que se listan. orden y estado vacíos:
	// los valores por defecto.
	ListarBacklog(proyectoID int64, orden, estado string) ([]string, error)
	Cerrar() error
}

type listadoBacklogWorld struct {
	sistema   sistemaListadoBacklog
	proyectos map[string]int64
	titulos   []string
	ultimoErr error
}

func registrarStepsHU07(ctx *godog.ScenarioContext) {
	w := &listadoBacklogWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaListadoBacklog()
		if err != nil {
			return c, err
		}
		*w = listadoBacklogWorld{sistema: sistema, proyectos: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto de backlog "([^"]*)"$`, w.proyectoDeBacklog)
	ctx.Step(`^en el backlog de "([^"]*)" se cargó "([^"]*)" con prioridad "([^"]*)"$`, w.seCargo)
	ctx.Step(`^en el backlog de "([^"]*)" se cargó "([^"]*)" con prioridad "([^"]*)" y estado "([^"]*)"$`, w.seCargoConEstado)

	// When
	ctx.Step(`^listo el backlog de "([^"]*)"$`, w.listo)
	ctx.Step(`^listo el backlog de "([^"]*)" en orden "([^"]*)"$`, w.listoEnOrden)
	ctx.Step(`^listo el backlog de "([^"]*)" con el estado "([^"]*)"$`, w.listoConEstado)
	ctx.Step(`^listo el backlog de un proyecto inexistente$`, w.listoProyectoInexistente)

	// Then
	ctx.Step(`^el backlog muestra en este orden "([^"]*)"$`, w.muestraEnOrden)
	ctx.Step(`^el backlog está vacío$`, w.estaVacio)
	ctx.Step(`^se rechaza el listado con el error "([^"]*)"$`, w.seRechazaListado)
}

// --- Given ---

func (w *listadoBacklogWorld) proyectoDeBacklog(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *listadoBacklogWorld) seCargo(proyecto, titulo, prioridad string) error {
	id, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearElemento(id, titulo, prioridad); err != nil {
		return fmt.Errorf("no se pudo cargar %q en %q: %w", titulo, proyecto, err)
	}
	return nil
}

func (w *listadoBacklogWorld) seCargoConEstado(proyecto, titulo, prioridad, estado string) error {
	if err := w.seCargo(proyecto, titulo, prioridad); err != nil {
		return err
	}
	if err := w.sistema.FijarEstadoElemento(w.proyectos[proyecto], titulo, estado); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", titulo, estado, err)
	}
	return nil
}

// --- When ---

func (w *listadoBacklogWorld) listar(proyecto, orden, estado string) error {
	id, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	w.titulos, w.ultimoErr = w.sistema.ListarBacklog(id, orden, estado)
	return nil
}

func (w *listadoBacklogWorld) listo(proyecto string) error { return w.listar(proyecto, "", "") }

func (w *listadoBacklogWorld) listoEnOrden(proyecto, orden string) error {
	return w.listar(proyecto, orden, "")
}

func (w *listadoBacklogWorld) listoConEstado(proyecto, estado string) error {
	return w.listar(proyecto, "", estado)
}

func (w *listadoBacklogWorld) listoProyectoInexistente() error {
	w.titulos, w.ultimoErr = w.sistema.ListarBacklog(idProyectoInexistente, "", "")
	return nil
}

// --- Then ---

func (w *listadoBacklogWorld) muestraEnOrden(titulos string) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba listar el backlog, pero hubo error: %w", w.ultimoErr)
	}
	if got := strings.Join(w.titulos, ", "); got != titulos {
		return fmt.Errorf("backlog: se esperaba %q, se obtuvo %q", titulos, got)
	}
	return nil
}

func (w *listadoBacklogWorld) estaVacio() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba listar el backlog, pero hubo error: %w", w.ultimoErr)
	}
	if len(w.titulos) != 0 {
		return fmt.Errorf("se esperaba un backlog vacío, se obtuvo %v", w.titulos)
	}
	return nil
}

func (w *listadoBacklogWorld) seRechazaListado(descripcion string) error {
	esperado, ok := erroresHU07()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el backlog se listó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *listadoBacklogWorld) proyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no preparó el proyecto %q", nombre)
	}
	return id, nil
}
