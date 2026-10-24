package domain

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func puntos(n int) *int { return &n }

// Traza: HU-05 CA-05.1
func TestNuevoElementoBacklog_DatosValidos_CreaPendiente(t *testing.T) {
	e, err := NuevoElementoBacklog(7, "  Login ", " Ingresar al sistema ", PrioridadAlta, puntos(5),
		[]string{" Ingresa con email ", "Error si la clave es incorrecta"})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := ElementoBacklog{
		ProyectoID:  7,
		Titulo:      "Login",
		Descripcion: "Ingresar al sistema",
		Prioridad:   PrioridadAlta,
		Estado:      EstadoPendiente,
		StoryPoints: puntos(5),
		Criterios:   []string{"Ingresa con email", "Error si la clave es incorrecta"},
	}
	if !reflect.DeepEqual(e, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", e, want)
	}
}

// Traza: HU-05 CA-05.2, CA-05.4, CA-05.7
func TestNuevoElementoBacklog_CasosLimite_SonValidos(t *testing.T) {
	tests := []struct {
		nombre string
		titulo string
		sp     *int
	}{
		{"sin story points", "Login", nil},
		{"un story point", "Login", puntos(1)},
		{"título de 100 caracteres", strings.Repeat("a", 100), nil},
		{"título de 100 runas con acentos", strings.Repeat("ñ", 100), nil},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			e, err := NuevoElementoBacklog(1, tt.titulo, "Descripción", PrioridadBaja, tt.sp, []string{"Criterio"})
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !reflect.DeepEqual(e.StoryPoints, tt.sp) {
				t.Errorf("StoryPoints = %v, se esperaba %v", e.StoryPoints, tt.sp)
			}
		})
	}
}

// Traza: HU-05 CA-05.4, CA-05.5, CA-05.6, CA-05.7, CA-05.8
func TestNuevoElementoBacklog_DatosInvalidos_RetornaError(t *testing.T) {
	type entrada struct {
		titulo, descripcion string
		prioridad           Prioridad
		sp                  *int
		criterios           []string
	}
	valida := entrada{"Login", "Descripción", PrioridadMedia, puntos(3), []string{"Criterio"}}
	con := func(cambiar func(*entrada)) entrada {
		e := valida
		cambiar(&e)
		return e
	}
	tests := []struct {
		nombre  string
		in      entrada
		wantErr error
	}{
		{"título vacío", con(func(e *entrada) { e.titulo = "" }), ErrTituloVacio},
		{"título de espacios", con(func(e *entrada) { e.titulo = " \t " }), ErrTituloVacio},
		{"título de 101 caracteres", con(func(e *entrada) { e.titulo = strings.Repeat("a", 101) }), ErrTituloDemasiadoLargo},
		{"descripción vacía", con(func(e *entrada) { e.descripcion = "   " }), ErrDescripcionVacia},
		{"prioridad en minúscula", con(func(e *entrada) { e.prioridad = "alta" }), ErrPrioridadInvalida},
		{"prioridad inexistente", con(func(e *entrada) { e.prioridad = "Urgente" }), ErrPrioridadInvalida},
		{"prioridad vacía", con(func(e *entrada) { e.prioridad = "" }), ErrPrioridadInvalida},
		{"cero story points", con(func(e *entrada) { e.sp = puntos(0) }), ErrStoryPointsInvalidos},
		{"story points negativos", con(func(e *entrada) { e.sp = puntos(-2) }), ErrStoryPointsInvalidos},
		{"sin criterios", con(func(e *entrada) { e.criterios = nil }), ErrSinCriterios},
		{"un criterio vacío", con(func(e *entrada) { e.criterios = []string{"Bien", "  "} }), ErrCriterioVacio},
		{"el título se informa antes que la prioridad", con(func(e *entrada) { e.titulo = ""; e.prioridad = "x" }), ErrTituloVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoElementoBacklog(1, tt.in.titulo, tt.in.descripcion, tt.in.prioridad, tt.in.sp, tt.in.criterios)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// El elemento no comparte memoria con la lista ni los story points que recibió.
func TestNuevoElementoBacklog_NoCompartePunteros(t *testing.T) {
	sp := 5
	criterios := []string{"Uno"}
	e, err := NuevoElementoBacklog(1, "Login", "Descripción", PrioridadAlta, &sp, criterios)
	if err != nil {
		t.Fatal(err)
	}
	sp, criterios[0] = 99, "Cambiado"
	if *e.StoryPoints != 5 || e.Criterios[0] != "Uno" {
		t.Errorf("el elemento cambió al modificar las entradas: %+v", e)
	}
}

// Traza: HU-05 RN-05.5, RN-05.8
func TestPrioridadesYEstados_SonLosDeLaSpec(t *testing.T) {
	if got := []Prioridad{PrioridadAlta, PrioridadMedia, PrioridadBaja}; !reflect.DeepEqual(got, []Prioridad{"Alta", "Media", "Baja"}) {
		t.Errorf("prioridades = %v", got)
	}
	got := []EstadoElemento{EstadoPendiente, EstadoEnProgreso, EstadoHecho}
	if !reflect.DeepEqual(got, []EstadoElemento{"Pendiente", "En progreso", "Hecho"}) {
		t.Errorf("estados = %v", got)
	}
}
