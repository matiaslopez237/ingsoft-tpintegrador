package bdd

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
)

// datosElemento son las entradas de HU-05 que arman los escenarios.
type datosElemento struct {
	Titulo      string
	Descripcion string
	Prioridad   string
	StoryPoints *int
	Criterios   []string
}

// elementoVisto es lo que los escenarios de HU-05 observan de un elemento creado.
type elementoVisto struct {
	ID          int64
	Titulo      string
	Estado      string
	Estimado    bool
	StoryPoints int
	Criterios   []string
}

// sistemaBacklog es la puerta de entrada de los escenarios al núcleo.
// La implementación real está en hu05_sistema_test.go.
type sistemaBacklog interface {
	CrearProyecto(nombre string) (int64, error)
	CrearElemento(proyectoID int64, d datosElemento) (elementoVisto, error)
	BuscarElemento(id int64) (elementoVisto, error)
	ContarElementos(proyectoID int64) (int, error)
	Cerrar() error
}

// idProyectoInexistente no corresponde a ningún proyecto de la base del escenario.
const idProyectoInexistente int64 = 999999

type backlogWorld struct {
	sistema   sistemaBacklog
	proyectos map[string]int64
	elementos map[string]elementoVisto
	ultimo    elementoVisto
	ultimoErr error
}

func registrarStepsHU05(ctx *godog.ScenarioContext) {
	w := &backlogWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaBacklog()
		if err != nil {
			return c, err
		}
		*w = backlogWorld{sistema: sistema, proyectos: map[string]int64{}, elementos: map[string]elementoVisto{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^existe el proyecto "([^"]*)"$`, w.existeElProyecto)
	ctx.Step(`^en "([^"]*)" existe el elemento "([^"]*)"$`, w.existeElElemento)

	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" con prioridad "([^"]*)", (\d+) story points y los criterios:$`, w.creoConCriterios)
	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" con prioridad "([^"]*)" y sin story points$`, w.creoSinStoryPoints)
	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" con prioridad "([^"]*)" y (-?\d+) story points$`, w.creoConStoryPoints)
	ctx.Step(`^creo en un proyecto inexistente el elemento "([^"]*)" con prioridad "([^"]*)" y sin story points$`, w.creoEnProyectoInexistente)
	ctx.Step(`^creo en "([^"]*)" un elemento con (.+) y prioridad "([^"]*)"$`, w.creoConTituloDescripto)
	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" con la descripción "([^"]*)"$`, w.creoConDescripcion)
	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" sin criterios de aceptación$`, w.creoSinCriterios)

	ctx.Step(`^el elemento se crea correctamente$`, w.seCreaCorrectamente)
	ctx.Step(`^el elemento queda en estado "([^"]*)"$`, w.quedaEnEstado)
	ctx.Step(`^el elemento tiene un identificador asignado$`, w.tieneIdentificador)
	ctx.Step(`^el elemento tiene (\d+) story points$`, w.tieneStoryPoints)
	ctx.Step(`^el elemento queda sin estimar$`, w.quedaSinEstimar)
	ctx.Step(`^el elemento tiene (\d+) criterios de aceptación$`, w.tieneCriterios)
	ctx.Step(`^se rechaza el elemento con el error "([^"]*)"$`, w.seRechazaConError)
	ctx.Step(`^el proyecto "([^"]*)" tiene (\d+) elementos en el backlog$`, w.proyectoTieneElementos)
	ctx.Step(`^los elementos "([^"]*)" y "([^"]*)" tienen identificadores distintos$`, w.identificadoresDistintos)
	ctx.Step(`^cada elemento se recupera por su identificador con sus criterios en orden$`, w.recuperablesPorID)
}

// datosValidos completa con valores válidos lo que el escenario no especifica.
func datosValidos(titulo, prioridad string) datosElemento {
	return datosElemento{
		Titulo:      titulo,
		Descripcion: "Necesidad de ejemplo del producto",
		Prioridad:   prioridad,
		Criterios:   []string{"Criterio de ejemplo"},
	}
}

// --- Given ---

func (w *backlogWorld) existeElProyecto(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *backlogWorld) existeElElemento(proyecto, titulo string) error {
	if err := w.crear(proyecto, datosValidos(titulo, "Media")); err != nil {
		return err
	}
	if w.ultimoErr != nil {
		return fmt.Errorf("preparar el elemento %q: %w", titulo, w.ultimoErr)
	}
	return nil
}

// --- When ---

func (w *backlogWorld) creoConCriterios(proyecto, titulo, prioridad string, puntos int, tabla *godog.Table) error {
	d := datosValidos(titulo, prioridad)
	d.StoryPoints = &puntos
	d.Criterios = nil
	for _, fila := range tabla.Rows[1:] { // la primera fila es el encabezado
		criterio := fila.Cells[0].Value
		if criterio == "<vacío>" {
			criterio = "   "
		}
		d.Criterios = append(d.Criterios, criterio)
	}
	return w.crear(proyecto, d)
}

func (w *backlogWorld) creoSinStoryPoints(proyecto, titulo, prioridad string) error {
	return w.crear(proyecto, datosValidos(titulo, prioridad))
}

func (w *backlogWorld) creoConStoryPoints(proyecto, titulo, prioridad string, puntos int) error {
	d := datosValidos(titulo, prioridad)
	d.StoryPoints = &puntos
	return w.crear(proyecto, d)
}

func (w *backlogWorld) creoEnProyectoInexistente(titulo, prioridad string) error {
	w.registrarResultado(w.sistema.CrearElemento(idProyectoInexistente, datosValidos(titulo, prioridad)))
	return nil
}

// creoConTituloDescripto arma el título a partir de su descripción en el escenario.
func (w *backlogWorld) creoConTituloDescripto(proyecto, descripcion, prioridad string) error {
	var titulo string
	switch {
	case descripcion == "un título vacío":
		titulo = ""
	case descripcion == "un título de espacios":
		titulo = "     "
	default:
		var letras int
		if _, err := fmt.Sscanf(descripcion, "un título de %d letras", &letras); err != nil {
			return fmt.Errorf("título desconocido en el escenario: %q", descripcion)
		}
		titulo = strings.Repeat("a", letras)
	}
	return w.crear(proyecto, datosValidos(titulo, prioridad))
}

func (w *backlogWorld) creoConDescripcion(proyecto, titulo, descripcion string) error {
	d := datosValidos(titulo, "Alta")
	d.Descripcion = descripcion
	return w.crear(proyecto, d)
}

func (w *backlogWorld) creoSinCriterios(proyecto, titulo string) error {
	d := datosValidos(titulo, "Alta")
	d.Criterios = nil
	return w.crear(proyecto, d)
}

func (w *backlogWorld) crear(proyecto string, d datosElemento) error {
	id, ok := w.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no creó el proyecto %q", proyecto)
	}
	w.registrarResultado(w.sistema.CrearElemento(id, d))
	return nil
}

func (w *backlogWorld) registrarResultado(e elementoVisto, err error) {
	w.ultimo, w.ultimoErr = e, err
	if err == nil {
		w.elementos[e.Titulo] = e
	}
}

// --- Then ---

func (w *backlogWorld) seCreaCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba crear el elemento, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *backlogWorld) quedaEnEstado(estado string) error {
	if w.ultimo.Estado != estado {
		return fmt.Errorf("estado: se esperaba %q, se obtuvo %q", estado, w.ultimo.Estado)
	}
	return nil
}

func (w *backlogWorld) tieneIdentificador() error {
	if w.ultimo.ID <= 0 {
		return fmt.Errorf("se esperaba un identificador asignado, se obtuvo %d", w.ultimo.ID)
	}
	return nil
}

func (w *backlogWorld) tieneStoryPoints(puntos int) error {
	if err := w.seCreaCorrectamente(); err != nil {
		return err
	}
	if !w.ultimo.Estimado || w.ultimo.StoryPoints != puntos {
		return fmt.Errorf("story points: se esperaban %d, se obtuvo %+v", puntos, w.ultimo)
	}
	return nil
}

func (w *backlogWorld) quedaSinEstimar() error {
	if w.ultimo.Estimado {
		return fmt.Errorf("se esperaba un elemento sin estimar, tiene %d story points", w.ultimo.StoryPoints)
	}
	return nil
}

func (w *backlogWorld) tieneCriterios(cantidad int) error {
	if len(w.ultimo.Criterios) != cantidad {
		return fmt.Errorf("criterios: se esperaban %d, hay %d", cantidad, len(w.ultimo.Criterios))
	}
	return nil
}

func (w *backlogWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU05()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el elemento se creó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *backlogWorld) proyectoTieneElementos(proyecto string, esperado int) error {
	n, err := w.sistema.ContarElementos(w.proyectos[proyecto])
	if err != nil {
		return err
	}
	if n != esperado {
		return fmt.Errorf("elementos de %q: se esperaban %d, hay %d", proyecto, esperado, n)
	}
	return nil
}

func (w *backlogWorld) identificadoresDistintos(a, b string) error {
	ea, okA := w.elementos[a]
	eb, okB := w.elementos[b]
	if !okA || !okB {
		return fmt.Errorf("no se crearon los elementos %q y %q", a, b)
	}
	if ea.ID == eb.ID {
		return fmt.Errorf("ambos elementos tienen el identificador %d", ea.ID)
	}
	return nil
}

func (w *backlogWorld) recuperablesPorID() error {
	for titulo, e := range w.elementos {
		encontrado, err := w.sistema.BuscarElemento(e.ID)
		if err != nil {
			return fmt.Errorf("buscar %q por ID %d: %w", titulo, e.ID, err)
		}
		if !reflect.DeepEqual(encontrado, e) {
			return fmt.Errorf("se recuperó %+v, se esperaba %+v", encontrado, e)
		}
	}
	return nil
}
