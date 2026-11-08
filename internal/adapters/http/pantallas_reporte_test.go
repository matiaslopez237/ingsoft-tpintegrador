package http

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type reporteProyectoFunc func(context.Context, int64) (domain.Reporte, error)

func (f reporteProyectoFunc) Ejecutar(ctx context.Context, id int64) (domain.Reporte, error) {
	return f(ctx, id)
}

type reporteSprintFunc func(context.Context, int64, int64) (domain.Reporte, error)

func (f reporteSprintFunc) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.Reporte, error) {
	return f(ctx, proyectoID, sprintID)
}

// casosConReporte suma a casosFalsos los casos de uso de HU-26, que devuelven el reporte indicado
// (o el error) para alfa (ID 7), y ErrProyectoNoEncontrado para cualquier otro.
func casosConReporte(reporte domain.Reporte, err error) CasosDeUso {
	cu := casosFalsos()
	cu.ArmarReporteProyecto = reporteProyectoFunc(func(_ context.Context, id int64) (domain.Reporte, error) {
		if id != alfa.ID {
			return domain.Reporte{}, domain.ErrProyectoNoEncontrado
		}
		return reporte, err
	})
	cu.ArmarReporteSprint = reporteSprintFunc(func(_ context.Context, proyectoID, sprintID int64) (domain.Reporte, error) {
		if proyectoID != alfa.ID {
			return domain.Reporte{}, domain.ErrProyectoNoEncontrado
		}
		if sprintID != 9 {
			return domain.Reporte{}, domain.ErrSprintNoEncontrado
		}
		return reporte, err
	})
	return cu
}

func reporteDeEjemplo() domain.Reporte {
	sp5, sp8 := 5, 8
	uno := int64(1)
	return domain.Reporte{
		Proyecto: alfa,
		Historias: []domain.ElementoBacklog{
			{Titulo: "Login", Prioridad: domain.PrioridadAlta, Estado: domain.EstadoHecho, StoryPoints: &sp5},
			{Titulo: "Logout", Prioridad: domain.PrioridadMedia, Estado: domain.EstadoPendiente, StoryPoints: &sp8},
		},
		HistoriasCompletadas: 1, StoryPointsCompletados: 5, StoryPointsTotal: 13,
		Esfuerzo:     domain.TotalHoras{HorasEstimadas: 10, HorasReales: 8},
		HayVelocidad: true, Velocidad: 6,
		Metricas: []domain.FilaMetrica{
			{Sprint: domain.Sprint{Numero: 1}, HayPorcentaje: true, PorcentajeCompletado: 75, DefectosDetectados: 2, DefectosResueltos: 1},
		},
		Defectos: []domain.Defecto{
			{Descripcion: "Falla el login", Severidad: domain.SeveridadAlta, Estado: domain.DefectoAbierto, SprintDeteccionID: &uno},
		},
	}
}

// Traza: HU-26 CA-26.1, CA-26.3, CA-26.4, CA-26.5, CA-26.9
func TestPantallaReporteProyecto_MuestraLasCincoSecciones(t *testing.T) {
	cu := casosConReporte(reporteDeEjemplo(), nil)

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/reporte", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Login", "Logout", "1 de 2 historias completadas", "5 de 13 story points",
		"Velocidad del equipo", "6", "75", "2", "1", "Falla el login")
}

// Traza: HU-26 CA-26.2
func TestPantallaReporteSprint_Muestra(t *testing.T) {
	cu := casosConReporte(reporteDeEjemplo(), nil)

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/sprints/9/reporte", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Login")
}

// Traza: HU-26 CA-26.15
func TestPantallaReporteProyecto_SinDatos_MuestraMensajes(t *testing.T) {
	cu := casosConReporte(domain.Reporte{Proyecto: alfa}, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/reporte", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "no tiene historias", "no tiene sprints cerrados", "no tiene defectos")
}

// Traza: HU-26 sección 7
func TestPantallaReporteProyecto_Errores(t *testing.T) {
	tests := []struct {
		nombre string
		ruta   string
		err    error
		want   int
		msg    string
	}{
		{"proyecto inexistente", "/web/proyectos/99/reporte", nil, http.StatusNotFound, "proyecto no encontrado"},
		{"id inválido", "/web/proyectos/abc/reporte", nil, http.StatusBadRequest, "número"},
		{"falla interna", "/web/proyectos/7/reporte", errors.New("disco lleno"), http.StatusInternalServerError, "error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConReporte(domain.Reporte{}, tt.err), http.MethodGet, tt.ruta, nil)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.300s", rec.Code, tt.want, rec.Body)
			}
			contiene(t, rec, tt.msg)
			noContiene(t, rec.Body.String(), "disco lleno")
		})
	}
}

// Traza: HU-26 sección 7
func TestPantallaReporteSprint_Errores(t *testing.T) {
	tests := []struct {
		nombre string
		ruta   string
		err    error
		want   int
		msg    string
	}{
		{"proyecto inexistente", "/web/proyectos/99/sprints/9/reporte", nil, http.StatusNotFound, "proyecto no encontrado"},
		{"sprint inexistente", "/web/proyectos/7/sprints/999/reporte", nil, http.StatusNotFound, "sprint no encontrado"},
		{"id de sprint inválido", "/web/proyectos/7/sprints/abc/reporte", nil, http.StatusBadRequest, "número"},
		{"falla interna", "/web/proyectos/7/sprints/9/reporte", errors.New("disco lleno"), http.StatusInternalServerError, "error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConReporte(domain.Reporte{}, tt.err), http.MethodGet, tt.ruta, nil)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.300s", rec.Code, tt.want, rec.Body)
			}
			contiene(t, rec, tt.msg)
			noContiene(t, rec.Body.String(), "disco lleno")
		})
	}
}

// Traza: HU-26 sección 5 (P-3): enlace desde el detalle del proyecto.
func TestPantallaProyectoDetalle_EnlazaElReporte(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/7", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	contiene(t, rec, `href="/web/proyectos/7/reporte"`)
}
