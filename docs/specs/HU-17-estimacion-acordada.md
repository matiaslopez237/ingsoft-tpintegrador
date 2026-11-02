# HU-17: Registrar la estimación acordada

**Historia:** Como Product Owner, quiero registrar la estimación final acordada por el equipo, para dejarla asociada a la historia.
**Issue:** #28  ·  **Sprint:** 3  ·  **Story points:** 2
**Estado de la spec:** Aprobada

## 1. Objetivo
Cerrar el ciclo del Planning Poker (HU-14 a HU-16): cuando el equipo termina de votar y discutir, registrar el valor acordado. Ese valor pasa a ser los story points del elemento del backlog, y la ronda conserva lo que se acordó en ella, así queda trazable de qué ronda salió la estimación.

## 2. Entradas

**Registrar la estimación** (`POST /proyectos/{id}/rondas/{rondaID}/estimacion`):

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID`, `rondaID` | `int64` | Sí | La ronda existe y es de ese proyecto (RN-14.7). |
| `valor` | `int` | Sí | Uno de 1, 2, 3, 5, 8, 13 o 21 (RN-17.2). |

```json
{ "valor": 8 }
```

## 3. Salidas esperadas
- La ronda guarda la estimación acordada y los story points del elemento pasan a ser ese valor. Los demás datos del elemento no cambian. Si hay un error, no se guarda nada.
- Se devuelve el resumen de la ronda de HU-14 y HU-15, que suma la estimación acordada.
- **En HTTP:** `200` con el JSON del resumen, que suma `estimacion_acordada` cuando la ronda tiene una (si no, el campo no aparece); `400` si el dato no es válido o no se puede registrar; `404` si no existe el proyecto o la ronda, o si la ronda es de otro proyecto.
- **Pantallas** (ADR-0003):
  - **Pantalla de la ronda:** si se puede registrar (RN-17.3), un formulario con la escala para elegir el valor (campo `estimacion`) y el botón "Registrar estimación". Al guardar se vuelve a la misma pantalla (POST-redirect-GET, `303`); si hay un error, se muestra en ella. Si la ronda ya tiene una estimación, se muestra "Estimación acordada: N".
  - **Pantalla del elemento**, "Historial de rondas" (HU-16): las rondas con una estimación muestran "acordada: N".

## 4. Reglas de negocio
- **RN-17.1:** Registrar la estimación guarda el valor en la ronda y lo copia a los story points del elemento, en una sola operación. Si el elemento ya tenía story points, se reemplazan.
- **RN-17.2:** El valor es uno de la escala 1, 2, 3, 5, 8, 13 o 21 (RN-13.1, RN-14.4). No hace falta que alguien lo haya votado: después de discutir, el equipo puede acordar otro valor.
- **RN-17.3:** Solo se registra sobre la **última ronda** del elemento, **cerrada** y **sin estimación**, y si el elemento **no está `Hecho`** (RN-14.1).
- **RN-17.4:** Cada ronda tiene a lo sumo una estimación, y no se corrige. Para cambiarla, se vota otra ronda (HU-16) y se registra ahí. La anterior queda en el historial.
- **RN-17.5:** La estimación de la ronda es un registro de lo que se acordó. Si después se editan los story points del elemento (HU-06), el elemento muestra el valor nuevo y la ronda conserva el acordado. HU-06 no se bloquea.
- **RN-17.6:** Registrar una estimación no impide abrir otra ronda (RN-16.4).

## 5. Restricciones
- RN-17.2 a RN-17.4 viven en el dominio (`internal/domain/planning_poker_estimacion.go`): `Ronda.RegistrarEstimacion(valor, ultima, elemento)` aplica las reglas, y `PuedeRegistrarEstimacion` le dice a la pantalla si ofrece el formulario, sin repetirlas.
- La escala es la única del dominio (RN-13.4). Por eso esta historia se mergea **después de HU-13 (#288)**, que unifica la de los votos y la de los story points.
- El caso de uso `RegistrarEstimacion` (`internal/app`) busca el proyecto, la ronda, la última ronda del elemento y el elemento, aplica el dominio y guarda con un puerto angosto propio, `RepositorioEstimacion`.
- **Persistencia** (SQLite, ADR-0002): columna `estimacion_acordada` en `rondas` (`NULL` si no tiene), agregada también a las bases existentes con `agregarColumnas`.
- **Pedidos simultáneos:** guardar es una transacción con las condiciones en los mismos `UPDATE`: la ronda se actualiza si sigue siendo la última, cerrada y sin estimación, y el elemento si no está `Hecho`. Del elemento solo se actualiza `story_points`, para no pisar otros campos (lección de DEF-01). Si una condición falla, no se guarda nada y se informa el error de esa regla.
- `net/http` (ADR-0001): la ruta JSON se registra con `RegistrarEstimacionAcordada`, que se suma a `NuevoRouter`, a `CasosDeUso` y a `router_test.go`. La acción de pantalla pasa por el chequeo de mismo origen de los formularios.
- Fuera de alcance: controlar quién registra (el sistema no tiene roles), corregir una estimación ya registrada y mostrar cuántas veces se reestimó.

## 6. Casos límite
- Valor que nadie votó, pero de la escala → se registra (RN-17.2).
- Valor 4, 0 o 34 → error.
- Elemento sin estimar → queda con el valor; elemento con 5 puntos → pasa al valor nuevo.
- Elemento con un valor anterior a la escala (RN-13.3) → se reemplaza por el acordado.
- Ronda abierta → error.
- Ronda 1 cerrada con la ronda 2 ya creada, abierta o cerrada → error en la ronda 1, porque no es la última.
- Registrar dos veces en la misma ronda → la segunda da error y conserva la primera.
- Elemento que pasó a `Hecho` después de cerrar la ronda → error.
- Ronda 1 con estimación 5 y ronda 2 con estimación 8 → el elemento queda con 8 y el historial muestra los dos acuerdos.
- Editar los story points a mano después (HU-06) → el elemento cambia y la ronda sigue mostrando el acordado.
- Dos pedidos simultáneos sobre la misma ronda → se guarda uno y el otro recibe `ErrEstimacionYaRegistrada`.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01) | "proyecto no encontrado" (HTTP 404) |
| No existe la ronda, o es de otro proyecto | `domain.ErrRondaNoEncontrada` (HU-14) | "ronda de estimación no encontrada" (HTTP 404) |
| La ronda sigue abierta | `domain.ErrEstimacionRondaAbierta` | "la estimación se registra cuando votaron todos" (HTTP 400) |
| La ronda no es la última del elemento | `domain.ErrRondaNoEsLaUltima` | "solo se puede registrar la estimación en la última ronda" (HTTP 400) |
| El elemento está `Hecho` | `domain.ErrElementoTerminado` (HU-14) | "no se puede estimar un elemento terminado" (HTTP 400) |
| La ronda ya tiene estimación | `domain.ErrEstimacionYaRegistrada` | "la ronda ya tiene una estimación acordada" (HTTP 400) |
| El valor no es de la escala | `domain.ErrEstimacionInvalida` | "la estimación debe ser 1, 2, 3, 5, 8, 13 o 21" (HTTP 400) |
| Cuerpo JSON mal formado o con campos desconocidos | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso |

Si hay varios errores, se informa el primero en el orden de la tabla. Cualquier otro error es una falla interna: `500` con un mensaje genérico, sin el detalle.

## 8. Criterios de aceptación
- **CA-17.1:** Dada la última ronda de "Login", cerrada, cuando registro 8, entonces la ronda queda con la estimación 8 y "Login" con 8 story points. (RN-17.1)
- **CA-17.2:** Dado un elemento con 5 story points y otros datos cargados, cuando registro 8, entonces pasa a 8 y sus demás datos no cambian. (RN-17.1)
- **CA-17.3:** Dada una ronda con votos 3 y 13, cuando registro 8, entonces se acepta aunque nadie lo haya votado; cuando registro 4, obtengo `ErrEstimacionInvalida` y no se guarda nada. (RN-17.2)
- **CA-17.4:** No se puede registrar en una ronda abierta (`ErrEstimacionRondaAbierta`), en una que no es la última (`ErrRondaNoEsLaUltima`) ni con el elemento `Hecho` (`ErrElementoTerminado`). (RN-17.3)
- **CA-17.5:** Dada una ronda con estimación, cuando registro otra, entonces obtengo `ErrEstimacionYaRegistrada` y se conserva la primera. (RN-17.4)
- **CA-17.6:** Dada la ronda 1 con estimación 5, cuando se vota la ronda 2 y registro 8, entonces el elemento queda con 8 y el historial muestra 5 en la ronda 1 y 8 en la ronda 2. (RN-17.4, RN-17.6)
- **CA-17.7:** Dada una ronda con estimación 8, cuando edito a mano los story points del elemento a 13 (HU-06), entonces el elemento queda con 13 y la ronda sigue con 8. (RN-17.5)
- **CA-17.8:** Dado un proyecto o una ronda inexistentes, o una ronda de otro proyecto, entonces obtengo `ErrProyectoNoEncontrado` o `ErrRondaNoEncontrada`. (Sección 7)
- **CA-17.9:** La API responde `200` con `estimacion_acordada`; los errores de validación, `400`; lo inexistente, `404`. Una ronda sin estimación no trae el campo. (Sección 3)
- **CA-17.10:** La pantalla de la última ronda, cerrada y sin estimación, ofrece el formulario, y al guardar muestra "Estimación acordada: N" y ya no lo ofrece. Una ronda abierta, una que no es la última o una de un elemento `Hecho` no lo ofrece. Los errores se muestran en la misma pantalla. El historial del elemento muestra "acordada: N". (Sección 3)

## Decisiones tomadas
Confirmadas por Matías López (Agile Enabler) el 2026-11-02:
- **P-1:** cualquier valor de la escala, aunque nadie lo haya votado (RN-17.2).
- **P-2:** solo sobre la última ronda, cerrada (RN-17.3).
- **P-3:** una estimación por ronda; para cambiarla se vota otra ronda (RN-17.4).
- **P-4:** la edición manual (HU-06) no se bloquea, y la ronda conserva lo acordado (RN-17.5).
- **P-5:** un elemento `Hecho` no recibe estimación, como RN-14.1, que no permite estimarlo.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-02: creada, con P-1 a P-5 confirmadas. Matías aprobó la spec.
