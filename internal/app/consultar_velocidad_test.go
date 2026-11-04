package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: HU-22 CA-22.2: la velocidad es el promedio de los sprints cerrados del proyecto.
func TestConsultarVelocidad_PromedioDeLosCerrados(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{
		sprintCerradoConPuntos(1, 1, 1, 8),
		sprintCerradoConPuntos(2, 1, 2, 6),
		sprintCerradoConPuntos(3, 1, 3, 4),
		sprintGuardado(4, 1, 4, 26, domain.SprintAbierto), // no cuenta
	}

	got, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := VelocidadProyecto{Velocidad: 6, HayDatos: true, SprintsCerrados: 3}
	if got != want {
		t.Errorf("velocidad = %+v, se esperaba %+v", got, want)
	}
}

// Traza: HU-22 CA-22.4: el promedio se informa con decimales.
func TestConsultarVelocidad_ConDecimales(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{
		sprintCerradoConPuntos(1, 1, 1, 5),
		sprintCerradoConPuntos(2, 1, 2, 4),
	}

	got, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := VelocidadProyecto{Velocidad: 4.5, HayDatos: true, SprintsCerrados: 2}
	if got != want {
		t.Errorf("velocidad = %+v, se esperaba %+v", got, want)
	}
}

// Traza: HU-22 CA-22.3: sin sprints cerrados, no hay velocidad (ni sin sprints, ni con todos abiertos).
func TestConsultarVelocidad_SinCerrados_NoHayDatos(t *testing.T) {
	tests := []struct {
		name    string
		sprints []domain.Sprint
	}{
		{"sin ningún sprint", nil},
		{"con sprints, pero ninguno cerrado", []domain.Sprint{sprintGuardado(1, 1, 1, 1, domain.SprintAbierto)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoSprintsFalso(proyectoActivo(1))
			f.sprints = tt.sprints

			got, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 1)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got.HayDatos {
				t.Errorf("velocidad = %+v, se esperaba HayDatos = false", got)
			}
		})
	}
}

// Traza: HU-22 CA-22.5: un único sprint cerrado con 0 puntos da velocidad 0, no "sin datos".
func TestConsultarVelocidad_UnSoloSprintEnCero_DaVelocidadCero(t *testing.T) {
	f := nuevoSprintsFalso(proyectoActivo(1))
	f.sprints = []domain.Sprint{sprintCerradoConPuntos(1, 1, 1, 0)}

	got, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	want := VelocidadProyecto{Velocidad: 0, HayDatos: true, SprintsCerrados: 1}
	if got != want {
		t.Errorf("velocidad = %+v, se esperaba %+v", got, want)
	}
}

// Traza: HU-22 CA-22.6: con el proyecto inexistente no se leen sus sprints.
func TestConsultarVelocidad_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoSprintsFalso()

	_, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if f.listoProyect != 0 {
		t.Errorf("se listaron los sprints del proyecto %d; no se esperaba listar nada", f.listoProyect)
	}
}

// Traza: HU-22 sección 7: una falla al leer se propaga, sin perder la causa.
func TestConsultarVelocidad_FallaDelRepositorio_SePropaga(t *testing.T) {
	disco := errors.New("disco lleno")
	tests := []struct {
		name  string
		falla func(*sprintsFalso)
	}{
		{"al buscar el proyecto", func(f *sprintsFalso) { f.errBuscar = disco }},
		{"al listar los sprints", func(f *sprintsFalso) { f.errListar = disco }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := nuevoSprintsFalso(proyectoActivo(1))
			tt.falla(f)

			_, err := NuevoConsultarVelocidad(f, f).Ejecutar(context.Background(), 1)

			if !errors.Is(err, disco) {
				t.Errorf("error = %v, se esperaba que envolviera %v", err, disco)
			}
		})
	}
}

// sprintCerradoConPuntos arma un sprint Cerrado del proyecto, con sus puntos completados (los
// planificados no le importan a la velocidad).
func sprintCerradoConPuntos(id, proyectoID int64, numero int, completados int) domain.Sprint {
	s := sprintGuardado(id, proyectoID, numero, numero, domain.SprintCerrado)
	s.PuntosCompletados = completados
	return s
}
