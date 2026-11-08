package http

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// La fuente implementa los dos puertos que usan los casos de uso de exportación de HU-27.
var (
	_ app.FuenteDocumentoProyecto = (*fuenteDocumentoReporte)(nil)
	_ app.FuenteDocumentoSprint   = (*fuenteDocumentoReporte)(nil)
)

// pedidoArmador es lo que recibió el último caso de uso de HU-26 llamado.
type pedidoArmador struct {
	ctx                  context.Context
	proyectoID, sprintID int64
	llamado              bool
}

// armadoresQueAnotan devuelve los dos casos de uso de HU-26, que anotan el contexto y los IDs que
// reciben y responden con el reporte o el error indicados.
func armadoresQueAnotan(reporte domain.Reporte, err error) (*pedidoArmador, reporteProyectoFunc, reporteSprintFunc) {
	pedido := &pedidoArmador{}
	proyecto := reporteProyectoFunc(func(ctx context.Context, id int64) (domain.Reporte, error) {
		*pedido = pedidoArmador{ctx: ctx, proyectoID: id, llamado: true}
		return reporte, err
	})
	sprint := reporteSprintFunc(func(ctx context.Context, proyectoID, sprintID int64) (domain.Reporte, error) {
		*pedido = pedidoArmador{ctx: ctx, proyectoID: proyectoID, sprintID: sprintID, llamado: true}
		return reporte, err
	})
	return pedido, proyecto, sprint
}

// Traza: HU-27 CA-27.25 (RN-27.16; P-11): la fuente arma el reporte del proyecto con el caso de
// uso de HU-26, con el contexto y el ID del pedido, y lo convierte con aDocumentoPDF.
func TestFuenteDocumentoReporte_DocumentoDeProyecto(t *testing.T) {
	reporte := reporteParaDocumento()
	pedido, proyecto, sprint := armadoresQueAnotan(reporte, nil)
	fuente := NuevaFuenteDocumentoReporte(proyecto, sprint)
	ctx := context.WithValue(context.Background(), claveContextoPDF{}, "pedido")

	doc, err := fuente.DocumentoDeProyecto(ctx, 7)

	if err != nil {
		t.Fatalf("DocumentoDeProyecto() error = %v", err)
	}
	if !pedido.llamado || pedido.proyectoID != 7 || pedido.sprintID != 0 {
		t.Errorf("caso de uso llamado = %v con proyecto %d y sprint %d; se esperaba solo el del proyecto, con 7",
			pedido.llamado, pedido.proyectoID, pedido.sprintID)
	}
	if pedido.ctx == nil || pedido.ctx.Value(claveContextoPDF{}) != "pedido" {
		t.Error("el caso de uso de HU-26 no recibió el contexto del pedido")
	}
	if want := aDocumentoPDF(aDatosReporte(reporte)); !reflect.DeepEqual(doc, want) {
		t.Errorf("DocumentoDeProyecto() = %+v, se esperaba la conversión del reporte %+v", doc, want)
	}
}

// Traza: HU-27 CA-27.25 (RN-27.16; P-11): lo mismo para el reporte de un Sprint, con los dos IDs.
func TestFuenteDocumentoReporte_DocumentoDeSprint(t *testing.T) {
	reporte := reporteDeSprintDeEjemplo()
	pedido, proyecto, sprint := armadoresQueAnotan(reporte, nil)
	fuente := NuevaFuenteDocumentoReporte(proyecto, sprint)
	ctx := context.WithValue(context.Background(), claveContextoPDF{}, "pedido")

	doc, err := fuente.DocumentoDeSprint(ctx, 7, 9)

	if err != nil {
		t.Fatalf("DocumentoDeSprint() error = %v", err)
	}
	if !pedido.llamado || pedido.proyectoID != 7 || pedido.sprintID != 9 {
		t.Errorf("caso de uso llamado = %v con proyecto %d y sprint %d; se esperaba 7 y 9",
			pedido.llamado, pedido.proyectoID, pedido.sprintID)
	}
	if pedido.ctx == nil || pedido.ctx.Value(claveContextoPDF{}) != "pedido" {
		t.Error("el caso de uso de HU-26 no recibió el contexto del pedido")
	}
	if want := aDocumentoPDF(aDatosReporte(reporte)); !reflect.DeepEqual(doc, want) {
		t.Errorf("DocumentoDeSprint() = %+v, se esperaba la conversión del reporte %+v", doc, want)
	}
}

// Traza: HU-27 CA-27.17, CA-27.18, CA-27.20 (RN-27.13, sección 7; P-11): la fuente devuelve el error del
// caso de uso de HU-26 sin cambios y sin documento, así que errors.Is sigue reconociendo los "no
// encontrado" (que el adaptador HTTP responde con 404) y cualquier otra falla sigue siendo interna.
func TestFuenteDocumentoReporte_ErrorDelCasoDeUso_SinCambios(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
		es     error
	}{
		{"proyecto inexistente", fmt.Errorf("buscar proyecto 99: %w", domain.ErrProyectoNoEncontrado), domain.ErrProyectoNoEncontrado},
		{"sprint inexistente", domain.ErrSprintNoEncontrado, domain.ErrSprintNoEncontrado},
		{"falla de persistencia", errors.New("disco lleno"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, proyecto, sprint := armadoresQueAnotan(reporteParaDocumento(), tt.err)
			fuente := NuevaFuenteDocumentoReporte(proyecto, sprint)

			docProyecto, errProyecto := fuente.DocumentoDeProyecto(context.Background(), 7)
			docSprint, errSprint := fuente.DocumentoDeSprint(context.Background(), 7, 9)

			for nombre, got := range map[string]error{"proyecto": errProyecto, "sprint": errSprint} {
				if got != tt.err { // el mismo error, sin envolver: no alcanza con errors.Is
					t.Errorf("%s: error = %v, se esperaba el del caso de uso sin cambios: %v", nombre, got, tt.err)
				}
				if tt.es != nil && !errors.Is(got, tt.es) {
					t.Errorf("%s: errors.Is(%v, %v) = false", nombre, got, tt.es)
				}
			}
			if !reflect.DeepEqual(docProyecto, domain.DocumentoPDF{}) || !reflect.DeepEqual(docSprint, domain.DocumentoPDF{}) {
				t.Errorf("con error, el documento tiene que quedar vacío: %+v, %+v", docProyecto, docSprint)
			}
		})
	}
}
