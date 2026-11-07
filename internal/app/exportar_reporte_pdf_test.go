package app

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// exportarFalso implementa todos los puertos de los casos de uso de HU-27, con lo que cada test prepara.
// La fuente es un doble del armado de HU-26 (SUPUESTO: la real la escribe HU-26).
type exportarFalso struct {
	proyectos map[int64]bool
	sprints   map[int64]domain.Sprint
	documento domain.DocumentoPDF

	errProyectos, errSprints, errFuente, errGenerador error

	buscoSprint, llamadasFuente, llamadasGenerador int
	recibido                                       domain.DocumentoPDF
	contextos                                      []context.Context
}

var relojDePrueba = time.Date(2026, 11, 7, 17, 5, 0, 0, time.UTC)

func nuevoExportarFalso() *exportarFalso {
	return &exportarFalso{
		proyectos: map[int64]bool{3: true, 4: true},
		sprints: map[int64]domain.Sprint{
			12: {ID: 12, ProyectoID: 3},
			13: {ID: 13, ProyectoID: 4},
		},
		documento: domain.DocumentoPDF{
			Titulo:    "Reporte del proyecto Alfa",
			Secciones: []domain.SeccionPDF{{Titulo: "Resumen", Bloques: []domain.BloquePDF{{Parrafo: "Activo"}}}},
		},
	}
}

func (f *exportarFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyectos != nil {
		return domain.Proyecto{}, f.errProyectos
	}
	if !f.proyectos[id] {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return domain.Proyecto{ID: id}, nil
}

func (f *exportarFalso) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	f.buscoSprint++
	if f.errSprints != nil {
		return domain.Sprint{}, f.errSprints
	}
	s, ok := f.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return s, nil
}

func (f *exportarFalso) DocumentoDeProyecto(ctx context.Context, _ int64) (domain.DocumentoPDF, error) {
	return f.devolverDocumento(ctx)
}

func (f *exportarFalso) DocumentoDeSprint(ctx context.Context, _, _ int64) (domain.DocumentoPDF, error) {
	return f.devolverDocumento(ctx)
}

func (f *exportarFalso) devolverDocumento(ctx context.Context) (domain.DocumentoPDF, error) {
	f.llamadasFuente++
	f.contextos = append(f.contextos, ctx)
	if f.errFuente != nil {
		return domain.DocumentoPDF{}, f.errFuente
	}
	return f.documento, nil
}

func (f *exportarFalso) Generar(ctx context.Context, doc domain.DocumentoPDF) ([]byte, error) {
	f.llamadasGenerador++
	f.contextos = append(f.contextos, ctx)
	f.recibido = doc
	if f.errGenerador != nil {
		return nil, f.errGenerador
	}
	return []byte("%PDF-de prueba"), nil
}

func (f *exportarFalso) reloj() time.Time { return relojDePrueba }

func exportarProyecto(f *exportarFalso, proyectoID int64) (ArchivoPDF, error) {
	return NuevoExportarReporteProyecto(f, f, f, f.reloj).Ejecutar(context.Background(), proyectoID)
}

func exportarSprint(f *exportarFalso, proyectoID, sprintID int64) (ArchivoPDF, error) {
	return NuevoExportarReporteSprint(f, f, f, f, f.reloj).Ejecutar(context.Background(), proyectoID, sprintID)
}

// CA-27.1, CA-27.14: el generador recibe el documento con la emisión en UTC-3 y el archivo lleva sus
// bytes y el nombre de RN-27.12.
func TestExportarReporteProyecto_GeneraElPDFDelDocumento(t *testing.T) {
	f := nuevoExportarFalso()
	archivo, err := exportarProyecto(f, 3)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if archivo.Nombre != "reporte-proyecto-3-2026-11-07.pdf" {
		t.Errorf("nombre = %q", archivo.Nombre)
	}
	if !bytes.Equal(archivo.Contenido, []byte("%PDF-de prueba")) {
		t.Errorf("contenido = %q, se esperaban los bytes del generador", archivo.Contenido)
	}
	want := f.documento
	want.Emision = domain.HoraDeEmision(relojDePrueba)
	if !reflect.DeepEqual(f.recibido, want) {
		t.Errorf("el generador recibió %+v, se esperaba %+v", f.recibido, want)
	}
	if _, offset := f.recibido.Emision.Zone(); offset != -3*60*60 {
		t.Errorf("la emisión no está en UTC-3: %v", f.recibido.Emision)
	}
}

// CA-27.15: el reporte de un sprint del proyecto lleva el nombre con los dos IDs.
func TestExportarReporteSprint_GeneraElPDFDelDocumento(t *testing.T) {
	f := nuevoExportarFalso()
	archivo, err := exportarSprint(f, 3, 12)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if archivo.Nombre != "reporte-proyecto-3-sprint-12-2026-11-07.pdf" {
		t.Errorf("nombre = %q", archivo.Nombre)
	}
	if f.llamadasGenerador != 1 || !bytes.Equal(archivo.Contenido, []byte("%PDF-de prueba")) {
		t.Errorf("generador llamado %d veces, contenido %q", f.llamadasGenerador, archivo.Contenido)
	}
}

// CA-27.8, CA-27.9: el caso de uso normaliza el documento y completa el aviso con la cantidad de
// reemplazos (RN-27.8, RN-27.9).
func TestExportarReporte_NormalizaYCompletaElAviso(t *testing.T) {
	f := nuevoExportarFalso()
	f.documento.Titulo = "Sprint 😀 final"
	f.documento.Secciones[0].Titulo = "Equipo 😀😀"
	if _, err := exportarProyecto(f, 3); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if f.recibido.Titulo != "Sprint ? final" || f.recibido.Secciones[0].Titulo != "Equipo ?" {
		t.Errorf("no se normalizó: %q, %q", f.recibido.Titulo, f.recibido.Secciones[0].Titulo)
	}
	if want := domain.TextoAvisoReemplazos(2); f.recibido.Aviso != want {
		t.Errorf("aviso = %q, se esperaba %q", f.recibido.Aviso, want)
	}
}

// RN-27.9: el aviso lo escribe siempre el caso de uso; si la fuente trajera uno, no se usa.
func TestExportarReporte_SinReemplazos_AvisoVacioAunqueLaFuenteTraigaUno(t *testing.T) {
	f := nuevoExportarFalso()
	f.documento.Aviso = "aviso que HU-26 no debería poner"
	if _, err := exportarSprint(f, 3, 12); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if f.recibido.Aviso != "" {
		t.Errorf("aviso = %q, se esperaba vacío", f.recibido.Aviso)
	}
}

// CA-27.11: un documento mal armado no llega al generador; el título se valida después de normalizar.
func TestExportarReporte_DocumentoInvalido_NoLlegaAlGenerador(t *testing.T) {
	casos := map[string]string{"título vacío": "", "título de invisibles": "​​"}
	for nombre, titulo := range casos {
		t.Run(nombre, func(t *testing.T) {
			f := nuevoExportarFalso()
			f.documento.Titulo = titulo
			archivo, err := exportarProyecto(f, 3)
			if !errors.Is(err, domain.ErrDocumentoInvalido) {
				t.Errorf("error = %v, se esperaba ErrDocumentoInvalido", err)
			}
			sinArchivoNiGenerador(t, f, archivo)
		})
	}
}

// CA-27.17: con un proyecto inexistente no se busca el sprint, ni se arma ni se genera nada.
func TestExportarReporte_ProyectoInexistente(t *testing.T) {
	exportaciones := map[string]func(f *exportarFalso) (ArchivoPDF, error){
		"proyecto": func(f *exportarFalso) (ArchivoPDF, error) { return exportarProyecto(f, 99) },
		"sprint":   func(f *exportarFalso) (ArchivoPDF, error) { return exportarSprint(f, 99, 12) },
	}
	for nombre, exportar := range exportaciones {
		t.Run(nombre, func(t *testing.T) {
			f := nuevoExportarFalso()
			archivo, err := exportar(f)
			if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("error = %v, se esperaba ErrProyectoNoEncontrado", err)
			}
			if f.buscoSprint != 0 || f.llamadasFuente != 0 {
				t.Errorf("buscó el sprint %d veces y armó %d documentos", f.buscoSprint, f.llamadasFuente)
			}
			sinArchivoNiGenerador(t, f, archivo)
		})
	}
}

// CA-27.18: un sprint inexistente o de otro proyecto es "sprint no encontrado" (RN-27.13).
func TestExportarReporteSprint_SprintInexistenteODeOtroProyecto(t *testing.T) {
	for nombre, sprintID := range map[string]int64{"inexistente": 99, "de otro proyecto": 13} {
		t.Run(nombre, func(t *testing.T) {
			f := nuevoExportarFalso()
			archivo, err := exportarSprint(f, 3, sprintID)
			if !errors.Is(err, domain.ErrSprintNoEncontrado) {
				t.Errorf("error = %v, se esperaba ErrSprintNoEncontrado", err)
			}
			if f.llamadasFuente != 0 {
				t.Errorf("armó %d documentos", f.llamadasFuente)
			}
			sinArchivoNiGenerador(t, f, archivo)
		})
	}
}

// CA-27.20: una falla de la fuente o del generador llega envuelta, sin archivo.
func TestExportarReporte_FallaDeLaFuenteODelGenerador_SeEnvuelveSinArchivo(t *testing.T) {
	falla := errors.New("falla de prueba")
	casos := map[string]func(f *exportarFalso){
		"fuente":    func(f *exportarFalso) { f.errFuente = falla },
		"generador": func(f *exportarFalso) { f.errGenerador = falla },
	}
	for nombre, preparar := range casos {
		for exportacion, exportar := range map[string]func(f *exportarFalso) (ArchivoPDF, error){
			"proyecto": func(f *exportarFalso) (ArchivoPDF, error) { return exportarProyecto(f, 3) },
			"sprint":   func(f *exportarFalso) (ArchivoPDF, error) { return exportarSprint(f, 3, 12) },
		} {
			t.Run(nombre+"/"+exportacion, func(t *testing.T) {
				f := nuevoExportarFalso()
				preparar(f)
				archivo, err := exportar(f)
				if !errors.Is(err, falla) {
					t.Errorf("error = %v, se esperaba la falla envuelta", err)
				}
				if err != nil && err == falla {
					t.Errorf("la falla no se envolvió con contexto: %v", err)
				}
				if archivo.Nombre != "" || archivo.Contenido != nil {
					t.Errorf("devolvió un archivo: %+v", archivo)
				}
			})
		}
	}
}

// Una falla al buscar el proyecto o el sprint se devuelve tal cual, como en HU-23.
func TestExportarReporte_FallaDeLosLectores_SePropaga(t *testing.T) {
	falla := errors.New("falla de la base")
	f := nuevoExportarFalso()
	f.errProyectos = falla
	if _, err := exportarProyecto(f, 3); !errors.Is(err, falla) {
		t.Errorf("proyecto: error = %v", err)
	}
	f = nuevoExportarFalso()
	f.errSprints = falla
	if _, err := exportarSprint(f, 3, 12); !errors.Is(err, falla) {
		t.Errorf("sprint: error = %v", err)
	}
	if f.llamadasFuente != 0 || f.llamadasGenerador != 0 {
		t.Errorf("siguió después de la falla: fuente %d, generador %d", f.llamadasFuente, f.llamadasGenerador)
	}
}

// La fuente y el generador reciben el contexto del pedido, para que cancelarlo corte el trabajo.
func TestExportarReporte_LosPuertosRecibenElContexto(t *testing.T) {
	type clave struct{}
	ctx := context.WithValue(context.Background(), clave{}, "pedido")
	f := nuevoExportarFalso()
	if _, err := NuevoExportarReporteSprint(f, f, f, f, f.reloj).Ejecutar(ctx, 3, 12); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(f.contextos) != 2 {
		t.Fatalf("se esperaban 2 llamadas con contexto, hubo %d", len(f.contextos))
	}
	for _, c := range f.contextos {
		if c.Value(clave{}) != "pedido" {
			t.Errorf("un puerto no recibió el contexto del pedido")
		}
	}
}

func sinArchivoNiGenerador(t *testing.T, f *exportarFalso, archivo ArchivoPDF) {
	t.Helper()
	if f.llamadasGenerador != 0 {
		t.Errorf("el generador se llamó %d veces", f.llamadasGenerador)
	}
	if archivo.Nombre != "" || archivo.Contenido != nil {
		t.Errorf("devolvió un archivo: %+v", archivo)
	}
}
