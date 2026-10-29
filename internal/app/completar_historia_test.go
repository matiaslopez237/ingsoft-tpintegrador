package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// completarFalso implementa LectorProyectos y RepositorioCompletar en memoria y permite simular
// fallas. Proyecto 100 con el elemento 1 en el sprint abierto 10; proyecto 200 sin elementos.
type completarFalso struct {
	proyectos        map[int64]domain.Proyecto
	elementos        map[int64]domain.ElementoBacklog
	sprints          map[int64]sprintFalso
	errProyecto      error
	errElemento      error
	errSprint        error
	errGuardar       error
	vecesBuscoSprint int
	vecesGuardo      int
}

func nuevoCompletarFalso() *completarFalso {
	sprintID := int64(10)
	return &completarFalso{
		proyectos: map[int64]domain.Proyecto{100: proyectoFalso(100), 200: proyectoFalso(200)},
		elementos: map[int64]domain.ElementoBacklog{1: elementoFalso(1, 100, &sprintID)},
		sprints:   map[int64]sprintFalso{10: {proyectoID: 100, abierto: true}},
	}
}

func (f *completarFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *completarFalso) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if f.errElemento != nil {
		return domain.ElementoBacklog{}, f.errElemento
	}
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *completarFalso) BuscarSprint(_ context.Context, id int64) (proyectoID int64, abierto bool, err error) {
	f.vecesBuscoSprint++
	if f.errSprint != nil {
		return 0, false, f.errSprint
	}
	s, ok := f.sprints[id]
	if !ok {
		return 0, false, domain.ErrSprintNoEncontrado
	}
	return s.proyectoID, s.abierto, nil
}

func (f *completarFalso) GuardarCompletado(_ context.Context, e *domain.ElementoBacklog) error {
	if f.errGuardar != nil {
		return f.errGuardar
	}
	f.vecesGuardo++
	f.elementos[e.ID] = *e
	return nil
}

func completarCon(f *completarFalso) *CompletarHistoria {
	return NuevoCompletarHistoria(f, f)
}

// Traza: HU-10 CA-10.1, CA-10.2 (RN-10.1 a RN-10.3)
func TestCompletarHistoria_SprintAbierto_GuardaComoHecho(t *testing.T) {
	tests := []struct {
		nombre string
		estado domain.EstadoElemento
	}{
		{"desde Pendiente", domain.EstadoPendiente},
		{"desde En progreso", domain.EstadoEnProgreso},
		{"ya Hecho, sin error (RN-10.3)", domain.EstadoHecho},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoCompletarFalso()
			e := f.elementos[1]
			e.Estado = tt.estado
			f.elementos[1] = e

			err := completarCon(f).Ejecutar(context.Background(), DatosCompletar{ProyectoID: 100, ElementoID: 1})

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			guardado := f.elementos[1]
			if guardado.Estado != domain.EstadoHecho {
				t.Errorf("Estado = %q, se esperaba %q", guardado.Estado, domain.EstadoHecho)
			}
			if guardado.SprintID == nil || *guardado.SprintID != 10 {
				t.Errorf("SprintID = %v, se esperaba que siga en 10", guardado.SprintID)
			}
			if f.vecesGuardo != 1 {
				t.Errorf("se guardó %d veces, se esperaba 1", f.vecesGuardo)
			}
		})
	}
}

// Traza: HU-10 CA-10.3: sin sprint no hace falta consultar ningún sprint.
func TestCompletarHistoria_SinSprint_RetornaErrorYNoGuarda(t *testing.T) {
	f := nuevoCompletarFalso()
	f.elementos[1] = elementoFalso(1, 100, nil)

	err := completarCon(f).Ejecutar(context.Background(), DatosCompletar{ProyectoID: 100, ElementoID: 1})

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
	if f.vecesBuscoSprint != 0 {
		t.Errorf("se consultó el sprint %d veces, se esperaba 0", f.vecesBuscoSprint)
	}
	if f.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", f.vecesGuardo)
	}
}

// Traza: HU-10 CA-10.4 (P-8: también si ya está Hecho)
func TestCompletarHistoria_SprintCerrado_RetornaErrorYNoGuarda(t *testing.T) {
	for _, estado := range []domain.EstadoElemento{domain.EstadoEnProgreso, domain.EstadoHecho} {
		t.Run(string(estado), func(t *testing.T) {
			f := nuevoCompletarFalso()
			f.sprints[10] = sprintFalso{proyectoID: 100, abierto: false}
			e := f.elementos[1]
			e.Estado = estado
			f.elementos[1] = e

			err := completarCon(f).Ejecutar(context.Background(), DatosCompletar{ProyectoID: 100, ElementoID: 1})

			if !errors.Is(err, domain.ErrSprintCerrado) {
				t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
			}
			if f.vecesGuardo != 0 {
				t.Errorf("no se esperaba ningún guardado, hubo %d", f.vecesGuardo)
			}
			if got := f.elementos[1].Estado; got != estado {
				t.Errorf("Estado = %q, no debería haber cambiado (%q)", got, estado)
			}
		})
	}
}

// Traza: HU-10 CA-10.5 (RN-10.4), sección 7: el primer error en el orden de la tabla.
func TestCompletarHistoria_ProyectoOElementoInexistente_RetornaErrorYNoGuarda(t *testing.T) {
	tests := []struct {
		nombre      string
		datos       DatosCompletar
		errEsperado error
	}{
		{"proyecto inexistente", DatosCompletar{ProyectoID: 999, ElementoID: 1}, domain.ErrProyectoNoEncontrado},
		{"proyecto y elemento inexistentes: primero el proyecto", DatosCompletar{ProyectoID: 999, ElementoID: 99}, domain.ErrProyectoNoEncontrado},
		{"elemento inexistente", DatosCompletar{ProyectoID: 100, ElementoID: 99}, domain.ErrElementoNoEncontrado},
		{"elemento de otro proyecto (RN-10.4)", DatosCompletar{ProyectoID: 200, ElementoID: 1}, domain.ErrElementoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoCompletarFalso()

			err := completarCon(f).Ejecutar(context.Background(), tt.datos)

			if !errors.Is(err, tt.errEsperado) {
				t.Errorf("error = %v, se esperaba %v", err, tt.errEsperado)
			}
			if f.vecesGuardo != 0 {
				t.Errorf("no se esperaba ningún guardado, hubo %d", f.vecesGuardo)
			}
		})
	}
}

func TestCompletarHistoria_FallaDelRepositorio_SePropagaYNoGuarda(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre  string
		preparo func(f *completarFalso)
	}{
		{"al buscar el proyecto", func(f *completarFalso) { f.errProyecto = errBase }},
		{"al buscar el elemento", func(f *completarFalso) { f.errElemento = errBase }},
		{"al buscar el sprint", func(f *completarFalso) { f.errSprint = errBase }},
		{"al guardar", func(f *completarFalso) { f.errGuardar = errBase }},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoCompletarFalso()
			tt.preparo(f)

			err := completarCon(f).Ejecutar(context.Background(), DatosCompletar{ProyectoID: 100, ElementoID: 1})

			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
			if f.vecesGuardo != 0 {
				t.Errorf("no se esperaba ningún guardado, hubo %d", f.vecesGuardo)
			}
		})
	}
}
