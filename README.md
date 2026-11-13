# Software Metrics & Estimation

[![CI/CD](https://github.com/matiaslopez237/ingsoft-tpintegrador/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/matiaslopez237/ingsoft-tpintegrador/actions/workflows/ci.yml)

Trabajo Práctico Integrador · Ingeniería y Calidad de Software 2026 · UTN FRSR.

Aplicación web en **Go** para estimar, planificar, seguir y medir proyectos de software. Ya permite gestionar proyectos, integrantes, product backlog y sprints, estimar con Planning Poker, registrar esfuerzo y defectos, y ver métricas de velocidad, horas y desviación. El dashboard con gráficos y los reportes exportables a PDF están en desarrollo (Sprint 4).

Metodología: Scrum · SDD · BDD (godog) · TDD.

**Versión publicada:** <https://ingsoft-tpintegrador.onrender.com>. Está en el plan gratuito de Render: si nadie la usó en un rato, la primera visita tarda cerca de un minuto, y los datos se borran en cada despliegue.

## Puesta en marcha

**Requisitos:** Go, en la versión que indica `go.mod`, y Git. No hace falta instalar una base de datos ni un compilador de C: SQLite viene incluido en el programa ([ADR-0002](docs/adr/0002-persistencia.md)).

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/matiaslopez237/ingsoft-tpintegrador.git
   cd ingsoft-tpintegrador
   ```

2. Levantar el servidor. La primera vez Go descarga las dependencias:

   ```bash
   go run ./cmd/server
   ```

3. Abrir <http://localhost:8080> en el navegador, que lleva a la lista de proyectos.

La base se crea sola, en el archivo `metrics.db` de la carpeta desde donde se corre el servidor. Para empezar de cero, se detiene el servidor y se borra ese archivo.

### Configuración

| Variable | Para qué sirve | Por defecto |
|---|---|---|
| `PORT` | Puerto del servidor | `8080` |
| `DB_PATH` | Archivo de la base SQLite | `metrics.db` |
| `CARGAR_EJEMPLO` | Con `true` y la base vacía, arranca con el proyecto de ejemplo (este mismo TP). Si la carga falla, el servidor no arranca y pide borrar la base | Apagada (Render la tiene en `true`) |

Por ejemplo, en Linux o macOS: `PORT=9090 DB_PATH=prueba.db go run ./cmd/server`. En PowerShell: `$env:PORT="9090"; go run ./cmd/server`.

### Rutas útiles

- `/web/proyectos`: la interfaz web.
- `/salud`: responde si el servidor está en pie.
- `/proyectos/...`: la API JSON, que usa los mismos casos de uso que las pantallas. Las rutas de cada historia están en su spec, en [`docs/specs/`](docs/specs/).

## Tests y calidad

```bash
go test ./...                                 # tests unitarios y escenarios BDD
BDD_TAGS=@HU-08 go test ./test/bdd/           # solo los escenarios de una historia
go test ./... -coverprofile=coverage.out      # cobertura
go tool cover -func=coverage.out              # resumen de cobertura
go vet ./...
gofmt -l .                                    # tiene que no listar nada
```

- **Escenarios BDD:** están en [`features/`](features/), y sus pasos en [`test/bdd/`](test/bdd/).
- **Detector de carreras (`-race`):** el CI lo corre en Linux. En Windows necesita un compilador de C, así que localmente se puede omitir.
- **Suite completa:** puede tardar más de los 10 minutos que `go test` deja por paquete, sobre todo en máquinas lentas. Si corta por tiempo, agregá `-timeout 20m`, como el CI: `go test -timeout 20m ./...`.
- **`BDD_TAGS` en Windows:** el ejemplo usa la sintaxis de Linux y macOS. En PowerShell: `$env:BDD_TAGS="@HU-08"; go test ./test/bdd/`.

## Documentación

| Documento | Qué tiene |
|---|---|
| [`docs/manual-usuario.md`](docs/manual-usuario.md) | Cómo se usa cada pantalla, con capturas del proyecto de ejemplo. |
| [`docs/arquitectura.md`](docs/arquitectura.md) | Capas, puertos y adaptadores, cómo viaja un pedido y decisiones técnicas. |
| [`docs/adr/`](docs/adr/) | Decisiones de arquitectura, una por archivo. |
| [`docs/specs/`](docs/specs/) | Spec SDD de cada historia: reglas de negocio, errores y criterios de aceptación. |
| [`docs/trazabilidad.md`](docs/trazabilidad.md) | Historia → spec → criterios → BDD → tests → código → commits. |
| [`docs/informe-metricas.md`](docs/informe-metricas.md) | Cobertura, tests, velocidad, estimación y defectos al cierre del Sprint 4. |
| [`docs/scrum/`](docs/scrum/) | Product Backlog y actas de cada sprint. |

## Integración y despliegue continuos

Cada PR a `main` corre el CI: build, `go vet`, `gofmt` y los tests unitarios y BDD, con `-race` y una cobertura mínima del 80 %. `main` está protegida: hace falta un PR, una aprobación y el CI en verde.

Al mergear a `main`, si el CI pasa, el pipeline despliega en Render (https://ingsoft-tpintegrador.onrender.com) con el *deploy hook* guardado en el secret `RENDER_DEPLOY_HOOK_URL`. Detalle en [ADR-0004](docs/adr/0004-ci-cd.md).

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

## Equipo

| Integrante | Rol |
|---|---|
| Matías López | Agile Enabler |
| Ezequiel Assenza | Product Builder |
| Nicolás Andino | Product Builder |
| Lucas Orellana | Product Builder |

Product Architect: profesores de la cátedra.
