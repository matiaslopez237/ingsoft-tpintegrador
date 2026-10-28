package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// sprintFalso es lo que el repositorio falso guarda de un sprint para los tests.
type sprintFalso struct {
	proyectoID int64
	abierto    bool
}

// repoAsignacionFalso guarda elementos y sprints en memoria y permite simular fallas.
type repoAsignacionFalso struct {
	elementos      map[int64]domain.ElementoBacklog
	sprints        map[int64]sprintFalso
	errBuscarElem  error
	errBuscarSprnt error
	errGuardar     error
	vecesGuardo    int
}

func nuevoRepoAsignacionFalso() *repoAsignacionFalso {
	return &repoAsignacionFalso{elementos: map[int64]domain.ElementoBacklog{}, sprints: map[int64]sprintFalso{}}
}

func (r *repoAsignacionFalso) BuscarElementoPorID(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if r.errBuscarElem != nil {
		return domain.ElementoBacklog{}, r.errBuscarElem
	}
	e, ok := r.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (r *repoAsignacionFalso) BuscarSprint(_ context.Context, id int64) (proyectoID int64, abierto bool, err error) {
	if r.errBuscarSprnt != nil {
		return 0, false, r.errBuscarSprnt
	}
	s, ok := r.sprints[id]
	if !ok {
		return 0, false, domain.ErrSprintNoEncontrado
	}
	return s.proyectoID, s.abierto, nil
}

func (r *repoAsignacionFalso) GuardarAsignacion(_ context.Context, e *domain.ElementoBacklog) error {
	if r.errGuardar != nil {
		return r.errGuardar
	}
	r.vecesGuardo++
	r.elementos[e.ID] = *e
	return nil
}

func elementoFalso(id, proyectoID int64, sprintID *int64) domain.ElementoBacklog {
	return domain.ElementoBacklog{
		ID: id, ProyectoID: proyectoID, Titulo: "Login de usuarios", Estado: domain.EstadoPendiente, SprintID: sprintID,
	}
}

// --- AsignarHistoria ---

// Traza: HU-09 CA-09.1
func TestAsignarHistoria_DatosValidos_GuardaElSprintID(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	guardado := repo.elementos[1]
	if guardado.SprintID == nil || *guardado.SprintID != 10 {
		t.Errorf("SprintID = %v, se esperaba 10", guardado.SprintID)
	}
	if repo.vecesGuardo != 1 {
		t.Errorf("se guardó %d veces, se esperaba 1", repo.vecesGuardo)
	}
}

// Traza: HU-09 CA-09.3: el sprint previo (20) ya está cerrado, se puede reasignar.
func TestAsignarHistoria_SprintPrevioCerrado_Reasigna(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	sprintPrevio := int64(20)
	repo.elementos[1] = elementoFalso(1, 100, &sprintPrevio)
	repo.sprints[20] = sprintFalso{proyectoID: 100, abierto: false}
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if guardado := repo.elementos[1]; guardado.SprintID == nil || *guardado.SprintID != 10 {
		t.Errorf("SprintID = %v, se esperaba 10", guardado.SprintID)
	}
}

// Traza: HU-09 CA-09.2: el sprint previo (20) sigue activo, no se puede reasignar.
func TestAsignarHistoria_YaAsignadoAOtroSprintActivo_NoGuarda(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	sprintPrevio := int64(20)
	repo.elementos[1] = elementoFalso(1, 100, &sprintPrevio)
	repo.sprints[20] = sprintFalso{proyectoID: 100, abierto: true}
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if !errors.Is(err, domain.ErrElementoYaAsignado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoYaAsignado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-09 CA-09.7
func TestAsignarHistoria_SprintDeOtroProyecto_NoGuarda(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	repo.sprints[10] = sprintFalso{proyectoID: 999, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if !errors.Is(err, domain.ErrSprintDeOtroProyecto) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintDeOtroProyecto)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-09 CA-09.8
func TestAsignarHistoria_SprintDestinoCerrado_NoGuarda(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: false}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-09, sección 7: elemento inexistente
func TestAsignarHistoria_ElementoInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 99, SprintID: 10})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-09, sección 7: sprint inexistente
func TestAsignarHistoria_SprintInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 99})

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

func TestAsignarHistoria_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	repoBuscarElem := nuevoRepoAsignacionFalso()
	repoBuscarElem.errBuscarElem = errBase

	repoBuscarSprnt := nuevoRepoAsignacionFalso()
	repoBuscarSprnt.elementos[1] = elementoFalso(1, 100, nil)
	repoBuscarSprnt.errBuscarSprnt = errBase

	repoGuardar := nuevoRepoAsignacionFalso()
	repoGuardar.elementos[1] = elementoFalso(1, 100, nil)
	repoGuardar.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	repoGuardar.errGuardar = errBase

	tests := []struct {
		nombre string
		repo   *repoAsignacionFalso
	}{
		{"al buscar el elemento", repoBuscarElem},
		{"al buscar el sprint", repoBuscarSprnt},
		{"al guardar", repoGuardar},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			err := NuevoAsignarHistoria(tt.repo).Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})
			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}

// --- DesasignarHistoria ---

// Traza: HU-09 CA-09.4
func TestDesasignarHistoria_SprintActivo_Desasigna(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	sprintID := int64(10)
	repo.elementos[1] = elementoFalso(1, 100, &sprintID)
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ProyectoID: 100, ElementoID: 1})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if guardado := repo.elementos[1]; guardado.SprintID != nil {
		t.Errorf("SprintID = %v, se esperaba nil", guardado.SprintID)
	}
}

// Traza: HU-09 CA-09.5
func TestDesasignarHistoria_SprintCerrado_RetornaErrorYNoGuarda(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	sprintID := int64(10)
	repo.elementos[1] = elementoFalso(1, 100, &sprintID)
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: false}
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ProyectoID: 100, ElementoID: 1})

	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

// Traza: HU-09 CA-09.6
func TestDesasignarHistoria_SinSprintAsignado_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ProyectoID: 100, ElementoID: 1})

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
}

// Traza: HU-09, sección 7: elemento inexistente
func TestDesasignarHistoria_ElementoInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ProyectoID: 100, ElementoID: 99})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-09, sección 7, RN-09.7 (igual que RN-18.1): un elemento de otro proyecto se trata
// como inexistente para la ruta anidada /proyectos/{id}/backlog/{elementoID}/sprint.
func TestAsignarHistoria_ElementoDeOtroProyecto_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 200, nil) // el elemento es del proyecto 200
	repo.sprints[10] = sprintFalso{proyectoID: 100, abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ProyectoID: 100, ElementoID: 1, SprintID: 10})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}

func TestDesasignarHistoria_ElementoDeOtroProyecto_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	sprintID := int64(10)
	repo.elementos[1] = elementoFalso(1, 200, &sprintID) // el elemento es del proyecto 200
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ProyectoID: 100, ElementoID: 1})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
	if repo.vecesGuardo != 0 {
		t.Errorf("no se esperaba ningún guardado, hubo %d", repo.vecesGuardo)
	}
}
