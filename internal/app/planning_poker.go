package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioRondas es el puerto de persistencia de las rondas de Planning Poker (HU-14).
type RepositorioRondas interface {
	// UltimaRonda devuelve la última ronda del elemento, o nil si no tuvo ninguna.
	UltimaRonda(ctx context.Context, elementoID int64) (*domain.Ronda, error)
	// GuardarRonda inserta la ronda con sus votantes y le asigna su ID. Si otro pedido abrió una
	// ronda sobre el mismo elemento en el medio, devuelve domain.ErrRondaYaAbierta.
	GuardarRonda(ctx context.Context, r *domain.Ronda) error
	// BuscarRonda devuelve domain.ErrRondaNoEncontrada si no existe.
	BuscarRonda(ctx context.Context, id int64) (domain.Ronda, error)
	// GuardarVoto guarda el voto que el votante tiene en r y cierra la ronda si ya votaron todos.
	// Si la ronda se cerró en el medio, devuelve domain.ErrRondaCerrada.
	GuardarVoto(ctx context.Context, r domain.Ronda, votante string) error
}

// DatosIniciarRonda son las entradas para abrir una ronda (HU-14, sección 2).
type DatosIniciarRonda struct {
	ProyectoID int64
	ElementoID int64
}

// DatosVoto son las entradas para votar en una ronda (HU-14, sección 2).
type DatosVoto struct {
	ProyectoID int64
	RondaID    int64
	Integrante string
	Valor      int
}

// IniciarRonda es el caso de uso de HU-14 para abrir una ronda sobre un elemento del backlog.
type IniciarRonda struct {
	proyectos LectorProyectos
	elementos LectorElementos
	rondas    RepositorioRondas
}

func NuevoIniciarRonda(proyectos LectorProyectos, elementos LectorElementos, rondas RepositorioRondas) *IniciarRonda {
	return &IniciarRonda{proyectos: proyectos, elementos: elementos, rondas: rondas}
}

// Ejecutar busca el proyecto (sus integrantes son los votantes) y el elemento, abre la ronda en el
// dominio y la guarda. Devuelve el resumen, sin votos. Implementa HU-14 (RN-14.1, RN-14.2, RN-14.7).
func (uc *IniciarRonda) Ejecutar(ctx context.Context, d DatosIniciarRonda) (domain.ResumenRonda, error) {
	p, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}
	e, err := buscarElementoDelProyecto(ctx, uc.elementos, d.ProyectoID, d.ElementoID) // RN-14.7, como RN-06.6
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	anterior, err := uc.rondas.UltimaRonda(ctx, e.ID)
	if err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("buscar la última ronda del elemento %d: %w", e.ID, err)
	}
	r, err := domain.NuevaRonda(e, p.Integrantes, anterior)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	if err := uc.rondas.GuardarRonda(ctx, &r); err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("guardar ronda del elemento %d: %w", e.ID, err)
	}
	return r.Resumen(), nil
}

// Votar es el caso de uso de HU-14 para registrar el voto de un integrante.
type Votar struct {
	proyectos LectorProyectos
	rondas    RepositorioRondas
}

func NuevoVotar(proyectos LectorProyectos, rondas RepositorioRondas) *Votar {
	return &Votar{proyectos: proyectos, rondas: rondas}
}

// Ejecutar busca la ronda en el proyecto, aplica el voto en el dominio, lo guarda y devuelve el
// resumen como quedó guardado, sin los valores. Si algo falla, no se guarda nada.
// Implementa HU-14 (RN-14.2 a RN-14.7).
func (uc *Votar) Ejecutar(ctx context.Context, d DatosVoto) (domain.ResumenRonda, error) {
	r, err := rondaDelProyecto(ctx, uc.proyectos, uc.rondas, d.ProyectoID, d.RondaID)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	if err := r.Votar(d.Integrante, d.Valor); err != nil {
		return domain.ResumenRonda{}, err
	}
	if err := uc.rondas.GuardarVoto(ctx, r, d.Integrante); err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("guardar voto en la ronda %d: %w", d.RondaID, err)
	}
	// Se relee porque otro voto simultáneo pudo cerrar la ronda.
	guardada, err := uc.rondas.BuscarRonda(ctx, d.RondaID)
	if err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("releer ronda %d: %w", d.RondaID, err)
	}
	return guardada.Resumen(), nil
}

// ConsultarRonda es el caso de uso de HU-14 para ver una ronda, sin los valores de los votos.
type ConsultarRonda struct {
	proyectos LectorProyectos
	elementos LectorElementos
	rondas    RepositorioRondas
}

func NuevoConsultarRonda(proyectos LectorProyectos, elementos LectorElementos, rondas RepositorioRondas) *ConsultarRonda {
	return &ConsultarRonda{proyectos: proyectos, elementos: elementos, rondas: rondas}
}

// Ejecutar devuelve el resumen de la ronda del proyecto. Implementa HU-14 (RN-14.3, RN-14.7).
func (uc *ConsultarRonda) Ejecutar(ctx context.Context, proyectoID, rondaID int64) (domain.ResumenRonda, error) {
	r, err := rondaDelProyecto(ctx, uc.proyectos, uc.rondas, proyectoID, rondaID)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	return r.Resumen(), nil
}

// UltimaDelElemento devuelve el resumen de la última ronda del elemento, o nil si no tuvo
// ninguna. Lo usa la pantalla del elemento para ofrecer iniciar una ronda o ir a la abierta.
func (uc *ConsultarRonda) UltimaDelElemento(ctx context.Context, proyectoID, elementoID int64) (*domain.ResumenRonda, error) {
	e, err := buscarElementoDelProyecto(ctx, uc.elementos, proyectoID, elementoID)
	if err != nil {
		return nil, err
	}
	r, err := uc.rondas.UltimaRonda(ctx, e.ID)
	if err != nil || r == nil {
		return nil, err
	}
	res := r.Resumen()
	return &res, nil
}

// lectorRonda es lo que necesita rondaDelProyecto, compartida con HU-17.
type lectorRonda interface {
	BuscarRonda(ctx context.Context, id int64) (domain.Ronda, error)
}

// rondaDelProyecto verifica que exista el proyecto y devuelve la ronda si es suya; una de otro
// proyecto se trata como inexistente (RN-14.7).
func rondaDelProyecto(ctx context.Context, proyectos LectorProyectos, rondas lectorRonda, proyectoID, rondaID int64) (domain.Ronda, error) {
	if _, err := proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return domain.Ronda{}, fmt.Errorf("buscar proyecto %d: %w", proyectoID, err)
	}
	r, err := rondas.BuscarRonda(ctx, rondaID)
	if err != nil {
		return domain.Ronda{}, fmt.Errorf("buscar ronda %d: %w", rondaID, err)
	}
	if r.ProyectoID != proyectoID {
		return domain.Ronda{}, fmt.Errorf("ronda %d en el proyecto %d: %w", rondaID, proyectoID, domain.ErrRondaNoEncontrada)
	}
	return r, nil
}
