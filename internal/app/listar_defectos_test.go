package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// listadoDefectosFalso implementa los dos puertos de ListarDefectos en memoria.
type listadoDefectosFalso struct {
	proyectos   map[int64]domain.Proyecto
	defectos    map[int64][]domain.Defecto // por proyecto
	errProyecto error
	errListar   error
}

func nuevoListadoDefectosFalso() *listadoDefectosFalso {
	return &listadoDefectosFalso{
		proyectos: map[int64]domain.Proyecto{1: proyectoFalso(1)},
		defectos:  map[int64][]domain.Defecto{},
	}
}

func (f *listadoDefectosFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *listadoDefectosFalso) ListarDefectos(_ context.Context, proyectoID int64) ([]domain.Defecto, error) {
	if f.errListar != nil {
		return nil, f.errListar
	}
	return f.defectos[proyectoID], nil
}

func listarDefectosCon(f *listadoDefectosFalso) *ListarDefectos {
	return NuevoListarDefectos(f, f)
}

// Traza: HU-21 CA-21.8
func TestListarDefectos_ProyectoConDefectos_LosDevuelve(t *testing.T) {
	f := nuevoListadoDefectosFalso()
	f.defectos[1] = []domain.Defecto{
		{ID: 100, ProyectoID: 1, Estado: domain.DefectoAbierto},
		{ID: 101, ProyectoID: 1, Estado: domain.DefectoResuelto},
	}

	got, err := listarDefectosCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("se esperaban 2 defectos, hay %d", len(got))
	}
}

// Traza: HU-21 CA-21.8 (lista vacía si no tiene ninguno)
func TestListarDefectos_ProyectoSinDefectos_DevuelveVacio(t *testing.T) {
	f := nuevoListadoDefectosFalso()

	got, err := listarDefectosCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("se esperaban 0 defectos, hay %d", len(got))
	}
}

// Traza: HU-21 CA-21.8
func TestListarDefectos_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoListadoDefectosFalso()

	_, err := listarDefectosCon(f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}
