package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// exportadorReporteProyecto y exportadorReporteSprint son lo que las descargas necesitan de los casos de
// uso ExportarReporteProyecto y ExportarReporteSprint (HU-27).
type (
	exportadorReporteProyecto interface {
		Ejecutar(ctx context.Context, proyectoID int64) (app.ArchivoPDF, error)
	}
	exportadorReporteSprint interface {
		Ejecutar(ctx context.Context, proyectoID, sprintID int64) (app.ArchivoPDF, error)
	}
)

// erroresReportePDF son los errores de HU-27 que se informan con 404 (spec, sección 7). Se muestra el
// mensaje del dominio y no el que devolvió el caso de uso, que puede venir envuelto.
var erroresReportePDF = []error{domain.ErrProyectoNoEncontrado, domain.ErrSprintNoEncontrado}

// nombreArchivoPorDefecto se usa si el nombre del caso de uso queda vacío al sanearlo.
const nombreArchivoPorDefecto = "reporte.pdf"

// RegistrarExportarReportePDF agrega al mux las dos descargas de HU-27. Son rutas /web/ porque las pide
// el navegador al seguir el enlace de la pantalla del reporte: los errores se muestran como página (P-7).
// Son GET de solo lectura, así que no pasan por el chequeo de origen de los formularios.
func RegistrarExportarReportePDF(mux *http.ServeMux, proyecto exportadorReporteProyecto, sprint exportadorReporteSprint) {
	p := nuevasPantallas(CasosDeUso{}) // solo para la página de error con el layout común
	mux.HandleFunc("GET /web/proyectos/{id}/reporte/pdf", p.descargarReporteProyecto(proyecto))
	mux.HandleFunc("GET /web/proyectos/{id}/sprints/{sprintID}/reporte/pdf", p.descargarReporteSprint(sprint))
}

// descargarReporteProyecto atiende GET /web/proyectos/{id}/reporte/pdf: 200 con el PDF, 400 si el ID no
// es válido, 404 si el proyecto no existe y 500 ante cualquier otra falla. Implementa HU-27 (CA-27.14).
func (p *pantallas) descargarReporteProyecto(uc exportadorReporteProyecto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, ok := p.idDeRuta(w, r)
		if !ok {
			return
		}
		archivo, err := uc.Ejecutar(r.Context(), proyectoID)
		if err != nil {
			p.errorDeReportePDF(w, err, fmt.Sprintf("proyecto %d", proyectoID))
			return
		}
		responderPDF(w, archivo)
	}
}

// descargarReporteSprint atiende GET /web/proyectos/{id}/sprints/{sprintID}/reporte/pdf. El ID del
// proyecto se valida antes que el del sprint (sección 7). Implementa HU-27 (CA-27.15).
func (p *pantallas) descargarReporteSprint(uc exportadorReporteSprint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, ok := p.idDeRuta(w, r)
		if !ok {
			return
		}
		sprintID, err := sprintIDDesdeRuta(r)
		if err != nil {
			p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
			return
		}
		archivo, err := uc.Ejecutar(r.Context(), proyectoID, sprintID)
		if err != nil {
			p.errorDeReportePDF(w, err, fmt.Sprintf("sprint %d del proyecto %d", sprintID, proyectoID))
			return
		}
		responderPDF(w, archivo)
	}
}

// errorDeReportePDF muestra la página 404 si no existe el proyecto o el sprint, o la de falla interna,
// que deja el detalle solo en el log. alcance dice qué reporte se estaba exportando.
func (p *pantallas) errorDeReportePDF(w http.ResponseWriter, err error, alcance string) {
	for _, noEncontrado := range erroresReportePDF {
		if errors.Is(err, noEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{noEncontrado.Error()})
			return
		}
	}
	p.fallaInterna(w, "exportar el reporte PDF del "+alcance, err)
}

// responderPDF escribe la descarga. El PDF ya está entero en memoria (RN-27.15), así que Content-Length
// es exacto y nunca sale un archivo a medias.
func responderPDF(w http.ResponseWriter, archivo app.ArchivoPDF) {
	h := w.Header()
	h.Set("Content-Type", "application/pdf")
	h.Set("Content-Disposition", `attachment; filename="`+nombreArchivoSeguro(archivo.Nombre)+`"`)
	h.Set("Content-Length", strconv.Itoa(len(archivo.Contenido)))
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(archivo.Contenido); err != nil {
		log.Printf("escribir el reporte PDF %s: %v", archivo.Nombre, err)
	}
}

// nombreArchivoSeguro deja solo los caracteres de RN-27.12 (minúsculas sin acentos, dígitos, "-" y ".")
// y reemplaza el resto por "-". El caso de uso ya da un nombre así; esto es una defensa del adaptador
// para que comillas o saltos de línea nunca corten la cabecera ni agreguen otras.
func nombreArchivoSeguro(nombre string) string {
	seguro := strings.Map(func(c rune) rune {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '.' {
			return c
		}
		return '-'
	}, nombre)
	if seguro == "" {
		return nombreArchivoPorDefecto
	}
	return seguro
}
