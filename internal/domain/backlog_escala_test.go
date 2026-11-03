package domain

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func crearConPuntos(sp int) (ElementoBacklog, error) {
	return NuevoElementoBacklog(1, "Login", "Ingresar", PrioridadAlta, &sp, []string{"Con email"})
}

// Traza: HU-13 CA-13.1 (RN-13.1)
func TestNuevoElementoBacklog_ValoresDeLaEscala(t *testing.T) {
	for _, sp := range []int{1, 2, 3, 5, 8, 13, 21} {
		e, err := crearConPuntos(sp)
		if err != nil || e.StoryPoints == nil || *e.StoryPoints != sp {
			t.Errorf("%d story points: %v, %v; se esperaba que se acepten", sp, e.StoryPoints, err)
		}
	}
}

// Traza: HU-13 CA-13.2 (RN-13.1): fuera de la escala, con el mensaje de la escala.
func TestNuevoElementoBacklog_FueraDeLaEscala(t *testing.T) {
	for _, sp := range []int{4, 6, 10, 34, 0, -3} {
		_, err := crearConPuntos(sp)
		if !errors.Is(err, ErrStoryPointsInvalidos) {
			t.Errorf("%d story points: error = %v, se esperaba %v", sp, err, ErrStoryPointsInvalidos)
		}
	}
	if !strings.Contains(ErrStoryPointsInvalidos.Error(), "1, 2, 3, 5, 8, 13 o 21") {
		t.Errorf("el mensaje %q no menciona la escala", ErrStoryPointsInvalidos)
	}
}

// Traza: HU-13 CA-13.3 (RN-13.1)
func TestEditar_StoryPointsFueraDeLaEscala(t *testing.T) {
	diez, trece := 10, 13
	if _, err := login(t).Editar(CambiosElemento{CambiarStoryPoints: true, StoryPoints: &diez}); !errors.Is(err, ErrStoryPointsInvalidos) {
		t.Errorf("editar a 10: error = %v, se esperaba %v", err, ErrStoryPointsInvalidos)
	}
	e, err := login(t).Editar(CambiosElemento{CambiarStoryPoints: true, StoryPoints: &trece})
	if err != nil || *e.StoryPoints != 13 {
		t.Errorf("editar a 13: %v, %v", e.StoryPoints, err)
	}
}

// Traza: HU-13 CA-13.5 (RN-13.3): un valor guardado antes de la escala se respeta si no se envían
// story points, y se valida si se envían.
func TestEditar_ValorAnteriorALaEscala(t *testing.T) {
	viejo := login(t)
	cuatro := 4
	viejo.StoryPoints = &cuatro // así quedó un elemento cargado antes de HU-13

	titulo := "Ingreso"
	e, err := viejo.Editar(CambiosElemento{Titulo: &titulo})
	if err != nil || e.StoryPoints == nil || *e.StoryPoints != 4 || e.Titulo != "Ingreso" {
		t.Errorf("editar el título: %+v, %v; se esperaba conservar los 4 story points", e, err)
	}
	*e.StoryPoints = 99
	if *viejo.StoryPoints != 4 {
		t.Error("la copia editada comparte los story points con el original")
	}

	otroCuatro := 4
	if _, err := viejo.Editar(CambiosElemento{CambiarStoryPoints: true, StoryPoints: &otroCuatro}); !errors.Is(err, ErrStoryPointsInvalidos) {
		t.Errorf("enviar 4: error = %v, se esperaba %v", err, ErrStoryPointsInvalidos)
	}
	if e, err := viejo.Editar(CambiosElemento{CambiarStoryPoints: true}); err != nil || e.StoryPoints != nil {
		t.Errorf("quitar la estimación: %v, %v", e.StoryPoints, err)
	}
}

// Traza: HU-13 RN-13.4: la escala, en orden, para las listas de las pantallas.
func TestEscalaStoryPoints(t *testing.T) {
	got := EscalaStoryPoints()
	if want := []int{1, 2, 3, 5, 8, 13, 21}; !slices.Equal(got, want) {
		t.Fatalf("EscalaStoryPoints() = %v, se esperaba %v", got, want)
	}
	got[0] = 99
	if EscalaStoryPoints()[0] != 1 {
		t.Error("modificar la lista devuelta cambió la escala")
	}
}
