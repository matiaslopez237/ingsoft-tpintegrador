package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// textoAclaracionSprintCerrado explica qué cuenta en un sprint Cerrado (RN-23.8).
const textoAclaracionSprintCerrado = "Solo cuentan los elementos que siguen en este sprint. " +
	"Los que salieron de él no suman acá, pero sí en el total del proyecto."

// datosHoras es el bloque "Horas y desviación" de los detalles del sprint y del proyecto (HU-23): solo
// texto ya formateado. Los números los calcula el dominio.
type datosHoras struct {
	Elementos                                                          int
	HorasEstimadas, HorasReales, DesviacionHoras, DesviacionPorcentaje string
	Avisos                                                             []string
}

// vistaHoras da formato al total. Implementa HU-23 (CA-23.17, CA-23.18, sección 3).
func vistaHoras(t domain.TotalHoras) datosHoras {
	v := datosHoras{
		Elementos: t.Elementos, HorasEstimadas: textoHoras(t.HorasEstimadas), HorasReales: textoHoras(t.HorasReales),
	}
	v.DesviacionHoras, v.DesviacionPorcentaje = textosDesviacion(t.DesviacionHoras, t.DesviacionPorcentaje)
	switch {
	case t.Elementos == 0:
		v.Avisos = []string{"Todavía no hay elementos."}
	case t.ElementosSinEstimar > 0:
		v.Avisos = []string{avisoSinEstimar(t)}
	}
	return v
}

// avisoSinEstimar cuenta los elementos sin estimar y, si tienen horas reales, explica por qué la
// desviación no sale de restar las dos cifras del bloque (RN-23.5).
func avisoSinEstimar(t domain.TotalHoras) string {
	nombre := "elementos"
	if t.Elementos == 1 {
		nombre = "elemento"
	}
	cuenta := fmt.Sprintf("%d de %d %s sin estimar", t.ElementosSinEstimar, t.Elementos, nombre)
	if t.HorasRealesSinEstimar == 0 {
		return cuenta + "."
	}
	return fmt.Sprintf("%s: sus %s están en las horas reales, pero no en la desviación.", cuenta, textoHoras(t.HorasRealesSinEstimar))
}

// errorDeHoras muestra la página 404 si no existe el proyecto o el sprint, con el mensaje del dominio y
// no el envuelto por el caso de uso, o la de falla interna.
func (p *pantallas) errorDeHoras(w http.ResponseWriter, err error) {
	for _, noEncontrado := range erroresConsultaHoras {
		if errors.Is(err, noEncontrado) {
			p.renderizar(w, http.StatusNotFound, "error.html", datosError{noEncontrado.Error()})
			return
		}
	}
	p.fallaInterna(w, "consultar horas", err)
}
