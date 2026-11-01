package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// comparacionFalsa implementa los tres puertos de CompararEsfuerzo en memoria, permite simular fallas
// y cuenta las consultas. Proyecto 1: el elemento 10, con 8 h estimadas, y el 11, sin estimar.
// Proyecto 2: el elemento 20.
type comparacionFalsa struct {
	proyectos      map[int64]domain.Proyecto
	elementos      map[int64]domain.ElementoBacklog
	horas          map[int64]float64 // horas reales por elemento; sin entrada: 0 (RN-19.4)
	errProyecto    error
	errElemento    error
	errHoras       error
	vecesElemento  int
	consultasHoras []int64
}

func pHoras(h float64) *float64 { return &h }

func nuevaComparacionFalsa() *comparacionFalsa {
	return &comparacionFalsa{
		proyectos: map[int64]domain.Proyecto{
			1: {ID: 1, Nombre: "Alfa", Estado: domain.EstadoActivo},
			2: {ID: 2, Nombre: "Beta", Estado: domain.EstadoActivo},
		},
		elementos: map[int64]domain.ElementoBacklog{
			10: {ID: 10, ProyectoID: 1, Titulo: "Login", Estado: domain.EstadoPendiente, HorasEstimadas: pHoras(8)},
			11: {ID: 11, ProyectoID: 1, Titulo: "Reportes", Estado: domain.EstadoPendiente},
			20: {ID: 20, ProyectoID: 2, Titulo: "Pagos", Estado: domain.EstadoPendiente, HorasEstimadas: pHoras(3)},
		},
		horas: map[int64]float64{},
	}
}

func (f *comparacionFalsa) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *comparacionFalsa) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	f.vecesElemento++
	if f.errElemento != nil {
		return domain.ElementoBacklog{}, f.errElemento
	}
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *comparacionFalsa) HorasRealesDeElemento(_ context.Context, elementoID int64) (float64, error) {
	f.consultasHoras = append(f.consultasHoras, elementoID)
	if f.errHoras != nil {
		return 0, f.errHoras
	}
	return f.horas[elementoID], nil
}

func comparar(f *comparacionFalsa, proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error) {
	return NuevoCompararEsfuerzo(f, f, f).Ejecutar(context.Background(), proyectoID, elementoID)
}

// Traza: HU-19 CA-19.1, CA-19.2, CA-19.3, CA-19.4, CA-19.5, CA-19.6 (RN-19.4 a RN-19.7, RN-19.9)
func TestCompararEsfuerzo_DevuelveLaComparacionDelElemento(t *testing.T) {
	tests := []struct {
		nombre     string
		elementoID int64
		reales     float64
		want       domain.ComparacionEsfuerzo
	}{
		{"desviación positiva", 10, 10, domain.ComparacionEsfuerzo{
			HorasEstimadas: pHoras(8), HorasReales: 10, DesviacionHoras: pHoras(2), DesviacionPorcentaje: pHoras(25),
		}},
		{"desviación negativa", 10, 6, domain.ComparacionEsfuerzo{
			HorasEstimadas: pHoras(8), HorasReales: 6, DesviacionHoras: pHoras(-2), DesviacionPorcentaje: pHoras(-25),
		}},
		{"sin desviación", 10, 8, domain.ComparacionEsfuerzo{
			HorasEstimadas: pHoras(8), HorasReales: 8, DesviacionHoras: pHoras(0), DesviacionPorcentaje: pHoras(0),
		}},
		{"sin esfuerzo", 10, 0, domain.ComparacionEsfuerzo{
			HorasEstimadas: pHoras(8), HorasReales: 0, DesviacionHoras: pHoras(-8), DesviacionPorcentaje: pHoras(-100),
		}},
		{"sin estimar, con esfuerzo", 11, 10, domain.ComparacionEsfuerzo{HorasReales: 10}},
		{"sin estimar y sin esfuerzo", 11, 0, domain.ComparacionEsfuerzo{HorasReales: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevaComparacionFalsa()
			f.horas[tt.elementoID] = tt.reales

			got, err := comparar(f, 1, tt.elementoID)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("se obtuvo %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

// Las horas estimadas salen del elemento pedido, y las reales se piden solo de ese elemento.
// Traza: HU-19 CA-19.11 (RN-19.4), sección 5: caso de uso
func TestCompararEsfuerzo_UsaLasHorasDelElementoPedido(t *testing.T) {
	f := nuevaComparacionFalsa()
	f.elementos[12] = domain.ElementoBacklog{ID: 12, ProyectoID: 1, Titulo: "Ayuda", HorasEstimadas: pHoras(4)}
	f.horas[10] = 100 // otro elemento del mismo proyecto: no tiene que sumar
	f.horas[12] = 5

	got, err := comparar(f, 1, 12)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.HorasEstimadas == nil || *got.HorasEstimadas != 4 || got.HorasReales != 5 {
		t.Errorf("se obtuvo %+v, se esperaba 4 estimadas y 5 reales", got)
	}
	if !reflect.DeepEqual(f.consultasHoras, []int64{12}) {
		t.Errorf("horas consultadas de los elementos %v, se esperaba solo el 12", f.consultasHoras)
	}
}

// Traza: HU-19 CA-19.12 (RN-19.11)
func TestCompararEsfuerzo_CualquierEstadoDeElementoYProyecto(t *testing.T) {
	sprint := int64(3)
	for _, estadoProyecto := range []domain.Estado{domain.EstadoActivo, domain.EstadoEnPausa,
		domain.EstadoFinalizado, domain.EstadoPostergado, domain.EstadoCancelado} {
		for _, estadoElemento := range []domain.EstadoElemento{domain.EstadoPendiente, domain.EstadoEnProgreso, domain.EstadoHecho} {
			for _, sprintID := range []*int64{nil, &sprint} {
				f := nuevaComparacionFalsa()
				p := f.proyectos[1]
				p.Estado = estadoProyecto
				f.proyectos[1] = p
				e := f.elementos[10]
				e.Estado, e.SprintID = estadoElemento, sprintID
				f.elementos[10] = e
				f.horas[10] = 2

				got, err := comparar(f, 1, 10)

				if err != nil {
					t.Errorf("proyecto %s, elemento %s, sprint %v: error inesperado: %v", estadoProyecto, estadoElemento, sprintID, err)
					continue
				}
				if got.HorasReales != 2 || got.DesviacionHoras == nil || *got.DesviacionHoras != -6 {
					t.Errorf("proyecto %s, elemento %s, sprint %v: se obtuvo %+v", estadoProyecto, estadoElemento, sprintID, got)
				}
			}
		}
	}
}

// Los errores de validación se devuelven tal cual (por eso == y no errors.Is), y sin consultar las
// horas reales: la consulta no hace falta si no hay nada que comparar.
// Traza: HU-19 CA-19.9, CA-19.10 (RN-19.10)
func TestCompararEsfuerzo_ValidacionFallida_RetornaErrorYNoConsultaHoras(t *testing.T) {
	tests := []struct {
		nombre     string
		proyectoID int64
		elementoID int64
		wantErr    error
	}{
		{"proyecto inexistente", 99, 10, domain.ErrProyectoNoEncontrado},
		{"elemento inexistente", 1, 99, domain.ErrElementoNoEncontrado},
		{"elemento de otro proyecto", 1, 20, domain.ErrElementoNoEncontrado},
		{"otro proyecto, visto desde el segundo", 2, 10, domain.ErrElementoNoEncontrado},
		{"proyecto y elemento inexistentes: primero el proyecto", 99, 99, domain.ErrProyectoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevaComparacionFalsa()

			got, err := comparar(f, tt.proyectoID, tt.elementoID)

			if err != tt.wantErr {
				t.Errorf("error = %v, se esperaba exactamente %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, domain.ComparacionEsfuerzo{}) {
				t.Errorf("con error se devolvió %+v, se esperaba una comparación vacía", got)
			}
			if len(f.consultasHoras) != 0 {
				t.Errorf("se consultaron las horas de %v, no se esperaba ninguna consulta", f.consultasHoras)
			}
		})
	}
}

// Traza: HU-19 sección 7: primero se busca el proyecto y después el elemento.
func TestCompararEsfuerzo_ProyectoInexistente_NoBuscaElElemento(t *testing.T) {
	f := nuevaComparacionFalsa()

	if _, err := comparar(f, 99, 10); err == nil {
		t.Fatal("se esperaba un error")
	}
	if f.vecesElemento != 0 {
		t.Errorf("se buscó el elemento %d veces, no se esperaba ninguna", f.vecesElemento)
	}
}

// Las fallas de infraestructura se propagan y no se confunden con "no encontrado". Al buscar, salen
// como vienen; al leer las horas, envueltas con %w.
// Traza: HU-19 sección 7: falla interna al consultar
func TestCompararEsfuerzo_FallaDeInfraestructura_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre        string
		preparar      func(f *comparacionFalsa)
		consultaHoras int
		envuelta      bool
	}{
		{"al buscar el proyecto", func(f *comparacionFalsa) { f.errProyecto = errBase }, 0, false},
		{"al buscar el elemento", func(f *comparacionFalsa) { f.errElemento = errBase }, 0, false},
		{"al leer las horas", func(f *comparacionFalsa) { f.errHoras = errBase }, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevaComparacionFalsa()
			tt.preparar(f)

			got, err := comparar(f, 1, 10)

			if !errors.Is(err, errBase) {
				t.Fatalf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
			if tt.envuelta == (err == errBase) {
				t.Errorf("error = %q; ¿envuelto? se esperaba %t", err, tt.envuelta)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrElementoNoEncontrado) {
				t.Errorf("una falla interna no puede parecer \"no encontrado\": %v", err)
			}
			if !reflect.DeepEqual(got, domain.ComparacionEsfuerzo{}) {
				t.Errorf("con error se devolvió %+v, se esperaba una comparación vacía", got)
			}
			if len(f.consultasHoras) != tt.consultaHoras {
				t.Errorf("consultas de horas = %v, se esperaban %d", f.consultasHoras, tt.consultaHoras)
			}
		})
	}
}
