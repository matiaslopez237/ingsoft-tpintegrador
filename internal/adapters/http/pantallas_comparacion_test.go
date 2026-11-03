package http

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

const rutaEditarWeb = rutaElementoWeb + "/editar"

// comparacionFija usa la función real del dominio: la pantalla muestra lo que el caso de uso devuelve.
func comparacionFija(estimadas *float64, reales float64) *comparadorEsfuerzoFalso {
	return &comparadorEsfuerzoFalso{comparacion: domain.CompararEsfuerzo(estimadas, reales)}
}

// casosConComparacion arma la pantalla del elemento con esa comparación.
func casosConComparacion(uc comparadorEsfuerzo) CasosDeUso {
	cu := casosConEsfuerzo(&registradorEsfuerzoFalso{})
	cu.CompararEsfuerzo = uc
	return cu
}

// seccionComparacion devuelve el HTML del bloque "Esfuerzo estimado y real".
func seccionComparacion(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	_, resto, ok := strings.Cut(rec.Body.String(), `id="comparacion"`)
	if !ok {
		t.Fatalf("la página no tiene el bloque de comparación; status = %d, cuerpo: %.400s", rec.Code, rec.Body)
	}
	seccion, _, ok := strings.Cut(resto, "</section>")
	if !ok {
		t.Fatalf("el bloque de comparación no cierra; cuerpo: %.400s", resto)
	}
	return seccion
}

// valorDeComparacion arma el par que tiene que aparecer en el bloque.
func valorDeComparacion(rotulo, valor string) string {
	return "<dt>" + rotulo + "</dt><dd>" + valor + "</dd>"
}

// bloqueMuestra verifica las cuatro cifras del bloque tal como las ve quien usa la pantalla:
// html/template escribe el + como &#43; y el navegador lo muestra como +.
func bloqueMuestra(t *testing.T, s, estimadas, reales, horas, porcentaje string) {
	t.Helper()
	seccionContiene(t, html.UnescapeString(s),
		valorDeComparacion("Horas estimadas", estimadas),
		valorDeComparacion("Horas reales", reales),
		valorDeComparacion("Desviación en horas", horas),
		valorDeComparacion("Desviación en porcentaje", porcentaje))
}

// formConHoras es el formulario de Login tal como vuelve de la pantalla, con esas horas originales y tipeadas.
func formConHoras(original, tipeadas string) url.Values {
	form := formElementoOriginal()
	form.Set("original_horas_estimadas", original)
	form.Set("horas_estimadas", tipeadas)
	return form
}

// elementoConHoras es Login (7) del proyecto Alfa (7) con esas horas estimadas guardadas.
func elementoConHoras(horas *float64) func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
	return func(context.Context, int64, int64) (domain.ElementoBacklog, error) {
		e := elementoLogin()
		e.ProyectoID = 7
		e.HorasEstimadas = horas
		return e, nil
	}
}

const envolturaEditar = "editar elemento 7:"

// editarQueEnvuelve valida con el dominio real y envuelve los errores con %w.
func editarQueEnvuelve(f *editorElementosFalso) editorElementos {
	return editarElementoFunc(func(_ context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error) {
		f.recibido = d
		e, err := elementoLogin().Editar(d.Cambios)
		if err != nil {
			return e, fmt.Errorf("%s %w", envolturaEditar, err)
		}
		return e, nil
	})
}

// Formato solo de presentación: coma decimal, un decimal en el porcentaje y signo en las positivas.
// Traza: HU-19 CA-19.18 (RN-19.5, RN-19.6, RN-19.7, RN-19.8, RN-19.9)
func TestPantallaComparacion_MuestraLasCuatroCifras(t *testing.T) {
	const sinEstimar = "Sin estimar"
	tests := []struct {
		nombre                             string
		estimadas                          *float64
		reales                             float64
		wantEst, wantReal, wantHs, wantPor string
	}{
		{"con estimación y esfuerzo, de más", pHoras(8), 10, "8 h", "10 h", "+2 h", "+25,0 %"},
		{"con estimación y esfuerzo, de menos", pHoras(8), 6, "8 h", "6 h", "-2 h", "-25,0 %"},
		{"sin desviación: sin signo", pHoras(8), 8, "8 h", "8 h", "0 h", "0,0 %"},
		{"sin esfuerzo: 0 h", pHoras(8), 0, "8 h", "0 h", "-8 h", "-100,0 %"},
		{"sin estimar, con esfuerzo", nil, 10, sinEstimar, "10 h", sinEstimar, sinEstimar},
		{"sin estimar y sin esfuerzo", nil, 0, sinEstimar, "0 h", sinEstimar, sinEstimar},
		{"coma decimal en las horas", pHoras(7.5), 10.5, "7,5 h", "10,5 h", "+3 h", "+40,0 %"},
		{"un decimal en el porcentaje", pHoras(3), 4, "3 h", "4 h", "+1 h", "+33,3 %"},
		{"desviación con media hora", pHoras(40), 40.5, "40 h", "40,5 h", "+0,5 h", "+1,3 %"},
		{"negativa con media hora", pHoras(40), 39.5, "40 h", "39,5 h", "-0,5 h", "-1,3 %"},
		{"el redondeo da cero: sin signo", pHoras(2000), 1999.5, "2000 h", "1999,5 h", "-0,5 h", "0,0 %"},
		{"estimado 0: sin porcentaje", pHoras(0), 10, "0 h", "10 h", "+10 h", "No aplica"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := comparacionFija(tt.estimadas, tt.reales)

			rec := pedirPantalla(t, casosConComparacion(uc), http.MethodGet, rutaElementoWeb, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			s := seccionComparacion(t, rec)
			seccionContiene(t, s, `<h2 id="titulo-comparacion">Esfuerzo estimado y real</h2>`)
			bloqueMuestra(t, s, tt.wantEst, tt.wantReal, tt.wantHs, tt.wantPor)
			if uc.proyectoID != 7 || uc.elementoID != 7 {
				t.Errorf("el caso de uso recibió el proyecto %d y el elemento %d", uc.proyectoID, uc.elementoID)
			}
		})
	}
}

// El bloque muestra lo guardado, no lo que se acaba de tipear, en todos los caminos que vuelven a la
// pantalla con un error. Lo tipeado en horas estimadas (20) sigue en el formulario.
// Traza: HU-19 CA-19.18 (sección 5: rutas y pantallas)
func TestPantallaComparacion_AlVolverConError_MuestraLoGuardado(t *testing.T) {
	formEsfuerzoMalo := formEsfuerzoValido()
	formEsfuerzoMalo.Set("horas", "1.25")
	tests := []struct {
		nombre, ruta string
		form         url.Values
		preparar     func(cu *CasosDeUso)
		wantCode     int
		wantMsg      string
		wantForm     string
	}{
		{"editar", rutaEditarWeb, formConHoras("8", "20"), func(cu *CasosDeUso) {
			cu.EditarElementoBacklog = &editorElementosFalso{err: domain.ErrTituloDuplicado}
		}, http.StatusBadRequest, domain.ErrTituloDuplicado.Error(), `name="horas_estimadas" type="number" step="0.5" value="20"`},
		{"eliminar", rutaElementoWeb + "/eliminar", url.Values{}, func(cu *CasosDeUso) {
			cu.EliminarElementoBacklog = &eliminadorElementosFalso{err: domain.ErrElementoConEsfuerzo}
		}, http.StatusConflict, domain.ErrElementoConEsfuerzo.Error(), `name="horas_estimadas"`},
		{"registrar esfuerzo", rutaEsfuerzoWeb, formEsfuerzoMalo, nil,
			http.StatusBadRequest, domain.ErrHorasNoMultiploMediaHora.Error(), `name="horas_estimadas"`},
		{"iniciar ronda", rutaElementoWeb + "/rondas", url.Values{}, func(cu *CasosDeUso) {
			poker := conRonda(t)
			cu.IniciarRonda, cu.Votar, cu.ConsultarRonda = poker, votadorPokerFalso{poker}, consultorPokerFalso{poker}
		}, http.StatusBadRequest, domain.ErrRondaYaAbierta.Error(), `name="horas_estimadas"`},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu := casosConComparacion(comparacionFija(pHoras(8), 10))
			if tt.preparar != nil {
				tt.preparar(&cu)
			}

			rec := pedirPantalla(t, cu, http.MethodPost, tt.ruta, tt.form)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, se esperaba %d; cuerpo: %.400s", rec.Code, tt.wantCode, rec.Body)
			}
			contiene(t, rec, tt.wantMsg, tt.wantForm)
			bloqueMuestra(t, seccionComparacion(t, rec), "8 h", "10 h", "+2 h", "+25,0 %")
		})
	}
}

// El formulario muestra las horas guardadas, con punto como exige un input numérico, y deja los
// límites en el dominio: sin min ni max.
// Traza: HU-19 CA-19.18 (RN-19.2, RN-19.3)
func TestPantallaElemento_CampoHorasEstimadas(t *testing.T) {
	tests := []struct {
		nombre string
		horas  *float64
		valor  string
	}{
		{"sin estimar", nil, ""},
		{"con horas", pHoras(8), "8"},
		{"con media hora", pHoras(7.5), "7.5"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			cu := casosConComparacion(comparacionFija(tt.horas, 0))
			cu.ConsultarElementoBacklog = consultarElementoFunc(elementoConHoras(tt.horas))

			rec := pedirPantalla(t, cu, http.MethodGet, rutaElementoWeb, nil)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, `<label for="horas_estimadas">Horas estimadas`,
				`<input id="horas_estimadas" name="horas_estimadas" type="number" step="0.5" value="`+tt.valor+`">`,
				`<input type="hidden" name="original_horas_estimadas" value="`+tt.valor+`">`)
		})
	}
}

// Traza: HU-19 CA-19.13, CA-19.18 (RN-19.3)
func TestPantallaElemento_CargarYQuitarHorasEstimadas(t *testing.T) {
	tests := []struct {
		nombre, original, tipeadas string
		want                       *float64
	}{
		{"cargar", "", "8", pHoras(8)},
		{"cargar con media hora", "", "7.5", pHoras(7.5)},
		{"cambiar", "8", "10", pHoras(10)},
		{"con espacios", "", " 8 ", pHoras(8)},
		{"vacío las quita", "8", "", nil},
		{"solo espacios las quita", "8", "   ", nil},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			editor := &editorElementosFalso{}
			cu := casosConComparacion(comparacionFija(nil, 0))
			cu.EditarElementoBacklog = editarQueAcepta(editor)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formConHoras(tt.original, tt.tipeadas))

			if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != rutaElementoWeb {
				t.Fatalf("status = %d, Location = %q; cuerpo: %.400s", rec.Code, rec.Header().Get("Location"), rec.Body)
			}
			c := editor.recibido.Cambios
			if !c.CambiarHorasEstimadas {
				t.Fatalf("CambiarHorasEstimadas = false; cambios: %+v", c)
			}
			switch {
			case tt.want == nil && c.HorasEstimadas != nil:
				t.Errorf("HorasEstimadas = %v, se esperaba nil", *c.HorasEstimadas)
			case tt.want != nil && (c.HorasEstimadas == nil || *c.HorasEstimadas != *tt.want):
				t.Errorf("HorasEstimadas = %v, se esperaba %v", c.HorasEstimadas, *tt.want)
			}
		})
	}
}

// Si no cambió nada no se llama al caso de uso; si cambió otro campo, las horas no se envían.
// Traza: HU-19 CA-19.13 (RN-19.3; HU-06 sección 5: solo los campos cambiados)
func TestPantallaElemento_HorasSinCambios_NoSeEnvian(t *testing.T) {
	for _, tt := range []struct{ nombre, original, tipeadas string }{
		{"iguales", "8", "8"},
		{"con espacios", "8", " 8 "},
		{"las dos vacías", "", ""},
	} {
		t.Run(tt.nombre, func(t *testing.T) {
			editor := &editorElementosFalso{}
			cu := casosConComparacion(comparacionFija(nil, 0))
			cu.EditarElementoBacklog = editarQueAcepta(editor)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formConHoras(tt.original, tt.tipeadas))

			if rec.Code != http.StatusSeeOther {
				t.Fatalf("status = %d; cuerpo: %.400s", rec.Code, rec.Body)
			}
			if editor.recibido.ProyectoID != 0 {
				t.Errorf("no se esperaba llamar al caso de uso: %+v", editor.recibido)
			}

			form := formConHoras(tt.original, tt.tipeadas)
			form.Set("prioridad", "Media")
			pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, form)

			c := editor.recibido.Cambios
			if c.Prioridad == nil || c.CambiarHorasEstimadas || c.HorasEstimadas != nil {
				t.Errorf("cambios enviados: %+v", c)
			}
		})
	}
}

// Traza: HU-19 CA-19.15 (RN-19.1): cada campo viaja con su propia marca.
func TestPantallaElemento_HorasYStoryPoints_NoSePisan(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosConComparacion(comparacionFija(nil, 0))
	cu.EditarElementoBacklog = editarQueAcepta(editor)

	soloHoras := formConHoras("8", "10")
	pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, soloHoras)
	if c := editor.recibido.Cambios; !c.CambiarHorasEstimadas || c.CambiarStoryPoints || c.StoryPoints != nil {
		t.Errorf("cambiando las horas: %+v", c)
	}

	soloPuntos := formConHoras("8", "8")
	soloPuntos.Set("story_points", "13")
	pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, soloPuntos)
	if c := editor.recibido.Cambios; c.CambiarHorasEstimadas || c.HorasEstimadas != nil || !c.CambiarStoryPoints {
		t.Errorf("cambiando los story points: %+v", c)
	}
}

// El doble usa el dominio real: NaN e Inf los parsea ParseFloat y los rechaza la validación verdadera.
// El mensaje sale sin envolver, y el formulario conserva lo tipeado y el valor original.
// Traza: HU-19 CA-19.14, CA-19.18 (RN-19.2)
func TestPantallaElemento_HorasInvalidas_VuelvenAlFormularioCon400(t *testing.T) {
	tests := []struct {
		tipeadas string
		want     error
	}{
		{"0", domain.ErrHorasInvalidas},
		{"-1", domain.ErrHorasInvalidas},
		{"NaN", domain.ErrHorasInvalidas},
		{"1.25", domain.ErrHorasNoMultiploMediaHora},
		{"1000.5", domain.ErrHorasEstimadasExcedenMaximo},
		{"Inf", domain.ErrHorasEstimadasExcedenMaximo},
	}
	for _, tt := range tests {
		t.Run(tt.tipeadas, func(t *testing.T) {
			editor := &editorElementosFalso{}
			cu := casosConComparacion(comparacionFija(pHoras(8), 10))
			cu.EditarElementoBacklog = editarQueEnvuelve(editor)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formConHoras("8", tt.tipeadas))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, tt.want.Error(),
				`name="horas_estimadas" type="number" step="0.5" value="`+tt.tipeadas+`"`,
				`name="original_horas_estimadas" value="8"`)
			if strings.Contains(rec.Body.String(), envolturaEditar) {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", rec.Body)
			}
		})
	}
}

// Lo que no es un número no llega al caso de uso. La coma tampoco se acepta, como en HU-18.
// Traza: HU-19 CA-19.14, CA-19.18 (RN-19.2; sección 7)
func TestPantallaElemento_HorasQueNoSonNumero_NoLlamanAlCasoDeUso(t *testing.T) {
	for _, tipeadas := range []string{"muchas", "1,5", "1e999"} {
		t.Run(tipeadas, func(t *testing.T) {
			editor := &editorElementosFalso{}
			cu := casosConComparacion(comparacionFija(nil, 0))
			cu.EditarElementoBacklog = editarQueAcepta(editor)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formConHoras("", tipeadas))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
			}
			contiene(t, rec, "las horas estimadas deben ser un número", `value="`+tipeadas+`"`)
			if editor.recibido.ProyectoID != 0 {
				t.Errorf("no se esperaba llamar al caso de uso: %+v", editor.recibido)
			}
		})
	}
}

// Traza: HU-19 CA-19.18, ADR-0003: lo tipeado vuelve al formulario escapado.
func TestPantallaElemento_HorasTipeadas_SeEscapan(t *testing.T) {
	cu := casosConComparacion(comparacionFija(nil, 0))

	rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formConHoras("", `"><script>alert(1)</script>`))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "<script>") {
		t.Errorf("la página contiene <script> sin escapar; cuerpo: %s", rec.Body)
	}
}

// El mismo intento de inyección en todos los campos del formulario del elemento, tipeados y ocultos.
// Lo tipeado difiere de lo original para que haya cambios, y story_points no es un número: así la
// página vuelve con 400 y se renderiza todo. Sin el control positivo, el test también pasaría con una
// redirección o con una plantilla que dejara de mostrar los campos.
// Traza: HU-19 CA-19.18, ADR-0003 (HU-06 CA-06.10)
func TestPantallaElemento_TodosLosCampos_SeEscapan(t *testing.T) {
	const (
		ataque        = `"><script>alert(1)</script>`
		ataqueEscrito = `&#34;&gt;&lt;script&gt;alert(1)&lt;/script&gt;`
	)
	campos := []string{"titulo", "descripcion", "prioridad", "story_points", "criterios", "estado", "horas_estimadas"}
	form := url.Values{}
	for _, campo := range campos {
		form.Set(campo, ataque+"-tipeado")
		form.Set("original_"+campo, ataque+"-original")
	}
	cu := casosConComparacion(comparacionFija(nil, 0))

	rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, form)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "<script>") {
		t.Errorf("la página contiene <script> sin escapar; cuerpo: %s", rec.Body)
	}
	for _, campo := range campos {
		contiene(t, rec, `name="original_`+campo+`" value="`+ataqueEscrito+`-original"`)
	}
	// prioridad y estado son desplegables: lo tipeado no vuelve a la página. story_points también es
	// un desplegable (HU-13), pero un valor fuera de la escala vuelve como opción elegida: en el
	// atributo y en el texto.
	tipeado := ataqueEscrito + "-tipeado"
	contiene(t, rec,
		`name="titulo" type="text" maxlength="100" value="`+tipeado+`"`,
		`name="descripcion" rows="3" required>`+tipeado+`</textarea>`,
		`<option value="`+tipeado+`" selected>`+tipeado+` (fuera de la escala)</option>`,
		`name="criterios" rows="4" required>`+tipeado+`</textarea>`,
		`name="horas_estimadas" type="number" step="0.5" value="`+tipeado+`"`,
		`<h1>`+ataqueEscrito+`-original</h1>`,
		`<title>`+ataqueEscrito+`-original`)
}

// Un POST sin los campos ocultos envía todos los campos, como los story points (DEF-01).
// Traza: HU-19 CA-19.13 (RN-19.3; HU-06 CA-06.10)
func TestPantallaElemento_SinOriginales_EnviaLasHoras(t *testing.T) {
	editor := &editorElementosFalso{}
	cu := casosConComparacion(comparacionFija(nil, 0))
	cu.EditarElementoBacklog = editarQueAcepta(editor)
	form := url.Values{
		"titulo": {"Login"}, "descripcion": {"Ingresar"}, "prioridad": {"Alta"},
		"story_points": {"5"}, "criterios": {"Con email"}, "estado": {"Pendiente"}, "horas_estimadas": {"8"},
	}

	pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, form)

	c := editor.recibido.Cambios
	if !c.CambiarHorasEstimadas || c.HorasEstimadas == nil || *c.HorasEstimadas != 8 {
		t.Errorf("cambios: %+v", c)
	}
}

// editarSobreGuardado aplica los cambios con el dominio real a Login con 5 story points y 8 h
// estimadas, y deja el resultado en guardado.
func editarSobreGuardado(guardado *domain.ElementoBacklog) editorElementos {
	return editarElementoFunc(func(_ context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error) {
		e := elementoLogin()
		e.HorasEstimadas = pHoras(8)
		editado, err := e.Editar(d.Cambios)
		if err != nil {
			return domain.ElementoBacklog{}, err
		}
		*guardado = editado
		return editado, nil
	})
}

// formSinEstimaciones es un POST que no viene de la pantalla: sin original_*, sin story_points y sin
// horas_estimadas.
func formSinEstimaciones() url.Values {
	return url.Values{"titulo": {"Login"}, "descripcion": {"Ingresar"}, "prioridad": {"Alta"}, "criterios": {"Con email"}, "estado": {"Pendiente"}}
}

// Un campo horas_estimadas ausente significa "no se envió": no cambia las horas, aunque falten los
// ocultos. Presente y vacío las quita si difiere del original; presente con valor las cambia.
// Traza: HU-19 CA-19.19 (RN-19.12), CA-19.13, CA-19.18 (RN-19.3)
func TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran(t *testing.T) {
	conOriginales := func(modificar func(url.Values)) url.Values {
		form := formElementoOriginal()
		form.Set("prioridad", "Media") // para que siempre haya algún cambio que guardar
		modificar(form)
		return form
	}
	tests := []struct {
		nombre string
		form   url.Values
		want   *float64
	}{
		{"sin original_* y sin el campo: se conservan", formSinEstimaciones(), pHoras(8)},
		{"con original_* y sin el campo ni su oculto: se conservan", conOriginales(func(url.Values) {}), pHoras(8)},
		{"con su oculto y sin el campo: se conservan", conOriginales(func(f url.Values) { f.Set("original_horas_estimadas", "8") }), pHoras(8)},
		{"sin original_* y vacío: se quitan", func() url.Values {
			f := formSinEstimaciones()
			f.Set("horas_estimadas", "")
			return f
		}(), nil},
		{"con su oculto y vacío: se quitan", conOriginales(func(f url.Values) {
			f.Set("original_horas_estimadas", "8")
			f.Set("horas_estimadas", "")
		}), nil},
		{"con su oculto vacío y vacío: no difieren, se conservan", conOriginales(func(f url.Values) {
			f.Set("original_horas_estimadas", "")
			f.Set("horas_estimadas", "")
		}), pHoras(8)},
		{"sin original_* y con valor: se cambian", func() url.Values {
			f := formSinEstimaciones()
			f.Set("horas_estimadas", "10")
			return f
		}(), pHoras(10)},
		{"con su oculto y con valor: se cambian", conOriginales(func(f url.Values) {
			f.Set("original_horas_estimadas", "8")
			f.Set("horas_estimadas", "10")
		}), pHoras(10)},
		{"sin su oculto y con valor: se cambian", conOriginales(func(f url.Values) { f.Set("horas_estimadas", "10") }), pHoras(10)},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			var guardado domain.ElementoBacklog
			cu := casosConComparacion(comparacionFija(nil, 0))
			cu.EditarElementoBacklog = editarSobreGuardado(&guardado)

			rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, tt.form)

			if rec.Code != http.StatusSeeOther {
				t.Fatalf("status = %d, se esperaba 303; cuerpo: %.400s", rec.Code, rec.Body)
			}
			switch {
			case tt.want == nil && guardado.HorasEstimadas != nil:
				t.Errorf("horas guardadas = %v, se esperaba sin estimar", *guardado.HorasEstimadas)
			case tt.want != nil && (guardado.HorasEstimadas == nil || *guardado.HorasEstimadas != *tt.want):
				t.Errorf("horas guardadas = %v, se esperaba %v", guardado.HorasEstimadas, *tt.want)
			}
		})
	}
}

// La corrección es solo para las horas: story_points mantiene el comportamiento de HU-06 y DEF-01.
// Un POST sin original_* y sin el campo sigue quitándolos.
// Traza: HU-19 CA-19.19 (RN-19.12), CA-19.15 (RN-19.1), HU-06 CA-06.10
func TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06(t *testing.T) {
	var guardado domain.ElementoBacklog
	cu := casosConComparacion(comparacionFija(nil, 0))
	cu.EditarElementoBacklog = editarSobreGuardado(&guardado)

	rec := pedirPantalla(t, cu, http.MethodPost, rutaEditarWeb, formSinEstimaciones())

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, se esperaba 303; cuerpo: %.400s", rec.Code, rec.Body)
	}
	if guardado.StoryPoints != nil {
		t.Errorf("story points = %d, se esperaba que siguieran quitándose", *guardado.StoryPoints)
	}
}

// Traza: HU-19 CA-19.9, CA-19.10, CA-19.18 (RN-19.10; sección 7)
func TestPantallaComparacion_NoEncontrados_Devuelven404(t *testing.T) {
	tests := []struct {
		nombre string
		err    error
		want   error
	}{
		{"proyecto inexistente", domain.ErrProyectoNoEncontrado, domain.ErrProyectoNoEncontrado},
		{"elemento inexistente", domain.ErrElementoNoEncontrado, domain.ErrElementoNoEncontrado},
		{"elemento, envuelto", fmt.Errorf("buscar elemento 7: %w", domain.ErrElementoNoEncontrado), domain.ErrElementoNoEncontrado},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := pedirPantalla(t, casosConComparacion(&comparadorEsfuerzoFalso{err: tt.err}), http.MethodGet, rutaElementoWeb, nil)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, tt.want.Error())
			if strings.Contains(rec.Body.String(), "buscar elemento") {
				t.Errorf("la página muestra el error envuelto; cuerpo: %.400s", rec.Body)
			}
		})
	}
}

// Traza: HU-19 CA-19.18 (sección 7): cualquier otro error es 500 genérico, también al volver de un
// error del formulario: la falla interna tiene prioridad sobre el mensaje de validación.
func TestPantallaComparacion_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	tests := []struct {
		nombre, metodo, ruta string
		form                 url.Values
	}{
		{"al abrir la pantalla", http.MethodGet, rutaElementoWeb, nil},
		{"al volver con un error de edición", http.MethodPost, rutaEditarWeb, formConHoras("", "muchas")},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &comparadorEsfuerzoFalso{err: errors.New("consultar horas reales: database is locked")}

			rec := pedirPantalla(t, casosConComparacion(uc), tt.metodo, tt.ruta, tt.form)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, se esperaba 500; cuerpo: %.400s", rec.Code, rec.Body)
			}
			esPaginaDeError(t, rec, "Ocurrió un error interno")
			if strings.Contains(rec.Body.String(), "database") {
				t.Errorf("la página muestra el detalle interno; cuerpo: %.400s", rec.Body)
			}
		})
	}
}
