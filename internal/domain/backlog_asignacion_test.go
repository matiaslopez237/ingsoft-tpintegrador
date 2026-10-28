package domain

import (
	"errors"
	"testing"
)

func elementoDePrueba(t *testing.T, proyectoID int64) ElementoBacklog {
	t.Helper()
	e, err := NuevoElementoBacklog(proyectoID, "Login de usuarios", "Descripción de prueba",
		PrioridadMedia, nil, []string{"Criterio de prueba"})
	if err != nil {
		t.Fatalf("no se pudo crear el elemento de prueba: %v", err)
	}
	return e
}

// Traza: HU-09 CA-09.1
func TestElementoBacklog_AsignarASprint_SinSprintPrevio_Asigna(t *testing.T) {
	e := elementoDePrueba(t, 1)

	err := e.AsignarASprint(10, 1, true, false)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.SprintID == nil || *e.SprintID != 10 {
		t.Errorf("SprintID = %v, se esperaba 10", e.SprintID)
	}
}

// Traza: HU-09 CA-09.1 (caso límite, RN-09.3)
func TestElementoBacklog_AsignarASprint_MismoSprintQueYaTiene_NoEsError(t *testing.T) {
	e := elementoDePrueba(t, 1)
	if err := e.AsignarASprint(10, 1, true, false); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	// El sprint actual (10) sigue activo, pero reasignar al mismo sprint no debe fallar.
	err := e.AsignarASprint(10, 1, true, true)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.SprintID == nil || *e.SprintID != 10 {
		t.Errorf("SprintID = %v, se esperaba 10", e.SprintID)
	}
}

// Traza: HU-09 CA-09.2, RN-09.2
func TestElementoBacklog_AsignarASprint_YaAsignadoAOtroSprintActivo_RetornaError(t *testing.T) {
	e := elementoDePrueba(t, 1)
	if err := e.AsignarASprint(10, 1, true, false); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	err := e.AsignarASprint(20, 1, true, true)

	if !errors.Is(err, ErrElementoYaAsignado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrElementoYaAsignado)
	}
	if e.SprintID == nil || *e.SprintID != 10 {
		t.Errorf("SprintID = %v, no debería haber cambiado (se esperaba 10)", e.SprintID)
	}
}

// Traza: HU-09 CA-09.3, RN-09.2
func TestElementoBacklog_AsignarASprint_SprintPrevioYaCerrado_Reasigna(t *testing.T) {
	e := elementoDePrueba(t, 1)
	if err := e.AsignarASprint(10, 1, true, false); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	// El sprint actual (10) ya no está activo (se cerró).
	err := e.AsignarASprint(20, 1, true, false)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.SprintID == nil || *e.SprintID != 20 {
		t.Errorf("SprintID = %v, se esperaba 20", e.SprintID)
	}
}

// Traza: HU-09 CA-09.7, RN-09.5
func TestElementoBacklog_AsignarASprint_SprintDeOtroProyecto_RetornaError(t *testing.T) {
	e := elementoDePrueba(t, 1)

	err := e.AsignarASprint(10, 2, true, false)

	if !errors.Is(err, ErrSprintDeOtroProyecto) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSprintDeOtroProyecto)
	}
	if e.SprintID != nil {
		t.Errorf("SprintID = %v, no debería haberse asignado", e.SprintID)
	}
}

// Traza: HU-09 CA-09.8, RN-09.6
func TestElementoBacklog_AsignarASprint_SprintDestinoCerrado_RetornaError(t *testing.T) {
	e := elementoDePrueba(t, 1)

	err := e.AsignarASprint(10, 1, false, false)

	if !errors.Is(err, ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSprintCerrado)
	}
	if e.SprintID != nil {
		t.Errorf("SprintID = %v, no debería haberse asignado", e.SprintID)
	}
}

// Traza: HU-09 CA-09.4
func TestElementoBacklog_DesasignarDeSprint_SprintActivo_Desasigna(t *testing.T) {
	e := elementoDePrueba(t, 1)
	if err := e.AsignarASprint(10, 1, true, false); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	err := e.DesasignarDeSprint(true)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.SprintID != nil {
		t.Errorf("SprintID = %v, se esperaba nil", e.SprintID)
	}
}

// Traza: HU-09 CA-09.5, RN-09.4
func TestElementoBacklog_DesasignarDeSprint_SprintYaCerrado_RetornaError(t *testing.T) {
	e := elementoDePrueba(t, 1)
	if err := e.AsignarASprint(10, 1, true, false); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	err := e.DesasignarDeSprint(false)

	if !errors.Is(err, ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSprintCerrado)
	}
	if e.SprintID == nil || *e.SprintID != 10 {
		t.Errorf("SprintID = %v, no debería haber cambiado (se esperaba 10)", e.SprintID)
	}
}

// Traza: HU-09 CA-09.6
func TestElementoBacklog_DesasignarDeSprint_SinSprintAsignado_RetornaError(t *testing.T) {
	e := elementoDePrueba(t, 1)

	err := e.DesasignarDeSprint(true)

	if !errors.Is(err, ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, ErrElementoSinSprint)
	}
}
