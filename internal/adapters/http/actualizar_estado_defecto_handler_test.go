package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type actualizadorEstadoDefectoFalso struct {
	errResolver, errCerrar           error
	recibidoResolver, recibidoCerrar app.DatosActualizarEstadoDefecto
	llamadoResolver, llamadoCerrar   bool
}

func (f *actualizadorEstadoDefectoFalso) Resolver(_ context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	f.llamadoResolver, f.recibidoResolver = true, d
	return domain.Defecto{}, f.errResolver
}

func (f *actualizadorEstadoDefectoFalso) Cerrar(_ context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error) {
	f.llamadoCerrar, f.recibidoCerrar = true, d
	return domain.Defecto{}, f.errCerrar
}

func postEstadoDefecto(t *testing.T, uc actualizadorEstadoDefecto, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarActualizarEstadoDefecto(mux, uc)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, ruta, nil))
	return rec
}

// Traza: HU-21 CA-21.7: pedido válido a resolver, 204 sin cuerpo.
func TestPostResolverDefecto_PedidoValido_Devuelve204SinCuerpo(t *testing.T) {
	uc := &actualizadorEstadoDefectoFalso{}

	rec := postEstadoDefecto(t, uc, "/proyectos/7/defectos/3/resolver")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("cuerpo = %q, se esperaba vacío", rec.Body)
	}
	if !uc.llamadoResolver || uc.recibidoResolver != (app.DatosActualizarEstadoDefecto{ProyectoID: 7, DefectoID: 3}) {
		t.Errorf("datos recibidos por Resolver: %+v (llamado: %v)", uc.recibidoResolver, uc.llamadoResolver)
	}
}

// Traza: HU-21 CA-21.7: pedido válido a cerrar, 204 sin cuerpo.
func TestPostCerrarDefecto_PedidoValido_Devuelve204SinCuerpo(t *testing.T) {
	uc := &actualizadorEstadoDefectoFalso{}

	rec := postEstadoDefecto(t, uc, "/proyectos/7/defectos/3/cerrar")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204; cuerpo: %s", rec.Code, rec.Body)
	}
	if !uc.llamadoCerrar || uc.recibidoCerrar != (app.DatosActualizarEstadoDefecto{ProyectoID: 7, DefectoID: 3}) {
		t.Errorf("datos recibidos por Cerrar: %+v (llamado: %v)", uc.recibidoCerrar, uc.llamadoCerrar)
	}
}

// Traza: HU-21 CA-21.7: proyecto o defecto inexistente, 404.
func TestPostResolverDefecto_NoEncontrado_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
	}{
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado},
		{"defecto inexistente", domain.ErrDefectoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := postEstadoDefecto(t, &actualizadorEstadoDefectoFalso{errResolver: tt.err}, "/proyectos/7/defectos/3/resolver")

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.err.Error() {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.err.Error())
			}
		})
	}
}

// Traza: HU-21 CA-21.7: estado inválido para la acción, 400.
func TestPostActualizarEstadoDefecto_EstadoInvalido_Devuelve400(t *testing.T) {
	rec := postEstadoDefecto(t, &actualizadorEstadoDefectoFalso{errResolver: domain.ErrDefectoNoAbierto}, "/proyectos/7/defectos/3/resolver")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := mensajeDeError(t, rec); got != domain.ErrDefectoNoAbierto.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", got, domain.ErrDefectoNoAbierto.Error())
	}

	rec = postEstadoDefecto(t, &actualizadorEstadoDefectoFalso{errCerrar: domain.ErrDefectoNoResuelto}, "/proyectos/7/defectos/3/cerrar")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
	}
	if got := mensajeDeError(t, rec); got != domain.ErrDefectoNoResuelto.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", got, domain.ErrDefectoNoResuelto.Error())
	}
}

// IDs inválidos en la ruta: 400, sin llamar al caso de uso.
func TestPostActualizarEstadoDefecto_IDInvalido_Devuelve400(t *testing.T) {
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"id de proyecto no numérico", "/proyectos/abc/defectos/3/resolver", "el id del proyecto debe ser un número entero positivo"},
		{"id de defecto no numérico", "/proyectos/7/defectos/abc/resolver", "el id del defecto debe ser un número entero positivo"},
		{"id de defecto en cero", "/proyectos/7/defectos/0/cerrar", "el id del defecto debe ser un número entero positivo"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &actualizadorEstadoDefectoFalso{}

			rec := postEstadoDefecto(t, uc, tt.ruta)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
			if uc.llamadoResolver || uc.llamadoCerrar {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}

// Cualquier otro error es una falla interna: 500, sin mostrar el detalle.
func TestPostActualizarEstadoDefecto_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	rec := postEstadoDefecto(t, &actualizadorEstadoDefectoFalso{errResolver: errors.New("resolver defecto: database is locked")}, "/proyectos/7/defectos/3/resolver")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, se esperaba 500; cuerpo: %s", rec.Code, rec.Body)
	}
}
