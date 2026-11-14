# Sprint 4: Review

> **Borrador para validar con el equipo.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Las observaciones de la demo y las conclusiones son una propuesta armada a partir de lo que pasó en el sprint: se revisan y se corrigen en la presentación final. El feedback del Product Architect se completa con lo que digan los profesores.

**Fecha:** [completar: día de la presentación]  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana y los profesores (Product Architect)
**Fechas del sprint:** del 5 al 13 de noviembre de 2026. Empezó hasta el 8; se extendió el 9/11 hasta el 10, y el 11/11 hasta el 13. Se cerró el 12/11.
**Sprint Goal:** ver el avance del proyecto de un vistazo, compartirlo en un reporte en PDF y dejar el sistema documentado para la entrega. **¿Cumplido?** Sí: se terminaron las 7 historias.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-24 % de story points completados y defectos | 3 | Eze | Terminada (6/11) | #366 | Se mostró el porcentaje de story points completados y los defectos detectados y resueltos de cada sprint. |
| HU-25 Dashboard con el estado del proyecto y gráficos | 8 | Nico | Terminada (8/11) | #375 | Se mostró el Dashboard del proyecto de ejemplo, con los gráficos de velocidad, horas y defectos. |
| HU-26 Reporte de un proyecto o sprint | 5 | Eze | Terminada (8/11) | #378 | Se generó el reporte del proyecto y de un sprint, con sus cinco secciones. |
| HU-27 Exportar el reporte a PDF | 5 | Lucas | Terminada (8/11) | #379 | Se descargó el reporte en PDF, con acentos y tablas. |
| HU-29 Documentación técnica | 2 | Matías | Terminada (6/11) | #365 | Se recorrió `docs/arquitectura.md` y los ADR. |
| HU-28 Manual breve de usuario | 2 | Matías | Terminada (12/11) | #384 | Se recorrió `docs/manual-usuario.md`, con las capturas del proyecto de ejemplo. |
| HU-31 Proyecto de ejemplo precargado (sumada durante el sprint) | 3 | Matías | Terminada (11/11) | #377 | La versión publicada en Render arranca con el proyecto de ejemplo, con los cuatro sprints del TP. |

**SP planificados:** 25 al inicio, 28 con HU-31  ·  **SP completados:** 28 de 28 (100 %)  ·  **Velocidad del sprint:** 28  ·  **Velocidad del equipo:** 27,3 (promedio de los Sprints 2 a 4: 26, 28 y 28)

**Evidencia de proceso:**
- Specs SDD, escenarios BDD y TDD con commits RED → GREEN en las historias con código; HU-28 y HU-29 son de documentación y no los necesitan.
- **Decisiones de arquitectura:** ADR-0005 (gráficos en SVG armados en Go, sin JavaScript) y ADR-0006 (PDF con `fpdf` y la fuente DejaVu), los dos aceptados en la revisión de su PR.
- **Despliegue continuo activo desde el 5/11:** cada merge a `main` con el CI en verde se despliega solo en Render. Desde HU-31, la versión publicada arranca con el proyecto de ejemplo.
- **CI (#376):** el timeout de los tests subió a 20 minutos, porque la suite ya tardaba casi 9.

**Defectos detectados en el sprint:** 1 (DEF-05; DEF-04 se detectó entre el Sprint 3 y el 4)  ·  **Resueltos:** 2 (DEF-04 y DEF-05). Al cierre no queda ningún defecto abierto.
- **DEF-05** (#380 → #381, resuelto el 11/11): la lista de defectos, el reporte y el PDF mostraban el ID del sprint en lugar de su número; además, el Resumen del reporte no tenía la cantidad de sprints y la velocidad salía sin redondear. Se vio en la revisión de HU-27.
- **DEF-04** (#296 → #382, resuelto el 12/11): la carrera en las migraciones con arranques simultáneos, que venía del Sprint 3.

**Feedback del Product Architect:** [completar]

**Cambios al Product Backlog:**
- Se sumó HU-31 (proyecto de ejemplo precargado, 3 SP) durante el sprint, para que la versión publicada muestre un proyecto completo.
- HU-24 se expresó en story points en lugar de historias, porque el cierre de un sprint guarda los puntos y no la cantidad de historias.
- La conversión del reporte a PDF pasó de HU-26 a HU-27.
