package bdd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU01 traduce la descripción del escenario al error definido en la spec.
func erroresHU01() map[string]error {
	return map[string]error{
		"nombre vacío":                domain.ErrNombreVacio,
		"nombre demasiado largo":      domain.ErrNombreDemasiadoLargo,
		"nombre duplicado":            domain.ErrNombreDuplicado,
		"fecha de inicio obligatoria": domain.ErrFechaInicioObligatoria,
		"fecha de fin anterior":       domain.ErrFechaFinAnterior,
	}
}

// sistemaReal conecta los escenarios al caso de uso CrearProyecto con SQLite real,
// en una base nueva por escenario.
type sistemaReal struct {
	crear *app.CrearProyecto
	repo  *persistencia.RepositorioProyectos
}

func nuevoSistemaProyectos() (sistemaProyectos, error) {
	dir, err := os.MkdirTemp("", "bdd-hu01-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		return nil, err
	}
	repo := persistencia.NuevoRepositorioProyectos(db)
	return sistemaReal{crear: app.NuevoCrearProyecto(repo), repo: repo}, nil
}

func (s sistemaReal) Crear(nombre string, inicio, fin *time.Time) (proyectoVisto, error) {
	datos := app.DatosNuevoProyecto{Nombre: nombre, FechaFin: fin}
	if inicio != nil {
		datos.FechaInicio = *inicio
	}
	p, err := s.crear.Ejecutar(context.Background(), datos)
	if err != nil {
		return proyectoVisto{}, err
	}
	return verProyecto(p), nil
}

func (s sistemaReal) BuscarPorID(id int64) (proyectoVisto, error) {
	p, err := s.repo.BuscarPorID(context.Background(), id)
	if err != nil {
		return proyectoVisto{}, err
	}
	return verProyecto(p), nil
}

func (s sistemaReal) Cantidad() (int, error) {
	return s.repo.Contar(context.Background())
}

func verProyecto(p domain.Proyecto) proyectoVisto {
	return proyectoVisto{ID: p.ID, Nombre: p.Nombre, Estado: string(p.Estado), TieneFin: p.FechaFin != nil}
}
