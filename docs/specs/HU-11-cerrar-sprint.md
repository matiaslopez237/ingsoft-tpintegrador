# HU-11: Cerrar un Sprint

**Historia:** Como Scrum Master, quiero cerrar un Sprint, para dar por finalizado el período de trabajo.
**Issue:** #22 (a confirmar)  ·  **Sprint:** 2  ·  **Story points:** 3
**Estado de la spec:** Aprobada

## 1. Objetivo
Dar por terminado un sprint de un proyecto (HU-08), pasándolo de `Abierto` a `Cerrado`. Al cerrarlo, el sprint registra sus story points planificados y completados, y las historias sin terminar vuelven al Product Backlog. Un sprint cerrado queda como registro histórico: ya no se le asignan, quitan ni completan historias (lo controlan HU-09 y HU-10), y es la base de las métricas por sprint (Épica 7).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `sprintID` | `int64` | Sí | Debe corresponder a un sprint de ese proyecto (RN-11.4). |

No hay cuerpo. **En HTTP:** `POST /proyectos/{id}/sprints/{sprintID}/cerrar`.

## 3. Salidas esperadas
- El sprint queda en estado `Cerrado`. Su número, fechas y Sprint Goal no cambian.
- Las historias sin terminar (`Pendiente` o `En progreso`) quedan sin sprint, de vuelta en el Product Backlog; las `Hecho` siguen asignadas al sprint cerrado (RN-11.3).
- El sprint guarda sus story points planificados y completados en el momento del cierre (RN-11.6).
- Si hay un error, no se modifica nada.
- **En HTTP:** `204` sin cuerpo; `404` si no existe el proyecto o el sprint; `400` si el sprint ya estaba cerrado.
- **En la pantalla** de sprints de HU-08 (`/web/proyectos/{id}/sprints`): cada sprint `Abierto` tiene un botón "Cerrar sprint" (`POST /web/proyectos/{id}/sprints/{sprintID}/cerrar`), que usa el mismo caso de uso que el endpoint JSON. Al cerrar, se vuelve a la lista con el sprint en `Cerrado` y sin botón (POST-redirect-GET, `303`); si hay un error, se muestra su mensaje (sección 7).

## 4. Reglas de negocio
- **RN-11.1:** Cerrar un sprint es pasarlo de `Abierto` a `Cerrado` (RN-08.8). Solo se cierra un sprint `Abierto`: cerrar uno ya `Cerrado` es un error, no una operación sin efecto (a diferencia de RN-10.3).
- **RN-11.2:** Un sprint se puede cerrar en cualquier momento: antes, durante o después de sus fechas, y aunque tenga historias sin terminar o no tenga ninguna. Tampoco importa el estado del proyecto: se puede cerrar aunque no esté `Activo` (P-10).
- **RN-11.3:** Al cerrar, las historias sin terminar del sprint (`Pendiente` o `En progreso`) vuelven al Product Backlog: quedan sin sprint, con su estado, prioridad y story points sin cambios, y se pueden planificar de nuevo en otro sprint por HU-09. Las historias `Hecho` siguen asignadas al sprint cerrado.
- **RN-11.4:** El sprint tiene que pertenecer al proyecto de la ruta; si es de otro proyecto, se trata como inexistente (igual que RN-18.1, RN-09.7 y RN-10.4).
- **RN-11.5:** Un sprint cerrado no se reabre.
- **RN-11.6:** Al cerrar, el sprint registra sus **story points planificados** (la suma de los de todas las historias asignadas en el momento del cierre, terminadas o no) y sus **story points completados** (la suma de los de las historias `Hecho`). Una historia sin estimar suma 0. Los valores quedan fijos: no cambian si después las historias se reasignan o se editan. Sirven de base a HU-22 (velocidad); mientras el sprint está `Abierto` valen 0.

## 5. Restricciones
- RN-11.1 vive en el dominio (`internal/domain/sprint.go`), sin dependencias externas: el método `Cerrar()` de `Sprint` hace el cambio de estado, y `PuedeCerrarse()` dice si un sprint se puede cerrar. La pantalla usa `PuedeCerrarse()` para decidir si muestra el botón "Cerrar sprint" (CA-11.5): la regla de visibilidad no se repite en el adaptador ni en la plantilla. RN-11.3 y RN-11.6 también viven en el dominio: la función `ResumirCierre(historias)` recibe las historias asignadas al sprint y devuelve los puntos planificados, los completados y qué historias vuelven al Product Backlog.
- La existencia del proyecto y del sprint (RN-11.4) se resuelve en el caso de uso `CerrarSprint` (`internal/app`). Para el sprint usa un puerto angosto propio de HU-11, `RepositorioCierre`, con `BuscarSprintPorID` y `CerrarSprint` (segregación de interfaces, como `RepositorioAsignacion` en HU-09). No se ensancha `RepositorioSprints` de HU-08, porque lo implementan el repositorio SQLite y los fakes y sistemas BDD de HU-08, y sumarle métodos obligaría a tocar esa historia. El repositorio SQLite de sprints implementa `RepositorioCierre` sin cambiar lo existente. Los nombres no chocan con el `BuscarSprint` del puerto de HU-09. `CerrarSprint` persiste siempre el estado `Cerrado` (RN-11.1) y del sprint usa solo el ID; además aplica RN-11.3 y RN-11.6 con `ResumirCierre` y deja en el sprint recibido los puntos registrados.
- **Persistencia del resumen:** la tabla `sprints` suma las columnas `puntos_planificados` y `puntos_completados` (enteros, 0 por omisión). Como `CREATE TABLE IF NOT EXISTS` no modifica una tabla existente, una base creada antes de HU-11 las recibe con `ALTER TABLE` al arrancar, igual que `sprint_id` en HU-09.
- **Dependencias:** el comportamiento de cerrar solo depende en código del sprint de HU-08; no depende de HU-09 ni de HU-10. El plan del Sprint 2 la pone después de HU-10 por orden de trabajo, no porque use su código. Lo que sí necesita HU-09 y HU-10 (ya mergeadas en `main`) es RN-11.3 y RN-11.6 (CA-11.6 y CA-11.7): las historias se asignan a un sprint con el `SprintID` de HU-09 y se terminan con el estado `Hecho` de HU-10. El efecto del cierre sobre las historias (no asignar, quitar ni completar en un sprint cerrado) lo controlan HU-09 y HU-10 leyendo el estado del sprint.
- **Pedidos simultáneos:** el cambio de estado se verifica en el mismo `UPDATE` que lo guarda (`WHERE estado = 'Abierto'`), como en HU-08. Si dos pedidos cierran el mismo sprint a la vez, uno lo cierra y el otro recibe `ErrSprintYaCerrado`. El cambio de estado, la lectura de las historias, el resumen y la devolución al backlog van en una sola transacción, que empieza por ese `UPDATE`: desde ahí SQLite no deja escribir a otro pedido hasta terminar, así que una asignación (HU-09) o un completado (HU-10) simultáneos quedan enteros antes o después del cierre, nunca en el medio. Si algo falla, no se guarda nada.
- **Nombres de errores:** `ErrSprintCerrado` y `ErrSprintNoEncontrado` los define HU-09, ya mergeada en `main`; los dos viven en el dominio (`domain.ErrSprintNoEncontrado`, en `internal/domain/sprint.go`). HU-11 no los redefine: usa `ErrSprintYaCerrado`, nuevo, y para el sprint inexistente reutiliza `domain.ErrSprintNoEncontrado`.
- Persistencia en SQLite (ADR-0002), `net/http` (ADR-0001) y pantallas con `html/template` (ADR-0003). La ruta JSON se registra con su propia función `RegistrarCerrarSprint` y se suma a `NuevoRouter` y a `router_test.go`; la de pantalla se suma a `RegistrarPantallas`, como las de las demás historias, y pasa por el chequeo de mismo origen de los formularios.
- Fuera de alcance: reabrir un sprint, mostrar el resumen en pantalla, velocidad y demás métricas (Épica 7, HU-22), y controlar quién puede cerrar (el sistema no tiene roles todavía).

## 6. Casos límite
- Sprint `Abierto` cuya fecha de fin todavía no llegó → se cierra (RN-11.2).
- Sprint `Abierto` sin historias asignadas → se cierra.
- Sprint `Abierto` sin historias asignadas → se cierra con 0 puntos planificados y 0 completados (RN-11.6).
- Sprint con historias sin terminar → se cierra y esas historias quedan sin sprint (RN-11.3); sus puntos cuentan como planificados (RN-11.6).
- Sprint con todas sus historias `Hecho` → se cierra, ninguna vuelve al backlog, y planificados = completados.
- Historia asignada sin estimar → suma 0 en planificados y en completados; si no está `Hecho`, vuelve al backlog igual.
- Sprint ya `Cerrado` → `ErrSprintYaCerrado`, no cambia nada.
- Sprint que existe pero en otro proyecto → sprint no encontrado.
- Dos pedidos simultáneos sobre el mismo sprint → uno lo cierra y el otro recibe `ErrSprintYaCerrado` (sección 5).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el sprint, o es de otro proyecto | `domain.ErrSprintNoEncontrado` (HU-09, se reutiliza) | "sprint no encontrado" (HTTP 404) |
| El sprint ya está Cerrado | `domain.ErrSprintYaCerrado` (nuevo) | "el sprint ya fue cerrado" (HTTP 400) |

Si hay varios errores, se informa el primero en el orden de la tabla. Tanto el endpoint JSON como la pantalla muestran el mensaje del error del dominio de esta tabla, nunca el texto envuelto por el caso de uso: si dos pedidos cierran a la vez, `ErrSprintYaCerrado` llega envuelto ("cerrar sprint: el sprint ya fue cerrado") y se responde igual el `400` con "el sprint ya fue cerrado". En la pantalla, el `400` vuelve a mostrar la lista de sprints con el mensaje, y los `404` muestran la página de error. Cualquier otro error es una falla interna: `500` con un mensaje genérico, sin el detalle. El mensaje de `ErrSprintYaCerrado` es distinto del de `ErrSprintCerrado` de HU-09 ("el sprint ya está cerrado"), que se usa cuando se intenta operar con historias sobre un sprint cerrado.

Los errores `ErrSprintYaCerrado` (HU-11) y `ErrSprintCerrado` (HU-09) describen la misma situación; se mantienen separados por el orden de los merges y se pueden unificar después de HU-09.

## 8. Criterios de aceptación
- **CA-11.1:** Dado un sprint `Abierto` con historias sin terminar y con su fecha de fin todavía sin llegar, cuando lo cierro, entonces queda `Cerrado`, con sus demás datos sin cambios. (RN-11.1, RN-11.2)
- **CA-11.2:** Dado un sprint ya `Cerrado`, cuando intento cerrarlo, entonces obtengo `ErrSprintYaCerrado` y sigue `Cerrado`. (RN-11.1, RN-11.5)
- **CA-11.3:** Dado un ID de proyecto o de sprint que no existe, o un sprint de otro proyecto, cuando lo cierro, entonces obtengo `ErrProyectoNoEncontrado` o `ErrSprintNoEncontrado`, y no cambia nada. (RN-11.4)
- **CA-11.4:** Dado un pedido válido a `POST /proyectos/{id}/sprints/{sprintID}/cerrar`, entonces la respuesta es `204`; con proyecto o sprint inexistente, `404`; con el sprint ya cerrado, `400`. (Sección 3)
- **CA-11.5:** Dada la pantalla de sprints de un proyecto, entonces solo los sprints `Abierto` muestran el botón "Cerrar sprint", y al usarlo el sprint aparece `Cerrado` y sin botón. (Sección 3)
- **CA-11.6:** Dado un sprint `Abierto` con historias `Hecho` y sin terminar, cuando lo cierro, entonces las sin terminar quedan sin sprint (de vuelta en el Product Backlog, con su estado sin cambios) y las `Hecho` siguen asignadas al sprint. (RN-11.3)
- **CA-11.7:** Dado un sprint `Abierto` con historias de 5 (`Hecho`), 3 (`Pendiente`) y sin estimar (`Hecho`), cuando lo cierro, entonces registra 8 story points planificados y 5 completados, y esos valores no cambian si después la historia de 3 se asigna a otro sprint. (RN-11.6)

CA-11.6 y CA-11.7 se verifican con escenarios sobre SQLite real (las historias se asignan con HU-09 y se completan con HU-10), con los tests del dominio de `ResumirCierre` y con los del repositorio SQLite.

## Decisiones tomadas
P-1 a P-10 las tomó el autor (Lucas Orellana) el 2026-10-27. Las confirmó Matías López (Agile Enabler) el 2026-10-29, que cambió P-4 y agregó P-11 a P-13.
- **P-1:** cerrar = pasar de `Abierto` a `Cerrado` (RN-11.1).
- **P-2:** cerrar un sprint ya `Cerrado` es un error, no una operación sin efecto (RN-11.1).
- **P-3:** se puede cerrar en cualquier momento, aunque no haya llegado la fecha de fin y aunque queden historias sin terminar (RN-11.2).
- **P-4:** las historias sin terminar vuelven al Product Backlog, sin sprint, y se planifican de nuevo con HU-09 (RN-11.3). Es lo que indica la Guía Scrum para el trabajo no terminado. Se descartaron dos alternativas: dejarlas en el sprint cerrado (la propuesta original), porque quedaban asignadas a un sprint que ya terminó, y moverlas solas al sprint siguiente, porque ese sprint puede no existir todavía y además se saltearía la Sprint Planning.
- **P-5:** un sprint cerrado no se reabre (RN-11.5).
- **P-6:** un sprint de otro proyecto se trata como inexistente (RN-11.4).
- **P-7:** ruta JSON `POST /proyectos/{id}/sprints/{sprintID}/cerrar` (`204`) y botón "Cerrar sprint" en la lista de sprints de HU-08, para no depender del detalle de HU-09.
- **P-8:** pedidos simultáneos: el estado se verifica en el mismo `UPDATE` (sección 5).
- **P-9:** fuera de alcance: mostrar el resumen en pantalla, velocidad y métricas, y control de quién cierra (sección 5).
- **P-10:** se puede cerrar un sprint de un proyecto que no está Activo, porque cerrar no suma trabajo nuevo (RN-08.2 solo limita crear).
- **P-11:** al cerrar, el sprint guarda sus puntos planificados y completados (RN-11.6). Como por P-4 las historias sin terminar salen del sprint, después ya no se podría saber cuántos puntos se habían planificado, y HU-22 calcularía mal la velocidad.
- **P-12:** las historias sin estimar suman 0 (RN-11.6). No hay otro valor posible sin inventar una estimación.
- **P-13:** las historias `Hecho` siguen en el sprint cerrado (RN-11.3): son el trabajo que se hizo en ese sprint.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-10-27: creada (borrador), con las decisiones P-1 a P-9 pendientes de confirmar.
- 2026-10-27: ajustes: puerto angosto, CA-11.6.
- 2026-10-27: decisiones P-1 a P-10 tomadas por el autor, abiertas a objeción en el PR.
- 2026-10-29: actualizada al estado del código tras la revisión: `ErrSprintNoEncontrado` vive en el dominio y HU-09 ya está en `main`; CA-11.6 se verifica con el escenario `@CA-11.6` sobre SQLite; se documentan `Sprint.PuedeCerrarse()`, el contrato de `CerrarSprint` del puerto y los mensajes de error (mensaje del dominio aunque llegue envuelto). Sigue en borrador hasta que la apruebe el equipo.
- 2026-10-29: Matías López confirma P-1 a P-10, salvo P-4: las historias sin terminar vuelven al Product Backlog (RN-11.3, CA-11.6). Se agregan RN-11.6, CA-11.7 y P-11 a P-13 (resumen de puntos). Spec aprobada.
