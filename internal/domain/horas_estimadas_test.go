package domain

import (
	"errors"
	"math"
	"testing"
)

// Traza: HU-19 CA-19.14 (RN-19.2)
func TestValidarHorasEstimadas_Validas(t *testing.T) {
	for _, horas := range []float64{0.5, 1, 8, 24.5, 999.5, 1000} {
		if err := validarHorasEstimadas(horas); err != nil {
			t.Errorf("%v: error inesperado: %v", horas, err)
		}
	}
}

// Los casos marcados prueban el orden de RN-19.2: mayor a 0, múltiplo de 0,5 y máximo 1000.
// Traza: HU-19 CA-19.14 (RN-19.2; spec, sección 6: NaN e Inf)
func TestValidarHorasEstimadas_Invalidas(t *testing.T) {
	tests := []struct {
		horas   float64
		wantErr error
	}{
		{0, ErrHorasInvalidas},
		{-0.5, ErrHorasInvalidas},
		{-1, ErrHorasInvalidas},
		{-1.25, ErrHorasInvalidas}, // orden: negativo y además no múltiplo
		{math.NaN(), ErrHorasInvalidas},
		{math.Inf(-1), ErrHorasInvalidas},
		{1.25, ErrHorasNoMultiploMediaHora},
		{0.1, ErrHorasNoMultiploMediaHora},
		{1000.25, ErrHorasNoMultiploMediaHora}, // orden: no múltiplo y además mayor al máximo
		{1000.5, ErrHorasEstimadasExcedenMaximo},
		{1e300, ErrHorasEstimadasExcedenMaximo}, // desbordaría el entero de medias horas de la base
		{math.Inf(1), ErrHorasEstimadasExcedenMaximo},
	}
	for _, tt := range tests {
		if err := validarHorasEstimadas(tt.horas); !errors.Is(err, tt.wantErr) {
			t.Errorf("%v: error = %v, se esperaba %v", tt.horas, err, tt.wantErr)
		}
	}
}

// Traza: HU-19 CA-19.14 (RN-19.2; spec, sección 7: mensaje del error nuevo)
func TestErrHorasEstimadasExcedenMaximo_Mensaje(t *testing.T) {
	want := "las horas estimadas no pueden superar las 1000 horas"
	if got := ErrHorasEstimadasExcedenMaximo.Error(); got != want {
		t.Errorf("mensaje = %q, se esperaba %q", got, want)
	}
}

// El tope de 24 h de HU-18 no se aplica a las estimaciones (RN-19.2), y el de 1000 h de HU-19 no se
// aplica a los registros: separar las comprobaciones comunes no tiene que cruzarlos.
// Traza: HU-19 CA-19.14 (RN-19.2), HU-18 CA-18.10 (RN-18.8)
func TestValidarHoras_CadaUnaConSuTope(t *testing.T) {
	if err := validarHorasEstimadas(24.5); err != nil {
		t.Errorf("estimadas 24,5: error inesperado: %v", err)
	}
	if err := validarHoras(24); err != nil {
		t.Errorf("registro de 24: error inesperado: %v", err)
	}
	if err := validarHoras(24.5); !errors.Is(err, ErrHorasExcedenMaximo) {
		t.Errorf("registro de 24,5: error = %v, se esperaba ErrHorasExcedenMaximo", err)
	}
	if err := validarHoras(math.Inf(1)); !errors.Is(err, ErrHorasExcedenMaximo) {
		t.Errorf("registro de +Inf: error = %v, se esperaba ErrHorasExcedenMaximo", err)
	}
}
