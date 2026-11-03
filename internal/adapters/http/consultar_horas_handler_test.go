package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// consultorHorasProyectoFalso devuelve el total que se le indica y guarda lo que recibió.
type consultorHorasProyectoFalso struct {
	total      domain.TotalHoras
	err        error
	proyectoID int64
	ctx        context.Context
	llamado    bool
}

func (f *consultorHorasProyectoFalso) Ejecutar(ctx context.Context, proyectoID int64) (domain.TotalHoras, error) {
	f.llamado, f.ctx = true, ctx
	f.proyectoID = proyectoID
	return f.total, f.err
}

// consultorHorasSprintFalso devuelve el total que se le indica y guarda lo que recibió.
type consultorHorasSprintFalso struct {
	total      domain.TotalHoras
	err        error
	proyectoID int64
	sprintID   int64
	ctx        context.Context
	llamado    bool
}

func (f *consultorHorasSprintFalso) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.TotalHoras, error) {
	f.llamado, f.ctx = true, ctx
	f.proyectoID, f.sprintID = proyectoID, sprintID
	return f.total, f.err
}

const (
	rutaHorasProyecto = "/proyectos/7/horas"
	rutaHorasSprint   = "/proyectos/7/sprints/3/horas"
)

func getHoras(t *testing.T, proyecto *consultorHorasProyectoFalso, sprint *consultorHorasSprintFalso, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	RegistrarConsultarHoras(mux, proyecto, sprint)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
	return rec
}

func estimado(estimadas, reales float64) domain.HorasElemento {
	return domain.HorasElemento{HorasEstimadas: &estimadas, HorasReales: reales}
}

func sinEstimar(reales float64) domain.HorasElemento {
	return domain.HorasElemento{HorasReales: reales}
}

// casoHoras es un total calculado con la función real del dominio y el JSON literal que le corresponde.
type casoHoras struct {
	nombre                                                     string
	elementos                                                  []domain.HorasElemento
	estimadas, reales, realesSinEstimar, desvHoras, desvPorcen string
	cantidad, cantidadSinEstimar                               string
}

// casosHoras son los mismos para las dos rutas. Los valores son el JSON literal: 25,0 sale como el
// número 25, y una desviación que no existe es null, no 0.
func casosHoras() []casoHoras {
	return []casoHoras{
		{"sprint normal", []domain.HorasElemento{estimado(8, 10), estimado(4, 3)}, "12", "13", "0", "1", "8.3", "2", "0"},
		{"proyecto con tres elementos", []domain.HorasElemento{estimado(8, 10), estimado(4, 3), estimado(20, 15)}, "32", "28", "0", "-4", "-12.5", "3", "0"},
		{"sobre las sumas", []domain.HorasElemento{estimado(10, 20), estimado(90, 90)}, "100", "110", "0", "10", "10", "2", "0"},
		{"desvíos que se compensan: 0 y no null", []domain.HorasElemento{estimado(8, 10), estimado(8, 6)}, "16", "16", "0", "0", "0", "2", "0"},
		{"horas con decimal", []domain.HorasElemento{estimado(20, 20.5), estimado(20, 20)}, "40", "40.5", "0", "0.5", "1.3", "2", "0"},
		{"negativos con decimal", []domain.HorasElemento{estimado(20, 19.5), estimado(20, 20)}, "40", "39.5", "0", "-0.5", "-1.3", "2", "0"},
		{"con un elemento sin estimar", []domain.HorasElemento{estimado(8, 10), sinEstimar(5)}, "8", "15", "5", "2", "25", "2", "1"},
		{"ningún elemento estimado: null", []domain.HorasElemento{sinEstimar(5), sinEstimar(3)}, "0", "8", "8", "null", "null", "2", "2"},
		{"sin elementos: null", nil, "0", "0", "0", "null", "null", "0", "0"},
		{"estimado sin esfuerzo", []domain.HorasElemento{estimado(8, 0)}, "8", "0", "0", "-8", "-100", "1", "0"},
	}
}

func comprobarCamposHoras(t *testing.T, rec *httptest.ResponseRecorder, want map[string]string) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, se esperaba 200; cuerpo: %s", rec.Code, rec.Body)
	}
	got := camposJSON(t, rec)
	if len(got) != len(want) {
		t.Errorf("campos = %v, se esperaban %v", got, want)
	}
	for campo, valor := range want {
		if got[campo] != valor {
			t.Errorf("%s = %s, se esperaba %s", campo, got[campo], valor)
		}
	}
}

func (c casoHoras) campos() map[string]string {
	return map[string]string{
		"horas_estimadas": c.estimadas, "horas_reales": c.reales, "horas_reales_sin_estimar": c.realesSinEstimar,
		"desviacion_horas": c.desvHoras, "desviacion_porcentaje": c.desvPorcen,
		"elementos": c.cantidad, "elementos_sin_estimar": c.cantidadSinEstimar,
	}
}

// Traza: HU-23 CA-23.16, CA-23.1, CA-23.2, CA-23.3, CA-23.4, CA-23.5, CA-23.6, CA-23.7, CA-23.8, CA-23.9 (sección 3)
func TestGetHorasProyecto_Devuelve200ConElTotal(t *testing.T) {
	for _, tt := range casosHoras() {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &consultorHorasProyectoFalso{total: domain.TotalizarHoras(tt.elementos)}

			rec := getHoras(t, uc, nil, rutaHorasProyecto)

			want := tt.campos()
			want["proyecto_id"] = "7"
			comprobarCamposHoras(t, rec, want)
			if uc.proyectoID != 7 {
				t.Errorf("el caso de uso recibió el proyecto %d", uc.proyectoID)
			}
		})
	}
}

// Traza: HU-23 CA-23.16, CA-23.1, CA-23.2, CA-23.3, CA-23.4, CA-23.5, CA-23.6, CA-23.7, CA-23.8, CA-23.9 (sección 3)
func TestGetHorasSprint_Devuelve200ConElTotal(t *testing.T) {
	for _, tt := range casosHoras() {
		t.Run(tt.nombre, func(t *testing.T) {
			uc := &consultorHorasSprintFalso{total: domain.TotalizarHoras(tt.elementos)}

			rec := getHoras(t, nil, uc, rutaHorasSprint)

			want := tt.campos()
			want["proyecto_id"], want["sprint_id"] = "7", "3"
			comprobarCamposHoras(t, rec, want)
			if uc.proyectoID != 7 || uc.sprintID != 3 {
				t.Errorf("el caso de uso recibió el proyecto %d y el sprint %d", uc.proyectoID, uc.sprintID)
			}
		})
	}
}

// Con 2000 estimadas y media hora de menos, el porcentaje redondeado sería -0: el dominio lo evita
// (RN-19.6), y acá se verifica lo que le llega al cliente.
// Traza: HU-23 CA-23.16, CA-23.4 (RN-23.4)
func TestGetHoras_NuncaEscribeCeroNegativo(t *testing.T) {
	uc := &consultorHorasSprintFalso{total: domain.TotalizarHoras([]domain.HorasElemento{estimado(1000, 999.5), estimado(1000, 1000)})}

	got := camposJSON(t, getHoras(t, nil, uc, rutaHorasSprint))

	if got["desviacion_porcentaje"] != "0" || got["desviacion_horas"] != "-0.5" {
		t.Errorf("desviación = %s h y %s %%, se esperaba -0.5 y 0", got["desviacion_horas"], got["desviacion_porcentaje"])
	}
}

// Traza: HU-23 CA-23.16 (sección 3)
func TestGetHoras_RespondeJSON(t *testing.T) {
	falla := errors.New("base caída")
	tests := []struct {
		nombre string
		rec    func(t *testing.T) *httptest.ResponseRecorder
	}{
		{"200 de proyecto", func(t *testing.T) *httptest.ResponseRecorder {
			return getHoras(t, &consultorHorasProyectoFalso{}, nil, rutaHorasProyecto)
		}},
		{"200 de sprint", func(t *testing.T) *httptest.ResponseRecorder {
			return getHoras(t, nil, &consultorHorasSprintFalso{}, rutaHorasSprint)
		}},
		{"400", func(t *testing.T) *httptest.ResponseRecorder {
			return getHoras(t, nil, &consultorHorasSprintFalso{}, "/proyectos/7/sprints/abc/horas")
		}},
		{"404", func(t *testing.T) *httptest.ResponseRecorder {
			return getHoras(t, &consultorHorasProyectoFalso{err: domain.ErrProyectoNoEncontrado}, nil, rutaHorasProyecto)
		}},
		{"500", func(t *testing.T) *httptest.ResponseRecorder {
			return getHoras(t, nil, &consultorHorasSprintFalso{err: falla}, rutaHorasSprint)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if ct := tt.rec(t).Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, se esperaba application/json", ct)
			}
		})
	}
}

// El mensaje es el del dominio, aunque el error llegue envuelto: se reconoce con errors.Is.
// Traza: HU-23 CA-23.14, CA-23.15, CA-23.16 (RN-23.11; sección 7)
func TestGetHoras_NoEncontrado_Devuelve404(t *testing.T) {
	tests := []struct {
		nombre  string
		rec     func(t *testing.T, err error) *httptest.ResponseRecorder
		err     error
		wantMsg string
	}{
		{"proyecto: proyecto inexistente", getHorasProyectoConError, domain.ErrProyectoNoEncontrado, "proyecto no encontrado"},
		{"proyecto: proyecto inexistente, envuelto", getHorasProyectoConError, fmt.Errorf("buscar proyecto 7: %w", domain.ErrProyectoNoEncontrado), "proyecto no encontrado"},
		{"sprint: proyecto inexistente", getHorasSprintConError, domain.ErrProyectoNoEncontrado, "proyecto no encontrado"},
		{"sprint: sprint inexistente", getHorasSprintConError, domain.ErrSprintNoEncontrado, "sprint no encontrado"},
		{"sprint: sprint inexistente, envuelto", getHorasSprintConError, fmt.Errorf("buscar sprint 3: %w", domain.ErrSprintNoEncontrado), "sprint no encontrado"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			rec := tt.rec(t, tt.err)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, se esperaba 404; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
		})
	}
}

func getHorasProyectoConError(t *testing.T, err error) *httptest.ResponseRecorder {
	t.Helper()
	return getHoras(t, &consultorHorasProyectoFalso{err: err}, nil, rutaHorasProyecto)
}

func getHorasSprintConError(t *testing.T, err error) *httptest.ResponseRecorder {
	t.Helper()
	return getHoras(t, nil, &consultorHorasSprintFalso{err: err}, rutaHorasSprint)
}

// IDs inválidos en la ruta: 400, como en los demás handlers, sin llamar al caso de uso. El proyecto se
// lee primero, y "anteriores" no se confunde con la ruta de HU-12.
// Traza: HU-23 CA-23.16 (sección 6, sección 7)
func TestGetHoras_IDInvalido_Devuelve400(t *testing.T) {
	const (
		idProyecto = "el id del proyecto debe ser un número entero positivo"
		idSprint   = "el id del sprint debe ser un número entero positivo"
	)
	tests := []struct {
		nombre  string
		ruta    string
		wantMsg string
	}{
		{"proyecto: no numérico", "/proyectos/abc/horas", idProyecto},
		{"proyecto: en cero", "/proyectos/0/horas", idProyecto},
		{"proyecto: negativo", "/proyectos/-1/horas", idProyecto},
		{"sprint, id de proyecto no numérico", "/proyectos/abc/sprints/3/horas", idProyecto},
		{"sprint, id de proyecto en cero", "/proyectos/0/sprints/3/horas", idProyecto},
		{"sprint, id de proyecto negativo", "/proyectos/-1/sprints/3/horas", idProyecto},
		{"sprint: no numérico", "/proyectos/7/sprints/abc/horas", idSprint},
		{"sprint: en cero", "/proyectos/7/sprints/0/horas", idSprint},
		{"sprint: negativo", "/proyectos/7/sprints/-1/horas", idSprint},
		{"sprint: anteriores", "/proyectos/7/sprints/anteriores/horas", idSprint},
		{"los dos inválidos: primero el proyecto", "/proyectos/abc/sprints/abc/horas", idProyecto},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			proyecto, sprint := &consultorHorasProyectoFalso{}, &consultorHorasSprintFalso{}

			rec := getHoras(t, proyecto, sprint, tt.ruta)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400; cuerpo: %s", rec.Code, rec.Body)
			}
			if got := mensajeDeError(t, rec); got != tt.wantMsg {
				t.Errorf("mensaje = %q, se esperaba %q", got, tt.wantMsg)
			}
			if proyecto.llamado || sprint.llamado {
				t.Error("no se esperaba que se llamara al caso de uso")
			}
		})
	}
}

// Cualquier otro error es una falla interna: 500, sin mostrar el detalle. El texto no se mira: uno que
// solo se parece a "sprint no encontrado" sigue siendo una falla interna.
// Traza: HU-23 CA-23.16 (sección 7)
func TestGetHoras_FallaInterna_Devuelve500SinDetalle(t *testing.T) {
	const generico = "error interno al consultar las horas"
	fallas := []error{
		fmt.Errorf("horas del sprint 3: %w", errors.New("database is locked")),
		errors.New("sprint no encontrado"),
		errors.New("proyecto no encontrado"),
	}
	for _, falla := range fallas {
		for nombre, get := range map[string]func(*testing.T, error) *httptest.ResponseRecorder{
			"proyecto": getHorasProyectoConError, "sprint": getHorasSprintConError,
		} {
			t.Run(nombre+": "+falla.Error(), func(t *testing.T) {
				rec := get(t, falla)

				if rec.Code != http.StatusInternalServerError {
					t.Fatalf("status = %d, se esperaba 500; cuerpo: %s", rec.Code, rec.Body)
				}
				if got := mensajeDeError(t, rec); got != generico {
					t.Errorf("mensaje = %q: no tiene que mostrar el detalle interno", got)
				}
			})
		}
	}
}

// El caso de uso recibe el contexto del pedido, para que cancelarlo corte la consulta.
// Traza: HU-23 (sección 5: contexto)
func TestGetHoras_ElCasoDeUsoRecibeElContextoDelPedido(t *testing.T) {
	type clave struct{}
	ctx := context.WithValue(context.Background(), clave{}, "pedido")
	proyecto, sprint := &consultorHorasProyectoFalso{}, &consultorHorasSprintFalso{}
	mux := http.NewServeMux()
	RegistrarConsultarHoras(mux, proyecto, sprint)

	for _, ruta := range []string{rutaHorasProyecto, rutaHorasSprint} {
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, ruta, nil).WithContext(ctx))
	}

	if proyecto.ctx != ctx {
		t.Errorf("proyecto: el caso de uso recibió %v, se esperaba el contexto del pedido", proyecto.ctx)
	}
	if sprint.ctx != ctx {
		t.Errorf("sprint: el caso de uso recibió %v, se esperaba el contexto del pedido", sprint.ctx)
	}
}
