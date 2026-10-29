package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type consultarElementoFunc func(context.Context, int64, int64) (domain.ElementoBacklog, error)

func (f consultarElementoFunc) Ejecutar(ctx context.Context, proyectoID, elementoID int64) (domain.ElementoBacklog, error) {
	return f(ctx, proyectoID, elementoID)
}

// casosBacklogFalsos agrega a casosFalsos el elemento Login (7) del proyecto Alfa (7).
func casosBacklogFalsos(editor *editorElementosFalso, eliminador *eliminadorElementosFalso) CasosDeUso {
	cu := casosFalsos()
	cu.ConsultarElementoBacklog = consultarElementoFunc(func(_ context.Context, pid, eid int64) (domain.ElementoBacklog, error) {
		if pid != 7 || eid != 7 {
			return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
		}
		e := elementoLogin()
		e.ProyectoID = 7
		return e, nil
	})
	cu.EditarElementoBacklog = editor
	cu.EliminarElementoBacklog = eliminador
	return cu
}

// formElementoOriginal es el formulario tal como lo devuelve la pantalla de Login, sin cambios.
func formElementoOriginal() url.Values {
	v := url.Values{
		"titulo": {"Login"}, "descripcion": {"Ingresar"}, "prioridad": {"Alta"},
		"story_points": {"5"}, "criterios": {"Con email"}, "estado": {"Pendiente"},
	}
	for k, val := range v {
		v["original_"+k] = val
	}
	return v
}

// Traza: HU-06 CA-06.10
func TestPantallaElemento_MuestraElFormulario(t *testing.T) {
	rec := pedirPantalla(t, casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{}),
		http.MethodGet, "/web/proyectos/7/backlog/7", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec, "<header", `name="titulo"`, `value="Login"`, `name="original_titulo" value="Login"`,
		`<option value="Alta" selected>`, `<option value="Pendiente" selected>`, `value="5"`,
		">Con email</textarea>", `action="/web/proyectos/7/backlog/7/eliminar"`, `href="/web/proyectos/7/backlog"`)

	rec = pedirPantalla(t, casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{}),
		http.MethodGet, "/web/proyectos/7/backlog/99", nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("inexistente: status = %d", rec.Code)
	}
}

// Traza: HU-06 CA-06.10 (sección 5: solo los campos cambiados, como DEF-01)
func TestPantallaElemento_EditarEnviaSoloLoQueCambio(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosBacklogFalsos(editor, &eliminadorElementosFalso{})
	cu.EditarElementoBacklog = editarQueAcepta(editor)
	form := formElementoOriginal()
	form.Set("prioridad", "Media")
	form.Set("criterios", "Uno\r\n\r\n  Dos  \r\n")

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/editar", form)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/backlog/7" {
		t.Fatalf("status = %d, Location = %q; cuerpo %.300s", rec.Code, rec.Header().Get("Location"), rec.Body)
	}
	c := editor.recibido.Cambios
	if c.Titulo != nil || c.Descripcion != nil || c.CambiarStoryPoints || c.Estado != nil {
		t.Errorf("se enviaron campos que no cambiaron: %+v", c)
	}
	if c.Prioridad == nil || *c.Prioridad != domain.PrioridadMedia || c.Criterios == nil ||
		len(*c.Criterios) != 2 || (*c.Criterios)[0] != "Uno" || (*c.Criterios)[1] != "Dos" {
		t.Errorf("cambios: prioridad %v, criterios %v", c.Prioridad, c.Criterios)
	}
}

// editarQueAcepta registra lo recibido y responde con el elemento sin validar el proyecto.
func editarQueAcepta(f *editorElementosFalso) editorElementos {
	return editarElementoFunc(func(_ context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error) {
		f.recibido = d
		return elementoLogin().Editar(d.Cambios)
	})
}

type editarElementoFunc func(context.Context, app.DatosEditarElemento) (domain.ElementoBacklog, error)

func (f editarElementoFunc) Ejecutar(ctx context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error) {
	return f(ctx, d)
}

// Traza: HU-06 CA-06.10: story points vacío quita la estimación.
func TestPantallaElemento_StoryPointsVacioQuitaLaEstimacion(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosBacklogFalsos(editor, &eliminadorElementosFalso{})
	cu.EditarElementoBacklog = editarQueAcepta(editor)
	form := formElementoOriginal()
	form.Set("story_points", "")

	pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/editar", form)

	c := editor.recibido.Cambios
	if !c.CambiarStoryPoints || c.StoryPoints != nil {
		t.Errorf("CambiarStoryPoints=%v StoryPoints=%v", c.CambiarStoryPoints, c.StoryPoints)
	}
}

// Traza: HU-06 CA-06.10: los errores se muestran en la misma pantalla.
func TestPantallaElemento_Errores(t *testing.T) {
	tests := []struct {
		nombre   string
		ruta     string
		form     url.Values
		editor   *editorElementosFalso
		elim     *eliminadorElementosFalso
		wantCode int
		wantMsg  string
	}{
		{"título vacío", "/web/proyectos/7/backlog/7/editar", func() url.Values { f := formElementoOriginal(); f.Set("titulo", ""); return f }(),
			&editorElementosFalso{err: domain.ErrTituloVacio}, &eliminadorElementosFalso{}, http.StatusBadRequest, domain.ErrTituloVacio.Error()},
		{"story points no numéricos", "/web/proyectos/7/backlog/7/editar", func() url.Values { f := formElementoOriginal(); f.Set("story_points", "muchos"); return f }(),
			&editorElementosFalso{}, &eliminadorElementosFalso{}, http.StatusBadRequest, "story_points debe ser un número entero"},
		{"eliminar con esfuerzo", "/web/proyectos/7/backlog/7/eliminar", url.Values{},
			&editorElementosFalso{}, &eliminadorElementosFalso{err: domain.ErrElementoConEsfuerzo}, http.StatusConflict, domain.ErrElementoConEsfuerzo.Error()},
		{"eliminar en un sprint activo", "/web/proyectos/7/backlog/7/eliminar", url.Values{}, // RN-06.8
			&editorElementosFalso{}, &eliminadorElementosFalso{err: domain.ErrElementoEnSprintActivo}, http.StatusConflict, domain.ErrElementoEnSprintActivo.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosBacklogFalsos(tt.editor, tt.elim), http.MethodPost, tt.ruta, tt.form)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo %.300s", rec.Code, tt.wantCode, rec.Body)
			}
			contiene(t, rec, tt.wantMsg, `name="titulo"`)
		})
	}
}

// Traza: HU-06 CA-06.10 (CA-06.8)
func TestPantallaElemento_Eliminar(t *testing.T) {
	eliminador := &eliminadorElementosFalso{}
	rec := pedirPantalla(t, casosBacklogFalsos(&editorElementosFalso{}, eliminador),
		http.MethodPost, "/web/proyectos/7/backlog/7/eliminar", url.Values{})

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/backlog" {
		t.Errorf("status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}
	if eliminador.recibido != (app.DatosEliminarElemento{ProyectoID: 7, ElementoID: 7}) {
		t.Errorf("datos recibidos: %+v", eliminador.recibido)
	}
}

// Agregar espacios alrededor del título no cuenta como cambio: el dominio los recorta igual.
func TestPantallaElemento_EspaciosAlrededorNoSonUnCambio(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosBacklogFalsos(editor, &eliminadorElementosFalso{})
	cu.EditarElementoBacklog = editarQueAcepta(editor)
	form := formElementoOriginal()
	form.Set("titulo", "  Login ")
	form.Set("estado", "Hecho")

	pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/editar", form)

	if c := editor.recibido.Cambios; c.Titulo != nil || c.Estado == nil {
		t.Errorf("cambios: título %v, estado %v", c.Titulo, c.Estado)
	}
}

// Eliminar pide confirmación antes de enviar el formulario.
func TestPantallaElemento_EliminarPideConfirmacion(t *testing.T) {
	rec := pedirPantalla(t, casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{}),
		http.MethodGet, "/web/proyectos/7/backlog/7", nil)
	contiene(t, rec, "<details", "<summary>Eliminar elemento</summary>", "Sí, eliminar «Login»")
}

// Traza: HU-06 CA-06.10. Guardar la pantalla sin cambiar nada vuelve al elemento sin mensaje y sin
// llamar al caso de uso (no es un error del usuario), igual que el formulario de proyecto (DEF-01).
func TestPantallaElemento_SinCambios_VuelveSinError(t *testing.T) {
	llamado := false
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})
	cu.EditarElementoBacklog = editarElementoFunc(func(context.Context, app.DatosEditarElemento) (domain.ElementoBacklog, error) {
		llamado = true
		return domain.ElementoBacklog{}, domain.ErrSinCambios
	})

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/editar", formElementoOriginal())

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/backlog/7" {
		t.Errorf("status = %d, Location = %q; se esperaba volver al elemento", rec.Code, rec.Header().Get("Location"))
	}
	if llamado {
		t.Error("no se esperaba llamar al caso de uso si no cambió ningún campo")
	}
}

// Traza: HU-06 CA-06.10. La lista del backlog (HU-07) enlaza cada elemento a su pantalla de
// edición y permite eliminarlo, con confirmación.
func TestPantallaBacklog_EnlazaAEditarYEliminarCadaElemento(t *testing.T) {
	rec := pedirPantalla(t, casosConBacklog(nil, nil, &llamadasBacklog{}), http.MethodGet, "/web/proyectos/7/backlog", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	contiene(t, rec,
		`href="/web/proyectos/7/backlog/2"`, `href="/web/proyectos/7/backlog/1"`,
		`action="/web/proyectos/7/backlog/2/eliminar"`, `action="/web/proyectos/7/backlog/1/eliminar"`,
		`aria-label="Editar Login"`, "Sí, eliminar «Login»")
}

// Traza: HU-06 CA-06.10. Después de eliminar se vuelve a la lista del backlog, y la pantalla
// del elemento enlaza al backlog.
func TestPantallaElemento_VuelveAlBacklog(t *testing.T) {
	cu := casosBacklogFalsos(&editorElementosFalso{}, &eliminadorElementosFalso{})

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/eliminar", url.Values{})
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/web/proyectos/7/backlog" {
		t.Errorf("eliminar: status = %d, Location = %q", rec.Code, rec.Header().Get("Location"))
	}

	rec = pedirPantalla(t, cu, http.MethodGet, "/web/proyectos/7/backlog/7", nil)
	contiene(t, rec, `href="/web/proyectos/7/backlog">Volver al backlog`)
}

// Traza: HU-06 CA-06.10. Un POST sin los campos ocultos original_* (que no viene de la pantalla)
// envía todos los campos, igual que el formulario de proyecto (cambiosDeForm, DEF-01). Antes
// los comparaba contra valores vacíos y podía guardar cero cambios sin avisar.
func TestPantallaElemento_SinValoresOriginales_EnviaTodosLosCampos(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosBacklogFalsos(editor, &eliminadorElementosFalso{})
	cu.EditarElementoBacklog = editarQueAcepta(editor)
	form := url.Values{
		"titulo": {"Login"}, "descripcion": {"Ingresar"}, "prioridad": {"Alta"},
		"story_points": {""}, "criterios": {"Con email"}, "estado": {"Pendiente"},
	}

	rec := pedirPantalla(t, cu, http.MethodPost, "/web/proyectos/7/backlog/7/editar", form)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d; cuerpo %.300s", rec.Code, rec.Body)
	}
	c := editor.recibido.Cambios
	if c.Titulo == nil || c.Descripcion == nil || c.Prioridad == nil || !c.CambiarStoryPoints ||
		c.Criterios == nil || c.Estado == nil {
		t.Errorf("se esperaban todos los campos en los cambios: %+v", c)
	}
}
