package bdd

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"
)

// sistemaDefecto es la puerta de entrada de los escenarios de HU-20 al núcleo.
// La implementación real está en hu20_sistema_test.go.
type sistemaDefecto interface {
	CrearProyecto(nombre string, inicio time.Time) (int64, error)
	CrearElemento(proyectoID int64, titulo string) (int64, error)
	CrearSprint(proyectoID int64, inicio, fin time.Time) error
	RegistrarDefecto(d datosDefecto) (defectoVisto, error)
	Cerrar() error
}

// datosDefecto son las entradas de HU-20 (spec, sección 2).
type datosDefecto struct {
	ProyectoID  int64
	ElementoID  int64
	Descripcion string
	Severidad   string
}

// defectoVisto es lo que los escenarios miran de un defecto (spec, sección 3).
// SprintDeteccionID queda en nil cuando no hay Sprint vigente (RN-20.5).
type defectoVisto struct {
	ID                int64
	ProyectoID        int64
	ElementoID        int64
	Descripcion       string
	Severidad         string
	Estado            string
	SprintDeteccionID *int64
}

type registrarDefectoWorld struct {
	sistema     sistemaDefecto
	hoy         time.Time // lo lee el reloj falso del sistema
	proyectos   map[string]int64
	elementos   map[string]int64
	proyectoID  int64 // el proyecto del Background: ahí se registra si no se indica otro
	registrados []defectoVisto
	ultimo      defectoVisto
	ultimoErr   error
}

func registrarStepsHU20(ctx *godog.ScenarioContext) {
	w := &registrarDefectoWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = registrarDefectoWorld{proyectos: map[string]int64{}, elementos: map[string]int64{}}
		sistema, err := nuevoSistemaDefecto(w.reloj)
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

	ctx.Step(`^la fecha de hoy es "([^"]*)"$`, w.hoyEs)
	ctx.Step(`^el proyecto "([^"]*)" con inicio "([^"]*)" tiene el elemento "([^"]*)"$`, w.proyectoConElemento)
	ctx.Step(`^el proyecto "([^"]*)" tiene un sprint con fechas "([^"]*)" a "([^"]*)"$`, w.proyectoTieneSprint)

	ctx.Step(`^registro el defecto "([^"]*)" con severidad "([^"]*)" en "([^"]*)"$`, w.registro)
	ctx.Step(`^registro el defecto "([^"]*)" con severidad "([^"]*)" en un proyecto inexistente$`, w.registroEnProyectoInexistente)
	ctx.Step(`^registro el defecto "([^"]*)" con severidad "([^"]*)" en un elemento inexistente$`, w.registroEnElementoInexistente)

	ctx.Step(`^el defecto se registra correctamente$`, w.seRegistraCorrectamente)
	ctx.Step(`^el defecto tiene un identificador asignado$`, w.tieneIdentificador)
	ctx.Step(`^el defecto tiene la descripción "([^"]*)", la severidad "([^"]*)", el estado "([^"]*)" y el Sprint de detección asignado$`, w.defectoTieneDatosConSprint)
	ctx.Step(`^el defecto tiene la descripción "([^"]*)", la severidad "([^"]*)", el estado "([^"]*)" y sin Sprint de detección$`, w.defectoTieneDatosSinSprint)
	ctx.Step(`^se rechaza el defecto con el error "([^"]*)"$`, w.seRechazaConError)
	ctx.Step(`^los dos defectos tienen identificadores distintos$`, w.identificadoresDistintos)
}

// reloj es el reloj falso que recibe el caso de uso: devuelve el "hoy" del escenario.
func (w *registrarDefectoWorld) reloj() time.Time {
	return w.hoy
}

// --- Given ---

// hoyEs fija "hoy" al mediodía, en la zona del servidor, igual que HU-18 (RN-18.4): así la
// comparación por día no depende de la hora.
func (w *registrarDefectoWorld) hoyEs(fecha string) error {
	d, err := parsearFechaLocal(fecha)
	if err != nil {
		return err
	}
	w.hoy = d.Add(12 * time.Hour)
	return nil
}

func (w *registrarDefectoWorld) proyectoConElemento(proyecto, inicio, elemento string) error {
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

func (w *registrarDefectoWorld) proyectoTieneSprint(proyecto, inicio, fin string) error {
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

// --- When ---

func (w *registrarDefectoWorld) registro(descripcion, severidad, elemento string) error {
	d, err := w.datos(descripcion, severidad, elemento)
	if err != nil {
		return err
	}
	w.ejecutar(d)
	return nil
}

func (w *registrarDefectoWorld) registroEnProyectoInexistente(descripcion, severidad, elemento string) error {
	d, err := w.datos(descripcion, severidad, elemento)
	if err != nil {
		return err
	}
	d.ProyectoID = idInexistente
	w.ejecutar(d)
	return nil
}

func (w *registrarDefectoWorld) registroEnElementoInexistente(descripcion, severidad string) error {
	w.ejecutar(datosDefecto{
		ProyectoID: w.proyectoID, ElementoID: idInexistente,
		Descripcion: descripcion, Severidad: severidad,
	})
	return nil
}

// datos arma las entradas sobre el proyecto del Background. La descripción y la severidad pasan
// tal cual, aunque sean vacías o inválidas: validarlas es trabajo del dominio.
func (w *registrarDefectoWorld) datos(descripcion, severidad, elemento string) (datosDefecto, error) {
	elemID, err := w.idElemento(elemento)
	if err != nil {
		return datosDefecto{}, err
	}
	return datosDefecto{
		ProyectoID: w.proyectoID, ElementoID: elemID,
		Descripcion: descripcion, Severidad: severidad,
	}, nil
}

func (w *registrarDefectoWorld) ejecutar(d datosDefecto) {
	w.ultimo, w.ultimoErr = w.sistema.RegistrarDefecto(d)
	if w.ultimoErr == nil {
		w.registrados = append(w.registrados, w.ultimo)
	}
}

// --- Then ---

func (w *registrarDefectoWorld) seRegistraCorrectamente() error {
	if w.ultimoErr != nil {
		return fmt.Errorf("se esperaba registrar el defecto, pero hubo error: %w", w.ultimoErr)
	}
	return nil
}

func (w *registrarDefectoWorld) tieneIdentificador() error {
	if err := w.seRegistraCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.ID <= 0 {
		return fmt.Errorf("se esperaba un identificador asignado, se obtuvo %d", w.ultimo.ID)
	}
	return nil
}

func (w *registrarDefectoWorld) defectoTieneDatosConSprint(descripcion, severidad, estado string) error {
	if err := w.datosCoinciden(descripcion, severidad, estado); err != nil {
		return err
	}
	if w.ultimo.SprintDeteccionID == nil {
		return fmt.Errorf("se esperaba un Sprint de detección asignado, no se asignó ninguno")
	}
	return nil
}

func (w *registrarDefectoWorld) defectoTieneDatosSinSprint(descripcion, severidad, estado string) error {
	if err := w.datosCoinciden(descripcion, severidad, estado); err != nil {
		return err
	}
	if w.ultimo.SprintDeteccionID != nil {
		return fmt.Errorf("se esperaba sin Sprint de detección, se asignó el %d", *w.ultimo.SprintDeteccionID)
	}
	return nil
}

func (w *registrarDefectoWorld) datosCoinciden(descripcion, severidad, estado string) error {
	if err := w.seRegistraCorrectamente(); err != nil {
		return err
	}
	if w.ultimo.Descripcion != descripcion || w.ultimo.Severidad != severidad || w.ultimo.Estado != estado {
		return fmt.Errorf("defecto: se esperaba descripción %q, severidad %q y estado %q; se obtuvo descripción %q, severidad %q y estado %q",
			descripcion, severidad, estado, w.ultimo.Descripcion, w.ultimo.Severidad, w.ultimo.Estado)
	}
	return nil
}

func (w *registrarDefectoWorld) seRechazaConError(descripcion string) error {
	esperado, ok := erroresHU20()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if w.ultimoErr == nil {
		return fmt.Errorf("se esperaba el error %q, pero el defecto se registró", descripcion)
	}
	if !esperado.coincide(w.ultimoErr) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, w.ultimoErr)
	}
	return nil
}

func (w *registrarDefectoWorld) identificadoresDistintos() error {
	if len(w.registrados) != 2 {
		return fmt.Errorf("se esperaban 2 defectos, hay %d", len(w.registrados))
	}
	if a, b := w.registrados[0].ID, w.registrados[1].ID; a == b {
		return fmt.Errorf("los dos defectos tienen el mismo identificador %d", a)
	}
	return nil
}

// --- Auxiliares ---

func (w *registrarDefectoWorld) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

func (w *registrarDefectoWorld) idElemento(titulo string) (int64, error) {
	id, ok := w.elementos[titulo]
	if !ok {
		return 0, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return id, nil
}
