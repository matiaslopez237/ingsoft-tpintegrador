package domain

import (
	"fmt"
	"time"
)

// zonaEmision es la hora de Argentina, UTC-3 sin horario de verano (RN-27.10). Es fija y no
// time.LoadLocation, que necesitaría la base de husos del sistema o time/tzdata (spec, sección 5).
var zonaEmision = time.FixedZone("UTC-3", -3*60*60)

// HoraDeEmision devuelve el mismo instante en la hora de Argentina (RN-27.10).
func HoraDeEmision(t time.Time) time.Time {
	return t.In(zonaEmision)
}

// NombreArchivoReporteProyecto es el nombre del PDF del reporte de un proyecto (RN-27.12): solo el ID y
// la fecha de emisión en UTC-3, sin el nombre del proyecto.
func NombreArchivoReporteProyecto(proyectoID int64, emision time.Time) string {
	return fmt.Sprintf("reporte-proyecto-%d-%s.pdf", proyectoID, HoraDeEmision(emision).Format(time.DateOnly))
}

// NombreArchivoReporteSprint es el nombre del PDF del reporte de un sprint (RN-27.12).
func NombreArchivoReporteSprint(proyectoID, sprintID int64, emision time.Time) string {
	return fmt.Sprintf("reporte-proyecto-%d-sprint-%d-%s.pdf", proyectoID, sprintID, HoraDeEmision(emision).Format(time.DateOnly))
}
