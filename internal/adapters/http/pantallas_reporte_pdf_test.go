package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const (
	rutaPDFProyecto = "/web/proyectos/3/reporte/pdf"
	rutaPDFSprint   = "/web/proyectos/3/sprints/12/reporte/pdf"
)

// contenidoPDFFalso tiene bytes que no son texto (un 0 y un 0xFF) para comprobar que el cuerpo sale
// tal cual, sin conversiones.
var contenidoPDFFalso = []byte("%PDF-1.7\n\x00\xff contenido binario\n%%EOF\n")

// exportadorReporteProyectoFalso es el doble de ExportarReporteProyecto: devuelve lo que se le indica y
// guarda con qué se lo llamó.
type exportadorReporteProyectoFalso struct {
	archivo    app.ArchivoPDF
	err        error
	llamado    bool
	ctx        context.Context
	proyectoID int64
}

func (f *exportadorReporteProyectoFalso) Ejecutar(ctx context.Context, proyectoID int64) (app.ArchivoPDF, error) {
	f.llamado, f.ctx, f.proyectoID = true, ctx, proyectoID
	return f.archivo, f.err
}

// exportadorReporteSprintFalso es el doble de ExportarReporteSprint.
type exportadorReporteSprintFalso struct {
	archivo              app.ArchivoPDF
	err                  error
	llamado              bool
	ctx                  context.Context
	proyectoID, sprintID int64
}

func (f *exportadorReporteSprintFalso) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (app.ArchivoPDF, error) {
	f.llamado, f.ctx, f.proyectoID, f.sprintID = true, ctx, proyectoID, sprintID
	return f.archivo, f.err
}

// exportadoresConError arma los dos dobles con el mismo error.
func exportadoresConError(err error) (*exportadorReporteProyectoFalso, *exportadorReporteSprintFalso) {
	return &exportadorReporteProyectoFalso{err: err}, &exportadorReporteSprintFalso{err: err}
}

// pedirReportePDF registra solo las rutas de HU-27 y hace el GET con el contexto indicado.
func pedirReportePDF(t *testing.T, ctx context.Context, proyecto exportadorReporteProyecto, sprint exportadorReporteSprint, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarExportarReportePDF(mux, proyecto, sprint)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequestWithContext(ctx, http.MethodGet, ruta, nil))
	return rec
}

// rutasReportePDF son las dos descargas, para los casos que valen igual en las dos.
var rutasReportePDF = []struct{ nombre, ruta string }{
	{"proyecto", rutaPDFProyecto},
	{"sprint", rutaPDFSprint},
}

// esDescargaPDF verifica las cabeceras y el cuerpo de una descarga exitosa.
func esDescargaPDF(t *testing.T, rec *httptest.ResponseRecorder, nombreArchivo string) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	h := rec.Header()
	want := map[string]string{
		"Content-Type":           "application/pdf",
		"Content-Disposition":    `attachment; filename="` + nombreArchivo + `"`,
		"Content-Length":         strconv.Itoa(len(contenidoPDFFalso)),
		"X-Content-Type-Options": "nosniff",
	}
	for k, v := range want {
		if got := h.Values(k); len(got) != 1 || got[0] != v {
			t.Errorf("%s = %q, se esperaba [%q]", k, got, v)
		}
	}
	if !bytes.Equal(rec.Body.Bytes(), contenidoPDFFalso) {
		t.Errorf("el cuerpo no es el PDF del caso de uso: %q", rec.Body.Bytes())
	}
}

// esPaginaDeErrorPDF verifica que el error se respondió como página HTML (P-7) y no como descarga.
func esPaginaDeErrorPDF(t *testing.T, rec *httptest.ResponseRecorder, code int, mensaje string) {
	t.Helper()
	if rec.Code != code {
		t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, se esperaba la página HTML", ct)
	}
	if cd := rec.Header().Get("Content-Disposition"); cd != "" {
		t.Errorf("una página de error no se descarga; Content-Disposition = %q", cd)
	}
	contiene(t, rec, "No se pudo completar el pedido", mensaje, `<a href="/web/proyectos">`)
	seccionNoContiene(t, rec.Body.String(), "<script", "%PDF")
}

// Traza: HU-27 CA-27.14, CA-27.24 (sección 3, RN-27.12, RN-27.15): cabeceras de descarga y los bytes
// del caso de uso tal cual, con el Content-Length de su tamaño.
func TestReportePDFProyecto_Descarga(t *testing.T) {
	uc := &exportadorReporteProyectoFalso{archivo: app.ArchivoPDF{Nombre: "reporte-proyecto-3-2026-11-07.pdf", Contenido: contenidoPDFFalso}}

	rec := pedirReportePDF(t, context.Background(), uc, &exportadorReporteSprintFalso{}, rutaPDFProyecto)

	esDescargaPDF(t, rec, "reporte-proyecto-3-2026-11-07.pdf")
	if uc.proyectoID != 3 {
		t.Errorf("el caso de uso recibió el proyecto %d, se esperaba 3", uc.proyectoID)
	}
}

// Traza: HU-27 CA-27.15, CA-27.24 (sección 3, RN-27.12, RN-27.15). Los IDs llegan cada uno en su lugar.
func TestReportePDFSprint_Descarga(t *testing.T) {
	proyecto := &exportadorReporteProyectoFalso{}
	uc := &exportadorReporteSprintFalso{archivo: app.ArchivoPDF{Nombre: "reporte-proyecto-3-sprint-12-2026-11-07.pdf", Contenido: contenidoPDFFalso}}

	rec := pedirReportePDF(t, context.Background(), proyecto, uc, rutaPDFSprint)

	esDescargaPDF(t, rec, "reporte-proyecto-3-sprint-12-2026-11-07.pdf")
	if uc.proyectoID != 3 || uc.sprintID != 12 {
		t.Errorf("el caso de uso recibió proyecto %d y sprint %d, se esperaba 3 y 12", uc.proyectoID, uc.sprintID)
	}
	if proyecto.llamado {
		t.Error("la descarga del sprint no tiene que pedir el reporte del proyecto")
	}
}

type claveContextoPDF struct{}

// Traza: HU-27 CA-27.14, CA-27.15 (sección 5): el caso de uso recibe el contexto del pedido, así una
// cancelación del navegador llega a la fuente y al generador.
func TestReportePDF_PasaElContextoDelPedido(t *testing.T) {
	ctx := context.WithValue(context.Background(), claveContextoPDF{}, "pedido")
	archivo := app.ArchivoPDF{Nombre: "r.pdf", Contenido: contenidoPDFFalso}
	proyecto, sprint := &exportadorReporteProyectoFalso{archivo: archivo}, &exportadorReporteSprintFalso{archivo: archivo}

	pedirReportePDF(t, ctx, proyecto, sprint, rutaPDFProyecto)
	pedirReportePDF(t, ctx, proyecto, sprint, rutaPDFSprint)

	if proyecto.ctx == nil || proyecto.ctx.Value(claveContextoPDF{}) != "pedido" {
		t.Error("el reporte del proyecto no recibió el contexto del pedido")
	}
	if sprint.ctx == nil || sprint.ctx.Value(claveContextoPDF{}) != "pedido" {
		t.Error("el reporte del sprint no recibió el contexto del pedido")
	}
}

// Traza: HU-27 CA-27.19 (sección 7): 400 con la página de error y sin llamar al caso de uso. Con los
// dos IDs inválidos, el mensaje es el del proyecto, que se valida primero.
func TestReportePDF_IDsInvalidos_Devuelven400SinLlamar(t *testing.T) {
	const (
		idProyecto = "el id del proyecto debe ser un número entero positivo"
		idSprint   = "el id del sprint debe ser un número entero positivo"
	)
	tests := []struct{ nombre, ruta, wantMsg string }{
		{"proyecto: no numérico", "/web/proyectos/abc/reporte/pdf", idProyecto},
		{"proyecto: cero", "/web/proyectos/0/reporte/pdf", idProyecto},
		{"proyecto: negativo", "/web/proyectos/-3/reporte/pdf", idProyecto},
		{"proyecto: fuera de rango", "/web/proyectos/99999999999999999999/reporte/pdf", idProyecto},
		{"sprint: proyecto no numérico", "/web/proyectos/abc/sprints/12/reporte/pdf", idProyecto},
		{"sprint: proyecto en cero", "/web/proyectos/0/sprints/12/reporte/pdf", idProyecto},
		{"sprint: no numérico", "/web/proyectos/3/sprints/abc/reporte/pdf", idSprint},
		{"sprint: cero", "/web/proyectos/3/sprints/0/reporte/pdf", idSprint},
		{"sprint: negativo", "/web/proyectos/3/sprints/-12/reporte/pdf", idSprint},
		{"sprint: los dos inválidos", "/web/proyectos/-1/sprints/abc/reporte/pdf", idProyecto},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			proyecto, sprint := exportadoresConError(errors.New("no se esperaba llamar al caso de uso"))

			rec := pedirReportePDF(t, context.Background(), proyecto, sprint, tt.ruta)

			esPaginaDeErrorPDF(t, rec, http.StatusBadRequest, tt.wantMsg)
			if proyecto.llamado || sprint.llamado {
				t.Error("no se esperaba llamar al caso de uso")
			}
		})
	}
}

// Traza: HU-27 CA-27.17, CA-27.18 (RN-27.13, sección 7): 404 con el mensaje del dominio, aunque el error
// llegue envuelto, y sin el texto que lo envuelve.
func TestReportePDF_NoEncontrados_Devuelven404(t *testing.T) {
	tests := []struct {
		nombre, ruta string
		err, want    error
	}{
		{"proyecto: proyecto", rutaPDFProyecto, domain.ErrProyectoNoEncontrado, domain.ErrProyectoNoEncontrado},
		{"proyecto: proyecto, envuelto", rutaPDFProyecto, fmt.Errorf("buscar registro 3: %w", domain.ErrProyectoNoEncontrado), domain.ErrProyectoNoEncontrado},
		{"sprint: proyecto", rutaPDFSprint, domain.ErrProyectoNoEncontrado, domain.ErrProyectoNoEncontrado},
		{"sprint: proyecto, envuelto", rutaPDFSprint, fmt.Errorf("buscar registro 3: %w", domain.ErrProyectoNoEncontrado), domain.ErrProyectoNoEncontrado},
		{"sprint: sprint", rutaPDFSprint, domain.ErrSprintNoEncontrado, domain.ErrSprintNoEncontrado},
		{"sprint: sprint, envuelto", rutaPDFSprint, fmt.Errorf("buscar registro 12: %w", domain.ErrSprintNoEncontrado), domain.ErrSprintNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			proyecto, sprint := exportadoresConError(tt.err)

			rec := pedirReportePDF(t, context.Background(), proyecto, sprint, tt.ruta)

			esPaginaDeErrorPDF(t, rec, http.StatusNotFound, tt.want.Error())
			seccionNoContiene(t, rec.Body.String(), "buscar registro")
		})
	}
}

// capturarLog redirige el log estándar durante el test y devuelve lo escrito.
func capturarLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	anterior := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(anterior) })
	return &buf
}

// Traza: HU-27 CA-27.20 (sección 7): cualquier otro error es 500 con la página genérica; el detalle va
// solo al log. Un texto parecido a "no encontrado" que no es el error del dominio sigue siendo 500.
func TestReportePDF_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	errores := []struct{ nombre, detalle string }{
		{"de la fuente", "armar el reporte: database is locked"},
		{"del generador", "generar el PDF: fuente tipográfica dañada"},
		{"parece no encontrado", "sprint no encontrado en la caché interna"},
	}
	for _, r := range rutasReportePDF {
		for _, e := range errores {
			t.Run(r.nombre+": "+e.nombre, func(t *testing.T) {
				logs := capturarLog(t)
				proyecto, sprint := exportadoresConError(fmt.Errorf("reporte: %w", errors.New(e.detalle)))

				rec := pedirReportePDF(t, context.Background(), proyecto, sprint, r.ruta)

				esPaginaDeErrorPDF(t, rec, http.StatusInternalServerError, "Ocurrió un error interno. Probá de nuevo más tarde.")
				seccionNoContiene(t, rec.Body.String(), e.detalle, "reporte:")
				if !strings.Contains(logs.String(), e.detalle) {
					t.Errorf("el detalle no quedó en el log; log: %q", logs)
				}
			})
		}
	}
}

// Traza: HU-27 CA-27.20 (RN-27.11, sección 7): un documento inválido es un error de quien lo armó, no del
// usuario: 500 con la página genérica, sin el texto del error, que queda en el log.
func TestReportePDF_DocumentoInvalido_Devuelve500SinDetalle(t *testing.T) {
	for _, r := range rutasReportePDF {
		t.Run(r.nombre, func(t *testing.T) {
			logs := capturarLog(t)
			proyecto, sprint := exportadoresConError(fmt.Errorf("reporte del proyecto 3: %w", domain.ErrDocumentoInvalido))

			rec := pedirReportePDF(t, context.Background(), proyecto, sprint, r.ruta)

			esPaginaDeErrorPDF(t, rec, http.StatusInternalServerError, "Ocurrió un error interno. Probá de nuevo más tarde.")
			seccionNoContiene(t, rec.Body.String(), domain.ErrDocumentoInvalido.Error(), "reporte del proyecto")
			if !strings.Contains(logs.String(), domain.ErrDocumentoInvalido.Error()) {
				t.Errorf("el detalle no quedó en el log; log: %q", logs)
			}
		})
	}
}

// Traza: HU-27 CA-27.16 (RN-27.12): aunque el caso de uso ya da un nombre seguro, el adaptador no deja
// que un nombre con comillas o saltos de línea corte la cabecera o agregue otras.
func TestReportePDF_NombreConCaracteresPeligrosos_NoInyectaCabeceras(t *testing.T) {
	nombres := []string{
		"a\"; filename*=utf-8''malo.exe",
		"a.pdf\r\nSet-Cookie: sesion=robada",
		"a.pdf\nX-Inyectada: 1",
		"Ñandú \"final\".pdf",
		"",
	}
	for _, nombre := range nombres {
		t.Run(strconv.Quote(nombre), func(t *testing.T) {
			uc := &exportadorReporteProyectoFalso{archivo: app.ArchivoPDF{Nombre: nombre, Contenido: contenidoPDFFalso}}

			rec := pedirReportePDF(t, context.Background(), uc, &exportadorReporteSprintFalso{}, rutaPDFProyecto)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			for _, k := range []string{"Set-Cookie", "X-Inyectada"} {
				if v := rec.Header().Values(k); len(v) > 0 {
					t.Errorf("se inyectó la cabecera %s: %q", k, v)
				}
			}
			cd := rec.Header().Values("Content-Disposition")
			if len(cd) != 1 {
				t.Fatalf("Content-Disposition = %q, se esperaba una sola", cd)
			}
			valor, ok := strings.CutPrefix(cd[0], `attachment; filename="`)
			valor, ok2 := strings.CutSuffix(valor, `"`)
			if !ok || !ok2 || valor == "" {
				t.Fatalf("Content-Disposition = %q, se esperaba attachment con un nombre entre comillas", cd[0])
			}
			for _, c := range valor {
				if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '.') {
					t.Errorf("el nombre %q tiene el carácter %q, fuera de los de RN-27.12", valor, c)
				}
			}
		})
	}
}

// Traza: HU-27 (sección 5): las rutas de la descarga conviven con las de las pantallas en un mismo mux,
// como van a quedar en NuevoRouter (ServeMux entra en pánico si dos patrones chocan).
func TestReportePDF_ConviveConLasPantallas(t *testing.T) {
	mux := http.NewServeMux()
	RegistrarPantallas(mux, casosFalsos())
	uc := &exportadorReporteProyectoFalso{archivo: app.ArchivoPDF{Nombre: "reporte-proyecto-3-2026-11-07.pdf", Contenido: contenidoPDFFalso}}
	RegistrarExportarReportePDF(mux, uc, &exportadorReporteSprintFalso{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, rutaPDFProyecto, nil))

	esDescargaPDF(t, rec, "reporte-proyecto-3-2026-11-07.pdf")
}
