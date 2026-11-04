package domain

import (
	"errors"
	"testing"
	"time"
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

// Traza: HU-21 CA-21.1 (RN-21.2)
func TestDefecto_Iniciar_DesdeAbierto_QuedaEnProgreso(t *testing.T) {
	d := Defecto{Estado: DefectoAbierto}

	err := d.Iniciar()

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Estado != DefectoEnProgreso {
		t.Errorf("Estado = %q, se esperaba %q", d.Estado, DefectoEnProgreso)
	}
}

// Traza: HU-21 CA-21.5 (RN-21.2, RN-21.5: solo se inicia un defecto Abierto, sin retroceder ni repetir)
func TestDefecto_Iniciar_NoAbierto_RetornaError(t *testing.T) {
	for _, estado := range []EstadoDefecto{DefectoEnProgreso, DefectoResuelto, DefectoCerrado} {
		d := Defecto{Estado: estado}
		err := d.Iniciar()
		if !errors.Is(err, ErrDefectoNoAbierto) {
			t.Errorf("estado %q: error = %v, se esperaba ErrDefectoNoAbierto", estado, err)
		}
		if d.Estado != estado {
			t.Errorf("estado %q: Estado cambió a %q, no debía cambiar", estado, d.Estado)
		}
	}
}

// Traza: HU-21 CA-21.2 (RN-21.3)
func TestDefecto_Resolver_DesdeEnProgreso_QuedaResuelto(t *testing.T) {
	d := Defecto{Estado: DefectoEnProgreso}

	err := d.Resolver()

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Estado != DefectoResuelto {
		t.Errorf("Estado = %q, se esperaba %q", d.Estado, DefectoResuelto)
	}
}

// Traza: HU-21 CA-21.6 (RN-21.3, RN-21.5: solo se resuelve un defecto En progreso, no se saltea ni repite)
func TestDefecto_Resolver_NoEnProgreso_RetornaError(t *testing.T) {
	for _, estado := range []EstadoDefecto{DefectoAbierto, DefectoResuelto, DefectoCerrado} {
		d := Defecto{Estado: estado}
		err := d.Resolver()
		if !errors.Is(err, ErrDefectoNoEnProgreso) {
			t.Errorf("estado %q: error = %v, se esperaba ErrDefectoNoEnProgreso", estado, err)
		}
		if d.Estado != estado {
			t.Errorf("estado %q: Estado cambió a %q, no debía cambiar", estado, d.Estado)
		}
	}
}

// Traza: HU-21 CA-21.4 (RN-21.4)
func TestDefecto_Cerrar_DesdeResuelto_QuedaCerrado(t *testing.T) {
	d := Defecto{Estado: DefectoResuelto}

	err := d.Cerrar()

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Estado != DefectoCerrado {
		t.Errorf("Estado = %q, se esperaba %q", d.Estado, DefectoCerrado)
	}
}

// Traza: HU-21 CA-21.7 (RN-21.4, RN-21.5: solo se cierra un defecto Resuelto, no se saltea ni retrocede)
func TestDefecto_Cerrar_NoResuelto_RetornaError(t *testing.T) {
	for _, estado := range []EstadoDefecto{DefectoAbierto, DefectoEnProgreso, DefectoCerrado} {
		d := Defecto{Estado: estado}
		err := d.Cerrar()
		if !errors.Is(err, ErrDefectoNoResuelto) {
			t.Errorf("estado %q: error = %v, se esperaba ErrDefectoNoResuelto", estado, err)
		}
		if d.Estado != estado {
			t.Errorf("estado %q: Estado cambió a %q, no debía cambiar", estado, d.Estado)
		}
	}
}

// Traza: HU-21 (PuedeIniciarse/PuedeResolverse/PuedeCerrarse los usa la pantalla para decidir qué
// botón mostrar, sin repetir la regla, igual que Sprint.PuedeCerrarse en HU-11)
func TestDefecto_PuedeIniciarseResolverseYCerrarse(t *testing.T) {
	tests := []struct {
		estado          EstadoDefecto
		puedeIniciarse  bool
		puedeResolverse bool
		puedeCerrarse   bool
	}{
		{DefectoAbierto, true, false, false},
		{DefectoEnProgreso, false, true, false},
		{DefectoResuelto, false, false, true},
		{DefectoCerrado, false, false, false},
	}
	for _, tt := range tests {
		d := Defecto{Estado: tt.estado}
		if got := d.PuedeIniciarse(); got != tt.puedeIniciarse {
			t.Errorf("estado %q: PuedeIniciarse() = %v, se esperaba %v", tt.estado, got, tt.puedeIniciarse)
		}
		if got := d.PuedeResolverse(); got != tt.puedeResolverse {
			t.Errorf("estado %q: PuedeResolverse() = %v, se esperaba %v", tt.estado, got, tt.puedeResolverse)
		}
		if got := d.PuedeCerrarse(); got != tt.puedeCerrarse {
			t.Errorf("estado %q: PuedeCerrarse() = %v, se esperaba %v", tt.estado, got, tt.puedeCerrarse)
		}
	}
}

// diaSprint arma una fecha a medianoche, en la zona del servidor, igual que los sprints de HU-08.
func diaSprint(anio int, mes time.Month, d int) time.Time {
	return time.Date(anio, mes, d, 0, 0, 0, 0, time.Local)
}

// Traza: HU-20 CA-20.1 (RN-20.5)
func TestSprintVigente_SprintAbiertoConHoyEnRango_DevuelveSuID(t *testing.T) {
	sprints := []Sprint{
		{ID: 5, Estado: SprintAbierto, FechaInicio: diaSprint(2026, time.October, 20), FechaFin: diaSprint(2026, time.October, 30)},
	}
	hoy := diaSprint(2026, time.October, 25)

	id, ok := SprintVigente(sprints, hoy)

	if !ok || id != 5 {
		t.Errorf("SprintVigente = %d, %v; se esperaba 5, true", id, ok)
	}
}

// Traza: HU-20 CA-20.2 (RN-20.5: el día de inicio y el de fin son inclusive)
func TestSprintVigente_HoyEnLosLimites_DevuelveSuID(t *testing.T) {
	sprint := Sprint{ID: 9, Estado: SprintAbierto, FechaInicio: diaSprint(2026, time.October, 20), FechaFin: diaSprint(2026, time.October, 30)}
	tests := []struct {
		nombre string
		hoy    time.Time
	}{
		{"primer día", diaSprint(2026, time.October, 20)},
		{"último día", diaSprint(2026, time.October, 30)},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			id, ok := SprintVigente([]Sprint{sprint}, tt.hoy)
			if !ok || id != 9 {
				t.Errorf("SprintVigente = %d, %v; se esperaba 9, true", id, ok)
			}
		})
	}
}

// Traza: HU-20 CA-20.2
func TestSprintVigente_SinSprints_NoHayNinguno(t *testing.T) {
	id, ok := SprintVigente(nil, diaSprint(2026, time.October, 25))

	if ok || id != 0 {
		t.Errorf("SprintVigente = %d, %v; se esperaba 0, false", id, ok)
	}
}

// Traza: HU-20 CA-20.2 (límite: un sprint que ya terminó, o que todavía no empezó, no califica)
func TestSprintVigente_HoyFueraDelRango_NoHayNinguno(t *testing.T) {
	sprint := Sprint{ID: 9, Estado: SprintAbierto, FechaInicio: diaSprint(2026, time.October, 20), FechaFin: diaSprint(2026, time.October, 30)}
	tests := []struct {
		nombre string
		hoy    time.Time
	}{
		{"un día antes de empezar", diaSprint(2026, time.October, 19)},
		{"un día después de terminar", diaSprint(2026, time.October, 31)},
		{"el proyecto todavía no tiene este sprint", diaSprint(2026, time.September, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			id, ok := SprintVigente([]Sprint{sprint}, tt.hoy)
			if ok || id != 0 {
				t.Errorf("SprintVigente = %d, %v; se esperaba 0, false", id, ok)
			}
		})
	}
}

// Traza: HU-20 RN-20.5 (un sprint Cerrado no se asigna, aunque sus fechas incluyan hoy)
func TestSprintVigente_SprintCerrado_NoCalifica(t *testing.T) {
	sprints := []Sprint{
		{ID: 5, Estado: SprintCerrado, FechaInicio: diaSprint(2026, time.October, 20), FechaFin: diaSprint(2026, time.October, 30)},
	}

	id, ok := SprintVigente(sprints, diaSprint(2026, time.October, 25))

	if ok || id != 0 {
		t.Errorf("SprintVigente = %d, %v; se esperaba 0, false", id, ok)
	}
}
