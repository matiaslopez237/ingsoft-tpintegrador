package bdd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// dashboardVista es lo que los escenarios de HU-25 observan del Dashboard.
type dashboardVista struct {
	Pendientes, EnProgreso, Hechos, Total int

	HaySprintEnCurso                                  bool
	NumeroEnCurso, PlanificadosEnCurso, HechosEnCurso int

	Cerrados []cerradoVista // del más antiguo al más reciente (RN-25.5)

	Velocidad    float64
	HayVelocidad bool

	HorasEstimadas, HorasReales float64
	Desviacion                  *float64 // nil: ningún elemento estimado en horas (RN-23.6)
}

// cerradoVista es un sprint cerrado en el Dashboard.
type cerradoVista struct {
	Numero, Planificados, Completados int
	Porcentaje                        float64
	HayPorcentaje                     bool
	Detectados, Resueltos             int
}

// sistemaDashboard es la puerta de entrada de los escenarios de HU-25 al núcleo. La implementa
// sistemaDashboardReal, sobre SQLite (hu25_sistema_test.go).
type sistemaDashboard interface {
	CrearProyecto(nombre string) (int64, error)
	// FijarHoy fija la fecha que el caso de uso usa como "hoy" (RN-25.3).
	FijarHoy(hoy time.Time)
	// CrearElemento prepara un elemento con su estado, sus story points (nil: sin estimar) y su
	// sprint (0: sin sprint). No pasa por HU-05, HU-09 ni HU-10.
	CrearElemento(proyectoID int64, estado domain.EstadoElemento, puntos *int, sprintID int64) (int64, error)
	// CrearSprint inserta un sprint con sus fechas y su resumen de cierre, sin pasar por HU-08 ni HU-11.
	CrearSprint(proyectoID int64, estado domain.EstadoSprint, inicio, fin time.Time, planificados, completados int) (int64, error)
	// CrearDefecto inserta un defecto con su sprint de detección y de resolución (0: sin sprint).
	CrearDefecto(proyectoID, elementoID int64, estado domain.EstadoDefecto, deteccionID, resolucionID int64) error
	// CrearElementoConHoras prepara un elemento con horas estimadas (nil: sin estimar) y un registro
	// de esfuerzo con las horas reales. No pasa por HU-18 ni HU-19.
	CrearElementoConHoras(proyectoID int64, estimadas *float64, reales float64) error
	// CompletarHistoria marca la historia como Hecho con el caso de uso de HU-10.
	CompletarHistoria(proyectoID, elementoID int64) error
	ConsultarDashboard(proyectoID int64) (dashboardVista, error)
	Cerrar() error
}

// erroresHU25 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU25() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

type dashboardWorld struct {
	sistema   sistemaDashboard
	proyectos map[string]int64
	sprints   map[string][]int64 // IDs en orden de creación: el sprint N es el índice N-1
	enCurso   map[string]int64   // sprint abierto de cada proyecto
	historias map[string]map[int]int64
	ultima    dashboardVista
	ultimoErr error
}

func registrarStepsHU25(ctx *godog.ScenarioContext) {
	w := &dashboardWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = dashboardWorld{} // el sistema se crea recién en el primer step de HU-25
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para el dashboard$`, w.proyectoParaDashboard)
	ctx.Step(`^hoy es "([^"]*)" para el dashboard$`, w.hoyEs)
	ctx.Step(`^"([^"]*)" tiene (\d+) elementos? "([^"]*)" en el backlog$`, w.tieneElementos)
	ctx.Step(`^"([^"]*)" tiene un sprint abierto del "([^"]*)" al "([^"]*)"$`, w.tieneSprintAbierto)
	ctx.Step(`^"([^"]*)" tiene un sprint cerrado del "([^"]*)" al "([^"]*)" con (\d+) puntos planificados y (\d+) completados$`, w.tieneSprintCerrado)
	ctx.Step(`^el sprint en curso de "([^"]*)" tiene una historia "([^"]*)" de (\d+) puntos$`, w.historiaEnCurso)
	ctx.Step(`^el sprint en curso de "([^"]*)" tiene una historia "([^"]*)" sin estimar$`, w.historiaEnCursoSinEstimar)
	ctx.Step(`^se detectaron (\d+) defectos en el sprint (\d+) de "([^"]*)", y 1 se resolvió en el sprint (\d+)$`, w.defectosDetectadosYResuelto)
	ctx.Step(`^"([^"]*)" tiene un elemento estimado en (\d+) horas con (\d+) horas reales registradas$`, w.elementoConHoras)
	ctx.Step(`^"([^"]*)" tiene un elemento sin estimar en horas con (\d+) horas reales registradas$`, w.elementoSinEstimarConHoras)

	// When
	ctx.Step(`^consulto el dashboard de "([^"]*)"$`, w.consultoDashboard)
	ctx.Step(`^consulto el dashboard de un proyecto inexistente$`, w.consultoDashboardDeInexistente)
	ctx.Step(`^se completa la historia de (\d+) puntos del sprint en curso de "([^"]*)"$`, w.seCompletaHistoria)

	// Then
	ctx.Step(`^el backlog del dashboard tiene (\d+) pendientes, (\d+) en progreso y (\d+) hechos, de (\d+) en total$`, w.backlogTiene)
	ctx.Step(`^el sprint en curso del dashboard es el (\d+), con (\d+) puntos planificados y (\d+) completados$`, w.sprintEnCursoEs)
	ctx.Step(`^el dashboard no tiene sprint en curso$`, w.sinSprintEnCurso)
	ctx.Step(`^los sprints cerrados del dashboard son el (\d+)(?: y el (\d+))?, en ese orden$`, w.cerradosSon)
	ctx.Step(`^el dashboard no tiene sprints cerrados$`, w.sinCerrados)
	ctx.Step(`^el sprint (\d+) del dashboard tiene (\d+) puntos planificados, (\d+) completados y (\d+)% completado$`, w.cerradoTienePuntos)
	ctx.Step(`^el sprint (\d+) del dashboard tiene (\d+) defectos detectados y (\d+) resueltos$`, w.cerradoTieneDefectos)
	ctx.Step(`^la velocidad del dashboard es (\d+)$`, w.velocidadEs)
	ctx.Step(`^la velocidad del dashboard no tiene datos$`, w.velocidadSinDatos)
	ctx.Step(`^el dashboard tiene (\d+) horas estimadas y (\d+) reales, con una desviación del (-?\d+)%$`, w.horasSon)
	ctx.Step(`^el dashboard no tiene datos de desviación de horas$`, w.sinDesviacion)
	ctx.Step(`^se rechaza la consulta del dashboard con el error "([^"]*)"$`, w.seRechazaDashboard)
}

// --- Given ---

func (w *dashboardWorld) proyectoParaDashboard(nombre string) error {
	if w.sistema == nil {
		sistema, err := nuevoSistemaDashboard()
		if err != nil {
			return err
		}
		w.sistema = sistema
		w.proyectos = map[string]int64{}
		w.sprints = map[string][]int64{}
		w.enCurso = map[string]int64{}
		w.historias = map[string]map[int]int64{}
	}
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *dashboardWorld) hoyEs(fecha string) error {
	hoy, err := time.Parse(time.DateOnly, fecha)
	if err != nil {
		return fmt.Errorf("fecha inválida en el escenario %q: %w", fecha, err)
	}
	w.sistema.FijarHoy(hoy)
	return nil
}

func (w *dashboardWorld) tieneElementos(proyecto string, cantidad int, estado string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	for range cantidad {
		if _, err := w.sistema.CrearElemento(id, domain.EstadoElemento(estado), nil, 0); err != nil {
			return fmt.Errorf("no se pudo preparar un elemento %q de %q: %w", estado, proyecto, err)
		}
	}
	return nil
}

func (w *dashboardWorld) tieneSprintAbierto(proyecto, inicio, fin string) error {
	sid, err := w.crearSprint(proyecto, domain.SprintAbierto, inicio, fin, 0, 0)
	if err != nil {
		return err
	}
	w.enCurso[proyecto] = sid
	return nil
}

func (w *dashboardWorld) tieneSprintCerrado(proyecto, inicio, fin string, planificados, completados int) error {
	_, err := w.crearSprint(proyecto, domain.SprintCerrado, inicio, fin, planificados, completados)
	return err
}

func (w *dashboardWorld) historiaEnCurso(proyecto, estado string, puntos int) error {
	return w.crearHistoriaEnCurso(proyecto, estado, &puntos)
}

func (w *dashboardWorld) historiaEnCursoSinEstimar(proyecto, estado string) error {
	return w.crearHistoriaEnCurso(proyecto, estado, nil)
}

func (w *dashboardWorld) defectosDetectadosYResuelto(cantidad, deteccion int, proyecto string, resolucion int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sDeteccion, err := w.idSprint(proyecto, deteccion)
	if err != nil {
		return err
	}
	sResolucion, err := w.idSprint(proyecto, resolucion)
	if err != nil {
		return err
	}
	eid, err := w.sistema.CrearElemento(pid, domain.EstadoPendiente, nil, 0)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el elemento de los defectos: %w", err)
	}
	for i := range cantidad {
		estado, resuelto := domain.DefectoAbierto, int64(0)
		if i == 0 {
			estado, resuelto = domain.DefectoResuelto, sResolucion
		}
		if err := w.sistema.CrearDefecto(pid, eid, estado, sDeteccion, resuelto); err != nil {
			return fmt.Errorf("no se pudo preparar el defecto %d: %w", i+1, err)
		}
	}
	return nil
}

func (w *dashboardWorld) elementoConHoras(proyecto string, estimadas, reales int) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	e := float64(estimadas)
	return w.sistema.CrearElementoConHoras(id, &e, float64(reales))
}

func (w *dashboardWorld) elementoSinEstimarConHoras(proyecto string, reales int) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	return w.sistema.CrearElementoConHoras(id, nil, float64(reales))
}

// --- When ---

func (w *dashboardWorld) consultoDashboard(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultima, w.ultimoErr = w.sistema.ConsultarDashboard(id)
	return nil
}

func (w *dashboardWorld) consultoDashboardDeInexistente() error {
	w.ultima, w.ultimoErr = w.sistema.ConsultarDashboard(idInexistente)
	return nil
}

func (w *dashboardWorld) seCompletaHistoria(puntos int, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	eid, ok := w.historias[proyecto][puntos]
	if !ok {
		return fmt.Errorf("no se preparó una historia de %d puntos en el sprint en curso de %q", puntos, proyecto)
	}
	if err := w.sistema.CompletarHistoria(pid, eid); err != nil {
		return fmt.Errorf("no se pudo completar la historia: %w", err)
	}
	return nil
}

// --- Then ---

func (w *dashboardWorld) backlogTiene(pendientes, enProgreso, hechos, total int) error {
	if err := w.consultado(); err != nil {
		return err
	}
	v := w.ultima
	if v.Pendientes != pendientes || v.EnProgreso != enProgreso || v.Hechos != hechos || v.Total != total {
		return fmt.Errorf("backlog = %d pendientes, %d en progreso, %d hechos, %d en total; se esperaba %d, %d, %d, %d",
			v.Pendientes, v.EnProgreso, v.Hechos, v.Total, pendientes, enProgreso, hechos, total)
	}
	return nil
}

func (w *dashboardWorld) sprintEnCursoEs(numero, planificados, completados int) error {
	if err := w.consultado(); err != nil {
		return err
	}
	v := w.ultima
	if !v.HaySprintEnCurso {
		return errors.New("se esperaba un sprint en curso, pero el dashboard no tiene")
	}
	if v.NumeroEnCurso != numero || v.PlanificadosEnCurso != planificados || v.HechosEnCurso != completados {
		return fmt.Errorf("sprint en curso = %d con %d planificados y %d completados; se esperaba %d con %d y %d",
			v.NumeroEnCurso, v.PlanificadosEnCurso, v.HechosEnCurso, numero, planificados, completados)
	}
	return nil
}

func (w *dashboardWorld) sinSprintEnCurso() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if w.ultima.HaySprintEnCurso {
		return fmt.Errorf("se esperaba que no hubiera sprint en curso, pero está el %d", w.ultima.NumeroEnCurso)
	}
	return nil
}

func (w *dashboardWorld) cerradosSon(primero int, segundo string) error {
	if err := w.consultado(); err != nil {
		return err
	}
	want := []int{primero}
	if segundo != "" {
		var n int
		if _, err := fmt.Sscan(segundo, &n); err != nil {
			return fmt.Errorf("número de sprint inválido en el escenario %q: %w", segundo, err)
		}
		want = append(want, n)
	}
	got := make([]int, len(w.ultima.Cerrados))
	for i, c := range w.ultima.Cerrados {
		got[i] = c.Numero
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		return fmt.Errorf("sprints cerrados = %v, se esperaba %v", got, want)
	}
	return nil
}

func (w *dashboardWorld) sinCerrados() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if n := len(w.ultima.Cerrados); n != 0 {
		return fmt.Errorf("se esperaba que no hubiera sprints cerrados, pero hay %d", n)
	}
	return nil
}

func (w *dashboardWorld) cerradoTienePuntos(numero, planificados, completados, porcentaje int) error {
	c, err := w.cerrado(numero)
	if err != nil {
		return err
	}
	if c.Planificados != planificados || c.Completados != completados || !c.HayPorcentaje || c.Porcentaje != float64(porcentaje) {
		return fmt.Errorf("sprint %d = %d planificados, %d completados, %v%% (con datos: %v); se esperaba %d, %d, %d%%",
			numero, c.Planificados, c.Completados, c.Porcentaje, c.HayPorcentaje, planificados, completados, porcentaje)
	}
	return nil
}

func (w *dashboardWorld) cerradoTieneDefectos(numero, detectados, resueltos int) error {
	c, err := w.cerrado(numero)
	if err != nil {
		return err
	}
	if c.Detectados != detectados || c.Resueltos != resueltos {
		return fmt.Errorf("sprint %d = %d defectos detectados y %d resueltos; se esperaba %d y %d",
			numero, c.Detectados, c.Resueltos, detectados, resueltos)
	}
	return nil
}

func (w *dashboardWorld) velocidadEs(velocidad int) error {
	if err := w.consultado(); err != nil {
		return err
	}
	if !w.ultima.HayVelocidad || w.ultima.Velocidad != float64(velocidad) {
		return fmt.Errorf("velocidad = %v (con datos: %v), se esperaba %d", w.ultima.Velocidad, w.ultima.HayVelocidad, velocidad)
	}
	return nil
}

func (w *dashboardWorld) velocidadSinDatos() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if w.ultima.HayVelocidad {
		return fmt.Errorf("se esperaba la velocidad sin datos, pero vale %v", w.ultima.Velocidad)
	}
	return nil
}

func (w *dashboardWorld) horasSon(estimadas, reales, desviacion int) error {
	if err := w.consultado(); err != nil {
		return err
	}
	v := w.ultima
	if v.HorasEstimadas != float64(estimadas) || v.HorasReales != float64(reales) || v.Desviacion == nil || *v.Desviacion != float64(desviacion) {
		return fmt.Errorf("horas = %v estimadas, %v reales, desviación %v; se esperaba %d, %d y %d%%",
			v.HorasEstimadas, v.HorasReales, textoDesviacion(v.Desviacion), estimadas, reales, desviacion)
	}
	return nil
}

func (w *dashboardWorld) sinDesviacion() error {
	if err := w.consultado(); err != nil {
		return err
	}
	if w.ultima.Desviacion != nil {
		return fmt.Errorf("se esperaba la desviación sin datos, pero vale %v%%", *w.ultima.Desviacion)
	}
	return nil
}

func (w *dashboardWorld) seRechazaDashboard(descripcion string) error {
	esperado, ok := erroresHU25()[descripcion]
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

func (w *dashboardWorld) consultado() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el dashboard, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *dashboardWorld) cerrado(numero int) (cerradoVista, error) {
	if err := w.consultado(); err != nil {
		return cerradoVista{}, err
	}
	for _, c := range w.ultima.Cerrados {
		if c.Numero == numero {
			return c, nil
		}
	}
	return cerradoVista{}, fmt.Errorf("el sprint %d no aparece entre los sprints cerrados del dashboard", numero)
}

func (w *dashboardWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *dashboardWorld) idSprint(proyecto string, numero int) (int64, error) {
	ids := w.sprints[proyecto]
	if numero < 1 || numero > len(ids) {
		return 0, fmt.Errorf("el sprint %d de %q no se preparó en el escenario", numero, proyecto)
	}
	return ids[numero-1], nil
}

func (w *dashboardWorld) crearSprint(proyecto string, estado domain.EstadoSprint, inicio, fin string, planificados, completados int) (int64, error) {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return 0, err
	}
	desde, err := time.Parse(time.DateOnly, inicio)
	if err != nil {
		return 0, fmt.Errorf("fecha inválida en el escenario %q: %w", inicio, err)
	}
	hasta, err := time.Parse(time.DateOnly, fin)
	if err != nil {
		return 0, fmt.Errorf("fecha inválida en el escenario %q: %w", fin, err)
	}
	sid, err := w.sistema.CrearSprint(pid, estado, desde, hasta, planificados, completados)
	if err != nil {
		return 0, fmt.Errorf("no se pudo preparar el sprint de %q: %w", proyecto, err)
	}
	w.sprints[proyecto] = append(w.sprints[proyecto], sid)
	return sid, nil
}

func (w *dashboardWorld) crearHistoriaEnCurso(proyecto, estado string, puntos *int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, ok := w.enCurso[proyecto]
	if !ok {
		return fmt.Errorf("%q no tiene un sprint abierto preparado en el escenario", proyecto)
	}
	eid, err := w.sistema.CrearElemento(pid, domain.EstadoElemento(estado), puntos, sid)
	if err != nil {
		return fmt.Errorf("no se pudo preparar la historia del sprint en curso: %w", err)
	}
	if puntos != nil {
		if w.historias[proyecto] == nil {
			w.historias[proyecto] = map[int]int64{}
		}
		w.historias[proyecto][*puntos] = eid
	}
	return nil
}

func textoDesviacion(d *float64) string {
	if d == nil {
		return "sin datos"
	}
	return fmt.Sprintf("%v%%", *d)
}
