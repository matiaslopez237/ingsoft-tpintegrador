# ADR-0001: Framework para la interfaz web

**Date**: 2026-10-22
**Status**: accepted
**Deciders**: Equipo completo (Matías López, Ezequiel Assenza, Nicolás Andino, Lucas Orellana)

## Context

El TP requiere una aplicación web con el núcleo y las reglas de negocio en Go. La arquitectura hexagonal ya definida en `CLAUDE.md` aísla la interfaz web en `internal/adapters/http`, separada del dominio (`internal/domain`) y los casos de uso (`internal/app`). Hace falta elegir cómo se exponen los endpoints HTTP antes de poder empezar a implementar HU-01.

El equipo arranca con Go por primera vez en un proyecto de este tamaño, y la cátedra evalúa explícitamente "Calidad del código Go, Modularidad y mantenibilidad".

## Decision

Usar el paquete estándar `net/http` de Go (sin framework externo), aprovechando el soporte de rutas con método HTTP y parámetros de ruta disponible desde Go 1.22 (ej. `mux.HandleFunc("POST /proyectos/{id}", handler)`).

## Alternatives Considered

### Alternative 1: chi (github.com/go-chi/chi)
- **Pros**: router liviano, middlewares encadenables, muy usado en la comunidad Go.
- **Cons**: dependencia externa a aprender; su valor diferencial frente al router estándar de Go 1.22+ es menor que hace unos años.
- **Why not**: con la arquitectura hexagonal ya aislando la capa HTTP, no se justifica sumar una dependencia para un beneficio marginal.

### Alternative 2: gin (github.com/gin-gonic/gin)
- **Pros**: framework completo, mucha documentación y ejemplos.
- **Cons**: más "mágico" (binding automático, su propio sistema de contexto); agrega una curva de aprendizaje y oculta parte de lo que pasa en cada request.
- **Why not**: para un equipo que recién aprende Go, entender el comportamiento de cada endpoint explícitamente pesa más que la comodidad que ofrece gin.

## Consequences

### Positive
- Cero dependencias externas para la capa HTTP; `go build` no requiere nada adicional.
- El equipo entiende exactamente qué hace cada handler, sin abstracciones intermedias.
- Si en el futuro hace falta migrar a chi/gin, el cambio queda contenido en `internal/adapters/http` sin tocar el dominio ni los casos de uso.

### Negative
- Si se necesitan middlewares compartidos entre muchas rutas (logging, CORS, auth), hay que escribirlos a mano en vez de usar los de una librería.

## Notes

Confirmado por los 4 integrantes del equipo. Cierra la Tarea 0.1 de la Épica 0 (issue #81).
