package persistencia

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// completarConSprintViejo simula pedidos simultáneos en HU-10: BuscarSprint informa el estado del
// sprint tal como estaba al principio (Abierto), aunque otro pedido ya lo haya cerrado. El resto va
// al repositorio real, igual que backlogConLecturaVieja en asignacion_concurrencia_test.go.
type completarConSprintViejo struct {
	*RepositorioBacklog
	proyectoID int64
}

func (r completarConSprintViejo) BuscarSprint(context.Context, int64) (int64, bool, error) {
	return r.proyectoID, true, nil
}

// Traza: HU-10, sección 5 y 6: el sprint se cierra entre la validación y el guardado →
// ErrSprintCerrado y el elemento no cambia, porque RN-10.2 se verifica en el propio UPDATE.
func TestCompletar_SprintCerradoEntreLecturaYGuardado_Falla(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	cerrarSprintDePrueba(t, repo.db, *e.SprintID) // otro pedido (HU-11) lo cierra en el medio
	vieja := completarConSprintViejo{RepositorioBacklog: repo, proyectoID: ids[0]}
	proyectos := NuevoRepositorioProyectos(repo.db)

	err := app.NuevoCompletarHistoria(proyectos, vieja).Ejecutar(ctx, app.DatosCompletar{
		ProyectoID: ids[0], ElementoID: e.ID,
	})

	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoPendiente {
		t.Errorf("Estado = %q, se esperaba que siguiera %q", got.Estado, domain.EstadoPendiente)
	}
}

// completarConLecturaVieja simula que todo lo que lee el caso de uso es la foto del principio: el
// elemento en su sprint original y ese sprint Abierto. GuardarCompletado va al repositorio real.
type completarConLecturaVieja struct {
	completarConSprintViejo
	foto domain.ElementoBacklog
}

func (r completarConLecturaVieja) BuscarElementoPorID(context.Context, int64) (domain.ElementoBacklog, error) {
	return r.foto, nil
}

// Traza: HU-10, sección 5 (revisión M3): el elemento estaba en el sprint S; entre la validación y
// el guardado S se cerró y el elemento se movió a otro sprint Abierto T (HU-09). Por el caso de uso
// real → ErrElementoSinSprint y el elemento no cambia, aunque T esté Abierto.
func TestCompletar_MovidoAOtroSprintEntreLecturaYGuardado_Falla(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	cerrarSprintDePrueba(t, repo.db, *e.SprintID) // otro pedido (HU-11) cierra S
	otro := guardarSprintDePrueba(t, NuevoRepositorioSprints(repo.db), sprintDePrueba(ids[0], 28, 29))
	movido := e
	movido.SprintID = &otro.ID
	if err := repo.GuardarAsignacion(ctx, &movido); err != nil { // otro pedido (HU-09) lo mueve a T
		t.Fatalf("mover al otro sprint: %v", err)
	}
	vieja := completarConLecturaVieja{
		completarConSprintViejo: completarConSprintViejo{RepositorioBacklog: repo, proyectoID: ids[0]},
		foto:                    e,
	}
	proyectos := NuevoRepositorioProyectos(repo.db)

	err := app.NuevoCompletarHistoria(proyectos, vieja).Ejecutar(ctx, app.DatosCompletar{
		ProyectoID: ids[0], ElementoID: e.ID,
	})

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
	got := buscarElementoDePrueba(t, repo, e.ID)
	if got.Estado != domain.EstadoPendiente {
		t.Errorf("Estado = %q, se esperaba que siguiera %q", got.Estado, domain.EstadoPendiente)
	}
	if got.SprintID == nil || *got.SprintID != otro.ID {
		t.Errorf("SprintID = %v, se esperaba que siguiera en el sprint %d", got.SprintID, otro.ID)
	}
}
