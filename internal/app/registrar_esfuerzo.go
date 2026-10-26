package app

import (
	"context"
	"fmt"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// LectorElementos es el puerto para buscar un elemento del backlog por su ID (HU-18).
type LectorElementos interface {
	// BuscarElementoPorID devuelve domain.ErrElementoNoEncontrado si no existe.
	BuscarElementoPorID(ctx context.Context, id int64) (domain.ElementoBacklog, error)
}

// RepositorioEsfuerzo es el puerto para guardar los registros de esfuerzo (HU-18).
type RepositorioEsfuerzo interface {
	// GuardarRegistro persiste el registro y le asigna su ID (RN-18.10).
	GuardarRegistro(ctx context.Context, r *domain.RegistroEsfuerzo) error
}

// DatosRegistroEsfuerzo son las entradas de HU-18.
type DatosRegistroEsfuerzo struct {
	ProyectoID int64
	ElementoID int64
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
}

// RegistrarEsfuerzo es el caso de uso de HU-18. Lee el proyecto con el mismo puerto de solo
// lectura que HU-04 (LectorProyectos): lo único que necesita es buscarlo con sus integrantes.
type RegistrarEsfuerzo struct {
	proyectos LectorProyectos
	elementos LectorElementos
	registros RepositorioEsfuerzo
	reloj     func() time.Time
}

// NuevoRegistrarEsfuerzo recibe el reloj para que "hoy" se pueda fijar en los tests (spec, sección 5).
func NuevoRegistrarEsfuerzo(proyectos LectorProyectos, elementos LectorElementos,
	registros RepositorioEsfuerzo, reloj func() time.Time) *RegistrarEsfuerzo {
	return &RegistrarEsfuerzo{proyectos: proyectos, elementos: elementos, registros: registros, reloj: reloj}
}

// Ejecutar busca el proyecto y el elemento, valida el registro en el dominio, verifica que el
// integrante pertenezca al proyecto y lo guarda. Los errores salen en el orden de la sección 7
// y envueltos con %w, así errors.Is sigue reconociendo los de la spec. Si algo falla, no se guarda nada.
// Implementa HU-18 (RN-18.1, RN-18.3, RN-18.10).
func (uc *RegistrarEsfuerzo) Ejecutar(ctx context.Context, d DatosRegistroEsfuerzo) (domain.RegistroEsfuerzo, error) {
	p, err := uc.proyectos.BuscarPorID(ctx, d.ProyectoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("buscar proyecto %d: %w", d.ProyectoID, err)
	}

	e, err := uc.elementos.BuscarElementoPorID(ctx, d.ElementoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("buscar elemento %d: %w", d.ElementoID, err)
	}
	if e.ProyectoID != p.ID { // RN-18.1: un elemento de otro proyecto se trata como inexistente
		return domain.RegistroEsfuerzo{}, fmt.Errorf("elemento %d en el proyecto %d: %w", d.ElementoID, p.ID, domain.ErrElementoNoEncontrado)
	}

	r, err := domain.NuevoRegistroEsfuerzo(p.ID, e.ID, d.Integrante, d.Fecha, d.Actividad, d.Horas, uc.reloj())
	if err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("registrar esfuerzo: %w", err)
	}

	nombre, ok := nombreEnProyecto(p, r.Integrante)
	if !ok {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("integrante %q: %w", r.Integrante, domain.ErrIntegranteNoEncontrado)
	}
	r.Integrante = nombre

	if err := uc.registros.GuardarRegistro(ctx, &r); err != nil {
		return domain.RegistroEsfuerzo{}, fmt.Errorf("guardar registro de esfuerzo: %w", err)
	}
	return r, nil
}

// nombreEnProyecto busca al integrante con domain.ClaveNombre, como HU-03, y devuelve su nombre
// tal como figura en el proyecto (" ana " → "Ana"). Implementa HU-18 (RN-18.3).
func nombreEnProyecto(p domain.Proyecto, nombre string) (string, bool) {
	clave := domain.ClaveNombre(nombre)
	for _, i := range p.Integrantes {
		if domain.ClaveNombre(i.Nombre) == clave {
			return i.Nombre, true
		}
	}
	return "", false
}
