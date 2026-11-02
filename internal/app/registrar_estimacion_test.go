package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// estimacionFalsa envuelve el fake de HU-14 (pokerFalso) y le suma el guardado de HU-17, sin tocarlo.
type estimacionFalsa struct {
	*pokerFalso
	guardadas  []domain.Ronda
	errGuardar error
}

func (f *estimacionFalsa) GuardarEstimacion(_ context.Context, r domain.Ronda) error {
	if f.errGuardar != nil {
		return f.errGuardar
	}
	f.guardadas = append(f.guardadas, r)
	return nil
}

// conRondaCerrada devuelve el fake con la ronda 1 de Login cerrada (Ana 3, Beto 13).
func conRondaCerrada(t *testing.T) *estimacionFalsa {
	t.Helper()
	f := nuevoPokerFalso(t, "Ana", "Beto")
	rondaCerradaEn(t, f, 7, map[string]int{"Ana": 3, "Beto": 13})
	return &estimacionFalsa{pokerFalso: f}
}

func registrarEstimacion(f *estimacionFalsa, proyectoID, rondaID int64, valor int) (domain.ResumenRonda, error) {
	return NuevoRegistrarEstimacion(f, f, f).Ejecutar(context.Background(),
		DatosEstimacion{ProyectoID: proyectoID, RondaID: rondaID, Valor: valor})
}

// Traza: HU-17 CA-17.1, CA-17.3 (RN-17.1, RN-17.2): guarda la ronda con la estimación, aunque nadie
// haya votado ese valor, y devuelve el resumen con ella.
func TestRegistrarEstimacion_Valida_LaGuarda(t *testing.T) {
	f := conRondaCerrada(t)

	res, err := registrarEstimacion(f, 1, 1, 8)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.EstimacionAcordada == nil || *res.EstimacionAcordada != 8 {
		t.Errorf("EstimacionAcordada del resumen = %v, se esperaba 8", res.EstimacionAcordada)
	}
	if len(f.guardadas) != 1 || f.guardadas[0].ID != 1 || *f.guardadas[0].EstimacionAcordada != 8 {
		t.Errorf("guardadas = %+v, se esperaba la ronda 1 con la estimación 8", f.guardadas)
	}
}

// Traza: HU-17 CA-17.3 a CA-17.5, CA-17.8 (sección 7): si algo falla, no se guarda nada.
func TestRegistrarEstimacion_Errores_NoGuarda(t *testing.T) {
	tests := []struct {
		name                string
		preparar            func(*testing.T, *estimacionFalsa)
		proyectoID, rondaID int64
		valor               int
		want                error
	}{
		{"proyecto inexistente", nil, 99, 1, 8, domain.ErrProyectoNoEncontrado},
		{"ronda inexistente", nil, 1, 99, 8, domain.ErrRondaNoEncontrada},
		{"ronda de otro proyecto", nil, 2, 1, 8, domain.ErrRondaNoEncontrada},
		{"fuera de la escala", nil, 1, 1, 4, domain.ErrEstimacionInvalida},
		{"elemento Hecho", func(t *testing.T, f *estimacionFalsa) {
			e := f.elementos[7]
			e.Estado = domain.EstadoHecho
			f.elementos[7] = e
		}, 1, 1, 8, domain.ErrElementoTerminado},
		{"ya tiene estimación", func(t *testing.T, f *estimacionFalsa) {
			v := 5
			f.rondas[0].EstimacionAcordada = &v
		}, 1, 1, 8, domain.ErrEstimacionYaRegistrada},
		{"no es la última", func(t *testing.T, f *estimacionFalsa) {
			rondaCerradaEn(t, f.pokerFalso, 7, map[string]int{"Ana": 5, "Beto": 8})
		}, 1, 1, 8, domain.ErrRondaNoEsLaUltima},
		{"ronda abierta", func(t *testing.T, f *estimacionFalsa) {
			f.rondas[0].Estado = domain.RondaAbierta
		}, 1, 1, 8, domain.ErrEstimacionRondaAbierta},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := conRondaCerrada(t)
			if tt.preparar != nil {
				tt.preparar(t, f)
			}

			_, err := registrarEstimacion(f, tt.proyectoID, tt.rondaID, tt.valor)

			if !errors.Is(err, tt.want) {
				t.Errorf("error = %v, se esperaba %v", err, tt.want)
			}
			if len(f.guardadas) != 0 {
				t.Errorf("no se tendría que guardar nada: %+v", f.guardadas)
			}
		})
	}
}

// Traza: HU-17 sección 5: si otro pedido registró en el medio, el error del guardado llega sin perderse.
func TestRegistrarEstimacion_FallaAlGuardar_SePropaga(t *testing.T) {
	f := conRondaCerrada(t)
	f.errGuardar = domain.ErrEstimacionYaRegistrada

	_, err := registrarEstimacion(f, 1, 1, 8)

	if !errors.Is(err, domain.ErrEstimacionYaRegistrada) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrEstimacionYaRegistrada)
	}
}

// Traza: HU-17 sección 7: una falla al leer las rondas se propaga sin perder la causa.
func TestRegistrarEstimacion_FallaAlLeer_SePropaga(t *testing.T) {
	f := conRondaCerrada(t)
	disco := errors.New("disco lleno")
	f.errRonda = disco

	_, err := registrarEstimacion(f, 1, 1, 8)

	if !errors.Is(err, disco) {
		t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
	}
}
