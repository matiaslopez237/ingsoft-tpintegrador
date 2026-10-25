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

// Traza: HU-18 CA-18.5
func TestNuevoRegistroEsfuerzo_SinIntegrante_RetornaError(t *testing.T) {
	for _, integrante := range []string{"", "   "} {
		_, err := NuevoRegistroEsfuerzo(1, 1, integrante, hoyEsfuerzo, "Desarrollo", 2, hoyEsfuerzo)
		if !errors.Is(err, ErrNombreIntegranteVacio) {
			t.Errorf("integrante %q: error = %v, se esperaba ErrNombreIntegranteVacio", integrante, err)
		}
	}
}
