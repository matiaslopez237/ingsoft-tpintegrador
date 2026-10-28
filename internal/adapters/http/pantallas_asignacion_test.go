package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

var idSprintDeAlfa = sprintDeAlfa.ID

var elementoAsignadoAlSprint = domain.ElementoBacklog{
	ID: 10, ProyectoID: alfa.ID, Titulo: "Login de usuarios", Prioridad: domain.PrioridadAlta,
	Estado: domain.EstadoPendiente, Criterios: []string{"Criterio"}, SprintID: &idSprintDeAlfa,
}

var elementoSinSprint = domain.ElementoBacklog{
	ID: 11, ProyectoID: alfa.ID, Titulo: "Logout", Prioridad: domain.PrioridadMedia,
	Estado: domain.EstadoPendiente, Criterios: []string{"Criterio"},
}

// casosConAsignacion suma a casosConSprints los casos de uso que necesita la pantalla de detalle
// del sprint: listar el backlog (con el elemento ya asignado y uno disponible) y asignar/desasignar.
func casosConAsignacion(errAsignar, errDesasignar error) CasosDeUso {
	cu := casosConSprints(nil, nil)
	cu.ListarBacklog = listarBacklogFunc(func(_ context.Context, d app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
		if d.ProyectoID != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		return []domain.ElementoBacklog{elementoAsignadoAlSprint, elementoSinSprint}, nil
	})
	cu.AsignarHistoria = &asignarHistoriaFalso{err: errAsignar}
	cu.DesasignarHistoria = &desasignarHistoriaFalso{err: errDesasignar}
	return cu
}

const rutaSprintDetalle = "/web/proyectos/7/sprints/1"

// Traza: HU-09, Tarea 9.8: la pantalla muestra el sprint, lo asignado y lo disponible para asignar.
func TestPantallaSprintDetalle_MuestraAsignadosYDisponibles(t *testing.T) {
	rec := pedirPantalla(t, casosConAsignacion(nil, nil), http.MethodGet, rutaSprintDetalle, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Login de usuarios", "Logout", "Gestionar el backlog",
		`action="/web/proyectos/7/sprints/1/asignar"`, `action="/web/proyectos/7/sprints/1/quitar"`)
}

// Traza: HU-09, Tarea 9.8
func TestPantallaSprintDetalle_SprintInexistente_Devuelve404(t *testing.T) {
	rec := pedirPantalla(t, casosConAsignacion(nil, nil), http.MethodGet, "/web/proyectos/7/sprints/99", nil)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-09 CA-09.1
func TestPantallaSprintDetalle_Asignar_RedirigeAlDetalle(t *testing.T) {
	uc := &asignarHistoriaFalso{}
	cu := casosConAsignacion(nil, nil)
	cu.AsignarHistoria = uc

	rec := pedirPantalla(t, cu, http.MethodPost, rutaSprintDetalle+"/asignar", url.Values{"elemento_id": {"11"}})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaSprintDetalle {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 al detalle del sprint", rec.Code, rec.Header().Get("Location"))
	}
	if uc.recibido.ProyectoID != 7 || uc.recibido.ElementoID != 11 || uc.recibido.SprintID != 1 {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-09 CA-09.2, CA-09.7, CA-09.8
func TestPantallaSprintDetalle_ErrorAlAsignar_MuestraElMensaje(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
	}{
		{"elemento ya asignado", domain.ErrElementoYaAsignado},
		{"sprint cerrado", domain.ErrSprintCerrado},
		{"sprint de otro proyecto", domain.ErrSprintDeOtroProyecto},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConAsignacion(tt.err, nil), http.MethodPost, rutaSprintDetalle+"/asignar",
				url.Values{"elemento_id": {"11"}})

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, tt.err.Error())
		})
	}
}

// Traza: HU-09 CA-09.4
func TestPantallaSprintDetalle_Quitar_RedirigeAlDetalle(t *testing.T) {
	uc := &desasignarHistoriaFalso{}
	cu := casosConAsignacion(nil, nil)
	cu.DesasignarHistoria = uc

	rec := pedirPantalla(t, cu, http.MethodPost, rutaSprintDetalle+"/quitar", url.Values{"elemento_id": {"10"}})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaSprintDetalle {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 al detalle del sprint", rec.Code, rec.Header().Get("Location"))
	}
	if uc.recibido.ProyectoID != 7 || uc.recibido.ElementoID != 10 {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-09 CA-09.5, CA-09.6
func TestPantallaSprintDetalle_ErrorAlQuitar_MuestraElMensaje(t *testing.T) {
	rec := pedirPantalla(t, casosConAsignacion(nil, domain.ErrSprintCerrado), http.MethodPost, rutaSprintDetalle+"/quitar",
		url.Values{"elemento_id": {"10"}})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, domain.ErrSprintCerrado.Error())
}
