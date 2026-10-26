package bdd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// sistemaEsfuerzo es la puerta de entrada de los escenarios de HU-18 al núcleo.
// La implementación real está en hu18_sistema_test.go.
type sistemaEsfuerzo interface {
	CrearProyecto(nombre string, inicio time.Time) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre string) error
	QuitarIntegrante(proyectoID int64, nombre string) error
	CambiarEstado(proyectoID int64, estado string) error
	CrearElemento(proyectoID int64, titulo string) (int64, error)
	RegistrarEsfuerzo(d datosEsfuerzo) (registroVisto, error)
	BuscarRegistro(id int64) (registroVisto, error)
	ContarRegistros(elementoID int64) (int, error)
	Cerrar() error
}

// datosEsfuerzo son las entradas de HU-18 (spec, sección 2). Fecha cero = "sin fecha".
type datosEsfuerzo struct {
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
}

// registroVisto es lo que los escenarios miran de un registro de esfuerzo (spec, sección 3).
type registroVisto struct {
	ID         int64
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      string // AAAA-MM-DD
	Actividad  string
	Horas      float64
}

// idInexistente no lo asigna SQLite en una base nueva por escenario.
const idInexistente int64 = 999999

type registrarEsfuerzoWorld struct {
	sistema     sistemaEsfuerzo
	hoy         time.Time // lo lee el reloj falso del sistema
	proyectos   map[string]int64
	elementos   map[string]int64
	proyectoID  int64 // el proyecto del Background: ahí se registra si no se indica otro
	registrados []registroVisto
	ultimo      registroVisto
	ultimoErr   error
}

func registrarStepsHU18(ctx *godog.ScenarioContext) {
	w := &registrarEsfuerzoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = registrarEsfuerzoWorld{proyectos: map[string]int64{}, elementos: map[string]int64{}}
		sistema, err := nuevoSistemaEsfuerzo(w.reloj)
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

	ctx.Step(`^hoy es "([^"]*)"$`, w.hoyEs)
	ctx.Step(`^el proyecto "([^"]*)" con inicio "([^"]*)" tiene la integrante "([^"]*)" y el elemento "([^"]*)"$`, w.proyectoConIntegranteYElemento)
	ctx.Step(`^"([^"]*)" no tiene esfuerzo registrado$`, w.noTieneEsfuerzo)
	ctx.Step(`^"([^"]*)" ya registró (-?[\d,]+) horas de "([^"]*)" en "([^"]*)" con fecha "([^"]*)"$`, w.yaRegistro)
	ctx.Step(`^se quita a "([^"]*)" del proyecto "([^"]*)"$`, w.seQuitaDelProyecto)
	ctx.Step(`^el proyecto "([^"]*)" pasó al estado "([^"]*)"$`, w.pasoAlEstado)

	ctx.Step(`^"([^"]*)" registra (-?[\d,]+) horas de "([^"]*)" en "([^"]*)" con fecha "([^"]*)"$`, w.registra)
	ctx.Step(`^"([^"]*)" registra (-?[\d,]+) horas de "([^"]*)" en "([^"]*)" sin fecha$`, w.registraSinFecha)
	ctx.Step(`^"([^"]*)" registra (-?[\d,]+) horas de "([^"]*)" en "([^"]*)" de un proyecto inexistente$`, w.registraEnProyectoInexistente)
	ctx.Step(`^"([^"]*)" registra (-?[\d,]+) horas de "([^"]*)" en un elemento inexistente con fecha "([^"]*)"$`, w.registraEnElementoInexistente)

	ctx.Step(`^el esfuerzo se registra correctamente$`, w.seRegistraCorrectamente)
	ctx.Step(`^el registro tiene un identificador asignado$`, w.tieneIdentificador)
	ctx.Step(`^el registro tiene la integrante "([^"]*)", la fecha "([^"]*)", la actividad "([^"]*)" y (-?[\d,]+) horas$`, w.registroTieneDatos)
	ctx.Step(`^se rechaza el esfuerzo con el error "([^"]*)"$`, w.seRechazaConError)
	ctx.Step(`^"([^"]*)" tiene (\d+) registros de esfuerzo$`, w.elementoTieneRegistros)
	ctx.Step(`^los dos registros tienen identificadores distintos$`, w.identificadoresDistintos)
	ctx.Step(`^cada registro se recupera por su identificador con sus datos$`, w.recuperablesPorID)
	ctx.Step(`^el registro de "([^"]*)" se sigue recuperando por su identificador$`, w.registroSeSigueRecuperando)
}

// reloj es el reloj falso que recibe el caso de uso: devuelve el "hoy" del escenario.
func (w *registrarEsfuerzoWorld) reloj() time.Time {
	return w.hoy
}

// --- Given ---

// hoyEs fija "hoy" al mediodía, en la zona del servidor (RN-18.4): así la comparación por día
// no depende de la hora.
func (w *registrarEsfuerzoWorld) hoyEs(fecha string) error {
	d, err := parsearFechaLocal(fecha)
	if err != nil {
		return err
	}
	w.hoy = d.Add(12 * time.Hour)
	return nil
}

func (w *registrarEsfuerzoWorld) proyectoConIntegranteYElemento(proyecto, inicio, integrante, elemento string) error {
	fi, err := parsearFecha(inicio)
	if err != nil {
		return err
	}
	id, err := w.sistema.CrearProyecto(proyecto, fi)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
	}
	if err := w.sistema.AgregarIntegrante(id, integrante); err != nil {
		return fmt.Errorf("no se pudo preparar el integrante %q: %w", integrante, err)
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

func (w *registrarEsfuerzoWorld) noTieneEsfuerzo(elemento string) error {
	return w.elementoTieneRegistros(elemento, 0)
}

func (w *registrarEsfuerzoWorld) yaRegistro(integrante, horas, actividad, elemento, fecha string) error {
	d, err := w.datos(integrante, horas, actividad, elemento)
	if err != nil {
		return err
	}
	if d.Fecha, err = parsearFechaLocal(fecha); err != nil {
		return err
	}
	r, err := w.sistema.RegistrarEsfuerzo(d)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el esfuerzo de %q: %w", integrante, err)
	}
	w.registrados = append(w.registrados, r)
	return nil
}

// seQuitaDelProyecto usa QuitarIntegrante de HU-03, no la base directamente.
func (w *registrarEsfuerzoWorld) seQuitaDelProyecto(integrante, proyecto string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.QuitarIntegrante(id, integrante); err != nil {
		return fmt.Errorf("no se pudo quitar a %q del proyecto %q: %w", integrante, proyecto, err)
	}
	return nil
}

// pasoAlEstado usa ModificarProyecto de HU-02, no la base directamente.
func (w *registrarEsfuerzoWorld) pasoAlEstado(proyecto, estado string) error {
	id, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CambiarEstado(id, estado); err != nil {
		return fmt.Errorf("no se pudo pasar el proyecto %q al estado %q: %w", proyecto, estado, err)
	}
	return nil
}

// --- When ---

func (w *registrarEsfuerzoWorld) registra(integrante, horas, actividad, elemento, fecha string) error {
	d, err := w.datos(integrante, horas, actividad, elemento)
	if err != nil {
		return err
	}
	if d.Fecha, err = parsearFechaLocal(fecha); err != nil {
		return err
	}
	w.ejecutar(d)
	return nil
}

func (w *registrarEsfuerzoWorld) registraSinFecha(integrante, horas, actividad, elemento string) error {
	d, err := w.datos(integrante, horas, actividad, elemento)
	if err != nil {
		return err
	}
	w.ejecutar(d) // d.Fecha queda en cero: "sin fecha"
	return nil
}

func (w *registrarEsfuerzoWorld) registraEnProyectoInexistente(integrante, horas, actividad, elemento string) error {
	d, err := w.datos(integrante, horas, actividad, elemento)
	if err != nil {
		return err
	}
	d.ProyectoID = idInexistente
	d.Fecha = w.hoy
	w.ejecutar(d)
	return nil
}

func (w *registrarEsfuerzoWorld) registraEnElementoInexistente(integrante, horas, actividad, fecha string) error {
	h, err := parsearHoras(horas)
	if err != nil {
		return err
	}
	f, err := parsearFechaLocal(fecha)
	if err != nil {
		return err
	}
	w.ejecutar(datosEsfuerzo{
		ProyectoID: w.proyectoID, ElementoID: idInexistente,
		Integrante: integrante, Fecha: f, Actividad: actividad, Horas: h,
	})
	return nil
}

// datos arma las entradas sobre el proyecto del Background. Integrante y actividad pasan tal
// cual, aunque sean vacíos o solo espacios: validarlos es trabajo del dominio.
func (w *registrarEsfuerzoWorld) datos(integrante, horas, actividad, elemento string) (datosEsfuerzo, error) {
	h, err := parsearHoras(horas)
	if err != nil {
		return datosEsfuerzo{}, err
	}
	elemID, err := w.idElemento(elemento)
	if err != nil {
		return datosEsfuerzo{}, err
	}
	return datosEsfuerzo{
		ProyectoID: w.proyectoID, ElementoID: elemID,
		Integrante: integrante, Actividad: actividad, Horas: h,
	}, nil
}

func (w *registrarEsfuerzoWorld) ejecutar(d datosEsfuerzo) {
	w.ultimo, w.ultimoErr = w.sistema.RegistrarEsfuerzo(d)
	if w.ultimoErr == nil {
		w.registrados = append(w.registrados, w.ultimo)
	}
}

// --- Then ---

func (w *registrarEsfuerzoWorld) seRegistraCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba registrar el esfuerzo, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) tieneIdentificador() error {
	if err := w.seRegistraCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.ID <= 0 {
		return fmt.Errorf("se esperaba un identificador asignado, se obtuvo %d", w.ultimo.ID)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) registroTieneDatos(integrante, fecha, actividad, horas string) error {
	if err := w.seRegistraCorrectamente(); err != nil {
		return err
	}
	h, err := parsearHoras(horas)
	if err != nil {
		return err
	}
	want := registroVisto{
		ID: w.ultimo.ID, ProyectoID: w.proyectoID, ElementoID: w.ultimo.ElementoID,
		Integrante: integrante, Fecha: fecha, Actividad: actividad, Horas: h,
	}
	if w.ultimo != want {
		return fmt.Errorf("registro: se esperaba %+v, se obtuvo %+v", want, w.ultimo)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU18()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el esfuerzo se registró", descripcion)
	}
	if !errors.Is(w.ultimoErr, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) elementoTieneRegistros(elemento string, esperado int) error {
	id, err := w.idElemento(elemento)
	if err != nil {
		return err
	}
	n, err := w.sistema.ContarRegistros(id)
	if err != nil {
		return fmt.Errorf("no se pudieron contar los registros de %q: %w", elemento, err)
	}
	if n != esperado {
		return fmt.Errorf("registros de %q: se esperaba %d, hay %d", elemento, esperado, n)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) identificadoresDistintos() error {
	if len(w.registrados) != 2 {
		return fmt.Errorf("se esperaban 2 registros, hay %d", len(w.registrados))
	}
	if a, b := w.registrados[0].ID, w.registrados[1].ID; a == b {
		return fmt.Errorf("los dos registros tienen el mismo identificador %d", a)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) recuperablesPorID() error {
	if len(w.registrados) == 0 {
		return errors.New("no hay registros para recuperar")
	}
	for _, r := range w.registrados {
		if err := w.recuperaIgual(r); err != nil {
			return err
		}
	}
	return nil
}

// registroSeSigueRecuperando cubre RN-18.3: los registros sobreviven a quitar al integrante.
func (w *registrarEsfuerzoWorld) registroSeSigueRecuperando(integrante string) error {
	encontrados := 0
	for _, r := range w.registrados {
		if r.Integrante == integrante {
			encontrados++
			if err := w.recuperaIgual(r); err != nil {
				return err
			}
		}
	}
	if encontrados == 0 {
		return fmt.Errorf("el escenario no preparó registros de %q", integrante)
	}
	return nil
}

func (w *registrarEsfuerzoWorld) recuperaIgual(r registroVisto) error {
	got, err := w.sistema.BuscarRegistro(r.ID)
	if err != nil {
		return fmt.Errorf("no se pudo recuperar el registro %d: %w", r.ID, err)
	}
	if got != r {
		return fmt.Errorf("registro %d: se esperaba %+v, se recuperó %+v", r.ID, r, got)
	}
	return nil
}

// --- Auxiliares ---

func (w *registrarEsfuerzoWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *registrarEsfuerzoWorld) idElemento(titulo string) (int64, error) {
	id, ok := w.elementos[titulo]
	if !ok {
		return 0, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return id, nil
}

// parsearHoras acepta la coma decimal de la spec ("0,5", "-1,25").
func parsearHoras(s string) (float64, error) {
	h, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("horas inválidas en el escenario %q: %w", s, err)
	}
	return h, nil
}

// parsearFechaLocal interpreta la fecha como un día de la zona del servidor (RN-18.4).
// parsearFecha (HU-01) usa UTC, que acá correría el día según la zona.
func parsearFechaLocal(s string) (time.Time, error) {
	t, err := time.ParseInLocation(time.DateOnly, s, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("fecha inválida en el escenario %q: %w", s, err)
	}
	return t, nil
}
