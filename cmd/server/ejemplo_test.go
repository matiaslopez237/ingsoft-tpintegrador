package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/ejemplo"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
)

// Traza: HU-31 RN-31.5: cargarEjemplo arma los casos de uso que dependen de "hoy" con el reloj de la
// carga, no con time.Now. Si usara la hora real, DEF-01 (26/10) no caería en el Sprint 2 y DEF-04
// (3/11, entre sprints) tendría el sprint abierto del día.
func TestCargarEjemplo_UsaElRelojDelEjemplo(t *testing.T) {
	ctx := context.Background()
	db, err := persistencia.AbrirDB(filepath.Join(t.TempDir(), "ejemplo.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	r, err := cargarEjemplo(ctx, db, "true")
	if err != nil || r != ejemplo.Cargado {
		t.Fatalf("cargarEjemplo = %v, %v; se esperaba Cargado sin error", r, err)
	}

	proyectos := persistencia.NuevoRepositorioProyectos(db)
	sprints, err := app.NuevoListarSprints(proyectos, persistencia.NuevoRepositorioSprints(db)).Ejecutar(ctx, 1)
	if err != nil {
		t.Fatalf("listar sprints: %v", err)
	}
	numero := map[int64]int{}
	for _, s := range sprints {
		numero[s.ID] = s.Numero
	}
	defectos, err := app.NuevoListarDefectos(proyectos, persistencia.NuevoRepositorioDefectos(db)).Ejecutar(ctx, 1)
	if err != nil {
		t.Fatalf("listar defectos: %v", err)
	}
	deteccion := map[string]int{}
	for _, d := range defectos {
		clave := strings.SplitN(d.Descripcion, ":", 2)[0]
		if d.SprintDeteccionID != nil {
			deteccion[clave] = numero[*d.SprintDeteccionID]
		}
	}
	if deteccion["DEF-01"] != 2 {
		t.Errorf("DEF-01 se detectó en el sprint %d, se esperaba el 2", deteccion["DEF-01"])
	}
	if n, ok := deteccion["DEF-04"]; ok {
		t.Errorf("DEF-04 tiene sprint de detección %d, se esperaba ninguno", n)
	}
}

// Traza: HU-31 P-2, sección 6: con Ctrl+C (o SIGTERM) durante la carga, el contexto se cancela y la
// carga se corta con un error, en lugar de seguir hasta el final.
func TestCargarEjemplo_ContextoCancelado_CortaLaCarga(t *testing.T) {
	db, err := persistencia.AbrirDB(filepath.Join(t.TempDir(), "ejemplo.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	ctx, cancelar := context.WithCancel(context.Background())
	cancelar()

	_, err = cargarEjemplo(ctx, db, "true")

	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, se esperaba que la carga se corte por el contexto cancelado", err)
	}
}

// Traza: HU-31 P-2, sección 7: si la carga falla, la base queda con el ejemplo a medias y el próximo
// arranque no lo carga (RN-31.1). El mensaje dice qué archivo borrar.
func TestMensajeCargaFallida_PideBorrarLaBase(t *testing.T) {
	err := errors.New("cargar ejemplo: cerrar el sprint 2: falla")

	got := mensajeCargaFallida("/datos/metrics.db", err)

	for _, want := range []string{"cerrar el sprint 2: falla", "a medias", "borrá /datos/metrics.db antes de volver a arrancar"} {
		if !strings.Contains(got, want) {
			t.Errorf("mensaje = %q, falta %q", got, want)
		}
	}
}
