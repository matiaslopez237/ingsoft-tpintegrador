package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Adaptadores de funciones a las interfaces de los casos de uso de HU-08, como en pantallas_test.go.
type (
	crearSprintFunc   func(context.Context, app.DatosNuevoSprint) (domain.Sprint, error)
	listarSprintsFunc func(context.Context, int64) ([]domain.Sprint, error)
)

func (f crearSprintFunc) Ejecutar(ctx context.Context, d app.DatosNuevoSprint) (domain.Sprint, error) {
	return f(ctx, d)
}
func (f listarSprintsFunc) Ejecutar(ctx context.Context, id int64) ([]domain.Sprint, error) {
	return f(ctx, id)
}

var sprintDeAlfa = domain.Sprint{ID: 1, ProyectoID: 7, Numero: 1,
	FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC),
	Objetivo: "Gestionar el backlog", Estado: domain.SprintAbierto}

// casosConSprints suma a casosFalsos los casos de uso de HU-08: alfa (ID 7) tiene un sprint, y
// crear devuelve el error indicado (nil: lo crea). recibido guarda lo que llegó a CrearSprint.
func casosConSprints(errCrear error, recibido *app.DatosNuevoSprint) CasosDeUso {
	cu := casosFalsos()
	cu.ListarSprints = listarSprintsFunc(func(_ context.Context, id int64) ([]domain.Sprint, error) {
		if id != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		return []domain.Sprint{sprintDeAlfa}, nil
	})
	cu.CrearSprint = crearSprintFunc(func(_ context.Context, d app.DatosNuevoSprint) (domain.Sprint, error) {
		if recibido != nil {
			*recibido = d
		}
		if d.ProyectoID != alfa.ID {
			return domain.Sprint{}, domain.ErrProyectoNoEncontrado
		}
		return domain.Sprint{ID: 2, ProyectoID: d.ProyectoID, Numero: 2}, errCrear
	})
	return cu
}

func formSprint(inicio, fin, objetivo string) url.Values {
	return url.Values{"fecha_inicio": {inicio}, "fecha_fin": {fin}, "objetivo": {objetivo}}
}

// Traza: HU-08 CA-08.9: la pantalla lista los sprints del proyecto y ofrece el formulario.
func TestPantallaSprints_ListaYFormulario(t *testing.T) {
	rec := pedirPantalla(t, casosConSprints(nil, nil), http.MethodGet, "/web/proyectos/7/sprints", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Alfa", "2026-10-26", "2026-10-27", "Gestionar el backlog", "Abierto",
		`action="/web/proyectos/7/sprints"`, `name="fecha_inicio"`, `name="fecha_fin"`, `name="objetivo"`)
}

// Traza: HU-08 CA-08.9: un proyecto sin sprints muestra un mensaje en lugar de una tabla vacía.
func TestPantallaSprints_SinSprints_MuestraMensaje(t *testing.T) {
	cu := casosConSprints(nil, nil)
	cu.ListarSprints = listarSprintsFunc(func(context.Context, int64) ([]domain.Sprint, error) { return nil, nil })

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/sprints", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
	}
	contiene(t, rec, "todavía no tiene sprints")
}

// Traza: HU-08 CA-08.3
func TestPantallaSprints_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := pedirPantalla(t, casosConSprints(nil, nil), http.MethodGet, "/web/proyectos/99/sprints", nil)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-08 CA-08.1: crear un sprint vuelve a la pantalla de sprints (POST/redirección/GET).
func TestPantallaSprints_Crear_RedirigeALaLista(t *testing.T) {
	var recibido app.DatosNuevoSprint

	rec := pedirPantalla(t, casosConSprints(nil, &recibido), http.MethodPost, "/web/proyectos/7/sprints",
		formSprint("2026-10-28", "2026-10-29", "Planning Poker"))

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/sprints" {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a la pantalla de sprints", rec.Code, rec.Header().Get("Location"))
	}
	if recibido.ProyectoID != 7 || recibido.FechaInicio.Format(time.DateOnly) != "2026-10-28" ||
		recibido.FechaFin.Format(time.DateOnly) != "2026-10-29" || recibido.Objetivo != "Planning Poker" {
		t.Errorf("datos recibidos por el caso de uso: %+v", recibido)
	}
}

// Traza: HU-08 CA-08.4 a CA-08.8: los errores de validación se muestran en la misma pantalla, con
// lo que se había escrito en el formulario.
func TestPantallaSprints_ErrorDeValidacion_MuestraElMensajeYConservaElFormulario(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
	}{
		{"superpuesto", domain.ErrSprintSuperpuesto},
		{"fuera del proyecto", domain.ErrSprintFueraDelProyecto},
		{"Sprint Goal vacío", domain.ErrObjetivoVacio},
		{"proyecto no activo", domain.ErrProyectoNoActivo},
		{"sin fecha de fin", domain.ErrFechaFinObligatoria},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConSprints(tt.err, nil), http.MethodPost, "/web/proyectos/7/sprints",
				formSprint("2026-10-27", "2026-10-28", "Mi objetivo"))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, tt.err.Error(), `value="2026-10-27"`, `value="2026-10-28"`, `value="Mi objetivo"`,
				"Gestionar el backlog") // la lista de sprints sigue visible
		})
	}
}

// Traza: HU-08 sección 7: una fecha con formato inválido se informa sin llamar al caso de uso.
func TestPantallaSprints_FechaInvalida_NoLlamaAlCasoDeUso(t *testing.T) {
	var recibido app.DatosNuevoSprint

	rec := pedirPantalla(t, casosConSprints(nil, &recibido), http.MethodPost, "/web/proyectos/7/sprints",
		formSprint("26/10/2026", "2026-10-28", "x"))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400", rec.Code)
	}
	contiene(t, rec, "fecha_inicio debe tener el formato AAAA-MM-DD")
	if recibido.ProyectoID != 0 {
		t.Errorf("no se esperaba llamar al caso de uso, recibió %+v", recibido)
	}
}

// Traza: HU-08 CA-08.3
func TestPantallaSprints_CrearEnProyectoInexistente_Devuelve404(t *testing.T) {
	rec := pedirPantalla(t, casosConSprints(nil, nil), http.MethodPost, "/web/proyectos/99/sprints",
		formSprint("2026-10-27", "2026-10-28", "x"))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-08: desde el detalle del proyecto se llega a sus sprints.
func TestPantallaDetalle_EnlazaALosSprints(t *testing.T) {
	rec := pedirPantalla(t, casosConSprints(nil, nil), http.MethodGet, "/web/proyectos/7", nil)

	contiene(t, rec, `href="/web/proyectos/7/sprints"`)
}
