# Product Backlog — Software Metrics & Estimation

Trabajo Práctico Integrador · Ingeniería y Calidad de Software 2026 · UTN FRSR.

10 épicas, 30 historias de usuario, **98 story points** en total.

Todo esto está cargado en GitHub como Issues: las épicas son issues padre, y cada historia es una sub-issue de su épica. Ver `https://github.com/matiaslopez237/ingsoft-tpintegrador/issues`.

**Estimaciones.** La estimación inicial (columna *Inicial*) fue una simulación hecha por Claude al armar el backlog. Al cerrar el Sprint 1 (2026-10-25) se re-estimaron las historias pendientes con lo que ya se sabía del proyecto (la base técnica está hecha y desde el Sprint 2 cada historia incluye su pantalla), y **el equipo adoptó esas estimaciones** en lugar de votarlas en Planning Poker (columna *Pts*). Las historias del Sprint 1 conservan su estimación original. La estimación inicial sumaba 117 en el resumen, pero las épicas sumaban 115: era un error de suma.

---

## Calendario

| Sprint | Fechas | Objetivo | Pts |
|---|---|---|---|
| 1 — MVP | hasta el 25 oct | Gestión de proyectos (Épica 1) + crear elementos de backlog + registrar esfuerzo | 19 |
| 2 — Interfaz, backlog y sprints | 26 y 27 oct | Gestionar el backlog y armar sprints de un proyecto desde una interfaz web | 26 |
| 3 — Estimación y métricas | 28 y 29 oct | Planning Poker, esfuerzo estimado vs. real y primeras métricas | 28 |
| 4 — Visualización y cierre | 30 y 31 oct | Dashboard, reportes en PDF y documentación | 25 |
| — | 1 nov | Revisión general y correcciones (sin historias nuevas) | — |
| — | 2 nov | **Entrega** | — |

Sprints de 2 días. Al cerrar el Sprint 2 se mide la velocidad real: si quedó lejos de lo planificado, se posterga lo último de cada lista (el orden de cada sprint es su prioridad).

---

## Épica 1 — Gestión de proyectos

> Crear, modificar y consultar proyectos; registrar integrantes del equipo.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-01 | Como Product Builder, quiero crear un proyecto con nombre, fecha de inicio y fin, para poder empezar a planificar el trabajo. | 3 | 3 | 1 ✓ |
| HU-02 | Como Product Builder, quiero modificar los datos de un proyecto existente, para poder corregir o actualizar su información. | 2 | 2 | 1 ✓ |
| HU-03 | Como Product Builder, quiero registrar los integrantes de un proyecto, para saber quién participa. | 3 | 3 | 1 ✓ |
| HU-04 | Como Product Builder, quiero consultar el estado general de un proyecto, para ver su progreso. | 3 | 3 | 1 ✓ |

**Subtotal: 11 pts**

---

## Épica 2 — Product Backlog

> Administrar el Product Backlog de cada proyecto: crear, editar, eliminar y priorizar elementos.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-05 | Como Product Builder, quiero crear un elemento de backlog con id, título, descripción, prioridad, estado, story points y criterios de aceptación, para registrar una necesidad del producto. | 5 | 5 | 1 ✓ |
| HU-06 | Como Product Builder, quiero editar o eliminar un elemento del backlog, para mantenerlo actualizado. | 3 | 3 | 2 |
| HU-07 | Como Product Owner, quiero listar y ordenar el backlog por prioridad, para saber qué se debe entregar primero. | 3 | 2 | 2 |

**Subtotal: 10 pts**

---

## Épica 3 — Gestión de Sprints

> Crear Sprints, definir Sprint Goal, asignar historias, marcar completadas, cerrar y consultar Sprints anteriores.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-08 | Como Scrum Master, quiero crear un Sprint con fecha de inicio, fin y un Sprint Goal, para organizar el trabajo del equipo. | 3 | 3 | 2 |
| HU-09 | Como equipo, quiero asignar historias del Product Backlog a un Sprint, para definir el Sprint Backlog. | 5 | 5 | 2 |
| HU-10 | Como Developer, quiero marcar una historia como completada dentro de un Sprint, para reflejar el avance real. | 3 | 2 | 2 |
| HU-11 | Como Scrum Master, quiero cerrar un Sprint, para dar por finalizado el período de trabajo. | 3 | 3 | 2 |
| HU-12 | Como Scrum Master, quiero consultar Sprints anteriores de un proyecto, para revisar el historial. | 2 | 2 | 3 |

**Subtotal: 15 pts**

---

## Épica 4 — Estimación (Planning Poker)

> Estimar historias con story points e implementar el mecanismo de Planning Poker (votación oculta, rondas, acuerdo final).

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-13 | Como Product Builder, quiero asignar story points a una historia, para reflejar su tamaño relativo. | 2 | 1 | 3 |
| HU-14 | Como Product Builder, quiero registrar mi estimación de forma oculta hasta que todos voten, para evitar sesgos por ver el voto de otros. | 8 | 5 | 3 |
| HU-15 | Como equipo, quiero ver todas las estimaciones reveladas y si hay diferencias significativas, para saber si hace falta discutir. | 5 | 3 | 3 |
| HU-16 | Como equipo, quiero iniciar una nueva ronda de votación sobre la misma historia, para converger en una estimación. | 3 | 2 | 3 |
| HU-17 | Como Product Owner, quiero registrar la estimación final acordada por el equipo, para dejarla asociada a la historia. | 2 | 2 | 3 |

**Subtotal: 13 pts**

---

## Épica 5 — Registro de esfuerzo

> Registrar esfuerzo por integrante/fecha/actividad/horas y comparar esfuerzo estimado vs. real.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-18 | Como Product Builder, quiero registrar integrante, fecha, actividad y horas trabajadas en una historia o tarea, para dejar constancia del esfuerzo real. | 3 | 3 | 1 |
| HU-19 | Como Scrum Master, quiero comparar el esfuerzo estimado contra el real de una historia, para evaluar la precisión de las estimaciones. | 5 | 3 | 3 |

**Subtotal: 6 pts**

---

## Épica 6 — Gestión de defectos

> Registrar y actualizar defectos, vinculados a historias y Sprints de detección/resolución.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-20 | Como Product Builder, quiero registrar un defecto con descripción, severidad, estado, historia relacionada y Sprint de detección, para darle seguimiento. | 3 | 3 | 2 |
| HU-21 | Como Product Builder, quiero actualizar el estado de un defecto y registrar su Sprint de resolución, para reflejar su progreso. | 2 | 2 | 3 |

**Subtotal: 5 pts**

---

## Épica 7 — Métricas

> Calcular story points planificados/completados, velocidad, horas, desviación, % completadas y defectos.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-22 | Como Scrum Master, quiero ver los story points planificados, completados y la velocidad del equipo por Sprint, para medir el rendimiento. | 5 | 5 | 3 |
| HU-23 | Como Scrum Master, quiero ver las horas estimadas, las horas reales y su desviación, para evaluar la planificación. | 5 | 3 | 3 |
| HU-24 | Como Scrum Master, quiero ver el porcentaje de historias completadas y la cantidad de defectos detectados/resueltos, para evaluar la calidad del Sprint. | 5 | 3 | 4 |

**Subtotal: 11 pts**

---

## Épica 8 — Dashboard

> Panel con estado del proyecto y representaciones gráficas de las métricas.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-25 | Como cualquier rol del equipo, quiero un Dashboard con el estado del proyecto y gráficos de las métricas principales, para tener una visión rápida del avance. | 8 | 8 | 4 |

**Subtotal: 8 pts**

---

## Épica 9 — Reportes

> Generar y exportar reportes de proyecto/Sprint (historias, estimaciones, esfuerzo, métricas, defectos).

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-26 | Como Product Owner, quiero generar un reporte de un proyecto o Sprint con historias, estimaciones, esfuerzo, métricas y defectos, para compartirlo con los interesados. | 8 | 5 | 4 |
| HU-27 | Como Product Owner, quiero exportar el reporte a un archivo PDF, para entregarlo como documento formal. | 5 | 5 | 4 |

**Subtotal: 10 pts**

---

## Épica 10 — Documentación

> Manual de usuario y documentación técnica del proyecto.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-28 | Como usuario del sistema, quiero contar con un manual breve de uso, para poder operar el sistema sin ayuda del equipo de desarrollo. | 3 | 2 | 4 |
| HU-29 | Como integrante del equipo, quiero documentación técnica de la arquitectura y decisiones del proyecto, para que cualquiera pueda entender el sistema. | 5 | 2 | 4 |

**Subtotal: 4 pts**

---

## Historias transversales

> Trabajo que no pertenece a una sola épica.

| ID | Historia | Inicial | Pts | Sprint |
|---|---|---|---|---|
| HU-30 | Como usuario del sistema, quiero una interfaz web con una base común (diseño, navegación y pantallas de proyectos e integrantes), para operar el sistema desde el navegador. Ver ADR-0003. | — | 5 | 2 |

Se sumó al backlog al planificar el Sprint 2 (2026-10-25): al cerrar el Sprint 1 la aplicación solo respondía JSON. Desde el Sprint 2, cada historia incluye su propia pantalla.

**Subtotal: 5 pts**

---

## Resumen

| Épica | Inicial | Pts |
|---|---|---|
| 1. Gestión de proyectos | 11 | 11 |
| 2. Product Backlog | 11 | 10 |
| 3. Gestión de Sprints | 16 | 15 |
| 4. Estimación (Planning Poker) | 20 | 13 |
| 5. Registro de esfuerzo | 8 | 6 |
| 6. Gestión de defectos | 5 | 5 |
| 7. Métricas | 15 | 11 |
| 8. Dashboard | 8 | 8 |
| 9. Reportes | 13 | 10 |
| 10. Documentación | 8 | 4 |
| Transversales (HU-30) | — | 5 |
| **Total** | **115** | **98** |

## Orden dentro de cada sprint

El orden de las historias de cada sprint es su prioridad. Se decidió por dependencias (qué historia necesita que exista otra antes) y por la progresión del plan de trabajo del TP:

- **Sprint 2:** HU-30, HU-08, HU-06, HU-07, HU-09, HU-10, HU-11, HU-20.
- **Sprint 3:** HU-22, HU-14, HU-15, HU-16, HU-17, HU-19, HU-23, HU-21, HU-12, HU-13.
- **Sprint 4:** HU-24, HU-25, HU-26, HU-27, HU-29, HU-28.

Dependencias principales: las pantallas dependen de HU-30; HU-09 a HU-11 dependen de HU-08; las métricas (Épica 7) dependen de sprints cerrados (HU-11), esfuerzo (HU-18, HU-19) y defectos (HU-20, HU-21); el Dashboard y los reportes dependen de las métricas.
