# Sprint 2 — Planning

**Fechas:** 26 y 27 de octubre de 2026 (2 días)
**Sprint Goal:** gestionar el backlog y armar sprints de un proyecto desde una interfaz web.
**Compromiso:** 26 story points (velocidad del Sprint 1: 19 puntos).

## Sprint Backlog

El orden es la prioridad: si al final del sprint no se llega con todo, quedan afuera las últimas.

| # | Historia | Pts | Depende de | Responsable |
|---|---|---|---|---|
| 1 | HU-30 Interfaz web base: diseño, navegación y pantallas de proyectos e integrantes (ADR-0003) | 5 | — | |
| 2 | HU-08 Crear un sprint con fechas y Sprint Goal | 3 | — | |
| 3 | HU-06 Editar o eliminar un elemento del backlog | 3 | HU-05 ✓ | |
| 4 | HU-07 Listar y ordenar el backlog por prioridad | 2 | HU-05 ✓ | |
| 5 | HU-09 Asignar historias del backlog a un sprint | 5 | HU-08 | |
| 6 | HU-10 Marcar una historia como completada en el sprint | 2 | HU-09 | |
| 7 | HU-11 Cerrar un sprint | 3 | HU-10 | |
| 8 | HU-20 Registrar un defecto | 3 | HU-08 | |
| | **Total** | **26** | | |

Todas las historias con pantalla dependen además de HU-30.

## Reglas para no bloquearnos

1. **HU-30 y HU-08 se mergean primero**, si se puede el primer día a la mañana. HU-30 trae el diseño base que usan todas las pantallas, y HU-08 crea la entidad Sprint que usan HU-09, HU-10, HU-11 y HU-20.
2. **Quien depende de otra historia arranca por lo que no depende de nadie:** la spec, los escenarios BDD y el dominio. La persistencia, el HTTP y la pantalla se conectan cuando la otra historia está en `main`.
3. **Antes de empezar, `git checkout main` y `git pull`**, y la rama de la historia sale de ahí. Si alguien mergea algo que necesitás, traelo con `git merge origin/main` en tu rama.
4. **PR chicos y seguido:** mejor mergear el dominio y el caso de uso apenas estén en verde que esperar a tener la historia entera.

## Definition of Done (agregado del Sprint 2)

Además de lo que ya pide el CLAUDE.md, **cada historia incluye su pantalla**: un formulario o una vista en la interfaz web que use el mismo caso de uso que el endpoint JSON (ADR-0003).

## Decisiones de la planning

- Las estimaciones son las re-estimadas al cerrar el Sprint 1 (ver `docs/scrum/product-backlog.md`); el equipo las adoptó en lugar de votarlas en Planning Poker.
- Cada integrante elige sus historias. La columna *Responsable* se completa en la planning.

## Ceremonias

- **Planning:** al arrancar el sprint, unos 30 minutos.
- **Daily:** un mensaje en el grupo al terminar el primer día: qué hice, qué hago mañana y si estoy trabado.
- **Review y retrospectiva:** juntas, al cerrar el sprint, unos 30 minutos. Ahí se mide la velocidad real y, si hace falta, se recorta el Sprint 3.
