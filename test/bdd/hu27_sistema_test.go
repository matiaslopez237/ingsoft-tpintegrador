package bdd

import (
	"context"
	"errors"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu27Sistema conecta los escenarios de HU-27 a los casos de uso ExportarReporteProyecto y
// ExportarReporteSprint (spec HU-27, sección 5). Sus puertos se reemplazan por dobles en memoria:
//   - los lectores de proyectos y de sprints, y la fuente de documentos (FuenteDocumentoProyecto y
//     FuenteDocumentoSprint), en hu27Fuente. SUPUESTO: la fuente real la escribe HU-26, que todavía no
//     existe; acá cada proyecto y cada sprint tienen un DocumentoPDF armado por el escenario;
//   - el generador de PDF (GeneradorPDF), que guarda lo que recibe y devuelve bytes fijos.
//
// No usa SQLite: HU-27 no guarda nada y el reporte de HU-26 no existe todavía.
type hu27Sistema struct {
	fuente    *hu27Fuente
	generador *hu27Generador
	ahora     time.Time // lo que devuelve el reloj de los casos de uso; lo cambia el escenario
	proyecto  *app.ExportarReporteProyecto
	sprint    *app.ExportarReporteSprint
}

func hu27NuevoSistema() *hu27Sistema {
	s := &hu27Sistema{
		fuente:    &hu27Fuente{proyectos: map[int64]*domain.DocumentoPDF{}, sprints: map[int64]hu27SprintEnMemoria{}},
		generador: &hu27Generador{},
	}
	reloj := func() time.Time { return s.ahora }
	// La fuente en memoria también hace de lector de proyectos y de sprints, con los mismos datos.
	s.proyecto = app.NuevoExportarReporteProyecto(s.fuente, s.fuente, s.generador, reloj)
	s.sprint = app.NuevoExportarReporteSprint(s.fuente, s.fuente, s.fuente, s.generador, reloj)
	return s
}

func (s *hu27Sistema) ExportarProyecto(proyectoID int64) (app.ArchivoPDF, error) {
	return s.proyecto.Ejecutar(context.Background(), proyectoID)
}

func (s *hu27Sistema) ExportarSprint(proyectoID, sprintID int64) (app.ArchivoPDF, error) {
	return s.sprint.Ejecutar(context.Background(), proyectoID, sprintID)
}

// errHU27FallaInterna es la falla que el escenario le pide a un doble. No es un error de la spec.
var errHU27FallaInterna = errors.New("falla de prueba de HU-27")

// --- Fuente de documentos en memoria ---

type hu27SprintEnMemoria struct {
	proyectoID int64
	doc        *domain.DocumentoPDF
}

// hu27Fuente da los documentos y hace de lector de proyectos y de sprints (los puertos de HU-04 y HU-23
// que reutiliza el caso de uso). El caso de uso comprueba la existencia con los lectores (RN-27.13); la
// fuente igual devuelve esos errores, como haría la de HU-26, aunque con los lectores no se llega.
type hu27Fuente struct {
	proyectos map[int64]*domain.DocumentoPDF
	sprints   map[int64]hu27SprintEnMemoria
	falla     error
}

func (f *hu27Fuente) BuscarPorID(_ context.Context, proyectoID int64) (domain.Proyecto, error) {
	if _, ok := f.proyectos[proyectoID]; !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return domain.Proyecto{ID: proyectoID}, nil
}

func (f *hu27Fuente) BuscarSprintPorID(_ context.Context, sprintID int64) (domain.Sprint, error) {
	sp, ok := f.sprints[sprintID]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return domain.Sprint{ID: sprintID, ProyectoID: sp.proyectoID}, nil
}

func (f *hu27Fuente) DocumentoDeProyecto(_ context.Context, proyectoID int64) (domain.DocumentoPDF, error) {
	if f.falla != nil {
		return domain.DocumentoPDF{}, f.falla
	}
	doc, ok := f.proyectos[proyectoID]
	if !ok {
		return domain.DocumentoPDF{}, domain.ErrProyectoNoEncontrado
	}
	return hu27Copiar(*doc), nil
}

func (f *hu27Fuente) DocumentoDeSprint(_ context.Context, proyectoID, sprintID int64) (domain.DocumentoPDF, error) {
	if f.falla != nil {
		return domain.DocumentoPDF{}, f.falla
	}
	if _, ok := f.proyectos[proyectoID]; !ok {
		return domain.DocumentoPDF{}, domain.ErrProyectoNoEncontrado
	}
	sp, ok := f.sprints[sprintID]
	if !ok || sp.proyectoID != proyectoID {
		return domain.DocumentoPDF{}, domain.ErrSprintNoEncontrado
	}
	return hu27Copiar(*sp.doc), nil
}

// --- Generador de PDF doble ---

// hu27Generador guarda una copia del documento que recibe, para que los escenarios lo miren, y devuelve
// bytes que dependen del título, para comprobar que el caso de uso entrega lo que devolvió el generador.
type hu27Generador struct {
	llamadas int
	recibido domain.DocumentoPDF
	devuelto []byte
	falla    error
}

func (g *hu27Generador) Generar(_ context.Context, doc domain.DocumentoPDF) ([]byte, error) {
	g.llamadas++
	g.recibido = hu27Copiar(doc)
	if g.falla != nil {
		return nil, g.falla
	}
	g.devuelto = []byte("%PDF-doble de HU-27: " + doc.Titulo)
	return append([]byte(nil), g.devuelto...), nil
}

// hu27Copiar copia el documento entero, para que ni el caso de uso ni el escenario cambien el original
// del otro.
func hu27Copiar(d domain.DocumentoPDF) domain.DocumentoPDF {
	c := d
	c.Secciones = nil
	for _, s := range d.Secciones {
		cs := domain.SeccionPDF{Titulo: s.Titulo}
		for _, b := range s.Bloques {
			cb := domain.BloquePDF{Parrafo: b.Parrafo}
			if b.Tabla != nil {
				t := domain.TablaPDF{Columnas: append([]string(nil), b.Tabla.Columnas...)}
				for _, f := range b.Tabla.Filas {
					t.Filas = append(t.Filas, append([]string(nil), f...))
				}
				cb.Tabla = &t
			}
			cs.Bloques = append(cs.Bloques, cb)
		}
		c.Secciones = append(c.Secciones, cs)
	}
	return c
}
