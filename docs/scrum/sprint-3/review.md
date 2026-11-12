# Sprint 3: Review

> **Borrador.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Lo que está entre corchetes se completa en la reunión.

**Fecha:** [completar]  ·  **Asistentes:** [completar]
**Fechas del sprint:** del 29 al 31 de octubre de 2026 (planificadas). Las historias se mergearon entre el 1 y el 5 de noviembre.
**Sprint Goal:** estimar historias en equipo con Planning Poker y medir el trabajo: velocidad, esfuerzo estimado contra real y estado de los defectos. **¿Cumplido?** Sí: las 10 historias están en `main`.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-22 Story points y velocidad del equipo | 5 | Eze | Terminada (4/11) | #292 | [completar] |
| HU-14 Estimaciones individuales ocultas | 5 | Nico | Terminada (1/11) | #286 | [completar] |
| HU-15 Mostrar las estimaciones y detectar diferencias | 3 | Nico | Terminada (1/11) | #287 | [completar] |
| HU-16 Repetir una ronda de votación | 2 | Matías | Terminada (2/11) | #291 | [completar] |
| HU-17 Registrar la estimación acordada | 2 | Matías | Terminada (3/11) | #295 | [completar] |
| HU-19 Comparar esfuerzo estimado contra real | 3 | Lucas | Terminada (4/11) | #297 | [completar] |
| HU-23 Horas estimadas, reales y desviación | 3 | Lucas | Terminada (5/11) | #298 | [completar] |
| HU-21 Actualizar el estado de un defecto | 2 | Eze | Terminada (4/11) | #294 | [completar] |
| HU-12 Consultar sprints anteriores | 2 | Matías | Terminada (1/11) | #285 | [completar] |
| HU-13 Escala de story points (Fibonacci) | 1 | Nico | Terminada (3/11) | #288 | [completar] |

**SP planificados:** 28  ·  **SP completados:** 28  ·  **Velocidad:** 28 (Sprint 1: 19, Sprint 2: 26)

**Evidencia de proceso:**
- Specs SDD, escenarios BDD con godog y TDD con commits RED → GREEN en las 10 historias, con la trazabilidad de cada una en `docs/trazabilidad.md`.
- **CI/CD (#281, 1/11):** tests con `-race`, cobertura mínima del 80 %, `main` protegida (PR, una aprobación y CI en verde) y despliegue en Render (ADR-0004).
- Pantalla para registrar esfuerzo de HU-18 (#290), que había quedado del Sprint 1.

**Defectos detectados en el sprint:** 2 (DEF-02 y DEF-03)  ·  **Resueltos:** 2
- **DEF-02** (#205 → #261): eliminar un elemento con defectos decía que tenía esfuerzo registrado.
- **DEF-03** (#283 → #284): la pantalla de un sprint cerrado mostraba las acciones de asignar y quitar.
- Después del sprint apareció **DEF-04** (#296, 3/11): la carrera en las migraciones, que se corrige en el Sprint 4.

**Feedback del Product Architect:** [completar]

**Cambios al Product Backlog:**
- La entrega se postergó (ya no es el 2/11), así que el Sprint 3 usó un día más que los anteriores y el Sprint 4 quedó "a definir".
- HU-19 sumó las horas estimadas como un dato propio del elemento (decisión P-10, confirmada en el PR #297).
