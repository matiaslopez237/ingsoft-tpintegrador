package bdd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// proyectoHU02 es lo que los escenarios de HU-02 necesitan observar de un proyecto.
// Las fechas van como texto "AAAA-MM-DD" ("" si no tiene fecha de fin) para poder comparar con ==.
type proyectoHU02 struct {
	ID     int64
	Nombre string
	Inicio string
	Fin    string
	Estado string
}

// cambiosProyecto describe qué campos se envían en una modificación (edición parcial, RN-02.1).
// nil = el campo no se envía. Los Quitar* representan enviar el campo explícitamente vacío.
type cambiosProyecto struct {
	Nombre       *string
	FechaInicio  *time.Time
	QuitarInicio bool
	FechaFin     *time.Time
	QuitarFin    bool
	Estado       *string
}

// sistemaModificacion es la puerta de entrada de los escenarios de HU-02 al núcleo.
// La implementación está en hu02_sistema_test.go.
type sistemaModificacion interface {
	Crear(nombre string, inicio, fin *time.Time) (proyectoHU02, error)
	Modificar(id int64, c cambiosProyecto) (proyectoHU02, error)
	BuscarPorID(id int64) (proyectoHU02, error)
	// FijarEstado prepara un escenario con un estado dado, sin pasar por la funcionalidad bajo prueba.
	FijarEstado(id int64, estado string) error
	Cerrar() error
}

type modificarProyectoWorld struct {
	sistema    sistemaModificacion
	porNombre  map[string]proyectoHU02 // cómo quedó cada proyecto al preparar el escenario
	objetivoID int64                   // ID del proyecto que se intentó modificar
	ultimo     proyectoHU02
	ultimoErr  error
}

func registrarStepsHU02(ctx *godog.ScenarioContext) {
	w := &modificarProyectoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaModificacion()
		if err != nil {
			return c, err
		}
		*w = modificarProyectoWorld{sistema: sistema, porNombre: map[string]proyectoHU02{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^hay un proyecto "([^"]*)" con inicio "([^"]*)" y fin "([^"]*)"$`, w.hayUnProyecto)
	ctx.Step(`^el proyecto "([^"]*)" está en estado "([^"]*)"$`, w.estaEnEstado)

	// When
	ctx.Step(`^modifico el proyecto "([^"]*)" con el nombre "([^"]*)"$`, w.modificoNombre)
	ctx.Step(`^modifico el proyecto "([^"]*)" con un nombre de (\d+) caracteres$`, w.modificoNombreDeLargo)
	ctx.Step(`^modifico el proyecto "([^"]*)" con la fecha de inicio "([^"]*)"$`, w.modificoFechaInicio)
	ctx.Step(`^modifico el proyecto "([^"]*)" con la fecha de fin "([^"]*)"$`, w.modificoFechaFin)
	ctx.Step(`^modifico el proyecto "([^"]*)" quitando la fecha de inicio$`, w.modificoQuitandoInicio)
	ctx.Step(`^modifico el proyecto "([^"]*)" quitando la fecha de fin$`, w.modificoQuitandoFin)
	ctx.Step(`^modifico el proyecto "([^"]*)" con el estado "([^"]*)"$`, w.modificoEstado)
	ctx.Step(`^modifico el proyecto "([^"]*)" sin indicar ningún dato$`, w.modificoSinDatos)
	ctx.Step(`^modifico el proyecto con identificador (\d+) con el nombre "([^"]*)"$`, w.modificoPorIDInexistente)

	// Then
	ctx.Step(`^la modificación se realiza correctamente$`, w.seRealizaCorrectamente)
	ctx.Step(`^el proyecto modificado se llama "([^"]*)"$`, w.modificadoSeLlama)
	ctx.Step(`^el proyecto modificado tiene inicio "([^"]*)" y fin "([^"]*)"$`, w.modificadoTieneFechas)
	ctx.Step(`^el proyecto modificado no tiene fecha de fin$`, w.modificadoSinFechaFin)
	ctx.Step(`^el proyecto modificado está en estado "([^"]*)"$`, w.modificadoEnEstado)
	ctx.Step(`^se rechaza la modificación con el error "([^"]*)"$`, w.seRechazaModificacion)
	ctx.Step(`^el proyecto "([^"]*)" no cambió$`, w.noCambio)
	ctx.Step(`^el proyecto modificado conserva su identificador$`, w.conservaIdentificador)
	ctx.Step(`^al recuperarlo por su identificador se llama "([^"]*)"$`, w.alRecuperarloSeLlama)
}

// --- Given ---

func (w *modificarProyectoWorld) hayUnProyecto(nombre, inicio, fin string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	p, err := w.sistema.Crear(nombre, &fi, &ff)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.porNombre[nombre] = p
	return nil
}

func (w *modificarProyectoWorld) estaEnEstado(nombre, estado string) error {
	p, err := w.proyecto(nombre)
	if err != nil {
		return err
	}
	if err := w.sistema.FijarEstado(p.ID, estado); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", nombre, estado, err)
	}
	p.Estado = estado
	w.porNombre[nombre] = p
	return nil
}

// --- When ---

func (w *modificarProyectoWorld) modificoNombre(proyecto, nombre string) error {
	return w.modificar(proyecto, cambiosProyecto{Nombre: &nombre})
}

func (w *modificarProyectoWorld) modificoNombreDeLargo(proyecto string, largo int) error {
	nombre := strings.Repeat("a", largo)
	return w.modificar(proyecto, cambiosProyecto{Nombre: &nombre})
}

func (w *modificarProyectoWorld) modificoFechaInicio(proyecto, fecha string) error {
	f, err := parsearFecha(fecha)
	if err != nil {
		return err
	}
	return w.modificar(proyecto, cambiosProyecto{FechaInicio: &f})
}

func (w *modificarProyectoWorld) modificoFechaFin(proyecto, fecha string) error {
	f, err := parsearFecha(fecha)
	if err != nil {
		return err
	}
	return w.modificar(proyecto, cambiosProyecto{FechaFin: &f})
}

func (w *modificarProyectoWorld) modificoQuitandoInicio(proyecto string) error {
	return w.modificar(proyecto, cambiosProyecto{QuitarInicio: true})
}

func (w *modificarProyectoWorld) modificoQuitandoFin(proyecto string) error {
	return w.modificar(proyecto, cambiosProyecto{QuitarFin: true})
}

func (w *modificarProyectoWorld) modificoEstado(proyecto, estado string) error {
	return w.modificar(proyecto, cambiosProyecto{Estado: &estado})
}

func (w *modificarProyectoWorld) modificoSinDatos(proyecto string) error {
	return w.modificar(proyecto, cambiosProyecto{})
}

func (w *modificarProyectoWorld) modificoPorIDInexistente(id int64, nombre string) error {
	w.objetivoID = id
	w.ultimo, w.ultimoErr = w.sistema.Modificar(id, cambiosProyecto{Nombre: &nombre})
	return nil
}

// modificar aplica los cambios al proyecto preparado con ese nombre y guarda el resultado.
// Un error de la modificación no hace fallar el step: lo verifican los Then.
func (w *modificarProyectoWorld) modificar(proyecto string, c cambiosProyecto) error {
	p, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	w.objetivoID = p.ID
	w.ultimo, w.ultimoErr = w.sistema.Modificar(p.ID, c)
	return nil
}

// --- Then ---

func (w *modificarProyectoWorld) seRealizaCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba modificar el proyecto, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *modificarProyectoWorld) modificadoSeLlama(nombre string) error {
	if err := w.seRealizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Nombre != nombre {
		return fmt.Errorf("nombre: se esperaba %q, se obtuvo %q", nombre, w.ultimo.Nombre)
	}
	return nil
}

func (w *modificarProyectoWorld) modificadoTieneFechas(inicio, fin string) error {
	if err := w.seRealizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Inicio != inicio || w.ultimo.Fin != fin {
		return fmt.Errorf("fechas: se esperaba inicio %q y fin %q, se obtuvo inicio %q y fin %q",
			inicio, fin, w.ultimo.Inicio, w.ultimo.Fin)
	}
	return nil
}

func (w *modificarProyectoWorld) modificadoSinFechaFin() error {
	if err := w.seRealizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Fin != "" {
		return fmt.Errorf("se esperaba un proyecto sin fecha de fin, tiene %q", w.ultimo.Fin)
	}
	return nil
}

func (w *modificarProyectoWorld) modificadoEnEstado(estado string) error {
	if err := w.seRealizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Estado != estado {
		return fmt.Errorf("estado: se esperaba %q, se obtuvo %q", estado, w.ultimo.Estado)
	}
	return nil
}

func (w *modificarProyectoWorld) seRechazaModificacion(descripcion string) error {
	esperado, ok := erroresHU02()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la modificación se realizó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *modificarProyectoWorld) noCambio(nombre string) error {
	antes, err := w.proyecto(nombre)
	if err != nil {
		return err
	}
	ahora, err := w.sistema.BuscarPorID(antes.ID)
	if err != nil {
		return fmt.Errorf("buscar %q por ID %d: %w", nombre, antes.ID, err)
	}
	if ahora != antes {
		return fmt.Errorf("el proyecto %q cambió: antes %+v, ahora %+v", nombre, antes, ahora)
	}
	return nil
}

func (w *modificarProyectoWorld) conservaIdentificador() error {
	if err := w.seRealizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.ID != w.objetivoID {
		return fmt.Errorf("identificador: se esperaba %d, se obtuvo %d", w.objetivoID, w.ultimo.ID)
	}
	return nil
}

func (w *modificarProyectoWorld) alRecuperarloSeLlama(nombre string) error {
	p, err := w.sistema.BuscarPorID(w.objetivoID)
	if err != nil {
		return fmt.Errorf("buscar por ID %d: %w", w.objetivoID, err)
	}
	if p.Nombre != nombre {
		return fmt.Errorf("al recuperarlo, nombre: se esperaba %q, se obtuvo %q", nombre, p.Nombre)
	}
	return nil
}

func (w *modificarProyectoWorld) proyecto(nombre string) (proyectoHU02, error) {
	p, ok := w.porNombre[nombre]
	if !ok {
		return proyectoHU02{}, fmt.Errorf("el escenario no preparó el proyecto %q", nombre)
	}
	return p, nil
}
