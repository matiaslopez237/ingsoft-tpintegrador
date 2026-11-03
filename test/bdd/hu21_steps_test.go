package bdd

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"
)

// sistemaEstadoDefecto es la puerta de entrada de los escenarios de HU-21 al núcleo.
// La implementación real está en hu21_sistema_test.go.
type sistemaEstadoDefecto interface {
	CrearProyecto(nombre string, inicio time.Time) (int64, error)
	CrearElemento(proyectoID int64, titulo string) (int64, error)
	CrearSprint(proyectoID int64, inicio, fin time.Time) error
	CrearDefecto(d datosDefecto) (defectoVisto, error)
	IniciarDefecto(d datosActualizarDefecto) (defectoVisto, error)
	ResolverDefecto(d datosActualizarDefecto) (defectoVisto, error)
	CerrarDefecto(d datosActualizarDefecto) (defectoVisto, error)
	Cerrar() error
}

// datosActualizarDefecto son las entradas de HU-21 (spec, sección 2).
type datosActualizarDefecto struct {
	ProyectoID int64
	DefectoID  int64
}

type actualizarEstadoDefectoWorld struct {
	sistema    sistemaEstadoDefecto
	hoy        time.Time // lo lee el reloj falso del sistema
	proyectos  map[string]int64
	elementos  map[string]int64
	proyectoID int64 // el proyecto del Background: ahí se preparó el defecto si no se indica otro
	defectoID  int64 // el defecto del Background
	ultimo     defectoVisto
	ultimoErr  error
}

func registrarStepsHU21(ctx *godog.ScenarioContext) {
	w := &actualizarEstadoDefectoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = actualizarEstadoDefectoWorld{proyectos: map[string]int64{}, elementos: map[string]int64{}}
		sistema, err := nuevoSistemaEstadoDefecto(w.reloj)
		if err != nil {
			return c, err
		}
		w.sistema = sistema
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil { // el Before falló: no hay base que cerrar y godog ya informa su error
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^la fecha de referencia es "([^"]*)"$`, w.hoyEs)
	ctx.Step(`^el proyecto "([^"]*)" iniciado el "([^"]*)" ya tiene el elemento "([^"]*)"$`, w.proyectoConElemento)
	ctx.Step(`^el proyecto "([^"]*)" cuenta con un sprint vigente entre "([^"]*)" y "([^"]*)"$`, w.proyectoTieneSprint)
	ctx.Step(`^hay un defecto "([^"]*)" con severidad "([^"]*)" en "([^"]*)"$`, w.hayUnDefecto)

	ctx.Step(`^inicié el defecto$`, w.inicieElDefecto)
	ctx.Step(`^resolví el defecto$`, w.resolviElDefecto)
	ctx.Step(`^cerré el defecto$`, w.cerreElDefecto)
	ctx.Step(`^inicio el defecto$`, w.inicioElDefecto)
	ctx.Step(`^resuelvo el defecto$`, w.resuelvoElDefecto)
	ctx.Step(`^cierro el defecto$`, w.cierroElDefecto)
	ctx.Step(`^resuelvo el defecto en un proyecto inexistente$`, w.resuelvoEnProyectoInexistente)
	ctx.Step(`^resuelvo un defecto inexistente$`, w.resuelvoDefectoInexistente)
	ctx.Step(`^resuelvo el defecto de "([^"]*)" indicando el proyecto "([^"]*)"$`, w.resuelvoDefectoDeOtroProyecto)

	ctx.Step(`^el defecto se actualiza correctamente$`, w.seActualizaCorrectamente)
	ctx.Step(`^el defecto queda en estado "([^"]*)"$`, w.estadoCoincide)
	ctx.Step(`^el defecto queda en estado "([^"]*)" con el Sprint de resolución asignado$`, w.defectoQuedaConSprint)
	ctx.Step(`^el defecto queda en estado "([^"]*)" sin Sprint de resolución$`, w.defectoQuedaSinSprint)
	ctx.Step(`^se rechaza la actualización con el error "([^"]*)"$`, w.seRechazaConError)
}

// reloj es el reloj falso que recibe el caso de uso: devuelve el "hoy" del escenario.
func (w *actualizarEstadoDefectoWorld) reloj() time.Time {
	return w.hoy
}

// --- Given ---

// hoyEs fija "hoy" al mediodía, en la zona del servidor, igual que HU-18/HU-20 (RN-18.4, RN-20.5):
// así la comparación por día no depende de la hora.
func (w *actualizarEstadoDefectoWorld) hoyEs(fecha string) error {
	d, err := parsearFechaLocal(fecha)
	if err != nil {
		return err
	}
	w.hoy = d.Add(12 * time.Hour)
	return nil
}

func (w *actualizarEstadoDefectoWorld) proyectoConElemento(proyecto, inicio, elemento string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	id, err := w.sistema.CrearProyecto(proyecto, fi)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
	}
	elemID, err := w.sistema.CrearElemento(id, elemento)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el elemento %q: %w", elemento, err)
	}
	w.proyectos[proyecto] = id
	w.elementos[elemento] = elemID
	if w.proyectoID == 0 {
		w.proyectoID = id
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) proyectoTieneSprint(proyecto, inicio, fin string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	ff, err := parsearFecha(fin)
	if err != nil {
		return err
	}
	if err := w.sistema.CrearSprint(id, fi, ff); err != nil {
		return fmt.Errorf("no se pudo preparar el sprint del proyecto %q: %w", proyecto, err)
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) hayUnDefecto(descripcion, severidad, elemento string) error {
	elemID, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	d, err := w.sistema.CrearDefecto(datosDefecto{
		ProyectoID: w.proyectoID, ElementoID: elemID, Descripcion: descripcion, Severidad: severidad,
	})
	if err != nil {
		return fmt.Errorf("no se pudo preparar el defecto %q: %w", descripcion, err)
	}
	w.defectoID = d.ID
	return nil
}

// inicieElDefecto prepara un defecto ya En progreso: a diferencia de "inicio", si falla corta el
// escenario (es un paso de preparación, no la acción bajo prueba), igual que "cerré" en HU-11.
func (w *actualizarEstadoDefectoWorld) inicieElDefecto() error {
	w.ejecutarIniciar(w.datosActuales())
	return w.ultimoErr
}

// resolviElDefecto prepara un defecto ya Resuelto: a diferencia de "resuelvo", si falla corta el
// escenario (es un paso de preparación, no la acción bajo prueba), igual que "cerré" en HU-11.
func (w *actualizarEstadoDefectoWorld) resolviElDefecto() error {
	w.ejecutarResolver(w.datosActuales())
	return w.ultimoErr
}

func (w *actualizarEstadoDefectoWorld) cerreElDefecto() error {
	w.ejecutarCerrar(w.datosActuales())
	return w.ultimoErr
}

// --- When ---

func (w *actualizarEstadoDefectoWorld) inicioElDefecto() error {
	w.ejecutarIniciar(w.datosActuales())
	return nil
}

func (w *actualizarEstadoDefectoWorld) resuelvoElDefecto() error {
	w.ejecutarResolver(w.datosActuales())
	return nil
}

func (w *actualizarEstadoDefectoWorld) cierroElDefecto() error {
	w.ejecutarCerrar(w.datosActuales())
	return nil
}

func (w *actualizarEstadoDefectoWorld) resuelvoEnProyectoInexistente() error {
	w.ejecutarResolver(datosActualizarDefecto{ProyectoID: idInexistente, DefectoID: w.defectoID})
	return nil
}

func (w *actualizarEstadoDefectoWorld) resuelvoDefectoInexistente() error {
	w.ejecutarResolver(datosActualizarDefecto{ProyectoID: w.proyectoID, DefectoID: idInexistente})
	return nil
}

// resuelvoDefectoDeOtroProyecto pide resolver el defecto del Background (que pertenece a
// proyectoDefecto) pero indicando en la ruta proyectoPedido, un proyecto distinto (RN-21.6).
func (w *actualizarEstadoDefectoWorld) resuelvoDefectoDeOtroProyecto(proyectoDefecto, proyectoPedido string) error {
	idDefecto, err := w.idProyecto(proyectoDefecto)
	if err != nil {
		return err
	}
	if idDefecto != w.proyectoID {
		return fmt.Errorf("el defecto del Background pertenece al proyecto %d, no a %q", w.proyectoID, proyectoDefecto)
	}
	idPedido, err := w.idProyecto(proyectoPedido)
	if err != nil {
		return err
	}
	w.ejecutarResolver(datosActualizarDefecto{ProyectoID: idPedido, DefectoID: w.defectoID})
	return nil
}

func (w *actualizarEstadoDefectoWorld) ejecutarIniciar(d datosActualizarDefecto) {
	w.ultimo, w.ultimoErr = w.sistema.IniciarDefecto(d)
}

func (w *actualizarEstadoDefectoWorld) ejecutarResolver(d datosActualizarDefecto) {
	w.ultimo, w.ultimoErr = w.sistema.ResolverDefecto(d)
}

func (w *actualizarEstadoDefectoWorld) ejecutarCerrar(d datosActualizarDefecto) {
	w.ultimo, w.ultimoErr = w.sistema.CerrarDefecto(d)
}

func (w *actualizarEstadoDefectoWorld) datosActuales() datosActualizarDefecto {
	return datosActualizarDefecto{ProyectoID: w.proyectoID, DefectoID: w.defectoID}
}

// --- Then ---

func (w *actualizarEstadoDefectoWorld) seActualizaCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba actualizar el defecto, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) defectoQuedaConSprint(estado string) error {
	if err := w.estadoCoincide(estado); err != nil {
		return err
	}
	if w.ultimo.SprintResolucionID == nil {
		return fmt.Errorf("se esperaba un Sprint de resolución asignado, no se asignó ninguno")
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) defectoQuedaSinSprint(estado string) error {
	if err := w.estadoCoincide(estado); err != nil {
		return err
	}
	if w.ultimo.SprintResolucionID != nil {
		return fmt.Errorf("se esperaba sin Sprint de resolución, se asignó el %d", *w.ultimo.SprintResolucionID)
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) estadoCoincide(estado string) error {
	if err := w.seActualizaCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Estado != estado {
		return fmt.Errorf("se esperaba el estado %q, se obtuvo %q", estado, w.ultimo.Estado)
	}
	return nil
}

func (w *actualizarEstadoDefectoWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU21()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el defecto se actualizó", descripcion)
	}
	if !esperado.coincide(w.ultimoErr) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

// --- Auxiliares ---

func (w *actualizarEstadoDefectoWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *actualizarEstadoDefectoWorld) idElemento(titulo string) (int64, error) {
	id, ok := w.elementos[titulo]
	if !ok {
		return 0, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return id, nil
}
