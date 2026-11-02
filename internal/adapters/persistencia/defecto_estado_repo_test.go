package persistencia

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir los puertos de HU-21 (se verifica al compilar).
var (
	_ app.RepositorioEstadoDefecto = (*RepositorioDefectos)(nil)
	_ app.LectorDefectos           = (*RepositorioDefectos)(nil)
)

// defectoGuardado crea y guarda un defecto Abierto sobre el elemento de la base.
func (b baseDefecto) defectoGuardado(t *testing.T, severidad domain.Severidad) domain.Defecto {
	t.Helper()
	d := b.defecto(t, severidad)
	if err := b.repo.GuardarDefecto(context.Background(), &d); err != nil {
		t.Fatalf("guardar defecto de prueba: %v", err)
	}
	return d
}

// Traza: HU-21 (BuscarDefectoPorID lo usa el caso de uso ActualizarEstadoDefecto)
func TestRepositorioDefectos_BuscarDefectoPorID_Existente_LoDevuelve(t *testing.T) {
	b := nuevaBaseDefecto(t)
	d := b.defectoGuardado(t, domain.SeveridadAlta)

	got, err := b.repo.BuscarDefectoPorID(context.Background(), d.ID)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got != d {
		t.Errorf("se obtuvo %+v, se esperaba %+v", got, d)
	}
}

func TestRepositorioDefectos_BuscarDefectoPorID_Inexistente_RetornaError(t *testing.T) {
	b := nuevaBaseDefecto(t)

	_, err := b.repo.BuscarDefectoPorID(context.Background(), 999)

	if !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoEncontrado)
	}
}

// Traza: HU-21 CA-21.1 (RN-21.5: se guarda el Sprint de resolución)
func TestRepositorioDefectos_ActualizarEstado_ConSprintResolucion_LoGuarda(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d := b.defectoGuardado(t, domain.SeveridadAlta)
	s := domain.Sprint{ProyectoID: b.proyectoID, FechaInicio: dia(20), FechaFin: dia(30),
		Objetivo: "Sprint de prueba", Estado: domain.SprintAbierto}
	if err := b.sprints.GuardarSprint(ctx, &s); err != nil {
		t.Fatalf("guardar sprint: %v", err)
	}
	d.Estado, d.SprintResolucionID = domain.DefectoResuelto, &s.ID

	if err := b.repo.ActualizarEstadoDefecto(ctx, &d, domain.DefectoAbierto); err != nil {
		t.Fatal(err)
	}

	got, err := b.repo.BuscarDefectoPorID(ctx, d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Estado != domain.DefectoResuelto {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.DefectoResuelto)
	}
	if got.SprintResolucionID == nil || *got.SprintResolucionID != s.ID {
		t.Errorf("SprintResolucionID = %v, se esperaba un puntero a %d", got.SprintResolucionID, s.ID)
	}
}

// Traza: HU-21, sección 5 (pedidos simultáneos: el estado se verifica en el mismo UPDATE)
func TestRepositorioDefectos_ActualizarEstado_YaCambiadoPorOtroPedido_RetornaError(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d := b.defectoGuardado(t, domain.SeveridadAlta)
	// otro pedido ya lo resolvió: en la base sigue Abierto como "anterior" esperado por este pedido.
	ya := d
	ya.Estado = domain.DefectoResuelto
	if err := b.repo.ActualizarEstadoDefecto(ctx, &ya, domain.DefectoAbierto); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	cambio := d
	cambio.Estado = domain.DefectoResuelto
	err := b.repo.ActualizarEstadoDefecto(ctx, &cambio, domain.DefectoAbierto)

	if !errors.Is(err, domain.ErrDefectoNoAbierto) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoAbierto)
	}
}

// Traza: HU-21 CA-21.8
func TestRepositorioDefectos_ListarDefectos_DelProyecto_LosDevuelve(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d1 := b.defectoGuardado(t, domain.SeveridadAlta)
	d2 := b.defectoGuardado(t, domain.SeveridadBaja)

	got, err := b.repo.ListarDefectos(ctx, b.proyectoID)

	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("se esperaban 2 defectos, hay %d", len(got))
	}
	ids := map[int64]bool{got[0].ID: true, got[1].ID: true}
	if !ids[d1.ID] || !ids[d2.ID] {
		t.Errorf("IDs devueltos = %v, se esperaban %d y %d", ids, d1.ID, d2.ID)
	}
}

func TestRepositorioDefectos_ListarDefectos_ProyectoSinDefectos_DevuelveVacio(t *testing.T) {
	b := nuevaBaseDefecto(t)

	got, err := b.repo.ListarDefectos(context.Background(), b.proyectoID)

	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("se esperaban 0 defectos, hay %d", len(got))
	}
}
