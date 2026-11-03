package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// historialFalso envuelve el fake de HU-14 (pokerFalso) y le suma el puerto de HU-16, sin tocarlo.
// listo cuenta las veces que se leyó el historial.
type historialFalso struct {
	*pokerFalso
	listo int
}

func (f *historialFalso) ListarRondas(_ context.Context, elementoID int64) ([]domain.Ronda, error) {
	f.listo++
	if f.errRonda != nil {
		return nil, f.errRonda
	}
	var del []domain.Ronda
	for _, r := range f.rondas {
		if r.ElementoID == elementoID {
			del = append(del, f.copia(r))
		}
	}
	return del, nil
}

// rondaCerradaEn abre una ronda sobre el elemento y la cierra con un voto de cada votante.
func rondaCerradaEn(t *testing.T, f *pokerFalso, elementoID int64, votos map[string]int) {
	t.Helper()
	iniciar, votar, _ := casosPoker(f)
	r, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: elementoID})
	if err != nil {
		t.Fatalf("iniciar ronda: %v", err)
	}
	for nombre, valor := range votos {
		if _, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: r.ID, Integrante: nombre, Valor: valor}); err != nil {
			t.Fatalf("votar: %v", err)
		}
	}
}

// Traza: HU-16 CA-16.1, CA-16.4 (RN-16.2, RN-16.3): las rondas del elemento, de la más reciente a la
// más antigua, sin las de otro elemento; la anterior conserva sus votos.
func TestListarRondas_DevuelveElHistorialDelElemento(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana", "Beto")
	f.elementos[8] = elementoDePrueba(t, 8, 1, "Registro")
	rondaCerradaEn(t, f, 7, map[string]int{"Ana": 3, "Beto": 13})
	rondaCerradaEn(t, f, 8, map[string]int{"Ana": 5, "Beto": 5}) // de otro elemento
	iniciar, _, _ := casosPoker(f)
	if _, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}); err != nil {
		t.Fatalf("iniciar la ronda 2: %v", err)
	}

	got, err := NuevoListarRondas(f, f, &historialFalso{pokerFalso: f}).Ejecutar(context.Background(), 1, 7)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 2 || got[0].Numero != 2 || got[1].Numero != 1 {
		t.Fatalf("historial = %+v, se esperaban las rondas 2 y 1 del elemento 7", got)
	}
	if got[0].Revelacion != nil {
		t.Errorf("la ronda 2, abierta, no debería revelar votos")
	}
	if r := got[1].Revelacion; r == nil || r.Minimo != 3 || r.Maximo != 13 {
		t.Errorf("revelación de la ronda 1 = %+v, se esperaba de 3 a 13", r)
	}
}

// Traza: HU-16 CA-16.5: sin rondas es una lista vacía, no un error.
func TestListarRondas_SinRondas_ListaVacia(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana")

	got, err := NuevoListarRondas(f, f, &historialFalso{pokerFalso: f}).Ejecutar(context.Background(), 1, 7)

	if err != nil || len(got) != 0 {
		t.Errorf("historial = %+v, error = %v; se esperaba una lista vacía sin error", got, err)
	}
}

// Traza: HU-16 CA-16.6 (RN-16.5): proyecto o elemento inexistentes, o elemento de otro proyecto. En
// ningún caso se lee el historial.
func TestListarRondas_Errores(t *testing.T) {
	tests := []struct {
		name                   string
		proyectoID, elementoID int64
		want                   error
	}{
		{"proyecto inexistente", 99, 7, domain.ErrProyectoNoEncontrado},
		{"elemento inexistente", 1, 99, domain.ErrElementoNoEncontrado},
		{"elemento de otro proyecto", 2, 7, domain.ErrElementoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoPokerFalso(t, "Ana")
			h := &historialFalso{pokerFalso: f}

			_, err := NuevoListarRondas(f, f, h).Ejecutar(context.Background(), tt.proyectoID, tt.elementoID)

			if !errors.Is(err, tt.want) {
				t.Errorf("error = %v, se esperaba %v", err, tt.want)
			}
			if h.listo != 0 {
				t.Error("no se tendría que leer el historial")
			}
		})
	}
}

// Traza: HU-16 sección 7: una falla al leer las rondas se propaga sin perder la causa.
func TestListarRondas_FallaDelRepositorio_SePropaga(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana")
	disco := errors.New("disco lleno")
	f.errRonda = disco

	_, err := NuevoListarRondas(f, f, &historialFalso{pokerFalso: f}).Ejecutar(context.Background(), 1, 7)

	if !errors.Is(err, disco) {
		t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
	}
}
