package persistencia

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: HU-11, sección 5 y P-8: si varios pedidos cierran el mismo sprint a la vez, uno lo cierra y
// los demás reciben ErrSprintYaCerrado. Todos llegan al UPDATE con el sprint ya cerrado en memoria,
// como lo deja el caso de uso después de validar: la única barrera es la condición del UPDATE.
func TestRepositorioSprints_CerrarSprint_Simultaneos_SoloUnoGana(t *testing.T) {
	const pedidos = 10
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	repo := NuevoRepositorioSprints(proyectos.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 26, 27))

	var (
		wg      sync.WaitGroup
		largada = make(chan struct{}) // las goroutines arrancan juntas al cerrarlo
		errs    = make([]error, pedidos)
	)
	for i := range pedidos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cerrado := s
			cerrado.Estado = domain.SprintCerrado
			<-largada
			errs[i] = repo.CerrarSprint(ctx, &cerrado)
		}()
	}
	close(largada)
	wg.Wait()

	exitos := 0
	for i, err := range errs {
		switch {
		case err == nil:
			exitos++
		case !errors.Is(err, domain.ErrSprintYaCerrado):
			t.Errorf("pedido %d: error = %v, se esperaba nil o %v", i, err, domain.ErrSprintYaCerrado)
		}
	}
	if exitos != 1 {
		t.Errorf("cerraron el sprint %d pedidos, se esperaba exactamente 1", exitos)
	}
	if got, err := repo.BuscarSprintPorID(ctx, s.ID); err != nil || got.Estado != domain.SprintCerrado {
		t.Errorf("sprint = %+v (error %v), se esperaba Cerrado", got, err)
	}
}

// Traza: HU-11, sección 5: un cierre simultáneo con varios completados (HU-10) deja todo coherente.
// Cada historia termina Hecho y en el sprint (se completó antes del cierre) o Pendiente y sin sprint
// (el cierre llegó primero y el completado falló), y los puntos completados son los de las Hecho.
func TestRepositorioSprints_CerrarSprint_SimultaneoConCompletados_QuedaCoherente(t *testing.T) {
	const historias = 8
	ctx := context.Background()
	backlog, ids := nuevoBacklogDePrueba(t, "Alfa")
	repo := NuevoRepositorioSprints(backlog.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 26, 27))
	hids := make([]int64, historias)
	for i := range historias {
		hids[i] = historiaEnSprintDePrueba(t, backlog, ids[0], s.ID, fmt.Sprintf("Historia %d", i), puntosDePrueba(i+1), false)
	}

	var (
		wg        sync.WaitGroup
		largada   = make(chan struct{})
		errCierre error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		cerrado := s
		cerrado.Estado = domain.SprintCerrado
		<-largada
		errCierre = repo.CerrarSprint(ctx, &cerrado)
	}()
	for _, id := range hids {
		e := buscarElementoDePrueba(t, backlog, id) // fuera de la goroutine: usa t.Fatalf
		e.Estado = domain.EstadoHecho
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-largada
			_ = backlog.GuardarCompletado(ctx, &e) // puede fallar si el cierre llegó antes
		}()
	}
	close(largada)
	wg.Wait()

	if errCierre != nil {
		t.Fatalf("cerrar: %v", errCierre)
	}
	planificados, completados := 0, 0
	for i, id := range hids {
		e := buscarElementoDePrueba(t, backlog, id)
		planificados += i + 1
		switch {
		case e.Estado == domain.EstadoHecho && e.SprintID != nil && *e.SprintID == s.ID:
			completados += i + 1
		case e.Estado == domain.EstadoPendiente && e.SprintID == nil:
		default:
			t.Errorf("%q: Estado = %q, SprintID = %v; se esperaba Hecho en el sprint o Pendiente sin sprint", e.Titulo, e.Estado, e.SprintID)
		}
	}
	got, err := repo.BuscarSprintPorID(ctx, s.ID)
	if err != nil {
		t.Fatalf("buscar sprint: %v", err)
	}
	if got.PuntosPlanificados != planificados || got.PuntosCompletados != completados {
		t.Errorf("puntos = %d y %d, se esperaba %d y %d", got.PuntosPlanificados, got.PuntosCompletados, planificados, completados)
	}
}
