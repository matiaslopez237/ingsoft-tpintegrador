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
	"strconv"
	"strings"
	"time"

	adapterhttp "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/http"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU06 traduce la descripción del escenario al error definido en la spec.
func erroresHU06() map[string]error {
	return map[string]error{
		"elemento no encontrado": domain.ErrElementoNoEncontrado,
		"sin cambios":            domain.ErrSinCambios,
		"título vacío":           domain.ErrTituloVacio,
		"descripción vacía":      domain.ErrDescripcionVacia,
		"prioridad inválida":     domain.ErrPrioridadInvalida,
		"story points inválidos": domain.ErrStoryPointsInvalidos,
		"sin criterios":          domain.ErrSinCriterios,
		"estado inválido":        domain.ErrEstadoElementoInvalido,
		"título duplicado":       domain.ErrTituloDuplicado,
		"elemento con esfuerzo":  domain.ErrElementoConEsfuerzo,
	}
}

// sistemaEdicionReal usa los casos de uso reales sobre SQLite, en una base nueva por escenario,
// y levanta un servidor de prueba con el mismo router que main.go para los escenarios de pantalla.
type sistemaEdicionReal struct {
	crearProyecto *app.CrearProyecto
	agregar       *app.AgregarIntegrante
	crear         *app.CrearElementoBacklog
	editar        *app.EditarElementoBacklog
	eliminar      *app.EliminarElementoBacklog
	registrar     *app.RegistrarEsfuerzo
	repo          *persistencia.RepositorioBacklog
	servidor      *httptest.Server
	cliente       *http.Client
	db            *sql.DB
	dir           string
}

func nuevoSistemaEdicionBacklog() (sistemaEdicionBacklog, error) {
	dir, err := os.MkdirTemp("", "bdd-hu06-")
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
	reloj := func() time.Time { return time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC) }
	s := sistemaEdicionReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		agregar:       app.NuevoAgregarIntegrante(proyectos),
		crear:         app.NuevoCrearElementoBacklog(backlog),
		editar:        app.NuevoEditarElementoBacklog(backlog),
		eliminar:      app.NuevoEliminarElementoBacklog(backlog),
		registrar:     app.NuevoRegistrarEsfuerzo(proyectos, backlog, persistencia.NuevoRepositorioEsfuerzo(db), reloj),
		repo:          backlog,
		db:            db,
		dir:           dir,
	}
	s.servidor = httptest.NewServer(adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		ConsultarProyecto:        app.NuevoConsultarProyecto(proyectos),
		ListarProyectos:          app.NuevoListarProyectos(proyectos),
		EditarElementoBacklog:    s.editar,
		EliminarElementoBacklog:  s.eliminar,
		ConsultarElementoBacklog: app.NuevoConsultarElementoBacklog(backlog),
	}))
	s.cliente = s.servidor.Client()
	s.cliente.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return s, nil
}

// Cerrar apaga el servidor, cierra la base del escenario y borra su carpeta temporal.
func (s sistemaEdicionReal) Cerrar() error {
	s.servidor.Close()
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaEdicionReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

func (s sistemaEdicionReal) CrearElemento(proyectoID int64, titulo, prioridad string, storyPoints int) (int64, error) {
	e, err := s.crear.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID: proyectoID, Titulo: titulo, Descripcion: "Necesidad de prueba",
		Prioridad: domain.Prioridad(prioridad), StoryPoints: &storyPoints, Criterios: []string{"Criterio uno"},
	})
	return e.ID, err
}

func (s sistemaEdicionReal) RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64) error {
	ctx := context.Background()
	if err := s.agregar.Ejecutar(ctx, app.DatosIntegrante{ProyectoID: proyectoID, Nombre: integrante, Rol: "Desarrollo"}); err != nil {
		return err
	}
	_, err := s.registrar.Ejecutar(ctx, app.DatosRegistroEsfuerzo{
		ProyectoID: proyectoID, ElementoID: elementoID, Integrante: integrante,
		Fecha: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), Actividad: "Desarrollo", Horas: horas,
	})
	return err
}

func (s sistemaEdicionReal) EditarElemento(proyectoID, elementoID int64, campos map[string]string) error {
	c, err := cambiosDeCampos(campos)
	if err != nil {
		return err
	}
	_, err = s.editar.Ejecutar(context.Background(), app.DatosEditarElemento{ProyectoID: proyectoID, ElementoID: elementoID, Cambios: c})
	return err
}

// cambiosDeCampos traduce los campos del escenario a los cambios del dominio.
func cambiosDeCampos(campos map[string]string) (domain.CambiosElemento, error) {
	var c domain.CambiosElemento
	for nombre, valor := range campos {
		v := valor
		switch nombre {
		case "titulo":
			c.Titulo = &v
		case "descripcion":
			c.Descripcion = &v
		case "prioridad":
			p := domain.Prioridad(v)
			c.Prioridad = &p
		case "estado":
			e := domain.EstadoElemento(v)
			c.Estado = &e
		case "story_points":
			c.CambiarStoryPoints = true
			if v != "null" {
				n, err := strconv.Atoi(v)
				if err != nil {
					return c, err
				}
				c.StoryPoints = &n
			}
		case "criterios":
			lista := []string{}
			if v != "" {
				for _, criterio := range strings.Split(v, ";") {
					lista = append(lista, strings.TrimSpace(criterio))
				}
			}
			c.Criterios = &lista
		default:
			return c, errors.New("campo desconocido en el escenario: " + nombre)
		}
	}
	return c, nil
}

func (s sistemaEdicionReal) EliminarElemento(proyectoID, elementoID int64) error {
	return s.eliminar.Ejecutar(context.Background(), app.DatosEliminarElemento{ProyectoID: proyectoID, ElementoID: elementoID})
}

func (s sistemaEdicionReal) BuscarElemento(id int64) (elementoEditable, error) {
	e, err := s.repo.BuscarElementoPorID(context.Background(), id)
	if err != nil {
		return elementoEditable{}, err
	}
	v := elementoEditable{
		ID: e.ID, Titulo: e.Titulo, Descripcion: e.Descripcion, Prioridad: string(e.Prioridad),
		StoryPoints: "sin estimar", Criterios: strings.Join(e.Criterios, "; "), Estado: string(e.Estado),
	}
	if e.StoryPoints != nil {
		v.StoryPoints = strconv.Itoa(*e.StoryPoints)
	}
	return v, nil
}

func (s sistemaEdicionReal) ContarElementos(proyectoID int64) (int, error) {
	return s.repo.ContarElementos(context.Background(), proyectoID)
}

func (s sistemaEdicionReal) Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error) {
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
