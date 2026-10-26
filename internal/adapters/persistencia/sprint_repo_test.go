package persistencia

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioSprints = (*RepositorioSprints)(nil)

func dia(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.UTC) }

func sprintDePrueba(proyectoID int64, desde, hasta int) domain.Sprint {
	return domain.Sprint{ProyectoID: proyectoID, FechaInicio: dia(desde), FechaFin: dia(hasta),
		Objetivo: "Gestionar el backlog", Estado: domain.SprintAbierto}
}

func guardarSprintDePrueba(t *testing.T, repo *RepositorioSprints, s domain.Sprint) domain.Sprint {
	t.Helper()
	if err := repo.GuardarSprint(context.Background(), &s); err != nil {
		t.Fatalf("guardar sprint: %v", err)
	}
	return s
}

// Traza: HU-08 CA-08.1, CA-08.2, CA-08.9: se guardan con ID y número, y se listan por número.
func TestRepositorioSprints_GuardarYListar(t *testing.T) {
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	repo := NuevoRepositorioSprints(proyectos.db)

	// Se guardan en orden inverso de fechas: el número sigue el orden de creación.
	primero := guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 28, 29))
	segundo := guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 26, 27))

	if primero.ID == 0 || segundo.ID == 0 || primero.Numero != 1 || segundo.Numero != 2 {
		t.Fatalf("guardados = %+v y %+v, se esperaban ID asignado y números 1 y 2", primero, segundo)
	}
	got, err := repo.ListarSprints(ctx, p.ID)
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(got) != 2 || got[0] != primero || got[1] != segundo {
		t.Errorf("listados = %+v, se esperaba %+v y %+v en ese orden", got, primero, segundo)
	}
}

// Traza: HU-08 CA-08.2: cada proyecto tiene su propia numeración.
func TestRepositorioSprints_NumeracionPorProyecto(t *testing.T) {
	proyectos := nuevoRepoDePrueba(t)
	alfa := guardarDePrueba(t, proyectos, "Alfa")
	beta := guardarDePrueba(t, proyectos, "Beta")
	repo := NuevoRepositorioSprints(proyectos.db)

	guardarSprintDePrueba(t, repo, sprintDePrueba(alfa.ID, 26, 27))
	s := guardarSprintDePrueba(t, repo, sprintDePrueba(beta.ID, 26, 27))

	if s.Numero != 1 {
		t.Errorf("número del primer sprint de Beta = %d, se esperaba 1", s.Numero)
	}
}

// Traza: HU-08 sección 5 y RN-08.5: la base rechaza un sprint superpuesto aunque el dominio no lo
// haya visto (otro pedido lo creó en el medio), y no guarda nada.
func TestRepositorioSprints_SuperpuestoAlGuardar_RetornaError(t *testing.T) {
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	repo := NuevoRepositorioSprints(proyectos.db)
	guardarSprintDePrueba(t, repo, sprintDePrueba(p.ID, 26, 27))

	nuevo := sprintDePrueba(p.ID, 27, 28)
	err := repo.GuardarSprint(ctx, &nuevo)

	if !errors.Is(err, domain.ErrSprintSuperpuesto) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintSuperpuesto)
	}
	if got, _ := repo.ListarSprints(ctx, p.ID); len(got) != 1 {
		t.Errorf("hay %d sprints, se esperaba solo el primero", len(got))
	}
}

// Traza: HU-08 RN-08.1: si el proyecto ya no existe al guardar, no se crea el sprint.
func TestRepositorioSprints_ProyectoInexistente_RetornaError(t *testing.T) {
	repo := NuevoRepositorioSprints(nuevoRepoDePrueba(t).db)
	s := sprintDePrueba(99, 26, 27)

	err := repo.GuardarSprint(context.Background(), &s)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}

// sprintsSinVer simula pedidos simultáneos: ListarSprints no ve los sprints que otro pedido acaba
// de crear (los dos leyeron antes de que cualquiera guardara). Guardar va al repositorio real.
type sprintsSinVer struct{ *RepositorioSprints }

func (sprintsSinVer) ListarSprints(context.Context, int64) ([]domain.Sprint, error) { return nil, nil }

// Traza: HU-08 sección 5: dos sprints creados a la vez sin superponerse reciben números distintos.
func TestCrearSprint_DosSimultaneos_RecibenNumerosDistintos(t *testing.T) {
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	sinVer := sprintsSinVer{NuevoRepositorioSprints(proyectos.db)}
	uc := app.NuevoCrearSprint(proyectos, sinVer)

	a, errA := uc.Ejecutar(ctx, app.DatosNuevoSprint{ProyectoID: p.ID, FechaInicio: dia(26), FechaFin: dia(27), Objetivo: "A"})
	b, errB := uc.Ejecutar(ctx, app.DatosNuevoSprint{ProyectoID: p.ID, FechaInicio: dia(28), FechaFin: dia(29), Objetivo: "B"})

	if errA != nil || errB != nil {
		t.Fatalf("errores inesperados: %v, %v", errA, errB)
	}
	if a.Numero != 1 || b.Numero != 2 {
		t.Errorf("números = %d y %d, se esperaban 1 y 2", a.Numero, b.Numero)
	}
}

// Traza: HU-08 sección 5 y RN-08.5: si dos pedidos crean a la vez sprints superpuestos, el segundo
// se rechaza al guardar.
func TestCrearSprint_DosSimultaneosSuperpuestos_SegundoEsRechazado(t *testing.T) {
	ctx := context.Background()
	proyectos := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, proyectos, "Alfa")
	sinVer := sprintsSinVer{NuevoRepositorioSprints(proyectos.db)}
	uc := app.NuevoCrearSprint(proyectos, sinVer)

	if _, err := uc.Ejecutar(ctx, app.DatosNuevoSprint{ProyectoID: p.ID, FechaInicio: dia(26), FechaFin: dia(27), Objetivo: "A"}); err != nil {
		t.Fatalf("primer sprint: %v", err)
	}
	_, err := uc.Ejecutar(ctx, app.DatosNuevoSprint{ProyectoID: p.ID, FechaInicio: dia(27), FechaFin: dia(28), Objetivo: "B"})

	if !errors.Is(err, domain.ErrSprintSuperpuesto) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintSuperpuesto)
	}
}
