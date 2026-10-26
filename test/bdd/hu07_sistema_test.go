package bdd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// errListarBacklogNoImplementado hace fallar los escenarios mientras ListarBacklog no existe (rojo esperado).
var errListarBacklogNoImplementado = errors.New("pendiente: ListarBacklog todavía no está implementado")

// erroresHU07 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU07() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

// sistemaListadoBacklogReal prepara proyectos y elementos con los casos de uso de HU-01 y HU-05
// sobre SQLite, en una base nueva por escenario. Listar el backlog todavía no está implementado.
type sistemaListadoBacklogReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	db            *sql.DB
	dir           string
}

func nuevoSistemaListadoBacklog() (sistemaListadoBacklog, error) {
	dir, err := os.MkdirTemp("", "bdd-hu07-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return sistemaListadoBacklogReal{
		crearProyecto: app.NuevoCrearProyecto(persistencia.NuevoRepositorioProyectos(db)),
		crearElemento: app.NuevoCrearElementoBacklog(persistencia.NuevoRepositorioBacklog(db)),
		db:            db,
		dir:           dir,
	}, nil
}

func (s sistemaListadoBacklogReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (s sistemaListadoBacklogReal) CrearElemento(proyectoID int64, titulo, prioridad string) error {
	_, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID: proyectoID, Titulo: titulo, Descripcion: "Descripción de " + titulo,
		Prioridad: domain.Prioridad(prioridad), Criterios: []string{"Criterio de " + titulo},
	})
	return err
}

func (s sistemaListadoBacklogReal) FijarEstadoElemento(proyectoID int64, titulo, estado string) error {
	res, err := s.db.Exec(`UPDATE backlog_elementos SET estado = ? WHERE proyecto_id = ? AND titulo_clave = ?`,
		estado, proyectoID, domain.ClaveNombre(titulo))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n != 1 {
		return fmt.Errorf("no se encontró el elemento %q", titulo)
	}
	return nil
}

func (s sistemaListadoBacklogReal) ListarBacklog(int64, string, string) ([]string, error) {
	return nil, errListarBacklogNoImplementado
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaListadoBacklogReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
