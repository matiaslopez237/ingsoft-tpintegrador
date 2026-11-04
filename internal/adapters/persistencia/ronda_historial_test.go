package persistencia

import (
	"context"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador también cumple el puerto angosto de HU-16 (se verifica al compilar).
var _ app.LectorHistorialRondas = (*RepositorioRondas)(nil)

// Traza: HU-16 CA-16.1, CA-16.4 (RN-16.2): todas las rondas del elemento, con sus votantes y votos
// tal como se guardaron, sin las de otro elemento. El orden lo decide el dominio (HistorialRondas).
func TestRepositorioRondas_ListarRondas(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r1 := b.abrir(t, nil, "Ana", "Beto")
	if err := b.votar(t, &r1, "Ana", 3); err != nil {
		t.Fatal(err)
	}
	if err := b.votar(t, &r1, "Beto", 13); err != nil {
		t.Fatal(err)
	}
	r1, _ = b.repo.BuscarRonda(ctx, r1.ID) // como quedó guardada: Cerrada
	r2 := b.abrir(t, &r1, "Ana", "Beto")
	otro := guardarElemento(t, b.backlog, b.elemento.ProyectoID, "Registro")
	ajena, err := domain.NuevaRonda(otro, []domain.Integrante{{Nombre: "Ana"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.repo.GuardarRonda(ctx, &ajena); err != nil {
		t.Fatal(err)
	}

	got, err := b.repo.ListarRondas(ctx, b.elemento.ID)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("rondas = %+v, se esperaban 2 (sin la del otro elemento)", got)
	}
	porNumero := map[int]domain.Ronda{got[0].Numero: got[0], got[1].Numero: got[1]}
	if !reflect.DeepEqual(porNumero[1], r1) {
		t.Errorf("ronda 1 = %+v, se esperaba %+v", porNumero[1], r1)
	}
	if !reflect.DeepEqual(porNumero[2], r2) {
		t.Errorf("ronda 2 = %+v, se esperaba %+v", porNumero[2], r2)
	}
}

// Traza: HU-16 CA-16.5: un elemento sin rondas da una lista vacía, sin error.
func TestRepositorioRondas_ListarRondas_SinRondas(t *testing.T) {
	b := nuevaBasePoker(t)

	got, err := b.repo.ListarRondas(context.Background(), b.elemento.ID)

	if err != nil || len(got) != 0 {
		t.Errorf("rondas = %+v, error = %v; se esperaba una lista vacía sin error", got, err)
	}
}
