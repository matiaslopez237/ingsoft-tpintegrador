package app

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// pokerFalso guarda proyectos, elementos y rondas en memoria y permite simular fallas. Implementa
// LectorProyectos, LectorElementos y RepositorioRondas.
type pokerFalso struct {
	proyectos map[int64]domain.Proyecto
	elementos map[int64]domain.ElementoBacklog
	rondas    []domain.Ronda
	votos     []string // votantes de cada GuardarVoto, en orden
	errBuscar error    // al buscar el proyecto
	errRonda  error    // al leer o guardar rondas
}

func nuevoPokerFalso(t *testing.T, votantes ...string) *pokerFalso {
	t.Helper()
	p := domain.Proyecto{ID: 1, Nombre: "Alfa"}
	for _, n := range votantes {
		p.Integrantes = append(p.Integrantes, domain.Integrante{Nombre: n, Roles: []string{"Desarrollo"}})
	}
	return &pokerFalso{
		proyectos: map[int64]domain.Proyecto{1: p, 2: {ID: 2, Nombre: "Beta"}},
		elementos: map[int64]domain.ElementoBacklog{7: elementoDePrueba(t, 7, 1, "Login")},
	}
}

func (f *pokerFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errBuscar != nil {
		return domain.Proyecto{}, f.errBuscar
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *pokerFalso) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *pokerFalso) UltimaRonda(_ context.Context, elementoID int64) (*domain.Ronda, error) {
	if f.errRonda != nil {
		return nil, f.errRonda
	}
	for i := len(f.rondas) - 1; i >= 0; i-- {
		if f.rondas[i].ElementoID == elementoID {
			r := f.copia(f.rondas[i])
			return &r, nil
		}
	}
	return nil, nil
}

func (f *pokerFalso) GuardarRonda(_ context.Context, r *domain.Ronda) error {
	if f.errRonda != nil {
		return f.errRonda
	}
	r.ID = int64(len(f.rondas) + 1)
	f.rondas = append(f.rondas, f.copia(*r))
	return nil
}

func (f *pokerFalso) BuscarRonda(_ context.Context, id int64) (domain.Ronda, error) {
	if f.errRonda != nil {
		return domain.Ronda{}, f.errRonda
	}
	if id < 1 || int(id) > len(f.rondas) {
		return domain.Ronda{}, domain.ErrRondaNoEncontrada
	}
	return f.copia(f.rondas[id-1]), nil
}

func (f *pokerFalso) GuardarVoto(_ context.Context, r domain.Ronda, votante string) error {
	if f.errRonda != nil {
		return f.errRonda
	}
	f.votos = append(f.votos, votante)
	f.rondas[r.ID-1] = f.copia(r)
	return nil
}

// copia evita que la ronda guardada comparta los votantes con la del caso de uso.
func (f *pokerFalso) copia(r domain.Ronda) domain.Ronda {
	r.Votantes = slices.Clone(r.Votantes)
	return r
}

func casosPoker(f *pokerFalso) (*IniciarRonda, *Votar, *ConsultarRonda) {
	return NuevoIniciarRonda(f, f, f), NuevoVotar(f, f), NuevoConsultarRonda(f, f, f)
}

// Traza: HU-14 CA-14.1 (RN-14.1, RN-14.2)
func TestIniciarRonda_GuardaLaRondaAbiertaConLosIntegrantes(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana", "Beto")
	iniciar, _, _ := casosPoker(f)

	res, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: 7})

	if err != nil {
		t.Fatal(err)
	}
	want := []domain.VotanteResumen{{Nombre: "Ana"}, {Nombre: "Beto"}}
	if res.ID != 1 || res.Numero != 1 || res.Estado != domain.RondaAbierta || !slices.Equal(res.Votantes, want) {
		t.Errorf("resumen = %+v", res)
	}
	if len(f.rondas) != 1 || f.rondas[0].ElementoID != 7 {
		t.Errorf("rondas guardadas = %+v", f.rondas)
	}
}

// Traza: HU-14 CA-14.6, CA-14.7 (RN-14.1, RN-14.2, RN-14.7)
func TestIniciarRonda_Errores(t *testing.T) {
	tests := []struct {
		nombre    string
		votantes  []string
		preparar  func(*pokerFalso)
		datos     DatosIniciarRonda
		wantErr   error
		hayRondas int
	}{
		{"proyecto inexistente", []string{"Ana"}, nil, DatosIniciarRonda{ProyectoID: 9, ElementoID: 7}, domain.ErrProyectoNoEncontrado, 0},
		{"elemento inexistente", []string{"Ana"}, nil, DatosIniciarRonda{ProyectoID: 1, ElementoID: 9}, domain.ErrElementoNoEncontrado, 0},
		{"elemento de otro proyecto", []string{"Ana"}, nil, DatosIniciarRonda{ProyectoID: 2, ElementoID: 7}, domain.ErrElementoNoEncontrado, 0},
		{"elemento terminado", []string{"Ana"}, func(f *pokerFalso) {
			e := f.elementos[7]
			e.Estado = domain.EstadoHecho
			f.elementos[7] = e
		}, DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}, domain.ErrElementoTerminado, 0},
		{"sin integrantes", nil, nil, DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}, domain.ErrRondaSinVotantes, 0},
		{"ya hay una ronda abierta", []string{"Ana"}, func(f *pokerFalso) {
			f.rondas = []domain.Ronda{{ID: 1, ProyectoID: 1, ElementoID: 7, Numero: 1, Estado: domain.RondaAbierta}}
		}, DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}, domain.ErrRondaYaAbierta, 1},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoPokerFalso(t, tt.votantes...)
			if tt.preparar != nil {
				tt.preparar(f)
			}
			iniciar, _, _ := casosPoker(f)
			_, err := iniciar.Ejecutar(context.Background(), tt.datos)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if len(f.rondas) != tt.hayRondas {
				t.Errorf("se guardaron rondas: %+v", f.rondas)
			}
		})
	}
}

// Traza: HU-14 CA-14.6 (RN-14.1)
func TestIniciarRonda_ConLaAnteriorCerrada_AbreLaSiguiente(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana")
	f.rondas = []domain.Ronda{{ID: 1, ProyectoID: 1, ElementoID: 7, Numero: 1, Estado: domain.RondaCerrada}}
	iniciar, _, _ := casosPoker(f)

	res, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: 7})

	if err != nil || res.Numero != 2 || res.ID != 2 {
		t.Errorf("resumen = %+v, err = %v; se esperaba la ronda 2", res, err)
	}
}

// iniciada devuelve el falso con una ronda abierta sobre el elemento 7 del proyecto 1.
func iniciada(t *testing.T, votantes ...string) *pokerFalso {
	t.Helper()
	f := nuevoPokerFalso(t, votantes...)
	iniciar, _, _ := casosPoker(f)
	if _, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}); err != nil {
		t.Fatal(err)
	}
	return f
}

// Traza: HU-14 CA-14.2, CA-14.3 (RN-14.3, RN-14.6)
func TestVotar_GuardaElVotoYDevuelveElResumenSinValores(t *testing.T) {
	f := iniciada(t, "Ana", "Beto")
	_, votar, _ := casosPoker(f)

	res, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 5})
	if err != nil {
		t.Fatal(err)
	}
	if res.Estado != domain.RondaAbierta || !slices.Equal(res.Votantes, []domain.VotanteResumen{{Nombre: "Ana", YaVoto: true}, {Nombre: "Beto"}}) {
		t.Errorf("resumen tras el primer voto = %+v", res)
	}

	res, err = votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Beto", Valor: 8})
	if err != nil {
		t.Fatal(err)
	}
	if res.Estado != domain.RondaCerrada {
		t.Errorf("resumen tras votar todos = %+v, se esperaba cerrada", res)
	}
	if !slices.Equal(f.votos, []string{"Ana", "Beto"}) {
		t.Errorf("votos guardados = %v", f.votos)
	}
}

// Traza: HU-14 CA-14.5, CA-14.7 (RN-14.2, RN-14.4, RN-14.6, RN-14.7)
func TestVotar_Errores(t *testing.T) {
	tests := []struct {
		nombre  string
		datos   DatosVoto
		wantErr error
	}{
		{"proyecto inexistente", DatosVoto{ProyectoID: 9, RondaID: 1, Integrante: "Ana", Valor: 5}, domain.ErrProyectoNoEncontrado},
		{"ronda inexistente", DatosVoto{ProyectoID: 1, RondaID: 9, Integrante: "Ana", Valor: 5}, domain.ErrRondaNoEncontrada},
		{"ronda de otro proyecto", DatosVoto{ProyectoID: 2, RondaID: 1, Integrante: "Ana", Valor: 5}, domain.ErrRondaNoEncontrada},
		{"no es votante", DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Dani", Valor: 5}, domain.ErrNoEsVotante},
		{"voto inválido", DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 4}, domain.ErrVotoInvalido},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := iniciada(t, "Ana", "Beto")
			_, votar, _ := casosPoker(f)
			_, err := votar.Ejecutar(context.Background(), tt.datos)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if len(f.votos) != 0 {
				t.Errorf("no se esperaba guardar votos: %v", f.votos)
			}
		})
	}
}

// Traza: HU-14 CA-14.5 (RN-14.6)
func TestVotar_RondaCerrada(t *testing.T) {
	f := iniciada(t, "Ana")
	_, votar, _ := casosPoker(f)
	if _, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 5}); err != nil {
		t.Fatal(err)
	}

	_, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 3})

	if !errors.Is(err, domain.ErrRondaCerrada) || len(f.votos) != 1 {
		t.Errorf("err = %v, votos = %v", err, f.votos)
	}
}

// Traza: HU-14 CA-14.7 (RN-14.3, RN-14.7)
func TestConsultarRonda(t *testing.T) {
	f := iniciada(t, "Ana")
	_, _, consultar := casosPoker(f)

	res, err := consultar.Ejecutar(context.Background(), 1, 1)
	if err != nil || res.ID != 1 || res.Estado != domain.RondaAbierta {
		t.Errorf("resumen = %+v, err = %v", res, err)
	}
	for _, d := range []struct {
		proyecto, ronda int64
		wantErr         error
	}{{9, 1, domain.ErrProyectoNoEncontrado}, {2, 1, domain.ErrRondaNoEncontrada}, {1, 9, domain.ErrRondaNoEncontrada}} {
		if _, err := consultar.Ejecutar(context.Background(), d.proyecto, d.ronda); !errors.Is(err, d.wantErr) {
			t.Errorf("proyecto %d, ronda %d: error = %v, se esperaba %v", d.proyecto, d.ronda, err, d.wantErr)
		}
	}
}

// Traza: HU-14 CA-14.9: la pantalla del elemento muestra su última ronda, si tiene.
func TestConsultarRonda_UltimaDelElemento(t *testing.T) {
	f := nuevoPokerFalso(t, "Ana")
	_, _, consultar := casosPoker(f)

	if res, err := consultar.UltimaDelElemento(context.Background(), 1, 7); err != nil || res != nil {
		t.Errorf("sin rondas: resumen = %+v, err = %v; se esperaba nil", res, err)
	}

	f = iniciada(t, "Ana")
	_, _, consultar = casosPoker(f)
	res, err := consultar.UltimaDelElemento(context.Background(), 1, 7)
	if err != nil || res == nil || res.ID != 1 {
		t.Errorf("con una ronda: resumen = %+v, err = %v", res, err)
	}
	if _, err := consultar.UltimaDelElemento(context.Background(), 2, 7); !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("elemento de otro proyecto: error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

func TestPlanningPoker_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	f := iniciada(t, "Ana", "Beto")
	iniciar, votar, consultar := casosPoker(f)
	f.errRonda = errBase

	if _, err := iniciar.Ejecutar(context.Background(), DatosIniciarRonda{ProyectoID: 1, ElementoID: 7}); !errors.Is(err, errBase) {
		t.Errorf("iniciar: error = %v", err)
	}
	if _, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 5}); !errors.Is(err, errBase) {
		t.Errorf("votar: error = %v", err)
	}
	if _, err := consultar.Ejecutar(context.Background(), 1, 1); !errors.Is(err, errBase) {
		t.Errorf("consultar: error = %v", err)
	}
	f.errRonda, f.errBuscar = nil, errBase
	if _, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 5}); !errors.Is(err, errBase) {
		t.Errorf("votar con falla al buscar el proyecto: error = %v", err)
	}
}

// Traza: HU-15 CA-15.1, CA-15.4 (RN-15.1): ConsultarRonda devuelve la revelación solo con la ronda
// cerrada. No hizo falta cambiar el caso de uso: la decide Ronda.Resumen.
func TestConsultarRonda_RevelaSoloConLaRondaCerrada(t *testing.T) {
	f := iniciada(t, "Ana", "Beto")
	_, votar, consultar := casosPoker(f)
	if _, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Ana", Valor: 3}); err != nil {
		t.Fatal(err)
	}
	if res, err := consultar.Ejecutar(context.Background(), 1, 1); err != nil || res.Revelacion != nil {
		t.Fatalf("ronda abierta: revelación = %+v, err = %v; se esperaba nil", res.Revelacion, err)
	}

	cerrada, err := votar.Ejecutar(context.Background(), DatosVoto{ProyectoID: 1, RondaID: 1, Integrante: "Beto", Valor: 8})
	if err != nil {
		t.Fatal(err)
	}
	res, err := consultar.Ejecutar(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range []domain.ResumenRonda{cerrada, res} {
		if r.Revelacion == nil || r.Revelacion.Minimo != 3 || r.Revelacion.Maximo != 8 || !r.Revelacion.DiferenciaSignificativa {
			t.Errorf("revelación = %+v, se esperaba mínimo 3, máximo 8 y diferencia", r.Revelacion)
		}
	}
}
