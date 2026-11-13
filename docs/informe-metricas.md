# Informe de métricas y cobertura

Estado del proyecto al cierre del Sprint 4 (12 de noviembre de 2026).

**De dónde salen los datos:**
- **Cobertura:** el perfil `coverage.out` del CI de `main` en el commit `68eac1f` ([corrida 43239161841](https://github.com/matiaslopez237/ingsoft-tpintegrador/actions/runs/43239161841)), con `go test -race -covermode=atomic -coverprofile=coverage.out ./...`.
- **Tests, commits y PRs:** el repositorio en ese mismo commit y GitHub.
- **Sprints y story points:** el [Product Backlog](scrum/product-backlog.md) y las plannings de [`docs/scrum/`](scrum/).
- **Defectos:** los issues DEF-01 a DEF-05.

## 1. Cobertura

**Total: 90,6 % de las sentencias** (3.933 de 4.340). El CI exige un mínimo del 80 % y falla si no se llega.

| Paquete | Sentencias cubiertas | Cobertura |
|---|---|---|
| `internal/domain` | 523 / 524 | 99,8 % |
| `internal/app` | 575 / 598 | 96,2 % |
| `internal/adapters/http` | 1.909 / 2.044 | 93,4 % |
| `internal/adapters/pdf` | 110 / 110 | 100 % |
| `internal/adapters/persistencia` | 721 / 838 | 86,0 % |
| `internal/adapters/ejemplo` | 86 / 106 | 81,1 % |
| `cmd/server` | 9 / 120 | 7,5 % |
| **Total** | **3.933 / 4.340** | **90,6 %** |

La cobertura es más alta en el centro de la arquitectura: el dominio y los casos de uso tienen casi todo cubierto por tests unitarios.

**Dónde baja, y por qué:**
- **`cmd/server` (7,5 %):** casi todo es `main`, que arma las dependencias y levanta el servidor. Los tests cubren las funciones que se pueden probar por separado: la carga del proyecto de ejemplo y su mensaje de error. El arranque completo se prueba al desplegar en Render.
- **`internal/adapters/ejemplo` (81,1 %):** los métodos `String` de los pasos (solo arman mensajes de error) y el contenido del proyecto de ejemplo, que los tests del paquete reemplazan por uno chico. El contenido real lo recorre la suite BDD de HU-31.
- **`internal/adapters/persistencia` (86,0 %):** `AbrirDB` tiene el 48 %. Lo que falta son las ramas de error al abrir la base o al migrarla, que son difíciles de provocar sin un disco que falle.

**Una aclaración sobre la medición:** el perfil cuenta las sentencias que cubre cada paquete con sus propios tests. La suite BDD (`test/bdd`) recorre los casos de uso y los adaptadores, pero no suma a esos números porque el CI no usa `-coverpkg`. Por eso la cifra es conservadora.

## 2. Tests

### Tests unitarios y de integración en Go

916 funciones `Test`, en 198 archivos. Muchas son tablas con varios casos, así que los casos probados son más.

| Capa | Funciones `Test` | Archivos |
|---|---|---|
| `internal/domain` | 171 | 26 |
| `internal/app` | 188 | 29 |
| `internal/adapters/http` | 372 | 51 |
| `internal/adapters/persistencia` | 154 | 29 |
| `internal/adapters/pdf` | 21 | 2 |
| `internal/adapters/ejemplo` | 5 | 2 |
| `cmd/server` | 4 | 2 |
| `test/bdd` (los *step definitions* de godog) | 1 | 57 |

Los tests de `persistencia` corren contra SQLite real, en un archivo temporal, así que son de integración.

### Escenarios BDD

- **29 archivos `.feature`**, uno por historia con spec, con **378 escenarios**, de los cuales 44 son *Scenario Outline*.
- Por tipo de caso (algunos escenarios tienen más de una etiqueta): 96 `@normal`, 75 `@alternativo`, 82 `@limite` y 139 `@error`.
- Más de la mitad de los escenarios prueba casos límite y errores, no solo el camino feliz.

### Tamaño del código

- **Producción:** 10.033 líneas de Go, sin comentarios ni líneas en blanco.
- **Tests:** 33.158 líneas.
- Hay unas 3,3 líneas de test por cada línea de producción.

### TDD y control de cambios

- **Commits TDD en `main`:** 217 `[RED]`, 195 `[GREEN]` y 35 `[REFACTOR]`. Hay más RED que GREEN porque a veces dos tests en rojo se resolvieron con una sola implementación.
- **Total:** 764 commits sin contar los merges, y 48 PRs mergeados.
- **Reglas de `main`:** está protegida. Cada cambio entra por un PR con la aprobación de otro integrante y el CI en verde (build, `go vet`, `gofmt`, tests unitarios y BDD con `-race`, y la cobertura mínima).

## 3. Métricas del proceso (Scrum)

### Velocidad

| Sprint | Fechas | Objetivo | Planificados | Completados |
|---|---|---|---|---|
| 1 — MVP | 22 al 25 de octubre | Gestión de proyectos, backlog y esfuerzo | 19 | 19 |
| 2 — Interfaz, backlog y sprints | 26 y 27 de octubre | Backlog y sprints desde la interfaz web | 26 | 26 |
| 3 — Estimación y métricas | 29 al 31 de octubre | Planning Poker, esfuerzo y primeras métricas | 28 | 28 |
| 4 — Visualización y cierre | 5 al 13 de noviembre (se cerró el 12) | Dashboard, reportes en PDF y documentación | 28 | 28 |
| **Total** | | | **101** | **101** |

- **Cumplimiento:** el 100 % de los puntos en todos los sprints.
- **Velocidad del equipo:** 27,3 puntos por sprint. Es el promedio de los últimos tres sprints cerrados (26, 28 y 28), con el mismo cálculo que usa la aplicación (HU-22).
- **Tendencia:** la velocidad subió en los primeros tres sprints, mientras el equipo armaba la base, y se mantuvo en 28.
- **Duración del Sprint 4:** fue más largo que los anteriores porque la entrega del TP se postergó. Tuvo más días pero los mismos puntos, así que en él la velocidad por día bajó.

### Estimación

- **Story points del backlog:** el backlog empezó con 115 puntos estimados. Después del Planning Poker quedó en 101: los mismos 93 de las historias originales más 8 de HU-30 y HU-31, que se sumaron después.
- **Historias con estimación inicial:** de las 29, 12 cambiaron en el Planning Poker, todas a menos puntos. En total bajaron de 115 a 93 (−19 %).
- **Dónde más bajó:** en la Épica 4 (Planning Poker, de 20 a 13), en la 7 (Métricas, de 15 a 11) y en la 10 (Documentación, de 8 a 4).

### Horas

El proyecto de ejemplo de la aplicación (HU-31) muestra estas horas del proyecto:

| | Horas |
|---|---|
| Estimadas | 151,5 |
| Reales | 166 |
| Desviación | +14,5 h (+9,6 %) |

La desviación se calcula como en HU-23: (reales − estimadas) / estimadas.

## 4. Defectos

| Defecto | Descripción | Historia | Severidad | Detectado | Sprint de detección | Resuelto | Sprint de resolución | Días hasta resolverlo | Cómo se encontró |
|---|---|---|---|---|---|---|---|---|---|
| DEF-01 | El formulario de edición de proyecto pisa cambios simultáneos | HU-30 | Alta | 26/10 | 2 | 26/10 | 2 | 0 | Revisión del PR #106 |
| DEF-02 | Eliminar un elemento con defectos informa que tiene esfuerzo registrado | HU-06 | Baja | 29/10 | 3 | 29/10 | 3 | 0 | Revisión del PR #185 |
| DEF-03 | La pantalla de un sprint cerrado muestra las acciones de asignar y quitar | HU-09 | Media | 30/10 | 3 | 1/11 | entre el 3 y el 4 | 2 | Prueba en el navegador de HU-12 |
| DEF-04 | Carrera en las migraciones de columnas con arranques simultáneos | HU-11 | Media | 3/11 | entre el 3 y el 4 | 12/11 | 4 | 9 | Al agregar una migración en HU-19 |
| DEF-05 | La lista de defectos y el reporte muestran el ID del sprint en lugar de su número | HU-26 | Baja | 8/11 | 4 | 11/11 | 4 | 3 | Revisión del PR #379, con el servidor real |

- **Total:** 5 defectos, todos cerrados. Ninguno quedó abierto para la entrega.
- **Por severidad:** 1 alta, 2 medias y 2 bajas.
- **Densidad:** 0,05 defectos por story point entregado (5 en 101), o 0,5 por cada mil líneas de producción.
- **Tiempo de resolución:** 2,8 días en promedio. DEF-01 y DEF-02 se resolvieron el mismo día. DEF-04 fue el que más tardó: no bloqueaba nada, así que se dejó para el Sprint 4.
- **Cómo se encontraron:** 3 de los 5 aparecieron en la revisión de un PR. Los otros 2 aparecieron al probar o extender otra historia.
- **Cómo se corrigieron:** cada corrección empezó con un test en rojo que reproducía el defecto, y ese test quedó en la suite como regresión.

## 5. Conclusiones

- **Cobertura y tests:** la cobertura (90,6 %) supera con margen el mínimo del 80 %, y es más alta en el dominio y los casos de uso, que es donde están las reglas de negocio. Cada criterio de aceptación tiene al menos un escenario BDD o un test Go ([trazabilidad](trazabilidad.md)).
- **Velocidad:** el equipo completó todo lo que planificó en los cuatro sprints. La velocidad se estabilizó en 28 puntos, lo que hace confiable planificar con ella.
- **Estimación:** el Planning Poker bajó la estimación inicial un 19 %. Las historias se veían más grandes de lo que fueron, sobre todo las de estimación y métricas.
- **Defectos:** pocos y de severidad mayormente media o baja. La revisión de PRs fue el filtro que más defectos encontró.

## Cómo reproducir los números

```bash
go test -race -timeout 20m -covermode=atomic -coverprofile=coverage.out ./...   # cobertura, como el CI
go tool cover -func=coverage.out | tail -1                                     # total
go tool cover -html=coverage.out -o coverage.html                              # detalle por línea
grep -rhE "^func Test" --include=*_test.go . | wc -l                           # funciones Test
git log --oneline | grep -c "\[RED\]"                                          # commits RED (igual con GREEN y REFACTOR)
```

El CI también publica `coverage.out` y `coverage.html` como artefacto (`cobertura`) en cada corrida.
