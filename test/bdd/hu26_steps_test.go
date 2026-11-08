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

// sistemaReporte es la puerta de entrada de los escenarios de HU-26 al núcleo. La implementa
// sistemaReporteReal, sobre SQLite (hu26_sistema_test.go).
type sistemaReporte interface {
	CrearProyecto(nombre string) (int64, error)
	CrearSprint(proyectoID int64, estado domain.EstadoSprint, planificados, completados int) (int64, error)
	// CrearHistoria inserta una historia del backlog directo en la base. sprintID en 0 es "sin
	// asignar".
	CrearHistoria(proyectoID int64, titulo string, estado domain.EstadoElemento, storyPoints *int, horasEstimadas *float64, sprintID int64) (int64, error)
	RegistrarHoras(proyectoID, elementoID int64, horas float64) error
	// CrearDefecto inserta un defecto directo, con el Sprint de detección y el de resolución que
	// indique el escenario (0 = sin Sprint).
	CrearDefecto(proyectoID int64, sprintDeteccionID, sprintResolucionID int64) error
	ConsultarReporteProyecto(proyectoID int64) (domain.Reporte, error)
	ConsultarReporteSprint(proyectoID, sprintID int64) (domain.Reporte, error)
	Cerrar() error
}

// erroresHU26 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU26() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"sprint no encontrado":   domain.ErrSprintNoEncontrado,
	}
}

type reporteWorld struct {
	sistema   sistemaReporte
	proyectos map[string]int64
	sprints   map[string][]int64          // sprints[proyecto] = IDs en orden de creación (numero = índice+1)
	historias map[string]map[string]int64 // historias[proyecto][titulo] = ID
	ultimo    domain.Reporte
	ultimoErr error
}

func registrarStepsHU26(ctx *godog.ScenarioContext) {
	w := &reporteWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaReporte()
		if err != nil {
			return c, err
		}
		*w = reporteWorld{sistema: sistema, proyectos: map[string]int64{}, sprints: map[string][]int64{}, historias: map[string]map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para el reporte$`, w.proyectoParaReporte)
	ctx.Step(`^"([^"]*)" tiene una historia "([^"]*)" (Pendiente|Hecha) con (\d+) story points asignada al sprint (\d+)$`, w.tieneHistoriaAsignada)
	ctx.Step(`^"([^"]*)" tiene una historia "([^"]*)" (Pendiente|Hecha) con (\d+) story points$`, w.tieneHistoria)
	ctx.Step(`^"([^"]*)" tiene una historia "([^"]*)" (Pendiente|Hecha) con (\d+) horas estimadas$`, w.tieneHistoriaConHoras)
	ctx.Step(`^se registraron (\d+) horas reales en "([^"]*)" de "([^"]*)"$`, w.registraronHorasReales)
	ctx.Step(`^"([^"]*)" tiene para el reporte (?:un|otro) sprint cerrado con (\d+) puntos planificados y (\d+) completados$`, w.tieneSprintCerrado)
	ctx.Step(`^"([^"]*)" tiene para el reporte un sprint todavía abierto$`, w.tieneSprintAbierto)
	ctx.Step(`^se detectó un defecto en el sprint (\d+) de "([^"]*)" y se resolvió en el sprint (\d+)$`, w.detectoYResolvioDefecto)
	ctx.Step(`^se detectó un defecto en el sprint (\d+) de "([^"]*)"$`, w.detectoDefecto)

	// When
	ctx.Step(`^consulto el reporte de "([^"]*)"$`, w.consultoReporte)
	ctx.Step(`^consulto el reporte del sprint (\d+) de "([^"]*)"$`, w.consultoReporteDeSprint)
	ctx.Step(`^consulto el reporte de un proyecto inexistente$`, w.consultoReporteDeProyectoInexistente)
	ctx.Step(`^consulto el reporte de un sprint inexistente de "([^"]*)"$`, w.consultoReporteDeSprintInexistente)

	// Then
	ctx.Step(`^el reporte tiene (\d+) de (\d+) historias completadas y (\d+) de (\d+) story points$`, w.reporteTieneResumen)
	ctx.Step(`^el reporte tiene (\d+) historias$`, w.reporteTieneNHistorias)
	ctx.Step(`^el esfuerzo del reporte es (\d+(?:,\d+)?) horas estimadas y (\d+(?:,\d+)?) horas reales$`, w.esfuerzoEs)
	ctx.Step(`^el reporte tiene velocidad (-?\d+(?:,\d+)?) sobre (\d+) sprints cerrados$`, w.tieneVelocidad)
	ctx.Step(`^el reporte no tiene velocidad$`, w.sinVelocidad)
	ctx.Step(`^el reporte no tiene métricas de sprints$`, w.sinMetricas)
	ctx.Step(`^el sprint (\d+) del reporte tiene (-?\d+(?:,\d+)?)% completado$`, w.sprintTienePorcentaje)
	ctx.Step(`^el sprint (\d+) del reporte no tiene datos de porcentaje$`, w.sprintSinDatos)
	ctx.Step(`^el sprint (\d+) del reporte tiene (\d+) defectos detectados$`, w.sprintTieneDefectosDetectados)
	ctx.Step(`^el reporte tiene (\d+) defectos$`, w.reporteTieneNDefectos)
	ctx.Step(`^se rechaza la consulta del reporte con el error "([^"]*)"$`, w.seRechazaConsulta)
}

// --- Given ---

func (w *reporteWorld) proyectoParaReporte(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *reporteWorld) tieneHistoria(proyecto, titulo, estadoTxt string, storyPoints int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	return w.crearHistoria(pid, proyecto, titulo, estadoTxt, &storyPoints, nil, 0)
}

func (w *reporteWorld) tieneHistoriaAsignada(proyecto, titulo, estadoTxt string, storyPoints, numeroSprint int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(proyecto, numeroSprint)
	if err != nil {
		return err
	}
	return w.crearHistoria(pid, proyecto, titulo, estadoTxt, &storyPoints, nil, sid)
}

func (w *reporteWorld) tieneHistoriaConHoras(proyecto, titulo, estadoTxt string, horas int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	h := float64(horas)
	return w.crearHistoria(pid, proyecto, titulo, estadoTxt, nil, &h, 0)
}

func (w *reporteWorld) crearHistoria(proyectoID int64, proyecto, titulo, estadoTxt string, storyPoints *int, horasEstimadas *float64, sprintID int64) error {
	estado := domain.EstadoPendiente
	if estadoTxt == "Hecha" {
		estado = domain.EstadoHecho
	}
	id, err := w.sistema.CrearHistoria(proyectoID, titulo, estado, storyPoints, horasEstimadas, sprintID)
	if err != nil {
		return fmt.Errorf("no se pudo preparar la historia %q de %q: %w", titulo, proyecto, err)
	}
	if w.historias[proyecto] == nil {
		w.historias[proyecto] = map[string]int64{}
	}
	w.historias[proyecto][titulo] = id
	return nil
}

func (w *reporteWorld) registraronHorasReales(horas int, titulo, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	eid, err := w.idHistoria(proyecto, titulo)
	if err != nil {
		return err
	}
	if err := w.sistema.RegistrarHoras(pid, eid, float64(horas)); err != nil {
		return fmt.Errorf("no se pudieron registrar horas reales en %q: %w", titulo, err)
	}
	return nil
}

func (w *reporteWorld) tieneSprintCerrado(proyecto string, planificados, completados int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.sistema.CrearSprint(pid, domain.SprintCerrado, planificados, completados)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint cerrado del proyecto %q: %w", proyecto, err)
	}
	w.sprints[proyecto] = append(w.sprints[proyecto], sid)
	return nil
}

func (w *reporteWorld) tieneSprintAbierto(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.sistema.CrearSprint(pid, domain.SprintAbierto, 0, 0)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint abierto del proyecto %q: %w", proyecto, err)
	}
	w.sprints[proyecto] = append(w.sprints[proyecto], sid)
	return nil
}

func (w *reporteWorld) detectoDefecto(numero int, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(proyecto, numero)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearDefecto(pid, sid, 0); err != nil {
		return fmt.Errorf("no se pudo preparar un defecto detectado en el sprint %d de %q: %w", numero, proyecto, err)
	}
	return nil
}

func (w *reporteWorld) detectoYResolvioDefecto(numeroDeteccion int, proyecto string, numeroResolucion int) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sidDet, err := w.idSprint(proyecto, numeroDeteccion)
	if err != nil {
		return err
	}
	sidRes, err := w.idSprint(proyecto, numeroResolucion)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearDefecto(pid, sidDet, sidRes); err != nil {
		return fmt.Errorf("no se pudo preparar el defecto detectado/resuelto en %q: %w", proyecto, err)
	}
	return nil
}

// --- When ---

func (w *reporteWorld) consultoReporte(proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultimo, w.ultimoErr = w.sistema.ConsultarReporteProyecto(id)
	return nil
}

func (w *reporteWorld) consultoReporteDeSprint(numero int, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	sid, err := w.idSprint(proyecto, numero)
	if err != nil {
		return err
	}
	w.ultimo, w.ultimoErr = w.sistema.ConsultarReporteSprint(pid, sid)
	return nil
}

func (w *reporteWorld) consultoReporteDeProyectoInexistente() error {
	w.ultimo, w.ultimoErr = w.sistema.ConsultarReporteProyecto(idInexistente)
	return nil
}

func (w *reporteWorld) consultoReporteDeSprintInexistente(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultimo, w.ultimoErr = w.sistema.ConsultarReporteSprint(pid, idInexistente)
	return nil
}

// --- Then ---

func (w *reporteWorld) reporteTieneNHistorias(n int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	if len(w.ultimo.Historias) != n {
		return fmt.Errorf("el reporte tiene %d historias, se esperaban %d", len(w.ultimo.Historias), n)
	}
	return nil
}

func (w *reporteWorld) reporteTieneResumen(completadas, totalHistorias, spCompletados, spTotal int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	if len(w.ultimo.Historias) != totalHistorias {
		return fmt.Errorf("el reporte tiene %d historias, se esperaban %d", len(w.ultimo.Historias), totalHistorias)
	}
	if w.ultimo.HistoriasCompletadas != completadas {
		return fmt.Errorf("HistoriasCompletadas = %d, se esperaban %d", w.ultimo.HistoriasCompletadas, completadas)
	}
	if w.ultimo.StoryPointsCompletados != spCompletados || w.ultimo.StoryPointsTotal != spTotal {
		return fmt.Errorf("story points = %d de %d, se esperaban %d de %d",
			w.ultimo.StoryPointsCompletados, w.ultimo.StoryPointsTotal, spCompletados, spTotal)
	}
	return nil
}

func (w *reporteWorld) esfuerzoEs(estimadas, reales string) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	wantEst, err := parseComaReporte(estimadas)
	if err != nil {
		return err
	}
	wantReal, err := parseComaReporte(reales)
	if err != nil {
		return err
	}
	if w.ultimo.Esfuerzo.HorasEstimadas != wantEst || w.ultimo.Esfuerzo.HorasReales != wantReal {
		return fmt.Errorf("esfuerzo = %v estimadas, %v reales; se esperaban %v y %v",
			w.ultimo.Esfuerzo.HorasEstimadas, w.ultimo.Esfuerzo.HorasReales, wantEst, wantReal)
	}
	return nil
}

func (w *reporteWorld) tieneVelocidad(velocidad string, cantidadSprints int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	if !w.ultimo.HayVelocidad {
		return fmt.Errorf("se esperaba velocidad con datos, pero no hay")
	}
	want, err := parseComaReporte(velocidad)
	if err != nil {
		return err
	}
	if w.ultimo.Velocidad != want || len(w.ultimo.Metricas) != cantidadSprints {
		return fmt.Errorf("velocidad = %v sobre %d sprints, se esperaba %v sobre %d",
			w.ultimo.Velocidad, len(w.ultimo.Metricas), want, cantidadSprints)
	}
	return nil
}

func (w *reporteWorld) sinVelocidad() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	if w.ultimo.HayVelocidad {
		return fmt.Errorf("se esperaba sin velocidad, pero hay %v", w.ultimo.Velocidad)
	}
	return nil
}

func (w *reporteWorld) sinMetricas() error {
	if len(w.ultimo.Metricas) != 0 {
		return fmt.Errorf("se esperaban sin métricas, hay %+v", w.ultimo.Metricas)
	}
	return nil
}

func (w *reporteWorld) sprintTienePorcentaje(numero int, porcentaje string) error {
	m, err := w.metrica(numero)
	if err != nil {
		return err
	}
	if !m.HayPorcentaje {
		return fmt.Errorf("se esperaba porcentaje con datos en el sprint %d, pero no hay", numero)
	}
	want, err := parseComaReporte(porcentaje)
	if err != nil {
		return err
	}
	if m.PorcentajeCompletado != want {
		return fmt.Errorf("el sprint %d tiene %v%% completado, se esperaba %v%%", numero, m.PorcentajeCompletado, want)
	}
	return nil
}

func (w *reporteWorld) sprintSinDatos(numero int) error {
	m, err := w.metrica(numero)
	if err != nil {
		return err
	}
	if m.HayPorcentaje {
		return fmt.Errorf("se esperaba que el sprint %d no tuviera datos de porcentaje", numero)
	}
	return nil
}

func (w *reporteWorld) sprintTieneDefectosDetectados(numero, cantidad int) error {
	m, err := w.metrica(numero)
	if err != nil {
		return err
	}
	if m.DefectosDetectados != cantidad {
		return fmt.Errorf("el sprint %d tiene %d defectos detectados, se esperaban %d", numero, m.DefectosDetectados, cantidad)
	}
	return nil
}

func (w *reporteWorld) reporteTieneNDefectos(n int) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el reporte, pero hubo error: %w", w.ultimoErr)
	}
	if len(w.ultimo.Defectos) != n {
		return fmt.Errorf("el reporte tiene %d defectos, se esperaban %d", len(w.ultimo.Defectos), n)
	}
	return nil
}

func (w *reporteWorld) seRechazaConsulta(descripcion string) error {
	esperado, ok := erroresHU26()[descripcion]
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

func (w *reporteWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *reporteWorld) idSprint(proyecto string, numero int) (int64, error) {
	ids, ok := w.sprints[proyecto]
	if !ok || numero < 1 || numero > len(ids) {
		return 0, fmt.Errorf("el sprint %d de %q no se preparó en el escenario", numero, proyecto)
	}
	return ids[numero-1], nil
}

func (w *reporteWorld) idHistoria(proyecto, titulo string) (int64, error) {
	id, ok := w.historias[proyecto][titulo]
	if !ok {
		return 0, fmt.Errorf("la historia %q de %q no se preparó en el escenario", titulo, proyecto)
	}
	return id, nil
}

func (w *reporteWorld) metrica(numeroSprint int) (domain.FilaMetrica, error) {
	for _, m := range w.ultimo.Metricas {
		if m.Sprint.Numero == numeroSprint {
			return m, nil
		}
	}
	return domain.FilaMetrica{}, fmt.Errorf("no se encontró el sprint %d entre las métricas del reporte", numeroSprint)
}

func parseComaReporte(s string) (float64, error) {
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("número inválido en el escenario %q: %w", s, err)
	}
	return v, nil
}
