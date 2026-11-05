package bdd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu23Sistema es la puerta de entrada de los escenarios de HU-23 al núcleo. La implementa
// hu23SistemaReal, sobre SQLite (hu23_sistema_test.go). Todos los identificadores de paquete llevan el
// prefijo hu23 para no chocar con los de otras historias.
type hu23Sistema interface {
	CrearProyecto(nombre string) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre string) error
	QuitarIntegrante(proyectoID int64, nombre string) error
	// CambiarEstadoProyecto deja un proyecto preparado en otro estado (con HU-02).
	CambiarEstadoProyecto(proyectoID int64, estado string) error
	// CrearSprint deja un sprint Abierto; numero (1, 2, 3…) define sus fechas, que no se superponen.
	CrearSprint(proyectoID int64, numero int) (int64, error)
	// CrearElemento crea el elemento con HU-05 y, si horasEstimadas no es nil, se las carga con HU-06.
	CrearElemento(proyectoID int64, titulo string, horasEstimadas *float64) (int64, error)
	// RegistrarEsfuerzo registra horas con el caso de uso de HU-18.
	RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64, fecha time.Time) error
	EditarHorasEstimadas(proyectoID, elementoID int64, horas float64) error
	EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error
	PasarAEnProgreso(proyectoID, elementoID int64) error
	// Asignar usa HU-09, Completar usa HU-10 y CerrarSprint usa HU-11.
	Asignar(proyectoID, elementoID, sprintID int64) error
	Completar(proyectoID, elementoID int64) error
	CerrarSprint(proyectoID, sprintID int64) error
	HorasDeProyecto(proyectoID int64) (domain.TotalHoras, error)
	HorasDeSprint(proyectoID, sprintID int64) (domain.TotalHoras, error)
	Elemento(elementoID int64) (domain.ElementoBacklog, error)
	Sprint(sprintID int64) (domain.Sprint, error)
	Cerrar() error
}

// hu23Errores traduce la descripción del escenario al error definido en la spec (sección 7).
func hu23Errores() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"sprint no encontrado":   domain.ErrSprintNoEncontrado,
	}
}

// hu23SinEstimar es lo que dice la columna "estimadas" de un elemento sin horas estimadas.
const hu23SinEstimar = "sin estimar"

// hu23Elemento y hu23Sprint son lo que preparó el escenario. Los nombres no se repiten en un escenario.
type hu23Elemento struct {
	proyectoID, id int64
}

type hu23Sprint struct {
	proyectoID, id int64
}

type hu23World struct {
	sistema     hu23Sistema
	proyectos   map[string]int64
	integrante  map[int64]string // proyecto → la integrante del Background, que carga las reales de las tablas
	sprints     map[string]hu23Sprint
	numeros     map[int64]int // proyecto → cuántos sprints tiene, para numerar el próximo
	elementos   map[string]hu23Elemento
	total       domain.TotalHoras
	errConsulta error // el de la última consulta
	consulto    bool  // para no dar por buena una consulta que el escenario nunca hizo
}

func registrarStepsHU23(ctx *godog.ScenarioContext) {
	w := &hu23World{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = hu23World{} // si falla la preparación, el After no cierra el sistema del escenario anterior
		sistema, err := hu23NuevoSistema()
		if err != nil {
			return c, err
		}
		*w = hu23World{
			sistema:    sistema,
			proyectos:  map[string]int64{},
			integrante: map[int64]string{},
			sprints:    map[string]hu23Sprint{},
			numeros:    map[int64]int{},
			elementos:  map[string]hu23Elemento{},
		}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	const horas = `(-?\d+(?:,\d+)?)` // con coma decimal, como en la spec

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para totalizar horas de HU-23, con la integrante "([^"]*)"$`, w.proyectoConIntegrante)
	ctx.Step(`^el proyecto "([^"]*)" tiene abierto el sprint "([^"]*)" de HU-23$`, w.sprintAbierto)
	ctx.Step(`^el sprint "([^"]*)" de HU-23 tiene estos elementos:$`, w.elementosEnSprint)
	ctx.Step(`^el proyecto "([^"]*)" de HU-23 tiene estos elementos sin sprint:$`, w.elementosSinSprint)
	ctx.Step(`^"([^"]*)" de HU-23 se completa$`, w.seCompleta)
	ctx.Step(`^"([^"]*)" de HU-23 pasa a En progreso$`, w.pasaAEnProgreso)
	ctx.Step(`^"([^"]*)" de HU-23 tiene (\d+) story points$`, w.tienePuntos)
	ctx.Step(`^"([^"]*)" de HU-23 pasa a tener `+horas+` horas estimadas$`, w.pasaATenerHoras)
	ctx.Step(`^"([^"]*)" de HU-23 se planifica en el sprint "([^"]*)"$`, w.sePlanifica)
	ctx.Step(`^el sprint "([^"]*)" de HU-23 se cierra$`, w.seCierra)
	ctx.Step(`^"([^"]*)" registró `+horas+` horas para HU-23 en "([^"]*)"$`, w.registro)
	ctx.Step(`^"([^"]*)" registró `+horas+` horas para HU-23 en "([^"]*)" con fecha fuera del sprint$`, w.registroFueraDelSprint)
	ctx.Step(`^"([^"]*)" también integra el proyecto "([^"]*)" de HU-23$`, w.tambienIntegra)
	ctx.Step(`^"([^"]*)" dejó de integrar el proyecto "([^"]*)" de HU-23$`, w.dejoDeIntegrar)
	ctx.Step(`^el proyecto "([^"]*)" de HU-23 pasa a estar "([^"]*)"$`, w.proyectoPasaAEstar)

	// When
	ctx.Step(`^consulto las horas de HU-23 del sprint "([^"]*)"$`, w.consultoSprint)
	ctx.Step(`^consulto las horas de HU-23 del proyecto "([^"]*)"$`, w.consultoProyecto)
	ctx.Step(`^consulto las horas de HU-23 de un proyecto inexistente$`, w.consultoProyectoInexistente)
	ctx.Step(`^consulto las horas de HU-23 del sprint "([^"]*)" desde un proyecto inexistente$`, w.consultoSprintDesdeProyectoInexistente)
	ctx.Step(`^consulto las horas de HU-23 de un sprint inexistente de "([^"]*)"$`, w.consultoSprintInexistente)
	ctx.Step(`^consulto las horas de HU-23 del sprint "([^"]*)" desde el proyecto "([^"]*)"$`, w.consultoSprintDesdeProyecto)

	// Then
	ctx.Step(`^las horas de HU-23 informan `+horas+` estimadas, `+horas+` reales y `+horas+` reales sin estimar$`, w.informan)
	ctx.Step(`^la desviación de HU-23 es de `+horas+` horas y de `+horas+` %$`, w.desviacion)
	ctx.Step(`^las horas de HU-23 cuentan (\d+) elementos? y (\d+) sin estimar$`, w.cuentan)
	ctx.Step(`^las horas de HU-23 están sin estimar y sin desviaciones$`, w.sinEstimarYSinDesviaciones)
	ctx.Step(`^se rechaza la consulta de horas de HU-23 con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^el sprint "([^"]*)" de HU-23 conserva (\d+) puntos planificados y (\d+) completados$`, w.conservaPuntos)
	ctx.Step(`^"([^"]*)" de HU-23 sigue Pendiente en el sprint "([^"]*)"$`, w.siguePendienteEnSprint)
	ctx.Step(`^el sprint "([^"]*)" de HU-23 sigue Abierto$`, w.sigueAbierto)
}

// --- Given ---

func (w *hu23World) proyectoConIntegrante(proyecto, integrante string) error {
	pid, err := w.sistema.CrearProyecto(proyecto)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
	}
	if err := w.sistema.AgregarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo preparar a %q en %q: %w", integrante, proyecto, err)
	}
	w.proyectos[proyecto] = pid
	w.integrante[pid] = integrante
	return nil
}

func (w *hu23World) sprintAbierto(proyecto, sprint string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.numeros[pid]++
	id, err := w.sistema.CrearSprint(pid, w.numeros[pid])
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint %q de %q: %w", sprint, proyecto, err)
	}
	w.sprints[sprint] = hu23Sprint{proyectoID: pid, id: id}
	return nil
}

func (w *hu23World) elementosEnSprint(sprint string, tabla *godog.Table) error {
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	return w.prepararElementos(s.proyectoID, tabla, func(elementoID int64) error {
		return w.sistema.Asignar(s.proyectoID, elementoID, s.id)
	})
}

func (w *hu23World) elementosSinSprint(proyecto string, tabla *godog.Table) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	return w.prepararElementos(pid, tabla, func(int64) error { return nil })
}

// prepararElementos crea cada fila de la tabla (elemento, estimadas, reales), la ubica con ubicar y
// le registra las reales a nombre de la integrante del proyecto.
func (w *hu23World) prepararElementos(proyectoID int64, tabla *godog.Table, ubicar func(elementoID int64) error) error {
	if len(tabla.Rows) < 2 || len(tabla.Rows[0].Cells) != 3 {
		return errors.New("la tabla de elementos tiene que tener elemento, estimadas y reales")
	}
	for _, fila := range tabla.Rows[1:] {
		titulo, estimadas, reales := fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value
		var horasEstimadas *float64
		if estimadas != hu23SinEstimar {
			h, err := hu23Horas(estimadas)
			if err != nil {
				return err
			}
			horasEstimadas = &h
		}
		id, err := w.sistema.CrearElemento(proyectoID, titulo, horasEstimadas)
		if err != nil {
			return fmt.Errorf("no se pudo preparar el elemento %q: %w", titulo, err)
		}
		w.elementos[titulo] = hu23Elemento{proyectoID: proyectoID, id: id}
		if err := ubicar(id); err != nil {
			return fmt.Errorf("no se pudo ubicar el elemento %q: %w", titulo, err)
		}
		if err := w.acumular(proyectoID, id, reales); err != nil {
			return fmt.Errorf("no se pudo cargar el esfuerzo de %q: %w", titulo, err)
		}
	}
	return nil
}

// acumular reparte las horas en registros de hasta 24 h (RN-18.8). Con 0 horas no registra nada: un
// registro de 0 h no es válido (RN-18.6).
func (w *hu23World) acumular(proyectoID, elementoID int64, horas string) error {
	resto, err := hu23Horas(horas)
	if err != nil {
		return err
	}
	for resto > 0 {
		parte := min(resto, domain.HorasMaximasPorRegistro)
		if err := w.sistema.RegistrarEsfuerzo(proyectoID, elementoID, w.integrante[proyectoID], parte, hu23Hoy); err != nil {
			return err
		}
		resto -= parte
	}
	return nil
}

func (w *hu23World) seCompleta(titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if err := w.sistema.Completar(e.proyectoID, e.id); err != nil {
		return fmt.Errorf("no se pudo completar %q: %w", titulo, err)
	}
	return nil
}

func (w *hu23World) pasaAEnProgreso(titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if err := w.sistema.PasarAEnProgreso(e.proyectoID, e.id); err != nil {
		return fmt.Errorf("no se pudo pasar %q a En progreso: %w", titulo, err)
	}
	return nil
}

func (w *hu23World) tienePuntos(titulo string, puntos int) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if err := w.sistema.EditarStoryPoints(e.proyectoID, e.id, puntos); err != nil {
		return fmt.Errorf("no se pudieron cargar %d story points en %q: %w", puntos, titulo, err)
	}
	return nil
}

func (w *hu23World) pasaATenerHoras(titulo, horas string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	h, err := hu23Horas(horas)
	if err != nil {
		return err
	}
	if err := w.sistema.EditarHorasEstimadas(e.proyectoID, e.id, h); err != nil {
		return fmt.Errorf("no se pudieron cargar %s horas estimadas en %q: %w", horas, titulo, err)
	}
	return nil
}

func (w *hu23World) sePlanifica(titulo, sprint string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	if err := w.sistema.Asignar(e.proyectoID, e.id, s.id); err != nil {
		return fmt.Errorf("no se pudo planificar %q en %q: %w", titulo, sprint, err)
	}
	return nil
}

func (w *hu23World) seCierra(sprint string) error {
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	if err := w.sistema.CerrarSprint(s.proyectoID, s.id); err != nil {
		return fmt.Errorf("no se pudo cerrar %q: %w", sprint, err)
	}
	return nil
}

func (w *hu23World) registro(integrante, horas, titulo string) error {
	return w.registrar(integrante, horas, titulo, hu23Hoy)
}

func (w *hu23World) registroFueraDelSprint(integrante, horas, titulo string) error {
	return w.registrar(integrante, horas, titulo, hu23FechaFueraDelSprint)
}

func (w *hu23World) registrar(integrante, horas, titulo string, fecha time.Time) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	h, err := hu23Horas(horas)
	if err != nil {
		return err
	}
	if err := w.sistema.RegistrarEsfuerzo(e.proyectoID, e.id, integrante, h, fecha); err != nil {
		return fmt.Errorf("no se pudieron registrar %s horas de %q en %q: %w", horas, integrante, titulo, err)
	}
	return nil
}

func (w *hu23World) tambienIntegra(integrante, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.AgregarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo agregar a %q en %q: %w", integrante, proyecto, err)
	}
	return nil
}

func (w *hu23World) dejoDeIntegrar(integrante, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.QuitarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo quitar a %q de %q: %w", integrante, proyecto, err)
	}
	return nil
}

func (w *hu23World) proyectoPasaAEstar(proyecto, estado string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CambiarEstadoProyecto(pid, estado); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", proyecto, estado, err)
	}
	return nil
}

// --- When ---

func (w *hu23World) consultoSprint(sprint string) error {
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	w.consultarSprint(s.proyectoID, s.id)
	return nil
}

func (w *hu23World) consultoProyecto(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.consultarProyecto(pid)
	return nil
}

func (w *hu23World) consultoProyectoInexistente() error {
	w.consultarProyecto(idInexistente)
	return nil
}

func (w *hu23World) consultoSprintDesdeProyectoInexistente(sprint string) error {
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	w.consultarSprint(idInexistente, s.id)
	return nil
}

func (w *hu23World) consultoSprintInexistente(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.consultarSprint(pid, idInexistente)
	return nil
}

func (w *hu23World) consultoSprintDesdeProyecto(sprint, proyecto string) error {
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.consultarSprint(pid, s.id)
	return nil
}

func (w *hu23World) consultarSprint(proyectoID, sprintID int64) {
	w.total, w.errConsulta = w.sistema.HorasDeSprint(proyectoID, sprintID)
	w.consulto = true
}

func (w *hu23World) consultarProyecto(proyectoID int64) {
	w.total, w.errConsulta = w.sistema.HorasDeProyecto(proyectoID)
	w.consulto = true
}

// --- Then ---

func (w *hu23World) informan(estimadas, reales, realesSinEstimar string) error {
	t, err := w.resultado()
	if err != nil {
		return err
	}
	if err := hu23IgualValor("horas estimadas", t.HorasEstimadas, estimadas); err != nil {
		return err
	}
	if err := hu23IgualValor("horas reales", t.HorasReales, reales); err != nil {
		return err
	}
	return hu23IgualValor("horas reales sin estimar", t.HorasRealesSinEstimar, realesSinEstimar)
}

func (w *hu23World) desviacion(horas, porcentaje string) error {
	t, err := w.resultado()
	if err != nil {
		return err
	}
	if t.SinEstimar() {
		return fmt.Errorf("el total quedó sin estimar: %s", hu23Ver(t))
	}
	if err := hu23Igual("desviación en horas", t.DesviacionHoras, horas); err != nil {
		return err
	}
	return hu23Igual("desviación en porcentaje", t.DesviacionPorcentaje, porcentaje)
}

func (w *hu23World) cuentan(elementos, sinEstimar int) error {
	t, err := w.resultado()
	if err != nil {
		return err
	}
	if t.Elementos != elementos || t.ElementosSinEstimar != sinEstimar {
		return fmt.Errorf("se esperaban %d elementos y %d sin estimar: %s", elementos, sinEstimar, hu23Ver(t))
	}
	return nil
}

func (w *hu23World) sinEstimarYSinDesviaciones() error {
	t, err := w.resultado()
	if err != nil {
		return err
	}
	if !t.SinEstimar() {
		return fmt.Errorf("se esperaba sin estimar: %s", hu23Ver(t))
	}
	if t.DesviacionHoras != nil || t.DesviacionPorcentaje != nil {
		return fmt.Errorf("sin estimar no debería informar desviaciones: %s", hu23Ver(t))
	}
	return nil
}

func (w *hu23World) seRechaza(descripcion string) error {
	if !w.consulto {
		return errors.New("el escenario no hizo ninguna consulta")
	}
	esperado, ok := hu23Errores()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.errConsulta == nil {
		return fmt.Errorf("se esperaba el error %q, pero no hubo error", descripcion)
	}
	if !errors.Is(w.errConsulta, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.errConsulta)
	}
	return nil
}

func (w *hu23World) conservaPuntos(sprint string, planificados, completados int) error {
	got, err := w.leerSprint(sprint)
	if err != nil {
		return err
	}
	if got.PuntosPlanificados != planificados || got.PuntosCompletados != completados {
		return fmt.Errorf("%q: se esperaban %d planificados y %d completados, tiene %d y %d",
			sprint, planificados, completados, got.PuntosPlanificados, got.PuntosCompletados)
	}
	return nil
}

func (w *hu23World) siguePendienteEnSprint(titulo, sprint string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	s, err := w.sprint(sprint)
	if err != nil {
		return err
	}
	got, err := w.sistema.Elemento(e.id)
	if err != nil {
		return fmt.Errorf("no se pudo leer el elemento %q: %w", titulo, err)
	}
	if got.Estado != domain.EstadoPendiente || got.SprintID == nil || *got.SprintID != s.id {
		return fmt.Errorf("%q: se esperaba Pendiente en %q, está %q con sprint %v", titulo, sprint, got.Estado, got.SprintID)
	}
	return nil
}

func (w *hu23World) sigueAbierto(sprint string) error {
	got, err := w.leerSprint(sprint)
	if err != nil {
		return err
	}
	if got.Estado != domain.SprintAbierto {
		return fmt.Errorf("%q: se esperaba Abierto, está %q", sprint, got.Estado)
	}
	return nil
}

// --- Auxiliares ---

// resultado devuelve el total de la última consulta, o un error si no la hubo o falló.
func (w *hu23World) resultado() (domain.TotalHoras, error) {
	if !w.consulto {
		return domain.TotalHoras{}, errors.New("el escenario no hizo ninguna consulta")
	}
	if w.errConsulta != nil {
		return domain.TotalHoras{}, fmt.Errorf("se esperaba consultar, pero hubo error: %w", w.errConsulta)
	}
	return w.total, nil
}

func (w *hu23World) leerSprint(nombre string) (domain.Sprint, error) {
	s, err := w.sprint(nombre)
	if err != nil {
		return domain.Sprint{}, err
	}
	got, err := w.sistema.Sprint(s.id)
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("no se pudo leer el sprint %q: %w", nombre, err)
	}
	return got, nil
}

func (w *hu23World) elemento(titulo string) (hu23Elemento, error) {
	e, ok := w.elementos[titulo]
	if !ok {
		return hu23Elemento{}, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return e, nil
}

func (w *hu23World) sprint(nombre string) (hu23Sprint, error) {
	s, ok := w.sprints[nombre]
	if !ok {
		return hu23Sprint{}, fmt.Errorf("el sprint %q no se preparó en el escenario", nombre)
	}
	return s, nil
}

func (w *hu23World) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

// hu23Horas lee un número con coma decimal, como se escribe en el .feature.
func hu23Horas(s string) (float64, error) {
	h, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("número mal escrito en el escenario: %q", s)
	}
	return h, nil
}

// hu23Igual compara un valor opcional con el del escenario. La comparación es exacta a propósito: las
// horas son múltiplos de 0,5 y el porcentaje llega redondeado a un decimal (spec, sección 5).
func hu23Igual(campo string, got *float64, esperado string) error {
	if got == nil {
		return fmt.Errorf("%s: se esperaba %s, no tiene valor", campo, esperado)
	}
	return hu23IgualValor(campo, *got, esperado)
}

func hu23IgualValor(campo string, got float64, esperado string) error {
	e, err := hu23Horas(esperado)
	if err != nil {
		return err
	}
	if got != e {
		return fmt.Errorf("%s: se esperaba %s, es %v", campo, esperado, got)
	}
	return nil
}

// hu23Ver muestra un total legible en los mensajes: las desviaciones son punteros.
func hu23Ver(t domain.TotalHoras) string {
	ver := func(v *float64) string {
		if v == nil {
			return "sin valor"
		}
		return strconv.FormatFloat(*v, 'f', -1, 64)
	}
	return fmt.Sprintf("estimadas %v, reales %v (%v sin estimar), desviación %s h y %s %%, %d elementos (%d sin estimar)",
		t.HorasEstimadas, t.HorasReales, t.HorasRealesSinEstimar, ver(t.DesviacionHoras), ver(t.DesviacionPorcentaje),
		t.Elementos, t.ElementosSinEstimar)
}
