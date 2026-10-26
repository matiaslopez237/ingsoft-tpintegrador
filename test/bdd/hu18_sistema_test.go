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

// erroresHU18 traduce la descripción del escenario al error definido en la spec.
func erroresHU18() map[string]error {
	return map[string]error{
		"proyecto no encontrado":          domain.ErrProyectoNoEncontrado,
		"elemento no encontrado":          domain.ErrElementoNoEncontrado,
		"integrante vacío":                domain.ErrNombreIntegranteVacio,
		"integrante no encontrado":        domain.ErrIntegranteNoEncontrado,
		"fecha obligatoria":               domain.ErrFechaEsfuerzoObligatoria,
		"fecha futura":                    domain.ErrFechaEsfuerzoFutura,
		"actividad vacía":                 domain.ErrActividadVacia,
		"horas inválidas":                 domain.ErrHorasInvalidas,
		"horas no múltiplo de media hora": domain.ErrHorasNoMultiploMediaHora,
		"horas exceden el máximo":         domain.ErrHorasExcedenMaximo,
	}
}

// sistemaEsfuerzoReal prepara los datos con los casos de uso de HU-01, HU-02, HU-03 y HU-05 y
// registra con RegistrarEsfuerzo (HU-18), sobre SQLite real, en una base nueva por escenario.
type sistemaEsfuerzoReal struct {
	crearProyecto *app.CrearProyecto
	modificar     *app.ModificarProyecto
	agregar       *app.AgregarIntegrante
	quitar        *app.QuitarIntegrante
	crearElemento *app.CrearElementoBacklog
	registrar     *app.RegistrarEsfuerzo
	esfuerzo      *persistencia.RepositorioEsfuerzo
	db            *sql.DB
	dir           string
}

func nuevoSistemaEsfuerzo(reloj func() time.Time) (sistemaEsfuerzo, error) {
	dir, err := os.MkdirTemp("", "bdd-hu18-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	return sistemaEsfuerzoReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		modificar:     app.NuevoModificarProyecto(proyectos),
		agregar:       app.NuevoAgregarIntegrante(proyectos),
		quitar:        app.NuevoQuitarIntegrante(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(backlog),
		registrar:     app.NuevoRegistrarEsfuerzo(proyectos, backlog, esfuerzo, reloj),
		esfuerzo:      esfuerzo,
		db:            db,
		dir:           dir,
	}, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s sistemaEsfuerzoReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}

func (s sistemaEsfuerzoReal) CrearProyecto(nombre string, inicio time.Time) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: inicio,
	})
	return p.ID, err
}

// AgregarIntegrante usa un rol fijo: los escenarios de HU-18 no miran los roles.
func (s sistemaEsfuerzoReal) AgregarIntegrante(proyectoID int64, nombre string) error {
	return s.agregar.Ejecutar(context.Background(), app.DatosIntegrante{
		ProyectoID: proyectoID, Nombre: nombre, Rol: "Desarrolladora",
	})
}

func (s sistemaEsfuerzoReal) QuitarIntegrante(proyectoID int64, nombre string) error {
	return s.quitar.Ejecutar(context.Background(), app.DatosQuitarIntegrante{
		ProyectoID: proyectoID, Nombre: nombre,
	})
}

func (s sistemaEsfuerzoReal) CambiarEstado(proyectoID int64, estado string) error {
	e := domain.Estado(estado)
	_, err := s.modificar.Ejecutar(context.Background(), app.DatosModificarProyecto{
		ProyectoID: proyectoID, Cambios: domain.CambiosProyecto{Estado: &e},
	})
	return err
}

// CrearElemento completa con datos fijos lo que HU-05 exige y los escenarios de HU-18 no miran.
func (s sistemaEsfuerzoReal) CrearElemento(proyectoID int64, titulo string) (int64, error) {
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Elemento para registrar esfuerzo",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Se puede registrar esfuerzo"},
	})
	return e.ID, err
}

func (s sistemaEsfuerzoReal) RegistrarEsfuerzo(d datosEsfuerzo) (registroVisto, error) {
	r, err := s.registrar.Ejecutar(context.Background(), app.DatosRegistroEsfuerzo{
		ProyectoID: d.ProyectoID,
		ElementoID: d.ElementoID,
		Integrante: d.Integrante,
		Fecha:      d.Fecha,
		Actividad:  d.Actividad,
		Horas:      d.Horas,
	})
	if err != nil {
		return registroVisto{}, err
	}
	return verRegistro(r), nil
}

func (s sistemaEsfuerzoReal) BuscarRegistro(id int64) (registroVisto, error) {
	r, err := s.esfuerzo.BuscarRegistroPorID(context.Background(), id)
	if err != nil {
		return registroVisto{}, err
	}
	return verRegistro(r), nil
}

func (s sistemaEsfuerzoReal) ContarRegistros(elementoID int64) (int, error) {
	return s.esfuerzo.ContarRegistros(context.Background(), elementoID)
}

func verRegistro(r domain.RegistroEsfuerzo) registroVisto {
	return registroVisto{
		ID:         r.ID,
		ProyectoID: r.ProyectoID,
		ElementoID: r.ElementoID,
		Integrante: r.Integrante,
		Fecha:      r.Fecha.Format(time.DateOnly),
		Actividad:  r.Actividad,
		Horas:      r.Horas,
	}
}
