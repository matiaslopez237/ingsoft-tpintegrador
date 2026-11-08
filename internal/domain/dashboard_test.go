package domain

import (
	"reflect"
	"testing"
)

// elementoDashboard arma un elemento del backlog con su estado, sus story points (nil: sin estimar)
// y su sprint (nil: sin asignar).
func elementoDashboard(id int64, estado EstadoElemento, puntos *int, sprintID *int64) ElementoBacklog {
	return ElementoBacklog{ID: id, ProyectoID: 7, Titulo: "Elemento", Estado: estado, StoryPoints: puntos, SprintID: sprintID}
}

func puntosDashboard(n int) *int { return &n }

func sprintDashboard(id int64) *int64 { return &id }

// Traza: HU-25 CA-25.1, CA-25.11, RN-25.2
func TestContarBacklog(t *testing.T) {
	tests := []struct {
		name      string
		elementos []ElementoBacklog
		want      ConteoBacklog
	}{
		{
			name: "CA-25.11: sin elementos, todo en 0",
			want: ConteoBacklog{},
		},
		{
			name: "CA-25.1: cuenta cada estado y el total",
			elementos: []ElementoBacklog{
				elementoDashboard(1, EstadoPendiente, nil, nil),
				elementoDashboard(2, EstadoPendiente, nil, nil),
				elementoDashboard(3, EstadoEnProgreso, nil, nil),
				elementoDashboard(4, EstadoHecho, nil, nil),
				elementoDashboard(5, EstadoHecho, nil, nil),
				elementoDashboard(6, EstadoHecho, nil, nil),
			},
			want: ConteoBacklog{Pendientes: 2, EnProgreso: 1, Hechos: 3, Total: 6},
		},
		{
			name: "RN-25.2: cuentan estén o no en un sprint",
			elementos: []ElementoBacklog{
				elementoDashboard(1, EstadoPendiente, nil, sprintDashboard(4)),
				elementoDashboard(2, EstadoHecho, puntosDashboard(3), nil),
			},
			want: ConteoBacklog{Pendientes: 1, Hechos: 1, Total: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContarBacklog(tt.elementos); got != tt.want {
				t.Errorf("ContarBacklog = %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

// Traza: HU-25 CA-25.2, CA-25.3, RN-25.3, RN-25.4
func TestResumirSprintEnCurso(t *testing.T) {
	hoy := fecha(t, "2026-03-10")
	enCurso := sprintEn(t, 2, "2026-03-01", "2026-03-14", SprintAbierto)
	tests := []struct {
		name      string
		sprints   []Sprint
		elementos []ElementoBacklog
		want      SprintEnCurso
		wantHay   bool
	}{
		{
			name:    "CA-25.3: sin sprints no hay sprint en curso",
			wantHay: false,
		},
		{
			name: "CA-25.2: suma todos los puntos del sprint y los Hecho como completados",
			sprints: []Sprint{
				sprintEn(t, 1, "2026-02-01", "2026-02-14", SprintCerrado),
				enCurso,
			},
			elementos: []ElementoBacklog{
				elementoDashboard(1, EstadoPendiente, puntosDashboard(3), sprintDashboard(2)),
				elementoDashboard(2, EstadoHecho, puntosDashboard(5), sprintDashboard(2)),
				elementoDashboard(3, EstadoEnProgreso, puntosDashboard(2), sprintDashboard(2)),
			},
			want:    SprintEnCurso{Sprint: enCurso, PuntosPlanificados: 10, PuntosCompletados: 5},
			wantHay: true,
		},
		{
			name:    "RN-25.4: un elemento sin estimar suma 0 (RN-11.6)",
			sprints: []Sprint{enCurso},
			elementos: []ElementoBacklog{
				elementoDashboard(1, EstadoHecho, puntosDashboard(5), sprintDashboard(2)),
				elementoDashboard(2, EstadoHecho, nil, sprintDashboard(2)),
			},
			want:    SprintEnCurso{Sprint: enCurso, PuntosPlanificados: 5, PuntosCompletados: 5},
			wantHay: true,
		},
		{
			name:    "RN-25.4: los elementos de otro sprint o sin sprint no cuentan",
			sprints: []Sprint{enCurso},
			elementos: []ElementoBacklog{
				elementoDashboard(1, EstadoHecho, puntosDashboard(5), sprintDashboard(2)),
				elementoDashboard(2, EstadoHecho, puntosDashboard(8), sprintDashboard(1)),
				elementoDashboard(3, EstadoPendiente, puntosDashboard(3), nil),
			},
			want:    SprintEnCurso{Sprint: enCurso, PuntosPlanificados: 5, PuntosCompletados: 5},
			wantHay: true,
		},
		{
			name:    "límite: un sprint en curso sin elementos tiene 0 y 0",
			sprints: []Sprint{enCurso},
			want:    SprintEnCurso{Sprint: enCurso},
			wantHay: true,
		},
		{
			name:    "CA-25.3 (P-1): un sprint abierto que ya venció no es el sprint en curso",
			sprints: []Sprint{sprintEn(t, 1, "2026-02-01", "2026-02-14", SprintAbierto)},
			wantHay: false,
		},
		{
			name:    "RN-25.3: un sprint abierto que todavía no empezó no es el sprint en curso",
			sprints: []Sprint{sprintEn(t, 3, "2026-03-15", "2026-03-28", SprintAbierto)},
			wantHay: false,
		},
		{
			name:    "RN-25.3: un sprint cerrado cuyas fechas incluyen hoy no es el sprint en curso",
			sprints: []Sprint{sprintEn(t, 2, "2026-03-01", "2026-03-14", SprintCerrado)},
			wantHay: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, hay := ResumirSprintEnCurso(tt.sprints, tt.elementos, hoy)
			if hay != tt.wantHay {
				t.Fatalf("hay sprint en curso = %v, se esperaba %v", hay, tt.wantHay)
			}
			if hay && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResumirSprintEnCurso = %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

// Traza: HU-25 CA-25.4, RN-25.5
func TestSprintsCerradosCronologicos(t *testing.T) {
	tests := []struct {
		name    string
		sprints []Sprint
		want    []int
	}{
		{
			name: "CA-25.11: sin sprints",
			want: []int{},
		},
		{
			name: "CA-25.4: solo los cerrados, del más antiguo al más reciente",
			sprints: []Sprint{
				sprintEn(t, 2, "2026-01-15", "2026-01-28", SprintCerrado),
				sprintEn(t, 3, "2026-03-01", "2026-03-14", SprintAbierto),
				sprintEn(t, 1, "2026-01-01", "2026-01-14", SprintCerrado),
			},
			want: []int{1, 2},
		},
		{
			name: "RN-25.5 (P-2): todos los cerrados, sin límite de cantidad",
			sprints: []Sprint{
				sprintEn(t, 1, "2025-01-01", "2025-01-14", SprintCerrado),
				sprintEn(t, 2, "2025-02-01", "2025-02-14", SprintCerrado),
				sprintEn(t, 3, "2025-03-01", "2025-03-14", SprintCerrado),
				sprintEn(t, 4, "2025-04-01", "2025-04-14", SprintCerrado),
				sprintEn(t, 5, "2025-05-01", "2025-05-14", SprintCerrado),
				sprintEn(t, 6, "2025-06-01", "2025-06-14", SprintCerrado),
				sprintEn(t, 7, "2025-07-01", "2025-07-14", SprintCerrado),
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numeros(SprintsCerradosCronologicos(tt.sprints)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SprintsCerradosCronologicos = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

// Traza: HU-25 RN-25.5: no modifica la lista que recibe.
func TestSprintsCerradosCronologicos_NoModificaLaLista(t *testing.T) {
	sprints := []Sprint{
		sprintEn(t, 2, "2026-01-15", "2026-01-28", SprintCerrado),
		sprintEn(t, 1, "2026-01-01", "2026-01-14", SprintCerrado),
	}

	SprintsCerradosCronologicos(sprints)

	if got := numeros(sprints); !reflect.DeepEqual(got, []int{2, 1}) {
		t.Errorf("la lista recibida quedó %v, se esperaba sin cambios [2 1]", got)
	}
}
