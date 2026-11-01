package bdd

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// erroresHU14 traduce la descripción del escenario al error definido en la spec.
func erroresHU14() map[string]error {
	return map[string]error{
		"elemento no encontrado": domain.ErrElementoNoEncontrado,
		"ronda no encontrada":    domain.ErrRondaNoEncontrada,
		"elemento terminado":     domain.ErrElementoTerminado,
		"sin votantes":           domain.ErrRondaSinVotantes,
		"ronda ya abierta":       domain.ErrRondaYaAbierta,
		"ronda cerrada":          domain.ErrRondaCerrada,
		"no es votante":          domain.ErrNoEsVotante,
		"voto inválido":          domain.ErrVotoInvalido,
	}
}

// sistemaPlanningPokerReal usa los casos de uso reales sobre SQLite, en una base nueva por
// escenario, y levanta un servidor de prueba con el mismo router que main.go.
type sistemaPlanningPokerReal struct {
	crearProyecto *app.CrearProyecto
	agregar       *app.AgregarIntegrante
	crear         *app.CrearElementoBacklog
	editar        *app.EditarElementoBacklog
	servidor      *httptest.Server
	cliente       *http.Client
	db            *sql.DB
	dir           string
}

func nuevoSistemaPlanningPoker() (sistemaPlanningPoker, error) {
	dir, err := os.MkdirTemp("", "bdd-hu14-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	s := sistemaPlanningPokerReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		agregar:       app.NuevoAgregarIntegrante(proyectos),
		crear:         app.NuevoCrearElementoBacklog(backlog),
		editar:        app.NuevoEditarElementoBacklog(backlog),
		db:            db,
		dir:           dir,
	}
	s.servidor = httptest.NewServer(adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		ConsultarProyecto:        app.NuevoConsultarProyecto(proyectos),
		ListarProyectos:          app.NuevoListarProyectos(proyectos),
		EditarElementoBacklog:    s.editar,
		EliminarElementoBacklog:  app.NuevoEliminarElementoBacklog(backlog),
		ConsultarElementoBacklog: app.NuevoConsultarElementoBacklog(backlog),
	}))
	s.cliente = s.servidor.Client()
	s.cliente.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return s, nil
}

// Cerrar apaga el servidor, cierra la base del escenario y borra su carpeta temporal.
func (s sistemaPlanningPokerReal) Cerrar() error {
	s.servidor.Close()
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaPlanningPokerReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

func (s sistemaPlanningPokerReal) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrollo"})
}

// CrearElemento crea el elemento con HU-05 y, si el escenario pide otro estado, lo cambia con HU-06.
func (s sistemaPlanningPokerReal) CrearElemento(proyectoID int64, titulo, estado string) (int64, error) {
	ctx := context.Background()
	e, err := s.crear.Ejecutar(ctx, app.DatosNuevoElemento{
		ProyectoID: proyectoID, Titulo: titulo, Descripcion: "Necesidad de prueba",
		Prioridad: domain.PrioridadAlta, Criterios: []string{"Criterio uno"},
	})
	if err != nil || estado == string(domain.EstadoPendiente) {
		return e.ID, err
	}
	nuevo := domain.EstadoElemento(estado)
	_, err = s.editar.Ejecutar(ctx, app.DatosEditarElemento{
		ProyectoID: proyectoID, ElementoID: e.ID, Cambios: domain.CambiosElemento{Estado: &nuevo},
	})
	return e.ID, err
}

// IniciarRonda, Votar y ConsultarRonda quedan pendientes hasta que existan los casos de uso de
// HU-14 (tarea 14.5): mientras tanto, los escenarios quedan en rojo.
func (s sistemaPlanningPokerReal) IniciarRonda(proyectoID, elementoID int64) (rondaVista, error) {
	return rondaVista{}, godog.ErrPending
}

func (s sistemaPlanningPokerReal) Votar(proyectoID, rondaID int64, integrante string, valor int) (rondaVista, error) {
	return rondaVista{}, godog.ErrPending
}

func (s sistemaPlanningPokerReal) ConsultarRonda(proyectoID, rondaID int64) (rondaVista, error) {
	return rondaVista{}, godog.ErrPending
}

func (s sistemaPlanningPokerReal) ResumenJSON(proyectoID, rondaID int64) (string, error) {
	return "", godog.ErrPending
}

func (s sistemaPlanningPokerReal) Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error) {
	var cuerpo io.Reader
	if form != nil {
		cuerpo = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(metodo, s.servidor.URL+ruta, cuerpo)
	if err != nil {
		return respuestaWeb{}, err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	res, err := s.cliente.Do(req)
	if err != nil {
		return respuestaWeb{}, err
	}
	defer res.Body.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(res.Body); err != nil {
		return respuestaWeb{}, err
	}
	return respuestaWeb{Estado: res.StatusCode, Ubicacion: res.Header.Get("Location"), Cuerpo: b.String()}, nil
}
