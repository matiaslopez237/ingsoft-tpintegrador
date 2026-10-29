package bdd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu11Sistema es la puerta de entrada de los escenarios de HU-11 al núcleo. La implementa
// hu11SistemaReal, sobre SQLite (hu11_sistema_test.go). Todos los identificadores de paquete llevan
// el prefijo hu11 para no chocar con los de HU-09 y HU-10 al mergear.
type hu11Sistema interface {
	CrearProyecto(nombre string) (int64, error)
	// FijarEstadoProyecto deja un proyecto preparado en otro estado (con HU-02).
	FijarEstadoProyecto(proyectoID int64, estado domain.Estado) error
	// CrearSprint prepara un sprint con su número, fechas y estado, sin pasar por HU-08: es
	// preparación del escenario.
	CrearSprint(proyectoID int64, numero int, inicio, fin time.Time, estado domain.EstadoSprint) (int64, error)
	// CrearHistoriaAsignada prepara una historia sin terminar del proyecto, con esos story points
	// (nil: sin estimar), asignada al sprint (HU-05 y HU-09). Es para CA-11.6 y CA-11.7.
	CrearHistoriaAsignada(proyectoID, sprintID int64, titulo string, sp *int) (int64, error)
	// CompletarHistoria deja la historia Hecho con el caso de uso de HU-10.
	CompletarHistoria(proyectoID, historiaID int64) error
	// AsignarHistoria asigna la historia a otro sprint con el caso de uso de HU-09.
	AsignarHistoria(proyectoID, historiaID, sprintID int64) error
	// EstadoDeHistoria devuelve el estado de la historia.
	EstadoDeHistoria(historiaID int64) (domain.EstadoElemento, error)
	// SprintDeHistoria devuelve el sprint al que está asignada la historia, si tiene.
	SprintDeHistoria(historiaID int64) (sprintID int64, asignada bool, err error)
	EjecutarCierre(proyectoID, sprintID int64) error
	Sprint(sprintID int64) (domain.Sprint, error)
	Cerrar() error
}

// hu11Errores traduce la descripción del escenario al error definido en la spec (sección 7).
// ErrSprintNoEncontrado es el que HU-09 dejó en el dominio (spec HU-11, sección 5).
func hu11Errores() map[string]error {
	return map[string]error{
		"proyecto no encontrado": domain.ErrProyectoNoEncontrado,
		"sprint no encontrado":   domain.ErrSprintNoEncontrado,
		"sprint ya cerrado":      domain.ErrSprintYaCerrado,
	}
}

type hu11World struct {
	sistema    hu11Sistema
	proyectos  map[string]int64
	sprints    map[string]int64         // clave: proyecto y número (ver hu11ClaveSprint)
	preparados map[string]domain.Sprint // el sprint tal como se preparó, para ver que no cambió
	historias  map[string]int64         // clave: proyecto y título (ver hu11ClaveHistoria)
	ultimoErr  error
}

func registrarStepsHU11(ctx *godog.ScenarioContext) {
	w := &hu11World{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = hu11World{} // si falla la base, el After no cierra el sistema del escenario anterior
		sistema, err := hu11NuevoSistema()
		if err != nil {
			return c, err
		}
		*w = hu11World{
			sistema:    sistema,
			proyectos:  map[string]int64{},
			sprints:    map[string]int64{},
			preparados: map[string]domain.Sprint{},
			historias:  map[string]int64{},
		}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^el proyecto "([^"]*)" con el sprint (\d+) (abierto|cerrado) del "([^"]*)" al "([^"]*)"$`, w.proyectoConSprint)
	ctx.Step(`^cerré el sprint (\d+) de "([^"]*)"$`, w.cerreSprint)
	ctx.Step(`^el proyecto "([^"]*)" está "([^"]*)"$`, w.proyectoEsta)
	ctx.Step(`^la historia "([^"]*)" de "([^"]*)" está asignada al sprint (\d+)$`, w.historiaAsignada)
	ctx.Step(`^la historia "([^"]*)" de "([^"]*)", (de \d+ puntos|sin estimar), está asignada al sprint (\d+)$`, w.historiaConPuntosAsignada)
	ctx.Step(`^la historia "([^"]*)" de "([^"]*)" está completada$`, w.historiaCompletada)

	// When
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)"$`, w.cierroSprint)
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)" desde el proyecto "([^"]*)"$`, w.cierroSprintDesdeProyecto)
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)" desde un proyecto inexistente$`, w.cierroSprintDesdeProyectoInexistente)
	ctx.Step(`^cierro un sprint inexistente de "([^"]*)"$`, w.cierroSprintInexistente)
	ctx.Step(`^asigno la historia "([^"]*)" de "([^"]*)" al sprint (\d+)$`, w.asignoHistoria)

	// Then
	ctx.Step(`^el cierre del sprint termina sin errores$`, w.cierreSinErrores)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" está (Abierto|Cerrado)$`, w.sprintEsta)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" conserva sus datos$`, w.sprintConservaDatos)
	ctx.Step(`^se rechaza el cierre del sprint con el error "([^"]*)"$`, w.seRechazaCierre)
	ctx.Step(`^la historia "([^"]*)" de "([^"]*)" sigue asignada al sprint (\d+)$`, w.historiaSigueAsignada)
	ctx.Step(`^la historia "([^"]*)" de "([^"]*)" volvió al backlog sin sprint y sigue "([^"]*)"$`, w.historiaVolvioAlBacklog)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" registra (\d+) puntos planificados y (\d+) completados$`, w.sprintRegistraPuntos)
}

// --- Given ---

func (w *hu11World) proyectoConSprint(proyecto string, numero int, estado, inicio, fin string) error {
	pid, ok := w.proyectos[proyecto]
	if !ok {
		var err error
		if pid, err = w.sistema.CrearProyecto(proyecto); err != nil {
			return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
		}
		w.proyectos[proyecto] = pid
	}
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	e := domain.SprintAbierto
	if estado == "cerrado" {
		e = domain.SprintCerrado
	}
	sid, err := w.sistema.CrearSprint(pid, numero, fi, ff, e)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el sprint %d de %q: %w", numero, proyecto, err)
	}
	s, err := w.sistema.Sprint(sid)
	if err != nil {
		return fmt.Errorf("no se pudo leer el sprint %d de %q: %w", numero, proyecto, err)
	}
	clave := hu11ClaveSprint(proyecto, numero)
	w.sprints[clave] = sid
	w.preparados[clave] = s
	return nil
}

func (w *hu11World) cerreSprint(numero int, proyecto string) error {
	pid, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.EjecutarCierre(pid, sid); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint %d de %q cerrado: %w", numero, proyecto, err)
	}
	return nil
}

func (w *hu11World) proyectoEsta(proyecto, estado string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	var e domain.Estado
	switch estado {
	case "En pausa":
		e = domain.EstadoEnPausa
	case "Finalizado":
		e = domain.EstadoFinalizado
	case "Postergado":
		e = domain.EstadoPostergado
	case "Cancelado":
		e = domain.EstadoCancelado
	default:
		return fmt.Errorf("estado de proyecto desconocido en el escenario: %q", estado)
	}
	if err := w.sistema.FijarEstadoProyecto(pid, e); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", proyecto, estado, err)
	}
	return nil
}

func (w *hu11World) historiaAsignada(titulo, proyecto string, numero int) error {
	return w.prepararHistoria(titulo, proyecto, numero, nil)
}

// historiaConPuntosAsignada recibe "de N puntos" o "sin estimar".
func (w *hu11World) historiaConPuntosAsignada(titulo, proyecto, puntos string, numero int) error {
	if puntos == "sin estimar" {
		return w.prepararHistoria(titulo, proyecto, numero, nil)
	}
	var sp int
	if _, err := fmt.Sscanf(puntos, "de %d puntos", &sp); err != nil {
		return fmt.Errorf("story points mal escritos en el escenario: %q", puntos)
	}
	return w.prepararHistoria(titulo, proyecto, numero, &sp)
}

func (w *hu11World) historiaCompletada(titulo, proyecto string) error {
	pid, hid, err := w.idsHistoria(titulo, proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CompletarHistoria(pid, hid); err != nil {
		return fmt.Errorf("no se pudo completar la historia %q de %q: %w", titulo, proyecto, err)
	}
	return nil
}

func (w *hu11World) prepararHistoria(titulo, proyecto string, numero int, sp *int) error {
	pid, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	hid, err := w.sistema.CrearHistoriaAsignada(pid, sid, titulo, sp)
	if err != nil {
		return fmt.Errorf("no se pudo preparar la historia %q asignada al sprint %d de %q: %w", titulo, numero, proyecto, err)
	}
	w.historias[hu11ClaveHistoria(proyecto, titulo)] = hid
	return nil
}

// --- When ---

func (w *hu11World) cierroSprint(numero int, proyecto string) error {
	return w.cierroSprintDesdeProyecto(numero, proyecto, proyecto)
}

// asignoHistoria replanifica con HU-09 una historia que volvió al backlog. Si falla, falla el paso:
// lo que se verifica después es el resumen del sprint cerrado, no la asignación.
func (w *hu11World) asignoHistoria(titulo, proyecto string, numero int) error {
	pid, hid, err := w.idsHistoria(titulo, proyecto)
	if err != nil {
		return err
	}
	_, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.AsignarHistoria(pid, hid, sid); err != nil {
		return fmt.Errorf("no se pudo asignar la historia %q al sprint %d de %q: %w", titulo, numero, proyecto, err)
	}
	return nil
}

func (w *hu11World) cierroSprintDesdeProyecto(numero int, proyecto, ruta string) error {
	_, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	pid, err := w.idProyecto(ruta)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.EjecutarCierre(pid, sid)
	return nil
}

func (w *hu11World) cierroSprintDesdeProyectoInexistente(numero int, proyecto string) error {
	_, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.EjecutarCierre(idInexistente, sid)
	return nil
}

func (w *hu11World) cierroSprintInexistente(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.ultimoErr = w.sistema.EjecutarCierre(pid, idInexistente)
	return nil
}

// --- Then ---

func (w *hu11World) cierreSinErrores() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba cerrar el sprint, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *hu11World) sprintEsta(numero int, proyecto, estado string) error {
	s, err := w.sprint(numero, proyecto)
	if err != nil {
		return err
	}
	if s.Estado != domain.EstadoSprint(estado) {
		return fmt.Errorf("estado del sprint %d de %q: se esperaba %q, es %q", numero, proyecto, estado, s.Estado)
	}
	return nil
}

// sprintConservaDatos compara todo menos el estado con el sprint tal como se preparó.
func (w *hu11World) sprintConservaDatos(numero int, proyecto string) error {
	s, err := w.sprint(numero, proyecto)
	if err != nil {
		return err
	}
	antes := w.preparados[hu11ClaveSprint(proyecto, numero)]
	s.Estado, antes.Estado = "", ""
	if s != antes {
		return fmt.Errorf("el sprint %d de %q cambió: era %+v, es %+v", numero, proyecto, antes, s)
	}
	return nil
}

func (w *hu11World) seRechazaCierre(descripcion string) error {
	esperado, ok := hu11Errores()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el sprint se cerró", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *hu11World) historiaSigueAsignada(titulo, proyecto string, numero int) error {
	_, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return err
	}
	hid, ok := w.historias[hu11ClaveHistoria(proyecto, titulo)]
	if !ok {
		return fmt.Errorf("la historia %q de %q no se preparó en el escenario", titulo, proyecto)
	}
	got, asignada, err := w.sistema.SprintDeHistoria(hid)
	if err != nil {
		return fmt.Errorf("no se pudo leer la historia %q de %q: %w", titulo, proyecto, err)
	}
	if !asignada || got != sid {
		return fmt.Errorf("la historia %q de %q: se esperaba asignada al sprint %d (ID %d), está en %d (asignada: %t)",
			titulo, proyecto, numero, sid, got, asignada)
	}
	return nil
}

func (w *hu11World) historiaVolvioAlBacklog(titulo, proyecto, estado string) error {
	_, hid, err := w.idsHistoria(titulo, proyecto)
	if err != nil {
		return err
	}
	got, asignada, err := w.sistema.SprintDeHistoria(hid)
	if err != nil {
		return fmt.Errorf("no se pudo leer la historia %q de %q: %w", titulo, proyecto, err)
	}
	if asignada {
		return fmt.Errorf("la historia %q de %q: se esperaba sin sprint, sigue asignada al sprint con ID %d", titulo, proyecto, got)
	}
	actual, err := w.sistema.EstadoDeHistoria(hid)
	if err != nil {
		return fmt.Errorf("no se pudo leer la historia %q de %q: %w", titulo, proyecto, err)
	}
	if string(actual) != estado {
		return fmt.Errorf("la historia %q de %q: estado %q, se esperaba %q", titulo, proyecto, actual, estado)
	}
	return nil
}

func (w *hu11World) sprintRegistraPuntos(numero int, proyecto string, planificados, completados int) error {
	s, err := w.sprint(numero, proyecto)
	if err != nil {
		return err
	}
	if s.PuntosPlanificados != planificados || s.PuntosCompletados != completados {
		return fmt.Errorf("sprint %d de %q: %d planificados y %d completados, se esperaba %d y %d",
			numero, proyecto, s.PuntosPlanificados, s.PuntosCompletados, planificados, completados)
	}
	return nil
}

// --- Auxiliares ---

// idsHistoria devuelve el ID del proyecto y el de una historia que preparó el escenario.
func (w *hu11World) idsHistoria(titulo, proyecto string) (proyectoID, historiaID int64, err error) {
	if proyectoID, err = w.idProyecto(proyecto); err != nil {
		return 0, 0, err
	}
	historiaID, ok := w.historias[hu11ClaveHistoria(proyecto, titulo)]
	if !ok {
		return 0, 0, fmt.Errorf("la historia %q de %q no se preparó en el escenario", titulo, proyecto)
	}
	return proyectoID, historiaID, nil
}

// hu11ClaveSprint identifica un sprint del escenario por su proyecto y su número.
func hu11ClaveSprint(proyecto string, numero int) string {
	return fmt.Sprintf("%s#%d", proyecto, numero)
}

// hu11ClaveHistoria identifica una historia del escenario por su proyecto y su título.
func hu11ClaveHistoria(proyecto, titulo string) string {
	return proyecto + "#" + titulo
}

func (w *hu11World) sprint(numero int, proyecto string) (domain.Sprint, error) {
	_, sid, err := w.ids(numero, proyecto)
	if err != nil {
		return domain.Sprint{}, err
	}
	s, err := w.sistema.Sprint(sid)
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("no se pudo leer el sprint %d de %q: %w", numero, proyecto, err)
	}
	return s, nil
}

func (w *hu11World) ids(numero int, proyecto string) (proyectoID, sprintID int64, err error) {
	if proyectoID, err = w.idProyecto(proyecto); err != nil {
		return 0, 0, err
	}
	sprintID, ok := w.sprints[hu11ClaveSprint(proyecto, numero)]
	if !ok {
		return 0, 0, fmt.Errorf("el sprint %d de %q no se preparó en el escenario", numero, proyecto)
	}
	return proyectoID, sprintID, nil
}

func (w *hu11World) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}
