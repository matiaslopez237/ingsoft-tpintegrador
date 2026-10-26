package bdd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// sprintVisto es lo que los escenarios de HU-08 observan de un sprint. Las fechas van como texto
// AAAA-MM-DD para poder comparar con ==.
type sprintVisto struct {
	Numero   int
	Inicio   string
	Fin      string
	Objetivo string
	Estado   string
}

// sistemaSprints es la puerta de entrada de los escenarios de HU-08 al núcleo.
// La implementación real está en hu08_sistema_test.go.
type sistemaSprints interface {
	CrearProyecto(nombre string, inicio time.Time, fin *time.Time) (int64, error)
	// FijarEstadoProyecto prepara un escenario con un proyecto en otro estado (con HU-02).
	FijarEstadoProyecto(proyectoID int64, estado string) error
	// CrearSprint recibe las fechas como puntero: nil es "sin fecha".
	CrearSprint(proyectoID int64, inicio, fin *time.Time, objetivo string) (sprintVisto, error)
	ListarSprints(proyectoID int64) ([]sprintVisto, error)
	Cerrar() error
}

type sprintsWorld struct {
	sistema   sistemaSprints
	proyectos map[string]int64
	ultimo    sprintVisto
	ultimoErr error
}

func registrarStepsHU08(ctx *godog.ScenarioContext) {
	w := &sprintsWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaSprints()
		if err != nil {
			return c, err
		}
		*w = sprintsWorld{sistema: sistema, proyectos: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^un proyecto "([^"]*)" que va del "([^"]*)" al "([^"]*)"$`, w.proyectoQueVa)
	ctx.Step(`^un proyecto "([^"]*)" que empieza el "([^"]*)" y no tiene fecha de fin$`, w.proyectoSinFin)
	ctx.Step(`^el proyecto "([^"]*)" pasa al estado "([^"]*)"$`, w.proyectoPasaAlEstado)
	ctx.Step(`^"([^"]*)" ya tiene un sprint del "([^"]*)" al "([^"]*)"$`, w.yaTieneSprint)

	// When
	ctx.Step(`^creo en "([^"]*)" un sprint del "([^"]*)" al "([^"]*)" con el Sprint Goal "([^"]*)"$`, w.creoSprint)
	ctx.Step(`^creo en "([^"]*)" un sprint del "([^"]*)" al "([^"]*)" con un Sprint Goal de (\d+) caracteres$`, w.creoSprintConObjetivoDeLargo)
	ctx.Step(`^creo un sprint en un proyecto inexistente$`, w.creoSprintEnProyectoInexistente)
	ctx.Step(`^creo en "([^"]*)" un sprint sin fecha de (inicio|fin)$`, w.creoSprintSinFecha)

	// Then
	ctx.Step(`^el sprint se crea con el número (\d+)$`, w.seCreaConNumero)
	ctx.Step(`^el sprint queda en estado "([^"]*)"$`, w.quedaEnEstado)
	ctx.Step(`^el sprint va del "([^"]*)" al "([^"]*)"$`, w.vaDelAl)
	ctx.Step(`^el sprint tiene el Sprint Goal "([^"]*)"$`, w.tieneObjetivo)
	ctx.Step(`^se rechaza el sprint con el error "([^"]*)"$`, w.seRechazaSprint)
	ctx.Step(`^"([^"]*)" tiene (\d+) sprints?$`, w.tieneSprints)
	ctx.Step(`^los sprints de "([^"]*)" son los números "([^"]*)"$`, w.sprintsSonLosNumeros)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" va del "([^"]*)" al "([^"]*)"$`, w.sprintNumeroVaDelAl)
}

// --- Given ---

func (w *sprintsWorld) proyectoQueVa(nombre, inicio, fin string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	return w.prepararProyecto(nombre, fi, &ff)
}

func (w *sprintsWorld) proyectoSinFin(nombre, inicio string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	return w.prepararProyecto(nombre, fi, nil)
}

func (w *sprintsWorld) prepararProyecto(nombre string, inicio time.Time, fin *time.Time) error {
	id, err := w.sistema.CrearProyecto(nombre, inicio, fin)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *sprintsWorld) proyectoPasaAlEstado(nombre, estado string) error {
	id, err := w.proyecto(nombre)
	if err != nil {
		return err
	}
	if err := w.sistema.FijarEstadoProyecto(id, estado); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", nombre, estado, err)
	}
	return nil
}

func (w *sprintsWorld) yaTieneSprint(nombre, inicio, fin string) error {
	id, err := w.proyecto(nombre)
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
	if _, err := w.sistema.CrearSprint(id, &fi, &ff, "Sprint preparado"); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint de %q: %w", nombre, err)
	}
	return nil
}

// --- When ---

func (w *sprintsWorld) creoSprint(nombre, inicio, fin, objetivo string) error {
	id, err := w.proyecto(nombre)
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
	w.ultimo, w.ultimoErr = w.sistema.CrearSprint(id, &fi, &ff, objetivo)
	return nil
}

func (w *sprintsWorld) creoSprintConObjetivoDeLargo(nombre, inicio, fin string, largo int) error {
	return w.creoSprint(nombre, inicio, fin, strings.Repeat("a", largo))
}

func (w *sprintsWorld) creoSprintEnProyectoInexistente() error {
	fi, ff := time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)
	w.ultimo, w.ultimoErr = w.sistema.CrearSprint(idProyectoInexistente, &fi, &ff, "Objetivo")
	return nil
}

func (w *sprintsWorld) creoSprintSinFecha(nombre, cual string) error {
	id, err := w.proyecto(nombre)
	if err != nil {
		return err
	}
	fi, ff := time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)
	inicio, fin := &fi, &ff
	if cual == "inicio" {
		inicio = nil
	} else {
		fin = nil
	}
	w.ultimo, w.ultimoErr = w.sistema.CrearSprint(id, inicio, fin, "Objetivo")
	return nil
}

// --- Then ---

func (w *sprintsWorld) creado() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba crear el sprint, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *sprintsWorld) seCreaConNumero(numero int) error {
	if err := w.creado(); err != nil {
		return err
	}
	if w.ultimo.Numero != numero {
		return fmt.Errorf("número: se esperaba %d, se obtuvo %d", numero, w.ultimo.Numero)
	}
	return nil
}

func (w *sprintsWorld) quedaEnEstado(estado string) error {
	if err := w.creado(); err != nil {
		return err
	}
	if w.ultimo.Estado != estado {
		return fmt.Errorf("estado: se esperaba %q, se obtuvo %q", estado, w.ultimo.Estado)
	}
	return nil
}

func (w *sprintsWorld) vaDelAl(inicio, fin string) error {
	if err := w.creado(); err != nil {
		return err
	}
	if w.ultimo.Inicio != inicio || w.ultimo.Fin != fin {
		return fmt.Errorf("fechas: se esperaba %s al %s, se obtuvo %s al %s", inicio, fin, w.ultimo.Inicio, w.ultimo.Fin)
	}
	return nil
}

func (w *sprintsWorld) tieneObjetivo(objetivo string) error {
	if err := w.creado(); err != nil {
		return err
	}
	if w.ultimo.Objetivo != objetivo {
		return fmt.Errorf("Sprint Goal: se esperaba %q, se obtuvo %q", objetivo, w.ultimo.Objetivo)
	}
	return nil
}

func (w *sprintsWorld) seRechazaSprint(descripcion string) error {
	esperado, ok := erroresHU08()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el sprint se creó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *sprintsWorld) tieneSprints(nombre string, cantidad int) error {
	sprints, err := w.sprintsDe(nombre)
	if err != nil {
		return err
	}
	if len(sprints) != cantidad {
		return fmt.Errorf("%q tiene %d sprints, se esperaban %d", nombre, len(sprints), cantidad)
	}
	return nil
}

// sprintsSonLosNumeros compara los números en el orden en que se listan, por ejemplo "1, 2".
func (w *sprintsWorld) sprintsSonLosNumeros(nombre, numeros string) error {
	sprints, err := w.sprintsDe(nombre)
	if err != nil {
		return err
	}
	var got []string
	for _, s := range sprints {
		got = append(got, fmt.Sprint(s.Numero))
	}
	if strings.Join(got, ", ") != numeros {
		return fmt.Errorf("sprints de %q: se esperaba %q, se obtuvo %q", nombre, numeros, strings.Join(got, ", "))
	}
	return nil
}

func (w *sprintsWorld) sprintNumeroVaDelAl(numero int, nombre, inicio, fin string) error {
	sprints, err := w.sprintsDe(nombre)
	if err != nil {
		return err
	}
	for _, s := range sprints {
		if s.Numero == numero {
			if s.Inicio != inicio || s.Fin != fin {
				return fmt.Errorf("sprint %d: se esperaba %s al %s, se obtuvo %s al %s", numero, inicio, fin, s.Inicio, s.Fin)
			}
			return nil
		}
	}
	return fmt.Errorf("%q no tiene un sprint número %d", nombre, numero)
}

func (w *sprintsWorld) sprintsDe(nombre string) ([]sprintVisto, error) {
	id, err := w.proyecto(nombre)
	if err != nil {
		return nil, err
	}
	sprints, err := w.sistema.ListarSprints(id)
	if err != nil {
		return nil, fmt.Errorf("listar sprints de %q: %w", nombre, err)
	}
	return sprints, nil
}

func (w *sprintsWorld) proyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no preparó el proyecto %q", nombre)
	}
	return id, nil
}
