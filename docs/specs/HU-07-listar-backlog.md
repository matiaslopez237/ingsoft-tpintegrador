# HU-07: Listar y ordenar el backlog por prioridad

**Historia:** Como Product Owner, quiero listar y ordenar el backlog por prioridad, para saber qué se debe entregar primero.
**Issue:** #18  ·  **Sprint:** 2  ·  **Story points:** 2
**Estado de la spec:** Borrador

## 1. Objetivo
Mostrar el Product Backlog de un proyecto (los elementos creados en HU-05) ordenado por prioridad, con la opción de invertir el orden y de ver solo los elementos de un estado. Incluye la pantalla del backlog, que también permite agregar elementos con el caso de uso de HU-05, que hasta ahora no tenía pantalla.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente. |
| `orden` | `string` | No | `desc` (por defecto) o `asc`. |
| `estado` | `string` | No | Si se envía: `Pendiente`, `En progreso` o `Hecho` (RN-05.8). Vacío: todos los estados. |

**En HTTP** (`GET /proyectos/{id}/backlog`), `orden` y `estado` van como parámetros de la URL: `GET /proyectos/1/backlog?orden=asc&estado=Pendiente`.

## 3. Salidas esperadas
- La lista de elementos del proyecto, completos (ID, título, descripción, prioridad, estado, story points y criterios), en el orden de RN-07.2. Si no hay elementos, o ninguno tiene el estado pedido, la lista está vacía.
- No se modifica ningún dato.

En HTTP se responde `200` con un array de elementos en el mismo formato que la respuesta de HU-05 (`[]` si no hay ninguno), `400` si `orden` o `estado` no son válidos y `404` si el proyecto no existe.

## 4. Reglas de negocio
- **RN-07.1:** Se listan los elementos de un proyecto existente, y solo los de ese proyecto.
- **RN-07.2:** El orden por defecto es por prioridad descendente: primero `Alta`, después `Media` y al final `Baja`, porque es "lo que se debe entregar primero". Con orden ascendente se invierte: primero `Baja` y al final `Alta`.
- **RN-07.3:** A igual prioridad, los elementos van en el orden en que se crearon, en los dos sentidos. Así el orden es estable: no cambia entre una consulta y otra si no se cargan elementos nuevos.
- **RN-07.4:** Por defecto se listan los elementos de todos los estados. Con el filtro de estado, solo los de ese estado.
- **RN-07.5:** Un proyecto sin elementos, o un filtro que no coincide con ninguno, devuelve una lista vacía, no un error.

## 5. Restricciones
- El orden (RN-07.2, RN-07.3) y el filtro (RN-07.4) son reglas de negocio: viven en el dominio (`internal/domain`) y se desarrollan con TDD. La persistencia solo devuelve los elementos del proyecto.
- La existencia del proyecto (RN-07.1) se resuelve en el caso de uso a través del puerto `LectorProyectos` (HU-04).
- **Pantalla** (`/web/proyectos/{id}/backlog`, ADR-0003): la lista con el orden y el filtro, y un formulario para agregar elementos con el caso de uso de HU-05 (sus reglas no cambian). En el formulario los criterios de aceptación se escriben uno por línea; las líneas vacías no son criterios y se ignoran. Al agregar un elemento, la pantalla conserva el orden y el filtro que estaban elegidos.
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001).

## 6. Casos límite
- Proyecto sin elementos → lista vacía.
- Todos los elementos con la misma prioridad → en orden de creación, en los dos sentidos.
- Filtro por un estado que ningún elemento tiene → lista vacía.
- Orden o estado escritos con otras mayúsculas (`ASC`, `pendiente`) → error: se aceptan exactamente los valores de la sección 2.
- Elemento sin story points → aparece igual, sin estimar.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `ErrProyectoNoEncontrado` (HU-01) | "proyecto no encontrado" (HTTP 404) |
| `orden` distinto de `asc` y `desc` | `ErrOrdenInvalido` | "el orden debe ser asc o desc" (HTTP 400) |
| `estado` que no es uno de los tres de RN-05.8 | `ErrEstadoElementoInvalido` | "el estado debe ser Pendiente, En progreso o Hecho" (HTTP 400) |

Si hay varios errores, se informa el primero en el orden de la tabla.

## 8. Criterios de aceptación
- **CA-07.1:** Dado un proyecto con elementos de prioridad Baja, Alta y Media, cuando listo su backlog, entonces aparecen en el orden Alta, Media, Baja. (RN-07.2)
- **CA-07.2:** Dado el mismo proyecto, cuando listo su backlog en orden ascendente, entonces aparecen en el orden Baja, Media, Alta. (RN-07.2)
- **CA-07.3:** Dados dos elementos con la misma prioridad, entonces aparecen en el orden en que se crearon, tanto en orden descendente como ascendente. (RN-07.3)
- **CA-07.4:** Dado un proyecto con elementos en distintos estados, cuando listo su backlog sin filtro, entonces aparecen todos; y cuando filtro por un estado, solo los de ese estado. (RN-07.4)
- **CA-07.5:** Dado un proyecto sin elementos, o un filtro que no coincide con ninguno, cuando listo su backlog, entonces obtengo una lista vacía. (RN-07.5)
- **CA-07.6:** Dado otro proyecto con elementos, cuando listo el backlog de un proyecto, entonces no aparecen los elementos del otro. (RN-07.1)
- **CA-07.7:** Dado un ID de proyecto que no existe, cuando listo su backlog, entonces obtengo `ErrProyectoNoEncontrado`. (RN-07.1)
- **CA-07.8:** Cuando pido un orden o un estado que no son válidos, entonces obtengo `ErrOrdenInvalido` o `ErrEstadoElementoInvalido`. (Sección 7)

## Decisiones tomadas
Definidas por Matías López (2026-10-26):
- A igual prioridad, orden de creación.
- Por defecto se listan todos los estados, incluido `Hecho`.
- La pantalla del backlog incluye el formulario para agregar elementos (HU-05 no tenía pantalla).

Tomadas de los criterios del issue #18: orden ascendente y descendente, y filtro por estado. El orden por defecto descendente (Alta primero) sale de la historia: "saber qué se debe entregar primero".

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-26: creada (borrador).
- 2026-10-26: al agregar un elemento, la pantalla conserva el orden y el filtro elegidos (sección 5).
