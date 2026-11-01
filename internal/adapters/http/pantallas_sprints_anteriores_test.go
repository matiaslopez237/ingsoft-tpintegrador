package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const rutaAnterioresDeAlfa = "/web/proyectos/7/sprints/anteriores"

// casosConAnteriores suma a casosConSprints el caso de uso de HU-12, que devuelve los sprints
// indicados para alfa (ID 7) y ErrProyectoNoEncontrado para cualquier otro proyecto.
func casosConAnteriores(sprints []domain.Sprint, err error) CasosDeUso {
	cu := casosConSprints(nil, nil)
	cu.ListarSprintsAnteriores = listarSprintsFunc(func(_ context.Context, id int64) ([]domain.Sprint, error) {
		if id != alfa.ID {
			return nil, fmt.Errorf("buscar proyecto %d: %w", id, domain.ErrProyectoNoEncontrado)
		}
		return sprints, err
	})
	return cu
}

func sprintAnteriorDeAlfa(id int64, numero, diaInicio int, objetivo string, planificados, completados int) domain.Sprint {
	inicio := time.Date(2026, 10, diaInicio, 0, 0, 0, 0, time.UTC)
	return domain.Sprint{ID: id, ProyectoID: alfa.ID, Numero: numero, FechaInicio: inicio, FechaFin: inicio.AddDate(0, 0, 11),
		Objetivo: objetivo, Estado: domain.SprintCerrado, PuntosPlanificados: planificados, PuntosCompletados: completados}
}

// Traza: HU-12 CA-12.7: una fila por sprint cerrado, en el orden del caso de uso, con sus datos y
// un enlace a su detalle (HU-09).
func TestPantallaSprintsAnteriores_MuestraLosSprintsConSusDatos(t *testing.T) {
	cu := casosConAnteriores([]domain.Sprint{
		sprintAnteriorDeAlfa(2, 2, 13, "Interfaz web", 8, 5),
		sprintAnteriorDeAlfa(1, 1, 1, "Armar el MVP", 13, 10),
	}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Sprints anteriores", "Alfa",
		"2026-10-13", "2026-10-24", "Interfaz web", "2026-10-01", "2026-10-12", "Armar el MVP",
		`href="/web/proyectos/7/sprints/2"`, `href="/web/proyectos/7/sprints/1"`)
	cuerpo := rec.Body.String()
	if strings.Index(cuerpo, "Interfaz web") > strings.Index(cuerpo, "Armar el MVP") {
		t.Error("el sprint 2 tendría que aparecer antes que el 1, como los devuelve el caso de uso")
	}
	for _, celda := range []string{"<td>8</td>", "<td>5</td>", "<td>13</td>", "<td>10</td>"} {
		if !strings.Contains(cuerpo, celda) {
			t.Errorf("falta la celda %s con los puntos", celda)
		}
	}
}

// Traza: HU-12 CA-12.7 (CA-12.2): sin sprints cerrados, un mensaje en lugar de una tabla vacía.
func TestPantallaSprintsAnteriores_SinSprints_MuestraMensaje(t *testing.T) {
	rec := pedirPantalla(t, casosConAnteriores(nil, nil), http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	contiene(t, rec, "todavía no tiene sprints cerrados")
	noContiene(t, rec.Body.String(), "<table")
}

// Traza: HU-12 sección 7: proyecto inexistente (404), ID inválido (400) y falla interna (500, sin
// el detalle).
func TestPantallaSprintsAnteriores_Errores(t *testing.T) {
	tests := []struct {
		nombre string
		ruta   string
		err    error
		want   int
		msg    string
	}{
		{"proyecto inexistente", "/web/proyectos/99/sprints/anteriores", nil, http.StatusNotFound, "proyecto no encontrado"},
		{"id inválido", "/web/proyectos/abc/sprints/anteriores", nil, http.StatusBadRequest, "número"},
		{"falla interna", rutaAnterioresDeAlfa, errors.New("disco lleno"), http.StatusInternalServerError, "error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConAnteriores(nil, tt.err), http.MethodGet, tt.ruta, nil)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.300s", rec.Code, tt.want, rec.Body)
			}
			contiene(t, rec, tt.msg)
			noContiene(t, rec.Body.String(), "disco lleno")
		})
	}
}

// Traza: HU-12 CA-12.7: se llega desde el detalle del proyecto y desde la lista de sprints (HU-08).
func TestPantallaSprintsAnteriores_EnlazadaDesdeElProyectoYLosSprints(t *testing.T) {
	cu := casosConAnteriores(nil, nil)
	for _, ruta := range []string{"/web/proyectos/7", "/web/proyectos/7/sprints"} {
		t.Run(ruta, func(t *testing.T) {
			rec := pedirPantalla(t, cu, http.MethodGet, ruta, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200", rec.Code)
			}
			contiene(t, rec, `href="`+rutaAnterioresDeAlfa+`"`)
		})
	}
}
