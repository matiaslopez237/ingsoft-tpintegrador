package domain

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

// proyectoAlfa es el proyecto de partida de los tests de HU-02: guardado, con fin, en estado Activo.
func proyectoAlfa(t *testing.T) Proyecto {
	t.Helper()
	return Proyecto{
		ID:          7,
		Nombre:      "Alfa",
		FechaInicio: fecha(t, "2026-03-01"),
		FechaFin:    ptr(fecha(t, "2026-11-30")),
		Estado:      EstadoActivo,
	}
}

func texto(s string) *string  { return &s }
func estado(e Estado) *Estado { return &e }

// Traza: HU-02 CA-02.8, RN-02.8
func TestProyecto_Modificar_SinCambios_RetornaError(t *testing.T) {
	_, err := proyectoAlfa(t).Modificar(CambiosProyecto{})

	if !errors.Is(err, ErrSinCambios) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSinCambios)
	}
}

// Traza: HU-02 CA-02.1, CA-02.10, RN-02.1, RN-02.2, RN-02.9
func TestProyecto_Modificar_SoloNombre_ConservaElResto(t *testing.T) {
	original := proyectoAlfa(t)

	p, err := original.Modificar(CambiosProyecto{Nombre: texto("  Alfa Renovado  ")})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.Nombre != "Alfa Renovado" {
		t.Errorf("Nombre = %q, se esperaba %q", p.Nombre, "Alfa Renovado")
	}
	if p.ID != original.ID {
		t.Errorf("ID = %d, se esperaba %d", p.ID, original.ID)
	}
	if !p.FechaInicio.Equal(original.FechaInicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", p.FechaInicio, original.FechaInicio)
	}
	if p.FechaFin == nil || !p.FechaFin.Equal(*original.FechaFin) {
		t.Errorf("FechaFin = %v, se esperaba %v", p.FechaFin, original.FechaFin)
	}
	if p.Estado != original.Estado {
		t.Errorf("Estado = %q, se esperaba %q", p.Estado, original.Estado)
	}
}

// Traza: HU-02 CA-02.2, CA-02.4, CA-02.7, CA-02.11, RN-02.2, RN-02.4, RN-02.5, RN-02.6
func TestProyecto_Modificar_DatosInvalidos_RetornaError(t *testing.T) {
	tests := []struct {
		nombre  string
		cambios CambiosProyecto
		wantErr error
	}{
		{"nombre vacío", CambiosProyecto{Nombre: texto("")}, ErrNombreVacio},
		{"nombre de solo espacios", CambiosProyecto{Nombre: texto("   ")}, ErrNombreVacio},
		{"nombre de 101 caracteres", CambiosProyecto{Nombre: texto(strings.Repeat("a", 101))}, ErrNombreDemasiadoLargo},
		{"quitar la fecha de inicio", CambiosProyecto{FechaInicio: &time.Time{}}, ErrFechaInicioObligatoria},
		{"fin anterior al inicio guardado", CambiosProyecto{CambiarFechaFin: true, FechaFin: ptr(fecha(t, "2026-02-28"))}, ErrFechaFinAnterior},
		{"inicio posterior al fin guardado", CambiosProyecto{FechaInicio: ptr(fecha(t, "2026-12-01"))}, ErrFechaFinAnterior},
		{"estado que no existe", CambiosProyecto{Estado: estado("Cerrado")}, ErrEstadoInvalido},
		{"estado en minúscula", CambiosProyecto{Estado: estado("activo")}, ErrEstadoInvalido},
		{"nombre vacío tiene prioridad sobre fechas y estado", CambiosProyecto{
			Nombre: texto(""), FechaInicio: &time.Time{}, Estado: estado("Cerrado"),
		}, ErrNombreVacio},
		{"fechas tienen prioridad sobre estado", CambiosProyecto{
			FechaInicio: &time.Time{}, Estado: estado("Cerrado"),
		}, ErrFechaInicioObligatoria},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := proyectoAlfa(t).Modificar(tt.cambios)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-02 CA-02.2, CA-02.3, CA-02.4, CA-02.5, CA-02.6, RN-02.5, RN-02.6, RN-02.7
func TestProyecto_Modificar_CasosValidos(t *testing.T) {
	tests := []struct {
		nombre     string
		cambios    CambiosProyecto
		wantNombre string
		wantInicio string
		wantFin    string // "" = sin fecha de fin
		wantEstado Estado
	}{
		{"nombre de exactamente 100 caracteres",
			CambiosProyecto{Nombre: texto(strings.Repeat("a", 100))},
			strings.Repeat("a", 100), "2026-03-01", "2026-11-30", EstadoActivo},
		{"cambiar solo las mayúsculas del nombre",
			CambiosProyecto{Nombre: texto("ALFA")},
			"ALFA", "2026-03-01", "2026-11-30", EstadoActivo},
		{"fin igual al inicio",
			CambiosProyecto{CambiarFechaFin: true, FechaFin: ptr(fecha(t, "2026-03-01"))},
			"Alfa", "2026-03-01", "2026-03-01", EstadoActivo},
		{"quitar la fecha de fin",
			CambiosProyecto{CambiarFechaFin: true, FechaFin: nil},
			"Alfa", "2026-03-01", "", EstadoActivo},
		{"mover inicio y fin juntos",
			CambiosProyecto{FechaInicio: ptr(fecha(t, "2027-01-10")), CambiarFechaFin: true, FechaFin: ptr(fecha(t, "2027-06-30"))},
			"Alfa", "2027-01-10", "2027-06-30", EstadoActivo},
		{"finalizar",
			CambiosProyecto{Estado: estado(EstadoFinalizado)},
			"Alfa", "2026-03-01", "2026-11-30", EstadoFinalizado},
		{"mismo estado que ya tiene",
			CambiosProyecto{Estado: estado(EstadoActivo)},
			"Alfa", "2026-03-01", "2026-11-30", EstadoActivo},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p, err := proyectoAlfa(t).Modificar(tt.cambios)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if p.Nombre != tt.wantNombre {
				t.Errorf("Nombre = %q, se esperaba %q", p.Nombre, tt.wantNombre)
			}
			if got := p.FechaInicio.Format(time.DateOnly); got != tt.wantInicio {
				t.Errorf("FechaInicio = %s, se esperaba %s", got, tt.wantInicio)
			}
			gotFin := ""
			if p.FechaFin != nil {
				gotFin = p.FechaFin.Format(time.DateOnly)
			}
			if gotFin != tt.wantFin {
				t.Errorf("FechaFin = %q, se esperaba %q", gotFin, tt.wantFin)
			}
			if p.Estado != tt.wantEstado {
				t.Errorf("Estado = %q, se esperaba %q", p.Estado, tt.wantEstado)
			}
		})
	}
}

// Traza: HU-02 CA-02.6, RN-02.6: se puede pasar a cualquier estado desde cualquier estado.
func TestProyecto_Modificar_CualquierTransicionDeEstado(t *testing.T) {
	estados := []Estado{EstadoActivo, EstadoEnPausa, EstadoFinalizado, EstadoPostergado, EstadoCancelado}
	for _, desde := range estados {
		for _, hacia := range estados {
			t.Run(string(desde)+" a "+string(hacia), func(t *testing.T) {
				p := proyectoAlfa(t)
				p.Estado = desde

				got, err := p.Modificar(CambiosProyecto{Estado: estado(hacia)})

				if err != nil {
					t.Fatalf("error inesperado: %v", err)
				}
				if got.Estado != hacia {
					t.Errorf("Estado = %q, se esperaba %q", got.Estado, hacia)
				}
			})
		}
	}
}

// Modificar no debe alterar el proyecto original, ni cuando falla ni cuando funciona.
// Traza: HU-02 sección 3 ("si hay un error, no se modifica nada").
func TestProyecto_Modificar_NoAlteraElOriginal(t *testing.T) {
	original := proyectoAlfa(t)
	copia := original
	finOriginal := *original.FechaFin

	_, _ = original.Modificar(CambiosProyecto{Nombre: texto("")})
	_, _ = original.Modificar(CambiosProyecto{Nombre: texto("Beta"), CambiarFechaFin: true, FechaFin: ptr(fecha(t, "2026-12-31"))})

	if original.Nombre != copia.Nombre || original.Estado != copia.Estado || !original.FechaInicio.Equal(copia.FechaInicio) {
		t.Errorf("el original cambió: %+v, se esperaba %+v", original, copia)
	}
	if original.FechaFin == nil || !original.FechaFin.Equal(finOriginal) {
		t.Errorf("la fecha de fin del original cambió: %v, se esperaba %v", original.FechaFin, finOriginal)
	}
}

// Traza: HU-02 RN-02.6 (los cinco estados de RN-01.7 son los únicos válidos)
func TestEstadoValido(t *testing.T) {
	validos := []Estado{EstadoActivo, EstadoEnPausa, EstadoFinalizado, EstadoPostergado, EstadoCancelado}
	for _, e := range validos {
		if !EstadoValido(e) {
			t.Errorf("EstadoValido(%q) = false, se esperaba true", e)
		}
	}
	for _, e := range []Estado{"", "Cerrado", "activo", "ACTIVO", " Activo"} {
		if EstadoValido(e) {
			t.Errorf("EstadoValido(%q) = true, se esperaba false", e)
		}
	}
}

// Traza: HU-02 CA-02.12, RN-02.10: modificar no borra los integrantes ni sus roles, y la copia no
// comparte la lista con el original.
func TestProyecto_Modificar_ConservaLosIntegrantes(t *testing.T) {
	original := proyectoAlfa(t)
	original.Integrantes = []Integrante{
		{Nombre: "Ana", Roles: []string{"Desarrolladora", "Tester"}},
		{Nombre: "Beto", Roles: []string{"Scrum Master"}},
	}

	p, err := original.Modificar(CambiosProyecto{Nombre: texto("Alfa Renovado")})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(p.Integrantes) != 2 ||
		p.Integrantes[0].Nombre != "Ana" || !slices.Equal(p.Integrantes[0].Roles, []string{"Desarrolladora", "Tester"}) ||
		p.Integrantes[1].Nombre != "Beto" || !slices.Equal(p.Integrantes[1].Roles, []string{"Scrum Master"}) {
		t.Fatalf("Integrantes = %+v, se esperaban los mismos del original", p.Integrantes)
	}
	p.Integrantes[0].Roles[0] = "Cambiado"
	p.Integrantes[1].Nombre = "Cambiado"
	if original.Integrantes[0].Roles[0] != "Desarrolladora" || original.Integrantes[1].Nombre != "Beto" {
		t.Errorf("cambiar la copia alteró el original: %+v", original.Integrantes)
	}
}
