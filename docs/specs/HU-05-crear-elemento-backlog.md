# HU-05: Crear un elemento de backlog

**Historia:** Como Product Builder, quiero crear un elemento de backlog con id, título, descripción, prioridad, estado, story points y criterios de aceptación, para registrar una necesidad del producto.
**Issue:** #16  ·  **Sprint:** 1  ·  **Story points:** 5
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir cargar en el Product Backlog de un proyecto las necesidades del producto (historias, tareas o pedidos), con lo mínimo para priorizarlas y estimarlas. Es la base de la Épica 2: HU-06 las edita o elimina, HU-07 las ordena por prioridad y los sprints las toman para planificar.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `titulo` | `string` | Sí | Se recortan espacios. No vacío. Máximo 100 caracteres (runas). Único dentro del proyecto, sin distinguir mayúsculas. |
| `descripcion` | `string` | Sí | Se recortan espacios. No vacía. Sin largo máximo. |
| `prioridad` | `domain.Prioridad` | Sí | Exactamente `Alta`, `Media` o `Baja` (distingue mayúsculas). |
| `storyPoints` | `*int` | No | Si se envía: entero mayor a 0. Si no se envía, el elemento queda sin estimar. |
| `criterios` | `[]string` | Sí | Al menos 1 criterio de aceptación. Cada uno se recorta y no puede quedar vacío. Se conserva el orden. |

## 3. Salidas esperadas
- Se devuelve el elemento creado con `ID` (único, autogenerado), `ProyectoID`, `Titulo` y `Descripcion` recortados, `Prioridad`, `Estado` = `Pendiente`, `StoryPoints` (o sin estimar) y `Criterios` en el orden en que se enviaron.
- El elemento queda persistido y se puede recuperar por su `ID`.
- Si hay un error, no se guarda nada.

## 4. Reglas de negocio
- **RN-05.1:** Todo elemento pertenece a un proyecto existente.
- **RN-05.2:** El título es obligatorio (un título de solo espacios es vacío) y tiene como máximo 100 caracteres, contados después de recortar.
- **RN-05.3:** No puede haber dos elementos con el mismo título dentro del mismo proyecto, comparando sin distinguir mayúsculas ni espacios alrededor. En proyectos distintos sí se puede repetir.
- **RN-05.4:** La descripción es obligatoria (una de solo espacios es vacía).
- **RN-05.5:** La prioridad es obligatoria y toma uno de tres valores: `Alta`, `Media` o `Baja`.
- **RN-05.6:** Los story points son opcionales (se pueden estimar después). Si se informan, son un entero positivo.
- **RN-05.7:** El elemento tiene al menos un criterio de aceptación, y ninguno puede estar vacío.
- **RN-05.8:** Los estados posibles de un elemento son `Pendiente`, `En progreso` y `Hecho`. Todo elemento nuevo queda en `Pendiente`. Los cambios de estado se definen en otras historias.
- **RN-05.9:** El `ID` es único y lo genera el sistema (entero autoincremental).

## 5. Restricciones
- Las validaciones de campos (RN-05.2, 05.4 a 05.8) viven en el dominio (`internal/domain`) sin dependencias externas. Reutilizan `domain.ClaveNombre` (HU-01) para comparar títulos.
- La existencia del proyecto (RN-05.1) y la unicidad del título (RN-05.3) se resuelven en el caso de uso a través de un puerto de repositorio. La unicidad también la garantiza la base (`UNIQUE (proyecto_id, titulo_clave)`), por si dos pedidos llegan a la vez.
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001): `POST /proyectos/{id}/backlog`.
- No hay autenticación ni roles en esta historia.

## 6. Casos límite
- Sin story points → válido, queda sin estimar.
- Story points = 1 → válido; = 0 o negativo → error.
- Título de exactamente 100 caracteres → válido; de 101 → error.
- Título con espacios alrededor → se guarda recortado.
- Título repetido con otras mayúsculas en el mismo proyecto (`"Login"` y `"LOGIN"`) → error de duplicado.
- Mismo título en otro proyecto → válido.
- Un solo criterio → válido; lista vacía → error; un criterio de solo espacios entre otros válidos → error.
- Prioridad `"alta"` (minúscula) o `"Urgente"` → error.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `ErrProyectoNoEncontrado` | "Proyecto no encontrado." (HTTP 404) |
| Título vacío o solo espacios | `ErrTituloVacio` | "El título es obligatorio." |
| Título de más de 100 caracteres | `ErrTituloDemasiadoLargo` | "El título no puede superar los 100 caracteres." |
| Descripción vacía | `ErrDescripcionVacia` | "La descripción es obligatoria." |
| Prioridad que no es Alta, Media ni Baja | `ErrPrioridadInvalida` | "La prioridad debe ser Alta, Media o Baja." |
| Story points 0 o negativos | `ErrStoryPointsInvalidos` | "Los story points deben ser un entero positivo." |
| Sin criterios de aceptación | `ErrSinCriterios` | "Debe haber al menos un criterio de aceptación." |
| Un criterio vacío | `ErrCriterioVacio` | "Los criterios de aceptación no pueden estar vacíos." |
| Ya hay un elemento con ese título en el proyecto | `ErrTituloDuplicado` | "Ya existe un elemento con ese título en el proyecto." |

Errores de datos: HTTP 400. Si hay varios errores, se informa el primero en el orden de la tabla: primero se busca el proyecto, después el dominio valida los campos y por último el caso de uso verifica el título duplicado.

## 8. Criterios de aceptación
- **CA-05.1:** Dado un proyecto existente, cuando creo un elemento con todos los datos válidos, entonces se crea en estado `Pendiente`, con un ID asignado y los datos enviados. (RN-05.8, RN-05.9)
- **CA-05.2:** Dado un proyecto existente, cuando creo un elemento sin story points, entonces se crea sin estimar. (RN-05.6)
- **CA-05.3:** Dado un ID de proyecto que no existe, cuando creo un elemento, entonces obtengo `ErrProyectoNoEncontrado`. (RN-05.1)
- **CA-05.4:** Cuando el título está vacío, tiene solo espacios o supera los 100 caracteres, entonces obtengo `ErrTituloVacio` o `ErrTituloDemasiadoLargo`; con exactamente 100 se crea. (RN-05.2)
- **CA-05.5:** Cuando la descripción está vacía, entonces obtengo `ErrDescripcionVacia`. (RN-05.4)
- **CA-05.6:** Cuando la prioridad no es `Alta`, `Media` ni `Baja`, entonces obtengo `ErrPrioridadInvalida`. (RN-05.5)
- **CA-05.7:** Cuando los story points son 0 o negativos, entonces obtengo `ErrStoryPointsInvalidos`. (RN-05.6)
- **CA-05.8:** Cuando no envío criterios, o alguno está vacío, entonces obtengo `ErrSinCriterios` o `ErrCriterioVacio`. (RN-05.7)
- **CA-05.9:** Dado un elemento "Login" en un proyecto, cuando creo otro "login" en el mismo proyecto, entonces obtengo `ErrTituloDuplicado`; si lo creo en otro proyecto, se crea. (RN-05.3)
- **CA-05.10:** Dados dos elementos creados, entonces tienen IDs distintos y cada uno se recupera por su ID con sus criterios en orden. (RN-05.9)

En todos los casos de error, no se crea ningún elemento.

## Decisiones tomadas
Las propuso Nicolás Andino al refinar la historia y las decidió el equipo completo en clase (2026-10-23):
- Prioridad en tres niveles: Alta, Media y Baja.
- Story points: cualquier entero positivo (no se restringe a Fibonacci).
- Estados del elemento: Pendiente, En progreso y Hecho.
- Criterios de aceptación como lista (al menos uno). Título único por proyecto y de hasta 100 caracteres, como el nombre del proyecto.
- El ID es global (autoincremental de SQLite), no un número por proyecto.

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-23: creada (borrador).
- 2026-10-23: el equipo confirma en clase las decisiones tomadas; la historia se suma al Sprint 1.
