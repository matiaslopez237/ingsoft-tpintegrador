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
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const (
	rutaElementoWeb = "/web/proyectos/7/backlog/7"
	rutaEsfuerzoWeb = "/web/proyectos/7/backlog/7/esfuerzo"
	envolturaCaso   = "registrar esfuerzo:"
)

// hoyPantalla es casi medianoche en una zona con otro día en UTC (ya es 2 de noviembre): la fecha
// inicial tiene que ser la de la zona del reloj, no la de UTC (RN-18.4).
var hoyPantalla = time.Date(2026, time.November, 1, 23, 30, 0, 0, time.FixedZone("ART", -3*60*60))

const valorHoyPantalla = `name="fecha" type="date" value="2026-11-01"`

// registradorEnvuelto envuelve con %w los errores del doble, como hace el caso de uso real.
type registradorEnvuelto struct{ *registradorEsfuerzoFalso }

func (f registradorEnvuelto) Ejecutar(ctx context.Context, d app.DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error) {
	r, err := f.registradorEsfuerzoFalso.Ejecutar(ctx, d)
	if err != nil {
		return r, fmt.Errorf("%s %w", envolturaCaso, err)
	}
	return r, nil
}

// casosConEsfuerzo suma a casosBacklogFalsos el caso de uso de HU-18: alfa tiene a Ana y a Caro, y
// el elemento Login (7).
func casosConEsfuerzo(uc registradorEsfuerzo) CasosDeUso {
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})
	cu.RegistrarEsfuerzo = uc
	cu.ConsultarProyecto = consultarFunc(func(_ context.Context, id int64) (domain.Proyecto, error) {
		if id != alfa.ID {
			return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
		}
		p := alfa
		p.Integrantes = []domain.Integrante{{Nombre: "Ana"}, {Nombre: "Caro"}}
		return p, nil
	})
	return cu
}

// pedirConReloj arma las pantallas del elemento con el reloj fijo en hoyPantalla.
func pedirConReloj(t *testing.T, cu CasosDeUso, metodo, ruta string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	p := nuevasPantallas(cu)
	p.reloj = func() time.Time { return hoyPantalla }
	mux := http.NewServeMux()
	p.registrarPantallasElemento(mux)
	return pedirAlMux(t, mux, metodo, ruta, form)
}

func formEsfuerzoValido() url.Values {
	return url.Values{"integrante": {"Ana"}, "fecha": {"2026-10-20"}, "actividad": {"Desarrollo"}, "horas": {"1.5"}}
}

// seccionEsfuerzo devuelve el HTML de la sección "Registrar esfuerzo" de la pantalla.
func seccionEsfuerzo(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	_, resto, ok := strings.Cut(rec.Body.String(), `id="esfuerzo"`)
	if !ok {
		t.Fatalf("la página no tiene la sección de esfuerzo; status = %d, cuerpo: %.400s", rec.Code, rec.Body)
	}
	seccion, _, ok := strings.Cut(resto, "</section>")
	if !ok {
		t.Fatalf("la sección de esfuerzo no cierra; cuerpo: %.400s", resto)
	}
	return seccion
}

func seccionContiene(t *testing.T, seccion string, textos ...string) {
	t.Helper()
	for _, s := range textos {
		if !strings.Contains(seccion, s) {
			t.Errorf("la sección no contiene %q; sección: %s", s, seccion)
		}
	}
}

func seccionNoContiene(t *testing.T, seccion string, textos ...string) {
	t.Helper()
	for _, s := range textos {
		if strings.Contains(strings.ToLower(seccion), strings.ToLower(s)) {
			t.Errorf("la sección no debería contener %q; sección: %s", s, seccion)
		}
	}
}

// esPaginaDeError verifica que se respondió con error.html y no con la pantalla del elemento.
func esPaginaDeError(t *testing.T, rec *httptest.ResponseRecorder, mensaje string) {
	t.Helper()
	contiene(t, rec, "No se pudo completar el pedido", mensaje)
	if strings.Contains(rec.Body.String(), `name="integrante"`) {
		t.Errorf("se esperaba la página de error, no el formulario; cuerpo: %.400s", rec.Body)
	}
}

// Traza: HU-18 CA-18.16: formulario inicial, con la fecha de hoy en la zona del reloj.
func TestPantallaEsfuerzo_FormularioInicial(t *testing.T) {
	rec := pedirConReloj(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodGet, rutaElementoWeb, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	s := seccionEsfuerzo(t, rec)
	seccionContiene(t, s, "Registrar esfuerzo", `action="/web/proyectos/7/backlog/7/esfuerzo#esfuerzo"`,
		`<option value="">Elegí un integrante</option>`, `<option value="Ana">Ana</option>`, `<option value="Caro">Caro</option>`,
		valorHoyPantalla, `name="actividad" type="text" value=""`, `name="horas" type="number" step="0.5" value=""`)
	seccionNoContiene(t, s, " selected", "<table", "total", "registros", "Esfuerzo registrado.", "no tiene integrantes")
	// Los límites quedan solo en el dominio (P-12).
	seccionNoContiene(t, s, ` min=`, ` max=`)
}

// Traza: HU-18 CA-18.17
func TestPantallaBacklog_EnlazaARegistrarEsfuerzoEnCadaElemento(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet, "/web/proyectos/7/backlog", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec,
		`href="/web/proyectos/7/backlog/2#esfuerzo"`, `aria-label="Registrar esfuerzo en Login"`,
		`href="/web/proyectos/7/backlog/1#esfuerzo"`, `aria-label="Registrar esfuerzo en Reportes"`)
}

// Traza: HU-18 CA-18.18: POST-redirect-GET; recargar el GET no registra otro.
func TestPantallaEsfuerzo_EnvioValido_RedirigeYConfirma(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}
	cu := casosConEsfuerzo(uc)

	rec := pedirPantalla(t, cu, http.MethodPost, rutaEsfuerzoWeb, formEsfuerzoValido())

	const destino = rutaElementoWeb + "?esfuerzo=registrado#esfuerzo"
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != destino {
		t.Fatalf("status = %d, Location = %q; se esperaba 303 a %q; cuerpo: %.300s", rec.Code, rec.Header().Get("Location"), destino, rec.Body)
	}
	if uc.llamadas != 1 {
		t.Errorf("llamadas al caso de uso = %d, se esperaba 1", uc.llamadas)
	}
	d := uc.recibido
	if d.ProyectoID != 7 || d.ElementoID != 7 || d.Integrante != "Ana" || d.Actividad != "Desarrollo" || d.Horas != 1.5 {
		t.Errorf("datos recibidos por el caso de uso: %+v", d)
	}

	// El navegador no envía el fragmento al servidor.
	siguiente, _, _ := strings.Cut(destino, "#")
	for range 2 { // el segundo pedido es la recarga
		get := pedirConReloj(t, cu, http.MethodGet, siguiente, nil)
		if get.Code != http.StatusOK {
			t.Fatalf("GET: status = %d", get.Code)
		}
		s := seccionEsfuerzo(t, get)
		seccionContiene(t, s, `class="exito" role="status"`, "Esfuerzo registrado.", valorHoyPantalla, `name="actividad" type="text" value=""`)
		seccionNoContiene(t, s, " selected")
	}
	if uc.llamadas != 1 {
		t.Errorf("después de recargar: llamadas al caso de uso = %d, se esperaba 1", uc.llamadas)
	}
}

// Traza: HU-18 CA-18.18: el aviso sale solo con el valor que usa la redirección.
func TestPantallaEsfuerzo_OtroValorDelParametro_NoConfirma(t *testing.T) {
	rec := pedirConReloj(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodGet, rutaElementoWeb+"?esfuerzo=otro", nil)

	seccionNoContiene(t, seccionEsfuerzo(t, rec), "Esfuerzo registrado.")
}

// Traza: HU-18 RN-18.4, sección 5: la fecha se interpreta como un día de la zona del servidor. Con
// fechaDeForm llegaría en UTC y, en una zona como la de Argentina, sería otro instante.
func TestPantallaEsfuerzo_FechaEnLaZonaDelServidor(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}

	pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, rutaEsfuerzoWeb, formEsfuerzoValido())

	want := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.Local)
	if f := uc.recibido.Fecha; !f.Equal(want) || f.Location() != time.Local {
		t.Errorf("fecha recibida = %v (zona %v), se esperaba %v en time.Local", f, f.Location(), want)
	}
}

// Traza: HU-18 sección 9: las horas se leen sin los espacios de alrededor.
func TestPantallaEsfuerzo_HorasConEspacios_SeLeen(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}
	form := formEsfuerzoValido()
	form.Set("horas", " 2 ")

	rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, rutaEsfuerzoWeb, form)

	if rec.Code != http.StatusSeeOther || uc.recibido.Horas != 2 {
		t.Errorf("status = %d, horas recibidas = %v", rec.Code, uc.recibido.Horas)
	}
}

// Traza: HU-18 CA-18.19 (RN-18.2, RN-18.4 a RN-18.8): errores del dominio con el caso de uso
// devolviéndolos envueltos: mensajeDeValidacion no conoce los de HU-18 y respondería 500.
func TestPantallaEsfuerzo_ErroresDelDominio_VuelvenAlFormulario(t *testing.T) {
	tests := []struct {
		nombre, campo, valor string
		want                 error
		conservado           string
	}{
		{"integrante vacío", "integrante", "", domain.ErrNombreIntegranteVacio, `name="fecha" type="date" value="2026-10-20"`},
		{"fecha vacía", "fecha", "", domain.ErrFechaEsfuerzoObligatoria, `name="fecha" type="date" value=""`},
		{"fecha futura", "fecha", "2026-10-26", domain.ErrFechaEsfuerzoFutura, `name="fecha" type="date" value="2026-10-26"`},
		{"actividad vacía", "actividad", "   ", domain.ErrActividadVacia, `name="actividad" type="text" value="   "`},
		{"horas 0", "horas", "0", domain.ErrHorasInvalidas, `name="horas" type="number" step="0.5" value="0"`},
		{"horas 1,25", "horas", "1.25", domain.ErrHorasNoMultiploMediaHora, `name="horas" type="number" step="0.5" value="1.25"`},
		{"horas 24,5", "horas", "24.5", domain.ErrHorasExcedenMaximo, `name="horas" type="number" step="0.5" value="24.5"`},
		{"horas NaN", "horas", "NaN", domain.ErrHorasInvalidas, `name="horas" type="number" step="0.5" value="NaN"`},
		{"horas Inf", "horas", "Inf", domain.ErrHorasExcedenMaximo, `name="horas" type="number" step="0.5" value="Inf"`},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			form := formEsfuerzoValido()
			form.Set(tt.campo, tt.valor)

			rec := pedirPantalla(t, casosConEsfuerzo(registradorEnvuelto{&registradorEsfuerzoFalso{}}), http.MethodPost, rutaEsfuerzoWeb, form)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			s := seccionEsfuerzo(t, rec)
			seccionContiene(t, s, tt.want.Error(), tt.conservado)
			if strings.Contains(rec.Body.String(), envolturaCaso) {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", rec.Body)
			}
			// El formulario de edición muestra los valores guardados del elemento.
			contiene(t, rec, `name="original_titulo" value="Login"`, `name="titulo" type="text" maxlength="100" value="Login"`)
		})
	}
}

// Traza: HU-18 CA-18.19: al volver con un error se conserva todo lo tipeado.
func TestPantallaEsfuerzo_ConError_ConservaLoTipeado(t *testing.T) {
	form := formEsfuerzoValido()
	form.Set("horas", "1.25")

	rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodPost, rutaEsfuerzoWeb, form)

	seccionContiene(t, seccionEsfuerzo(t, rec), `<option value="Ana" selected>Ana</option>`, `<option value="Caro">Caro</option>`,
		`name="fecha" type="date" value="2026-10-20"`, `name="actividad" type="text" value="Desarrollo"`,
		`name="horas" type="number" step="0.5" value="1.25"`)
}

// Traza: HU-18 CA-18.19, ADR-0003: lo tipeado vuelve al formulario escapado, sin inyectar HTML.
func TestPantallaEsfuerzo_LoTipeado_SeEscapa(t *testing.T) {
	form := formEsfuerzoValido()
	form.Set("actividad", `"><script>alert(1)</script>`)
	form.Set("horas", "1.25")

	rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodPost, rutaEsfuerzoWeb, form)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	seccionContiene(t, seccionEsfuerzo(t, rec), `name="actividad" type="text" value="&#34;&gt;&lt;script&gt;alert(1)&lt;/script&gt;"`)
	if strings.Contains(rec.Body.String(), "<script>") {
		t.Errorf("la página contiene <script> sin escapar; cuerpo: %s", rec.Body)
	}
}

// Traza: HU-18 CA-18.18, ADR-0003: el parámetro de la confirmación se compara y nunca se muestra.
func TestPantallaEsfuerzo_ParametroDeConfirmacion_NoInyecta(t *testing.T) {
	ruta := rutaElementoWeb + "?esfuerzo=" + url.QueryEscape("<script>alert(1)</script>")

	rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodGet, ruta, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if cuerpo := rec.Body.String(); strings.Contains(cuerpo, "<script>") || strings.Contains(cuerpo, "alert(1)") {
		t.Errorf("la página refleja el parámetro; cuerpo: %s", cuerpo)
	}
	seccionNoContiene(t, seccionEsfuerzo(t, rec), "Esfuerzo registrado.")
}

// Traza: HU-18 CA-18.20: errores de lectura del formulario, sin llamar al caso de uso. Si fallan la
// fecha y las horas, se informa la fecha.
func TestPantallaEsfuerzo_ErroresDeFormato_NoLlamanAlCasoDeUso(t *testing.T) {
	const msgFecha, msgHoras = "la fecha debe tener el formato AAAA-MM-DD", "las horas deben ser un número"
	tests := []struct {
		nombre, fecha, horas, want string
		noWant                     string
	}{
		{"fecha con otro formato", "20/10/2026", "1.5", msgFecha, msgHoras},
		{"horas no numéricas", "2026-10-20", "muchas", msgHoras, msgFecha},
		{"horas con coma", "2026-10-20", "1,5", msgHoras, msgFecha},
		{"fallan las dos", "20/10/2026", "muchas", msgFecha, msgHoras},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorEsfuerzoFalso{}
			form := formEsfuerzoValido()
			form.Set("fecha", tt.fecha)
			form.Set("horas", tt.horas)

			rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, rutaEsfuerzoWeb, form)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			s := seccionEsfuerzo(t, rec)
			seccionContiene(t, s, tt.want, `value="`+tt.fecha+`"`, `value="`+tt.horas+`"`, `value="Desarrollo"`)
			seccionNoContiene(t, s, tt.noWant)
			if uc.llamadas != 0 {
				t.Errorf("llamadas al caso de uso = %d, se esperaba 0", uc.llamadas)
			}
		})
	}
}

// Traza: HU-18 CA-18.20: las horas vacías valen 0 y responden como en JSON.
func TestPantallaEsfuerzo_HorasVacias_MensajeDeHorasInvalidas(t *testing.T) {
	form := formEsfuerzoValido()
	form.Set("horas", "")

	rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodPost, rutaEsfuerzoWeb, form)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	seccionContiene(t, seccionEsfuerzo(t, rec), domain.ErrHorasInvalidas.Error())
}

// Traza: HU-18 CA-18.21 (RN-18.3): Caro ya no está en el proyecto: el doble la rechaza y alfa solo
// tiene a Ana. El mensaje se muestra sin envolver y Caro no se ofrece.
func TestPantallaEsfuerzo_IntegranteQuitado_400SinElDesplegable(t *testing.T) {
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})
	cu.RegistrarEsfuerzo = registradorEnvuelto{&registradorEsfuerzoFalso{}}
	form := formEsfuerzoValido()
	form.Set("integrante", "Caro")

	rec := pedirPantalla(t, cu, http.MethodPost, rutaEsfuerzoWeb, form)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	s := seccionEsfuerzo(t, rec)
	seccionContiene(t, s, "integrante no encontrado en el proyecto", `<option value="Ana">Ana</option>`,
		`name="fecha" type="date" value="2026-10-20"`, `name="actividad" type="text" value="Desarrollo"`,
		`name="horas" type="number" step="0.5" value="1.5"`)
	seccionNoContiene(t, s, `value="Caro"`, " selected", envolturaCaso)
}

// Traza: HU-18 CA-18.22 (RN-18.1): proyecto o elemento inexistente, o de otro proyecto.
func TestPantallaEsfuerzo_NoEncontrados_Devuelven404ConLaPaginaDeError(t *testing.T) {
	tests := []struct {
		nombre, ruta string
		err          error
		want         error
	}{
		{"proyecto inexistente", "/web/proyectos/404/backlog/7/esfuerzo", nil, domain.ErrProyectoNoEncontrado},
		{"elemento inexistente", "/web/proyectos/7/backlog/404/esfuerzo", nil, domain.ErrElementoNoEncontrado},
		{"elemento de otro proyecto", rutaEsfuerzoWeb, fmt.Errorf("elemento 7 en el proyecto 7: %w", domain.ErrElementoNoEncontrado), domain.ErrElementoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{err: tt.err}), http.MethodPost, tt.ruta, formEsfuerzoValido())

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, tt.want.Error())
			if strings.Contains(rec.Body.String(), "elemento 7 en el proyecto") {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", rec.Body)
			}
		})
	}
}

// Traza: HU-18 CA-18.22: un error del formulario sobre un elemento que ya no existe da 404, porque
// no hay elemento con el que volver a mostrar la pantalla.
func TestPantallaEsfuerzo_ErrorDeFormatoConElementoInexistente_Devuelve404(t *testing.T) {
	form := formEsfuerzoValido()
	form.Set("horas", "muchas")

	rec := pedirPantalla(t, casosConEsfuerzo(&registradorEsfuerzoFalso{}), http.MethodPost, "/web/proyectos/7/backlog/404/esfuerzo", form)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
	}
	esPaginaDeError(t, rec, domain.ErrElementoNoEncontrado.Error())
}

// Traza: HU-18 CA-18.22 (sección 7: primero el proyecto): un error de formato no llega al caso de
// uso, pero con un proyecto inexistente se informa el proyecto, no el elemento.
func TestPantallaEsfuerzo_ErrorDeFormatoConProyectoInexistente_Devuelve404(t *testing.T) {
	tests := []struct{ nombre, campo, valor string }{
		{"horas no numéricas", "horas", "muchas"},
		{"fecha con otro formato", "fecha", "20/10/2026"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorEsfuerzoFalso{}
			form := formEsfuerzoValido()
			form.Set(tt.campo, tt.valor)

			rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, "/web/proyectos/999/backlog/1/esfuerzo", form)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, domain.ErrProyectoNoEncontrado.Error())
			if uc.llamadas != 0 {
				t.Errorf("llamadas al caso de uso = %d, se esperaba 0", uc.llamadas)
			}
		})
	}
}

// Traza: HU-18 CA-18.22, sección 9: errorDeElemento reconoce ErrProyectoNoEncontrado. La pantalla
// del elemento ahora consulta el proyecto; sin esto respondería 500.
func TestPantallaElemento_ProyectoInexistente_Devuelve404(t *testing.T) {
	cu := casosConEsfuerzo(&registradorEsfuerzoFalso{})
	cu.ConsultarProyecto = consultarFunc(func(context.Context, int64) (domain.Proyecto, error) {
		return domain.Proyecto{}, fmt.Errorf("buscar proyecto 7: %w", domain.ErrProyectoNoEncontrado)
	})

	rec := pedirPantalla(t, cu, http.MethodGet, rutaElementoWeb, nil)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
	}
	esPaginaDeError(t, rec, domain.ErrProyectoNoEncontrado.Error())
}

// Traza: HU-18 CA-18.22: IDs inválidos en la ruta, página de error 400 sin llamar al caso de uso.
func TestPantallaEsfuerzo_IDInvalido_Devuelve400(t *testing.T) {
	tests := []struct{ nombre, ruta, want string }{
		{"proyecto", "/web/proyectos/abc/backlog/7/esfuerzo", "el id del proyecto debe ser un número entero positivo"},
		{"elemento", "/web/proyectos/7/backlog/abc/esfuerzo", "el id del elemento debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &registradorEsfuerzoFalso{}

			rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, tt.ruta, formEsfuerzoValido())

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, tt.want)
			if uc.llamadas != 0 {
				t.Errorf("llamadas al caso de uso = %d, se esperaba 0", uc.llamadas)
			}
		})
	}
}

// Traza: HU-18 CA-18.23: un envío desde otro sitio se rechaza, sin llamar al caso de uso.
func TestPantallaEsfuerzo_DesdeOtroSitio_SeRechaza(t *testing.T) {
	tests := []struct{ encabezado, valor string }{
		{"Sec-Fetch-Site", "cross-site"},
		{"Origin", "http://malicioso.example"},
	}
	for _, tt := range tests {
		t.Run(tt.encabezado, func(t *testing.T) {
			uc := &registradorEsfuerzoFalso{}
			mux := http.NewServeMux()
			RegistrarPantallas(mux, casosConEsfuerzo(uc))
			req := httptest.NewRequest(http.MethodPost, rutaEsfuerzoWeb, strings.NewReader(formEsfuerzoValido().Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set(tt.encabezado, tt.valor)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, se esperaba 403", rec.Code)
			}
			if uc.llamadas != 0 {
				t.Errorf("llamadas al caso de uso = %d, se esperaba 0", uc.llamadas)
			}
		})
	}
}

// Traza: HU-18 CA-18.23: una falla interna muestra la página genérica, sin el detalle.
func TestPantallaEsfuerzo_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	uc := &registradorEsfuerzoFalso{err: errors.New("guardar registro: database is locked")}

	rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, rutaEsfuerzoWeb, formEsfuerzoValido())

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500; cuerpo: %.400s", rec.Code, rec.Body)
	}
	esPaginaDeError(t, rec, "Ocurrió un error interno")
	if strings.Contains(rec.Body.String(), "database") {
		t.Errorf("la página muestra el detalle interno; cuerpo: %.400s", rec.Body)
	}
}

// Traza: HU-18 sección 9: una falla al consultar los integrantes también es una falla interna.
func TestPantallaElemento_FallaAlConsultarElProyecto_Devuelve500SinDetalle(t *testing.T) {
	cu := casosConEsfuerzo(&registradorEsfuerzoFalso{})
	cu.ConsultarProyecto = consultarFunc(func(context.Context, int64) (domain.Proyecto, error) {
		return domain.Proyecto{}, errors.New("database is locked")
	})

	rec := pedirPantalla(t, cu, http.MethodGet, rutaElementoWeb, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "database") {
		t.Errorf("la página muestra el detalle interno; cuerpo: %.400s", rec.Body)
	}
}

// Traza: HU-18 sección 9: un formulario demasiado grande da 413, sin llamar al caso de uso.
func TestPantallaEsfuerzo_FormularioGrande_Devuelve413(t *testing.T) {
	uc := &registradorEsfuerzoFalso{}
	form := formEsfuerzoValido()
	form.Set("actividad", strings.Repeat("a", tamanoMaximoFormulario+1))

	rec := pedirPantalla(t, casosConEsfuerzo(uc), http.MethodPost, rutaEsfuerzoWeb, form)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, se esperaba 413", rec.Code)
	}
	if uc.llamadas != 0 {
		t.Errorf("llamadas al caso de uso = %d, se esperaba 0", uc.llamadas)
	}
}

// Traza: HU-18 CA-18.24: sin integrantes, el formulario avisa y enlaza al detalle del proyecto.
func TestPantallaEsfuerzo_ProyectoSinIntegrantes_AvisaYEnlaza(t *testing.T) {
	cu := casosConEsfuerzo(&registradorEsfuerzoFalso{})
	cu.ConsultarProyecto = consultarFunc(func(context.Context, int64) (domain.Proyecto, error) {
		p := alfa
		p.Integrantes = nil
		return p, nil
	})

	rec := pedirPantalla(t, cu, http.MethodGet, rutaElementoWeb, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	seccionContiene(t, seccionEsfuerzo(t, rec), "no tiene integrantes", `href="/web/proyectos/7"`, `name="integrante"`)
}

// Traza: HU-18 CA-18.25: los errores de editar y eliminar siguen arriba, y la sección de esfuerzo
// aparece con sus valores iniciales.
func TestPantallaElemento_ErroresDeEditarYEliminar_MuestranLaSeccionDeEsfuerzo(t *testing.T) {
	formTituloVacio := formElementoOriginal()
	formTituloVacio.Set("titulo", "")
	tests := []struct {
		nombre, ruta string
		form         url.Values
		editor       *editorElementosFalso
		eliminador   *eliminadorElementosFalso
		wantCode     int
		wantMsg      string
	}{
		{"editar", rutaElementoWeb + "/editar", formTituloVacio, &editorElementosFalso{err: domain.ErrTituloVacio}, &eliminadorElementosFalso{},
			http.StatusBadRequest, domain.ErrTituloVacio.Error()},
		{"eliminar", rutaElementoWeb + "/eliminar", url.Values{}, &editorElementosFalso{}, &eliminadorElementosFalso{err: domain.ErrElementoConEsfuerzo},
			http.StatusConflict, domain.ErrElementoConEsfuerzo.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu := casosConEsfuerzo(&registradorEsfuerzoFalso{})
			cu.EditarElementoBacklog = tt.editor
			cu.EliminarElementoBacklog = tt.eliminador

			rec := pedirConReloj(t, cu, http.MethodPost, tt.ruta, tt.form)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantCode, rec.Body)
			}
			s := seccionEsfuerzo(t, rec)
			seccionContiene(t, s, `<option value="Ana">Ana</option>`, `<option value="Caro">Caro</option>`, valorHoyPantalla,
				`name="actividad" type="text" value=""`, `name="horas" type="number" step="0.5" value=""`)
			seccionNoContiene(t, s, tt.wantMsg, `class="alerta"`)
			arriba, _, _ := strings.Cut(rec.Body.String(), `id="esfuerzo"`)
			if !strings.Contains(arriba, tt.wantMsg) {
				t.Errorf("el mensaje %q no se muestra arriba de la sección; cuerpo: %.400s", tt.wantMsg, rec.Body)
			}
		})
	}
}
