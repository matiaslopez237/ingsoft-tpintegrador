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
var _ app.RepositorioDefectos = (*RepositorioDefectos)(nil)

// baseDefecto es una base SQLite real con un proyecto y el elemento "Login".
type baseDefecto struct {
	db         *sql.DB
	repo       *RepositorioDefectos
	sprints    *RepositorioSprints
	proyectoID int64
	elementoID int64
}

func nuevaBaseDefecto(t *testing.T) baseDefecto {
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
	e := elementoDePrueba(t, p.ID, "Login", nil)
	if err := NuevoRepositorioBacklog(db).GuardarElemento(ctx, &e); err != nil {
		t.Fatalf("guardar elemento: %v", err)
	}
	return baseDefecto{
		db: db, repo: NuevoRepositorioDefectos(db), sprints: NuevoRepositorioSprints(db),
		proyectoID: p.ID, elementoID: e.ID,
	}
}

// defecto arma un defecto válido sobre el elemento de la base.
func (b baseDefecto) defecto(t *testing.T, severidad domain.Severidad) domain.Defecto {
	t.Helper()
	d, err := domain.NuevoDefecto(b.proyectoID, b.elementoID, "No carga la pantalla de login", severidad)
	if err != nil {
		t.Fatalf("crear defecto de prueba: %v", err)
	}
	return d
}

// contarDefectos cuenta los defectos del elemento directo en la base: solo lo necesitan los
// tests, así que no es parte del repositorio.
func contarDefectos(t *testing.T, db *sql.DB, elementoID int64) int {
	t.Helper()
	var n int
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM defectos WHERE elemento_id = ?`, elementoID).Scan(&n)
	if err != nil {
		t.Fatalf("contar defectos del elemento %d: %v", elementoID, err)
	}
	return n
}

// Traza: HU-20 CA-20.1, CA-20.7
func TestRepositorioDefectos_Guardar_AsignaIDsDistintos(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d1 := b.defecto(t, domain.SeveridadAlta)
	d2 := b.defecto(t, domain.SeveridadAlta)

	if err := b.repo.GuardarDefecto(ctx, &d1); err != nil {
		t.Fatal(err)
	}
	if err := b.repo.GuardarDefecto(ctx, &d2); err != nil {
		t.Fatal(err)
	}

	if d1.ID == 0 || d1.ID == d2.ID {
		t.Errorf("IDs = %d y %d, se esperaban asignados y distintos", d1.ID, d2.ID)
	}
	if n := contarDefectos(t, b.db, b.elementoID); n != 2 {
		t.Errorf("defectos del elemento = %d, se esperaban 2", n)
	}
}

// Traza: HU-20 CA-20.1 (RN-20.5: se guarda el Sprint de detección cuando hay uno)
func TestRepositorioDefectos_ConSprintDeteccion_LoGuarda(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	s := domain.Sprint{ProyectoID: b.proyectoID, FechaInicio: dia(20), FechaFin: dia(30),
		Objetivo: "Sprint de prueba", Estado: domain.SprintAbierto}
	if err := b.sprints.GuardarSprint(ctx, &s); err != nil {
		t.Fatalf("guardar sprint: %v", err)
	}
	d := b.defecto(t, domain.SeveridadMedia)
	d.SprintDeteccionID = &s.ID

	if err := b.repo.GuardarDefecto(ctx, &d); err != nil {
		t.Fatal(err)
	}

	var sprintID sql.NullInt64
	err := b.db.QueryRowContext(ctx, `SELECT sprint_deteccion_id FROM defectos WHERE id = ?`, d.ID).Scan(&sprintID)
	if err != nil {
		t.Fatal(err)
	}
	if !sprintID.Valid || sprintID.Int64 != s.ID {
		t.Errorf("sprint_deteccion_id = %v, se esperaba %d", sprintID, s.ID)
	}
}

// Traza: HU-20 CA-20.2 (RN-20.5: sin Sprint vigente, queda NULL)
func TestRepositorioDefectos_SinSprintDeteccion_GuardaNull(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d := b.defecto(t, domain.SeveridadBaja)

	if err := b.repo.GuardarDefecto(ctx, &d); err != nil {
		t.Fatal(err)
	}

	var sprintID sql.NullInt64
	err := b.db.QueryRowContext(ctx, `SELECT sprint_deteccion_id FROM defectos WHERE id = ?`, d.ID).Scan(&sprintID)
	if err != nil {
		t.Fatal(err)
	}
	if sprintID.Valid {
		t.Errorf("sprint_deteccion_id = %v, se esperaba NULL", sprintID)
	}
}

// Traza: HU-20 CA-20.5 (si el elemento ya no existe al guardar, la clave foránea lo rechaza)
func TestRepositorioDefectos_ElementoInexistente_RetornaErrElementoNoEncontrado(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d := b.defecto(t, domain.SeveridadAlta)
	d.ElementoID = 999

	err := b.repo.GuardarDefecto(ctx, &d)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if n := contarDefectos(t, b.db, 999); n != 0 {
		t.Errorf("defectos del elemento = %d, no se esperaban defectos", n)
	}
}

// Guarda la descripción, la severidad y el estado tal cual.
func TestRepositorioDefectos_GuardaDescripcionSeveridadYEstado(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseDefecto(t)
	d := b.defecto(t, domain.SeveridadCritica)

	if err := b.repo.GuardarDefecto(ctx, &d); err != nil {
		t.Fatal(err)
	}

	var descripcion, severidad, estado string
	err := b.db.QueryRowContext(ctx, `SELECT descripcion, severidad, estado FROM defectos WHERE id = ?`, d.ID).
		Scan(&descripcion, &severidad, &estado)
	if err != nil {
		t.Fatal(err)
	}
	if descripcion != d.Descripcion || severidad != "Crítica" || estado != "Abierto" {
		t.Errorf("descripcion=%q severidad=%q estado=%q; se esperaba %q, \"Crítica\", \"Abierto\"",
			descripcion, severidad, estado, d.Descripcion)
	}
}
