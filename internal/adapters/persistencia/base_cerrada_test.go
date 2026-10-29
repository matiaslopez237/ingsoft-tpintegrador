package persistencia

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: hueco B-3 de la trazabilidad (ramas defensivas de persistencia sin test). Con la base
// cerrada, toda consulta falla: cada método tiene que devolver el error, envuelto pero sin
// perder la causa, en lugar de tragárselo y devolver datos vacíos como si no hubiera pasado nada.
func TestRepositorios_BaseCerrada_DevuelvenElError(t *testing.T) {
	db, err := AbrirDB(filepath.Join(t.TempDir(), "cerrada.db"))
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	db.Close()

	ctx := context.Background()
	proyectos, backlog := NuevoRepositorioProyectos(db), NuevoRepositorioBacklog(db)
	sprints, esfuerzo, defectos := NuevoRepositorioSprints(db), NuevoRepositorioEsfuerzo(db), NuevoRepositorioDefectos(db)
	nombre, titulo, sprintID := "Beta", "Otro", int64(1)
	dia := time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC)
	proyecto := domain.Proyecto{ID: 1, Nombre: "Alfa", FechaInicio: dia, Estado: domain.EstadoActivo}
	elemento := domain.ElementoBacklog{ID: 1, ProyectoID: 1, Titulo: "Login", Descripcion: "d",
		Prioridad: domain.PrioridadAlta, Estado: domain.EstadoPendiente, Criterios: []string{"c"}, SprintID: &sprintID}
	sprint := domain.Sprint{ID: 1, ProyectoID: 1, FechaInicio: dia, FechaFin: dia, Objetivo: "o", Estado: domain.SprintAbierto}

	metodos := map[string]func() error{
		// Proyectos (HU-01 a HU-04, HU-30)
		"ExisteNombre":    func() error { _, err := proyectos.ExisteNombre(ctx, "alfa"); return err },
		"Guardar":         func() error { p := proyecto; return proyectos.Guardar(ctx, &p) },
		"Actualizar":      func() error { return proyectos.Actualizar(ctx, proyecto, domain.CambiosProyecto{Nombre: &nombre}) },
		"BuscarPorID":     func() error { _, err := proyectos.BuscarPorID(ctx, 1); return err },
		"ListarProyectos": func() error { _, err := proyectos.ListarProyectos(ctx); return err },
		"AgregarIntegrante": func() error {
			return proyectos.AgregarIntegrante(ctx, 1, domain.Integrante{Nombre: "Ana", Roles: []string{"Dev"}})
		},
		"AgregarRol":       func() error { return proyectos.AgregarRol(ctx, 1, "Ana", "QA") },
		"QuitarIntegrante": func() error { return proyectos.QuitarIntegrante(ctx, 1, "Ana") },
		"Contar":           func() error { _, err := proyectos.Contar(ctx); return err },
		// Backlog (HU-05 a HU-07, HU-09, HU-10)
		"ExisteProyecto":      func() error { _, err := backlog.ExisteProyecto(ctx, 1); return err },
		"ExisteTitulo":        func() error { _, err := backlog.ExisteTitulo(ctx, 1, "login"); return err },
		"GuardarElemento":     func() error { e := elemento; return backlog.GuardarElemento(ctx, &e) },
		"BuscarElementoPorID": func() error { _, err := backlog.BuscarElementoPorID(ctx, 1); return err },
		"ContarElementos":     func() error { _, err := backlog.ContarElementos(ctx, 1); return err },
		"ListarElementos":     func() error { _, err := backlog.ListarElementos(ctx, 1); return err },
		"ActualizarElemento": func() error {
			return backlog.ActualizarElemento(ctx, elemento, domain.CambiosElemento{Titulo: &titulo})
		},
		"TieneEsfuerzo":     func() error { _, err := backlog.TieneEsfuerzo(ctx, 1); return err },
		"TieneDefectos":     func() error { _, err := backlog.TieneDefectos(ctx, 1); return err },
		"EliminarElemento":  func() error { return backlog.EliminarElemento(ctx, 1) },
		"BuscarSprint":      func() error { _, _, err := backlog.BuscarSprint(ctx, 1); return err },
		"GuardarAsignacion": func() error { e := elemento; return backlog.GuardarAsignacion(ctx, &e) },
		"GuardarCompletado": func() error { e := elemento; return backlog.GuardarCompletado(ctx, &e) },
		// Sprints (HU-08, HU-11)
		"GuardarSprint":     func() error { s := sprint; return sprints.GuardarSprint(ctx, &s) },
		"ListarSprints":     func() error { _, err := sprints.ListarSprints(ctx, 1); return err },
		"BuscarSprintPorID": func() error { _, err := sprints.BuscarSprintPorID(ctx, 1); return err },
		"CerrarSprint":      func() error { s := sprint; return sprints.CerrarSprint(ctx, &s) },
		// Esfuerzo (HU-18) y defectos (HU-20)
		"GuardarRegistro": func() error {
			r := domain.RegistroEsfuerzo{ProyectoID: 1, ElementoID: 1, Integrante: "Ana", Fecha: dia, Actividad: "Dev", Horas: 1}
			return esfuerzo.GuardarRegistro(ctx, &r)
		},
		"BuscarRegistroPorID": func() error { _, err := esfuerzo.BuscarRegistroPorID(ctx, 1); return err },
		"GuardarDefecto": func() error {
			d := domain.Defecto{ProyectoID: 1, ElementoID: 1, Descripcion: "d", Severidad: domain.SeveridadAlta, Estado: domain.DefectoAbierto}
			return defectos.GuardarDefecto(ctx, &d)
		},
	}
	for nombre, metodo := range metodos {
		t.Run(nombre, func(t *testing.T) {
			err := metodo()
			if err == nil {
				t.Fatal("se esperaba un error con la base cerrada, pero devolvió nil")
			}
			if !strings.Contains(err.Error(), "database is closed") {
				t.Errorf("error = %v; se esperaba que conservara la causa (base cerrada)", err)
			}
		})
	}
}

// Traza: hueco B-3. AbrirDB informa el error si no puede abrir el archivo o si el archivo no es
// una base SQLite, en lugar de arrancar el servidor con una base inutilizable.
func TestAbrirDB_RutaInvalida_RetornaError(t *testing.T) {
	dir := t.TempDir()
	noEsBase := filepath.Join(dir, "texto.db")
	if err := os.WriteFile(noEsBase, []byte(strings.Repeat("esto no es una base SQLite\n", 100)), 0o600); err != nil {
		t.Fatal(err)
	}

	casos := map[string]string{
		"carpeta inexistente":    filepath.Join(dir, "no", "existe", "metrics.db"),
		"archivo que no es base": noEsBase,
	}
	for nombre, ruta := range casos {
		t.Run(nombre, func(t *testing.T) {
			db, err := AbrirDB(ruta)
			if err == nil {
				db.Close()
				t.Fatal("se esperaba un error, pero AbrirDB devolvió la base")
			}
		})
	}
}
