package bdd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"

	adapterhttp "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/http"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// mensajeEscala es el texto con que se informa la escala de story points (HU-13, sección 7).
const mensajeEscala = "1, 2, 3, 5, 8, 13 o 21"

// escalaWorld usa los casos de uso reales de HU-05 y HU-06 sobre SQLite, en una base nueva por
// escenario, y un servidor de prueba con el mismo router que main.go para el escenario de la API.
type escalaWorld struct {
	db        *sql.DB
	dir       string
	servidor  *httptest.Server
	crear     *app.CrearElementoBacklog
	editar    *app.EditarElementoBacklog
	backlog   *persistencia.RepositorioBacklog
	proyectos map[string]int64
	elementos map[string]int64 // título → ID
	ultimoErr error
	respuesta struct {
		estado int
		cuerpo string
	}
}

func registrarStepsHU13(ctx *godog.ScenarioContext) {
	w := &escalaWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		return c, w.preparar()
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		w.servidor.Close()
		return c, errors.Join(w.db.Close(), os.RemoveAll(w.dir))
	})

	ctx.Step(`^el proyecto "([^"]*)" para estimar con story points$`, w.proyecto)
	ctx.Step(`^en "([^"]*)" existe el elemento "([^"]*)" estimado en (\d+) story points$`, w.existeEstimado)
	ctx.Step(`^en "([^"]*)" existe el elemento "([^"]*)" guardado antes de la escala con (\d+) story points$`, w.existeAntesDeLaEscala)

	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" estimado en (-?\d+) story points$`, w.creoEstimado)
	ctx.Step(`^creo en "([^"]*)" el elemento "([^"]*)" sin estimar$`, w.creoSinEstimar)
	ctx.Step(`^cambio los story points de "([^"]*)" a (-?\d+)$`, w.cambioStoryPoints)
	ctx.Step(`^le quito la estimación a "([^"]*)"$`, w.quitoEstimacion)
	ctx.Step(`^cambio el título de "([^"]*)" a "([^"]*)"$`, w.cambioTitulo)
	ctx.Step(`^por la API creo en "([^"]*)" el elemento "([^"]*)" con (-?\d+) story points$`, w.creoPorAPI)

	ctx.Step(`^el elemento estimado "([^"]*)" tiene (\d+) story points$`, w.tienePuntos)
	ctx.Step(`^el elemento estimado "([^"]*)" queda sin estimar$`, w.quedaSinEstimar)
	ctx.Step(`^se rechaza la estimación con el error de la escala$`, w.rechazada)
	ctx.Step(`^el backlog de "([^"]*)" no tiene el elemento "([^"]*)"$`, w.noTieneElemento)
	ctx.Step(`^la API responde 400 con el mensaje de la escala$`, w.apiRechaza)
}

func (w *escalaWorld) preparar() error {
	dir, err := os.MkdirTemp("", "bdd-hu13-")
	if err != nil {
		return err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	*w = escalaWorld{
		db: db, dir: dir, backlog: backlog,
		crear:     app.NuevoCrearElementoBacklog(backlog),
		editar:    app.NuevoEditarElementoBacklog(backlog),
		proyectos: map[string]int64{}, elementos: map[string]int64{},
	}
	w.servidor = httptest.NewServer(adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		CrearProyecto:         app.NuevoCrearProyecto(proyectos),
		ConsultarProyecto:     app.NuevoConsultarProyecto(proyectos),
		CrearElementoBacklog:  w.crear,
		EditarElementoBacklog: w.editar,
	}))
	return nil
}

// --- Given ---

func (w *escalaWorld) proyecto(nombre string) error {
	p, err := app.NuevoCrearProyecto(persistencia.NuevoRepositorioProyectos(w.db)).Ejecutar(context.Background(),
		app.DatosNuevoProyecto{Nombre: nombre, FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		return fmt.Errorf("preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = p.ID
	return nil
}

func (w *escalaWorld) crearElemento(proyecto, titulo string, puntos *int) error {
	e, err := w.crear.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID: w.proyectos[proyecto], Titulo: titulo, Descripcion: "Necesidad de prueba",
		Prioridad: domain.PrioridadAlta, StoryPoints: puntos, Criterios: []string{"Criterio uno"},
	})
	w.ultimoErr = err
	if err == nil {
		w.elementos[titulo] = e.ID
	}
	return nil
}

func (w *escalaWorld) existeEstimado(proyecto, titulo string, puntos int) error {
	if err := w.crearElemento(proyecto, titulo, &puntos); err != nil || w.ultimoErr != nil {
		return fmt.Errorf("preparar %q con %d story points: %w", titulo, puntos, errors.Join(err, w.ultimoErr))
	}
	return nil
}

// existeAntesDeLaEscala crea el elemento sin estimar y le escribe los story points directo en la
// base, como quedó un elemento cargado antes de HU-13 (RN-13.3).
func (w *escalaWorld) existeAntesDeLaEscala(proyecto, titulo string, puntos int) error {
	if err := w.crearElemento(proyecto, titulo, nil); err != nil || w.ultimoErr != nil {
		return fmt.Errorf("preparar %q: %w", titulo, errors.Join(err, w.ultimoErr))
	}
	_, err := w.db.Exec(`UPDATE backlog_elementos SET story_points = ? WHERE id = ?`, puntos, w.elementos[titulo])
	return err
}

// --- When ---

func (w *escalaWorld) creoEstimado(proyecto, titulo string, puntos int) error {
	return w.crearElemento(proyecto, titulo, &puntos)
}

func (w *escalaWorld) creoSinEstimar(proyecto, titulo string) error {
	return w.crearElemento(proyecto, titulo, nil)
}

func (w *escalaWorld) editarElemento(titulo string, c domain.CambiosElemento) error {
	id, ok := w.elementos[titulo]
	if !ok {
		return fmt.Errorf("el escenario no conoce el elemento %q", titulo)
	}
	e, err := w.backlog.BuscarElementoPorID(context.Background(), id)
	if err != nil {
		return err
	}
	_, w.ultimoErr = w.editar.Ejecutar(context.Background(), app.DatosEditarElemento{ProyectoID: e.ProyectoID, ElementoID: id, Cambios: c})
	return nil
}

func (w *escalaWorld) cambioStoryPoints(titulo string, puntos int) error {
	return w.editarElemento(titulo, domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &puntos})
}

func (w *escalaWorld) quitoEstimacion(titulo string) error {
	return w.editarElemento(titulo, domain.CambiosElemento{CambiarStoryPoints: true})
}

func (w *escalaWorld) cambioTitulo(titulo, nuevo string) error {
	if err := w.editarElemento(titulo, domain.CambiosElemento{Titulo: &nuevo}); err != nil {
		return err
	}
	w.elementos[nuevo] = w.elementos[titulo]
	return nil
}

func (w *escalaWorld) creoPorAPI(proyecto, titulo string, puntos int) error {
	cuerpo := fmt.Sprintf(`{"titulo":%q,"descripcion":"Necesidad de prueba","prioridad":"Alta","story_points":%d,"criterios":["Criterio uno"]}`, titulo, puntos)
	res, err := http.Post(fmt.Sprintf("%s/proyectos/%d/backlog", w.servidor.URL, w.proyectos[proyecto]), "application/json", strings.NewReader(cuerpo))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	w.respuesta.estado, w.respuesta.cuerpo = res.StatusCode, string(b)
	return err
}

// --- Then ---

func (w *escalaWorld) elemento(titulo string) (domain.ElementoBacklog, error) {
	if w.ultimoErr != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("no se esperaba un error: %w", w.ultimoErr)
	}
	id, ok := w.elementos[titulo]
	if !ok {
		return domain.ElementoBacklog{}, fmt.Errorf("no se creó el elemento %q", titulo)
	}
	return w.backlog.BuscarElementoPorID(context.Background(), id)
}

func (w *escalaWorld) tienePuntos(titulo string, puntos int) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if e.Titulo != titulo || e.StoryPoints == nil || *e.StoryPoints != puntos {
		return fmt.Errorf("%q tiene título %q y story points %v, se esperaba %d", titulo, e.Titulo, e.StoryPoints, puntos)
	}
	return nil
}

func (w *escalaWorld) quedaSinEstimar(titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	if e.StoryPoints != nil {
		return fmt.Errorf("%q tiene %d story points, se esperaba sin estimar", titulo, *e.StoryPoints)
	}
	return nil
}

func (w *escalaWorld) rechazada() error {
	if !errors.Is(w.ultimoErr, domain.ErrStoryPointsInvalidos) || !strings.Contains(w.ultimoErr.Error(), mensajeEscala) {
		return fmt.Errorf("se esperaba el error de la escala (%q), se obtuvo: %v", mensajeEscala, w.ultimoErr)
	}
	w.ultimoErr = nil // el Then siguiente puede mirar el elemento
	return nil
}

func (w *escalaWorld) noTieneElemento(proyecto, titulo string) error {
	n, err := w.backlog.ContarElementos(context.Background(), w.proyectos[proyecto])
	if err != nil {
		return err
	}
	if _, creado := w.elementos[titulo]; creado || n != 0 {
		return fmt.Errorf("el backlog de %q tiene %d elementos, se esperaba que no estuviera %q", proyecto, n, titulo)
	}
	return nil
}

func (w *escalaWorld) apiRechaza() error {
	if w.respuesta.estado != http.StatusBadRequest || !strings.Contains(w.respuesta.cuerpo, mensajeEscala) {
		return fmt.Errorf("la API respondió %d: %s; se esperaba 400 con %q", w.respuesta.estado, w.respuesta.cuerpo, mensajeEscala)
	}
	return nil
}
