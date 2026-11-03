package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// horasFalsas implementa los puertos de los dos casos de uso en memoria, permite simular fallas y
// registra las consultas. Proyecto 1: el sprint 10. Proyecto 2: el sprint 20. Las horas del proyecto y
// las del sprint son distintas, para que se note si el caso de uso lee las equivocadas.
type horasFalsas struct {
	proyectos         map[int64]domain.Proyecto
	sprints           map[int64]domain.Sprint
	horasPorProyecto  map[int64][]domain.HorasElemento
	horasPorSprint    map[int64][]domain.HorasElemento
	errProyecto       error
	errSprint         error
	errHoras          error
	vecesSprint       int
	consultasProyecto []int64
	consultasSprint   []int64
	ctxHoras          context.Context // el contexto que recibió el lector de horas
}

func nuevasHorasFalsas() *horasFalsas {
	return &horasFalsas{
		proyectos: map[int64]domain.Proyecto{
			1: {ID: 1, Nombre: "Alfa", Estado: domain.EstadoActivo},
			2: {ID: 2, Nombre: "Beta", Estado: domain.EstadoActivo},
		},
		sprints: map[int64]domain.Sprint{
			10: {ID: 10, ProyectoID: 1, Numero: 1, Estado: domain.SprintAbierto},
			20: {ID: 20, ProyectoID: 2, Numero: 1, Estado: domain.SprintAbierto},
		},
		horasPorProyecto: map[int64][]domain.HorasElemento{
			1: {{HorasEstimadas: pHoras(8), HorasReales: 10}, {HorasEstimadas: pHoras(4), HorasReales: 3}, {HorasEstimadas: pHoras(20), HorasReales: 15}},
		},
		horasPorSprint: map[int64][]domain.HorasElemento{
			10: {{HorasEstimadas: pHoras(8), HorasReales: 10}, {HorasEstimadas: pHoras(4), HorasReales: 3}},
			20: {{HorasEstimadas: pHoras(6), HorasReales: 7}},
		},
	}
}

func (f *horasFalsas) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *horasFalsas) BuscarSprintPorID(_ context.Context, id int64) (domain.Sprint, error) {
	f.vecesSprint++
	if f.errSprint != nil {
		return domain.Sprint{}, f.errSprint
	}
	s, ok := f.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return s, nil
}

func (f *horasFalsas) HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error) {
	f.ctxHoras = ctx
	f.consultasProyecto = append(f.consultasProyecto, proyectoID)
	if f.errHoras != nil {
		return nil, f.errHoras
	}
	return f.horasPorProyecto[proyectoID], nil
}

func (f *horasFalsas) HorasDeSprint(ctx context.Context, sprintID int64) ([]domain.HorasElemento, error) {
	f.ctxHoras = ctx
	f.consultasSprint = append(f.consultasSprint, sprintID)
	if f.errHoras != nil {
		return nil, f.errHoras
	}
	return f.horasPorSprint[sprintID], nil
}

func horasDeProyecto(f *horasFalsas, proyectoID int64) (domain.TotalHoras, error) {
	return NuevoConsultarHorasProyecto(f, f).Ejecutar(context.Background(), proyectoID)
}

func horasDeSprint(f *horasFalsas, proyectoID, sprintID int64) (domain.TotalHoras, error) {
	return NuevoConsultarHorasSprint(f, f, f).Ejecutar(context.Background(), proyectoID, sprintID)
}

// Traza: HU-23 CA-23.1 (RN-23.1, RN-23.4, RN-23.7; sección 5: casos de uso)
func TestConsultarHorasSprint_DevuelveElTotalDeLosElementosDelSprint(t *testing.T) {
	f := nuevasHorasFalsas()

	got, err := horasDeSprint(f, 1, 10)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := domain.TotalHoras{
		HorasEstimadas: 12, HorasReales: 13, DesviacionHoras: pHoras(1), DesviacionPorcentaje: pHoras(8.3), Elementos: 2,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", got, want)
	}
	if !reflect.DeepEqual(f.consultasSprint, []int64{10}) || len(f.consultasProyecto) != 0 {
		t.Errorf("horas consultadas: sprints %v y proyectos %v; se esperaba solo el sprint 10", f.consultasSprint, f.consultasProyecto)
	}
}

// Traza: HU-23 CA-23.2 (RN-23.1, RN-23.10; sección 5: casos de uso)
func TestConsultarHorasProyecto_DevuelveElTotalDeLosElementosDelProyecto(t *testing.T) {
	f := nuevasHorasFalsas()

	got, err := horasDeProyecto(f, 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := domain.TotalHoras{
		HorasEstimadas: 32, HorasReales: 28, DesviacionHoras: pHoras(-4), DesviacionPorcentaje: pHoras(-12.5), Elementos: 3,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", got, want)
	}
	if !reflect.DeepEqual(f.consultasProyecto, []int64{1}) || len(f.consultasSprint) != 0 {
		t.Errorf("horas consultadas: proyectos %v y sprints %v; se esperaba solo el proyecto 1", f.consultasProyecto, f.consultasSprint)
	}
	if f.vecesSprint != 0 {
		t.Errorf("se buscó un sprint %d veces, no se esperaba ninguna", f.vecesSprint)
	}
}

// El lector de horas recibe el contexto del pedido, para que cancelarlo corte la consulta.
// Traza: HU-23 (sección 5: contexto)
func TestConsultarHoras_ElLectorRecibeElContextoDelPedido(t *testing.T) {
	type clave struct{}
	ctx := context.WithValue(context.Background(), clave{}, "pedido")
	tests := []struct {
		nombre    string
		consultar func(f *horasFalsas) (domain.TotalHoras, error)
	}{
		{"sprint", func(f *horasFalsas) (domain.TotalHoras, error) {
			return NuevoConsultarHorasSprint(f, f, f).Ejecutar(ctx, 1, 10)
		}},
		{"proyecto", func(f *horasFalsas) (domain.TotalHoras, error) {
			return NuevoConsultarHorasProyecto(f, f).Ejecutar(ctx, 1)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevasHorasFalsas()

			if _, err := tt.consultar(f); err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if f.ctxHoras != ctx {
				t.Errorf("el lector recibió %v, se esperaba el contexto del pedido", f.ctxHoras)
			}
		})
	}
}

// El caso de uso le deja al dominio los elementos sin estimar: no los filtra antes de totalizar.
// Traza: HU-23 CA-23.6 (RN-23.5)
func TestConsultarHoras_ConElementosSinEstimar_LosPasaAlDominio(t *testing.T) {
	filas := []domain.HorasElemento{{HorasEstimadas: pHoras(8), HorasReales: 10}, {HorasReales: 5}}
	want := domain.TotalHoras{
		HorasEstimadas: 8, HorasReales: 15, HorasRealesSinEstimar: 5,
		DesviacionHoras: pHoras(2), DesviacionPorcentaje: pHoras(25), Elementos: 2, ElementosSinEstimar: 1,
	}
	f := nuevasHorasFalsas()
	f.horasPorSprint[10] = filas
	f.horasPorProyecto[1] = filas

	deSprint, err := horasDeSprint(f, 1, 10)
	if err != nil {
		t.Fatalf("sprint: error inesperado: %v", err)
	}
	deProyecto, err := horasDeProyecto(f, 1)
	if err != nil {
		t.Fatalf("proyecto: error inesperado: %v", err)
	}

	if !reflect.DeepEqual(deSprint, want) {
		t.Errorf("sprint: se obtuvo %+v, se esperaba %+v", deSprint, want)
	}
	if !reflect.DeepEqual(deProyecto, want) {
		t.Errorf("proyecto: se obtuvo %+v, se esperaba %+v", deProyecto, want)
	}
}

// Sin elementos no es un error (RN-23.6): el lector devuelve una lista vacía.
// Traza: HU-23 CA-23.7, CA-23.8 (RN-23.6)
func TestConsultarHoras_SinElementos_NoEsUnError(t *testing.T) {
	f := nuevasHorasFalsas()
	f.horasPorSprint[10] = nil
	f.horasPorProyecto[1] = nil

	deSprint, err := horasDeSprint(f, 1, 10)
	if err != nil {
		t.Fatalf("sprint: error inesperado: %v", err)
	}
	deProyecto, err := horasDeProyecto(f, 1)
	if err != nil {
		t.Fatalf("proyecto: error inesperado: %v", err)
	}

	if !deSprint.SinEstimar() || !reflect.DeepEqual(deSprint, domain.TotalHoras{}) {
		t.Errorf("sprint: se obtuvo %+v, se esperaba un total vacío", deSprint)
	}
	if !deProyecto.SinEstimar() || !reflect.DeepEqual(deProyecto, domain.TotalHoras{}) {
		t.Errorf("proyecto: se obtuvo %+v, se esperaba un total vacío", deProyecto)
	}
}

// Traza: HU-23 CA-23.13 (RN-23.12)
func TestConsultarHoras_CualquierEstadoDeSprintYProyecto(t *testing.T) {
	for _, estadoProyecto := range []domain.Estado{domain.EstadoActivo, domain.EstadoEnPausa,
		domain.EstadoFinalizado, domain.EstadoPostergado, domain.EstadoCancelado} {
		for _, estadoSprint := range []domain.EstadoSprint{domain.SprintAbierto, domain.SprintCerrado} {
			f := nuevasHorasFalsas()
			p := f.proyectos[1]
			p.Estado = estadoProyecto
			f.proyectos[1] = p
			s := f.sprints[10]
			s.Estado = estadoSprint
			f.sprints[10] = s

			deSprint, err := horasDeSprint(f, 1, 10)
			if err != nil {
				t.Errorf("proyecto %s, sprint %s: error inesperado: %v", estadoProyecto, estadoSprint, err)
			} else if deSprint.Elementos != 2 {
				t.Errorf("proyecto %s, sprint %s: se obtuvo %+v", estadoProyecto, estadoSprint, deSprint)
			}
			deProyecto, err := horasDeProyecto(f, 1)
			if err != nil {
				t.Errorf("proyecto %s: error inesperado: %v", estadoProyecto, err)
			} else if deProyecto.Elementos != 3 {
				t.Errorf("proyecto %s: se obtuvo %+v", estadoProyecto, deProyecto)
			}
		}
	}
}

// Los errores de validación se devuelven tal cual (por eso == y no errors.Is), y sin consultar las
// horas: un sprint de otro proyecto da el mismo error que uno inexistente, sin revelar que existe.
// Traza: HU-23 CA-23.14, CA-23.15 (RN-23.11; sección 7: primero el proyecto y después el sprint)
func TestConsultarHorasSprint_ValidacionFallida_RetornaErrorYNoConsultaHoras(t *testing.T) {
	tests := []struct {
		nombre     string
		proyectoID int64
		sprintID   int64
		wantErr    error
	}{
		{"proyecto inexistente", 99, 10, domain.ErrProyectoNoEncontrado},
		{"sprint inexistente", 1, 99, domain.ErrSprintNoEncontrado},
		{"sprint de otro proyecto", 1, 20, domain.ErrSprintNoEncontrado},
		{"otro proyecto, visto desde el segundo", 2, 10, domain.ErrSprintNoEncontrado},
		{"proyecto y sprint inexistentes: primero el proyecto", 99, 99, domain.ErrProyectoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevasHorasFalsas()

			got, err := horasDeSprint(f, tt.proyectoID, tt.sprintID)

			if err != tt.wantErr {
				t.Errorf("error = %v, se esperaba exactamente %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, domain.TotalHoras{}) {
				t.Errorf("con error se devolvió %+v, se esperaba un total vacío", got)
			}
			if len(f.consultasSprint) != 0 || len(f.consultasProyecto) != 0 {
				t.Errorf("se consultaron horas: sprints %v y proyectos %v; no se esperaba ninguna consulta", f.consultasSprint, f.consultasProyecto)
			}
		})
	}
}

// Traza: HU-23 CA-23.14 (RN-23.11)
func TestConsultarHorasProyecto_ProyectoInexistente_RetornaErrorYNoConsultaHoras(t *testing.T) {
	f := nuevasHorasFalsas()

	got, err := horasDeProyecto(f, 99)

	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("error = %v, se esperaba exactamente %v", err, domain.ErrProyectoNoEncontrado)
	}
	if !reflect.DeepEqual(got, domain.TotalHoras{}) {
		t.Errorf("con error se devolvió %+v, se esperaba un total vacío", got)
	}
	if len(f.consultasProyecto) != 0 {
		t.Errorf("se consultaron las horas de %v, no se esperaba ninguna consulta", f.consultasProyecto)
	}
}

// Traza: HU-23 CA-23.14 (RN-23.11; sección 7: primero se busca el proyecto y después el sprint)
func TestConsultarHorasSprint_ProyectoInexistente_NoBuscaElSprint(t *testing.T) {
	f := nuevasHorasFalsas()

	if _, err := horasDeSprint(f, 99, 10); err == nil {
		t.Fatal("se esperaba un error")
	}
	if f.vecesSprint != 0 {
		t.Errorf("se buscó el sprint %d veces, no se esperaba ninguna", f.vecesSprint)
	}
}

// Las fallas de infraestructura se propagan y no se confunden con "no encontrado". Al buscar, salen
// como vienen; al leer las horas, envueltas con %w.
// Traza: HU-23 (sección 7: falla interna al consultar)
func TestConsultarHoras_FallaDeInfraestructura_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	consultarSprint := func(f *horasFalsas) (domain.TotalHoras, error) { return horasDeSprint(f, 1, 10) }
	consultarProyecto := func(f *horasFalsas) (domain.TotalHoras, error) { return horasDeProyecto(f, 1) }
	tests := []struct {
		nombre    string
		consultar func(f *horasFalsas) (domain.TotalHoras, error)
		preparar  func(f *horasFalsas)
		consultas int
		envuelta  bool
	}{
		{"sprint: al buscar el proyecto", consultarSprint, func(f *horasFalsas) { f.errProyecto = errBase }, 0, false},
		{"sprint: al buscar el sprint", consultarSprint, func(f *horasFalsas) { f.errSprint = errBase }, 0, false},
		{"sprint: al leer las horas", consultarSprint, func(f *horasFalsas) { f.errHoras = errBase }, 1, true},
		{"proyecto: al buscar el proyecto", consultarProyecto, func(f *horasFalsas) { f.errProyecto = errBase }, 0, false},
		{"proyecto: al leer las horas", consultarProyecto, func(f *horasFalsas) { f.errHoras = errBase }, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevasHorasFalsas()
			tt.preparar(f)

			got, err := tt.consultar(f)

			if !errors.Is(err, errBase) {
				t.Fatalf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
			if tt.envuelta == (err == errBase) {
				t.Errorf("error = %q; ¿envuelto? se esperaba %t", err, tt.envuelta)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrSprintNoEncontrado) {
				t.Errorf("una falla interna no puede parecer \"no encontrado\": %v", err)
			}
			if !reflect.DeepEqual(got, domain.TotalHoras{}) {
				t.Errorf("con error se devolvió %+v, se esperaba un total vacío", got)
			}
			if consultas := len(f.consultasSprint) + len(f.consultasProyecto); consultas != tt.consultas {
				t.Errorf("consultas de horas = %d, se esperaban %d", consultas, tt.consultas)
			}
		})
	}
}
