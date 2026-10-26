package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// listadoBacklogFalso guarda proyectos y elementos en memoria y permite simular fallas.
// Cumple LectorProyectos y ListadorElementos.
type listadoBacklogFalso struct {
	proyectos map[int64]domain.Proyecto
	elementos []domain.ElementoBacklog
	errBuscar error
	errListar error
}

func (f *listadoBacklogFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errBuscar != nil {
		return domain.Proyecto{}, f.errBuscar
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *listadoBacklogFalso) ListarElementos(_ context.Context, proyectoID int64) ([]domain.ElementoBacklog, error) {
	if f.errListar != nil {
		return nil, f.errListar
	}
	var del []domain.ElementoBacklog
	for _, e := range f.elementos {
		if e.ProyectoID == proyectoID {
			del = append(del, e)
		}
	}
	return del, nil
}

func nuevoListadoBacklogFalso() *listadoBacklogFalso {
	return &listadoBacklogFalso{
		proyectos: map[int64]domain.Proyecto{1: {ID: 1, Nombre: "Alfa", FechaInicio: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Estado: domain.EstadoActivo}},
		elementos: []domain.ElementoBacklog{
			{ID: 1, ProyectoID: 1, Titulo: "Reportes", Prioridad: domain.PrioridadBaja, Estado: domain.EstadoHecho},
			{ID: 2, ProyectoID: 2, Titulo: "De otro proyecto", Prioridad: domain.PrioridadAlta, Estado: domain.EstadoPendiente},
			{ID: 3, ProyectoID: 1, Titulo: "Login", Prioridad: domain.PrioridadAlta, Estado: domain.EstadoPendiente},
		},
	}
}

func titulosDe(elementos []domain.ElementoBacklog) []string {
	var t []string
	for _, e := range elementos {
		t = append(t, e.Titulo)
	}
	return t
}

// Traza: HU-07 CA-07.1, CA-07.2, CA-07.4, CA-07.6: lista solo los del proyecto, en el orden del dominio.
func TestListarBacklog_DevuelveLosDelProyectoOrdenados(t *testing.T) {
	tests := []struct {
		nombre string
		datos  DatosListarBacklog
		want   string
	}{
		{"por defecto", DatosListarBacklog{ProyectoID: 1}, "[Login Reportes]"},
		{"ascendente", DatosListarBacklog{ProyectoID: 1, Orden: "asc"}, "[Reportes Login]"},
		{"filtro por estado", DatosListarBacklog{ProyectoID: 1, Estado: "Hecho"}, "[Reportes]"},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoListadoBacklogFalso()

			got, err := NuevoListarBacklog(f, f).Ejecutar(context.Background(), tt.datos)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if s := fmt.Sprint(titulosDe(got)); s != tt.want {
				t.Errorf("elementos = %s, se esperaba %s", s, tt.want)
			}
		})
	}
}

// Traza: HU-07 CA-07.7, CA-07.8, sección 7: el proyecto inexistente se informa antes que un orden inválido.
func TestListarBacklog_Errores(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre  string
		datos   DatosListarBacklog
		ajuste  func(*listadoBacklogFalso)
		wantErr error
	}{
		{"proyecto inexistente", DatosListarBacklog{ProyectoID: 99}, nil, domain.ErrProyectoNoEncontrado},
		{"proyecto inexistente y orden inválido", DatosListarBacklog{ProyectoID: 99, Orden: "ASC"}, nil, domain.ErrProyectoNoEncontrado},
		{"orden inválido", DatosListarBacklog{ProyectoID: 1, Orden: "ASC"}, nil, domain.ErrOrdenInvalido},
		{"estado inválido", DatosListarBacklog{ProyectoID: 1, Estado: "pendiente"}, nil, domain.ErrEstadoElementoInvalido},
		{"falla al buscar el proyecto", DatosListarBacklog{ProyectoID: 1}, func(f *listadoBacklogFalso) { f.errBuscar = errBase }, errBase},
		{"falla al listar los elementos", DatosListarBacklog{ProyectoID: 1}, func(f *listadoBacklogFalso) { f.errListar = errBase }, errBase},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoListadoBacklogFalso()
			if tt.ajuste != nil {
				tt.ajuste(f)
			}

			_, err := NuevoListarBacklog(f, f).Ejecutar(context.Background(), tt.datos)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}
