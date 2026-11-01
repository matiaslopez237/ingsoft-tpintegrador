package bdd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// votoVisto y revelacionVista son lo que los escenarios de HU-15 observan de una ronda cerrada.
type (
	votoVisto struct {
		Nombre string
		Valor  int
	}
	revelacionVista struct {
		Votos          []votoVisto
		Minimo, Maximo int
		Diferencia     bool
	}
)

// registrarStepsHU15 agrega los pasos de HU-15 sobre el mismo mundo de HU-14: la ronda se arma y
// se vota con los pasos de HU-14 (registrarStepsHU14 llama a esta función).
func registrarStepsHU15(ctx *godog.ScenarioContext, w *planningPokerWorld) {
	ctx.Step(`^la ronda revela los votos "([^"]*)"$`, w.revelaLosVotos)
	ctx.Step(`^el voto mínimo es (\d+) y el máximo es (\d+)$`, w.minimoYMaximo)
	ctx.Step(`^(hay|no hay) una diferencia significativa$`, w.diferencia)
	ctx.Step(`^la ronda no revela ningún voto$`, w.noRevela)
	ctx.Step(`^la API revela los votos "([^"]*)" con mínimo (\d+), máximo (\d+) y diferencia significativa$`, w.apiRevela)
	ctx.Step(`^la pantalla de la ronda muestra los votos "([^"]*)"$`, w.pantallaMuestraVotos)
	ctx.Step(`^la pantalla de la ronda avisa que hay diferencias significativas entre (\d+) y (\d+)$`, w.pantallaAvisaDiferencia)
}

// votosDe lee una lista del escenario como "Ana: 3, Beto: 8".
func votosDe(lista string) ([]votoVisto, error) {
	var votos []votoVisto
	for _, par := range nombres(lista) {
		nombre, valor, ok := strings.Cut(par, ":")
		v, err := strconv.Atoi(strings.TrimSpace(valor))
		if !ok || err != nil {
			return nil, fmt.Errorf("voto mal escrito en el escenario: %q", par)
		}
		votos = append(votos, votoVisto{Nombre: strings.TrimSpace(nombre), Valor: v})
	}
	return votos, nil
}

// revelacion consulta la ronda guardada y devuelve su revelación, o un error si no la tiene.
func (w *planningPokerWorld) revelacion() (*revelacionVista, error) {
	if err := w.sinError(); err != nil {
		return nil, err
	}
	r, err := w.sistema.ConsultarRonda(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return nil, err
	}
	if r.Revelacion == nil {
		return nil, fmt.Errorf("la ronda %q no revela los votos", r.Estado)
	}
	return r.Revelacion, nil
}

func (w *planningPokerWorld) revelaLosVotos(lista string) error {
	want, err := votosDe(lista)
	if err != nil {
		return err
	}
	rev, err := w.revelacion()
	if err != nil {
		return err
	}
	if !slices.Equal(rev.Votos, want) {
		return fmt.Errorf("votos revelados = %+v, se esperaba %+v", rev.Votos, want)
	}
	return nil
}

func (w *planningPokerWorld) minimoYMaximo(minimo, maximo int) error {
	rev, err := w.revelacion()
	if err != nil {
		return err
	}
	if rev.Minimo != minimo || rev.Maximo != maximo {
		return fmt.Errorf("mínimo %d y máximo %d, se esperaba %d y %d", rev.Minimo, rev.Maximo, minimo, maximo)
	}
	return nil
}

func (w *planningPokerWorld) diferencia(hay string) error {
	rev, err := w.revelacion()
	if err != nil {
		return err
	}
	if want := hay == "hay"; rev.Diferencia != want {
		return fmt.Errorf("diferencia significativa = %v, se esperaba %v (votos %+v)", rev.Diferencia, want, rev.Votos)
	}
	return nil
}

func (w *planningPokerWorld) noRevela() error {
	r, err := w.sistema.ConsultarRonda(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return err
	}
	if r.Revelacion != nil {
		return fmt.Errorf("con la ronda %q se revelan los votos: %+v", r.Estado, r.Revelacion)
	}
	return nil
}

// apiRevela pide la ronda al endpoint real (CA-15.5).
func (w *planningPokerWorld) apiRevela(lista string, minimo, maximo int) error {
	want, err := votosDe(lista)
	if err != nil {
		return err
	}
	cuerpo, err := w.sistema.ResumenJSON(w.ronda.ProyectoID, w.ronda.ID)
	if err != nil {
		return err
	}
	var res struct {
		Revelacion *struct {
			Votos []struct {
				Nombre string `json:"nombre"`
				Valor  int    `json:"valor"`
			} `json:"votos"`
			Minimo                  int  `json:"minimo"`
			Maximo                  int  `json:"maximo"`
			DiferenciaSignificativa bool `json:"diferencia_significativa"`
		} `json:"revelacion"`
	}
	if err := json.Unmarshal([]byte(cuerpo), &res); err != nil {
		return fmt.Errorf("leer la ronda %q: %w", cuerpo, err)
	}
	if res.Revelacion == nil {
		return fmt.Errorf("la API no revela los votos: %s", cuerpo)
	}
	var got []votoVisto
	for _, v := range res.Revelacion.Votos {
		got = append(got, votoVisto{Nombre: v.Nombre, Valor: v.Valor})
	}
	if !slices.Equal(got, want) || res.Revelacion.Minimo != minimo || res.Revelacion.Maximo != maximo || !res.Revelacion.DiferenciaSignificativa {
		return fmt.Errorf("revelación de la API = %+v, se esperaba %+v, %d, %d y diferencia", res.Revelacion, want, minimo, maximo)
	}
	return nil
}

func (w *planningPokerWorld) pantallaMuestraVotos(lista string) error {
	votos, err := votosDe(lista)
	if err != nil {
		return err
	}
	for _, v := range votos {
		if !filaVotante(v.Nombre, strconv.Itoa(v.Valor)).MatchString(w.pantalla.Cuerpo) {
			return fmt.Errorf("la pantalla no muestra que %q votó %d: %.800s", v.Nombre, v.Valor, w.pantalla.Cuerpo)
		}
	}
	return nil
}

var avisoDiferencia = regexp.MustCompile(`Hay diferencias significativas[^<]*`)

func (w *planningPokerWorld) pantallaAvisaDiferencia(minimo, maximo int) error {
	aviso := avisoDiferencia.FindString(w.pantalla.Cuerpo)
	if aviso == "" {
		return fmt.Errorf("la pantalla no avisa la diferencia: %.800s", w.pantalla.Cuerpo)
	}
	if want := fmt.Sprintf("entre %d y %d", minimo, maximo); !strings.Contains(aviso, want) {
		return fmt.Errorf("el aviso %q no dice %q", aviso, want)
	}
	return nil
}
