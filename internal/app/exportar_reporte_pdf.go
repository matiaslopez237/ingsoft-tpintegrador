package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// FuenteDocumentoProyecto arma el documento del reporte de un proyecto. SUPUESTO: la implementa HU-26,
// que convierte su reporte al modelo neutral (spec HU-27, sección 5).
type FuenteDocumentoProyecto interface {
	DocumentoDeProyecto(ctx context.Context, proyectoID int64) (domain.DocumentoPDF, error)
}

// FuenteDocumentoSprint arma el documento del reporte de un sprint. SUPUESTO: la implementa HU-26.
type FuenteDocumentoSprint interface {
	DocumentoDeSprint(ctx context.Context, proyectoID, sprintID int64) (domain.DocumentoPDF, error)
}

// GeneradorPDF genera los bytes del PDF a partir de un documento ya normalizado y validado. Lo
// implementa el adaptador internal/adapters/pdf (ADR-0006).
type GeneradorPDF interface {
	Generar(ctx context.Context, doc domain.DocumentoPDF) ([]byte, error)
}

// ArchivoPDF es lo que se descarga: el nombre de RN-27.12 y los bytes del PDF.
type ArchivoPDF struct {
	Nombre    string
	Contenido []byte
}

// ExportarReporteProyecto es el caso de uso de HU-27 para el reporte de un proyecto.
type ExportarReporteProyecto struct {
	proyectos LectorProyectos
	fuente    FuenteDocumentoProyecto
	generador GeneradorPDF
	reloj     func() time.Time
}

// NuevoExportarReporteProyecto recibe el reloj para que la emisión se pueda fijar en los tests, igual
// que HU-20.
func NuevoExportarReporteProyecto(proyectos LectorProyectos, fuente FuenteDocumentoProyecto,
	generador GeneradorPDF, reloj func() time.Time) *ExportarReporteProyecto {
	return &ExportarReporteProyecto{proyectos: proyectos, fuente: fuente, generador: generador, reloj: reloj}
}

// Ejecutar verifica que el proyecto exista (RN-27.13), arma el documento, lo prepara (RN-27.8 a
// RN-27.11) y genera el PDF. Ante cualquier error no devuelve ningún archivo.
func (uc *ExportarReporteProyecto) Ejecutar(ctx context.Context, proyectoID int64) (ArchivoPDF, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return ArchivoPDF{}, err
	}
	doc, err := uc.fuente.DocumentoDeProyecto(ctx, proyectoID)
	if err != nil {
		return ArchivoPDF{}, fmt.Errorf("armar el reporte del proyecto %d: %w", proyectoID, err)
	}
	contenido, emision, err := generarPDF(ctx, uc.generador, uc.reloj, doc)
	if err != nil {
		return ArchivoPDF{}, fmt.Errorf("reporte del proyecto %d: %w", proyectoID, err)
	}
	return ArchivoPDF{Nombre: domain.NombreArchivoReporteProyecto(proyectoID, emision), Contenido: contenido}, nil
}

// ExportarReporteSprint es el caso de uso de HU-27 para el reporte de un sprint.
type ExportarReporteSprint struct {
	proyectos LectorProyectos
	sprints   LectorSprintPorID
	fuente    FuenteDocumentoSprint
	generador GeneradorPDF
	reloj     func() time.Time
}

func NuevoExportarReporteSprint(proyectos LectorProyectos, sprints LectorSprintPorID, fuente FuenteDocumentoSprint,
	generador GeneradorPDF, reloj func() time.Time) *ExportarReporteSprint {
	return &ExportarReporteSprint{proyectos: proyectos, sprints: sprints, fuente: fuente, generador: generador, reloj: reloj}
}

// Ejecutar verifica que el proyecto exista y que el sprint sea suyo (RN-27.13), en el orden de la tabla
// de errores de la spec, y después hace lo mismo que el del proyecto.
func (uc *ExportarReporteSprint) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (ArchivoPDF, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return ArchivoPDF{}, err
	}
	s, err := uc.sprints.BuscarSprintPorID(ctx, sprintID)
	if err != nil {
		return ArchivoPDF{}, err
	}
	if s.ProyectoID != proyectoID { // RN-27.13: de otro proyecto es "inexistente"
		return ArchivoPDF{}, domain.ErrSprintNoEncontrado
	}
	doc, err := uc.fuente.DocumentoDeSprint(ctx, proyectoID, sprintID)
	if err != nil {
		return ArchivoPDF{}, fmt.Errorf("armar el reporte del sprint %d: %w", sprintID, err)
	}
	contenido, emision, err := generarPDF(ctx, uc.generador, uc.reloj, doc)
	if err != nil {
		return ArchivoPDF{}, fmt.Errorf("reporte del sprint %d: %w", sprintID, err)
	}
	return ArchivoPDF{Nombre: domain.NombreArchivoReporteSprint(proyectoID, sprintID, emision), Contenido: contenido}, nil
}

// generarPDF es lo común a los dos casos de uso: pone la emisión (RN-27.10), normaliza y completa el
// aviso (RN-27.8, RN-27.9), valida (RN-27.11) y genera. Devuelve la emisión, que da la fecha del nombre
// del archivo. Las reglas viven en el dominio; acá solo se ordenan.
func generarPDF(ctx context.Context, generador GeneradorPDF, reloj func() time.Time, doc domain.DocumentoPDF) ([]byte, time.Time, error) {
	doc.Emision = domain.HoraDeEmision(reloj())
	doc, n := domain.NormalizarDocumento(doc)
	doc.Aviso = domain.TextoAvisoReemplazos(n)
	if err := doc.Validar(); err != nil {
		return nil, time.Time{}, err
	}
	contenido, err := generador.Generar(ctx, doc)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("generar el PDF: %w", err)
	}
	return contenido, doc.Emision, nil
}
