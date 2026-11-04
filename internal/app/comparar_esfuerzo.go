package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorHorasReales es el puerto para consultar las horas registradas de un elemento (HU-19). Es
// angosto a propósito: no se ensancha RepositorioEsfuerzo de HU-18 (spec HU-19, sección 5).
type LectorHorasReales interface {
	// HorasRealesDeElemento devuelve la suma de las horas registradas; 0 si no hay registros (RN-19.4).
	HorasRealesDeElemento(ctx context.Context, elementoID int64) (float64, error)
}

// CompararEsfuerzo es el caso de uso de HU-19. Reutiliza los puertos de lectura de HU-04 y HU-18.
// Implementa HU-19 (RN-19.4, RN-19.10, RN-19.11).
type CompararEsfuerzo struct {
	proyectos LectorProyectos
	elementos LectorElementos
	horas     LectorHorasReales
}

func NuevoCompararEsfuerzo(proyectos LectorProyectos, elementos LectorElementos, horas LectorHorasReales) *CompararEsfuerzo {
	return &CompararEsfuerzo{proyectos: proyectos, elementos: elementos, horas: horas}
}

// Ejecutar verifica que el proyecto exista y que el elemento sea suyo, consulta las horas reales y
// las compara con las estimadas del elemento. Los errores se informan en el orden de la tabla de la
// spec. El estado del elemento y el del proyecto no importan (RN-19.11).
func (uc *CompararEsfuerzo) Ejecutar(ctx context.Context, proyectoID, elementoID int64) (domain.ComparacionEsfuerzo, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, proyectoID); err != nil {
		return domain.ComparacionEsfuerzo{}, err
	}
	e, err := uc.elementos.BuscarElementoPorID(ctx, elementoID)
	if err != nil {
		return domain.ComparacionEsfuerzo{}, err
	}
	if e.ProyectoID != proyectoID { // RN-19.10: de otro proyecto es "inexistente"
		return domain.ComparacionEsfuerzo{}, domain.ErrElementoNoEncontrado
	}

	reales, err := uc.horas.HorasRealesDeElemento(ctx, elementoID)
	if err != nil {
		return domain.ComparacionEsfuerzo{}, fmt.Errorf("consultar horas reales del elemento %d: %w", elementoID, err)
	}
	return domain.CompararEsfuerzo(e.HorasEstimadas, reales), nil
}
