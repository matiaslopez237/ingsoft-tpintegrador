---
name: sdd-spec
description: Escribe o revisa la especificación SDD de una historia de usuario con el formato exigido por la cátedra (objetivo, entradas, salidas, reglas de negocio, restricciones, casos límite, errores, criterios de aceptación). Usar antes de implementar cualquier funcionalidad, cuando el usuario dice "spec", "especificación", "SDD", o pide implementar una historia que todavía no tiene spec en docs/specs/.
---

# Especificación SDD

En este proyecto **ninguna funcionalidad se implementa sin spec**. La spec es el contrato del que salen los escenarios BDD (`bdd-go`), los tests (`tdd-workflow`) y la matriz de trazabilidad (`trazabilidad`).

## Dónde va

- Una spec por historia: `docs/specs/HU-XX-nombre-corto.md` (ej. `docs/specs/HU-04-planning-poker.md`).
- El ID `HU-XX` es el mismo que el de la historia en GitHub Projects (título del issue: `HU-04: Planning Poker`).
- Se versiona en git junto al código. Si la spec cambia después de implementada, actualizala en el mismo commit/PR que cambia el código y anotá el cambio en "Historial".

## Proceso

1. **Leé la historia** (issue o lo que diga el usuario) y las specs relacionadas en `docs/specs/`.
2. **Si faltan definiciones de negocio, preguntá** antes de escribir (usá `grill-me` si son muchas). No inventes reglas de negocio: marcá lo dudoso en "Preguntas abiertas".
3. **Escribí la spec con la plantilla de abajo, sin saltear secciones.** Si una sección no aplica, escribí "No aplica" y por qué.
4. **Numerá los criterios de aceptación** como `CA-XX.N` (ej. `CA-04.1`). Cada CA debe ser verificable con un escenario BDD o un test; si no se puede verificar, reescribilo.
5. **Cubrí casos límite y errores de verdad**: valores vacíos, cero, negativos, máximos, duplicados, estados inválidos (ej. cerrar un sprint ya cerrado), datos de otro proyecto.
6. Mostrale la spec al usuario y esperá su OK antes de pasar a BDD/TDD.

## Plantilla

```markdown
# HU-XX: <Título de la historia>

**Historia:** Como <rol>, quiero <acción>, para <beneficio>.
**Issue:** #<número en GitHub>  ·  **Sprint:** <N>  ·  **Story points:** <SP>
**Estado de la spec:** Borrador | Aprobada | Implementada

## 1. Objetivo
Qué problema resuelve y para quién, en 2-4 líneas.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|

## 3. Salidas esperadas
Qué devuelve / qué estado queda / qué ve el usuario.

## 4. Reglas de negocio
- **RN-XX.1:** ...
- **RN-XX.2:** ...

## 5. Restricciones
Técnicas o de dominio (ej. "el núcleo no depende del framework web", "solo el Agile Enabler puede cerrar un sprint").

## 6. Casos límite
- ...

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| ... | `ErrSprintCerrado` | ... |

## 8. Criterios de aceptación
- **CA-XX.1:** Dado ..., cuando ..., entonces ...
- **CA-XX.2:** ...

## Preguntas abiertas
- ...

## Historial
- AAAA-MM-DD: creada.
```

## Revisión de una spec existente

Chequeá: las 8 secciones presentes; CAs numerados y verificables; cada regla de negocio cubierta por al menos un CA; errores con nombre de error Go concreto; casos límite reales (no genéricos). Reportá lo que falta en una lista corta.
