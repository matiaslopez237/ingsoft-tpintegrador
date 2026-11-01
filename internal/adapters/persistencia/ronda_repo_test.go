package persistencia

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto de HU-14 (se verifica al compilar).
var _ app.RepositorioRondas = (*RepositorioRondas)(nil)

// basePoker es una base SQLite real con un proyecto y el elemento "Login".
type basePoker struct {
	repo     *RepositorioRondas
	backlog  *RepositorioBacklog
	elemento domain.ElementoBacklog
}

func nuevaBasePoker(t *testing.T) basePoker {
	t.Helper()
	db, err := AbrirDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	p := domain.Proyecto{Nombre: "Alfa", FechaInicio: dia(1), Estado: domain.EstadoActivo}
	if err := NuevoRepositorioProyectos(db).Guardar(context.Background(), &p); err != nil {
		t.Fatal(err)
	}
	backlog := NuevoRepositorioBacklog(db)
	return basePoker{repo: NuevoRepositorioRondas(db), backlog: backlog, elemento: guardarElemento(t, backlog, p.ID, "Login")}
}

// abrir guarda una ronda abierta sobre "Login" con los votantes dados.
func (b basePoker) abrir(t *testing.T, anterior *domain.Ronda, votantes ...string) domain.Ronda {
	t.Helper()
	var integrantes []domain.Integrante
	for _, n := range votantes {
		integrantes = append(integrantes, domain.Integrante{Nombre: n})
	}
	r, err := domain.NuevaRonda(b.elemento, integrantes, anterior)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.repo.GuardarRonda(context.Background(), &r); err != nil {
		t.Fatalf("guardar ronda: %v", err)
	}
	return r
}

// votar aplica el voto en el dominio y lo guarda.
func (b basePoker) votar(t *testing.T, r *domain.Ronda, integrante string, valor int) error {
	t.Helper()
	if err := r.Votar(integrante, valor); err != nil {
		t.Fatal(err)
	}
	return b.repo.GuardarVoto(context.Background(), *r, integrante)
}

// Traza: HU-14 CA-14.1 (RN-14.1, RN-14.2)
func TestRepositorioRondas_GuardarYBuscar(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana", "Beto")
	if r.ID == 0 {
		t.Fatal("GuardarRonda no asignó el ID")
	}

	got, err := b.repo.BuscarRonda(ctx, r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ProyectoID != r.ProyectoID || got.ElementoID != b.elemento.ID || got.Numero != 1 || got.Estado != domain.RondaAbierta {
		t.Errorf("ronda = %+v", got)
	}
	if len(got.Votantes) != 2 || got.Votantes[0].Nombre != "Ana" || got.Votantes[1].Nombre != "Beto" ||
		got.Votantes[0].Voto != nil || got.Votantes[1].Voto != nil {
		t.Errorf("votantes = %+v, se esperaba Ana y Beto sin votar, en orden", got.Votantes)
	}
	if _, err := b.repo.BuscarRonda(ctx, 99); !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Errorf("ronda inexistente: error = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
}

// Traza: HU-14 CA-14.6 (RN-14.1)
func TestRepositorioRondas_UltimaRonda(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	if r, err := b.repo.UltimaRonda(ctx, b.elemento.ID); err != nil || r != nil {
		t.Fatalf("sin rondas: %+v, %v; se esperaba nil", r, err)
	}

	primera := b.abrir(t, nil, "Ana")
	if err := b.votar(t, &primera, "Ana", 5); err != nil {
		t.Fatal(err)
	}
	segunda := b.abrir(t, &primera, "Ana")

	ultima, err := b.repo.UltimaRonda(ctx, b.elemento.ID)
	if err != nil || ultima == nil || ultima.ID != segunda.ID || ultima.Numero != 2 {
		t.Errorf("última ronda = %+v, %v; se esperaba la número 2", ultima, err)
	}
}

// Traza: HU-14 RN-14.1 garantizada también en la base: si dos pedidos abren una ronda sobre el
// mismo elemento a la vez, el segundo recibe ErrRondaYaAbierta.
func TestRepositorioRondas_UnaSolaAbiertaPorElemento(t *testing.T) {
	b := nuevaBasePoker(t)
	b.abrir(t, nil, "Ana")

	otra, err := domain.NuevaRonda(b.elemento, []domain.Integrante{{Nombre: "Ana"}}, nil) // sin ver la abierta
	if err != nil {
		t.Fatal(err)
	}
	otra.Numero = 2
	if err := b.repo.GuardarRonda(context.Background(), &otra); !errors.Is(err, domain.ErrRondaYaAbierta) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrRondaYaAbierta)
	}
}

// Traza: HU-14 CA-14.2 a CA-14.4 (RN-14.5, RN-14.6)
func TestRepositorioRondas_GuardarVoto(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana", "Beto")

	if err := b.votar(t, &r, "ana", 5); err != nil {
		t.Fatal(err)
	}
	if err := b.votar(t, &r, "Ana", 8); err != nil { // cambia su voto
		t.Fatal(err)
	}
	got, err := b.repo.BuscarRonda(ctx, r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Estado != domain.RondaAbierta || got.Votantes[0].Voto == nil || *got.Votantes[0].Voto != 8 || got.Votantes[1].Voto != nil {
		t.Errorf("tras los votos de Ana: %+v", got)
	}

	if err := b.votar(t, &r, "Beto", 3); err != nil {
		t.Fatal(err)
	}
	if got, err = b.repo.BuscarRonda(ctx, r.ID); err != nil || got.Estado != domain.RondaCerrada {
		t.Errorf("tras votar todos: %+v, %v; se esperaba cerrada", got, err)
	}
}

// Traza: HU-14 RN-14.6: un voto que llega después de que otro pedido cerró la ronda no se guarda.
func TestRepositorioRondas_GuardarVoto_RondaYaCerrada(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana")
	vieja := r // la que leyó el otro pedido, todavía abierta
	vieja.Votantes = append([]domain.Votante(nil), r.Votantes...)
	if err := b.votar(t, &r, "Ana", 5); err != nil {
		t.Fatal(err)
	}

	if err := b.votar(t, &vieja, "Ana", 8); !errors.Is(err, domain.ErrRondaCerrada) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrRondaCerrada)
	}
	got, err := b.repo.BuscarRonda(ctx, r.ID)
	if err != nil || *got.Votantes[0].Voto != 5 {
		t.Errorf("el voto guardado tendría que seguir siendo 5: %+v, %v", got.Votantes, err)
	}
}

// Traza: HU-14 sección 5: dos votos simultáneos que completan la ronda quedan los dos y la ronda
// se cierra una vez.
func TestRepositorioRondas_VotosSimultaneos(t *testing.T) {
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana", "Beto")
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i, n := range []string{"Ana", "Beto"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			copia := r
			copia.Votantes = append([]domain.Votante(nil), r.Votantes...)
			if err := copia.Votar(n, 5); err != nil {
				errs[i] = err
				return
			}
			errs[i] = b.repo.GuardarVoto(context.Background(), copia, n)
		}()
	}
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		t.Fatal(err)
	}
	got, err := b.repo.BuscarRonda(context.Background(), r.ID)
	if err != nil || got.Estado != domain.RondaCerrada || got.Votantes[0].Voto == nil || got.Votantes[1].Voto == nil {
		t.Errorf("ronda = %+v, %v; se esperaba cerrada con los dos votos", got, err)
	}
}

// Traza: HU-14 sección 5: eliminar el elemento (HU-06) borra sus rondas y no lo impide.
func TestRepositorioRondas_SeBorranConElElemento(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana", "Beto")
	if err := b.votar(t, &r, "Ana", 5); err != nil {
		t.Fatal(err)
	}

	if err := b.backlog.EliminarElemento(ctx, b.elemento.ID); err != nil {
		t.Fatalf("eliminar el elemento con una ronda: %v", err)
	}
	if _, err := b.repo.BuscarRonda(ctx, r.ID); !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Errorf("la ronda tendría que haberse borrado: %v", err)
	}
}

// Traza: HU-14 RN-14.5: guardar exige que la ronda traiga el voto del votante (lo aplica el dominio
// antes); si no, es un error interno y no "no es votante".
func TestRepositorioRondas_GuardarVoto_SinVotoAplicado(t *testing.T) {
	b := nuevaBasePoker(t)
	r := b.abrir(t, nil, "Ana")

	err := b.repo.GuardarVoto(context.Background(), r, "Ana")

	if err == nil || errors.Is(err, domain.ErrNoEsVotante) || errors.Is(err, domain.ErrRondaCerrada) {
		t.Errorf("error = %v, se esperaba un error interno", err)
	}
	if err := b.repo.GuardarVoto(context.Background(), r, "Dani"); !errors.Is(err, domain.ErrNoEsVotante) {
		t.Errorf("votante desconocido: error = %v, se esperaba %v", err, domain.ErrNoEsVotante)
	}
}
