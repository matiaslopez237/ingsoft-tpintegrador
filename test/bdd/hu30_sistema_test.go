package bdd

import (
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

	adapterhttp "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/http"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
)

// sistemaWebReal levanta un servidor HTTP de prueba con el mismo router que main.go, sobre
// SQLite real, en una base nueva por escenario. Los pedidos no siguen redirecciones, así los
// escenarios pueden verificarlas.
type sistemaWebReal struct {
	servidor *httptest.Server
	cliente  *http.Client
	repo     *persistencia.RepositorioProyectos
	crear    *app.CrearProyecto
	agregar  *app.AgregarIntegrante
	db       *sql.DB
	dir      string
}

func nuevoSistemaWeb() (sistemaWeb, error) {
	dir, err := os.MkdirTemp("", "bdd-hu30-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repo := persistencia.NuevoRepositorioProyectos(db)
	mux := adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		CrearProyecto:     app.NuevoCrearProyecto(repo),
		ModificarProyecto: app.NuevoModificarProyecto(repo),
		AgregarIntegrante: app.NuevoAgregarIntegrante(repo),
		QuitarIntegrante:  app.NuevoQuitarIntegrante(repo),
		AgregarRol:        app.NuevoAgregarRol(repo),
		ConsultarProyecto: app.NuevoConsultarProyecto(repo),
		ListarProyectos:   app.NuevoListarProyectos(repo),
	})
	servidor := httptest.NewServer(mux)
	cliente := servidor.Client()
	cliente.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return sistemaWebReal{
		servidor: servidor,
		cliente:  cliente,
		repo:     repo,
		crear:    app.NuevoCrearProyecto(repo),
		agregar:  app.NuevoAgregarIntegrante(repo),
		db:       db,
		dir:      dir,
	}, nil
}

// Cerrar apaga el servidor, cierra la base del escenario y borra su carpeta temporal.
func (s sistemaWebReal) Cerrar() error {
	s.servidor.Close()
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaWebReal) GuardarProyecto(nombre string) (int64, error) {
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	p, err := s.crear.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    &fin,
	})
	return p.ID, err
}

func (s sistemaWebReal) AgregarIntegrante(proyectoID int64, nombre, rol string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{ProyectoID: proyectoID, Nombre: nombre, Rol: rol})
}

func (s sistemaWebReal) Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error) {
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
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return respuestaWeb{}, err
	}
	return respuestaWeb{Estado: res.StatusCode, Ubicacion: res.Header.Get("Location"), Cuerpo: string(b)}, nil
}

func (s sistemaWebReal) ContarProyectos() (int, error) {
	return s.repo.Contar(context.Background())
}

func (s sistemaWebReal) BuscarProyecto(id int64) (proyectoGuardado, error) {
	p, err := s.repo.BuscarPorID(context.Background(), id)
	if err != nil {
		return proyectoGuardado{}, err
	}
	g := proyectoGuardado{Nombre: p.Nombre, Estado: string(p.Estado), TieneFin: p.FechaFin != nil}
	if p.FechaFin != nil {
		g.Fin = p.FechaFin.Format(time.DateOnly)
	}
	return g, nil
}
