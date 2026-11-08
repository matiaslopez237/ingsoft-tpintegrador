# ADR-0005: Gráficos del Dashboard como SVG armado en Go

**Date**: 2026-11-06
**Status**: accepted
**Deciders**: Equipo completo (Matías López, Ezequiel Assenza, Nicolás Andino, Lucas Orellana)

## Context

HU-25 pide un Dashboard con al menos 3 métricas clave **de forma gráfica**. Hasta ahora la interfaz no usa JavaScript (ADR-0003), y desde HU-30 las pantallas responden con `Content-Security-Policy: default-src 'self'` (HU-29 la documentó), que bloquea:
- scripts y estilos de otros sitios (un CDN, por ejemplo);
- scripts y estilos escritos dentro de la página (`<script>…</script>`, atributos `onclick`, `style="…"`).

Los gráficos que necesita el Dashboard son simples: barras por sprint (puntos, porcentaje, defectos) y una comparación de dos valores (horas estimadas contra reales). No hacen falta zoom, animaciones ni actualización en vivo.

## Decision

Los gráficos se dibujan como **SVG armado en el servidor**, con `html/template`, dentro de la misma página HTML:
- El adaptador HTTP calcula la geometría (alto y posición de cada barra, escala del eje) a partir de los datos que devuelve el caso de uso, y la plantilla escribe los elementos `<svg>`, `<rect>` y `<text>`.
- Los colores y tamaños de letra van en `estilos.css` con clases (`class="barra-planificados"`), no con atributos `style`, para respetar la CSP.
- Cada gráfico tiene `role="img"`, un `<title>` y una tabla con los mismos datos debajo, para lectores de pantalla y para quien quiera el número exacto.
- La geometría se testea con tests de Go comunes; la pantalla, con `httptest`, igual que el resto.

## Alternatives Considered

### Alternative 1: Chart.js (u otra librería JS) copiada en `/static`
- **Pros**: gráficos más vistosos, con tooltips y leyendas interactivas sin trabajo extra.
- **Cons**: suma JavaScript propio (leer los datos de un endpoint JSON y armar el gráfico), una dependencia de ~200 KB que hay que actualizar a mano, y una parte de la pantalla que no se puede probar con `httptest`.
- **Why not**: la interactividad no la pide ninguna historia, y sumar JavaScript va contra la decisión de ADR-0003 sin una razón de peso.

### Alternative 2: Librería de gráficos en Go que genera imágenes (por ejemplo `go-chart`)
- **Pros**: gráficos completos sin JavaScript.
- **Cons**: una dependencia externa más, imágenes PNG que no escalan ni se leen con lector de pantalla, y una ruta aparte para servir cada imagen.
- **Why not**: para barras simples, el SVG a mano es menos código que adaptar la librería, y el resultado es más accesible.

### Alternative 3: Barras hechas con HTML y CSS (`<div>` con ancho variable)
- **Pros**: no requiere SVG.
- **Cons**: el ancho de cada barra tendría que ir en un atributo `style`, que la CSP bloquea; habría que generar una clase por porcentaje.
- **Why not**: SVG permite poner la medida en atributos (`height`, `y`) sin tocar la CSP.

## Consequences

### Positive
- Sin JavaScript ni dependencias nuevas: un solo `go build`, como hasta ahora.
- La CSP queda igual; no hace falta abrir excepciones.
- Todo se testea en Go, incluida la escala de los gráficos.
- El Dashboard muestra los datos vigentes en cada pedido, sin cache del lado del cliente.

### Negative
- Los gráficos no son interactivos (sin tooltips ni zoom). Se compensa con la tabla de datos debajo de cada uno.
- La geometría de las barras (escala, márgenes, etiquetas) es código nuestro que hay que mantener.

### Risks
- Si una historia futura necesita gráficos interactivos, habría que revisar esta decisión. Se puede sumar una librería en `/static` sin rehacer lo existente, porque los datos ya salen de los casos de uso.

## Notes

Propuesto por Nicolás Andino al empezar HU-25 (tarea 25.1, #309). Cierra el pendiente que dejaron ADR-0003 (actualización del 2026-11-05) y `docs/arquitectura.md`.

**Actualización 2026-11-07:** aceptado. Lo aprobó Matías López (Agile Enabler) el 2026-11-07, en la revisión del PR #375.
