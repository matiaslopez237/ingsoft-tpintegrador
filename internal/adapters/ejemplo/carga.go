package ejemplo

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// contenido es lo que carga el ejemplo: el proyecto, su equipo, su backlog y los pasos que pasaron,
// en orden cronológico (RN-31.4).
type contenido struct {
	proyecto    string
	inicio      time.Time
	integrantes []integrante
	historias   []historia
	pasos       []paso
}

type integrante struct {
	nombre, rol string
}

// historia es un elemento del backlog. Su título es "<id>: <titulo>", así se la encuentra por su ID.
type historia struct {
	id, titulo, descripcion string
	prioridad               domain.Prioridad
	puntos                  int
	criterios               []string
	horasEstimadas          float64 // 0: sin estimar en horas
}

// paso es algo que pasó en el proyecto en un día dado. La carga fija el reloj en ese día antes de
// aplicarlo (RN-31.5), y si falla, el error lo nombra con String (sección 7).
type paso interface {
	dia() time.Time
	aplicar(ctx context.Context, c *carga) error
	String() string
}

// carga guarda lo que se va creando, para que los pasos se refieran a las cosas por su nombre.
type carga struct {
	cu         CasosDeUso
	proyectoID int64
	historias  map[string]int64 // por ID de historia ("HU-01")
	sprints    map[int]int64    // por número de sprint
	defectos   map[string]int64 // por clave ("DEF-01")
}

// cargar crea el proyecto, su equipo y su backlog, y aplica los pasos en orden. Todo pasa por los
// casos de uso (RN-31.2); el primer error corta la carga y nombra el paso que falló (sección 7).
func cargar(ctx context.Context, cu CasosDeUso, c contenido) error {
	cu.Reloj.hoy = c.inicio
	p, err := cu.CrearProyecto.Ejecutar(ctx, app.DatosNuevoProyecto{Nombre: c.proyecto, FechaInicio: c.inicio})
	if err != nil {
		return fallo("crear el proyecto", err)
	}
	k := &carga{cu: cu, proyectoID: p.ID, historias: map[string]int64{}, sprints: map[int]int64{}, defectos: map[string]int64{}}

	for _, i := range c.integrantes {
		if err := cu.AgregarIntegrante.Ejecutar(ctx, app.DatosIntegrante{ProyectoID: p.ID, Nombre: i.nombre, Rol: i.rol}); err != nil {
			return fallo("agregar a "+i.nombre, err)
		}
	}
	for _, h := range c.historias {
		if err := k.crearHistoria(ctx, h); err != nil {
			return fallo("crear "+h.id, err)
		}
	}
	for _, paso := range c.pasos {
		cu.Reloj.hoy = paso.dia()
		if err := paso.aplicar(ctx, k); err != nil {
			return fallo(paso.String(), err)
		}
	}
	return nil
}

func fallo(paso string, err error) error {
	return fmt.Errorf("cargar ejemplo: %s: %w", paso, err)
}

// crearHistoria crea el elemento y, si tiene horas estimadas, las carga con la edición, porque el
// alta no las recibe (RN-19.3).
func (k *carga) crearHistoria(ctx context.Context, h historia) error {
	puntos := h.puntos
	e, err := k.cu.CrearElemento.Ejecutar(ctx, app.DatosNuevoElemento{
		ProyectoID: k.proyectoID, Titulo: h.id + ": " + h.titulo, Descripcion: h.descripcion,
		Prioridad: h.prioridad, StoryPoints: &puntos, Criterios: h.criterios,
	})
	if err != nil {
		return err
	}
	k.historias[h.id] = e.ID
	if h.horasEstimadas == 0 {
		return nil
	}
	horas := h.horasEstimadas
	_, err = k.cu.EditarElemento.Ejecutar(ctx, app.DatosEditarElemento{
		ProyectoID: k.proyectoID, ElementoID: e.ID,
		Cambios: domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: &horas},
	})
	return err
}

type crearSprint struct {
	numero      int
	inicio, fin time.Time
	objetivo    string
}

type asignar struct {
	el       time.Time
	historia string
	sprint   int
}

type completar struct {
	el       time.Time
	historia string
}

type cerrarSprint struct {
	el     time.Time
	sprint int
}

type esfuerzo struct {
	el                              time.Time
	historia, integrante, actividad string
	horas                           float64
}

type defecto struct {
	el                           time.Time
	clave, historia, descripcion string
	severidad                    domain.Severidad
}

// cambioDefecto pasa un defecto al estado siguiente: "iniciar", "resolver" o "cerrar" (HU-21).
type cambioDefecto struct {
	el            time.Time
	clave, accion string
}

// ronda abre una ronda de Planning Poker, registra los votos de todos los integrantes y, si
// acordada no es 0, la estimación acordada (HU-14 a HU-17).
type ronda struct {
	el       time.Time
	historia string
	votos    []voto
	acordada int
}

type voto struct {
	integrante string
	valor      int
}

func (p crearSprint) dia() time.Time   { return p.inicio }
func (p asignar) dia() time.Time       { return p.el }
func (p completar) dia() time.Time     { return p.el }
func (p cerrarSprint) dia() time.Time  { return p.el }
func (p esfuerzo) dia() time.Time      { return p.el }
func (p defecto) dia() time.Time       { return p.el }
func (p cambioDefecto) dia() time.Time { return p.el }
func (p ronda) dia() time.Time         { return p.el }

func (p crearSprint) aplicar(ctx context.Context, k *carga) error {
	s, err := k.cu.CrearSprint.Ejecutar(ctx, app.DatosNuevoSprint{ProyectoID: k.proyectoID, FechaInicio: p.inicio, FechaFin: p.fin, Objetivo: p.objetivo})
	if err != nil {
		return err
	}
	if s.Numero != p.numero {
		return fmt.Errorf("se creó el sprint %d", s.Numero) // el contenido está desordenado
	}
	k.sprints[p.numero] = s.ID
	return nil
}

func (p asignar) aplicar(ctx context.Context, k *carga) error {
	return k.cu.AsignarHistoria.Ejecutar(ctx, app.DatosAsignacion{ProyectoID: k.proyectoID, ElementoID: k.historias[p.historia], SprintID: k.sprints[p.sprint]})
}

func (p completar) aplicar(ctx context.Context, k *carga) error {
	return k.cu.CompletarHistoria.Ejecutar(ctx, app.DatosCompletar{ProyectoID: k.proyectoID, ElementoID: k.historias[p.historia]})
}

func (p cerrarSprint) aplicar(ctx context.Context, k *carga) error {
	return k.cu.CerrarSprint.Ejecutar(ctx, app.DatosCerrarSprint{ProyectoID: k.proyectoID, SprintID: k.sprints[p.sprint]})
}

func (p esfuerzo) aplicar(ctx context.Context, k *carga) error {
	_, err := k.cu.RegistrarEsfuerzo.Ejecutar(ctx, app.DatosRegistroEsfuerzo{
		ProyectoID: k.proyectoID, ElementoID: k.historias[p.historia], Integrante: p.integrante,
		Fecha: p.el, Actividad: p.actividad, Horas: p.horas,
	})
	return err
}

func (p defecto) aplicar(ctx context.Context, k *carga) error {
	d, err := k.cu.RegistrarDefecto.Ejecutar(ctx, app.DatosDefecto{
		ProyectoID: k.proyectoID, ElementoID: k.historias[p.historia],
		Descripcion: p.clave + ": " + p.descripcion, Severidad: p.severidad,
	})
	if err != nil {
		return err
	}
	k.defectos[p.clave] = d.ID
	return nil
}

func (p cambioDefecto) aplicar(ctx context.Context, k *carga) error {
	d := app.DatosActualizarEstadoDefecto{ProyectoID: k.proyectoID, DefectoID: k.defectos[p.clave]}
	var err error
	switch p.accion {
	case "iniciar":
		_, err = k.cu.ActualizarEstadoDefecto.Iniciar(ctx, d)
	case "resolver":
		_, err = k.cu.ActualizarEstadoDefecto.Resolver(ctx, d)
	case "cerrar":
		_, err = k.cu.ActualizarEstadoDefecto.Cerrar(ctx, d)
	default:
		err = fmt.Errorf("acción desconocida %q", p.accion)
	}
	return err
}

func (p ronda) aplicar(ctx context.Context, k *carga) error {
	r, err := k.cu.IniciarRonda.Ejecutar(ctx, app.DatosIniciarRonda{ProyectoID: k.proyectoID, ElementoID: k.historias[p.historia]})
	if err != nil {
		return err
	}
	for _, v := range p.votos {
		if _, err := k.cu.Votar.Ejecutar(ctx, app.DatosVoto{ProyectoID: k.proyectoID, RondaID: r.ID, Integrante: v.integrante, Valor: v.valor}); err != nil {
			return fmt.Errorf("voto de %s: %w", v.integrante, err)
		}
	}
	if p.acordada == 0 {
		return nil
	}
	_, err = k.cu.RegistrarEstimacion.Ejecutar(ctx, app.DatosEstimacion{ProyectoID: k.proyectoID, RondaID: r.ID, Valor: p.acordada})
	return err
}

func (p crearSprint) String() string  { return fmt.Sprintf("crear el sprint %d", p.numero) }
func (p asignar) String() string      { return fmt.Sprintf("asignar %s al sprint %d", p.historia, p.sprint) }
func (p completar) String() string    { return "completar " + p.historia }
func (p cerrarSprint) String() string { return fmt.Sprintf("cerrar el sprint %d", p.sprint) }
func (p esfuerzo) String() string {
	return fmt.Sprintf("registrar esfuerzo de %s en %s", p.integrante, p.historia)
}
func (p defecto) String() string       { return "registrar " + p.clave }
func (p cambioDefecto) String() string { return p.accion + " " + p.clave }
func (p ronda) String() string         { return "votar una ronda de " + p.historia }
