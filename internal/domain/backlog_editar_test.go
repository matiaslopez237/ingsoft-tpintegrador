package domain

import (
	"errors"
	"reflect"
	"testing"
)

func login(t *testing.T) ElementoBacklog {
	t.Helper()
	e, err := NuevoElementoBacklog(1, "Login", "Ingresar al sistema", PrioridadAlta, puntos(5), []string{"Con email"})
	if err != nil {
		t.Fatal(err)
	}
	e.ID = 7
	return e
}

// Traza: HU-06 CA-06.1 (RN-06.1)
func TestEditar_SoloPrioridad_ElRestoNoCambia(t *testing.T) {
	original := login(t)
	baja := PrioridadBaja

	e, err := original.Editar(CambiosElemento{Prioridad: &baja})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := original
	want.Prioridad = PrioridadBaja
	if !reflect.DeepEqual(e, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", e, want)
	}
}

// Traza: HU-06 CA-06.2 (RN-06.1, RN-06.4, RN-06.9)
func TestEditar_TodosLosCampos(t *testing.T) {
	enProgreso := EstadoEnProgreso
	media := PrioridadMedia
	criterios := []string{" Uno ", "Dos"}

	e, err := login(t).Editar(CambiosElemento{
		Titulo:             texto("  Inicio de sesión "),
		Descripcion:        texto("Entrar"),
		Prioridad:          &media,
		CambiarStoryPoints: true,
		StoryPoints:        puntos(8),
		Criterios:          &criterios,
		Estado:             &enProgreso,
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := ElementoBacklog{
		ID: 7, ProyectoID: 1, Titulo: "Inicio de sesión", Descripcion: "Entrar", Prioridad: PrioridadMedia,
		Estado: EstadoEnProgreso, StoryPoints: puntos(8), Criterios: []string{"Uno", "Dos"},
	}
	if !reflect.DeepEqual(e, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", e, want)
	}
}

// Traza: HU-06 CA-06.3 (RN-06.3)
func TestEditar_StoryPointsNil_QuedaSinEstimar(t *testing.T) {
	e, err := login(t).Editar(CambiosElemento{CambiarStoryPoints: true, StoryPoints: nil})
	if err != nil || e.StoryPoints != nil {
		t.Errorf("se obtuvo %v, %v; se esperaba un elemento sin estimar", e.StoryPoints, err)
	}
}

// Traza: HU-06 CA-06.2 (RN-06.4): el estado cambia en cualquier sentido.
func TestEditar_EstadoEnCualquierSentido(t *testing.T) {
	for _, desde := range []EstadoElemento{EstadoPendiente, EstadoEnProgreso, EstadoHecho} {
		for _, hacia := range []EstadoElemento{EstadoPendiente, EstadoEnProgreso, EstadoHecho} {
			e := login(t)
			e.Estado = desde
			h := hacia
			got, err := e.Editar(CambiosElemento{Estado: &h})
			if err != nil || got.Estado != hacia {
				t.Errorf("%s → %s: se obtuvo %q, %v", desde, hacia, got.Estado, err)
			}
		}
	}
}

// Traza: HU-06 CA-06.4, CA-06.6 (RN-06.2, RN-06.5)
func TestEditar_DatosInvalidos_RetornaError(t *testing.T) {
	cero, vacios := 0, []string{}
	urgente, cerrado, minuscula := Prioridad("Urgente"), EstadoElemento("Cerrado"), EstadoElemento("hecho")
	tests := []struct {
		nombre  string
		c       CambiosElemento
		wantErr error
	}{
		{"sin cambios", CambiosElemento{}, ErrSinCambios},
		{"título vacío", CambiosElemento{Titulo: texto("  ")}, ErrTituloVacio},
		{"descripción vacía", CambiosElemento{Descripcion: texto("")}, ErrDescripcionVacia},
		{"prioridad inválida", CambiosElemento{Prioridad: &urgente}, ErrPrioridadInvalida},
		{"story points en cero", CambiosElemento{CambiarStoryPoints: true, StoryPoints: &cero}, ErrStoryPointsInvalidos},
		{"criterios vacíos", CambiosElemento{Criterios: &vacios}, ErrSinCriterios},
		{"estado inexistente", CambiosElemento{Estado: &cerrado}, ErrEstadoElementoInvalido},
		{"estado en minúscula", CambiosElemento{Estado: &minuscula}, ErrEstadoElementoInvalido},
		{"los campos de HU-05 se informan antes que el estado", CambiosElemento{Titulo: texto(""), Estado: &cerrado}, ErrTituloVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			original := login(t)
			_, err := original.Editar(tt.c)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// La copia editada no comparte memoria con el original ni con los cambios.
func TestEditar_NoModificaElOriginal(t *testing.T) {
	original := login(t)
	criterios := []string{"Nuevo"}

	e, err := original.Editar(CambiosElemento{Criterios: &criterios})
	if err != nil {
		t.Fatal(err)
	}
	criterios[0] = "Cambiado"
	*e.StoryPoints = 99

	if original.Criterios[0] != "Con email" || *original.StoryPoints != 5 || e.Criterios[0] != "Nuevo" {
		t.Errorf("se comparte memoria: original %+v, editado %+v", original, e)
	}
}
