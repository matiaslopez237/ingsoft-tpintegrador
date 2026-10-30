# Sprint 2: Review

> **Borrador.** Los datos (historias, puntos, PR y fechas) salen del repositorio y del tablero. Lo que está entre corchetes se completa en la reunión.

**Fecha:** [completar]  ·  **Asistentes:** [completar]
**Fechas del sprint:** planificado para el 26 y 27 de octubre de 2026; la última historia se mergeó el 29/10.
**Sprint Goal:** gestionar el backlog y armar sprints de un proyecto desde una interfaz web. **¿Cumplido?** Sí, con dos días de atraso: las 8 historias están en `main`, cada una con su pantalla.

| Historia | SP | Responsable | Estado | PR | Demo / observaciones |
|---|---|---|---|---|---|
| HU-30 Interfaz web base | 5 | Nico | Terminada (26/10) | #106 | [completar] |
| HU-08 Crear un sprint con Sprint Goal | 3 | Matías | Terminada (26/10) | #180 | [completar] |
| HU-06 Editar o eliminar un elemento del backlog | 3 | Nico | Terminada (27/10) | #182 | [completar] |
| HU-07 Listar y ordenar el backlog por prioridad | 2 | Matías | Terminada (26/10) | #183 | [completar] |
| HU-09 Asignar historias del backlog a un sprint | 5 | Eze | Terminada (28/10) | #184 | [completar] |
| HU-10 Marcar una historia como completada | 2 | Lucas | Terminada (29/10) | #186 | [completar] |
| HU-11 Cerrar un sprint | 3 | Lucas | Terminada (29/10) | #187 | Al cerrar, las historias sin terminar vuelven al Product Backlog y el sprint guarda sus puntos planificados y completados (cambio decidido en la revisión del PR). |
| HU-20 Registrar un defecto | 3 | Eze | Terminada (29/10) | #185 | [completar] |

**SP planificados:** 26  ·  **SP completados:** 26  ·  **Velocidad:** 26 (Sprint 1: 19)

**Evidencia de proceso:**
- 56 commits `[RED]`, 54 `[GREEN]` y 12 `[REFACTOR]` entre el 26 y el 29/10.
- ADR-0003: interfaz web con `html/template`.
- Pedidos simultáneos resueltos verificando la regla en el mismo `UPDATE` o `INSERT` (HU-08 a HU-11), con tests de concurrencia.
- 221 escenarios BDD en verde al cerrar el sprint.

**Defectos detectados en el sprint:** 2  ·  **Resueltos:** 2
| Defecto | Detectado en | Resuelto en |
|---|---|---|
| DEF-01: el formulario de edición de proyecto pisaba cambios simultáneos (#179) | revisión del PR #106 (HU-30), 26/10 | PR #181, 26/10 |
| DEF-02: eliminar un elemento con defectos informaba "con esfuerzo registrado" (#205) | revisión del PR #185 (HU-20), 29/10 | PR #261, 29/10 |

**Feedback del Product Architect:** la cátedra va a hacer hincapié en la integración y el despliegue continuos en la presentación. Se atiende en el PR #281 (ADR-0004). [completar con otros comentarios]

**Cambios al Product Backlog:**
- HU-11: cambió la decisión P-4 (las historias sin terminar vuelven al Product Backlog) y se agregó el resumen de puntos al cerrar (RN-11.6), que va a usar HU-22.
- Las tareas del Sprint 3 se cargaron como sub-issues de cada historia.
- La entrega se postergó: el Sprint 3 pasa a durar 3 días (29 al 31/10) y el Sprint 4 queda a definir.
