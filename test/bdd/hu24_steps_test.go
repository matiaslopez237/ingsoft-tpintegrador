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

// calidadVista es lo que los escenarios de HU-24 observan de la calidad de un Sprint.
type calidadVista struct {
	Numero               int
	PorcentajeCompletado float64
	HayPorcentaje        bool
	DefectosDetectados   int
	DefectosResueltos    int
}

// sistemaCalidadSprint es la puerta de entrada de los escenarios de HU-24 al núcleo. La implementa
// sistemaCalidadSprintReal, sobre SQLite (hu24_sistema_test.go).
type sistemaCalidadSprint interface {
	CrearProyecto(nombre string) (int64, error)
	CrearElemento(proyectoID int64) (int64, error)
	// CrearSprint prepara un sprint cerrado (con sus puntos planificados y completados) o abierto.
	// Es preparación del escenario: no pasa por HU-08 ni por HU-11, igual que en HU-22.
	CrearSprint(proyectoID int64, estado domain.EstadoSprint, planificados, completados int) (int64, error)
	// CrearDefecto inserta un defecto directo, con el Sprint de detección y el de resolución que
	// indique el escenario (0 = sin Sprint). No pasa por HU-20 ni por HU-21.
	CrearDefecto(proyectoID, elementoID int64, estado domain.EstadoDefecto, sprintDeteccionID, sprintResolucionID int64) error
	ConsultarCalidad(proyectoID int64) ([]calidadVista, error)
	Cerrar() error
}

// erroresHU24 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU24() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

type calidadWorld struct {
	sistema   sistemaCalidadSprint
	proyectos map[string]int64
	elementos map[string]int64
	sprints   map[string][]int64 // sprints[proyecto] = IDs en orden de creación (numero = índice+1)
	ultima    []calidadVista
	ultimoErr error
}

func registrarStepsHU24(ctx *godog.ScenarioContext) {
	w := &calidadWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaCalidadSprint()
		if err != nil {
			return c, err
		}
		*w = calidadWorld{sistema: sistema, proyectos: map[string]int64{}, elementos: map[string]int64{}, sprints: map[string][]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para la calidad del sprint$`, w.proyectoParaCalidad)
	ctx.Step(`^"([^"]*)" tiene (?:un|otro) sprint cerrado con (\d+) puntos planificados y (\d+) completados$`, w.tieneSprintCerrado)
	ctx.Step(`^"([^"]*)" tiene un sprint todavía abierto$`, w.tieneSprintAbierto)
	ctx.Step(`^se detect(?:ó|aron) (\d+) defectos? en el sprint (\d+) de "([^"]*)"$`, w.detectaronDefectos)
	ctx.Step(`^se detectó un defecto en "([^"]*)" sin ningún sprint$`, w.detectoDefectoSinSprint)
	ctx.Step(`^se resolvió un defecto en el sprint (\d+) de "([^"]*)", y después se cerró$`, w.resolvioYCerroDefecto)

	// When
	ctx.Step(`^consulto la calidad de los sprints de "([^"]*)"$`, w.consultoCalidad)
	ctx.Step(`^consulto la calidad de los sprints de un proyecto inexistente$`, w.consultoCalidadDeInexistente)

	// Then
	ctx.Step(`^el sprint (\d+) tiene (-?\d+(?:,\d+)?)% completado$`, w.sprintTienePorcentaje)
	ctx.Step(`^el sprint (\d+) no tiene datos de porcentaje$`, w.sprintSinDatosDePorcentaje)
	ctx.Step(`^el sprint (\d+) tiene (\d+) defectos detectados$`, w.sprintTieneDefectosDetectados)
	ctx.Step(`^el sprint (\d+) tiene (\d+) defectos resueltos$`, w.sprintTieneDefectosResueltos)
	ctx.Step(`^solo aparece (\d+) sprint en la calidad$`, w.soloAparecenNSprints)
	ctx.Step(`^se rechaza la consulta de calidad con el error "([^"]*)"$`, w.seRechazaConsultaCalidad)
}

// --- Given ---

func (w *calidadWorld) proyectoParaCalidad(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *calidadWorld) tieneSprintCerrado(proyecto string, planificados, completados int) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.sistema.CrearSprint(id, domain.SprintCerrado, planificados, completados)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint cerrado del proyecto %q: %w", proyecto, err)
	}
	w.sprints[proyecto] = append(w.sprints[proyecto], sid)
	return nil
}

func (w *calidadWorld) tieneSprintAbierto(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.sistema.CrearSprint(id, domain.SprintAbierto, 0, 0)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint abierto del proyecto %q: %w", proyecto, err)
	}
	w.sprints[proyecto] = append(w.sprints[proyecto], sid)
	return nil
}

func (w *calidadWorld) detectaronDefectos(cantidad, numero int, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(proyecto, numero)
	if err != nil {
		return err
	}
	eid, err := w.idElemento(proyecto, pid)
	if err != nil {
		return err
	}
	for i := 0; i < cantidad; i++ {
		if err := w.sistema.CrearDefecto(pid, eid, domain.DefectoAbierto, sid, 0); err != nil {
			return fmt.Errorf("no se pudo preparar un defecto detectado en el sprint %d de %q: %w", numero, proyecto, err)
		}
	}
	return nil
}

func (w *calidadWorld) detectoDefectoSinSprint(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	eid, err := w.idElemento(proyecto, pid)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearDefecto(pid, eid, domain.DefectoAbierto, 0, 0); err != nil {
		return fmt.Errorf("no se pudo preparar un defecto sin sprint en %q: %w", proyecto, err)
	}
	return nil
}

func (w *calidadWorld) resolvioYCerroDefecto(numero int, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(proyecto, numero)
	if err != nil {
		return err
	}
	eid, err := w.idElemento(proyecto, pid)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearDefecto(pid, eid, domain.DefectoCerrado, 0, sid); err != nil {
		return fmt.Errorf("no se pudo preparar un defecto resuelto en el sprint %d de %q: %w", numero, proyecto, err)
	}
	return nil
}

// --- When ---

func (w *calidadWorld) consultoCalidad(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultima, w.ultimoErr = w.sistema.ConsultarCalidad(id)
	return nil
}

func (w *calidadWorld) consultoCalidadDeInexistente() error {
	w.ultima, w.ultimoErr = w.sistema.ConsultarCalidad(idInexistente)
	return nil
}

// --- Then ---

func (w *calidadWorld) sprintTienePorcentaje(numero int, porcentaje string) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la calidad, pero hubo error: %w", w.ultimoErr)
	}
	v, err := w.porSprintNumero(numero)
	if err != nil {
		return err
	}
	if !v.HayPorcentaje {
		return fmt.Errorf("se esperaba porcentaje con datos en el sprint %d, pero no hay", numero)
	}
	want, err := strconv.ParseFloat(strings.ReplaceAll(porcentaje, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("porcentaje inválido en el escenario %q: %w", porcentaje, err)
	}
	if v.PorcentajeCompletado != want {
		return fmt.Errorf("el sprint %d tiene %v%% completado, se esperaba %v%%", numero, v.PorcentajeCompletado, want)
	}
	return nil
}

func (w *calidadWorld) sprintSinDatosDePorcentaje(numero int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la calidad, pero hubo error: %w", w.ultimoErr)
	}
	v, err := w.porSprintNumero(numero)
	if err != nil {
		return err
	}
	if v.HayPorcentaje {
		return fmt.Errorf("se esperaba que el sprint %d no tuviera datos de porcentaje, pero tiene %v%%", numero, v.PorcentajeCompletado)
	}
	return nil
}

func (w *calidadWorld) sprintTieneDefectosDetectados(numero, cantidad int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la calidad, pero hubo error: %w", w.ultimoErr)
	}
	v, err := w.porSprintNumero(numero)
	if err != nil {
		return err
	}
	if v.DefectosDetectados != cantidad {
		return fmt.Errorf("el sprint %d tiene %d defectos detectados, se esperaban %d", numero, v.DefectosDetectados, cantidad)
	}
	return nil
}

func (w *calidadWorld) sprintTieneDefectosResueltos(numero, cantidad int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la calidad, pero hubo error: %w", w.ultimoErr)
	}
	v, err := w.porSprintNumero(numero)
	if err != nil {
		return err
	}
	if v.DefectosResueltos != cantidad {
		return fmt.Errorf("el sprint %d tiene %d defectos resueltos, se esperaban %d", numero, v.DefectosResueltos, cantidad)
	}
	return nil
}

func (w *calidadWorld) soloAparecenNSprints(cantidad int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar la calidad, pero hubo error: %w", w.ultimoErr)
	}
	if len(w.ultima) != cantidad {
		return fmt.Errorf("se esperaban %d sprints en la calidad, hay %d", cantidad, len(w.ultima))
	}
	return nil
}

func (w *calidadWorld) seRechazaConsultaCalidad(descripcion string) error {
	esperado, ok := erroresHU24()[descripcion]
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

func (w *calidadWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *calidadWorld) idSprint(proyecto string, numero int) (int64, error) {
	ids, ok := w.sprints[proyecto]
	if !ok || numero < 1 || numero > len(ids) {
		return 0, fmt.Errorf("el sprint %d de %q no se preparó en el escenario", numero, proyecto)
	}
	return ids[numero-1], nil
}

// idElemento crea, la primera vez que hace falta, el único elemento de prueba del proyecto (lo
// reutilizan todos los defectos del escenario: su título no importa para HU-24).
func (w *calidadWorld) idElemento(proyecto string, proyectoID int64) (int64, error) {
	if id, ok := w.elementos[proyecto]; ok {
		return id, nil
	}
	id, err := w.sistema.CrearElemento(proyectoID)
	if err != nil {
		return 0, fmt.Errorf("no se pudo preparar el elemento de %q: %w", proyecto, err)
	}
	w.elementos[proyecto] = id
	return id, nil
}

func (w *calidadWorld) porSprintNumero(numero int) (calidadVista, error) {
	for _, v := range w.ultima {
		if v.Numero == numero {
			return v, nil
		}
	}
	return calidadVista{}, fmt.Errorf("no se encontró el sprint %d entre los resultados de calidad", numero)
}
