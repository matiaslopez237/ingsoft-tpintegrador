# Software Metrics & Estimation

Trabajo Práctico Integrador · Ingeniería y Calidad de Software 2026 · UTN FRSR.

Aplicación web en **Go** para estimar, planificar, seguir y medir proyectos de software: proyectos, product backlog, sprints, planning poker, registro de esfuerzo, defectos, métricas, dashboard y reportes exportables a PDF.

Metodología: Scrum · SDD · BDD (godog) · TDD.

## Trabajar con Claude Code

El repo trae un kit de skills en `.claude/skills/` y reglas del proyecto en `CLAUDE.md`. Claude Code los carga solo al abrir una sesión en esta carpeta, así que todo el equipo trabaja con las mismas reglas.

Flujo por historia (detalle en `CLAUDE.md`):

| Paso | Comando |
|---|---|
| Refinar la historia | `/grill-me` |
| Cargarla al tablero / planning / retro | `/scrum-artefactos` |
| Especificar (SDD) | `/sdd-spec HU-XX` |
| Escenarios BDD | `/bdd-go HU-XX` |
| Implementar con TDD | `/tdd-workflow HU-XX` |
| Revisar | "revisá el cambio con el agente go-reviewer" |
| Verificar trazabilidad | `/trazabilidad HU-XX` |

Convención de commits: `test(HU-XX): [RED] ...` → `feat(HU-XX): [GREEN] ...` → `refactor(HU-XX): [REFACTOR] ...`

## Comandos

```bash
go test ./...                              # tests unitarios + BDD
go test ./... -coverprofile=coverage.out   # cobertura
go tool cover -func=coverage.out
```

## Equipo

| Integrante | Rol |
|---|---|
| | Agile Enabler |
| | Product Builder |
| | Product Builder |
| | Product Builder |
| | Product Builder |

Product Architect: profesores de la cátedra.
