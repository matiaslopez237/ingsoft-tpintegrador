package http

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Adaptadores de funciones a las interfaces de los casos de uso de HU-05 y HU-07.
type (
	listarBacklogFunc func(context.Context, app.DatosListarBacklog) ([]domain.ElementoBacklog, error)
	crearElementoFunc func(context.Context, app.DatosNuevoElemento) (domain.ElementoBacklog, error)
)

func (f listarBacklogFunc) Ejecutar(ctx context.Context, d app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
	return f(ctx, d)
}
func (f crearElementoFunc) Ejecutar(ctx context.Context, d app.DatosNuevoElemento) (domain.ElementoBacklog, error) {
	return f(ctx, d)
}

// llamadasBacklog guarda lo que recibieron los casos de uso en una prueba.
type llamadasBacklog struct {
	listado app.DatosListarBacklog
	creado  *app.DatosNuevoElemento
}

// casosConBacklog suma a casosFalsos los casos de uso de HU-05 y HU-07: el backlog de alfa (ID 7)
// tiene dos elementos; listar devuelve errListar y crear devuelve errCrear (nil: funcionan).
func casosConBacklog(errListar, errCrear error, llamadas *llamadasBacklog) CasosDeUso {
	cinco := 5
	cu := casosFalsos()
	cu.ListarBacklog = listarBacklogFunc(func(_ context.Context, d app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
		llamadas.listado = d
		if d.ProyectoID != alfa.ID {
			return nil, domain.ErrProyectoNoEncontrado
		}
		if errListar != nil {
			return nil, errListar
		}
		return []domain.ElementoBacklog{
			{ID: 2, ProyectoID: 7, Titulo: "Login", Prioridad: domain.PrioridadAlta, Estado: domain.EstadoPendiente,
				StoryPoints: &cinco, Criterios: []string{"Con usuario y clave"}},
			{ID: 1, ProyectoID: 7, Titulo: "Reportes", Prioridad: domain.PrioridadBaja, Estado: domain.EstadoHecho,
				Criterios: []string{"Exporta a PDF"}},
		}, nil
	})
	cu.CrearElementoBacklog = crearElementoFunc(func(_ context.Context, d app.DatosNuevoElemento) (domain.ElementoBacklog, error) {
		llamadas.creado = &d
		if d.ProyectoID != alfa.ID {
			return domain.ElementoBacklog{}, domain.ErrProyectoNoEncontrado
		}
		return domain.ElementoBacklog{ID: 3, ProyectoID: d.ProyectoID}, errCrear
	})
	return cu
}

func formElemento(titulo, descripcion, prioridad, storyPoints, criterios string) url.Values {
	return url.Values{"titulo": {titulo}, "descripcion": {descripcion}, "prioridad": {prioridad},
		"story_points": {storyPoints}, "criterios": {criterios}}
}

// Traza: HU-07 CA-07.1: la pantalla lista el backlog, con los controles de orden y filtro y el
// formulario para agregar elementos.
func TestPantallaBacklog_ListaControlesYFormulario(t *testing.T) {
	var ll llamadasBacklog

	rec := pedirPantalla(t, casosConBacklog(nil, nil, &ll), http.MethodGet, "/web/proyectos/7/backlog", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	contiene(t, rec, "Alfa", "Login", "Alta", "Pendiente", "Con usuario y clave", "Reportes", "Sin estimar",
		`name="orden"`, `name="estado"`, `action="/web/proyectos/7/backlog"`,
		`name="titulo"`, `name="descripcion"`, `name="prioridad"`, `name="story_points"`, `name="criterios"`)
	if ll.listado != (app.DatosListarBacklog{ProyectoID: 7}) {
		t.Errorf("listado pedido = %+v, se esperaba el orden y el estado por defecto", ll.listado)
	}
}

// Traza: HU-07 CA-07.2, CA-07.4: el orden y el filtro de la URL llegan al caso de uso y quedan elegidos.
func TestPantallaBacklog_OrdenYFiltro(t *testing.T) {
	var ll llamadasBacklog

	rec := pedirPantalla(t, casosConBacklog(nil, nil, &ll), http.MethodGet, "/web/proyectos/7/backlog?orden=asc&estado=Hecho", nil)

	if ll.listado != (app.DatosListarBacklog{ProyectoID: 7, Orden: "asc", Estado: "Hecho"}) {
		t.Errorf("listado pedido = %+v", ll.listado)
	}
	contiene(t, rec, `value="asc" selected`, `value="Hecho" selected`)
}

// Traza: HU-07 CA-07.8: un orden inválido en la URL muestra el error.
func TestPantallaBacklog_OrdenInvalido_Devuelve400(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(domain.ErrOrdenInvalido, nil, &llamadasBacklog{}), http.MethodGet,
		"/web/proyectos/7/backlog?orden=ASC", nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400", rec.Code)
	}
	contiene(t, rec, domain.ErrOrdenInvalido.Error())
}

// Traza: HU-07 CA-07.7
func TestPantallaBacklog_ProyectoInexistente_Devuelve404(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet, "/web/proyectos/99/backlog", nil)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-07 CA-07.5
func TestPantallaBacklog_SinElementos_MuestraMensaje(t *testing.T) {
	cu := casosConBacklog(nil, nil, &llamadasBacklog{})
	cu.ListarBacklog = listarBacklogFunc(func(context.Context, app.DatosListarBacklog) ([]domain.ElementoBacklog, error) {
		return nil, nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog", nil)

	contiene(t, rec, "No hay elementos")
}

// Traza: HU-07 sección 5 (formulario de HU-05): crear un elemento vuelve al backlog; los criterios
// van uno por línea y las líneas vacías se ignoran.
func TestPantallaBacklog_CrearElemento_RedirigeAlBacklog(t *testing.T) {
	var ll llamadasBacklog

	rec := pedirPantalla(t, casosConBacklog(nil, nil, &ll), http.MethodPost, "/web/proyectos/7/backlog",
		formElemento("Perfil", "Editar mis datos", "Media", "3", "Cambia el nombre\r\n\r\n  Cambia la foto  \n"))

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/backlog" {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 al backlog", rec.Code, rec.Header().Get("Location"))
	}
	d := ll.creado
	if d == nil || d.ProyectoID != 7 || d.Titulo != "Perfil" || d.Descripcion != "Editar mis datos" ||
		d.Prioridad != domain.PrioridadMedia || d.StoryPoints == nil || *d.StoryPoints != 3 ||
		!slices.Equal(d.Criterios, []string{"Cambia el nombre", "Cambia la foto"}) {
		t.Errorf("datos recibidos por el caso de uso: %+v", d)
	}
}

// Traza: HU-05 RN-05.6: sin story points el elemento queda sin estimar.
func TestPantallaBacklog_CrearSinStoryPoints(t *testing.T) {
	var ll llamadasBacklog

	pedirPantalla(t, casosConBacklog(nil, nil, &ll), http.MethodPost, "/web/proyectos/7/backlog",
		formElemento("Perfil", "x", "Media", "", "Un criterio"))

	if ll.creado == nil || ll.creado.StoryPoints != nil {
		t.Errorf("datos recibidos = %+v, se esperaba sin story points", ll.creado)
	}
}

// Traza: HU-07 sección 5: story points que no son un número se informan sin llamar al caso de uso.
func TestPantallaBacklog_StoryPointsNoNumericos_NoLlamaAlCasoDeUso(t *testing.T) {
	var ll llamadasBacklog

	rec := pedirPantalla(t, casosConBacklog(nil, nil, &ll), http.MethodPost, "/web/proyectos/7/backlog",
		formElemento("Perfil", "x", "Media", "tres", "Un criterio"))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400", rec.Code)
	}
	contiene(t, rec, "story_points debe ser un número entero")
	if ll.creado != nil {
		t.Errorf("no se esperaba llamar al caso de uso, recibió %+v", ll.creado)
	}
}

// Traza: HU-05 sección 7 desde la pantalla: el error se muestra con lo que se había escrito.
func TestPantallaBacklog_ErrorAlCrear_MuestraElMensajeYConservaElFormulario(t *testing.T) {
	tests := []error{domain.ErrTituloDuplicado, domain.ErrTituloVacio, domain.ErrPrioridadInvalida,
		domain.ErrStoryPointsInvalidos, domain.ErrSinCriterios, domain.ErrDescripcionVacia}
	for _, errCrear := range tests {
		t.Run(errCrear.Error(), func(t *testing.T) {
			rec := pedirPantalla(t, casosConBacklog(nil, errCrear, &llamadasBacklog{}), http.MethodPost, "/web/proyectos/7/backlog",
				formElemento("Login", "Mi descripción", "Alta", "8", "Mi criterio"))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, errCrear.Error(), `value="Login"`, "Mi descripción", `value="8"`, "Mi criterio",
				"Reportes") // la lista sigue visible
		})
	}
}

// Traza: HU-07 CA-07.7
func TestPantallaBacklog_CrearEnProyectoInexistente_Devuelve404(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodPost, "/web/proyectos/99/backlog",
		formElemento("Login", "x", "Alta", "", "a"))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-07: desde el detalle del proyecto se llega a su backlog.
func TestPantallaDetalle_EnlazaAlBacklog(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet, "/web/proyectos/7", nil)

	contiene(t, rec, `href="/web/proyectos/7/backlog"`)
}

// Traza: HU-07 sección 5: el formulario de alta conserva el orden y el filtro elegidos.
func TestPantallaBacklog_FormularioConservaOrdenYFiltro(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet,
		"/web/proyectos/7/backlog?orden=asc&estado=Hecho", nil)

	contiene(t, rec, `action="/web/proyectos/7/backlog?orden=asc&amp;estado=Hecho"`)
}

// Traza: HU-07 sección 5: después de agregar un elemento se vuelve al backlog con el mismo orden y filtro.
func TestPantallaBacklog_CrearVuelveConElMismoOrdenYFiltro(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodPost,
		"/web/proyectos/7/backlog?orden=asc&estado=Hecho", formElemento("Perfil", "x", "Media", "", "Un criterio"))

	if got := rec.Header().Get("Location"); got != "/web/proyectos/7/backlog?estado=Hecho&orden=asc" {
		t.Errorf("Location = %q, se esperaba volver con el orden y el filtro", got)
	}
}

// Traza: HU-07 sección 5: si el alta falla, la lista se sigue mostrando con el orden y el filtro elegidos.
func TestPantallaBacklog_ErrorAlCrearConservaOrdenYFiltro(t *testing.T) {
	var ll llamadasBacklog

	pedirPantalla(t, casosConBacklog(nil, domain.ErrTituloDuplicado, &ll), http.MethodPost,
		"/web/proyectos/7/backlog?orden=asc&estado=Hecho", formElemento("Login", "x", "Alta", "", "Un criterio"))

	if ll.listado != (app.DatosListarBacklog{ProyectoID: 7, Orden: "asc", Estado: "Hecho"}) {
		t.Errorf("listado pedido = %+v, se esperaba con el orden y el filtro", ll.listado)
	}
}
