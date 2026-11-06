package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const rutaAnterioresDeAlfa = "/web/proyectos/7/sprints/anteriores"

type consultarVelocidadFunc func(context.Context, int64) (app.VelocidadProyecto, error)

func (f consultarVelocidadFunc) Ejecutar(ctx context.Context, id int64) (app.VelocidadProyecto, error) {
	return f(ctx, id)
}

type consultarCalidadSprintsFunc func(context.Context, int64) ([]app.CalidadSprint, error)

func (f consultarCalidadSprintsFunc) Ejecutar(ctx context.Context, id int64) ([]app.CalidadSprint, error) {
	return f(ctx, id)
}

// casosConAnteriores suma a casosConSprints los casos de uso de HU-12, HU-22 y HU-24, que devuelven
// los sprints indicados (la velocidad y la calidad que les corresponde, calculadas con el dominio
// real, igual que los casos de uso de producción) para alfa (ID 7), y ErrProyectoNoEncontrado para
// cualquier otro. defectos es opcional: solo lo usan los escenarios de HU-24.
func casosConAnteriores(sprints []domain.Sprint, err error, defectos ...domain.Defecto) CasosDeUso {
	cu := casosConSprints(nil, nil)
	cu.ListarSprintsAnteriores = listarSprintsFunc(func(_ context.Context, id int64) ([]domain.Sprint, error) {
		if id != alfa.ID {
			return nil, fmt.Errorf("buscar proyecto %d: %w", id, domain.ErrProyectoNoEncontrado)
		}
		return sprints, err
	})
	cu.ConsultarVelocidad = consultarVelocidadFunc(func(_ context.Context, id int64) (app.VelocidadProyecto, error) {
		if id != alfa.ID {
			return app.VelocidadProyecto{}, fmt.Errorf("buscar proyecto %d: %w", id, domain.ErrProyectoNoEncontrado)
		}
		if err != nil {
			return app.VelocidadProyecto{}, err
		}
		cerrados := domain.SprintsAnteriores(sprints)
		velocidad, cantidad, hayDatos := domain.Velocidad(cerrados)
		return app.VelocidadProyecto{Velocidad: velocidad, HayDatos: hayDatos, SprintsCerrados: cantidad}, nil
	})
	cu.ConsultarCalidadSprints = consultarCalidadSprintsFunc(func(_ context.Context, id int64) ([]app.CalidadSprint, error) {
		if id != alfa.ID {
			return nil, fmt.Errorf("buscar proyecto %d: %w", id, domain.ErrProyectoNoEncontrado)
		}
		if err != nil {
			return nil, err
		}
		cerrados := domain.SprintsAnteriores(sprints)
		resumen := make([]app.CalidadSprint, len(cerrados))
		for i, s := range cerrados {
			porcentaje, hayPorcentaje := domain.PorcentajeCompletado(s)
			detectados, resueltos := domain.ContarDefectosDeSprint(defectos, s.ID)
			resumen[i] = app.CalidadSprint{
				Sprint: s, PorcentajeCompletado: porcentaje, HayPorcentaje: hayPorcentaje,
				DefectosDetectados: detectados, DefectosResueltos: resueltos,
			}
		}
		return resumen, nil
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

// Traza: HU-22 CA-22.8: con sprints cerrados, se muestra la velocidad del equipo.
func TestPantallaSprintsAnteriores_ConSprintsCerrados_MuestraLaVelocidad(t *testing.T) {
	cu := casosConAnteriores([]domain.Sprint{
		sprintAnteriorDeAlfa(2, 2, 13, "Interfaz web", 8, 5),
		sprintAnteriorDeAlfa(1, 1, 1, "Armar el MVP", 13, 10),
	}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Velocidad del equipo", "7,5") // promedio de 5 y 10 completados
}

// Traza: HU-22 CA-22.8 (RN-22.3): sin sprints cerrados, se avisa que no hay datos, no un 0.
func TestPantallaSprintsAnteriores_SinSprintsCerrados_AvisaSinVelocidad(t *testing.T) {
	rec := pedirPantalla(t, casosConAnteriores(nil, nil), http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	contiene(t, rec, "Todavía no hay datos de velocidad")
}

// Traza: HU-24 CA-24.9: cada fila muestra su porcentaje completado y sus defectos detectados y
// resueltos.
func TestPantallaSprintsAnteriores_MuestraElPorcentajeYLosDefectos(t *testing.T) {
	sprint := sprintAnteriorDeAlfa(1, 1, 1, "Armar el MVP", 8, 6)
	cu := casosConAnteriores([]domain.Sprint{sprint}, nil,
		domain.Defecto{ProyectoID: alfa.ID, SprintDeteccionID: &sprint.ID},
		domain.Defecto{ProyectoID: alfa.ID, SprintDeteccionID: &sprint.ID, SprintResolucionID: &sprint.ID, Estado: domain.DefectoCerrado},
	)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "75,0 %", "2", "1")
}

// Traza: HU-24 CA-24.1, CA-24.9: el porcentaje se redondea a un decimal, igual que en el JSON.
func TestPantallaSprintsAnteriores_PorcentajeConDecimales_SeRedondeaAUnDecimal(t *testing.T) {
	cu := casosConAnteriores([]domain.Sprint{sprintAnteriorDeAlfa(1, 1, 1, "Con decimales", 16, 1)}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "6,3 %")
	noContiene(t, rec.Body.String(), "6,25 %")
}

// Traza: HU-24 CA-24.2, CA-24.9: un sprint cerrado sin puntos planificados avisa que no hay
// porcentaje, no un 0%.
func TestPantallaSprintsAnteriores_SinPuntosPlanificados_AvisaSinPorcentaje(t *testing.T) {
	cu := casosConAnteriores([]domain.Sprint{sprintAnteriorDeAlfa(1, 1, 1, "Sin plan", 0, 0)}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaAnterioresDeAlfa, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "<td>—</td>")
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
