package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Lo que los handlers y las pantallas de HU-14 necesitan de los casos de uso.
type (
	iniciadorRondas interface {
		Ejecutar(ctx context.Context, d app.DatosIniciarRonda) (domain.ResumenRonda, error)
	}
	votadorRondas interface {
		Ejecutar(ctx context.Context, d app.DatosVoto) (domain.ResumenRonda, error)
	}
	consultorRondas interface {
		Ejecutar(ctx context.Context, proyectoID, rondaID int64) (domain.ResumenRonda, error)
		UltimaDelElemento(ctx context.Context, proyectoID, elementoID int64) (*domain.ResumenRonda, error)
	}
)

// pedidoVoto es el JSON de POST /proyectos/{id}/rondas/{rondaID}/votos.
type pedidoVoto struct {
	Integrante string `json:"integrante"`
	Valor      int    `json:"valor"`
}

// respuestaRonda es el resumen de una ronda en JSON (HU-14, sección 3): de cada votante, solo si
// votó (RN-14.3); con la ronda cerrada, además, la revelación de HU-15.
type (
	respuestaVotante struct {
		Nombre string `json:"nombre"`
		Voto   bool   `json:"voto"`
	}
	respuestaRonda struct {
		ID         int64              `json:"id"`
		ProyectoID int64              `json:"proyecto_id"`
		ElementoID int64              `json:"elemento_id"`
		Numero     int                `json:"numero"`
		Estado     string             `json:"estado"`
		Votantes   []respuestaVotante `json:"votantes"`
		Revelacion *respuestaRevelada `json:"revelacion,omitempty"` // HU-15: solo con la ronda cerrada
	}
	respuestaVoto struct {
		Nombre string `json:"nombre"`
		Valor  int    `json:"valor"`
	}
	respuestaRevelada struct {
		Votos                   []respuestaVoto `json:"votos"`
		Minimo                  int             `json:"minimo"`
		Maximo                  int             `json:"maximo"`
		DiferenciaSignificativa bool            `json:"diferencia_significativa"`
	}
)

// erroresNoEncontradoRonda y erroresDeValidacionRonda son los errores de HU-14 que se responden
// como 404 y como 400 (spec, sección 7).
var (
	erroresNoEncontradoRonda = []error{domain.ErrProyectoNoEncontrado, domain.ErrElementoNoEncontrado, domain.ErrRondaNoEncontrada}
	erroresDeValidacionRonda = []error{domain.ErrElementoTerminado, domain.ErrRondaSinVotantes, domain.ErrRondaYaAbierta,
		domain.ErrRondaCerrada, domain.ErrNoEsVotante, domain.ErrVotoInvalido}
)

// RegistrarPlanningPoker agrega las rutas de HU-14 al mux.
func RegistrarPlanningPoker(mux *http.ServeMux, iniciar iniciadorRondas, votar votadorRondas, consultar consultorRondas) {
	mux.HandleFunc("POST /proyectos/{id}/backlog/{elementoID}/rondas", handlerIniciarRonda(iniciar))
	mux.HandleFunc("GET /proyectos/{id}/rondas/{rondaID}", handlerConsultarRonda(consultar))
	mux.HandleFunc("POST /proyectos/{id}/rondas/{rondaID}/votos", handlerVotar(votar))
}

// idsRondaDesdeRuta lee el proyecto y la ronda de la ruta.
func idsRondaDesdeRuta(r *http.Request) (int64, int64, error) {
	proyectoID, err := proyectoIDDesdeRuta(r)
	if err != nil {
		return 0, 0, err
	}
	rondaID, err := strconv.ParseInt(r.PathValue("rondaID"), 10, 64)
	if err != nil || rondaID <= 0 {
		return 0, 0, errors.New("el id de la ronda debe ser un número entero positivo")
	}
	return proyectoID, rondaID, nil
}

// handlerIniciarRonda atiende POST /proyectos/{id}/backlog/{elementoID}/rondas: 201 con el resumen
// de la ronda abierta. Implementa HU-14 (CA-14.1, CA-14.6 a CA-14.8).
func handlerIniciarRonda(uc iniciadorRondas) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, elementoID, err := idsElementoDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		res, err := uc.Ejecutar(r.Context(), app.DatosIniciarRonda{ProyectoID: proyectoID, ElementoID: elementoID})
		responderRonda(w, http.StatusCreated, res, err, "iniciar ronda")
	}
}

// handlerConsultarRonda atiende GET /proyectos/{id}/rondas/{rondaID}: 200 con el resumen.
// Implementa HU-14 (CA-14.7, CA-14.8).
func handlerConsultarRonda(uc consultorRondas) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, rondaID, err := idsRondaDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		res, err := uc.Ejecutar(r.Context(), proyectoID, rondaID)
		responderRonda(w, http.StatusOK, res, err, "consultar ronda")
	}
}

// handlerVotar atiende POST /proyectos/{id}/rondas/{rondaID}/votos: 200 con el resumen como quedó,
// sin el valor; un cuerpo mal formado da 400 sin llamar al caso de uso. Implementa HU-14
// (CA-14.2 a CA-14.5, CA-14.8).
func handlerVotar(uc votadorRondas) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectoID, rondaID, err := idsRondaDesdeRuta(r)
		if err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{err.Error()})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, tamanoMaximoCuerpo)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var pedido pedidoVoto
		if err := dec.Decode(&pedido); err != nil {
			escribirJSON(w, http.StatusBadRequest, respuestaError{"el cuerpo debe ser un JSON válido"})
			return
		}
		res, err := uc.Ejecutar(r.Context(), app.DatosVoto{
			ProyectoID: proyectoID, RondaID: rondaID, Integrante: pedido.Integrante, Valor: pedido.Valor,
		})
		responderRonda(w, http.StatusOK, res, err, "votar")
	}
}

// responderRonda escribe el resumen con el estado indicado, o el error con su estado (sección 7);
// ante una falla interna responde 500 sin el detalle.
func responderRonda(w http.ResponseWriter, estado int, res domain.ResumenRonda, err error, operacion string) {
	if err == nil {
		escribirJSON(w, estado, aRespuestaRonda(res))
		return
	}
	for _, e := range erroresNoEncontradoRonda {
		if errors.Is(err, e) {
			escribirJSON(w, http.StatusNotFound, respuestaError{e.Error()})
			return
		}
	}
	for _, e := range erroresDeValidacionRonda {
		if errors.Is(err, e) {
			escribirJSON(w, http.StatusBadRequest, respuestaError{e.Error()})
			return
		}
	}
	log.Printf("%s: %v", operacion, err)
	escribirJSON(w, http.StatusInternalServerError, respuestaError{"error interno al " + operacion})
}

func aRespuestaRonda(r domain.ResumenRonda) respuestaRonda {
	votantes := make([]respuestaVotante, len(r.Votantes))
	for i, v := range r.Votantes {
		votantes[i] = respuestaVotante{Nombre: v.Nombre, Voto: v.YaVoto}
	}
	res := respuestaRonda{ID: r.ID, ProyectoID: r.ProyectoID, ElementoID: r.ElementoID, Numero: r.Numero,
		Estado: string(r.Estado), Votantes: votantes}
	if rev := r.Revelacion; rev != nil { // HU-15 (RN-15.1)
		res.Revelacion = &respuestaRevelada{Votos: make([]respuestaVoto, len(rev.Votos)), Minimo: rev.Minimo,
			Maximo: rev.Maximo, DiferenciaSignificativa: rev.DiferenciaSignificativa}
		for i, v := range rev.Votos {
			res.Revelacion.Votos[i] = respuestaVoto{Nombre: v.Nombre, Valor: v.Valor}
		}
	}
	return res
}
