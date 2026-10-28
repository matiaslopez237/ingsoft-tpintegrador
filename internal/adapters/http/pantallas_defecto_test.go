package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type registrarDefectoFunc func(context.Context, app.DatosDefecto) (domain.Defecto, error)

func (f registrarDefectoFunc) Ejecutar(ctx context.Context, d app.DatosDefecto) (domain.Defecto, error) {
	return f(ctx, d)
}

// llamadasDefecto guarda lo que recibió RegistrarDefecto en una prueba.
type llamadasDefecto struct {
	recibido app.DatosDefecto
}

// casosConDefecto suma a casosConBacklog (el backlog de alfa: "Login" ID 2, "Reportes" ID 1) el
// caso de uso de HU-20: usa el dominio real, el elemento 1 queda con el Sprint de detección 9, y
// cualquier otro queda sin Sprint.
func casosConDefecto(llamadas *llamadasDefecto) CasosDeUso {
	cu := casosConBacklog(nil, nil, &llamadasBacklog{})
	cu.RegistrarDefecto = registrarDefectoFunc(func(_ context.Context, d app.DatosDefecto) (domain.Defecto, error) {
		llamadas.recibido = d
		if d.ProyectoID != alfa.ID {
			return domain.Defecto{}, domain.ErrProyectoNoEncontrado
		}
		def, err := domain.NuevoDefecto(d.ProyectoID, d.ElementoID, d.Descripcion, d.Severidad)
		if err != nil {
			return domain.Defecto{}, err
		}
		def.ID = 1
		if d.ElementoID == 1 {
			nueve := int64(9)
			def.SprintDeteccionID = &nueve
		}
		return def, nil
	})
	return cu
}

func formDefectoValido() url.Values {
	return url.Values{"elemento_id": {"2"}, "descripcion": {"No carga la pantalla de login"}, "severidad": {"Alta"}}
}

// Traza: HU-20 CA-20.8
func TestPantallaDefecto_MuestraElFormulario(t *testing.T) {
	rec := pedirPantalla(t, casosConDefecto(&llamadasDefecto{}), http.MethodGet, "/web/proyectos/7/defectos/nuevo", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec, `<option value="2">Login</option>`, `<option value="1">Reportes</option>`,
		`name="descripcion"`, `<option value="Crítica">`, `<option value="Alta">`,
		`<option value="Media">`, `<option value="Baja">`, `action="/web/proyectos/7/defectos/nuevo"`)
}

func TestPantallaDefecto_ProyectoInexistente_Es404(t *testing.T) {
	rec := pedirPantalla(t, casosConDefecto(&llamadasDefecto{}), http.MethodGet, "/web/proyectos/99/defectos/nuevo", nil)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, se esperaba 404", rec.Code)
	}
}

// Traza: HU-20 CA-20.1, CA-20.8 (RN-20.5: se informa el Sprint de detección asignado)
func TestPantallaDefecto_Crear_ConSprintDeteccion_MuestraConfirmacion(t *testing.T) {
	llamadas := &llamadasDefecto{}
	form := formDefectoValido()
	form.Set("elemento_id", "1") // tiene Sprint de detección 9, según casosConDefecto

	rec := pedirPantalla(t, casosConDefecto(llamadas), http.MethodPost, "/web/proyectos/7/defectos/nuevo", form)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec, "registr", "Sprint", "9")
	want := app.DatosDefecto{ProyectoID: 7, ElementoID: 1, Descripcion: "No carga la pantalla de login", Severidad: domain.SeveridadAlta}
	if llamadas.recibido != want {
		t.Errorf("el caso de uso recibió %+v, se esperaba %+v", llamadas.recibido, want)
	}
}

// Traza: HU-20 CA-20.2, CA-20.8 (RN-20.5: sin Sprint vigente, se avisa que no se asignó ninguno)
func TestPantallaDefecto_Crear_SinSprintDeteccion_AvisaQueNoSeAsigno(t *testing.T) {
	rec := pedirPantalla(t, casosConDefecto(&llamadasDefecto{}), http.MethodPost, "/web/proyectos/7/defectos/nuevo", formDefectoValido())

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec, "No se asignó ningún Sprint")
}

// Traza: HU-20 CA-20.3, CA-20.4, CA-20.8
func TestPantallaDefecto_Crear_DatosInvalidos_MuestraElError(t *testing.T) {
	tests := []struct {
		nombre  string
		cambiar func(url.Values)
		wantMsg string
	}{
		{"descripción vacía", func(f url.Values) { f.Set("descripcion", "  ") }, domain.ErrDescripcionDefectoVacia.Error()},
		{"severidad inválida", func(f url.Values) { f.Set("severidad", "urgente") }, domain.ErrSeveridadInvalida.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			form := formDefectoValido()
			tt.cambiar(form)

			rec := pedirPantalla(t, casosConDefecto(&llamadasDefecto{}), http.MethodPost, "/web/proyectos/7/defectos/nuevo", form)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo %.300s", rec.Code, rec.Body)
			}
			contiene(t, rec, tt.wantMsg)
		})
	}
}
