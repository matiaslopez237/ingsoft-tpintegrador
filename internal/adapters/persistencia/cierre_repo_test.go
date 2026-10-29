package persistencia

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador de sprints también cumple el puerto angosto de HU-11 (se verifica al compilar).
var _ app.RepositorioCierre = (*RepositorioSprints)(nil)

// Traza: HU-11, sección 5: BuscarSprintPorID devuelve el sprint tal como se guardó.
func TestRepositorioSprints_BuscarSprintPorID_Existente(t *testing.T) {
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	repo := NuevoRepositorioSprints(proyectos.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 26, 27))

	got, err := repo.BuscarSprintPorID(context.Background(), s.ID)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got != s {
		t.Errorf("sprint = %+v, se esperaba %+v", got, s)
	}
}

// Traza: HU-11 CA-11.3
func TestRepositorioSprints_BuscarSprintPorID_Inexistente(t *testing.T) {
	repo := NuevoRepositorioSprints(nuevoRepoDePrueba(t).db)

	_, err := repo.BuscarSprintPorID(context.Background(), 99)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
}

// historiaEnSprintDePrueba guarda una historia con esos story points (nil: sin estimar), la asigna
// al sprint (Abierto) y, si hecha, la completa. Devuelve su ID.
func historiaEnSprintDePrueba(t *testing.T, backlog *RepositorioBacklog, proyectoID, sprintID int64,
	titulo string, sp *int, hecha bool) int64 {
	t.Helper()
	ctx := context.Background()
	e := elementoDePrueba(t, proyectoID, titulo, sp)
	if err := backlog.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar %q: %v", titulo, err)
	}
	e.SprintID = &sprintID
	if err := backlog.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("asignar %q: %v", titulo, err)
	}
	if hecha {
		e.Estado = domain.EstadoHecho
		if err := backlog.GuardarCompletado(ctx, &e); err != nil {
			t.Fatalf("completar %q: %v", titulo, err)
		}
	}
	return e.ID
}

func puntosDePrueba(n int) *int { return &n }

// Traza: HU-11 CA-11.1, CA-11.6, CA-11.7: el sprint queda Cerrado con sus demás datos y su resumen
// de puntos (RN-11.6); las historias sin terminar vuelven al backlog y las Hecho quedan (RN-11.3).
// Las historias de otro sprint no se tocan.
func TestRepositorioSprints_CerrarSprint_Abierto_QuedaCerradoConResumenYDevuelveLasSinTerminar(t *testing.T) {
	ctx := context.Background()
	backlog, ids := nuevoBacklogDePrueba(t, "Alfa")
	repo := NuevoRepositorioSprints(backlog.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 26, 27))
	otro := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 28, 29))
	login := historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, "Login", puntosDePrueba(5), true)
	reportes := historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, "Reportes", puntosDePrueba(3), false)
	ayuda := historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, "Ayuda", nil, true)
	ajena := historiaEnSprintDePrueba(t, backlog, ids[0], otro.ID, "Ajena", puntosDePrueba(8), false)

	cerrado := s
	cerrado.Estado = domain.SprintCerrado
	if err := repo.CerrarSprint(ctx, &cerrado); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	want := s
	want.Estado, want.PuntosPlanificados, want.PuntosCompletados = domain.SprintCerrado, 8, 5
	if cerrado != want {
		t.Errorf("sprint recibido = %+v, se esperaba %+v", cerrado, want)
	}
	got, err := repo.BuscarSprintPorID(ctx, s.ID)
	if err != nil {
		t.Fatalf("buscar sprint: %v", err)
	}
	if got != want {
		t.Errorf("sprint guardado = %+v, se esperaba %+v", got, want)
	}
	if e := buscarElementoDePrueba(t, backlog, reportes); e.SprintID != nil || e.Estado != domain.EstadoPendiente {
		t.Errorf("Reportes: SprintID = %v, Estado = %q; se esperaba sin sprint y Pendiente", e.SprintID, e.Estado)
	}
	for _, id := range []int64{login, ayuda} {
		if e := buscarElementoDePrueba(t, backlog, id); e.SprintID == nil || *e.SprintID != s.ID {
			t.Errorf("%q: SprintID = %v, se esperaba que siguiera en %d", e.Titulo, e.SprintID, s.ID)
		}
	}
	if e := buscarElementoDePrueba(t, backlog, ajena); e.SprintID == nil || *e.SprintID != otro.ID {
		t.Errorf("Ajena: SprintID = %v, se esperaba que siguiera en %d", e.SprintID, otro.ID)
	}
}

// Traza: HU-11 RN-11.6: los puntos registrados no cambian si después una historia devuelta se
// asigna a otro sprint.
func TestRepositorioSprints_CerrarSprint_ResumenFijoAunqueSeReasigne(t *testing.T) {
	ctx := context.Background()
	backlog, ids := nuevoBacklogDePrueba(t, "Alfa")
	repo := NuevoRepositorioSprints(backlog.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 26, 27))
	siguiente := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 28, 29))
	reportes := historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, "Reportes", puntosDePrueba(3), false)
	cerrado := s
	cerrado.Estado = domain.SprintCerrado
	if err := repo.CerrarSprint(ctx, &cerrado); err != nil {
		t.Fatalf("cerrar: %v", err)
	}

	e := buscarElementoDePrueba(t, backlog, reportes)
	e.SprintID = &siguiente.ID
	if err := backlog.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("reasignar: %v", err)
	}

	got, err := repo.BuscarSprintPorID(ctx, s.ID)
	if err != nil {
		t.Fatalf("buscar sprint: %v", err)
	}
	if got.PuntosPlanificados != 3 || got.PuntosCompletados != 0 {
		t.Errorf("puntos = %d planificados, %d completados; se esperaba 3 y 0", got.PuntosPlanificados, got.PuntosCompletados)
	}
}

// Traza: HU-11 CA-11.2: tampoco devuelve al backlog las historias que quedaron en el sprint.
func TestRepositorioSprints_CerrarSprint_YaCerrado_RetornaErrorSinCambios(t *testing.T) {
	ctx := context.Background()
	backlog, ids := nuevoBacklogDePrueba(t, "Alfa")
	repo := NuevoRepositorioSprints(backlog.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 26, 27))
	login := historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, "Login", puntosDePrueba(5), false)
	cerrarSprintDePrueba(t, backlog.db, s.ID)
	antes, err := repo.BuscarSprintPorID(ctx, s.ID)
	if err != nil {
		t.Fatalf("buscar sprint: %v", err)
	}

	cerrado := s
	cerrado.Estado = domain.SprintCerrado
	err = repo.CerrarSprint(ctx, &cerrado)

	if !errors.Is(err, domain.ErrSprintYaCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintYaCerrado)
	}
	if got, _ := repo.BuscarSprintPorID(ctx, s.ID); got != antes {
		t.Errorf("sprint = %+v, se esperaba sin cambios %+v", got, antes)
	}
	if e := buscarElementoDePrueba(t, backlog, login); e.SprintID == nil || *e.SprintID != s.ID {
		t.Errorf("SprintID = %v, se esperaba que siguiera en %d", e.SprintID, s.ID)
	}
}

// Traza: HU-11 CA-11.3
func TestRepositorioSprints_CerrarSprint_Inexistente_RetornaError(t *testing.T) {
	repo := NuevoRepositorioSprints(nuevoRepoDePrueba(t).db)
	s := domain.Sprint{ID: 99, Estado: domain.SprintCerrado}

	err := repo.CerrarSprint(context.Background(), &s)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
}
