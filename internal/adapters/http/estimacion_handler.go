package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// registradorEstimacion es lo que el handler y la pantalla necesitan del caso de uso de HU-17.
type registradorEstimacion interface {
	Ejecutar(ctx context.Context, d app.DatosEstimacion) (domain.ResumenRonda, error)
}

// pedidoEstimacion es el JSON de POST /proyectos/{id}/rondas/{rondaID}/estimacion. Sin valor, queda
// en 0 y el dominio lo rechaza por estar fuera de la escala.
type pedidoEstimacion struct {
	Valor int `json:"valor"`
}

// RegistrarEstimacionAcordada agrega al mux la ruta de HU-17.
func RegistrarEstimacionAcordada(mux *http.ServeMux, uc registradorEstimacion) {
	mux.HandleFunc("POST /proyectos/{id}/rondas/{rondaID}/estimacion", handlerRegistrarEstimacion(uc))
}

// handlerRegistrarEstimacion atiende POST /proyectos/{id}/rondas/{rondaID}/estimacion: 200 con el
// resumen de la ronda; un pedido mal formado da 400 sin llamar al caso de uso. Los errores se
// responden con responderRonda de HU-14. Implementa HU-17 (CA-17.9).
func handlerRegistrarEstimacion(uc registradorEstimacion) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, rondaID, err := idsRondaDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoEstimacion
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		res, err := uc.Ejecutar(r.Context(), app.DatosEstimacion{ProyectoID: proyectoID, RondaID: rondaID, Valor: pedido.Valor})
		responderRonda(w, http.StatusOK, res, err, "registrar la estimación")
	}
}
