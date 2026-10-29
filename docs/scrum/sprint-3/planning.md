# Sprint 3 — Planning

**Fechas:** 29 de octubre (noche) y 30 de octubre de 2026. El Sprint 2 se terminó de mergear el 29 a la noche, un día después de lo previsto en el calendario del Product Backlog. Las fechas del Sprint 3 y del Sprint 4 se confirman en el grupo.
**Sprint Goal:** estimar historias en equipo con Planning Poker y medir el trabajo: velocidad, esfuerzo estimado contra real y estado de los defectos.
**Compromiso:** 28 story points (velocidad del Sprint 2: 26 puntos, todas las historias terminadas).

## Sprint Backlog

El orden es la prioridad: si al final del sprint no se llega con todo, quedan afuera las últimas.

| # | Historia | Pts | Depende de | Responsable |
|---|---|---|---|---|
| 1 | HU-22 Story points planificados, completados y velocidad por sprint | 5 | HU-11 ✓ | |
| 2 | HU-14 Registrar estimaciones individuales ocultas (Planning Poker) | 5 | HU-03 ✓, HU-05 ✓ | |
| 3 | HU-15 Mostrar las estimaciones y detectar diferencias | 3 | HU-14 | |
| 4 | HU-16 Repetir una ronda de votación | 2 | HU-15 | |
| 5 | HU-17 Registrar la estimación acordada | 2 | HU-16 | |
| 6 | HU-19 Comparar esfuerzo estimado contra real | 3 | HU-18 ✓ | |
| 7 | HU-23 Horas estimadas, reales y desviación | 3 | HU-19 | |
| 8 | HU-21 Actualizar el estado de un defecto | 2 | HU-20 ✓ | |
| 9 | HU-12 Consultar sprints anteriores | 2 | HU-11 ✓ | |
| 10 | HU-13 Escala de story points (Fibonacci) | 1 | HU-05 ✓, HU-06 ✓ | |
| | **Total** | **28** | | |

Cada historia tiene sus tareas como sub-issues en GitHub, con el mismo formato que en el Sprint 2: spec, escenarios BDD, steps, TDD del dominio, TDD del caso de uso, persistencia, HTTP, pantalla, y revisión y trazabilidad.

## Cómo repartir

- **HU-14 a HU-17 son el Planning Poker** y tocan los mismos archivos. Conviene que las haga una sola persona, o dos en orden mergeando seguido.
- **HU-19 y HU-23** van juntas: HU-23 reutiliza el cálculo de desviación de HU-19.
- **HU-22, HU-21, HU-12 y HU-13** son independientes entre sí y se pueden arrancar ya.

## Preguntas que tiene que resolver la spec antes de programar

- **HU-19 y HU-23:** hoy el sistema no guarda horas estimadas: HU-05 tiene story points y HU-18 registra horas reales. Hay que decidir de dónde salen las horas estimadas (por ejemplo, un campo nuevo en el elemento del backlog).
- **HU-23:** al cerrar un sprint, las historias sin terminar vuelven al backlog (HU-11). Hay que decidir si sus horas cuentan para el sprint en el que se trabajaron.
- **HU-13:** HU-05 y HU-06 ya aceptan cualquier entero positivo como story points. Hay que definir la escala y qué pasa con los elementos ya guardados con un valor fuera de ella. Afecta a HU-17.
- **HU-15:** cuál es el umbral de "diferencia significativa".
- **HU-21:** no existe una lista de defectos. La pantalla tiene que sumarla para poder elegir el defecto.
- **HU-22:** usar los puntos que HU-11 guarda al cerrar cada sprint (RN-11.6), no recalcularlos desde las historias.

## Reglas para no bloquearnos

Las mismas del Sprint 2:
1. Antes de empezar, `git checkout main` y `git pull`; la rama de la historia sale de ahí.
2. Quien depende de otra historia arranca por lo que no depende de nadie: la spec, los escenarios BDD y el dominio.
3. PR chicos y seguido.
4. Cada historia incluye su pantalla (ADR-0003).

## Ceremonias

- **Planning:** este documento. Cada integrante elige sus historias y completa la columna *Responsable*.
- **Daily:** un mensaje en el grupo al terminar el día: qué hice, qué hago y si estoy trabado.
- **Review y retrospectiva:** al cerrar el sprint, unos 30 minutos. Pendiente también la del Sprint 2.
