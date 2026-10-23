package bdd

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// erroresHU02 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU02() map[string]error {
	return map[string]error{
		"proyecto no encontrado":      domain.ErrProyectoNoEncontrado,
		"sin cambios":                 domain.ErrSinCambios,
		"estado inválido":             domain.ErrEstadoInvalido,
		"nombre vacío":                domain.ErrNombreVacio,
		"nombre demasiado largo":      domain.ErrNombreDemasiadoLargo,
		"nombre duplicado":            domain.ErrNombreDuplicado,
		"fecha de inicio obligatoria": domain.ErrFechaInicioObligatoria,
		"fecha de fin anterior":       domain.ErrFechaFinAnterior,
	}
}

// sistemaModificacionReal conecta los escenarios a los casos de uso de HU-01 (crear) y HU-02
// (modificar) con SQLite real, en una base nueva por escenario.
type sistemaModificacionReal struct {
	crear     *app.CrearProyecto
	modificar *app.ModificarProyecto
	repo      *persistencia.RepositorioProyectos
	db        *sql.DB
	dir       string
}

func nuevoSistemaModificacion() (sistemaModificacion, error) {
	dir, err := os.MkdirTemp("", "bdd-hu02-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repo := persistencia.NuevoRepositorioProyectos(db)
	return sistemaModificacionReal{
		crear:     app.NuevoCrearProyecto(repo),
		modificar: app.NuevoModificarProyecto(repo),
		repo:      repo,
		db:        db,
		dir:       dir,
	}, nil
}

func (s sistemaModificacionReal) Crear(nombre string, inicio, fin *time.Time) (proyectoHU02, error) {
	datos := app.DatosNuevoProyecto{Nombre: nombre, FechaFin: fin}
	if inicio != nil {
		datos.FechaInicio = *inicio
	}
	p, err := s.crear.Ejecutar(context.Background(), datos)
	if err != nil {
		return proyectoHU02{}, err
	}
	return verProyectoHU02(p), nil
}

func (s sistemaModificacionReal) Modificar(id int64, c cambiosProyecto) (proyectoHU02, error) {
	cambios := domain.CambiosProyecto{Nombre: c.Nombre, FechaInicio: c.FechaInicio}
	if c.QuitarInicio {
		cambios.FechaInicio = &time.Time{} // la fecha cero equivale a quitarla
	}
	if c.FechaFin != nil || c.QuitarFin {
		cambios.CambiarFechaFin = true
		cambios.FechaFin = c.FechaFin
	}
	if c.Estado != nil {
		e := domain.Estado(*c.Estado)
		cambios.Estado = &e
	}
	p, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{ProyectoID: id, Cambios: cambios})
	if err != nil {
		return proyectoHU02{}, err
	}
	return verProyectoHU02(p), nil
}

func (s sistemaModificacionReal) BuscarPorID(id int64) (proyectoHU02, error) {
	p, err := s.repo.BuscarPorID(context.Background(), id)
	if err != nil {
		return proyectoHU02{}, err
	}
	return verProyectoHU02(p), nil
}

func (s sistemaModificacionReal) FijarEstado(id int64, estado string) error {
	_, err := s.db.Exec(`UPDATE proyectos SET estado = ? WHERE id = ?`, estado, id)
	return err
}

func (s sistemaModificacionReal) AgregarIntegrante(id int64, nombre, rol string) error {
	return app.NuevoAgregarIntegrante(s.repo).Ejecutar(context.Background(),
		app.DatosIntegrante{ProyectoID: id, Nombre: nombre, Rol: rol})
}

func (s sistemaModificacionReal) AgregarRol(id int64, nombre, rol string) error {
	return app.NuevoAgregarRol(s.repo).Ejecutar(context.Background(),
		app.DatosRol{ProyectoID: id, Nombre: nombre, Rol: rol})
}

func (s sistemaModificacionReal) Roles(id int64, nombre string) ([]string, error) {
	p, err := s.repo.BuscarPorID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	for _, integ := range p.Integrantes {
		if domain.ClaveNombre(integ.Nombre) == domain.ClaveNombre(nombre) {
			return integ.Roles, nil
		}
	}
	return nil, domain.ErrIntegranteNoEncontrado
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaModificacionReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func verProyectoHU02(p domain.Proyecto) proyectoHU02 {
	v := proyectoHU02{
		ID:     p.ID,
		Nombre: p.Nombre,
		Inicio: p.FechaInicio.Format(time.DateOnly),
		Estado: string(p.Estado),
	}
	if p.FechaFin != nil {
		v.Fin = p.FechaFin.Format(time.DateOnly)
	}
	return v
}
