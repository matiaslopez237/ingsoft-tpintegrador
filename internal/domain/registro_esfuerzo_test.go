package domain

import (
	"errors"
	"testing"
	"time"
)

// dia arma una fecha a la hora indicada, en la zona del servidor (RN-18.4).
func dia(anio int, mes time.Month, d, hora int) time.Time {
	return time.Date(anio, mes, d, hora, 0, 0, 0, time.Local)
}

// hoyEsfuerzo es el "hoy" fijo de los tests: el dominio lo recibe por parámetro.
var hoyEsfuerzo = dia(2026, time.October, 25, 12)

// Traza: HU-18 CA-18.1
func TestNuevoRegistroEsfuerzo_Valido(t *testing.T) {
	r, err := NuevoRegistroEsfuerzo(7, 3, "Ana", dia(2026, time.October, 24, 0), "Desarrollo", 2, hoyEsfuerzo)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := RegistroEsfuerzo{
		ProyectoID: 7,
		ElementoID: 3,
		Integrante: "Ana",
		Fecha:      dia(2026, time.October, 24, 0),
		Actividad:  "Desarrollo",
		Horas:      2,
	}
	if r != want {
		t.Errorf("se obtuvo %+v, se esperaba %+v", r, want)
	}
}

// Traza: HU-18 CA-18.1 (spec, sección 3: integrante y actividad recortados)
func TestNuevoRegistroEsfuerzo_RecortaIntegranteYActividad(t *testing.T) {
	r, err := NuevoRegistroEsfuerzo(1, 1, "  Ana ", hoyEsfuerzo, "  Desarrollo  ", 2, hoyEsfuerzo)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if r.Integrante != "Ana" || r.Actividad != "Desarrollo" {
		t.Errorf("Integrante = %q, Actividad = %q; se esperaba \"Ana\" y \"Desarrollo\"", r.Integrante, r.Actividad)
	}
}

// Traza: HU-18 CA-18.4
func TestNuevoRegistroEsfuerzo_SinActividad_RetornaError(t *testing.T) {
	for _, actividad := range []string{"", "   ", " \t "} {
		_, err := NuevoRegistroEsfuerzo(1, 1, "Ana", hoyEsfuerzo, actividad, 2, hoyEsfuerzo)
		if !errors.Is(err, ErrActividadVacia) {
			t.Errorf("actividad %q: error = %v, se esperaba ErrActividadVacia", actividad, err)
		}
	}
}

// Traza: HU-18 CA-18.6
func TestNuevoRegistroEsfuerzo_SinFecha_RetornaError(t *testing.T) {
	_, err := NuevoRegistroEsfuerzo(1, 1, "Ana", time.Time{}, "Desarrollo", 2, hoyEsfuerzo)

	if !errors.Is(err, ErrFechaEsfuerzoObligatoria) {
		t.Errorf("error = %v, se esperaba ErrFechaEsfuerzoObligatoria", err)
	}
}

// Traza: HU-18 CA-18.6
func TestNuevoRegistroEsfuerzo_FechaFutura_RetornaError(t *testing.T) {
	tests := []struct {
		nombre string
		fecha  time.Time
		hoy    time.Time
	}{
		{"mañana", dia(2026, time.October, 26, 0), hoyEsfuerzo},
		{"mañana a primera hora, hoy a última hora", dia(2026, time.October, 26, 0), dia(2026, time.October, 25, 23)},
		{"el año que viene", dia(2027, time.January, 1, 0), hoyEsfuerzo},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoRegistroEsfuerzo(1, 1, "Ana", tt.fecha, "Desarrollo", 2, tt.hoy)
			if !errors.Is(err, ErrFechaEsfuerzoFutura) {
				t.Errorf("error = %v, se esperaba ErrFechaEsfuerzoFutura", err)
			}
		})
	}
}

// Traza: HU-18 CA-18.6 (RN-18.4: se compara por día, no por hora) y CA-18.14
func TestNuevoRegistroEsfuerzo_FechaHoyOPasada_EsValidaYGuardaSoloElDia(t *testing.T) {
	tests := []struct {
		nombre string
		fecha  time.Time
		hoy    time.Time
		want   time.Time
	}{
		{"hoy, a una hora posterior a la actual", dia(2026, time.October, 25, 23), dia(2026, time.October, 25, 1), dia(2026, time.October, 25, 0)},
		{"hoy, a medianoche", dia(2026, time.October, 25, 0), hoyEsfuerzo, dia(2026, time.October, 25, 0)},
		{"ayer", dia(2026, time.October, 24, 18), hoyEsfuerzo, dia(2026, time.October, 24, 0)},
		{"hace meses", dia(2026, time.February, 15, 0), hoyEsfuerzo, dia(2026, time.February, 15, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			r, err := NuevoRegistroEsfuerzo(1, 1, "Ana", tt.fecha, "Desarrollo", 2, tt.hoy)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !r.Fecha.Equal(tt.want) {
				t.Errorf("Fecha = %v, se esperaba %v", r.Fecha, tt.want)
			}
		})
	}
}

// Traza: HU-18 CA-18.6 (RN-18.4: se comparan días de calendario, aunque las zonas difieran)
func TestNuevoRegistroEsfuerzo_FechaYHoyEnZonasDistintas_ComparaElDiaDeCalendario(t *testing.T) {
	menosTres := time.FixedZone("UTC-3", -3*60*60)
	fecha := time.Date(2026, time.October, 25, 0, 0, 0, 0, menosTres) // 03:00 UTC del 25
	hoy := time.Date(2026, time.October, 25, 1, 0, 0, 0, time.UTC)    // 01:00 UTC del 25

	_, err := NuevoRegistroEsfuerzo(1, 1, "Ana", fecha, "Desarrollo", 2, hoy)

	if err != nil {
		t.Errorf("la fecha es el mismo día que hoy, error inesperado: %v", err)
	}
}

// Traza: HU-18 CA-18.3, CA-18.9, CA-18.10
// El orden es > 0, después múltiplo de 0,5 y al final máximo 24: -1,25 y 25,3 lo verifican.
func TestNuevoRegistroEsfuerzo_HorasNegativas_RetornaError(t *testing.T) {
	tests := []struct {
		horas   float64
		wantErr error
	}{
		{0, ErrHorasInvalidas},
		{-2, ErrHorasInvalidas},
		{-0.5, ErrHorasInvalidas},
		{-1.25, ErrHorasInvalidas},
		{1.25, ErrHorasNoMultiploMediaHora},
		{0.3, ErrHorasNoMultiploMediaHora},
		{25.3, ErrHorasNoMultiploMediaHora},
		{24.5, ErrHorasExcedenMaximo},
		{100, ErrHorasExcedenMaximo},
	}
	for _, tt := range tests {
		_, err := NuevoRegistroEsfuerzo(1, 1, "Ana", hoyEsfuerzo, "Desarrollo", tt.horas, hoyEsfuerzo)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("horas %v: error = %v, se esperaba %v", tt.horas, err, tt.wantErr)
		}
	}
}

// Traza: HU-18 CA-18.2, CA-18.10
func TestNuevoRegistroEsfuerzo_HorasLimite_SonValidas(t *testing.T) {
	for _, horas := range []float64{0.5, 1, 7.5, 24} {
		r, err := NuevoRegistroEsfuerzo(1, 1, "Ana", hoyEsfuerzo, "Desarrollo", horas, hoyEsfuerzo)
		if err != nil {
			t.Errorf("horas %v: error inesperado: %v", horas, err)
			continue
		}
		if r.Horas != horas {
			t.Errorf("Horas = %v, se esperaba %v", r.Horas, horas)
		}
	}
}

// Traza: HU-18 sección 7 (si hay varios errores, se informa el primero en el orden de la tabla)
func TestNuevoRegistroEsfuerzo_VariosErrores_InformaElPrimero(t *testing.T) {
	manana := dia(2026, time.October, 26, 0)
	tests := []struct {
		nombre     string
		integrante string
		fecha      time.Time
		actividad  string
		horas      float64
		wantErr    error
	}{
		{"todo inválido", "", time.Time{}, "", -1, ErrNombreIntegranteVacio},
		{"sin fecha antes que actividad y horas", "Ana", time.Time{}, "", -1, ErrFechaEsfuerzoObligatoria},
		{"fecha futura antes que actividad y horas", "Ana", manana, "", -1, ErrFechaEsfuerzoFutura},
		{"actividad antes que horas", "Ana", hoyEsfuerzo, " ", -1, ErrActividadVacia},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoRegistroEsfuerzo(1, 1, tt.integrante, tt.fecha, tt.actividad, tt.horas, hoyEsfuerzo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-18 CA-18.5
func TestNuevoRegistroEsfuerzo_SinIntegrante_RetornaError(t *testing.T) {
	for _, integrante := range []string{"", "   "} {
		_, err := NuevoRegistroEsfuerzo(1, 1, integrante, hoyEsfuerzo, "Desarrollo", 2, hoyEsfuerzo)
		if !errors.Is(err, ErrNombreIntegranteVacio) {
			t.Errorf("integrante %q: error = %v, se esperaba ErrNombreIntegranteVacio", integrante, err)
		}
	}
}
