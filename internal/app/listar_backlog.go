package app

import (
	"context"
	"fmt"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// ListadorElementos es el puerto de solo lectura que necesita ListarBacklog (HU-07).
type ListadorElementos interface {
	// ListarElementos devuelve los elementos del proyecto con sus criterios, en cualquier orden.
	ListarElementos(ctx context.Context, proyectoID int64) ([]domain.ElementoBacklog, error)
}

// DatosListarBacklog son las entradas de HU-07. Orden y Estado vacíos toman los valores por defecto.
type DatosListarBacklog struct {
	ProyectoID int64
	Orden      string // "desc" (por defecto) o "asc"
	Estado     string // vacío: todos los estados
}

// ListarBacklog es el caso de uso de HU-07. Implementa HU-07 (RN-07.1).
type ListarBacklog struct {
	proyectos LectorProyectos
	elementos ListadorElementos
}

func NuevoListarBacklog(proyectos LectorProyectos, elementos ListadorElementos) *ListarBacklog {
	return &ListarBacklog{proyectos: proyectos, elementos: elementos}
}

// Ejecutar verifica que el proyecto exista, trae sus elementos y los ordena y filtra con el
// dominio (RN-07.2 a RN-07.5). El proyecto inexistente se informa antes que un orden o un
// estado inválidos (spec, sección 7).
func (uc *ListarBacklog) Ejecutar(ctx context.Context, d DatosListarBacklog) ([]domain.ElementoBacklog, error) {
	if _, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID); err != nil {
		return nil, err
	}
	elementos, err := uc.elementos.ListarElementos(ctx, d.ProyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar elementos del proyecto: %w", err)
	}
	return domain.OrdenarBacklog(elementos, d.Orden, d.Estado)
}
