# HU-24: Calcular % de story points completados y defectos por Sprint

**Historia:** Como Scrum Master, quiero ver el porcentaje de story points completados y la cantidad de defectos detectados/resueltos de cada Sprint, para evaluar la calidad del Sprint.
**Issue:** #35  ·  **Sprint:** 4  ·  **Story points:** 3
**Estado de la spec:** Aprobada

## 1. Objetivo
Dar al Scrum Master, por cada Sprint ya `Cerrado` de un proyecto, dos datos de calidad que hoy no están juntos en ningún lado: qué porcentaje de los story points planificados se completó, y cuántos defectos se detectaron y se resolvieron durante ese Sprint. Sirve para leer de un vistazo si un Sprint fue "bueno" (completó lo planificado, sin dejar muchos defectos) o no, sin tener que cruzar a mano la tabla de sprints con la lista de defectos.

**Alcance real de esta historia:** no agrega persistencia nueva. Reutiliza datos que ya guardan otras historias: los puntos planificados/completados de cada Sprint (HU-11, RN-11.6) y el Sprint de detección/resolución de cada defecto (HU-20 RN-20.5, HU-21 RN-21.6). Lo único que agrega HU-24 es el cálculo del porcentaje y el conteo de defectos por Sprint, igual criterio que HU-22 con la velocidad del equipo.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |

Es una consulta de solo lectura: no hay más entradas. En HTTP, `proyectoID` va en la ruta (`GET /proyectos/{id}/calidad`).

## 3. Salidas esperadas
Por cada Sprint `Cerrado` del proyecto (mismo conjunto y orden que `SprintsAnteriores` de HU-12: del más reciente al más antiguo por fecha de inicio), se informa:
- El **porcentaje de story points completados**: `PuntosCompletados / PuntosPlanificados × 100` (RN-11.6), redondeado a **un decimal** (mismo criterio que la desviación de HU-19/HU-23, RN-19.6). Si el Sprint se cerró con 0 puntos planificados, el porcentaje queda **indefinido**, nunca `0` (RN-24.2).
- La **cantidad de defectos detectados** en ese Sprint: defectos del proyecto cuyo `SprintDeteccionID` es el de ese Sprint (RN-20.5).
- La **cantidad de defectos resueltos** en ese Sprint: defectos del proyecto cuyo `SprintResolucionID` es el de ese Sprint (RN-21.6), sin importar si después se cerraron.

Los Sprints `Abierto` no entran en esta salida (mismo criterio que la pantalla de sprints anteriores de HU-12): sus puntos valen 0 por diseño (RN-11.6, mientras no se cierran) y mostrar un porcentaje ahí sería engañoso, no un dato real todavía.

- **En HTTP** (`GET /proyectos/{id}/calidad`): `200` con un array, uno por Sprint cerrado (lista vacía si no tiene ninguno), cada uno con `sprint_id`, `numero`, `porcentaje_completado` (`null` si no hay datos), `defectos_detectados`, `defectos_resueltos`; `400` si el ID no es un número; `404` si el proyecto no existe.
- **En la pantalla:** se agregan tres columnas a la tabla ya existente de sprints anteriores (HU-12, `/web/proyectos/{id}/sprints/anteriores`): "% completado", "Defectos detectados" y "Defectos resueltos". No hay pantalla nueva (mismo patrón que HU-22 con "Velocidad del equipo").

## 4. Reglas de negocio
- **RN-24.1:** El porcentaje de story points completados de un Sprint `Cerrado` es `PuntosCompletados / PuntosPlanificados × 100`, usando los valores que HU-11 calculó y guardó al cerrarlo (RN-11.6). HU-24 no los recalcula ni los vuelve a sumar desde las historias. Se redondea a **un decimal**, con los medios alejándose de cero (mismo criterio y función que `redondearUnDecimal` de HU-19, RN-19.6), para que la API y la pantalla siempre muestren el mismo valor.
- **RN-24.2:** Si `PuntosPlanificados` es `0`, el porcentaje queda **indefinido**: no se informa como `0%`, porque `0%` sería un dato real (se planificó algo y no se completó nada), distinto de "no se planificó nada". Mismo criterio que RN-22.3 (velocidad sin Sprints cerrados).
- **RN-24.3:** Los defectos **detectados** de un Sprint son los que tienen ese Sprint como `SprintDeteccionID` (RN-20.5). Un defecto sin Sprint de detección (`nil`, porque no había ningún Sprint `Abierto` vigente el día que se detectó) no cuenta en ningún Sprint.
- **RN-24.4:** Los defectos **resueltos** de un Sprint son los que tienen ese Sprint como `SprintResolucionID` (RN-21.6), sin importar el estado actual del defecto al momento de la consulta: `SprintResolucionID` se fija una sola vez, al resolver (`Resolver()`, HU-21), y no cambia si después se cierra. Un defecto sin Sprint de resolución (`nil`, todavía no resuelto, o resuelto sin ningún Sprint `Abierto` vigente ese día) no cuenta en ningún Sprint.
- **RN-24.5:** Solo se calculan estos datos para Sprints `Cerrado`; un Sprint `Abierto` no aparece en esta consulta (mismo filtro que `SprintsAnteriores`, HU-12).

## 5. Restricciones
- El cálculo vive en el dominio (`internal/domain/sprint_historial.go`, junto a `SprintsAnteriores` y `Velocidad` de HU-12/HU-22), como funciones puras sin dependencias externas:
  - `PorcentajeCompletado(s Sprint) (porcentaje float64, hayDatos bool)` (RN-24.1, RN-24.2).
  - `ContarDefectosDeSprint(defectos []Defecto, sprintID int64) (detectados, resueltos int)` (RN-24.3, RN-24.4).
- Caso de uso nuevo `ConsultarCalidadSprints` (`internal/app/consultar_calidad_sprints.go`), con la misma forma que `ConsultarVelocidad` (HU-22): busca el proyecto (`LectorProyectos`, reutilizado), lista sus Sprints (`RepositorioSprints.ListarSprints`, reutilizado de HU-08) y sus defectos (`LectorDefectos.ListarDefectos`, reutilizado de HU-21), filtra los Sprints con `domain.SprintsAnteriores` y arma, para cada uno, su porcentaje y sus conteos de defectos. **No agrega ningún puerto de persistencia nuevo.**
- El porcentaje **nunca puede superar el 100%**: `ResumirCierre` (HU-11) calcula `PuntosPlanificados` sumando los story points de **todas** las historias que tiene el Sprint al cerrarlo, y `PuntosCompletados` sumando solo las que están `Hecho`, un subconjunto de esas mismas historias. Por diseño, `PuntosCompletados` nunca es mayor que `PuntosPlanificados`.
- La pantalla reutiliza el handler de sprints anteriores de HU-12 (`pantallas_sprints.go`): agrega el llamado a `ConsultarCalidadSprints` y los datos correspondientes a cada fila de `datosSprintsAnteriores`, sin crear una pantalla ni una plantilla nueva.
- Endpoint JSON propio (`GET /proyectos/{id}/calidad`), separado del de HU-12 (`GET /proyectos/{id}/sprints/anteriores`) y del de HU-22 (`GET /proyectos/{id}/velocidad`), para no tocar los contratos que ya exponen esas historias.
- Persistencia en SQLite (ADR-0002), pero sin cambios de esquema: no hay tabla ni columna nueva.

## 6. Casos límite
- Sprint cerrado sin ningún defecto detectado ni resuelto → `0` y `0` (son conteos reales; distinto del porcentaje, que si no hay puntos planificados queda indefinido, no en `0`).
- Defecto detectado en un Sprint y resuelto en un Sprint posterior → cuenta como "detectado" en el primero y como "resuelto" en el segundo (pueden ser Sprints distintos).
- Defecto detectado y resuelto dentro del mismo Sprint → cuenta en ambas columnas de ese Sprint.
- Defecto sin `SprintDeteccionID` (detectado sin ningún Sprint `Abierto` vigente ese día) → no cuenta como detectado en ningún Sprint, pero el defecto existe igual (lo sigue mostrando HU-21).
- Defecto `Resuelto` o `Cerrado` pero con `SprintResolucionID` nulo (resuelto sin ningún Sprint `Abierto` vigente ese día) → no cuenta como resuelto en ningún Sprint.
- Defecto todavía `Abierto` o `En progreso` → no cuenta en la columna de resueltos de ningún Sprint (`SprintResolucionID` sigue `nil` hasta que se resuelve).
- Sprint cerrado con puntos planificados mayor a cero y puntos completados en `0` → `0%` (dato real: se planificó algo y no se completó nada; distinto del caso "sin datos" de RN-24.2).
- Sprint cerrado con puntos completados igual a los planificados → `100%` (el máximo posible; ver sección 5, nunca se supera).
- Sprint cerrado con un porcentaje que no redondea exacto (por ejemplo 1 de 16 puntos, 6,25%) → se redondea a un decimal, con los medios alejándose de cero: `6,3%`, no `6,25%` ni `6,333333333333334%`.
- Proyecto sin ningún Sprint cerrado → lista vacía, igual que HU-12 y HU-22.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| ID de proyecto no numérico o no positivo | — (adaptador HTTP, igual que el resto de las historias) | HTTP 400; no se llama al caso de uso. |

No hay más condiciones de error: es una consulta de solo lectura sin más datos de entrada que el ID del proyecto.

## 8. Criterios de aceptación
- **CA-24.1:** Dado un Sprint `Cerrado` con 8 puntos planificados y 6 completados, cuando se consulta su calidad, entonces el porcentaje es 75; con 16 planificados y 1 completado, el porcentaje es 6,3 (redondeado a un decimal, con los medios alejándose de cero: 6,25 redondea a 6,3, no a 6,2). (RN-24.1)
- **CA-24.2:** Dado un Sprint `Cerrado` con 0 puntos planificados, cuando se consulta su calidad, entonces el porcentaje queda indefinido (no `0`). (RN-24.2)
- **CA-24.3:** Dado un Sprint con dos defectos cuyo `SprintDeteccionID` es el de ese Sprint y uno cuyo `SprintDeteccionID` es de otro Sprint, cuando se consulta su calidad, entonces los defectos detectados de ese Sprint son 2. (RN-24.3)
- **CA-24.4:** Dado un Sprint con un defecto cuyo `SprintResolucionID` es el de ese Sprint y está `Cerrado`, cuando se consulta su calidad, entonces ese defecto cuenta como resuelto en ese Sprint. (RN-24.4)
- **CA-24.5:** Dado un defecto sin `SprintDeteccionID` ni `SprintResolucionID`, cuando se consulta la calidad de cualquier Sprint del proyecto, entonces ese defecto no cuenta ni como detectado ni como resuelto en ninguno. (RN-24.3, RN-24.4)
- **CA-24.6:** Dado un proyecto con Sprints `Abierto` y `Cerrado`, cuando se consulta su calidad, entonces solo aparecen los Sprints `Cerrado`. (RN-24.5)
- **CA-24.7:** Dado un ID de proyecto que no existe, cuando se consulta su calidad, entonces se obtiene `ErrProyectoNoEncontrado`. (HTTP 404)
- **CA-24.8:** Dado un pedido válido a `GET /proyectos/{id}/calidad`, entonces la respuesta es `200` con un array (uno por Sprint cerrado, con `porcentaje_completado` en `null` cuando no hay datos); con un proyecto sin Sprints cerrados, un array vacío; con un ID inválido, `400`; con un proyecto inexistente, `404`.
- **CA-24.9:** Dada la pantalla de sprints anteriores (HU-12) de un proyecto con Sprints cerrados, cuando se abre, entonces cada fila muestra su porcentaje completado (o que no hay datos), sus defectos detectados y sus defectos resueltos, sin romper lo que ya mostraba HU-12/HU-22.

## Decisiones confirmadas
- **P-1 — Pantalla a extender:** se extiende la pantalla de sprints anteriores de HU-12 (solo Sprints `Cerrado`), no la de todos los sprints (HU-08, que incluye Abiertos). Un Sprint abierto no tiene puntos completados reales todavía (RN-11.6), así que mostrar ahí un porcentaje habría exigido otra regla para "sprint en curso, sin cerrar". Confirmada por Matías López (Agile Enabler) el 2026-11-05, en la revisión del PR [HU-24: Calcular % de story points completados y defectos por Sprint #366](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/366). Antes, confirmada con Ezequiel en el chat.
- **P-2 — Defectos sin Sprint asignado:** no se cuentan aparte (no hay una fila o dato "sin sprint"). El AC del issue pide los datos explícitamente "por Sprint"; un defecto sin Sprint de detección o de resolución simplemente no entra en el desglose. Confirmada por Matías López (Agile Enabler) el 2026-11-05, en la revisión del PR #366. Antes, confirmada con Ezequiel en el chat.
- **P-3 — Redondeo del porcentaje:** se redondea a **un decimal**, con los medios alejándose de cero, igual criterio que la desviación de HU-19/HU-23 (`redondearUnDecimal`, RN-19.6). Se calcula en el dominio (`PorcentajeCompletado`, RN-24.1), no en la plantilla, para que la API y la pantalla muestren siempre el mismo valor. Pedido por Matías López (Agile Enabler) el 2026-11-05, en la revisión del PR #366: la versión original (sin redondear) mostraba `66,66666666666666%` en pantalla y `66.66666666666666` en el JSON.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-05: creada (borrador). P-1 a P-3 confirmadas con Ezequiel en el chat, antes de escribir la spec.
- 2026-11-05: Aprobada. En la revisión del PR #366, Matías López confirmó P-1 y P-2, y pidió redondear el porcentaje a un decimal en el dominio (P-3, antes abierta); se corrigió también la sección 6, que daba como caso límite un porcentaje mayor a 100% que en realidad no puede pasar (RN-11.6: `PuntosCompletados` siempre es un subconjunto de `PuntosPlanificados`).
