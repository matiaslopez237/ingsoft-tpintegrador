package domain

import (
	"errors"
	"slices"
)

// Errores de HU-14 (spec, sección 7).
var (
	ErrRondaNoEncontrada = errors.New("ronda de estimación no encontrada")
	ErrElementoTerminado = errors.New("no se puede estimar un elemento terminado")
	ErrRondaSinVotantes  = errors.New("el proyecto no tiene integrantes para votar")
	ErrRondaYaAbierta    = errors.New("el elemento ya tiene una ronda abierta")
	ErrRondaCerrada      = errors.New("la ronda ya está cerrada")
	ErrNoEsVotante       = errors.New("el integrante no participa de esta ronda")
	ErrVotoInvalido      = errors.New("el voto debe ser 1, 2, 3, 5, 8, 13 o 21")
	ErrRondaAbierta      = errors.New("los votos se ven cuando votaron todos") // HU-15
)

// EstadoRonda es el estado de una ronda de Planning Poker. Implementa HU-14 (RN-14.6).
type EstadoRonda string

const (
	RondaAbierta EstadoRonda = "Abierta"
	RondaCerrada EstadoRonda = "Cerrada"
)

// valoresVoto es la escala de Fibonacci del Planning Poker (RN-14.4).
var valoresVoto = []int{1, 2, 3, 5, 8, 13, 21}

// ValoresVoto devuelve una copia de la escala, en orden, para las opciones de la pantalla.
func ValoresVoto() []int { return slices.Clone(valoresVoto) }

// Votante es un integrante convocado a la ronda. Voto queda en nil hasta que vota.
type Votante struct {
	Nombre string
	Voto   *int
}

// Ronda es una ronda de Planning Poker sobre un elemento del backlog. Los valores de los votos
// solo los usa el núcleo: lo que se muestra es Resumen, que no los trae (RN-14.3). Implementa HU-14.
type Ronda struct {
	ID         int64 // lo asigna la persistencia al guardar
	ProyectoID int64
	ElementoID int64
	Numero     int // 1, 2, 3… por elemento (RN-14.1)
	Estado     EstadoRonda
	Votantes   []Votante
	// EstimacionAcordada es lo que el equipo acordó en esta ronda (HU-17); nil si no se registró.
	EstimacionAcordada *int
}

// NuevaRonda abre una ronda sobre el elemento, con los integrantes del proyecto como votantes.
// anterior es la última ronda del elemento, o nil si no tuvo ninguna. Los errores se informan en
// el orden de la spec. Implementa HU-14 (RN-14.1, RN-14.2).
func NuevaRonda(e ElementoBacklog, integrantes []Integrante, anterior *Ronda) (Ronda, error) {
	switch {
	case e.Estado == EstadoHecho:
		return Ronda{}, ErrElementoTerminado
	case len(integrantes) == 0:
		return Ronda{}, ErrRondaSinVotantes
	case anterior != nil && anterior.Estado == RondaAbierta:
		return Ronda{}, ErrRondaYaAbierta
	}
	numero := 1
	if anterior != nil {
		numero = anterior.Numero + 1
	}
	votantes := make([]Votante, len(integrantes))
	for i, in := range integrantes {
		votantes[i] = Votante{Nombre: in.Nombre}
	}
	return Ronda{ProyectoID: e.ProyectoID, ElementoID: e.ID, Numero: numero, Estado: RondaAbierta, Votantes: votantes}, nil
}

// Votar registra el voto del integrante (el nombre se compara con ClaveNombre, como RN-03.3). Si
// ya había votado, queda el último valor; si con este votaron todos, la ronda se cierra. Si hay un
// error, la ronda no cambia. Implementa HU-14 (RN-14.2, RN-14.4 a RN-14.6).
func (r *Ronda) Votar(integrante string, valor int) error {
	if r.Estado != RondaAbierta {
		return ErrRondaCerrada
	}
	clave := ClaveNombre(integrante)
	i := slices.IndexFunc(r.Votantes, func(v Votante) bool { return ClaveNombre(v.Nombre) == clave })
	switch {
	case i < 0:
		return ErrNoEsVotante
	case !slices.Contains(valoresVoto, valor):
		return ErrVotoInvalido
	}
	r.Votantes[i].Voto = &valor
	if !slices.ContainsFunc(r.Votantes, func(v Votante) bool { return v.Voto == nil }) {
		r.Estado = RondaCerrada
	}
	return nil
}

// VotanteResumen dice de un votante solo si ya votó, nunca cuánto (RN-14.3).
type VotanteResumen struct {
	Nombre string
	YaVoto bool
}

// ResumenRonda es lo que se puede mostrar de una ronda: sin los valores de los votos mientras está
// abierta (RN-14.3); cerrada, con la revelación (RN-15.1).
type ResumenRonda struct {
	ID         int64
	ProyectoID int64
	ElementoID int64
	Numero     int
	Estado     EstadoRonda
	Votantes   []VotanteResumen
	Revelacion *Revelacion // nil mientras la ronda está abierta (HU-15)
	// EstimacionAcordada es lo que se acordó en esta ronda (HU-17); nil si no se registró.
	EstimacionAcordada *int
}

// Resumen devuelve los datos de la ronda que se pueden mostrar: los votos, solo si está cerrada.
// Implementa HU-14 (RN-14.3) y HU-15 (RN-15.1).
func (r Ronda) Resumen() ResumenRonda {
	votantes := make([]VotanteResumen, len(r.Votantes))
	for i, v := range r.Votantes {
		votantes[i] = VotanteResumen{Nombre: v.Nombre, YaVoto: v.Voto != nil}
	}
	res := ResumenRonda{ID: r.ID, ProyectoID: r.ProyectoID, ElementoID: r.ElementoID, Numero: r.Numero,
		Estado: r.Estado, Votantes: votantes}
	if rev, err := r.Revelar(); err == nil {
		res.Revelacion = &rev
	}
	if r.EstimacionAcordada != nil { // una copia: el resumen no comparte memoria con la ronda
		v := *r.EstimacionAcordada
		res.EstimacionAcordada = &v
	}
	return res
}

// VotoRevelado es el voto de un integrante, visible con la ronda cerrada (HU-15).
type VotoRevelado struct {
	Nombre string
	Valor  int
}

// Revelacion es lo que se muestra de una ronda cerrada (HU-15): los votos, el menor y el mayor, y
// si la diferencia entre esos dos es significativa.
type Revelacion struct {
	Votos                   []VotoRevelado
	Minimo, Maximo          int
	DiferenciaSignificativa bool
}

// Revelar devuelve los votos de la ronda en el orden de los votantes, el mínimo, el máximo y si el
// máximo es más del doble del mínimo. Con la ronda abierta, devuelve ErrRondaAbierta.
// Implementa HU-15 (RN-15.1 a RN-15.3).
func (r Ronda) Revelar() (Revelacion, error) {
	if r.Estado != RondaCerrada {
		return Revelacion{}, ErrRondaAbierta
	}
	rev := Revelacion{Votos: make([]VotoRevelado, 0, len(r.Votantes))}
	for _, v := range r.Votantes {
		if v.Voto == nil { // no pasa: la ronda se cierra cuando votaron todos (RN-14.6)
			continue
		}
		valor := *v.Voto
		rev.Votos = append(rev.Votos, VotoRevelado{Nombre: v.Nombre, Valor: valor})
		if len(rev.Votos) == 1 || valor < rev.Minimo {
			rev.Minimo = valor
		}
		rev.Maximo = max(rev.Maximo, valor)
	}
	rev.DiferenciaSignificativa = rev.Maximo > 2*rev.Minimo
	return rev, nil
}
