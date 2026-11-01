package domain

import (
	"math"
	"testing"
)

func pHoras(h float64) *float64 { return &h }

// verHoras muestra un valor opcional en los mensajes de error.
func verHoras(h *float64) any {
	if h == nil {
		return "nil"
	}
	return *h
}

func mismasHoras(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// Traza: HU-19 CA-19.1 a CA-19.8 (RN-19.4 a RN-19.9; spec, tabla de la sección 3 y casos de la sección 6)
func TestCompararEsfuerzo(t *testing.T) {
	tests := []struct {
		nombre     string
		estimadas  *float64
		reales     float64
		desviacion *float64
		porcentaje *float64
		sinEstimar bool
	}{
		{"se trabajó más de lo estimado", pHoras(8), 10, pHoras(2), pHoras(25), false},
		{"se trabajó menos de lo estimado", pHoras(8), 6, pHoras(-2), pHoras(-25), false},
		{"coincide", pHoras(8), 8, pHoras(0), pHoras(0), false},
		{"sin esfuerzo registrado", pHoras(8), 0, pHoras(-8), pHoras(-100), false},
		{"sin estimar, con esfuerzo", nil, 10, nil, nil, true},
		{"sin estimar y sin esfuerzo", nil, 0, nil, nil, true},
		{"estimado 0, con esfuerzo", pHoras(0), 10, pHoras(10), nil, false},
		{"estimado 0 y sin esfuerzo", pHoras(0), 0, pHoras(0), nil, false},
		{"redondea un periódico hacia arriba", pHoras(3), 4, pHoras(1), pHoras(33.3), false},
		{"redondea un periódico negativo", pHoras(3), 2, pHoras(-1), pHoras(-33.3), false},
		{"un medio positivo se aleja de cero", pHoras(40), 40.5, pHoras(0.5), pHoras(1.3), false},
		{"un medio negativo se aleja de cero", pHoras(40), 39.5, pHoras(-0.5), pHoras(-1.3), false},
		{"el porcentaje no tiene tope", pHoras(0.5), 24, pHoras(23.5), pHoras(4700), false},
		{"las horas reales no tienen tope", pHoras(8), 56, pHoras(48), pHoras(600), false},
		{"porcentaje con decimal en las horas", pHoras(8), 5, pHoras(-3), pHoras(-37.5), false},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			got := CompararEsfuerzo(tt.estimadas, tt.reales)

			if !mismasHoras(got.HorasEstimadas, tt.estimadas) {
				t.Errorf("HorasEstimadas = %v, se esperaba %v", verHoras(got.HorasEstimadas), verHoras(tt.estimadas))
			}
			if got.HorasReales != tt.reales {
				t.Errorf("HorasReales = %v, se esperaba %v", got.HorasReales, tt.reales)
			}
			if !mismasHoras(got.DesviacionHoras, tt.desviacion) {
				t.Errorf("DesviacionHoras = %v, se esperaba %v", verHoras(got.DesviacionHoras), verHoras(tt.desviacion))
			}
			if !mismasHoras(got.DesviacionPorcentaje, tt.porcentaje) {
				t.Errorf("DesviacionPorcentaje = %v, se esperaba %v", verHoras(got.DesviacionPorcentaje), verHoras(tt.porcentaje))
			}
			if got.SinEstimar() != tt.sinEstimar {
				t.Errorf("SinEstimar() = %t, se esperaba %t", got.SinEstimar(), tt.sinEstimar)
			}
		})
	}
}

// Dividir antes de multiplicar da un valor apenas menor que el medio exacto y redondea para el otro
// lado (51,5 sobre 40 daría 28,7). Son empates reales del dominio: las horas van en medias horas.
// Traza: HU-19 CA-19.8 (RN-19.6; spec, sección 5: precisión)
func TestCompararEsfuerzo_MultiplicaAntesDeDividir(t *testing.T) {
	tests := []struct {
		estimadas, reales, porcentaje float64
	}{
		{40, 51.5, 28.8},
		{40, 28.5, -28.8},
		{40, 60.5, 51.3},
		{40, 19.5, -51.3},
		{40, 65.5, 63.8},
		{40, 14.5, -63.8},
	}
	for _, tt := range tests {
		got := CompararEsfuerzo(&tt.estimadas, tt.reales)
		if got.DesviacionPorcentaje == nil || *got.DesviacionPorcentaje != tt.porcentaje {
			t.Errorf("%v estimadas y %v reales: porcentaje = %v, se esperaba %v",
				tt.estimadas, tt.reales, verHoras(got.DesviacionPorcentaje), tt.porcentaje)
		}
	}
}

// -0 == 0 en Go, así que solo math.Signbit distingue el caso. Con 2000 estimadas (una suma de HU-23)
// y media hora de menos, el porcentaje es -0,025: redondeado a un decimal, un cero negativo.
// Traza: HU-19 CA-19.8 (RN-19.6: si el redondeo da cero, el resultado es 0, nunca -0)
func TestCompararEsfuerzo_NuncaDevuelveCeroNegativo(t *testing.T) {
	tests := []struct {
		nombre            string
		estimadas, reales float64
	}{
		{"el redondeo de un valor negativo da cero", 2000, 1999.5},
		{"el redondeo de un valor positivo da cero", 2000, 2000.5},
		{"sin desviación", 8, 8},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			got := CompararEsfuerzo(&tt.estimadas, tt.reales)

			if got.DesviacionPorcentaje == nil {
				t.Fatal("DesviacionPorcentaje = nil, se esperaba 0")
			}
			if p := *got.DesviacionPorcentaje; p != 0 || math.Signbit(p) {
				t.Errorf("DesviacionPorcentaje = %v (signo negativo: %t), se esperaba +0", p, math.Signbit(p))
			}
		})
	}
}

// El resultado no comparte el puntero recibido, igual que los story points en NuevoElementoBacklog.
// Traza: HU-19 (spec, sección 5: función del dominio reutilizable)
func TestCompararEsfuerzo_NoCompartePunteroRecibido(t *testing.T) {
	estimadas := 8.0

	got := CompararEsfuerzo(&estimadas, 10)
	estimadas = 99

	if got.HorasEstimadas == nil || *got.HorasEstimadas != 8 {
		t.Errorf("HorasEstimadas = %v, se esperaba 8", verHoras(got.HorasEstimadas))
	}
}
