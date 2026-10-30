# Sprint 1: Review

> **Borrador.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Lo que está entre corchetes se completa en la reunión.

**Fecha:** [completar]  ·  **Asistentes:** [completar]
**Fechas del sprint:** hasta el 25 de octubre de 2026
**Sprint Goal:** gestión de proyectos (Épica 1), crear elementos del backlog y registrar esfuerzo (MVP). **¿Cumplido?** Sí: las 6 historias están en `main`.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-01 Crear un proyecto | 3 | Nico | Terminada | #86 | [completar] |
| HU-02 Modificar un proyecto | 2 | Matías | Terminada | #88 | [completar] |
| HU-03 Registrar integrantes de un proyecto | 3 | Eze | Terminada | #87 | [completar] |
| HU-04 Consultar el estado de un proyecto | 3 | Lucas | Terminada | #89 | [completar] |
| HU-18 Registrar esfuerzo en una historia o tarea | 3 | Lucas | Terminada (mergeada el 26/10) | #104 | [completar] |
| HU-05 Crear un elemento de backlog | 5 | Nico | Terminada | #92 | [completar] |

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
