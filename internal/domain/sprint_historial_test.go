package domain

import (
	"reflect"
	"testing"
)

// sprintEn arma un sprint guardado del proyecto 7, con su número, fechas y estado.
func sprintEn(t *testing.T, numero int, inicio, fin string, estado EstadoSprint) Sprint {
	t.Helper()
	return Sprint{ID: int64(numero), ProyectoID: 7, Numero: numero, FechaInicio: fecha(t, inicio), FechaFin: fecha(t, fin),
		Objetivo: "Objetivo", Estado: estado}
}

func numeros(sprints []Sprint) []int {
	ns := []int{}
	for _, s := range sprints {
		ns = append(ns, s.Numero)
	}
	return ns
}

// Traza: HU-12 CA-12.1 a CA-12.3, RN-12.1, RN-12.3
func TestSprintsAnteriores(t *testing.T) {
	tests := []struct {
		name    string
		sprints []Sprint
		want    []int
	}{
		{
			name: "sin sprints",
			want: []int{},
		},
		{
			name: "CA-12.1: solo los cerrados, del más reciente al más antiguo",
			sprints: []Sprint{
				sprintEn(t, 1, "2026-10-01", "2026-10-12", SprintCerrado),
				sprintEn(t, 2, "2026-10-13", "2026-10-25", SprintCerrado),
				sprintEn(t, 3, "2026-10-26", "2026-11-08", SprintAbierto),
			},
			want: []int{2, 1},
		},
		{
			name: "CA-12.2: los abiertos no aparecen, aunque estén vencidos",
			sprints: []Sprint{
				sprintEn(t, 1, "2020-01-01", "2020-01-14", SprintAbierto),
				sprintEn(t, 2, "2026-10-26", "2026-11-08", SprintAbierto),
			},
			want: []int{},
		},
		{
			name: "CA-12.3: el orden sigue la fecha de inicio, no el número",
			sprints: []Sprint{
				sprintEn(t, 1, "2026-10-13", "2026-10-25", SprintCerrado),
				sprintEn(t, 2, "2026-10-01", "2026-10-12", SprintCerrado),
				sprintEn(t, 3, "2026-10-26", "2026-11-08", SprintCerrado),
			},
			want: []int{3, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SprintsAnteriores(tt.sprints)

			if !reflect.DeepEqual(numeros(got), tt.want) {
				t.Errorf("SprintsAnteriores() = números %v, se esperaba %v", numeros(got), tt.want)
			}
		})
	}
}

// Traza: HU-12 RN-12.2: cada sprint se devuelve con todos sus datos, incluidos los puntos que
// registró el cierre (HU-11), sin recalcularlos.
func TestSprintsAnteriores_ConservaLosDatos(t *testing.T) {
	s := sprintEn(t, 1, "2026-10-01", "2026-10-12", SprintCerrado)
	s.Objetivo, s.PuntosPlanificados, s.PuntosCompletados = "Armar el MVP", 13, 10

	got := SprintsAnteriores([]Sprint{s})

	if len(got) != 1 || got[0] != s {
		t.Errorf("SprintsAnteriores() = %+v, se esperaba [%+v]", got, s)
	}
}

// Traza: HU-12: ordenar no modifica la lista recibida (la sigue usando la pantalla de HU-08).
func TestSprintsAnteriores_NoModificaLaLista(t *testing.T) {
	sprints := []Sprint{
		sprintEn(t, 1, "2026-10-01", "2026-10-12", SprintCerrado),
		sprintEn(t, 2, "2026-10-13", "2026-10-25", SprintCerrado),
	}
	antes := append([]Sprint(nil), sprints...)

	SprintsAnteriores(sprints)

	if !reflect.DeepEqual(sprints, antes) {
		t.Errorf("la lista cambió: %v, se esperaba %v", numeros(sprints), numeros(antes))
	}
}

// sprintCerradoCon arma un sprint Cerrado con los puntos completados indicados (los planificados
// no le importan a la velocidad). Implementa HU-22.
func sprintCerradoCon(t *testing.T, numero, completados int) Sprint {
	t.Helper()
	s := sprintEn(t, numero, "2026-10-01", "2026-10-12", SprintCerrado)
	s.PuntosCompletados = completados
	return s
}

// Traza: HU-22 CA-22.2 a CA-22.5
func TestVelocidad(t *testing.T) {
	tests := []struct {
		name          string
		cerrados      []Sprint
		wantVelocidad float64
		wantHayDatos  bool
	}{
		{
			name:         "CA-22.3: sin sprints cerrados, no hay velocidad",
			cerrados:     nil,
			wantHayDatos: false,
		},
		{
			name:          "CA-22.2: promedio de varios sprints cerrados",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 8), sprintCerradoCon(t, 2, 6), sprintCerradoCon(t, 3, 4)},
			wantVelocidad: 6,
			wantHayDatos:  true,
		},
		{
			name:          "CA-22.4: el promedio se informa con decimales",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 5), sprintCerradoCon(t, 2, 4)},
			wantVelocidad: 4.5,
			wantHayDatos:  true,
		},
		{
			name:          "CA-22.5: un único sprint cerrado con 0 puntos da velocidad 0, no 'sin datos'",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 0)},
			wantVelocidad: 0,
			wantHayDatos:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			velocidad, hayDatos := Velocidad(tt.cerrados)

			if hayDatos != tt.wantHayDatos {
				t.Fatalf("hayDatos = %v, se esperaba %v", hayDatos, tt.wantHayDatos)
			}
			if hayDatos && velocidad != tt.wantVelocidad {
				t.Errorf("Velocidad() = %v, se esperaba %v", velocidad, tt.wantVelocidad)
			}
		})
	}
}

// Traza: HU-22 RN-22.2 (un sprint Abierto no cuenta): Velocidad confía en que ya le pasaron solo
// los cerrados (como hace SprintsAnteriores), así que si se le pasa uno Abierto lo cuenta igual;
// este test documenta esa responsabilidad del llamador, no una validación de Velocidad.
func TestVelocidad_NoFiltraPorEstado(t *testing.T) {
	abierto := sprintEn(t, 1, "2026-10-01", "2026-10-12", SprintAbierto)
	abierto.PuntosCompletados = 99 // no debería pasar en la práctica (RN-11.6: un sprint Abierto vale 0)

	velocidad, hayDatos := Velocidad([]Sprint{abierto})

	if !hayDatos || velocidad != 99 {
		t.Errorf("Velocidad() = %v, %v; se esperaba 99, true (no filtra por estado)", velocidad, hayDatos)
	}
}
