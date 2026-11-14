# Sprint 1: Acta de Retrospectiva

> **Borrador para validar con el equipo.** Los "Hechos del sprint" salen del repositorio y del tablero. Las opiniones y las acciones son una propuesta armada a partir de esos hechos: cada integrante las revisa y las corrige en la presentación final, donde se hace la retrospectiva con los profesores.

**Fecha:** 25 de octubre de 2026  ·  **Facilitador (Agile Enabler):** Matías López  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Dinámica usada:** Start / Stop / Continue

## Hechos del sprint (para la conversación)
- Se terminaron las 6 historias (19 SP), incluida HU-05, que se sumó con el sprint en marcha.
- HU-18 se mergeó el 26/10, con el Sprint 2 ya empezado.
- Los ADR de framework web y persistencia se tomaron al principio y no hubo que cambiarlos.
- Las tareas de HU-05 se cargaron en el tablero después de hacerlas, a partir de los commits.
- Hubo PR con varias historias juntas (por ejemplo, #86).

## Qué salió bien
- Tomar las decisiones de arquitectura al principio (ADR-0001 y ADR-0002) nos dejó trabajar en paralelo sin pisarnos.
- El flujo spec → BDD → TDD funcionó desde la primera historia, y el CI atajó errores de formato y de tests antes del merge.
- Pudimos sumar HU-05 con el sprint en marcha y terminarla igual.

## Qué salió mal / a mejorar
- Algunos PR juntaron más de una historia (#86), y eso hizo más difícil revisarlos.
- Las tareas de HU-05 se cargaron en el tablero después de hacerlas, así que el tablero no mostró el avance real.
- HU-18 se mergeó un día después del cierre.
- Sin interfaz web, la demo se hizo con la API, y cuesta mostrarla.

## Acciones de mejora
| Acción | Responsable | Para cuándo | Estado |
|---|---|---|---|
| Una historia por PR | Todos | Sprint 2 | Cumplida |
| Cargar las tareas en el tablero antes de empezar cada historia | Matías | Sprint 2 | Cumplida en el Sprint 3 (sub-issues) |
| Sumar una interfaz web base para operar el sistema desde el navegador (HU-30) | Nico | Sprint 2 | Cumplida |

## Seguimiento de acciones de la retro anterior
No aplica: es la primera retrospectiva.
