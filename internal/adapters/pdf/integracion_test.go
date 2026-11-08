package pdf

import (
	"context"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// proyectoExistente es el lector de proyectos de la prueba de punta a punta: el proyecto siempre existe.
type proyectoExistente struct{}

func (proyectoExistente) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	return domain.Proyecto{ID: id, Nombre: "Alfa"}, nil
}

// fuenteFija devuelve siempre el mismo documento, como lo armaría HU-26 (sin Emision ni Aviso).
type fuenteFija struct{ doc domain.DocumentoPDF }

func (f fuenteFija) DocumentoDeProyecto(context.Context, int64) (domain.DocumentoPDF, error) {
	return f.doc, nil
}

// Traza: HU-27 CA-27.8 (RN-27.8, RN-27.9, sección 5): de punta a punta, con el caso de uso y el generador
// reales y sin dobles del generador. El emoji y los ideogramas llegan sin normalizar desde la fuente; el
// PDF que se descarga dice "?" en su lugar y lleva el aviso con la cantidad de reemplazos.
func TestExportarConGeneradorReal_Emojis_PDFConSignoYAviso(t *testing.T) {
	reloj := func() time.Time { return time.Date(2026, 11, 7, 17, 5, 0, 0, time.UTC) }
	casos := []struct {
		nombre, titulo, wantTitulo, wantAviso string
	}{
		{"un emoji", "Sprint 😀 final", "Sprint ? final",
			"Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?»."},
		{"emojis seguidos e ideogramas", "Equipo 😀😀 y 中文", "Equipo ? y ?",
			"Aviso: 2 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?»."},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			fuente := fuenteFija{doc: domain.DocumentoPDF{
				Titulo:    c.titulo,
				Secciones: []domain.SeccionPDF{{Titulo: "Resumen", Bloques: []domain.BloquePDF{{Parrafo: "Sin novedades."}}}},
			}}
			uc := app.NuevoExportarReporteProyecto(proyectoExistente{}, fuente, NuevoGenerador(), reloj)

			archivo, err := uc.Ejecutar(context.Background(), 3)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if archivo.Nombre != "reporte-proyecto-3-2026-11-07.pdf" {
				t.Errorf("nombre = %q", archivo.Nombre)
			}
			textos := paginas(t, archivo.Contenido)
			if len(textos) != 1 {
				t.Fatalf("páginas = %d, se esperaba 1", len(textos))
			}
			contieneEnOrden(t, textos[0], c.wantTitulo, "Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3).",
				c.wantAviso, "Resumen", "Sin novedades.")
		})
	}
}
