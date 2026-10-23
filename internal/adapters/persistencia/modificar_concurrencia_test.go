package persistencia

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// pedidoConLecturaVieja simula un pedido que leyó el proyecto antes de que otro guardara: su
// primera lectura devuelve la foto vieja y las siguientes, lo que hay en la base. Va uno por pedido.
type pedidoConLecturaVieja struct {
	*RepositorioProyectos
	foto   domain.Proyecto
	yaLeyo bool
}

func (r *pedidoConLecturaVieja) BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error) {
	if r.yaLeyo {
		return r.RepositorioProyectos.BuscarPorID(ctx, id)
	}
	r.yaLeyo = true
	return r.foto, nil
}

func fechaDePrueba(s string) *time.Time {
	f, _ := time.Parse(time.DateOnly, s)
	return &f
}

// Traza: HU-02 sección 5 (pedidos simultáneos): dos modificaciones de campos distintos no se pisan.
func TestModificar_CamposDistintosSimultaneos_SeConservanAmbos(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, repo, "Alfa")
	finalizado := domain.EstadoFinalizado

	if _, err := app.NuevoModificarProyecto(&pedidoConLecturaVieja{RepositorioProyectos: repo, foto: p}).Ejecutar(ctx, app.DatosModificarProyecto{
		ProyectoID: p.ID, Cambios: domain.CambiosProyecto{Estado: &finalizado},
	}); err != nil {
		t.Fatalf("cambiar estado: %v", err)
	}
	devuelto, err := app.NuevoModificarProyecto(&pedidoConLecturaVieja{RepositorioProyectos: repo, foto: p}).Ejecutar(ctx, app.DatosModificarProyecto{
		ProyectoID: p.ID, Cambios: domain.CambiosProyecto{CambiarFechaFin: true, FechaFin: fechaDePrueba("2026-12-15")},
	})
	if err != nil {
		t.Fatalf("cambiar fecha de fin: %v", err)
	}
	// Sección 3: se devuelve el proyecto como quedó guardado, no la foto vieja con el cambio encima.
	if devuelto.Estado != domain.EstadoFinalizado {
		t.Errorf("proyecto devuelto con estado %q, se esperaba el guardado (Finalizado)", devuelto.Estado)
	}

	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.Estado != domain.EstadoFinalizado || got.FechaFin == nil || got.FechaFin.Format(time.DateOnly) != "2026-12-15" {
		t.Errorf("guardado = estado %q, fin %v; se esperaba Finalizado y 2026-12-15 (un cambio pisó al otro)", got.Estado, got.FechaFin)
	}
}

// Traza: HU-02 sección 5 y RN-02.5: si dos pedidos cambian una fecha cada uno y juntos dejarían la
// fecha de fin antes que la de inicio, el segundo se rechaza al guardar y no cambia nada.
func TestModificar_FechasCruzadasSimultaneas_SegundoEsRechazado(t *testing.T) {
	ctx := context.Background()
	repo := nuevoRepoDePrueba(t)
	p := guardarDePrueba(t, repo, "Alfa") // inicio 2026-03-01, fin 2026-11-30

	if _, err := app.NuevoModificarProyecto(&pedidoConLecturaVieja{RepositorioProyectos: repo, foto: p}).Ejecutar(ctx, app.DatosModificarProyecto{
		ProyectoID: p.ID, Cambios: domain.CambiosProyecto{FechaInicio: fechaDePrueba("2026-06-01")},
	}); err != nil {
		t.Fatalf("cambiar inicio: %v", err)
	}
	// Contra la foto vieja (inicio 2026-03-01) este fin es válido, pero contra lo guardado no.
	_, err := app.NuevoModificarProyecto(&pedidoConLecturaVieja{RepositorioProyectos: repo, foto: p}).Ejecutar(ctx, app.DatosModificarProyecto{
		ProyectoID: p.ID, Cambios: domain.CambiosProyecto{CambiarFechaFin: true, FechaFin: fechaDePrueba("2026-05-01")},
	})

	if !errors.Is(err, domain.ErrFechaFinAnterior) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrFechaFinAnterior)
	}
	got, err := repo.BuscarPorID(ctx, p.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if got.FechaFin == nil || got.FechaFin.Format(time.DateOnly) != "2026-11-30" {
		t.Errorf("fecha de fin guardada = %v, no debería haber cambiado", got.FechaFin)
	}
}
