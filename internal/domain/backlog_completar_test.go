package domain

import (
	"errors"
	"reflect"
	"testing"
)

// elementoEnSprint arma un elemento de prueba con un sprint asignado (nil: sin sprint) y un estado.
func elementoEnSprint(t *testing.T, sprintID *int64, estado EstadoElemento) ElementoBacklog {
	t.Helper()
	e := elementoDePrueba(t, 1)
	e.ID = 7
	sp := 3
	e.StoryPoints = &sp
	e.SprintID = sprintID
	e.Estado = estado
	return e
}

// copiaProfunda copia el elemento sin compartir punteros ni slices, para que la comparación
// detecte también cambios en *SprintID, *StoryPoints o Criterios.
func copiaProfunda(e ElementoBacklog) ElementoBacklog {
	c := e
	if e.SprintID != nil {
		id := *e.SprintID
		c.SprintID = &id
	}
	if e.StoryPoints != nil {
		sp := *e.StoryPoints
		c.StoryPoints = &sp
	}
	c.Criterios = append([]string(nil), e.Criterios...)
	return c
}

// Traza: HU-10 CA-10.1, CA-10.2 (RN-10.1 a RN-10.3)
func TestElementoBacklog_Completar_Exito(t *testing.T) {
	tests := []struct {
		nombre string
		estado EstadoElemento
	}{
		{"desde Pendiente", EstadoPendiente},
		{"desde En progreso", EstadoEnProgreso},
		{"ya Hecho no cambia nada (RN-10.3)", EstadoHecho},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			sprintID := int64(10)
			e := elementoEnSprint(t, &sprintID, tt.estado)
			esperado := copiaProfunda(e)
			esperado.Estado = EstadoHecho

			err := e.Completar(true)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			// RN-10.1: solo cambia el estado; el SprintID y el resto de los datos se conservan.
			if !reflect.DeepEqual(e, esperado) {
				t.Errorf("elemento = %+v, se esperaba %+v", e, esperado)
			}
		})
	}
}

// Traza: HU-10 CA-10.3, CA-10.4 (RN-10.2, P-8)
func TestElementoBacklog_Completar_Errores_NoCambiaElElemento(t *testing.T) {
	sprintID := int64(10)
	tests := []struct {
		nombre        string
		sprintID      *int64
		estado        EstadoElemento
		sprintAbierto bool
		errEsperado   error
	}{
		{"sin sprint, Pendiente", nil, EstadoPendiente, true, ErrElementoSinSprint},
		{"sin sprint, Hecho", nil, EstadoHecho, true, ErrElementoSinSprint},
		{"sprint cerrado, Pendiente", &sprintID, EstadoPendiente, false, ErrSprintCerrado},
		{"sprint cerrado, En progreso", &sprintID, EstadoEnProgreso, false, ErrSprintCerrado},
		{"sprint cerrado, Hecho (P-8: antes que RN-10.3)", &sprintID, EstadoHecho, false, ErrSprintCerrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			e := elementoEnSprint(t, tt.sprintID, tt.estado)
			antes := copiaProfunda(e)

			err := e.Completar(tt.sprintAbierto)

			if !errors.Is(err, tt.errEsperado) {
				t.Errorf("error = %v, se esperaba %v", err, tt.errEsperado)
			}
			if !reflect.DeepEqual(e, antes) {
				t.Errorf("elemento = %+v, no debería haber cambiado (%+v)", e, antes)
			}
		})
	}
}
