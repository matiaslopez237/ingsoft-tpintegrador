package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// calidadSprintsFalso implementa los tres puertos que necesita ConsultarCalidadSprints
// (LectorProyectos, RepositorioSprints y LectorDefectos) en memoria.
type calidadSprintsFalso struct {
	proyectos   map[int64]domain.Proyecto
	sprints     []domain.Sprint
	defectos    []domain.Defecto
	errBuscar   error
	errSprints  error
	errDefectos error
}

func nuevoCalidadSprintsFalso(proyectos ...domain.Proyecto) *calidadSprintsFalso {
	m := map[int64]domain.Proyecto{}
	for _, p := range proyectos {
		m[p.ID] = p
	}
	return &calidadSprintsFalso{proyectos: m}
}

func (f *calidadSprintsFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errBuscar != nil {
		return domain.Proyecto{}, f.errBuscar
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *calidadSprintsFalso) ListarSprints(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
	if f.errSprints != nil {
		return nil, f.errSprints
	}
	var del []domain.Sprint
	for _, s := range f.sprints {
		if s.ProyectoID == proyectoID {
			del = append(del, s)
		}
	}
	return del, nil
}

// GuardarSprint no lo usa ConsultarCalidadSprints (es de solo lectura); existe para cumplir
// RepositorioSprints.
func (f *calidadSprintsFalso) GuardarSprint(_ context.Context, _ *domain.Sprint) error {
	return errors.New("GuardarSprint no se esperaba en estos tests")
}

func (f *calidadSprintsFalso) ListarDefectos(_ context.Context, proyectoID int64) ([]domain.Defecto, error) {
	if f.errDefectos != nil {
		return nil, f.errDefectos
	}
	var del []domain.Defecto
	for _, d := range f.defectos {
		if d.ProyectoID == proyectoID {
			del = append(del, d)
		}
	}
	return del, nil
}

func calidadCon(f *calidadSprintsFalso) *ConsultarCalidadSprints {
	return NuevoConsultarCalidadSprints(f, f, f)
}

// Traza: HU-24 CA-24.1, CA-24.6: solo los sprints cerrados, con su porcentaje.
func TestConsultarCalidadSprints_SprintsCerrados_InformaElPorcentaje(t *testing.T) {
	f := nuevoCalidadSprintsFalso(proyectoActivo(1))
	cerrado := sprintCerradoConPuntos(1, 1, 1, 6)
	cerrado.PuntosPlanificados = 8
	f.sprints = []domain.Sprint{
		cerrado,
		sprintGuardado(2, 1, 2, 26, domain.SprintAbierto), // no aparece
	}

	got, err := calidadCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("se esperaba 1 sprint, hay %d", len(got))
	}
	if !got[0].HayPorcentaje || got[0].PorcentajeCompletado != 75 {
		t.Errorf("sprint = %+v, se esperaba 75%% con datos", got[0])
	}
}

// Traza: HU-24 CA-24.2: un sprint cerrado sin puntos planificados no tiene porcentaje.
func TestConsultarCalidadSprints_SinPuntosPlanificados_SinPorcentaje(t *testing.T) {
	f := nuevoCalidadSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{sprintCerradoConPuntos(1, 1, 1, 0)}

	got, err := calidadCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 1 || got[0].HayPorcentaje {
		t.Errorf("sprint = %+v, se esperaba sin datos de porcentaje", got)
	}
}

// Traza: HU-24 CA-24.3, CA-24.4, CA-24.5: defectos detectados y resueltos, por sprint.
func TestConsultarCalidadSprints_CuentaLosDefectosPorSprint(t *testing.T) {
	f := nuevoCalidadSprintsFalso(proyectoActivo(1))
	sprint1 := sprintCerradoConPuntos(1, 1, 1, 5)
	sprint1.PuntosPlanificados = 5
	sprint2 := sprintCerradoConPuntos(2, 1, 2, 5)
	sprint2.PuntosPlanificados = 5
	f.sprints = []domain.Sprint{sprint1, sprint2}
	f.defectos = []domain.Defecto{
		{ProyectoID: 1, SprintDeteccionID: idInt64(1)},                                 // detectado en el sprint 1
		{ProyectoID: 1, SprintDeteccionID: idInt64(2)},                                 // detectado en el sprint 2
		{ProyectoID: 1, Estado: domain.DefectoCerrado, SprintResolucionID: idInt64(1)}, // resuelto en el sprint 1
		{ProyectoID: 1}, // sin sprint: no cuenta en ninguno
	}

	got, err := calidadCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	porSprint := map[int64]CalidadSprint{}
	for _, c := range got {
		porSprint[c.Sprint.ID] = c
	}
	if d := porSprint[1]; d.DefectosDetectados != 1 || d.DefectosResueltos != 1 {
		t.Errorf("sprint 1 = %+v, se esperaban 1 detectado y 1 resuelto", d)
	}
	if d := porSprint[2]; d.DefectosDetectados != 1 || d.DefectosResueltos != 0 {
		t.Errorf("sprint 2 = %+v, se esperaban 1 detectado y 0 resueltos", d)
	}
}

// Traza: HU-24 CA-24.7: con el proyecto inexistente no se leen sus sprints ni sus defectos.
func TestConsultarCalidadSprints_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoCalidadSprintsFalso()

	_, err := calidadCon(f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}

// Traza: HU-24 sección 7: una falla al leer se propaga, sin perder la causa.
func TestConsultarCalidadSprints_FallaDelRepositorio_SePropaga(t *testing.T) {
	disco := errors.New("disco lleno")
	tests := []struct {
		name  string
		falla func(*calidadSprintsFalso)
	}{
		{"al buscar el proyecto", func(f *calidadSprintsFalso) { f.errBuscar = disco }},
		{"al listar los sprints", func(f *calidadSprintsFalso) { f.errSprints = disco }},
		{"al listar los defectos", func(f *calidadSprintsFalso) { f.errDefectos = disco }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoCalidadSprintsFalso(proyectoActivo(1))
			tt.falla(f)

			_, err := calidadCon(f).Ejecutar(context.Background(), 1)

			if !errors.Is(err, disco) {
				t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
			}
		})
	}
}

// idInt64 arma un puntero al ID de sprint indicado, para los defectos de prueba.
func idInt64(id int64) *int64 { return &id }
