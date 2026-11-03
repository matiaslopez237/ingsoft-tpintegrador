package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// estadoDefectoFalso implementa los tres puertos de ActualizarEstadoDefecto en memoria y permite
// simular fallas. El proyecto 1 tiene el defecto 100 (Abierto); el proyecto 2 tiene el defecto 200.
type estadoDefectoFalso struct {
	proyectos    map[int64]domain.Proyecto
	defectos     map[int64]domain.Defecto
	sprints      []domain.Sprint
	errProyecto  error
	errBuscar    error
	errSprints   error
	errActualiza error
	vecesGuarda  int
}

func nuevoEstadoDefectoFalso() *estadoDefectoFalso {
	return &estadoDefectoFalso{
		proyectos: map[int64]domain.Proyecto{
			1: proyectoFalso(1),
			2: proyectoFalso(2),
		},
		defectos: map[int64]domain.Defecto{
			100: {ID: 100, ProyectoID: 1, ElementoID: 10, Descripcion: "No carga el login", Severidad: domain.SeveridadAlta, Estado: domain.DefectoAbierto},
			200: {ID: 200, ProyectoID: 2, ElementoID: 20, Descripcion: "Falla el registro", Severidad: domain.SeveridadBaja, Estado: domain.DefectoAbierto},
		},
	}
}

func (f *estadoDefectoFalso) BuscarPorID(_ context.Context, id int64) (domain.Proyecto, error) {
	if f.errProyecto != nil {
		return domain.Proyecto{}, f.errProyecto
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *estadoDefectoFalso) ListarSprints(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
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

// GuardarSprint no la usa ActualizarEstadoDefecto; está solo para cumplir RepositorioSprints.
func (f *estadoDefectoFalso) GuardarSprint(_ context.Context, s *domain.Sprint) error {
	f.sprints = append(f.sprints, *s)
	return nil
}

func (f *estadoDefectoFalso) BuscarDefectoPorID(_ context.Context, id int64) (domain.Defecto, error) {
	if f.errBuscar != nil {
		return domain.Defecto{}, f.errBuscar
	}
	d, ok := f.defectos[id]
	if !ok {
		return domain.Defecto{}, domain.ErrDefectoNoEncontrado
	}
	return d, nil
}

func (f *estadoDefectoFalso) ActualizarEstadoDefecto(_ context.Context, d *domain.Defecto, _ domain.EstadoDefecto) error {
	f.vecesGuarda++
	if f.errActualiza != nil {
		return f.errActualiza
	}
	f.defectos[d.ID] = *d
	return nil
}

// hoyFalsoEstadoDefecto es lo que devuelve el reloj falso: el caso de uso no lee la hora real.
var hoyFalsoEstadoDefecto = time.Date(2026, time.October, 25, 12, 0, 0, 0, time.Local)

func relojFalsoEstadoDefecto() time.Time { return hoyFalsoEstadoDefecto }

func actualizarEstadoDefectoCon(f *estadoDefectoFalso) *ActualizarEstadoDefecto {
	return NuevoActualizarEstadoDefecto(f, f, f, relojFalsoEstadoDefecto)
}

func diaEstadoDefecto(anio int, mes time.Month, d int) time.Time {
	return time.Date(anio, mes, d, 0, 0, 0, 0, time.Local)
}

// --- Iniciar ---

// Traza: HU-21 CA-21.1
func TestActualizarEstadoDefecto_Iniciar_DesdeAbierto_QuedaEnProgreso(t *testing.T) {
	f := nuevoEstadoDefectoFalso()

	d, err := actualizarEstadoDefectoCon(f).Iniciar(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if d.Estado != domain.DefectoEnProgreso {
		t.Errorf("Estado = %q, se esperaba %q", d.Estado, domain.DefectoEnProgreso)
	}
	if f.defectos[100].Estado != domain.DefectoEnProgreso {
		t.Errorf("no quedó guardado: Estado = %q", f.defectos[100].Estado)
	}
}

// Traza: HU-21 CA-21.5 (RN-21.2: solo se inicia un defecto Abierto)
func TestActualizarEstadoDefecto_Iniciar_NoAbierto_RetornaErrorYNoGuarda(t *testing.T) {
	for _, estado := range []domain.EstadoDefecto{domain.DefectoEnProgreso, domain.DefectoResuelto, domain.DefectoCerrado} {
		t.Run(string(estado), func(t *testing.T) {
			f := nuevoEstadoDefectoFalso()
			d := f.defectos[100]
			d.Estado = estado
			f.defectos[100] = d

			_, err := actualizarEstadoDefectoCon(f).Iniciar(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

			if !errors.Is(err, domain.ErrDefectoNoAbierto) {
				t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoAbierto)
			}
			if f.vecesGuarda != 0 {
				t.Errorf("se llamó a ActualizarEstadoDefecto %d veces, no se esperaba ninguna", f.vecesGuarda)
			}
		})
	}
}

// --- Resolver ---

// Traza: HU-21 CA-21.2 (RN-21.6: se asigna el Sprint Abierto vigente hoy)
func TestActualizarEstadoDefecto_Resolver_ConSprintVigente_AsignaSuID(t *testing.T) {
	f := nuevoEstadoDefectoFalso()
	d := f.defectos[100]
	d.Estado = domain.DefectoEnProgreso
	f.defectos[100] = d
	f.sprints = []domain.Sprint{
		{ID: 5, ProyectoID: 1, Estado: domain.SprintAbierto,
			FechaInicio: diaEstadoDefecto(2026, time.October, 20), FechaFin: diaEstadoDefecto(2026, time.October, 30)},
	}

	got, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.Estado != domain.DefectoResuelto {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.DefectoResuelto)
	}
	if got.SprintResolucionID == nil || *got.SprintResolucionID != 5 {
		t.Errorf("SprintResolucionID = %v, se esperaba un puntero a 5", got.SprintResolucionID)
	}
	if f.defectos[100].Estado != domain.DefectoResuelto {
		t.Errorf("no quedó guardado: Estado = %q", f.defectos[100].Estado)
	}
}

// Traza: HU-21 CA-21.3 (RN-21.6: sin sprint vigente, queda sin asignar)
func TestActualizarEstadoDefecto_Resolver_SinSprintVigente_QuedaSinAsignar(t *testing.T) {
	f := nuevoEstadoDefectoFalso()
	d := f.defectos[100]
	d.Estado = domain.DefectoEnProgreso
	f.defectos[100] = d
	f.sprints = []domain.Sprint{ // ya terminó: no incluye "hoy" (2026-10-25)
		{ID: 5, ProyectoID: 1, Estado: domain.SprintAbierto,
			FechaInicio: diaEstadoDefecto(2026, time.September, 1), FechaFin: diaEstadoDefecto(2026, time.September, 15)},
	}

	got, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.Estado != domain.DefectoResuelto {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.DefectoResuelto)
	}
	if got.SprintResolucionID != nil {
		t.Errorf("SprintResolucionID = %v, se esperaba nil", *got.SprintResolucionID)
	}
}

// Traza: HU-21 CA-21.6 (RN-21.3: solo se resuelve un defecto En progreso)
func TestActualizarEstadoDefecto_Resolver_NoEnProgreso_RetornaErrorYNoGuarda(t *testing.T) {
	for _, estado := range []domain.EstadoDefecto{domain.DefectoAbierto, domain.DefectoResuelto, domain.DefectoCerrado} {
		t.Run(string(estado), func(t *testing.T) {
			f := nuevoEstadoDefectoFalso()
			d := f.defectos[100]
			d.Estado = estado
			f.defectos[100] = d

			_, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

			if !errors.Is(err, domain.ErrDefectoNoEnProgreso) {
				t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoEnProgreso)
			}
			if f.vecesGuarda != 0 {
				t.Errorf("se llamó a ActualizarEstadoDefecto %d veces, no se esperaba ninguna", f.vecesGuarda)
			}
		})
	}
}

// --- Cerrar ---

// Traza: HU-21 CA-21.4 (RN-21.4: el Sprint de resolución no cambia al cerrar)
func TestActualizarEstadoDefecto_Cerrar_DesdeResuelto_QuedaCerradoConElMismoSprint(t *testing.T) {
	f := nuevoEstadoDefectoFalso()
	sprintID := int64(5)
	d := f.defectos[100]
	d.Estado, d.SprintResolucionID = domain.DefectoResuelto, &sprintID
	f.defectos[100] = d

	got, err := actualizarEstadoDefectoCon(f).Cerrar(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.Estado != domain.DefectoCerrado {
		t.Errorf("Estado = %q, se esperaba %q", got.Estado, domain.DefectoCerrado)
	}
	if got.SprintResolucionID == nil || *got.SprintResolucionID != sprintID {
		t.Errorf("SprintResolucionID = %v, se esperaba que siguiera en %d", got.SprintResolucionID, sprintID)
	}
}

// Traza: HU-21 CA-21.7 (RN-21.4, RN-21.5: solo se cierra un defecto Resuelto, no se saltea ni retrocede)
func TestActualizarEstadoDefecto_Cerrar_NoResuelto_RetornaErrorYNoGuarda(t *testing.T) {
	for _, estado := range []domain.EstadoDefecto{domain.DefectoAbierto, domain.DefectoEnProgreso, domain.DefectoCerrado} {
		t.Run(string(estado), func(t *testing.T) {
			f := nuevoEstadoDefectoFalso()
			d := f.defectos[100]
			d.Estado = estado
			f.defectos[100] = d

			_, err := actualizarEstadoDefectoCon(f).Cerrar(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

			if !errors.Is(err, domain.ErrDefectoNoResuelto) {
				t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoResuelto)
			}
			if f.vecesGuarda != 0 {
				t.Errorf("se llamó a ActualizarEstadoDefecto %d veces, no se esperaba ninguna", f.vecesGuarda)
			}
		})
	}
}

// --- Comunes a Iniciar, Resolver y Cerrar (RN-21.7) ---

// Traza: HU-21 CA-21.8
func TestActualizarEstadoDefecto_ProyectoInexistente_RetornaError(t *testing.T) {
	f := nuevoEstadoDefectoFalso()

	_, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 99, DefectoID: 100})

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if f.vecesGuarda != 0 {
		t.Errorf("se llamó a ActualizarEstadoDefecto %d veces, no se esperaba ninguna", f.vecesGuarda)
	}
}

// Traza: HU-21 CA-21.8
func TestActualizarEstadoDefecto_DefectoInexistente_RetornaError(t *testing.T) {
	f := nuevoEstadoDefectoFalso()

	_, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 999})

	if !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoEncontrado)
	}
}

// Traza: HU-21 CA-21.8, RN-21.7: un defecto de otro proyecto se trata como inexistente.
func TestActualizarEstadoDefecto_DefectoDeOtroProyecto_RetornaError(t *testing.T) {
	f := nuevoEstadoDefectoFalso()

	_, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 200})

	if !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Errorf("error = %v, se esperaba %v", err, domain.ErrDefectoNoEncontrado)
	}
	if got := f.defectos[200].Estado; got != domain.DefectoAbierto {
		t.Errorf("el defecto de otro proyecto cambió: Estado = %q", got)
	}
}

func TestActualizarEstadoDefecto_FallaDelRepositorio_SePropaga(t *testing.T) {
	errBase := errors.New("base caída")
	tests := []struct {
		nombre string
		fallar func(*estadoDefectoFalso)
	}{
		{"al buscar el proyecto", func(f *estadoDefectoFalso) { f.errProyecto = errBase }},
		{"al buscar el defecto", func(f *estadoDefectoFalso) { f.errBuscar = errBase }},
		{"al listar los sprints", func(f *estadoDefectoFalso) { f.errSprints = errBase }},
		{"al guardar", func(f *estadoDefectoFalso) { f.errActualiza = errBase }},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			f := nuevoEstadoDefectoFalso()
			d := f.defectos[100]
			d.Estado = domain.DefectoEnProgreso // para llegar a listar sprints y guardar
			f.defectos[100] = d
			tt.fallar(f)

			_, err := actualizarEstadoDefectoCon(f).Resolver(context.Background(), DatosActualizarEstadoDefecto{ProyectoID: 1, DefectoID: 100})

			if !errors.Is(err, errBase) {
				t.Errorf("error = %v, se esperaba que envuelva %v", err, errBase)
			}
		})
	}
}
