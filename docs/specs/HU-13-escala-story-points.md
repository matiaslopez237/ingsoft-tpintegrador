# HU-13: Estimar una historia con story points de la escala Fibonacci

**Historia:** Como Product Builder, quiero asignar story points a una historia, para reflejar su tamaño relativo.
**Issue:** #24  ·  **Sprint:** 3  ·  **Story points:** 1
**Estado de la spec:** Borrador

## 1. Objetivo
Que los story points de un elemento del backlog sean un tamaño relativo de la escala de Fibonacci, la misma que usa el Planning Poker (HU-14), en lugar de cualquier entero positivo. HU-05 (crear) y HU-06 (editar) ya permiten cargarlos; esta historia restringe sus valores y hace que las pantallas los ofrezcan como una lista. Así la estimación acordada en Planning Poker (HU-17) siempre es un valor válido.

## 2. Entradas
No hay entradas nuevas: son los story points de HU-05 (`POST /proyectos/{id}/backlog` y el alta en la pantalla del backlog) y de HU-06 (`PATCH /proyectos/{id}/backlog/{elementoID}` y la pantalla del elemento).

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `story_points` | `*int` | No | Si se envía un número: 1, 2, 3, 5, 8, 13 o 21 (RN-13.1). Sin enviar, o `null` al editar: sin estimar, como hasta ahora. |

## 3. Salidas esperadas
- Con un valor de la escala, el elemento se crea o se edita igual que en HU-05 y HU-06.
- Con un valor fuera de la escala, la operación falla con `ErrStoryPointsInvalidos` y no se guarda nada (HTTP 400; en las pantallas, el error en el mismo formulario).
- **Pantallas:** en el alta (pantalla del backlog, HU-07) y en la edición (pantalla del elemento, HU-06), los story points se eligen de una lista: "Sin estimar" y los valores de la escala. Si el elemento tenía un valor anterior fuera de la escala, la lista lo incluye como seleccionado, marcado "(fuera de la escala)", para que guardar sin tocarlo no lo cambie.

## 4. Reglas de negocio
- **RN-13.1:** Los story points de un elemento son uno de 1, 2, 3, 5, 8, 13 o 21 (escala de Fibonacci, la misma de RN-14.4). Reemplaza la regla "entero positivo" de RN-05.6.
- **RN-13.2:** Un elemento puede seguir sin estimar (RN-05.6, RN-06.3).
- **RN-13.3:** Los elementos guardados antes de esta historia con un valor fuera de la escala (por ejemplo, 4) se respetan: se leen y se muestran igual, y se pueden editar otros campos sin tocar los story points. La escala se valida solo cuando se envían story points (al crear, o al editar enviando `story_points`).
- **RN-13.4:** La escala es una sola en el dominio: la usan la validación de los story points y las listas de las pantallas.

## 5. Restricciones
- RN-13.1 vive en el dominio (`internal/domain/backlog.go`): `NuevoElementoBacklog` valida la escala y `EscalaStoryPoints()` devuelve sus valores. `ElementoBacklog.Editar` valida los story points solo si se envían (RN-13.3).
- Se reutiliza `ErrStoryPointsInvalidos` de HU-05 con el mensaje nuevo, "los story points deben ser 1, 2, 3, 5, 8, 13 o 21", en vez de crear otro error: es la misma regla con otros valores, y así los adaptadores de HU-05 y HU-06 ya la responden con 400 sin cambios.
- Casos de uso y persistencia: sin cambios. Crear y editar ya pasan por el dominio, y la base guarda un entero.
- **Datos existentes:** no se migran ni se redondean (RN-13.3).
- **HU-14 y HU-17:** los votos del Planning Poker (HU-14, RN-14.4) usan la misma `EscalaStoryPoints()`: hay una sola escala en el dominio (RN-13.4). HU-17 guarda la estimación acordada como story points, así que siempre será un valor de la escala.
- Fuera de alcance: otras escalas (talles de remera, potencias de 2) y cambiar la escala por proyecto.

## 6. Casos límite
- 1 y 21, los extremos de la escala → válidos.
- 4, 6, 10 o 34 → error.
- 0 y negativos → error, igual que antes.
- Elemento guardado con 4: editar su título → se guarda y sigue con 4. Enviar `story_points: 4` al editarlo → error (si se envía, tiene que ser de la escala).
- Elemento con 4: quitarle la estimación (`null`) o cambiarla a 5 → válido.
- En la pantalla, un elemento con 4 muestra "4 (fuera de la escala)" seleccionado; guardar sin cambiar nada no lo modifica.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Story points fuera de la escala (incluye 0 y negativos) | `domain.ErrStoryPointsInvalidos` (HU-05, mensaje nuevo) | "los story points deben ser 1, 2, 3, 5, 8, 13 o 21" (HTTP 400) |

El orden de los errores de HU-05 y HU-06 no cambia.

## 8. Criterios de aceptación
- **CA-13.1:** Cuando creo un elemento con 1, 2, 3, 5, 8, 13 o 21 story points, entonces queda guardado con ese valor. (RN-13.1)
- **CA-13.2:** Cuando creo un elemento con 4, 10, 34 o 0 story points, entonces obtengo `ErrStoryPointsInvalidos` y no se crea. (RN-13.1)
- **CA-13.3:** Dado un elemento con 5 story points, cuando los edito a 13 se guardan; cuando los edito a 10, obtengo `ErrStoryPointsInvalidos` y sigue con 5. (RN-13.1)
- **CA-13.4:** Un elemento se puede crear sin estimar y quitarle la estimación al editarlo. (RN-13.2)
- **CA-13.5:** Dado un elemento guardado antes con 4 story points, cuando edito su título, entonces se guarda y conserva los 4; cuando edito sus story points a 4, obtengo `ErrStoryPointsInvalidos`. (RN-13.3)
- **CA-13.6:** La API responde 400 con el mensaje de la escala al crear o editar con un valor fuera de ella. (Sección 3)
- **CA-13.7:** En las pantallas de alta y de edición, los story points se eligen de una lista con "Sin estimar" y los valores de la escala; un valor anterior fuera de la escala aparece seleccionado y marcado. (Sección 3, RN-13.3)

## Decisiones tomadas
Definidas por Nicolás Andino al refinar la historia (2026-11-01):
- Escala de Fibonacci 1, 2, 3, 5, 8, 13, 21, la misma del Planning Poker (RN-13.1).
- Los valores guardados antes fuera de la escala se respetan, sin migrarlos (RN-13.3).
- Se arranca desde `main`, sin esperar a HU-14; la escala de los votos se unifica cuando estén las dos (RN-13.4).

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-11-01: creada (borrador).
- 2026-11-01: con HU-14 en `main`, los votos usan `EscalaStoryPoints()` y se quitó la escala propia de HU-14 (RN-13.4), a pedido de la revisión del PR #288.
