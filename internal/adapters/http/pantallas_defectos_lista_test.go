package http

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type listarDefectosFunc func(context.Context, int64) ([]domain.Defecto, error)

func (f listarDefectosFunc) Ejecutar(ctx context.Context, id int64) ([]domain.Defecto, error) {
	return f(ctx, id)
}

// actualizarEstadoDefectoFake implementa actualizadorEstadoDefecto con tres closures.
type actualizarEstadoDefectoFake struct {
	iniciar, resolver, cerrar func(context.Context, app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
}

func (f actualizarEstadoDefectoFake) Iniciar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	return f.iniciar(ctx, d)
}

func (f actualizarEstadoDefectoFake) Resolver(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	return f.resolver(ctx, d)
}

func (f actualizarEstadoDefectoFake) Cerrar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	return f.cerrar(ctx, d)
}

// llamadasEstadoDefecto guarda lo que recibieron Iniciar, Resolver y Cerrar en una prueba.
type llamadasEstadoDefecto struct {
	recibidoIniciar, recibidoResolver, recibidoCerrar app.DatosActualizarEstadoDefecto
}

// defectosDeAlfa son los cuatro defectos de ejemplo del elemento "Login" (ID 2) del proyecto alfa
// (ID 7): uno en cada estado.
func defectosDeAlfa() []domain.Defecto {
	deteccion, resolucion := int64(9), int64(9)
	return []domain.Defecto{
		{ID: 1, ProyectoID: 7, ElementoID: 2, Descripcion: "No carga el login", Severidad: domain.SeveridadAlta, Estado: domain.DefectoAbierto, SprintDeteccionID: &deteccion},
		{ID: 2, ProyectoID: 7, ElementoID: 2, Descripcion: "Tarda en cargar", Severidad: domain.SeveridadMedia, Estado: domain.DefectoEnProgreso, SprintDeteccionID: &deteccion},
		{ID: 3, ProyectoID: 7, ElementoID: 2, Descripcion: "Falla el logout", Severidad: domain.SeveridadMedia, Estado: domain.DefectoResuelto, SprintResolucionID: &resolucion},
		{ID: 4, ProyectoID: 7, ElementoID: 2, Descripcion: "Error de tipeo", Severidad: domain.SeveridadBaja, Estado: domain.DefectoCerrado, SprintResolucionID: &resolucion},
	}
}

// casosConListaDefectos suma a casosConBacklog (elemento "Login" ID 2 en alfa) los casos de uso de
// HU-21: listar devuelve defectos/errListar, iniciar/resolver/cerrar devuelven lo que indiquen las
// closures (por defecto, éxito) y guardan lo recibido en llamadas.
func casosConListaDefectos(defectos []domain.Defecto, errListar error, llamadas *llamadasEstadoDefecto) CasosDeUso {
	cu := casosConBacklog(nil, nil, &llamadasBacklog{})
	cu.ListarDefectos = listarDefectosFunc(func(_ context.Context, id int64) ([]domain.Defecto, error) {
		if id != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		return defectos, errListar
	})
	cu.ActualizarEstadoDefecto = actualizarEstadoDefectoFake{
		iniciar: func(_ context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			llamadas.recibidoIniciar = d
			return domain.Defecto{}, nil
		},
		resolver: func(_ context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			llamadas.recibidoResolver = d
			return domain.Defecto{}, nil
		},
		cerrar: func(_ context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			llamadas.recibidoCerrar = d
			return domain.Defecto{}, nil
		},
	}
	return cu
}

// Traza: HU-21 CA-21.11
func TestPantallaDefectos_ListaConSusDatosYBotonesSegunEstado(t *testing.T) {
	rec := pedirPantalla(t, casosConListaDefectos(defectosDeAlfa(), nil, &llamadasEstadoDefecto{}), http.MethodGet, "/web/proyectos/7/defectos", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec,
		"No carga el login", "Login", "Alta", "Abierto", "#9",
		"Tarda en cargar", "En progreso",
		"Falla el logout", "Resuelto",
		"Error de tipeo", "Cerrado",
		`action="/web/proyectos/7/defectos/1/iniciar"`,
		`action="/web/proyectos/7/defectos/2/resolver"`,
		`action="/web/proyectos/7/defectos/3/cerrar"`,
	)
	if strings.Contains(rec.Body.String(), `action="/web/proyectos/7/defectos/4/`) {
		t.Errorf("el defecto Cerrado no debería mostrar ningún botón; cuerpo: %.400s", rec.Body)
	}
}

// Traza: HU-21 CA-21.9
func TestPantallaDefectos_ProyectoSinDefectos_MuestraVacio(t *testing.T) {
	rec := pedirPantalla(t, casosConListaDefectos(nil, nil, &llamadasEstadoDefecto{}), http.MethodGet, "/web/proyectos/7/defectos", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "vacio")
}

func TestPantallaDefectos_ProyectoInexistente_Es404(t *testing.T) {
	rec := pedirPantalla(t, casosConListaDefectos(nil, nil, &llamadasEstadoDefecto{}), http.MethodGet, "/web/proyectos/99/defectos", nil)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-21 CA-21.11: al usar el botón, vuelve a la lista (POST-redirect-GET).
func TestPantallaIniciarDefecto_PedidoValido_RedirigeALaLista(t *testing.T) {
	llamadas := &llamadasEstadoDefecto{}

	rec := pedirPantalla(t, casosConListaDefectos(defectosDeAlfa(), nil, llamadas), http.MethodPost, "/web/proyectos/7/defectos/1/iniciar", url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/defectos" {
		t.Errorf("status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}
	if llamadas.recibidoIniciar != (app.DatosActualizarEstadoDefecto{ProyectoID: 7, DefectoID: 1}) {
		t.Errorf("datos recibidos por Iniciar: %+v", llamadas.recibidoIniciar)
	}
}

// Traza: HU-21 CA-21.11: al usar el botón, vuelve a la lista (POST-redirect-GET).
func TestPantallaResolverDefecto_PedidoValido_RedirigeALaLista(t *testing.T) {
	llamadas := &llamadasEstadoDefecto{}

	rec := pedirPantalla(t, casosConListaDefectos(defectosDeAlfa(), nil, llamadas), http.MethodPost, "/web/proyectos/7/defectos/2/resolver", url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/defectos" {
		t.Errorf("status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}
	if llamadas.recibidoResolver != (app.DatosActualizarEstadoDefecto{ProyectoID: 7, DefectoID: 2}) {
		t.Errorf("datos recibidos por Resolver: %+v", llamadas.recibidoResolver)
	}
}

// Traza: HU-21 CA-21.11: al usar el botón, vuelve a la lista (POST-redirect-GET).
func TestPantallaCerrarDefecto_PedidoValido_RedirigeALaLista(t *testing.T) {
	llamadas := &llamadasEstadoDefecto{}

	rec := pedirPantalla(t, casosConListaDefectos(defectosDeAlfa(), nil, llamadas), http.MethodPost, "/web/proyectos/7/defectos/3/cerrar", url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/defectos" {
		t.Errorf("status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}
	if llamadas.recibidoCerrar != (app.DatosActualizarEstadoDefecto{ProyectoID: 7, DefectoID: 3}) {
		t.Errorf("datos recibidos por Cerrar: %+v", llamadas.recibidoCerrar)
	}
}

// Traza: HU-21 CA-21.11: con un error, se muestra la lista con el mensaje.
func TestPantallaIniciarDefecto_DefectoNoAbierto_MuestraError(t *testing.T) {
	cu := casosConListaDefectos(defectosDeAlfa(), nil, &llamadasEstadoDefecto{})
	cu.ActualizarEstadoDefecto = actualizarEstadoDefectoFake{
		iniciar: func(context.Context, app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			return domain.Defecto{}, domain.ErrDefectoNoAbierto
		},
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/defectos/2/iniciar", url.Values{})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, domain.ErrDefectoNoAbierto.Error())
}

// Traza: HU-21 CA-21.11: con un error, se muestra la lista con el mensaje.
func TestPantallaResolverDefecto_DefectoNoEnProgreso_MuestraError(t *testing.T) {
	cu := casosConListaDefectos(defectosDeAlfa(), nil, &llamadasEstadoDefecto{})
	cu.ActualizarEstadoDefecto = actualizarEstadoDefectoFake{
		resolver: func(context.Context, app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			return domain.Defecto{}, domain.ErrDefectoNoEnProgreso
		},
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/defectos/1/resolver", url.Values{})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, domain.ErrDefectoNoEnProgreso.Error())
}

// Traza: HU-21 CA-21.11: con un error, se muestra la lista con el mensaje.
func TestPantallaCerrarDefecto_DefectoNoResuelto_MuestraError(t *testing.T) {
	cu := casosConListaDefectos(defectosDeAlfa(), nil, &llamadasEstadoDefecto{})
	cu.ActualizarEstadoDefecto = actualizarEstadoDefectoFake{
		cerrar: func(context.Context, app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			return domain.Defecto{}, domain.ErrDefectoNoResuelto
		},
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/defectos/1/cerrar", url.Values{})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, domain.ErrDefectoNoResuelto.Error())
}

func TestPantallaActualizarEstadoDefecto_DefectoInexistente_Es404(t *testing.T) {
	cu := casosConListaDefectos(defectosDeAlfa(), nil, &llamadasEstadoDefecto{})
	cu.ActualizarEstadoDefecto = actualizarEstadoDefectoFake{
		iniciar: func(context.Context, app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
			return domain.Defecto{}, domain.ErrDefectoNoEncontrado
		},
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/defectos/99/iniciar", url.Values{})

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}
