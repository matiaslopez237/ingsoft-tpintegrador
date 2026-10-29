package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sprintCerradoDeAlfa es el sprint 2 de alfa, ya Cerrado.
var sprintCerradoDeAlfa = domain.Sprint{ID: 2, ProyectoID: 7, Numero: 2,
	FechaInicio: sprintDeAlfa.FechaInicio.AddDate(0, 0, 2), FechaFin: sprintDeAlfa.FechaFin.AddDate(0, 0, 2),
	Objetivo: "Planning Poker", Estado: domain.SprintCerrado}

// casosConCierre suma a casosConSprints el caso de uso de HU-11: alfa tiene el sprint 1 Abierto y
// el 2 Cerrado. cerrar es el doble de cerrarSprintFalso (cerrar_sprint_handler_test.go).
func casosConCierre(cerrar *cerrarSprintFalso) CasosDeUso {
	cu := casosConSprints(nil, nil)
	cu.ListarSprints = listarSprintsFunc(func(_ context.Context, id int64) ([]domain.Sprint, error) {
		if id != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		return []domain.Sprint{sprintDeAlfa, sprintCerradoDeAlfa}, nil
	})
	cu.CerrarSprint = cerrar
	return cu
}

const rutaCerrarSprintWeb = "/web/proyectos/7/sprints/1/cerrar"

// Traza: HU-11 CA-11.5: solo los sprints Abiertos muestran el botón "Cerrar sprint".
func TestPantallaSprints_BotonCerrar_SoloEnSprintsAbiertos(t *testing.T) {
	rec := pedirPantalla(t, casosConCierre(&cerrarSprintFalso{}), http.MethodGet, "/web/proyectos/7/sprints", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, `action="/web/proyectos/7/sprints/1/cerrar"`, "Cerrar sprint")
	if strings.Contains(rec.Body.String(), `action="/web/proyectos/7/sprints/2/cerrar"`) {
		t.Error("el sprint 2 está Cerrado: no tiene que mostrar el botón para cerrarlo")
	}
}

// Traza: HU-11 CA-11.5: cerrar desde la pantalla vuelve a la lista de sprints (POST/redirección/GET).
func TestPantallaSprints_Cerrar_RedirigeALaLista(t *testing.T) {
	cerrar := &cerrarSprintFalso{}

	rec := pedirPantalla(t, casosConCierre(cerrar), http.MethodPost, rutaCerrarSprintWeb, url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/sprints" {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a la pantalla de sprints", rec.Code, rec.Header().Get("Location"))
	}
	if cerrar.recibido != (app.DatosCerrarSprint{ProyectoID: 7, SprintID: 1}) {
		t.Errorf("datos recibidos por el caso de uso: %+v", cerrar.recibido)
	}
}

// cerrarSprintFunc adapta una función a la interfaz del caso de uso de HU-11.
type cerrarSprintFunc func(context.Context, app.DatosCerrarSprint) error

func (f cerrarSprintFunc) Ejecutar(ctx context.Context, d app.DatosCerrarSprint) error {
	return f(ctx, d)
}

// Traza: HU-11 CA-11.5: después de cerrar desde la pantalla, la lista muestra el sprint Cerrado y
// sin botón. El doble guarda los sprints y los cierra con domain.Sprint.Cerrar, así la lista que se
// ve tras la redirección refleja el cierre, como con la base real.
func TestPantallaSprints_Cerrar_DePuntaAPunta_QuedaCerradoYSinBoton(t *testing.T) {
	guardados := []domain.Sprint{sprintDeAlfa}
	cu := casosConSprints(nil, nil)
	cu.ListarSprints = listarSprintsFunc(func(_ context.Context, id int64) ([]domain.Sprint, error) {
		if id != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		return append([]domain.Sprint(nil), guardados...), nil
	})
	cu.CerrarSprint = cerrarSprintFunc(func(_ context.Context, d app.DatosCerrarSprint) error {
		for i := range guardados {
			if guardados[i].ID == d.SprintID && guardados[i].ProyectoID == d.ProyectoID {
				return guardados[i].Cerrar()
			}
		}
		return domain.ErrSprintNoEncontrado
	})
	mux := http.NewServeMux()
	RegistrarPantallas(mux, cu)
	pedir := func(metodo, ruta string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(metodo, ruta, strings.NewReader(""))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}
	const lista, accion = "/web/proyectos/7/sprints", `action="/web/proyectos/7/sprints/1/cerrar"`

	// Punto de partida: Abierto y con botón (si no, lo que sigue podría pasar sin probar nada).
	antes := pedir(http.MethodGet, lista)
	contiene(t, antes, "Abierto", accion)
	if strings.Contains(antes.Body.String(), "Cerrado") {
		t.Fatal("antes de cerrar, la lista ya muestra un sprint Cerrado")
	}

	post := pedir(http.MethodPost, rutaCerrarSprintWeb)
	if post.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, se esperaba 303; cuerpo: %.400s", post.Code, post.Body)
	}
	despues := pedir(http.MethodGet, post.Header().Get("Location"))

	if despues.Code != http.StatusOK {
		t.Fatalf("status = %d al seguir la redirección, se esperaba 200", despues.Code)
	}
	contiene(t, despues, "Cerrado")
	if cuerpo := despues.Body.String(); strings.Contains(cuerpo, accion) || strings.Contains(cuerpo, "Cerrar sprint") {
		t.Errorf("el sprint cerrado sigue mostrando el botón \"Cerrar sprint\"; cuerpo: %.600s", cuerpo)
	}
}

// Cerrar un sprint enviado desde otro sitio se rechaza (CSRF), sin llamar al caso de uso.
func TestPantallaSprints_Cerrar_DesdeOtroSitio_SeRechaza(t *testing.T) {
	tests := []struct{ encabezado, valor string }{
		{"Sec-Fetch-Site", "cross-site"},
		{"Origin", "http://malicioso.example"},
	}
	for _, tt := range tests {
		t.Run(tt.encabezado, func(t *testing.T) {
			cerrar := &cerrarSprintFalso{}
			mux := http.NewServeMux()
			RegistrarPantallas(mux, casosConCierre(cerrar))
			req := httptest.NewRequest(http.MethodPost, rutaCerrarSprintWeb, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set(tt.encabezado, tt.valor)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, se esperaba 403", rec.Code)
			}
			if cerrar.llamado {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}

// Traza: HU-11 CA-11.5, sección 7: los errores muestran el mensaje del dominio, aunque el caso de
// uso lo devuelva envuelto (dos pedidos que cierran a la vez, spec sección 5). Sprint ya cerrado
// es un error de validación: vuelve a la lista con 400. Lo inexistente es la página 404.
func TestPantallaSprints_Cerrar_Errores(t *testing.T) {
	tests := []struct {
		nombre     string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"ya cerrado", domain.ErrSprintYaCerrado, http.StatusBadRequest, domain.ErrSprintYaCerrado.Error()},
		{"ya cerrado, envuelto", fmt.Errorf("cerrar sprint: %w", domain.ErrSprintYaCerrado), http.StatusBadRequest, domain.ErrSprintYaCerrado.Error()},
		{"sprint inexistente", domain.ErrSprintNoEncontrado, http.StatusNotFound, domain.ErrSprintNoEncontrado.Error()},
		{"sprint inexistente, envuelto", fmt.Errorf("buscar sprint 1: %w", domain.ErrSprintNoEncontrado), http.StatusNotFound, domain.ErrSprintNoEncontrado.Error()},
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, http.StatusNotFound, domain.ErrProyectoNoEncontrado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConCierre(&cerrarSprintFalso{err: tt.err}), http.MethodPost, rutaCerrarSprintWeb, url.Values{})

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantStatus, rec.Body)
			}
			contiene(t, rec, tt.wantMsg)
			if cuerpo := rec.Body.String(); strings.Contains(cuerpo, "cerrar sprint:") || strings.Contains(cuerpo, "buscar sprint 1:") {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", cuerpo)
			}
		})
	}
}

// Traza: HU-11 CA-11.5: al fallar por sprint ya cerrado, la lista de sprints sigue visible.
func TestPantallaSprints_Cerrar_YaCerrado_SigueMostrandoLaLista(t *testing.T) {
	rec := pedirPantalla(t, casosConCierre(&cerrarSprintFalso{err: domain.ErrSprintYaCerrado}), http.MethodPost, rutaCerrarSprintWeb, url.Values{})

	contiene(t, rec, "Gestionar el backlog", "Planning Poker")
}

// Una falla interna muestra la página de error genérica, sin el detalle.
func TestPantallaSprints_Cerrar_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	cerrar := &cerrarSprintFalso{err: errors.New("cerrar sprint: database is locked")}

	rec := pedirPantalla(t, casosConCierre(cerrar), http.MethodPost, rutaCerrarSprintWeb, url.Values{})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "database") {
		t.Errorf("la página muestra el detalle interno; cuerpo: %.400s", rec.Body)
	}
}

// IDs inválidos en la ruta: página de error 400, como las otras acciones de pantalla.
func TestPantallaSprints_Cerrar_IDInvalido_Devuelve400(t *testing.T) {
	tests := []struct{ nombre, ruta, wantMsg string }{
		{"proyecto", "/web/proyectos/abc/sprints/1/cerrar", "el id del proyecto debe ser un número entero positivo"},
		{"sprint", "/web/proyectos/7/sprints/abc/cerrar", "el id del sprint debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cerrar := &cerrarSprintFalso{}

			rec := pedirPantalla(t, casosConCierre(cerrar), http.MethodPost, tt.ruta, url.Values{})

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400", rec.Code)
			}
			contiene(t, rec, tt.wantMsg)
			if cerrar.llamado {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}
