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

// erroresHU07 traduce la descripción del escenario al error definido en la spec (sección 7).
func erroresHU07() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"orden inválido":         domain.ErrOrdenInvalido,
		"estado inválido":        domain.ErrEstadoElementoInvalido,
	}
}

// sistemaListadoBacklogReal prepara proyectos y elementos con los casos de uso de HU-01 y HU-05,
// y lo lista con el de HU-07, sobre SQLite, en una base nueva por escenario.
type sistemaListadoBacklogReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	listar        *app.ListarBacklog
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
	proyectos, backlog := persistencia.NuevoRepositorioProyectos(db), persistencia.NuevoRepositorioBacklog(db)
	return sistemaListadoBacklogReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(backlog),
		listar:        app.NuevoListarBacklog(proyectos, backlog),
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

func (s sistemaListadoBacklogReal) ListarBacklog(proyectoID int64, orden, estado string) ([]string, error) {
	elementos, err := s.listar.Ejecutar(context.Background(), app.DatosListarBacklog{ProyectoID: proyectoID, Orden: orden, Estado: estado})
	if err != nil {
		return nil, err
	}
	titulos := make([]string, len(elementos))
	for i, e := range elementos {
		titulos[i] = e.Titulo
	}
	return titulos, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaListadoBacklogReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
