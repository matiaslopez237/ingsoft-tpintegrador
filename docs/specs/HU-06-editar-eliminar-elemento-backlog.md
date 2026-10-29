# HU-06: Editar o eliminar un elemento del backlog

**Historia:** Como Product Builder, quiero editar o eliminar un elemento del backlog, para mantenerlo actualizado.
**Issue:** #17  ·  **Sprint:** 2  ·  **Story points:** 3
**Estado de la spec:** Aprobada (PR #182)

## 1. Objetivo
Mantener al día el Product Backlog de un proyecto: corregir o completar los datos de un elemento creado en HU-05 (por ejemplo, estimarlo o cambiarle la prioridad) y sacar los que ya no hacen falta, sin perder información que usan otras historias (el esfuerzo de HU-18).

## 2. Entradas

**Editar** (`PATCH /proyectos/{id}/backlog/{elementoID}`):

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID`, `elementoID` | `int64` | Sí | El elemento existe y pertenece a ese proyecto. |
| `titulo` | `*string` | No | Si se envía: reglas de HU-05 (RN-05.2, RN-05.3). |
| `descripcion` | `*string` | No | Si se envía: no vacía (RN-05.4). |
| `prioridad` | `*domain.Prioridad` | No | Si se envía: `Alta`, `Media` o `Baja` (RN-05.5). |
| `story_points` | opcional y anulable | No | Si se envía un número: entero positivo (RN-05.6). Si se envía `null`: el elemento queda sin estimar. |
| `criterios` | `*[]string` | No | Si se envía: reemplaza la lista; al menos uno y ninguno vacío (RN-05.7). |
| `estado` | `*domain.EstadoElemento` | No | Si se envía: `Pendiente`, `En progreso` o `Hecho`. |

Se distingue "no enviado" (no cambia) de "enviado". En `story_points`, enviar `null` quita la estimación. En los demás campos, `null` se toma como "no enviado" (igual que en HU-02).

**Eliminar** (`DELETE /proyectos/{id}/backlog/{elementoID}`): `proyectoID` y `elementoID`.

## 3. Salidas esperadas
- **Editar:** el elemento queda guardado con los campos enviados cambiados y el resto igual. Se devuelve el elemento completo como quedó guardado. El `ID` y el proyecto no cambian.
- **Eliminar:** el elemento y sus criterios dejan de existir (HTTP 204). Buscarlo por ID da `ErrElementoNoEncontrado`.
- Si hay un error, no se modifica ni se elimina nada.

## 4. Reglas de negocio
- **RN-06.1:** La edición es parcial: solo cambian los campos enviados.
- **RN-06.2:** Los campos editados cumplen las mismas reglas que al crear (RN-05.2 a RN-05.7). El título sigue siendo único dentro del proyecto sin distinguir mayúsculas; renombrar un elemento a su mismo título con otras mayúsculas está permitido.
- **RN-06.3:** Se puede quitar la estimación enviando los story points en `null`.
- **RN-06.4:** El estado se puede cambiar libremente entre `Pendiente`, `En progreso` y `Hecho`, desde cualquiera de ellos.
- **RN-06.5:** Al editar debe enviarse al menos un campo. Enviar un campo con el valor que ya tiene no es un error.
- **RN-06.6:** El elemento tiene que pertenecer al proyecto indicado; un elemento de otro proyecto se trata como inexistente (igual que RN-18.1).
- **RN-06.7:** No se puede eliminar un elemento que tiene esfuerzo registrado (HU-18): esas horas se conservan para las métricas.
- **RN-06.8:** No se puede eliminar un elemento asignado a un Sprint activo. "Activo" es un Sprint en estado `Abierto`, el mismo criterio que HU-09 (RN-09.2). Un elemento asignado a un Sprint `Cerrado`, o sin Sprint, se puede eliminar (si cumple RN-06.7).
- **RN-06.9:** El `ID` del elemento no cambia.

## 5. Restricciones
- La validación de la edición vive en el dominio (`ElementoBacklog.Editar`) y reutiliza `NuevoElementoBacklog`, sin duplicar reglas.
- RN-06.8 vive en el dominio (`ElementoBacklog.ValidarEliminacion`), que recibe como `bool` si su Sprint está activo, igual que las reglas de HU-09; el caso de uso lo consulta con `BuscarSprint` (puerto de HU-09).
- La unicidad del título y la verificación de esfuerzo se resuelven en los casos de uso (`EditarElementoBacklog`, `EliminarElementoBacklog`) a través de puertos. La base también las garantiza: `UNIQUE (proyecto_id, titulo_clave)` y la clave foránea de `registros_esfuerzo`.
- **Pedidos simultáneos:** al guardar se actualizan solo los campos enviados (lección de DEF-01), así dos ediciones de campos distintos no se pisan. Los criterios, si se envían, se reemplazan en una transacción.
- Persistencia en SQLite (ADR-0002), API con `net/http` (ADR-0001) y pantalla con `html/template` (ADR-0003). La pantalla también envía solo los campos que el usuario cambió.

## 6. Casos límite
- Enviar solo `prioridad`: el resto no cambia.
- `story_points: null` → sin estimar; `story_points: 0` → error.
- Renombrar "Login" a "login" (el mismo elemento) → válido. Renombrar a "Registro" si ya existe "registro" en el proyecto → error de duplicado. Si "Registro" existe en otro proyecto → válido.
- Enviar `criterios: []` → error de al menos un criterio.
- Pasar de `Hecho` a `Pendiente` → válido. Estado `"hecho"` en minúscula o `"Cerrado"` → error.
- Pedido sin campos → error.
- Elemento de otro proyecto (por ruta) → no encontrado.
- Eliminar un elemento con un registro de esfuerzo → error; sin esfuerzo → se elimina con sus criterios.
- Eliminar dos veces el mismo elemento → la segunda da no encontrado.
- Eliminar un elemento asignado a un Sprint `Abierto` → error; asignado a un Sprint `Cerrado` → se elimina.

## 7. Condiciones de error
| Situación | Error (Go) | HTTP |
|---|---|---|
| El elemento no existe o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05) | 404 |
| No se envió ningún campo | `domain.ErrSinCambios` (HU-02, se reutiliza) | 400 |
| Título, descripción, prioridad, story points o criterios inválidos | los de HU-05 (`ErrTituloVacio`, `ErrTituloDemasiadoLargo`, `ErrDescripcionVacia`, `ErrPrioridadInvalida`, `ErrStoryPointsInvalidos`, `ErrSinCriterios`, `ErrCriterioVacio`) | 400 |
| Estado que no es Pendiente, En progreso ni Hecho | `domain.ErrEstadoElementoInvalido` (el mismo de HU-07) | 400 |
| Otro elemento del proyecto ya tiene ese título | `domain.ErrTituloDuplicado` (HU-05) | 400 |
| Eliminar un elemento asignado a un Sprint activo | `ErrElementoEnSprintActivo` | 409 |
| Eliminar un elemento con esfuerzo registrado | `ErrElementoConEsfuerzo` | 409 |
| Cuerpo JSON mal formado, con campos desconocidos o `story_points` que no es un entero | — (adaptador) | 400. Es un error de formato: se informa antes de buscar el elemento, como las fechas mal formadas en HU-30. |

Si hay varios errores, se informa el primero en el orden de la tabla: primero se busca el elemento, después el dominio valida los campos (en el orden de HU-05 y al final el estado) y por último se verifica el título duplicado. Al eliminar: primero se busca el elemento, después se verifica el Sprint activo y por último el esfuerzo.

## 8. Criterios de aceptación
- **CA-06.1:** Dado un elemento, cuando envío solo una nueva prioridad, entonces cambia la prioridad y el resto queda igual. (RN-06.1)
- **CA-06.2:** Dado un elemento, cuando edito título, descripción, story points, criterios y estado juntos, entonces se guardan todos y al recuperarlo tiene los valores nuevos y el mismo ID. (RN-06.1, RN-06.4, RN-06.9)
- **CA-06.3:** Dado un elemento estimado, cuando envío `story_points` en `null`, entonces queda sin estimar. (RN-06.3)
- **CA-06.4:** Cuando envío datos inválidos (título vacío, story points 0, criterios vacíos o un estado inexistente), entonces obtengo el error correspondiente y el elemento no cambia. (RN-06.2)
- **CA-06.5:** Dados "Login" y "Registro" en un proyecto, cuando renombro "Login" a "registro", entonces obtengo `ErrTituloDuplicado`; cuando lo renombro a "LOGIN", se guarda. (RN-06.2)
- **CA-06.6:** Cuando no envío ningún campo, entonces obtengo `ErrSinCambios`. (RN-06.5)
- **CA-06.7:** Dado un elemento de otro proyecto o un ID inexistente, cuando lo edito o lo elimino, entonces obtengo `ErrElementoNoEncontrado`. (RN-06.6)
- **CA-06.8:** Dado un elemento sin esfuerzo, cuando lo elimino, entonces deja de existir junto con sus criterios. (RN-06.7)
- **CA-06.9:** Dado un elemento con esfuerzo registrado, cuando lo elimino, entonces obtengo `ErrElementoConEsfuerzo` y el elemento sigue existiendo. (RN-06.7)
- **CA-06.10:** Desde la pantalla del elemento puedo editarlo (se guardan solo los campos que cambié) y eliminarlo, y los errores se muestran en la misma pantalla. Si guardo sin cambiar nada, vuelvo al elemento sin mensaje de error, igual que en el formulario de proyecto (DEF-01). Desde la lista del backlog (HU-07) cada elemento tiene "Editar" y "Eliminar" (con confirmación), y al eliminar vuelvo a esa lista. (Restricciones)
- **CA-06.11:** Dado un elemento asignado a un Sprint Abierto, cuando lo elimino, entonces obtengo `ErrElementoEnSprintActivo` y el elemento sigue existiendo; si su Sprint está Cerrado, se elimina. (RN-06.8)

## Decisiones tomadas
Definidas por Nicolás Andino al refinar la historia (2026-10-26):
- No se elimina un elemento con esfuerzo registrado (RN-06.7).
- No se elimina un elemento en un Sprint activo (RN-06.8, criterio del issue #17). Se implementó después de HU-09, que define cuándo un Sprint está activo.
- El estado se edita libremente (RN-06.4).
- Edición parcial, como HU-02 (RN-06.1).

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-26: creada (borrador).
- 2026-10-26: aclarados, a partir de la revisión con go-reviewer, el `null` en campos que no son `story_points` y que el formato inválido de `story_points` se informa antes de buscar el elemento. La pantalla pide confirmación antes de eliminar.
- 2026-10-26: guardar la pantalla sin cambios vuelve al elemento sin error, alineado con DEF-01 (la API sigue devolviendo `ErrSinCambios`).
- 2026-10-26: integrada con la lista del backlog de HU-07: botones Editar y Eliminar en cada elemento, y al eliminar se vuelve a la lista. `ErrEstadoElementoInvalido` es el mismo de HU-07.
- 2026-10-26: a partir de la revisión del PR #182, un POST a la pantalla sin los campos ocultos `original_*` envía todos los campos (como DEF-01), en lugar de compararlos contra valores vacíos.
- 2026-10-28: implementada RN-06.8 (no eliminar un elemento en un Sprint activo) ahora que HU-09 está en `main`; nuevo error `ErrElementoEnSprintActivo` (409) y CA-06.11. La spec pasa a Aprobada (PR #182 mergeado).
