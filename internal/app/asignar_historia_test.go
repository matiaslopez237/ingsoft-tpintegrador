package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// repoAsignacionFalso guarda elementos y sprints en memoria y permite simular fallas.
type repoAsignacionFalso struct {
	elementos      map[int64]domain.ElementoBacklog
	sprints        map[int64]InfoSprint
	errBuscarElem  error
	errBuscarSprnt error
	errGuardar     error
	vecesGuardo    int
}

func nuevoRepoAsignacionFalso() *repoAsignacionFalso {
	return &repoAsignacionFalso{elementos: map[int64]domain.ElementoBacklog{}, sprints: map[int64]InfoSprint{}}
}

func (r *repoAsignacionFalso) BuscarElemento(_ context.Context, id int64) (domain.ElementoBacklog, error) {
	if r.errBuscarElem != nil {
		return domain.ElementoBacklog{}, r.errBuscarElem
	}
	e, ok := r.elementos[id]
	if !ok {
		return domain.ElementoBacklog{}, domain.ErrElementoNoEncontrado
	}
	return e, nil
}

func (r *repoAsignacionFalso) BuscarSprint(_ context.Context, id int64) (InfoSprint, error) {
	if r.errBuscarSprnt != nil {
		return InfoSprint{}, r.errBuscarSprnt
	}
	s, ok := r.sprints[id]
	if !ok {
		return InfoSprint{}, ErrSprintNoEncontrado
	}
	return s, nil
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
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})

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
	repo.sprints[20] = InfoSprint{ProyectoID: 100, Abierto: false}
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})

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
	repo.sprints[20] = InfoSprint{ProyectoID: 100, Abierto: true}
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})

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
	repo.sprints[10] = InfoSprint{ProyectoID: 999, Abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})

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
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: false}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})

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
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 99, SprintID: 10})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}

// Traza: HU-09, sección 7: sprint inexistente
func TestAsignarHistoria_SprintInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	repo.elementos[1] = elementoFalso(1, 100, nil)
	uc := NuevoAsignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 99})

	if !errors.Is(err, ErrSprintNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSprintNoEncontrado)
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
	repoGuardar.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
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
			err := NuevoAsignarHistoria(tt.repo).Ejecutar(context.Background(), DatosAsignacion{ElementoID: 1, SprintID: 10})
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
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: true}
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ElementoID: 1})

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
	repo.sprints[10] = InfoSprint{ProyectoID: 100, Abierto: false}
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ElementoID: 1})

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

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ElementoID: 1})

	if !errors.Is(err, domain.ErrElementoSinSprint) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoSinSprint)
	}
}

// Traza: HU-09, sección 7: elemento inexistente
func TestDesasignarHistoria_ElementoInexistente_RetornaError(t *testing.T) {
	repo := nuevoRepoAsignacionFalso()
	uc := NuevoDesasignarHistoria(repo)

	err := uc.Ejecutar(context.Background(), DatosDesasignacion{ElementoID: 99})

	if !errors.Is(err, domain.ErrElementoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrElementoNoEncontrado)
	}
}
