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
	// FijarEstadoProyecto deja un proyecto preparado en otro estado, sin pasar por HU-02.
	FijarEstadoProyecto(proyectoID int64, estado domain.Estado) error
	// CrearSprint prepara un sprint con su número, fechas y estado, sin pasar por HU-08: es
	// preparación del escenario.
	CrearSprint(proyectoID int64, numero int, inicio, fin time.Time, estado domain.EstadoSprint) (int64, error)
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

	// When
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)"$`, w.cierroSprint)
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)" desde el proyecto "([^"]*)"$`, w.cierroSprintDesdeProyecto)
	ctx.Step(`^cierro el sprint (\d+) de "([^"]*)" desde un proyecto inexistente$`, w.cierroSprintDesdeProyectoInexistente)
	ctx.Step(`^cierro un sprint inexistente de "([^"]*)"$`, w.cierroSprintInexistente)

	// Then
	ctx.Step(`^el cierre del sprint termina sin errores$`, w.cierreSinErrores)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" está (Abierto|Cerrado)$`, w.sprintEsta)
	ctx.Step(`^el sprint (\d+) de "([^"]*)" conserva sus datos$`, w.sprintConservaDatos)
	ctx.Step(`^se rechaza el cierre del sprint con el error "([^"]*)"$`, w.seRechazaCierre)
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

// --- When ---

func (w *hu11World) cierroSprint(numero int, proyecto string) error {
	return w.cierroSprintDesdeProyecto(numero, proyecto, proyecto)
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

// --- Auxiliares ---

// hu11ClaveSprint identifica un sprint del escenario por su proyecto y su número.
func hu11ClaveSprint(proyecto string, numero int) string {
	return fmt.Sprintf("%s#%d", proyecto, numero)
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
