package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// defectoFalso implementa los cuatro puertos de RegistrarDefecto en memoria y permite simular
// fallas. Proyecto 1 tiene el elemento 10; proyecto 2 tiene el elemento 20.
type defectoFalso struct {
	proyectos   map[int64]domain.Proyecto
	elementos   map[int64]domain.ElementoBacklog
	sprints     []domain.Sprint
	guardados   []domain.Defecto
	siguienteID int64
	errProyecto error
	errElemento error
	errSprints  error
	errGuardar  error
}

func nuevoDefectoFalso() *defectoFalso {
	return &defectoFalso{
		proyectos: map[int64]domain.Proyecto{
			1: proyectoFalso(1),
			2: proyectoFalso(2),
		},
		elementos: map[int64]domain.ElementoBacklog{
			10: {ID: 10, ProyectoID: 1, Titulo: "Login"},
			20: {ID: 20, ProyectoID: 2, Titulo: "Registro"},
		},
	}
}

func (f *defectoFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *defectoFalso) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if f.errElemento != nil {
		return domain.ElementoBacklog{}, f.errElemento
	}
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *defectoFalso) ListarSprints(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
	if f.errSprints != nil {
		return nil, f.errSprints
	}
	var del []domain.Sprint
	for _, s := range f.sprints {
		if s.ProyectoID == proyectoID {
			del = append(del, s)
		}
	}
	return del, nil
}

// GuardarSprint no la usa RegistrarDefecto; está solo para cumplir RepositorioSprints.
func (f *defectoFalso) GuardarSprint(_ context.Context, s *domain.Sprint) error {
	f.sprints = append(f.sprints, *s)
	return nil
}

func (f *defectoFalso) GuardarDefecto(_ context.Context, d *domain.Defecto) error {
	if f.errGuardar != nil {
		return f.errGuardar
	}
	f.siguienteID++
	d.ID = f.siguienteID
	f.guardados = append(f.guardados, *d)
	return nil
}

// hoyFalsoDefecto es lo que devuelve el reloj falso: el caso de uso no lee la hora real.
var hoyFalsoDefecto = time.Date(2026, time.October, 25, 12, 0, 0, 0, time.Local)

func relojFalsoDefecto() time.Time { return hoyFalsoDefecto }

func registrarDefectoCon(f *defectoFalso) *RegistrarDefecto {
	return NuevoRegistrarDefecto(f, f, f, f, relojFalsoDefecto)
}

// datosDefecto es un defecto válido sobre el elemento 10 del proyecto 1.
func datosDefectoValidos() DatosDefecto {
	return DatosDefecto{
		ProyectoID:  1,
		ElementoID:  10,
		Descripcion: "No carga la pantalla de login en Firefox",
		Severidad:   domain.SeveridadAlta,
	}
}

func diaDefecto(anio int, mes time.Month, d int) time.Time {
	return time.Date(anio, mes, d, 0, 0, 0, 0, time.Local)
}

// Traza: HU-20 CA-20.1, CA-20.7
func TestRegistrarDefecto_DatosValidos_GuardaConID(t *testing.T) {
	f := nuevoDefectoFalso()
	uc := registrarDefectoCon(f)

	a, errA := uc.Ejecutar(context.Background(), datosDefectoValidos())
	b, errB := uc.Ejecutar(context.Background(), datosDefectoValidos())

	if errA != nil || errB != nil {
		t.Fatalf("errores inesperados: %v, %v", errA, errB)
	}
	if a.ID == 0 || a.ID == b.ID {
		t.Errorf("IDs = %d y %d, se esperaban asignados y distintos", a.ID, b.ID)
	}
	want := domain.Defecto{
		ID: a.ID, ProyectoID: 1, ElementoID: 10,
		Descripcion: "No carga la pantalla de login en Firefox", Severidad: domain.SeveridadAlta,
		Estado: domain.DefectoAbierto, SprintDeteccionID: nil,
	}
	if a != want {
		t.Errorf("se obtuvo %+v, se esperaba %+v", a, want)
	}
	if len(f.guardados) != 2 {
		t.Errorf("se esperaban 2 defectos guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-20 CA-20.1 (RN-20.5: se asigna el Sprint Abierto vigente hoy)
func TestRegistrarDefecto_ConSprintVigente_AsignaSuID(t *testing.T) {
	f := nuevoDefectoFalso()
	f.sprints = []domain.Sprint{
		{ID: 5, ProyectoID: 1, Estado: domain.SprintAbierto,
			FechaInicio: diaDefecto(2026, time.October, 20), FechaFin: diaDefecto(2026, time.October, 30)},
	}

	d, err := registrarDefectoCon(f).Ejecutar(context.Background(), datosDefectoValidos())

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.SprintDeteccionID == nil || *d.SprintDeteccionID != 5 {
		t.Errorf("SprintDeteccionID = %v, se esperaba un puntero a 5", d.SprintDeteccionID)
	}
}

// Traza: HU-20 CA-20.2 (RN-20.5: sin sprint vigente, queda sin asignar)
func TestRegistrarDefecto_SinSprintVigente_QuedaSinAsignar(t *testing.T) {
	f := nuevoDefectoFalso()
	f.sprints = []domain.Sprint{ // ya terminó: no incluye "hoy" (2026-10-25)
		{ID: 5, ProyectoID: 1, Estado: domain.SprintAbierto,
			FechaInicio: diaDefecto(2026, time.September, 1), FechaFin: diaDefecto(2026, time.September, 15)},
	}

	d, err := registrarDefectoCon(f).Ejecutar(context.Background(), datosDefectoValidos())

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.SprintDeteccionID != nil {
		t.Errorf("SprintDeteccionID = %v, se esperaba nil", *d.SprintDeteccionID)
	}
}

// Traza: HU-20 CA-20.5
func TestRegistrarDefecto_ElementoInexistente_RetornaErrorYNoGuarda(t *testing.T) {
	f := nuevoDefectoFalso()
	d := datosDefectoValidos()
	d.ElementoID = 99

	_, err := registrarDefectoCon(f).Ejecutar(context.Background(), d)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if len(f.guardados) != 0 {
		t.Errorf("no se esperaban defectos guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-20 CA-20.5
func TestRegistrarDefecto_ProyectoInexistente_RetornaErrorYNoGuarda(t *testing.T) {
	f := nuevoDefectoFalso()
	d := datosDefectoValidos()
	d.ProyectoID = 99
	d.Descripcion = "" // aunque los datos sean inválidos, primero se informa el proyecto (sección 7)

	_, err := registrarDefectoCon(f).Ejecutar(context.Background(), d)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if len(f.guardados) != 0 {
		t.Errorf("no se esperaban defectos guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-20 CA-20.5 (RN-20.1: un elemento de otro proyecto se trata como inexistente)
func TestRegistrarDefecto_ElementoDeOtroProyecto_RetornaErrorYNoGuarda(t *testing.T) {
	f := nuevoDefectoFalso()
	d := datosDefectoValidos()
	d.ElementoID = 20 // es del proyecto 2

	_, err := registrarDefectoCon(f).Ejecutar(context.Background(), d)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if len(f.guardados) != 0 {
		t.Errorf("no se esperaban defectos guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-20 CA-20.3, CA-20.4
func TestRegistrarDefecto_DatosInvalidos_NoLleganAlRepositorio(t *testing.T) {
	tests := []struct {
		nombre  string
		cambiar func(*DatosDefecto)
		wantErr error
	}{
		{"descripción vacía", func(d *DatosDefecto) { d.Descripcion = "  " }, domain.ErrDescripcionDefectoVacia},
		{"severidad inválida", func(d *DatosDefecto) { d.Severidad = "urgente" }, domain.ErrSeveridadInvalida},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoDefectoFalso()
			d := datosDefectoValidos()
			tt.cambiar(&d)

			_, err := registrarDefectoCon(f).Ejecutar(context.Background(), d)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if len(f.guardados) != 0 {
				t.Errorf("no se esperaban defectos guardados, hay %d", len(f.guardados))
			}
		})
	}
}

func TestRegistrarDefecto_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		fallar func(*defectoFalso)
	}{
		{"al buscar el proyecto", func(f *defectoFalso) { f.errProyecto = errBase }},
		{"al buscar el elemento", func(f *defectoFalso) { f.errElemento = errBase }},
		{"al listar los sprints", func(f *defectoFalso) { f.errSprints = errBase }},
		{"al guardar", func(f *defectoFalso) { f.errGuardar = errBase }},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoDefectoFalso()
			tt.fallar(f)

			_, err := registrarDefectoCon(f).Ejecutar(context.Background(), datosDefectoValidos())

			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
			if len(f.guardados) != 0 {
				t.Errorf("no se esperaban defectos guardados, hay %d", len(f.guardados))
			}
		})
	}
}
