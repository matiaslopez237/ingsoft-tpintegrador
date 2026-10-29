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

// Traza: HU-11 CA-11.1, CA-11.6: el sprint queda Cerrado con sus demás datos, y las historias
// asignadas siguen asignadas a él (RN-11.3).
func TestRepositorioSprints_CerrarSprint_Abierto_QuedaCerradoYConservaHistorias(t *testing.T) {
	ctx := context.Background()
	backlog, ids := nuevoBacklogDePrueba(t, "Alfa")
	repo := NuevoRepositorioSprints(backlog.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(ids[0], 26, 27))
	e := elementoDePrueba(t, ids[0], "Login", nil)
	if err := backlog.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	e.SprintID = &s.ID
	if err := backlog.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("asignar elemento: %v", err)
	}

	cerrado := s
	cerrado.Estado = domain.SprintCerrado
	if err := repo.CerrarSprint(ctx, &cerrado); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	got, err := repo.BuscarSprintPorID(ctx, s.ID)
	if err != nil {
		t.Fatalf("buscar sprint: %v", err)
	}
	if got != cerrado {
		t.Errorf("sprint = %+v, se esperaba %+v", got, cerrado)
	}
	elem, err := backlog.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("buscar elemento: %v", err)
	}
	if elem.SprintID == nil || *elem.SprintID != s.ID {
		t.Errorf("SprintID = %v, se esperaba que siguiera %d", elem.SprintID, s.ID)
	}
}

// Traza: HU-11 CA-11.2
func TestRepositorioSprints_CerrarSprint_YaCerrado_RetornaErrorSinCambios(t *testing.T) {
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	repo := NuevoRepositorioSprints(proyectos.db)
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 26, 27))
	cerrarSprintDePrueba(t, proyectos.db, s.ID)
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
