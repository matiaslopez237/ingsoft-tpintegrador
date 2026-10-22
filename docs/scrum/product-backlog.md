# Product Backlog — Software Metrics & Estimation

Trabajo Práctico Integrador · Ingeniería y Calidad de Software 2026 · UTN FRSR.

10 épicas, 29 historias de usuario, **117 story points** en total. Las estimaciones son una **simulación inicial** (hecha por Claude, no por el equipo) — hay que confirmarlas o ajustarlas en una sesión real de Planning Poker.

Todo esto ya está cargado en GitHub como Issues: las épicas son issues padre, y cada historia es una sub-issue de su épica. Ver `https://github.com/matiaslopez237/ingsoft-tpintegrador/issues`.

**Sprint 1 (MVP) — en curso:** HU-01, HU-02, HU-03, HU-04, HU-18 (14 pts). Detalle de tareas en `docs/scrum/sprint-1/tareas-sprint-1.md`.

---

## Épica 1 — Gestión de proyectos

> Crear, modificar y consultar proyectos; registrar integrantes del equipo.

| ID | Historia | Pts |
|---|---|---|
| HU-01 | Como Product Builder, quiero crear un proyecto con nombre, fecha de inicio y fin, para poder empezar a planificar el trabajo. | 3 |
| HU-02 | Como Product Builder, quiero modificar los datos de un proyecto existente, para poder corregir o actualizar su información. | 2 |
| HU-03 | Como Product Builder, quiero registrar los integrantes de un proyecto, para saber quién participa. | 3 |
| HU-04 | Como Product Builder, quiero consultar el estado general de un proyecto, para ver su progreso. | 3 |

**Subtotal: 11 pts**

---

## Épica 2 — Product Backlog

> Administrar el Product Backlog de cada proyecto: crear, editar, eliminar y priorizar elementos.

| ID | Historia | Pts |
|---|---|---|
| HU-05 | Como Product Builder, quiero crear un elemento de backlog con id, título, descripción, prioridad, estado, story points y criterios de aceptación, para registrar una necesidad del producto. | 5 |
| HU-06 | Como Product Builder, quiero editar o eliminar un elemento del backlog, para mantenerlo actualizado. | 3 |
| HU-07 | Como Product Owner, quiero listar y ordenar el backlog por prioridad, para saber qué se debe entregar primero. | 3 |

**Subtotal: 11 pts**

---

## Épica 3 — Gestión de Sprints

> Crear Sprints, definir Sprint Goal, asignar historias, marcar completadas, cerrar y consultar Sprints anteriores.

| ID | Historia | Pts |
|---|---|---|
| HU-08 | Como Scrum Master, quiero crear un Sprint con fecha de inicio, fin y un Sprint Goal, para organizar el trabajo del equipo. | 3 |
| HU-09 | Como equipo, quiero asignar historias del Product Backlog a un Sprint, para definir el Sprint Backlog. | 5 |
| HU-10 | Como Developer, quiero marcar una historia como completada dentro de un Sprint, para reflejar el avance real. | 3 |
| HU-11 | Como Scrum Master, quiero cerrar un Sprint, para dar por finalizado el período de trabajo. | 3 |
| HU-12 | Como Scrum Master, quiero consultar Sprints anteriores de un proyecto, para revisar el historial. | 2 |

**Subtotal: 16 pts**

---

## Épica 4 — Estimación (Planning Poker)

> Estimar historias con story points e implementar el mecanismo de Planning Poker (votación oculta, rondas, acuerdo final).

| ID | Historia | Pts |
|---|---|---|
| HU-13 | Como Product Builder, quiero asignar story points a una historia, para reflejar su tamaño relativo. | 2 |
| HU-14 | Como Product Builder, quiero registrar mi estimación de forma oculta hasta que todos voten, para evitar sesgos por ver el voto de otros. | 8 |
| HU-15 | Como equipo, quiero ver todas las estimaciones reveladas y si hay diferencias significativas, para saber si hace falta discutir. | 5 |
| HU-16 | Como equipo, quiero iniciar una nueva ronda de votación sobre la misma historia, para converger en una estimación. | 3 |
| HU-17 | Como Product Owner, quiero registrar la estimación final acordada por el equipo, para dejarla asociada a la historia. | 2 |

**Subtotal: 20 pts**

---

## Épica 5 — Registro de esfuerzo

> Registrar esfuerzo por integrante/fecha/actividad/horas y comparar esfuerzo estimado vs. real.

| ID | Historia | Pts |
|---|---|---|
| HU-18 | Como Product Builder, quiero registrar integrante, fecha, actividad y horas trabajadas en una historia o tarea, para dejar constancia del esfuerzo real. | 3 |
| HU-19 | Como Scrum Master, quiero comparar el esfuerzo estimado contra el real de una historia, para evaluar la precisión de las estimaciones. | 5 |

**Subtotal: 8 pts**

---

## Épica 6 — Gestión de defectos

> Registrar y actualizar defectos, vinculados a historias y Sprints de detección/resolución.

| ID | Historia | Pts |
|---|---|---|
| HU-20 | Como Product Builder, quiero registrar un defecto con descripción, severidad, estado, historia relacionada y Sprint de detección, para darle seguimiento. | 3 |
| HU-21 | Como Product Builder, quiero actualizar el estado de un defecto y registrar su Sprint de resolución, para reflejar su progreso. | 2 |

**Subtotal: 5 pts**

---

## Épica 7 — Métricas

> Calcular story points planificados/completados, velocidad, horas, desviación, % completadas y defectos.

| ID | Historia | Pts |
|---|---|---|
| HU-22 | Como Scrum Master, quiero ver los story points planificados, completados y la velocidad del equipo por Sprint, para medir el rendimiento. | 5 |
| HU-23 | Como Scrum Master, quiero ver las horas estimadas, las horas reales y su desviación, para evaluar la planificación. | 5 |
| HU-24 | Como Scrum Master, quiero ver el porcentaje de historias completadas y la cantidad de defectos detectados/resueltos, para evaluar la calidad del Sprint. | 5 |

**Subtotal: 15 pts**

---

## Épica 8 — Dashboard

> Panel con estado del proyecto y representaciones gráficas de las métricas.

| ID | Historia | Pts |
|---|---|---|
| HU-25 | Como cualquier rol del equipo, quiero un Dashboard con el estado del proyecto y gráficos de las métricas principales, para tener una visión rápida del avance. | 8 |

**Subtotal: 8 pts**

---

## Épica 9 — Reportes

> Generar y exportar reportes de proyecto/Sprint (historias, estimaciones, esfuerzo, métricas, defectos).

| ID | Historia | Pts |
|---|---|---|
| HU-26 | Como Product Owner, quiero generar un reporte de un proyecto o Sprint con historias, estimaciones, esfuerzo, métricas y defectos, para compartirlo con los interesados. | 8 |
| HU-27 | Como Product Owner, quiero exportar el reporte a un archivo PDF, para entregarlo como documento formal. | 5 |

**Subtotal: 13 pts**

---

## Épica 10 — Documentación

> Manual de usuario y documentación técnica del proyecto.

| ID | Historia | Pts |
|---|---|---|
| HU-28 | Como usuario del sistema, quiero contar con un manual breve de uso, para poder operar el sistema sin ayuda del equipo de desarrollo. | 3 |
| HU-29 | Como integrante del equipo, quiero documentación técnica de la arquitectura y decisiones del proyecto, para que cualquiera pueda entender el sistema. | 5 |

**Subtotal: 8 pts**

---

## Resumen

| Épica | Pts |
|---|---|
| 1. Gestión de proyectos | 11 |
| 2. Product Backlog | 11 |
| 3. Gestión de Sprints | 16 |
| 4. Estimación (Planning Poker) | 20 |
| 5. Registro de esfuerzo | 8 |
| 6. Gestión de defectos | 5 |
| 7. Métricas | 15 |
| 8. Dashboard | 8 |
| 9. Reportes | 13 |
| 10. Documentación | 8 |
| **Total** | **117** |

## Progresión sugerida por Sprint (según el plan de trabajo del TP)

- **Sprint 1 — MVP:** Épica 1 completa + HU-18 (14 pts).
- **Sprint 2 — Interfaz:** Épica 2, resto de Épica 3 (Sprints), interfaz usable sobre lo del Sprint 1.
- **Sprint 3 — Funcionalidad y Calidad:** Épica 4 (Planning Poker), Épica 6 (Defectos), Épica 7 (Métricas), Épica 8 (Dashboard/visualización).
- **Sprint 4 — Cierre:** Épica 9 (Reportes + exportación a PDF), Épica 10 (Documentación), pulido general.

Esta distribución es una propuesta, no está cerrada — se ajusta en cada Sprint Planning según cómo vaya rindiendo el equipo.
