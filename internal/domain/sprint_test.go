package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// proyectoParaSprints es un proyecto activo del 1/10 al 31/12, ya guardado.
func proyectoParaSprints(t *testing.T) Proyecto {
	t.Helper()
	return Proyecto{
		ID:          7,
		Nombre:      "Alfa",
		FechaInicio: fecha(t, "2026-10-01"),
		FechaFin:    ptr(fecha(t, "2026-12-31")),
		Estado:      EstadoActivo,
	}
}

// sprintDel arma un sprint ya guardado del proyecto, para usar como existente.
func sprintDel(t *testing.T, proyectoID int64, inicio, fin string) Sprint {
	t.Helper()
	return Sprint{ID: 1, ProyectoID: proyectoID, Numero: 1, FechaInicio: fecha(t, inicio), FechaFin: fecha(t, fin),
		Objetivo: "Existente", Estado: SprintAbierto}
}

// Traza: HU-08 CA-08.1, RN-08.6, RN-08.8
func TestNuevoSprint_Valido(t *testing.T) {
	p := proyectoParaSprints(t)

	s, err := NuevoSprint(p, fecha(t, "2026-10-26"), fecha(t, "2026-10-27"), "  Gestionar el backlog  ", nil)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if s.ProyectoID != p.ID || !s.FechaInicio.Equal(fecha(t, "2026-10-26")) || !s.FechaFin.Equal(fecha(t, "2026-10-27")) {
		t.Errorf("sprint = %+v, se esperaba el proyecto %d del 26 al 27", s, p.ID)
	}
	if s.Objetivo != "Gestionar el backlog" {
		t.Errorf("Objetivo = %q, se esperaba recortado", s.Objetivo)
	}
	if s.Estado != SprintAbierto {
		t.Errorf("Estado = %q, se esperaba %q", s.Estado, SprintAbierto)
	}
	if s.ID != 0 || s.Numero != 0 {
		t.Errorf("ID = %d, Numero = %d: los asigna la persistencia al guardar (RN-08.7, RN-08.9)", s.ID, s.Numero)
	}
}

// Traza: HU-08 CA-08.4 a CA-08.8, RN-08.2 a RN-08.6, sección 7
func TestNuevoSprint_Invalido_RetornaError(t *testing.T) {
	cero := time.Time{}
	tests := []struct {
		nombre         string
		estado         Estado
		inicio, fin    string
		sinInicio      bool
		sinFin         bool
		objetivo       string
		existenteDesde string // vacío: el proyecto no tiene sprints
		existenteHasta string
		wantErr        error
	}{
		{nombre: "proyecto en pausa", estado: EstadoEnPausa, inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x", wantErr: ErrProyectoNoActivo},
		{nombre: "proyecto finalizado", estado: EstadoFinalizado, inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x", wantErr: ErrProyectoNoActivo},
		{nombre: "proyecto postergado", estado: EstadoPostergado, inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x", wantErr: ErrProyectoNoActivo},
		{nombre: "proyecto cancelado", estado: EstadoCancelado, inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x", wantErr: ErrProyectoNoActivo},
		{nombre: "sin fecha de inicio", sinInicio: true, fin: "2026-10-27", objetivo: "x", wantErr: ErrFechaInicioObligatoria},
		{nombre: "sin fecha de fin", inicio: "2026-10-26", sinFin: true, objetivo: "x", wantErr: ErrFechaFinObligatoria},
		{nombre: "fin anterior al inicio", inicio: "2026-10-28", fin: "2026-10-27", objetivo: "x", wantErr: ErrFechaFinAnterior},
		{nombre: "Sprint Goal vacío", inicio: "2026-10-26", fin: "2026-10-27", objetivo: "", wantErr: ErrObjetivoVacio},
		{nombre: "Sprint Goal de solo espacios", inicio: "2026-10-26", fin: "2026-10-27", objetivo: "   ", wantErr: ErrObjetivoVacio},
		{nombre: "Sprint Goal de 201 caracteres", inicio: "2026-10-26", fin: "2026-10-27", objetivo: strings.Repeat("a", 201), wantErr: ErrObjetivoDemasiadoLargo},
		{nombre: "empieza antes que el proyecto", inicio: "2026-09-30", fin: "2026-10-02", objetivo: "x", wantErr: ErrSprintFueraDelProyecto},
		{nombre: "termina después del fin del proyecto", inicio: "2026-12-30", fin: "2027-01-01", objetivo: "x", wantErr: ErrSprintFueraDelProyecto},
		{nombre: "comparte el último día de otro", inicio: "2026-10-27", fin: "2026-10-28", objetivo: "x", existenteDesde: "2026-10-26", existenteHasta: "2026-10-27", wantErr: ErrSprintSuperpuesto},
		{nombre: "comparte el primer día de otro", inicio: "2026-10-25", fin: "2026-10-26", objetivo: "x", existenteDesde: "2026-10-26", existenteHasta: "2026-10-27", wantErr: ErrSprintSuperpuesto},
		{nombre: "contiene a otro", inicio: "2026-10-20", fin: "2026-10-30", objetivo: "x", existenteDesde: "2026-10-26", existenteHasta: "2026-10-27", wantErr: ErrSprintSuperpuesto},
		{nombre: "está contenido en otro", inicio: "2026-10-26", fin: "2026-10-26", objetivo: "x", existenteDesde: "2026-10-26", existenteHasta: "2026-10-27", wantErr: ErrSprintSuperpuesto},
		{nombre: "mismas fechas que otro", inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x", existenteDesde: "2026-10-26", existenteHasta: "2026-10-27", wantErr: ErrSprintSuperpuesto},
		// Sección 7: si hay varios errores se informa el primero de la tabla.
		{nombre: "proyecto no activo y Sprint Goal vacío", estado: EstadoCancelado, inicio: "2026-10-26", fin: "2026-10-27", objetivo: "", wantErr: ErrProyectoNoActivo},
		{nombre: "Sprint Goal vacío y fuera del proyecto", inicio: "2026-09-01", fin: "2026-09-02", objetivo: "", wantErr: ErrObjetivoVacio},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p := proyectoParaSprints(t)
			if tt.estado != "" {
				p.Estado = tt.estado
			}
			inicio, fin := cero, cero
			if !tt.sinInicio {
				inicio = fecha(t, tt.inicio)
			}
			if !tt.sinFin {
				fin = fecha(t, tt.fin)
			}
			var existentes []Sprint
			if tt.existenteDesde != "" {
				existentes = []Sprint{sprintDel(t, p.ID, tt.existenteDesde, tt.existenteHasta)}
			}

			_, err := NuevoSprint(p, inicio, fin, tt.objetivo, existentes)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-08 CA-08.5 a CA-08.8 (casos límite válidos), RN-08.3 a RN-08.6
func TestNuevoSprint_CasosLimiteValidos(t *testing.T) {
	tests := []struct {
		nombre      string
		inicio, fin string
		objetivo    string
		sinFinProy  bool
		existentes  []Sprint
	}{
		{nombre: "sprint de un día", inicio: "2026-10-26", fin: "2026-10-26", objetivo: "x"},
		{nombre: "Sprint Goal de 200 caracteres", inicio: "2026-10-26", fin: "2026-10-27", objetivo: strings.Repeat("a", 200)},
		{nombre: "empieza el mismo día que el proyecto", inicio: "2026-10-01", fin: "2026-10-02", objetivo: "x"},
		{nombre: "termina el mismo día que el proyecto", inicio: "2026-12-30", fin: "2026-12-31", objetivo: "x"},
		{nombre: "proyecto sin fecha de fin", inicio: "2027-05-01", fin: "2027-05-14", objetivo: "x", sinFinProy: true},
		{nombre: "empieza el día siguiente al fin de otro", inicio: "2026-10-28", fin: "2026-10-29", objetivo: "x",
			existentes: []Sprint{{ProyectoID: 7, FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)}}},
		{nombre: "termina el día anterior al inicio de otro", inicio: "2026-10-24", fin: "2026-10-25", objetivo: "x",
			existentes: []Sprint{{ProyectoID: 7, FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)}}},
		{nombre: "los sprints de otro proyecto no cuentan", inicio: "2026-10-26", fin: "2026-10-27", objetivo: "x",
			existentes: []Sprint{{ProyectoID: 99, FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)}}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p := proyectoParaSprints(t)
			if tt.sinFinProy {
				p.FechaFin = nil
			}
			if _, err := NuevoSprint(p, fecha(t, tt.inicio), fecha(t, tt.fin), tt.objetivo, tt.existentes); err != nil {
				t.Errorf("error inesperado: %v", err)
			}
		})
	}
}

// Traza: HU-08 RN-08.5: las fechas se comparan por día aunque vengan con hora o en otra zona.
func TestNuevoSprint_ComparaPorDiaCalendario(t *testing.T) {
	p := proyectoParaSprints(t)
	tarde := time.FixedZone("UTC-3", -3*60*60)
	existente := []Sprint{{ProyectoID: p.ID,
		FechaInicio: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), FechaFin: time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC)}}

	// El 28 a las 22 h en UTC-3 ya es 29 en UTC, pero como día es el 28: no se superpone.
	_, err := NuevoSprint(p, time.Date(2026, 10, 28, 22, 0, 0, 0, tarde), time.Date(2026, 10, 29, 22, 0, 0, 0, tarde), "x", existente)
	if err != nil {
		t.Errorf("error inesperado: %v", err)
	}
	// El 27 a las 23 h en UTC-3 es el mismo día que el fin del existente: se superpone.
	_, err = NuevoSprint(p, time.Date(2026, 10, 27, 23, 0, 0, 0, tarde), time.Date(2026, 10, 29, 0, 0, 0, 0, tarde), "x", existente)
	if !errors.Is(err, ErrSprintSuperpuesto) {
		t.Errorf("error = %v, se esperaba %v", err, ErrSprintSuperpuesto)
	}
}
