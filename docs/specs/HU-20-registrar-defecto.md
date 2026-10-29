# HU-20: Registrar un defecto

**Historia:** Como Product Builder, quiero registrar un defecto con descripción, severidad, estado, historia relacionada y Sprint de detección, para darle seguimiento.
**Issue:** #31  ·  **Sprint:** 2  ·  **Story points:** 3
**Estado de la spec:** Aprobada

## 1. Objetivo
Dejar constancia de un defecto encontrado sobre un elemento del backlog ya existente: qué pasa, qué tan grave es y en qué Sprint se detectó. Es el dato base de la Épica 6 (gestión de defectos), que después usan el seguimiento, las métricas y el dashboard. Reutiliza los elementos de backlog de HU-05 y los sprints de HU-08. En esta historia se exponen tanto el endpoint JSON como la pantalla (a diferencia de HU-18, que dejó la pantalla para HU-30: acá la pantalla es la tarea 20.8 de esta misma historia).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento de backlog (HU-05) del mismo proyecto. Es la "historia relacionada"; se acepta cualquier elemento del backlog, sin restricción de estado ni prioridad. |
| `descripcion` | `string` | Sí | Se recortan espacios. No vacía. Texto libre. |
| `severidad` | `string` | Sí | Una de: `Crítica`, `Alta`, `Media`, `Baja`. |

El **estado** y el **Sprint de detección** no son entradas: los asigna el sistema (ver sección 4).

**En HTTP** (`POST /proyectos/{id}/defectos`), el proyecto va en la ruta y el resto en un cuerpo JSON, con los nombres en snake_case, como en HU-18:

```json
{
  "elemento_id": 7,
  "descripcion": "No carga la pantalla de login en Firefox",
  "severidad": "Alta"
}
```

## 3. Salidas esperadas
- Se devuelve el defecto creado con `ID` (único, autogenerado), `ProyectoID`, `ElementoID`, `Descripcion` recortada, `Severidad`, `Estado` (siempre `Abierto`) y `SprintDeteccionID` (el `ID` del Sprint vigente hoy, o `nil` si no hay ninguno).
- El defecto queda persistido y se puede recuperar por su `ID`.
- Si hay un error, no se guarda nada.
- **En HTTP:** `201` con el defecto creado en JSON (`id`, `proyecto_id`, `elemento_id`, `descripcion`, `severidad`, `estado`, `sprint_deteccion_id` — este último `null` si no hay sprint vigente); `400` con el error si los datos del cuerpo son inválidos; `404` solo si no existe el proyecto o el elemento, que son los datos que identifican entidades existentes.
- **Pantalla** (`/web/proyectos/{id}/defectos/nuevo`, tarea 20.8): formulario con el elemento del backlog (elegido de una lista del proyecto), la descripción y la severidad; al guardar, confirma la creación y muestra el Sprint de detección asignado (o que no se asignó ninguno).

## 4. Reglas de negocio
- **RN-20.1:** Todo defecto pertenece a un elemento de backlog existente del proyecto indicado (la "historia relacionada"). Un elemento de otro proyecto se trata como inexistente.
- **RN-20.2:** La descripción es obligatoria (una de solo espacios es vacía) y es texto libre.
- **RN-20.3:** La severidad es obligatoria y tiene que ser una de: Crítica, Alta, Media o Baja.
- **RN-20.4:** El estado de un defecto nuevo siempre es `Abierto`; no lo elige quien lo registra. Cambiar el estado queda fuera de esta historia (Épica 6, historias futuras).
- **RN-20.5:** El Sprint de detección se asigna automáticamente: es el Sprint del proyecto en estado `Abierto` cuyo rango de fechas (día completo, igual que RN-08.5) incluye la fecha de hoy. Si ningún Sprint del proyecto cumple esa condición, el defecto queda sin Sprint de detección (`nil`). Nunca lo elige quien registra el defecto.
- **RN-20.6:** El `ID` es único y lo genera el sistema (entero autoincremental), igual que en HU-05/HU-18.

## 5. Restricciones
- Las validaciones de campos (RN-20.2, RN-20.3) viven en el dominio (`internal/domain/defecto.go`, constructor `NuevoDefecto`), sin dependencias externas.
- "Hoy" no se lee del reloj dentro del dominio: el caso de uso recibe un reloj (`func() time.Time`), igual que HU-18, para que los tests sean deterministas. La comparación de fechas usa el mismo criterio de "día calendario" que HU-08 (RN-08.5) y HU-18 (RN-18.4).
- La existencia del proyecto y del elemento (RN-20.1) se resuelve en el caso de uso `RegistrarDefecto` (`internal/app/registrar_defecto.go`) reutilizando los puertos `LectorProyectos` (HU-04) y `LectorElementos` (HU-18).
- La búsqueda del Sprint vigente (RN-20.5) reutiliza el puerto `RepositorioSprints.ListarSprints` (HU-08); como RN-08.5 impide que los sprints de un proyecto se superpongan en fechas, a lo sumo un Sprint `Abierto` puede incluir el día de hoy.
- Persistencia en SQLite (ADR-0002); `net/http` (ADR-0001), ruta propuesta `POST /proyectos/{id}/defectos` (no anidada bajo el elemento, porque el elemento se indica en el cuerpo como `elemento_id`, no en la URL — ver P-3).
- Pantalla con `html/template` (ADR-0003), igual que el resto de las historias del Sprint 2.
- Fuera de alcance: editar, eliminar, listar o cambiar el estado de un defecto; vincularlo manualmente a un Sprint distinto del detectado automáticamente. Quedan para historias futuras de la Épica 6 ("Gestión de defectos"). No hay autenticación ni roles.

## 6. Casos límite
- Elemento sin ningún defecto registrado todavía → el primero es válido.
- Descripción `" "` (solo espacios) → error.
- Severidad `"urgente"` o vacía → error (no es ninguna de las 4 válidas).
- Proyecto sin ningún Sprint → `SprintDeteccionID` queda `nil`, el defecto se crea igual.
- Proyecto con un Sprint `Abierto` cuyas fechas no incluyen hoy (ya terminó o todavía no empezó) → `SprintDeteccionID` queda `nil`.
- Proyecto con un Sprint `Abierto` cuyas fechas incluyen hoy → se asigna el `ID` de ese Sprint.
- Elemento que existe pero en otro proyecto → error de elemento no encontrado.
- Proyecto en un estado distinto de Activo → se puede registrar igual, mismo criterio que HU-18 (RN-18.4) — ver P-1.
- Dos defectos sobre el mismo elemento, con la misma descripción y severidad → ambos válidos, con IDs distintos (RN-20.6).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05, se reutiliza) | "elemento de backlog no encontrado" (HTTP 404) |
| Descripción vacía o solo espacios | `ErrDescripcionDefectoVacia` | "la descripción es obligatoria" (HTTP 400) |
| Severidad vacía o fuera de la lista | `ErrSeveridadInvalida` | "la severidad debe ser Crítica, Alta, Media o Baja" (HTTP 400) |
| Cuerpo JSON mal formado | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso. |

`ErrDescripcionDefectoVacia` es un error propio de HU-20 (no se reutiliza `domain.ErrDescripcionVacia` de HU-05): son campos de entidades distintas — ver P-2.

Si hay varios errores al registrar, se informa el primero en este orden: primero se buscan el proyecto y el elemento, después el dominio valida los campos en el orden de la tabla (descripción antes que severidad).

## 8. Criterios de aceptación
- **CA-20.1:** Dado un proyecto con el elemento "Login" y un Sprint Abierto cuyas fechas incluyen hoy, cuando registro un defecto con descripción "No carga la pantalla" y severidad Alta vinculado a "Login", entonces se crea con ID asignado, Estado Abierto y el SprintDeteccionID de ese Sprint. (RN-20.1, RN-20.4, RN-20.5, RN-20.6)
- **CA-20.2:** Dado un proyecto sin ningún Sprint Abierto vigente hoy (sin sprints, o con sprints cuyas fechas no incluyen hoy), cuando registro un defecto, entonces se crea con SprintDeteccionID nulo. (RN-20.5)
- **CA-20.3:** Cuando la descripción está vacía o tiene solo espacios, entonces obtengo ErrDescripcionDefectoVacia. (RN-20.2)
- **CA-20.4:** Cuando la severidad no es Crítica, Alta, Media ni Baja, entonces obtengo ErrSeveridadInvalida. (RN-20.3)
- **CA-20.5:** Dado un ID de proyecto o de elemento que no existe, o un elemento de otro proyecto, cuando registro un defecto, entonces obtengo ErrProyectoNoEncontrado o ErrElementoNoEncontrado. (RN-20.1)
- **CA-20.6:** Dado un pedido válido a `POST /proyectos/{id}/defectos`, entonces la respuesta es 201 con el defecto en JSON; con datos inválidos, 400; con proyecto o elemento inexistente, 404.
- **CA-20.7:** Dados dos defectos sobre el mismo elemento, entonces ambos se crean con IDs distintos. (RN-20.6)
- **CA-20.8:** Dada la pantalla de registrar un defecto, cuando completo el elemento, la descripción y la severidad y guardo, entonces veo la confirmación con el Sprint de detección asignado (o la indicación de que no se asignó ninguno); con datos inválidos, veo el error sin perder lo ya completado.

En todos los casos de error, no se crea ningún defecto.

## Decisiones tomadas
Las acordó el equipo (Ezequiel) en el chat el 2026-10-26, a partir del borrador de esta spec:
- **D-1 (severidad):** 4 niveles — Crítica, Alta, Media, Baja — distinta de la escala de 3 niveles de `Prioridad` en `ElementoBacklog`.
- **D-2 (estado inicial):** siempre `Abierto`, no lo elige quien registra el defecto; cambiarlo queda para una historia futura.
- **D-3 (historia relacionada):** cualquier elemento del backlog del proyecto, sin restricción de estado ni prioridad.
- **D-4 (Sprint de detección):** el Sprint `Abierto` cuyas fechas incluyen hoy, porque hoy todos los sprints quedan en estado `Abierto` (HU-11, que los cierra, todavía no existe) y hace falta mirar el rango de fechas para saber cuál está vigente.

## Decisiones confirmadas
Las preguntas P-1 a P-3 del borrador las confirmó Matías López (Agile Enabler) el 2026-10-29, tal como estaban propuestas:
- **P-1 (proyecto no Activo):** se puede registrar un defecto aunque el proyecto no esté Activo, igual que HU-18 (RN-18.4).
- **P-2 (error propio):** se usa `ErrDescripcionDefectoVacia` y no el `ErrDescripcionVacia` de HU-05. Se reutiliza un error solo si se refiere a la misma entidad, como `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado`.
- **P-3 (ruta HTTP):** `POST /proyectos/{id}/defectos`, con el elemento en el cuerpo como `elemento_id`.

Nota: D-4 se escribió cuando HU-11 todavía no existía. Con HU-11 en `main` la regla sigue igual: un sprint cerrado antes de su fecha de fin deja de ser el vigente.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-10-26: creada (borrador), con las decisiones D-1 a D-4 confirmadas en el chat; P-1 a P-3 pendientes de confirmar con el equipo.
- 2026-10-29: P-1 a P-3 confirmadas por Matías López; spec aprobada.
