# Sprint 3: Acta de Retrospectiva

> **Borrador para validar con el equipo.** Los "Hechos del sprint" salen del repositorio y del tablero. Las opiniones y las acciones son una propuesta armada a partir de esos hechos: cada integrante las revisa y las corrige en la presentación final, donde se hace la retrospectiva con los profesores.

**Fecha:** 5 de noviembre de 2026  ·  **Facilitador (Agile Enabler):** Matías López  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Dinámica usada:** Start / Stop / Continue

## Hechos del sprint (para la conversación)
- Se terminaron las 10 historias (28 SP, la velocidad más alta hasta ahora), pero ninguna dentro de las fechas planificadas (29 al 31/10): se mergearon entre el 1 y el 5/11.
- Las historias del Planning Poker (HU-14 a HU-17) dependían unas de otras y tocaban los mismos archivos; HU-17 entró antes que HU-13 y hubo que adaptar HU-13 al mergear.
- Varias decisiones de las specs las tomaron los autores y se confirmaron en la revisión del PR (HU-19, HU-21, HU-22, HU-23).
- Desde el 29/10 `main` está protegida, y desde el 1/11 el CI corre con `-race` y una cobertura mínima del 80 %.
- Se detectaron y corrigieron 2 defectos (DEF-02 y DEF-03), los dos en la revisión de PR; un tercero (DEF-04) quedó para el Sprint 4.

## Qué salió bien
- Llegamos a 28 puntos, la velocidad más alta, y completamos todo lo planificado.
- El CI quedó completo: `-race`, cobertura mínima del 80 % y `main` protegida. Eso nos dio confianza para mergear.
- Las tareas cargadas como sub-issues hicieron que el tablero mostrara el avance real.

## Qué salió mal / a mejorar
- Ninguna historia entró en las fechas del sprint (29 al 31/10): se mergearon entre el 1 y el 5/11.
- Las historias del Planning Poker dependían unas de otras, y HU-13 tuvo que adaptarse porque HU-17 entró antes.
- Las decisiones de las specs se siguieron confirmando en la revisión del PR (HU-19, HU-21, HU-22 y HU-23).
- DEF-04 se detectó al agregar una migración, pero no se corrigió en ese momento.

## Acciones de mejora
| Acción | Responsable | Para cuándo | Estado |
|---|---|---|---|
| Mergear las historias en el orden de sus dependencias, y actualizar las ramas apenas entra un merge | Todos | Sprint 4 | Parcial |
| Corregir DEF-04 con TDD | Matías | Sprint 4 | Cumplida (12/11, #382) |
| Activar el despliegue continuo en Render | Matías | Sprint 4 | Cumplida (5/11) |
| Estimar el Sprint 4 con margen, porque incluye documentación además de código | Todos | Planning del Sprint 4 | Cumplida (25 puntos, por debajo de la velocidad) |

## Seguimiento de acciones de la retro anterior
| Acción | ¿Se cumplió? | Comentario |
|---|---|---|
| Proteger `main` | Sí | Desde el 29/10 exige PR, una aprobación y el CI en verde. |
| Confirmar las preguntas de la spec antes de implementar | No del todo | Varias decisiones se siguieron confirmando en la revisión del PR. |
| Avisar antes de mergear archivos compartidos | Parcial | Se avisó en algunos casos, pero los conflictos siguieron, sobre todo en el Planning Poker. |
| Convención de commits `tipo(HU-XX)` | Parcial | La mayoría la sigue, pero todavía hay commits como `docs(trazabilidad)` o `test(bdd)`. |
