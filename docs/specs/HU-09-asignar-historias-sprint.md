# HU-09: Asignar historias del backlog a un Sprint

**Historia:** Como equipo, quiero asignar historias del Product Backlog a un Sprint, para definir el Sprint Backlog.
**Issue:** #20  ·  **Sprint:** 2  ·  **Story points:** 5
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir vincular un elemento del backlog (HU-05) a un sprint del mismo proyecto (HU-08), para formar el Sprint Backlog, y desvincularlo si hace falta mientras el sprint sigue abierto.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe existir un proyecto con ese id; el elemento tiene que pertenecerle (RN-09.7) |
| `elementoID` | `int64` | Sí | Debe existir un elemento de backlog con ese id |
| `sprintID` | `int64` | Sí (al asignar) | Debe existir un sprint con ese id, del mismo proyecto que el elemento |

**En HTTP:** `POST /proyectos/{id}/backlog/{elementoID}/sprint` con `{ "sprint_id": 2 }` para asignar; `DELETE /proyectos/{id}/backlog/{elementoID}/sprint` para desasignar. La ruta sigue la misma forma que usa HU-18 para los elementos del backlog (`/proyectos/{id}/backlog/{elementoID}/esfuerzo`), no `/elementos/` como se había escrito en el borrador original de esta spec.

## 3. Salidas esperadas
- **Asignar:** el elemento queda con `SprintID` apuntando al sprint indicado.
- **Desasignar:** el elemento queda con `SprintID` en `nil` (sin sprint).
- No hay un "dato" adicional que devolver más allá de confirmar la operación; consultar qué elementos tiene un sprint es responsabilidad de HU-09 también a través de `GET /proyectos/{id}/sprints/{sprintID}/elementos` (ver Preguntas abiertas) o de HU-07 (listar backlog) si ya muestra el sprint asignado.

## 4. Reglas de negocio
- **RN-09.1:** Un elemento de backlog está asignado a lo sumo a un sprint a la vez.
- **RN-09.2 (CA del issue #20):** No se puede asignar un elemento a un sprint si ya está asignado a **otro** sprint que sigue **Abierto**. Si no tiene sprint asignado, o el que tiene ya está **Cerrado**, se puede asignar (o reasignar) sin problema.
- **RN-09.3:** Asignar un elemento al sprint al que ya está asignado no es un error (operación sin efecto).
- **RN-09.4 (CA del issue #20):** Se puede desasignar un elemento de su sprint mientras ese sprint sigue **Abierto**. Una vez que el sprint está **Cerrado**, no se puede desasignar (queda como registro histórico de qué se completó en ese sprint).
- **RN-09.5:** El sprint de destino tiene que pertenecer al mismo proyecto que el elemento de backlog; no se puede asignar un elemento de un proyecto a un sprint de otro.
- **RN-09.6:** Solo se puede asignar a un sprint que esté **Abierto** (no tiene sentido sumar trabajo nuevo a un sprint ya cerrado).
- **RN-09.7:** El elemento tiene que pertenecer al proyecto de la ruta (`/proyectos/{id}/backlog/{elementoID}/sprint`); si es de otro proyecto, se trata como inexistente (`ErrElementoNoEncontrado`). Misma regla que RN-18.1 de HU-18, para la misma forma de ruta anidada.

## 5. Restricciones
- El núcleo (`internal/domain`) **no depende de `domain.Sprint`** (de HU-08, todavía no mergeada a `main` al momento de escribir esta spec). Las reglas que necesitan saber si un sprint está activo (RN-09.2, RN-09.4, RN-09.6) reciben ese dato como un `bool` simple, calculado por el caso de uso. Esto permite implementar dominio y caso de uso sin esperar el merge de HU-08; solo la persistencia, el HTTP y la pantalla (Tareas 9.6 a 9.8) necesitan que exista la tabla real de sprints.
- El puerto de persistencia que necesita el caso de uso para consultar el estado del sprint es angosto a propósito (p. ej. `SprintEstaActivo(ctx, sprintID) (bool, error)`), en vez de devolver un `domain.Sprint` completo, para no acoplar con el tipo de HU-08.

## 6. Casos límite
- Elemento nunca asignado (`SprintID` nil): se puede asignar directo, sin pasar por la validación de "ya asignado".
- Reasignar un elemento al mismo sprint que ya tiene: válido, no cambia nada (RN-09.3).
- Desasignar un elemento que no tiene sprint asignado: error (no hay nada que desasignar).
- Asignar/desasignar sobre un sprint que ya está Cerrado: error en ambos casos (RN-09.4, RN-09.6).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Elemento inexistente | `domain.ErrElementoNoEncontrado` (ya definido en HU-05, se reutiliza) | 404 |
| Sprint inexistente | `ErrSprintNoEncontrado` (nuevo; puede coincidir con el de HU-08 si ya existe cuando se implemente la persistencia) | 404 |
| Sprint de otro proyecto | `ErrSprintDeOtroProyecto` | 400 |
| Elemento ya asignado a otro sprint activo | `ErrElementoYaAsignado` | 400 |
| Asignar a un sprint Cerrado | `ErrSprintCerrado` | 400 |
| Desasignar un elemento sin sprint | `ErrElementoSinSprint` | 400 |
| Desasignar de un sprint ya Cerrado | `ErrSprintCerrado` (mismo error que al asignar) | 400 |

## 8. Criterios de aceptación
- **CA-09.1:** Dado un elemento de backlog sin sprint asignado y un sprint Abierto del mismo proyecto, cuando lo asigno, entonces el elemento queda con ese `SprintID`.
- **CA-09.2:** Dado un elemento ya asignado a un sprint Abierto, cuando intento asignarlo a otro sprint, entonces la operación falla con `ErrElementoYaAsignado` y no cambia nada. *(CA del issue #20)*
- **CA-09.3:** Dado un elemento asignado a un sprint que ya está Cerrado, cuando lo asigno a otro sprint Abierto, entonces se reasigna correctamente.
- **CA-09.4:** Dado un elemento asignado a un sprint Abierto, cuando lo desasigno, entonces queda sin sprint. *(CA del issue #20)*
- **CA-09.5:** Dado un elemento asignado a un sprint que ya está Cerrado, cuando intento desasignarlo, entonces la operación falla con `ErrSprintCerrado` y no cambia nada.
- **CA-09.6:** Dado un elemento sin sprint asignado, cuando intento desasignarlo, entonces la operación falla con `ErrElementoSinSprint`.
- **CA-09.7:** Dado un sprint de otro proyecto, cuando intento asignarle un elemento, entonces la operación falla con `ErrSprintDeOtroProyecto`.
- **CA-09.8:** Dado un sprint ya Cerrado, cuando intento asignarle un elemento, entonces la operación falla con `ErrSprintCerrado`.

## Preguntas abiertas
- ¿Hace falta un endpoint/pantalla para listar qué elementos tiene asignado un sprint (`GET /proyectos/{id}/sprints/{sprintID}/elementos`), o alcanza con que HU-07 (listar backlog) muestre el sprint de cada elemento? Por ahora se asume que no hace falta un endpoint nuevo de solo lectura para esta historia; si el equipo lo pide, se agrega como extensión menor.
- ¿El estado del elemento (`Pendiente`/`En progreso`/`Hecho`) restringe si se lo puede asignar o desasignar? Por ahora se asume que no — cualquier estado del elemento se puede asignar/desasignar, ya que no está en los criterios de aceptación del issue.
- RN-09.5 y RN-09.6 (mismo proyecto, sprint Abierto) no estaban en los 2 criterios del issue #20; se agregaron por ser restricciones de integridad razonables. Si el equipo no está de acuerdo, avisar antes de implementar Tarea 9.4.

## Historial
- 2026-10-26: creada (borrador).
