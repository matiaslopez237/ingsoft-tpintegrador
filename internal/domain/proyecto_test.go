package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func fecha(t *testing.T, s string) time.Time {
	t.Helper()
	f, err := time.Parse(time.DateOnly, s)
	if err != nil {
		t.Fatalf("fecha inválida %q: %v", s, err)
	}
	return f
}

func ptr(f time.Time) *time.Time { return &f }

// CA-01.1, RN-01.6
func TestNuevoProyecto_DatosValidos_CreaProyectoActivo(t *testing.T) {
	inicio, fin := fecha(t, "2026-03-01"), fecha(t, "2026-11-30")

	p, err := NuevoProyecto("TP Integrador", inicio, &fin)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.Nombre != "TP Integrador" {
		t.Errorf("Nombre = %q, se esperaba %q", p.Nombre, "TP Integrador")
	}
	if !p.FechaInicio.Equal(inicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", p.FechaInicio, inicio)
	}
	if p.FechaFin == nil || !p.FechaFin.Equal(fin) {
		t.Errorf("FechaFin = %v, se esperaba %v", p.FechaFin, fin)
	}
	if p.Estado != EstadoActivo {
		t.Errorf("Estado = %q, se esperaba %q", p.Estado, EstadoActivo)
	}
	if p.ID != 0 {
		t.Errorf("ID = %d, se esperaba 0 (lo asigna la persistencia)", p.ID)
	}
}

// CA-01.2, CA-01.3, CA-01.5, CA-01.6
func TestNuevoProyecto_DatosInvalidos_RetornaError(t *testing.T) {
	inicio := fecha(t, "2026-03-02")
	tests := []struct {
		nombre    string
		nombreArg string
		inicio    time.Time
		fin       *time.Time
		wantErr   error
	}{
		{"nombre vacío", "", inicio, nil, ErrNombreVacio},
		{"nombre de solo espacios", "   \t ", inicio, nil, ErrNombreVacio},
		{"nombre de 101 caracteres", strings.Repeat("a", 101), inicio, nil, ErrNombreDemasiadoLargo},
		{"nombre de 101 runas con acentos", strings.Repeat("ñ", 101), inicio, nil, ErrNombreDemasiadoLargo},
		{"sin fecha de inicio", "Alfa", time.Time{}, nil, ErrFechaInicioObligatoria},
		{"fecha de fin anterior a inicio", "Alfa", inicio, ptr(fecha(t, "2026-03-01")), ErrFechaFinAnterior},
		{"nombre vacío tiene prioridad sobre fechas", "", time.Time{}, nil, ErrNombreVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoProyecto(tt.nombreArg, tt.inicio, tt.fin)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// CA-01.3, CA-01.7, CA-01.8
func TestNuevoProyecto_CasosLimite_SonValidos(t *testing.T) {
	inicio := fecha(t, "2026-03-01")
	tests := []struct {
		nombre     string
		nombreArg  string
		fin        *time.Time
		wantNombre string
	}{
		{"fechas iguales", "Alfa", ptr(inicio), "Alfa"},
		{"sin fecha de fin", "Alfa", nil, "Alfa"},
		{"nombre de exactamente 100 caracteres", strings.Repeat("a", 100), nil, strings.Repeat("a", 100)},
		{"nombre de 100 runas con acentos", strings.Repeat("ñ", 100), nil, strings.Repeat("ñ", 100)},
		{"nombre con espacios alrededor se recorta", "  Alfa  ", nil, "Alfa"},
		{"los espacios no cuentan para el máximo", " " + strings.Repeat("a", 100) + " ", nil, strings.Repeat("a", 100)},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p, err := NuevoProyecto(tt.nombreArg, inicio, tt.fin)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if p.Nombre != tt.wantNombre {
				t.Errorf("Nombre = %q, se esperaba %q", p.Nombre, tt.wantNombre)
			}
			if (p.FechaFin == nil) != (tt.fin == nil) {
				t.Errorf("FechaFin = %v, se esperaba %v", p.FechaFin, tt.fin)
			}
		})
	}
}

// RN-01.3: la clave de unicidad ignora espacios y mayúsculas.
func TestClaveNombre_IgnoraMayusculasYEspacios(t *testing.T) {
	tests := []struct{ a, b string }{
		{"Alfa", "alfa"},
		{"Alfa", "  ALFA "},
		{"Gestión Año", "GESTIÓN AÑO"},
	}
	for _, tt := range tests {
		if ClaveNombre(tt.a) != ClaveNombre(tt.b) {
			t.Errorf("ClaveNombre(%q) = %q y ClaveNombre(%q) = %q, se esperaban iguales",
				tt.a, ClaveNombre(tt.a), tt.b, ClaveNombre(tt.b))
		}
	}
	if ClaveNombre("Alfa") == ClaveNombre("Beta") {
		t.Error("nombres distintos no deberían tener la misma clave")
	}
}

// RN-01.7
func TestEstados_SonLosDefinidosEnLaSpec(t *testing.T) {
	want := []Estado{"Activo", "En pausa", "Finalizado", "Postergado", "Cancelado"}
	got := []Estado{EstadoActivo, EstadoEnPausa, EstadoFinalizado, EstadoPostergado, EstadoCancelado}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("estado %d = %q, se esperaba %q", i, got[i], want[i])
		}
	}
}
