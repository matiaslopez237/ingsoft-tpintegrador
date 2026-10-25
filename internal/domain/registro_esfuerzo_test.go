package domain

import (
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
