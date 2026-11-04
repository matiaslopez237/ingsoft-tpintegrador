package domain

import (
	"errors"
	"math"
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

// Editar reconstruye el elemento con NuevoElementoBacklog, que no recibe el sprint. En producción no
// se notaba porque el repositorio guarda solo los campos enviados y el caso de uso relee el elemento.
// Traza: HU-19 CA-19.13 (RN-19.3, RN-06.1: lo que no se envía conserva su valor)
func TestEditar_ConservaElSprint(t *testing.T) {
	original := login(t)
	sprint := int64(3)
	original.SprintID = &sprint
	baja := PrioridadBaja

	e, err := original.Editar(CambiosElemento{Prioridad: &baja})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.SprintID == nil || *e.SprintID != 3 {
		t.Fatalf("SprintID = %v, se esperaba 3", e.SprintID)
	}
	*e.SprintID = 99
	if *original.SprintID != 3 {
		t.Errorf("se comparte memoria: SprintID del original = %d, se esperaba 3", *original.SprintID)
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

// loginConHoras es un elemento con sprint y horas estimadas, ya guardado.
func loginConHoras(t *testing.T, horas float64) ElementoBacklog {
	t.Helper()
	e := login(t)
	sprint := int64(3)
	e.SprintID = &sprint
	e.HorasEstimadas = &horas
	return e
}

// Traza: HU-19 CA-19.15 (RN-19.3)
func TestNuevoElementoBacklog_NaceSinHorasEstimadas(t *testing.T) {
	if e := login(t); e.HorasEstimadas != nil {
		t.Errorf("HorasEstimadas = %v, se esperaba nil", *e.HorasEstimadas)
	}
}

// Traza: HU-19 CA-19.13 (RN-19.3)
func TestEditar_CargarHorasEstimadas(t *testing.T) {
	original := login(t)

	e, err := original.Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: pHoras(8)})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := original
	want.HorasEstimadas = pHoras(8)
	if !reflect.DeepEqual(e, want) {
		t.Errorf("se obtuvo %+v, se esperaba %+v", e, want)
	}
}

// Traza: HU-19 CA-19.13 (RN-19.3)
func TestEditar_CambiarHorasEstimadas(t *testing.T) {
	e, err := loginConHoras(t, 8).Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: pHoras(10)})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mismasHoras(e.HorasEstimadas, pHoras(10)) {
		t.Errorf("HorasEstimadas = %v, se esperaba 10", verHoras(e.HorasEstimadas))
	}
}

// Enviarlas en nil es quitarlas, no "sin cambios": es la única forma de volver a sin estimar.
// Traza: HU-19 CA-19.13 (RN-19.3)
func TestEditar_QuitarHorasEstimadas(t *testing.T) {
	e, err := loginConHoras(t, 8).Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: nil})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.HorasEstimadas != nil {
		t.Errorf("HorasEstimadas = %v, se esperaba nil", *e.HorasEstimadas)
	}
}

// Traza: HU-19 CA-19.13, CA-19.15 (RN-06.1, RN-19.1, RN-19.3)
func TestEditar_OtrosCampos_ConservanHorasEstimadasYSprint(t *testing.T) {
	baja := PrioridadBaja
	tests := []struct {
		nombre string
		c      CambiosElemento
	}{
		{"prioridad", CambiosElemento{Prioridad: &baja}},
		{"cambiar story points", CambiosElemento{CambiarStoryPoints: true, StoryPoints: puntos(13)}},
		{"quitar story points", CambiosElemento{CambiarStoryPoints: true, StoryPoints: nil}},
		// Como los story points: sin CambiarHorasEstimadas, el valor no se envió.
		{"horas sin indicar que se envían", CambiosElemento{Prioridad: &baja, HorasEstimadas: pHoras(99)}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			e, err := loginConHoras(t, 8).Editar(tt.c)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !mismasHoras(e.HorasEstimadas, pHoras(8)) {
				t.Errorf("HorasEstimadas = %v, se esperaba 8", verHoras(e.HorasEstimadas))
			}
			if e.SprintID == nil || *e.SprintID != 3 {
				t.Errorf("SprintID = %v, se esperaba 3", e.SprintID)
			}
		})
	}
}

// Traza: HU-19 CA-19.15 (RN-19.1): las horas tampoco cambian los story points.
func TestEditar_HorasEstimadas_ConservanStoryPoints(t *testing.T) {
	e, err := login(t).Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: pHoras(8)})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.StoryPoints == nil || *e.StoryPoints != 5 {
		t.Errorf("StoryPoints = %v, se esperaba 5", e.StoryPoints)
	}
}

// Traza: HU-19 CA-19.14 (RN-19.2)
func TestEditar_HorasEstimadasEnLosExtremos(t *testing.T) {
	for _, horas := range []float64{0.5, 1000} {
		e, err := login(t).Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: pHoras(horas)})
		if err != nil {
			t.Errorf("%v: error inesperado: %v", horas, err)
			continue
		}
		if !mismasHoras(e.HorasEstimadas, pHoras(horas)) {
			t.Errorf("HorasEstimadas = %v, se esperaba %v", verHoras(e.HorasEstimadas), horas)
		}
	}
}

// Si hay error, el elemento no cambia.
// Traza: HU-19 CA-19.14 (RN-19.2)
func TestEditar_HorasEstimadasInvalidas_RetornaError(t *testing.T) {
	tests := []struct {
		horas   float64
		wantErr error
	}{
		{0, ErrHorasInvalidas},
		{-1, ErrHorasInvalidas},
		{math.NaN(), ErrHorasInvalidas},
		{1.25, ErrHorasNoMultiploMediaHora},
		{1000.5, ErrHorasEstimadasExcedenMaximo},
		{math.Inf(1), ErrHorasEstimadasExcedenMaximo},
	}
	for _, tt := range tests {
		original := loginConHoras(t, 8)

		e, err := original.Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: pHoras(tt.horas)})

		if !errors.Is(err, tt.wantErr) {
			t.Errorf("%v: error = %v, se esperaba %v", tt.horas, err, tt.wantErr)
		}
		if !reflect.DeepEqual(e, ElementoBacklog{}) {
			t.Errorf("%v: con error se devolvió %+v, se esperaba un elemento vacío", tt.horas, e)
		}
		if !mismasHoras(original.HorasEstimadas, pHoras(8)) {
			t.Errorf("%v: el original cambió: HorasEstimadas = %v", tt.horas, verHoras(original.HorasEstimadas))
		}
	}
}

// Los errores de HU-05 y HU-06 se informan antes que el de las horas, que es el último en sumarse.
// Traza: HU-19 CA-19.14 (RN-19.2; spec HU-06, sección 7)
func TestEditar_HorasEstimadasInvalidas_SeInformanAlFinal(t *testing.T) {
	cerrado := EstadoElemento("Cerrado")
	tests := []struct {
		nombre  string
		c       CambiosElemento
		wantErr error
	}{
		{"título", CambiosElemento{Titulo: texto(""), CambiarHorasEstimadas: true, HorasEstimadas: pHoras(0)}, ErrTituloVacio},
		{"estado", CambiosElemento{Estado: &cerrado, CambiarHorasEstimadas: true, HorasEstimadas: pHoras(0)}, ErrEstadoElementoInvalido},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if _, err := login(t).Editar(tt.c); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// La copia editada no comparte las horas con el original ni con los cambios.
// Traza: HU-19 CA-19.13 (RN-19.3)
func TestEditar_HorasEstimadas_NoCompartenMemoria(t *testing.T) {
	enviadas := pHoras(10)
	original := loginConHoras(t, 8)

	e, err := original.Editar(CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: enviadas})
	if err != nil {
		t.Fatal(err)
	}
	*enviadas = 99
	if !mismasHoras(e.HorasEstimadas, pHoras(10)) {
		t.Errorf("se comparte memoria con los cambios: HorasEstimadas = %v, se esperaba 10", verHoras(e.HorasEstimadas))
	}

	baja := PrioridadBaja
	sinTocar, err := original.Editar(CambiosElemento{Prioridad: &baja})
	if err != nil {
		t.Fatal(err)
	}
	*sinTocar.HorasEstimadas = 99
	if !mismasHoras(original.HorasEstimadas, pHoras(8)) {
		t.Errorf("se comparte memoria con el original: HorasEstimadas = %v, se esperaba 8", verHoras(original.HorasEstimadas))
	}
}
