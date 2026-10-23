package domain

import (
	"errors"
	"slices"
	"testing"
)

func proyectoDePrueba(t *testing.T) *Proyecto {
	t.Helper()
	p, err := NuevoProyecto("TP Integrador", fecha(t, "2026-03-01"), nil)
	if err != nil {
		t.Fatalf("no se pudo crear el proyecto de prueba: %v", err)
	}
	return &p
}

// Traza: HU-03 CA-03.1, RN-03.2
func TestProyecto_AgregarIntegrante_Valido(t *testing.T) {
	p := proyectoDePrueba(t)

	if err := p.AgregarIntegrante("Ana", "Desarrolladora"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if err := p.AgregarIntegrante("Beto", "Tester"); err != nil {
		t.Fatalf("error inesperado al agregar un segundo integrante: %v", err)
	}

	if len(p.Integrantes) != 2 {
		t.Fatalf("Integrantes = %+v, se esperaban 2", p.Integrantes)
	}
	if got := p.Integrantes[0]; got.Nombre != "Ana" || !slices.Equal(got.Roles, []string{"Desarrolladora"}) {
		t.Errorf("Integrantes[0] = %+v, se esperaba Ana/Desarrolladora", got)
	}
	if got := p.Integrantes[1]; got.Nombre != "Beto" || !slices.Equal(got.Roles, []string{"Tester"}) {
		t.Errorf("Integrantes[1] = %+v, se esperaba Beto/Tester", got)
	}
}

// Traza: HU-03 CA-03.2, CA-03.5, CA-03.6, RN-03.1, RN-03.3, RN-03.4
func TestProyecto_AgregarIntegrante_Invalido_RetornaError(t *testing.T) {
	tests := []struct {
		nombreCaso string
		preparar   func(p *Proyecto)
		nombreArg  string
		rolArg     string
		wantErr    error
	}{
		{"nombre vacío", nil, "", "Desarrolladora", ErrNombreIntegranteVacio},
		{"nombre de solo espacios", nil, "   ", "Desarrolladora", ErrNombreIntegranteVacio},
		{"rol vacío", nil, "Ana", "", ErrRolIntegranteVacio},
		{
			"nombre duplicado exacto",
			func(p *Proyecto) { _ = p.AgregarIntegrante("Ana", "Desarrolladora") },
			"Ana", "Tester", ErrIntegranteDuplicado,
		},
		{
			"nombre duplicado con mayúsculas y espacios",
			func(p *Proyecto) { _ = p.AgregarIntegrante("Ana", "Desarrolladora") },
			"  ANA  ", "Tester", ErrIntegranteDuplicado,
		},
	}
	for _, tt := range tests {
		t.Run(tt.nombreCaso, func(t *testing.T) {
			p := proyectoDePrueba(t)
			if tt.preparar != nil {
				tt.preparar(p)
			}
			err := p.AgregarIntegrante(tt.nombreArg, tt.rolArg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-03 CA-03.3
func TestProyecto_QuitarIntegrante_Existente(t *testing.T) {
	p := proyectoDePrueba(t)
	_ = p.AgregarIntegrante("Ana", "Desarrolladora")
	_ = p.AgregarIntegrante("Beto", "Tester")

	if err := p.QuitarIntegrante("Ana"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(p.Integrantes) != 1 || p.Integrantes[0].Nombre != "Beto" {
		t.Fatalf("Integrantes = %+v, se esperaba solo Beto", p.Integrantes)
	}
}

// Traza: HU-03 CA-03.3 (caso límite: queda vacío)
func TestProyecto_QuitarIntegrante_UnicoRestante_DejaListaVacia(t *testing.T) {
	p := proyectoDePrueba(t)
	_ = p.AgregarIntegrante("Ana", "Desarrolladora")

	if err := p.QuitarIntegrante("Ana"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(p.Integrantes) != 0 {
		t.Fatalf("Integrantes = %+v, se esperaba lista vacía", p.Integrantes)
	}
}

// Traza: HU-03 CA-03.4
func TestProyecto_QuitarIntegrante_Inexistente_RetornaError(t *testing.T) {
	p := proyectoDePrueba(t)
	_ = p.AgregarIntegrante("Ana", "Desarrolladora")

	err := p.QuitarIntegrante("Caro")

	if !errors.Is(err, ErrIntegranteNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrIntegranteNoEncontrado)
	}
	if len(p.Integrantes) != 1 {
		t.Errorf("Integrantes = %+v, no debería haber cambiado", p.Integrantes)
	}
}

// Traza: HU-03 CA-03.7, RN-03.5
func TestProyecto_AgregarRol_SumaElRolSinDuplicarAlIntegrante(t *testing.T) {
	p := proyectoDePrueba(t)
	_ = p.AgregarIntegrante("Ana", "Desarrolladora")

	if err := p.AgregarRol("  ana ", "  Tester "); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(p.Integrantes) != 1 {
		t.Fatalf("Integrantes = %+v, Ana debería figurar una sola vez", p.Integrantes)
	}
	if got, want := p.Integrantes[0].Roles, []string{"Desarrolladora", "Tester"}; !slices.Equal(got, want) {
		t.Errorf("Roles = %v, se esperaba %v", got, want)
	}
}

// Traza: HU-03 CA-03.8, CA-03.9, CA-03.10, RN-03.4, RN-03.6
func TestProyecto_AgregarRol_Invalido_RetornaErrorYNoCambia(t *testing.T) {
	tests := []struct {
		nombreCaso string
		nombreArg  string
		rolArg     string
		wantErr    error
	}{
		{"rol duplicado con otras mayúsculas", "Ana", "desarrolladora", ErrRolDuplicado},
		{"rol duplicado con espacios", "Ana", "  Desarrolladora  ", ErrRolDuplicado},
		{"integrante inexistente", "Caro", "Tester", ErrIntegranteNoEncontrado},
		{"rol vacío", "Ana", "", ErrRolIntegranteVacio},
		{"rol de solo espacios", "Ana", "   ", ErrRolIntegranteVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombreCaso, func(t *testing.T) {
			p := proyectoDePrueba(t)
			_ = p.AgregarIntegrante("Ana", "Desarrolladora")

			err := p.AgregarRol(tt.nombreArg, tt.rolArg)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if got, want := p.Integrantes[0].Roles, []string{"Desarrolladora"}; !slices.Equal(got, want) {
				t.Errorf("Roles = %v, no debería haber cambiado (%v)", got, want)
			}
		})
	}
}
