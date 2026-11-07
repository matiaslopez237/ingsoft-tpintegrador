package bdd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu27Errores traduce la descripción del escenario al error definido en la spec (sección 7).
func hu27Errores() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"sprint no encontrado":   domain.ErrSprintNoEncontrado,
	}
}

// hu27UTCMenos3 es el huso con el que se compara la emisión (RN-27.10): 3 horas menos que UTC.
const hu27UTCMenos3 = -3 * 60 * 60

type hu27World struct {
	sistema   *hu27Sistema
	archivo   app.ArchivoPDF
	errExport error // el de la última exportación
	exporto   bool  // para no dar por buena una exportación que el escenario nunca hizo
}

func registrarStepsHU27(ctx *godog.ScenarioContext) {
	w := &hu27World{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = hu27World{sistema: hu27NuevoSistema()}
		return c, nil
	})

	// Given
	ctx.Step(`^el reloj de HU-27 marca el (\d{4}-\d{2}-\d{2}) a las (\d{2}:\d{2}) UTC$`, w.relojMarca)
	ctx.Step(`^el proyecto (\d+) de HU-27 tiene un reporte titulado "([^"]*)"$`, w.proyectoConReporte)
	ctx.Step(`^el proyecto (\d+) de HU-27 tiene el sprint (\d+) con un reporte titulado "([^"]*)"$`, w.sprintConReporte)
	ctx.Step(`^el título del reporte de HU-27 del proyecto (\d+) es "([^"]*)"$`, w.tituloDelReporte)
	ctx.Step(`^el reporte de HU-27 del proyecto (\d+) tiene el subtítulo "([^"]*)"$`, w.subtituloDelReporte)
	ctx.Step(`^el reporte de HU-27 del proyecto (\d+) tiene la sección "([^"]*)" con el párrafo "([^"]*)"$`, w.seccionConParrafo)
	ctx.Step(`^el reporte de HU-27 del proyecto (\d+) tiene la sección "([^"]*)" con esta tabla:$`, w.seccionConTabla)
	ctx.Step(`^el reporte de HU-27 del proyecto (\d+) tiene la sección "([^"]*)" sin contenido$`, w.seccionSinContenido)
	ctx.Step(`^el reporte de HU-27 del proyecto (\d+) es inválido por tener (.+)$`, w.reporteInvalido)
	ctx.Step(`^(la fuente de documentos|el generador de PDF) de HU-27 falla$`, w.componenteFalla)

	// When
	ctx.Step(`^exporto a PDF el reporte de HU-27 del proyecto (\d+)$`, w.exportoProyecto)
	ctx.Step(`^exporto a PDF el reporte de HU-27 del sprint (\d+) del proyecto (\d+)$`, w.exportoSprint)

	// Then
	ctx.Step(`^el documento de HU-27 tiene el título "([^"]*)"$`, w.tieneTitulo)
	ctx.Step(`^el documento de HU-27 tiene el título "([^"]*)" y el subtítulo "([^"]*)"$`, w.tieneTituloYSubtitulo)
	ctx.Step(`^el documento de HU-27 tiene el subtítulo "([^"]*)"$`, w.tieneSubtitulo)
	ctx.Step(`^el documento de HU-27 no tiene subtítulo$`, w.noTieneSubtitulo)
	ctx.Step(`^el documento de HU-27 se emite el (\d{2}/\d{2}/\d{4}) a las (\d{2}:\d{2}) en UTC-3$`, w.seEmite)
	ctx.Step(`^el documento de HU-27 tiene estas secciones, en este orden:$`, w.tieneSecciones)
	ctx.Step(`^el documento de HU-27 no tiene secciones$`, w.noTieneSecciones)
	ctx.Step(`^la sección "([^"]*)" del documento de HU-27 tiene el párrafo "([^"]*)"$`, w.seccionTieneParrafo)
	ctx.Step(`^la sección "([^"]*)" del documento de HU-27 tiene esta tabla:$`, w.seccionTieneTabla)
	ctx.Step(`^la sección "([^"]*)" del documento de HU-27 no tiene contenido$`, w.seccionNoTieneContenido)
	ctx.Step(`^el documento de HU-27 no tiene aviso de reemplazos$`, w.noTieneAviso)
	ctx.Step(`^el documento de HU-27 tiene el aviso "([^"]*)"$`, w.tieneAviso)
	ctx.Step(`^el archivo PDF de HU-27 se llama "([^"]*)"$`, w.archivoSeLlama)
	ctx.Step(`^el archivo PDF de HU-27 tiene los bytes que devolvió el generador$`, w.archivoTieneBytesDelGenerador)
	ctx.Step(`^se rechaza la exportación de HU-27 con el error "([^"]*)"$`, w.seRechazaConError)
	ctx.Step(`^se rechaza la exportación de HU-27 por documento inválido$`, w.seRechazaPorDocumentoInvalido)
	ctx.Step(`^la exportación de HU-27 falla por un error interno$`, w.fallaPorErrorInterno)
	ctx.Step(`^el generador de PDF de HU-27 no se usó$`, w.generadorNoSeUso)
}

// --- Given ---

func (w *hu27World) relojMarca(fecha, hora string) error {
	t, err := time.ParseInLocation("2006-01-02 15:04", fecha+" "+hora, time.UTC)
	if err != nil {
		return fmt.Errorf("fecha del reloj inválida en el escenario: %w", err)
	}
	w.sistema.ahora = t
	return nil
}

func (w *hu27World) proyectoConReporte(proyectoID int64, titulo string) error {
	w.sistema.fuente.proyectos[proyectoID] = &domain.DocumentoPDF{Titulo: titulo}
	return nil
}

func (w *hu27World) sprintConReporte(proyectoID, sprintID int64, titulo string) error {
	if _, ok := w.sistema.fuente.proyectos[proyectoID]; !ok {
		return fmt.Errorf("el escenario no preparó el proyecto %d", proyectoID)
	}
	w.sistema.fuente.sprints[sprintID] = hu27SprintEnMemoria{proyectoID: proyectoID, doc: &domain.DocumentoPDF{Titulo: titulo}}
	return nil
}

// reporte devuelve el documento que el escenario está armando para el proyecto.
func (w *hu27World) reporte(proyectoID int64) (*domain.DocumentoPDF, error) {
	doc, ok := w.sistema.fuente.proyectos[proyectoID]
	if !ok {
		return nil, fmt.Errorf("el escenario no preparó el proyecto %d", proyectoID)
	}
	return doc, nil
}

func (w *hu27World) tituloDelReporte(proyectoID int64, titulo string) error {
	doc, err := w.reporte(proyectoID)
	if err != nil {
		return err
	}
	doc.Titulo, err = hu27Decodificar(titulo)
	return err
}

func (w *hu27World) subtituloDelReporte(proyectoID int64, subtitulo string) error {
	doc, err := w.reporte(proyectoID)
	if err != nil {
		return err
	}
	doc.Subtitulo = subtitulo
	return nil
}

func (w *hu27World) agregarSeccion(proyectoID int64, s domain.SeccionPDF) error {
	doc, err := w.reporte(proyectoID)
	if err != nil {
		return err
	}
	doc.Secciones = append(doc.Secciones, s)
	return nil
}

func (w *hu27World) seccionConParrafo(proyectoID int64, titulo, parrafo string) error {
	return w.agregarSeccion(proyectoID, domain.SeccionPDF{Titulo: titulo, Bloques: []domain.BloquePDF{{Parrafo: parrafo}}})
}

func (w *hu27World) seccionConTabla(proyectoID int64, titulo string, tabla *godog.Table) error {
	t, err := hu27Tabla(tabla)
	if err != nil {
		return err
	}
	return w.agregarSeccion(proyectoID, domain.SeccionPDF{Titulo: titulo, Bloques: []domain.BloquePDF{{Tabla: &t}}})
}

func (w *hu27World) seccionSinContenido(proyectoID int64, titulo string) error {
	return w.agregarSeccion(proyectoID, domain.SeccionPDF{Titulo: titulo})
}

// reporteInvalido arma cada uno de los documentos que RN-27.11 rechaza.
func (w *hu27World) reporteInvalido(proyectoID int64, defecto string) error {
	doc, err := w.reporte(proyectoID)
	if err != nil {
		return err
	}
	parrafo := domain.BloquePDF{Parrafo: "Texto"}
	tabla := func(columnas []string, filas ...[]string) domain.SeccionPDF {
		return domain.SeccionPDF{Titulo: "Tabla", Bloques: []domain.BloquePDF{{Tabla: &domain.TablaPDF{Columnas: columnas, Filas: filas}}}}
	}
	switch defecto {
	case "un título vacío":
		doc.Titulo = ""
	case "un título solo de espacios":
		doc.Titulo = "   "
	case "un título con solo caracteres invisibles":
		doc.Titulo = "​​"
	case "una sección sin título":
		doc.Secciones = append(doc.Secciones, domain.SeccionPDF{Titulo: "", Bloques: []domain.BloquePDF{parrafo}})
	case "un bloque con párrafo y tabla":
		doc.Secciones = append(doc.Secciones, domain.SeccionPDF{Titulo: "Mixto", Bloques: []domain.BloquePDF{
			{Parrafo: "Texto", Tabla: &domain.TablaPDF{Columnas: []string{"A"}}},
		}})
	case "una tabla sin columnas":
		doc.Secciones = append(doc.Secciones, tabla(nil))
	case "una fila con menos celdas que columnas":
		doc.Secciones = append(doc.Secciones, tabla([]string{"A", "B"}, []string{"1"}))
	case "una fila con más celdas que columnas":
		doc.Secciones = append(doc.Secciones, tabla([]string{"A", "B"}, []string{"1", "2", "3"}))
	default:
		return fmt.Errorf("defecto desconocido en el escenario: %q", defecto)
	}
	return nil
}

func (w *hu27World) componenteFalla(componente string) error {
	switch componente {
	case "la fuente de documentos":
		w.sistema.fuente.falla = errHU27FallaInterna
	case "el generador de PDF":
		w.sistema.generador.falla = errHU27FallaInterna
	}
	return nil
}

// --- When ---

func (w *hu27World) exportoProyecto(proyectoID int64) error {
	w.archivo, w.errExport = w.sistema.ExportarProyecto(proyectoID)
	w.exporto = true
	return nil
}

func (w *hu27World) exportoSprint(sprintID, proyectoID int64) error {
	w.archivo, w.errExport = w.sistema.ExportarSprint(proyectoID, sprintID)
	w.exporto = true
	return nil
}

// --- Then ---

// documento es lo que recibió el generador: el documento ya normalizado, con la emisión y el aviso.
func (w *hu27World) documento() (domain.DocumentoPDF, error) {
	if !w.exporto {
		return domain.DocumentoPDF{}, errors.New("el escenario no exportó ningún reporte")
	}
	if w.errExport != nil {
		return domain.DocumentoPDF{}, fmt.Errorf("la exportación falló: %w", w.errExport)
	}
	if w.sistema.generador.llamadas != 1 {
		return domain.DocumentoPDF{}, fmt.Errorf("el generador se llamó %d veces, se esperaba 1", w.sistema.generador.llamadas)
	}
	return w.sistema.generador.recibido, nil
}

func (w *hu27World) tieneTitulo(titulo string) error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if doc.Titulo != titulo {
		return fmt.Errorf("título: se esperaba %q, se obtuvo %q", titulo, doc.Titulo)
	}
	return nil
}

func (w *hu27World) tieneTituloYSubtitulo(titulo, subtitulo string) error {
	if err := w.tieneTitulo(titulo); err != nil {
		return err
	}
	return w.tieneSubtitulo(subtitulo)
}

func (w *hu27World) tieneSubtitulo(subtitulo string) error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if doc.Subtitulo != subtitulo {
		return fmt.Errorf("subtítulo: se esperaba %q, se obtuvo %q", subtitulo, doc.Subtitulo)
	}
	return nil
}

func (w *hu27World) noTieneSubtitulo() error {
	return w.tieneSubtitulo("")
}

func (w *hu27World) seEmite(fecha, hora string) error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if _, offset := doc.Emision.Zone(); offset != hu27UTCMenos3 {
		return fmt.Errorf("la emisión está en un huso de %d segundos respecto de UTC, se esperaba UTC-3", offset)
	}
	if got := doc.Emision.Format("02/01/2006 15:04"); got != fecha+" "+hora {
		return fmt.Errorf("emisión: se esperaba %s %s, se obtuvo %s", fecha, hora, got)
	}
	return nil
}

func (w *hu27World) tieneSecciones(tabla *godog.Table) error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	var esperadas, obtenidas []string
	for _, fila := range tabla.Rows[1:] {
		esperadas = append(esperadas, fila.Cells[0].Value)
	}
	for _, s := range doc.Secciones {
		obtenidas = append(obtenidas, s.Titulo)
	}
	if !slices.Equal(esperadas, obtenidas) {
		return fmt.Errorf("secciones: se esperaban %q, se obtuvieron %q", esperadas, obtenidas)
	}
	return nil
}

func (w *hu27World) noTieneSecciones() error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if len(doc.Secciones) != 0 {
		return fmt.Errorf("se esperaba un documento sin secciones, tiene %d", len(doc.Secciones))
	}
	return nil
}

// seccion busca la sección por su título ya normalizado.
func (w *hu27World) seccion(titulo string) (domain.SeccionPDF, error) {
	doc, err := w.documento()
	if err != nil {
		return domain.SeccionPDF{}, err
	}
	for _, s := range doc.Secciones {
		if s.Titulo == titulo {
			return s, nil
		}
	}
	return domain.SeccionPDF{}, fmt.Errorf("el documento no tiene la sección %q", titulo)
}

func (w *hu27World) seccionTieneParrafo(titulo, parrafo string) error {
	s, err := w.seccion(titulo)
	if err != nil {
		return err
	}
	if len(s.Bloques) != 1 || s.Bloques[0].Tabla != nil || s.Bloques[0].Parrafo != parrafo {
		return fmt.Errorf("la sección %q no tiene solo el párrafo %q: %+v", titulo, parrafo, s.Bloques)
	}
	return nil
}

func (w *hu27World) seccionTieneTabla(titulo string, tabla *godog.Table) error {
	s, err := w.seccion(titulo)
	if err != nil {
		return err
	}
	esperada, err := hu27Tabla(tabla)
	if err != nil {
		return err
	}
	if len(s.Bloques) != 1 || s.Bloques[0].Tabla == nil {
		return fmt.Errorf("la sección %q no tiene solo una tabla: %+v", titulo, s.Bloques)
	}
	got := s.Bloques[0].Tabla
	if !slices.Equal(got.Columnas, esperada.Columnas) {
		return fmt.Errorf("columnas de %q: se esperaban %q, se obtuvieron %q", titulo, esperada.Columnas, got.Columnas)
	}
	if !slices.EqualFunc(got.Filas, esperada.Filas, slices.Equal[[]string]) {
		return fmt.Errorf("filas de %q: se esperaban %q, se obtuvieron %q", titulo, esperada.Filas, got.Filas)
	}
	return nil
}

func (w *hu27World) seccionNoTieneContenido(titulo string) error {
	s, err := w.seccion(titulo)
	if err != nil {
		return err
	}
	if len(s.Bloques) != 0 {
		return fmt.Errorf("la sección %q tiene %d bloques, se esperaba ninguno", titulo, len(s.Bloques))
	}
	return nil
}

func (w *hu27World) noTieneAviso() error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if doc.Aviso != "" {
		return fmt.Errorf("se esperaba un documento sin aviso, tiene %q", doc.Aviso)
	}
	return nil
}

func (w *hu27World) tieneAviso(aviso string) error {
	doc, err := w.documento()
	if err != nil {
		return err
	}
	if doc.Aviso != aviso {
		return fmt.Errorf("aviso: se esperaba %q, se obtuvo %q", aviso, doc.Aviso)
	}
	return nil
}

func (w *hu27World) archivoSeLlama(nombre string) error {
	if _, err := w.documento(); err != nil {
		return err
	}
	if w.archivo.Nombre != nombre {
		return fmt.Errorf("nombre del archivo: se esperaba %q, se obtuvo %q", nombre, w.archivo.Nombre)
	}
	return nil
}

func (w *hu27World) archivoTieneBytesDelGenerador() error {
	if _, err := w.documento(); err != nil {
		return err
	}
	if len(w.archivo.Contenido) == 0 || !bytes.Equal(w.archivo.Contenido, w.sistema.generador.devuelto) {
		return fmt.Errorf("el archivo no tiene los bytes del generador: %q", w.archivo.Contenido)
	}
	return nil
}

// sinArchivo comprueba que una exportación rechazada no devolvió ningún archivo.
func (w *hu27World) sinArchivo() error {
	if w.archivo.Nombre != "" || len(w.archivo.Contenido) != 0 {
		return fmt.Errorf("la exportación falló pero devolvió el archivo %q con %d bytes", w.archivo.Nombre, len(w.archivo.Contenido))
	}
	return nil
}

func (w *hu27World) rechazo() error {
	if !w.exporto {
		return errors.New("el escenario no exportó ningún reporte")
	}
	if w.errExport == nil {
		return errors.New("se esperaba un error, pero la exportación salió bien")
	}
	return w.sinArchivo()
}

func (w *hu27World) seRechazaConError(descripcion string) error {
	esperado, ok := hu27Errores()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if err := w.rechazo(); err != nil {
		return err
	}
	if !errors.Is(w.errExport, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.errExport)
	}
	return nil
}

func (w *hu27World) seRechazaPorDocumentoInvalido() error {
	if err := w.rechazo(); err != nil {
		return err
	}
	if !errors.Is(w.errExport, domain.ErrDocumentoInvalido) {
		return fmt.Errorf("se esperaba ErrDocumentoInvalido, se obtuvo: %w", w.errExport)
	}
	return nil
}

// fallaPorErrorInterno comprueba que la falla del doble llega envuelta y que no se confunde con un error
// de la spec, que el adaptador HTTP mostraría al usuario.
func (w *hu27World) fallaPorErrorInterno() error {
	if err := w.rechazo(); err != nil {
		return err
	}
	if !errors.Is(w.errExport, errHU27FallaInterna) {
		return fmt.Errorf("se esperaba la falla interna del doble, se obtuvo: %w", w.errExport)
	}
	for _, conocido := range []error{domain.ErrProyectoNoEncontrado, domain.ErrSprintNoEncontrado, domain.ErrDocumentoInvalido} {
		if errors.Is(w.errExport, conocido) {
			return fmt.Errorf("una falla interna no puede ser %q", conocido)
		}
	}
	return nil
}

func (w *hu27World) generadorNoSeUso() error {
	if n := w.sistema.generador.llamadas; n != 0 {
		return fmt.Errorf("el generador se llamó %d veces, se esperaba ninguna", n)
	}
	return nil
}

// --- Ayudantes ---

// hu27Tabla convierte una tabla del escenario en una TablaPDF: la primera fila son los encabezados.
func hu27Tabla(tabla *godog.Table) (domain.TablaPDF, error) {
	if len(tabla.Rows) == 0 {
		return domain.TablaPDF{}, errors.New("la tabla del escenario no tiene encabezados")
	}
	var t domain.TablaPDF
	for _, c := range tabla.Rows[0].Cells {
		t.Columnas = append(t.Columnas, c.Value)
	}
	for _, fila := range tabla.Rows[1:] {
		var celdas []string
		for _, c := range fila.Cells {
			celdas = append(celdas, c.Value)
		}
		t.Filas = append(t.Filas, celdas)
	}
	return t, nil
}

// hu27Escape es la notación del .feature para lo que no se puede escribir: {TAB}, {LF}, {U+XXXX} y
// {0xFF} (un byte suelto, que no es UTF-8 válido).
var hu27Escape = regexp.MustCompile(`\{(TAB|LF|U\+[0-9A-F]{4,6}|0x[0-9A-F]{2})\}`)

func hu27Decodificar(s string) (string, error) {
	var errDecodificar error
	r := hu27Escape.ReplaceAllStringFunc(s, func(m string) string {
		codigo := strings.Trim(m, "{}")
		switch {
		case codigo == "TAB":
			return "\t"
		case codigo == "LF":
			return "\n"
		case strings.HasPrefix(codigo, "U+"):
			n, err := strconv.ParseUint(codigo[2:], 16, 32)
			if err != nil {
				errDecodificar = err
			}
			return string(rune(n))
		default: // 0xFF
			n, err := strconv.ParseUint(codigo[2:], 16, 8)
			if err != nil {
				errDecodificar = err
			}
			return string([]byte{byte(n)})
		}
	})
	return r, errDecodificar
}
