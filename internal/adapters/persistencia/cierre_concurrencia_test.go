package persistencia

import (
	"context"
	"errors"
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
