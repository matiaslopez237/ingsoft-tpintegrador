# HU-08: Crear un Sprint con Sprint Goal

**Historia:** Como Scrum Master, quiero crear un Sprint con fecha de inicio, fin y un Sprint Goal, para organizar el trabajo del equipo.
**Issue:** #19  ·  **Sprint:** 2  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir crear los sprints de un proyecto: un período con fechas y un objetivo (Sprint Goal) al que después se le asignan historias del backlog (HU-09), se marcan completadas (HU-10) y se cierra (HU-11). Es la base de la Épica 3 y de las métricas por sprint (Épica 7).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente y en estado `Activo`. |
| `fechaInicio` | `time.Time` | Sí | Fecha válida (solo día). |
| `fechaFin` | `time.Time` | Sí | Fecha válida, no anterior a `fechaInicio`. |
| `objetivo` | `string` | Sí | El Sprint Goal. Se recortan espacios; no vacío; máximo 200 caracteres (runas). |

**En HTTP** (`POST /proyectos/{id}/sprints`), el proyecto va en la ruta y el resto en un cuerpo JSON con fechas `AAAA-MM-DD`:

```json
{ "fecha_inicio": "2026-10-26", "fecha_fin": "2026-10-27", "objetivo": "Gestionar el backlog desde la web" }
```

## 3. Salidas esperadas
- Se devuelve el sprint creado con `ID`, `ProyectoID`, `Numero`, `FechaInicio`, `FechaFin`, `Objetivo` (recortado) y `Estado` = `Abierto`.
- El sprint queda persistido y aparece en la lista de sprints del proyecto, ordenada por número.
- Si hay un error, no se guarda nada.

En HTTP se responde `201` con el sprint:

```json
{ "id": 1, "proyecto_id": 1, "numero": 1, "fecha_inicio": "2026-10-26", "fecha_fin": "2026-10-27", "objetivo": "Gestionar el backlog desde la web", "estado": "Abierto" }
```

La lista de sprints del proyecto se consulta con `GET /proyectos/{id}/sprints`: `200` con un array de sprints con el mismo formato, ordenados por número (`[]` si el proyecto no tiene sprints), o `404` si el proyecto no existe.

## 4. Reglas de negocio
- **RN-08.1:** Todo sprint pertenece a un proyecto existente.
- **RN-08.2:** Solo se pueden crear sprints en proyectos en estado `Activo` (RN-01.7). En un proyecto En pausa, Finalizado, Postergado o Cancelado no se planifica trabajo nuevo.
- **RN-08.3:** Las fechas de inicio y de fin son obligatorias, y la de fin no puede ser anterior a la de inicio. Un sprint de un solo día (inicio = fin) es válido.
- **RN-08.4:** El sprint tiene que estar dentro de las fechas del proyecto: no puede empezar antes del inicio del proyecto y, si el proyecto tiene fecha de fin, no puede terminar después de ella.
- **RN-08.5:** Los sprints de un mismo proyecto no se superponen. Las fechas son días completos: dos sprints que comparten un día se superponen. Un sprint que empieza el día siguiente al fin de otro es válido. Sprints de proyectos distintos no se comparan.
- **RN-08.6:** El Sprint Goal es obligatorio (uno de solo espacios es vacío) y tiene como máximo 200 caracteres, contados después de recortar.
- **RN-08.7:** Cada sprint tiene un número dentro de su proyecto, que asigna el sistema: el primero es 1 y cada nuevo es el siguiente al mayor existente. Proyectos distintos tienen su propia numeración.
- **RN-08.8:** Todo sprint nuevo queda en estado `Abierto`. HU-11 lo pasa a `Cerrado`.
- **RN-08.9:** El `ID` es único en todo el sistema y lo genera la persistencia (entero autoincremental), como en HU-01 y HU-05.

## 5. Restricciones
- Las validaciones de datos (RN-08.2 a RN-08.6, RN-08.8) viven en el dominio (`internal/domain`) sin dependencias externas. El dominio recibe el proyecto y los sprints que ya tiene, así puede validar el rango y la superposición sin acceder a la base.
- La existencia del proyecto (RN-08.1) se resuelve en el caso de uso a través de un puerto de repositorio.
- **Pedidos simultáneos:** el número (RN-08.7) y la no superposición (RN-08.5) se vuelven a garantizar en la base al guardar, en una sola operación. Si dos pedidos crean sprints a la vez, no pueden quedar dos con el mismo número ni dos superpuestos: el segundo se rechaza con `ErrSprintSuperpuesto`, o recibe el número siguiente si no se superpone.
- Si después se modifican las fechas del proyecto (HU-02), los sprints existentes no se revalidan: RN-08.4 se controla solo al crear.
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001) y pantallas con `html/template` (ADR-0003).

## 6. Casos límite
- Primer sprint de un proyecto → número 1. Primer sprint de otro proyecto → también número 1.
- Sprint de un día (inicio = fin) → válido.
- Sprint que empieza el mismo día que el proyecto, o termina el mismo día que el fin del proyecto → válido.
- Sprint que empieza un día antes que el proyecto, o termina un día después de su fin → error.
- Proyecto sin fecha de fin → el sprint puede terminar en cualquier fecha posterior a su inicio.
- Sprint que empieza el mismo día que termina otro → error de superposición.
- Sprint que empieza el día siguiente al fin de otro → válido.
- Sprint que contiene por completo a otro, o está contenido en otro → error de superposición.
- Sprint Goal de exactamente 200 caracteres → válido; de 201 → error.
- Sprint Goal con espacios alrededor → se guarda recortado.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `ErrProyectoNoEncontrado` (HU-01) | "proyecto no encontrado" (HTTP 404) |
| El proyecto no está Activo | `ErrProyectoNoActivo` | "solo se pueden crear sprints en proyectos activos" |
| Falta la fecha de inicio | `ErrFechaInicioObligatoria` (HU-01) | "la fecha de inicio es obligatoria" |
| Falta la fecha de fin | `ErrFechaFinObligatoria` | "la fecha de fin es obligatoria" |
| Fecha de fin anterior a la de inicio | `ErrFechaFinAnterior` (HU-01) | "la fecha de fin no puede ser anterior a la de inicio" |
| Sprint Goal vacío o de solo espacios | `ErrObjetivoVacio` | "el Sprint Goal es obligatorio" |
| Sprint Goal de más de 200 caracteres | `ErrObjetivoDemasiadoLargo` | "el Sprint Goal no puede superar los 200 caracteres" |
| Fechas fuera de las del proyecto | `ErrSprintFueraDelProyecto` | "las fechas del sprint tienen que estar dentro de las del proyecto" |
| Se superpone con otro sprint del proyecto | `ErrSprintSuperpuesto` | "las fechas se superponen con otro sprint del proyecto" |

Errores de datos: HTTP 400. Si hay varios errores, se informa el primero en el orden de la tabla.

## 8. Criterios de aceptación
- **CA-08.1:** Dado un proyecto activo sin sprints, cuando creo un sprint con fechas dentro del proyecto y un Sprint Goal, entonces se crea con número 1, estado `Abierto` y los datos enviados. (RN-08.7, RN-08.8)
- **CA-08.2:** Dado un proyecto con un sprint, cuando creo otro que no se superpone, entonces recibe el número 2; y el primer sprint de otro proyecto recibe el número 1. (RN-08.7)
- **CA-08.3:** Dado un ID de proyecto que no existe, cuando creo un sprint, entonces obtengo `ErrProyectoNoEncontrado`. (RN-08.1)
- **CA-08.4:** Dado un proyecto que no está Activo, cuando creo un sprint, entonces obtengo `ErrProyectoNoActivo`. (RN-08.2)
- **CA-08.5:** Cuando falta alguna fecha o la de fin es anterior a la de inicio, entonces obtengo `ErrFechaInicioObligatoria`, `ErrFechaFinObligatoria` o `ErrFechaFinAnterior`; un sprint de un día se crea. (RN-08.3)
- **CA-08.6:** Cuando el Sprint Goal está vacío o supera los 200 caracteres, entonces obtengo `ErrObjetivoVacio` o `ErrObjetivoDemasiadoLargo`; con exactamente 200 se crea. (RN-08.6)
- **CA-08.7:** Cuando el sprint empieza antes que el proyecto o termina después de su fin, entonces obtengo `ErrSprintFueraDelProyecto`; justo en los límites se crea. (RN-08.4)
- **CA-08.8:** Dado un sprint del 26 al 27, cuando creo otro que comparte al menos un día con él, entonces obtengo `ErrSprintSuperpuesto`; si empieza el 28, se crea. (RN-08.5)
- **CA-08.9:** Dados dos sprints creados, cuando consulto los sprints del proyecto, entonces aparecen los dos, ordenados por número. (Sección 3)

En todos los casos de error, no se crea ningún sprint.

## Decisiones tomadas
Definidas por Matías López (2026-10-26):
- Número automático por proyecto en lugar de un nombre libre.
- Los sprints de un proyecto no se superponen.
- Las fechas del sprint tienen que estar dentro de las del proyecto.
- Solo se crean sprints en proyectos Activos.

Propuestas para validar en el PR: estado inicial `Abierto` (HU-11 lo pasa a `Cerrado`), Sprint Goal de hasta 200 caracteres y sprint de un día válido.

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-26: creada (borrador).
- 2026-10-26: se agrega la consulta de la lista de sprints por HTTP (sección 3).
