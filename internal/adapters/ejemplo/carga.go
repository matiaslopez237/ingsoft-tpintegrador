package ejemplo

import (
	"context"
	"time"

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

// cargar crea el proyecto, su equipo y su backlog, y aplica los pasos en orden.
func cargar(ctx context.Context, cu CasosDeUso, c contenido) error {
	return nil
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

func (p crearSprint) aplicar(context.Context, *carga) error   { return nil }
func (p asignar) aplicar(context.Context, *carga) error       { return nil }
func (p completar) aplicar(context.Context, *carga) error     { return nil }
func (p cerrarSprint) aplicar(context.Context, *carga) error  { return nil }
func (p esfuerzo) aplicar(context.Context, *carga) error      { return nil }
func (p defecto) aplicar(context.Context, *carga) error       { return nil }
func (p cambioDefecto) aplicar(context.Context, *carga) error { return nil }
func (p ronda) aplicar(context.Context, *carga) error         { return nil }

func (p crearSprint) String() string   { return "" }
func (p asignar) String() string       { return "" }
func (p completar) String() string     { return "" }
func (p cerrarSprint) String() string  { return "" }
func (p esfuerzo) String() string      { return "" }
func (p defecto) String() string       { return "" }
func (p cambioDefecto) String() string { return "" }
func (p ronda) String() string         { return "" }
