package domain

import (
	"math"
	"testing"
)

// he arma un elemento con estimación; sinEst, uno sin estimar.
func he(estimadas, reales float64) HorasElemento {
	return HorasElemento{HorasEstimadas: &estimadas, HorasReales: reales}
}

func sinEst(reales float64) HorasElemento {
	return HorasElemento{HorasReales: reales}
}

// Traza: HU-23 CA-23.1, CA-23.2, CA-23.3, CA-23.4, CA-23.5, CA-23.6, CA-23.7, CA-23.8, CA-23.9 (RN-23.2, RN-23.3, RN-23.4, RN-23.5, RN-23.6; spec, tabla de la sección 3 y casos de la sección 6)
func TestTotalizarHoras(t *testing.T) {
	tests := []struct {
		nombre                              string
		elementos                           []HorasElemento
		estimadas, reales, realesSinEstimar float64
		desviacion, porcentaje              *float64
		cantidad, cantidadSinEstimar        int
		sinEstimar                          bool
	}{
		// CA-23.1, CA-23.2: se suman las horas de los elementos.
		{"suma las horas de un sprint", []HorasElemento{he(8, 10), he(4, 3)},
			12, 13, 0, pHoras(1), pHoras(8.3), 2, 0, false},
		{"suma las horas de un proyecto", []HorasElemento{he(8, 10), he(4, 3), he(20, 15)},
			32, 28, 0, pHoras(-4), pHoras(-12.5), 3, 0, false},
		{"desviación negativa", []HorasElemento{he(8, 6), he(4, 3)},
			12, 9, 0, pHoras(-3), pHoras(-25), 2, 0, false},
		// CA-23.3: el porcentaje de las sumas es 10,0; el promedio de los porcentajes daría 50,0.
		{"el porcentaje es el de las sumas y no el promedio", []HorasElemento{he(10, 20), he(90, 90)},
			100, 110, 0, pHoras(10), pHoras(10), 2, 0, false},
		// CA-23.4: los porcentajes de cada elemento (100 y 0) no dan el de la suma.
		{"redondea el periódico de la suma", []HorasElemento{he(1, 2), he(2, 2)},
			3, 4, 0, pHoras(1), pHoras(33.3), 2, 0, false},
		{"un medio positivo de la suma se aleja de cero", []HorasElemento{he(20, 20.5), he(20, 20)},
			40, 40.5, 0, pHoras(0.5), pHoras(1.3), 2, 0, false},
		{"un medio negativo de la suma se aleja de cero", []HorasElemento{he(20, 19.5), he(20, 20)},
			40, 39.5, 0, pHoras(-0.5), pHoras(-1.3), 2, 0, false},
		// CA-23.5: dos desvíos opuestos se compensan en la suma.
		{"los desvíos opuestos se compensan", []HorasElemento{he(8, 10), he(8, 6)},
			16, 16, 0, pHoras(0), pHoras(0), 2, 0, false},
		// CA-23.6: las reales de un elemento sin estimar van al total, no a la desviación.
		{"un elemento sin estimar no entra en la desviación", []HorasElemento{he(8, 10), sinEst(5)},
			8, 15, 5, pHoras(2), pHoras(25), 2, 1, false},
		{"no importa el orden de los elementos", []HorasElemento{sinEst(5), he(8, 10)},
			8, 15, 5, pHoras(2), pHoras(25), 2, 1, false},
		{"las horas sin estimar no compensan la desviación", []HorasElemento{he(8, 0), sinEst(5)},
			8, 5, 5, pHoras(-8), pHoras(-100), 2, 1, false},
		{"un elemento sin estimar y sin esfuerzo solo se cuenta", []HorasElemento{he(8, 10), sinEst(0)},
			8, 10, 0, pHoras(2), pHoras(25), 2, 1, false},
		// CA-23.7: con ningún elemento estimado no hay desviación; las reales conservan su valor.
		{"ningún elemento estimado, con esfuerzo", []HorasElemento{sinEst(5), sinEst(3)},
			0, 8, 8, nil, nil, 2, 2, true},
		{"ningún elemento estimado, sin esfuerzo", []HorasElemento{sinEst(0)},
			0, 0, 0, nil, nil, 1, 1, true},
		// CA-23.8: sin elementos todo vale 0 y no hay desviación.
		{"lista vacía", []HorasElemento{}, 0, 0, 0, nil, nil, 0, 0, true},
		{"lista nil", nil, 0, 0, 0, nil, nil, 0, 0, true},
		// CA-23.9: un elemento estimado y sin registros suma a las estimadas.
		{"estimado sin esfuerzo", []HorasElemento{he(8, 0)},
			8, 0, 0, pHoras(-8), pHoras(-100), 1, 0, false},
		// Spec, sección 5: las horas van en medias horas, que se suman sin error.
		{"suma medias horas", []HorasElemento{he(0.5, 1), he(1.5, 1), he(2.5, 3.5), sinEst(0.5)},
			4.5, 6, 0.5, pHoras(1), pHoras(22.2), 4, 1, false},
		// Spec, sección 6: un solo elemento da lo mismo que HU-19.
		{"un solo elemento", []HorasElemento{he(3, 4)},
			3, 4, 0, pHoras(1), pHoras(33.3), 1, 0, false},
		// Spec, sección 5: la suma no tiene tope; los 1000 h son por elemento.
		{"la suma supera el tope de un elemento", []HorasElemento{he(1000, 1000), he(1000, 1200)},
			2000, 2200, 0, pHoras(200), pHoras(10), 2, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			got := TotalizarHoras(tt.elementos)

			if got.HorasEstimadas != tt.estimadas {
				t.Errorf("HorasEstimadas = %v, se esperaba %v", got.HorasEstimadas, tt.estimadas)
			}
			if got.HorasReales != tt.reales {
				t.Errorf("HorasReales = %v, se esperaba %v", got.HorasReales, tt.reales)
			}
			if got.HorasRealesSinEstimar != tt.realesSinEstimar {
				t.Errorf("HorasRealesSinEstimar = %v, se esperaba %v", got.HorasRealesSinEstimar, tt.realesSinEstimar)
			}
			if !mismasHoras(got.DesviacionHoras, tt.desviacion) {
				t.Errorf("DesviacionHoras = %v, se esperaba %v", verHoras(got.DesviacionHoras), verHoras(tt.desviacion))
			}
			if !mismasHoras(got.DesviacionPorcentaje, tt.porcentaje) {
				t.Errorf("DesviacionPorcentaje = %v, se esperaba %v", verHoras(got.DesviacionPorcentaje), verHoras(tt.porcentaje))
			}
			if got.Elementos != tt.cantidad {
				t.Errorf("Elementos = %d, se esperaba %d", got.Elementos, tt.cantidad)
			}
			if got.ElementosSinEstimar != tt.cantidadSinEstimar {
				t.Errorf("ElementosSinEstimar = %d, se esperaba %d", got.ElementosSinEstimar, tt.cantidadSinEstimar)
			}
			if got.SinEstimar() != tt.sinEstimar {
				t.Errorf("SinEstimar() = %t, se esperaba %t", got.SinEstimar(), tt.sinEstimar)
			}
		})
	}
}

// Con 2000 estimadas y media hora de menos, el porcentaje es -0,025: redondeado, un cero negativo. El
// caso solo se alcanza con sumas, porque un elemento tiene como máximo 1000 h estimadas (RN-19.2).
// Traza: HU-23 CA-23.4 (RN-23.4: mismo redondeo que RN-19.6, nunca -0)
func TestTotalizarHoras_NuncaDevuelveCeroNegativo(t *testing.T) {
	got := TotalizarHoras([]HorasElemento{he(1000, 999.5), he(1000, 1000)})

	if got.DesviacionHoras == nil || *got.DesviacionHoras != -0.5 {
		t.Errorf("DesviacionHoras = %v, se esperaba -0.5", verHoras(got.DesviacionHoras))
	}
	if got.DesviacionPorcentaje == nil {
		t.Fatal("DesviacionPorcentaje = nil, se esperaba 0")
	}
	if p := *got.DesviacionPorcentaje; p != 0 || math.Signbit(p) {
		t.Errorf("DesviacionPorcentaje = %v (signo negativo: %t), se esperaba +0", p, math.Signbit(p))
	}
}

// El resultado es el de CompararEsfuerzo sobre las sumas de los elementos estimados: mismas desviaciones
// en todos los casos, también en los empates del redondeo (40 y 51,5 dan 28,8; RN-19.6).
// Traza: HU-23 CA-23.1 (RN-23.4, RN-23.5)
func TestTotalizarHoras_ReutilizaCompararEsfuerzo(t *testing.T) {
	estimadas, reales := 40.0, 51.5

	got := TotalizarHoras([]HorasElemento{he(15, 21.5), he(25, 30), sinEst(7)})
	esperado := CompararEsfuerzo(&estimadas, reales)

	if !mismasHoras(got.DesviacionHoras, esperado.DesviacionHoras) {
		t.Errorf("DesviacionHoras = %v, se esperaba %v", verHoras(got.DesviacionHoras), verHoras(esperado.DesviacionHoras))
	}
	if !mismasHoras(got.DesviacionPorcentaje, esperado.DesviacionPorcentaje) {
		t.Errorf("DesviacionPorcentaje = %v, se esperaba %v", verHoras(got.DesviacionPorcentaje), verHoras(esperado.DesviacionPorcentaje))
	}
}

// Totalizar no toca los elementos que recibe: el caso de uso se los pasa tal como los leyó.
// Traza: HU-23 (spec, sección 5: cálculo en el dominio)
func TestTotalizarHoras_NoModificaLosElementos(t *testing.T) {
	estimadas := 8.0
	elementos := []HorasElemento{{HorasEstimadas: &estimadas, HorasReales: 10}, sinEst(5)}

	got := TotalizarHoras(elementos)
	got.HorasEstimadas = 99

	if estimadas != 8 || elementos[0].HorasReales != 10 || elementos[1].HorasEstimadas != nil || elementos[1].HorasReales != 5 {
		t.Errorf("los elementos cambiaron: %+v (estimadas %v)", elementos, estimadas)
	}
}
