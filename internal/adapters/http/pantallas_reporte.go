package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// armadorReporteProyecto es lo que el handler necesita del caso de uso ArmarReporteProyecto (HU-26).
type armadorReporteProyecto interface {
	Ejecutar(ctx context.Context, proyectoID int64) (domain.Reporte, error)
}

// armadorReporteSprint es lo que el handler necesita del caso de uso ArmarReporteSprint (HU-26).
type armadorReporteSprint interface {
	Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.Reporte, error)
}

// Datos de la pantalla del reporte (HU-26): solo texto ya formateado, sin lógica de negocio.
type (
	vistaHistoriaReporte struct {
		Titulo, Prioridad, Estado, StoryPoints string
	}
	vistaMetricaReporte struct {
		Sprint                                int
		HayPorcentaje                         bool
		Porcentaje                            string
		DefectosDetectados, DefectosResueltos int
	}
	vistaDefectoReporte struct {
		Elemento, Descripcion, Severidad, Estado string
		SprintDeteccion, SprintResolucion        *int64
	}
	datosReporte struct {
		Proyecto                              vistaProyecto
		EsDeSprint                            bool
		NumeroSprint                          int
		SprintGoal                            string
		InicioSprint, FinSprint, EstadoSprint string

		Historias                                []vistaHistoriaReporte
		HistoriasCompletadas, HistoriasTotal     int
		StoryPointsCompletados, StoryPointsTotal int

		EsfuerzoEstimadas, EsfuerzoReales, EsfuerzoDesviacionHoras, EsfuerzoDesviacionPorcentaje string

		HayVelocidad bool
		Velocidad    string
		Metricas     []vistaMetricaReporte

		Defectos []vistaDefectoReporte

		URLPDF string // descarga del PDF (HU-27, CA-27.21)
	}
)

// registrarPantallaReporte agrega las pantallas de HU-26: el reporte de un proyecto y el de un
// Sprint. No tienen endpoint JSON (spec, sección 2): el único consumidor es la pantalla misma, y
// HU-27 llama al caso de uso directo, no por HTTP.
func (p *pantallas) registrarPantallaReporte(mux *http.ServeMux) {
	mux.HandleFunc("GET /web/proyectos/{id}/reporte", p.reporteProyecto)
	mux.HandleFunc("GET /web/proyectos/{id}/sprints/{sprintID}/reporte", p.reporteSprint)
}

// reporteProyecto muestra el reporte de un proyecto. Implementa HU-26 (CA-26.1, CA-26.3 a CA-26.6,
// CA-26.9, CA-26.15).
func (p *pantallas) reporteProyecto(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	reporte, err := p.cu.ArmarReporteProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		p.fallaInterna(w, "armar el reporte del proyecto", err)
		return
	}
	p.renderizar(w, http.StatusOK, "reporte.html", aDatosReporte(reporte))
}

// reporteSprint muestra el reporte de un Sprint, en cualquier estado (RN-26.10). Implementa HU-26
// (CA-26.2, CA-26.7, CA-26.8, CA-26.10).
func (p *pantallas) reporteSprint(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	sprintID, err := sprintIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return
	}
	reporte, err := p.cu.ArmarReporteSprint.Ejecutar(r.Context(), id, sprintID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			p.errorDeConsulta(w, err)
		case errors.Is(err, domain.ErrSprintNoEncontrado):
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrSprintNoEncontrado.Error()})
		default:
			p.fallaInterna(w, "armar el reporte del sprint", err)
		}
		return
	}
	p.renderizar(w, http.StatusOK, "reporte.html", aDatosReporte(reporte))
}

// aDatosReporte da formato al reporte (RN-26.1 a RN-26.9), reutilizando textoHoras, textosDesviacion
// y textoPorcentaje que ya usan HU-19, HU-22 y HU-23.
func aDatosReporte(r domain.Reporte) datosReporte {
	d := datosReporte{
		Proyecto:               aVistaProyecto(r.Proyecto),
		Historias:              aVistasHistoriasReporte(r.Historias),
		HistoriasCompletadas:   r.HistoriasCompletadas,
		HistoriasTotal:         len(r.Historias),
		StoryPointsCompletados: r.StoryPointsCompletados,
		StoryPointsTotal:       r.StoryPointsTotal,
		HayVelocidad:           r.HayVelocidad,
		Metricas:               aVistasMetricasReporte(r.Metricas),
		Defectos:               aVistasDefectosReporte(r.Defectos, r.Historias),
		URLPDF:                 fmt.Sprintf("/web/proyectos/%d/reporte/pdf", r.Proyecto.ID), // HU-27
	}
	if r.Sprint != nil {
		d.URLPDF = fmt.Sprintf("/web/proyectos/%d/sprints/%d/reporte/pdf", r.Proyecto.ID, r.Sprint.ID) // HU-27
		d.EsDeSprint = true
		d.NumeroSprint = r.Sprint.Numero
		d.SprintGoal = r.Sprint.Objetivo
		d.InicioSprint = r.Sprint.FechaInicio.Format(time.DateOnly)
		d.FinSprint = r.Sprint.FechaFin.Format(time.DateOnly)
		d.EstadoSprint = string(r.Sprint.Estado)
	}
	d.EsfuerzoEstimadas = textoSinEstimar
	if !r.Esfuerzo.SinEstimar() {
		d.EsfuerzoEstimadas = textoHoras(r.Esfuerzo.HorasEstimadas)
	}
	d.EsfuerzoReales = textoHoras(r.Esfuerzo.HorasReales)
	d.EsfuerzoDesviacionHoras, d.EsfuerzoDesviacionPorcentaje = textosDesviacion(r.Esfuerzo.DesviacionHoras, r.Esfuerzo.DesviacionPorcentaje)
	if r.HayVelocidad {
		d.Velocidad = conComa(strconv.FormatFloat(r.Velocidad, 'f', -1, 64))
	}
	return d
}

func aVistasHistoriasReporte(historias []domain.ElementoBacklog) []vistaHistoriaReporte {
	vistas := make([]vistaHistoriaReporte, len(historias))
	for i, h := range historias {
		sp := textoSinEstimar
		if h.StoryPoints != nil {
			sp = strconv.Itoa(*h.StoryPoints)
		}
		vistas[i] = vistaHistoriaReporte{Titulo: h.Titulo, Prioridad: string(h.Prioridad), Estado: string(h.Estado), StoryPoints: sp}
	}
	return vistas
}

func aVistasMetricasReporte(metricas []domain.FilaMetrica) []vistaMetricaReporte {
	vistas := make([]vistaMetricaReporte, len(metricas))
	for i, m := range metricas {
		v := vistaMetricaReporte{
			Sprint: m.Sprint.Numero, HayPorcentaje: m.HayPorcentaje,
			DefectosDetectados: m.DefectosDetectados, DefectosResueltos: m.DefectosResueltos,
		}
		if m.HayPorcentaje {
			v.Porcentaje = textoPorcentaje(m.PorcentajeCompletado)
		}
		vistas[i] = v
	}
	return vistas
}

// aVistasDefectosReporte arma las filas de defectos, con el título del elemento relacionado si está
// entre las historias del reporte (el reporte de un Sprint no trae las de otros Sprints).
func aVistasDefectosReporte(defectos []domain.Defecto, historias []domain.ElementoBacklog) []vistaDefectoReporte {
	titulos := make(map[int64]string, len(historias))
	for _, h := range historias {
		titulos[h.ID] = h.Titulo
	}
	vistas := make([]vistaDefectoReporte, len(defectos))
	for i, def := range defectos {
		elemento := titulos[def.ElementoID]
		if elemento == "" {
			elemento = "—"
		}
		vistas[i] = vistaDefectoReporte{
			Elemento: elemento, Descripcion: def.Descripcion, Severidad: string(def.Severidad), Estado: string(def.Estado),
			SprintDeteccion: def.SprintDeteccionID, SprintResolucion: def.SprintResolucionID,
		}
	}
	return vistas
}
