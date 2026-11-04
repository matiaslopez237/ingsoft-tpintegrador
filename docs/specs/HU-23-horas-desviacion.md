# HU-23: Horas estimadas, reales y desviación de un sprint o proyecto

**Historia:** Como Scrum Master, quiero ver las horas estimadas, las horas reales y su desviación, para evaluar la planificación.
**Issue:** #34  ·  **Sprint:** 3  ·  **Story points:** 3 (según el planning y el Product Backlog; la issue dice 5, ver "Preguntas abiertas")
**Estado de la spec:** Borrador

## 1. Objetivo
Mostrar, para un sprint o para un proyecto completo, cuántas horas se estimaron, cuántas se trabajaron y cuánto se desvió la planificación, en horas y en porcentaje. HU-19 hace esa comparación para un elemento; esta historia suma los elementos y aplica el mismo cálculo (`domain.CompararEsfuerzo`) sobre las sumas. No guarda datos nuevos: las horas estimadas y los registros de esfuerzo ya existen (HU-19 y HU-18).

## 2. Entradas
Son dos consultas sin cuerpo.

**Horas de un sprint.** **En HTTP:** `GET /proyectos/{id}/sprints/{sprintID}/horas`.

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |
| `sprintID` | `int64` | Sí | Entero positivo. Debe corresponder a un sprint de ese proyecto (RN-23.11). |

**Horas de un proyecto.** **En HTTP:** `GET /proyectos/{id}/horas`.

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |

## 3. Salidas esperadas

**Resultado** (`domain.TotalHoras`, el mismo para las dos consultas):

| Campo | Tipo (Go) | Significado |
|---|---|---|
| `HorasEstimadas` | `float64` | Suma de las horas estimadas de los elementos que tienen estimación; 0 si ninguno tiene. |
| `HorasReales` | `float64` | Suma de las horas reales de **todos** los elementos, estimados o no. |
| `HorasRealesSinEstimar` | `float64` | La parte de `HorasReales` que viene de elementos sin estimar. |
| `DesviacionHoras` | `*float64` | Reales de los elementos estimados − estimadas; `nil` si ningún elemento tiene estimación. |
| `DesviacionPorcentaje` | `*float64` | Con un decimal; `nil` si ningún elemento tiene estimación. |
| `Elementos` | `int` | Cantidad de elementos considerados. |
| `ElementosSinEstimar` | `int` | Cuántos de ellos no tienen horas estimadas. |

El resultado tiene además el método `SinEstimar() bool` (ningún elemento con estimación, o ningún elemento), para que la pantalla no repita la regla.

Cuando hay elementos sin estimar, la desviación **no** es `HorasReales − HorasEstimadas`, sino `HorasReales − HorasRealesSinEstimar − HorasEstimadas`. Por eso se informa `HorasRealesSinEstimar`: sin ese dato, los números que ve el usuario no cerrarían (RN-23.5).

| Caso (elementos: estimadas / reales) | Estimadas | Reales | Reales sin estimar | Desviación (h) | Desviación (%) | Elementos / sin estimar |
|---|---|---|---|---|---|---|
| 8 / 10 y 4 / 3 | 12 | 13 | 0 | 1 | 8,3 | 2 / 0 |
| 8 / 6 y 4 / 3 | 12 | 9 | 0 | −3 | −25,0 | 2 / 0 |
| 8 / 10 y — / 5 | 8 | 15 | 5 | 2 | 25,0 | 2 / 1 |
| 8 / 0 y — / 5 | 8 | 5 | 5 | −8 | −100,0 | 2 / 1 |
| — / 5 y — / 3 | 0 | 8 | 8 | — | — | 2 / 2 |
| 8 / 0 | 8 | 0 | 0 | −8 | −100,0 | 1 / 0 |
| Sin elementos | 0 | 0 | 0 | — | — | 0 / 0 |

- La consulta no modifica nada.
- **En HTTP:** `200` con el resultado en JSON (desviaciones en `null` cuando no hay valor); `404` si no existe el proyecto o el sprint; `400` si un ID no es un entero positivo. La respuesta del sprint lleva además `sprint_id`:

```json
{
  "proyecto_id": 1,
  "sprint_id": 2,
  "horas_estimadas": 8,
  "horas_reales": 15,
  "horas_reales_sin_estimar": 5,
  "desviacion_horas": 2,
  "desviacion_porcentaje": 25,
  "elementos": 2,
  "elementos_sin_estimar": 1
}
```

- **Formato de los números:** igual que en HU-19. En JSON son números sin formato; la coma decimal, el signo y el `25,0` son solo de la pantalla.

### Pantalla
Solo lectura, sin formularios. El mismo bloque "Horas y desviación" se muestra en dos pantallas existentes, en las dos después de la lista de datos (`dl.datos`) y antes de la primera sección:

- **Sprint:** en el detalle del sprint de HU-09 (`/web/proyectos/{id}/sprints/{sprintID}`), antes de "Elementos del Sprint Backlog". A esa pantalla se llega con "Ver detalle" desde la lista de sprints (HU-08) y desde sprints anteriores (HU-12).
- **Proyecto:** en el detalle del proyecto (`/web/proyectos/{id}`), antes de "Integrantes".

El bloque muestra:

- La cantidad de elementos, las horas estimadas, las horas reales y las dos desviaciones, con el formato de HU-19: coma decimal, `+2 h` y `+25,0 %` en las positivas, `-2 h` y `-25,0 %` con el guion común en las negativas.
- Con `SinEstimar()`, las dos desviaciones muestran "Sin estimar". Las horas estimadas muestran `0 h` y las reales, su valor.
- Con elementos sin estimar, un aviso: "1 de 2 elementos sin estimar: sus 5 h están en las horas reales, pero no en la desviación." Si esos elementos no tienen horas reales, solo los cuenta: "1 de 2 elementos sin estimar."
- Sin elementos, además de lo anterior: "Todavía no hay elementos." (`0 h` en estimadas y reales).
- **Solo en un sprint `Cerrado`,** una aclaración: "Solo cuentan los elementos que siguen en este sprint. Los que volvieron al backlog al cerrarlo no suman acá, pero sí en el total del proyecto." (RN-23.8)

El bloque también se muestra cuando estas pantallas vuelven con el error de un formulario (`mostrarSprintDetalle`, `mostrarDetalle`).

## 4. Reglas de negocio
- **RN-23.1:** Las horas se consultan por sprint o por proyecto. Cada consulta informa horas estimadas, horas reales, la desviación en horas y en porcentaje, y cuántos elementos se consideraron y cuántos están sin estimar.
- **RN-23.2:** Las horas estimadas son la suma de las horas estimadas (HU-19) de los elementos considerados. Un elemento sin estimar no suma, aunque tenga story points (RN-19.1).
- **RN-23.3:** Las horas reales son la suma de los registros de esfuerzo de **todos** los elementos considerados, estimados o no, de cualquier integrante y fecha, incluidos los de integrantes que después se quitaron del proyecto (RN-19.4).
- **RN-23.4:** La desviación se calcula una sola vez sobre las sumas, con `domain.CompararEsfuerzo`: mismas fórmulas, mismo signo y mismo redondeo a un decimal que RN-19.5 y RN-19.6. No se promedian ni se suman las desviaciones de cada elemento. Ejemplo: 10 / 20 (+100 %) y 90 / 90 (0 %) dan 100 / 110, es decir, +10,0 %, no +50 %.
- **RN-23.5:** La desviación compara solo a los elementos que tienen estimación: a `CompararEsfuerzo` se le pasan la suma de sus estimadas y la suma de **sus** reales. Las horas de los elementos sin estimar se muestran en el total de reales y aparte en `HorasRealesSinEstimar`, pero no entran en la desviación: contarlas compararía trabajo sin estimación contra una estimación que no lo incluye (mismo criterio que RN-19.7).
- **RN-23.6:** Si no hay ningún elemento con estimación (porque no hay elementos o porque ninguno está estimado), no es un error: `HorasEstimadas` es 0, `HorasReales` conserva su valor y no hay desviación (`SinEstimar()` verdadero). A `CompararEsfuerzo` se le pasa `nil`, no 0: pasarle 0 mostraría todo lo trabajado como desviación.
- **RN-23.7:** Los elementos de un sprint son los que tienen hoy ese `SprintID`, en cualquier estado.
- **RN-23.8:** Como el cierre deja sin sprint a los elementos sin terminar (RN-11.3), un elemento devuelto al backlog deja de contar para el sprint cerrado, incluidas las horas que se trabajaron durante ese sprint. Sus horas siguen contando en el proyecto (RN-23.10). Si se lo planifica en otro sprint (HU-09), **todas** sus horas pasan a contar ahí. Lo mismo pasa con un elemento de un sprint cerrado que se reasigna a otro sprint, algo que RN-09.2 permite. Es una limitación conocida (P-3) y la pantalla la aclara.
- **RN-23.9:** Los totales se calculan al consultar, con los valores actuales; no se congelan al cerrar el sprint, a diferencia de los puntos de RN-11.6. Si después del cierre se cambia una estimación (RN-19.3) o se registra esfuerzo (RN-18.4 no lo impide), el total del sprint cambia. No se mira la fecha de los registros: cuenta un registro con fecha fuera de las fechas del sprint si es de uno de sus elementos.
- **RN-23.10:** Los elementos de un proyecto son todos los suyos, estén o no en un sprint. Por RN-23.7, la suma de los sprints de un proyecto puede ser menor que el total del proyecto, y eso no es un error.
- **RN-23.11:** El proyecto tiene que existir y el sprint tiene que ser de ese proyecto; un sprint de otro proyecto se trata como inexistente (igual que RN-11.4).
- **RN-23.12:** Se puede consultar con el sprint `Abierto` o `Cerrado` y con el proyecto en cualquier estado (como RN-19.11). En elementos sin terminar, la desviación muestra cómo va hasta ahora: un sprint recién empezado puede dar cerca de −100 %.

## 5. Restricciones
- **Cálculo en el dominio:** `TotalizarHoras(elementos []HorasElemento) TotalHoras` vive en `internal/domain/total_horas.go`, sin dependencias externas. `HorasElemento` tiene `HorasEstimadas *float64` y `HorasReales float64`. Implementa RN-23.2 a RN-23.6 y llama a `CompararEsfuerzo`, que no cambia. La suma no se hace en SQL porque es regla de negocio y tiene que poder probarse con TDD. No modifica los elementos que recibe.
- **Precisión:** los valores vienen de medias horas enteras (`medias_horas`, `medias_horas_estimadas`), así que todos son múltiplos de 0,5, que `float64` representa de forma exacta: sumarlos no acumula error. La suma no tiene tope; el de 1000 h es por elemento (RN-19.2) y el de 24 h, por registro (RN-18.8).
- **Puerto de lectura:** los casos de uso definen `LectorHoras`, angosto, sin ensanchar `RepositorioEsfuerzo` ni `LectorHorasReales`, por el mismo motivo que HU-19:

  ```go
  // Devuelven una fila por elemento; un elemento sin registros viene con 0 horas reales.
  HorasDeProyecto(ctx context.Context, proyectoID int64) ([]domain.HorasElemento, error)
  HorasDeSprint(ctx context.Context, sprintID int64) ([]domain.HorasElemento, error)
  ```

  Lo implementa `persistencia.RepositorioEsfuerzo` con una sola consulta, con `LEFT JOIN` a `registros_esfuerzo` y `GROUP BY` por elemento. Tiene que ser `LEFT JOIN`: si no, un elemento estimado y sin esfuerzo no sumaría a las estimadas. Se usa una consulta por conjunto y no `HorasRealesDeElemento` por cada elemento, para no hacer N consultas. Las filas vienen ordenadas por el id del elemento: el total no depende del orden, pero así el resultado es estable.
- **Casos de uso:** `ConsultarHorasProyecto` y `ConsultarHorasSprint` (`internal/app/consultar_horas.go`). Reutilizan `LectorProyectos` (HU-04). Para el sprint, un puerto de un método con `BuscarSprintPorID`, que ya implementa `persistencia.RepositorioSprints` (HU-11) y devuelve `ErrSprintNoEncontrado`. Primero buscan el proyecto, después el sprint.
- **Contexto:** los handlers JSON pasan a los casos de uso el contexto del pedido, y los casos de uso se lo pasan al lector de horas, para que cancelar el pedido corte la consulta.
- **Sin persistencia nueva:** no hay tablas, columnas ni migraciones.
- **Rutas y pantallas:** `net/http` (ADR-0001), SQLite (ADR-0002) y `html/template` (ADR-0003). Las dos rutas JSON se registran con `RegistrarConsultarHoras` en `consultar_horas_handler.go` y se suman a `NuevoRouter`, a `CasosDeUso` y a `router_test.go`. Los dobles de `pantallas_test.go` también tienen que sumar los dos casos de uso, porque las dos pantallas los consultan siempre. El formato reutiliza `textoHoras`, `textoPorcentaje`, `conSigno` y `textoSinEstimar` de HU-19. La pantalla no calcula nada.
- **Dependencia:** usa `CompararEsfuerzo` y `HorasEstimadas` de HU-19, que se mergeó a `main` (PR #297). Esta rama salió de `hu-19-comparar-esfuerzo` y ya se integró con `main`.
- **Coordinación:** HU-22 (rama `hu-22-metricas-sprint`) usa la ruta `GET /proyectos/{id}/velocidad` y la pantalla de sprints anteriores, así que no se pisa con estas rutas ni con estas pantallas.
- **Fuera de alcance:**
  - **Para otras historias:** gráficos del dashboard (HU-25) y reportes y PDF (HU-26).
  - **Sin historia:** detalle de horas por elemento o por integrante dentro del total, totales en la tabla de la lista de sprints, congelar las horas al cerrar un sprint, historial de asignaciones de un elemento y controlar quién puede consultar.

## 6. Casos límite
- Sprint sin elementos → totales en 0, sin desviación, no es un error (RN-23.6).
- Proyecto sin elementos → igual que el sprint sin elementos.
- Todos los elementos sin estimar y con esfuerzo → estimadas 0, reales con su valor, sin desviación (RN-23.6).
- Elementos estimados sin esfuerzo y uno sin estimar con esfuerzo (8 / 0 y — / 5) → reales 5, desviación −8 h y −100,0 %: las 5 h no compensan la desviación (RN-23.5).
- Suma estimada 0 con elementos estimados → no puede pasar: cada estimación es mayor a 0 (RN-19.2). La rama "estimado 0" de `CompararEsfuerzo` (RN-19.8) y el texto "No aplica" no se alcanzan desde esta historia.
- Sumas iguales con desviaciones opuestas (8 / 10 y 8 / 6) → 16 / 16, 0 h y 0,0 %: los desvíos se compensan. Es lo que pide RN-23.4.
- Redondeo sobre la suma: 1 / 2 y 2 / 2 → 3 / 4, 33,3 %. 20 / 20,5 y 20 / 20 → 40 / 40,5, 1,3 %. 20 / 19,5 y 20 / 20 → 40 / 39,5, −1,3 %.
- Un solo elemento → el mismo resultado que HU-19 para ese elemento.
- Sprint cerrado con un elemento `Hecho` y otro devuelto al backlog → solo cuenta el `Hecho` en el sprint; los dos cuentan en el proyecto (RN-23.8, RN-23.10).
- Elemento devuelto al backlog y replanificado en otro sprint → todas sus horas cuentan en el sprint nuevo (RN-23.8).
- Elemento `Hecho` de un sprint cerrado reasignado a un sprint abierto (RN-09.2) → deja de contar en el cerrado y cuenta en el abierto (RN-23.8).
- Elemento de un sprint cerrado que pasa a `Pendiente` editándolo (RN-06.4) → sigue contando: lo que importa es el `SprintID` (RN-23.7).
- Se quita un elemento de un sprint abierto (RN-09.4) → deja de contar en ese sprint.
- Se elimina un elemento sin esfuerzo (HU-06) → deja de sumar sus estimadas. Uno con esfuerzo no se puede eliminar (RN-06.7).
- Registro de esfuerzo con fecha fuera de las del sprint → cuenta (RN-23.9).
- Registros de un integrante que ya no está en el proyecto → cuentan (RN-23.3).
- Elemento con story points y sin horas estimadas → sin estimar (RN-23.2).
- Sprint de otro proyecto → sprint no encontrado (RN-23.11).
- ID `0`, negativo, `abc` o `anteriores` en la ruta de horas del sprint → `400` (sección 7).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| ID de proyecto que no es un entero positivo | — (adaptador HTTP) | "el id del proyecto debe ser un número entero positivo" (HTTP 400); no se llama al caso de uso. |
| ID de sprint que no es un entero positivo | — (adaptador HTTP) | "el id del sprint debe ser un número entero positivo" (HTTP 400); no se llama al caso de uso. |
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el sprint, o es de otro proyecto | `domain.ErrSprintNoEncontrado` (HU-09, se reutiliza) | "sprint no encontrado" (HTTP 404) |
| Falla interna al consultar | — (se envuelve con `%w` la falla al leer las horas; las de las búsquedas salen tal cual, como en HU-19) | HTTP 500 con un mensaje genérico, sin el detalle. |

Sprint vacío, sin estimar y sin esfuerzo **no son errores** (RN-23.6).

Si hay varios errores, se informa el primero en el orden de la tabla. No se usa `ErrSprintDeOtroProyecto` de HU-09: acá un sprint de otro proyecto es inexistente (RN-23.11). En la pantalla, los `400` y los `404` muestran la página de error, y una falla interna, la página de falla sin el detalle.

## 8. Criterios de aceptación
- **CA-23.1 (sprint):** Dado un sprint con un elemento de 8 h estimadas y 10 h reales y otro de 4 h estimadas y 3 h reales, cuando consulto sus horas, entonces obtengo 12 estimadas, 13 reales, desviación de 1 h y de 8,3 %, 2 elementos y 0 sin estimar. (RN-23.1, RN-23.2, RN-23.3, RN-23.4, RN-23.7)
- **CA-23.2 (proyecto):** Dado un proyecto con elementos 8 / 10 y 4 / 3 en un sprint, y 20 / 15 sin sprint, cuando consulto sus horas, entonces obtengo 32 estimadas, 28 reales, desviación de −4 h y de −12,5 %, y 3 elementos. (RN-23.1, RN-23.4, RN-23.10)
- **CA-23.3 (sobre las sumas):** Dados elementos 10 / 20 y 90 / 90, cuando consulto, entonces la desviación es de 10 h y de 10,0 %, y no 50,0 % como daría el promedio. (RN-23.4)
- **CA-23.4 (redondeo):** Dados elementos 1 / 2 y 2 / 2, el porcentaje es 33,3; con 20 / 20,5 y 20 / 20, es 1,3; con 20 / 19,5 y 20 / 20, es −1,3. (RN-23.4)
- **CA-23.5 (desvíos que se compensan):** Dados elementos 8 / 10 y 8 / 6, cuando consulto, entonces la desviación es de 0 h y de 0,0 %. (RN-23.4)
- **CA-23.6 (elementos sin estimar):** Dados elementos 8 / 10 y uno sin estimar con 5 h reales, cuando consulto, entonces obtengo 8 estimadas, 15 reales, 5 reales sin estimar, desviación de 2 h y de 25,0 %, 2 elementos y 1 sin estimar. (RN-23.2, RN-23.3, RN-23.5)
- **CA-23.7 (ninguno estimado):** Dado un sprint con dos elementos sin estimar con 5 h y 3 h reales, cuando consulto, entonces obtengo 0 estimadas, 8 reales, sin desviaciones y `SinEstimar()` verdadero. No es un error. (RN-23.6)
- **CA-23.8 (sin elementos):** Dado un sprint o un proyecto sin elementos, cuando consulto, entonces obtengo todo en 0, sin desviaciones y `SinEstimar()` verdadero. No es un error. (RN-23.6)
- **CA-23.9 (sin esfuerzo):** Dado un sprint con un elemento 8 / 0, cuando consulto, entonces obtengo 8 estimadas, 0 reales, desviación de −8 h y de −100,0 %. Un elemento estimado y sin registros suma a las estimadas. (RN-23.2, RN-23.3)
- **CA-23.10 (cierre del sprint):** Dado un sprint con un elemento `Hecho` 8 / 10 y uno `En progreso` 4 / 3, cuando lo cierro y consulto, entonces el sprint muestra 8 / 10 y el proyecto sigue mostrando 12 / 13; cuando planifico el elemento devuelto en otro sprint, ese sprint muestra 4 / 3; cuando reasigno el elemento `Hecho` a ese mismo sprint, el sprint cerrado muestra 0 / 0 y el otro, 12 / 13. (RN-23.7, RN-23.8, RN-23.10)
- **CA-23.11 (valores actuales):** Dado un sprint cerrado con un elemento 8 / 10, cuando cambio su estimación a 10 y registro 2 h más, entonces el sprint muestra 10 / 12, y sus puntos planificados y completados (RN-11.6) no cambian. (RN-23.9)
- **CA-23.12 (qué suma):** Dado un sprint con un elemento que tiene registros de un integrante que después se quitó del proyecto y un registro con fecha fuera de las del sprint, y dados elementos de otro sprint y de otro proyecto, cuando consulto, entonces suman los registros de ese elemento y no los de los demás. (RN-23.3, RN-23.7, RN-23.9)
- **CA-23.13 (cualquier estado):** Dado un sprint `Abierto` con elementos `Pendiente`, en un proyecto que no está Activo, cuando consulto, entonces obtengo sus horas y no cambia nada. (RN-23.12)
- **CA-23.14 (proyecto inexistente):** Dado un ID de proyecto que no existe, cuando consulto las horas del proyecto o de un sprint, entonces obtengo `ErrProyectoNoEncontrado`. (RN-23.11)
- **CA-23.15 (sprint inexistente o de otro proyecto):** Dado un ID de sprint que no existe, o un sprint de otro proyecto, cuando consulto, entonces obtengo `ErrSprintNoEncontrado`. (RN-23.11)
- **CA-23.16 (HTTP):** Dado un pedido a cualquiera de las dos rutas, la respuesta es `200` con los campos de la sección 3 y `null` en las desviaciones que no existen; con proyecto o sprint inexistente, `404`; con un ID que no es un entero positivo, `400`, sin llamar al caso de uso; con una falla interna, `500` sin el detalle. (Sección 3, sección 7)
- **CA-23.17 (pantalla del sprint):** Dado el detalle de un sprint, entonces muestra el bloque antes de "Elementos del Sprint Backlog", con los formatos de la sección 3 y "Sin estimar" en las desviaciones cuando no hay estimación; con elementos sin estimar, muestra el aviso; solo con el sprint `Cerrado`, muestra la aclaración de RN-23.8. El bloque sigue visible cuando la pantalla vuelve con un error. (Sección 3)
- **CA-23.18 (pantalla del proyecto):** Dado el detalle de un proyecto, entonces muestra el bloque antes de "Integrantes", con los mismos formatos y sin la aclaración de RN-23.8. El bloque sigue visible cuando la pantalla vuelve con un error. (Sección 3)

Los escenarios mínimos corresponden a: feliz por sprint, CA-23.1; feliz por proyecto, CA-23.2; sobre las sumas, CA-23.3; sin estimar, CA-23.6; sin elementos, CA-23.8; cierre, CA-23.10; errores, CA-23.14 y CA-23.15.

## Decisiones tomadas
Tomadas por el autor (Lucas Orellana) el 2026-11-03, abiertas a objeción en el PR:

- **P-1 — Elementos sin horas estimadas:** no suman a las estimadas, sus horas reales sí suman al total de reales y la desviación se calcula solo sobre los elementos estimados; la respuesta informa cuántos quedaron sin estimar (RN-23.2, RN-23.3, RN-23.5). **Ajuste a la propuesta original:** se suma `HorasRealesSinEstimar`. Sin ese dato, con 8 / 10 y — / 5 la pantalla diría estimadas 8, reales 15 y desviación +2, y parecería un error de cuentas.
- **P-2 — Sprint vacío o sin horas estimadas:** no es un error; estimadas en 0 y desviación "sin estimar" (RN-23.6). **Ajuste a la propuesta original:** "totales en 0" vale para las estimadas y para el sprint vacío. Si ningún elemento está estimado pero hay esfuerzo, las reales conservan su valor, para no esconder horas trabajadas. **Diferencia con HU-19:** `horas_estimadas` sale como 0 y no como `null`, porque es una suma; lo que indica "sin estimar" es el `null` de las desviaciones.
- **P-3 — Qué elementos cuentan en un sprint:** los que tienen hoy ese `SprintID` (RN-23.7, opción A). Responde la pregunta del planning: las horas de un elemento devuelto al backlog **no** cuentan para el sprint en que se trabajaron, porque hoy el sistema no guarda ese dato. `registros_esfuerzo` no tiene sprint y el cierre pone `sprint_id = NULL` en los elementos sin terminar (`cierre_repo.go`).

  | Opción | Ventajas | Desventajas |
  |---|---|---|
  | **A. Contar por el `SprintID` actual** (elegida) | No toca HU-09, HU-11 ni la base. Un sprint cerrado queda con el trabajo terminado, así que su desviación compara trabajo completo. | Las horas de un elemento devuelto salen de ese sprint y pasan enteras al siguiente. Los totales no quedan fijos (RN-23.9). |
  | B. Guardar los totales al cerrar, como RN-11.6 | El sprint cerrado queda fijo e incluye lo no terminado. | Cambia la transacción de cierre de HU-11 y suma columnas. Los sprints ya cerrados no tendrían datos. Un sprint abierto se seguiría calculando con otro criterio. |
  | C. Atribuir las reales por la fecha del registro | Las horas quedan en el sprint en que se trabajaron. | Las estimadas no se pueden repartir por fecha, así que se compararía contra otro conjunto. Los registros fuera de todo sprint quedarían sin atribuir (RN-18.4 no valida fechas). |
  | D. Guardar el historial de asignaciones | Es la respuesta más completa. | Tabla nueva y cambios en HU-09 y HU-11; no entra en una historia de 3 puntos. |

  La limitación se documenta en RN-23.8 y se muestra en pantalla. No se pierden horas: el total del proyecto las incluye siempre. **Consecuencia:** los puntos planificados de un sprint cerrado (RN-11.6) incluyen lo no terminado y estas horas no, así que no se comparan entre sí.
- **P-4 — Proyecto:** suma todos los elementos, estén o no en un sprint (RN-23.10).
- **P-5 — Permisos:** igual que HU-19. **Ajuste a la propuesta original:** se proponía "miembros del proyecto", pero HU-19 no verifica eso (sección 5 de HU-19, `comparar_esfuerzo.go`) y el sistema no tiene sesión ni usuario, así que no hay a quién comprobar como miembro. Por lo tanto, cualquiera que llegue a la ruta puede consultar. No se modela el rol Scrum Master.
- **P-6 — Pantalla:** se incluye, de solo lectura, en el detalle del sprint y en el del proyecto (sección 3). Se eligió el detalle del sprint y no la tabla de `sprints.html` porque la tabla necesitaría una consulta por fila y siete datos más por sprint, y al detalle ya se llega desde la lista de sprints y desde sprints anteriores.
- **P-7 — Horas:** reales en medias horas enteras y estimadas con paso de 0,5 h, como en HU-18 y HU-19. La suma es exacta (sección 5).
- **P-8 — Reutilización del dominio:** `CompararEsfuerzo` no cambia. `TotalHoras` es un tipo nuevo porque `ComparacionEsfuerzo.HorasReales` serían solo las de los elementos estimados y no alcanza para informar el total ni los elementos sin estimar.

## Preguntas abiertas
- **P-9 — Story points: ¿3 o 5?** La issue #34 dice 5. El planning del Sprint 3 dice 3, y el Product Backlog tiene 5 como estimación inicial y 3 en la columna de puntos, así que parece que la issue quedó desactualizada. La spec usa 3. **Falta que el equipo lo confirme y actualice la issue**: cambia el compromiso del sprint, que suma 28 con 3 puntos.
- **P-10 — ¿Alcanza la opción A de P-3 para el equipo?** La pregunta la dejó abierta el planning. Si alguien necesita que un sprint cerrado muestre las horas de lo que no se terminó, hay que elegir B, C o D, que tocan HU-09 o HU-11 y agrandan la historia.
- **P-11 — Texto de HU-19 a corregir.** RN-19.8 y P-4 de HU-19 dicen que HU-23 le va a pasar sumas que pueden dar 0. Con RN-23.6 no pasa: sin estimaciones se pasa `nil`, y con estimaciones la suma es mayor a 0. La rama de estimado 0 sigue siendo correcta y queda como defensa. **Propuesta:** corregir esa frase de HU-19 cuando se apruebe esta spec; no se toca en este cambio.
- **P-12 — ¿Mostrar los totales también en la lista de sprints?** Quedan fuera por P-6. Si el Scrum Master quiere comparar sprints de un vistazo, podría resolverlo el dashboard (HU-25); hay que definirlo en esa historia.

## Historial
- 2026-11-03: creada (borrador). Decisiones P-1 a P-8 tomadas por el autor, abiertas a objeción en el PR; P-9 a P-12, pendientes del equipo.
- 2026-11-03: ajustada a la pantalla implementada, sin cambiar reglas ni criterios: el bloque se llama "Horas y desviación" y muestra también la cantidad de elementos. La spec sigue en borrador.
- 2026-11-03: el aviso de elementos sin estimar y sin horas reales solo los cuenta, sin hablar de 0 h. La spec sigue en borrador.
- 2026-11-03: se documentan en la sección 5 tres cosas que el código ya hacía y los tests verifican: el orden de las filas, el paso del contexto y que `TotalizarHoras` no modifica lo que recibe. La spec sigue en borrador.
- 2026-11-04: se actualiza la dependencia: HU-19 se mergeó a `main` (PR #297) y HU-23 ya se integró con `main`. No cambia ninguna regla, criterio ni decisión. La spec sigue en borrador.
