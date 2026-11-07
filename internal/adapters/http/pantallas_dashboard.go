package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// datosDashboard es lo que muestra la pantalla del Dashboard (HU-25).
type datosDashboard struct {
	Proyecto         vistaProyecto
	Backlog          vistaConteoBacklog
	HaySprintEnCurso bool
	SprintEnCurso    vistaSprintEnCurso
	// Graficos de los sprints cerrados: puntos, porcentaje y defectos. Vacío si no hay sprints cerrados.
	Graficos     []graficoBarras
	HayHoras     bool // false: ningún elemento estimado en horas (RN-23.6)
	GraficoHoras graficoBarras
	Desviacion   string
	// AvisoHoras explica, como el bloque de HU-23, que las horas reales de los elementos sin estimar
	// están en la barra de reales pero no en la desviación (RN-23.5). Vacío si todos están estimados.
	AvisoHoras string
}

type vistaConteoBacklog struct {
	Pendientes, EnProgreso, Hechos, Total int
}

type vistaSprintEnCurso struct {
	ID                        int64
	Numero                    int
	Objetivo                  string
	FechaInicio, FechaFin     string
	Planificados, Completados int
}

func (p *pantallas) registrarPantallasDashboard(mux *http.ServeMux) {
	mux.HandleFunc("GET /web/proyectos/{id}/dashboard", p.dashboard) // HU-25
}

// dashboard muestra el estado del proyecto y los gráficos (HU-25). Llama al caso de uso recién al
// atender el pedido, así cada vez que se abre muestra los datos vigentes (RN-25.1).
func (p *pantallas) dashboard(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	proyecto, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	d, err := p.cu.ConsultarDashboard.Ejecutar(r.Context(), id)
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		p.errorDeConsulta(w, err)
		return
	}
	if err != nil {
		p.fallaInterna(w, "consultar dashboard", err) // así el log dice qué falló
		return
	}
	p.renderizar(w, http.StatusOK, "dashboard.html", aDatosDashboard(aVistaProyecto(proyecto), d))
}

func aDatosDashboard(proyecto vistaProyecto, d app.Dashboard) datosDashboard {
	datos := datosDashboard{
		Proyecto: proyecto,
		Backlog: vistaConteoBacklog{
			Pendientes: d.Backlog.Pendientes, EnProgreso: d.Backlog.EnProgreso, Hechos: d.Backlog.Hechos, Total: d.Backlog.Total,
		},
		HaySprintEnCurso: d.HaySprintEnCurso,
	}
	if d.HaySprintEnCurso {
		s := d.SprintEnCurso.Sprint
		datos.SprintEnCurso = vistaSprintEnCurso{
			ID: s.ID, Numero: s.Numero, Objetivo: s.Objetivo,
			FechaInicio: s.FechaInicio.Format(time.DateOnly), FechaFin: s.FechaFin.Format(time.DateOnly),
			Planificados: d.SprintEnCurso.PuntosPlanificados, Completados: d.SprintEnCurso.PuntosCompletados,
		}
	}
	if len(d.SprintsCerrados) > 0 {
		datos.Graficos = graficosDeSprints(d.SprintsCerrados, d.Velocidad)
	}
	if !d.Horas.SinEstimar() {
		datos.HayHoras = true
		datos.GraficoHoras = armarGraficoBarras(especGrafico{
			ID: "horas", Titulo: "Horas estimadas contra reales", EncabezadoGrupos: "Horas",
			Descripcion: "Total de horas estimadas y reales del proyecto.",
			Series:      []serieGrafico{{Nombre: "Horas", Clase: "serie-1"}},
			Grupos: []grupoGrafico{
				{Etiqueta: "Estimadas", Valores: []*float64{&d.Horas.HorasEstimadas}},
				{Etiqueta: "Reales", Valores: []*float64{&d.Horas.HorasReales}},
			},
		})
		datos.Desviacion = textoPorcentaje(*d.Horas.DesviacionPorcentaje)
		if d.Horas.ElementosSinEstimar > 0 {
			datos.AvisoHoras = avisoSinEstimar(d.Horas)
		}
	}
	return datos
}

// graficosDeSprints arma los tres gráficos por sprint cerrado, del más antiguo al más reciente
// (RN-25.5): puntos (con la velocidad como línea), porcentaje completado y defectos.
func graficosDeSprints(cerrados []app.CalidadSprint, velocidad app.VelocidadProyecto) []graficoBarras {
	var puntos, porcentaje, defectos []grupoGrafico
	for _, c := range cerrados {
		etiqueta := "Sprint " + strconv.Itoa(c.Sprint.Numero)
		planificados, completados := float64(c.Sprint.PuntosPlanificados), float64(c.Sprint.PuntosCompletados)
		detectados, resueltos := float64(c.DefectosDetectados), float64(c.DefectosResueltos)
		var pct *float64
		if c.HayPorcentaje {
			v := c.PorcentajeCompletado
			pct = &v
		}
		puntos = append(puntos, grupoGrafico{Etiqueta: etiqueta, Valores: []*float64{&planificados, &completados}})
		porcentaje = append(porcentaje, grupoGrafico{Etiqueta: etiqueta, Valores: []*float64{pct}})
		defectos = append(defectos, grupoGrafico{Etiqueta: etiqueta, Valores: []*float64{&detectados, &resueltos}})
	}

	especPuntos := especGrafico{
		ID: "puntos", Titulo: "Puntos por sprint", EncabezadoGrupos: "Sprint",
		Descripcion: "Story points planificados y completados de cada sprint cerrado, del más antiguo al más reciente.",
		Series:      []serieGrafico{{Nombre: "Planificados", Clase: "serie-1"}, {Nombre: "Completados", Clase: "serie-2"}},
		Grupos:      puntos,
	}
	if velocidad.HayDatos {
		v := velocidad.Velocidad
		especPuntos.Referencia, especPuntos.EtiquetaReferencia = &v, "Velocidad"
		especPuntos.Descripcion += " La línea marca la velocidad del equipo."
	}
	return []graficoBarras{
		armarGraficoBarras(especPuntos),
		armarGraficoBarras(especGrafico{
			ID: "porcentaje", Titulo: "Porcentaje completado por sprint", EncabezadoGrupos: "Sprint",
			Descripcion: "Porcentaje de los story points planificados que se completó en cada sprint cerrado.",
			Series:      []serieGrafico{{Nombre: "% completado", Clase: "serie-2"}},
			Grupos:      porcentaje, Sufijo: "%", MinimoEscala: 100,
		}),
		armarGraficoBarras(especGrafico{
			ID: "defectos", Titulo: "Defectos por sprint", EncabezadoGrupos: "Sprint",
			Descripcion: "Defectos detectados y resueltos en cada sprint cerrado.",
			Series:      []serieGrafico{{Nombre: "Detectados", Clase: "serie-3"}, {Nombre: "Resueltos", Clase: "serie-2"}},
			Grupos:      defectos,
		}),
	}
}
