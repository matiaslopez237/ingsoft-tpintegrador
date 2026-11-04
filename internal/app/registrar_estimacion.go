package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// RepositorioEstimacion es el puerto angosto de HU-17: no se ensancha RepositorioRondas de HU-14,
// para no tocar sus fakes (spec HU-17, sección 5).
type RepositorioEstimacion interface {
	BuscarRonda(ctx context.Context, id int64) (domain.Ronda, error)
	UltimaRonda(ctx context.Context, elementoID int64) (*domain.Ronda, error)
	// GuardarEstimacion guarda en una transacción la estimación que r tiene y la copia a los story
	// points del elemento (RN-17.1). Vuelve a verificar RN-17.3 y RN-17.4 en los mismos UPDATE: si
	// otro pedido cambió algo en el medio, no guarda nada y devuelve el error de esa regla.
	GuardarEstimacion(ctx context.Context, r domain.Ronda) error
}

// DatosEstimacion son las entradas de HU-17 (spec, sección 2).
type DatosEstimacion struct {
	ProyectoID int64
	RondaID    int64
	Valor      int
}

// RegistrarEstimacion es el caso de uso de HU-17.
type RegistrarEstimacion struct {
	proyectos LectorProyectos
	elementos LectorElementos
	rondas    RepositorioEstimacion
}

func NuevoRegistrarEstimacion(proyectos LectorProyectos, elementos LectorElementos, rondas RepositorioEstimacion) *RegistrarEstimacion {
	return &RegistrarEstimacion{proyectos: proyectos, elementos: elementos, rondas: rondas}
}

// Ejecutar busca la ronda en el proyecto, su elemento y la última ronda del elemento, registra la
// estimación en el dominio y la guarda. Devuelve el resumen de la ronda. Si algo falla, no se guarda
// nada. Implementa HU-17 (RN-17.1 a RN-17.4).
func (uc *RegistrarEstimacion) Ejecutar(ctx context.Context, d DatosEstimacion) (domain.ResumenRonda, error) {
	r, err := rondaDelProyecto(ctx, uc.proyectos, uc.rondas, d.ProyectoID, d.RondaID)
	if err != nil {
		return domain.ResumenRonda{}, err
	}
	e, err := uc.elementos.BuscarElementoPorID(ctx, r.ElementoID)
	if err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("buscar elemento %d: %w", r.ElementoID, err)
	}
	ultima, err := uc.rondas.UltimaRonda(ctx, r.ElementoID)
	if err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("buscar la última ronda del elemento %d: %w", r.ElementoID, err)
	}
	if err := r.RegistrarEstimacion(d.Valor, ultima, e); err != nil {
		return domain.ResumenRonda{}, err
	}
	if err := uc.rondas.GuardarEstimacion(ctx, r); err != nil {
		return domain.ResumenRonda{}, fmt.Errorf("guardar estimación de la ronda %d: %w", r.ID, err)
	}
	return r.Resumen(), nil
}
