package http

import (
	"context"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Textos fijos de la plantilla reporte.html que también lleva el PDF (spec HU-27, P-11). Están repetidos
// a propósito: la plantilla no se toca, y TestADocumentoPDF_TextosSonLosDeLaPantalla falla si alguno deja
// de coincidir con la pantalla.
const (
	textoSinFechaFin       = "Sin fecha de fin"
	textoSinDatos          = "Sin datos"
	textoGuion             = "—" // sin elemento o sin Sprint
	textoSinHistorias      = "El reporte no tiene historias."
	textoSinVelocidad      = "Todavía no hay datos de velocidad (sin sprints cerrados)."
	textoSinSprintsCerrado = "El reporte no tiene sprints cerrados."
	textoSinDefectos       = "El reporte no tiene defectos."
)

// encabezadosDatoValor son los encabezados de las tablas de dato/valor (Resumen y Esfuerzo), que en la
// pantalla son listas de definiciones sin encabezados.
var encabezadosDatoValor = []string{"Dato", "Valor"}

// aDocumentoPDF convierte el reporte de HU-26, ya formateado para la pantalla, al documento que exporta
// HU-27: las cinco secciones de reporte.html, en el mismo orden, con los mismos títulos y los mismos
// textos (spec HU-27, RN-27.16 y P-11). No completa Emision ni Aviso, que son del caso de uso de HU-27.
func aDocumentoPDF(d datosReporte) domain.DocumentoPDF {
	return domain.DocumentoPDF{
		Titulo: tituloReporte(d),
		Secciones: []domain.SeccionPDF{
			{Titulo: "Resumen", Bloques: []domain.BloquePDF{tablaPDF(encabezadosDatoValor, filasResumen(d))}},
			{Titulo: "Historias y estimaciones", Bloques: bloquesHistorias(d)},
			{Titulo: "Esfuerzo", Bloques: []domain.BloquePDF{tablaPDF(encabezadosDatoValor, [][]string{
				{"Horas estimadas", d.EsfuerzoEstimadas},
				{"Horas reales", d.EsfuerzoReales},
				{"Desviación en horas", d.EsfuerzoDesviacionHoras},
				{"Desviación en porcentaje", d.EsfuerzoDesviacionPorcentaje},
			})}},
			{Titulo: "Métricas", Bloques: bloquesMetricas(d)},
			{Titulo: "Defectos", Bloques: bloquesDefectos(d)},
		},
	}
}

// tituloReporte es el <h1> de reporte.html.
func tituloReporte(d datosReporte) string {
	if d.EsDeSprint {
		return "Reporte del Sprint " + strconv.Itoa(d.NumeroSprint) + " de " + d.Proyecto.Nombre
	}
	return "Reporte de " + d.Proyecto.Nombre
}

func filasResumen(d datosReporte) [][]string {
	if d.EsDeSprint {
		return [][]string{
			{"Número", strconv.Itoa(d.NumeroSprint)},
			{"Sprint Goal", d.SprintGoal},
			{"Inicio", d.InicioSprint},
			{"Fin", d.FinSprint},
			{"Estado", d.EstadoSprint},
		}
	}
	fin := d.Proyecto.FechaFin
	if fin == "" {
		fin = textoSinFechaFin
	}
	return [][]string{
		{"Nombre", d.Proyecto.Nombre},
		{"Estado", d.Proyecto.Estado},
		{"Inicio", d.Proyecto.FechaInicio},
		{"Fin", fin},
	}
}

func bloquesHistorias(d datosReporte) []domain.BloquePDF {
	if len(d.Historias) == 0 {
		return []domain.BloquePDF{{Parrafo: textoSinHistorias}}
	}
	resumen := strconv.Itoa(d.HistoriasCompletadas) + " de " + strconv.Itoa(d.HistoriasTotal) +
		" historias completadas (" + strconv.Itoa(d.StoryPointsCompletados) + " de " +
		strconv.Itoa(d.StoryPointsTotal) + " story points)"
	filas := make([][]string, len(d.Historias))
	for i, h := range d.Historias {
		filas[i] = []string{h.Titulo, h.Prioridad, h.Estado, h.StoryPoints}
	}
	return []domain.BloquePDF{
		{Parrafo: resumen},
		tablaPDF([]string{"Título", "Prioridad", "Estado", "Story points"}, filas),
	}
}

// bloquesMetricas pone la velocidad antes de la tabla, solo en el reporte de un proyecto, como la
// pantalla de HU-26.
func bloquesMetricas(d datosReporte) []domain.BloquePDF {
	var bloques []domain.BloquePDF
	if !d.EsDeSprint {
		velocidad := textoSinVelocidad
		if d.HayVelocidad {
			velocidad = "Velocidad del equipo: " + d.Velocidad + " puntos por Sprint"
		}
		bloques = append(bloques, domain.BloquePDF{Parrafo: velocidad})
	}
	if len(d.Metricas) == 0 {
		return append(bloques, domain.BloquePDF{Parrafo: textoSinSprintsCerrado})
	}
	filas := make([][]string, len(d.Metricas))
	for i, m := range d.Metricas {
		porcentaje := textoSinDatos
		if m.HayPorcentaje {
			porcentaje = m.Porcentaje
		}
		filas[i] = []string{strconv.Itoa(m.Sprint), porcentaje, strconv.Itoa(m.DefectosDetectados), strconv.Itoa(m.DefectosResueltos)}
	}
	return append(bloques, tablaPDF([]string{"Sprint", "% completado", "Defectos detectados", "Defectos resueltos"}, filas))
}

func bloquesDefectos(d datosReporte) []domain.BloquePDF {
	if len(d.Defectos) == 0 {
		return []domain.BloquePDF{{Parrafo: textoSinDefectos}}
	}
	filas := make([][]string, len(d.Defectos))
	for i, def := range d.Defectos {
		filas[i] = []string{def.Elemento, def.Descripcion, def.Severidad, def.Estado,
			textoSprintDefecto(def.SprintDeteccion), textoSprintDefecto(def.SprintResolucion)}
	}
	return []domain.BloquePDF{tablaPDF(
		[]string{"Elemento", "Descripción", "Severidad", "Estado", "Sprint detección", "Sprint resolución"}, filas)}
}

// textoSprintDefecto es lo que muestra la plantilla en las columnas de Sprint de un defecto: "#" y el
// ID, o "—" si no tiene.
func textoSprintDefecto(id *int64) string {
	if id == nil {
		return textoGuion
	}
	return "#" + strconv.FormatInt(*id, 10)
}

func tablaPDF(columnas []string, filas [][]string) domain.BloquePDF {
	return domain.BloquePDF{Tabla: &domain.TablaPDF{Columnas: columnas, Filas: filas}}
}

// fuenteDocumentoReporte implementa app.FuenteDocumentoProyecto y app.FuenteDocumentoSprint (HU-27) sobre
// los casos de uso del reporte de HU-26: arma el reporte, le da el formato de la pantalla y lo convierte
// con aDocumentoPDF (spec HU-27, sección 5 y P-11).
type fuenteDocumentoReporte struct {
	proyecto armadorReporteProyecto
	sprint   armadorReporteSprint
}

// NuevaFuenteDocumentoReporte arma la fuente con los casos de uso ArmarReporteProyecto y
// ArmarReporteSprint de HU-26.
func NuevaFuenteDocumentoReporte(proyecto armadorReporteProyecto, sprint armadorReporteSprint) *fuenteDocumentoReporte {
	return &fuenteDocumentoReporte{proyecto: proyecto, sprint: sprint}
}

// DocumentoDeProyecto devuelve el documento del reporte de un proyecto. El error del caso de uso sale
// sin cambios: el caso de uso de HU-27 ya lo envuelve, y errors.Is tiene que seguir reconociéndolo.
func (f *fuenteDocumentoReporte) DocumentoDeProyecto(ctx context.Context, proyectoID int64) (domain.DocumentoPDF, error) {
	reporte, err := f.proyecto.Ejecutar(ctx, proyectoID)
	if err != nil {
		return domain.DocumentoPDF{}, err
	}
	return aDocumentoPDF(aDatosReporte(reporte)), nil
}

// DocumentoDeSprint devuelve el documento del reporte de un Sprint, con el mismo manejo de errores.
func (f *fuenteDocumentoReporte) DocumentoDeSprint(ctx context.Context, proyectoID, sprintID int64) (domain.DocumentoPDF, error) {
	reporte, err := f.sprint.Ejecutar(ctx, proyectoID, sprintID)
	if err != nil {
		return domain.DocumentoPDF{}, err
	}
	return aDocumentoPDF(aDatosReporte(reporte)), nil
}
