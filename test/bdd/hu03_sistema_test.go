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

// erroresHU03 traduce la descripción del escenario al error definido en la spec.
func erroresHU03() map[string]error {
	return map[string]error{
		"nombre vacío":             domain.ErrNombreIntegranteVacio,
		"rol vacío":                domain.ErrRolIntegranteVacio,
		"integrante duplicado":     domain.ErrIntegranteDuplicado,
		"integrante no encontrado": domain.ErrIntegranteNoEncontrado,
	}
}

// sistemaIntegrantesReal conecta los escenarios a los casos de uso de HU-03 con SQLite real,
// en una base nueva por escenario.
type sistemaIntegrantesReal struct {
	crear   *app.CrearProyecto
	agregar *app.AgregarIntegrante
	quitar  *app.QuitarIntegrante
	repo    *persistencia.RepositorioProyectos
	db      *sql.DB
	dir     string
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
	return sistemaIntegrantesReal{
		crear:   app.NuevoCrearProyecto(repo),
		agregar: app.NuevoAgregarIntegrante(repo),
		quitar:  app.NuevoQuitarIntegrante(repo),
		repo:    repo,
		db:      db,
		dir:     dir,
	}, nil
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
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: rol,
	})
}

func (s sistemaIntegrantesReal) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s sistemaIntegrantesReal) Integrantes(proyectoID int64) ([]string, error) {
	p, err := s.repo.BuscarPorID(context.Background(), proyectoID)
	if err != nil {
		return nil, err
	}
	nombres := make([]string, len(p.Integrantes))
	for i, integ := range p.Integrantes {
		nombres[i] = integ.Nombre
	}
	return nombres, nil
}

// errAgregarRolNoImplementado hace fallar los escenarios de roles mientras AgregarRol no existe (rojo esperado).
var errAgregarRolNoImplementado = errors.New("pendiente: AgregarRol todavía no está implementado")

func (s sistemaIntegrantesReal) AgregarRol(int64, string, string) error {
	return errAgregarRolNoImplementado
}

// Roles devuelve los roles del integrante en el orden en que se agregaron.
// Mientras el dominio tenga un solo rol por integrante, devuelve ese.
func (s sistemaIntegrantesReal) Roles(proyectoID int64, nombre string) ([]string, error) {
	p, err := s.repo.BuscarPorID(context.Background(), proyectoID)
	if err != nil {
		return nil, err
	}
	for _, integ := range p.Integrantes {
		if domain.ClaveNombre(integ.Nombre) == domain.ClaveNombre(nombre) {
			return []string{integ.Rol}, nil
		}
	}
	return nil, fmt.Errorf("%q no es integrante del proyecto", nombre)
}
