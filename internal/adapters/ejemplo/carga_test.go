package ejemplo

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultas son los casos de uso de lectura con los que el test verifica lo que cargó la carga.
type consultas struct {
	listar  *app.ListarProyectos
	ver     *app.ConsultarProyecto
	backlog *app.ListarBacklog
	sprints *app.ListarSprints
	defs    *app.ListarDefectos
	horas   *app.ConsultarHorasProyecto
	rondas  *app.ListarRondas
}

// casosReales arma los casos de uso sobre una base SQLite nueva, como cmd/server/main.go, con el
// reloj de la carga en los que dependen de la fecha (RN-31.5).
func casosReales(t *testing.T) (CasosDeUso, consultas) {
	t.Helper()
	db, err := persistencia.AbrirDB(filepath.Join(t.TempDir(), "ejemplo.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	proyectos := persistencia.NuevoRepositorioProyectos(db)
	backlog := persistencia.NuevoRepositorioBacklog(db)
	sprints := persistencia.NuevoRepositorioSprints(db)
	rondas := persistencia.NuevoRepositorioRondas(db)
	esfuerzo := persistencia.NuevoRepositorioEsfuerzo(db)
	defectos := persistencia.NuevoRepositorioDefectos(db)
	reloj := &Reloj{}
	return CasosDeUso{
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
	}, consultas{
		listar:  app.NuevoListarProyectos(proyectos),
		ver:     app.NuevoConsultarProyecto(proyectos),
		backlog: app.NuevoListarBacklog(proyectos, backlog),
		sprints: app.NuevoListarSprints(proyectos, sprints),
		defs:    app.NuevoListarDefectos(proyectos, defectos),
		horas:   app.NuevoConsultarHorasProyecto(proyectos, esfuerzo),
		rondas:  app.NuevoListarRondas(proyectos, backlog, rondas),
	}
}

func dia(mes time.Month, d int) time.Time {
	return time.Date(2026, mes, d, 0, 0, 0, 0, time.Local)
}

// contenidoChico tiene un sprint con lo que puede pasar en uno: asignar, votar dos rondas, registrar
// esfuerzo, detectar y resolver un defecto, completar una historia y cerrarlo con la otra sin
// terminar. Después, un defecto que aparece entre sprints.
func contenidoChico() contenido {
	return contenido{
		proyecto:    "Demo",
		inicio:      dia(time.October, 1),
		integrantes: []integrante{{"Ana", "Dev"}, {"Beto", "QA"}},
		historias: []historia{
			{id: "HU-01", titulo: "Uno", descripcion: "Primera", prioridad: domain.PrioridadAlta, puntos: 3, criterios: []string{"c1"}, horasEstimadas: 8},
			{id: "HU-02", titulo: "Dos", descripcion: "Segunda", prioridad: domain.PrioridadMedia, puntos: 2, criterios: []string{"c2"}},
		},
		pasos: []paso{
			crearSprint{numero: 1, inicio: dia(time.October, 1), fin: dia(time.October, 5), objetivo: "Arrancar"},
			asignar{el: dia(time.October, 1), historia: "HU-01", sprint: 1},
			asignar{el: dia(time.October, 1), historia: "HU-02", sprint: 1},
			ronda{el: dia(time.October, 1), historia: "HU-02", votos: []voto{{"Ana", 1}, {"Beto", 8}}},
			ronda{el: dia(time.October, 1), historia: "HU-02", votos: []voto{{"Ana", 2}, {"Beto", 2}}, acordada: 2},
			esfuerzo{el: dia(time.October, 2), historia: "HU-01", integrante: "Ana", actividad: "Desarrollo", horas: 10},
			defecto{el: dia(time.October, 3), clave: "DEF-01", historia: "HU-01", descripcion: "Falla", severidad: domain.SeveridadAlta},
			cambioDefecto{el: dia(time.October, 3), clave: "DEF-01", accion: "iniciar"},
			cambioDefecto{el: dia(time.October, 4), clave: "DEF-01", accion: "resolver"},
			completar{el: dia(time.October, 4), historia: "HU-01"},
			cerrarSprint{el: dia(time.October, 5), sprint: 1},
			defecto{el: dia(time.October, 7), clave: "DEF-02", historia: "HU-02", descripcion: "Otra", severidad: domain.SeveridadBaja},
		},
	}
}

// Traza: HU-31 (RN-31.2, RN-31.4, RN-31.5)
func TestCargar_AplicaLosPasosEnOrdenConElRelojDeCadaDia(t *testing.T) {
	ctx := context.Background()
	cu, q := casosReales(t)

	if err := cargar(ctx, cu, contenidoChico()); err != nil {
		t.Fatalf("cargar: %v", err)
	}

	ps, _ := q.listar.Ejecutar(ctx)
	if len(ps) != 1 {
		t.Fatalf("se esperaba un proyecto, hay %d", len(ps))
	}
	p, _ := q.ver.Ejecutar(ctx, ps[0].ID)
	if p.Nombre != "Demo" || len(p.Integrantes) != 2 || p.Integrantes[0].Roles[0] == "" {
		t.Fatalf("proyecto = %+v, se esperaba Demo con Ana y Beto", p)
	}

	es, _ := q.backlog.Ejecutar(ctx, app.DatosListarBacklog{ProyectoID: p.ID})
	porTitulo := map[string]domain.ElementoBacklog{}
	for _, e := range es {
		porTitulo[e.Titulo] = e
	}
	uno, dos := porTitulo["HU-01: Uno"], porTitulo["HU-02: Dos"]
	if uno.Estado != domain.EstadoHecho || *uno.StoryPoints != 3 || uno.HorasEstimadas == nil || *uno.HorasEstimadas != 8 {
		t.Errorf("HU-01 = %+v, se esperaba Hecho, 3 puntos y 8 h estimadas", uno)
	}
	if dos.Estado == domain.EstadoHecho || dos.SprintID != nil || dos.StoryPoints == nil || *dos.StoryPoints != 2 {
		t.Errorf("HU-02 = %+v, se esperaba sin terminar, de vuelta en el backlog y con 2 puntos", dos)
	}

	ss, _ := q.sprints.Ejecutar(ctx, p.ID)
	if len(ss) != 1 || ss[0].Estado != domain.SprintCerrado || ss[0].PuntosPlanificados != 5 || ss[0].PuntosCompletados != 3 {
		t.Fatalf("sprints = %+v, se esperaba el 1 cerrado con 5 planificados y 3 completados", ss)
	}

	rs, _ := q.rondas.Ejecutar(ctx, p.ID, dos.ID)
	slices.SortFunc(rs, func(a, b domain.ResumenRonda) int { return a.Numero - b.Numero }) // el historial viene de la más nueva a la más vieja
	if len(rs) != 2 || rs[1].EstimacionAcordada == nil || *rs[1].EstimacionAcordada != 2 {
		t.Errorf("rondas de HU-02 = %+v, se esperaban 2, la última con 2 acordado", rs)
	}

	h, _ := q.horas.Ejecutar(ctx, p.ID)
	if h.HorasReales != 10 {
		t.Errorf("horas reales = %v, se esperaban 10", h.HorasReales)
	}

	ds, _ := q.defs.Ejecutar(ctx, p.ID)
	var got []string
	for _, d := range ds {
		got = append(got, d.Descripcion+" "+string(d.Estado)+" "+sprintDe(d.SprintDeteccionID, ss)+" "+sprintDe(d.SprintResolucionID, ss))
	}
	slices.Sort(got)
	want := []string{"DEF-01: Falla Resuelto 1 1", "DEF-02: Otra Abierto - -"}
	if !slices.Equal(got, want) {
		t.Errorf("defectos = %v, se esperaba %v", got, want)
	}
}

func sprintDe(id *int64, ss []domain.Sprint) string {
	if id == nil {
		return "-"
	}
	for _, s := range ss {
		if s.ID == *id {
			return string(rune('0' + s.Numero))
		}
	}
	return "?"
}

// cerradorQueFalla simula que un caso de uso falla a mitad de la carga (CA-31.10).
type cerradorQueFalla struct{ err error }

func (f cerradorQueFalla) Ejecutar(context.Context, app.DatosCerrarSprint) error { return f.err }

// Traza: HU-31 CA-31.10 (sección 7)
func TestCargar_UnPasoQueFalla_DevuelveElErrorConElPaso(t *testing.T) {
	cu, _ := casosReales(t)
	errCierre := errors.New("falla al cerrar")
	cu.CerrarSprint = cerradorQueFalla{err: errCierre}

	err := cargar(context.Background(), cu, contenidoChico())

	if !errors.Is(err, errCierre) {
		t.Fatalf("error = %v, se esperaba que envuelva %v", err, errCierre)
	}
	if !strings.Contains(err.Error(), "cerrar el sprint 1") {
		t.Errorf("el error %q no nombra el paso que falló", err)
	}
}
