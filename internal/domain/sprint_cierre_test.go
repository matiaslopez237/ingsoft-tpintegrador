package domain

import (
	"reflect"
	"testing"
)

// historiaDelSprint arma una historia asignada al sprint 1, con su estado y story points (nil: sin
// estimar).
func historiaDelSprint(id int64, estado EstadoElemento, sp *int) ElementoBacklog {
	sprint := int64(1)
	return ElementoBacklog{ID: id, ProyectoID: 7, Estado: estado, StoryPoints: sp, SprintID: &sprint}
}

// Traza: HU-11 CA-11.6, CA-11.7, RN-11.3, RN-11.6 (P-11 a P-13)
func TestResumirCierre(t *testing.T) {
	tests := []struct {
		name      string
		historias []ElementoBacklog
		want      ResumenCierre
	}{
		{
			name: "sin historias: 0 y 0, ninguna vuelve",
			want: ResumenCierre{},
		},
		{
			name: "CA-11.7: 5 Hecho, 3 Pendiente y sin estimar Hecho",
			historias: []ElementoBacklog{
				historiaDelSprint(1, EstadoHecho, puntos(5)),
				historiaDelSprint(2, EstadoPendiente, puntos(3)),
				historiaDelSprint(3, EstadoHecho, nil),
			},
			want: ResumenCierre{PuntosPlanificados: 8, PuntosCompletados: 5, DevueltasAlBacklog: []int64{2}},
		},
		{
			name: "Pendiente y En progreso vuelven al backlog; las Hecho quedan",
			historias: []ElementoBacklog{
				historiaDelSprint(1, EstadoPendiente, puntos(2)),
				historiaDelSprint(2, EstadoHecho, puntos(8)),
				historiaDelSprint(3, EstadoEnProgreso, puntos(5)),
			},
			want: ResumenCierre{PuntosPlanificados: 15, PuntosCompletados: 8, DevueltasAlBacklog: []int64{1, 3}},
		},
		{
			name: "todas Hecho: planificados = completados",
			historias: []ElementoBacklog{
				historiaDelSprint(1, EstadoHecho, puntos(3)),
				historiaDelSprint(2, EstadoHecho, puntos(13)),
			},
			want: ResumenCierre{PuntosPlanificados: 16, PuntosCompletados: 16},
		},
		{
			name: "sin estimar y sin terminar: suma 0 pero vuelve igual",
			historias: []ElementoBacklog{
				historiaDelSprint(1, EstadoEnProgreso, nil),
			},
			want: ResumenCierre{DevueltasAlBacklog: []int64{1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResumirCierre(tt.historias)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResumirCierre() = %+v, se esperaba %+v", got, tt.want)
			}
		})
	}
}

// Traza: HU-11 RN-11.3: resumir no modifica las historias recibidas (de eso se ocupa la persistencia).
func TestResumirCierre_NoModificaLasHistorias(t *testing.T) {
	historias := []ElementoBacklog{historiaDelSprint(1, EstadoPendiente, puntos(3))}
	antes := historias[0]

	ResumirCierre(historias)

	if !reflect.DeepEqual(historias[0], antes) {
		t.Errorf("historia = %+v, se esperaba sin cambios %+v", historias[0], antes)
	}
}
