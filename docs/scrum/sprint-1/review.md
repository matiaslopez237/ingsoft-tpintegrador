# Sprint 1: Review

> **Borrador para validar con el equipo.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Las observaciones de la demo y las conclusiones son una propuesta armada a partir de lo que pasó en el sprint: se revisan y se corrigen en la presentación final. El feedback del Product Architect se completa con lo que digan los profesores.

**Fecha:** 25 de octubre de 2026  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Fechas del sprint:** hasta el 25 de octubre de 2026
**Sprint Goal:** gestión de proyectos (Épica 1), crear elementos del backlog y registrar esfuerzo (MVP). **¿Cumplido?** Sí: las 6 historias están en `main`.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-01 Crear un proyecto | 3 | Nico | Terminada | #86 | Se creó un proyecto con nombre y fechas, y se mostraron los errores de validación (nombre vacío, fin antes del inicio). |
| HU-02 Modificar un proyecto | 2 | Matías | Terminada | #88 | Se modificó el proyecto creado; las reglas de validación son las mismas que al crearlo. |
| HU-03 Registrar integrantes de un proyecto | 3 | Eze | Terminada | #87 | Se registraron los cuatro integrantes del equipo con su rol; no se puede repetir un integrante. |
| HU-04 Consultar el estado de un proyecto | 3 | Lucas | Terminada | #89 | Se consultó el estado del proyecto con sus integrantes y su backlog. |
| HU-18 Registrar esfuerzo en una historia o tarea | 3 | Lucas | Terminada (mergeada el 26/10) | #104 | Se registraron horas en una historia por la API; la pantalla quedó para más adelante (se hizo en el Sprint 3, #290). |
| HU-05 Crear un elemento de backlog | 5 | Nico | Terminada | #92 | Se creó un elemento de backlog con prioridad, story points y criterios de aceptación. |

**SP planificados:** 14 al inicio, 19 después de sumar HU-05  ·  **SP completados:** 19  ·  **Velocidad:** 19

**Evidencia de proceso:**
- Specs SDD, escenarios BDD con godog y TDD con commits RED → GREEN → REFACTOR en las 6 historias (39 commits `[RED]`, 38 `[GREEN]` y 5 `[REFACTOR]` hasta el 25/10).
- Decisiones de arquitectura: ADR-0001 (`net/http`) y ADR-0002 (SQLite).
- CI en GitHub Actions: build, `go vet`, `gofmt` y tests en cada PR.

**Defectos detectados en el sprint:** ninguno registrado  ·  **Resueltos:** —

**Feedback del Product Architect:** [completar]

**Cambios al Product Backlog:**
- 23/10: se sumó HU-05 (5 SP) al sprint, por decisión del equipo en clase.
- Al cerrar el sprint se re-estimó el Product Backlog: de 115 a 98 SP, con la velocidad del Sprint 1 como referencia.
- Se agregó HU-30 (interfaz web base, 5 SP) para el Sprint 2 y se armó el plan de los Sprints 2 a 4 (PR #102).
