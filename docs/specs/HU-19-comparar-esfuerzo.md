# HU-19: Comparar esfuerzo estimado contra real

**Historia:** Como Scrum Master, quiero comparar el esfuerzo estimado contra el real de una historia, para evaluar la precisión de las estimaciones.
**Issue:** #30  ·  **Sprint:** 3  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Mostrar, para un elemento del backlog, cuántas horas se estimaron, cuántas se trabajaron de verdad y cuánto se desvió la estimación, en horas y en porcentaje. Hoy el sistema tiene las horas reales (HU-18) pero no las estimadas: esta historia también define de dónde salen (P-10). El alcance es **un elemento**; los totales por sprint o por proyecto son de HU-23, que reutiliza el cálculo de desviación que esta historia deja en el dominio.

## 2. Entradas

**Comparar** (consulta, sin cuerpo). **En HTTP:** `GET /proyectos/{id}/backlog/{elementoID}/esfuerzo/comparacion`.

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento de backlog (HU-05) de ese proyecto (RN-19.10). |

**Cargar las horas estimadas** (se suma a la edición de HU-06, según P-10). **En HTTP:** `PATCH /proyectos/{id}/backlog/{elementoID}`, con el campo nuevo `horas_estimadas`. Como `story_points`, ausente no cambia nada y `null` quita la estimación:

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `horasEstimadas` | `*float64` | No | Si viene: mayor a 0, múltiplo de 0,5 y como máximo 1000, validado en ese orden (RN-19.2). |

```json
{ "horas_estimadas": 8 }
```

## 3. Salidas esperadas

**Resultado de la comparación** (`domain.ComparacionEsfuerzo`):

| Campo | Tipo (Go) | Significado |
|---|---|---|
| `HorasEstimadas` | `*float64` | Las del elemento; `nil` si está sin estimar. |
| `HorasReales` | `float64` | Suma de los registros de esfuerzo del elemento; 0 si no tiene ninguno. Nunca es `nil`. |
| `DesviacionHoras` | `*float64` | Real − estimado; `nil` si está sin estimar. |
| `DesviacionPorcentaje` | `*float64` | Con un decimal; `nil` si está sin estimar o si el estimado es 0. |

El resultado tiene además el método `SinEstimar() bool`, para que la pantalla no repita la regla.

| Caso | Estimadas | Reales | Desviación (h) | Desviación (%) |
|---|---|---|---|---|
| Se trabajó más de lo estimado | 8 | 10 | 2 | 25,0 |
| Se trabajó menos de lo estimado | 8 | 6 | −2 | −25,0 |
| Coincide | 8 | 8 | 0 | 0,0 |
| Sin esfuerzo registrado | 8 | 0 | −8 | −100,0 |
| Sin estimar, con esfuerzo | — | 10 | — | — |
| Sin estimar y sin esfuerzo | — | 0 | — | — |
| Estimado 0 (solo por la función del dominio, RN-19.8) | 0 | 10 | 10 | — |

- La consulta no modifica nada.
- **En HTTP:** `200` con la comparación en JSON (nulos cuando no hay valor); `404` si no existe el proyecto o el elemento:

```json
{
  "proyecto_id": 1,
  "elemento_id": 7,
  "horas_estimadas": 8,
  "horas_reales": 10,
  "desviacion_horas": 2,
  "desviacion_porcentaje": 25
}
```

- **Formato de los números:** en JSON son números, sin formato: 25,0 % sale como `25` y 33,3 % como `33.3`. El formato con coma y un decimal (`25,0`), como el de la tabla de arriba, es solo de la pantalla.
- **Horas estimadas en el elemento:** las respuestas JSON del elemento (crear, editar y listar el backlog) suman `horas_estimadas`, `null` si está sin estimar. No existe un GET JSON de un elemento: esa consulta es de las pantallas.
- **En la pantalla:** la pantalla del elemento de HU-06 (`/web/proyectos/{id}/backlog/{elementoID}`) suma un campo "Horas estimadas" al formulario y un bloque "Esfuerzo estimado y real" con los cuatro valores, después de "Registrar esfuerzo" y antes de "Eliminar". Las horas y los porcentajes se muestran con coma decimal; las desviaciones positivas llevan signo (`+2 h`, `+25,0 %`) y las negativas, el guion común (`-2 h`, `-25,0 %`). Sin estimar, se muestra "Sin estimar" en lugar de la estimación y de las dos desviaciones, incluido el porcentaje, junto con las horas reales. Con estimado 0, que solo es posible en la función del dominio (RN-19.8), el porcentaje muestra "No aplica". Sin esfuerzo, las horas reales se muestran como `0 h`.

## 4. Reglas de negocio
- **RN-19.1:** Las horas estimadas son un dato propio del elemento y son opcionales: un elemento sin ese dato está "sin estimar" en horas, aunque tenga story points. No se derivan de los story points (P-10).
- **RN-19.2:** Si se informan, las horas estimadas son mayores a 0, múltiplos de 0,5 y como máximo 1000 (P-12). El 0 no es una estimación válida: para no estimar, se deja el campo vacío.
- **RN-19.3:** Las horas estimadas se cargan, cambian o quitan editando el elemento (HU-06), en cualquier momento, aunque ya tenga esfuerzo registrado. Un elemento nuevo (HU-05) nace sin estimar. No se guarda historial: la comparación usa siempre el valor actual.
- **RN-19.4:** Las horas reales de un elemento son la suma de todos sus registros de esfuerzo (HU-18), de cualquier integrante y fecha, incluidos los de integrantes que después se quitaron del proyecto (RN-18.3). La suma no tiene tope: el de 24 h es por registro (RN-18.8). Sin registros, las horas reales son 0, y no es un error.
- **RN-19.5:** La desviación absoluta es `reales − estimadas`, en horas. Positiva: se trabajó más de lo estimado. Negativa: se trabajó menos.
- **RN-19.6:** La desviación porcentual es `(reales − estimadas) × 100 / estimadas`, redondeada a un decimal, con los medios alejándose de cero (1,25 → 1,3; −1,25 → −1,3). Lleva el mismo signo que la desviación absoluta. Si el redondeo da cero, el resultado es 0, nunca −0.
- **RN-19.7:** Un elemento sin estimar no es un error: se informan sus horas reales y ninguna desviación, ni absoluta ni porcentual. Tomar la estimación faltante como 0 mostraría todo lo trabajado como desviación.
- **RN-19.8:** Con un estimado igual a 0, la desviación absoluta se calcula y la porcentual no existe (no se divide por cero). Un elemento no puede tener estimado 0 (RN-19.2), pero la función del dominio resuelve el caso como defensa. Desde HU-23 no se alcanza: sin estimaciones le pasa `nil` (RN-23.6).
- **RN-19.9:** Un elemento con estimación y sin esfuerzo se compara igual: 0 horas reales y desviación de −100 %.
- **RN-19.10:** El proyecto tiene que existir y el elemento tiene que ser de ese proyecto; un elemento de otro proyecto se trata como inexistente (igual que RN-18.1).
- **RN-19.11:** La comparación se puede consultar con el elemento en cualquier estado, asignado o no a un sprint, y con el proyecto en cualquier estado. No se limita a los elementos `Hecho`: en uno sin terminar, la desviación muestra cómo va hasta ahora.
- **RN-19.12:** En el formulario de la pantalla del elemento, un campo `horas_estimadas` ausente del pedido (no la clave vacía) significa "no se envió" y no cambia las horas, aunque falte `original_horas_estimadas`. Presente y vacío, las quita si difiere del original. Presente con valor, las carga si difiere del original. Es coherente con el PATCH JSON, donde un campo ausente no cambia nada. `story_points` conserva el comportamiento de HU-06: un POST sin campos ocultos ni `story_points` lo quita. Esa asimetría es conocida y queda para el equipo.

## 5. Restricciones
- **Cálculo en el dominio y reutilizable (HU-23):** `CompararEsfuerzo(estimadas *float64, reales float64) ComparacionEsfuerzo` vive en `internal/domain/comparacion_esfuerzo.go`, sin dependencias externas, sin leer nada de la base y sin saber de elementos ni de sprints: recibe dos números, así HU-23 le puede pasar totales. Implementa RN-19.5 a RN-19.9. Recibe valores ya validados, así que no devuelve error.
- **Precisión:** en RN-19.6 se multiplica por 100 antes de dividir, para no acumular error de representación del `float64` antes de redondear.
- **Validación de las horas estimadas:** vive en el dominio, junto con la edición del elemento (`internal/domain/backlog_editar.go`). No se puede llamar a `validarHoras` de HU-18 tal como está, porque aplica el tope de 24 h; hay que separar las dos primeras comprobaciones y sumarle el tope propio.
- **Cambios en HU-05 y HU-06:** `ElementoBacklog` suma `HorasEstimadas *float64` y `CambiosElemento` suma `CambiarHorasEstimadas bool` y `HorasEstimadas *float64`, como los story points. La firma de `NuevoElementoBacklog` no cambia. El PATCH lee el campo como `json.RawMessage` para distinguir ausente de `null`. El formulario de la pantalla suma el campo oculto `original_horas_estimadas`, porque solo se guarda lo que cambió. Los errores nuevos se suman a `erroresDeValidacionEdicion`; si no, responderían `500`.
- **Persistencia:** `backlog_elementos` suma la columna `medias_horas_estimadas` (entero, `NULL` = sin estimar), en medias horas como `registros_esfuerzo`. Una base creada antes de HU-19 la recibe con `ALTER TABLE` al arrancar, igual que `sprint_id` en HU-09; sus elementos quedan sin estimar.
- **Consulta de horas reales, con puerto angosto:** el caso de uso define el puerto `LectorHorasReales`, con un solo método:

  ```go
  // HorasRealesDeElemento devuelve la suma de las horas registradas; 0 si no hay registros.
  HorasRealesDeElemento(ctx context.Context, elementoID int64) (float64, error)
  ```

  Devuelve solo el total: el detalle por integrante o por fecha queda fuera. Lo implementa `persistencia.RepositorioEsfuerzo` con `SELECT COALESCE(SUM(medias_horas), 0)` y `desdeMediasHoras`: se suman enteros, sin error de redondeo. No se ensancha `RepositorioEsfuerzo` de HU-18, porque lo implementan otros fakes de tests (segregación de interfaces, como `RepositorioCierre` en HU-11). Un elemento inexistente daría 0, pero no llega a esta consulta: el caso de uso ya lo verificó.
- **Caso de uso:** `CompararEsfuerzo` (`internal/app/comparar_esfuerzo.go`) reutiliza `LectorProyectos` (HU-04) y `LectorElementos` (HU-18). Busca primero el proyecto, porque `ConsultarElementoBacklog` no lo verifica y no alcanzaría para distinguir `ErrProyectoNoEncontrado`.
- **Rutas y pantallas:** `net/http` (ADR-0001), SQLite (ADR-0002) y `html/template` (ADR-0003). La ruta JSON se registra con `RegistrarCompararEsfuerzo` y se suma a `NuevoRouter`, a `CasosDeUso` y a `router_test.go`. La pantalla llama al mismo caso de uso y no hace ningún cálculo; el bloque también se muestra cuando el formulario vuelve con un error, con los valores guardados y no con los que se acaban de tipear.
- **Coordinación:** el campo nuevo toca los mismos archivos que HU-13 (validación de story points de HU-05 y HU-06). Conviene mergear seguido.
- **Fuera de alcance:**
  - **Para HU-23:** totales de horas estimadas, reales y desviación por sprint o por proyecto, y si las horas de una historia que volvió al backlog cuentan para el sprint en que se trabajaron (pregunta del planning). El porcentaje de un total se calcula con `CompararEsfuerzo` sobre las sumas, no promediando porcentajes.
  - **Para otras historias:** gráficos (HU-25) y reportes (HU-26).
  - **Sin historia:** historial de estimaciones, horas estimadas por integrante, listar, editar o eliminar registros de esfuerzo, cargar horas estimadas al crear el elemento y controlar quién puede consultar (el sistema no tiene roles todavía).

## 6. Casos límite
- Estimadas 8, reales 8 → desviación 0 h y 0,0 %.
- Estimadas 3, reales 4 → 1 h y 33,3 %. Estimadas 3, reales 2 → −1 h y −33,3 %.
- Estimadas 40, reales 40,5 → 0,5 h y 1,3 %. Estimadas 40, reales 39,5 → −0,5 h y −1,3 % (RN-19.6).
- Estimadas 0,5, reales 24 → 23,5 h y 4700,0 %: el porcentaje no tiene tope.
- Elemento con esfuerzo que supera las 24 h en total (24 + 24 + 8) → reales 56 (RN-19.4).
- Elemento con story points y sin horas estimadas → sin estimar; los story points no se usan (RN-19.1).
- Se cambian los story points de un elemento → sus horas estimadas no cambian.
- Se quita la estimación de un elemento con esfuerzo → pasa a sin estimar y conserva las horas reales.
- Se cambia la estimación después de registrar esfuerzo → la desviación se calcula con el valor nuevo (RN-19.3).
- Registros de un integrante que ya no está en el proyecto → siguen sumando (RN-19.4).
- Los registros de otro elemento no suman, aunque sean del mismo proyecto y del mismo integrante.
- Elemento `Pendiente` o que volvió al backlog al cerrar un sprint (RN-11.3), con esfuerzo → se compara igual (RN-19.11).
- Horas estimadas 0,5 → válido; 0 o negativas → error; 1,25 → error de múltiplo; 1000 → válido; 1000,5 → error de tope.
- Horas estimadas `"NaN"` o `"Inf"` en el formulario → las rechaza el dominio: `NaN` no es mayor a 0 e `Inf` supera el tope.
- Base creada antes de HU-19 → sus elementos se leen sin estimar y se pueden editar.
- Elemento que existe pero en otro proyecto → elemento no encontrado.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05, se reutiliza) | "elemento de backlog no encontrado" (HTTP 404) |
| Horas estimadas 0 o negativas | `domain.ErrHorasInvalidas` (HU-18, se reutiliza) | "las horas deben ser mayores a 0" (HTTP 400) |
| Horas estimadas que no son múltiplo de 0,5 | `domain.ErrHorasNoMultiploMediaHora` (HU-18, se reutiliza) | "las horas se registran en múltiplos de media hora" (HTTP 400) |
| Horas estimadas mayores a 1000 | `domain.ErrHorasEstimadasExcedenMaximo` (nuevo) | "las horas estimadas no pueden superar las 1000 horas" (HTTP 400) |
| IDs de la ruta que no son números, o `horas_estimadas` que no es un número | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso. |
| Horas estimadas que no son un número en el formulario de la pantalla | — (lectura del formulario; lo define el adaptador) | "las horas estimadas deben ser un número" (HTTP 400, en el formulario); no se llama al caso de uso. |
| Falla interna al consultar | — (solo se envuelve con `%w` la falla al leer las horas reales; las de las búsquedas salen tal cual) | HTTP 500 con un mensaje genérico, sin el detalle. |

Sin estimar, sin esfuerzo y estimado 0 **no son errores** (RN-19.7 a RN-19.9).

Si hay varios errores, se informa el primero en el orden de la tabla: primero se busca el proyecto y después el elemento. Las respuestas muestran el mensaje del error del dominio, nunca el texto envuelto por el caso de uso. En la pantalla, los `404` muestran la página de error y los `400` vuelven al formulario con el mensaje.

## 8. Criterios de aceptación
- **CA-19.1 (desviación positiva):** Dado un elemento con 8 h estimadas y registros de 6 h y 4 h, cuando comparo, entonces obtengo 8 estimadas, 10 reales, desviación de 2 h y de 25,0 %. (RN-19.4, RN-19.5, RN-19.6)
- **CA-19.2 (desviación negativa):** Dado un elemento con 8 h estimadas y un registro de 6 h, cuando comparo, entonces obtengo desviación de −2 h y de −25,0 %. (RN-19.5, RN-19.6)
- **CA-19.3 (sin desviación):** Dado un elemento con 8 h estimadas y registros que suman 8 h, cuando comparo, entonces obtengo desviación de 0 h y de 0,0 %. (RN-19.5, RN-19.6)
- **CA-19.4 (sin estimar):** Dado un elemento sin horas estimadas y con 10 h registradas, cuando comparo, entonces obtengo 10 reales, sin estimación ni desviaciones, y `SinEstimar()` es verdadero. No es un error. (RN-19.1, RN-19.7)
- **CA-19.5 (sin esfuerzo):** Dado un elemento con 8 h estimadas y sin registros, cuando comparo, entonces obtengo 0 reales, desviación de −8 h y de −100,0 %. (RN-19.4, RN-19.9)
- **CA-19.6 (sin estimar y sin esfuerzo):** Dado un elemento sin horas estimadas y sin registros, cuando comparo, entonces obtengo 0 reales, sin estimación ni desviaciones. (RN-19.4, RN-19.7)
- **CA-19.7 (estimado 0):** Cuando llamo a la función del dominio con estimado 0 y 10 h reales, entonces obtengo desviación de 10 h y sin porcentaje. (RN-19.8)
- **CA-19.8 (redondeo):** Cuando comparo 3 h estimadas con 4 h reales, entonces el porcentaje es 33,3; con 40 y 40,5, es 1,3; con 40 y 39,5, es −1,3. (RN-19.6)
- **CA-19.9 (inexistente):** Dado un ID de proyecto o de elemento que no existe, cuando comparo, entonces obtengo `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado`. (RN-19.10)
- **CA-19.10 (otro proyecto):** Dado un elemento de otro proyecto, cuando comparo, entonces obtengo `ErrElementoNoEncontrado`. (RN-19.10)
- **CA-19.11 (qué suma):** Dado un elemento con registros de un integrante que después se quitó del proyecto, y otro elemento del mismo proyecto con sus propios registros, cuando comparo el primero, entonces suman los registros del integrante quitado y no los del otro elemento. (RN-19.4)
- **CA-19.12 (cualquier estado):** Dado un elemento `Pendiente` sin sprint y con esfuerzo, en un proyecto que no está Activo, cuando comparo, entonces obtengo la comparación. (RN-19.11)
- **CA-19.13 (cargar la estimación):** Dado un elemento sin estimar, cuando edito sus horas estimadas a 8, entonces quedan guardadas; si después edito otro campo, siguen en 8; si las envío en `null`, queda sin estimar. Se puede hacer aunque el elemento tenga esfuerzo. (RN-19.3)
- **CA-19.14 (validar la estimación):** Cuando edito las horas estimadas con 0 o negativas, obtengo `ErrHorasInvalidas`; con 1,25, `ErrHorasNoMultiploMediaHora`; con 1000,5, `ErrHorasEstimadasExcedenMaximo`; con 0,5 o 1000, se guardan. Si hay error, el elemento no cambia. (RN-19.2)
- **CA-19.15 (independiente de los story points):** Dado un elemento con 5 story points y sin horas estimadas, cuando comparo, entonces está sin estimar; cuando cambio sus story points, sus horas estimadas no cambian. Un elemento recién creado nace sin estimar. (RN-19.1, RN-19.3)
- **CA-19.16 (base existente):** Dada una base creada antes de HU-19, cuando arranca el sistema, entonces sus elementos se leen sin estimar y se les pueden cargar horas estimadas. (Sección 5)
- **CA-19.17 (HTTP):** Dado un pedido a `GET /proyectos/{id}/backlog/{elementoID}/esfuerzo/comparacion`, entonces la respuesta es `200` con la comparación en JSON, con `null` en los valores que no existen; con proyecto o elemento inexistente, `404`; con un ID que no es un número, `400`. (Sección 3)
- **CA-19.18 (pantalla):** Dada la pantalla de un elemento con estimación y esfuerzo, entonces muestra, después de "Registrar esfuerzo" y antes de "Eliminar", las horas estimadas, las reales y las dos desviaciones, las negativas con el guion común; dado uno sin estimar, muestra "Sin estimar" en lugar de la estimación y de las dos desviaciones, y sus horas reales; dado uno sin esfuerzo, muestra `0 h`; con estimado 0, que solo es posible en la función del dominio, el porcentaje muestra "No aplica". Desde esa pantalla se pueden cargar y quitar las horas estimadas, y los errores se muestran con el mensaje del dominio. (Sección 3)
- **CA-19.19 (campo ausente en el formulario):** Dado un elemento con 8 h estimadas, cuando envío el formulario de la pantalla sin el campo `horas_estimadas`, aunque también falten los campos ocultos, entonces conserva las 8 h; cuando lo envío con el campo vacío, entonces queda sin estimar; cuando lo envío con 10, entonces quedan 10 h. Un POST sin campos ocultos ni `story_points` sigue quitando los story points, como en HU-06. (RN-19.12)

Los escenarios mínimos corresponden a: feliz y desviación positiva, CA-19.1; desviación negativa, CA-19.2; sin estimar, CA-19.4; sin esfuerzo, CA-19.5; elemento inexistente, CA-19.9; otro proyecto, CA-19.10.

## Decisiones tomadas
Tomadas por el autor (Lucas Orellana) el 2026-11-01, abiertas a objeción en el PR:
- **P-1:** el resultado tiene cuatro valores: estimadas, reales, desviación en horas (real − estimado) y en porcentaje (sección 3, RN-19.5, RN-19.6). Positivo significa que se trabajó más de lo estimado.
- **P-2:** sin estimar no es un error y no hay ninguna desviación, ni siquiera la absoluta (RN-19.7). En JSON se informa con `null`, como `story_points`, y en pantalla con el texto "Sin estimar". No se suma un campo de texto al JSON: el `null` ya lo dice.
- **P-3:** sin esfuerzo son 0 horas reales y la comparación se hace igual, con −100 % si hay estimación (RN-19.9). Se informa así porque es lo que dicen los datos; sin esfuerzo no quiere decir que la historia haya tomado 0 horas, y por eso RN-19.11 aclara que en un elemento sin terminar la desviación es parcial.
- **P-4:** el estimado 0 no se puede cargar en un elemento (RN-19.2), igual que los story points de HU-05, pero la función del dominio lo resuelve sin dividir por cero, como defensa (RN-19.8).
- **P-5:** el porcentaje lleva un decimal, redondeado con los medios alejándose de cero (RN-19.6). Un decimal alcanza para comparar sin aparentar más precisión que la de los datos, que vienen en medias horas. En JSON es un número; la coma decimal y el signo son solo de la pantalla.
- **P-6:** las horas reales salen de un puerto nuevo de un solo método, `LectorHorasReales`, que devuelve solo el total; no se ensancha ninguna interfaz existente (sección 5).
- **P-7:** el cálculo vive en una función del dominio que recibe números y no entidades, para que HU-23 la reutilice con totales (sección 5).
- **P-8:** el alcance es un elemento. Los totales por sprint o por proyecto y qué hacer con las historias que volvieron al backlog son de HU-23 (sección 5).
- **P-9:** se compara en cualquier estado del elemento, del sprint y del proyecto (RN-19.11), y usa la estimación actual, sin historial (RN-19.3). Es una consulta de solo lectura, y limitarla a `Hecho` no dejaría ver cómo viene una historia en curso.

## Decisiones confirmadas por el equipo
Confirmadas por Matías el 2026-11-03, en la revisión del PR #297. No están abiertas.

- **P-10 — De dónde salen las horas estimadas: campo propio `HorasEstimadas` (opción A).** Antes de HU-19, `ElementoBacklog` solo tenía `StoryPoints`.

  | Opción | Ventajas | Desventajas |
  |---|---|---|
  | **A. Campo nuevo `HorasEstimadas` en el elemento, editable con HU-06** | Es una estimación en horas de verdad, así que la desviación mide la precisión de esa estimación. No inventa ningún factor. Los story points siguen siendo tamaño relativo. Es lo que sugiere el planning. | Toca código de HU-05 y HU-06 (dominio, SQLite, PATCH y pantalla), justo en una historia de 3 puntos. Alguien tiene que cargar un dato más por elemento. |
  | **B. Conversión fija desde story points** (por ejemplo, 1 SP = 4 h) | No hay campo ni cambio de pantalla. Lo estimado en Planning Poker (HU-17) se usa directamente. | Inventa un factor sin datos que lo respalden. La desviación mediría sobre todo si el factor está bien, no la estimación. Mezcla tamaño relativo con duración. Cambiar el factor reescribe todas las desviaciones. |
  | **C. Factor por proyecto, configurable** | Cada equipo lo puede calibrar con su velocidad. No hay que cargar horas por elemento. | Mismo problema de fondo que B. Además suma un dato de proyecto con su edición y su pantalla: es más trabajo que A. |

  **Decisión: opción A**, por la que está escrita esta spec (RN-19.1 a RN-19.3, CA-19.13 a CA-19.16). Con B o C se quitarían esas reglas y criterios, y el dominio no cambiaría: `CompararEsfuerzo` sigue recibiendo un número. Afecta también a HU-23.

- **P-12 — Tope de las horas estimadas: 1000 h por elemento.** RN-18.8 (24 h) es por registro y no sirve para una estimación. El tope frena errores de tipeo y valores que desbordarían el entero de medias horas de la base. El número es arbitrario; la alternativa era otro valor, pero sin ningún tope `1e300` llegaría a la conversión a entero.

## Preguntas abiertas
P-11 y P-13 quedaron resueltas. P-10 y P-12 pasaron a "Decisiones confirmadas por el equipo".

- **P-11 — Dónde se muestra la comparación** (la DoD exige pantalla, ADR-0003). **Resuelta con lo implementado: decisión tomada por el autor, abierta a objeción en la PR.** Se muestra en la pantalla del elemento de HU-06, porque la comparación es por elemento y ahí mismo se cargan las horas estimadas, después de "Registrar esfuerzo" y antes de "Eliminar" (sección 3). La alternativa era una columna en el detalle del sprint (HU-09), que permite ver varias historias juntas pero deja afuera a los elementos sin sprint y se superpone con lo que va a mostrar HU-23.

- **P-13 — No había pantalla para cargar esfuerzo. Resuelta por la pantalla de HU-18** (issue #289, rama `hu-18-pantalla-esfuerzo`): el formulario "Registrar esfuerzo" está en la pantalla del elemento. HU-18 decía que su pantalla "se suma después con HU-30", pero HU-30 no la incluyó, y sin ella quien usara la interfaz vería siempre 0 horas reales en esta comparación. Se resolvió aparte de HU-19, como se proponía acá.

## Historial
- 2026-11-01: creada (borrador). Decisiones P-1 a P-9 tomadas por el autor, abiertas a objeción en el PR; P-10 a P-13, pendientes del equipo.
- 2026-11-01: ajustada al contrato JSON real, sin cambiar reglas ni criterios: en el ejemplo de la sección 3, `desviacion_porcentaje` es `25` (el formato con coma y un decimal es solo de la pantalla), y se aclara que no existe un GET JSON de un elemento.
- 2026-11-01: se corrige la fila de falla interna de la sección 7 (solo se envuelve la falla al leer las horas reales).
- 2026-11-01: ajustada a la pantalla implementada. Secciones 3 y CA-19.18: el bloque va después de "Registrar esfuerzo" y antes de "Eliminar", sin estimar el porcentaje muestra "Sin estimar", con estimado 0 muestra "No aplica" y las desviaciones negativas llevan el guion común. Sección 7: se quita lo pendiente de `errorDeElemento`, ya hecho por HU-18, y se suma el mensaje "las horas estimadas deben ser un número". P-11 queda resuelta por el autor, abierta a objeción en la PR; P-13, resuelta por la pantalla de HU-18 (issue #289); P-10 y P-12, implementadas según la propuesta, abiertas a objeción del equipo en la PR. La spec sigue en borrador.
- 2026-11-02: se agregan RN-19.12 y CA-19.19: un campo `horas_estimadas` ausente del formulario de la pantalla no cambia las horas. `story_points` conserva el comportamiento de HU-06; la asimetría queda para el equipo. La spec sigue en borrador.
- 2026-11-03: Matías confirma P-10 (campo propio `HorasEstimadas`) y P-12 (tope de 1000 h) en la revisión del PR #297; pasan a "Decisiones confirmadas por el equipo". Las reglas y los criterios no cambian.
- 2026-11-04: se corrige la frase de RN-19.8 y P-4 sobre las sumas de HU-23 (P-11 de HU-23). Las reglas y los criterios no cambian.
