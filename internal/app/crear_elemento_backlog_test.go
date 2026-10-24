package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// backlogFalso guarda los elementos en memoria y permite simular fallas.
type backlogFalso struct {
	proyectos   map[int64]bool
	guardados   []domain.ElementoBacklog
	siguienteID int64
	consultas   int
	errProyecto error
	errTitulo   error
	errGuardar  error
}

func nuevoBacklogFalso(proyectos ...int64) *backlogFalso {
	b := &backlogFalso{proyectos: map[int64]bool{}}
	for _, id := range proyectos {
		b.proyectos[id] = true
	}
	return b
}

func (b *backlogFalso) ExisteProyecto(_ context.Context, id int64) (bool, error) {
	b.consultas++
	return b.proyectos[id], b.errProyecto
}

func (b *backlogFalso) ExisteTitulo(_ context.Context, proyectoID int64, clave string) (bool, error) {
	b.consultas++
	if b.errTitulo != nil {
		return false, b.errTitulo
	}
	for _, e := range b.guardados {
		if e.ProyectoID == proyectoID && domain.ClaveNombre(e.Titulo) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (b *backlogFalso) GuardarElemento(_ context.Context, e *domain.ElementoBacklog) error {
	if b.errGuardar != nil {
		return b.errGuardar
	}
	b.siguienteID++
	e.ID = b.siguienteID
	b.guardados = append(b.guardados, *e)
	return nil
}

func datosElemento(proyectoID int64, titulo string) DatosNuevoElemento {
	sp := 5
	return DatosNuevoElemento{
		ProyectoID:  proyectoID,
		Titulo:      titulo,
		Descripcion: "Ingresar al sistema",
		Prioridad:   domain.PrioridadAlta,
		StoryPoints: &sp,
		Criterios:   []string{"Ingresa con email"},
	}
}

// Traza: HU-05 CA-05.1, CA-05.10
func TestCrearElementoBacklog_DatosValidos_GuardaPendienteConID(t *testing.T) {
	repo := nuevoBacklogFalso(1)
	uc := NuevoCrearElementoBacklog(repo)

	a, errA := uc.Ejecutar(context.Background(), datosElemento(1, "Login"))
	b, errB := uc.Ejecutar(context.Background(), datosElemento(1, "Registro"))

	if errA != nil || errB != nil {
		t.Fatalf("errores inesperados: %v, %v", errA, errB)
	}
	if a.ID == 0 || a.ID == b.ID {
		t.Errorf("IDs = %d y %d, se esperaban asignados y distintos", a.ID, b.ID)
	}
	if a.Estado != domain.EstadoPendiente || a.ProyectoID != 1 {
		t.Errorf("elemento inesperado: %+v", a)
	}
	if len(repo.guardados) != 2 {
		t.Errorf("se esperaban 2 elementos guardados, hay %d", len(repo.guardados))
	}
}

// Traza: HU-05 CA-05.3
func TestCrearElementoBacklog_ProyectoInexistente_RetornaErrorYNoGuarda(t *testing.T) {
	repo := nuevoBacklogFalso(1)

	// Aunque los datos sean inválidos, primero se informa que el proyecto no existe (sección 7).
	_, err := NuevoCrearElementoBacklog(repo).Ejecutar(context.Background(), datosElemento(2, ""))

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if len(repo.guardados) != 0 {
		t.Errorf("no se esperaban elementos guardados, hay %d", len(repo.guardados))
	}
}

// Traza: HU-05 CA-05.4
func TestCrearElementoBacklog_DatosInvalidos_NoConsultaTitulosNiGuarda(t *testing.T) {
	repo := nuevoBacklogFalso(1)

	_, err := NuevoCrearElementoBacklog(repo).Ejecutar(context.Background(), datosElemento(1, "  "))

	if !errors.Is(err, domain.ErrTituloVacio) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrTituloVacio)
	}
	if repo.consultas != 1 || len(repo.guardados) != 0 {
		t.Errorf("se esperaba solo la consulta del proyecto y nada guardado; consultas=%d guardados=%d",
			repo.consultas, len(repo.guardados))
	}
}

// Traza: HU-05 CA-05.9
func TestCrearElementoBacklog_TituloDuplicado(t *testing.T) {
	repo := nuevoBacklogFalso(1, 2)
	uc := NuevoCrearElementoBacklog(repo)
	if _, err := uc.Ejecutar(context.Background(), datosElemento(1, "Login")); err != nil {
		t.Fatalf("preparar: %v", err)
	}

	for _, titulo := range []string{"Login", "login", "  LOGIN "} {
		if _, err := uc.Ejecutar(context.Background(), datosElemento(1, titulo)); !errors.Is(err, domain.ErrTituloDuplicado) {
			t.Errorf("%q en el mismo proyecto: error = %v, se esperaba %v", titulo, err, domain.ErrTituloDuplicado)
		}
	}
	if _, err := uc.Ejecutar(context.Background(), datosElemento(2, "Login")); err != nil {
		t.Errorf("mismo título en otro proyecto: error inesperado %v", err)
	}
	if len(repo.guardados) != 2 {
		t.Errorf("se esperaban 2 elementos guardados, hay %d", len(repo.guardados))
	}
}

func TestCrearElementoBacklog_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		repo   *backlogFalso
	}{
		{"al buscar el proyecto", func() *backlogFalso { r := nuevoBacklogFalso(1); r.errProyecto = errBase; return r }()},
		{"al verificar el título", func() *backlogFalso { r := nuevoBacklogFalso(1); r.errTitulo = errBase; return r }()},
		{"al guardar", func() *backlogFalso { r := nuevoBacklogFalso(1); r.errGuardar = errBase; return r }()},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoCrearElementoBacklog(tt.repo).Ejecutar(context.Background(), datosElemento(1, "Login"))
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}
