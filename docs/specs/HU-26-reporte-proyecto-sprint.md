# HU-26: Generar reporte de proyecto o Sprint

**Historia:** Como Product Owner, quiero generar un reporte de un proyecto o Sprint con historias, estimaciones, esfuerzo, métricas y defectos, para compartirlo con los interesados.
**Issue:** #37  ·  **Sprint:** 4  ·  **Story points:** 5
**Estado de la spec:** Aprobada

## 1. Objetivo
Armar, para un proyecto o para un Sprint de ese proyecto, un reporte de solo lectura con cinco secciones (Resumen, Historias y estimaciones, Esfuerzo, Métricas, Defectos) que junta en un solo lugar datos que hoy están repartidos en varias pantallas (HU-04, HU-07, HU-22, HU-23, HU-24, HU-21), para que el Product Owner lo comparta con los interesados sin tener que recorrer el sistema.

**Alcance real de esta historia:** no agrega ningún cálculo nuevo. Cada sección reutiliza tal cual un caso de uso o una función del dominio que ya existe (sección 5). Lo único que agrega HU-26 es **juntar** esos datos en un reporte y mostrarlos en una pantalla propia.

**Exportar a PDF:** es responsabilidad de HU-27, no de esta historia (decisión de Matías López en la revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378); ver Historial). HU-26 no expone ningún puerto ni conversión para HU-27.

## 2. Entradas
**Reporte de un proyecto.** `GET /web/proyectos/{id}/reporte`

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |

**Reporte de un Sprint.** `GET /web/proyectos/{id}/sprints/{sprintID}/reporte`

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |
| `sprintID` | `int64` | Sí | Entero positivo. Debe corresponder a un Sprint de ese proyecto. |

Son dos consultas de solo lectura: no hay más entradas, ni formulario. El Sprint puede estar en cualquier estado, `Abierto` o `Cerrado` (RN-26.10; decisión confirmada con Ezequiel en el chat, igual criterio que HU-23 RN-23.12).

**Sin endpoint JSON propio:** a diferencia de HU-16, HU-17, HU-21 a HU-24, que exponen el mismo dato por JSON y por pantalla, HU-26 solo tiene pantalla HTML. El reporte no agrega ningún dato nuevo (cada número ya tiene su propio endpoint JSON en la historia que lo calculó). Si el equipo necesita el reporte en JSON más adelante, se agrega como una historia propia. Confirmada por Matías López (Agile Enabler) el 2026-11-08, en la revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378).

## 3. Salidas esperadas
Las dos pantallas (de proyecto y de Sprint) muestran el mismo reporte, con las mismas cinco secciones, en este orden:

1. **Resumen:** una tabla de dato/valor con los datos descriptivos.
   - Proyecto: Nombre, Estado, Fecha de inicio, Fecha de fin ("Sin fecha de fin" si no tiene), cantidad de Sprints Abiertos y Cerrados.
   - Sprint: Número, Sprint Goal, Fecha de inicio, Fecha de fin, Estado.
2. **Historias y estimaciones:** una tabla con Título, Prioridad, Estado y Story Points ("Sin estimar" si no tiene) de las historias del reporte (RN-26.1, RN-26.2), más una línea con el resumen: "X de Y historias completadas (Z de W story points)".
3. **Esfuerzo:** los mismos datos que ya da HU-23 (`domain.TotalHoras`) sobre esas mismas historias, en una tabla de dato/valor: Horas estimadas, Horas reales, Desviación en horas, Desviación en porcentaje ("Sin estimar" si ningún elemento tiene horas estimadas, RN-23.6).
4. **Métricas:** una tabla con una fila por Sprint:
   - Proyecto: una fila por cada Sprint `Cerrado` (RN-26.5), igual que HU-24, más la velocidad del equipo (HU-22) como un dato aparte antes de la tabla.
   - Sprint: una sola fila, la de ese Sprint, calculada igual que HU-24 pero sin exigir que esté `Cerrado` (RN-26.6). Si está `Cerrado`, antes de la tabla va una aclaración (RN-26.11) sobre por qué el resumen de historias (sección 2) y el % completado pueden parecer contradictorios cuando una historia volvió al backlog al cerrarse.
   - Columnas: Sprint (número), % completado ("Sin datos" si no hay, RN-24.2), Defectos detectados, Defectos resueltos.
5. **Defectos:** una tabla con Elemento, Descripción, Severidad, Estado, Sprint detección y Sprint resolución (mismas columnas que la pantalla de HU-21) de los defectos del reporte (RN-26.7, RN-26.8).

Una sección sin filas muestra un mensaje en vez de una tabla vacía (RN-26.9), igual criterio que el resto de las pantallas del proyecto (por ejemplo, "El proyecto todavía no tiene sprints cerrados." de HU-12).

**En HTTP:** `200` con la pantalla; `400` si el ID no es un número; `404` si el proyecto no existe, o si el Sprint no existe o es de otro proyecto.

## 4. Reglas de negocio
- **RN-26.1:** El reporte de un proyecto incluye **todas** las historias de su backlog (`app.ListarBacklog`, HU-07), estén o no asignadas a un Sprint, completadas o no.
- **RN-26.2:** El reporte de un Sprint incluye solo las historias cuyo `SprintID` es ese Sprint (mismo filtro que ya usa la pantalla de asignación de HU-09, `pantallas_asignacion.go`).
- **RN-26.3:** Una historia está **completada** si su `Estado` es `Hecho` (HU-10); las demás están **planificadas** (sin terminar todavía). El resumen de la sección 2 cuenta cuántas de las historias del reporte están completadas, y sus story points, sobre el total de historias y story points del reporte (una historia sin estimar no suma a ninguno de los dos totales de story points, solo al total de historias).
- **RN-26.4:** El esfuerzo del reporte (sección 3) se calcula con `domain.TotalizarHoras` (HU-23, RN-23.1 a RN-23.6) sobre las mismas historias de la sección 2: todas las del proyecto, o solo las del Sprint.
- **RN-26.5:** Las métricas del reporte de un **proyecto** son la velocidad del equipo (HU-22, `domain.Velocidad`) y, por cada Sprint `Cerrado`, su porcentaje completado y sus defectos detectados/resueltos (HU-24, `app.ConsultarCalidadSprints`, sin recalcular). Un Sprint `Abierto` no aparece en esta tabla (mismo criterio que HU-24, RN-24.5).
- **RN-26.6:** Las métricas del reporte de un **Sprint** son su propio porcentaje completado y sus defectos detectados/resueltos, calculados igual que HU-24 pero **sin exigir que esté `Cerrado`** (`domain.PorcentajeCompletado` y `domain.ContarDefectosDeSprint`, llamadas directo sobre ese Sprint): un Sprint `Abierto`, o uno `Cerrado` sin puntos planificados, muestra "Sin datos" en el porcentaje (RN-24.2), pero sus defectos detectados y resueltos hasta el momento se cuentan igual.
- **RN-26.7:** Los defectos del reporte de un **proyecto** son todos los defectos del proyecto (HU-21, `app.ListarDefectos`), sin filtrar por Sprint.
- **RN-26.8:** Los defectos del reporte de un **Sprint** son los que tienen ese Sprint como Sprint de detección o como Sprint de resolución (mismo criterio que `domain.ContarDefectosDeSprint`, HU-24): un defecto detectado en otro Sprint y resuelto en este cuenta igual; uno sin Sprint asignado no aparece en ningún reporte de Sprint (sí en el del proyecto, por RN-26.7).
- **RN-26.9:** Una sección sin datos (sin historias, sin Sprints cerrados, sin defectos) muestra un mensaje en vez de una tabla vacía; no es un error.
- **RN-26.10:** El reporte se arma para un Sprint en cualquier estado, `Abierto` o `Cerrado` (igual criterio que HU-23, RN-23.12): no hace falta esperar a que se cierre para pedirlo.
- **RN-26.11 (aclaración en el reporte de un Sprint Cerrado):** cuando el Sprint del reporte está `Cerrado`, la sección de Métricas muestra, antes de la tabla, una aclaración de por qué el resumen de historias (sección 2, RN-26.2 y RN-26.3) y el % completado (RN-26.6) pueden parecer contradictorios: las historias que volvieron al backlog al cerrar el Sprint (RN-26.2) ya no están en la sección de historias, pero sí cuentan en los puntos planificados que guardó el cierre (HU-11, RN-11.6), que es justamente lo que usa el % completado. Mismo caso que HU-23 (RN-23.8, opción A), con el texto adaptado a historias y story points en vez de horas. Un Sprint `Abierto`, y el reporte de un proyecto (que no tiene un único Sprint), no la muestran.

## 5. Restricciones
- **Dominio** (`internal/domain`, sin dependencias externas), con TDD:
  - `reporte.go`: `Reporte` (tipo armado por `ArmarReporteProyecto`/`ArmarReporteSprint`, con sus cinco secciones ya resueltas: `FilaResumen` (dato/valor), las historias (`[]ElementoBacklog`, ya existe), el esfuerzo (`TotalHoras`, ya existe de HU-23), las métricas (`[]FilaMetrica`, un tipo nuevo chico: Sprint, PorcentajeCompletado, HayPorcentaje, DefectosDetectados, DefectosResueltos — mismo contenido que `app.CalidadSprint` de HU-24, pero en el dominio, porque el dominio no puede depender de `internal/app`) y los defectos (`[]Defecto`, ya existe).
  - `ArmarReporteProyecto(proyecto Proyecto, historias []ElementoBacklog, esfuerzo TotalHoras, velocidad float64, hayVelocidad bool, metricas []FilaMetrica, defectos []Defecto) Reporte` y `ArmarReporteSprint(proyecto Proyecto, sprint Sprint, historias []ElementoBacklog, esfuerzo TotalHoras, metrica FilaMetrica, defectos []Defecto) Reporte`: funciones puras que arman el `Reporte` a partir de los datos ya calculados (el resumen de RN-26.3, los mensajes de "sin datos" de RN-26.9); no leen nada ni llaman a ningún puerto. Implementa HU-26 (RN-26.1 a RN-26.9).
- **Casos de uso** (`internal/app`), cada uno reutiliza los puertos que ya existen, sin agregar ninguno de persistencia nuevo:
  - `ArmarReporteProyecto`: busca el proyecto (`LectorProyectos`, HU-04), lista sus historias (`app.ListarBacklog`, HU-07), totaliza su esfuerzo (`LectorHoras` + `domain.TotalizarHoras`, HU-23), lista sus defectos (`app.ListarDefectos`, HU-21), y calcula su velocidad y su calidad por Sprint llamando directo a las mismas funciones del dominio que usan `app.ConsultarVelocidad` (HU-22) y `app.ConsultarCalidadSprints` (HU-24) — `domain.Velocidad`, `domain.PorcentajeCompletado`, `domain.ContarDefectosDeSprint` — en vez de esos dos casos de uso, para no repetir las consultas a los repositorios que ya hace `Ejecutar` (`BuscarPorID`, `ListarSprints`, `ListarDefectos`). **Nota (revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378)):** esto duplica la lógica de esos dos casos de uso en vez de reutilizarlos, que es lo que decía esta restricción originalmente; el PR [#375](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/375) (HU-25) agrega `calidadDe`, pensada justo para este cálculo — si entra antes que esta historia, conviene que HU-26 lo reutilice también. No bloquea esta historia. Arma el reporte con `domain.ArmarReporteProyecto`.
  - `ArmarReporteSprint`: busca el proyecto, busca el Sprint (`LectorSprintPorID`, HU-23) y verifica que sea suyo (RN-26.10, mismo patrón que `ConsultarHorasSprint`), filtra sus historias (RN-26.2) y sus defectos (RN-26.8) de las mismas listas completas del proyecto, calcula su métrica con `domain.PorcentajeCompletado` y `domain.ContarDefectosDeSprint` (sin pasar por `ConsultarCalidadSprints`, que exige `Cerrado`), y arma el reporte con `domain.ArmarReporteSprint`.
- **HTTP/pantalla:** sin endpoint JSON (sección 2). Dos rutas nuevas en `pantallas.go`, junto a las demás `/web/...`: `GET /web/proyectos/{id}/reporte` y `GET /web/proyectos/{id}/sprints/{sprintID}/reporte`. Reutilizan `proyectoIDDesdeRuta`, `sprintIDDesdeRuta` y `errorDeConsulta`. Una sola plantilla `reporte.html` para los dos casos (con el título que arma el dominio).
- **Enlace desde las pantallas existentes:** un enlace "Ver reporte" en el detalle del proyecto (HU-04) y en el detalle de un Sprint (HU-09), igual patrón que "Sprints anteriores" de HU-12.
- **Aclaración del Sprint Cerrado (RN-26.11):** una constante de texto nueva en `pantallas_reporte.go` (mismo patrón que `textoAclaracionSprintCerrado` de HU-23 en `pantallas_horas.go`), agregada a los datos de la pantalla solo cuando el Sprint del reporte está `Cerrado`; se muestra en `reporte.html`, en la sección de Métricas, antes de la tabla.
- **Sin persistencia nueva:** no hay tablas, columnas ni migraciones (tarea 26.6 no aplica, mismo caso que la 24.6 de HU-24).

## 6. Casos límite
- Proyecto sin ningún elemento en el backlog → sección de historias con el mensaje de "sin datos", esfuerzo en 0 (RN-23.2), métricas sin filas si tampoco tiene Sprints cerrados, defectos sin filas.
- Proyecto sin ningún Sprint cerrado → sin velocidad (RN-22.3) y sin filas en la tabla de métricas, pero el resto del reporte se arma igual.
- Sprint `Abierto` sin ninguna historia asignada todavía → historias y esfuerzo en 0, métrica con "Sin datos" en el porcentaje y 0 defectos (si todavía no se detectó ninguno ahí).
- Sprint `Cerrado` con 0 puntos planificados → "Sin datos" en el porcentaje (RN-24.2), igual que en la pantalla de HU-24.
- Defecto detectado en el Sprint 1 y resuelto en el Sprint 2 → aparece en el reporte de los dos Sprints (RN-26.8), y en el del proyecto una sola vez (RN-26.7).
- Defecto sin Sprint de detección ni de resolución → aparece en el reporte del proyecto, en ningún reporte de Sprint.
- Historia sin story points (`StoryPoints` `nil`) → cuenta para el total de historias del resumen (RN-26.3), no para el total de story points.
- Sprint de otro proyecto → se trata como inexistente (RN-26.10, mismo criterio que RN-23.11).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404, página de error) |
| No existe el Sprint, o es de otro proyecto | `domain.ErrSprintNoEncontrado` (HU-09, se reutiliza) | "sprint no encontrado" (HTTP 404, página de error) |
| ID de proyecto o de Sprint no numérico o no positivo | — (adaptador HTTP, igual que el resto de las historias) | HTTP 400; no se llama al caso de uso. |

No hay más condiciones de error: las dos son consultas de solo lectura.

## 8. Criterios de aceptación
- **CA-26.1:** Dado un proyecto con historias en distintos estados, cuando se pide su reporte, entonces la sección de historias las trae todas, con su título, prioridad, estado y story points. (RN-26.1)
- **CA-26.2:** Dado un Sprint con algunas historias asignadas y otras del proyecto sin asignar, cuando se pide el reporte de ese Sprint, entonces la sección de historias trae solo las asignadas a él. (RN-26.2)
- **CA-26.3:** Dado un reporte con 4 historias, 3 completadas (5, 8 y 3 story points) y 1 sin terminar (2 story points), cuando se arma, entonces el resumen dice "3 de 4 historias completadas (16 de 18 story points)". (RN-26.3)
- **CA-26.4:** Dado un reporte con historias con y sin horas estimadas, cuando se arma, entonces la sección de esfuerzo coincide con lo que calcula HU-23 sobre esas mismas historias. (RN-26.4)
- **CA-26.5:** Dado un proyecto con 2 Sprints cerrados, cuando se pide su reporte, entonces la sección de métricas trae la velocidad del equipo y una fila por cada Sprint cerrado, con el mismo porcentaje y los mismos defectos que da HU-24. (RN-26.5)
- **CA-26.6:** Dado un proyecto sin ningún Sprint cerrado, cuando se pide su reporte, entonces la sección de métricas no tiene velocidad ni filas, con el mensaje de "sin datos". (RN-26.5, RN-26.9)
- **CA-26.7:** Dado un Sprint `Cerrado` con puntos planificados, cuando se pide su reporte, entonces la métrica trae su porcentaje completado real. (RN-26.6)
- **CA-26.8:** Dado un Sprint `Abierto`, cuando se pide su reporte, entonces la métrica dice "Sin datos" en el porcentaje, pero cuenta igual sus defectos detectados y resueltos hasta el momento. (RN-26.6)
- **CA-26.9:** Dado un proyecto con defectos en distintos Sprints, cuando se pide su reporte, entonces la sección de defectos los trae todos. (RN-26.7)
- **CA-26.10:** Dado un defecto detectado en el Sprint 1 y resuelto en el Sprint 2, cuando se piden los reportes de esos dos Sprints, entonces aparece en los dos. (RN-26.8)
- **CA-26.11:** Dado un ID de proyecto que no existe, cuando se pide el reporte del proyecto o de un Sprint, entonces la respuesta es `404` con "proyecto no encontrado". (Sección 7)
- **CA-26.12:** Dado un ID de Sprint que no existe, o de un Sprint de otro proyecto, cuando se pide su reporte, entonces la respuesta es `404` con "sprint no encontrado". (RN-26.10, sección 7)
- **CA-26.13:** Dado un ID de proyecto o de Sprint inválido, cuando se pide un reporte, entonces la respuesta es `400` y no se llama al caso de uso. (Sección 7)
- **CA-26.14:** Dado un Sprint `Cerrado` con una historia que volvió al backlog al cerrarse, cuando se pide su reporte, entonces la sección de Métricas muestra la aclaración; dado un Sprint `Abierto`, cuando se pide su reporte, entonces no la muestra. (RN-26.11)
- **CA-26.15:** Dado un proyecto sin historias, sin Sprints cerrados y sin defectos, cuando se pide su reporte, entonces cada sección (salvo el resumen) muestra su mensaje de "sin datos", sin romper la pantalla. (RN-26.9)

## Preguntas abiertas
Ninguna: P-1 y P-3 quedaron confirmadas en la revisión del PR #378 (ver Historial); P-2 ya no aplica, porque la conversión a PDF pasó a HU-27.

## Historial
- 2026-11-07: creada (borrador). Alcance de las 5 secciones y "cualquier estado de Sprint" confirmados con Ezequiel en el chat antes de escribir la spec. P-1 a P-3 son decisiones propuestas, pendientes de confirmar con el equipo.
- 2026-11-08: aprobada por Matías López (Agile Enabler), en la revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378):
  - P-1 (sin endpoint JSON propio) y P-3 ("Ver reporte" en el detalle del proyecto y del Sprint) confirmadas tal cual estaban propuestas.
  - La conversión a PDF (antes RN-26.11 y CA-26.14) pasa a HU-27: se sacan de esta spec esa RN y ese CA, la parte de la sección 1 sobre el contrato con HU-27, la función `ADocumentoPDF` y los puertos `FuenteDocumento…` de la sección 5. P-2 ya no aplica, se saca.
  - Se agrega, en su lugar, una aclaración en el reporte de un Sprint `Cerrado` (nueva RN-26.11/CA-26.14, punto 1 de la revisión) para que el resumen de historias y el % completado no parezcan contradictorios cuando una historia volvió al backlog al cerrarse.
  - Se corrige la numeración de RN citada en las secciones 2 y 3 (había quedado corrida un número respecto de la sección 4, punto 3 de la revisión).
  - Se corrige la sección 5 sobre cómo calcula HU-26 la calidad por Sprint, para que describa lo que hace el código (duplica `ConsultarCalidadSprints`/`ConsultarVelocidad` en vez de llamarlos) en lugar de lo que decía la spec original (punto 2 de la revisión; no bloquea, se resuelve cuando se mergee HU-25 o en una corrección aparte).
- 2026-11-11: DEF-05 (#380): los defectos del reporte mostraban el ID del sprint en lugar de su número; el Resumen del proyecto no mostraba la cantidad de sprints abiertos y cerrados que pedía la sección 3; y la velocidad salía con todos los decimales ("24,333…"). Ahora muestra el número, las dos cantidades y la velocidad con un decimal, como HU-22. El reporte de un sprint lista también los sprints del proyecto, para los números de sus defectos. Las reglas y los criterios no cambian.
