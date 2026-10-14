# TP Integrador: Software Metrics & Estimation

Ingeniería y Calidad de Software 2026 · UTN FRSR. Aplicación **web** para estimar, planificar, seguir y medir proyectos de software (proyectos, backlog, sprints, planning poker, esfuerzo, defectos, métricas, dashboard y reportes con exportación a PDF).

## Reglas obligatorias de la cátedra

- **Núcleo y reglas de negocio en Go.** Nada de lógica de negocio en la interfaz.
- **Scrum** para la gestión (tablero en GitHub Projects).
- **SDD:** toda funcionalidad principal se especifica **antes** de implementarse.
- **BDD:** escenarios Given-When-Then (normales, alternativos, límite, errores), automatizados cuando se pueda.
- **TDD:** métricas, estimación, reglas de negocio y validaciones se desarrollan con RED → GREEN → REFACTOR, **con evidencia en el historial de git**.
- **Trazabilidad:** Historia → Spec → Criterios → BDD → Tests → Código Go.

## Reglas para Claude

1. **Explicá cada cambio.** El enunciado exige que el equipo comprenda y valide todo lo generado con IA, y en la defensa hay que justificar decisiones. Al terminar cada tarea: qué cambiaste, por qué, y cómo verificarlo. Si usás algo no obvio de Go, explicalo en una línea.
2. **No inventes reglas de negocio.** Si la spec no lo define, preguntá.
3. **No implementes sin spec.** Si la historia no tiene `docs/specs/HU-XX-*.md`, usá primero la skill `sdd-spec`.
4. **Test primero en el núcleo.** Nada de código de dominio sin un test que falle antes.
5. **No hagas `git push`, ni crees o modifiques issues o el tablero, sin confirmación explícita.** El repo y el tablero son compartidos por 5 personas.
6. Escribí en español (código: identificadores en español o inglés, pero consistentes con lo que ya exista).

## Flujo por historia de usuario

1. **Refinar** la historia (si hay dudas: `grill-me`). Historia en GitHub Projects como `HU-XX: título` (`scrum-artefactos`).
2. **Spec SDD** → `docs/specs/HU-XX-*.md` (`sdd-spec`). Esperar OK del equipo.
3. **Escenarios BDD** → `features/HU-XX-*.feature` + steps godog en rojo (`bdd-go`).
4. **TDD** del núcleo, ciclo por ciclo, con commits por etapa (`tdd-workflow` + `golang-testing`).
5. **Adaptadores e interfaz web** (HTTP, persistencia, vistas) sobre el núcleo ya testeado (`hexagonal-architecture`, `impeccable`, `frontend-a11y`).
6. **Revisión**: agente `go-reviewer` sobre el cambio.
7. **Trazabilidad** verificada antes de cerrar la historia (`trazabilidad`).

## Convenciones

**IDs:** `HU-XX` (historia) · `RN-XX.N` (regla de negocio) · `CA-XX.N` (criterio de aceptación) · `DEF-XX` (defecto). Detalle en la skill `trazabilidad`.

**Commits (Conventional Commits + historia + etapa TDD).** Este formato reemplaza a los formatos de ejemplo de `tdd-workflow` y `git-workflow`:
```
test(HU-04): [RED] votos ocultos hasta cerrar la ronda
feat(HU-04): [GREEN] ocultar votos hasta cerrar la ronda
refactor(HU-04): [REFACTOR] extraer cálculo de dispersión
docs(HU-04): spec de planning poker
fix(DEF-03): velocidad con sprint sin historias
```
No hacer squash de los commits RED/GREEN/REFACTOR: son la evidencia de TDD.

**Ramas:** `main` protegida; una rama por historia: `hu-04-planning-poker`. PR con revisión de al menos otro integrante.

## Estructura sugerida

```
cmd/server/            # main: arranca el servidor web
internal/domain/       # entidades y reglas de negocio (sin dependencias externas)
internal/app/          # casos de uso (orquestan el dominio, definen puertos)
internal/adapters/     # http (handlers/vistas), persistencia, exportación PDF
features/              # escenarios BDD (.feature)
test/bdd/              # runner y steps de godog
docs/specs/            # specs SDD por historia
docs/adr/              # decisiones de arquitectura (ADR)
docs/scrum/sprint-N/   # planning, review, retro (actas)
docs/trazabilidad.md   # matriz generada
```

Decisiones pendientes para Sprint 0, cada una se registra como ADR en `docs/adr/` (`architecture-decision-records`): framework/enfoque de la interfaz web, persistencia (ej. SQLite, PostgreSQL, archivos), librería de gráficos del dashboard y librería para exportar el PDF.

## Comandos

```bash
go test ./...                                   # todos los tests (unitarios + BDD)
go test ./test/bdd/ -v                          # solo escenarios BDD
go test ./... -coverprofile=coverage.out        # cobertura
go tool cover -func=coverage.out                # resumen de cobertura (para el informe)
go tool cover -html=coverage.out -o coverage.html
go vet ./...
gofmt -l .
```

## Definition of Done (por historia)

- Spec aprobada en `docs/specs/`.
- Todos los CA con escenario BDD (o `@manual` justificado) y test en Go, en verde.
- Commits RED → GREEN (→ REFACTOR) en el historial.
- `go vet` y `gofmt` limpios; revisado por `go-reviewer` y por otro integrante (PR).
- Matriz de trazabilidad actualizada.
- Issue cerrado y movido a "Done" en el tablero.

## Skills disponibles en `.claude/skills/`

| Para | Skill |
|---|---|
| Refinar historias / decisiones | `grill-me` |
| Scrum: historias, tablero, planning, review, retro | `scrum-artefactos` |
| Specs SDD | `sdd-spec` |
| BDD con godog | `bdd-go` |
| TDD con evidencia en git | `tdd-workflow` |
| Tests en Go, cobertura | `golang-testing` |
| Go idiomático | `golang-patterns` |
| Errores | `error-handling` |
| Arquitectura del núcleo | `hexagonal-architecture` |
| Decisiones (ADR) | `architecture-decision-records` |
| Git | `git-workflow` (usar la convención de commits de arriba) |
| Bugs | `diagnosing-bugs` |
| Interfaz web | `impeccable`, `frontend-a11y` |
| Trazabilidad y demo final | `trazabilidad` |

Agentes en `.claude/agents/`: `go-reviewer` (revisión de código Go) y `tdd-guide` (acompaña el ciclo TDD).
