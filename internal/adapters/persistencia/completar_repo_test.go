package persistencia

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioCompletar = (*RepositorioBacklog)(nil)

// elementoEnSprintDePrueba guarda un elemento Pendiente del proyecto y lo asigna a un sprint
// Abierto nuevo del mismo proyecto. Devuelve el elemento tal como quedó guardado.
func elementoEnSprintDePrueba(t *testing.T, repo *RepositorioBacklog, proyectoID int64) domain.ElementoBacklog {
	t.Helper()
	ctx := context.Background()
	e := elementoDePrueba(t, proyectoID, "Login", nil)
	if err := repo.GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	s := guardarSprintDePrueba(t, NuevoRepositorioSprints(repo.db), sprintDePrueba(proyectoID, 26, 27))
	e.SprintID = &s.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatalf("asignar elemento: %v", err)
	}
	return buscarElementoDePrueba(t, repo, e.ID)
}

func buscarElementoDePrueba(t *testing.T, repo *RepositorioBacklog, id int64) domain.ElementoBacklog {
	t.Helper()
	e, err := repo.BuscarElementoPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("buscar elemento %d: %v", id, err)
	}
	return e
}

// Traza: HU-10 CA-10.1 (RN-10.1)
func TestRepositorioBacklog_GuardarCompletado_GuardaHecho(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])

	e.Estado = domain.EstadoHecho
	if err := repo.GuardarCompletado(context.Background(), &e); err != nil {
		t.Fatalf("guardar completado: %v", err)
	}

	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoHecho {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.EstadoHecho)
	}
}

// Traza: HU-10 RN-10.1: el UPDATE cambia solo el estado; aunque el elemento recibido traiga otros
// campos distintos, lo guardado conserva los suyos.
func TestRepositorioBacklog_GuardarCompletado_NoCambiaOtrosCampos(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	antes := elementoEnSprintDePrueba(t, repo, ids[0])

	e := antes
	e.Estado = domain.EstadoHecho
	e.Titulo, e.Descripcion, e.Prioridad = "Otro", "Otra", domain.PrioridadAlta
	ocho := 8
	e.StoryPoints = &ocho
	if err := repo.GuardarCompletado(context.Background(), &e); err != nil {
		t.Fatalf("guardar completado: %v", err)
	}

	esperado := antes
	esperado.Estado = domain.EstadoHecho
	if got := buscarElementoDePrueba(t, repo, e.ID); !reflect.DeepEqual(got, esperado) {
		t.Errorf("elemento = %+v, se esperaba %+v", got, esperado)
	}
}

// Traza: HU-10 RN-10.3: completar un elemento ya Hecho en un sprint Abierto no es un error.
func TestRepositorioBacklog_GuardarCompletado_YaHecho_SinError(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	e.Estado = domain.EstadoHecho
	if err := repo.GuardarCompletado(ctx, &e); err != nil {
		t.Fatalf("primer completado: %v", err)
	}

	if err := repo.GuardarCompletado(ctx, &e); err != nil {
		t.Errorf("segundo completado: error = %v, se esperaba nil (RN-10.3)", err)
	}
	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoHecho {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.EstadoHecho)
	}
}

// Traza: HU-10, sección 7: elemento inexistente
func TestRepositorioBacklog_GuardarCompletado_Inexistente_RetornaError(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])

	e.ID = 999
	e.Estado = domain.EstadoHecho
	err := repo.GuardarCompletado(context.Background(), &e)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-10 CA-10.3 (RN-10.2), sección 5: el elemento ya no tiene sprint al guardar (se lo
// quitó HU-09 mientras tanto). El estado no cambia.
func TestRepositorioBacklog_GuardarCompletado_SinSprint_RetornaError(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	quitado := e
	quitado.SprintID = nil
	if err := repo.GuardarAsignacion(ctx, &quitado); err != nil {
		t.Fatalf("quitar del sprint: %v", err)
	}

	e.Estado = domain.EstadoHecho // e conserva el SprintID que se leyó antes de quitarlo
	err := repo.GuardarCompletado(ctx, &e)

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoPendiente {
		t.Errorf("Estado = %q, se esperaba que siguiera %q", got.Estado, domain.EstadoPendiente)
	}
}

// Traza: HU-10 CA-10.4 (RN-10.2): el sprint del elemento está Cerrado. El estado no cambia.
func TestRepositorioBacklog_GuardarCompletado_SprintCerrado_RetornaError(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	cerrarSprintDePrueba(t, repo.db, *e.SprintID)

	e.Estado = domain.EstadoHecho
	err := repo.GuardarCompletado(context.Background(), &e)

	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoPendiente {
		t.Errorf("Estado = %q, se esperaba que siguiera %q", got.Estado, domain.EstadoPendiente)
	}
}

// Traza: HU-10, sección 5. Caso no definido en la spec: entre la lectura y el guardado el sprint
// del elemento se cerró y el elemento se movió a otro sprint Abierto (HU-09, RN-09.2). Se trata
// como "ya no está en ese sprint": ErrElementoSinSprint, sin cambiar el estado.
func TestRepositorioBacklog_GuardarCompletado_MovidoAOtroSprint_RetornaError(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoEnSprintDePrueba(t, repo, ids[0])
	cerrarSprintDePrueba(t, repo.db, *e.SprintID)
	otro := guardarSprintDePrueba(t, NuevoRepositorioSprints(repo.db), sprintDePrueba(ids[0], 28, 29))
	movido := e
	movido.SprintID = &otro.ID
	if err := repo.GuardarAsignacion(ctx, &movido); err != nil {
		t.Fatalf("mover al otro sprint: %v", err)
	}

	e.Estado = domain.EstadoHecho // e conserva el sprint (ahora Cerrado) que se leyó antes
	err := repo.GuardarCompletado(ctx, &e)

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
	if got := buscarElementoDePrueba(t, repo, e.ID); got.Estado != domain.EstadoPendiente {
		t.Errorf("Estado = %q, se esperaba que siguiera %q", got.Estado, domain.EstadoPendiente)
	}
}
