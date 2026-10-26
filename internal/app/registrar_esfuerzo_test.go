package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// esfuerzoFalso implementa los tres puertos de RegistrarEsfuerzo en memoria y permite simular fallas.
// Proyecto 1: "Ana" y el elemento 10. Proyecto 2: "Beto" y el elemento 20.
type esfuerzoFalso struct {
	proyectos   map[int64]domain.Proyecto
	elementos   map[int64]domain.ElementoBacklog
	guardados   []domain.RegistroEsfuerzo
	siguienteID int64
	errProyecto error
	errElemento error
	errGuardar  error
}

func nuevoEsfuerzoFalso() *esfuerzoFalso {
	return &esfuerzoFalso{
		proyectos: map[int64]domain.Proyecto{
			1: proyectoFalso(1, domain.Integrante{Nombre: "Ana", Roles: []string{"Desarrolladora"}}),
			2: proyectoFalso(2, domain.Integrante{Nombre: "Beto", Roles: []string{"Tester"}}),
		},
		elementos: map[int64]domain.ElementoBacklog{
			10: {ID: 10, ProyectoID: 1, Titulo: "Login"},
			20: {ID: 20, ProyectoID: 2, Titulo: "Registro"},
		},
	}
}

func (f *esfuerzoFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *esfuerzoFalso) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if f.errElemento != nil {
		return domain.ElementoBacklog{}, f.errElemento
	}
	e, ok := f.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (f *esfuerzoFalso) GuardarRegistro(_ context.Context, r *domain.RegistroEsfuerzo) error {
	if f.errGuardar != nil {
		return f.errGuardar
	}
	f.siguienteID++
	r.ID = f.siguienteID
	f.guardados = append(f.guardados, *r)
	return nil
}

// hoyFalso es lo que devuelve el reloj falso: el caso de uso no lee la hora real.
var hoyFalso = time.Date(2026, time.October, 25, 12, 0, 0, 0, time.Local)

func relojFalso() time.Time { return hoyFalso }

func registrarEsfuerzoCon(f *esfuerzoFalso) *RegistrarEsfuerzo {
	return NuevoRegistrarEsfuerzo(f, f, f, relojFalso)
}

// datosEsfuerzo es un registro válido de "Ana" en el elemento 10 del proyecto 1.
func datosEsfuerzo() DatosRegistroEsfuerzo {
	return DatosRegistroEsfuerzo{
		ProyectoID: 1,
		ElementoID: 10,
		Integrante: "Ana",
		Fecha:      time.Date(2026, time.October, 25, 0, 0, 0, 0, time.Local),
		Actividad:  "Desarrollo",
		Horas:      2,
	}
}

// Traza: HU-18 CA-18.1, CA-18.12
func TestRegistrarEsfuerzo_DatosValidos_GuardaConID(t *testing.T) {
	f := nuevoEsfuerzoFalso()
	uc := registrarEsfuerzoCon(f)

	a, errA := uc.Ejecutar(context.Background(), datosEsfuerzo())
	b, errB := uc.Ejecutar(context.Background(), datosEsfuerzo())

	if errA != nil || errB != nil {
		t.Fatalf("errores inesperados: %v, %v", errA, errB)
	}
	if a.ID == 0 || a.ID == b.ID {
		t.Errorf("IDs = %d y %d, se esperaban asignados y distintos", a.ID, b.ID)
	}
	want := domain.RegistroEsfuerzo{
		ID: a.ID, ProyectoID: 1, ElementoID: 10, Integrante: "Ana",
		Fecha: datosEsfuerzo().Fecha, Actividad: "Desarrollo", Horas: 2,
	}
	if a != want {
		t.Errorf("se obtuvo %+v, se esperaba %+v", a, want)
	}
	if len(f.guardados) != 2 {
		t.Errorf("se esperaban 2 registros guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-18 CA-18.7
func TestRegistrarEsfuerzo_ElementoInexistente_RetornaErrorYNoGuarda(t *testing.T) {
	f := nuevoEsfuerzoFalso()
	d := datosEsfuerzo()
	d.ElementoID = 99

	_, err := registrarEsfuerzoCon(f).Ejecutar(context.Background(), d)

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if len(f.guardados) != 0 {
		t.Errorf("no se esperaban registros guardados, hay %d", len(f.guardados))
	}
}

// Traza: HU-18 CA-18.4, CA-18.6 (el "hoy" sale del reloj que recibe el caso de uso)
func TestRegistrarEsfuerzo_DatosInvalidos_NoLleganAlRepositorio(t *testing.T) {
	tests := []struct {
		nombre  string
		cambiar func(*DatosRegistroEsfuerzo)
		wantErr error
	}{
		{"actividad vacía", func(d *DatosRegistroEsfuerzo) { d.Actividad = "  " }, domain.ErrActividadVacia},
		{"mañana según el reloj", func(d *DatosRegistroEsfuerzo) { d.Fecha = hoyFalso.AddDate(0, 0, 1) }, domain.ErrFechaEsfuerzoFutura},
		{"horas negativas", func(d *DatosRegistroEsfuerzo) { d.Horas = -1 }, domain.ErrHorasInvalidas},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoEsfuerzoFalso()
			d := datosEsfuerzo()
			tt.cambiar(&d)

			_, err := registrarEsfuerzoCon(f).Ejecutar(context.Background(), d)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if len(f.guardados) != 0 {
				t.Errorf("no se esperaban registros guardados, hay %d", len(f.guardados))
			}
		})
	}
}

func TestRegistrarEsfuerzo_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		fallar func(*esfuerzoFalso)
	}{
		{"al buscar el elemento", func(f *esfuerzoFalso) { f.errElemento = errBase }},
		{"al guardar", func(f *esfuerzoFalso) { f.errGuardar = errBase }},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoEsfuerzoFalso()
			tt.fallar(f)

			_, err := registrarEsfuerzoCon(f).Ejecutar(context.Background(), datosEsfuerzo())

			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
			if len(f.guardados) != 0 {
				t.Errorf("no se esperaban registros guardados, hay %d", len(f.guardados))
			}
		})
	}
}
