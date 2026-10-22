# Tareas — Sprint 1 (MVP)

Alcance del sprint: Épica 1 completa (HU-01 a HU-04) + HU-18. Total: **14 story points**.

---

## Tarea 0 — Prerrequisitos técnicos (bloquean todo lo demás)

Estas no son de una historia puntual; habilitan que cualquiera pueda empezar a programar.

**Tarea 0.1 — Decidir framework/enfoque de la interfaz web**
- **Qué hacer:** reunión de los 4, evaluar `net/http` estándar vs. un router como `chi` o `gin`, y registrar la decisión como ADR (`docs/adr/0001-framework-web.md`, usando la skill `architecture-decision-records`).
- **Terminada cuando:** existe el archivo ADR con la decisión y el motivo, y los 4 integrantes están de acuerdo (o quedó asentado quién no y por qué).

**Tarea 0.2 — Decidir mecanismo de persistencia**
- **Qué hacer:** evaluar SQLite vs. otra opción simple, y registrar como ADR (`docs/adr/0002-persistencia.md`).
- **Terminada cuando:** existe el ADR, y se probó que se puede conectar desde Go (aunque sea con un `SELECT 1`).

**Tarea 0.3 — Esqueleto del proyecto**
- **Qué hacer:** `go mod init`, crear carpetas `cmd/server/`, `internal/domain/`, `internal/app/`, `internal/adapters/` (ya definidas en `CLAUDE.md`), un `main.go` que levante un servidor HTTP vacío.
- **Terminada cuando:** `go run ./cmd/server` levanta el servidor sin errores (aunque no tenga rutas todavía).

**Tarea 0.4 — Conexión a la base de datos**
- **Qué hacer:** configurar el driver elegido en 0.2, variables de entorno o archivo de config, una migración inicial vacía.
- **Terminada cuando:** el servidor arranca y logea "conectado a la base" sin error.

**Tarea 0.5 — Lint y test base**
- **Qué hacer:** correr `go vet ./...`, `gofmt -l .` y `go test ./...` sobre el esqueleto vacío.
- **Terminada cuando:** las 3 salidas están limpias (sin errores, sin archivos mal formateados, 0 tests pero sin fallar).

**Sin esto cerrado, ninguna tarea de HU-01 en adelante puede arrancar en serio.**

---

## HU-01: Crear un proyecto (3 pts)

**Tarea 1.1 — Spec SDD**
- **Archivo a crear:** `docs/specs/HU-01-crear-proyecto.md`
- **Qué tiene que contener:** Objetivo, Entradas (nombre, fecha_inicio, fecha_fin), Salidas (proyecto creado con id + estado), Reglas de negocio (nombre no vacío; fecha_fin ≥ fecha_inicio; estado inicial "Activo"), Restricciones, Casos límite (fechas iguales), Condiciones de error (nombre vacío, fechas invertidas), Criterios de aceptación.
- **Terminada cuando:** el archivo cubre las 8 secciones y al menos un compañero la revisó.

**Tarea 1.2 — Escenarios BDD**
- **Archivo a crear:** `features/HU-01-crear-proyecto.feature`
- **Escenarios mínimos:** crear con datos válidos (queda "Activo"); crear sin nombre (error); crear con fecha_fin anterior a fecha_inicio (error).
- **Terminada cuando:** el `.feature` tiene los 3 escenarios escritos en Given-When-Then (todavía no tienen que pasar).

**Tarea 1.3 — Steps de godog**
- **Archivo a crear:** `test/bdd/hu01_steps_test.go`
- **Qué hacer:** implementar en Go cada línea Given/When/Then de los 3 escenarios, conectándolas a `CrearProyecto` (que todavía no existe).
- **Terminada cuando:** `go test ./test/bdd/ -v` corre y los 3 escenarios **fallan** (rojo esperado, porque falta implementar).

**Tarea 1.4 — TDD: entidad `Proyecto`**
- **Archivos:** `internal/domain/proyecto.go` + `internal/domain/proyecto_test.go`
- **Tests a escribir primero:** `TestNuevoProyecto_DatosValidos_CreaProyectoActivo`, `TestNuevoProyecto_NombreVacio_RetornaError`, `TestNuevoProyecto_FechaFinAnteriorAInicio_RetornaError`, `TestNuevoProyecto_FechasIguales_EsValido`.
- **Terminada cuando:** los 4 tests están en verde, con al menos 2 commits (`test:` rojo, `feat:` verde).

**Tarea 1.5 — TDD: caso de uso `CrearProyecto`**
- **Archivo:** `internal/app/crear_proyecto.go` + test
- **Qué hace:** recibe los datos, llama a `domain.NuevoProyecto(...)`, y si es válido lo manda a guardar a través de una interfaz de repositorio (sin implementación real todavía).
- **Terminada cuando:** hay un test que verifica la orquestación usando un repositorio falso (mock/fake).

**Tarea 1.6 — Adaptador de persistencia**
- **Archivo:** `internal/adapters/persistencia/proyecto_repo.go`
- **Qué hace:** implementa la interfaz de repositorio de 1.5, guardando/leyendo de la base elegida en 0.2.
- **Terminada cuando:** hay un test de integración que crea un proyecto real contra la base (puede ser SQLite en memoria para test) y lo recupera igual a como se guardó.

**Tarea 1.7 — Adaptador HTTP**
- **Archivo:** `internal/adapters/http/proyecto_handler.go`
- **Qué hace:** expone `POST /proyectos`, parsea JSON, llama al caso de uso, devuelve 201 + proyecto creado o 400 + error.
- **Terminada cuando:** un `curl -X POST .../proyectos` con datos válidos devuelve 201, con nombre vacío devuelve 400, y los 3 escenarios BDD de la tarea 1.3 ahora pasan en verde.

**Tarea 1.8 — Revisión**
- **Qué hacer:** correr `go vet ./...` y `gofmt -l .`, pedir revisión al agente `go-reviewer` sobre el diff completo de HU-01.
- **Terminada cuando:** no quedan observaciones pendientes.

**Tarea 1.9 — Trazabilidad**
- **Qué hacer:** actualizar la matriz de trazabilidad (skill `trazabilidad`) agregando HU-01 con sus links a spec, feature, tests y commits.
- **Terminada cuando:** se puede seguir Historia → Spec → Criterios → BDD → Tests → Código sin buscar nada a mano.

---

## HU-02: Modificar un proyecto (2 pts)

**Tarea 2.1 — Spec SDD**
- **Archivo a crear:** `docs/specs/HU-02-modificar-proyecto.md`
- **Qué tiene que contener:** Objetivo, Entradas (id de proyecto, nuevos nombre/fechas), Salidas (proyecto actualizado), Reglas de negocio (mismas validaciones que HU-01: nombre no vacío, fecha_fin ≥ fecha_inicio), Restricciones (no se puede editar un proyecto inexistente), Casos límite, Condiciones de error, Criterios de aceptación.
- **Terminada cuando:** el archivo cubre las 8 secciones y fue revisado por un compañero.

**Tarea 2.2 — Escenarios BDD**
- **Archivo a crear:** `features/HU-02-modificar-proyecto.feature`
- **Escenarios mínimos:** editar con datos válidos (se actualiza); editar con fecha_fin inválida (error); editar un proyecto que no existe (error).
- **Terminada cuando:** los 3 escenarios están escritos en Given-When-Then.

**Tarea 2.3 — Steps de godog**
- **Archivo a crear:** `test/bdd/hu02_steps_test.go`
- **Terminada cuando:** corren los 3 escenarios y fallan en rojo (porque `ModificarProyecto` no existe aún).

**Tarea 2.4 — TDD: método de actualización en el dominio**
- **Archivos:** `internal/domain/proyecto.go` (agregar método `Actualizar(...)`) + test en `proyecto_test.go`
- **Tests a escribir:** `TestProyecto_Actualizar_DatosValidos`, `TestProyecto_Actualizar_FechaFinInvalida_RetornaError` (reutiliza la misma validación de fechas de HU-01, no la duplica).
- **Terminada cuando:** los tests están en verde y no se copió/pegó la lógica de validación de fechas (se llama a la misma función que usa HU-01).

**Tarea 2.5 — TDD: caso de uso `ModificarProyecto`**
- **Archivo:** `internal/app/modificar_proyecto.go` + test
- **Qué hace:** busca el proyecto por id en el repositorio, aplica `Actualizar(...)`, lo vuelve a guardar; si no existe, devuelve error "no encontrado".
- **Terminada cuando:** hay tests para el caso éxito y el caso "proyecto no encontrado", con repositorio falso.

**Tarea 2.6 — Adaptador de persistencia: actualizar**
- **Archivo:** `internal/adapters/persistencia/proyecto_repo.go` (agregar método `Actualizar`)
- **Terminada cuando:** hay un test de integración que crea un proyecto, lo actualiza, y verifica que la base refleja el cambio.

**Tarea 2.7 — Adaptador HTTP**
- **Archivo:** `internal/adapters/http/proyecto_handler.go` (agregar ruta)
- **Qué hace:** expone `PUT /proyectos/:id`, devuelve 200 + proyecto actualizado, 400 si los datos son inválidos, 404 si no existe.
- **Terminada cuando:** los 3 escenarios BDD de la tarea 2.3 pasan en verde.

**Tarea 2.8 — Revisión + trazabilidad**
- **Terminada cuando:** `go vet`/`gofmt` limpios, revisión de `go-reviewer` resuelta, y HU-02 agregada a la matriz de trazabilidad.

---

## HU-03: Registrar integrantes de un proyecto (3 pts)

**Tarea 3.1 — Spec SDD**
- **Archivo a crear:** `docs/specs/HU-03-registrar-integrantes.md`
- **Qué tiene que contener:** Objetivo, Entradas (id de proyecto, nombre del integrante, rol), Salidas (integrante agregado/quitado), Reglas de negocio (nombre obligatorio; un proyecto puede tener múltiples integrantes; no se duplica el mismo nombre dos veces en el mismo proyecto), Restricciones, Casos límite (quitar al único integrante), Condiciones de error, Criterios de aceptación.
- **Terminada cuando:** el archivo cubre las 8 secciones y fue revisado.

**Tarea 3.2 — Escenarios BDD**
- **Archivo a crear:** `features/HU-03-registrar-integrantes.feature`
- **Escenarios mínimos:** agregar integrante válido; agregar sin nombre (error); quitar un integrante existente; quitar un integrante que no pertenece al proyecto (error).
- **Terminada cuando:** los 4 escenarios están escritos.

**Tarea 3.3 — Steps de godog**
- **Archivo a crear:** `test/bdd/hu03_steps_test.go`
- **Terminada cuando:** los 4 escenarios corren y fallan en rojo.

**Tarea 3.4 — TDD: entidad `Integrante` y relación con `Proyecto`**
- **Archivos:** `internal/domain/integrante.go`, `internal/domain/proyecto.go` (agregar lista de integrantes) + tests
- **Tests a escribir:** `TestProyecto_AgregarIntegrante_Valido`, `TestProyecto_AgregarIntegrante_SinNombre_RetornaError`, `TestProyecto_QuitarIntegrante_Existente`, `TestProyecto_QuitarIntegrante_Inexistente_RetornaError`.
- **Terminada cuando:** los 4 tests están en verde.

**Tarea 3.5 — TDD: casos de uso `AgregarIntegrante` / `QuitarIntegrante`**
- **Archivo:** `internal/app/gestionar_integrantes.go` + tests
- **Terminada cuando:** ambos casos de uso tienen test con repositorio falso.

**Tarea 3.6 — Adaptador de persistencia**
- **Archivo:** `internal/adapters/persistencia/proyecto_repo.go` (persistir integrantes junto al proyecto, o tabla separada según el diseño elegido)
- **Terminada cuando:** hay un test de integración que agrega y quita un integrante contra la base real (o SQLite en memoria).

**Tarea 3.7 — Adaptador HTTP**
- **Archivo:** `internal/adapters/http/integrante_handler.go`
- **Qué hace:** expone `POST /proyectos/:id/integrantes` y `DELETE /proyectos/:id/integrantes/:nombre`.
- **Terminada cuando:** los 4 escenarios BDD de la tarea 3.3 pasan en verde.

**Tarea 3.8 — Revisión + trazabilidad**
- **Terminada cuando:** `go vet`/`gofmt` limpios, revisión resuelta, HU-03 en la matriz de trazabilidad.

---

## HU-04: Consultar el estado de un proyecto (3 pts)

**Tarea 4.1 — Spec SDD**
- **Archivo a crear:** `docs/specs/HU-04-consultar-proyecto.md`
- **Qué tiene que contener:** Objetivo, Entradas (id de proyecto), Salidas (datos del proyecto + lista de integrantes + estado), Reglas de negocio (ninguna de cálculo, es sólo lectura), Restricciones, Casos límite (proyecto sin integrantes todavía), Condiciones de error (id inexistente), Criterios de aceptación.
- **Terminada cuando:** el archivo cubre las 8 secciones y fue revisado.

**Tarea 4.2 — Escenarios BDD**
- **Archivo a crear:** `features/HU-04-consultar-proyecto.feature`
- **Escenarios mínimos:** consultar un proyecto existente (devuelve sus datos); consultar un proyecto inexistente (error 404).
- **Terminada cuando:** los 2 escenarios están escritos.

**Tarea 4.3 — Steps de godog**
- **Archivo a crear:** `test/bdd/hu04_steps_test.go`
- **Terminada cuando:** los 2 escenarios corren y fallan en rojo.

**Tarea 4.4 — TDD: caso de uso `ConsultarProyecto`**
- **Archivo:** `internal/app/consultar_proyecto.go` + test
- **Tests a escribir:** `TestConsultarProyecto_Existente_DevuelveDatos`, `TestConsultarProyecto_Inexistente_RetornaError`.
- **Terminada cuando:** ambos tests están en verde, usando repositorio falso.

**Tarea 4.5 — Adaptador HTTP**
- **Archivo:** `internal/adapters/http/proyecto_handler.go` (agregar ruta)
- **Qué hace:** expone `GET /proyectos/:id`, devuelve 200 + datos o 404 si no existe.
- **Terminada cuando:** los 2 escenarios BDD de la tarea 4.3 pasan en verde.

**Tarea 4.6 — Revisión + trazabilidad**
- **Terminada cuando:** `go vet`/`gofmt` limpios, revisión resuelta, HU-04 en la matriz de trazabilidad.

---

## HU-18: Registrar esfuerzo en una historia o tarea (3 pts)

**Tarea 18.1 — Spec SDD**
- **Archivo a crear:** `docs/specs/HU-18-registrar-esfuerzo.md`
- **Qué tiene que contener:** Objetivo, Entradas (integrante, fecha, actividad, horas, referencia a historia/tarea), Salidas (registro de esfuerzo guardado), Reglas de negocio (horas > 0; fecha no puede ser futura; todos los campos obligatorios), Restricciones, Casos límite (0.5 horas, historia sin esfuerzo registrado todavía), Condiciones de error (horas negativas, campo faltante), Criterios de aceptación.
- **Terminada cuando:** el archivo cubre las 8 secciones y fue revisado.

**Tarea 18.2 — Escenarios BDD**
- **Archivo a crear:** `features/HU-18-registrar-esfuerzo.feature`
- **Escenarios mínimos:** registrar esfuerzo válido; registrar con horas negativas (error); registrar sin actividad (error).
- **Terminada cuando:** los 3 escenarios están escritos.

**Tarea 18.3 — Steps de godog**
- **Archivo a crear:** `test/bdd/hu18_steps_test.go`
- **Terminada cuando:** los 3 escenarios corren y fallan en rojo.

**Tarea 18.4 — TDD: entidad `RegistroEsfuerzo`**
- **Archivos:** `internal/domain/registro_esfuerzo.go` + test
- **Tests a escribir:** `TestNuevoRegistroEsfuerzo_Valido`, `TestNuevoRegistroEsfuerzo_HorasNegativas_RetornaError`, `TestNuevoRegistroEsfuerzo_SinActividad_RetornaError`.
- **Terminada cuando:** los 3 tests están en verde.

**Tarea 18.5 — TDD: caso de uso `RegistrarEsfuerzo`**
- **Archivo:** `internal/app/registrar_esfuerzo.go` + test
- **Qué hace:** valida que la historia/tarea referenciada exista, crea el `RegistroEsfuerzo`, lo guarda.
- **Terminada cuando:** hay test para el caso éxito y el caso "historia inexistente", con repositorio falso.

**Tarea 18.6 — Adaptador de persistencia**
- **Archivo:** `internal/adapters/persistencia/esfuerzo_repo.go`
- **Terminada cuando:** hay un test de integración que guarda un registro y lo recupera igual, contra la base real o SQLite en memoria.

**Tarea 18.7 — Adaptador HTTP**
- **Archivo:** `internal/adapters/http/esfuerzo_handler.go`
- **Qué hace:** expone `POST /esfuerzo`, devuelve 201 + registro creado o 400 + error.
- **Terminada cuando:** los 3 escenarios BDD de la tarea 18.3 pasan en verde.

**Tarea 18.8 — Revisión + trazabilidad**
- **Terminada cuando:** `go vet`/`gofmt` limpios, revisión resuelta, HU-18 en la matriz de trazabilidad.

---

## Notas para la revisión en equipo

- **No está definido quién hace qué.** Conviene que cada HU la lleve una sola persona de punta a punta, para no pisarse commits.
- Las tareas 0.1/0.2 son decisiones de arquitectura que afectan a **todo** el proyecto, no sólo Sprint 1 — esa conversación va primero, con los 4 presentes.
- Los nombres de archivos/funciones son un ejemplo para que quede claro el "qué hacer", no una obligación — el equipo puede nombrarlos distinto al implementar, mientras quede la spec, el feature, los tests y la trazabilidad (eso sí lo pide la cátedra).
- Las estimaciones en horas son mi primera aproximación — ajústenlas si algo parece más/menos complejo de lo que parece acá.
- Esto está guardado en `docs/scrum/sprint-1/tareas-sprint-1.md`, todavía **no** subido (`git push`) ni cargado como sub-issues en GitHub.
