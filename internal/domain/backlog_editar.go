package domain

import "errors"

// Error nuevo de HU-06 (spec, sección 7). Los de los campos son los de HU-05, y el de estado
// inválido (ErrEstadoElementoInvalido) se comparte con HU-07 (backlog_listado.go).
var ErrElementoConEsfuerzo = errors.New("no se puede eliminar un elemento con esfuerzo registrado")

// ErrElementoConDefectos: no se elimina un elemento con defectos registrados (HU-20). DEF-02, RN-06.10.
var ErrElementoConDefectos = errors.New("no se puede eliminar un elemento con defectos registrados")

// HorasEstimadasMaximas es el tope de las horas estimadas de un elemento (RN-19.2).
const HorasEstimadasMaximas = 1000

// ErrHorasEstimadasExcedenMaximo: HU-19 (spec, sección 7). Los demás errores de las horas estimadas
// son los de HU-18.
var ErrHorasEstimadasExcedenMaximo = errors.New("las horas estimadas no pueden superar las 1000 horas")

// validarHorasEstimadas aplica RN-19.2 en este orden: > 0, múltiplo de 0,5 y máximo 1000.
func validarHorasEstimadas(horas float64) error {
	if err := validarMediasHoras(horas); err != nil {
		return err
	}
	if horas > HorasEstimadasMaximas { // también +Inf, que sí pasa la comprobación del múltiplo
		return ErrHorasEstimadasExcedenMaximo
	}
	return nil
}

// CambiosElemento son los datos de una edición parcial de un elemento (HU-06, RN-06.1).
// Un campo en nil no se envía y conserva su valor.
type CambiosElemento struct {
	Titulo      *string
	Descripcion *string
	Prioridad   *Prioridad
	// CambiarStoryPoints indica que se envían los story points. Con StoryPoints en nil, el
	// elemento queda sin estimar (RN-06.3).
	CambiarStoryPoints bool
	StoryPoints        *int
	// Criterios, si se envía, reemplaza la lista completa.
	Criterios *[]string
	Estado    *EstadoElemento
	// CambiarHorasEstimadas indica que se envían las horas estimadas. Con HorasEstimadas en nil, el
	// elemento queda sin estimar en horas (RN-19.3).
	CambiarHorasEstimadas bool
	HorasEstimadas        *float64
}

func (c CambiosElemento) vacios() bool {
	return c.Titulo == nil && c.Descripcion == nil && c.Prioridad == nil && !c.CambiarStoryPoints &&
		c.Criterios == nil && c.Estado == nil && !c.CambiarHorasEstimadas
}

// Editar devuelve una copia del elemento con los cambios aplicados; el original no cambia.
// Valida el resultado con las mismas reglas que al crear (NuevoElementoBacklog), después el estado y
// por último las horas estimadas, si se enviaron. La unicidad del título se verifica en el caso de
// uso, como en HU-05. Implementa HU-06 (RN-06.1 a RN-06.5, RN-06.9) y HU-19 (RN-19.2, RN-19.3).
func (e ElementoBacklog) Editar(c CambiosElemento) (ElementoBacklog, error) {
	if c.vacios() {
		return ElementoBacklog{}, ErrSinCambios
	}

	titulo, descripcion, prioridad := e.Titulo, e.Descripcion, e.Prioridad
	storyPoints, criterios, estado := e.StoryPoints, e.Criterios, e.Estado
	if c.Titulo != nil {
		titulo = *c.Titulo
	}
	if c.Descripcion != nil {
		descripcion = *c.Descripcion
	}
	if c.Prioridad != nil {
		prioridad = *c.Prioridad
	}
	if c.CambiarStoryPoints {
		storyPoints = c.StoryPoints
	}
	if c.Criterios != nil {
		criterios = *c.Criterios
	}
	if c.Estado != nil {
		estado = *c.Estado
	}

	// NuevoElementoBacklog valida y copia los story points y los criterios, así la copia
	// editada no comparte memoria con el original ni con c. Los story points se validan solo si se
	// envían: un valor guardado antes de la escala de HU-13 se respeta (RN-13.3).
	aValidar := storyPoints
	if !c.CambiarStoryPoints {
		aValidar = nil
	}
	editado, err := NuevoElementoBacklog(e.ProyectoID, titulo, descripcion, prioridad, aValidar, criterios)
	if err != nil {
		return ElementoBacklog{}, err
	}
	if !c.CambiarStoryPoints && storyPoints != nil {
		sp := *storyPoints
		editado.StoryPoints = &sp
	}
	if !estadoElementoValido(estado) {
		return ElementoBacklog{}, ErrEstadoElementoInvalido
	}
	horasEstimadas := e.HorasEstimadas
	if c.CambiarHorasEstimadas {
		horasEstimadas = c.HorasEstimadas
		if horasEstimadas != nil {
			if err := validarHorasEstimadas(*horasEstimadas); err != nil {
				return ElementoBacklog{}, err
			}
		}
	}
	editado.ID = e.ID
	editado.Estado = estado
	// NuevoElementoBacklog no recibe el sprint ni las horas estimadas: sin esto, la copia los perdería.
	editado.SprintID = copiarPuntero(e.SprintID)
	editado.HorasEstimadas = copiarPuntero(horasEstimadas)
	return editado, nil
}
