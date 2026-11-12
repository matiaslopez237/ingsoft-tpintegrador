# HU-31: Proyecto de ejemplo precargado

**Historia:** Como integrante del equipo, quiero que la aplicación arranque con un proyecto de ejemplo ya cargado, para mostrar cómo se ve un proyecto completo sin tener que cargarlo a mano.
**Issue:** #367  ·  **Sprint:** 4  ·  **Story points:** 3
**Estado de la spec:** Implementada

## 1. Objetivo
Una aplicación recién desplegada está vacía, y en el plan gratuito de Render la base se borra en cada despliegue (ADR-0004). Para mostrar el sistema hay que cargar a mano un proyecto, su backlog, sus sprints, el esfuerzo y los defectos antes de cada presentación. Esta historia carga sola un proyecto de ejemplo completo al arrancar, para que todas las pantallas y métricas tengan datos desde el primer momento. El ejemplo es este mismo TP: "Software Metrics & Estimation".

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| Variable de entorno `CARGAR_EJEMPLO` | `string` | No | Sin los espacios de alrededor, se interpreta con `strconv.ParseBool` (`true`, `1`, `TRUE`…). Ausente o vacía equivale a `false`. Un valor que `ParseBool` no reconoce no carga el ejemplo (RN-31.3). |
| Proyectos existentes en la base | `[]domain.Proyecto` | — | Se leen con el caso de uso `ListarProyectos` (HU-30) antes de cargar. |

No hay entradas del usuario: la carga ocurre al arrancar el servidor, antes de atender pedidos.

## 3. Salidas esperadas
- **Con la variable activada y la base sin proyectos:** queda cargado el proyecto de ejemplo con este contenido:
  - **Proyecto:** "Software Metrics & Estimation", `Activo`, con inicio el 2026-10-22 (P-1) y sin fecha de fin. El Sprint 1 va del 22 al 25 de octubre.
  - **Integrantes:** Matías López (Agile Enabler), Ezequiel Assenza, Nicolás Andino y Lucas Orellana (Product Builder).
  - **Product Backlog:** las historias HU-01 a HU-31 como elementos, con la descripción y los criterios de aceptación de su issue y los story points del Product Backlog. Como el Product Backlog no tiene prioridades, la prioridad sale del sprint: Alta en los Sprints 1 y 2, Media en el 3 y Baja en el 4.
  - **Sprints 1 a 3 `Cerrado`**, con las fechas y las historias reales de cada uno, todas `Hecho`: 19, 26 y 28 puntos planificados y completados.
  - **Sprint 4 `Abierto`** (del 2026-11-05 al 2026-11-13, con las extensiones del sprint real), con sus historias, HU-24 y HU-29 `Hecho` y el resto sin terminar (P-3).
  - **Esfuerzo:** horas estimadas y registros de esfuerzo en las historias de los sprints, para que HU-19 y HU-23 muestren desviaciones. Cada historia registra el esfuerzo a nombre del integrante que más commits hizo en ella.
  - **Planning Poker:** rondas con votos y estimación acordada en historias del Sprint 3, incluida una historia que necesitó una segunda ronda por diferencias.
  - **Defectos:** DEF-01 a DEF-04, sobre las historias donde aparecieron y con sus fechas reales, así que los sprints de detección y resolución salen de esas fechas (RN-20.5, RN-21.6):

    | Defecto | Detectado | Sprint de detección | Resuelto | Sprint de resolución | Estado |
    |---|---|---|---|---|---|
    | DEF-01 | 26/10 | 2 | 26/10 | 2 | `Cerrado` |
    | DEF-02 | 29/10 | 3 | 29/10 | 3 | `Cerrado` |
    | DEF-03 | 30/10 | 3 | 1/11 | ninguno: entre el Sprint 3 y el 4 | `Cerrado` |
    | DEF-04 | 3/11 | ninguno: entre el Sprint 3 y el 4 | — | — | `Abierto` |
- **En cualquier otro caso** (variable apagada, valor inválido o base con algún proyecto), la base no cambia.
- **En el log:** una línea que dice si se cargó el ejemplo o por qué no.
- No hay rutas ni pantallas nuevas: el ejemplo se ve en las pantallas que ya existen.

## 4. Reglas de negocio
- **RN-31.1:** El ejemplo se carga solo si `CARGAR_EJEMPLO` está activada **y** la base no tiene ningún proyecto. Con al menos un proyecto, no se carga nada, aunque ese proyecto no sea el de ejemplo.
- **RN-31.2:** La carga usa los **casos de uso existentes** (`CrearProyecto`, `AgregarIntegrante`, `CrearElementoBacklog`, `EditarElementoBacklog`, `CrearSprint`, `AsignarHistoria`, `CompletarHistoria`, `CerrarSprint`, `RegistrarEsfuerzo`, `RegistrarDefecto`, `ActualizarEstadoDefecto`, `IniciarRonda`, `Votar` y `RegistrarEstimacion`), nunca SQL directo. Así el ejemplo cumple todas las reglas de negocio de esas historias, y si una regla cambia, el test de la carga lo detecta.
- **RN-31.3:** Si `CARGAR_EJEMPLO` tiene un valor que no se reconoce como booleano, no se carga el ejemplo y se avisa en el log. No impide arrancar.
- **RN-31.4:** Los sprints se arman en orden cronológico, como pasaron: se crea el sprint, se le asignan sus historias, se registran el esfuerzo y los defectos con fechas dentro del sprint, se completan las historias y se cierra, antes de crear el siguiente. Así los puntos de cada cierre (RN-11.6) y los sprints de detección y resolución de los defectos (RN-20.5, RN-21.6) salen de las mismas reglas que en el uso normal.
- **RN-31.5:** Para los casos de uso que dependen de "hoy" (`RegistrarDefecto`, `ActualizarEstadoDefecto` y `RegistrarEsfuerzo`), la carga les pasa un reloj con la fecha del ejemplo que corresponde, dentro del sprint en curso. Los demás casos de uso no cambian.
- **RN-31.6:** La carga ocurre antes de que el servidor empiece a atender pedidos, así que no hay pedidos simultáneos que la afecten. Usa el mismo contexto que el apagado del servidor: con `Ctrl+C` o `SIGTERM` durante la carga, se corta.

## 5. Restricciones
- **Dónde va:** la carga es un adaptador que maneja la aplicación desde afuera, como el HTTP: `internal/adapters/ejemplo`. Recibe los casos de uso ya construidos y no depende de SQLite. El contenido del ejemplo son datos en Go, en el mismo paquete.
- **Conexión:** `cmd/server/main.go` lee `CARGAR_EJEMPLO` y llama a `cargarEjemplo` después de abrir la base y antes de `ListenAndServe`. Esa función arma los casos de uso de la carga, con su propio reloj. `render.yaml` define `CARGAR_EJEMPLO=true`. En las máquinas del equipo la variable no está definida, salvo que alguien la quiera.
- **Tests:** la carga completa se prueba contra una base SQLite real en una carpeta temporal, y se verifica el resultado con los casos de uso de consulta (velocidad de HU-22, horas de HU-23, defectos de HU-21 y rondas de HU-16). Así, si una regla de otra historia cambia y el ejemplo deja de cumplirla, el CI falla.
- **Fechas:** todas las fechas del ejemplo son anteriores a la fecha en que corre la aplicación, así que cumplen RN-18.4 (el esfuerzo no puede ser futuro).
- **Sin datos nuevos en el dominio:** no se agregan entidades, columnas ni reglas al dominio.
- **Fuera de alcance:** borrar o reiniciar el ejemplo desde la aplicación, cargar más de un proyecto de ejemplo y actualizar el ejemplo cuando cambie el TP real (se hace editando los datos en Go).

## 6. Casos límite
- Base vacía y variable en `true` → se carga el ejemplo completo.
- Base vacía y variable ausente, vacía o `false` → no se carga nada.
- Base vacía y variable `si` u otro valor no booleano → no se carga nada y queda un aviso en el log (RN-31.3).
- Base con un proyecto cualquiera (por ejemplo, uno creado a mano) y variable en `true` → no se carga nada (RN-31.1).
- Reiniciar el servidor con el ejemplo ya cargado → no se duplica, porque la base ya tiene un proyecto.
- Base de una versión anterior, a la que las migraciones le agregan columnas (ADR-0002) → sigue valiendo RN-31.1: solo importa si hay proyectos.
- El Sprint 4 sigue `Abierto` después de su fecha de fin (2026-11-13) → es válido: un sprint se cierra cuando el equipo lo decide (RN-11.2).
- `Ctrl+C` o `SIGTERM` durante la carga → la carga se corta y el servidor no arranca, igual que si un paso falla (sección 7).
- Arrancar de nuevo después de una carga que falló → la base tiene el proyecto a medias y RN-31.1 no vuelve a cargarlo. Por eso el mensaje de la falla pide borrar el archivo de `DB_PATH` antes de volver a arrancar. En Render no hace falta: cada despliegue arranca con la base vacía.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| `CARGAR_EJEMPLO` con un valor no booleano | — (no es un error de dominio) | No se carga y se registra en el log "CARGAR_EJEMPLO no es un valor booleano: no se carga el ejemplo". El servidor arranca igual. |
| No se pudo leer la lista de proyectos | El error del caso de uso, envuelto | Igual que una falla al cargar (fila siguiente). |
| Un caso de uso devuelve un error durante la carga, o `Ctrl+C`/`SIGTERM` la corta | El error del caso de uso (o `context.Canceled`), envuelto con el paso que falló (por ejemplo, "cargar ejemplo: cerrar el sprint 2: …") | El servidor no arranca (P-2). El log muestra el error y dice que la base quedó con el ejemplo a medias: "borrá <DB_PATH> antes de volver a arrancar". Con los tests de la carga en el CI, un error de un paso solo puede pasar si cambia una regla de otra historia sin actualizar el ejemplo. |

## 8. Criterios de aceptación
- **CA-31.1:** Dada una base vacía y `CARGAR_EJEMPLO=true`, cuando arranca la carga, entonces existe el proyecto "Software Metrics & Estimation", `Activo`, con 4 integrantes y sus roles. (RN-31.1)
- **CA-31.2:** Dado el ejemplo cargado, entonces el backlog tiene los elementos HU-01 a HU-31 con los story points del Product Backlog, y hay 4 sprints: del 1 al 3 `Cerrado`, con 19, 26 y 28 puntos planificados y completados, y el 4 `Abierto`. (RN-31.4)
- **CA-31.3:** Dado el ejemplo cargado, cuando se consulta la velocidad del equipo (HU-22), entonces es el promedio de 19, 26 y 28. (RN-31.4)
- **CA-31.4:** Dado el ejemplo cargado, cuando se consultan los defectos (HU-21), entonces están DEF-01 a DEF-04 con los estados y los sprints de detección y resolución de la tabla de la sección 3. (RN-31.4, RN-31.5)
- **CA-31.5:** Dado el ejemplo cargado, cuando se consultan las horas del proyecto (HU-23), entonces hay horas estimadas y reales y una desviación; y cuando se consulta el historial de rondas de la historia que se votó dos veces (HU-16), entonces tiene dos rondas, la última con estimación acordada. (RN-31.2, RN-31.5)
- **CA-31.6:** Dada una base con un proyecto y `CARGAR_EJEMPLO=true`, cuando arranca la carga, entonces la base no cambia. (RN-31.1)
- **CA-31.7:** Dada una base vacía y `CARGAR_EJEMPLO` ausente, vacía o `false`, cuando arranca el servidor, entonces no se carga nada. (RN-31.1)
- **CA-31.8:** Dada una base vacía y `CARGAR_EJEMPLO=si`, cuando arranca el servidor, entonces no se carga nada, el servidor arranca y el log avisa que el valor no es booleano. (RN-31.3)
- **CA-31.9:** Dado el ejemplo ya cargado, cuando se vuelve a ejecutar la carga, entonces no se duplica nada: sigue habiendo un solo proyecto. (RN-31.1)
- **CA-31.10:** Dado un caso de uso que falla durante la carga, entonces la carga devuelve un error que nombra el paso que falló. (Sección 7)
- **CA-31.11:** `render.yaml` define `CARGAR_EJEMPLO=true`, y después de desplegar, la URL de Render muestra el proyecto de ejemplo. (Sección 5)

## Decisiones tomadas
Confirmadas por Matías López (Agile Enabler) el 2026-11-05:
- **Carga automática y contenido:** la carga es automática al arrancar, con `CARGAR_EJEMPLO`, y el contenido es este mismo TP (decididas al crear la historia).
- **P-1 — Fechas:** el proyecto empieza el 2026-10-22 y el Sprint 1 va del 22 al 25 de octubre. El Sprint 0 (preparación) no se carga: no tuvo historias, y como la aplicación numera los sprints desde 1, cargarlo correría la numeración respecto de la real.
- **P-2 — Falla durante la carga:** el servidor no arranca y muestra el error. Un ejemplo a medias es peor que ninguno, y los tests de la carga en el CI detectan la falla antes de mergear.
- **P-3 — Sprint 4:** queda `Abierto`, con HU-24 y HU-29 `Hecho` y el resto sin terminar, para poder mostrar las acciones de un sprint abierto. Cuando el equipo cierre el Sprint 4 real, se cierra también en el ejemplo, en un cambio aparte que actualiza esta spec.

## Preguntas abiertas
Ninguna.

## Historial
- 2026-11-05: creada (borrador). P-1 a P-3, propuestas pendientes de confirmar por Matías. La carga automática con `CARGAR_EJEMPLO` y el uso de este TP como contenido los decidió Matías al crear la historia.
- 2026-11-05: Matías confirma P-1 a P-3 y aprueba la spec.
- 2026-11-05: los defectos del ejemplo usan sus fechas reales. DEF-03 se resolvió y DEF-04 se detectó entre el Sprint 3 y el 4, así que no tienen sprint de resolución y de detección, respectivamente.
- 2026-11-07: implementada. La prioridad de cada historia sale de su sprint, porque el Product Backlog no tiene prioridades, y el esfuerzo se registra a nombre de quien más commits hizo en cada historia. Los tests verifican los defectos (HU-21) y las rondas (HU-16), en lugar de la calidad de HU-24.
- 2026-11-08: en la revisión del PR #377, Nico señaló que una carga que falla deja el ejemplo a medias y el próximo arranque no lo vuelve a cargar. El mensaje de la falla pide borrar la base, `Ctrl+C` durante la carga la corta, y se agregaron los casos límite. CA-31.5 verifica los valores exactos de las horas.
- 2026-11-11: el Sprint 4 real se extendió hasta el 13 de noviembre, y el ejemplo lo refleja (P-3).
