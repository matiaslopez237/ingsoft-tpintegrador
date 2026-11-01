package domain

import (
	"reflect"
	"testing"
)

// rondaNumero le pone ID y número a una ronda armada con los helpers de HU-14.
func rondaNumero(r Ronda, numero int) Ronda {
	r.ID, r.Numero = int64(numero*10), numero
	return r
}

// Traza: HU-16 CA-16.4, RN-16.3: de la más reciente a la más antigua; las cerradas con la
// revelación de HU-15 y la abierta sin ningún valor.
func TestHistorialRondas_OrdenYRevelacion(t *testing.T) {
	abierta := rondaAbierta(t, "Ana", "Beto")
	if err := abierta.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}
	rondas := []Ronda{
		rondaNumero(cerrada(t, 3, 13), 1),
		rondaNumero(cerrada(t, 5, 8), 2),
		rondaNumero(abierta, 3),
	}

	got := HistorialRondas(rondas)

	var numeros []int
	for _, r := range got {
		numeros = append(numeros, r.Numero)
	}
	if !reflect.DeepEqual(numeros, []int{3, 2, 1}) {
		t.Fatalf("números = %v, se esperaba [3 2 1]", numeros)
	}
	if got[0].Revelacion != nil {
		t.Errorf("la ronda abierta no debería revelar votos: %+v", got[0].Revelacion)
	}
	if !got[0].Votantes[0].YaVoto || got[0].Votantes[1].YaVoto {
		t.Errorf("votantes de la ronda abierta = %+v, se esperaba que solo Ana figure como que votó", got[0].Votantes)
	}
	if r := got[1].Revelacion; r == nil || r.Minimo != 5 || r.Maximo != 8 {
		t.Errorf("revelación de la ronda 2 = %+v, se esperaba mínimo 5 y máximo 8", r)
	}
	if r := got[2].Revelacion; r == nil || r.Minimo != 3 || r.Maximo != 13 || !r.DiferenciaSignificativa {
		t.Errorf("revelación de la ronda 1 = %+v, se esperaba 3 a 13 con diferencia significativa", r)
	}
}

// Traza: HU-16 CA-16.5: sin rondas, una lista vacía.
func TestHistorialRondas_SinRondas(t *testing.T) {
	if got := HistorialRondas(nil); got == nil || len(got) != 0 {
		t.Errorf("HistorialRondas(nil) = %#v, se esperaba una lista vacía", got)
	}
}

// Traza: HU-16 RN-16.2: armar el historial no modifica las rondas recibidas.
func TestHistorialRondas_NoModificaLasRondas(t *testing.T) {
	rondas := []Ronda{rondaNumero(cerrada(t, 3, 13), 1), rondaNumero(cerrada(t, 5, 8), 2)}
	antes := []int{rondas[0].Numero, rondas[1].Numero}

	HistorialRondas(rondas)

	if rondas[0].Numero != antes[0] || rondas[1].Numero != antes[1] {
		t.Errorf("el orden de las rondas recibidas cambió: %d, %d", rondas[0].Numero, rondas[1].Numero)
	}
}

// Traza: HU-16 CA-16.8, RN-16.4: "Nueva ronda" solo en la última ronda, cerrada, de un elemento
// que no está Hecho; con o sin diferencia significativa.
func TestPuedeRepetirse(t *testing.T) {
	cerrada1 := rondaNumero(cerrada(t, 3, 13), 1).Resumen()
	sinDiferencia := rondaNumero(cerrada(t, 5, 5), 1).Resumen()
	abierta2 := rondaNumero(rondaAbierta(t, "Ana", "Beto"), 2).Resumen()
	cerrada2 := rondaNumero(cerrada(t, 5, 8), 2).Resumen()
	pendiente, hecho := elementoAEstimar(t, EstadoPendiente), elementoAEstimar(t, EstadoHecho)

	tests := []struct {
		name     string
		ronda    ResumenRonda
		ultima   *ResumenRonda
		elemento ElementoBacklog
		want     bool
	}{
		{"última y cerrada", cerrada1, &cerrada1, pendiente, true},
		{"última y cerrada, sin diferencia", sinDiferencia, &sinDiferencia, pendiente, true},
		{"abierta", abierta2, &abierta2, pendiente, false},
		{"cerrada pero no es la última (la 2 está abierta)", cerrada1, &abierta2, pendiente, false},
		{"cerrada pero no es la última (la 2 también se cerró)", cerrada1, &cerrada2, pendiente, false},
		{"elemento Hecho", cerrada1, &cerrada1, hecho, false},
		{"sin última conocida", cerrada1, nil, pendiente, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PuedeRepetirse(tt.ronda, tt.ultima, tt.elemento); got != tt.want {
				t.Errorf("PuedeRepetirse() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}
