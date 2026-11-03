# HU-14: Registrar estimaciones individuales ocultas

**Historia:** Como Product Builder, quiero registrar mi estimación de forma oculta hasta que todos voten, para evitar sesgos por ver el voto de otros.
**Issue:** #25  ·  **Sprint:** 3  ·  **Story points:** 5
**Estado de la spec:** Borrador

## 1. Objetivo
Arrancar el Planning Poker (Épica 4): abrir una ronda de estimación sobre un elemento del backlog (HU-05) y que cada integrante del proyecto (HU-03) registre su voto sin ver el de los demás. Mientras la ronda está abierta solo se sabe quién votó, no cuánto. Cuando votaron todos, la ronda se cierra. Mostrar los votos y detectar diferencias es HU-15; repetir la ronda, HU-16; y registrar la estimación acordada, HU-17.

## 2. Entradas

**Iniciar una ronda** (`POST /proyectos/{id}/backlog/{elementoID}/rondas`, sin cuerpo):

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento del backlog de ese proyecto que no esté `Hecho`. |

**Votar** (`POST /proyectos/{id}/rondas/{rondaID}/votos`):

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID`, `rondaID` | `int64` | Sí | La ronda existe y es de ese proyecto. |
| `integrante` | `string` | Sí | Uno de los votantes de la ronda (RN-14.2), sin distinguir mayúsculas ni espacios de alrededor (`domain.ClaveNombre`, como RN-03.3). |
| `valor` | `int` | Sí | Uno de 1, 2, 3, 5, 8, 13 o 21 (RN-14.4). |

```json
{ "integrante": "Ana", "valor": 5 }
```

**Consultar una ronda** (`GET /proyectos/{id}/rondas/{rondaID}`): `proyectoID` y `rondaID`.

## 3. Salidas esperadas
- Las tres operaciones devuelven el **resumen de la ronda**: `ID`, `ProyectoID`, `ElementoID`, `Numero` (1, 2, 3… por elemento), `Estado` (`Abierta` o `Cerrada`) y la lista de votantes, cada uno con su nombre y si ya votó. **Mientras la ronda está abierta, no incluye el valor de ningún voto** (RN-14.3); cerrada, HU-15 le suma los votos revelados.
- Iniciar: la ronda queda guardada en estado `Abierta`, con los integrantes del proyecto como votantes y ningún voto.
- Votar: el voto queda guardado. Si con él votaron todos, la ronda pasa a `Cerrada`.
- Si hay un error, no se guarda nada.
- **En HTTP:** iniciar responde `201`; votar y consultar, `200`; los errores, `400` o `404` (sección 7). JSON del resumen:

```json
{
  "id": 3, "proyecto_id": 1, "elemento_id": 7, "numero": 1, "estado": "Abierta",
  "votantes": [ { "nombre": "Ana", "voto": true }, { "nombre": "Beto", "voto": false } ]
}
```

- **Pantallas** (ADR-0003):
  - En la pantalla del elemento (HU-06), una sección "Planning Poker": si el elemento tiene una ronda abierta, un enlace a ella; si no, y el elemento no está `Hecho`, el botón "Iniciar ronda". Al iniciarla se va a la pantalla de la ronda (POST-redirect-GET, `303`).
  - Pantalla de la ronda (`/web/proyectos/{id}/rondas/{rondaID}`): el elemento, el número y el estado de la ronda y cada votante con "Votó" o "Falta votar". Con la ronda abierta, un formulario para elegir el integrante y el valor (solo los de RN-14.4) y votar; al votar se vuelve a la misma pantalla (`303`). Con la ronda cerrada, un aviso de que votaron todos y ningún formulario para votar (HU-17 suma el de registrar la estimación). Los errores se muestran en la misma pantalla.

## 4. Reglas de negocio
- **RN-14.1:** Una ronda se abre sobre un elemento del backlog que no está `Hecho`. Un elemento tiene a lo sumo una ronda `Abierta` a la vez; cuando se cierra, se puede abrir otra (HU-16), con el número siguiente.
- **RN-14.2:** Los votantes son los integrantes del proyecto en el momento de abrir la ronda. Si después se agrega o se quita un integrante del proyecto, la ronda no cambia. Un proyecto sin integrantes no puede abrir una ronda.
- **RN-14.3:** Mientras la ronda está abierta, nadie ve el valor de los votos: solo quién votó y quién falta. Revelarlos al cerrarse la ronda es HU-15 (RN-15.1).
- **RN-14.4:** Un voto es uno de 1, 2, 3, 5, 8, 13 o 21 (escala de Fibonacci del Planning Poker). Es la misma escala de los story points de HU-13 (`EscalaStoryPoints()`, RN-13.4).
- **RN-14.5:** Cada votante tiene un solo voto por ronda. Mientras la ronda está abierta, puede volver a votar y queda el último valor.
- **RN-14.6:** La ronda se cierra sola cuando votaron todos sus votantes. Una ronda cerrada no acepta más votos y no se reabre.
- **RN-14.7:** El elemento y la ronda tienen que pertenecer al proyecto de la ruta; si son de otro proyecto, se tratan como inexistentes (igual que RN-06.6 y RN-18.1).

## 5. Restricciones
- Las reglas RN-14.1 a RN-14.6 viven en el dominio (`internal/domain/planning_poker.go`), sin dependencias externas: `NuevaRonda` arma la ronda con los votantes (recibe el elemento y los integrantes), `Ronda.Votar` registra el voto y cierra la ronda, y `Ronda.Resumen` devuelve los datos que se pueden mostrar, sin los valores (RN-14.3). Los adaptadores muestran solo el resumen, así la regla de ocultar no se repite en la API ni en la pantalla.
- Los casos de uso `IniciarRonda`, `Votar` y `ConsultarRonda` (`internal/app/planning_poker.go`) buscan el proyecto (`LectorProyectos`, HU-04), el elemento (`LectorElementos`, HU-18) y las rondas por un puerto propio, `RepositorioRondas`.
- **Persistencia** (SQLite, ADR-0002): tabla `rondas` (proyecto, elemento, número, estado) y tabla `ronda_votantes` (ronda, nombre, clave del nombre y voto, `NULL` mientras no votó). Las dos se borran con el elemento (`ON DELETE CASCADE`): una ronda es parte de la estimación del elemento y no tiene que impedir eliminarlo (así no interfiere con RN-06.7 y RN-06.10, que bloquean por clave foránea).
- **Pedidos simultáneos:** RN-14.1 (una sola ronda abierta por elemento) también se garantiza en la base con un índice único sobre el elemento para las rondas abiertas. Votar guarda el voto y, si ya votaron todos, cierra la ronda en una sola transacción, y el voto solo se guarda si la ronda sigue abierta: si dos votantes votan a la vez, los dos votos quedan y la ronda se cierra una vez; un voto que llega después del cierre recibe `ErrRondaCerrada`.
- `net/http` (ADR-0001): cada historia registra sus rutas con su función `RegistrarX`, que se suma a `NuevoRouter` y a `router_test.go`. Las pantallas pasan por el chequeo de mismo origen de los formularios.
- Fuera de alcance: mostrar los valores y las diferencias (HU-15), reiniciar la ronda (HU-16), guardar la estimación acordada como story points (HU-17), y saber quién está votando: el sistema no tiene autenticación, así que en la pantalla se elige el integrante de una lista.

## 6. Casos límite
- Proyecto con un solo integrante → al votar él, la ronda se cierra.
- Elemento `En progreso` o `Pendiente` → se abre la ronda; `Hecho` → error.
- Iniciar una segunda ronda con otra todavía abierta → error; con la anterior cerrada → se abre la número 2.
- Votar con "ana" o " Ana " cuando el votante es "Ana" → se acepta (RN-14.2).
- Votar dos veces mientras la ronda está abierta → queda el último voto y el votante figura una vez (RN-14.5).
- Votar 4, 0 o 34 → error (RN-14.4).
- Un integrante agregado al proyecto después de abrir la ronda vota → error, no es votante (RN-14.2).
- Un integrante quitado del proyecto después de abrir la ronda → sigue siendo votante y la ronda espera su voto (RN-14.2).
- Votar en una ronda cerrada → error, aunque sea un votante que ya votó.
- Eliminar el elemento (HU-06) con una ronda abierta o cerrada → se eliminan sus rondas.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05) | "elemento de backlog no encontrado" (HTTP 404) |
| No existe la ronda, o es de otro proyecto | `domain.ErrRondaNoEncontrada` | "ronda de estimación no encontrada" (HTTP 404) |
| El elemento ya está `Hecho` | `domain.ErrElementoTerminado` | "no se puede estimar un elemento terminado" (HTTP 400) |
| El proyecto no tiene integrantes | `domain.ErrRondaSinVotantes` | "el proyecto no tiene integrantes para votar" (HTTP 400) |
| El elemento ya tiene una ronda abierta | `domain.ErrRondaYaAbierta` | "el elemento ya tiene una ronda abierta" (HTTP 400) |
| La ronda está cerrada | `domain.ErrRondaCerrada` | "la ronda ya está cerrada" (HTTP 400) |
| El integrante no es votante de la ronda | `domain.ErrNoEsVotante` | "el integrante no participa de esta ronda" (HTTP 400) |
| El valor no es de la escala | `domain.ErrVotoInvalido` | "el voto debe ser 1, 2, 3, 5, 8, 13 o 21" (HTTP 400) |
| Cuerpo JSON mal formado o con campos desconocidos | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso. |

Si hay varios errores, se informa el primero en este orden. Al iniciar: proyecto, elemento, elemento terminado, sin integrantes, ronda ya abierta. Al votar: proyecto, ronda, ronda cerrada, votante, valor.

## 8. Criterios de aceptación
- **CA-14.1:** Dado un proyecto con los integrantes Ana y Beto y el elemento "Login" Pendiente, cuando inicio una ronda sobre "Login", entonces queda abierta, con el número 1, Ana y Beto como votantes y ninguno votó. (RN-14.1, RN-14.2)
- **CA-14.2:** Dada una ronda abierta, cuando Ana vota 5, entonces la ronda sigue abierta y el resumen dice que Ana votó y Beto no, sin mostrar el valor. (RN-14.3, RN-14.6)
- **CA-14.3:** Dada una ronda en la que Ana ya votó, cuando Beto vota, entonces la ronda queda cerrada y ya no acepta votos. (RN-14.6)
- **CA-14.4:** Dada una ronda abierta, cuando Ana vota 5 y después 8, entonces Ana figura una sola vez como que votó y la ronda sigue abierta. (RN-14.5)
- **CA-14.5:** Cuando voto un valor que no es 1, 2, 3, 5, 8, 13 ni 21, entonces obtengo `ErrVotoInvalido`; cuando vota alguien que no es votante de la ronda (aunque se haya sumado después al proyecto), `ErrNoEsVotante`; y en una ronda cerrada, `ErrRondaCerrada`. En los tres casos no se guarda nada. (RN-14.2, RN-14.4, RN-14.6)
- **CA-14.6:** No se puede iniciar una ronda sobre un elemento `Hecho` (`ErrElementoTerminado`), en un proyecto sin integrantes (`ErrRondaSinVotantes`) ni si el elemento ya tiene una ronda abierta (`ErrRondaYaAbierta`). Con la ronda anterior cerrada, se abre la número 2. (RN-14.1, RN-14.2)
- **CA-14.7:** Dado un proyecto, un elemento o una ronda inexistentes o de otro proyecto, entonces obtengo `ErrProyectoNoEncontrado`, `ErrElementoNoEncontrado` o `ErrRondaNoEncontrada`. (RN-14.7)
- **CA-14.8:** La API responde `201` al iniciar, `200` al votar y al consultar, y mientras la ronda está abierta ninguna respuesta trae el valor de un voto; los errores de validación dan `400` y las entidades inexistentes, `404`. (Sección 3, RN-14.3)
- **CA-14.9:** Desde la pantalla del elemento inicio una ronda y llego a su pantalla; ahí veo quién votó y quién falta, voto eligiendo integrante y valor, y mientras la ronda está abierta no veo el valor de ningún voto. Con la ronda cerrada no hay formulario para votar. Los errores se muestran en la misma pantalla. (Sección 3)

## Decisiones tomadas
Definidas por Nicolás Andino al refinar la historia (2026-11-01):
- Los votantes son los integrantes del proyecto al abrir la ronda (RN-14.2). HU-03 no tiene integrantes "inactivos", así que "integrantes activos" del issue se interpreta así.
- Escala de Fibonacci 1, 2, 3, 5, 8, 13, 21 (RN-14.4), alineada con HU-13.
- Se puede cambiar el voto mientras la ronda está abierta (RN-14.5).
- Solo sobre elementos que no están `Hecho`, con una ronda abierta a la vez por elemento (RN-14.1).

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-11-01: creada (borrador).
- 2026-11-01: con HU-15, el resumen oculta los valores mientras la ronda está abierta y los revela al cerrarse (RN-14.3, CA-14.8 y CA-14.9 ajustados).
- 2026-11-01: la escala de los votos es la de los story points de HU-13 (`EscalaStoryPoints()`); se quitó `ValoresVoto`.
- 2026-11-02: aclarado, por HU-17, que la ronda cerrada no tiene formulario para votar (sí el de registrar la estimación).
