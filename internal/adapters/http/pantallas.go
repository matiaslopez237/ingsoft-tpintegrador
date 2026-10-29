package http

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// archivosWeb son las plantillas y los archivos estáticos de la interfaz, embebidos en el
// binario (ADR-0003): no hace falta copiar nada al desplegar.
//
//go:embed web
var archivosWeb embed.FS

// listadorProyectos es lo que la lista de proyectos necesita del caso de uso ListarProyectos.
type listadorProyectos interface {
	Ejecutar(ctx context.Context) ([]domain.Proyecto, error)
}

// erroresDeValidacionPantallas son los errores de dominio de HU-01 a HU-03 que las pantallas
// muestran en el mismo formulario con HTTP 400 (spec HU-30, sección 7).
var erroresDeValidacionPantallas = []error{
	domain.ErrNombreVacio,
	domain.ErrNombreDemasiadoLargo,
	domain.ErrNombreDuplicado,
	domain.ErrFechaInicioObligatoria,
	domain.ErrFechaFinAnterior,
	domain.ErrSinCambios,
	domain.ErrEstadoInvalido,
	domain.ErrNombreIntegranteVacio,
	domain.ErrRolIntegranteVacio,
	domain.ErrIntegranteDuplicado,
	domain.ErrIntegranteNoEncontrado,
	domain.ErrRolDuplicado,
	// HU-08: los de crear un sprint.
	domain.ErrProyectoNoActivo,
	domain.ErrFechaFinObligatoria,
	domain.ErrObjetivoVacio,
	domain.ErrObjetivoDemasiadoLargo,
	domain.ErrSprintFueraDelProyecto,
	domain.ErrSprintSuperpuesto,
	// HU-05 (alta desde la pantalla del backlog) y HU-07 (orden y filtro).
	domain.ErrTituloVacio,
	domain.ErrTituloDemasiadoLargo,
	domain.ErrDescripcionVacia,
	domain.ErrPrioridadInvalida,
	domain.ErrStoryPointsInvalidos,
	domain.ErrSinCriterios,
	domain.ErrCriterioVacio,
	domain.ErrTituloDuplicado,
	domain.ErrOrdenInvalido,
	domain.ErrEstadoElementoInvalido,
	// HU-09: asignar/desasignar un elemento a un sprint.
	domain.ErrSprintDeOtroProyecto,
	domain.ErrElementoYaAsignado,
	domain.ErrSprintCerrado,
	domain.ErrElementoSinSprint,
	// HU-11: cerrar un sprint que ya estaba cerrado.
	domain.ErrSprintYaCerrado,
}

// estadosProyecto son las opciones del formulario de edición (RN-01.7). Validarlas es tarea
// del dominio: acá solo se listan para el <select>.
var estadosProyecto = []domain.Estado{
	domain.EstadoActivo, domain.EstadoEnPausa, domain.EstadoFinalizado, domain.EstadoPostergado, domain.EstadoCancelado,
}

// errFormulario marca los errores de formato del formulario (fechas), que se muestran como los de dominio.
type errFormulario struct{ mensaje string }

func (e errFormulario) Error() string { return e.mensaje }

// Datos que reciben las plantillas: solo texto ya formateado, sin lógica de negocio.
type (
	vistaProyecto struct {
		ID          int64
		Nombre      string
		FechaInicio string
		FechaFin    string // vacío: sin fecha de fin
		Estado      string
		Integrantes []domain.Integrante
	}
	formProyecto struct {
		Nombre, FechaInicio, FechaFin, Estado string
	}
	formIntegrante struct {
		Nombre, Rol string
	}
	datosLista struct {
		Proyectos []vistaProyecto
	}
	datosFormProyecto struct {
		ID       int64 // 0: proyecto nuevo
		Form     formProyecto
		Original formProyecto // DEF-01: valores con los que se abrió el formulario de edición
		Error    string
		Estados  []domain.Estado
	}
	datosDetalle struct {
		Proyecto       vistaProyecto
		Error          string
		FormIntegrante formIntegrante
	}
	datosError struct {
		Mensaje string
	}
)

// pantallas arma las páginas HTML. Cada página es el layout común más su propia plantilla, que
// define los bloques "titulo" y "contenido". Las historias que agreguen pantallas suman su
// archivo a web/plantillas y lo nombran en paginas.
type pantallas struct {
	cu      CasosDeUso
	paginas map[string]*template.Template
}

var paginas = []string{
	"error.html", "proyectos.html", "proyecto_form.html", "proyecto_detalle.html",
	"sprints.html", "backlog.html", "sprint_detalle.html", "elemento.html", "defecto_form.html",
}

// funcionesPlantillas son las funciones disponibles en las plantillas. segmento escapa un texto
// para usarlo como parte de una ruta: html/template no escapa "/", "?" ni "#" dentro de una URL.
var funcionesPlantillas = template.FuncMap{"segmento": url.PathEscape}

// tamanoMaximoFormulario limita el cuerpo de los formularios de las pantallas (64 KiB).
const tamanoMaximoFormulario = 64 << 10

func nuevasPantallas(cu CasosDeUso) *pantallas {
	p := &pantallas{cu: cu, paginas: map[string]*template.Template{}}
	for _, nombre := range paginas {
		// Must: las plantillas están embebidas, así que un error acá es de programación y lo
		// detecta cualquier test que registre las pantallas.
		p.paginas[nombre] = template.Must(template.New("layout.html").Funcs(funcionesPlantillas).
			ParseFS(archivosWeb, "web/plantillas/layout.html", "web/plantillas/"+nombre))
	}
	return p
}

// RegistrarPantallas agrega al mux la interfaz web: inicio, archivos estáticos y pantallas de
// proyectos e integrantes. Implementa HU-30 (ADR-0003).
func RegistrarPantallas(mux *http.ServeMux, cu CasosDeUso) {
	p := nuevasPantallas(cu)
	estaticos, err := fs.Sub(archivosWeb, "web")
	if err != nil {
		panic(err) // la carpeta está embebida: no puede faltar
	}

	mux.Handle("GET /static/", http.FileServerFS(estaticos))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/proyectos", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /web/proyectos", p.lista)
	mux.HandleFunc("GET /web/proyectos/nuevo", p.formNuevo)
	mux.HandleFunc("POST /web/proyectos", p.formulario(p.crear))
	mux.HandleFunc("GET /web/proyectos/{id}", p.detalle)
	mux.HandleFunc("GET /web/proyectos/{id}/editar", p.formEditar)
	mux.HandleFunc("POST /web/proyectos/{id}/editar", p.formulario(p.editar))
	mux.HandleFunc("POST /web/proyectos/{id}/integrantes", p.formulario(p.agregarIntegrante))
	mux.HandleFunc("POST /web/proyectos/{id}/integrantes/{nombre}/quitar", p.formulario(p.quitarIntegrante))
	mux.HandleFunc("POST /web/proyectos/{id}/integrantes/{nombre}/roles", p.formulario(p.agregarRol))
	mux.HandleFunc("GET /web/proyectos/{id}/sprints", p.sprints)                                                       // HU-08
	mux.HandleFunc("POST /web/proyectos/{id}/sprints", p.formulario(p.crearSprint))                                    // HU-08
	mux.HandleFunc("GET /web/proyectos/{id}/backlog", p.backlog)                                                       // HU-07
	mux.HandleFunc("POST /web/proyectos/{id}/backlog", p.formulario(p.crearElemento))                                  // HU-07 (alta de HU-05)
	mux.HandleFunc("GET /web/proyectos/{id}/sprints/{sprintID}", p.sprintDetalle)                                      // HU-09
	mux.HandleFunc("POST /web/proyectos/{id}/sprints/{sprintID}/asignar", p.formulario(p.asignarElementoPantalla))     // HU-09
	mux.HandleFunc("POST /web/proyectos/{id}/sprints/{sprintID}/quitar", p.formulario(p.quitarElementoPantalla))       // HU-09
	mux.HandleFunc("POST /web/proyectos/{id}/sprints/{sprintID}/completar", p.formulario(p.completarElementoPantalla)) // HU-10
	mux.HandleFunc("POST /web/proyectos/{id}/sprints/{sprintID}/cerrar", p.formulario(p.cerrarSprintPantalla))         // HU-11
	p.registrarPantallasElemento(mux)                                                                                  // HU-06
	p.registrarPantallasDefecto(mux)                                                                                   // HU-20
}

// formulario envuelve las acciones POST de las pantallas: rechaza los envíos desde otro sitio
// (protección CSRF sin token, alcanza porque no hay sesiones) y limita el tamaño del cuerpo.
func (p *pantallas) formulario(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !mismoOrigen(r) {
			p.renderizar(w, http.StatusForbidden, "error.html", datosError{"el formulario se envió desde otro sitio"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoFormulario)
		if err := r.ParseForm(); err != nil {
			p.renderizar(w, http.StatusRequestEntityTooLarge, "error.html", datosError{"el formulario es demasiado grande"})
			return
		}
		h(w, r)
	}
}

// mismoOrigen indica si el pedido viene de una página de este mismo servidor. Los navegadores
// mandan Sec-Fetch-Site u Origin en los POST; un pedido sin ninguno de los dos (curl, tests) se acepta.
func mismoOrigen(r *http.Request) bool {
	switch r.Header.Get("Sec-Fetch-Site") {
	case "same-origin", "none":
		return true
	case "":
	default:
		return false
	}
	origen := r.Header.Get("Origin")
	if origen == "" {
		return true
	}
	u, err := url.Parse(origen)
	return err == nil && u.Host == r.Host
}

// --- Proyectos ---

func (p *pantallas) lista(w http.ResponseWriter, r *http.Request) {
	proyectos, err := p.cu.ListarProyectos.Ejecutar(r.Context())
	if err != nil {
		p.fallaInterna(w, "listar proyectos", err)
		return
	}
	datos := datosLista{}
	for _, pr := range proyectos {
		datos.Proyectos = append(datos.Proyectos, aVistaProyecto(pr))
	}
	p.renderizar(w, http.StatusOK, "proyectos.html", datos)
}

func (p *pantallas) formNuevo(w http.ResponseWriter, _ *http.Request) {
	p.renderizar(w, http.StatusOK, "proyecto_form.html", datosFormProyecto{})
}

func (p *pantallas) crear(w http.ResponseWriter, r *http.Request) {
	form := leerFormProyecto(r)
	pr, err := crearDesdeForm(r.Context(), p.cu.CrearProyecto, form)
	if err != nil {
		if mensaje, ok := mensajeDeValidacion(err); ok {
			p.renderizar(w, http.StatusBadRequest, "proyecto_form.html", datosFormProyecto{Form: form, Error: mensaje})
			return
		}
		p.fallaInterna(w, "crear proyecto", err)
		return
	}
	redirigirAlDetalle(w, r, pr.ID)
}

func crearDesdeForm(ctx context.Context, uc creadorProyectos, f formProyecto) (domain.Proyecto, error) {
	inicio, err := fechaDeForm(f.FechaInicio, "fecha_inicio")
	if err != nil {
		return domain.Proyecto{}, err
	}
	fin, err := fechaDeForm(f.FechaFin, "fecha_fin")
	if err != nil {
		return domain.Proyecto{}, err
	}
	d := app.DatosNuevoProyecto{Nombre: f.Nombre, FechaFin: fin}
	if inicio != nil {
		d.FechaInicio = *inicio
	}
	return uc.Ejecutar(ctx, d)
}

func (p *pantallas) detalle(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	p.mostrarDetalle(w, r, id, http.StatusOK, "", formIntegrante{})
}

// mostrarDetalle renderiza el detalle del proyecto, con el mensaje de error de una acción sobre
// sus integrantes si lo hay.
func (p *pantallas) mostrarDetalle(w http.ResponseWriter, r *http.Request, id int64, status int, mensaje string, f formIntegrante) {
	pr, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	p.renderizar(w, status, "proyecto_detalle.html", datosDetalle{Proyecto: aVistaProyecto(pr), Error: mensaje, FormIntegrante: f})
}

func (p *pantallas) formEditar(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	pr, err := p.cu.ConsultarProyecto.Ejecutar(r.Context(), id)
	if err != nil {
		p.errorDeConsulta(w, err)
		return
	}
	v := aVistaProyecto(pr)
	form := formProyecto{Nombre: v.Nombre, FechaInicio: v.FechaInicio, FechaFin: v.FechaFin, Estado: v.Estado}
	p.renderizar(w, http.StatusOK, "proyecto_form.html",
		datosFormProyecto{ID: id, Form: form, Original: form, Estados: estadosProyecto})
}

func (p *pantallas) editar(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	form := leerFormProyecto(r)
	original, conOriginal := leerOriginalProyecto(r)
	cambios, err := cambiosDeForm(form, original, conOriginal)
	if err == nil && conOriginal && cambios == (domain.CambiosProyecto{}) {
		// RN-30.3: si no cambió ningún campo, se vuelve al detalle sin mensaje (no es un error).
		redirigirAlDetalle(w, r, id)
		return
	}
	if err == nil {
		_, err = p.cu.ModificarProyecto.Ejecutar(r.Context(), app.DatosModificarProyecto{ProyectoID: id, Cambios: cambios})
	}
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			p.errorDeConsulta(w, err)
			return
		}
		if mensaje, ok := mensajeDeValidacion(err); ok {
			p.renderizar(w, http.StatusBadRequest, "proyecto_form.html",
				datosFormProyecto{ID: id, Form: form, Original: original, Error: mensaje, Estados: estadosProyecto})
			return
		}
		p.fallaInterna(w, "modificar proyecto", err)
		return
	}
	redirigirAlDetalle(w, r, id)
}

// cambiosDeForm arma la modificación con los campos que el usuario cambió respecto de los valores
// con los que se abrió el formulario (RN-30.3, DEF-01): así dos ediciones simultáneas de campos
// distintos no se pisan, igual que con la API (HU-02, sección 5). Si el pedido no trae los valores
// originales (un cliente que no es el formulario), se envían todos los campos.
// Una fecha de fin vacía la quita; una de inicio vacía la rechaza el dominio (RN-02.4).
func cambiosDeForm(f, original formProyecto, conOriginal bool) (domain.CambiosProyecto, error) {
	cambio := func(nuevo, viejo string) bool { return !conOriginal || nuevo != viejo }
	var c domain.CambiosProyecto
	if cambio(f.Nombre, original.Nombre) {
		nombre := f.Nombre
		c.Nombre = &nombre
	}
	if cambio(f.FechaInicio, original.FechaInicio) {
		inicio, err := fechaDeForm(f.FechaInicio, "fecha_inicio")
		if err != nil {
			return domain.CambiosProyecto{}, err
		}
		if inicio == nil {
			inicio = &time.Time{}
		}
		c.FechaInicio = inicio
	}
	if cambio(f.FechaFin, original.FechaFin) {
		fin, err := fechaDeForm(f.FechaFin, "fecha_fin")
		if err != nil {
			return domain.CambiosProyecto{}, err
		}
		c.CambiarFechaFin, c.FechaFin = true, fin
	}
	if cambio(f.Estado, original.Estado) {
		estado := domain.Estado(f.Estado)
		c.Estado = &estado
	}
	return c, nil
}

// --- Integrantes ---

func (p *pantallas) agregarIntegrante(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	f := formIntegrante{Nombre: r.PostFormValue("nombre"), Rol: r.PostFormValue("rol")}
	err := p.cu.AgregarIntegrante.Ejecutar(r.Context(), app.DatosIntegrante{ProyectoID: id, Nombre: f.Nombre, Rol: f.Rol})
	p.resultadoIntegrante(w, r, id, err, f, "agregar integrante")
}

func (p *pantallas) quitarIntegrante(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	err := p.cu.QuitarIntegrante.Ejecutar(r.Context(), app.DatosQuitarIntegrante{ProyectoID: id, Nombre: r.PathValue("nombre")})
	p.resultadoIntegrante(w, r, id, err, formIntegrante{}, "quitar integrante")
}

func (p *pantallas) agregarRol(w http.ResponseWriter, r *http.Request) {
	id, ok := p.idDeRuta(w, r)
	if !ok {
		return
	}
	err := p.cu.AgregarRol.Ejecutar(r.Context(),
		app.DatosRol{ProyectoID: id, Nombre: r.PathValue("nombre"), Rol: r.PostFormValue("rol")})
	p.resultadoIntegrante(w, r, id, err, formIntegrante{}, "sumar rol")
}

// resultadoIntegrante redirige al detalle si la acción salió bien, o vuelve a mostrarlo con el error.
func (p *pantallas) resultadoIntegrante(w http.ResponseWriter, r *http.Request, id int64, err error, f formIntegrante, accion string) {
	if err == nil {
		redirigirAlDetalle(w, r, id)
		return
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		p.errorDeConsulta(w, err)
		return
	}
	if mensaje, ok := mensajeDeValidacion(err); ok {
		p.mostrarDetalle(w, r, id, http.StatusBadRequest, mensaje, f)
		return
	}
	p.fallaInterna(w, accion, err)
}

// --- Ayudantes ---

// renderizar arma la página en memoria antes de escribirla: si la plantilla falla, se responde
// 500 en lugar de una página cortada con el estado original.
func (p *pantallas) renderizar(w http.ResponseWriter, status int, pagina string, datos any) {
	var buf bytes.Buffer
	if err := p.paginas[pagina].ExecuteTemplate(&buf, "layout", datos); err != nil {
		log.Printf("renderizar %s: %v", pagina, err)
		http.Error(w, "error interno al mostrar la página", http.StatusInternalServerError)
		return
	}
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'self'")
	h.Set("Referrer-Policy", "same-origin")
	w.WriteHeader(status)
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("escribir %s: %v", pagina, err)
	}
}

func (p *pantallas) idDeRuta(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := proyectoIDDesdeRuta(r)
	if err != nil {
		p.renderizar(w, http.StatusBadRequest, "error.html", datosError{err.Error()})
		return 0, false
	}
	return id, true
}

// errorDeConsulta muestra la página 404 si el proyecto no existe, o la de falla interna.
func (p *pantallas) errorDeConsulta(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		p.renderizar(w, http.StatusNotFound, "error.html", datosError{domain.ErrProyectoNoEncontrado.Error()})
		return
	}
	p.fallaInterna(w, "consultar proyecto", err)
}

func (p *pantallas) fallaInterna(w http.ResponseWriter, accion string, err error) {
	log.Printf("pantalla, %s: %v", accion, err)
	p.renderizar(w, http.StatusInternalServerError, "error.html", datosError{"Ocurrió un error interno. Probá de nuevo más tarde."})
}

// mensajeDeValidacion devuelve el mensaje a mostrar si err es un error de validación conocido.
func mensajeDeValidacion(err error) (string, bool) {
	var ef errFormulario
	if errors.As(err, &ef) {
		return ef.mensaje, true
	}
	for _, e := range erroresDeValidacionPantallas {
		if errors.Is(err, e) {
			return e.Error(), true // el texto del error de dominio, no el envuelto por el caso de uso
		}
	}
	return "", false
}

func leerFormProyecto(r *http.Request) formProyecto {
	return formProyecto{
		Nombre:      r.PostFormValue("nombre"),
		FechaInicio: r.PostFormValue("fecha_inicio"),
		FechaFin:    r.PostFormValue("fecha_fin"),
		Estado:      r.PostFormValue("estado"),
	}
}

// leerOriginalProyecto lee los valores con los que se abrió el formulario de edición (campos
// ocultos original_*). Devuelve false si el pedido no los trae.
func leerOriginalProyecto(r *http.Request) (formProyecto, bool) {
	if _, ok := r.PostForm["original_nombre"]; !ok {
		return formProyecto{}, false
	}
	return formProyecto{
		Nombre:      r.PostFormValue("original_nombre"),
		FechaInicio: r.PostFormValue("original_fecha_inicio"),
		FechaFin:    r.PostFormValue("original_fecha_fin"),
		Estado:      r.PostFormValue("original_estado"),
	}, true
}

// fechaDeForm convierte una fecha AAAA-MM-DD del formulario; vacía devuelve nil.
func fechaDeForm(valor, campo string) (*time.Time, error) {
	if valor == "" {
		return nil, nil
	}
	f, err := time.Parse(time.DateOnly, valor)
	if err != nil {
		return nil, errFormulario{campo + " debe tener el formato AAAA-MM-DD"}
	}
	return &f, nil
}

func aVistaProyecto(p domain.Proyecto) vistaProyecto {
	v := vistaProyecto{
		ID:          p.ID,
		Nombre:      p.Nombre,
		FechaInicio: p.FechaInicio.Format(time.DateOnly),
		Estado:      string(p.Estado),
		Integrantes: p.Integrantes,
	}
	if p.FechaFin != nil {
		v.FechaFin = p.FechaFin.Format(time.DateOnly)
	}
	return v
}

func redirigirAlDetalle(w http.ResponseWriter, r *http.Request, id int64) {
	http.Redirect(w, r, fmt.Sprintf("/web/proyectos/%d", id), http.StatusSeeOther)
}
