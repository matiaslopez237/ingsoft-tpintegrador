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

// ResumenRonda es lo que se puede mostrar de una ronda: sin los valores de los votos (RN-14.3).
type ResumenRonda struct {
	ID         int64
	ProyectoID int64
	ElementoID int64
	Numero     int
	Estado     EstadoRonda
	Votantes   []VotanteResumen
}

// Resumen devuelve los datos de la ronda que se pueden mostrar. Implementa HU-14 (RN-14.3).
func (r Ronda) Resumen() ResumenRonda {
	votantes := make([]VotanteResumen, len(r.Votantes))
	for i, v := range r.Votantes {
		votantes[i] = VotanteResumen{Nombre: v.Nombre, YaVoto: v.Voto != nil}
	}
	return ResumenRonda{ID: r.ID, ProyectoID: r.ProyectoID, ElementoID: r.ElementoID, Numero: r.Numero,
		Estado: r.Estado, Votantes: votantes}
}
