# HU-12: Consultar Sprints anteriores

**Historia:** Como Scrum Master, quiero consultar Sprints anteriores de un proyecto, para revisar el historial.
**Issue:** #23  ·  **Sprint:** 3  ·  **Story points:** 2
**Estado de la spec:** Borrador

## 1. Objetivo
Ver el historial de los sprints ya cerrados de un proyecto: qué se propuso cada uno (Sprint Goal), cuándo fue y cuántos story points planificó y completó. Sirve para revisar cómo viene el equipo antes de planificar el próximo sprint. Usa los sprints de HU-08 y el resumen de puntos que guarda HU-11 al cerrar un sprint (RN-11.6).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |

No hay filtros ni paginación. **En HTTP:** `GET /proyectos/{id}/sprints/anteriores`, sin cuerpo.

## 3. Salidas esperadas
- La lista de los sprints **cerrados** del proyecto (RN-12.1), del más reciente al más antiguo (RN-12.3). De cada uno: ID, número, fecha de inicio, fecha de fin, Sprint Goal, story points planificados y story points completados (RN-12.2).
- Si el proyecto no tiene sprints cerrados, la lista vacía. No es un error.
- La consulta no modifica nada.
- **En HTTP:** `200` con un arreglo JSON (`[]` si no hay ninguno). Cada elemento tiene `id`, `numero`, `fecha_inicio`, `fecha_fin` (AAAA-MM-DD), `objetivo`, `puntos_planificados` y `puntos_completados`. `404` si no existe el proyecto; `400` si el ID de la ruta no es un número.
- **Pantalla** `GET /web/proyectos/{id}/sprints/anteriores`, "Sprints anteriores": una tabla con los mismos datos y, en cada fila, un enlace al detalle del sprint (HU-09), que muestra sus historias terminadas. Si no hay sprints cerrados, un mensaje que lo dice. Se llega desde el detalle del proyecto y desde la lista de sprints (HU-08).

## 4. Reglas de negocio
- **RN-12.1:** Un sprint anterior es un sprint del proyecto en estado `Cerrado` (HU-11). Un sprint `Abierto` no aparece, aunque su fecha de fin ya haya pasado: hasta que se cierra, no tiene resumen de puntos.
- **RN-12.2:** De cada sprint se muestran el número, las fechas, el Sprint Goal y los story points planificados y completados que registró el cierre (RN-11.6). Se muestran tal como se guardaron: no se recalculan desde las historias.
- **RN-12.3:** Se ordenan por fecha de inicio, del más reciente al más antiguo. Se usa la fecha y no el número porque el número indica el orden en que se creó el sprint, no cuándo transcurrió. Como los sprints de un proyecto no se superponen (RN-08.5), no hay empates.

## 5. Restricciones
- RN-12.1 y RN-12.3 viven en el dominio (`internal/domain/sprint.go`), en una función pura `SprintsAnteriores(sprints []Sprint) []Sprint`, sin dependencias externas.
- El caso de uso `ListarSprintsAnteriores` (`internal/app`) verifica el proyecto con `LectorProyectos` y lee sus sprints con el `ListarSprints` del puerto `RepositorioSprints` de HU-08. No hace falta persistencia nueva: `ListarSprints` ya lee los puntos guardados por HU-11.
- `net/http` (ADR-0001). La ruta JSON se registra con su propia función `RegistrarSprintsAnteriores` y se suma a `NuevoRouter`, a `CasosDeUso` y a `router_test.go`. La pantalla, con `html/template` (ADR-0003), se suma a `RegistrarPantallas`. Las dos rutas terminan en el literal `anteriores`, que en el enrutador de Go tiene prioridad sobre `{sprintID}` de HU-09: no chocan.
- Es una consulta: no modifica nada, así que no hay pedidos simultáneos que resolver.
- Fuera de alcance: la velocidad del equipo y los gráficos (HU-22, Épica 7 y 8), filtrar o paginar el historial, y los sprints abiertos vencidos.

## 6. Casos límite
- Proyecto sin sprints → lista vacía.
- Proyecto con sprints, todos `Abierto` → lista vacía.
- Sprint `Abierto` con la fecha de fin ya pasada → no aparece (RN-12.1).
- Sprint cerrado sin historias → aparece con 0 puntos planificados y 0 completados.
- Sprint cerrado antes de que existiera el resumen de HU-11 (base de datos vieja) → aparece con 0 y 0, que es lo que dejó la migración.
- Sprint 2 creado con fechas anteriores a las del sprint 1 → el sprint 2 aparece después del 1 (RN-12.3).
- Sprints cerrados de otro proyecto → no aparecen.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404; en la pantalla, la página de error) |
| El ID de la ruta no es un número | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso |

Cualquier otro error es una falla interna: `500` con un mensaje genérico, sin el detalle.

## 8. Criterios de aceptación
- **CA-12.1:** Dado un proyecto con el sprint 1 y el sprint 2 cerrados y el sprint 3 abierto, cuando consulto los sprints anteriores, entonces obtengo el 2 y el 1, en ese orden, cada uno con sus fechas, su Sprint Goal y sus puntos planificados y completados. (RN-12.1, RN-12.2, RN-12.3)
- **CA-12.2:** Dado un proyecto sin sprints, o con todos sus sprints abiertos (aunque alguno ya haya vencido), cuando consulto los sprints anteriores, entonces obtengo una lista vacía, sin error. (RN-12.1)
- **CA-12.3:** Dado un proyecto cuyo sprint 2 transcurrió antes que el sprint 1, los dos cerrados, cuando consulto los sprints anteriores, entonces el sprint 1 aparece primero. (RN-12.3)
- **CA-12.4:** Dado un proyecto con un sprint cerrado y otro proyecto con sprints cerrados, cuando consulto los anteriores del primero, entonces solo veo el suyo. (RN-12.1)
- **CA-12.5:** Dado un ID de proyecto que no existe, cuando consulto los sprints anteriores, entonces obtengo `ErrProyectoNoEncontrado`. (Sección 7)
- **CA-12.6:** Dado un pedido a `GET /proyectos/{id}/sprints/anteriores`, entonces la respuesta es `200` con el arreglo JSON (`[]` si no hay sprints cerrados); con un proyecto inexistente, `404`; con un ID que no es un número, `400`. (Sección 3)
- **CA-12.7:** Dada la pantalla "Sprints anteriores" de un proyecto, entonces muestra una fila por sprint cerrado, con sus datos y un enlace a su detalle, o un mensaje si no hay ninguno. Se llega a ella desde el detalle del proyecto y desde la lista de sprints. (Sección 3)

## Decisiones tomadas
Confirmadas por Matías López (Agile Enabler) el 2026-10-30:
- **P-1:** un sprint anterior es uno `Cerrado`; los abiertos vencidos no aparecen (RN-12.1).
- **P-2:** se muestran el Sprint Goal, las fechas y los puntos planificados y completados de HU-11, con un enlace al detalle del sprint (RN-12.2).
- **P-3:** del más reciente al más antiguo (RN-12.3). "Más reciente" se mide por la fecha de inicio, no por el número.
- **P-4:** página propia "Sprints anteriores", sin tocar la pantalla de sprints de HU-08 (salvo por el enlace).

## Preguntas abiertas
Ninguna.

## Historial
- 2026-10-30: creada, con P-1 a P-4 confirmadas.
