package bdd

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
)

// elementoEditable es lo que los escenarios de HU-06 observan de un elemento del backlog,
// con cada campo como texto para compararlo con las tablas del .feature.
type elementoEditable struct {
	ID          int64
	Titulo      string
	Descripcion string
	Prioridad   string
	StoryPoints string // "sin estimar" o el número
	Criterios   string // separados por "; "
	Estado      string
}

func (e elementoEditable) campo(nombre string) (string, error) {
	switch nombre {
	case "titulo":
		return e.Titulo, nil
	case "descripcion":
		return e.Descripcion, nil
	case "prioridad":
		return e.Prioridad, nil
	case "story_points":
		return e.StoryPoints, nil
	case "criterios":
		return e.Criterios, nil
	case "estado":
		return e.Estado, nil
	}
	return "", fmt.Errorf("campo desconocido en el escenario: %q", nombre)
}

// sistemaEdicionBacklog es la puerta de entrada de los escenarios de HU-06 al núcleo y a la
// pantalla. La implementación real está en hu06_sistema_test.go.
type sistemaEdicionBacklog interface {
	CrearProyecto(nombre string) (int64, error)
	// CrearElemento crea un elemento con descripción "Necesidad de prueba" y un criterio.
	CrearElemento(proyectoID int64, titulo, prioridad string, storyPoints int) (int64, error)
	RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64) error
	// RegistrarDefecto registra un defecto en el elemento con el caso de uso de HU-20 (DEF-02).
	RegistrarDefecto(proyectoID, elementoID int64, severidad string) error
	// EditarElemento envía solo los campos indicados. story_points "null" quita la estimación;
	// criterios va separado por ";" y vacío es una lista vacía.
	EditarElemento(proyectoID, elementoID int64, campos map[string]string) error
	EliminarElemento(proyectoID, elementoID int64) error
	BuscarElemento(id int64) (elementoEditable, error)
	ContarElementos(proyectoID int64) (int, error)
	// Pedir hace un pedido a la interfaz web sin seguir redirecciones.
	Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error)
	Cerrar() error
}

// Expresiones para leer el formulario de la pantalla del elemento como lo enviaría el navegador.
var (
	inputDelElemento    = regexp.MustCompile(`<input[^>]*\sname="([^"]+)"[^>]*\svalue="([^"]*)"`)
	selectDelElemento   = regexp.MustCompile(`<select[^>]*\sname="([^"]+)"[\s\S]*?<option value="([^"]*)" selected>`)
	textareaDelElemento = regexp.MustCompile(`<textarea[^>]*\sname="([^"]+)"[^>]*>([^<]*)</textarea>`)
)

// leerFormularioElemento devuelve los <input>, los <select> elegidos y los <textarea> de la página.
func leerFormularioElemento(pagina string) url.Values {
	form := url.Values{}
	for _, re := range []*regexp.Regexp{inputDelElemento, selectDelElemento, textareaDelElemento} {
		for _, m := range re.FindAllStringSubmatch(pagina, -1) {
			form.Set(m[1], html.UnescapeString(m[2]))
		}
	}
	return form
}

type edicionBacklogWorld struct {
	sistema   sistemaEdicionBacklog
	proyectos map[string]int64
	elementos map[string]int64            // título → ID
	duenos    map[int64]int64             // ID del elemento → ID de su proyecto
	iniciales map[string]elementoEditable // título → cómo estaba antes del When
	ultimoErr error
	pantalla  respuestaWeb
}

func registrarStepsHU06(ctx *godog.ScenarioContext) {
	w := &edicionBacklogWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaEdicionBacklog()
		if err != nil {
			return c, err
		}
		*w = edicionBacklogWorld{sistema: sistema, proyectos: map[string]int64{},
			elementos: map[string]int64{}, duenos: map[int64]int64{}, iniciales: map[string]elementoEditable{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^el proyecto "([^"]*)" tiene en su backlog el elemento "([^"]*)" con prioridad "([^"]*)" y (\d+) story points$`, w.tieneEnSuBacklog)
	ctx.Step(`^el elemento "([^"]*)" está en estado "([^"]*)"$`, w.elementoEstaEnEstado)
	ctx.Step(`^"([^"]*)" registró (\d+) horas de esfuerzo en el elemento "([^"]*)"$`, w.registroEsfuerzo)
	ctx.Step(`^se registró un defecto de severidad "([^"]*)" en el elemento "([^"]*)"$`, w.registroDefecto)

	ctx.Step(`^edito el elemento "([^"]*)" cambiando "([^"]*)" a "([^"]*)"$`, w.editoCambiando)
	ctx.Step(`^edito el elemento "([^"]*)" con:$`, w.editoCon)
	ctx.Step(`^edito el elemento "([^"]*)" sin cambiar nada$`, w.editoSinCambios)
	ctx.Step(`^edito desde el proyecto "([^"]*)" el elemento "([^"]*)" cambiando "([^"]*)" a "([^"]*)"$`, w.editoDesdeProyecto)
	ctx.Step(`^elimino el elemento "([^"]*)"$`, w.elimino)
	ctx.Step(`^elimino desde el proyecto "([^"]*)" el elemento "([^"]*)"$`, w.eliminoDesdeProyecto)
	ctx.Step(`^desde la pantalla del elemento "([^"]*)" cambio "([^"]*)" a "([^"]*)"$`, w.desdePantallaCambio)
	ctx.Step(`^desde la pantalla elimino el elemento "([^"]*)"$`, w.desdePantallaElimino)
	ctx.Step(`^desde la pantalla del elemento "([^"]*)" guardo sin cambiar nada$`, w.desdePantallaGuardoSinCambios)

	ctx.Step(`^el elemento "([^"]*)" queda con:$`, w.quedaCon)
	ctx.Step(`^el elemento "([^"]*)" conserva el identificador de "([^"]*)"$`, w.conservaIdentificador)
	ctx.Step(`^se rechaza la edición del elemento con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^se rechaza la eliminación del elemento con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^el elemento "([^"]*)" no cambió$`, w.noCambio)
	ctx.Step(`^el elemento "([^"]*)" sigue existiendo$`, w.sigueExistiendo)
	ctx.Step(`^el elemento "([^"]*)" ya no existe$`, w.yaNoExiste)
	ctx.Step(`^el backlog de "([^"]*)" tiene (\d+) elementos$`, w.backlogTiene)
	ctx.Step(`^la pantalla del elemento responde (\d+) y muestra "([^"]*)"$`, w.pantallaMuestra)
	ctx.Step(`^la pantalla del elemento vuelve al elemento sin mostrar un error$`, w.pantallaVuelveSinError)
}

// --- Given ---

func (w *edicionBacklogWorld) proyecto(nombre string) (int64, error) {
	if id, ok := w.proyectos[nombre]; ok {
		return id, nil
	}
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return 0, fmt.Errorf("preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return id, nil
}

func (w *edicionBacklogWorld) tieneEnSuBacklog(proyecto, titulo, prioridad string, puntos int) error {
	pid, err := w.proyecto(proyecto)
	if err != nil {
		return err
	}
	id, err := w.sistema.CrearElemento(pid, titulo, prioridad, puntos)
	if err != nil {
		return fmt.Errorf("preparar el elemento %q: %w", titulo, err)
	}
	w.elementos[titulo] = id
	w.duenos[id] = pid
	return w.recordar(titulo)
}

// recordar guarda cómo está el elemento, para que "no cambió" lo compare después.
func (w *edicionBacklogWorld) recordar(titulo string) error {
	e, err := w.sistema.BuscarElemento(w.elementos[titulo])
	if err != nil {
		return err
	}
	w.iniciales[titulo] = e
	return nil
}

func (w *edicionBacklogWorld) elementoEstaEnEstado(titulo, estado string) error {
	if err := w.editoCambiando(titulo, "estado", estado); err != nil {
		return err
	}
	if w.ultimoErr != nil {
		return fmt.Errorf("preparar el estado de %q: %w", titulo, w.ultimoErr)
	}
	return w.recordar(titulo)
}

func (w *edicionBacklogWorld) registroEsfuerzo(integrante string, horas int, titulo string) error {
	pid, eid, err := w.ids("", titulo)
	if err != nil {
		return err
	}
	return w.sistema.RegistrarEsfuerzo(pid, eid, integrante, float64(horas))
}

func (w *edicionBacklogWorld) registroDefecto(severidad, titulo string) error {
	pid, eid, err := w.ids("", titulo)
	if err != nil {
		return err
	}
	return w.sistema.RegistrarDefecto(pid, eid, severidad)
}

// ids devuelve el proyecto indicado (o el del elemento si proyecto es "") y el ID del elemento.
func (w *edicionBacklogWorld) ids(proyecto, titulo string) (int64, int64, error) {
	eid, ok := w.elementos[titulo]
	if !ok {
		return 0, 0, fmt.Errorf("el escenario no conoce el elemento %q", titulo)
	}
	if proyecto == "" {
		return w.duenos[eid], eid, nil
	}
	pid, ok := w.proyectos[proyecto]
	if !ok {
		return 0, 0, fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	return pid, eid, nil
}

// --- When ---

func (w *edicionBacklogWorld) editar(proyecto, titulo string, campos map[string]string) error {
	pid, eid, err := w.ids(proyecto, titulo)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.EditarElemento(pid, eid, campos)
	if w.ultimoErr == nil {
		if nuevo, ok := campos["titulo"]; ok && nuevo != titulo {
			w.elementos[nuevo] = eid
		}
	}
	return nil
}

func (w *edicionBacklogWorld) editoCambiando(titulo, campo, valor string) error {
	return w.editar("", titulo, map[string]string{campo: valor})
}

func (w *edicionBacklogWorld) editoCon(titulo string, tabla *godog.Table) error {
	campos := map[string]string{}
	for _, fila := range tabla.Rows[1:] {
		campos[fila.Cells[0].Value] = fila.Cells[1].Value
	}
	return w.editar("", titulo, campos)
}

func (w *edicionBacklogWorld) editoSinCambios(titulo string) error {
	return w.editar("", titulo, map[string]string{})
}

func (w *edicionBacklogWorld) editoDesdeProyecto(proyecto, titulo, campo, valor string) error {
	return w.editar(proyecto, titulo, map[string]string{campo: valor})
}

func (w *edicionBacklogWorld) eliminar(proyecto, titulo string) error {
	pid, eid, err := w.ids(proyecto, titulo)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.EliminarElemento(pid, eid)
	return nil
}

func (w *edicionBacklogWorld) elimino(titulo string) error { return w.eliminar("", titulo) }

func (w *edicionBacklogWorld) eliminoDesdeProyecto(proyecto, titulo string) error {
	return w.eliminar(proyecto, titulo)
}

func (w *edicionBacklogWorld) rutaPantalla(titulo string) (string, error) {
	pid, eid, err := w.ids("", titulo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid), nil
}

// desdePantallaCambio abre la pantalla del elemento, cambia un campo del formulario de edición
// (los demás quedan como estaban) y lo envía, como lo haría el navegador.
func (w *edicionBacklogWorld) desdePantallaCambio(titulo, campo, valor string) error {
	ruta, err := w.rutaPantalla(titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if r.Estado != 200 {
		return fmt.Errorf("abrir la pantalla del elemento: estado %d", r.Estado)
	}
	form := leerFormularioElemento(r.Cuerpo)
	form.Set(campo, valor)
	w.pantalla, err = w.sistema.Pedir("POST", ruta+"/editar", form)
	return err
}

// desdePantallaGuardoSinCambios envía el formulario de la pantalla tal como se abrió.
func (w *edicionBacklogWorld) desdePantallaGuardoSinCambios(titulo string) error {
	ruta, err := w.rutaPantalla(titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	w.pantalla, err = w.sistema.Pedir("POST", ruta+"/editar", leerFormularioElemento(r.Cuerpo))
	return err
}

var rutaElemento = regexp.MustCompile(`^/web/proyectos/\d+/backlog/\d+$`)

// pantallaVuelveSinError verifica que el envío redirigió a la pantalla del elemento (sin página de error).
func (w *edicionBacklogWorld) pantallaVuelveSinError() error {
	if w.pantalla.Estado != 303 || !rutaElemento.MatchString(w.pantalla.Ubicacion) {
		return fmt.Errorf("se esperaba volver al elemento (303), estado %d, Location %q; cuerpo: %.300s",
			w.pantalla.Estado, w.pantalla.Ubicacion, w.pantalla.Cuerpo)
	}
	return nil
}

func (w *edicionBacklogWorld) desdePantallaElimino(titulo string) error {
	ruta, err := w.rutaPantalla(titulo)
	if err != nil {
		return err
	}
	w.pantalla, err = w.sistema.Pedir("POST", ruta+"/eliminar", url.Values{})
	if err != nil {
		return err
	}
	if w.pantalla.Estado != 303 {
		return fmt.Errorf("eliminar desde la pantalla: estado %d; cuerpo: %.300s", w.pantalla.Estado, w.pantalla.Cuerpo)
	}
	return nil
}

// --- Then ---

func (w *edicionBacklogWorld) quedaCon(titulo string, tabla *godog.Table) error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba editar el elemento, pero hubo error: %w", w.ultimoErr)
	}
	if w.pantalla.Estado != 0 && w.pantalla.Estado != 303 {
		return fmt.Errorf("la pantalla respondió %d; cuerpo: %.300s", w.pantalla.Estado, w.pantalla.Cuerpo)
	}
	e, err := w.sistema.BuscarElemento(w.elementos[titulo])
	if err != nil {
		return err
	}
	if e.Titulo != titulo {
		return fmt.Errorf("título: se esperaba %q, se obtuvo %q", titulo, e.Titulo)
	}
	for _, fila := range tabla.Rows[1:] {
		nombre, want := fila.Cells[0].Value, fila.Cells[1].Value
		got, err := e.campo(nombre)
		if err != nil {
			return err
		}
		if got != want {
			return fmt.Errorf("%s: se esperaba %q, se obtuvo %q", nombre, want, got)
		}
	}
	return nil
}

func (w *edicionBacklogWorld) conservaIdentificador(nuevo, viejo string) error {
	if w.elementos[nuevo] != w.iniciales[viejo].ID {
		return fmt.Errorf("el ID cambió: %d → %d", w.iniciales[viejo].ID, w.elementos[nuevo])
	}
	return nil
}

func (w *edicionBacklogWorld) seRechaza(descripcion string) error {
	esperado, ok := erroresHU06()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la operación se hizo", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *edicionBacklogWorld) noCambio(titulo string) error {
	antes := w.iniciales[titulo]
	ahora, err := w.sistema.BuscarElemento(antes.ID)
	if err != nil {
		return err
	}
	if ahora != antes {
		return fmt.Errorf("el elemento cambió: antes %+v, ahora %+v", antes, ahora)
	}
	return nil
}

func (w *edicionBacklogWorld) sigueExistiendo(titulo string) error {
	_, err := w.sistema.BuscarElemento(w.elementos[titulo])
	return err
}

func (w *edicionBacklogWorld) yaNoExiste(titulo string) error {
	_, err := w.sistema.BuscarElemento(w.elementos[titulo])
	if !errors.Is(err, erroresHU06()["elemento no encontrado"]) {
		return fmt.Errorf("se esperaba que el elemento no exista, se obtuvo: %v", err)
	}
	return nil
}

func (w *edicionBacklogWorld) backlogTiene(proyecto string, esperado int) error {
	n, err := w.sistema.ContarElementos(w.proyectos[proyecto])
	if err != nil {
		return err
	}
	if n != esperado {
		return fmt.Errorf("elementos de %q: se esperaban %d, hay %d", proyecto, esperado, n)
	}
	return nil
}

func (w *edicionBacklogWorld) pantallaMuestra(estado int, texto string) error {
	if w.pantalla.Estado != estado {
		return fmt.Errorf("estado: se esperaba %d, se obtuvo %d; cuerpo: %.300s", estado, w.pantalla.Estado, w.pantalla.Cuerpo)
	}
	if !strings.Contains(w.pantalla.Cuerpo, texto) {
		return fmt.Errorf("la pantalla no muestra %q", texto)
	}
	return nil
}
