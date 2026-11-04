package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

func pHoras(h float64) *float64 { return &h }

// verHoras muestra un valor opcional en los mensajes de error.
func verHoras(h *float64) any {
	if h == nil {
		return "nil"
	}
	return *h
}

func mismasHoras(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// editarYGuardar aplica los cambios en el dominio, los guarda y relee el elemento, como hace el caso
// de uso de edición.
func editarYGuardar(t *testing.T, repo *RepositorioBacklog, e domain.ElementoBacklog, c domain.CambiosElemento) domain.ElementoBacklog {
	t.Helper()
	ctx := context.Background()
	editado, err := e.Editar(c)
	if err != nil {
		t.Fatalf("editar: %v", err)
	}
	if err := repo.ActualizarElemento(ctx, editado, c); err != nil {
		t.Fatalf("actualizar: %v", err)
	}
	got, err := repo.BuscarElementoPorID(ctx, e.ID)
	if err != nil {
		t.Fatalf("releer: %v", err)
	}
	return got
}

func cambiarHoras(h *float64) domain.CambiosElemento {
	return domain.CambiosElemento{CambiarHorasEstimadas: true, HorasEstimadas: h}
}

// Antes de HU-19, un cambio solo de horas ejecutaba UPDATE ... SET id = id y las perdía sin error.
// Traza: HU-19 CA-19.13, CA-19.14 (RN-19.2, RN-19.3)
func TestRepositorioBacklog_HorasEstimadas_IdaYVuelta(t *testing.T) {
	for _, horas := range []float64{8, 0.5, 7.5, 1000} {
		repo, ids := nuevoBacklogDePrueba(t, "Alfa")
		e := guardarElemento(t, repo, ids[0], "Login")

		got := editarYGuardar(t, repo, e, cambiarHoras(pHoras(horas)))

		if !mismasHoras(got.HorasEstimadas, pHoras(horas)) {
			t.Errorf("%v h: se releyó %v", horas, verHoras(got.HorasEstimadas))
		}
	}
}

// Traza: HU-19 CA-19.13 (RN-19.3; sección 5: en medias horas, como registros_esfuerzo)
func TestRepositorioBacklog_HorasEstimadas_SeGuardanEnMediasHoras(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")
	editarYGuardar(t, repo, e, cambiarHoras(pHoras(7.5)))

	var medias int
	err := repo.db.QueryRow(`SELECT medias_horas_estimadas FROM backlog_elementos WHERE id = ?`, e.ID).Scan(&medias)

	if err != nil {
		t.Fatal(err)
	}
	if medias != 15 {
		t.Errorf("medias_horas_estimadas = %d, se esperaba 15", medias)
	}
}

// Traza: HU-19 CA-19.13 (RN-19.3)
func TestRepositorioBacklog_HorasEstimadas_CambiarYQuitar(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")

	e = editarYGuardar(t, repo, e, cambiarHoras(pHoras(8)))
	e = editarYGuardar(t, repo, e, cambiarHoras(pHoras(10)))
	if !mismasHoras(e.HorasEstimadas, pHoras(10)) {
		t.Fatalf("después de cambiar: HorasEstimadas = %v, se esperaba 10", verHoras(e.HorasEstimadas))
	}

	e = editarYGuardar(t, repo, e, cambiarHoras(nil))
	if e.HorasEstimadas != nil {
		t.Errorf("después de quitar: HorasEstimadas = %v, se esperaba nil", *e.HorasEstimadas)
	}
	var medias sql.NullInt64
	if err := repo.db.QueryRow(`SELECT medias_horas_estimadas FROM backlog_elementos WHERE id = ?`, e.ID).Scan(&medias); err != nil {
		t.Fatal(err)
	}
	if medias.Valid {
		t.Errorf("medias_horas_estimadas = %d, se esperaba NULL", medias.Int64)
	}
}

// Traza: HU-19 CA-19.13, CA-19.15 (RN-06.1, RN-19.1, RN-19.3)
func TestRepositorioBacklog_EditarOtrosCampos_NoTocaLasHorasEstimadas(t *testing.T) {
	baja, enProgreso := domain.PrioridadBaja, domain.EstadoEnProgreso
	titulo, descripcion, trece := "Inicio de sesión", "Entrar", 13
	criterios := []string{"Nuevo"}
	tests := []struct {
		nombre string
		c      domain.CambiosElemento
	}{
		{"título", domain.CambiosElemento{Titulo: &titulo}},
		{"descripción", domain.CambiosElemento{Descripcion: &descripcion}},
		{"prioridad", domain.CambiosElemento{Prioridad: &baja}},
		{"cambiar story points", domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &trece}},
		{"quitar story points", domain.CambiosElemento{CambiarStoryPoints: true}},
		{"criterios", domain.CambiosElemento{Criterios: &criterios}},
		{"estado", domain.CambiosElemento{Estado: &enProgreso}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			repo, ids := nuevoBacklogDePrueba(t, "Alfa")
			e := guardarElemento(t, repo, ids[0], "Login")
			e = editarYGuardar(t, repo, e, cambiarHoras(pHoras(8)))

			got := editarYGuardar(t, repo, e, tt.c)

			if !mismasHoras(got.HorasEstimadas, pHoras(8)) {
				t.Errorf("HorasEstimadas = %v, se esperaba 8", verHoras(got.HorasEstimadas))
			}
		})
	}
}

// Traza: HU-19 CA-19.15 (RN-19.1)
func TestRepositorioBacklog_EditarHorasEstimadas_NoTocaLosDemasCampos(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")

	got := editarYGuardar(t, repo, e, cambiarHoras(pHoras(8)))

	want := e
	want.HorasEstimadas = pHoras(8)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("se recuperó %+v, se esperaba %+v", got, want)
	}
}

// Dos personas parten del mismo elemento: una carga las horas y la otra, que todavía lo ve sin
// estimar, cambia la prioridad. No se pisan, en ninguno de los dos órdenes.
// Traza: HU-19 CA-19.13 (RN-19.3; sección 5, como HU-06 sección 5)
func TestRepositorioBacklog_HorasEstimadas_EdicionesSimultaneasNoSePisan(t *testing.T) {
	ctx := context.Background()
	baja := domain.PrioridadBaja
	conHoras, conPrioridad := cambiarHoras(pHoras(8)), domain.CambiosElemento{Prioridad: &baja}
	for _, tt := range []struct {
		nombre         string
		primero, luego domain.CambiosElemento
	}{
		{"primero las horas", conHoras, conPrioridad},
		{"primero la prioridad", conPrioridad, conHoras},
	} {
		t.Run(tt.nombre, func(t *testing.T) {
			repo, ids := nuevoBacklogDePrueba(t, "Alfa")
			e := guardarElemento(t, repo, ids[0], "Login")
			e1, _ := e.Editar(tt.primero)
			e2, _ := e.Editar(tt.luego)
			if err := repo.ActualizarElemento(ctx, e1, tt.primero); err != nil {
				t.Fatal(err)
			}
			if err := repo.ActualizarElemento(ctx, e2, tt.luego); err != nil {
				t.Fatal(err)
			}

			got, err := repo.BuscarElementoPorID(ctx, e.ID)

			if err != nil {
				t.Fatal(err)
			}
			if !mismasHoras(got.HorasEstimadas, pHoras(8)) || got.Prioridad != domain.PrioridadBaja {
				t.Errorf("se perdió un cambio: HorasEstimadas = %v, Prioridad = %q", verHoras(got.HorasEstimadas), got.Prioridad)
			}
		})
	}
}

// Traza: HU-19 CA-19.13 (RN-19.3; sección 3: el listado del backlog también trae las horas estimadas)
func TestRepositorioBacklog_ListarElementos_TraeLasHorasEstimadas(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	login := guardarElemento(t, repo, ids[0], "Login")
	guardarElemento(t, repo, ids[0], "Reportes")
	ayuda := guardarElemento(t, repo, ids[0], "Ayuda")
	editarYGuardar(t, repo, login, cambiarHoras(pHoras(8)))
	editarYGuardar(t, repo, ayuda, cambiarHoras(pHoras(0.5)))

	got, err := repo.ListarElementos(ctx, ids[0])

	if err != nil {
		t.Fatal(err)
	}
	want := map[string]*float64{"Login": pHoras(8), "Reportes": nil, "Ayuda": pHoras(0.5)}
	if len(got) != len(want) {
		t.Fatalf("elementos = %d, se esperaban %d", len(got), len(want))
	}
	for _, e := range got {
		if !mismasHoras(e.HorasEstimadas, want[e.Titulo]) {
			t.Errorf("%q: HorasEstimadas = %v, se esperaba %v", e.Titulo, verHoras(e.HorasEstimadas), verHoras(want[e.Titulo]))
		}
	}
}

// El dominio crea los elementos sin estimar, pero el repositorio guarda lo que recibe.
// Traza: HU-19 CA-19.15 (RN-19.3; sección 3)
func TestRepositorioBacklog_GuardarElemento_GuardaLasHorasEstimadas(t *testing.T) {
	ctx := context.Background()
	for _, tt := range []struct {
		nombre string
		horas  *float64
	}{
		{"recién creado, sin estimar", nil},
		{"con horas", pHoras(8)},
	} {
		t.Run(tt.nombre, func(t *testing.T) {
			repo, ids := nuevoBacklogDePrueba(t, "Alfa")
			e := elementoDePrueba(t, ids[0], "Login", nil)
			e.HorasEstimadas = tt.horas
			if err := repo.GuardarElemento(ctx, &e); err != nil {
				t.Fatal(err)
			}

			got, err := repo.BuscarElementoPorID(ctx, e.ID)

			if err != nil {
				t.Fatal(err)
			}
			if !mismasHoras(got.HorasEstimadas, tt.horas) {
				t.Errorf("HorasEstimadas = %v, se esperaba %v", verHoras(got.HorasEstimadas), verHoras(tt.horas))
			}
		})
	}
}

// Con el elemento asignado a un sprint (HU-09), editar conserva el sprint y las horas.
// Traza: HU-19 CA-19.13 (RN-19.3, RN-06.1)
func TestRepositorioBacklog_HorasEstimadas_SobrevivenAlEditarConSprint(t *testing.T) {
	ctx := context.Background()
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := guardarElemento(t, repo, ids[0], "Login")
	dia := time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC)
	s := domain.Sprint{ProyectoID: ids[0], FechaInicio: dia, FechaFin: dia, Objetivo: "o", Estado: domain.SprintAbierto}
	if err := NuevoRepositorioSprints(repo.db).GuardarSprint(ctx, &s); err != nil {
		t.Fatal(err)
	}
	e.SprintID = &s.ID
	if err := repo.GuardarAsignacion(ctx, &e); err != nil {
		t.Fatal(err)
	}
	e = editarYGuardar(t, repo, e, cambiarHoras(pHoras(8)))

	baja := domain.PrioridadBaja
	got := editarYGuardar(t, repo, e, domain.CambiosElemento{Prioridad: &baja})

	if !mismasHoras(got.HorasEstimadas, pHoras(8)) {
		t.Errorf("HorasEstimadas = %v, se esperaba 8", verHoras(got.HorasEstimadas))
	}
	if got.SprintID == nil || *got.SprintID != s.ID {
		t.Errorf("SprintID = %v, se esperaba %d", got.SprintID, s.ID)
	}
}

// Un cambio solo de horas también informa que el elemento no existe.
// Traza: HU-19 CA-19.13 (RN-19.3, RN-06.6)
func TestRepositorioBacklog_HorasEstimadas_ElementoInexistente(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")
	e := elementoDePrueba(t, ids[0], "Login", nil)
	e.ID = 999
	e.HorasEstimadas = pHoras(8)

	err := repo.ActualizarElemento(context.Background(), e, cambiarHoras(pHoras(8)))

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}
