package domain

import (
	"errors"
	"testing"
)

// Traza: HU-20 CA-20.1
func TestNuevoDefecto_Valido(t *testing.T) {
	d, err := NuevoDefecto(7, 3, "No carga la pantalla de login en Firefox", SeveridadAlta)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := Defecto{
		ProyectoID:  7,
		ElementoID:  3,
		Descripcion: "No carga la pantalla de login en Firefox",
		Severidad:   SeveridadAlta,
		Estado:      DefectoAbierto,
	}
	if d != want {
		t.Errorf("se obtuvo %+v, se esperaba %+v", d, want)
	}
}

// Traza: HU-20 CA-20.1 (spec, sección 3: la descripción se recorta)
func TestNuevoDefecto_RecortaDescripcion(t *testing.T) {
	d, err := NuevoDefecto(1, 1, "  No carga la pantalla  ", SeveridadMedia)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Descripcion != "No carga la pantalla" {
		t.Errorf("Descripcion = %q, se esperaba %q", d.Descripcion, "No carga la pantalla")
	}
}

// Traza: HU-20 RN-20.4: el estado de un defecto nuevo siempre es Abierto.
func TestNuevoDefecto_SiempreQuedaAbierto(t *testing.T) {
	d, err := NuevoDefecto(1, 1, "Falla el login", SeveridadBaja)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Estado != DefectoAbierto {
		t.Errorf("Estado = %q, se esperaba %q", d.Estado, DefectoAbierto)
	}
}

// Traza: HU-20 CA-20.3
func TestNuevoDefecto_SinDescripcion_RetornaError(t *testing.T) {
	for _, descripcion := range []string{"", "   ", " \t "} {
		_, err := NuevoDefecto(1, 1, descripcion, SeveridadAlta)
		if !errors.Is(err, ErrDescripcionDefectoVacia) {
			t.Errorf("descripción %q: error = %v, se esperaba ErrDescripcionDefectoVacia", descripcion, err)
		}
	}
}

// Traza: HU-20 CA-20.4
func TestNuevoDefecto_SeveridadInvalida_RetornaError(t *testing.T) {
	for _, severidad := range []Severidad{"", "urgente", "alta", "Alta "} {
		_, err := NuevoDefecto(1, 1, "Falla el login", severidad)
		if !errors.Is(err, ErrSeveridadInvalida) {
			t.Errorf("severidad %q: error = %v, se esperaba ErrSeveridadInvalida", severidad, err)
		}
	}
}

// Traza: HU-20 CA-20.4
func TestNuevoDefecto_SeveridadesValidas(t *testing.T) {
	for _, severidad := range []Severidad{SeveridadCritica, SeveridadAlta, SeveridadMedia, SeveridadBaja} {
		d, err := NuevoDefecto(1, 1, "Falla el login", severidad)
		if err != nil {
			t.Errorf("severidad %q: error inesperado: %v", severidad, err)
			continue
		}
		if d.Severidad != severidad {
			t.Errorf("Severidad = %q, se esperaba %q", d.Severidad, severidad)
		}
	}
}

// Traza: HU-20 sección 7 (si hay varios errores, se informa el primero en el orden de la tabla:
// descripción antes que severidad).
func TestNuevoDefecto_VariosErrores_InformaElPrimero(t *testing.T) {
	_, err := NuevoDefecto(1, 1, "", "urgente")
	if !errors.Is(err, ErrDescripcionDefectoVacia) {
		t.Errorf("error = %v, se esperaba ErrDescripcionDefectoVacia", err)
	}
}
