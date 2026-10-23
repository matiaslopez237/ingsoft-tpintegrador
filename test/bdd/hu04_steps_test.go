package bdd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/cucumber/godog"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sistemaConsulta es la puerta de entrada de los escenarios al núcleo.
// La implementación real está en hu04_sistema_test.go.
type sistemaConsulta interface {
	CrearProyecto(nombre string, inicio time.Time, fin *time.Time) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre, rol string) error
	AgregarRol(proyectoID int64, nombre, rol string) error
	Consultar(proyectoID int64) (domain.Proyecto, error)
	Cerrar() error
}

type consultarProyectoWorld struct {
	sistema    sistemaConsulta
	porNombre  map[string]int64
	proyectoID int64 // el último proyecto preparado: al que se suman los integrantes
	consultado domain.Proyecto
	ultimoErr  error
}

func registrarStepsHU04(ctx *godog.ScenarioContext) {
	w := &consultarProyectoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaConsulta()
		if err != nil {
			return c, err
		}
		*w = consultarProyectoWorld{sistema: sistema, porNombre: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^existe un proyecto "([^"]*)" con inicio "([^"]*)" y fin "([^"]*)"$`, w.existeConInicioYFin)
	ctx.Step(`^existe un proyecto "([^"]*)" con inicio "([^"]*)" y sin fecha de fin$`, w.existeSinFin)
	ctx.Step(`^existe un proyecto "([^"]*)" sin integrantes$`, w.existeSinIntegrantes)
	ctx.Step(`^"([^"]*)" se sumó al proyecto con el rol "([^"]*)"$`, w.seSumoConRol)
	ctx.Step(`^"([^"]*)" también sumó el rol "([^"]*)"$`, w.tambienSumoRol)

	ctx.Step(`^consulto el proyecto "([^"]*)"$`, w.consultoPorNombre)
	ctx.Step(`^consulto el proyecto con identificador (\d+)$`, w.consultoPorID)

	ctx.Step(`^veo el nombre "([^"]*)"$`, w.veoNombre)
	ctx.Step(`^veo el estado "([^"]*)"$`, w.veoEstado)
	ctx.Step(`^veo la fecha de inicio "([^"]*)" y la fecha de fin "([^"]*)"$`, w.veoFechas)
	ctx.Step(`^veo la fecha de inicio "([^"]*)" y ninguna fecha de fin$`, w.veoInicioSinFin)
	ctx.Step(`^veo los integrantes "([^"]*)" y "([^"]*)" en ese orden$`, w.veoIntegrantesEnOrden)
	ctx.Step(`^veo que "([^"]*)" tiene solo el rol "([^"]*)"$`, w.veoSoloElRol)
	ctx.Step(`^veo que "([^"]*)" tiene los roles "([^"]*)" y "([^"]*)" en ese orden$`, w.veoRolesEnOrden)
	ctx.Step(`^veo al integrante "([^"]*)" una sola vez$`, w.veoUnaSolaVez)
	ctx.Step(`^veo que el proyecto no tiene integrantes$`, w.veoSinIntegrantes)
	ctx.Step(`^se rechaza la consulta con el error "([^"]*)"$`, w.seRechazaConError)
}

// --- Given ---

func (w *consultarProyectoWorld) existeConInicioYFin(nombre, inicio, fin string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	return w.preparar(nombre, fi, &ff)
}

func (w *consultarProyectoWorld) existeSinFin(nombre, inicio string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	return w.preparar(nombre, fi, nil)
}

// existeSinIntegrantes usa fechas fijas: los escenarios de integrantes no las miran.
func (w *consultarProyectoWorld) existeSinIntegrantes(nombre string) error {
	fi, err := parsearFecha("2026-03-01")
	if err != nil {
		return err
	}
	return w.preparar(nombre, fi, nil)
}

func (w *consultarProyectoWorld) preparar(nombre string, inicio time.Time, fin *time.Time) error {
	id, err := w.sistema.CrearProyecto(nombre, inicio, fin)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", nombre, err)
	}
	w.porNombre[nombre] = id
	w.proyectoID = id
	return nil
}

func (w *consultarProyectoWorld) seSumoConRol(nombre, rol string) error {
	if err := w.sistema.AgregarIntegrante(w.proyectoID, nombre, rol); err != nil {
		return fmt.Errorf("no se pudo preparar el integrante %q: %w", nombre, err)
	}
	return nil
}

func (w *consultarProyectoWorld) tambienSumoRol(nombre, rol string) error {
	if err := w.sistema.AgregarRol(w.proyectoID, nombre, rol); err != nil {
		return fmt.Errorf("no se pudo preparar el rol %q de %q: %w", rol, nombre, err)
	}
	return nil
}

// --- When ---

func (w *consultarProyectoWorld) consultoPorNombre(nombre string) error {
	id, ok := w.porNombre[nombre]
	if !ok {
		return fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return w.consultoPorID(id)
}

func (w *consultarProyectoWorld) consultoPorID(id int64) error {
	w.consultado, w.ultimoErr = w.sistema.Consultar(id)
	return nil
}

// --- Then ---

// consultaOK falla el step si la consulta devolvió error: los Then que miran datos lo necesitan.
func (w *consultarProyectoWorld) consultaOK() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba consultar el proyecto, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *consultarProyectoWorld) veoNombre(nombre string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	if w.consultado.Nombre != nombre {
		return fmt.Errorf("nombre: se esperaba %q, se obtuvo %q", nombre, w.consultado.Nombre)
	}
	return nil
}

func (w *consultarProyectoWorld) veoEstado(estado string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	if string(w.consultado.Estado) != estado {
		return fmt.Errorf("estado: se esperaba %q, se obtuvo %q", estado, w.consultado.Estado)
	}
	return nil
}

func (w *consultarProyectoWorld) veoFechas(inicio, fin string) error {
	if err := w.veoInicio(inicio); err != nil {
		return err
	}
	if w.consultado.FechaFin == nil {
		return fmt.Errorf("fecha de fin: se esperaba %q, el proyecto no tiene", fin)
	}
	if got := w.consultado.FechaFin.Format(time.DateOnly); got != fin {
		return fmt.Errorf("fecha de fin: se esperaba %q, se obtuvo %q", fin, got)
	}
	return nil
}

func (w *consultarProyectoWorld) veoInicioSinFin(inicio string) error {
	if err := w.veoInicio(inicio); err != nil {
		return err
	}
	if w.consultado.FechaFin != nil {
		return fmt.Errorf("se esperaba un proyecto sin fecha de fin, tiene %s", w.consultado.FechaFin.Format(time.DateOnly))
	}
	return nil
}

func (w *consultarProyectoWorld) veoInicio(inicio string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	if got := w.consultado.FechaInicio.Format(time.DateOnly); got != inicio {
		return fmt.Errorf("fecha de inicio: se esperaba %q, se obtuvo %q", inicio, got)
	}
	return nil
}

// veoIntegrantesEnOrden compara en orden, sin ordenar: el orden de alta es parte de RN-04.2.
func (w *consultarProyectoWorld) veoIntegrantesEnOrden(a, b string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	var nombres []string
	for _, i := range w.consultado.Integrantes {
		nombres = append(nombres, i.Nombre)
	}
	if want := []string{a, b}; !slices.Equal(nombres, want) {
		return fmt.Errorf("integrantes: se esperaba %v, se obtuvo %v", want, nombres)
	}
	return nil
}

func (w *consultarProyectoWorld) veoSoloElRol(nombre, rol string) error {
	return w.rolesSonExactamente(nombre, rol)
}

func (w *consultarProyectoWorld) veoRolesEnOrden(nombre, a, b string) error {
	return w.rolesSonExactamente(nombre, a, b)
}

func (w *consultarProyectoWorld) rolesSonExactamente(nombre string, esperados ...string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	for _, i := range w.consultado.Integrantes {
		if i.Nombre == nombre {
			if !slices.Equal(i.Roles, esperados) {
				return fmt.Errorf("roles de %q: se esperaba %v, se obtuvo %v", nombre, esperados, i.Roles)
			}
			return nil
		}
	}
	return fmt.Errorf("%q no figura entre los integrantes consultados", nombre)
}

func (w *consultarProyectoWorld) veoUnaSolaVez(nombre string) error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	veces := 0
	for _, i := range w.consultado.Integrantes {
		if i.Nombre == nombre {
			veces++
		}
	}
	if veces != 1 {
		return fmt.Errorf("%q: se esperaba que figure 1 vez, figura %d", nombre, veces)
	}
	return nil
}

func (w *consultarProyectoWorld) veoSinIntegrantes() error {
	if err := w.consultaOK(); err != nil {
		return err
	}
	if n := len(w.consultado.Integrantes); n != 0 {
		return fmt.Errorf("se esperaba un proyecto sin integrantes, tiene %d", n)
	}
	return nil
}

func (w *consultarProyectoWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU04()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero la consulta funcionó", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}
