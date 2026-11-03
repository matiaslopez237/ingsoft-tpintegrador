package bdd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu19Sistema es la puerta de entrada de los escenarios de HU-19 al núcleo. La implementa
// hu19SistemaReal, sobre SQLite (hu19_sistema_test.go). Todos los identificadores de paquete llevan
// el prefijo hu19 para no chocar con los de otras historias.
type hu19Sistema interface {
	CrearProyecto(nombre string) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre string) error
	QuitarIntegrante(proyectoID int64, nombre string) error
	// CambiarEstadoProyecto deja un proyecto preparado en otro estado (con HU-02).
	CambiarEstadoProyecto(proyectoID int64, estado string) error
	// CrearElemento crea el elemento con HU-05 y, si horasEstimadas no es nil, se las carga con
	// HU-06: un elemento siempre nace sin estimar (RN-19.3).
	CrearElemento(proyectoID int64, titulo string, storyPoints *int, horasEstimadas *float64) (int64, error)
	// RegistrarEsfuerzo registra horas con el caso de uso de HU-18.
	RegistrarEsfuerzo(proyectoID, elementoID int64, integrante string, horas float64) error
	// EditarHorasEstimadas edita con HU-06; con horas en nil, el elemento queda sin estimar.
	EditarHorasEstimadas(proyectoID, elementoID int64, horas *float64) error
	EditarStoryPoints(proyectoID, elementoID int64, storyPoints int) error
	Comparar(proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error)
	Elemento(elementoID int64) (domain.ElementoBacklog, error)
	Cerrar() error
}

// hu19Errores traduce la descripción del escenario al error definido en la spec (sección 7).
func hu19Errores() map[string]error {
	return map[string]error{
		"proyecto no encontrado":            domain.ErrProyectoNoEncontrado,
		"elemento no encontrado":            domain.ErrElementoNoEncontrado,
		"horas inválidas":                   domain.ErrHorasInvalidas,
		"horas no múltiplo de media hora":   domain.ErrHorasNoMultiploMediaHora,
		"horas estimadas exceden el máximo": domain.ErrHorasEstimadasExcedenMaximo,
	}
}

// hu19Elemento es un elemento que preparó el escenario. Los títulos no se repiten dentro de un escenario.
type hu19Elemento struct {
	proyectoID, id int64
}

type hu19World struct {
	sistema    hu19Sistema
	proyectos  map[string]int64
	integrante map[int64]string // proyecto → la integrante del Background, que carga lo que "acumula"
	elementos  map[string]hu19Elemento
	comparada  domain.ComparacionEsfuerzo
	errCompara error // el de la última comparación
	comparo    bool  // para no dar por buena una comparación que el escenario nunca hizo
	errEdicion error // el de la última edición
	edito      bool  // ídem, para la edición
}

func registrarStepsHU19(ctx *godog.ScenarioContext) {
	w := &hu19World{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = hu19World{} // si falla la preparación, el After no cierra el sistema del escenario anterior
		sistema, err := hu19NuevoSistema()
		if err != nil {
			return c, err
		}
		*w = hu19World{
			sistema:    sistema,
			proyectos:  map[string]int64{},
			integrante: map[int64]string{},
			elementos:  map[string]hu19Elemento{},
		}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	const horas = `(-?\d+(?:,\d+)?)` // con coma decimal, como en la spec

	// Given
	ctx.Step(`^el proyecto "([^"]*)" para comparar esfuerzo, con la integrante "([^"]*)"$`, w.proyectoConIntegrante)
	ctx.Step(`^el proyecto "([^"]*)" tiene el elemento "([^"]*)" sin horas estimadas$`, w.elementoSinEstimar)
	ctx.Step(`^el proyecto "([^"]*)" tiene el elemento "([^"]*)" con `+horas+` horas estimadas$`, w.elementoConHoras)
	ctx.Step(`^el proyecto "([^"]*)" tiene el elemento "([^"]*)" con (\d+) story points y sin horas estimadas$`, w.elementoConPuntos)
	ctx.Step(`^"([^"]*)" tiene cargadas `+horas+` horas estimadas$`, w.tieneCargadas)
	ctx.Step(`^"([^"]*)" acumula `+horas+` horas de esfuerzo$`, w.acumula)
	ctx.Step(`^"([^"]*)" registró `+horas+` horas de esfuerzo en "([^"]*)"$`, w.registro)
	ctx.Step(`^"([^"]*)" también integra el proyecto "([^"]*)"$`, w.tambienIntegra)
	ctx.Step(`^"([^"]*)" dejó de integrar el proyecto "([^"]*)"$`, w.dejoDeIntegrar)
	ctx.Step(`^"([^"]*)" está Pendiente y sin sprint$`, w.estaPendienteSinSprint)
	ctx.Step(`^el proyecto "([^"]*)" que se compara pasa a estar "([^"]*)"$`, w.pasaAEstar)

	// When
	ctx.Step(`^comparo el esfuerzo de "([^"]*)"$`, w.comparoElemento)
	ctx.Step(`^comparo el esfuerzo de "([^"]*)" desde el proyecto "([^"]*)"$`, w.comparoDesdeProyecto)
	ctx.Step(`^comparo el esfuerzo de "([^"]*)" desde un proyecto inexistente$`, w.comparoDesdeProyectoInexistente)
	ctx.Step(`^comparo el esfuerzo de un elemento inexistente de "([^"]*)"$`, w.comparoElementoInexistente)
	ctx.Step(`^calculo la comparación en el dominio con `+horas+` horas estimadas y `+horas+` horas reales$`, w.calculoEnDominio)
	ctx.Step(`^edito las horas estimadas de "([^"]*)" a `+horas+`$`, w.editoHoras)
	ctx.Step(`^quito las horas estimadas de "([^"]*)"$`, w.quitoHoras)
	ctx.Step(`^edito los story points de "([^"]*)" a (\d+)$`, w.editoPuntos)
	ctx.Step(`^creo el elemento "([^"]*)" en el proyecto "([^"]*)" que se compara$`, w.creoElemento)

	// Then
	ctx.Step(`^la comparación informa `+horas+` horas estimadas, `+horas+` horas reales y una desviación de `+horas+` horas y de `+horas+` %$`, w.informa)
	ctx.Step(`^la comparación está sin estimar, con `+horas+` horas reales y sin desviaciones$`, w.informaSinEstimar)
	ctx.Step(`^la comparación no está sin estimar, con `+horas+` horas reales, una desviación de `+horas+` horas y sin porcentaje$`, w.informaSinPorcentaje)
	ctx.Step(`^se rechaza la comparación con el error "([^"]*)"$`, w.seRechazaComparacion)
	ctx.Step(`^la edición de las horas estimadas termina sin errores$`, w.edicionSinErrores)
	ctx.Step(`^se rechaza la edición de las horas estimadas con el error "([^"]*)"$`, w.seRechazaEdicion)
	ctx.Step(`^"([^"]*)" tiene `+horas+` horas estimadas$`, w.tieneHoras)
	ctx.Step(`^"([^"]*)" está sin horas estimadas$`, w.estaSinHoras)
	ctx.Step(`^"([^"]*)" quedó con (\d+) story points$`, w.tienePuntos)
}

// --- Given ---

func (w *hu19World) proyectoConIntegrante(proyecto, integrante string) error {
	pid, err := w.sistema.CrearProyecto(proyecto)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el proyecto %q: %w", proyecto, err)
	}
	if err := w.sistema.AgregarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo preparar a %q en %q: %w", integrante, proyecto, err)
	}
	w.proyectos[proyecto] = pid
	w.integrante[pid] = integrante
	return nil
}

func (w *hu19World) elementoSinEstimar(proyecto, titulo string) error {
	return w.prepararElemento(proyecto, titulo, nil, nil)
}

func (w *hu19World) elementoConHoras(proyecto, titulo, horas string) error {
	h, err := hu19Horas(horas)
	if err != nil {
		return err
	}
	return w.prepararElemento(proyecto, titulo, nil, &h)
}

func (w *hu19World) elementoConPuntos(proyecto, titulo string, puntos int) error {
	return w.prepararElemento(proyecto, titulo, &puntos, nil)
}

func (w *hu19World) prepararElemento(proyecto, titulo string, puntos *int, horas *float64) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	id, err := w.sistema.CrearElemento(pid, titulo, puntos, horas)
	if err != nil {
		return fmt.Errorf("no se pudo preparar el elemento %q de %q: %w", titulo, proyecto, err)
	}
	w.elementos[titulo] = hu19Elemento{proyectoID: pid, id: id}
	return nil
}

func (w *hu19World) tieneCargadas(titulo, horas string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	h, err := hu19Horas(horas)
	if err != nil {
		return err
	}
	if err := w.sistema.EditarHorasEstimadas(e.proyectoID, e.id, &h); err != nil {
		return fmt.Errorf("no se pudieron cargar %s horas estimadas en %q: %w", horas, titulo, err)
	}
	return nil
}

// acumula reparte las horas en registros de hasta 24 h (RN-18.8), a nombre de la integrante del
// proyecto. Con 0 horas no registra nada: un registro de 0 h no es válido (RN-18.6).
func (w *hu19World) acumula(titulo, horas string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	resto, err := hu19Horas(horas)
	if err != nil {
		return err
	}
	for resto > 0 {
		parte := min(resto, domain.HorasMaximasPorRegistro)
		if err := w.sistema.RegistrarEsfuerzo(e.proyectoID, e.id, w.integrante[e.proyectoID], parte); err != nil {
			return fmt.Errorf("no se pudieron registrar %v horas en %q: %w", parte, titulo, err)
		}
		resto -= parte
	}
	return nil
}

func (w *hu19World) registro(integrante, horas, titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	h, err := hu19Horas(horas)
	if err != nil {
		return err
	}
	if err := w.sistema.RegistrarEsfuerzo(e.proyectoID, e.id, integrante, h); err != nil {
		return fmt.Errorf("no se pudieron registrar %s horas de %q en %q: %w", horas, integrante, titulo, err)
	}
	return nil
}

func (w *hu19World) tambienIntegra(integrante, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.AgregarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo agregar a %q en %q: %w", integrante, proyecto, err)
	}
	return nil
}

func (w *hu19World) dejoDeIntegrar(integrante, proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.QuitarIntegrante(pid, integrante); err != nil {
		return fmt.Errorf("no se pudo quitar a %q de %q: %w", integrante, proyecto, err)
	}
	return nil
}

// estaPendienteSinSprint verifica la preparación: es lo que CA-19.12 da por cierto.
func (w *hu19World) estaPendienteSinSprint(titulo string) error {
	got, err := w.leerElemento(titulo)
	if err != nil {
		return err
	}
	if got.Estado != domain.EstadoPendiente || got.SprintID != nil {
		return fmt.Errorf("%q: se esperaba Pendiente y sin sprint, está %q con sprint %v", titulo, got.Estado, got.SprintID)
	}
	return nil
}

func (w *hu19World) pasaAEstar(proyecto, estado string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	if err := w.sistema.CambiarEstadoProyecto(pid, estado); err != nil {
		return fmt.Errorf("no se pudo dejar %q en estado %q: %w", proyecto, estado, err)
	}
	return nil
}

// --- When ---

func (w *hu19World) comparoElemento(titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	w.comparar(e.proyectoID, e.id)
	return nil
}

func (w *hu19World) comparoDesdeProyecto(titulo, proyecto string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.comparar(pid, e.id)
	return nil
}

func (w *hu19World) comparoDesdeProyectoInexistente(titulo string) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	w.comparar(idInexistente, e.id)
	return nil
}

func (w *hu19World) comparoElementoInexistente(proyecto string) error {
	pid, err := w.idProyecto(proyecto)
	if err != nil {
		return err
	}
	w.comparar(pid, idInexistente)
	return nil
}

func (w *hu19World) comparar(proyectoID, elementoID int64) {
	w.comparada, w.errCompara = w.sistema.Comparar(proyectoID, elementoID)
	w.comparo = true
}

// calculoEnDominio llama a la función del dominio sin pasar por el sistema: un elemento no puede
// tener estimado 0 (RN-19.2), así que el caso no se puede armar con datos del proyecto (RN-19.8).
func (w *hu19World) calculoEnDominio(estimadas, reales string) error {
	e, err := hu19Horas(estimadas)
	if err != nil {
		return err
	}
	r, err := hu19Horas(reales)
	if err != nil {
		return err
	}
	w.comparada = domain.CompararEsfuerzo(&e, r)
	w.comparo = true
	return nil
}

func (w *hu19World) editoHoras(titulo, horas string) error {
	h, err := hu19Horas(horas)
	if err != nil {
		return err
	}
	return w.editarHoras(titulo, &h)
}

func (w *hu19World) quitoHoras(titulo string) error {
	return w.editarHoras(titulo, nil)
}

func (w *hu19World) editarHoras(titulo string, horas *float64) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	w.errEdicion = w.sistema.EditarHorasEstimadas(e.proyectoID, e.id, horas)
	w.edito = true
	return nil
}

func (w *hu19World) editoPuntos(titulo string, puntos int) error {
	e, err := w.elemento(titulo)
	if err != nil {
		return err
	}
	w.errEdicion = w.sistema.EditarStoryPoints(e.proyectoID, e.id, puntos)
	w.edito = true
	return nil
}

func (w *hu19World) creoElemento(titulo, proyecto string) error {
	return w.prepararElemento(proyecto, titulo, nil, nil)
}

// --- Then ---

func (w *hu19World) informa(estimadas, reales, desviacion, porcentaje string) error {
	c, err := w.comparacion()
	if err != nil {
		return err
	}
	if c.SinEstimar() {
		return fmt.Errorf("la comparación quedó sin estimar: %s", hu19Ver(c))
	}
	if err := hu19Igual("horas estimadas", c.HorasEstimadas, estimadas); err != nil {
		return err
	}
	if err := hu19IgualValor("horas reales", c.HorasReales, reales); err != nil {
		return err
	}
	if err := hu19Igual("desviación en horas", c.DesviacionHoras, desviacion); err != nil {
		return err
	}
	return hu19Igual("desviación en porcentaje", c.DesviacionPorcentaje, porcentaje)
}

func (w *hu19World) informaSinEstimar(reales string) error {
	c, err := w.comparacion()
	if err != nil {
		return err
	}
	if !c.SinEstimar() {
		return fmt.Errorf("se esperaba sin estimar: %s", hu19Ver(c))
	}
	if c.HorasEstimadas != nil || c.DesviacionHoras != nil || c.DesviacionPorcentaje != nil {
		return fmt.Errorf("sin estimar no debería informar estimación ni desviaciones: %s", hu19Ver(c))
	}
	return hu19IgualValor("horas reales", c.HorasReales, reales)
}

func (w *hu19World) informaSinPorcentaje(reales, desviacion string) error {
	c, err := w.comparacion()
	if err != nil {
		return err
	}
	if c.SinEstimar() {
		return fmt.Errorf("no se esperaba sin estimar: %s", hu19Ver(c))
	}
	if c.DesviacionPorcentaje != nil {
		return fmt.Errorf("no se esperaba porcentaje: %s", hu19Ver(c))
	}
	if err := hu19IgualValor("horas reales", c.HorasReales, reales); err != nil {
		return err
	}
	return hu19Igual("desviación en horas", c.DesviacionHoras, desviacion)
}

func (w *hu19World) seRechazaComparacion(descripcion string) error {
	if !w.comparo {
		return errors.New("el escenario no hizo ninguna comparación")
	}
	return hu19EsError(w.errCompara, descripcion)
}

func (w *hu19World) edicionSinErrores() error {
	if !w.edito {
		return errors.New("el escenario no hizo ninguna edición")
	}
	if w.errEdicion != nil {
		return fmt.Errorf("se esperaba editar el elemento, pero hubo error: %w", w.errEdicion)
	}
	return nil
}

func (w *hu19World) seRechazaEdicion(descripcion string) error {
	if !w.edito {
		return errors.New("el escenario no hizo ninguna edición")
	}
	return hu19EsError(w.errEdicion, descripcion)
}

func (w *hu19World) tieneHoras(titulo, horas string) error {
	got, err := w.leerElemento(titulo)
	if err != nil {
		return err
	}
	return hu19Igual("horas estimadas de "+titulo, got.HorasEstimadas, horas)
}

func (w *hu19World) estaSinHoras(titulo string) error {
	got, err := w.leerElemento(titulo)
	if err != nil {
		return err
	}
	if got.HorasEstimadas != nil {
		return fmt.Errorf("%q: se esperaba sin horas estimadas, tiene %v", titulo, *got.HorasEstimadas)
	}
	return nil
}

func (w *hu19World) tienePuntos(titulo string, puntos int) error {
	got, err := w.leerElemento(titulo)
	if err != nil {
		return err
	}
	if got.StoryPoints == nil || *got.StoryPoints != puntos {
		return fmt.Errorf("%q: se esperaban %d story points, tiene %v", titulo, puntos, got.StoryPoints)
	}
	return nil
}

// --- Auxiliares ---

// comparacion devuelve el resultado de la última comparación, o un error si no la hubo o falló.
func (w *hu19World) comparacion() (domain.ComparacionEsfuerzo, error) {
	if !w.comparo {
		return domain.ComparacionEsfuerzo{}, errors.New("el escenario no hizo ninguna comparación")
	}
	if w.errCompara != nil {
		return domain.ComparacionEsfuerzo{}, fmt.Errorf("se esperaba comparar, pero hubo error: %w", w.errCompara)
	}
	return w.comparada, nil
}

func (w *hu19World) leerElemento(titulo string) (domain.ElementoBacklog, error) {
	e, err := w.elemento(titulo)
	if err != nil {
		return domain.ElementoBacklog{}, err
	}
	got, err := w.sistema.Elemento(e.id)
	if err != nil {
		return domain.ElementoBacklog{}, fmt.Errorf("no se pudo leer el elemento %q: %w", titulo, err)
	}
	return got, nil
}

func (w *hu19World) elemento(titulo string) (hu19Elemento, error) {
	e, ok := w.elementos[titulo]
	if !ok {
		return hu19Elemento{}, fmt.Errorf("el elemento %q no se preparó en el escenario", titulo)
	}
	return e, nil
}

func (w *hu19World) idProyecto(nombre string) (int64, error) {
	id, ok := w.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el proyecto %q no se preparó en el escenario", nombre)
	}
	return id, nil
}

// hu19Horas lee un número con coma decimal, como se escribe en el .feature.
func hu19Horas(s string) (float64, error) {
	h, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("número mal escrito en el escenario: %q", s)
	}
	return h, nil
}

// hu19Igual compara un valor opcional con el del escenario. La comparación es exacta a propósito:
// los valores llegan redondeados a un decimal (RN-19.6).
func hu19Igual(campo string, got *float64, esperado string) error {
	if got == nil {
		return fmt.Errorf("%s: se esperaba %s, no tiene valor", campo, esperado)
	}
	return hu19IgualValor(campo, *got, esperado)
}

func hu19IgualValor(campo string, got float64, esperado string) error {
	e, err := hu19Horas(esperado)
	if err != nil {
		return err
	}
	if got != e {
		return fmt.Errorf("%s: se esperaba %s, es %v", campo, esperado, got)
	}
	return nil
}

func hu19EsError(got error, descripcion string) error {
	esperado, ok := hu19Errores()[descripcion]
	if !ok {
		return fmt.Errorf("error desconocido en el escenario: %q", descripcion)
	}
	if got == nil {
		return fmt.Errorf("se esperaba el error %q, pero no hubo error", descripcion)
	}
	if !errors.Is(got, esperado) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo: %w", descripcion, got)
	}
	return nil
}

// hu19Ver muestra una comparación legible en los mensajes: sus campos son punteros.
func hu19Ver(c domain.ComparacionEsfuerzo) string {
	ver := func(v *float64) string {
		if v == nil {
			return "sin valor"
		}
		return strconv.FormatFloat(*v, 'f', -1, 64)
	}
	return fmt.Sprintf("estimadas %s, reales %v, desviación %s h y %s %%",
		ver(c.HorasEstimadas), c.HorasReales, ver(c.DesviacionHoras), ver(c.DesviacionPorcentaje))
}
