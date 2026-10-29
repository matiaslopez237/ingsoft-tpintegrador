# Software Metrics & Estimation

[![CI/CD](https://github.com/matiaslopez237/ingsoft-tpintegrador/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/matiaslopez237/ingsoft-tpintegrador/actions/workflows/ci.yml)

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

## Integración y despliegue continuos

Cada PR a `main` corre el CI (build, `go vet`, `gofmt` y los tests unitarios y BDD con `-race` y cobertura). Cuando el PR se mergea y el CI pasa en `main`, el pipeline despliega en Render. Detalle en [ADR-0004](docs/adr/0004-ci-cd.md).

## Equipo

| Integrante | Rol |
|---|---|
| Matías López | Agile Enabler |
| Ezequiel Assenza | Product Builder |
| Nicolás Andino | Product Builder |
| Lucas Orellana | Product Builder |

Product Architect: profesores de la cátedra.
