# ADR-0003: Cómo se construye la interfaz web

**Date**: 2026-10-25
**Status**: accepted
**Deciders**: Equipo completo (Matías López, Ezequiel Assenza, Nicolás Andino, Lucas Orellana)

## Context

El TP pide una aplicación **web**. Al cerrar el Sprint 1 la aplicación solo expone endpoints que responden JSON (ADR-0001 eligió `net/http` para las rutas, pero no cómo se dibujan las pantallas). El Sprint 2 es el de la interfaz, y desde ahí cada historia tiene que incluir su pantalla.

Restricciones que pesan en la decisión:
- La cátedra exige que el núcleo y las reglas de negocio estén en Go y que no haya lógica de negocio en la interfaz.
- Quedan 6 días de trabajo hasta la entrega (2 de noviembre), así que no hay margen para aprender otro stack.
- El equipo ya trabaja con arquitectura hexagonal: los casos de uso de `internal/app` se pueden llamar desde cualquier adaptador.

## Decision

Páginas HTML generadas en el servidor con `html/template`, el paquete estándar de Go.

- Los handlers de las pantallas viven en `internal/adapters/http`, junto a los de JSON, y llaman a **los mismos casos de uso**: no se duplica lógica ni validaciones.
- Las plantillas comparten un layout base (encabezado, navegación, estilos) que arma HU-30. Cada historia agrega sus propias plantillas y registra sus rutas con su función `RegistrarX`, como las rutas JSON (convención del CLAUDE.md).
- Estilos con un CSS propio y simple, servido como archivo estático. Sin frameworks de CSS ni pasos de build.
- JavaScript solo donde haga falta. El caso conocido son los gráficos del Dashboard (HU-25), que tendrá su propia decisión sobre la librería.
- Los formularios usan `POST` común con redirección después de guardar; los errores de validación se muestran en la misma pantalla con el mensaje del error de dominio.

## Alternatives Considered

### Alternative 1: SPA en JavaScript (React, Vue o similar) que consuma la API JSON
- **Pros**: interfaz más dinámica; la API JSON ya existe.
- **Cons**: suma otro lenguaje, un paso de build (npm, bundler), otra forma de testear y la tentación de repetir validaciones en el cliente.
- **Why not**: con 6 días hasta la entrega, aprender y mantener un segundo stack es el mayor riesgo del plan, y no aporta a lo que evalúa la cátedra (calidad del código Go y reglas de negocio en Go).

### Alternative 2: `html/template` con htmx para actualizar partes de la página
- **Pros**: interactividad sin escribir JavaScript propio.
- **Cons**: una dependencia más y otra forma de pensar las respuestas (fragmentos de HTML).
- **Why not**: no hace falta para lo que piden las historias. Si alguna pantalla lo necesita más adelante, se puede sumar sin cambiar esta decisión.

## Consequences

### Positive
- Todo queda en Go: un solo lenguaje, un solo `go build`, sin dependencias nuevas.
- La regla de "nada de lógica de negocio en la interfaz" se cumple por diseño: las pantallas solo arman datos para las plantillas y llaman a los casos de uso.
- `html/template` escapa el contenido automáticamente, lo que evita inyección de HTML (XSS).
- Los handlers de pantallas se testean con `httptest`, igual que los de JSON.

### Negative
- Cada interacción recarga la página completa; la interfaz es menos fluida que una SPA.
- Para gráficos (HU-25) va a hacer falta algo de JavaScript igual.

## Notes

Decidido por el equipo al planificar el Sprint 2 (2026-10-25). Lo implementa primero HU-30 (interfaz web base), que bloquea las pantallas de las demás historias.
