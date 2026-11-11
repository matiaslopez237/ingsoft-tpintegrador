package bdd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/ejemplo"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu31Sistema conecta los escenarios de HU-31 a una base SQLite nueva por escenario. La carga recibe
// los casos de uso armados como en cmd/server/main.go, con el reloj del ejemplo en los que dependen
// de la fecha (RN-31.5), y el resultado se verifica con los casos de uso de consulta de las otras
// historias, como lo vería la aplicación.
type hu31Sistema struct {
	dir     string
	db      *sql.DB
	carga   ejemplo.CasosDeUso
	crear   *app.CrearProyecto
	listar  *app.ListarProyectos
	ver     *app.ConsultarProyecto
	backlog *app.ListarBacklog
	sprints *app.ListarSprints
	vel     *app.ConsultarVelocidad
	defs    *app.ListarDefectos
	horas   *app.ConsultarHorasProyecto
	rondas  *app.ListarRondas
}

func hu31NuevoSistema() (*hu31Sistema, error) {
	dir, err := os.MkdirTemp("", "bdd-hu31-")
	if err != nil {
		return nil, fmt.Errorf("crear carpeta temporal: %w", err)
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, fmt.Errorf("abrir base: %w", err)
	}

	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	sprints := persistencia.NuevoRepositorioSprints(db)
	rondas := persistencia.NuevoRepositorioRondas(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	defectos := persistencia.NuevoRepositorioDefectos(db)
	reloj := &ejemplo.Reloj{}

	return &hu31Sistema{
		dir: dir,
		db:  db,
		carga: ejemplo.CasosDeUso{
			Reloj:                   reloj,
			ListarProyectos:         app.NuevoListarProyectos(proyectos),
			CrearProyecto:           app.NuevoCrearProyecto(proyectos),
			AgregarIntegrante:       app.NuevoAgregarIntegrante(proyectos),
			CrearElemento:           app.NuevoCrearElementoBacklog(backlog),
			EditarElemento:          app.NuevoEditarElementoBacklog(backlog),
			CrearSprint:             app.NuevoCrearSprint(proyectos, sprints),
			AsignarHistoria:         app.NuevoAsignarHistoria(backlog),
			CompletarHistoria:       app.NuevoCompletarHistoria(proyectos, backlog),
			CerrarSprint:            app.NuevoCerrarSprint(proyectos, sprints),
			RegistrarEsfuerzo:       app.NuevoRegistrarEsfuerzo(proyectos, backlog, esfuerzo, reloj.Ahora),
			RegistrarDefecto:        app.NuevoRegistrarDefecto(proyectos, backlog, sprints, defectos, reloj.Ahora),
			ActualizarEstadoDefecto: app.NuevoActualizarEstadoDefecto(proyectos, sprints, defectos, reloj.Ahora),
			IniciarRonda:            app.NuevoIniciarRonda(proyectos, backlog, rondas),
			Votar:                   app.NuevoVotar(proyectos, rondas),
			RegistrarEstimacion:     app.NuevoRegistrarEstimacion(proyectos, backlog, rondas),
		},
		crear:   app.NuevoCrearProyecto(proyectos),
		listar:  app.NuevoListarProyectos(proyectos),
		ver:     app.NuevoConsultarProyecto(proyectos),
		backlog: app.NuevoListarBacklog(proyectos, backlog),
		sprints: app.NuevoListarSprints(proyectos, sprints),
		vel:     app.NuevoConsultarVelocidad(proyectos, sprints),
		defs:    app.NuevoListarDefectos(proyectos, defectos),
		horas:   app.NuevoConsultarHorasProyecto(proyectos, esfuerzo),
		rondas:  app.NuevoListarRondas(proyectos, backlog, rondas),
	}, nil
}

// cargar corre la carga con el valor de CARGAR_EJEMPLO, como main.go al arrancar.
func (s *hu31Sistema) cargar(ctx context.Context, valor string) (ejemplo.Resultado, error) {
	return ejemplo.CargarSiCorresponde(ctx, s.carga, valor)
}

// crearProyecto prepara un proyecto creado a mano, antes de la carga (CA-31.6).
func (s *hu31Sistema) crearProyecto(ctx context.Context, nombre string) error {
	_, err := s.crear.Ejecutar(ctx, app.DatosNuevoProyecto{Nombre: nombre, FechaInicio: time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local)})
	return err
}

func (s *hu31Sistema) proyectos(ctx context.Context) ([]domain.Proyecto, error) {
	return s.listar.Ejecutar(ctx)
}

func (s *hu31Sistema) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
