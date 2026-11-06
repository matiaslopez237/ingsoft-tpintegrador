// Package ejemplo carga, al arrancar el servidor, un proyecto de ejemplo completo: este mismo TP.
// Es un adaptador que maneja la aplicación desde afuera, como el HTTP: usa los casos de uso de
// internal/app y no conoce la persistencia. Implementa HU-31.
package ejemplo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Resultado dice si se cargó el ejemplo o por qué no (sección 3 de la spec: se informa en el log).
type Resultado int

const (
	Cargado          Resultado = iota + 1 // se cargó el ejemplo
	OpcionApagada                         // CARGAR_EJEMPLO ausente, vacía o falsa (RN-31.1)
	ValorInvalido                         // CARGAR_EJEMPLO no es un booleano (RN-31.3)
	BaseConProyectos                      // la base ya tiene algún proyecto (RN-31.1)
)

// Reloj es el "hoy" que la carga les da a los casos de uso que dependen de la fecha (RN-31.5). Quien
// arma esos casos de uso para la carga les pasa Ahora como reloj.
type Reloj struct {
	hoy time.Time
}

// Ahora devuelve el día que la carga fijó por última vez.
func (r *Reloj) Ahora() time.Time {
	return r.hoy
}

// CasosDeUso son los casos de uso que usa la carga (RN-31.2). RegistrarEsfuerzo, RegistrarDefecto
// y ActualizarEstadoDefecto tienen que estar armados con Reloj.Ahora como reloj (RN-31.5).
type CasosDeUso struct {
	Reloj                   *Reloj
	ListarProyectos         listadorProyectos
	CrearProyecto           creadorProyecto
	AgregarIntegrante       agregadorIntegrante
	CrearElemento           creadorElemento
	EditarElemento          editorElemento
	CrearSprint             creadorSprint
	AsignarHistoria         asignadorHistoria
	CompletarHistoria       completadorHistoria
	CerrarSprint            cerradorSprint
	RegistrarEsfuerzo       registradorEsfuerzo
	RegistrarDefecto        registradorDefecto
	ActualizarEstadoDefecto actualizadorDefecto
	IniciarRonda            iniciadorRonda
	Votar                   votante
	RegistrarEstimacion     registradorEstimacion
}

type listadorProyectos interface {
	Ejecutar(ctx context.Context) ([]domain.Proyecto, error)
}

type creadorProyecto interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoProyecto) (domain.Proyecto, error)
}

type agregadorIntegrante interface {
	Ejecutar(ctx context.Context, d app.DatosIntegrante) error
}

type creadorElemento interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoElemento) (domain.ElementoBacklog, error)
}

type editorElemento interface {
	Ejecutar(ctx context.Context, d app.DatosEditarElemento) (domain.ElementoBacklog, error)
}

type creadorSprint interface {
	Ejecutar(ctx context.Context, d app.DatosNuevoSprint) (domain.Sprint, error)
}

type asignadorHistoria interface {
	Ejecutar(ctx context.Context, d app.DatosAsignacion) error
}

type completadorHistoria interface {
	Ejecutar(ctx context.Context, d app.DatosCompletar) error
}

type cerradorSprint interface {
	Ejecutar(ctx context.Context, d app.DatosCerrarSprint) error
}

type registradorEsfuerzo interface {
	Ejecutar(ctx context.Context, d app.DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error)
}

type registradorDefecto interface {
	Ejecutar(ctx context.Context, d app.DatosDefecto) (domain.Defecto, error)
}

type actualizadorDefecto interface {
	Iniciar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
	Resolver(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
	Cerrar(ctx context.Context, d app.DatosActualizarEstadoDefecto) (domain.Defecto, error)
}

type iniciadorRonda interface {
	Ejecutar(ctx context.Context, d app.DatosIniciarRonda) (domain.ResumenRonda, error)
}

type votante interface {
	Ejecutar(ctx context.Context, d app.DatosVoto) (domain.ResumenRonda, error)
}

type registradorEstimacion interface {
	Ejecutar(ctx context.Context, d app.DatosEstimacion) (domain.ResumenRonda, error)
}

// CargarSiCorresponde carga el ejemplo si valor (el de CARGAR_EJEMPLO) está activado y la base no
// tiene proyectos (RN-31.1, RN-31.3). Implementa HU-31.
func CargarSiCorresponde(ctx context.Context, cu CasosDeUso, valor string) (Resultado, error) {
	if strings.TrimSpace(valor) == "" {
		return OpcionApagada, nil // ausente o vacía equivale a false (sección 2)
	}
	activada, err := strconv.ParseBool(strings.TrimSpace(valor))
	if err != nil {
		return ValorInvalido, nil // RN-31.3: no es un error, solo no se carga
	}
	if !activada {
		return OpcionApagada, nil
	}
	proyectos, err := cu.ListarProyectos.Ejecutar(ctx)
	if err != nil {
		return 0, fmt.Errorf("cargar ejemplo: listar los proyectos: %w", err)
	}
	if len(proyectos) > 0 {
		return BaseConProyectos, nil
	}
	if err := cargar(ctx, cu, contenidoTP()); err != nil {
		return 0, err
	}
	return Cargado, nil
}
