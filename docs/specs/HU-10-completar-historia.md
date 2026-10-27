# HU-10: Marcar una historia como completada en el Sprint

**Historia:** Como Developer, quiero marcar una historia como completada dentro de un Sprint, para reflejar el avance real.
**Issue:** #21  ·  **Sprint:** 2  ·  **Story points:** 2
**Estado de la spec:** Aprobada

## 1. Objetivo
Registrar que un elemento del backlog (HU-05) asignado a un sprint (HU-09) quedó terminado, pasándolo al estado `Hecho`. Es el dato que después usan el cierre del sprint (HU-11) y las métricas de velocidad.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento de backlog de ese proyecto (RN-10.4). |

No hay cuerpo: el sprint es el que el elemento ya tiene asignado (`SprintID`, HU-09).

**En HTTP:** `POST /proyectos/{id}/backlog/{elementoID}/completar`, con la misma forma de ruta anidada que HU-18 y HU-09.

## 3. Salidas esperadas
- El elemento queda en estado `Hecho`. Su `SprintID` y el resto de sus datos no cambian.
- Si hay un error, no se modifica nada.
- **En HTTP:** `204` sin cuerpo; `404` si no existe el proyecto o el elemento; `400` para los demás errores de la sección 7.

## 4. Reglas de negocio
- **RN-10.1:** Marcar una historia como completada es pasar el elemento al estado `Hecho` (RN-05.8). Se puede completar desde `Pendiente` o desde `En progreso`.
- **RN-10.2:** Solo se completa un elemento asignado a un sprint (HU-09), y ese sprint tiene que estar `Abierto`.
- **RN-10.3:** Completar un elemento que ya está `Hecho` no es un error y no cambia nada.
- **RN-10.4:** El elemento tiene que pertenecer al proyecto de la ruta; si es de otro proyecto, se trata como inexistente, igual que RN-18.1 y RN-09.7.

## 5. Restricciones
- La regla RN-10.1 a RN-10.3 vive en el dominio (`internal/domain`) sin dependencias externas. Igual que HU-09, el dominio **no depende de `domain.Sprint`**: recibe si el sprint del elemento está abierto como un `bool`, calculado por el caso de uso con el puerto angosto de HU-09 (`SprintEstaActivo`).
- La existencia del proyecto y del elemento (RN-10.4) se resuelve en el caso de uso a través de puertos de repositorio.
- **Dependencia de HU-09:** el dominio y el caso de uso se pueden hacer ya (el `SprintID` del elemento y el puerto se definen en HU-09). La persistencia, el HTTP y la pantalla se conectan cuando HU-09 esté en `main`.
- **Pedidos simultáneos:** la regla se verifica también en el mismo `UPDATE` que guarda el estado (el elemento sigue en ese sprint y el sprint sigue `Abierto`), como en HU-02, HU-03 y HU-08. Si mientras tanto se quitó el elemento del sprint (HU-09) o se cerró el sprint (HU-11), no se guarda nada y se devuelve el error que corresponda (`ErrElementoSinSprint` o `ErrSprintCerrado`).
- **Pantalla** (ADR-0003): en el detalle del sprint de HU-09 (`/web/proyectos/{id}/sprints/{sprintID}`), cada elemento asignado que no está `Hecho` tiene un botón "Marcar como hecha" si el sprint está `Abierto`. Acción: `POST /web/proyectos/{id}/sprints/{sprintID}/completar` (campo `elemento_id`), mismo patrón que "Quitar" de HU-09; usa el mismo caso de uso que el endpoint JSON.
- Fuera de alcance: reabrir un elemento (volver de `Hecho` a otro estado) y pasar un elemento a `En progreso`.

## 6. Casos límite
- Elemento `Pendiente` o `En progreso` en un sprint Abierto → queda `Hecho`.
- Elemento ya `Hecho` en un sprint Abierto → sin efecto, sin error (RN-10.3).
- Elemento ya `Hecho` en un sprint Cerrado → `ErrSprintCerrado`: RN-10.2 se valida antes que RN-10.3.
- Elemento sin sprint (nunca asignado o quitado con HU-09) → `ErrElementoSinSprint`.
- Elemento que existe pero en otro proyecto → `ErrElementoNoEncontrado`.
- El sprint se cierra entre la validación y el guardado → `ErrSprintCerrado`, el elemento no cambia (sección 5).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05, se reutiliza) | "elemento de backlog no encontrado" (HTTP 404) |
| El elemento no está asignado a un sprint | `ErrElementoSinSprint` (HU-09, se reutiliza) | "el elemento no tiene un sprint asignado" (HTTP 400) |
| El sprint del elemento está Cerrado | `ErrSprintCerrado` (HU-09, se reutiliza) | "el sprint ya está cerrado" (HTTP 400) |

Si hay varios errores, se informa el primero en el orden de la tabla. Los textos de `ErrElementoSinSprint` y `ErrSprintCerrado` son los definidos por HU-09 en el dominio (`internal/domain/backlog_asignacion.go`); si cambian allá, se actualizan acá.

## 8. Criterios de aceptación
- **CA-10.1:** Dado un elemento `Pendiente` o `En progreso` asignado a un sprint Abierto, cuando lo marco como completado, entonces queda `Hecho` y sigue en ese sprint. (RN-10.1, RN-10.2)
- **CA-10.2:** Dado un elemento ya `Hecho` en un sprint Abierto, cuando lo marco como completado, entonces no hay error y sigue `Hecho`. (RN-10.3)
- **CA-10.3:** Dado un elemento sin sprint asignado, cuando lo marco como completado, entonces obtengo `ErrElementoSinSprint` y su estado no cambia. (RN-10.2)
- **CA-10.4:** Dado un elemento asignado a un sprint Cerrado, cuando lo marco como completado, entonces obtengo `ErrSprintCerrado` y su estado no cambia. (RN-10.2)
- **CA-10.5:** Dado un ID de proyecto o de elemento que no existe, o un elemento de otro proyecto, cuando lo marco como completado, entonces obtengo `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado`. (RN-10.4)
- **CA-10.6:** Dado un pedido válido a `POST /proyectos/{id}/backlog/{elementoID}/completar`, entonces la respuesta es `204`; con proyecto o elemento inexistente, `404`; sin sprint o con el sprint Cerrado, `400`. (Sección 3)

## Decisiones tomadas
Confirmadas por Matías López (Agile Enabler) el 2026-10-27. Eze y Nico pueden objetar en la revisión del PR.
- **P-1:** completar = pasar a `Hecho` (RN-10.1).
- **P-2:** el elemento tiene que estar en un sprint y el sprint tiene que estar Abierto (RN-10.2).
- **P-3:** completar algo ya `Hecho` no da error (RN-10.3).
- **P-4:** reabrir queda fuera de esta historia (sección 5).
- **P-5:** un elemento de otro proyecto se trata como inexistente (RN-10.4).
- **P-6:** ruta JSON `POST /proyectos/{id}/backlog/{elementoID}/completar` (`204`) y botón "Marcar como hecha" en el detalle del sprint con `POST /web/proyectos/{id}/sprints/{sprintID}/completar`.
- **P-7:** pedidos simultáneos: la regla se verifica en el mismo `UPDATE` que guarda (sección 5).
- **P-8:** elemento ya `Hecho` en un sprint Cerrado: se informa `ErrSprintCerrado` (RN-10.2 se valida antes que RN-10.3). La alternativa sería no dar error, por RN-10.3.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-10-26: creada (borrador), con las decisiones P-1 a P-8 pendientes de confirmar.
- 2026-10-27: Decisiones P-1 a P-8 confirmadas; spec aprobada.
