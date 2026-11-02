# HU-21: Actualizar el estado de un defecto

**Historia:** Como Product Builder, quiero actualizar el estado de un defecto y registrar su Sprint de resolución, para reflejar su progreso.
**Issue:** a confirmar  ·  **Sprint:** 3  ·  **Story points:** 2
**Estado de la spec:** Aprobada

## 1. Objetivo
Dar seguimiento a un defecto ya registrado (HU-20) a medida que se corrige: marcarlo `Resuelto` cuando se corrigió y `Cerrado` cuando se verificó y se da por terminado. Al resolverlo se asigna automáticamente el Sprint en el que se resolvió, igual que HU-20 asigna el Sprint de detección. Como todavía no existe una pantalla que liste los defectos de un proyecto, esta historia también la agrega (tarea 21.8): sin ella no habría forma de elegir qué defecto actualizar.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `defectoID` | `int64` | Sí | Debe corresponder a un defecto (HU-20) de ese proyecto. |

No hay cuerpo: son dos acciones concretas, igual que "cerrar sprint" (HU-11) y "completar historia" (HU-10), no una actualización genérica a un estado arbitrario.

**En HTTP:**
- `POST /proyectos/{id}/defectos/{defectoID}/resolver`
- `POST /proyectos/{id}/defectos/{defectoID}/cerrar`

## 3. Salidas esperadas
- **Resolver:** el defecto pasa de `Abierto` a `Resuelto` y se le asigna el `SprintResolucionID` (el Sprint `Abierto` vigente hoy, o `nil` si ninguno califica). El resto de sus datos no cambia.
- **Cerrar:** el defecto pasa de `Resuelto` a `Cerrado`. El `SprintResolucionID` no se toca (ya quedó fijado al resolver).
- Si hay un error, no se modifica nada.
- **En HTTP:** `204` sin cuerpo en ambas acciones; `404` si no existe el proyecto o el defecto; `400` si el defecto no está en el estado de origen requerido.
- **Listado de apoyo** (necesario porque HU-20 no expone ninguno): `GET /proyectos/{id}/defectos` → `200` con los defectos del proyecto en JSON.
- **Pantalla** (`/web/proyectos/{id}/defectos`, tarea 21.8): lista los defectos del proyecto con su elemento relacionado, descripción, severidad, estado y Sprints de detección/resolución (o "ninguno"). Cada defecto `Abierto` muestra el botón "Marcar resuelto"; cada `Resuelto` muestra "Cerrar"; los `Cerrado` no muestran botón. Al usar un botón, la lista se vuelve a mostrar con el estado actualizado (POST-redirect-GET, `303`, mismo patrón que HU-10/HU-11); si hay un error se muestra su mensaje.

## 4. Reglas de negocio
- **RN-21.1:** Un defecto tiene tres estados posibles: `Abierto` (HU-20), `Resuelto` y `Cerrado`.
- **RN-21.2:** Resolver un defecto es pasarlo de `Abierto` a `Resuelto`. Solo se puede resolver un defecto `Abierto`; intentarlo sobre uno `Resuelto` o `Cerrado` es un error, no una operación sin efecto (mismo criterio que RN-11.1: cerrar un sprint ya cerrado).
- **RN-21.3:** Cerrar un defecto es pasarlo de `Resuelto` a `Cerrado`. Solo se puede cerrar un defecto `Resuelto`; intentarlo sobre uno `Abierto` (saltearía `Resuelto`) o ya `Cerrado` es un error.
- **RN-21.4:** No se puede retroceder de estado (de `Cerrado` o `Resuelto` a uno anterior). Reabrir un defecto queda fuera de esta historia.
- **RN-21.5:** Al resolver, el Sprint de resolución se asigna automáticamente: es el Sprint del proyecto en estado `Abierto` cuyo rango de fechas (día completo) incluye la fecha de hoy, el mismo criterio que el Sprint de detección de HU-20 (RN-20.5). Si ninguno califica, `SprintResolucionID` queda `nil`. Nunca lo elige quien resuelve el defecto.
- **RN-21.6:** El defecto tiene que pertenecer al proyecto de la ruta; si es de otro proyecto, se trata como inexistente (mismo patrón que RN-20.1, RN-11.4, RN-10.4).

## 5. Restricciones
- RN-21.2 a RN-21.4 viven en el dominio (`internal/domain/defecto.go`), sin dependencias externas: `Defecto.PuedeResolverse()` / `Defecto.Resolver()` y `Defecto.PuedeCerrarse()` / `Defecto.Cerrar()`, mismo patrón que `Sprint.PuedeCerrarse()` / `Sprint.Cerrar()` de HU-11. El dominio no conoce los sprints del proyecto (igual que HU-20): `Resolver()` solo cambia el estado; asignar el `SprintResolucionID` es tarea del caso de uso, reutilizando `domain.SprintVigente` (ya existe, de HU-20).
- La existencia del proyecto y del defecto (RN-21.6), y la asignación del Sprint de resolución, se resuelven en un caso de uso nuevo `ActualizarEstadoDefecto` (`internal/app/actualizar_estado_defecto.go`) con dos métodos, `Resolver` y `Cerrar`, que comparten las mismas dependencias (puertos `LectorProyectos` de HU-04 y `RepositorioSprints` de HU-08, más un puerto propio para los defectos).
- Puerto angosto `RepositorioEstadoDefecto` (`BuscarDefectoPorID`, `ActualizarEstadoDefecto`), igual patrón de segregación de interfaces que `RepositorioCierre` en HU-11: no se ensancha el `RepositorioDefectos` de HU-20 (que solo tiene `GuardarDefecto`), para no tocar esa historia.
- El listado (`GET /proyectos/{id}/defectos`) tiene su propio caso de uso de solo lectura `ListarDefectos` (`internal/app/listar_defectos.go`), con el puerto `LectorDefectos.ListarDefectos(ctx, proyectoID)`. El dominio no sabe el título del elemento relacionado (solo guarda `ElementoID`); la pantalla lo resuelve con `ListarBacklog` (HU-07) y empareja por ID, igual que ya hace `pantallas_defecto.go` con el formulario de HU-20.
- **Persistencia:** la tabla `defectos` suma la columna `sprint_resolucion_id` (entero, `NULL` por omisión), agregada con `ALTER TABLE` al arrancar si falta (mismo mecanismo `agregarColumnas` que HU-09/HU-11). El repositorio SQLite de HU-20 (`defecto_repo.go`) implementa además `RepositorioEstadoDefecto` y `LectorDefectos`, sin cambiar lo que ya tiene.
- **Pedidos simultáneos:** el estado se verifica en el mismo `UPDATE` que lo guarda (`WHERE estado = 'Abierto'` o `'Resuelto'` según la acción), igual que HU-08/HU-11. Si dos pedidos actúan sobre el mismo defecto a la vez, uno tiene efecto y el otro recibe el error de estado correspondiente.
- Persistencia en SQLite (ADR-0002); `net/http` (ADR-0001); pantallas con `html/template` (ADR-0003).
- Fuera de alcance: reabrir un defecto, editar su descripción/severidad/elemento, eliminarlo, elegir manualmente el Sprint de resolución, filtros o búsqueda en la lista, y el endpoint JSON de un único defecto por ID (no hace falta: el caso de uso lo busca internamente).

## 6. Casos límite
- Defecto `Abierto` con un Sprint `Abierto` vigente hoy → se resuelve y el `SprintResolucionID` queda el de ese Sprint.
- Defecto `Abierto` sin ningún Sprint vigente hoy (sin sprints, o ninguno `Abierto` con fechas que incluyan hoy) → se resuelve con `SprintResolucionID` `nil`.
- Defecto `Resuelto` → se cierra; el `SprintResolucionID` no cambia.
- Defecto `Resuelto` o `Cerrado` → intentar resolverlo de nuevo es error (RN-21.2).
- Defecto `Abierto` → intentar cerrarlo directamente (saltear `Resuelto`) es error (RN-21.3).
- Defecto `Cerrado` → intentar cerrarlo de nuevo es error (RN-21.3).
- Defecto que existe pero en otro proyecto → se trata como inexistente (RN-21.6).
- Proyecto sin ningún defecto todavía → la lista se muestra vacía, sin error.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el defecto, o es de otro proyecto | `domain.ErrDefectoNoEncontrado` (nuevo) | "defecto no encontrado" (HTTP 404) |
| Resolver un defecto que no está `Abierto` | `domain.ErrDefectoNoAbierto` (nuevo) | "el defecto no está Abierto" (HTTP 400) |
| Cerrar un defecto que no está `Resuelto` | `domain.ErrDefectoNoResuelto` (nuevo) | "el defecto no está Resuelto" (HTTP 400) |

Si hay varios errores, se informa el primero en el orden de la tabla: primero se busca el proyecto, después el defecto, después se valida el estado de origen.

## 8. Criterios de aceptación
- **CA-21.1:** Dado un defecto `Abierto` y un Sprint `Abierto` cuyas fechas incluyen hoy, cuando lo resuelvo, entonces queda `Resuelto` con el `SprintResolucionID` de ese Sprint. (RN-21.2, RN-21.5)
- **CA-21.2:** Dado un defecto `Abierto` sin ningún Sprint vigente hoy, cuando lo resuelvo, entonces queda `Resuelto` con `SprintResolucionID` nulo. (RN-21.5)
- **CA-21.3:** Dado un defecto `Resuelto`, cuando lo cierro, entonces queda `Cerrado` y su `SprintResolucionID` no cambia. (RN-21.3)
- **CA-21.4:** Dado un defecto `Resuelto` o `Cerrado`, cuando intento resolverlo, entonces obtengo `ErrDefectoNoAbierto` y no cambia. (RN-21.2, RN-21.4)
- **CA-21.5:** Dado un defecto `Abierto` o `Cerrado`, cuando intento cerrarlo, entonces obtengo `ErrDefectoNoResuelto` y no cambia. (RN-21.3, RN-21.4)
- **CA-21.6:** Dado un ID de proyecto o de defecto que no existe, o un defecto de otro proyecto, cuando actualizo su estado, entonces obtengo `ErrProyectoNoEncontrado` o `ErrDefectoNoEncontrado`. (RN-21.6)
- **CA-21.7:** Dado un pedido válido a `POST /proyectos/{id}/defectos/{defectoID}/resolver` o `.../cerrar`, entonces la respuesta es `204`; con el defecto en un estado inválido para la acción, `400`; con proyecto o defecto inexistente, `404`.
- **CA-21.8:** Dado un pedido a `GET /proyectos/{id}/defectos`, entonces la respuesta es `200` con los defectos del proyecto en JSON (lista vacía si no tiene ninguno).
- **CA-21.9:** Dada la pantalla de defectos de un proyecto, entonces se listan todos con su elemento, descripción, severidad, estado y Sprints de detección/resolución; los `Abierto` muestran "Marcar resuelto", los `Resuelto` muestran "Cerrar", los `Cerrado` no muestran botón; al usar un botón la lista se actualiza. (Sección 3)

## Decisiones tomadas
Las acordó el equipo (Ezequiel) en el chat el 2026-11-01:
- **D-1 (estados):** tres estados — `Abierto`, `Resuelto`, `Cerrado` — sin un estado intermedio "En progreso". (RN-21.1)
- **D-2 (Sprint de resolución):** automático, mismo criterio que el Sprint de detección de HU-20. (RN-21.5)
- **D-3 (transiciones):** solo hacia adelante, sin saltos ni retrocesos: `Abierto`→`Resuelto`→`Cerrado`. Intentar cualquier otra transición es error, no una operación sin efecto (mismo criterio que RN-11.1, distinto del de RN-10.3). (RN-21.2 a RN-21.4)
- **D-4 (pantalla):** lista nueva en `/web/proyectos/{id}/defectos`, propia de esta historia.

Decisiones técnicas tomadas por el autor (Ezequiel), siguiendo los patrones ya usados en el proyecto, abiertas a objeción en la revisión del PR:
- **P-1 (rutas específicas, no un estado genérico):** dos rutas de acción (`.../resolver`, `.../cerrar`) sin cuerpo, igual que `.../cerrar` de HU-11 y `.../completar` de HU-10, en vez de un único endpoint que reciba el nuevo estado como dato. Evita inventar un mecanismo de transición genérico que no se usa en el resto del proyecto, y evita validar un string de estado arbitrario.
- **P-2 (sin reabrir):** fuera de alcance, no lo pide la historia (sección 5).
- **P-3 (puerto angosto):** `RepositorioEstadoDefecto` nuevo en vez de ensanchar `RepositorioDefectos` de HU-20 (sección 5), mismo criterio que `RepositorioCierre` de HU-11.
- **P-4 (listado):** se agrega `GET /proyectos/{id}/defectos` (JSON) además de la pantalla, porque la pantalla lo necesita igual y mantiene la paridad JSON/pantalla que tienen las demás historias con pantalla propia.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-01: creada (borrador), con las decisiones D-1 a D-4 confirmadas en el chat y P-1 a P-4 tomadas por el autor, abiertas a objeción en el PR.
- 2026-11-01: aprobada.
