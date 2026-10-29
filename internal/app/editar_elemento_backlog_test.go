package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// edicionFalsa guarda los elementos en memoria y permite simular fallas.
type edicionFalsa struct {
	elementos     map[int64]domain.ElementoBacklog
	conEsfuerzo   map[int64]bool
	abiertos      map[int64]bool // ID del sprint → si está Abierto (HU-09)
	actualizados  []domain.CambiosElemento
	eliminados    []int64
	errBuscar     error
	errTitulo     error
	errActualizar error
	errEsfuerzo   error
	errEliminar   error
	errSprint     error
}

func nuevaEdicionFalsa(elementos ...domain.ElementoBacklog) *edicionFalsa {
	f := &edicionFalsa{elementos: map[int64]domain.ElementoBacklog{}, conEsfuerzo: map[int64]bool{}, abiertos: map[int64]bool{}}
	for _, e := range elementos {
		f.elementos[e.ID] = e
	}
	return f
}

func (f *edicionFalsa) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if f.errBuscar != nil {
		return domain.ElementoBacklog{}, f.errBuscar
	}
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *edicionFalsa) ExisteTitulo(_ context.Context, proyectoID int64, clave string) (bool, error) {
	if f.errTitulo != nil {
		return false, f.errTitulo
	}
	for _, e := range f.elementos {
		if e.ProyectoID == proyectoID && domain.ClaveNombre(e.Titulo) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (f *edicionFalsa) ActualizarElemento(_ context.Context, e domain.ElementoBacklog, c domain.CambiosElemento) error {
	if f.errActualizar != nil {
		return f.errActualizar
	}
	f.actualizados = append(f.actualizados, c)
	f.elementos[e.ID] = e
	return nil
}

func (f *edicionFalsa) TieneEsfuerzo(_ context.Context, id int64) (bool, error) {
	return f.conEsfuerzo[id], f.errEsfuerzo
}

func (f *edicionFalsa) BuscarSprint(_ context.Context, id int64) (int64, bool, error) {
	if f.errSprint != nil {
		return 0, false, f.errSprint
	}
	return 1, f.abiertos[id], nil
}

func (f *edicionFalsa) EliminarElemento(_ context.Context, id int64) error {
	if f.errEliminar != nil {
		return f.errEliminar
	}
	f.eliminados = append(f.eliminados, id)
	delete(f.elementos, id)
	return nil
}

func elementoDePrueba(t *testing.T, id, proyectoID int64, titulo string) domain.ElementoBacklog {
	t.Helper()
	sp := 5
	e, err := domain.NuevoElementoBacklog(proyectoID, titulo, "Descripción", domain.PrioridadAlta, &sp, []string{"Criterio"})
	if err != nil {
		t.Fatal(err)
	}
	e.ID = id
	return e
}

func prioridad(p domain.Prioridad) *domain.Prioridad { return &p }
func titulo(s string) *string                        { return &s }

// Traza: HU-06 CA-06.1
func TestEditarElemento_DatosValidos_GuardaSoloLoEnviadoYDevuelveLoGuardado(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))
	cambios := domain.CambiosElemento{Prioridad: prioridad(domain.PrioridadBaja)}

	e, err := NuevoEditarElementoBacklog(repo).Ejecutar(context.Background(),
		DatosEditarElemento{ProyectoID: 1, ElementoID: 7, Cambios: cambios})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.ID != 7 || e.Prioridad != domain.PrioridadBaja || e.Titulo != "Login" {
		t.Errorf("elemento devuelto: %+v", e)
	}
	if len(repo.actualizados) != 1 || repo.actualizados[0].Prioridad == nil || repo.actualizados[0].Titulo != nil {
		t.Errorf("al repositorio se le pasaron los cambios %+v; se esperaba solo la prioridad", repo.actualizados)
	}
}

// Traza: HU-06 CA-06.7 (RN-06.6)
func TestEditarElemento_DeOtroProyectoOInexistente_NoEncontrado(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 2, "Pagos"))
	uc := NuevoEditarElementoBacklog(repo)
	cambios := domain.CambiosElemento{Prioridad: prioridad(domain.PrioridadBaja)}

	for _, d := range []DatosEditarElemento{
		{ProyectoID: 1, ElementoID: 7, Cambios: cambios},  // de otro proyecto
		{ProyectoID: 2, ElementoID: 99, Cambios: cambios}, // inexistente
	} {
		if _, err := uc.Ejecutar(context.Background(), d); !errors.Is(err, domain.ErrElementoNoEncontrado) {
			t.Errorf("%+v: error = %v, se esperaba %v", d, err, domain.ErrElementoNoEncontrado)
		}
	}
	if len(repo.actualizados) != 0 {
		t.Errorf("no se esperaban actualizaciones, hubo %d", len(repo.actualizados))
	}
}

// Traza: HU-06 CA-06.4, CA-06.6
func TestEditarElemento_DatosInvalidos_NoActualiza(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))
	uc := NuevoEditarElementoBacklog(repo)

	for _, tt := range []struct {
		c    domain.CambiosElemento
		want error
	}{
		{domain.CambiosElemento{}, domain.ErrSinCambios},
		{domain.CambiosElemento{Titulo: titulo("")}, domain.ErrTituloVacio},
	} {
		_, err := uc.Ejecutar(context.Background(), DatosEditarElemento{ProyectoID: 1, ElementoID: 7, Cambios: tt.c})
		if !errors.Is(err, tt.want) {
			t.Errorf("error = %v, se esperaba %v", err, tt.want)
		}
	}
	if len(repo.actualizados) != 0 {
		t.Errorf("no se esperaban actualizaciones, hubo %d", len(repo.actualizados))
	}
}

// Traza: HU-06 CA-06.5 (RN-06.2)
func TestEditarElemento_TituloDuplicado(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"), elementoDePrueba(t, 8, 1, "Registro"),
		elementoDePrueba(t, 9, 2, "Pagos"))
	uc := NuevoEditarElementoBacklog(repo)

	_, err := uc.Ejecutar(context.Background(), DatosEditarElemento{ProyectoID: 1, ElementoID: 7,
		Cambios: domain.CambiosElemento{Titulo: titulo("registro")}})
	if !errors.Is(err, domain.ErrTituloDuplicado) {
		t.Errorf("renombrar a un título existente: error = %v, se esperaba %v", err, domain.ErrTituloDuplicado)
	}

	for _, nuevo := range []string{"LOGIN", "Pagos"} { // el mismo título con otras mayúsculas, o uno de otro proyecto
		if _, err := uc.Ejecutar(context.Background(), DatosEditarElemento{ProyectoID: 1, ElementoID: 7,
			Cambios: domain.CambiosElemento{Titulo: titulo(nuevo)}}); err != nil {
			t.Errorf("renombrar a %q: error inesperado %v", nuevo, err)
		}
	}
}

func TestEditarElemento_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	cambios := domain.CambiosElemento{Titulo: titulo("Otro")}
	for nombre, preparar := range map[string]func(*edicionFalsa){
		"al buscar":           func(f *edicionFalsa) { f.errBuscar = errBase },
		"al verificar título": func(f *edicionFalsa) { f.errTitulo = errBase },
		"al actualizar":       func(f *edicionFalsa) { f.errActualizar = errBase },
	} {
		t.Run(nombre, func(t *testing.T) {
			repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))
			preparar(repo)
			_, err := NuevoEditarElementoBacklog(repo).Ejecutar(context.Background(),
				DatosEditarElemento{ProyectoID: 1, ElementoID: 7, Cambios: cambios})
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}

// Traza: HU-06 CA-06.8
func TestEliminarElemento_SinEsfuerzo_LoElimina(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))

	err := NuevoEliminarElementoBacklog(repo).Ejecutar(context.Background(), DatosEliminarElemento{ProyectoID: 1, ElementoID: 7})

	if err != nil || len(repo.eliminados) != 1 || repo.eliminados[0] != 7 {
		t.Errorf("err = %v, eliminados = %v", err, repo.eliminados)
	}
}

// Traza: HU-06 CA-06.9 (RN-06.7)
func TestEliminarElemento_ConEsfuerzo_NoLoElimina(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))
	repo.conEsfuerzo[7] = true

	err := NuevoEliminarElementoBacklog(repo).Ejecutar(context.Background(), DatosEliminarElemento{ProyectoID: 1, ElementoID: 7})

	if !errors.Is(err, domain.ErrElementoConEsfuerzo) || len(repo.eliminados) != 0 {
		t.Errorf("err = %v, eliminados = %v", err, repo.eliminados)
	}
}

// Traza: HU-06 CA-06.7
func TestEliminarElemento_DeOtroProyectoOInexistente_NoEncontrado(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 2, "Pagos"))
	uc := NuevoEliminarElementoBacklog(repo)

	for _, d := range []DatosEliminarElemento{{ProyectoID: 1, ElementoID: 7}, {ProyectoID: 2, ElementoID: 99}} {
		if err := uc.Ejecutar(context.Background(), d); !errors.Is(err, domain.ErrElementoNoEncontrado) {
			t.Errorf("%+v: error = %v, se esperaba %v", d, err, domain.ErrElementoNoEncontrado)
		}
	}
	if len(repo.eliminados) != 0 {
		t.Errorf("no se esperaban eliminaciones: %v", repo.eliminados)
	}
}

func TestEliminarElemento_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	for nombre, preparar := range map[string]func(*edicionFalsa){
		"al buscar":             func(f *edicionFalsa) { f.errBuscar = errBase },
		"al verificar esfuerzo": func(f *edicionFalsa) { f.errEsfuerzo = errBase },
		"al buscar el sprint":   func(f *edicionFalsa) { f.errSprint = errBase },
		"al eliminar":           func(f *edicionFalsa) { f.errEliminar = errBase },
	} {
		t.Run(nombre, func(t *testing.T) {
			repo := nuevaEdicionFalsa(enSprint(elementoDePrueba(t, 7, 1, "Login"), 10))
			preparar(repo)
			err := NuevoEliminarElementoBacklog(repo).Ejecutar(context.Background(), DatosEliminarElemento{ProyectoID: 1, ElementoID: 7})
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}

func enSprint(e domain.ElementoBacklog, sprintID int64) domain.ElementoBacklog {
	e.SprintID = &sprintID
	return e
}

// Traza: HU-06 CA-06.11 (RN-06.8). El Sprint activo se informa antes que el esfuerzo.
func TestEliminarElemento_EnSprintActivo_NoLoElimina(t *testing.T) {
	repo := nuevaEdicionFalsa(enSprint(elementoDePrueba(t, 7, 1, "Login"), 10))
	repo.abiertos[10] = true
	repo.conEsfuerzo[7] = true

	err := NuevoEliminarElementoBacklog(repo).Ejecutar(context.Background(), DatosEliminarElemento{ProyectoID: 1, ElementoID: 7})

	if !errors.Is(err, domain.ErrElementoEnSprintActivo) || len(repo.eliminados) != 0 {
		t.Errorf("err = %v, eliminados = %v", err, repo.eliminados)
	}
}

// Traza: HU-06 CA-06.11 (RN-06.8)
func TestEliminarElemento_EnSprintCerrado_LoElimina(t *testing.T) {
	repo := nuevaEdicionFalsa(enSprint(elementoDePrueba(t, 7, 1, "Login"), 10))

	err := NuevoEliminarElementoBacklog(repo).Ejecutar(context.Background(), DatosEliminarElemento{ProyectoID: 1, ElementoID: 7})

	if err != nil || len(repo.eliminados) != 1 {
		t.Errorf("err = %v, eliminados = %v", err, repo.eliminados)
	}
}
