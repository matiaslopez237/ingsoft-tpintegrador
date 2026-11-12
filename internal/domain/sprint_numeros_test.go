package domain

import "testing"

// Traza: DEF-05: el ID de un sprint es único en toda la base y su número es el de su proyecto
// (RN-08.7, RN-08.9). Las pantallas muestran el número.
func TestNumerosDeSprint_DevuelveElNumeroPorID(t *testing.T) {
	sprints := []Sprint{{ID: 4, Numero: 1}, {ID: 9, Numero: 2}, {ID: 12, Numero: 3}}

	got := NumerosDeSprint(sprints)

	want := map[int64]int{4: 1, 9: 2, 12: 3}
	if len(got) != len(want) {
		t.Fatalf("NumerosDeSprint = %v, se esperaba %v", got, want)
	}
	for id, n := range want {
		if got[id] != n {
			t.Errorf("sprint %d: número %d, se esperaba %d", id, got[id], n)
		}
	}
}

// Traza: DEF-05: un proyecto sin sprints no tiene números.
func TestNumerosDeSprint_SinSprints(t *testing.T) {
	if got := NumerosDeSprint(nil); len(got) != 0 {
		t.Errorf("NumerosDeSprint(nil) = %v, se esperaba vacío", got)
	}
}

// Traza: DEF-05 (spec HU-26, sección 3): el Resumen del reporte de un proyecto muestra cuántos
// sprints hay abiertos y cuántos cerrados.
func TestContarSprints(t *testing.T) {
	tests := []struct {
		name                       string
		sprints                    []Sprint
		wantAbiertos, wantCerrados int
	}{
		{name: "sin sprints", wantAbiertos: 0, wantCerrados: 0},
		{name: "solo abiertos", sprints: []Sprint{{Estado: SprintAbierto}, {Estado: SprintAbierto}}, wantAbiertos: 2},
		{name: "abiertos y cerrados", sprints: []Sprint{{Estado: SprintCerrado}, {Estado: SprintAbierto}, {Estado: SprintCerrado}, {Estado: SprintCerrado}},
			wantAbiertos: 1, wantCerrados: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abiertos, cerrados := ContarSprints(tt.sprints)
			if abiertos != tt.wantAbiertos || cerrados != tt.wantCerrados {
				t.Errorf("ContarSprints = %d abiertos y %d cerrados, se esperaba %d y %d", abiertos, cerrados, tt.wantAbiertos, tt.wantCerrados)
			}
		})
	}
}
