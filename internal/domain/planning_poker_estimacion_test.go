package domain

import (
	"errors"
	"testing"
)

// Traza: HU-17 CA-17.1, CA-17.3 (RN-17.1, RN-17.2): cualquier valor de la escala se registra, aunque
// nadie lo haya votado, y la ronda lo conserva.
func TestRegistrarEstimacion_ValoresDeLaEscala(t *testing.T) {
	for _, valor := range []int{1, 2, 3, 5, 8, 13, 21} {
		r := rondaNumero(cerrada(t, 3, 13), 1)
		ultima := r

		if err := r.RegistrarEstimacion(valor, &ultima, elementoAEstimar(t, EstadoPendiente)); err != nil {
			t.Fatalf("valor %d: error inesperado: %v", valor, err)
		}
		if r.EstimacionAcordada == nil || *r.EstimacionAcordada != valor {
			t.Errorf("valor %d: EstimacionAcordada = %v", valor, r.EstimacionAcordada)
		}
	}
}

// Traza: HU-17 CA-17.3 a CA-17.5 (RN-17.2 a RN-17.4), sección 7. Si hay un error, la ronda no cambia.
// Los errores se informan en el orden de la tabla de la spec.
func TestRegistrarEstimacion_Errores(t *testing.T) {
	yaRegistrada := func(t *testing.T) Ronda {
		r := rondaNumero(cerrada(t, 3, 13), 1)
		v := 5
		r.EstimacionAcordada = &v
		return r
	}
	sinEstimacion := func(t *testing.T) Ronda { return rondaNumero(cerrada(t, 3, 13), 1) }
	abierta := func(t *testing.T) Ronda { return rondaNumero(rondaAbierta(t, "Ana", "Beto"), 1) }
	otraUltima := func(t *testing.T) *Ronda {
		r := rondaNumero(cerrada(t, 5, 8), 2)
		return &r
	}

	tests := []struct {
		name   string
		ronda  func(*testing.T) Ronda
		ultima func(*testing.T, Ronda) *Ronda
		estado EstadoElemento
		valor  int
		want   error
	}{
		{"ronda abierta", abierta, esLaUltima, EstadoPendiente, 8, ErrEstimacionRondaAbierta},
		{"no es la última", sinEstimacion, func(t *testing.T, _ Ronda) *Ronda { return otraUltima(t) }, EstadoPendiente, 8, ErrRondaNoEsLaUltima},
		{"última desconocida", sinEstimacion, func(*testing.T, Ronda) *Ronda { return nil }, EstadoPendiente, 8, ErrRondaNoEsLaUltima},
		{"elemento Hecho", sinEstimacion, esLaUltima, EstadoHecho, 8, ErrElementoTerminado},
		{"ya registrada", yaRegistrada, esLaUltima, EstadoPendiente, 8, ErrEstimacionYaRegistrada},
		{"fuera de la escala: 4", sinEstimacion, esLaUltima, EstadoPendiente, 4, ErrEstimacionInvalida},
		{"fuera de la escala: 0", sinEstimacion, esLaUltima, EstadoPendiente, 0, ErrEstimacionInvalida},
		{"fuera de la escala: 34", sinEstimacion, esLaUltima, EstadoPendiente, 34, ErrEstimacionInvalida},
		{"abierta y valor inválido: primero la ronda", abierta, esLaUltima, EstadoPendiente, 4, ErrEstimacionRondaAbierta},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.ronda(t)
			antes := r.EstimacionAcordada

			err := r.RegistrarEstimacion(tt.valor, tt.ultima(t, r), elementoAEstimar(t, tt.estado))

			if !errors.Is(err, tt.want) {
				t.Errorf("error = %v, se esperaba %v", err, tt.want)
			}
			if r.EstimacionAcordada != antes {
				t.Errorf("EstimacionAcordada cambió a %v", r.EstimacionAcordada)
			}
		})
	}
}

func esLaUltima(_ *testing.T, r Ronda) *Ronda { return &r }

// Traza: HU-17 CA-17.10 (RN-17.3, RN-17.4): la pantalla ofrece el formulario solo si se puede
// registrar, con la misma regla.
func TestPuedeRegistrarEstimacion(t *testing.T) {
	cerrada1 := rondaNumero(cerrada(t, 3, 13), 1)
	abierta1 := rondaNumero(rondaAbierta(t, "Ana", "Beto"), 1)
	cerrada2 := rondaNumero(cerrada(t, 5, 8), 2)
	conEstimacion := rondaNumero(cerrada(t, 3, 13), 1)
	v := 5
	conEstimacion.EstimacionAcordada = &v
	resumen := func(r Ronda) *ResumenRonda { res := r.Resumen(); return &res }
	pendiente, hecho := elementoAEstimar(t, EstadoPendiente), elementoAEstimar(t, EstadoHecho)

	tests := []struct {
		name     string
		ronda    Ronda
		ultima   *ResumenRonda
		elemento ElementoBacklog
		want     bool
	}{
		{"última, cerrada y sin estimación", cerrada1, resumen(cerrada1), pendiente, true},
		{"abierta", abierta1, resumen(abierta1), pendiente, false},
		{"no es la última", cerrada1, resumen(cerrada2), pendiente, false},
		{"última desconocida", cerrada1, nil, pendiente, false},
		{"elemento Hecho", cerrada1, resumen(cerrada1), hecho, false},
		{"ya tiene estimación", conEstimacion, resumen(conEstimacion), pendiente, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PuedeRegistrarEstimacion(tt.ronda.Resumen(), tt.ultima, tt.elemento); got != tt.want {
				t.Errorf("PuedeRegistrarEstimacion() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

// Traza: HU-17 sección 3: el resumen trae la estimación, sin compartir memoria con la ronda.
func TestResumen_TraeLaEstimacionAcordada(t *testing.T) {
	r := rondaNumero(cerrada(t, 3, 13), 1)
	if r.Resumen().EstimacionAcordada != nil {
		t.Error("sin estimación, el resumen no debería traerla")
	}
	ultima := r
	if err := r.RegistrarEstimacion(8, &ultima, elementoAEstimar(t, EstadoPendiente)); err != nil {
		t.Fatal(err)
	}

	res := r.Resumen()
	if res.EstimacionAcordada == nil || *res.EstimacionAcordada != 8 {
		t.Fatalf("EstimacionAcordada del resumen = %v, se esperaba 8", res.EstimacionAcordada)
	}
	*res.EstimacionAcordada = 13
	if *r.EstimacionAcordada != 8 {
		t.Error("cambiar el resumen no debería cambiar la ronda")
	}
}
