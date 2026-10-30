package bdd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sprintAnteriorVisto es lo que los escenarios de HU-12 observan de un sprint del historial. Las
// fechas van como texto AAAA-MM-DD para poder comparar con ==.
type sprintAnteriorVisto struct {
	Numero       int
	Inicio       string
	Fin          string
	Objetivo     string
	Planificados int
	Completados  int
}

// sistemaHistorial es la puerta de entrada de los escenarios de HU-12 al núcleo. La implementa
// sistemaHistorialReal, sobre SQLite (hu12_sistema_test.go).
type sistemaHistorial interface {
	CrearProyecto(nombre string) (int64, error)
	// CrearSprint prepara un sprint con su número, fechas, estado y, si está cerrado, los puntos que
	// habría registrado HU-11. Es preparación del escenario: no pasa por HU-08 ni por HU-11.
	CrearSprint(proyectoID int64, numero int, inicio, fin time.Time, objetivo string,
		estado domain.EstadoSprint, planificados, completados int) error
	ConsultarAnteriores(proyectoID int64) ([]sprintAnteriorVisto, error)
	Cerrar() error
}

// erroresHU12 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU12() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

type historialWorld struct {
	sistema    sistemaHistorial
	proyectos  map[string]int64
	anteriores []sprintAnteriorVisto
	ultimoErr  error
}

func registrarStepsHU12(ctx *godog.ScenarioContext) {
	w := &historialWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = historialWorld{} // si falla la base, el After no cierra el sistema del escenario anterior
		sistema, err := nuevoSistemaHistorial()
		if err != nil {
			return c, err
		}
		*w = historialWorld{sistema: sistema, proyectos: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para el historial$`, w.proyectoParaHistorial)
	ctx.Step(`^"([^"]*)" tiene el sprint (\d+) cerrado del "([^"]*)" al "([^"]*)" con el Sprint Goal "([^"]*)", (\d+) puntos planificados y (\d+) completados$`, w.tieneSprintCerrado)
	ctx.Step(`^"([^"]*)" tiene el sprint (\d+) abierto del "([^"]*)" al "([^"]*)"$`, w.tieneSprintAbierto)

	// When
	ctx.Step(`^consulto los sprints anteriores de "([^"]*)"$`, w.consultoAnteriores)
	ctx.Step(`^consulto los sprints anteriores de un proyecto inexistente$`, w.consultoAnterioresDeInexistente)

	// Then
	ctx.Step(`^los sprints anteriores son los números "([^"]*)"$`, w.anterioresSonLosNumeros)
	ctx.Step(`^no hay sprints anteriores$`, w.noHayAnteriores)
	ctx.Step(`^el sprint anterior (\d+) va del "([^"]*)" al "([^"]*)" con el Sprint Goal "([^"]*)"$`, w.anteriorVaDelAl)
	ctx.Step(`^el sprint anterior (\d+) tiene (\d+) puntos planificados y (\d+) completados$`, w.anteriorTienePuntos)
	ctx.Step(`^se rechaza la consulta de sprints anteriores con el error "([^"]*)"$`, w.seRechazaConsultaAnteriores)
}

// --- Given ---

func (w *historialWorld) proyectoParaHistorial(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *historialWorld) tieneSprintCerrado(proyecto string, numero int, inicio, fin, objetivo string, planificados, completados int) error {
	return w.prepararSprint(proyecto, numero, inicio, fin, objetivo, domain.SprintCerrado, planificados, completados)
}

func (w *historialWorld) tieneSprintAbierto(proyecto string, numero int, inicio, fin string) error {
	return w.prepararSprint(proyecto, numero, inicio, fin, fmt.Sprintf("Objetivo del sprint %d", numero), domain.SprintAbierto, 0, 0)
}

func (w *historialWorld) prepararSprint(proyecto string, numero int, inicio, fin, objetivo string,
	estado domain.EstadoSprint, planificados, completados int) error {
	pid, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearSprint(pid, numero, fi, ff, objetivo, estado, planificados, completados); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint %d de %q: %w", numero, proyecto, err)
	}
	return nil
}

// --- When ---

func (w *historialWorld) consultoAnteriores(proyecto string) error {
	pid, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	w.anteriores, w.ultimoErr = w.sistema.ConsultarAnteriores(pid)
	return nil
}

func (w *historialWorld) consultoAnterioresDeInexistente() error {
	w.anteriores, w.ultimoErr = w.sistema.ConsultarAnteriores(idProyectoInexistente)
	return nil
}

// --- Then ---

func (w *historialWorld) consultado() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar los sprints anteriores, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *historialWorld) anterioresSonLosNumeros(numeros string) error {
	if err := w.consultado(); err != nil {
		return err
	}
	var got []string
	for _, s := range w.anteriores {
		got = append(got, fmt.Sprint(s.Numero))
	}
	if strings.Join(got, ", ") != numeros {
		return fmt.Errorf("sprints anteriores: se esperaba %q, se obtuvo %q", numeros, strings.Join(got, ", "))
	}
	return nil
}

func (w *historialWorld) noHayAnteriores() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if len(w.anteriores) != 0 {
		return fmt.Errorf("se esperaba ningún sprint anterior, se obtuvieron %d", len(w.anteriores))
	}
	return nil
}

func (w *historialWorld) anteriorVaDelAl(numero int, inicio, fin, objetivo string) error {
	s, err := w.anterior(numero)
	if err != nil {
		return err
	}
	if s.Inicio != inicio || s.Fin != fin || s.Objetivo != objetivo {
		return fmt.Errorf("sprint anterior %d: se esperaba %s al %s con %q, se obtuvo %s al %s con %q",
			numero, inicio, fin, objetivo, s.Inicio, s.Fin, s.Objetivo)
	}
	return nil
}

func (w *historialWorld) anteriorTienePuntos(numero, planificados, completados int) error {
	s, err := w.anterior(numero)
	if err != nil {
		return err
	}
	if s.Planificados != planificados || s.Completados != completados {
		return fmt.Errorf("sprint anterior %d: %d planificados y %d completados, se esperaba %d y %d",
			numero, s.Planificados, s.Completados, planificados, completados)
	}
	return nil
}

func (w *historialWorld) seRechazaConsultaAnteriores(descripcion string) error {
	esperado, ok := erroresHU12()[descripcion]
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

func (w *historialWorld) anterior(numero int) (sprintAnteriorVisto, error) {
	if err := w.consultado(); err != nil {
		return sprintAnteriorVisto{}, err
	}
	for _, s := range w.anteriores {
		if s.Numero == numero {
			return s, nil
		}
	}
	return sprintAnteriorVisto{}, fmt.Errorf("el sprint %d no está entre los sprints anteriores", numero)
}

func (w *historialWorld) proyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no preparó el proyecto %q", nombre)
	}
	return id, nil
}
