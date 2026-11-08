package domain

import (
	"regexp"
	"testing"
	"time"
)

func utc(anio int, mes time.Month, dia, hora, min int) time.Time {
	return time.Date(anio, mes, dia, hora, min, 0, 0, time.UTC)
}

// Traza: HU-27 CA-27.13: la emisión es el instante del reloj en hora de Argentina, UTC-3 fijo (RN-27.10).
func TestHoraDeEmision(t *testing.T) {
	casos := []struct {
		nombre string
		reloj  time.Time
		texto  string // dd/mm/aaaa HH:MM, como lo muestra el PDF
	}{
		{"mismo día", utc(2026, 11, 7, 17, 5), "07/11/2026 14:05"},
		{"el huso horario cambia el día", utc(2026, 11, 8, 2, 30), "07/11/2026 23:30"},
		{"medianoche en Argentina", utc(2026, 11, 8, 3, 0), "08/11/2026 00:00"},
		{"enero, sin horario de verano", utc(2027, 1, 15, 12, 0), "15/01/2027 09:00"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := HoraDeEmision(c.reloj)
			if _, offset := got.Zone(); offset != -3*60*60 {
				t.Errorf("huso = %d s, se esperaba UTC-3", offset)
			}
			if !got.Equal(c.reloj) {
				t.Errorf("HoraDeEmision cambió el instante: %v != %v", got, c.reloj)
			}
			if s := got.Format("02/01/2006 15:04"); s != c.texto {
				t.Errorf("emisión = %s, se esperaba %s", s, c.texto)
			}
		})
	}
}

// El reloj del servidor puede estar en cualquier huso: Render corre en UTC (sección 5).
func TestHoraDeEmision_RelojEnOtroHuso(t *testing.T) {
	madrid := time.FixedZone("UTC+1", 60*60)
	got := HoraDeEmision(time.Date(2026, 11, 7, 18, 5, 0, 0, madrid))
	if s := got.Format("02/01/2006 15:04"); s != "07/11/2026 14:05" {
		t.Errorf("emisión = %s, se esperaba 07/11/2026 14:05", s)
	}
}

var nombreSeguro = regexp.MustCompile(`^[a-z0-9.-]+$`)

// Traza: HU-27 CA-27.14, CA-27.15, CA-27.16: el nombre lleva los IDs y la fecha de emisión en UTC-3, y nada más
// (RN-27.12).
func TestNombreArchivoReporte(t *testing.T) {
	casos := []struct {
		nombre string
		got    string
		want   string
	}{
		{"proyecto", NombreArchivoReporteProyecto(3, utc(2026, 11, 7, 17, 5)), "reporte-proyecto-3-2026-11-07.pdf"},
		{"sprint", NombreArchivoReporteSprint(3, 12, utc(2026, 11, 7, 17, 5)), "reporte-proyecto-3-sprint-12-2026-11-07.pdf"},
		{"la fecha es la de Argentina", NombreArchivoReporteProyecto(3, utc(2026, 11, 8, 2, 30)), "reporte-proyecto-3-2026-11-07.pdf"},
		{"medianoche en Argentina", NombreArchivoReporteSprint(3, 12, utc(2026, 11, 8, 3, 0)), "reporte-proyecto-3-sprint-12-2026-11-08.pdf"},
		{"IDs grandes", NombreArchivoReporteSprint(9223372036854775807, 1, utc(2026, 11, 7, 17, 5)), "reporte-proyecto-9223372036854775807-sprint-1-2026-11-07.pdf"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if c.got != c.want {
				t.Errorf("nombre = %q, se esperaba %q", c.got, c.want)
			}
			if !nombreSeguro.MatchString(c.got) {
				t.Errorf("el nombre %q tiene caracteres fuera de [a-z0-9.-]", c.got)
			}
		})
	}
}
