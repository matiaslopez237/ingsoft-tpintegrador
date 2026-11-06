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

// Traza: HU-22 CA-22.2 a CA-22.5, CA-22.9
func TestVelocidad(t *testing.T) {
	tests := []struct {
		name          string
		cerrados      []Sprint
		wantVelocidad float64
		wantCantidad  int
		wantHayDatos  bool
	}{
		{
			name:         "CA-22.3: sin sprints cerrados, no hay velocidad",
			cerrados:     nil,
			wantCantidad: 0,
			wantHayDatos: false,
		},
		{
			name:          "CA-22.2: promedio de varios sprints cerrados",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 8), sprintCerradoCon(t, 2, 6), sprintCerradoCon(t, 3, 4)},
			wantVelocidad: 6,
			wantCantidad:  3,
			wantHayDatos:  true,
		},
		{
			name:          "CA-22.4: el promedio se informa con decimales",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 5), sprintCerradoCon(t, 2, 4)},
			wantVelocidad: 4.5,
			wantCantidad:  2,
			wantHayDatos:  true,
		},
		{
			name:          "CA-22.5: un único sprint cerrado con 0 puntos da velocidad 0, no 'sin datos'",
			cerrados:      []Sprint{sprintCerradoCon(t, 1, 0)},
			wantVelocidad: 0,
			wantCantidad:  1,
			wantHayDatos:  true,
		},
		{
			// cerrados ya viene del más reciente al más antiguo (como lo devuelve SprintsAnteriores):
			// el 4.º (2 puntos) es el más antiguo y queda afuera de la ventana de 3.
			name: "CA-22.9: con más de 3 sprints cerrados, solo se promedian los 3 más recientes",
			cerrados: []Sprint{
				sprintCerradoCon(t, 4, 8), sprintCerradoCon(t, 3, 6), sprintCerradoCon(t, 2, 4), sprintCerradoCon(t, 1, 2),
			},
			wantVelocidad: 6, // (8+6+4)/3, el sprint 1 (2 puntos) no entra
			wantCantidad:  3,
			wantHayDatos:  true,
		},
		{
			name:          "CA-22.9: con exactamente 3 sprints cerrados, la ventana no recorta nada",
			cerrados:      []Sprint{sprintCerradoCon(t, 3, 9), sprintCerradoCon(t, 2, 6), sprintCerradoCon(t, 1, 3)},
			wantVelocidad: 6,
			wantCantidad:  3,
			wantHayDatos:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			velocidad, cantidad, hayDatos := Velocidad(tt.cerrados)

			if hayDatos != tt.wantHayDatos {
				t.Fatalf("hayDatos = %v, se esperaba %v", hayDatos, tt.wantHayDatos)
			}
			if cantidad != tt.wantCantidad {
				t.Errorf("cantidad = %d, se esperaba %d", cantidad, tt.wantCantidad)
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

	velocidad, cantidad, hayDatos := Velocidad([]Sprint{abierto})

	if !hayDatos || velocidad != 99 || cantidad != 1 {
		t.Errorf("Velocidad() = %v, %d, %v; se esperaba 99, 1, true (no filtra por estado)", velocidad, cantidad, hayDatos)
	}
}

// Traza: HU-24 CA-24.1, CA-24.2
func TestPorcentajeCompletado(t *testing.T) {
	tests := []struct {
		name              string
		planificados      int
		completados       int
		wantPorcentaje    float64
		wantHayPorcentaje bool
	}{
		{name: "CA-24.1: porcentaje exacto", planificados: 8, completados: 6, wantPorcentaje: 75, wantHayPorcentaje: true},
		// No pasa con los datos reales de ResumirCierre (HU-11, RN-11.6): PuntosCompletados siempre es
		// un subconjunto de PuntosPlanificados. Documenta que la función, al ser pura, no lo valida.
		{name: "no ocurre con datos reales, pero la función no lo impide", planificados: 4, completados: 6, wantPorcentaje: 150, wantHayPorcentaje: true},
		{name: "CA-24.1: 0 completados sobre planificados da 0%, no 'sin datos'", planificados: 5, completados: 0, wantPorcentaje: 0, wantHayPorcentaje: true},
		{name: "CA-24.1: completados igual a planificados da 100%", planificados: 5, completados: 5, wantPorcentaje: 100, wantHayPorcentaje: true},
		{name: "CA-24.1: se redondea a un decimal, con los medios alejándose de cero", planificados: 16, completados: 1, wantPorcentaje: 6.3, wantHayPorcentaje: true},
		{name: "CA-24.2: sin puntos planificados, el porcentaje queda indefinido", planificados: 0, completados: 0, wantHayPorcentaje: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sprintCerradoCon(t, 1, tt.completados)
			s.PuntosPlanificados = tt.planificados

			porcentaje, hayPorcentaje := PorcentajeCompletado(s)

			if hayPorcentaje != tt.wantHayPorcentaje {
				t.Fatalf("hayPorcentaje = %v, se esperaba %v", hayPorcentaje, tt.wantHayPorcentaje)
			}
			if hayPorcentaje && porcentaje != tt.wantPorcentaje {
				t.Errorf("PorcentajeCompletado() = %v, se esperaba %v", porcentaje, tt.wantPorcentaje)
			}
		})
	}
}

// idSprintHistorial arma un puntero al ID de sprint indicado, para los defectos de prueba de
// ContarDefectosDeSprint (SprintDeteccionID y SprintResolucionID son *int64).
func idSprintHistorial(id int64) *int64 { return &id }

// Traza: HU-24 CA-24.3 a CA-24.5
func TestContarDefectosDeSprint(t *testing.T) {
	tests := []struct {
		name           string
		defectos       []Defecto
		sprintID       int64
		wantDetectados int
		wantResueltos  int
	}{
		{
			name:           "sin defectos",
			sprintID:       1,
			wantDetectados: 0,
			wantResueltos:  0,
		},
		{
			name: "CA-24.3: cuenta los detectados en ese sprint, no los de otro",
			defectos: []Defecto{
				{SprintDeteccionID: idSprintHistorial(1)},
				{SprintDeteccionID: idSprintHistorial(1)},
				{SprintDeteccionID: idSprintHistorial(2)},
			},
			sprintID:       1,
			wantDetectados: 2,
			wantResueltos:  0,
		},
		{
			name: "CA-24.4: cuenta los resueltos en ese sprint, aunque ya estén Cerrados",
			defectos: []Defecto{
				{Estado: DefectoCerrado, SprintDeteccionID: idSprintHistorial(1), SprintResolucionID: idSprintHistorial(1)},
			},
			sprintID:       1,
			wantDetectados: 1,
			wantResueltos:  1,
		},
		{
			name: "CA-24.5: un defecto sin sprint no cuenta en ninguno",
			defectos: []Defecto{
				{SprintDeteccionID: nil, SprintResolucionID: nil},
			},
			sprintID:       1,
			wantDetectados: 0,
			wantResueltos:  0,
		},
		{
			name: "un defecto todavía sin resolver no cuenta como resuelto",
			defectos: []Defecto{
				{Estado: DefectoEnProgreso, SprintDeteccionID: idSprintHistorial(1), SprintResolucionID: nil},
			},
			sprintID:       1,
			wantDetectados: 1,
			wantResueltos:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectados, resueltos := ContarDefectosDeSprint(tt.defectos, tt.sprintID)

			if detectados != tt.wantDetectados || resueltos != tt.wantResueltos {
				t.Errorf("ContarDefectosDeSprint() = %d detectados, %d resueltos; se esperaba %d, %d",
					detectados, resueltos, tt.wantDetectados, tt.wantResueltos)
			}
		})
	}
}
