package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir los puertos de los casos de uso (se verifica al compilar).
var (
	_ app.RepositorioEsfuerzo = (*RepositorioEsfuerzo)(nil)
	_ app.LectorHorasReales   = (*RepositorioEsfuerzo)(nil)
)

// baseEsfuerzo es una base SQLite real con un proyecto, la integrante "Ana" y el elemento "Login".
type baseEsfuerzo struct {
	db         *sql.DB
	repo       *RepositorioEsfuerzo
	proyectos  *RepositorioProyectos
	proyectoID int64
	elementoID int64
}

func nuevaBaseEsfuerzo(t *testing.T) baseEsfuerzo {
	t.Helper()
	ctx := context.Background()
	db, err := AbrirDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	proyectos := NuevoRepositorioProyectos(db)
	p, err := domain.NuevoProyecto("Alfa", time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := proyectos.Guardar(ctx, &p); err != nil {
		t.Fatalf("guardar proyecto: %v", err)
	}
	if err := proyectos.AgregarIntegrante(ctx, p.ID, domain.Integrante{Nombre: "Ana", Roles: []string{"Desarrolladora"}}); err != nil {
		t.Fatalf("agregar integrante: %v", err)
	}
	e := elementoDePrueba(t, p.ID, "Login", nil)
	if err := NuevoRepositorioBacklog(db).GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	return baseEsfuerzo{db: db, repo: NuevoRepositorioEsfuerzo(db), proyectos: proyectos, proyectoID: p.ID, elementoID: e.ID}
}

// registro arma un registro válido de "Ana" en el elemento de la base, con fecha en la zona del servidor.
func (b baseEsfuerzo) registro(t *testing.T, horas float64) domain.RegistroEsfuerzo {
	t.Helper()
	fecha := time.Date(2026, time.October, 24, 0, 0, 0, 0, time.Local)
	hoy := time.Date(2026, time.October, 25, 12, 0, 0, 0, time.Local)
	r, err := domain.NuevoRegistroEsfuerzo(b.proyectoID, b.elementoID, "Ana", fecha, "Desarrollo", horas, hoy)
	if err != nil {
		t.Fatalf("crear registro de prueba: %v", err)
	}
	return r
}

// contarRegistros cuenta los registros del elemento directo en la base: solo lo necesitan los
// tests, así que no es parte del repositorio.
func contarRegistros(t *testing.T, db *sql.DB, elementoID int64) int {
	t.Helper()
	var n int
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM registros_esfuerzo WHERE elemento_id = ?`, elementoID).Scan(&n)
	if err != nil {
		t.Fatalf("contar registros de esfuerzo del elemento %d: %v", elementoID, err)
	}
	return n
}

// mismoRegistro compara la fecha con Equal (el mismo instante) y el resto campo por campo.
func mismoRegistro(a, b domain.RegistroEsfuerzo) bool {
	fa, fb := a.Fecha, b.Fecha
	a.Fecha, b.Fecha = time.Time{}, time.Time{}
	return a == b && fa.Equal(fb)
}

// Traza: HU-18 CA-18.12 (sección 3: el registro se recupera por su ID)
func TestRepositorioEsfuerzo_GuardarYBuscar_DevuelveLoMismo(t *testing.T) {
	ctx := context.Background()
	for _, horas := range []float64{0.5, 2, 7.5, 24} {
		b := nuevaBaseEsfuerzo(t)
		r := b.registro(t, horas)

		if err := b.repo.GuardarRegistro(ctx, &r); err != nil {
			t.Fatalf("%v h: guardar: %v", horas, err)
		}
		got, err := b.repo.BuscarRegistroPorID(ctx, r.ID)

		if err != nil {
			t.Fatalf("%v h: buscar: %v", horas, err)
		}
		if r.ID == 0 || !mismoRegistro(got, r) {
			t.Errorf("%v h: se recuperó %+v, se esperaba %+v", horas, got, r)
		}
	}
}

// Traza: HU-18 CA-18.12 (RN-18.9, RN-18.10)
func TestRepositorioEsfuerzo_DosRegistrosIguales_TienenIDsDistintos(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t)
	r1, r2 := b.registro(t, 2), b.registro(t, 2)

	if err := b.repo.GuardarRegistro(ctx, &r1); err != nil {
		t.Fatal(err)
	}
	if err := b.repo.GuardarRegistro(ctx, &r2); err != nil {
		t.Fatal(err)
	}

	if r1.ID == 0 || r1.ID == r2.ID {
		t.Errorf("IDs = %d y %d, se esperaban asignados y distintos", r1.ID, r2.ID)
	}
	if n := contarRegistros(t, b.db, b.elementoID); n != 2 {
		t.Errorf("registros del elemento = %d, se esperaban 2", n)
	}
}

func TestRepositorioEsfuerzo_BuscarInexistente_RetornaErrRegistroEsfuerzoNoEncontrado(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)

	_, err := b.repo.BuscarRegistroPorID(context.Background(), 999)

	if !errors.Is(err, domain.ErrRegistroEsfuerzoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrRegistroEsfuerzoNoEncontrado)
	}
}

// Las horas se guardan como medias horas enteras (sin redondeo) y la fecha como AAAA-MM-DD.
func TestRepositorioEsfuerzo_GuardaHorasEnMediasHorasYFechaComoTexto(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t)
	r := b.registro(t, 7.5)
	if err := b.repo.GuardarRegistro(ctx, &r); err != nil {
		t.Fatal(err)
	}

	var (
		medias int
		fecha  string
	)
	err := b.db.QueryRowContext(ctx, `SELECT medias_horas, fecha FROM registros_esfuerzo WHERE id = ?`, r.ID).
		Scan(&medias, &fecha)

	if err != nil {
		t.Fatal(err)
	}
	if medias != 15 || fecha != "2026-10-24" {
		t.Errorf("medias_horas = %d, fecha = %q; se esperaba 15 y \"2026-10-24\"", medias, fecha)
	}
}

// Traza: HU-18 CA-18.7 (si el elemento ya no existe al guardar, la clave foránea lo rechaza)
func TestRepositorioEsfuerzo_ElementoInexistente_RetornaErrElementoNoEncontrado(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t)
	r := b.registro(t, 2)
	r.ElementoID = 999

	err := b.repo.GuardarRegistro(ctx, &r)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if n := contarRegistros(t, b.db, 999); n != 0 {
		t.Errorf("registros del elemento = %d, no se esperaban registros", n)
	}
}

// Traza: HU-18 CA-18.13 (RN-18.3: no hay clave foránea hacia integrantes)
func TestRepositorioEsfuerzo_QuitarAlIntegrante_ConservaSusRegistros(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t)
	r := b.registro(t, 3)
	if err := b.repo.GuardarRegistro(ctx, &r); err != nil {
		t.Fatal(err)
	}

	if err := b.proyectos.QuitarIntegrante(ctx, b.proyectoID, "Ana"); err != nil {
		t.Fatalf("quitar integrante (HU-03): %v", err)
	}
	got, err := b.repo.BuscarRegistroPorID(ctx, r.ID)

	if err != nil {
		t.Fatalf("el registro tendría que seguir existiendo: %v", err)
	}
	if got.Integrante != "Ana" {
		t.Errorf("Integrante = %q, se esperaba \"Ana\"", got.Integrante)
	}
}

// registrarHoras guarda un registro directo en la base del elemento y lo deja a nombre del integrante.
func (b baseEsfuerzo) registrarHoras(t *testing.T, elementoID int64, integrante string, dia int, horas float64) {
	t.Helper()
	fecha := time.Date(2026, time.October, dia, 0, 0, 0, 0, time.Local)
	hoy := time.Date(2026, time.October, 31, 12, 0, 0, 0, time.Local)
	r, err := domain.NuevoRegistroEsfuerzo(b.proyectoID, elementoID, integrante, fecha, "Desarrollo", horas, hoy)
	if err != nil {
		t.Fatalf("crear registro de prueba: %v", err)
	}
	if err := b.repo.GuardarRegistro(context.Background(), &r); err != nil {
		t.Fatalf("guardar registro de prueba: %v", err)
	}
}

// horasReales lee las horas del elemento y falla el test si hay un error.
func (b baseEsfuerzo) horasReales(t *testing.T, elementoID int64) float64 {
	t.Helper()
	got, err := b.repo.HorasRealesDeElemento(context.Background(), elementoID)
	if err != nil {
		t.Fatalf("horas reales del elemento %d: %v", elementoID, err)
	}
	return got
}

// Traza: HU-19 CA-19.5, CA-19.6 (RN-19.4: sin registros son 0 horas, y no es un error)
func TestRepositorioEsfuerzo_HorasReales_SinRegistros_DaCero(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)

	if got := b.horasReales(t, b.elementoID); got != 0 {
		t.Errorf("horas reales = %v, se esperaba 0", got)
	}
}

// Traza: HU-19 CA-19.1 (RN-19.4: de cualquier integrante y fecha)
func TestRepositorioEsfuerzo_HorasReales_SumaIntegrantesYFechasDistintos(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)
	if err := b.proyectos.AgregarIntegrante(context.Background(), b.proyectoID, domain.Integrante{Nombre: "Beto", Roles: []string{"Tester"}}); err != nil {
		t.Fatal(err)
	}
	b.registrarHoras(t, b.elementoID, "Ana", 20, 2)
	b.registrarHoras(t, b.elementoID, "Beto", 21, 0.5)
	b.registrarHoras(t, b.elementoID, "Ana", 21, 1.5)

	if got := b.horasReales(t, b.elementoID); got != 4 {
		t.Errorf("horas reales = %v, se esperaba 4", got)
	}
}

// Traza: HU-19 CA-19.1 (RN-19.4; sección 5): se suman enteros, y una suma impar de medias horas conserva el 0,5.
func TestRepositorioEsfuerzo_HorasReales_ConservaLaMediaHora(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)
	b.registrarHoras(t, b.elementoID, "Ana", 20, 0.5)
	b.registrarHoras(t, b.elementoID, "Ana", 20, 1)

	if got := b.horasReales(t, b.elementoID); got != 1.5 {
		t.Errorf("horas reales = %v, se esperaba 1,5", got)
	}
}

// Traza: HU-19 CA-19.11 (RN-19.4: el tope de 24 h es por registro, no para la suma)
func TestRepositorioEsfuerzo_HorasReales_LaSumaPuedeSuperar24(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)
	b.registrarHoras(t, b.elementoID, "Ana", 20, 24)
	b.registrarHoras(t, b.elementoID, "Ana", 21, 24)
	b.registrarHoras(t, b.elementoID, "Ana", 22, 8)

	if got := b.horasReales(t, b.elementoID); got != 56 {
		t.Errorf("horas reales = %v, se esperaba 56", got)
	}
}

// Traza: HU-19 CA-19.11 (RN-19.4)
func TestRepositorioEsfuerzo_HorasReales_NoSumaOtroElemento(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)
	otro := elementoDePrueba(t, b.proyectoID, "Reportes", nil)
	if err := NuevoRepositorioBacklog(b.db).GuardarElemento(context.Background(), &otro); err != nil {
		t.Fatal(err)
	}
	b.registrarHoras(t, b.elementoID, "Ana", 20, 2)
	b.registrarHoras(t, otro.ID, "Ana", 20, 5)

	if got := b.horasReales(t, b.elementoID); got != 2 {
		t.Errorf("horas reales del elemento = %v, se esperaba 2", got)
	}
	if got := b.horasReales(t, otro.ID); got != 5 {
		t.Errorf("horas reales del otro elemento = %v, se esperaba 5", got)
	}
}

// Traza: HU-19 CA-19.11 (RN-19.4, RN-18.3)
func TestRepositorioEsfuerzo_HorasReales_IntegranteQuitado_SigueSumando(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t)
	if err := b.proyectos.AgregarIntegrante(ctx, b.proyectoID, domain.Integrante{Nombre: "Beto", Roles: []string{"Tester"}}); err != nil {
		t.Fatal(err)
	}
	b.registrarHoras(t, b.elementoID, "Ana", 20, 2)
	b.registrarHoras(t, b.elementoID, "Beto", 20, 3)

	if err := b.proyectos.QuitarIntegrante(ctx, b.proyectoID, "Beto"); err != nil {
		t.Fatalf("quitar integrante (HU-03): %v", err)
	}

	if got := b.horasReales(t, b.elementoID); got != 5 {
		t.Errorf("horas reales = %v, se esperaba 5", got)
	}
}

// El caso de uso ya verificó el elemento; si no existe, la consulta da 0 y no un error (spec, sección 5).
// Traza: HU-19 CA-19.5 (RN-19.4; sección 5: consulta de horas reales)
func TestRepositorioEsfuerzo_HorasReales_ElementoInexistente_DaCero(t *testing.T) {
	b := nuevaBaseEsfuerzo(t)

	if got := b.horasReales(t, 999); got != 0 {
		t.Errorf("horas reales = %v, se esperaba 0", got)
	}
}
