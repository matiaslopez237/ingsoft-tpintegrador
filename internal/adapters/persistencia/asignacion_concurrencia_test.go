package persistencia

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// backlogConLecturaVieja simula pedidos simultáneos: BuscarElementoPorID siempre devuelve la
// misma foto tomada al principio, aunque la base ya haya cambiado. BuscarSprint y GuardarAsignacion
// van al repositorio real, igual que repoConLecturaVieja en integrantes_concurrencia_test.go.
type backlogConLecturaVieja struct {
	*RepositorioBacklog
	foto domain.ElementoBacklog
}

func (r backlogConLecturaVieja) BuscarElementoPorID(context.Context, int64) (domain.ElementoBacklog, error) {
	return r.foto, nil
}

// Traza: HU-09 (comentario de revisión): dos pedidos que asignan el mismo elemento a dos sprints
// Abiertos distintos "al mismo tiempo" (los dos leyeron el elemento sin sprint asignado) no pueden
// terminar los dos "ganando": RN-09.2 se tiene que cumplir aunque el UPDATE llegue sin validar de
// nuevo en memoria.
func TestAsignacion_DosAsignacionesSimultaneas_LaSegundaFalla(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	repoSprints := NuevoRepositorioSprints(repo.db)
	s1 := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))
	s2 := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 28, 29))
	vieja := backlogConLecturaVieja{RepositorioBacklog: repo, foto: e} // los dos "leyeron" sin sprint asignado

	if err := app.NuevoAsignarHistoria(vieja).Ejecutar(ctx, app.DatosAsignacion{
		ProyectoID: ids[0], ElementoID: e.ID, SprintID: s1.ID,
	}); err != nil {
		t.Fatalf("primera asignación: %v", err)
	}
	err := app.NuevoAsignarHistoria(vieja).Ejecutar(ctx, app.DatosAsignacion{
		ProyectoID: ids[0], ElementoID: e.ID, SprintID: s2.ID,
	})

	if !errors.Is(err, domain.ErrElementoYaAsignado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoYaAsignado)
	}
	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.SprintID == nil || *got.SprintID != s1.ID {
		t.Errorf("SprintID = %v, se esperaba %d (la primera asignación, sin pisar)", got.SprintID, s1.ID)
	}
}

// Traza: HU-09 (comentario de revisión): si el sprint del elemento se cierra entre la validación y
// el guardado, GuardarAsignacion no puede desasignar igual — ni siquiera si quien llama no volvió a
// consultar el estado del sprint antes de guardar (RN-09.4 se verifica en el propio UPDATE).
func TestRepositorioBacklog_GuardarAsignacion_DesasignarConSprintYaCerrado_Falla(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	repoSprints := NuevoRepositorioSprints(repo.db)
	s := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))
	e.SprintID = &s.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("preparar asignación: %v", err)
	}
	cerrarSprintDePrueba(t, repo.db, s.ID) // otro pedido (HU-11) lo cierra en el medio

	e.SprintID = nil
	err := repo.GuardarAsignacion(ctx, &e)

	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	got, errBuscar := repo.BuscarElementoPorID(ctx, e.ID)
	if errBuscar != nil {
		t.Fatalf("buscar: %v", errBuscar)
	}
	if got.SprintID == nil || *got.SprintID != s.ID {
		t.Errorf("SprintID = %v, se esperaba que siguiera %d (la desasignación no se tiene que aplicar)", got.SprintID, s.ID)
	}
}

// Traza: HU-09 RN-09.2 (un sprint Cerrado no bloquea moverse a otro Abierto): confirma que la
// condición atómica nueva no rompe este caso, que antes resolvía la validación en memoria.
func TestRepositorioBacklog_GuardarAsignacion_ReasignaDesdeSprintCerrado(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	repoSprints := NuevoRepositorioSprints(repo.db)
	cerrado := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 26, 27))
	e.SprintID = &cerrado.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("preparar asignación: %v", err)
	}
	cerrarSprintDePrueba(t, repo.db, cerrado.ID)
	abierto := guardarSprintDePrueba(t, repoSprints, sprintDePrueba(ids[0], 28, 29))

	e.SprintID = &abierto.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("guardar reasignación: %v", err)
	}

	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.SprintID == nil || *got.SprintID != abierto.ID {
		t.Errorf("SprintID = %v, se esperaba %d", got.SprintID, abierto.ID)
	}
}
