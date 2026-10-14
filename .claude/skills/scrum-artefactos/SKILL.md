---
name: scrum-artefactos
description: Ayuda con la gestión Scrum del TP - redactar historias de usuario, cargar el Product Backlog en GitHub Projects, armar el Sprint Backlog en la planning, y escribir actas de Sprint Review y Retrospectiva. Usar cuando el usuario menciona backlog, historia de usuario, sprint, planning, daily, review, retro, retrospectiva, acta o tablero.
---

# Artefactos Scrum

Roles del proyecto: **Product Architect** = profesores · **Agile Enabler** = un integrante · **Product Builders** = el resto del equipo.
Plan de sprints de la cátedra: Sprint 0 (preparación) · Sprint 1 (MVP) · Sprint 2 (interfaz) · Sprint 3 (funcionalidad, calidad y visualización) · Sprint 4 (cierre, exportar informe PDF, review final).

## Historias de usuario

- Formato: `Como <rol>, quiero <acción>, para <beneficio>.`
- Cada historia tiene los campos mínimos que pide el enunciado: **identificador (`HU-XX`), título, descripción, prioridad, estado, story points y criterios de aceptación**.
- Chequeá INVEST: independiente, negociable, valiosa, estimable, pequeña (entra en un sprint), testeable. Si no entra en un sprint, proponé cómo partirla.
- Los criterios de aceptación de la historia son el resumen; el detalle va en la spec SDD (`sdd-spec`).

## GitHub Projects (con `gh`)

**Siempre mostrale al usuario lo que vas a crear o modificar en GitHub y esperá confirmación**: el tablero es compartido con todo el equipo.

- Ver estado: `gh issue list --limit 200`, `gh project list --owner <owner>`, `gh project item-list <N> --owner <owner>`.
- Crear historia: `gh issue create --title "HU-XX: <título>" --body "<historia + CAs + SP>" --label "historia"`, y después `gh project item-add <N> --owner <owner> --url <url-del-issue>`.
- Defectos (requisito 6 del TP): issue con título `DEF-XX: <descripción>`, label `defecto` y en el cuerpo: severidad, estado, historia relacionada (`HU-XX`), sprint de detección y de resolución.
- Si `gh` pide permisos de proyecto: el usuario tiene que correr `gh auth refresh -s project` (no lo hagas por él).

## Sprint Planning → `docs/scrum/sprint-N/planning.md`

```markdown
# Sprint N: Planning
**Fechas:** AAAA-MM-DD a AAAA-MM-DD  ·  **Asistentes:** ...
**Sprint Goal:** <una frase>
**Capacidad:** <horas o SP disponibles>  ·  **Velocidad de referencia:** <SP del sprint anterior>

| Historia | SP | Responsable(s) | Tareas |
|---|---|---|---|
| HU-XX: ... | 5 | ... | spec · BDD · TDD núcleo · UI · review |

**Total comprometido:** X SP
**Riesgos / dependencias:** ...
```

## Daily → opcional, `docs/scrum/sprint-N/dailies.md`
Una línea por integrante y día: qué hice · qué voy a hacer · bloqueos.

## Sprint Review → `docs/scrum/sprint-N/review.md`

```markdown
# Sprint N: Review
**Fecha:** ...  ·  **Asistentes:** ...  ·  **Sprint Goal:** ... (¿cumplido? sí/parcial/no)

| Historia | SP | Estado | Demo / observaciones del cliente |
|---|---|---|---|

**SP planificados:** X  ·  **SP completados:** Y  ·  **Velocidad:** Y
**Defectos detectados en el sprint:** ...  ·  **Resueltos:** ...
**Feedback del Product Architect:** ...
**Cambios al Product Backlog:** ...
```

## Retrospectiva → `docs/scrum/sprint-N/retro.md` (es un entregable: "actas de retrospectivas")

```markdown
# Sprint N: Acta de Retrospectiva
**Fecha:** ...  ·  **Facilitador (Agile Enabler):** ...  ·  **Asistentes:** ...
**Dinámica usada:** (ej. Start/Stop/Continue, 4L, Velero)

## Qué salió bien
## Qué salió mal / a mejorar
## Acciones de mejora
| Acción | Responsable | Para cuándo | Estado |
|---|---|---|---|

## Seguimiento de acciones de la retro anterior
| Acción | ¿Se cumplió? | Comentario |
|---|---|---|
```

Para completar las secciones, podés apoyarte en datos reales: `git log --since=<inicio> --until=<fin> --format="%an %s"` (actividad por integrante y ciclos TDD), issues cerrados en el período y las métricas del sprint. No inventes opiniones del equipo: dejá las secciones de opinión con viñetas vacías para que las completen en la reunión, salvo que el usuario te las dicte.
