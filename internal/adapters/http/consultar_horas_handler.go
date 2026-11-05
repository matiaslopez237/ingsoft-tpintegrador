package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorHorasProyecto y consultorHorasSprint son lo que los handlers necesitan de los casos de uso
// ConsultarHorasProyecto y ConsultarHorasSprint.
type (
	consultorHorasProyecto interface {
		Ejecutar(ctx context.Context, proyectoID int64) (domain.TotalHoras, error)
	}
	consultorHorasSprint interface {
		Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.TotalHoras, error)
	}
)

// respuestaTotalHoras son los campos del total (spec HU-23, sección 3). Los punteros salen como null
// cuando no hay ningún elemento estimado (RN-23.6).
type respuestaTotalHoras struct {
	HorasEstimadas        float64  `json:"horas_estimadas"`
	HorasReales           float64  `json:"horas_reales"`
	HorasRealesSinEstimar float64  `json:"horas_reales_sin_estimar"`
	DesviacionHoras       *float64 `json:"desviacion_horas"`
	DesviacionPorcentaje  *float64 `json:"desviacion_porcentaje"`
	Elementos             int      `json:"elementos"`
	ElementosSinEstimar   int      `json:"elementos_sin_estimar"`
}

// respuestaHorasProyecto es el JSON de GET /proyectos/{id}/horas.
type respuestaHorasProyecto struct {
	ProyectoID int64 `json:"proyecto_id"`
	respuestaTotalHoras
}

// respuestaHorasSprint es el JSON de GET /proyectos/{id}/sprints/{sprintID}/horas.
type respuestaHorasSprint struct {
	ProyectoID int64 `json:"proyecto_id"`
	SprintID   int64 `json:"sprint_id"`
	respuestaTotalHoras
}

func aRespuestaTotalHoras(t domain.TotalHoras) respuestaTotalHoras {
	return respuestaTotalHoras{
		HorasEstimadas:        t.HorasEstimadas,
		HorasReales:           t.HorasReales,
		HorasRealesSinEstimar: t.HorasRealesSinEstimar,
		DesviacionHoras:       t.DesviacionHoras,
		DesviacionPorcentaje:  t.DesviacionPorcentaje,
		Elementos:             t.Elementos,
		ElementosSinEstimar:   t.ElementosSinEstimar,
	}
}

// erroresConsultaHoras son los errores de HU-23 que se informan al cliente (spec, sección 7). Se
// responde con el mensaje del dominio y no con el que devolvió el caso de uso, que puede venir envuelto.
var erroresConsultaHoras = []error{domain.ErrProyectoNoEncontrado, domain.ErrSprintNoEncontrado}

// RegistrarConsultarHoras agrega al mux las rutas de HU-23: las horas de un proyecto y las de un sprint.
func RegistrarConsultarHoras(mux *http.ServeMux, proyecto consultorHorasProyecto, sprint consultorHorasSprint) {
	mux.HandleFunc("GET /proyectos/{id}/horas", handlerHorasProyecto(proyecto))
	mux.HandleFunc("GET /proyectos/{id}/sprints/{sprintID}/horas", handlerHorasSprint(sprint))
}

// handlerHorasProyecto atiende GET /proyectos/{id}/horas: 200 con el total, 400 si el ID no es válido,
// 404 si el proyecto no existe y 500 ante una falla interna. Implementa HU-23 (CA-23.16).
func handlerHorasProyecto(uc consultorHorasProyecto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		total, err := uc.Ejecutar(r.Context(), proyectoID)
		if err != nil {
			responderErrorHoras(w, err, fmt.Sprintf("proyecto %d", proyectoID))
			return
		}

		escribirJSON(w, http.StatusOK, respuestaHorasProyecto{ProyectoID: proyectoID, respuestaTotalHoras: aRespuestaTotalHoras(total)})
	}
}

// handlerHorasSprint atiende GET /proyectos/{id}/sprints/{sprintID}/horas: 200 con el total, 400 si los
// IDs no son válidos, 404 si el proyecto o el sprint no existen o el sprint es de otro proyecto, y 500
// ante una falla interna. Implementa HU-23 (CA-23.16).
func handlerHorasSprint(uc consultorHorasSprint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, err := proyectoIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		sprintID, err := sprintIDDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}

		total, err := uc.Ejecutar(r.Context(), proyectoID, sprintID)
		if err != nil {
			responderErrorHoras(w, err, fmt.Sprintf("sprint %d del proyecto %d", sprintID, proyectoID))
			return
		}

		escribirJSON(w, http.StatusOK, respuestaHorasSprint{ProyectoID: proyectoID, SprintID: sprintID, respuestaTotalHoras: aRespuestaTotalHoras(total)})
	}
}

// responderErrorHoras responde 404 con el mensaje del dominio si el error está en la tabla, o 500 sin el
// detalle, que queda solo en el log. alcance dice qué se estaba consultando.
func responderErrorHoras(w http.ResponseWriter, err error, alcance string) {
	for _, e := range erroresConsultaHoras {
		if errors.Is(err, e) {
			escribirJSON(w, http.StatusNotFound, respuestaError{e.Error()})
			return
		}
	}
	log.Printf("consultar horas del %s: %v", alcance, err)
	escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al consultar las horas"})
}
