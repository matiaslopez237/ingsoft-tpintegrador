# Sprint 3: Review

> **Borrador para validar con el equipo.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Las observaciones de la demo y las conclusiones son una propuesta armada a partir de lo que pasó en el sprint: se revisan y se corrigen en la presentación final. El feedback del Product Architect se completa con lo que digan los profesores.

**Fecha:** 5 de noviembre de 2026  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Fechas del sprint:** del 29 al 31 de octubre de 2026 (planificadas). Las historias se mergearon entre el 1 y el 5 de noviembre.
**Sprint Goal:** estimar historias en equipo con Planning Poker y medir el trabajo: velocidad, esfuerzo estimado contra real y estado de los defectos. **¿Cumplido?** Sí: las 10 historias están en `main`.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-22 Story points y velocidad del equipo | 5 | Eze | Terminada (4/11) | #292 | Se mostró la velocidad del equipo y los puntos planificados y completados de cada sprint cerrado. |
| HU-14 Estimaciones individuales ocultas | 5 | Nico | Terminada (1/11) | #286 | Se votó una historia con Planning Poker: los votos quedan ocultos hasta que votan todos. |
| HU-15 Mostrar las estimaciones y detectar diferencias | 3 | Nico | Terminada (1/11) | #287 | Se revelaron los votos y el sistema marcó la diferencia significativa entre ellos. |
| HU-16 Repetir una ronda de votación | 2 | Matías | Terminada (2/11) | #291 | Se repitió la ronda sobre la misma historia, con el historial de la ronda anterior. |
| HU-17 Registrar la estimación acordada | 2 | Matías | Terminada (3/11) | #295 | Se registró la estimación acordada, que pasa a ser los story points de la historia. |
| HU-19 Comparar esfuerzo estimado contra real | 3 | Lucas | Terminada (4/11) | #297 | Se comparó el esfuerzo estimado contra el real de una historia, con su desviación. |
| HU-23 Horas estimadas, reales y desviación | 3 | Lucas | Terminada (5/11) | #298 | Se mostraron las horas estimadas, reales y la desviación de un sprint y del proyecto. |
| HU-21 Actualizar el estado de un defecto | 2 | Eze | Terminada (4/11) | #294 | Se cambió el estado de un defecto hasta cerrarlo, con su sprint de resolución. |
| HU-12 Consultar sprints anteriores | 2 | Matías | Terminada (1/11) | #285 | Se consultaron los sprints anteriores del proyecto, con sus historias. |
| HU-13 Escala de story points (Fibonacci) | 1 | Nico | Terminada (3/11) | #288 | Se estimó una historia con la escala Fibonacci; un valor fuera de la escala da error. |

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
