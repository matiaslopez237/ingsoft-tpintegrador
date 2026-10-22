# ADR-0002: Mecanismo de persistencia

**Date**: 2026-10-22
**Status**: accepted
**Deciders**: Equipo completo (Matías López, Ezequiel Assenza, Nicolás Andino, Lucas Orellana)

## Context

El sistema necesita persistir proyectos, integrantes, backlog, sprints, estimaciones, esfuerzo, defectos y los datos para calcular métricas. Varias de esas consultas requieren relacionar entidades entre sí (ej. story points completados por sprint, desviación de horas por historia), lo que pide algo más que archivos planos. El equipo no cuenta con infraestructura propia para levantar un servidor de base de datos, y cada integrante va a correr el proyecto en su propia máquina (Windows/Mac/Linux mezclados).

## Decision

Usar SQLite como mecanismo de persistencia, a través del driver `modernc.org/sqlite` (implementación de SQLite en Go puro, sin CGO).

## Alternatives Considered

### Alternative 1: PostgreSQL / MySQL (servidor de base de datos aparte)
- **Pros**: más robusto para producción real, mejor concurrencia.
- **Cons**: requiere instalar y levantar un servidor de base de datos en la máquina de cada integrante, o depender de un servicio externo.
- **Why not**: complejidad innecesaria para el alcance del TP; agrega un paso de configuración que puede fallar distinto en cada máquina.

### Alternative 2: Persistencia en archivos JSON (encoding/json)
- **Pros**: cero dependencias, muy simple de entender.
- **Cons**: las consultas de métricas (Épica 7) y reportes (Épica 9) necesitan relacionar varias entidades — con JSON esas consultas hay que programarlas a mano en vez de usar SQL.
- **Why not**: el costo de reimplementar joins/agregaciones a mano supera la simplicidad que aporta evitar SQL.

### Alternative 3: SQLite con `mattn/go-sqlite3`
- **Pros**: driver más antiguo y muy usado, buena documentación.
- **Cons**: usa CGO, por lo que cada integrante necesita un compilador de C instalado (gcc/mingw en Windows) para que `go build` funcione.
- **Why not**: fuente típica de "en mi máquina no compila" entre compañeros con sistemas distintos; se evita por completo con una alternativa en Go puro.

## Consequences

### Positive
- Persistencia como un único archivo (`proyecto.db`), sin servidor aparte que instalar ni levantar.
- `go build` funciona en cualquier máquina sin dependencias de compilador C.
- Soporta SQL relacional real para las consultas de métricas y reportes.

### Negative
- SQLite no es apto para alta concurrencia real — aceptable para el alcance académico del TP, donde el acceso simultáneo es bajo (4 usuarios como mucho).
- Si más adelante hiciera falta escalar a un entorno multiusuario real, requeriría migrar a un motor cliente-servidor.

## Notes

Confirmado por los 4 integrantes del equipo. Cierra la Tarea 0.2 de la Épica 0 (issue #82).
