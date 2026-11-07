package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// reporteFalso implementa los puertos que necesitan ArmarReporteProyecto y ArmarReporteSprint en
// memoria.
type reporteFalso struct {
	proyectos     map[int64]domain.Proyecto
	sprints       []domain.Sprint
	elementos     []domain.ElementoBacklog
	horasProyecto []domain.HorasElemento
	horasSprint   []domain.HorasElemento
	defectos      []domain.Defecto

	errProyecto, errSprints, errBuscarSprint, errElementos, errHoras, errDefectos error
}

func nuevoReporteFalso(proyectos ...domain.Proyecto) *reporteFalso {
	m := map[int64]domain.Proyecto{}
	for _, p := range proyectos {
		m[p.ID] = p
	}
	return &reporteFalso{proyectos: m}
}

func (f *reporteFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *reporteFalso) ListarSprints(_ context.Context, _ int64) ([]domain.Sprint, error) {
	if f.errSprints != nil {
		return nil, f.errSprints
	}
	return f.sprints, nil
}

func (f *reporteFalso) GuardarSprint(_ context.Context, _ *domain.Sprint) error {
	return errors.New("GuardarSprint no se esperaba en estos tests")
}

func (f *reporteFalso) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	if f.errBuscarSprint != nil {
		return domain.Sprint{}, f.errBuscarSprint
	}
	for _, s := range f.sprints {
		if s.ID == id {
			return s, nil
		}
	}
	return domain.Sprint{}, domain.ErrSprintNoEncontrado
}

func (f *reporteFalso) ListarElementos(_ context.Context, _ int64) ([]domain.ElementoBacklog, error) {
	if f.errElementos != nil {
		return nil, f.errElementos
	}
	return f.elementos, nil
}

func (f *reporteFalso) HorasDeProyecto(_ context.Context, _ int64) ([]domain.HorasElemento, error) {
	if f.errHoras != nil {
		return nil, f.errHoras
	}
	return f.horasProyecto, nil
}

func (f *reporteFalso) HorasDeSprint(_ context.Context, _ int64) ([]domain.HorasElemento, error) {
	if f.errHoras != nil {
		return nil, f.errHoras
	}
	return f.horasSprint, nil
}

func (f *reporteFalso) ListarDefectos(_ context.Context, _ int64) ([]domain.Defecto, error) {
	if f.errDefectos != nil {
		return nil, f.errDefectos
	}
	return f.defectos, nil
}

func reporteProyectoCon(f *reporteFalso) *ArmarReporteProyecto {
	return NuevoArmarReporteProyecto(f, f, f, f, f)
}

func reporteSprintCon(f *reporteFalso) *ArmarReporteSprint {
	return NuevoArmarReporteSprint(f, f, f, f, f)
}

// Traza: HU-26 CA-26.1
func TestArmarReporteProyecto_TraeTodasLasHistorias(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))
	f.elementos = []domain.ElementoBacklog{
		{ProyectoID: 1, Titulo: "Login", Estado: domain.EstadoHecho},
		{ProyectoID: 1, Titulo: "Logout", Estado: domain.EstadoPendiente},
	}

	r, err := reporteProyectoCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(r.Historias) != 2 {
		t.Errorf("Historias = %+v, se esperaban 2", r.Historias)
	}
	if r.Sprint != nil {
		t.Errorf("Sprint = %+v, se esperaba nil (reporte de proyecto)", r.Sprint)
	}
}

// Traza: HU-26 CA-26.5, CA-26.6
func TestArmarReporteProyecto_Metricas(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{
		sprintCerradoConPuntos(1, 1, 1, 6),
		sprintGuardado(2, 1, 2, 26, domain.SprintAbierto), // no aparece
	}
	f.sprints[0].PuntosPlanificados = 8

	r, err := reporteProyectoCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !r.HayVelocidad || r.Velocidad != 6 {
		t.Errorf("velocidad = %v/%v, se esperaba 6 con datos", r.Velocidad, r.HayVelocidad)
	}
	if len(r.Metricas) != 1 || !r.Metricas[0].HayPorcentaje || r.Metricas[0].PorcentajeCompletado != 75 {
		t.Errorf("Metricas = %+v, se esperaba 1 fila con 75%%", r.Metricas)
	}
}

// Traza: HU-26 CA-26.6: sin sprints cerrados, sin velocidad ni filas.
func TestArmarReporteProyecto_SinSprintsCerrados_SinMetricas(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))

	r, err := reporteProyectoCon(f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.HayVelocidad || len(r.Metricas) != 0 {
		t.Errorf("r = %+v, se esperaba sin velocidad y sin métricas", r)
	}
}

// Traza: HU-26 sección 7
func TestArmarReporteProyecto_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoReporteFalso()

	_, err := reporteProyectoCon(f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}

// Traza: HU-26 sección 7: una falla al leer se propaga, sin perder la causa.
func TestArmarReporteProyecto_FallaDelRepositorio_SePropaga(t *testing.T) {
	disco := errors.New("disco lleno")
	tests := []struct {
		name  string
		falla func(*reporteFalso)
	}{
		{"al buscar el proyecto", func(f *reporteFalso) { f.errProyecto = disco }},
		{"al listar los elementos", func(f *reporteFalso) { f.errElementos = disco }},
		{"al listar los sprints", func(f *reporteFalso) { f.errSprints = disco }},
		{"al leer las horas", func(f *reporteFalso) { f.errHoras = disco }},
		{"al listar los defectos", func(f *reporteFalso) { f.errDefectos = disco }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoReporteFalso(proyectoActivo(1))
			tt.falla(f)

			_, err := reporteProyectoCon(f).Ejecutar(context.Background(), 1)

			if !errors.Is(err, disco) {
				t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
			}
		})
	}
}

// Traza: HU-26 CA-26.2
func TestArmarReporteSprint_SoloHistoriasAsignadas(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))
	sid := int64(10)
	f.sprints = []domain.Sprint{{ID: 10, ProyectoID: 1, Numero: 1, Estado: domain.SprintCerrado, PuntosPlanificados: 5, PuntosCompletados: 5}}
	f.elementos = []domain.ElementoBacklog{
		{ProyectoID: 1, Titulo: "Login", SprintID: &sid},
		{ProyectoID: 1, Titulo: "Sin sprint"},
	}

	r, err := reporteSprintCon(f).Ejecutar(context.Background(), 1, 10)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(r.Historias) != 1 || r.Historias[0].Titulo != "Login" {
		t.Errorf("Historias = %+v, se esperaba solo \"Login\"", r.Historias)
	}
	if r.Sprint == nil || r.Sprint.ID != 10 {
		t.Errorf("Sprint = %+v, se esperaba el 10", r.Sprint)
	}
}

// Traza: HU-26 CA-26.7, CA-26.8
func TestArmarReporteSprint_AbiertoYCerrado(t *testing.T) {
	tests := []struct {
		name              string
		sprint            domain.Sprint
		wantHayPorcentaje bool
		wantPorcentaje    float64
	}{
		{"CA-26.7: cerrado con puntos, porcentaje real", domain.Sprint{ID: 1, ProyectoID: 1, PuntosPlanificados: 4, PuntosCompletados: 3, Estado: domain.SprintCerrado}, true, 75},
		{"CA-26.8: abierto, sin datos", domain.Sprint{ID: 1, ProyectoID: 1, Estado: domain.SprintAbierto}, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoReporteFalso(proyectoActivo(1))
			f.sprints = []domain.Sprint{tt.sprint}

			r, err := reporteSprintCon(f).Ejecutar(context.Background(), 1, 1)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if r.Metricas[0].HayPorcentaje != tt.wantHayPorcentaje {
				t.Fatalf("HayPorcentaje = %v, se esperaba %v", r.Metricas[0].HayPorcentaje, tt.wantHayPorcentaje)
			}
			if tt.wantHayPorcentaje && r.Metricas[0].PorcentajeCompletado != tt.wantPorcentaje {
				t.Errorf("Porcentaje = %v, se esperaba %v", r.Metricas[0].PorcentajeCompletado, tt.wantPorcentaje)
			}
		})
	}
}

// Traza: HU-26 CA-26.9, CA-26.10: defectos detectados o resueltos en el sprint, sin importar otros.
func TestArmarReporteSprint_Defectos(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))
	uno, dos := int64(1), int64(2)
	f.sprints = []domain.Sprint{{ID: 1, ProyectoID: 1, Estado: domain.SprintCerrado}}
	f.defectos = []domain.Defecto{
		{ID: 100, SprintDeteccionID: &uno},
		{ID: 101, SprintDeteccionID: &dos},
	}

	r, err := reporteSprintCon(f).Ejecutar(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(r.Defectos) != 1 || r.Defectos[0].ID != 100 {
		t.Errorf("Defectos = %+v, se esperaba solo el 100", r.Defectos)
	}
	if r.Metricas[0].DefectosDetectados != 1 {
		t.Errorf("DefectosDetectados = %d, se esperaba 1", r.Metricas[0].DefectosDetectados)
	}
}

// Traza: HU-26 sección 7
func TestArmarReporteSprint_Errores(t *testing.T) {
	f := nuevoReporteFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{{ID: 1, ProyectoID: 2}} // de otro proyecto

	tests := []struct {
		name     string
		sprintID int64
		proyecto int64
		wantErr  error
	}{
		{"proyecto inexistente", 1, 99, domain.ErrProyectoNoEncontrado},
		{"sprint inexistente", 999, 1, domain.ErrSprintNoEncontrado},
		{"sprint de otro proyecto", 1, 1, domain.ErrSprintNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := reporteSprintCon(f).Ejecutar(context.Background(), tt.proyecto, tt.sprintID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}
