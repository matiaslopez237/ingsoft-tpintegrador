package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador también cumple el puerto angosto de HU-17 (se verifica al compilar).
var _ app.RepositorioEstimacion = (*RepositorioRondas)(nil)

// rondaCerradaGuardada abre una ronda sobre "Login" y la cierra con Ana y Beto. Devuelve la ronda
// como quedó guardada.
func (b basePoker) rondaCerradaGuardada(t *testing.T, anterior *domain.Ronda) domain.Ronda {
	t.Helper()
	r := b.abrir(t, anterior, "Ana", "Beto")
	for _, nombre := range []string{"Ana", "Beto"} {
		if err := b.votar(t, &r, nombre, 5); err != nil {
			t.Fatal(err)
		}
	}
	guardada, err := b.repo.BuscarRonda(context.Background(), r.ID)
	if err != nil {
		t.Fatal(err)
	}
	return guardada
}

// conEstimacion devuelve la ronda con el valor, lista para GuardarEstimacion.
func conEstimacion(r domain.Ronda, valor int) domain.Ronda {
	r.EstimacionAcordada = &valor
	return r
}

// Traza: HU-17 CA-17.1, CA-17.2 (RN-17.1): la ronda guarda la estimación y el elemento pasa a tener
// esos story points, sin cambiar sus demás datos. El historial también la trae.
func TestRepositorioRondas_GuardarEstimacion(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.rondaCerradaGuardada(t, nil)
	antes := buscarElementoDePrueba(t, b.backlog, b.elemento.ID)

	if err := b.repo.GuardarEstimacion(ctx, conEstimacion(r, 8)); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	got, err := b.repo.BuscarRonda(ctx, r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.EstimacionAcordada == nil || *got.EstimacionAcordada != 8 {
		t.Errorf("EstimacionAcordada = %v, se esperaba 8", got.EstimacionAcordada)
	}
	historial, err := b.repo.ListarRondas(ctx, b.elemento.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(historial) != 1 || historial[0].EstimacionAcordada == nil || *historial[0].EstimacionAcordada != 8 {
		t.Errorf("historial = %+v, se esperaba la ronda con la estimación 8", historial)
	}
	despues := buscarElementoDePrueba(t, b.backlog, b.elemento.ID)
	if despues.StoryPoints == nil || *despues.StoryPoints != 8 {
		t.Errorf("StoryPoints = %v, se esperaba 8", despues.StoryPoints)
	}
	despues.StoryPoints = antes.StoryPoints
	if !reflect.DeepEqual(despues, antes) {
		t.Errorf("el elemento cambió: %+v, se esperaba %+v salvo los story points", despues, antes)
	}
}

// Traza: HU-17 CA-17.6 (RN-17.4): cada ronda conserva lo suyo y el elemento queda con la última.
func TestRepositorioRondas_GuardarEstimacion_RondaNueva(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r1 := b.rondaCerradaGuardada(t, nil)
	if err := b.repo.GuardarEstimacion(ctx, conEstimacion(r1, 5)); err != nil {
		t.Fatal(err)
	}
	r2 := b.rondaCerradaGuardada(t, &r1)

	if err := b.repo.GuardarEstimacion(ctx, conEstimacion(r2, 8)); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if got, _ := b.repo.BuscarRonda(ctx, r1.ID); got.EstimacionAcordada == nil || *got.EstimacionAcordada != 5 {
		t.Errorf("la ronda 1 tendría que conservar el 5: %v", got.EstimacionAcordada)
	}
	if e := buscarElementoDePrueba(t, b.backlog, b.elemento.ID); e.StoryPoints == nil || *e.StoryPoints != 8 {
		t.Errorf("StoryPoints = %v, se esperaba 8", e.StoryPoints)
	}
}

// Traza: HU-17 sección 5: las reglas se verifican de nuevo al guardar, porque otro pedido pudo cambiar
// algo después de la validación del dominio. Si alguna falla, no se guarda nada.
func TestRepositorioRondas_GuardarEstimacion_Errores(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name  string
		armar func(*testing.T, basePoker) domain.Ronda
		want  error
	}{
		{"ronda inexistente", func(*testing.T, basePoker) domain.Ronda { return domain.Ronda{ID: 999} }, domain.ErrRondaNoEncontrada},
		{"ronda abierta", func(t *testing.T, b basePoker) domain.Ronda { return b.abrir(t, nil, "Ana", "Beto") }, domain.ErrEstimacionRondaAbierta},
		{"no es la última", func(t *testing.T, b basePoker) domain.Ronda {
			r1 := b.rondaCerradaGuardada(t, nil)
			b.abrir(t, &r1, "Ana", "Beto")
			return r1
		}, domain.ErrRondaNoEsLaUltima},
		{"ya tiene estimación", func(t *testing.T, b basePoker) domain.Ronda {
			r := b.rondaCerradaGuardada(t, nil)
			if err := b.repo.GuardarEstimacion(ctx, conEstimacion(r, 5)); err != nil {
				t.Fatal(err)
			}
			return r
		}, domain.ErrEstimacionYaRegistrada},
		{"elemento Hecho", func(t *testing.T, b basePoker) domain.Ronda {
			r := b.rondaCerradaGuardada(t, nil)
			if _, err := b.repo.db.Exec(`UPDATE backlog_elementos SET estado = ? WHERE id = ?`, string(domain.EstadoHecho), b.elemento.ID); err != nil {
				t.Fatal(err)
			}
			return r
		}, domain.ErrElementoTerminado},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := nuevaBasePoker(t)
			r := tt.armar(t, b)
			puntosAntes := buscarElementoDePrueba(t, b.backlog, b.elemento.ID).StoryPoints
			estimacionAntes := sql.NullInt64{}
			_ = b.repo.db.QueryRow(`SELECT estimacion_acordada FROM rondas WHERE id = ?`, r.ID).Scan(&estimacionAntes)

			err := b.repo.GuardarEstimacion(ctx, conEstimacion(r, 8))

			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, se esperaba %v", err, tt.want)
			}
			if got := buscarElementoDePrueba(t, b.backlog, b.elemento.ID).StoryPoints; !reflect.DeepEqual(got, puntosAntes) {
				t.Errorf("StoryPoints = %v, se esperaba sin cambios %v", got, puntosAntes)
			}
			estimacionDespues := sql.NullInt64{}
			_ = b.repo.db.QueryRow(`SELECT estimacion_acordada FROM rondas WHERE id = ?`, r.ID).Scan(&estimacionDespues)
			if estimacionDespues != estimacionAntes {
				t.Errorf("la estimación de la ronda cambió de %v a %v", estimacionAntes, estimacionDespues)
			}
		})
	}
}

// Traza: HU-17: guardar una ronda que no trae la estimación es un error de programación.
func TestRepositorioRondas_GuardarEstimacion_SinValor(t *testing.T) {
	b := nuevaBasePoker(t)
	r := b.rondaCerradaGuardada(t, nil)

	if err := b.repo.GuardarEstimacion(context.Background(), r); err == nil {
		t.Error("se esperaba un error")
	}
}

// Traza: HU-17 sección 5: si varios pedidos registran a la vez, se guarda uno solo, y el valor de la
// ronda y el del elemento son el mismo.
func TestRepositorioRondas_GuardarEstimacion_Simultaneos(t *testing.T) {
	ctx := context.Background()
	b := nuevaBasePoker(t)
	r := b.rondaCerradaGuardada(t, nil)
	valores := []int{1, 2, 3, 5, 8, 13, 21, 1, 2, 3}

	var (
		wg      sync.WaitGroup
		largada = make(chan struct{})
		errs    = make([]error, len(valores))
	)
	for i, v := range valores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-largada
			errs[i] = b.repo.GuardarEstimacion(ctx, conEstimacion(r, v))
		}()
	}
	close(largada)
	wg.Wait()

	exitos := 0
	for i, err := range errs {
		switch {
		case err == nil:
			exitos++
		case !errors.Is(err, domain.ErrEstimacionYaRegistrada):
			t.Errorf("pedido %d: error = %v, se esperaba nil o %v", i, err, domain.ErrEstimacionYaRegistrada)
		}
	}
	if exitos != 1 {
		t.Errorf("registraron %d pedidos, se esperaba exactamente 1", exitos)
	}
	got, err := b.repo.BuscarRonda(ctx, r.ID)
	if err != nil {
		t.Fatal(err)
	}
	e := buscarElementoDePrueba(t, b.backlog, b.elemento.ID)
	if got.EstimacionAcordada == nil || e.StoryPoints == nil || *got.EstimacionAcordada != *e.StoryPoints {
		t.Errorf("ronda = %v y elemento = %v: se esperaba el mismo valor", got.EstimacionAcordada, e.StoryPoints)
	}
}

// Traza: HU-17 sección 5: una base con la tabla rondas de HU-14 recibe la columna, sin perder datos.
func TestAbrirDB_BaseExistenteSinEstimacionAcordada_AgregaLaColumna(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "vieja.db")
	crearEsquemaSinSprintID(t, ruta)
	vieja, err := sql.Open("sqlite", ruta)
	if err != nil {
		t.Fatal(err)
	}
	_, err = vieja.Exec(`
CREATE TABLE rondas (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	proyecto_id INTEGER NOT NULL REFERENCES proyectos (id),
	elemento_id INTEGER NOT NULL REFERENCES backlog_elementos (id) ON DELETE CASCADE,
	numero      INTEGER NOT NULL,
	estado      TEXT    NOT NULL,
	UNIQUE (elemento_id, numero)
);
INSERT INTO proyectos (nombre, nombre_clave, fecha_inicio, estado) VALUES ('Alfa', 'alfa', '2026-10-01', 'Activo');
INSERT INTO backlog_elementos (proyecto_id, titulo, titulo_clave, descripcion, prioridad, estado)
VALUES (1, 'Login', 'login', 'd', 'Alta', 'Pendiente');
INSERT INTO rondas (proyecto_id, elemento_id, numero, estado) VALUES (1, 1, 1, 'Cerrada');`)
	vieja.Close()
	if err != nil {
		t.Fatalf("crear rondas viejas: %v", err)
	}

	db, err := AbrirDB(ruta)
	if err != nil {
		t.Fatalf("abrir base vieja con la migración actual: %v", err)
	}
	defer db.Close()

	if !tieneColumna(t, db, "rondas", "estimacion_acordada") {
		t.Fatal("se esperaba la columna estimacion_acordada después de migrar")
	}
	r, err := NuevoRepositorioRondas(db).BuscarRonda(context.Background(), 1)
	if err != nil {
		t.Fatalf("leer la ronda vieja: %v", err)
	}
	if r.Numero != 1 || r.EstimacionAcordada != nil {
		t.Errorf("ronda vieja = %+v, se esperaba la misma, sin estimación", r)
	}
}
