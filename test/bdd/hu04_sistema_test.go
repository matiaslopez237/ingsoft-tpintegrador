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

// errConsultarNoImplementado hace fallar los escenarios mientras ConsultarProyecto no existe (rojo esperado).
var errConsultarNoImplementado = errors.New("pendiente: ConsultarProyecto todavía no está implementado")

// erroresHU04 traduce la descripción del escenario al error definido en la spec.
func erroresHU04() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
	}
}

// sistemaConsultaReal prepara los datos con los casos de uso de HU-01 y HU-03 sobre SQLite real,
// en una base nueva por escenario.
type sistemaConsultaReal struct {
	crear      *app.CrearProyecto
	agregar    *app.AgregarIntegrante
	agregarRol *app.AgregarRol
	db         *sql.DB
	dir        string
}

func nuevoSistemaConsulta() (sistemaConsulta, error) {
	dir, err := os.MkdirTemp("", "bdd-hu04-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repo := persistencia.NuevoRepositorioProyectos(db)
	return sistemaConsultaReal{
		crear:      app.NuevoCrearProyecto(repo),
		agregar:    app.NuevoAgregarIntegrante(repo),
		agregarRol: app.NuevoAgregarRol(repo),
		db:         db,
		dir:        dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaConsultaReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaConsultaReal) CrearProyecto(nombre string, inicio time.Time, fin *time.Time) (int64, error) {
	p, err := s.crear.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: inicio, FechaFin: fin,
	})
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (s sistemaConsultaReal) AgregarIntegrante(proyectoID int64, nombre, rol string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: rol,
	})
}

func (s sistemaConsultaReal) AgregarRol(proyectoID int64, nombre, rol string) error {
	return s.agregarRol.Ejecutar(context.Background(), app.DatosRol{
		ProyectoID: proyectoID, Nombre: nombre, Rol: rol,
	})
}

// Consultar se conecta al caso de uso ConsultarProyecto en la Tarea 4.4.
func (s sistemaConsultaReal) Consultar(int64) (domain.Proyecto, error) {
	return domain.Proyecto{}, errConsultarNoImplementado
}
