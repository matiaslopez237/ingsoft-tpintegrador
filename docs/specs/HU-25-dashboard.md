# HU-25: Ver panel con estado del proyecto y gráficos

**Historia:** Como cualquier rol del equipo, quiero un Dashboard con el estado del proyecto y gráficos de las métricas principales, para tener una visión rápida del avance.
**Issue:** #36  ·  **Sprint:** 4  ·  **Story points:** 8
**Estado de la spec:** Aprobada

## 1. Objetivo
Juntar en una sola pantalla lo que hoy está repartido entre el backlog, los sprints, los sprints anteriores, los defectos y las horas, para que cualquier integrante vea de un vistazo cómo viene el proyecto: qué hay en el backlog, cómo va el sprint en curso y cómo rindieron los sprints cerrados.

**Alcance real de esta historia:** todas las métricas ya existen (HU-11, HU-22, HU-23 y HU-24). HU-25 no calcula métricas nuevas: las junta, cuenta el backlog por estado y resume el sprint en curso con la regla que ya usa el cierre (RN-11.6), y las muestra como gráficos. **No agrega persistencia nueva** (la tarea 25.7 no aplica).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `hoy` | `time.Time` | Sí (lo da el reloj, no el usuario) | Fecha del pedido, para saber qué sprint está en curso (RN-25.3). Se inyecta igual que en HU-18 y HU-20, para fijarla en los tests. |

Es una consulta de solo lectura. En HTTP, `proyectoID` va en la ruta.

## 3. Salidas esperadas
**Estado del proyecto:**
- **Backlog por estado:** cuántos elementos del proyecto están `Pendiente`, `En progreso` y `Hecho`, y el total (RN-25.2).
- **Sprint en curso:** si hay uno (RN-25.3), su número, objetivo, fechas, y sus puntos planificados y completados hasta hoy, calculados con la misma regla del cierre (RN-25.4). Si no hay, se informa que no hay sprint en curso.

**Gráficos** (RN-25.1, RN-25.5):
1. **Puntos por sprint:** barras de puntos planificados y completados de cada Sprint `Cerrado`, y la velocidad del equipo (HU-22).
2. **% completado por sprint:** una barra por Sprint `Cerrado` con su porcentaje (HU-24). Un sprint sin puntos planificados figura como "sin datos", sin barra (RN-24.2).
3. **Defectos por sprint:** barras de defectos detectados y resueltos de cada Sprint `Cerrado` (HU-24).
4. **Horas estimadas contra reales** del proyecto, con la desviación (HU-23, RN-23.4 a RN-23.6). Si ningún elemento está estimado en horas, se informa "sin datos" en lugar del gráfico (RN-23.6). Si hay elementos sin estimar, debajo del gráfico va el mismo aviso que el bloque de horas de HU-23: cuántos elementos están sin estimar y que sus horas reales están en la barra de reales pero no en la desviación (RN-23.5).

**En la pantalla** (`GET /web/proyectos/{id}/dashboard`): las dos secciones de estado arriba y los 4 gráficos abajo. Cada gráfico es un SVG armado en el servidor (ADR-0005) con un título accesible y una tabla con los mismos datos. La pantalla del proyecto (HU-04) suma un botón "Dashboard", el primero de la fila (RN-25.6).

**En HTTP** (`GET /proyectos/{id}/dashboard`): `200` con un objeto JSON con `backlog` (conteos por estado y total), `sprint_en_curso` (`null` si no hay), `sprints_cerrados` (array, con puntos, porcentaje y defectos de cada uno), `velocidad` (`null` si no hay sprints cerrados) y `horas` (con `desviacion_porcentaje` en `null` si no hay datos); `400` si el ID no es un número; `404` si el proyecto no existe.

## 4. Reglas de negocio
- **RN-25.1:** El Dashboard muestra los datos vigentes del proyecto en el momento del pedido: cada vez que se abre, se calculan de nuevo a partir de lo guardado. No hay datos guardados aparte para el Dashboard ni cache.
- **RN-25.2:** Los elementos del backlog se cuentan por su estado actual (`Pendiente`, `En progreso`, `Hecho`, HU-05/HU-10), estén o no asignados a un sprint. El total es la suma de los tres.
- **RN-25.3:** El sprint en curso es el Sprint `Abierto` cuyas fechas incluyen `hoy`, con la misma regla que usa HU-20 para el sprint de detección de un defecto (`domain.SprintVigente`, RN-20.5). Como los sprints de un proyecto no se superponen (RN-08.5), hay como mucho uno. Un sprint `Abierto` cuyas fechas no incluyen hoy (ya vencido o todavía no empezado) no es el sprint en curso (P-1).
- **RN-25.4:** Los puntos del sprint en curso se calculan con la regla del cierre de HU-11 (`domain.ResumirCierre`, RN-11.6) sobre los elementos que tiene asignados hoy: planificados es la suma de los story points de todos, y completados la de los que están `Hecho`; un elemento sin estimar suma 0. No se guardan: se recalculan en cada pedido (RN-25.1).
- **RN-25.5:** Los gráficos por sprint incluyen todos los Sprints `Cerrado` del proyecto, del más antiguo al más reciente por fecha de inicio (de izquierda a derecha), sin límite de cantidad (P-2), y usan los valores que ya calculan HU-11, HU-22 y HU-24 sin recalcularlos.
- **RN-25.6:** Al Dashboard se llega con un solo clic desde la pantalla del proyecto (HU-04), que es la vista principal de cada proyecto.

## 5. Restricciones
- El cálculo vive en el núcleo: el conteo del backlog y el resumen del sprint en curso son funciones del dominio; el caso de uso `ConsultarDashboard` (`internal/app/consultar_dashboard.go`) junta los datos.
- `ConsultarDashboard` reutiliza los puertos que ya existen: `LectorProyectos`, `RepositorioSprints` (HU-08), `LectorDefectos` (HU-21), el listado del backlog (HU-07) y `LectorHoras` (HU-23). **No agrega ningún puerto de persistencia nuevo** ni cambia el esquema de SQLite.
- Los gráficos se dibujan como SVG en el servidor, sin JavaScript, respetando la CSP de las pantallas (ADR-0005). La escala de las barras es lógica de presentación: vive en el adaptador HTTP, con sus propios tests.
- Los colores no son la única forma de distinguir las series: cada gráfico tiene leyenda con texto y una tabla con los datos.

## 6. Casos límite
- Proyecto recién creado, sin elementos ni sprints → backlog con todo en 0, sin sprint en curso, los gráficos por sprint con el mensaje "todavía no hay sprints cerrados" y las horas "sin datos". La pantalla no falla.
- Proyecto con sprints pero ninguno `Cerrado` → los gráficos por sprint muestran "todavía no hay sprints cerrados"; la velocidad, "sin datos" (RN-22.3).
- Sprint `Abierto` cuyas fechas ya pasaron (no se cerró a tiempo) o todavía no empezaron → no es el sprint en curso (RN-25.3).
- Sprint en curso sin elementos asignados → se muestra con 0 planificados y 0 completados.
- Sprint en curso con elementos sin estimar → suman 0 puntos (RN-11.6).
- Sprint cerrado con 0 puntos planificados → en el gráfico de porcentaje figura "sin datos", sin barra; en el de puntos, con dos barras en 0 (RN-24.2).
- Todos los sprints cerrados con 0 completados y 0 defectos → los gráficos se dibujan con barras en 0, sin dividir por cero al calcular la escala.
- Muchos sprints cerrados → todos entran en los gráficos, con barras más finas (P-2).
- Elementos o sprints de otro proyecto → no cuentan.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404; en la pantalla, la página de error) |
| ID de proyecto no numérico o no positivo | — (adaptador HTTP, igual que el resto de las historias) | HTTP 400; no se llama al caso de uso. |

No hay más condiciones de error: es una consulta de solo lectura sin más datos de entrada que el ID del proyecto.

## 8. Criterios de aceptación
- **CA-25.1:** Dado un proyecto con 2 elementos `Pendiente`, 1 `En progreso` y 3 `Hecho`, cuando se consulta su Dashboard, entonces el backlog informa 2, 1 y 3, con total 6. (RN-25.2)
- **CA-25.2:** Dado un Sprint `Abierto` cuyas fechas incluyen hoy, con elementos de 3, 5 y 2 puntos y solo el de 5 `Hecho`, cuando se consulta el Dashboard, entonces el sprint en curso es ese, con 10 puntos planificados y 5 completados. (RN-25.3, RN-25.4)
- **CA-25.3:** Dado un proyecto sin ningún Sprint `Abierto` que incluya hoy (sin sprints, o con uno abierto ya vencido), cuando se consulta el Dashboard, entonces informa que no hay sprint en curso. (RN-25.3)
- **CA-25.4:** Dado un proyecto con dos Sprints `Cerrado`, cuando se consulta el Dashboard, entonces los datos por sprint traen los dos, del más antiguo al más reciente, con sus puntos planificados y completados, su porcentaje y sus defectos detectados y resueltos, iguales a los de HU-12 y HU-24; y la velocidad es la de HU-22. (RN-25.5)
- **CA-25.5:** Dado un proyecto con horas estimadas y reales cargadas, cuando se consulta el Dashboard, entonces las horas y la desviación son las de HU-23 para el proyecto; si ningún elemento está estimado, la desviación queda sin datos; y si hay elementos sin estimar, la pantalla muestra el aviso de HU-23 que explica que sus horas reales no entran en la desviación. (RN-23.4 a RN-23.6, de HU-23)
- **CA-25.6:** Dado que se abrió el Dashboard, cuando se completa una historia del sprint en curso y se vuelve a abrir, entonces los puntos completados ya incluyen esa historia. (RN-25.1)
- **CA-25.7:** Dado un ID de proyecto que no existe, cuando se consulta el Dashboard, entonces se obtiene `ErrProyectoNoEncontrado`. (HTTP 404)
- **CA-25.8:** Dado un pedido válido a `GET /proyectos/{id}/dashboard`, entonces la respuesta es `200` con el objeto de la sección 3 (`sprint_en_curso` y `velocidad` en `null` cuando no hay datos); con un ID inválido, `400`; con un proyecto inexistente, `404`.
- **CA-25.9:** Dada la pantalla del Dashboard de un proyecto con sprints cerrados, cuando se abre, entonces muestra al menos 3 gráficos SVG (puntos, porcentaje y defectos por sprint; y el de horas si hay datos), cada uno con un título accesible y una tabla con sus datos, sin scripts. (RN-25.5, ADR-0005)
- **CA-25.10:** Dada la pantalla de un proyecto (HU-04), cuando se abre, entonces tiene un enlace "Dashboard" a `/web/proyectos/{id}/dashboard`. (RN-25.6)
- **CA-25.11:** Dado un proyecto recién creado, sin elementos ni sprints, cuando se abre su Dashboard, entonces responde `200` con el backlog en 0, sin sprint en curso, y en lugar de los gráficos el aviso de que no hay datos todavía. (sección 6)

## Decisiones confirmadas
- **Gráficos como SVG armado en Go**, sin JavaScript (ADR-0005). Propuesto por Nicolás Andino el 2026-11-06; aprobado por Matías López (Agile Enabler) el 2026-11-07, en la revisión del PR #375.
- **Métricas graficadas:** puntos por sprint (con la velocidad), % completado por sprint, defectos por sprint y horas estimadas contra reales. Elegidas por Nicolás Andino el 2026-11-06.
- **Pantalla propia** (`/web/proyectos/{id}/dashboard`) con un botón en la pantalla del proyecto, en lugar de mostrar el panel dentro de HU-04. Elegido por Nicolás Andino el 2026-11-06.
- **Estado del proyecto:** backlog por estado y sprint en curso. Elegido por Nicolás Andino el 2026-11-06.
- **P-1 — Sprint abierto pero vencido:** no figura como sprint en curso; la sección dice "No hay sprint en curso". Es la misma regla de HU-20 (`domain.SprintVigente`). Propuesta por Nicolás Andino el 2026-11-06; aprobada por Matías López (Agile Enabler) el 2026-11-07, en la revisión del PR #375.
- **P-2 — Cantidad de sprints en los gráficos:** todos los sprints cerrados, sin límite. Propuesta por Nicolás Andino el 2026-11-06; aprobada por Matías López (Agile Enabler) el 2026-11-07, en la revisión del PR #375.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-06: creada (borrador). P-1 y P-2 confirmadas por Nicolás Andino; falta la aprobación del equipo en el PR.
- 2026-11-07: Aprobada. En la revisión del PR #375, Matías López (Agile Enabler) aprobó ADR-0005, P-1 y P-2, y pidió mostrar en el gráfico de horas el aviso de elementos sin estimar de HU-23 (sección 3 y CA-25.5). CA-25.5 pasa a citar las reglas de HU-23 (RN-23.4 a RN-23.6) en lugar de RN-25.5.
