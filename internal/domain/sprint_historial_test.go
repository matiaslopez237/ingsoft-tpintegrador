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
