package bdd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// proyectoVisto es lo que los escenarios de HU-01 necesitan observar de un proyecto creado.
type proyectoVisto struct {
	ID       int64
	Nombre   string
	Estado   string
	TieneFin bool
}

// sistemaProyectos es la puerta de entrada de los escenarios al núcleo.
// La implementación real está en hu01_sistema_test.go.
type sistemaProyectos interface {
	Crear(nombre string, inicio, fin *time.Time) (proyectoVisto, error)
	BuscarPorID(id int64) (proyectoVisto, error)
	Cantidad() (int, error)
	Cerrar() error
}

type crearProyectoWorld struct {
	sistema   sistemaProyectos
	porNombre map[string]proyectoVisto
	ultimo    proyectoVisto
	ultimoErr error
}

func registrarStepsHU01(ctx *godog.ScenarioContext) {
	w := &crearProyectoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaProyectos()
		if err != nil {
			return c, err
		}
		*w = crearProyectoWorld{sistema: sistema, porNombre: map[string]proyectoVisto{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^existe un proyecto "([^"]*)"$`, w.existeUnProyecto)
	ctx.Step(`^creo un proyecto "([^"]*)" con inicio "([^"]*)" y fin "([^"]*)"$`, w.creoConInicioYFin)
	ctx.Step(`^creo un proyecto "([^"]*)" con inicio "([^"]*)" y sin fecha de fin$`, w.creoSinFin)
	ctx.Step(`^creo un proyecto "([^"]*)" sin fecha de inicio$`, w.creoSinInicio)
	ctx.Step(`^creo un proyecto con un nombre de (\d+) caracteres$`, w.creoConNombreDeLargo)

	ctx.Step(`^el proyecto se crea correctamente$`, w.seCreaCorrectamente)
	ctx.Step(`^el proyecto queda en estado "([^"]*)"$`, w.quedaEnEstado)
	ctx.Step(`^el proyecto tiene un identificador asignado$`, w.tieneIdentificador)
	ctx.Step(`^el proyecto no tiene fecha de fin$`, w.noTieneFechaFin)
	ctx.Step(`^el nombre del proyecto es "([^"]*)"$`, w.elNombreEs)
	ctx.Step(`^los proyectos "([^"]*)" y "([^"]*)" tienen identificadores distintos$`, w.identificadoresDistintos)
	ctx.Step(`^cada proyecto puede recuperarse por su identificador$`, w.recuperablesPorID)
	ctx.Step(`^se rechaza la creación con el error "([^"]*)"$`, w.seRechazaConError)
	ctx.Step(`^no se creó ningún proyecto$`, w.noSeCreoNinguno)
	ctx.Step(`^solo existe (\d+) proyecto$`, w.soloExisten)
}

// --- Given / When ---

func (w *crearProyectoWorld) existeUnProyecto(nombre string) error {
	if err := w.creoConInicioYFin(nombre, "2026-03-01", "2026-11-30"); err != nil {
		return err
	}
	if w.ultimoErr != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, w.ultimoErr)
	}
	return nil
}

func (w *crearProyectoWorld) creoConInicioYFin(nombre, inicio, fin string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	w.crear(nombre, &fi, &ff)
	return nil
}

func (w *crearProyectoWorld) creoSinFin(nombre, inicio string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	w.crear(nombre, &fi, nil)
	return nil
}

func (w *crearProyectoWorld) creoSinInicio(nombre string) error {
	ff, err := parsearFecha("2026-11-30")
	if err != nil {
		return err
	}
	w.crear(nombre, nil, &ff)
	return nil
}

func (w *crearProyectoWorld) creoConNombreDeLargo(largo int) error {
	return w.creoConInicioYFin(strings.Repeat("a", largo), "2026-03-01", "2026-11-30")
}

func (w *crearProyectoWorld) crear(nombre string, inicio, fin *time.Time) {
	w.ultimo, w.ultimoErr = w.sistema.Crear(nombre, inicio, fin)
	if w.ultimoErr == nil {
		w.porNombre[w.ultimo.Nombre] = w.ultimo
	}
}

// --- Then ---

func (w *crearProyectoWorld) seCreaCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba crear el proyecto, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *crearProyectoWorld) quedaEnEstado(estado string) error {
	if w.ultimo.Estado != estado {
		return fmt.Errorf("estado: se esperaba %q, se obtuvo %q", estado, w.ultimo.Estado)
	}
	return nil
}

func (w *crearProyectoWorld) tieneIdentificador() error {
	if w.ultimo.ID <= 0 {
		return fmt.Errorf("se esperaba un identificador asignado, se obtuvo %d", w.ultimo.ID)
	}
	return nil
}

func (w *crearProyectoWorld) noTieneFechaFin() error {
	if w.ultimo.TieneFin {
		return errors.New("se esperaba un proyecto sin fecha de fin")
	}
	return nil
}

func (w *crearProyectoWorld) elNombreEs(nombre string) error {
	if w.ultimo.Nombre != nombre {
		return fmt.Errorf("nombre: se esperaba %q, se obtuvo %q", nombre, w.ultimo.Nombre)
	}
	return nil
}

func (w *crearProyectoWorld) identificadoresDistintos(a, b string) error {
	pa, okA := w.porNombre[a]
	pb, okB := w.porNombre[b]
	if !okA || !okB {
		return fmt.Errorf("no se crearon los proyectos %q y %q", a, b)
	}
	if pa.ID == pb.ID {
		return fmt.Errorf("ambos proyectos tienen el identificador %d", pa.ID)
	}
	return nil
}

func (w *crearProyectoWorld) recuperablesPorID() error {
	for nombre, p := range w.porNombre {
		encontrado, err := w.sistema.BuscarPorID(p.ID)
		if err != nil {
			return fmt.Errorf("buscar %q por ID %d: %w", nombre, p.ID, err)
		}
		if encontrado != p {
			return fmt.Errorf("se recuperó %+v, se esperaba %+v", encontrado, p)
		}
	}
	return nil
}

func (w *crearProyectoWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU01()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el proyecto se creó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *crearProyectoWorld) noSeCreoNinguno() error {
	return w.soloExisten(0)
}

func (w *crearProyectoWorld) soloExisten(esperado int) error {
	n, err := w.sistema.Cantidad()
	if err != nil {
		return err
	}
	if n != esperado {
		return fmt.Errorf("cantidad de proyectos: se esperaba %d, hay %d", esperado, n)
	}
	return nil
}

func parsearFecha(s string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("fecha inválida en el escenario %q: %w", s, err)
	}
	return t, nil
}
