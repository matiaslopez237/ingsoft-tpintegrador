package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type completarFunc func(context.Context, app.DatosCompletar) error

func (f completarFunc) Ejecutar(ctx context.Context, d app.DatosCompletar) error { return f(ctx, d) }

const (
	accionCompletar       = `action="/web/proyectos/7/sprints/1/completar"`
	botonCompletarDeLogin = `aria-label="Marcar Login de usuarios como hecha"`
)

// casosConCompletar arma la pantalla de detalle del sprint con el elemento asignado en el estado
// indicado, el sprint en el estado indicado y CompletarHistoria devolviendo errCompletar.
func casosConCompletar(estado domain.EstadoElemento, estadoSprint domain.EstadoSprint, errCompletar error) (CasosDeUso, *completarHistoriaFalso) {
	cu := casosConAsignacion(nil, nil)
	sprint := sprintDeAlfa
	sprint.Estado = estadoSprint
	cu.ListarSprints = listarSprintsFunc(func(context.Context, int64) ([]domain.Sprint, error) {
		return []domain.Sprint{sprint}, nil
	})
	elemento := elementoAsignadoAlSprint
	elemento.Estado = estado
	cu.ListarBacklog = listarBacklogFunc(func(context.Context, app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
		return []domain.ElementoBacklog{elemento, elementoSinSprint}, nil
	})
	uc := &completarHistoriaFalso{err: errCompletar}
	cu.CompletarHistoria = uc
	return cu, uc
}

func noContiene(t *testing.T, cuerpo, texto string) {
	t.Helper()
	if strings.Contains(cuerpo, texto) {
		t.Errorf("la página no debería contener %q; cuerpo: %.600s", texto, cuerpo)
	}
}

// Traza: HU-10 CA-10.7: con el sprint Abierto, lo que no está Hecho tiene el botón y se ve su estado.
func TestPantallaSprintDetalle_BotonCompletar_SprintAbiertoYNoHecho(t *testing.T) {
	for _, estado := range []domain.EstadoElemento{domain.EstadoPendiente, domain.EstadoEnProgreso} {
		t.Run(string(estado), func(t *testing.T) {
			cu, _ := casosConCompletar(estado, domain.SprintAbierto, nil)

			rec := pedirPantalla(t, cu, http.MethodGet, rutaSprintDetalle, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, "Estado", string(estado), accionCompletar, botonCompletarDeLogin, "Marcar como hecha")
		})
	}
}

// Traza: HU-10 CA-10.7: sin botón si el elemento ya está Hecho o si el sprint está Cerrado.
func TestPantallaSprintDetalle_SinBotonCompletar(t *testing.T) {
	tests := []struct {
		nombre       string
		estado       domain.EstadoElemento
		estadoSprint domain.EstadoSprint
	}{
		{"elemento Hecho con sprint Abierto", domain.EstadoHecho, domain.SprintAbierto},
		{"elemento Pendiente con sprint Cerrado", domain.EstadoPendiente, domain.SprintCerrado},
		{"elemento En progreso con sprint Cerrado", domain.EstadoEnProgreso, domain.SprintCerrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu, _ := casosConCompletar(tt.estado, tt.estadoSprint, nil)

			rec := pedirPantalla(t, cu, http.MethodGet, rutaSprintDetalle, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, string(tt.estado))
			noContiene(t, rec.Body.String(), accionCompletar)
		})
	}
}

// Traza: HU-10 CA-10.7 (CA-10.1 por pantalla)
func TestPantallaSprintDetalle_Completar_RedirigeAlDetalle(t *testing.T) {
	cu, uc := casosConCompletar(domain.EstadoPendiente, domain.SprintAbierto, nil)

	rec := pedirPantalla(t, cu, http.MethodPost, rutaSprintDetalle+"/completar", url.Values{"elemento_id": {"10"}})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaSprintDetalle {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 al detalle del sprint", rec.Code, rec.Header().Get("Location"))
	}
	if uc.recibido != (app.DatosCompletar{ProyectoID: 7, ElementoID: 10}) {
		t.Errorf("datos recibidos por el caso de uso: %+v", uc.recibido)
	}
}

// Traza: HU-10 CA-10.7: al volver al detalle después de completar, el elemento aparece Hecho y sin botón.
func TestPantallaSprintDetalle_Completar_ElElementoApareceHecho(t *testing.T) {
	cu, _ := casosConCompletar(domain.EstadoPendiente, domain.SprintAbierto, nil)
	elementos := []domain.ElementoBacklog{elementoAsignadoAlSprint, elementoSinSprint}
	cu.ListarBacklog = listarBacklogFunc(func(context.Context, app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
		return elementos, nil
	})
	cu.CompletarHistoria = completarFunc(func(_ context.Context, d app.DatosCompletar) error {
		for i := range elementos {
			if elementos[i].ID == d.ElementoID {
				elementos[i].Estado = domain.EstadoHecho
			}
		}
		return nil
	})

	post := pedirPantalla(t, cu, http.MethodPost, rutaSprintDetalle+"/completar", url.Values{"elemento_id": {"10"}})
	if post.Code != http.StatusSeeOther {
		t.Fatalf("status del POST = %d, se esperaba 303", post.Code)
	}
	rec := pedirPantalla(t, cu, http.MethodGet, post.Header().Get("Location"), nil)

	contiene(t, rec, "Login de usuarios", string(domain.EstadoHecho))
	noContiene(t, rec.Body.String(), accionCompletar)
}

// Traza: HU-10 CA-10.7 (errores de CA-10.3 a CA-10.5 por pantalla)
func TestPantallaSprintDetalle_ErrorAlCompletar(t *testing.T) {
	tests := []struct {
		nombre      string
		err         error
		form        url.Values
		wantStatus  int
		wantMensaje string
	}{
		{"elemento sin sprint", domain.ErrElementoSinSprint, url.Values{"elemento_id": {"10"}},
			http.StatusBadRequest, domain.ErrElementoSinSprint.Error()},
		{"sprint cerrado", fmt.Errorf("completar: %w", domain.ErrSprintCerrado), url.Values{"elemento_id": {"10"}},
			http.StatusBadRequest, domain.ErrSprintCerrado.Error()},
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, url.Values{"elemento_id": {"10"}},
			http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
		{"elemento inexistente", domain.ErrElementoNoEncontrado, url.Values{"elemento_id": {"10"}},
			http.StatusNotFound, domain.ErrElementoNoEncontrado.Error()},
		{"elemento_id inválido", nil, url.Values{"elemento_id": {"abc"}},
			http.StatusBadRequest, "elegí un elemento del backlog"},
		{"falla interna", errors.New("disco lleno"), url.Values{"elemento_id": {"10"}},
			http.StatusInternalServerError, "Ocurrió un error interno"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu, _ := casosConCompletar(domain.EstadoPendiente, domain.SprintAbierto, tt.err)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaSprintDetalle+"/completar", tt.form)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantStatus, rec.Body)
			}
			contiene(t, rec, tt.wantMensaje)
			noContiene(t, rec.Body.String(), "disco lleno")
		})
	}
}

// Traza: HU-10 sección 7: las pantallas reconocen los errores de validación de completar.
func TestMensajeDeValidacion_ErroresDeCompletar(t *testing.T) {
	for _, e := range []error{domain.ErrElementoSinSprint, domain.ErrSprintCerrado} {
		mensaje, ok := mensajeDeValidacion(fmt.Errorf("completar historia: %w", e))
		if !ok || mensaje != e.Error() {
			t.Errorf("mensajeDeValidacion(%v) = %q, %v; se esperaba %q, true", e, mensaje, ok, e.Error())
		}
	}
}

// Traza: DEF-03, HU-09 sección 3 (RN-09.6, CA-09.5, CA-09.8) y HU-10: con el sprint Cerrado se ven
// sus elementos y sus estados, pero ninguna acción (antes se mostraban y el sistema las rechazaba),
// y un aviso de que el sprint está cerrado.
func TestPantallaSprintDetalle_SprintCerrado_SinAcciones(t *testing.T) {
	cu, _ := casosConCompletar(domain.EstadoPendiente, domain.SprintCerrado, nil)

	rec := pedirPantalla(t, cu, http.MethodGet, rutaSprintDetalle, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, elementoAsignadoAlSprint.Titulo, "Pendiente", "El sprint está cerrado")
	cuerpo := rec.Body.String()
	for _, accion := range []string{"/sprints/1/quitar", "/sprints/1/asignar", "/sprints/1/completar", "Sumar un elemento"} {
		noContiene(t, cuerpo, accion)
	}
}
