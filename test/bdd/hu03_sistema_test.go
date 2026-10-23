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
)

// errIntegrantesNoImplementado hace fallar los escenarios de HU-03 mientras
// AgregarIntegrante/QuitarIntegrante no existen (rojo esperado, Tareas 3.4 a 3.7).
var errIntegrantesNoImplementado = errors.New("pendiente: gestión de integrantes todavía no está implementada")

// erroresHU03 traduce la descripción del escenario al error definido en la spec.
// Se completa cuando existan los errores del dominio (Tarea 3.4).
func erroresHU03() map[string]error {
	return map[string]error{}
}

// sistemaIntegrantesReal crea el proyecto con el caso de uso real de HU-01 (ya
// implementado) y deja la gestión de integrantes sin implementar todavía.
type sistemaIntegrantesReal struct {
	crear *app.CrearProyecto
	db    *sql.DB
	dir   string
}

func nuevoSistemaIntegrantes() (sistemaIntegrantes, error) {
	dir, err := os.MkdirTemp("", "bdd-hu03-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	repo := persistencia.NuevoRepositorioProyectos(db)
	return sistemaIntegrantesReal{crear: app.NuevoCrearProyecto(repo), db: db, dir: dir}, nil
}

func (s sistemaIntegrantesReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaIntegrantesReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crear.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre:      nombre,
		FechaInicio: time.Now(),
	})
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (s sistemaIntegrantesReal) AgregarIntegrante(proyectoID int64, nombre, rol string) error {
	return errIntegrantesNoImplementado
}

func (s sistemaIntegrantesReal) QuitarIntegrante(proyectoID int64, nombre string) error {
	return errIntegrantesNoImplementado
}

func (s sistemaIntegrantesReal) Integrantes(proyectoID int64) ([]string, error) {
	return nil, errIntegrantesNoImplementado
}
