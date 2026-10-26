package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Adaptadores de funciones a las interfaces de los casos de uso, para armar fakes con closures.
type (
	listarFunc    func(context.Context) ([]domain.Proyecto, error)
	crearFunc     func(context.Context, app.DatosNuevoProyecto) (domain.Proyecto, error)
	consultarFunc func(context.Context, int64) (domain.Proyecto, error)
	modificarFunc func(context.Context, app.DatosModificarProyecto) (domain.Proyecto, error)
	agregarFunc   func(context.Context, app.DatosIntegrante) error
	quitarFunc    func(context.Context, app.DatosQuitarIntegrante) error
	rolFunc       func(context.Context, app.DatosRol) error
)

func (f listarFunc) Ejecutar(ctx context.Context) ([]domain.Proyecto, error) { return f(ctx) }
func (f crearFunc) Ejecutar(ctx context.Context, d app.DatosNuevoProyecto) (domain.Proyecto, error) {
	return f(ctx, d)
}
func (f consultarFunc) Ejecutar(ctx context.Context, id int64) (domain.Proyecto, error) {
	return f(ctx, id)
}
func (f modificarFunc) Ejecutar(ctx context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error) {
	return f(ctx, d)
}
func (f agregarFunc) Ejecutar(ctx context.Context, d app.DatosIntegrante) error { return f(ctx, d) }
func (f quitarFunc) Ejecutar(ctx context.Context, d app.DatosQuitarIntegrante) error {
	return f(ctx, d)
}
func (f rolFunc) Ejecutar(ctx context.Context, d app.DatosRol) error { return f(ctx, d) }

var alfa = domain.Proyecto{
	ID:          7,
	Nombre:      "Alfa",
	FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
	Estado:      domain.EstadoActivo,
	Integrantes: []domain.Integrante{{Nombre: "Ana", Roles: []string{"Tester", "Desarrolladora"}}},
}

// casosFalsos devuelve casos de uso que responden con alfa (ID 7) o con proyecto no encontrado.
func casosFalsos() CasosDeUso {
	consultar := func(_ context.Context, id int64) (domain.Proyecto, error) {
		if id != alfa.ID {
			return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
		}
		return alfa, nil
	}
	return CasosDeUso{
		ListarProyectos: listarFunc(func(context.Context) ([]domain.Proyecto, error) {
			return []domain.Proyecto{alfa}, nil
		}),
		CrearProyecto: crearFunc(func(_ context.Context, d app.DatosNuevoProyecto) (domain.Proyecto, error) {
			p, err := domain.NuevoProyecto(d.Nombre, d.FechaInicio, d.FechaFin)
			p.ID = 8
			return p, err
		}),
		ConsultarProyecto: consultarFunc(consultar),
		ModificarProyecto: modificarFunc(func(ctx context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error) {
			p, err := consultar(ctx, d.ProyectoID)
			if err != nil {
				return p, err
			}
			return p.Modificar(d.Cambios)
		}),
		AgregarIntegrante: agregarFunc(func(_ context.Context, d app.DatosIntegrante) error {
			if domain.ClaveNombre(d.Nombre) == "ana" {
				return domain.ErrIntegranteDuplicado
			}
			return nil
		}),
		QuitarIntegrante: quitarFunc(func(context.Context, app.DatosQuitarIntegrante) error { return nil }),
		AgregarRol:       rolFunc(func(context.Context, app.DatosRol) error { return domain.ErrRolDuplicado }),
	}
}

func pedirPantalla(t *testing.T, cu CasosDeUso, metodo, ruta string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarPantallas(mux, cu)
	var req *http.Request
	if form == nil {
		req = httptest.NewRequest(metodo, ruta, nil)
	} else {
		req = httptest.NewRequest(metodo, ruta, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func contiene(t *testing.T, rec *httptest.ResponseRecorder, textos ...string) {
	t.Helper()
	for _, s := range textos {
		if !strings.Contains(rec.Body.String(), s) {
			t.Errorf("la página no contiene %q; cuerpo: %.400s", s, rec.Body)
		}
	}
}

// Traza: HU-30 CA-30.10
func TestPantallas_Inicio_RedirigeALaLista(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/", nil)
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos" {
		t.Errorf("status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}
}

// Traza: HU-30 CA-30.1, CA-30.2
func TestPantallas_Lista_TieneLayoutYProyectos(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	contiene(t, rec, "<header", "<nav", `href="/static/estilos.css"`, `href="/web/proyectos/7"`, ">Alfa<", "Activo")
}

// Traza: HU-30 CA-30.2 (lista vacía)
func TestPantallas_ListaVacia(t *testing.T) {
	cu := casosFalsos()
	cu.ListarProyectos = listarFunc(func(context.Context) ([]domain.Proyecto, error) { return nil, nil })
	contiene(t, pedirPantalla(t, cu, http.MethodGet, "/web/proyectos", nil), "Todavía no hay proyectos")
}

// Traza: HU-30 CA-30.1
func TestPantallas_EstilosSeSirven(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/static/estilos.css", nil)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Header().Get("Content-Type"), "text/css") {
		t.Errorf("status = %d, Content-Type = %q", rec.Code, rec.Header().Get("Content-Type"))
	}
}

// Traza: HU-30 CA-30.3
func TestPantallas_NuevoProyecto(t *testing.T) {
	cu := casosFalsos()
	if rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/nuevo", nil); rec.Code != http.StatusOK {
		t.Errorf("formulario: status = %d", rec.Code)
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos",
		url.Values{"nombre": {"Beta"}, "fecha_inicio": {"2026-03-01"}, "fecha_fin": {""}})
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/8" {
		t.Errorf("crear válido: status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}

	rec = pedirPantalla(t, cu, http.MethodPost, "/web/proyectos",
		url.Values{"nombre": {""}, "fecha_inicio": {"2026-03-01"}})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("nombre vacío: status = %d", rec.Code)
	}
	contiene(t, rec, domain.ErrNombreVacio.Error(), `value="2026-03-01"`)

	rec = pedirPantalla(t, cu, http.MethodPost, "/web/proyectos",
		url.Values{"nombre": {"Beta"}, "fecha_inicio": {"01/03/2026"}})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("fecha inválida: status = %d", rec.Code)
	}
	contiene(t, rec, "fecha_inicio debe tener el formato AAAA-MM-DD", `value="Beta"`)
}

// Traza: HU-30 CA-30.4
func TestPantallas_Detalle(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/7", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	contiene(t, rec, ">Alfa<", "2026-03-01", "Activo", ">Ana<", "Tester", "Desarrolladora")

	rec = pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/99", nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("inexistente: status = %d", rec.Code)
	}
	contiene(t, rec, domain.ErrProyectoNoEncontrado.Error(), "<header")

	if rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/abc", nil); rec.Code != http.StatusBadRequest {
		t.Errorf("id inválido: status = %d", rec.Code)
	}
}

// Traza: HU-30 CA-30.5 (RN-30.3)
func TestPantallas_EditarProyecto(t *testing.T) {
	var recibido app.DatosModificarProyecto
	cu := casosFalsos()
	modificar := cu.ModificarProyecto
	cu.ModificarProyecto = modificarFunc(func(ctx context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error) {
		recibido = d
		return modificar.Ejecutar(ctx, d)
	})

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/editar", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("formulario: status = %d", rec.Code)
	}
	contiene(t, rec, `value="Alfa"`, `value="2026-03-01"`, `<option value="En pausa"`)

	rec = pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/editar", url.Values{
		"nombre": {"Alfa 2"}, "fecha_inicio": {"2026-03-01"}, "fecha_fin": {""}, "estado": {"En pausa"}})
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7" {
		t.Errorf("editar: status = %d, Location = %q; cuerpo %.300s", rec.Code, rec.Header().Get("Location"), rec.Body)
	}
	c := recibido.Cambios
	if c.Nombre == nil || *c.Nombre != "Alfa 2" || c.FechaInicio == nil || !c.CambiarFechaFin || c.FechaFin != nil ||
		c.Estado == nil || *c.Estado != domain.EstadoEnPausa {
		t.Errorf("cambios enviados al caso de uso: %+v", c)
	}

	rec = pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/editar", url.Values{
		"nombre": {"Alfa"}, "fecha_inicio": {"2026-03-01"}, "fecha_fin": {"2026-02-01"}, "estado": {"Activo"}})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("fin anterior: status = %d", rec.Code)
	}
	contiene(t, rec, domain.ErrFechaFinAnterior.Error(), `value="2026-02-01"`)
}

// Traza: HU-30 CA-30.6, CA-30.7, CA-30.8
func TestPantallas_Integrantes(t *testing.T) {
	cu := casosFalsos()
	tests := []struct {
		nombre   string
		ruta     string
		form     url.Values
		wantCode int
		wantMsg  string
	}{
		{"agregar", "/web/proyectos/7/integrantes", url.Values{"nombre": {"Beto"}, "rol": {"Dev"}}, http.StatusSeeOther, ""},
		{"agregar repetido", "/web/proyectos/7/integrantes", url.Values{"nombre": {"ana"}, "rol": {"Dev"}}, http.StatusBadRequest, domain.ErrIntegranteDuplicado.Error()},
		{"quitar", "/web/proyectos/7/integrantes/Ana/quitar", url.Values{}, http.StatusSeeOther, ""},
		{"rol repetido", "/web/proyectos/7/integrantes/Ana/roles", url.Values{"rol": {"Tester"}}, http.StatusBadRequest, domain.ErrRolDuplicado.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, cu, http.MethodPost, tt.ruta, tt.form)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo %.300s", rec.Code, tt.wantCode, rec.Body)
			}
			if tt.wantCode == http.StatusSeeOther && rec.Header().Get("Location") != "/web/proyectos/7" {
				t.Errorf("Location = %q", rec.Header().Get("Location"))
			}
			if tt.wantMsg != "" {
				contiene(t, rec, tt.wantMsg, ">Alfa<")
			}
		})
	}
}

// Traza: HU-30 CA-30.9 (RN-30.4)
func TestPantallas_EscapaElContenido(t *testing.T) {
	cu := casosFalsos()
	cu.ListarProyectos = listarFunc(func(context.Context) ([]domain.Proyecto, error) {
		return []domain.Proyecto{{ID: 1, Nombre: "<script>alert(1)</script>", Estado: domain.EstadoActivo}}, nil
	})
	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos", nil)
	if strings.Contains(rec.Body.String(), "<script>alert(1)</script>") {
		t.Error("el nombre se mostró sin escapar")
	}
	contiene(t, rec, "&lt;script&gt;alert(1)&lt;/script&gt;")
}

func TestPantallas_FallaInterna_Devuelve500SinDetalles(t *testing.T) {
	cu := casosFalsos()
	cu.ListarProyectos = listarFunc(func(context.Context) ([]domain.Proyecto, error) {
		return nil, errors.New("disco lleno")
	})
	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos", nil)
	if rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), "disco lleno") {
		t.Errorf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
}

// Un nombre con caracteres reservados de URL llega entero a la ruta de quitar y de sumar rol.
func TestPantallas_NombreConCaracteresReservados(t *testing.T) {
	const raro = "a/b?c#d%e"
	cu := casosFalsos()
	cu.ConsultarProyecto = consultarFunc(func(context.Context, int64) (domain.Proyecto, error) {
		p := alfa
		p.Integrantes = []domain.Integrante{{Nombre: raro, Roles: []string{"Tester"}}}
		return p, nil
	})
	var quitado string
	cu.QuitarIntegrante = quitarFunc(func(_ context.Context, d app.DatosQuitarIntegrante) error {
		quitado = d.Nombre
		return nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7", nil)
	ruta := "/web/proyectos/7/integrantes/" + url.PathEscape(raro) + "/quitar"
	contiene(t, rec, `action="`+ruta+`"`)

	pedirPantalla(t, cu, http.MethodPost, ruta, url.Values{})
	if quitado != raro {
		t.Errorf("nombre recibido = %q, se esperaba %q", quitado, raro)
	}
}

// Los formularios enviados desde otro sitio se rechazan (CSRF).
func TestPantallas_RechazaFormulariosDeOtroSitio(t *testing.T) {
	tests := []struct {
		nombre, encabezado, valor string
		want                      int
	}{
		{"Sec-Fetch-Site cross-site", "Sec-Fetch-Site", "cross-site", http.StatusForbidden},
		{"Origin de otro host", "Origin", "http://malicioso.example", http.StatusForbidden},
		{"Sec-Fetch-Site same-origin", "Sec-Fetch-Site", "same-origin", http.StatusSeeOther},
		{"Origin del mismo host", "Origin", "http://example.com", http.StatusSeeOther},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			mux := http.NewServeMux()
			RegistrarPantallas(mux, casosFalsos())
			req := httptest.NewRequest(http.MethodPost, "/web/proyectos/7/integrantes/Ana/quitar", strings.NewReader(""))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set(tt.encabezado, tt.valor)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status = %d, se esperaba %d", rec.Code, tt.want)
			}
		})
	}
}

func TestPantallas_FormularioDemasiadoGrande(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodPost, "/web/proyectos",
		url.Values{"nombre": {strings.Repeat("a", tamanoMaximoFormulario)}})
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, se esperaba 413", rec.Code)
	}
}

func TestPantallas_CabecerasDeSeguridad(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos", nil)
	for encabezado, want := range map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"Content-Security-Policy": "default-src 'self'",
	} {
		if got := rec.Header().Get(encabezado); got != want {
			t.Errorf("%s = %q, se esperaba %q", encabezado, got, want)
		}
	}
}

// Traza: DEF-01. El formulario de edición lleva los valores con los que se abrió.
func TestPantallas_EditarProyecto_FormularioLlevaLosValoresOriginales(t *testing.T) {
	rec := pedirPantalla(t, casosFalsos(), http.MethodGet, "/web/proyectos/7/editar", nil)
	contiene(t, rec,
		`name="original_nombre" value="Alfa"`,
		`name="original_fecha_inicio" value="2026-03-01"`,
		`name="original_fecha_fin" value=""`,
		`name="original_estado" value="Activo"`)
}

// Traza: DEF-01. Con los valores originales en el formulario, solo se envían al caso de uso los
// campos que el usuario cambió: así dos ediciones simultáneas de campos distintos no se pisan.
func TestPantallas_EditarProyecto_SoloEnviaLosCamposCambiados(t *testing.T) {
	var recibido app.DatosModificarProyecto
	cu := casosFalsos()
	cu.ModificarProyecto = modificarFunc(func(_ context.Context, d app.DatosModificarProyecto) (domain.Proyecto, error) {
		recibido = d
		return alfa, nil
	})

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/editar", url.Values{
		"nombre": {"Alfa"}, "fecha_inicio": {"2026-03-01"}, "fecha_fin": {"2026-12-15"}, "estado": {"Activo"},
		"original_nombre": {"Alfa"}, "original_fecha_inicio": {"2026-03-01"}, "original_fecha_fin": {""}, "original_estado": {"Activo"},
	})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	c := recibido.Cambios
	if c.Nombre != nil || c.FechaInicio != nil || c.Estado != nil {
		t.Errorf("se enviaron campos que no cambiaron: %+v", c)
	}
	if !c.CambiarFechaFin || c.FechaFin == nil || c.FechaFin.Format(time.DateOnly) != "2026-12-15" {
		t.Errorf("fecha de fin enviada: CambiarFechaFin=%v FechaFin=%v", c.CambiarFechaFin, c.FechaFin)
	}
}

// Traza: DEF-01, RN-30.3. Guardar el formulario sin cambiar nada vuelve al detalle sin mensaje y
// sin llamar al caso de uso (no es un error del usuario).
func TestPantallas_EditarProyecto_SinCambios_VuelveAlDetalleSinError(t *testing.T) {
	llamado := false
	cu := casosFalsos()
	cu.ModificarProyecto = modificarFunc(func(context.Context, app.DatosModificarProyecto) (domain.Proyecto, error) {
		llamado = true
		return domain.Proyecto{}, domain.ErrSinCambios
	})

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/editar", url.Values{
		"nombre": {"Alfa"}, "fecha_inicio": {"2026-03-01"}, "fecha_fin": {""}, "estado": {"Activo"},
		"original_nombre": {"Alfa"}, "original_fecha_inicio": {"2026-03-01"}, "original_fecha_fin": {""}, "original_estado": {"Activo"},
	})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7" {
		t.Errorf("status = %d, Location = %q; se esperaba volver al detalle", rec.Code, rec.Header().Get("Location"))
	}
	if llamado {
		t.Error("no se esperaba llamar al caso de uso si no cambió ningún campo")
	}
}
