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

// errModificarNoImplementado hace fallar los escenarios mientras ModificarProyecto no existe (rojo esperado).
var errModificarNoImplementado = errors.New("pendiente: ModificarProyecto todavía no está implementado")

// erroresHU02 traduce la descripción del escenario al error definido en la spec (sección 7).
// "estado inválido" (ErrEstadoInvalido) y "sin cambios" (ErrSinCambios) se agregan
// cuando existan esos errores en el dominio (Tarea 2.4).
func erroresHU02() map[string]error {
	return map[string]error{
		"proyecto no encontrado":      domain.ErrProyectoNoEncontrado,
		"nombre vacío":                domain.ErrNombreVacio,
		"nombre demasiado largo":      domain.ErrNombreDemasiadoLargo,
		"nombre duplicado":            domain.ErrNombreDuplicado,
		"fecha de inicio obligatoria": domain.ErrFechaInicioObligatoria,
		"fecha de fin anterior":       domain.ErrFechaFinAnterior,
	}
}

// sistemaModificacionReal prepara los escenarios con el código real de HU-01 (crear y buscar)
// sobre SQLite, en una base nueva por escenario. Modificar todavía no está implementado.
type sistemaModificacionReal struct {
	crear *app.CrearProyecto
	repo  *persistencia.RepositorioProyectos
	db    *sql.DB
	dir   string
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
	return sistemaModificacionReal{crear: app.NuevoCrearProyecto(repo), repo: repo, db: db, dir: dir}, nil
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

func (s sistemaModificacionReal) Modificar(int64, cambiosProyecto) (proyectoHU02, error) {
	return proyectoHU02{}, errModificarNoImplementado
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
