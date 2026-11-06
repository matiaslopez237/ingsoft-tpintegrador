package bdd

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/ejemplo"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// hu31World guarda el estado de un escenario de HU-31: el sistema y el último resultado de la carga.
type hu31World struct {
	sistema   *hu31Sistema
	resultado ejemplo.Resultado
}

func registrarStepsHU31(ctx *godog.ScenarioContext) {
	w := &hu31World{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		*w = hu31World{} // si falla la preparación, el After no cierra el sistema del escenario anterior
		sistema, err := hu31NuevoSistema()
		if err != nil {
			return c, err
		}
		w.sistema = sistema
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.sistema == nil {
			return c, nil
		}
		return c, w.sistema.Cerrar()
	})

	// Given
	ctx.Step(`^una base de datos vacía para el ejemplo$`, w.baseVacia)
	ctx.Step(`^la base del ejemplo ya tiene el proyecto "([^"]*)"$`, w.baseConProyecto)

	// When
	ctx.Step(`^arranca la carga del ejemplo con CARGAR_EJEMPLO "([^"]*)"$`, w.arrancaLaCarga)

	// Then
	ctx.Step(`^la carga del ejemplo informa que se cargó$`, w.informa(ejemplo.Cargado))
	ctx.Step(`^la carga del ejemplo informa que no se cargó porque la base ya tiene proyectos$`, w.informa(ejemplo.BaseConProyectos))
	ctx.Step(`^la carga del ejemplo informa que no se cargó porque la opción está apagada$`, w.informa(ejemplo.OpcionApagada))
	ctx.Step(`^la carga del ejemplo informa que el valor no es booleano$`, w.informa(ejemplo.ValorInvalido))
	ctx.Step(`^existe el proyecto de ejemplo "([^"]*)" en estado "([^"]*)"$`, w.existeProyecto)
	ctx.Step(`^el proyecto de ejemplo tiene estos integrantes:$`, w.integrantes)
	ctx.Step(`^el backlog del ejemplo tiene (\d+) elementos, de "HU-(\d+)" a "HU-(\d+)", que suman (\d+) story points$`, w.backlog)
	ctx.Step(`^los sprints del ejemplo son:$`, w.sprints)
	ctx.Step(`^en el sprint (\d+) del ejemplo están hechos solo "([^"]*)" y "([^"]*)"$`, w.hechosEnSprint)
	ctx.Step(`^la velocidad del ejemplo es el promedio de (\d+), (\d+) y (\d+)$`, w.velocidad)
	ctx.Step(`^los defectos del ejemplo son:$`, w.defectos)
	ctx.Step(`^las horas del proyecto de ejemplo tienen horas estimadas, horas reales y desviación$`, w.horas)
	ctx.Step(`^el historial de rondas de "([^"]*)" en el ejemplo tiene (\d+) rondas, la última con estimación acordada (\d+)$`, w.rondas)
	ctx.Step(`^la base del ejemplo tiene solo el proyecto "([^"]*)"$`, w.soloProyecto)
	ctx.Step(`^la base del ejemplo no tiene proyectos$`, w.sinProyectos)
}

func (w *hu31World) baseVacia() error {
	return nil // el Before ya abrió una base nueva
}

func (w *hu31World) baseConProyecto(ctx context.Context, nombre string) error {
	return w.sistema.crearProyecto(ctx, nombre)
}

func (w *hu31World) arrancaLaCarga(ctx context.Context, valor string) error {
	r, err := w.sistema.cargar(ctx, valor)
	if err != nil {
		return fmt.Errorf("la carga del ejemplo falló: %w", err)
	}
	w.resultado = r
	return nil
}

func (w *hu31World) informa(esperado ejemplo.Resultado) func() error {
	return func() error {
		if w.resultado != esperado {
			return fmt.Errorf("la carga informó %v, se esperaba %v", w.resultado, esperado)
		}
		return nil
	}
}

// proyectoDelEjemplo devuelve el único proyecto de la base: el de ejemplo.
func (w *hu31World) proyectoDelEjemplo(ctx context.Context) (domain.Proyecto, error) {
	ps, err := w.sistema.proyectos(ctx)
	if err != nil {
		return domain.Proyecto{}, err
	}
	if len(ps) != 1 {
		return domain.Proyecto{}, fmt.Errorf("se esperaba un solo proyecto, hay %d", len(ps))
	}
	return w.sistema.ver.Ejecutar(ctx, ps[0].ID)
}

func (w *hu31World) existeProyecto(ctx context.Context, nombre, estado string) error {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return err
	}
	if p.Nombre != nombre || string(p.Estado) != estado {
		return fmt.Errorf("el proyecto es %q en estado %q, se esperaba %q en estado %q", p.Nombre, p.Estado, nombre, estado)
	}
	return nil
}

func (w *hu31World) integrantes(ctx context.Context, tabla *godog.Table) error {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return err
	}
	var esperados, obtenidos []string
	for _, fila := range tabla.Rows[1:] {
		esperados = append(esperados, fila.Cells[0].Value+": "+fila.Cells[1].Value)
	}
	for _, i := range p.Integrantes {
		obtenidos = append(obtenidos, i.Nombre+": "+strings.Join(i.Roles, ", "))
	}
	slices.Sort(esperados)
	slices.Sort(obtenidos)
	if !slices.Equal(esperados, obtenidos) {
		return fmt.Errorf("los integrantes son %v, se esperaba %v", obtenidos, esperados)
	}
	return nil
}

// elementosDelEjemplo devuelve el backlog del proyecto de ejemplo.
func (w *hu31World) elementosDelEjemplo(ctx context.Context) (domain.Proyecto, []domain.ElementoBacklog, error) {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return domain.Proyecto{}, nil, err
	}
	es, err := w.sistema.backlog.Ejecutar(ctx, app.DatosListarBacklog{ProyectoID: p.ID})
	return p, es, err
}

// elementoPorHistoria busca el elemento cuyo título empieza con el ID de la historia ("HU-14").
func elementoPorHistoria(es []domain.ElementoBacklog, historia string) (domain.ElementoBacklog, bool) {
	for _, e := range es {
		if strings.HasPrefix(e.Titulo, historia+":") {
			return e, true
		}
	}
	return domain.ElementoBacklog{}, false
}

func (w *hu31World) backlog(ctx context.Context, cantidad, desde, hasta, puntos int) error {
	_, es, err := w.elementosDelEjemplo(ctx)
	if err != nil {
		return err
	}
	if len(es) != cantidad {
		return fmt.Errorf("el backlog tiene %d elementos, se esperaban %d", len(es), cantidad)
	}
	for n := desde; n <= hasta; n++ {
		if _, ok := elementoPorHistoria(es, fmt.Sprintf("HU-%02d", n)); !ok {
			return fmt.Errorf("falta el elemento de HU-%02d", n)
		}
	}
	suma := 0
	for _, e := range es {
		if e.StoryPoints != nil {
			suma += *e.StoryPoints
		}
	}
	if suma != puntos {
		return fmt.Errorf("el backlog suma %d story points, se esperaban %d", suma, puntos)
	}
	return nil
}

// sprintsDelEjemplo devuelve los sprints del proyecto de ejemplo, por número.
func (w *hu31World) sprintsDelEjemplo(ctx context.Context) ([]domain.Sprint, error) {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return nil, err
	}
	ss, err := w.sistema.sprints.Ejecutar(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(ss, func(a, b domain.Sprint) int { return a.Numero - b.Numero })
	return ss, nil
}

func (w *hu31World) sprints(ctx context.Context, tabla *godog.Table) error {
	ss, err := w.sprintsDelEjemplo(ctx)
	if err != nil {
		return err
	}
	var esperados, obtenidos []string
	for _, fila := range tabla.Rows[1:] {
		var celdas []string
		for _, c := range fila.Cells {
			celdas = append(celdas, c.Value)
		}
		esperados = append(esperados, strings.Join(celdas, " | "))
	}
	for _, s := range ss {
		obtenidos = append(obtenidos, fmt.Sprintf("%d | %s | %s | %s | %d | %d", s.Numero, s.Estado,
			s.FechaInicio.Format("2006-01-02"), s.FechaFin.Format("2006-01-02"), s.PuntosPlanificados, s.PuntosCompletados))
	}
	if !slices.Equal(esperados, obtenidos) {
		return fmt.Errorf("los sprints son:\n%s\nse esperaba:\n%s", strings.Join(obtenidos, "\n"), strings.Join(esperados, "\n"))
	}
	return nil
}

func (w *hu31World) hechosEnSprint(ctx context.Context, numero int, a, b string) error {
	ss, err := w.sprintsDelEjemplo(ctx)
	if err != nil {
		return err
	}
	i := slices.IndexFunc(ss, func(s domain.Sprint) bool { return s.Numero == numero })
	if i < 0 {
		return fmt.Errorf("no existe el sprint %d", numero)
	}
	_, es, err := w.elementosDelEjemplo(ctx)
	if err != nil {
		return err
	}
	var hechos []string
	for _, e := range es {
		if e.SprintID != nil && *e.SprintID == ss[i].ID && e.Estado == domain.EstadoHecho {
			hechos = append(hechos, strings.SplitN(e.Titulo, ":", 2)[0])
		}
	}
	slices.Sort(hechos)
	if esperados := []string{a, b}; !slices.Equal(hechos, esperados) {
		return fmt.Errorf("en el sprint %d están hechos %v, se esperaba %v", numero, hechos, esperados)
	}
	return nil
}

func (w *hu31World) velocidad(ctx context.Context, a, b, c int) error {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return err
	}
	v, err := w.sistema.vel.Ejecutar(ctx, p.ID)
	if err != nil {
		return err
	}
	esperada := float64(a+b+c) / 3
	if !v.HayDatos || v.SprintsCerrados != 3 || math.Abs(v.Velocidad-esperada) > 1e-9 {
		return fmt.Errorf("la velocidad es %v sobre %d sprints (hay datos: %v), se esperaba %v sobre 3", v.Velocidad, v.SprintsCerrados, v.HayDatos, esperada)
	}
	return nil
}

func (w *hu31World) defectos(ctx context.Context, tabla *godog.Table) error {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return err
	}
	ds, err := w.sistema.defs.Ejecutar(ctx, p.ID)
	if err != nil {
		return err
	}
	ss, err := w.sprintsDelEjemplo(ctx)
	if err != nil {
		return err
	}
	numero := func(id *int64) string {
		if id == nil {
			return "ninguno"
		}
		for _, s := range ss {
			if s.ID == *id {
				return strconv.Itoa(s.Numero)
			}
		}
		return "desconocido"
	}
	var esperados, obtenidos []string
	for _, fila := range tabla.Rows[1:] {
		esperados = append(esperados, fmt.Sprintf("%s | %s | %s | %s", fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value, fila.Cells[3].Value))
	}
	for _, d := range ds {
		obtenidos = append(obtenidos, fmt.Sprintf("%s | %s | %s | %s", strings.SplitN(d.Descripcion, ":", 2)[0], d.Estado,
			numero(d.SprintDeteccionID), numero(d.SprintResolucionID)))
	}
	slices.Sort(esperados)
	slices.Sort(obtenidos)
	if !slices.Equal(esperados, obtenidos) {
		return fmt.Errorf("los defectos son:\n%s\nse esperaba:\n%s", strings.Join(obtenidos, "\n"), strings.Join(esperados, "\n"))
	}
	return nil
}

func (w *hu31World) horas(ctx context.Context) error {
	p, err := w.proyectoDelEjemplo(ctx)
	if err != nil {
		return err
	}
	t, err := w.sistema.horas.Ejecutar(ctx, p.ID)
	if err != nil {
		return err
	}
	if t.HorasEstimadas <= 0 || t.HorasReales <= 0 || t.DesviacionHoras == nil {
		return fmt.Errorf("las horas son %v estimadas y %v reales, con desviación %v", t.HorasEstimadas, t.HorasReales, t.DesviacionHoras)
	}
	return nil
}

func (w *hu31World) rondas(ctx context.Context, historia string, cantidad, acordada int) error {
	p, es, err := w.elementosDelEjemplo(ctx)
	if err != nil {
		return err
	}
	e, ok := elementoPorHistoria(es, historia)
	if !ok {
		return fmt.Errorf("no existe el elemento de %s", historia)
	}
	rs, err := w.sistema.rondas.Ejecutar(ctx, p.ID, e.ID)
	if err != nil {
		return err
	}
	if len(rs) != cantidad {
		return fmt.Errorf("%s tiene %d rondas, se esperaban %d", historia, len(rs), cantidad)
	}
	slices.SortFunc(rs, func(a, b domain.ResumenRonda) int { return a.Numero - b.Numero })
	if ultima := rs[len(rs)-1]; ultima.EstimacionAcordada == nil || *ultima.EstimacionAcordada != acordada {
		return fmt.Errorf("la última ronda de %s tiene estimación %v, se esperaba %d", historia, ultima.EstimacionAcordada, acordada)
	}
	return nil
}

func (w *hu31World) soloProyecto(ctx context.Context, nombre string) error {
	ps, err := w.sistema.proyectos(ctx)
	if err != nil {
		return err
	}
	if len(ps) != 1 || ps[0].Nombre != nombre {
		var nombres []string
		for _, p := range ps {
			nombres = append(nombres, p.Nombre)
		}
		return fmt.Errorf("la base tiene los proyectos %v, se esperaba solo %q", nombres, nombre)
	}
	return nil
}

func (w *hu31World) sinProyectos(ctx context.Context) error {
	ps, err := w.sistema.proyectos(ctx)
	if err != nil {
		return err
	}
	if len(ps) != 0 {
		return fmt.Errorf("la base tiene %d proyectos, se esperaba ninguno", len(ps))
	}
	return nil
}
