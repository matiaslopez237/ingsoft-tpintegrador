package persistencia

import (
	"database/sql"
	"path/filepath"
	"sync"
	"testing"
)

// Traza: DEF-04: si dos arranques abren a la vez una base vieja, los dos ven en PRAGMA table_info que
// falta la columna, y el segundo en llegar al ALTER TABLE la encuentra ya agregada por el primero.
// Este test reproduce ese segundo arranque: la columna ya existe cuando se la quiere agregar. No es
// un error: la base quedó como se quería.
func TestAgregarColumna_YaLaAgregoOtroArranque_NoFalla(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "carrera.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, nueva INTEGER)`); err != nil {
		t.Fatalf("crear tabla: %v", err)
	}

	if err := agregarColumna(db, "t", "nueva", "INTEGER"); err != nil {
		t.Errorf("agregarColumna sobre una columna que ya agregó otro arranque = %v, se esperaba nil", err)
	}
}

// Traza: DEF-04: otro error del ALTER TABLE (por ejemplo, una tabla que no existe) se sigue informando.
func TestAgregarColumna_OtroError_SeInforma(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "carrera.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	defer db.Close()

	if err := agregarColumna(db, "no_existe", "nueva", "INTEGER"); err == nil {
		t.Error("agregarColumna sobre una tabla que no existe = nil, se esperaba un error")
	}
}

// Traza: DEF-04: varios arranques simultáneos sobre la misma base vieja terminan todos bien, y la base
// queda migrada.
func TestAbrirDB_ArranquesSimultaneosSobreUnaBaseVieja_NoFallan(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "vieja.db")
	crearEsquemaSinSprintID(t, ruta)

	const arranques = 6
	errs := make([]error, arranques)
	var wg sync.WaitGroup
	for i := range arranques {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db, err := AbrirDB(ruta)
			if err == nil {
				err = db.Close()
			}
			errs[i] = err
		}()
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("arranque %d: %v", i, err)
		}
	}
	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir la base migrada: %v", err)
	}
	defer db.Close()
	if !tieneColumna(t, db, "backlog_elementos", "sprint_id") {
		t.Error("se esperaba que la base quedara migrada")
	}
}
