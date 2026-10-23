package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// repoIntegrantesFalso guarda proyectos en memoria y permite simular fallas.
type repoIntegrantesFalso struct {
	porID       map[int64]domain.Proyecto
	errBuscar   error
	errGuardar  error
	vecesGuardo int
}

func nuevoRepoIntegrantesFalso(proyectos ...domain.Proyecto) *repoIntegrantesFalso {
	r := &repoIntegrantesFalso{porID: map[int64]domain.Proyecto{}}
	for _, p := range proyectos {
		r.porID[p.ID] = p
	}
	return r
}

func (r *repoIntegrantesFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if r.errBuscar != nil {
		return domain.Proyecto{}, r.errBuscar
	}
	p, ok := r.porID[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (r *repoIntegrantesFalso) GuardarIntegrantes(_ context.Context, p *domain.Proyecto) error {
	if r.errGuardar != nil {
		return r.errGuardar
	}
	r.vecesGuardo++
	r.porID[p.ID] = *p
	return nil
}

func proyectoFalso(id int64, integrantes ...domain.Integrante) domain.Proyecto {
	return domain.Proyecto{ID: id, Nombre: "TP Integrador", Estado: domain.EstadoActivo, Integrantes: integrantes}
}

// Traza: HU-03 CA-03.1
func TestAgregarIntegrante_DatosValidos_GuardaEnElProyecto(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso(proyectoFalso(1))
	uc := NuevoAgregarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosIntegrante{ProyectoID: 1, Nombre: "Ana", Rol: "Desarrolladora"})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	guardado := repo.porID[1]
	if len(guardado.Integrantes) != 1 || guardado.Integrantes[0].Nombre != "Ana" {
		t.Errorf("Integrantes = %+v, se esperaba 1 integrante Ana", guardado.Integrantes)
	}
	if repo.vecesGuardo != 1 {
		t.Errorf("se guardó %d veces, se esperaba 1", repo.vecesGuardo)
	}
}

// Traza: HU-03 (error condición), sección 7: proyecto inexistente
func TestAgregarIntegrante_ProyectoInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso()
	uc := NuevoAgregarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosIntegrante{ProyectoID: 99, Nombre: "Ana", Rol: "Desarrolladora"})

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-03 CA-03.2, CA-03.5, CA-03.6: si el dominio rechaza los datos, no se guarda nada.
func TestAgregarIntegrante_DatosInvalidos_NoGuarda(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso(proyectoFalso(1))
	uc := NuevoAgregarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosIntegrante{ProyectoID: 1, Nombre: "", Rol: "Desarrolladora"})

	if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

func TestAgregarIntegrante_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		repo   *repoIntegrantesFalso
	}{
		{"al buscar el proyecto", &repoIntegrantesFalso{errBuscar: errBase}},
		{"al guardar", &repoIntegrantesFalso{porID: map[int64]domain.Proyecto{1: proyectoFalso(1)}, errGuardar: errBase}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			err := NuevoAgregarIntegrante(tt.repo).Ejecutar(context.Background(),
				DatosIntegrante{ProyectoID: 1, Nombre: "Ana", Rol: "Desarrolladora"})
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}

// Traza: HU-03 CA-03.3
func TestQuitarIntegrante_Existente_GuardaSinEl(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso(proyectoFalso(1, domain.Integrante{Nombre: "Ana", Rol: "Desarrolladora"}))
	uc := NuevoQuitarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosQuitarIntegrante{ProyectoID: 1, Nombre: "Ana"})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(repo.porID[1].Integrantes) != 0 {
		t.Errorf("Integrantes = %+v, se esperaba lista vacía", repo.porID[1].Integrantes)
	}
	if repo.vecesGuardo != 1 {
		t.Errorf("se guardó %d veces, se esperaba 1", repo.vecesGuardo)
	}
}

// Traza: HU-03 CA-03.4
func TestQuitarIntegrante_Inexistente_RetornaErrorYNoGuarda(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso(proyectoFalso(1, domain.Integrante{Nombre: "Ana", Rol: "Desarrolladora"}))
	uc := NuevoQuitarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosQuitarIntegrante{ProyectoID: 1, Nombre: "Caro"})

	if !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrIntegranteNoEncontrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-03 (error condición), sección 7: proyecto inexistente
func TestQuitarIntegrante_ProyectoInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoIntegrantesFalso()
	uc := NuevoQuitarIntegrante(repo)

	err := uc.Ejecutar(context.Background(), DatosQuitarIntegrante{ProyectoID: 99, Nombre: "Ana"})

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
}
