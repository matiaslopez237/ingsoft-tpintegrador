package domain

import (
	"errors"
	"slices"
	"testing"
)

func elementoAEstimar(t *testing.T, estado EstadoElemento) ElementoBacklog {
	t.Helper()
	e, err := NuevoElementoBacklog(1, "Login", "Descripción", PrioridadAlta, nil, []string{"Criterio"})
	if err != nil {
		t.Fatal(err)
	}
	e.ID, e.Estado = 7, estado
	return e
}

func integrantes(nombres ...string) []Integrante {
	var r []Integrante
	for _, n := range nombres {
		r = append(r, Integrante{Nombre: n, Roles: []string{"Desarrollo"}})
	}
	return r
}

func rondaAbierta(t *testing.T, nombres ...string) Ronda {
	t.Helper()
	r, err := NuevaRonda(elementoAEstimar(t, EstadoPendiente), integrantes(nombres...), nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// votaron devuelve los nombres que ya votaron según el resumen.
func votaron(r Ronda) []string {
	var v []string
	for _, x := range r.Resumen().Votantes {
		if x.YaVoto {
			v = append(v, x.Nombre)
		}
	}
	return v
}

// Traza: HU-14 CA-14.1 (RN-14.1, RN-14.2)
func TestNuevaRonda_AbiertaConLosIntegrantesComoVotantes(t *testing.T) {
	for _, estado := range []EstadoElemento{EstadoPendiente, EstadoEnProgreso} {
		r, err := NuevaRonda(elementoAEstimar(t, estado), integrantes("Ana", "Beto"), nil)
		if err != nil {
			t.Fatalf("%s: error inesperado: %v", estado, err)
		}
		res := r.Resumen()
		if res.ProyectoID != 1 || res.ElementoID != 7 || res.Numero != 1 || res.Estado != RondaAbierta {
			t.Errorf("%s: resumen = %+v", estado, res)
		}
		want := []VotanteResumen{{Nombre: "Ana"}, {Nombre: "Beto"}}
		if !slices.Equal(res.Votantes, want) {
			t.Errorf("%s: votantes = %+v, se esperaba %+v", estado, res.Votantes, want)
		}
	}
}

// Traza: HU-14 CA-14.6 (RN-14.1, RN-14.2). Los errores se informan en el orden de la spec.
func TestNuevaRonda_Errores(t *testing.T) {
	abierta := rondaAbierta(t, "Ana")
	tests := []struct {
		nombre   string
		estado   EstadoElemento
		votantes []Integrante
		anterior *Ronda
		wantErr  error
	}{
		{"elemento terminado", EstadoHecho, integrantes("Ana"), nil, ErrElementoTerminado},
		{"sin integrantes", EstadoPendiente, nil, nil, ErrRondaSinVotantes},
		{"ya tiene una ronda abierta", EstadoPendiente, integrantes("Ana"), &abierta, ErrRondaYaAbierta},
		{"el elemento terminado se informa antes que la falta de integrantes", EstadoHecho, nil, nil, ErrElementoTerminado},
		{"la falta de integrantes se informa antes que la ronda abierta", EstadoPendiente, nil, &abierta, ErrRondaSinVotantes},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevaRonda(elementoAEstimar(t, tt.estado), tt.votantes, tt.anterior)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-14 CA-14.6 (RN-14.1): con la ronda anterior cerrada se abre la siguiente.
func TestNuevaRonda_ConLaAnteriorCerrada_TomaElNumeroSiguiente(t *testing.T) {
	anterior := rondaAbierta(t, "Ana")
	if err := anterior.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}

	r, err := NuevaRonda(elementoAEstimar(t, EstadoPendiente), integrantes("Ana"), &anterior)

	if err != nil || r.Numero != 2 || r.Estado != RondaAbierta {
		t.Errorf("ronda = %+v, err = %v; se esperaba la número 2 abierta", r, err)
	}
}

// Traza: HU-14 CA-14.2 (RN-14.3, RN-14.6)
func TestVotar_UnVotoNoCierraLaRondaNiMuestraElValor(t *testing.T) {
	r := rondaAbierta(t, "Ana", "Beto")

	if err := r.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}

	res := r.Resumen()
	want := []VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto"}}
	if res.Estado != RondaAbierta || !slices.Equal(res.Votantes, want) {
		t.Errorf("resumen = %+v, se esperaba abierta con %+v", res, want)
	}
}

// Traza: HU-14 CA-14.3 (RN-14.6)
func TestVotar_CuandoVotanTodosLaRondaSeCierra(t *testing.T) {
	r := rondaAbierta(t, "Ana", "Beto")
	if err := r.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}

	if err := r.Votar("Beto", 8); err != nil {
		t.Fatal(err)
	}

	if r.Estado != RondaCerrada || !slices.Equal(votaron(r), []string{"Ana", "Beto"}) {
		t.Errorf("ronda = %+v, se esperaba cerrada con los dos votos", r)
	}
	if err := r.Votar("Ana", 3); !errors.Is(err, ErrRondaCerrada) {
		t.Errorf("votar con la ronda cerrada: error = %v, se esperaba %v", err, ErrRondaCerrada)
	}
}

// Traza: HU-14 CA-14.4 (RN-14.5, RN-14.2)
func TestVotar_CambiarElVotoYNombreSinMayusculas(t *testing.T) {
	r := rondaAbierta(t, "Ana", "Beto")
	if err := r.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}

	if err := r.Votar(" ana ", 8); err != nil {
		t.Fatal(err)
	}

	if r.Estado != RondaAbierta || !slices.Equal(votaron(r), []string{"Ana"}) {
		t.Errorf("ronda = %+v, se esperaba abierta con un solo voto de Ana", r)
	}
	if v := r.Votantes[0].Voto; v == nil || *v != 8 {
		t.Errorf("voto de Ana = %v, se esperaba el último (8)", v)
	}
}

// Traza: HU-14 CA-14.5 (RN-14.2, RN-14.4, RN-14.6)
func TestVotar_Errores(t *testing.T) {
	cerrada := rondaAbierta(t, "Ana")
	if err := cerrada.Votar("Ana", 5); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		nombre     string
		ronda      Ronda
		integrante string
		valor      int
		wantErr    error
	}{
		{"valor cero", rondaAbierta(t, "Ana", "Beto"), "Ana", 0, ErrVotoInvalido},
		{"valor que no es de Fibonacci", rondaAbierta(t, "Ana", "Beto"), "Ana", 4, ErrVotoInvalido},
		{"valor mayor que 21", rondaAbierta(t, "Ana", "Beto"), "Ana", 34, ErrVotoInvalido},
		{"valor negativo", rondaAbierta(t, "Ana", "Beto"), "Ana", -1, ErrVotoInvalido},
		{"no es votante", rondaAbierta(t, "Ana", "Beto"), "Dani", 5, ErrNoEsVotante},
		{"ronda cerrada", cerrada, "Ana", 3, ErrRondaCerrada},
		{"la ronda cerrada se informa antes que el votante", cerrada, "Dani", 3, ErrRondaCerrada},
		{"el votante se informa antes que el valor", rondaAbierta(t, "Ana"), "Dani", 4, ErrNoEsVotante},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			antes := votaron(tt.ronda)
			err := tt.ronda.Votar(tt.integrante, tt.valor)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if !slices.Equal(votaron(tt.ronda), antes) {
				t.Errorf("con error no tendría que cambiar quién votó: %v → %v", antes, votaron(tt.ronda))
			}
		})
	}
}

// Traza: HU-14 RN-14.4: los valores de la escala, en orden, para el formulario de la pantalla.
func TestValoresVoto(t *testing.T) {
	want := []int{1, 2, 3, 5, 8, 13, 21}
	got := ValoresVoto()
	if !slices.Equal(got, want) {
		t.Fatalf("ValoresVoto() = %v, se esperaba %v", got, want)
	}
	got[0] = 99
	if ValoresVoto()[0] != 1 {
		t.Error("modificar la lista devuelta cambió la escala")
	}
}

// Traza: HU-14 RN-14.3: el resumen no comparte memoria con la ronda.
func TestResumen_NoCompartenLosVotantes(t *testing.T) {
	r := rondaAbierta(t, "Ana")
	res := r.Resumen()
	res.Votantes[0].Nombre = "Otra"
	if r.Votantes[0].Nombre != "Ana" {
		t.Error("modificar el resumen cambió la ronda")
	}
}
