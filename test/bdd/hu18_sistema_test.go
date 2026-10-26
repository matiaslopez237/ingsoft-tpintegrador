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

// errHU18SinImplementar es el rojo de la tarea 18.3: los escenarios compilan y corren, pero el
// caso de uso todavía no existe. Se va cuando las tareas 18.4 a 18.6 conecten RegistrarEsfuerzo.
var errHU18SinImplementar = errors.New("HU-18 todavía no está implementada: falta conectar el caso de uso RegistrarEsfuerzo (tareas 18.4 a 18.6)")

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

// sistemaEsfuerzoReal prepara los datos con los casos de uso de HU-01, HU-02, HU-03 y HU-05,
// sobre SQLite real, en una base nueva por escenario. Registrar, buscar y contar esfuerzo
// todavía devuelven errHU18SinImplementar.
type sistemaEsfuerzoReal struct {
	crearProyecto *app.CrearProyecto
	modificar     *app.ModificarProyecto
	agregar       *app.AgregarIntegrante
	quitar        *app.QuitarIntegrante
	crearElemento *app.CrearElementoBacklog
	reloj         func() time.Time // el reloj falso del escenario; lo va a recibir RegistrarEsfuerzo
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
	return sistemaEsfuerzoReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		modificar:     app.NuevoModificarProyecto(proyectos),
		agregar:       app.NuevoAgregarIntegrante(proyectos),
		quitar:        app.NuevoQuitarIntegrante(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(persistencia.NuevoRepositorioBacklog(db)),
		reloj:         reloj,
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

func (s sistemaEsfuerzoReal) RegistrarEsfuerzo(datosEsfuerzo) (registroVisto, error) {
	return registroVisto{}, errHU18SinImplementar
}

func (s sistemaEsfuerzoReal) BuscarRegistro(int64) (registroVisto, error) {
	return registroVisto{}, errHU18SinImplementar
}

func (s sistemaEsfuerzoReal) ContarRegistros(int64) (int, error) {
	return 0, errHU18SinImplementar
}
