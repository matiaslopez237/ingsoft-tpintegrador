package app

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// lectorFalso implementa LectorProyectos en memoria. Solo tiene BuscarPorID: si el caso de uso
// intentara escribir, no compilaría (HU-04, sección 5: solo lectura).
type lectorFalso struct {
	porID     map[int64]domain.Proyecto
	errBuscar error
}

func (r *lectorFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if r.errBuscar != nil {
		return domain.Proyecto{}, r.errBuscar
	}
	p, ok := r.porID[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

// proyectoConsultado tiene todos los datos que devuelve la consulta: fechas, estado e
// integrantes con roles, con el orden de alta distinto del alfabético.
func proyectoConsultado(id int64) domain.Proyecto {
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	return domain.Proyecto{
		ID:          id,
		Nombre:      "TP Integrador",
		FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    &fin,
		Estado:      domain.EstadoActivo,
		Integrantes: []domain.Integrante{
			{Nombre: "Beto", Roles: []string{"Tester"}},
			{Nombre: "Ana", Roles: []string{"Tester", "Desarrolladora"}},
		},
	}
}

// Traza: HU-04 CA-04.1, CA-04.2
func TestConsultarProyecto_Existente_DevuelveDatosEIntegrantes(t *testing.T) {
	repo := &lectorFalso{porID: map[int64]domain.Proyecto{
		7: proyectoConsultado(7),
		8: {ID: 8, Nombre: "Otro proyecto"},
	}}
	uc := NuevoConsultarProyecto(repo)

	got, err := uc.Ejecutar(context.Background(), 7)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := proyectoConsultado(7); !reflect.DeepEqual(got, want) {
		t.Errorf("proyecto = %+v, se esperaba %+v", got, want)
	}
}

// Traza: HU-04 CA-04.5
func TestConsultarProyecto_Inexistente_RetornaErrProyectoNoEncontrado(t *testing.T) {
	repo := &lectorFalso{porID: map[int64]domain.Proyecto{7: proyectoConsultado(7)}}
	uc := NuevoConsultarProyecto(repo)

	got, err := uc.Ejecutar(context.Background(), 9999)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if !reflect.DeepEqual(got, domain.Proyecto{}) {
		t.Errorf("proyecto = %+v, no se esperaba ningún dato", got)
	}
}

// Traza: HU-04 sección 7 (falla interna); no tiene un CA propio en la spec.
func TestConsultarProyecto_FallaDelRepositorio_SeEnvuelveYNoEsNoEncontrado(t *testing.T) {
	errBase := errors.New("base caída")
	uc := NuevoConsultarProyecto(&lectorFalso{errBuscar: errBase})

	_, err := uc.Ejecutar(context.Background(), 7)

	if !errors.Is(err, errBase) {
		t.Fatalf("error = %v, se esperaba que envuelva %v", err, errBase)
	}
	if err == errBase {
		t.Errorf("error = %v, se esperaba envuelto con contexto (%%w), no el error original", err)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, una falla interna no debe confundirse con %v", err, domain.ErrProyectoNoEncontrado)
	}
}
