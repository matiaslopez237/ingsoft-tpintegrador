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
