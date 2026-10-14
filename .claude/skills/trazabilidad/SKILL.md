---
name: trazabilidad
description: Genera o verifica la matriz de trazabilidad Historia → Spec SDD → Criterios de aceptación → Escenarios BDD → Tests → Código Go → commits TDD. Usar cuando el usuario pide trazabilidad, la matriz, preparar la demo/presentación final, o verificar que una historia está completa.
---

# Trazabilidad

La cátedra pide demostrar en la presentación final el recorrido completo de al menos una funcionalidad:

**Historia de Usuario → Especificación SDD → Criterios de Aceptación → Escenarios BDD → Tests → Código Go**

La trazabilidad se apoya en **IDs que aparecen en todos los artefactos**. Si un artefacto no lleva su ID, no se puede trazar.

## Convención de IDs (obligatoria en todo el repo)

| Artefacto | Dónde | Cómo lleva el ID |
|---|---|---|
| Historia | GitHub Projects (issue) | Título `HU-XX: ...` |
| Spec SDD | `docs/specs/HU-XX-*.md` | Nombre del archivo y título |
| Reglas de negocio | dentro de la spec | `RN-XX.N` |
| Criterios de aceptación | dentro de la spec | `CA-XX.N` |
| Escenarios BDD | `features/HU-XX-*.feature` | Tags `@HU-XX` (feature) y `@CA-XX.N` (scenario) |
| Tests unitarios Go | `*_test.go` | Nombre del subtest: `t.Run("CA-XX.N: descripción", ...)` o comentario `// Traza: HU-XX CA-XX.N` arriba del test |
| Código Go | funciones/tipos principales | Comentario de doc: `// Implementa HU-XX (RN-XX.1, RN-XX.2).` |
| Commits TDD | git | `test(HU-XX): [RED] ...` · `feat(HU-XX): [GREEN] ...` · `refactor(HU-XX): [REFACTOR] ...` |

## Generar la matriz

1. Listá las historias: specs en `docs/specs/` (y, si el usuario lo pide, issues con `gh issue list --search "HU-" --limit 200`).
2. Para cada `HU-XX`, buscá con Grep:
   - CAs y RNs en la spec.
   - Escenarios con `@CA-XX.` en `features/`.
   - Tests con `CA-XX.` o `HU-XX` en `**/*_test.go`.
   - Código con `Implementa HU-XX`.
   - Commits: `git log --oneline --grep="(HU-XX)"`, separando RED / GREEN / REFACTOR.
3. Corré los tests de la historia si es posible (`go test ./... -run <patrón>` y `go test ./test/bdd/ -v`) para marcar si pasan.
4. Escribí `docs/trazabilidad.md` con:
   - Una tabla resumen por historia: spec ✓/✗ · CAs cubiertos por BDD (x/y) · CAs cubiertos por tests (x/y) · código ✓/✗ · ciclo TDD en git (RED/GREEN/REFACTOR).
   - Una sección por historia con el detalle CA por CA y links relativos a los archivos (`features/...feature`, `internal/...go`) y hashes de commits.
5. Al final, una lista de **huecos**: CAs sin escenario, CAs sin test, código sin ID, historias sin commits RED.

No inventes coberturas: si algo no se encontró, va con ✗.

## Verificar una historia antes de cerrarla

Una historia está trazable si: tiene spec aprobada, cada CA tiene escenario BDD (o justificación `@manual`) y test, el código lleva `Implementa HU-XX`, y existe al menos un ciclo de commits RED → GREEN (REFACTOR si hubo). Reportá qué falta.

## Preparar la demo final

Elegí la historia con la matriz más completa y armá un guion corto en `docs/demo-trazabilidad.md`: el issue → la spec (sección 8) → el `.feature` con sus tags → el test → la función Go → los commits RED/GREEN/REFACTOR en orden → `go test` corriendo en verde.
