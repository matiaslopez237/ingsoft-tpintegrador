package bdd

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// respuestaWeb es lo que los escenarios de HU-30 observan de un pedido a las pantallas.
type respuestaWeb struct {
	Estado    int
	Ubicacion string // encabezado Location de una redirección
	Cuerpo    string
}

// proyectoGuardado es lo que los escenarios de HU-30 verifican en la base después de editar.
type proyectoGuardado struct {
	Nombre   string
	Estado   string
	TieneFin bool
}

// sistemaWeb es la puerta de entrada de los escenarios a las pantallas, sobre un servidor HTTP
// real. Los datos de partida se cargan con los casos de uso. La implementación está en
// hu30_sistema_test.go.
type sistemaWeb interface {
	GuardarProyecto(nombre string) (int64, error)
	AgregarIntegrante(proyectoID int64, nombre, rol string) error
	// Pedir hace el pedido sin seguir redirecciones; form en nil es un pedido sin cuerpo.
	Pedir(metodo, ruta string, form url.Values) (respuestaWeb, error)
	ContarProyectos() (int, error)
	BuscarProyecto(id int64) (proyectoGuardado, error)
	Cerrar() error
}

var rutaDetalle = regexp.MustCompile(`^/web/proyectos/(\d+)$`)

type webWorld struct {
	sistema   sistemaWeb
	proyectos map[string]int64
	ultima    respuestaWeb
}

func registrarStepsHU30(ctx *godog.ScenarioContext) {
	w := &webWorld{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		sistema, err := nuevoSistemaWeb()
		if err != nil {
			return c, err
		}
		*w = webWorld{sistema: sistema, proyectos: map[string]int64{}}
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return c, w.sistema.Cerrar()
	})

	ctx.Step(`^están guardados los proyectos:$`, w.estanGuardadosLosProyectos)
	ctx.Step(`^está guardado el proyecto "([^"]*)"$`, w.estaGuardadoElProyecto)
	ctx.Step(`^está guardado el proyecto "([^"]*)" con el integrante "([^"]*)" con rol "([^"]*)"$`, w.estaGuardadoConIntegrante)

	ctx.Step(`^abro la página "([^"]*)"$`, w.abroLaPagina)
	ctx.Step(`^abro el detalle del proyecto "([^"]*)"$`, w.abroElDetalle)
	ctx.Step(`^envío el formulario de nuevo proyecto con nombre "([^"]*)", inicio "([^"]*)" y fin "([^"]*)"$`, w.envioNuevoProyecto)
	ctx.Step(`^envío el formulario de edición de "([^"]*)" con nombre "([^"]*)", inicio "([^"]*)", fin "([^"]*)" y estado "([^"]*)"$`, w.envioEdicion)
	ctx.Step(`^agrego desde la web al integrante "([^"]*)" con rol "([^"]*)" en "([^"]*)"$`, w.agregoIntegrante)
	ctx.Step(`^quito desde la web al integrante "([^"]*)" de "([^"]*)"$`, w.quitoIntegrante)
	ctx.Step(`^le sumo desde la web el rol "([^"]*)" a "([^"]*)" en "([^"]*)"$`, w.sumoRol)

	ctx.Step(`^me redirige a "([^"]*)"$`, w.meRedirigeA)
	ctx.Step(`^me redirige al detalle del proyecto "([^"]*)"$`, w.meRedirigeAlDetalle)
	ctx.Step(`^la respuesta tiene estado (\d+)$`, w.respuestaTieneEstado)
	ctx.Step(`^la página tiene el layout común$`, w.tieneLayoutComun)
	ctx.Step(`^la hoja de estilos se sirve en "([^"]*)"$`, w.hojaDeEstilos)
	ctx.Step(`^veo el texto "([^"]*)"$`, w.paginaContiene)
	ctx.Step(`^la página contiene "([^"]*)"$`, w.paginaContiene)
	ctx.Step(`^la página no contiene "([^"]*)"$`, w.paginaNoContiene)
	ctx.Step(`^veo los proyectos en este orden:$`, w.veoEnOrden)
	ctx.Step(`^hay (\d+) proyectos guardados$`, w.hayProyectosGuardados)
	ctx.Step(`^el detalle de "([^"]*)" muestra al integrante "([^"]*)" con rol "([^"]*)"$`, w.detalleMuestraIntegrante)
	ctx.Step(`^el detalle de "([^"]*)" no muestra al integrante "([^"]*)"$`, w.detalleNoMuestraIntegrante)
	ctx.Step(`^el proyecto "([^"]*)" queda en estado "([^"]*)" y sin fecha de fin$`, w.quedaEnEstadoSinFin)
}

// --- Given ---

func (w *webWorld) estanGuardadosLosProyectos(tabla *godog.Table) error {
	for _, fila := range tabla.Rows[1:] { // la primera fila es el encabezado
		if err := w.estaGuardadoElProyecto(fila.Cells[0].Value); err != nil {
			return err
		}
	}
	return nil
}

func (w *webWorld) estaGuardadoElProyecto(nombre string) error {
	id, err := w.sistema.GuardarProyecto(nombre)
	if err != nil {
		return fmt.Errorf("preparar el proyecto %q: %w", nombre, err)
	}
	w.proyectos[nombre] = id
	return nil
}

func (w *webWorld) estaGuardadoConIntegrante(proyecto, integrante, rol string) error {
	if err := w.estaGuardadoElProyecto(proyecto); err != nil {
		return err
	}
	if err := w.sistema.AgregarIntegrante(w.proyectos[proyecto], integrante, rol); err != nil {
		return fmt.Errorf("preparar el integrante %q: %w", integrante, err)
	}
	return nil
}

// --- When ---

func (w *webWorld) pedir(metodo, ruta string, form url.Values) error {
	r, err := w.sistema.Pedir(metodo, ruta, form)
	if err != nil {
		return err
	}
	w.ultima = r
	return nil
}

func (w *webWorld) abroLaPagina(ruta string) error {
	return w.pedir("GET", ruta, nil)
}

func (w *webWorld) abroElDetalle(proyecto string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	return w.abroLaPagina(fmt.Sprintf("/web/proyectos/%d", id))
}

func (w *webWorld) envioNuevoProyecto(nombre, inicio, fin string) error {
	return w.pedir("POST", "/web/proyectos", url.Values{
		"nombre": {nombre}, "fecha_inicio": {inicio}, "fecha_fin": {fin},
	})
}

func (w *webWorld) envioEdicion(proyecto, nombre, inicio, fin, estado string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	return w.pedir("POST", fmt.Sprintf("/web/proyectos/%d/editar", id), url.Values{
		"nombre": {nombre}, "fecha_inicio": {inicio}, "fecha_fin": {fin}, "estado": {estado},
	})
}

func (w *webWorld) agregoIntegrante(nombre, rol, proyecto string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	return w.pedir("POST", fmt.Sprintf("/web/proyectos/%d/integrantes", id), url.Values{
		"nombre": {nombre}, "rol": {rol},
	})
}

func (w *webWorld) quitoIntegrante(nombre, proyecto string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	return w.pedir("POST", fmt.Sprintf("/web/proyectos/%d/integrantes/%s/quitar", id, url.PathEscape(nombre)), url.Values{})
}

func (w *webWorld) sumoRol(rol, nombre, proyecto string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	return w.pedir("POST", fmt.Sprintf("/web/proyectos/%d/integrantes/%s/roles", id, url.PathEscape(nombre)), url.Values{
		"rol": {rol},
	})
}

func (w *webWorld) idDe(proyecto string) (int64, error) {
	id, ok := w.proyectos[proyecto]
	if !ok {
		return 0, fmt.Errorf("el escenario no conoce el proyecto %q", proyecto)
	}
	return id, nil
}

// --- Then ---

func (w *webWorld) meRedirigeA(ruta string) error {
	if w.ultima.Estado != 303 && w.ultima.Estado != 302 {
		return fmt.Errorf("se esperaba una redirección, estado %d; cuerpo: %.300s", w.ultima.Estado, w.ultima.Cuerpo)
	}
	if w.ultima.Ubicacion != ruta {
		return fmt.Errorf("redirección: se esperaba %q, se obtuvo %q", ruta, w.ultima.Ubicacion)
	}
	return nil
}

// meRedirigeAlDetalle verifica la redirección al detalle, recuerda el ID del proyecto con ese
// nombre y abre el detalle para comprobar que lo muestra.
func (w *webWorld) meRedirigeAlDetalle(proyecto string) error {
	if w.ultima.Estado != 303 {
		return fmt.Errorf("se esperaba la redirección 303, estado %d; cuerpo: %.300s", w.ultima.Estado, w.ultima.Cuerpo)
	}
	m := rutaDetalle.FindStringSubmatch(w.ultima.Ubicacion)
	if m == nil {
		return fmt.Errorf("se esperaba una redirección al detalle, se obtuvo %q", w.ultima.Ubicacion)
	}
	id, _ := strconv.ParseInt(m[1], 10, 64)
	w.proyectos[proyecto] = id
	if err := w.abroLaPagina(w.ultima.Ubicacion); err != nil {
		return err
	}
	return w.paginaContiene(proyecto)
}

func (w *webWorld) respuestaTieneEstado(estado int) error {
	if w.ultima.Estado != estado {
		return fmt.Errorf("estado: se esperaba %d, se obtuvo %d; cuerpo: %.300s", estado, w.ultima.Estado, w.ultima.Cuerpo)
	}
	return nil
}

func (w *webWorld) tieneLayoutComun() error {
	for _, parte := range []string{"<header", "<nav", `href="/web/proyectos"`, `href="/static/estilos.css"`} {
		if err := w.paginaContiene(parte); err != nil {
			return fmt.Errorf("layout común: %w", err)
		}
	}
	return nil
}

func (w *webWorld) hojaDeEstilos(ruta string) error {
	r, err := w.sistema.Pedir("GET", ruta, nil)
	if err != nil {
		return err
	}
	if r.Estado != 200 || r.Cuerpo == "" {
		return fmt.Errorf("hoja de estilos %q: estado %d, %d bytes", ruta, r.Estado, len(r.Cuerpo))
	}
	return nil
}

func (w *webWorld) paginaContiene(texto string) error {
	if !strings.Contains(w.ultima.Cuerpo, texto) {
		return fmt.Errorf("la página no contiene %q; cuerpo: %.500s", texto, w.ultima.Cuerpo)
	}
	return nil
}

func (w *webWorld) paginaNoContiene(texto string) error {
	if strings.Contains(w.ultima.Cuerpo, texto) {
		return fmt.Errorf("la página no debería contener %q", texto)
	}
	return nil
}

func (w *webWorld) veoEnOrden(tabla *godog.Table) error {
	desde := 0
	for _, fila := range tabla.Rows[1:] {
		nombre := fila.Cells[0].Value
		i := strings.Index(w.ultima.Cuerpo[desde:], ">"+nombre+"<")
		if i < 0 {
			return fmt.Errorf("no aparece %q después de la posición %d (orden incorrecto o falta)", nombre, desde)
		}
		desde += i + len(nombre)
	}
	return nil
}

func (w *webWorld) hayProyectosGuardados(esperado int) error {
	n, err := w.sistema.ContarProyectos()
	if err != nil {
		return err
	}
	if n != esperado {
		return fmt.Errorf("proyectos guardados: se esperaban %d, hay %d", esperado, n)
	}
	return nil
}

func (w *webWorld) detalleMuestraIntegrante(proyecto, integrante, rol string) error {
	if err := w.abroElDetalle(proyecto); err != nil {
		return err
	}
	if err := w.paginaContiene(">" + integrante + "<"); err != nil {
		return err
	}
	return w.paginaContiene(rol)
}

func (w *webWorld) detalleNoMuestraIntegrante(proyecto, integrante string) error {
	if err := w.abroElDetalle(proyecto); err != nil {
		return err
	}
	return w.paginaNoContiene(">" + integrante + "<")
}

func (w *webWorld) quedaEnEstadoSinFin(proyecto, estado string) error {
	id, err := w.idDe(proyecto)
	if err != nil {
		return err
	}
	p, err := w.sistema.BuscarProyecto(id)
	if err != nil {
		return err
	}
	if p.Nombre != proyecto || p.Estado != estado || p.TieneFin {
		return fmt.Errorf("se esperaba %q en estado %q sin fecha de fin, se obtuvo %+v", proyecto, estado, p)
	}
	return nil
}
