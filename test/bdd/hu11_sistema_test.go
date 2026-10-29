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

// hu11SistemaReal conecta los escenarios de HU-11 al caso de uso CerrarSprint con SQLite real, en
// una base nueva por escenario. Los proyectos se preparan con los casos de uso de HU-01 y HU-02, y
// el sprint se lee con el mismo repositorio que usa el caso de uso.
type hu11SistemaReal struct {
	crearProyecto     *app.CrearProyecto
	modificarProyecto *app.ModificarProyecto
	cerrar            *app.CerrarSprint
	repoSprints       *persistencia.RepositorioSprints
	db                *sql.DB
	dir               string
}

func hu11NuevoSistema() (hu11Sistema, error) {
	dir, err := os.MkdirTemp("", "bdd-hu11-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repoProyectos := persistencia.NuevoRepositorioProyectos(db)
	repoSprints := persistencia.NuevoRepositorioSprints(db)
	return &hu11SistemaReal{
		crearProyecto:     app.NuevoCrearProyecto(repoProyectos),
		modificarProyecto: app.NuevoModificarProyecto(repoProyectos),
		cerrar:            app.NuevoCerrarSprint(repoProyectos, repoSprints),
		repoSprints:       repoSprints,
		db:                db,
		dir:               dir,
	}, nil
}

func (s *hu11SistemaReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

func (s *hu11SistemaReal) FijarEstadoProyecto(proyectoID int64, estado domain.Estado) error {
	_, err := s.modificarProyecto.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &estado},
	})
	return err
}

// CrearSprint inserta el sprint directamente en la base, con el número, las fechas y el estado del
// escenario: CrearSprint de HU-08 calcula el número solo y siempre crea sprints Abiertos, y los
// escenarios necesitan, por ejemplo, un sprint 2 que ya esté Cerrado.
func (s *hu11SistemaReal) CrearSprint(proyectoID int64, numero int, inicio, fin time.Time, estado domain.EstadoSprint) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado) VALUES (?, ?, ?, ?, ?, ?)`,
		proyectoID, numero, inicio.Format(time.DateOnly), fin.Format(time.DateOnly),
		fmt.Sprintf("Objetivo del sprint %d", numero), string(estado))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *hu11SistemaReal) EjecutarCierre(proyectoID, sprintID int64) error {
	return s.cerrar.Ejecutar(context.Background(), app.DatosCerrarSprint{ProyectoID: proyectoID, SprintID: sprintID})
}

func (s *hu11SistemaReal) Sprint(sprintID int64) (domain.Sprint, error) {
	return s.repoSprints.BuscarSprintPorID(context.Background(), sprintID)
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *hu11SistemaReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
