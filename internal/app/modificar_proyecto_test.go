package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// repoModificacionFalso guarda proyectos en memoria y permite simular fallas.
type repoModificacionFalso struct {
	porID          map[int64]domain.Proyecto
	errBuscar      error
	errExiste      error
	errActualizar  error
	vecesActualizo int
}

func nuevoRepoModificacionFalso(proyectos ...domain.Proyecto) *repoModificacionFalso {
	r := &repoModificacionFalso{porID: map[int64]domain.Proyecto{}}
	for _, p := range proyectos {
		r.porID[p.ID] = p
	}
	return r
}

func (r *repoModificacionFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if r.errBuscar != nil {
		return domain.Proyecto{}, r.errBuscar
	}
	p, ok := r.porID[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (r *repoModificacionFalso) ExisteNombre(_ context.Context, clave string) (bool, error) {
	if r.errExiste != nil {
		return false, r.errExiste
	}
	for _, p := range r.porID {
		if domain.ClaveNombre(p.Nombre) == clave {
			return true, nil
		}
	}
	return false, nil
}

func (r *repoModificacionFalso) Actualizar(_ context.Context, p domain.Proyecto) error {
	if r.errActualizar != nil {
		return r.errActualizar
	}
	r.vecesActualizo++
	r.porID[p.ID] = p
	return nil
}

func proyectoGuardado(id int64, nombre string) domain.Proyecto {
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	return domain.Proyecto{
		ID:          id,
		Nombre:      nombre,
		FechaInicio: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    &fin,
		Estado:      domain.EstadoActivo,
	}
}

func ptrTexto(s string) *string                { return &s }
func ptrEstado(e domain.Estado) *domain.Estado { return &e }

// Traza: HU-02 CA-02.1, CA-02.10
func TestModificarProyecto_Valido_ActualizaYDevuelveElProyecto(t *testing.T) {
	repo := nuevoRepoModificacionFalso(proyectoGuardado(1, "Alfa"))
	uc := NuevoModificarProyecto(repo)

	p, err := uc.Ejecutar(context.Background(), DatosModificarProyecto{
		ProyectoID: 1,
		Cambios:    domain.CambiosProyecto{Nombre: ptrTexto("Alfa Renovado")},
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.ID != 1 || p.Nombre != "Alfa Renovado" || p.Estado != domain.EstadoActivo {
		t.Errorf("proyecto devuelto = %+v, se esperaba ID 1, Alfa Renovado, Activo", p)
	}
	if repo.vecesActualizo != 1 || repo.porID[1].Nombre != "Alfa Renovado" {
		t.Errorf("se actualizó %d veces, guardado = %+v", repo.vecesActualizo, repo.porID[1])
	}
}

// Traza: HU-02 CA-02.3, RN-02.3: renombrar el mismo proyecto con otras mayúsculas está permitido.
func TestModificarProyecto_MismoNombreConOtrasMayusculas_Actualiza(t *testing.T) {
	repo := nuevoRepoModificacionFalso(proyectoGuardado(1, "Alfa"), proyectoGuardado(2, "Beta"))

	p, err := NuevoModificarProyecto(repo).Ejecutar(context.Background(), DatosModificarProyecto{
		ProyectoID: 1,
		Cambios:    domain.CambiosProyecto{Nombre: ptrTexto("alfa")},
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.Nombre != "alfa" {
		t.Errorf("Nombre = %q, se esperaba %q", p.Nombre, "alfa")
	}
}

// Traza: HU-02 CA-02.3, CA-02.9, sección 7: si el pedido se rechaza, no se actualiza nada.
func TestModificarProyecto_Invalido_RetornaErrorYNoActualiza(t *testing.T) {
	antes := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		nombreCaso string
		datos      DatosModificarProyecto
		wantErr    error
	}{
		{"proyecto inexistente", DatosModificarProyecto{ProyectoID: 99, Cambios: domain.CambiosProyecto{Nombre: ptrTexto("Gamma")}}, domain.ErrProyectoNoEncontrado},
		{"nombre de otro proyecto", DatosModificarProyecto{ProyectoID: 2, Cambios: domain.CambiosProyecto{Nombre: ptrTexto("  ALFA ")}}, domain.ErrNombreDuplicado},
		{"sin cambios", DatosModificarProyecto{ProyectoID: 1}, domain.ErrSinCambios},
		{"estado inválido", DatosModificarProyecto{ProyectoID: 1, Cambios: domain.CambiosProyecto{Estado: ptrEstado("Cerrado")}}, domain.ErrEstadoInvalido},
		// Sección 7: el duplicado se informa después de los errores del dominio.
		{
			"nombre duplicado y fecha de fin anterior",
			DatosModificarProyecto{ProyectoID: 2, Cambios: domain.CambiosProyecto{Nombre: ptrTexto("Alfa"), CambiarFechaFin: true, FechaFin: &antes}},
			domain.ErrFechaFinAnterior,
		},
	}
	for _, tt := range tests {
		t.Run(tt.nombreCaso, func(t *testing.T) {
			repo := nuevoRepoModificacionFalso(proyectoGuardado(1, "Alfa"), proyectoGuardado(2, "Beta"))

			_, err := NuevoModificarProyecto(repo).Ejecutar(context.Background(), tt.datos)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if repo.vecesActualizo != 0 {
				t.Errorf("no se esperaba ninguna actualización, hubo %d", repo.vecesActualizo)
			}
		})
	}
}

// Traza: HU-02 RN-02.3: si no se cambia el nombre, no hace falta consultar si está repetido.
func TestModificarProyecto_SinNombre_NoConsultaDuplicados(t *testing.T) {
	repo := nuevoRepoModificacionFalso(proyectoGuardado(1, "Alfa"))
	repo.errExiste = errors.New("no debería consultarse")

	p, err := NuevoModificarProyecto(repo).Ejecutar(context.Background(), DatosModificarProyecto{
		ProyectoID: 1,
		Cambios:    domain.CambiosProyecto{Estado: ptrEstado(domain.EstadoFinalizado)},
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if p.Estado != domain.EstadoFinalizado {
		t.Errorf("Estado = %q, se esperaba %q", p.Estado, domain.EstadoFinalizado)
	}
}

func TestModificarProyecto_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		repo   *repoModificacionFalso
	}{
		{"al buscar el proyecto", &repoModificacionFalso{errBuscar: errBase}},
		{"al verificar el nombre", &repoModificacionFalso{porID: map[int64]domain.Proyecto{1: proyectoGuardado(1, "Alfa")}, errExiste: errBase}},
		{"al actualizar", &repoModificacionFalso{porID: map[int64]domain.Proyecto{1: proyectoGuardado(1, "Alfa")}, errActualizar: errBase}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := NuevoModificarProyecto(tt.repo).Ejecutar(context.Background(), DatosModificarProyecto{
				ProyectoID: 1,
				Cambios:    domain.CambiosProyecto{Nombre: ptrTexto("Gamma")},
			})
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}
