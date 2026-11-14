# Sprint 2: Acta de Retrospectiva

> **Borrador para validar con el equipo.** Los "Hechos del sprint" salen del repositorio y del tablero. Las opiniones y las acciones son una propuesta armada a partir de esos hechos: cada integrante las revisa y las corrige en la presentación final, donde se hace la retrospectiva con los profesores.

**Fecha:** 29 de octubre de 2026  ·  **Facilitador (Agile Enabler):** Matías López  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Dinámica usada:** Start / Stop / Continue

## Hechos del sprint (para la conversación)
- Se terminaron las 8 historias (26 SP, contra 19 del Sprint 1), pero dos días después de lo planificado: el 26 y 27/10 se cerraron 4 historias, y las otras 4 entre el 28 y el 29/10.
- Las historias que dependían de otras (HU-09 → HU-10 → HU-11) se encadenaron, y cada merge generó conflictos en los mismos archivos (`router.go`, `pantallas.go`, `main.go`, `bdd_test.go`, `trazabilidad.md`).
- Varias decisiones de las specs (HU-10, HU-11, HU-20) las tomó el autor y se confirmaron recién en la revisión del PR; en HU-11 eso obligó a cambiar código ya hecho.
- Los PR #261 y #262 se mergearon sin una aprobación registrada en GitHub. Desde el 29/10 `main` está protegida: exige PR, una aprobación y el CI en verde.
- Al mergear, no siempre se cerraron los issues: HU-09 y las tareas de HU-30 quedaron abiertas en Backlog hasta el 29/10.
- Se detectaron y corrigieron 2 defectos, los dos en la revisión de PR.
- Algunos commits no siguieron la convención (`docs(spec)` o `test(bdd)` en lugar de `docs(HU-XX)`).

## Qué salió bien
- Subimos de 19 a 26 puntos, y cada historia ya tiene su pantalla.
- Las revisiones de PR encontraron los dos defectos del sprint antes de que llegaran a la entrega.
- Los tests de concurrencia nos obligaron a pensar los pedidos simultáneos (HU-08 a HU-11), y quedaron resueltos en la base de datos.

## Qué salió mal / a mejorar
- Terminamos dos días después de lo planificado: las historias encadenadas (HU-09 → HU-10 → HU-11) se esperaban unas a otras.
- Cada merge dejaba en conflicto a los PR siguientes, porque todos tocábamos `router.go`, `main.go`, `bdd_test.go` y la trazabilidad.
- Algunas decisiones de las specs se confirmaron recién en la revisión del PR, y en HU-11 hubo que rehacer código.
- Dos PR se mergearon sin una aprobación registrada, y no siempre se cerraron los issues al mergear.

## Acciones de mejora
| Acción | Responsable | Para cuándo | Estado |
|---|---|---|---|
| Proteger `main`: PR, una aprobación y el CI en verde | Matías | 29/10 | Cumplida (29/10) |
| Confirmar las preguntas abiertas de la spec antes de implementar, no en la revisión del PR | Todos | Sprint 3 | Parcial |
| Avisar en el grupo antes de mergear un PR que toca archivos compartidos, para que los demás actualicen su rama | Todos | Sprint 3 | Parcial |
| Usar la convención de commits `tipo(HU-XX)` en todos los commits | Todos | Sprint 3 | Parcial |

## Seguimiento de acciones de la retro anterior
| Acción | ¿Se cumplió? | Comentario |
|---|---|---|
| Una historia por PR | Sí | Cada historia del Sprint 2 tuvo su propio PR (#180 a #187). |
| Cargar las tareas en el tablero antes de empezar | Parcial | Se cumplió en el Sprint 3, con las tareas como sub-issues de cada historia. |
| Interfaz web base (HU-30) | Sí | Mergeada el 26/10 (#106); desde ahí cada historia suma su pantalla. |
