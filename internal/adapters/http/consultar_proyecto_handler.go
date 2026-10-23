package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorProyectos es lo que el handler necesita del caso de uso ConsultarProyecto.
type consultorProyectos interface {
	Ejecutar(ctx context.Context, id int64) (domain.Proyecto, error)
}

// respuestaConsultaProyecto es el JSON de GET /proyectos/{id} (HU-04, sección 3). Embebe
// respuestaProyecto de HU-01, así los datos del proyecto salen con el mismo formato que en POST.
type respuestaConsultaProyecto struct {
	respuestaProyecto
	Integrantes []respuestaIntegrante `json:"integrantes"`
}

type respuestaIntegrante struct {
	Nombre string   `json:"nombre"`
	Roles  []string `json:"roles"`
}

// RegistrarConsultarProyecto agrega GET /proyectos/{id} al mux. Implementa HU-04.
func RegistrarConsultarProyecto(mux *http.ServeMux, uc consultorProyectos) {
	mux.HandleFunc("GET /proyectos/{id}", handlerConsultarProyecto(uc))
}

// handlerConsultarProyecto atiende GET /proyectos/{id}: 200 con el proyecto, 400 si el id no es
// válido, 404 si el proyecto no existe y 500 ante una falla interna.
// Implementa HU-04 (expone ConsultarProyecto por HTTP).
func handlerConsultarProyecto(uc consultorProyectos) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		p, err := uc.Ejecutar(r.Context(), id)
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			// El texto del error centinela, no err.Error(): el caso de uso lo devuelve envuelto.
			escribirJSON(w, http.StatusNotFound, respuestaError{domain.ErrProyectoNoEncontrado.Error()})
		case err != nil:
			log.Printf("consultar proyecto: %v", err)
			escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al consultar el proyecto"})
		default:
			escribirJSON(w, http.StatusOK, aRespuestaConsultaProyecto(p))
		}
	}
}

// aRespuestaConsultaProyecto arma la respuesta con integrantes siempre como array: un slice nil
// se codifica como null y la spec exige [] (CA-04.3).
func aRespuestaConsultaProyecto(p domain.Proyecto) respuestaConsultaProyecto {
	integrantes := make([]respuestaIntegrante, 0, len(p.Integrantes))
	for _, i := range p.Integrantes {
		integrantes = append(integrantes, respuestaIntegrante{Nombre: i.Nombre, Roles: i.Roles})
	}
	return respuestaConsultaProyecto{respuestaProyecto: aRespuestaProyecto(p), Integrantes: integrantes}
}
