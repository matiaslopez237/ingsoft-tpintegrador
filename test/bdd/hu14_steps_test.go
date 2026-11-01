package bdd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/cucumber/godog"
)

// votanteVisto y rondaVista son lo que los escenarios de HU-14 observan del resumen de una
// ronda: nunca el valor de un voto (RN-14.3).
type (
	votanteVisto struct {
		Nombre string
		Voto   bool
	}
	rondaVista struct {
		ID, ProyectoID int64
		Numero         int
		Estado         string
		Votantes       []votanteVisto
	}
)

// sistemaPlanningPoker es la puerta de entrada de los escenarios de HU-14 al núcleo y a la
// interfaz web. La implementación real está en hu14_sistema_test.go.
type sistemaPlanningPoker interface {
	CrearProyecto(nombre string) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre string) error
	// CrearElemento crea un elemento del backlog y lo deja en el estado indicado.
	CrearElemento(proyectoID int64, titulo, estado string) (int64, error)
	IniciarRonda(proyectoID, elementoID int64) (rondaVista, error)
	Votar(proyectoID, rondaID int64, integrante string, valor int) (rondaVista, error)
	ConsultarRonda(proyectoID, rondaID int64) (rondaVista, error)
	// ResumenJSON devuelve el cuerpo de GET /proyectos/{id}/rondas/{rondaID}.
	ResumenJSON(proyectoID, rondaID int64) (string, error)
	// Pedir hace un pedido a la interfaz web sin seguir redirecciones.
	Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error)
	Cerrar() error
}

type planningPokerWorld struct {
	sistema   sistemaPlanningPoker
	proyectos map[string]int64
	elementos map[string]int64 // título → ID
	duenos    map[int64]int64  // ID del elemento → ID de su proyecto
	ronda     rondaVista
	ultimoErr error
	pantalla  respuestaWeb
}

func registrarStepsHU14(ctx *godog.ScenarioContext) {
	w := &planningPokerWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaPlanningPoker()
		if err != nil {
			return c, err
		}
		*w = planningPokerWorld{sistema: sistema, proyectos: map[string]int64{},
			elementos: map[string]int64{}, duenos: map[int64]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^el proyecto de Planning Poker "([^"]*)" con los integrantes "([^"]*)"$`, w.proyectoConIntegrantes)
	ctx.Step(`^el proyecto de Planning Poker "([^"]*)" sin integrantes$`, w.proyectoSinIntegrantes)
	ctx.Step(`^en "([^"]*)" hay un elemento a estimar "([^"]*)" en estado "([^"]*)"$`, w.elementoAEstimar)
	ctx.Step(`^hay una ronda abierta sobre "([^"]*)"$`, w.hayRondaAbierta)
	ctx.Step(`^"([^"]*)" ya votó (-?\d+) en la ronda$`, w.yaVoto)
	ctx.Step(`^"([^"]*)" se sumó a "([^"]*)" después de abrir la ronda$`, w.seSumoDespues)

	ctx.Step(`^inicio una ronda de estimación sobre "([^"]*)"$`, w.inicioRonda)
	ctx.Step(`^inicio desde el proyecto "([^"]*)" una ronda de estimación sobre "([^"]*)"$`, w.inicioRondaDesde)
	ctx.Step(`^"([^"]*)" vota (-?\d+) en la ronda$`, w.vota)
	ctx.Step(`^"([^"]*)" vota (-?\d+) en la ronda desde el proyecto "([^"]*)"$`, w.votaDesde)
	ctx.Step(`^desde la pantalla del elemento "([^"]*)" inicio una ronda de estimación$`, w.desdePantallaInicio)
	ctx.Step(`^desde la pantalla de la ronda "([^"]*)" vota (-?\d+)$`, w.desdePantallaVota)

	ctx.Step(`^la ronda queda "([^"]*)" con el número (\d+)$`, w.rondaQueda)
	ctx.Step(`^los votantes de la ronda son "([^"]*)"$`, w.votantesSon)
	ctx.Step(`^en la ronda ya votó "([^"]*)"$`, w.yaVotaron)
	ctx.Step(`^en la ronda falta votar "([^"]*)"$`, w.faltanVotar)
	ctx.Step(`^el resumen de la ronda no muestra ningún valor$`, w.resumenSinValores)
	ctx.Step(`^se rechaza el voto con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^se rechaza la ronda con el error "([^"]*)"$`, w.seRechaza)
	ctx.Step(`^llego a la pantalla de la ronda número (\d+) de "([^"]*)"$`, w.llegoAPantallaRonda)
	ctx.Step(`^la pantalla de la ronda muestra que falta votar "([^"]*)"$`, w.pantallaFaltan)
	ctx.Step(`^la pantalla de la ronda muestra que votaron "([^"]*)"$`, w.pantallaVotaron)
	ctx.Step(`^la pantalla de la ronda dice que está cerrada$`, w.pantallaCerrada)
	ctx.Step(`^la pantalla de la ronda no muestra los valores "([^"]*)"$`, w.pantallaSinValores)
}

// nombres separa una lista del escenario como "Ana, Beto".
func nombres(lista string) []string {
	var r []string
	for _, n := range strings.Split(lista, ",") {
		if n = strings.TrimSpace(n); n != "" {
			r = append(r, n)
		}
	}
	return r
}

// --- Given ---

func (w *planningPokerWorld) proyectoSinIntegrantes(nombre string) error {
	id, err := w.sistema.CrearProyecto(nombre)
	if err != nil {
		return fmt.Errorf("preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *planningPokerWorld) proyectoConIntegrantes(nombre, integrantes string) error {
	if err := w.proyectoSinIntegrantes(nombre); err != nil {
		return err
	}
	for _, n := range nombres(integrantes) {
		if err := w.sistema.AgregarIntegrante(w.proyectos[nombre], n); err != nil {
			return fmt.Errorf("agregar a %q: %w", n, err)
		}
	}
	return nil
}

func (w *planningPokerWorld) elementoAEstimar(proyecto, titulo, estado string) error {
	pid, ok := w.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	id, err := w.sistema.CrearElemento(pid, titulo, estado)
	if err != nil {
		return fmt.Errorf("preparar el elemento %q: %w", titulo, err)
	}
	w.elementos[titulo] = id
	w.duenos[id] = pid
	return nil
}

func (w *planningPokerWorld) elemento(titulo string) (pid, eid int64, err error) {
	eid, ok := w.elementos[titulo]
	if !ok {
		return 0, 0, fmt.Errorf("el escenario no conoce el elemento %q", titulo)
	}
	return w.duenos[eid], eid, nil
}

func (w *planningPokerWorld) hayRondaAbierta(titulo string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if w.ronda, err = w.sistema.IniciarRonda(pid, eid); err != nil {
		return fmt.Errorf("preparar la ronda sobre %q: %w", titulo, err)
	}
	return nil
}

func (w *planningPokerWorld) yaVoto(integrante string, valor int) error {
	r, err := w.sistema.Votar(w.ronda.ProyectoID, w.ronda.ID, integrante, valor)
	if err != nil {
		return fmt.Errorf("preparar el voto de %q: %w", integrante, err)
	}
	w.ronda = r
	return nil
}

func (w *planningPokerWorld) seSumoDespues(integrante, proyecto string) error {
	return w.sistema.AgregarIntegrante(w.proyectos[proyecto], integrante)
}

// --- When ---

func (w *planningPokerWorld) iniciar(pid int64, titulo string) error {
	_, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	r, err := w.sistema.IniciarRonda(pid, eid)
	if errors.Is(err, godog.ErrPending) {
		return err
	}
	w.ultimoErr = err
	if err == nil {
		w.ronda = r
	}
	return nil
}

func (w *planningPokerWorld) inicioRonda(titulo string) error {
	pid, _, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	return w.iniciar(pid, titulo)
}

func (w *planningPokerWorld) inicioRondaDesde(proyecto, titulo string) error {
	return w.iniciar(w.proyectos[proyecto], titulo)
}

func (w *planningPokerWorld) votar(pid int64, integrante string, valor int) error {
	r, err := w.sistema.Votar(pid, w.ronda.ID, integrante, valor)
	if errors.Is(err, godog.ErrPending) {
		return err
	}
	w.ultimoErr = err
	if err == nil {
		w.ronda = r
	}
	return nil
}

func (w *planningPokerWorld) vota(integrante string, valor int) error {
	return w.votar(w.ronda.ProyectoID, integrante, valor)
}

func (w *planningPokerWorld) votaDesde(integrante string, valor int, proyecto string) error {
	return w.votar(w.proyectos[proyecto], integrante, valor)
}

// desdePantallaInicio abre la pantalla del elemento y envía el formulario "Iniciar ronda".
func (w *planningPokerWorld) desdePantallaInicio(titulo string) error {
	pid, eid, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	ruta := fmt.Sprintf("/web/proyectos/%d/backlog/%d", pid, eid)
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if !strings.Contains(r.Cuerpo, ruta+"/rondas") {
		return fmt.Errorf("la pantalla del elemento no ofrece iniciar una ronda (estado %d)", r.Estado)
	}
	w.pantalla, err = w.sistema.Pedir("POST", ruta+"/rondas", url.Values{})
	return err
}

// desdePantallaVota envía el formulario de la pantalla de la ronda y sigue la redirección.
func (w *planningPokerWorld) desdePantallaVota(integrante string, valor int) error {
	ruta := fmt.Sprintf("/web/proyectos/%d/rondas/%d", w.ronda.ProyectoID, w.ronda.ID)
	r, err := w.sistema.Pedir("POST", ruta+"/votar", url.Values{"integrante": {integrante}, "valor": {fmt.Sprint(valor)}})
	if err != nil {
		return err
	}
	if r.Estado != 303 {
		return fmt.Errorf("votar desde la pantalla: estado %d; cuerpo: %.300s", r.Estado, r.Cuerpo)
	}
	w.pantalla, err = w.sistema.Pedir("GET", r.Ubicacion, nil)
	return err
}

// --- Then ---

func (w *planningPokerWorld) sinError() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("no se esperaba un error: %w", w.ultimoErr)
	}
	return nil
}

func (w *planningPokerWorld) rondaQueda(estado string, numero int) error {
	if err := w.sinError(); err != nil {
		return err
	}
	if w.ronda.Estado != estado || w.ronda.Numero != numero {
		return fmt.Errorf("la ronda está %q con el número %d, se esperaba %q con el %d", w.ronda.Estado, w.ronda.Numero, estado, numero)
	}
	return nil
}

func (w *planningPokerWorld) votantesSon(lista string) error {
	r, err := w.sistema.ConsultarRonda(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return err
	}
	var got []string
	for _, v := range r.Votantes {
		got = append(got, v.Nombre)
	}
	if want := nombres(lista); !slices.Equal(got, want) {
		return fmt.Errorf("votantes = %v, se esperaba %v", got, want)
	}
	return nil
}

// votaron verifica, en la ronda guardada, que cada nombre de la lista votó (o no, según voto).
func (w *planningPokerWorld) votaron(lista string, voto bool) error {
	r, err := w.sistema.ConsultarRonda(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return err
	}
	for _, n := range nombres(lista) {
		i := slices.IndexFunc(r.Votantes, func(v votanteVisto) bool { return v.Nombre == n })
		if i < 0 {
			return fmt.Errorf("%q no es votante de la ronda: %+v", n, r.Votantes)
		}
		if r.Votantes[i].Voto != voto {
			return fmt.Errorf("%q: votó = %v, se esperaba %v", n, r.Votantes[i].Voto, voto)
		}
	}
	return nil
}

func (w *planningPokerWorld) yaVotaron(lista string) error { return w.votaron(lista, true) }
func (w *planningPokerWorld) faltanVotar(lista string) error {
	return w.votaron(lista, false)
}

// resumenSinValores pide el resumen a la API (CA-14.8) y verifica que sus votantes solo traigan
// el nombre y si votaron (RN-14.3).
func (w *planningPokerWorld) resumenSinValores() error {
	cuerpo, err := w.sistema.ResumenJSON(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return err
	}
	var resumen struct {
		Votantes []map[string]any `json:"votantes"`
	}
	if err := json.Unmarshal([]byte(cuerpo), &resumen); err != nil {
		return fmt.Errorf("leer el resumen %q: %w", cuerpo, err)
	}
	if len(resumen.Votantes) == 0 {
		return fmt.Errorf("el resumen no trae votantes: %s", cuerpo)
	}
	for _, v := range resumen.Votantes {
		for clave := range v {
			if clave != "nombre" && clave != "voto" {
				return fmt.Errorf("el resumen muestra %q de un votante: %s", clave, cuerpo)
			}
		}
	}
	return nil
}

func (w *planningPokerWorld) seRechaza(descripcion string) error {
	esperado, ok := erroresHU14()[descripcion]
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

var rutaPantallaRonda = regexp.MustCompile(`^/web/proyectos/(\d+)/rondas/(\d+)$`)

func (w *planningPokerWorld) llegoAPantallaRonda(numero int, titulo string) error {
	if w.pantalla.Estado != 303 || !rutaPantallaRonda.MatchString(w.pantalla.Ubicacion) {
		return fmt.Errorf("se esperaba una redirección a la pantalla de la ronda, se obtuvo %d %q; cuerpo: %.300s",
			w.pantalla.Estado, w.pantalla.Ubicacion, w.pantalla.Cuerpo)
	}
	r, err := w.sistema.Pedir("GET", w.pantalla.Ubicacion, nil)
	if err != nil {
		return err
	}
	w.pantalla = r
	if !strings.Contains(r.Cuerpo, fmt.Sprintf("Ronda %d", numero)) || !strings.Contains(r.Cuerpo, titulo) {
		return fmt.Errorf("la pantalla no muestra la ronda %d de %q: %.400s", numero, titulo, r.Cuerpo)
	}
	return nil
}

// filaVotante es la fila de un votante en la pantalla de la ronda, con su nombre y su estado.
func filaVotante(nombre, estado string) *regexp.Regexp {
	return regexp.MustCompile(`<td>` + regexp.QuoteMeta(nombre) + `</td>\s*<td>` + regexp.QuoteMeta(estado) + `</td>`)
}

func (w *planningPokerWorld) pantallaConVotantes(lista, estado string) error {
	for _, n := range nombres(lista) {
		if !filaVotante(n, estado).MatchString(w.pantalla.Cuerpo) {
			return fmt.Errorf("la pantalla no muestra a %q con %q: %.600s", n, estado, w.pantalla.Cuerpo)
		}
	}
	return nil
}

func (w *planningPokerWorld) pantallaFaltan(lista string) error {
	return w.pantallaConVotantes(lista, "Falta votar")
}

func (w *planningPokerWorld) pantallaVotaron(lista string) error {
	return w.pantallaConVotantes(lista, "Votó")
}

func (w *planningPokerWorld) pantallaCerrada() error {
	if !strings.Contains(w.pantalla.Cuerpo, "Votaron todos") {
		return fmt.Errorf("la pantalla no dice que la ronda está cerrada: %.600s", w.pantalla.Cuerpo)
	}
	if strings.Contains(w.pantalla.Cuerpo, `name="valor"`) {
		return errors.New("con la ronda cerrada, la pantalla sigue ofreciendo votar")
	}
	return nil
}

// pantallaSinValores verifica que ningún valor votado aparezca como texto de una celda o de un
// elemento de la pantalla (los números de la ronda o del proyecto pueden coincidir, por eso se
// buscan entre etiquetas).
func (w *planningPokerWorld) pantallaSinValores(lista string) error {
	for _, v := range nombres(lista) {
		if strings.Contains(w.pantalla.Cuerpo, ">"+v+"<") {
			return fmt.Errorf("la pantalla muestra el valor %s", v)
		}
	}
	return nil
}
