# Sprint 4: Acta de Retrospectiva

> **Borrador para validar con el equipo.** Los "Hechos del sprint" salen del repositorio y del tablero. Las opiniones y las acciones son una propuesta armada a partir de esos hechos: cada integrante las revisa y las corrige en la presentación final, donde se hace la retrospectiva con los profesores.

**Fecha:** [completar: día de la presentación]  ·  **Facilitador (Agile Enabler):** Matías López  ·  **Asistentes:** Matías López (Agile Enabler), Nicolás Andino, Ezequiel Assenza y Lucas Orellana
**Dinámica usada:** Start / Stop / Continue, con los profesores

## Hechos del sprint (para la conversación)
- El sprint se extendió dos veces (del 8 al 10 y del 10 al 13/11). Al 8/11 había 23 de 28 puntos terminados, y al 11/11, 26. Se cerró el 12/11 con los 28, un día antes del fin.
- HU-25, HU-26, HU-27 y HU-31 tocaban los mismos archivos (`main.go`, `router.go`, `bdd_test.go` y la trazabilidad): cada merge dejó en conflicto a los PR siguientes, y HU-27 tuvo que esperar a HU-26 porque estaba armada sobre su rama.
- En tres PR (HU-24, HU-25 y HU-26), los `Closes` no cerraron las tareas, porque GitHub solo vincula el primer número de una lista, y hubo que cerrarlas a mano o corregir la descripción.
- Las revisiones encontraron problemas reales antes de mergear: el redondeo del porcentaje en HU-24, el aviso de horas sin estimar en HU-25, la aclaración del sprint cerrado en HU-26, la carga a medias del ejemplo en HU-31 y DEF-05 en la revisión de HU-27.
- La suite de tests llegó a tardar casi 9 minutos en el CI, cerca del límite de 10, y se subió el timeout (#376).
- El despliegue en Render quedó activo el 5/11, y desde entonces cada merge se despliega solo.
- Los dos defectos pendientes (DEF-04 y DEF-05) se corrigieron antes del cierre, con TDD, y el sprint terminó sin defectos abiertos.
- La review y la retro no se hicieron en una reunión aparte: se dejan para la presentación final con los profesores.

## Qué salió bien
- Terminamos las 7 historias (28 de 28 puntos) y cerramos el proyecto sin defectos abiertos.
- Las revisiones de PR siguieron encontrando problemas reales antes del merge (el redondeo de HU-24, el aviso de horas de HU-25, la carga a medias de HU-31 y DEF-05).
- Con el despliegue continuo, cada merge quedaba publicado, y el proyecto de ejemplo nos deja mostrar la aplicación completa.
- La trazabilidad y el informe de métricas quedaron al día con el cierre.

## Qué salió mal / a mejorar
- El sprint se extendió dos veces: la documentación y las revisiones llevaron más de lo estimado.
- HU-25, HU-26, HU-27 y HU-31 tocaban los mismos archivos, y cada merge dejaba en conflicto a los PR siguientes; HU-27 tuvo que esperar a HU-26.
- Los `Closes` con varios números no cerraron las tareas, y hubo que cerrarlas a mano.
- La suite de tests creció hasta casi 9 minutos en el CI.

## Acciones de mejora
| Acción | Responsable | Para cuándo | Estado |
|---|---|---|---|
| Escribir un `closes #N` por cada issue en la descripción del PR | Todos | Desde ya | Cumplida (desde el PR #381) |
| Subir el timeout de los tests en el CI | Matías | 8/11 | Cumplida (#376) |
| Para un próximo proyecto: separar el registro de rutas y de steps BDD por historia, para que los PR no choquen en los mismos archivos | Todos | Próximo proyecto | Pendiente |
| Para un próximo proyecto: estimar la documentación y las revisiones como trabajo del sprint | Todos | Próximo proyecto | Pendiente |

## Seguimiento de acciones de la retro anterior
| Acción | ¿Se cumplió? | Comentario |
|---|---|---|
| Mergear en el orden de las dependencias y actualizar las ramas enseguida | Parcial | Se avisaron los conflictos en las revisiones (HU-26, HU-31), pero siguieron apareciendo. |
| Corregir DEF-04 | Sí | Resuelto el 12/11 con un ciclo RED → GREEN (#382). |
| Despliegue continuo en Render | Sí | Activo desde el 5/11. |
| Estimar el Sprint 4 con margen | Sí | Se comprometieron 25 puntos y, con HU-31, se terminaron 28; igual hubo que extender el sprint. |
