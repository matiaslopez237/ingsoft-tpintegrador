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

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.RepositorioEsfuerzo = (*RepositorioEsfuerzo)(nil)

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
	if n, err := b.repo.ContarRegistros(ctx, b.elementoID); err != nil || n != 2 {
		t.Errorf("ContarRegistros = %d, %v; se esperaban 2", n, err)
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
	if n, err := b.repo.ContarRegistros(ctx, 999); err != nil || n != 0 {
		t.Errorf("ContarRegistros = %d, %v; no se esperaban registros", n, err)
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
