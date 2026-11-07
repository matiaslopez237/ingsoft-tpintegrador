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

// sistemaDashboardReal prepara el proyecto con el caso de uso de HU-01 y los elementos con el de
// HU-05; los estados, los sprints, los defectos y las horas los escribe directo en la base (igual que
// HU-22 y HU-24), y completa historias con el caso de uso de HU-10. Consulta el Dashboard con el caso
// de uso de HU-25, sobre SQLite real, en una base nueva por escenario.
type sistemaDashboardReal struct {
	crearProyecto *app.CrearProyecto
	crearElemento *app.CrearElementoBacklog
	completar     *app.CompletarHistoria
	consultar     *app.ConsultarDashboard
	db            *sql.DB
	dir           string
	hoy           time.Time
	siguiente     map[int64]int // próximo número de sprint de cada proyecto
	elementos     int           // para que cada título sea único (RN-05.3)
}

func nuevoSistemaDashboard() (sistemaDashboard, error) {
	dir, err := os.MkdirTemp("", "bdd-hu25-")
	if err != nil {
		return nil, err
	}
	db, err := persistencia.AbrirDB(filepath.Join(dir, "bdd.db"))
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	proyectos := persistencia.NuevoRepositorioProyectos(db)
	elementos := persistencia.NuevoRepositorioBacklog(db)
	s := &sistemaDashboardReal{
		crearProyecto: app.NuevoCrearProyecto(proyectos),
		crearElemento: app.NuevoCrearElementoBacklog(elementos),
		completar:     app.NuevoCompletarHistoria(proyectos, elementos),
		db:            db,
		dir:           dir,
		hoy:           time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
		siguiente:     map[int64]int{},
	}
	// El reloj lee s.hoy en cada consulta, así el step "hoy es" lo puede fijar después de crearlo.
	s.consultar = app.NuevoConsultarDashboard(proyectos, persistencia.NuevoRepositorioSprints(db), elementos,
		persistencia.NuevoRepositorioDefectos(db), persistencia.NuevoRepositorioEsfuerzo(db), func() time.Time { return s.hoy })
	return s, nil
}

func (s *sistemaDashboardReal) CrearProyecto(nombre string) (int64, error) {
	p, err := s.crearProyecto.Ejecutar(context.Background(), app.DatosNuevoProyecto{
		Nombre: nombre, FechaInicio: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	return p.ID, err
}

func (s *sistemaDashboardReal) FijarHoy(hoy time.Time) { s.hoy = hoy }

// CrearElemento da de alta el elemento con HU-05 y le fija estado, story points y sprint directo en la
// base: los escenarios de HU-25 miran el conteo, no cómo se llega a cada estado.
func (s *sistemaDashboardReal) CrearElemento(proyectoID int64, estado domain.EstadoElemento, puntos *int, sprintID int64) (int64, error) {
	s.elementos++
	e, err := s.crearElemento.Ejecutar(context.Background(), app.DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      fmt.Sprintf("Elemento del dashboard %d", s.elementos),
		Descripcion: "Elemento de prueba",
		Prioridad:   domain.PrioridadAlta,
		Criterios:   []string{"Sirve para el dashboard"},
	})
	if err != nil {
		return 0, err
	}
	var sp any
	if puntos != nil {
		sp = *puntos
	}
	_, err = s.db.Exec(`UPDATE backlog_elementos SET estado = ?, story_points = ?, sprint_id = ? WHERE id = ?`,
		string(estado), sp, nuloSiCero(sprintID), e.ID)
	return e.ID, err
}

// CrearSprint inserta el sprint con el número siguiente del proyecto, sus fechas y su resumen de
// cierre, sin pasar por HU-08 ni HU-11 (igual criterio que HU-22 y HU-24).
func (s *sistemaDashboardReal) CrearSprint(proyectoID int64, estado domain.EstadoSprint, inicio, fin time.Time, planificados, completados int) (int64, error) {
	numero := s.siguiente[proyectoID] + 1
	s.siguiente[proyectoID] = numero
	r, err := s.db.Exec(
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, objetivo, estado, puntos_planificados, puntos_completados)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, numero, inicio.Format(time.DateOnly), fin.Format(time.DateOnly), "Objetivo de prueba",
		string(estado), planificados, completados)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

// CrearDefecto inserta el defecto directo en la base, sin pasar por HU-20 ni HU-21.
func (s *sistemaDashboardReal) CrearDefecto(proyectoID, elementoID int64, estado domain.EstadoDefecto, deteccionID, resolucionID int64) error {
	_, err := s.db.Exec(
		`INSERT INTO defectos (proyecto_id, elemento_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id)
		 VALUES (?, ?, 'Defecto de prueba', 'Alta', ?, ?, ?)`,
		proyectoID, elementoID, string(estado), nuloSiCero(deteccionID), nuloSiCero(resolucionID))
	return err
}

// CrearElementoConHoras guarda las horas en medias horas, como la persistencia de HU-18 y HU-19.
func (s *sistemaDashboardReal) CrearElementoConHoras(proyectoID int64, estimadas *float64, reales float64) error {
	eid, err := s.CrearElemento(proyectoID, domain.EstadoPendiente, nil, 0)
	if err != nil {
		return err
	}
	if estimadas != nil {
		if _, err := s.db.Exec(`UPDATE backlog_elementos SET medias_horas_estimadas = ? WHERE id = ?`, int(*estimadas*2), eid); err != nil {
			return err
		}
	}
	_, err = s.db.Exec(
		`INSERT INTO registros_esfuerzo (proyecto_id, elemento_id, integrante, fecha, actividad, medias_horas)
		 VALUES (?, ?, 'Ana', ?, 'Desarrollo', ?)`,
		proyectoID, eid, s.hoy.Format(time.DateOnly), int(reales*2))
	return err
}

func (s *sistemaDashboardReal) CompletarHistoria(proyectoID, elementoID int64) error {
	return s.completar.Ejecutar(context.Background(), app.DatosCompletar{ProyectoID: proyectoID, ElementoID: elementoID})
}

func (s *sistemaDashboardReal) ConsultarDashboard(proyectoID int64) (dashboardVista, error) {
	d, err := s.consultar.Ejecutar(context.Background(), proyectoID)
	if err != nil {
		return dashboardVista{}, err
	}
	v := dashboardVista{
		Pendientes: d.Backlog.Pendientes, EnProgreso: d.Backlog.EnProgreso, Hechos: d.Backlog.Hechos, Total: d.Backlog.Total,
		HaySprintEnCurso: d.HaySprintEnCurso, NumeroEnCurso: d.SprintEnCurso.Sprint.Numero,
		PlanificadosEnCurso: d.SprintEnCurso.PuntosPlanificados, HechosEnCurso: d.SprintEnCurso.PuntosCompletados,
		Velocidad: d.Velocidad.Velocidad, HayVelocidad: d.Velocidad.HayDatos,
		HorasEstimadas: d.Horas.HorasEstimadas, HorasReales: d.Horas.HorasReales, Desviacion: d.Horas.DesviacionPorcentaje,
	}
	for _, c := range d.SprintsCerrados {
		v.Cerrados = append(v.Cerrados, cerradoVista{
			Numero: c.Sprint.Numero, Planificados: c.Sprint.PuntosPlanificados, Completados: c.Sprint.PuntosCompletados,
			Porcentaje: c.PorcentajeCompletado, HayPorcentaje: c.HayPorcentaje,
			Detectados: c.DefectosDetectados, Resueltos: c.DefectosResueltos,
		})
	}
	return v, nil
}

// Cerrar cierra la base del escenario y borra su carpeta temporal.
func (s *sistemaDashboardReal) Cerrar() error {
	return errors.Join(s.db.Close(), os.RemoveAll(s.dir))
}
