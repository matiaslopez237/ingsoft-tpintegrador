# HU-22: Calcular story points y velocidad del equipo

**Historia:** Como Scrum Master, quiero ver los story points planificados, completados y la velocidad del equipo por Sprint, para medir el rendimiento.
**Issue:** #33  ·  **Sprint:** 3  ·  **Story points:** 5
**Estado de la spec:** Borrador

## 1. Objetivo
Dar al Scrum Master una medida de cuánto entrega el equipo por Sprint: los puntos planificados contra los completados de cada Sprint, y la velocidad del equipo (el promedio de puntos completados en los Sprints ya cerrados), para planificar sprints futuros con un número real en vez de una estimación a ojo.

**Alcance real de esta historia:** HU-11 (cierre de sprint) ya calcula y guarda los puntos planificados y completados de cada Sprint al cerrarlo (RN-11.6), y HU-12 (sprints anteriores) ya los muestra en su pantalla y su endpoint JSON. Lo único que falta, y es lo que agrega HU-22, es la **velocidad del equipo**: el promedio de puntos completados entre los Sprints Cerrados del proyecto. No se agrega persistencia nueva: se reutilizan los datos que ya guarda HU-11.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |

Es una consulta de solo lectura: no hay más entradas. En HTTP, `proyectoID` va en la ruta (`GET /proyectos/{id}/velocidad`).

## 3. Salidas esperadas
- La **velocidad**: el promedio de `PuntosCompletados` de todos los Sprints `Cerrado` del proyecto (RN-11.6), con decimales (no se trunca a entero).
- La **cantidad de Sprints cerrados** que entraron en el promedio, para que quien lo lea sepa sobre cuántos datos se calculó.
- Si el proyecto no tiene ningún Sprint `Cerrado`, no hay velocidad: se informa explícitamente que no hay datos todavía, nunca el valor `0` (RN-22.3).
- Los puntos planificados y completados de un Sprint en particular **no son una salida nueva de HU-22**: ya los devuelve `GET /proyectos/{id}/sprints/anteriores` y los muestra la pantalla de HU-12 (sección 1).
- **En HTTP** (`GET /proyectos/{id}/velocidad`): `200` con `{"velocidad": 6, "sprints_cerrados": 3}` (Go serializa un `float64` entero sin el punto decimal), o `{"velocidad": null, "sprints_cerrados": 0}` si no hay Sprints cerrados; `400` si el ID no es un número; `404` si el proyecto no existe.
- **En la pantalla:** se agrega un dato "Velocidad del equipo" en la pantalla ya existente de sprints anteriores (HU-12, `/web/proyectos/{id}/sprints/anteriores`), antes de la tabla. No hay pantalla nueva.

## 4. Reglas de negocio
- **RN-22.1:** Los story points planificados y completados de un Sprint son los que HU-11 calculó y guardó al cerrarlo (RN-11.6). HU-22 no los recalcula ni los vuelve a sumar desde las historias: los lee tal cual están guardados en el Sprint.
- **RN-22.2:** La velocidad del equipo de un proyecto es el promedio de `PuntosCompletados` de **todos** los Sprints en estado `Cerrado` del proyecto (sin límite de cantidad ni ventana de los "últimos N"). Un Sprint `Abierto` no cuenta (sus puntos valen 0 por RN-11.6, pero además no está cerrado).
- **RN-22.3:** Si el proyecto no tiene ningún Sprint `Cerrado`, la velocidad queda **indefinida**: no se informa como `0`, porque `0` sería un promedio real (de un Sprint cerrado con 0 puntos completados), un dato distinto de "no hay Sprints cerrados todavía".
- **RN-22.4:** El promedio se calcula con precisión decimal (`float64`), sin redondear a entero, porque los puntos completados de los Sprints cerrados casi nunca dividen exacto.

## 5. Restricciones
- El cálculo vive en el dominio (`internal/domain/sprint_historial.go`, junto a `SprintsAnteriores` de HU-12), como una función pura sin dependencias externas: `Velocidad(cerrados []Sprint) (velocidad float64, hayDatos bool)`. Recibe los Sprints ya filtrados por `SprintsAnteriores` (HU-12): HU-22 no vuelve a filtrar por estado.
- Caso de uso nuevo `ConsultarVelocidad` (`internal/app/consultar_velocidad.go`), con la misma forma que `ListarSprintsAnteriores` (HU-12): busca el proyecto (`LectorProyectos`, reutilizado), lista sus Sprints (`RepositorioSprints.ListarSprints`, reutilizado de HU-08), los filtra con `domain.SprintsAnteriores` y calcula la velocidad con `domain.Velocidad`. **No agrega ningún puerto de persistencia nuevo.**
- La pantalla reutiliza el handler de sprints anteriores de HU-12 (`pantallas_sprints.go`): agrega el llamado a `ConsultarVelocidad` y un dato más a `datosSprintsAnteriores`, sin crear una pantalla ni una plantilla nueva.
- Endpoint JSON propio (`GET /proyectos/{id}/velocidad`), separado del de HU-12 (`GET /proyectos/{id}/sprints/anteriores`), para no tocar el contrato que ya expone esa historia.
- Persistencia en SQLite (ADR-0002), pero sin cambios de esquema: no hay tabla ni columna nueva.

## 6. Casos límite
- Proyecto sin ningún Sprint (ni abiertos ni cerrados) → velocidad indefinida, igual que sin Sprints cerrados.
- Proyecto con Sprints, pero todos `Abierto` (ninguno cerrado todavía) → velocidad indefinida.
- Un solo Sprint cerrado con 0 puntos completados → velocidad = 0 (es un dato real: hay 1 Sprint cerrado y completó 0 puntos; distinto del caso "sin datos").
- Sprints cerrados con puntos completados que no dividen exacto (ej. 5 y 4) → velocidad = 4,5, con decimales.
- Sprints cerrados con algunos en 0 puntos completados (por ejemplo, un Sprint donde no se terminó nada) → esos ceros entran en el promedio igual, bajan la velocidad; no se excluyen.
- Proyecto con muchos Sprints cerrados (10 o más) → el promedio los toma a todos, sin ventana.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| ID de proyecto no numérico o no positivo | — (adaptador HTTP, igual que el resto de las historias) | HTTP 400; no se llama al caso de uso. |

No hay más condiciones de error: es una consulta de solo lectura sin más datos de entrada que el ID del proyecto.

## 8. Criterios de aceptación
- **CA-22.1:** Dado un Sprint `Cerrado` con 8 puntos planificados y 5 completados (ya guardados por HU-11 al cerrarlo), cuando se consulta, entonces se informan esos mismos valores sin recalcularlos. (RN-22.1; ya lo cubren el endpoint y la pantalla de HU-12 — HU-22 no agrega código para esto, solo lo señala como ya cumplido.)
- **CA-22.2:** Dado un proyecto con 3 Sprints cerrados que completaron 8, 6 y 4 puntos, cuando se consulta la velocidad, entonces es 6 (el promedio) sobre 3 Sprints. (RN-22.2)
- **CA-22.3:** Dado un proyecto sin ningún Sprint `Cerrado`, cuando se consulta la velocidad, entonces se informa que no hay datos todavía, no `0`. (RN-22.3)
- **CA-22.4:** Dado un proyecto con Sprints cerrados que completaron 5 y 4 puntos, cuando se consulta la velocidad, entonces es 4,5. (RN-22.4)
- **CA-22.5:** Dado un Sprint `Cerrado` con 0 puntos completados como único Sprint cerrado del proyecto, cuando se consulta la velocidad, entonces es 0 (no "sin datos"). (RN-22.2, RN-22.3)
- **CA-22.6:** Dado un ID de proyecto que no existe, cuando se consulta la velocidad, entonces se obtiene `ErrProyectoNoEncontrado`. (HTTP 404)
- **CA-22.7:** Dado un pedido válido a `GET /proyectos/{id}/velocidad`, entonces la respuesta es `200` con `velocidad` y `sprints_cerrados`; con un proyecto sin Sprints cerrados, `velocidad` es `null` y `sprints_cerrados` es `0`; con un ID inválido, `400`; con un proyecto inexistente, `404`.
- **CA-22.8:** Dada la pantalla de sprints anteriores (HU-12) de un proyecto con Sprints cerrados, cuando se abre, entonces muestra "Velocidad del equipo" con el promedio; con un proyecto sin Sprints cerrados, muestra que todavía no hay datos de velocidad (sin romper la pantalla existente de HU-12, que ya avisa "El proyecto todavía no tiene sprints cerrados").

## Preguntas abiertas
- **P-1 — Ventana del promedio:** se decidió con el equipo (Ezequiel) en el chat que la velocidad promedia **todos** los Sprints cerrados, sin ventana de "últimos N" (más simple y no inventa un número arbitrario). Pendiente de confirmar con el resto del equipo si en algún momento se prefiere una ventana móvil (ej. últimos 3), como es más común en Scrum.
- **P-2 — Redondeo para mostrar en pantalla:** la spec no fija cuántos decimales se muestran en la pantalla (ej. "4,5" vs "4,50"); queda como detalle de implementación de la plantilla, no una regla de negocio.

## Historial
- 2026-11-01: creada (borrador). Decisiones P-1 (ventana de todos los Sprints cerrados), dónde se muestra (pantalla de HU-12, sin pantalla nueva) y qué pasa sin Sprints cerrados (indefinido, no 0) confirmadas con Ezequiel en el chat.
