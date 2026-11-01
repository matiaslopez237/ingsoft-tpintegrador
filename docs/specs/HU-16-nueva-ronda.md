# HU-16: Repetir una ronda de votación

**Historia:** Como equipo, quiero iniciar una nueva ronda de votación sobre la misma historia, para converger en una estimación.
**Issue:** #27  ·  **Sprint:** 3  ·  **Story points:** 2
**Estado de la spec:** Borrador

## 1. Objetivo
Cuando una ronda de Planning Poker termina con estimaciones muy distintas (HU-15), el equipo discute y vuelve a votar. Esta historia permite abrir otra ronda sobre el mismo elemento desde la pantalla de la ronda cerrada y ver el historial de todas las rondas del elemento, con sus votos, para comparar cómo fueron convergiendo. Reutiliza las rondas y los votos de HU-14 y la revelación de HU-15.

## 2. Entradas

**Abrir una ronda nueva:** la misma operación de HU-14 (`POST /proyectos/{id}/backlog/{elementoID}/rondas`, sin cuerpo). HU-16 no agrega entradas: cuando la última ronda del elemento está cerrada, esa operación abre la siguiente.

**Consultar el historial** (`GET /proyectos/{id}/backlog/{elementoID}/rondas`):

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento del backlog de ese proyecto (RN-16.5). |

## 3. Salidas esperadas
- **Ronda nueva:** queda `Abierta`, con el número siguiente, los integrantes actuales del proyecto como votantes y ningún voto (RN-14.2). Las rondas anteriores no cambian (RN-16.2).
- **Historial:** todas las rondas del elemento, de la más reciente a la más antigua (RN-16.3). Cada una con el resumen de HU-14 y HU-15: número, estado, votantes y, si está cerrada, la revelación (el voto de cada integrante, el mínimo, el máximo y si hubo diferencia significativa). De una ronda abierta solo se ve quién votó, nunca los valores (RN-14.3). Un elemento sin rondas da una lista vacía, sin error.
- La consulta no modifica nada.
- **En HTTP:** `200` con un arreglo JSON (`[]` si no hay rondas), donde cada elemento tiene la forma del resumen de una ronda de HU-14 y HU-15 (`id`, `proyecto_id`, `elemento_id`, `numero`, `estado`, `votantes` y, si está cerrada, `revelacion`). `404` si no existe el proyecto o el elemento, o si el elemento es de otro proyecto; `400` si un ID de la ruta no es un número.
- **Pantallas** (ADR-0003):
  - **Pantalla de la ronda** (HU-14): si la ronda está cerrada, es la última del elemento y el elemento no está `Hecho`, el botón "Nueva ronda" (RN-16.4). Al usarlo se abre la ronda siguiente y se va a su pantalla (POST-redirect-GET, `303`). Si no se puede (por ejemplo, otro pedido ya la abrió), se vuelve a la pantalla de la ronda con el error.
  - **Pantalla del elemento** (HU-06), sección "Planning Poker": debajo de lo que ya muestra HU-14 (el enlace a la ronda abierta o el botón "Iniciar ronda"), el historial: una fila por ronda, de la más reciente a la más antigua, con el número, el estado, el mínimo y el máximo y si hubo diferencia (en las cerradas), y un enlace a la pantalla de esa ronda, que muestra el voto de cada integrante. Sin rondas, no se muestra la lista.

## 4. Reglas de negocio
- **RN-16.1:** Sobre un elemento se puede abrir una ronda nueva cuando su última ronda está `Cerrada` y el elemento no está `Hecho`. Toma el número siguiente y los integrantes actuales como votantes, sin votos. No se puede con una ronda todavía abierta: hay que esperar a que voten todos (RN-14.1, RN-14.6).
- **RN-16.2:** Abrir una ronda nueva no modifica las anteriores: sus votos quedan guardados y se siguen viendo en el historial.
- **RN-16.3:** El historial de un elemento son todas sus rondas, de la más reciente (número mayor) a la más antigua. De las cerradas se muestran los votos revelados, el mínimo, el máximo y si hubo diferencia significativa (HU-15); de la abierta, solo quién votó (RN-14.3).
- **RN-16.4:** La pantalla de una ronda ofrece "Nueva ronda" solo si la ronda está `Cerrada`, es la última del elemento y el elemento no está `Hecho`. Se ofrece siempre que se cumpla eso, haya habido o no diferencia significativa: repetir es decisión del equipo.
- **RN-16.5:** El elemento tiene que pertenecer al proyecto de la ruta; si es de otro proyecto, se trata como inexistente (como RN-14.7).

## 5. Restricciones
- RN-16.1 ya la implementa `domain.NuevaRonda` de HU-14 (recibe la última ronda del elemento): HU-16 no cambia esa regla, la cubre con escenarios propios.
- RN-16.3 vive en el dominio (`internal/domain/planning_poker.go`), en una función pura `HistorialRondas(rondas []Ronda) []ResumenRonda`, que ordena y arma el resumen de cada ronda con `Ronda.Resumen` de HU-14 (que ya oculta los votos de la ronda abierta y suma la revelación de HU-15). Así la regla de ocultar no se repite.
- RN-16.4 también vive en el dominio: `PuedeRepetirse(ronda, ultima, elemento) bool`. La pantalla solo la consulta.
- El caso de uso `ListarRondas` (`internal/app/planning_poker.go`) verifica el proyecto (`LectorProyectos`) y el elemento (`LectorElementos`, RN-16.5) y lee las rondas por un **puerto angosto propio**, `LectorHistorialRondas`, con `ListarRondas(ctx, elementoID)`. No se ensancha `RepositorioRondas` de HU-14, para no tocar sus fakes (como hizo HU-11 con `RepositorioCierre`). El repositorio SQLite de rondas lo implementa.
- Es una consulta y una pantalla: abrir la ronda reutiliza el caso de uso y la persistencia de HU-14, que ya resuelve dos pedidos simultáneos (índice único de rondas abiertas, RN-14.1).
- `net/http` (ADR-0001). La ruta JSON se registra con su función `RegistrarHistorialRondas`, que se suma a `NuevoRouter`, a `CasosDeUso` y a `router_test.go`. El botón "Nueva ronda" usa la acción de pantalla de HU-14 (`POST /web/proyectos/{id}/backlog/{elementoID}/rondas`).
- Fuera de alcance: cancelar una ronda abierta, guardar la estimación acordada (HU-17) y comparar rondas con gráficos.

## 6. Casos límite
- Elemento sin rondas → historial vacío.
- Una sola ronda, abierta → el historial la muestra sin valores; no hay "Nueva ronda".
- Ronda cerrada sin diferencia significativa (todos votaron 5) → igual se ofrece "Nueva ronda" (RN-16.4).
- Ronda 1 cerrada, ronda 2 abierta → la pantalla de la ronda 1 no ofrece "Nueva ronda" (no es la última); el historial muestra 2 (sin valores) y 1 (con valores).
- Elemento que pasó a `Hecho` después de cerrar la ronda → no se ofrece "Nueva ronda" y abrirla da `ErrElementoTerminado` (RN-14.1).
- Integrante agregado al proyecto después de la ronda 1 → es votante de la ronda 2, no de la 1 (RN-14.2).
- Dos pedidos simultáneos de "Nueva ronda" → se abre una sola; el otro recibe `ErrRondaYaAbierta` (HU-14).
- Elemento de otro proyecto → no encontrado.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05) | "elemento de backlog no encontrado" (HTTP 404) |
| Abrir una ronda con otra abierta | `domain.ErrRondaYaAbierta` (HU-14) | "el elemento ya tiene una ronda abierta" (HTTP 400; en la pantalla, el mensaje en la ronda) |
| Abrir una ronda sobre un elemento `Hecho` | `domain.ErrElementoTerminado` (HU-14) | "no se puede estimar un elemento terminado" (HTTP 400) |
| Un ID de la ruta no es un número | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso |

HU-16 no define errores nuevos. Cualquier otro error es una falla interna: `500` con un mensaje genérico, sin el detalle.

## 8. Criterios de aceptación
- **CA-16.1:** Dada la ronda 1 de "Login" cerrada con Ana 3 y Beto 13, cuando inicio una ronda nueva, entonces se abre la ronda 2 con Ana y Beto como votantes y ningún voto, y la ronda 1 sigue cerrada con los mismos votos. (RN-16.1, RN-16.2)
- **CA-16.2:** Dada una ronda cerrada en la que todos votaron 5 (sin diferencia significativa), cuando inicio una ronda nueva, entonces se abre igual. (RN-16.4)
- **CA-16.3:** Dada una ronda abierta, cuando intento iniciar otra, entonces obtengo `ErrRondaYaAbierta` y no se crea ninguna. (RN-16.1)
- **CA-16.4:** Dadas las rondas 1 y 2 cerradas y la 3 abierta de un elemento, cuando consulto el historial, entonces obtengo 3, 2 y 1, en ese orden; la 1 y la 2 con sus votos, mínimo, máximo y diferencia, y la 3 sin ningún valor. (RN-16.3)
- **CA-16.5:** Dado un elemento sin rondas, cuando consulto el historial, entonces obtengo una lista vacía, sin error. (RN-16.3)
- **CA-16.6:** Dado un proyecto o un elemento inexistentes, o un elemento de otro proyecto, cuando consulto el historial, entonces obtengo `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado`. (RN-16.5)
- **CA-16.7:** Dado un pedido a `GET /proyectos/{id}/backlog/{elementoID}/rondas`, entonces la respuesta es `200` con el arreglo JSON (`[]` si no hay rondas), sin valores de la ronda abierta; con un proyecto o elemento inexistente, `404`; con un ID que no es un número, `400`. (Sección 3)
- **CA-16.8:** La pantalla de la última ronda, cerrada, ofrece "Nueva ronda" y al usarla se va a la ronda siguiente; la de una ronda cerrada que no es la última, o la de una ronda abierta, no la ofrece. La pantalla del elemento muestra el historial de la más reciente a la más antigua, con un enlace a cada ronda. (RN-16.3, RN-16.4, sección 3)

## Decisiones tomadas
Confirmadas por Matías López (Agile Enabler) el 2026-11-01:
- **P-1:** no se puede abrir una ronda nueva con otra abierta: hay que esperar a que voten todos (RN-16.1, igual que RN-14.1). No hay estado "Cancelada".
- **P-2:** "Nueva ronda" se ofrece siempre que la última ronda esté cerrada, haya habido o no diferencia significativa (RN-16.4).
- **P-3:** el historial muestra de cada ronda cerrada los votos, el mínimo, el máximo y si hubo diferencia; de la abierta, solo quién votó (RN-16.3).
- **P-4:** el historial se ve en la sección "Planning Poker" de la pantalla del elemento, de la ronda más reciente a la más antigua.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-01: creada, con P-1 a P-4 confirmadas.
