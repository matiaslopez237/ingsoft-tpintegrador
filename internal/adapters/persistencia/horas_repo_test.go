package persistencia

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir los puertos de HU-23 (se verifica al compilar).
var (
	_ app.LectorHoras       = (*RepositorioEsfuerzo)(nil)
	_ app.LectorSprintPorID = (*RepositorioSprints)(nil)
)

// elementoConHoras guarda un elemento en el proyecto de la base, con esas horas estimadas (nil: sin estimar).
func (b baseEsfuerzo) elementoConHoras(t *testing.T, titulo string, estimadas *float64) domain.ElementoBacklog {
	t.Helper()
	e := elementoDePrueba(t, b.proyectoID, titulo, nil)
	e.HorasEstimadas = estimadas
	return guardarElementoDePrueba(t, NuevoRepositorioBacklog(b.db), e)
}

// sprint guarda un sprint Abierto en el proyecto de la base, del día desde al día hasta de octubre.
func (b baseEsfuerzo) sprint(t *testing.T, desde, hasta int) domain.Sprint {
	t.Helper()
	return guardarSprintDePrueba(t, NuevoRepositorioSprints(b.db), sprintDePrueba(b.proyectoID, desde, hasta))
}

// asignar deja el elemento en el sprint.
func (b baseEsfuerzo) asignar(t *testing.T, e domain.ElementoBacklog, s domain.Sprint) domain.ElementoBacklog {
	t.Helper()
	e.SprintID = &s.ID
	if err := NuevoRepositorioBacklog(b.db).GuardarAsignacion(context.Background(), &e); err != nil {
		t.Fatalf("asignar %q: %v", e.Titulo, err)
	}
	return e
}

// otroProyecto devuelve la misma base, parada en un proyecto nuevo y sin elementos.
func (b baseEsfuerzo) otroProyecto(t *testing.T, nombre string) baseEsfuerzo {
	t.Helper()
	p := guardarDePrueba(t, b.proyectos, nombre)
	b.proyectoID = p.ID
	return b
}

// sinElementos devuelve la misma base, parada en un proyecto sin elementos.
func (b baseEsfuerzo) sinElementos(t *testing.T) baseEsfuerzo {
	t.Helper()
	return b.otroProyecto(t, "Vacío")
}

func (b baseEsfuerzo) horasDeProyecto(t *testing.T) []domain.HorasElemento {
	t.Helper()
	got, err := b.repo.HorasDeProyecto(context.Background(), b.proyectoID)
	if err != nil {
		t.Fatalf("horas del proyecto %d: %v", b.proyectoID, err)
	}
	return got
}

func (b baseEsfuerzo) horasDeSprint(t *testing.T, sprintID int64) []domain.HorasElemento {
	t.Helper()
	got, err := b.repo.HorasDeSprint(context.Background(), sprintID)
	if err != nil {
		t.Fatalf("horas del sprint %d: %v", sprintID, err)
	}
	return got
}

func fila(estimadas *float64, reales float64) domain.HorasElemento {
	return domain.HorasElemento{HorasEstimadas: estimadas, HorasReales: reales}
}

// mismasFilas compara las filas en orden: los punteros, por el valor al que apuntan.
func mismasFilas(t *testing.T, got, want []domain.HorasElemento) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%d filas, se esperaban %d: %s", len(got), len(want), verFilas(got))
	}
	for i := range want {
		if !mismasHoras(got[i].HorasEstimadas, want[i].HorasEstimadas) || got[i].HorasReales != want[i].HorasReales {
			t.Errorf("filas = %s, se esperaba %s", verFilas(got), verFilas(want))
			return
		}
	}
}

// verFilas muestra las filas en los mensajes de error: los punteros no se leen en %v.
func verFilas(filas []domain.HorasElemento) string {
	partes := make([]string, len(filas))
	for i, f := range filas {
		partes[i] = fmt.Sprintf("%v / %v", verHoras(f.HorasEstimadas), f.HorasReales)
	}
	return "[" + strings.Join(partes, "; ") + "]"
}

// Traza: HU-23 CA-23.9 (RN-23.2; sección 5: LEFT JOIN): sin registros, el elemento no desaparece.
func TestRepositorioEsfuerzo_HorasDeProyecto_ElementoSinRegistros_ApareceConCero(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	b.elementoConHoras(t, "Login", pHoras(8))
	b.elementoConHoras(t, "Reportes", nil)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 0), fila(nil, 0)})
}

// Traza: HU-23 CA-23.1 (RN-23.3): de cualquier integrante y fecha, una sola fila por elemento.
func TestRepositorioEsfuerzo_HorasDeProyecto_SumaIntegrantesYFechas_UnaFilaPorElemento(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	login := b.elementoConHoras(t, "Login", pHoras(8))
	b.registrarHoras(t, login.ID, "Ana", 20, 2)
	b.registrarHoras(t, login.ID, "Beto", 21, 3)
	b.registrarHoras(t, login.ID, "Ana", 22, 5)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 10)})
}

// Traza: HU-23 CA-23.6 (RN-23.2, RN-23.3): el sin estimar viene con nil y conserva sus reales.
func TestRepositorioEsfuerzo_HorasDeProyecto_EstimadoYSinEstimar(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	login := b.elementoConHoras(t, "Login", pHoras(8))
	reportes := b.elementoConHoras(t, "Reportes", nil)
	b.registrarHoras(t, login.ID, "Ana", 20, 10)
	b.registrarHoras(t, reportes.ID, "Ana", 20, 5)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 10), fila(nil, 5)})
}

// Traza: HU-23 CA-23.4 (sección 5: precisión): las medias horas se conservan en las dos columnas.
func TestRepositorioEsfuerzo_HorasDeProyecto_ConservaLasMediasHoras(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	login := b.elementoConHoras(t, "Login", pHoras(7.5))
	b.registrarHoras(t, login.ID, "Ana", 20, 0.5)
	b.registrarHoras(t, login.ID, "Ana", 21, 1)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(7.5), 1.5)})
}

// Traza: HU-23 CA-23.12 (RN-23.3): los registros de un elemento no se suman a otro.
func TestRepositorioEsfuerzo_HorasDeProyecto_NoMezclaLosRegistrosDeOtroElemento(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	login := b.elementoConHoras(t, "Login", pHoras(8))
	reportes := b.elementoConHoras(t, "Reportes", pHoras(4))
	b.registrarHoras(t, login.ID, "Ana", 20, 2)
	b.registrarHoras(t, reportes.ID, "Ana", 20, 5)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 2), fila(pHoras(4), 5)})
}

// Las filas vienen por el id del elemento, aunque los registros se hayan cargado en otro orden.
// Traza: HU-23 (sección 5: puerto de lectura)
func TestRepositorioEsfuerzo_HorasDeProyecto_OrdenaPorElemento(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	primero := b.elementoConHoras(t, "Primero", pHoras(1))
	segundo := b.elementoConHoras(t, "Segundo", pHoras(2))
	tercero := b.elementoConHoras(t, "Tercero", pHoras(3))
	b.registrarHoras(t, tercero.ID, "Ana", 20, 3)
	b.registrarHoras(t, primero.ID, "Ana", 21, 1)
	b.registrarHoras(t, segundo.ID, "Ana", 22, 2)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(1), 1), fila(pHoras(2), 2), fila(pHoras(3), 3)})
}

// Traza: HU-23 CA-23.2, CA-23.12 (RN-23.10): con sprint y sin sprint, y nada de otro proyecto.
func TestRepositorioEsfuerzo_HorasDeProyecto_IncluyeConYSinSprint_NoMezclaOtroProyecto(t *testing.T) {
	alfa := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := alfa.sprint(t, 26, 27)
	enSprint := alfa.asignar(t, alfa.elementoConHoras(t, "Login", pHoras(8)), s)
	sinSprint := alfa.elementoConHoras(t, "Pagos", pHoras(20))
	alfa.registrarHoras(t, enSprint.ID, "Ana", 26, 10)
	alfa.registrarHoras(t, sinSprint.ID, "Ana", 26, 15)

	beta := alfa.otroProyecto(t, "Beta")
	deBeta := beta.elementoConHoras(t, "Login", pHoras(6)) // mismo título que en Alfa
	beta.registrarHoras(t, deBeta.ID, "Caro", 26, 7)

	mismasFilas(t, alfa.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 10), fila(pHoras(20), 15)})
	mismasFilas(t, beta.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(6), 7)})
}

// Traza: HU-23 CA-23.8 (RN-23.6): sin elementos es una lista vacía, no un error.
func TestRepositorioEsfuerzo_HorasDeProyecto_SinElementos_ListaVacia(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)

	if got := b.horasDeProyecto(t); len(got) != 0 {
		t.Errorf("filas = %s, se esperaba una lista vacía", verFilas(got))
	}
}

// Traza: HU-23 CA-23.1, CA-23.12 (RN-23.7): solo los elementos del sprint, también los que no tienen
// registros; no los de otro sprint, ni los sin sprint, ni los de otro proyecto.
func TestRepositorioEsfuerzo_HorasDeSprint_SoloLosElementosDelSprint(t *testing.T) {
	alfa := nuevaBaseEsfuerzo(t).sinElementos(t)
	s, otro := alfa.sprint(t, 20, 21), alfa.sprint(t, 22, 23)
	login := alfa.asignar(t, alfa.elementoConHoras(t, "Login", pHoras(8)), s)
	alfa.asignar(t, alfa.elementoConHoras(t, "Ayuda", pHoras(1)), s)
	alfa.asignar(t, alfa.elementoConHoras(t, "SinEstimar", nil), s)
	deOtroSprint := alfa.asignar(t, alfa.elementoConHoras(t, "Reportes", pHoras(4)), otro)
	sinSprint := alfa.elementoConHoras(t, "Pagos", pHoras(20))
	alfa.registrarHoras(t, login.ID, "Ana", 20, 10)
	alfa.registrarHoras(t, deOtroSprint.ID, "Ana", 22, 9)
	alfa.registrarHoras(t, sinSprint.ID, "Ana", 20, 15)

	beta := alfa.otroProyecto(t, "Beta")
	sBeta := beta.sprint(t, 20, 21)
	deBeta := beta.asignar(t, beta.elementoConHoras(t, "Mapas", pHoras(6)), sBeta)
	beta.registrarHoras(t, deBeta.ID, "Caro", 20, 7)

	mismasFilas(t, alfa.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(pHoras(8), 10), fila(pHoras(1), 0), fila(nil, 0)})
	mismasFilas(t, alfa.horasDeSprint(t, otro.ID), []domain.HorasElemento{fila(pHoras(4), 9)})
	mismasFilas(t, beta.horasDeSprint(t, sBeta.ID), []domain.HorasElemento{fila(pHoras(6), 7)})
}

// La parte de CA-23.12 que el sistema en memoria no podía probar: la consulta no filtra por fecha, y
// los registros de un integrante quitado siguen sumando.
// Traza: HU-23 CA-23.12 (RN-23.3, RN-23.9)
func TestRepositorioEsfuerzo_HorasDeSprint_NoFiltraPorFechaNiPorIntegrante(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := b.sprint(t, 20, 21)
	login := b.asignar(t, b.elementoConHoras(t, "Login", pHoras(8)), s)
	if err := b.proyectos.AgregarIntegrante(ctx, b.proyectoID, domain.Integrante{Nombre: "Beto", Roles: []string{"Tester"}}); err != nil {
		t.Fatal(err)
	}
	b.registrarHoras(t, login.ID, "Ana", 1, 2)   // antes del sprint
	b.registrarHoras(t, login.ID, "Beto", 20, 3) // dentro del sprint
	b.registrarHoras(t, login.ID, "Ana", 30, 4)  // después del sprint
	if err := b.proyectos.QuitarIntegrante(ctx, b.proyectoID, "Beto"); err != nil {
		t.Fatalf("quitar integrante (HU-03): %v", err)
	}

	mismasFilas(t, b.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(pHoras(8), 9)})
	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 9)})
}

// Traza: HU-23 CA-23.8 (RN-23.6): sin elementos es una lista vacía; el sprint lo verificó el caso de uso.
func TestRepositorioEsfuerzo_HorasDeSprint_SinElementos_ListaVacia(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := b.sprint(t, 20, 21)

	if got := b.horasDeSprint(t, s.ID); len(got) != 0 {
		t.Errorf("filas = %s, se esperaba una lista vacía", verFilas(got))
	}
}

// Los story points no se convierten en horas, aunque el elemento tenga esfuerzo.
// Traza: HU-23 (sección 6; RN-23.2)
func TestRepositorioEsfuerzo_Horas_ElementoConStoryPointsYSinHoras_SigueSinEstimar(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := b.sprint(t, 20, 21)
	e := elementoDePrueba(t, b.proyectoID, "Login", puntosDePrueba(5))
	login := b.asignar(t, guardarElementoDePrueba(t, NuevoRepositorioBacklog(b.db), e), s)
	b.registrarHoras(t, login.ID, "Ana", 20, 3)

	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(nil, 3)})
	mismasFilas(t, b.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(nil, 3)})
}

// Con el caso de uso real de HU-09: el elemento quitado deja de contar en el sprint y sigue en el proyecto.
// Traza: HU-23 (sección 6; RN-23.7, RN-23.10)
func TestRepositorioEsfuerzo_HorasDeSprint_ElementoQuitado_DejaDeContar(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := b.sprint(t, 20, 21)
	login := b.asignar(t, b.elementoConHoras(t, "Login", pHoras(8)), s)
	b.asignar(t, b.elementoConHoras(t, "Reportes", pHoras(4)), s)
	b.registrarHoras(t, login.ID, "Ana", 20, 10)
	mismasFilas(t, b.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(pHoras(8), 10), fila(pHoras(4), 0)})

	err := app.NuevoDesasignarHistoria(NuevoRepositorioBacklog(b.db)).Ejecutar(context.Background(),
		app.DatosDesasignacion{ProyectoID: b.proyectoID, ElementoID: login.ID})

	if err != nil {
		t.Fatalf("quitar del sprint (HU-09): %v", err)
	}
	mismasFilas(t, b.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(pHoras(4), 0)})
	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 10), fila(pHoras(4), 0)})
}

// Con el caso de uso real de HU-06: el elemento eliminado deja de sumar sus estimadas. Que uno con
// esfuerzo no se pueda eliminar (RN-06.7) ya lo prueban los tests de HU-06.
// Traza: HU-23 (sección 6; RN-23.2, RN-23.10)
func TestRepositorioEsfuerzo_HorasDeProyecto_ElementoEliminado_DejaDeSumar(t *testing.T) {
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	login := b.elementoConHoras(t, "Login", pHoras(8))
	reportes := b.elementoConHoras(t, "Reportes", pHoras(4))
	b.registrarHoras(t, reportes.ID, "Ana", 20, 3)
	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 0), fila(pHoras(4), 3)})

	err := app.NuevoEliminarElementoBacklog(NuevoRepositorioBacklog(b.db)).Ejecutar(context.Background(),
		app.DatosEliminarElemento{ProyectoID: b.proyectoID, ElementoID: login.ID})

	if err != nil {
		t.Fatalf("eliminar elemento (HU-06): %v", err)
	}
	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(4), 3)})
}

// Con el cierre real de HU-11: el elemento sin terminar sale del sprint y sigue en el proyecto.
// Traza: HU-23 CA-23.10 (RN-23.7, RN-23.8, RN-23.10)
func TestRepositorioEsfuerzo_HorasDeSprint_TrasElCierre_SoloLoQueSigueAsignado(t *testing.T) {
	ctx := context.Background()
	b := nuevaBaseEsfuerzo(t).sinElementos(t)
	s := b.sprint(t, 20, 21)
	hecho := b.asignar(t, b.elementoConHoras(t, "Login", pHoras(8)), s)
	devuelto := b.asignar(t, b.elementoConHoras(t, "Reportes", pHoras(4)), s)
	b.registrarHoras(t, hecho.ID, "Ana", 20, 10)
	b.registrarHoras(t, devuelto.ID, "Ana", 20, 3)
	hecho.Estado = domain.EstadoHecho
	if err := NuevoRepositorioBacklog(b.db).GuardarCompletado(ctx, &hecho); err != nil {
		t.Fatalf("completar: %v", err)
	}
	if err := NuevoRepositorioSprints(b.db).CerrarSprint(ctx, &s); err != nil {
		t.Fatalf("cerrar sprint: %v", err)
	}

	mismasFilas(t, b.horasDeSprint(t, s.ID), []domain.HorasElemento{fila(pHoras(8), 10)})
	mismasFilas(t, b.horasDeProyecto(t), []domain.HorasElemento{fila(pHoras(8), 10), fila(pHoras(4), 3)})
}
