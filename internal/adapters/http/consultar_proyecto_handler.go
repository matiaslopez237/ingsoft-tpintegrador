package http

import (
	"context"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorProyectos es lo que el handler necesita del caso de uso ConsultarProyecto.
type consultorProyectos interface {
	Ejecutar(ctx context.Context, id int64) (domain.Proyecto, error)
}

// handlerConsultarProyecto atiende GET /proyectos/{id}.
// Pendiente: se implementa en la etapa GREEN de la Tarea 4.5.
func handlerConsultarProyecto(uc consultorProyectos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
