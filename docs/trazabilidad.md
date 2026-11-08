# Matriz de trazabilidad

Historia → Spec SDD → Criterios de aceptación → Escenarios BDD → Tests Go → Código Go → Commits TDD.

Cómo se busca cada eslabón: tags `@HU-XX` / `@CA-XX.N` en `features/`, comentarios `// Traza: HU-XX ...` en los tests, `Implementa HU-XX` en el código y `git log --grep="(HU-XX)"`.

## Resumen

| Historia | Spec | CA con BDD | CA con test Go | Código | TDD en git | BDD en verde |
|---|---|---|---|---|---|---|
| HU-01 Crear un proyecto | ✓ (Borrador, pendiente de aprobación en el PR) | 9/9 | 9/9 | ✓ | RED → GREEN ×2, REFACTOR ✓ | ✓ 14/14 escenarios |
| HU-02 Modificar un proyecto | ✓ (Borrador, pendiente de aprobación en el PR) | 12/12 | 12/12 | ✓ | RED → GREEN ×7 | ✓ 24/24 escenarios |
| HU-03 Registrar integrantes | ✓ (Borrador, pendiente de aprobación en el PR) | 10/10 | 10/10 | ✓ | RED → GREEN ×10 | ✓ 13/13 escenarios |
| HU-04 Consultar el estado de un proyecto | ✓ (Borrador, pendiente de aprobación en el PR) | 5/5 | 5/5 | ✓ | RED → GREEN ×3, REFACTOR ×2 | ✓ 6/6 escenarios |
| HU-05 Crear un elemento de backlog | ✓ (Borrador, pendiente de aprobación en el PR) | 10/10 | 10/10 | ✓ | RED → GREEN ×2, REFACTOR ✓ | ✓ 19/19 escenarios |
| HU-06 Editar o eliminar un elemento del backlog | ✓ (Borrador, pendiente de aprobación en el PR; RN-06.8 pendiente de HU-09) | 10/10 | 10/10 | ✓ | RED → GREEN ×6, REFACTOR ✓ | ✓ 22/22 escenarios |
| HU-07 Listar y ordenar el backlog | ✓ (Borrador, pendiente de aprobación en el PR) | 8/8 | 8/8 | ✓ | RED → GREEN ×6, REFACTOR ✓ | ✓ 14/14 escenarios |
| HU-08 Crear un Sprint con Sprint Goal | ✓ (Borrador, pendiente de aprobación en el PR) | 9/9 | 9/9 | ✓ | RED → GREEN ×5, REFACTOR ✓ | ✓ 28/28 escenarios |
| HU-18 Registrar esfuerzo en una historia o tarea | ✓ (Borrador; P-6, la ruta JSON, pendiente de confirmar; P-14 y P-15, de la pantalla, abiertas) | 14/25 (CA-18.15 se cubre con tests del handler; CA-18.16 a CA-18.25, con tests del adaptador) | 24/25 (CA-18.11 solo en BDD) | ✓ | RED → GREEN ×16, REFACTOR ×5 | ✓ 28/28 escenarios |
| HU-30 Interfaz web base | ✓ (Borrador, pendiente de aprobación en el PR) | 10/10 | 10/10 | ✓ | RED → GREEN ×3, REFACTOR ✓ | ✓ 17/17 escenarios |
| HU-09 Asignar historias del backlog a un Sprint | ✓ (Borrador, pendiente de aprobación en el PR) | 8/8 | 8/8 | ✓ | RED → GREEN ×8, REFACTOR ✓ | ✓ 8/8 escenarios |
| HU-10 Marcar una historia como completada en el Sprint | ✓ (Aprobada; P-9 abierta a objeción en el PR) | 5/7 (CA-10.6 y CA-10.7 se cubren con tests del adaptador) | 7/7 | ✓ | RED → GREEN ×6, REFACTOR ×5 | ✓ 9/9 escenarios |
| HU-11 Cerrar un Sprint | ✓ (Aprobada) | 5/7 (CA-11.4 y CA-11.5 se cubren con tests del adaptador) | 7/7 | ✓ | RED → GREEN ×6, REFACTOR ×4 | ✓ 17/17 escenarios |
| HU-20 Registrar un defecto | ✓ (Aprobada) | 6/8 (CA-20.6 y CA-20.8 se cubren con tests del handler y de la pantalla) | 8/8 | ✓ | RED → GREEN ×6 | ✓ 18/18 escenarios |
| HU-12 Consultar Sprints anteriores | ✓ (Implementada) | 5/7 (CA-12.6 y CA-12.7 se cubren con tests del adaptador) | 7/7 | ✓ | RED → GREEN ×4, REFACTOR ✓ | ✓ 7/7 escenarios |
| HU-13 Estimar una historia con story points de la escala Fibonacci | ✓ (Borrador, pendiente de aprobación en el PR) | 6/7 (CA-13.7 se cubre con tests de la pantalla) | 7/7 | ✓ | RED → GREEN ×2, REFACTOR ✓ | ✓ 19/19 escenarios |
| HU-14 Registrar estimaciones individuales ocultas | ✓ (Borrador, pendiente de aprobación en el PR) | 9/9 | 9/9 | ✓ | RED → GREEN ×5, REFACTOR ✓ | ✓ 20/20 escenarios |
| HU-15 Mostrar las estimaciones y detectar diferencias | ✓ (Borrador, pendiente de aprobación en el PR) | 6/6 | 6/6 | ✓ | RED → GREEN ×3 | ✓ 10/10 escenarios |
| HU-22 Calcular story points y velocidad del equipo | ✓ (Aprobada) | 6/9 (CA-22.1 ya la cubre HU-12; CA-22.7 y CA-22.8 se cubren con tests del handler y de la pantalla) | 9/9 | ✓ | RED → GREEN ×6 | ✓ 9/9 escenarios |
| HU-16 Repetir una ronda de votación | ✓ (Implementada) | 7/8 (CA-16.7 se cubre con tests del handler) | 8/8 | ✓ | RED → GREEN ×5 | ✓ 8/8 escenarios |
| HU-17 Registrar la estimación acordada | ✓ (Implementada) | 9/10 (CA-17.9 se cubre con tests del handler) | 10/10 | ✓ | RED → GREEN ×5 | ✓ 17/17 escenarios |
| HU-19 Comparar esfuerzo estimado contra real | ✓ (Borrador; P-10 y P-12 confirmadas por el equipo, PR #297) | 15/19 (CA-19.16 se cubre con tests de persistencia; CA-19.17, del handler; CA-19.18 y CA-19.19, de la pantalla) | 19/19 | ✓ | RED → GREEN ×11, REFACTOR ×2 | ✓ 34/34 escenarios (355/355 en el BDD completo; pasos: 1755/1755) |
| HU-21 Actualizar el estado de un defecto | ✓ (Aprobada) | 8/11 (CA-21.9 a CA-21.11 se cubren con tests del adaptador y de la pantalla) | 11/11 | ✓ | RED → GREEN ×6 | ✓ 14/14 escenarios |
| HU-23 Horas estimadas, reales y desviación de un sprint o proyecto | ✓ (Implementada; decisiones confirmadas en el PR #298) | 15/18 (CA-23.16 se cubre con tests del handler; CA-23.17 y CA-23.18, de la pantalla) | 17/18 (CA-23.11 solo en BDD) | ✓ | RED → GREEN ×6, REFACTOR ×5 | ✓ 29/29 escenarios (407/407 en el BDD completo; pasos: 2130/2130) |
| HU-24 Calcular % de story points completados y defectos por Sprint | ✓ (Aprobada; P-1 a P-3 confirmadas por Matías López en la revisión del PR #366) | 7/9 (CA-24.8 se cubre con tests del handler; CA-24.9, de la pantalla) | 9/9 | ✓ | RED → GREEN ×7 | ✓ 9/9 escenarios (416/416 en el BDD completo; pasos: 2172/2172) |
| HU-27 Exportar el reporte a PDF | ✓ (Borrador; que la conversión sea de HU-27 lo confirmó Matías López en #38 y #378; P-1 a P-11, abiertas a objeción en el PR; ADR-0006 en `proposed`) | 15/26 (CA-27.3 a CA-27.6, CA-27.12, CA-27.19, CA-27.24, CA-27.25 y CA-27.26 se cubren con tests del adaptador; CA-27.21 con tests de la pantalla y revisión manual; CA-27.23 con un test de integración) | 26/26 | ✓ | RED → GREEN ×14, REFACTOR ×2 | ✓ 38/38 escenarios (252 pasos; BDD completo en verde en 500 s) |
| HU-29 Documentación técnica del proyecto | — (documentación, sin funcionalidad nueva: ver la sección de HU-29) | — | — | — (sin código Go) | — (sin TDD: no hay lógica) | — (3/3 criterios verificados en la revisión) |
| HU-25 Ver panel con estado del proyecto y gráficos | ✓ (Aprobada; ADR-0005, P-1 y P-2 aprobadas por Matías López en la revisión del PR #375) | 8/11 (CA-25.8 se cubre con tests del handler; CA-25.9 y CA-25.10, de la pantalla) | 11/11 | ✓ | RED → GREEN ×7, REFACTOR ✓ | ✓ 11/11 escenarios |

---

## HU-01: Crear un proyecto

- **Issue:** [#12](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/12) · Sprint 1 · 3 SP
- **Spec:** [docs/specs/HU-01-crear-proyecto.md](specs/HU-01-crear-proyecto.md)
- **Feature:** [features/HU-01-crear-proyecto.feature](../features/HU-01-crear-proyecto.feature) (steps en [test/bdd/hu01_steps_test.go](../test/bdd/hu01_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-01.1 Crear con datos válidos → Activo con ID | RN-01.6, RN-01.8 | Crear un proyecto con datos válidos `@normal` | `TestNuevoProyecto_DatosValidos_CreaProyectoActivo`, `TestCrearProyecto_DatosValidos_GuardaYDevuelveConID`, `TestPostProyectos_DatosValidos_Devuelve201ConProyecto` | `domain.NuevoProyecto`, `app.CrearProyecto` |
| CA-01.2 Nombre vacío o solo espacios | RN-01.1 | No se puede crear un proyecto sin nombre `@error` · Un nombre de solo espacios cuenta como vacío `@error` | `TestNuevoProyecto_DatosInvalidos_RetornaError` (nombre vacío / solo espacios), `TestCrearProyecto_DatosInvalidos_NoGuarda`, `TestPostProyectos_DatosInvalidos_Devuelve400` | `domain.NuevoProyecto` |
| CA-01.3 Máximo 100 caracteres | RN-01.2 | Un nombre de exactamente 100 caracteres es válido `@limite` · Un nombre de más de 100 caracteres es rechazado `@error` | `TestNuevoProyecto_DatosInvalidos_RetornaError` (101), `TestNuevoProyecto_CasosLimite_SonValidos` (100) | `domain.NuevoProyecto` |
| CA-01.4 Nombre duplicado (sin distinguir mayúsculas) | RN-01.3 | No se puede repetir el nombre de un proyecto `@error` (Alfa / alfa / ALFA) | `TestClaveNombre_IgnoraMayusculasYEspacios`, `TestCrearProyecto_NombreDuplicado_RetornaErrorYNoGuarda`, `TestRepositorioProyectos_NombreDuplicado`, `TestPostProyectos_DatosInvalidos_Devuelve400` | `domain.ClaveNombre`, `app.CrearProyecto`, `persistencia.RepositorioProyectos` (UNIQUE `nombre_clave`) |
| CA-01.5 Fecha de inicio obligatoria | RN-01.4 | La fecha de inicio es obligatoria `@error` | `TestNuevoProyecto_DatosInvalidos_RetornaError`, `TestPostProyectos_DatosInvalidos_Devuelve400` | `domain.NuevoProyecto` |
| CA-01.6 Fin anterior a inicio | RN-01.5 | La fecha de fin no puede ser anterior a la de inicio `@error` | `TestNuevoProyecto_DatosInvalidos_RetornaError`, `TestPostProyectos_DatosInvalidos_Devuelve400` | `domain.NuevoProyecto` |
| CA-01.7 Fin igual a inicio o sin fin | RN-01.5 | Crear un proyecto sin fecha de fin `@alternativo` · La fecha de fin puede ser igual a la de inicio `@limite` | `TestNuevoProyecto_CasosLimite_SonValidos`, `TestPostProyectos_SinFechaFin_DevuelveFechaFinNula` | `domain.NuevoProyecto` |
| CA-01.8 Nombre recortado | RN-01.1 | El nombre se guarda sin los espacios de alrededor `@alternativo` | `TestNuevoProyecto_CasosLimite_SonValidos` | `domain.NuevoProyecto` |
| CA-01.9 IDs distintos y recuperables | RN-01.8 | Cada proyecto recibe un identificador distinto `@normal` | `TestCrearProyecto_DosProyectos_TienenIDsDistintos`, `TestRepositorioProyectos_GuardarYBuscar_DevuelveLoMismo`, `TestRepositorioProyectos_IDsAutoincrementalesDistintos` | `persistencia.RepositorioProyectos.Guardar` / `BuscarPorID` |

RN-01.7 (los 5 estados posibles) no tiene CA propio en esta historia; lo verifica `TestEstados_SonLosDefinidosEnLaSpec`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/proyecto.go](../internal/domain/proyecto.go) | [proyecto_test.go](../internal/domain/proyecto_test.go) |
| Caso de uso | [internal/app/crear_proyecto.go](../internal/app/crear_proyecto.go) | [crear_proyecto_test.go](../internal/app/crear_proyecto_test.go) |
| Persistencia | [internal/adapters/persistencia/proyecto_repo.go](../internal/adapters/persistencia/proyecto_repo.go), [db.go](../internal/adapters/persistencia/db.go) | [proyecto_repo_test.go](../internal/adapters/persistencia/proyecto_repo_test.go) |
| HTTP | [internal/adapters/http/proyecto_handler.go](../internal/adapters/http/proyecto_handler.go) (`POST /proyectos`) | [proyecto_handler_test.go](../internal/adapters/http/proyecto_handler_test.go) |
| BDD | [test/bdd/hu01_sistema_test.go](../test/bdd/hu01_sistema_test.go) | 14 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 1.1 Spec | `e2cf2d9` docs(HU-01): [1.1] spec de crear proyecto | — |
| 1.2 Escenarios BDD | `1025c3a` test(HU-01): [1.2] escenarios BDD | — |
| 1.3 Steps godog | `2f7f05c` test(HU-01): [1.3] steps de godog (en rojo) | BDD en rojo |
| 1.4 Entidad Proyecto | `c7cfa05` test(HU-01): [1.4][RED] | **RED** |
| | `5b48a7f` feat(HU-01): [1.4][GREEN] | **GREEN** |
| 1.5 Caso de uso | `40b564e` test(HU-01): [1.5][RED] | **RED** |
| | `e03080a` feat(HU-01): [1.5][GREEN] | **GREEN** |
| 1.6 Persistencia | `18ded3f` feat(HU-01): [1.6] repositorio SQLite | adaptador + test de integración |
| 1.7 HTTP | `ccd6c54` feat(HU-01): [1.7] POST /proyectos y BDD en verde | BDD en verde |
| 1.8 Revisión | `f88db54` refactor(HU-01): [1.8] correcciones go-reviewer | **REFACTOR** |
| 1.9 Trazabilidad | este commit | — |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.
- Los adaptadores (1.6 y 1.7) se commitearon junto con sus tests, sin un commit RED aparte. La cátedra exige TDD para reglas de negocio, cálculos y validaciones, y esos están en dominio y caso de uso (1.4 y 1.5), que sí tienen RED → GREEN.

---

## HU-02: Modificar un proyecto

- **Issue:** [#13](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/13) · Sprint 1 · 2 SP
- **Spec:** [docs/specs/HU-02-modificar-proyecto.md](specs/HU-02-modificar-proyecto.md)
- **Feature:** [features/HU-02-modificar-proyecto.feature](../features/HU-02-modificar-proyecto.feature) (steps en [test/bdd/hu02_steps_test.go](../test/bdd/hu02_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-02.1 Cambiar solo el nombre | RN-02.1, RN-02.2 | Cambiar solo el nombre no toca el resto de los datos `@normal` · El nuevo nombre se guarda sin los espacios de alrededor `@alternativo` | `TestProyecto_Modificar_SoloNombre_ConservaElResto`, `TestModificarProyecto_Valido_ActualizaYDevuelveElProyecto`, `TestPatchProyecto_SoloNombre_Devuelve200ConElProyecto` | `domain.Proyecto.Modificar`, `app.ModificarProyecto`, `http.handlerModificarProyecto` |
| CA-02.2 Nombre vacío o largo | RN-02.2 | Un nombre vacío es rechazado `@error` · Un nombre de solo espacios cuenta como vacío `@error` · Un nombre de exactamente 100 caracteres es válido `@limite` · Un nombre de más de 100 caracteres es rechazado `@error` | `TestProyecto_Modificar_DatosInvalidos_RetornaError`, `TestProyecto_Modificar_CasosValidos`, `TestPatchProyecto_DatosInvalidos_Devuelve400` (nombre vacío) | `domain.Proyecto.Modificar` (reutiliza `domain.NuevoProyecto`) |
| CA-02.3 Nombre de otro proyecto | RN-02.3 | No se puede usar el nombre de otro proyecto `@error` (3 ejemplos) · Un proyecto puede cambiar solo las mayúsculas de su propio nombre `@alternativo` | `TestModificarProyecto_Invalido_RetornaErrorYNoActualiza` (nombre de otro proyecto), `TestModificarProyecto_MismoNombreConOtrasMayusculas_Actualiza`, `TestRepositorioProyectos_Actualizar_NombreDeOtro_RetornaDuplicado`, `TestRepositorioProyectos_Actualizar_MismoNombreConOtrasMayusculas`, `TestPatchProyecto_DatosInvalidos_Devuelve400` (nombre duplicado) | `app.ModificarProyecto`, `persistencia.RepositorioProyectos.Actualizar` |
| CA-02.4 Fechas en orden | RN-02.5 | La nueva fecha de fin no puede ser anterior a la de inicio guardada `@error` · La nueva fecha de inicio no puede ser posterior a la de fin guardada `@error` · La nueva fecha de fin puede ser igual a la de inicio `@limite` | `TestProyecto_Modificar_DatosInvalidos_RetornaError`, `TestProyecto_Modificar_CasosValidos`, `TestPatchProyecto_DatosInvalidos_Devuelve400` (fecha de fin anterior) | `domain.Proyecto.Modificar` |
| CA-02.5 Quitar la fecha de fin | RN-02.5 | Se puede quitar la fecha de fin `@alternativo` | `TestProyecto_Modificar_CasosValidos`, `TestRepositorioProyectos_Actualizar_Persiste`, `TestPatchProyecto_TraduceLosCampos` (fecha_fin en null) | `domain.Proyecto.Modificar`, `http.pedidoModificarProyecto.aCambios` |
| CA-02.6 Cambiar el estado | RN-02.6, RN-02.7 | Finalizar un proyecto activo `@normal` · Reactivar un proyecto cancelado `@alternativo` · Enviar el mismo estado que ya tiene no es un error `@limite` | `TestProyecto_Modificar_CasosValidos`, `TestProyecto_Modificar_CualquierTransicionDeEstado`, `TestModificarProyecto_SinNombre_NoConsultaDuplicados`, `TestPatchProyecto_TraduceLosCampos` (estado) | `domain.Proyecto.Modificar`, `domain.EstadoValido` |
| CA-02.7 Estado inexistente | RN-02.6 | Un estado que no existe es rechazado `@error` (2 ejemplos) | `TestProyecto_Modificar_DatosInvalidos_RetornaError`, `TestEstadoValido`, `TestModificarProyecto_Invalido_RetornaErrorYNoActualiza` (estado inválido), `TestPatchProyecto_DatosInvalidos_Devuelve400` (estado inválido) | `domain.EstadoValido` |
| CA-02.8 Sin cambios | RN-02.8 | Hay que indicar al menos un dato para modificar `@error` | `TestProyecto_Modificar_SinCambios_RetornaError`, `TestModificarProyecto_Invalido_RetornaErrorYNoActualiza` (sin cambios), `TestPatchProyecto_DatosInvalidos_Devuelve400` (sin cambios) | `domain.Proyecto.Modificar` |
| CA-02.9 Proyecto inexistente | RN-02.9 | No se puede modificar un proyecto que no existe `@error` | `TestModificarProyecto_Invalido_RetornaErrorYNoActualiza` (proyecto inexistente), `TestRepositorioProyectos_Actualizar_Inexistente_RetornaNoEncontrado`, `TestPatchProyecto_Inexistente_Devuelve404` | `app.ModificarProyecto`, `persistencia.RepositorioProyectos.Actualizar` |
| CA-02.10 ID y persistencia | RN-02.9 | La modificación queda guardada y el identificador no cambia `@normal` | `TestProyecto_Modificar_SoloNombre_ConservaElResto`, `TestRepositorioProyectos_Actualizar_Persiste` | `domain.Proyecto.Modificar`, `persistencia.RepositorioProyectos.Actualizar` |
| CA-02.11 Quitar la fecha de inicio | RN-02.4 | No se puede quitar la fecha de inicio `@error` | `TestProyecto_Modificar_DatosInvalidos_RetornaError`, `TestPatchProyecto_DatosInvalidos_Devuelve400` (quitar fecha de inicio) | `domain.Proyecto.Modificar` |
| CA-02.12 Conservar integrantes | RN-02.10 | Cambiar el nombre conserva los integrantes y sus roles `@normal` | `TestProyecto_Modificar_ConservaLosIntegrantes`, `TestRepositorioProyectos_Actualizar_ConservaLosIntegrantes` | `domain.Proyecto.Modificar`, `persistencia.RepositorioProyectos.Actualizar` (solo cambia la tabla `proyectos`) |

Pedidos simultáneos (sección 5 de la spec; sin CA propio porque no se puede expresar en Gherkin): [modificar_concurrencia_test.go](../internal/adapters/persistencia/modificar_concurrencia_test.go) (`TestModificar_CamposDistintosSimultaneos_SeConservanAmbos`, `TestModificar_FechasCruzadasSimultaneas_SegundoEsRechazado`), `TestRepositorioProyectos_Actualizar_SoloGuardaLosCamposEnviados`, `TestRepositorioProyectos_Actualizar_FinAnteriorAlInicioGuardado_RetornaError`.

Orden de los errores (sección 7 de la spec): `TestModificarProyecto_Invalido_RetornaErrorYNoActualiza` (nombre duplicado y fecha de fin anterior) verifica que los errores del dominio se informan antes que el duplicado. Que el original no cambie ante un error: `TestProyecto_Modificar_NoAlteraElOriginal` y los escenarios con "el proyecto no cambió".

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/proyecto_modificar.go](../internal/domain/proyecto_modificar.go) (`CambiosProyecto`, `Modificar`, `EstadoValido`) | [proyecto_modificar_test.go](../internal/domain/proyecto_modificar_test.go) |
| Caso de uso | [internal/app/modificar_proyecto.go](../internal/app/modificar_proyecto.go) | [modificar_proyecto_test.go](../internal/app/modificar_proyecto_test.go) |
| Persistencia | [internal/adapters/persistencia/proyecto_repo.go](../internal/adapters/persistencia/proyecto_repo.go) (`Actualizar`: solo los campos enviados, con RN-02.5 verificada al guardar) | [modificar_repo_test.go](../internal/adapters/persistencia/modificar_repo_test.go), [modificar_concurrencia_test.go](../internal/adapters/persistencia/modificar_concurrencia_test.go) |
| HTTP | [internal/adapters/http/modificar_proyecto_handler.go](../internal/adapters/http/modificar_proyecto_handler.go) (`PATCH /proyectos/{id}`) | [modificar_proyecto_handler_test.go](../internal/adapters/http/modificar_proyecto_handler_test.go) |
| BDD | [test/bdd/hu02_sistema_test.go](../test/bdd/hu02_sistema_test.go), [hu02_steps_test.go](../test/bdd/hu02_steps_test.go) | 24 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 2.1 Spec | `c7041fd` docs(HU-02): [2.1] spec de modificar proyecto | — |
| | `d7e5e90` docs(HU-02): [2.1] agregar CA-02.11 | — |
| 2.2 Escenarios BDD | `4f2f060` test(HU-02): [2.2] escenarios BDD | — |
| 2.3 Steps godog | `1f2f41a` test(HU-02): [2.3] steps de godog (en rojo) | BDD en rojo |
| 2.4 Dominio | `74a59be` test(HU-02): [2.4][RED] Modificar en Proyecto | **RED** |
| | `ab31b30` feat(HU-02): [2.4][GREEN] Modificar reutilizando las validaciones de HU-01 | **GREEN** |
| | `b40892e` test(HU-02): [2.4] mapear ErrSinCambios y ErrEstadoInvalido en BDD | — |
| | `1c88932` docs(HU-02): RN-02.10 y CA-02.12 | — |
| | `707b448` test(HU-02): [RED] modificar conserva los integrantes | **RED** |
| | `aeda9f0` feat(HU-02): [GREEN] modificar conserva los integrantes | **GREEN** |
| 2.5 Caso de uso | `d253400` docs(HU-02): orden de los errores | — |
| | `4061ca5` test(HU-02): [2.5][RED] caso de uso ModificarProyecto | **RED** |
| | `0c10a6e` feat(HU-02): [2.5][GREEN] caso de uso ModificarProyecto | **GREEN** |
| 2.6 Persistencia | `d225791` test(HU-02): [2.6][RED] persistir la modificación | **RED** |
| | `2eadb9d` feat(HU-02): [2.6][GREEN] persistir la modificación | **GREEN** |
| | `9d2119e` test(HU-02): escenario de CA-02.12 y BDD en verde | BDD en verde |
| 2.7 HTTP | `d7fcf47` test(HU-02): [2.7][RED] endpoint PATCH /proyectos/{id} | **RED** |
| | `8eada14` feat(HU-02): [2.7][GREEN] endpoint PATCH /proyectos/{id} | **GREEN** |
| 2.8 Revisión + trazabilidad | `1cf6cf2` docs(HU-02): [2.8] trazabilidad | — |
| Revisión: pedidos simultáneos | `8f16f49` docs(HU-02): guardar solo los campos enviados | — |
| | `79afcc1` test(HU-02): [RED] modificaciones simultáneas | **RED** |
| | `10d03a1` feat(HU-02): [GREEN] guardar solo los campos enviados y verificar las fechas al guardar | **GREEN** |
| | `f310577` test(HU-02): [RED] devolver el proyecto como quedó guardado | **RED** |
| | `fc18ad4` feat(HU-02): [GREEN] devolver el proyecto como quedó guardado | **GREEN** |
| | este commit | trazabilidad actualizada |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.
- El endpoint es `PATCH /proyectos/{id}` y no `PUT` como decía la Tarea 2.7: la spec define una edición parcial, y en HTTP eso es `PATCH`.

---

## HU-03: Registrar integrantes de un proyecto

- **Issue:** #14 · Sprint 1 · 3 SP
- **Spec:** [docs/specs/HU-03-registrar-integrantes.md](specs/HU-03-registrar-integrantes.md)
- **Feature:** [features/HU-03-registrar-integrantes.feature](../features/HU-03-registrar-integrantes.feature) (steps en [test/bdd/hu03_steps_test.go](../test/bdd/hu03_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-03.1 Agregar integrante válido | RN-03.2 | Agregar un integrante válido `@normal` · Agregar varios integrantes distintos al mismo proyecto `@alternativo` | `TestProyecto_AgregarIntegrante_Valido`, `TestAgregarIntegrante_DatosValidos_GuardaEnElProyecto`, `TestPostIntegrantes_DatosValidos_Devuelve204`, `TestRepositorioProyectos_AgregarIntegrante_Persiste` | `domain.Proyecto.AgregarIntegrante`, `app.AgregarIntegrante`, `http.handlerAgregarIntegrante` |
| CA-03.2 Agregar sin nombre | RN-03.1 | No se puede agregar un integrante sin nombre `@error` | `TestProyecto_AgregarIntegrante_Invalido_RetornaError` (nombre vacío / solo espacios), `TestAgregarIntegrante_DatosInvalidos_NoGuarda`, `TestPostIntegrantes_DatosInvalidos_Devuelve400` (nombre vacío) | `domain.Proyecto.AgregarIntegrante` |
| CA-03.3 Quitar integrante existente | — (acción principal) | Quitar un integrante existente quedando otros en el proyecto `@normal` · Quitar al único integrante restante deja la lista vacía `@limite` · Quitar a un integrante con varios roles le quita todos `@limite` | `TestProyecto_QuitarIntegrante_Existente`, `TestProyecto_QuitarIntegrante_UnicoRestante_DejaListaVacia`, `TestQuitarIntegrante_Existente_GuardaSinEl`, `TestDeleteIntegrantes_Existente_Devuelve204`, `TestRepositorioProyectos_QuitarIntegrante_Persiste`, `TestRepositorioProyectos_QuitarIntegrante_UnicoRestante_QuedaVacio`, `TestRepositorioProyectos_QuitarIntegrante_BorraSusRoles` | `domain.Proyecto.QuitarIntegrante`, `app.QuitarIntegrante`, `http.handlerQuitarIntegrante` |
| CA-03.4 Quitar integrante que no pertenece | — (error de dominio) | No se puede quitar un integrante que no pertenece al proyecto `@error` | `TestProyecto_QuitarIntegrante_Inexistente_RetornaError`, `TestQuitarIntegrante_Inexistente_RetornaErrorYNoGuarda`, `TestDeleteIntegrantes_Inexistente_Devuelve404` | `domain.Proyecto.QuitarIntegrante` |
| CA-03.5 Agregar nombre duplicado | RN-03.3 | No se puede agregar un integrante con nombre duplicado `@error` | `TestProyecto_AgregarIntegrante_Invalido_RetornaError` (duplicado exacto / con mayúsculas y espacios), `TestRepositorioProyectos_DuplicadoTrasRecargar_RetornaError`, `TestPostIntegrantes_DatosInvalidos_Devuelve400` (nombre duplicado) | `domain.Proyecto.AgregarIntegrante`, `domain.ClaveNombre` |
| CA-03.6 Agregar sin rol | RN-03.4 | No se puede agregar un integrante sin rol `@error` | `TestProyecto_AgregarIntegrante_Invalido_RetornaError` (rol vacío), `TestPostIntegrantes_DatosInvalidos_Devuelve400` (rol vacío) | `domain.Proyecto.AgregarIntegrante` |
| CA-03.7 Sumar un rol | RN-03.5 | Sumarle un rol a un integrante que ya está en el proyecto `@normal` | `TestProyecto_AgregarRol_SumaElRolSinDuplicarAlIntegrante`, `TestAgregarRol_Valido_GuardaElRolNormalizado`, `TestRepositorioProyectos_AgregarRol_Persiste`, `TestPostRoles_DatosValidos_Devuelve204` | `domain.Proyecto.AgregarRol`, `app.AgregarRol`, `persistencia.RepositorioProyectos.AgregarRol`, `http.handlerAgregarRol` |
| CA-03.8 Rol repetido | RN-03.6 | No se puede repetir un rol en el mismo integrante `@error` | `TestProyecto_AgregarRol_Invalido_RetornaErrorYNoCambia` (duplicado), `TestAgregarRol_Invalido_RetornaErrorYNoGuarda` (duplicado), `TestRepositorioProyectos_AgregarRol_Duplicado_RetornaError`, `TestPostRoles_DatosInvalidos_Devuelve400` (duplicado) | `domain.Proyecto.AgregarRol`, `domain.ClaveNombre` |
| CA-03.9 Rol a quien no es integrante | — (error de dominio) | No se puede agregar un rol a quien no es integrante del proyecto `@error` | `TestProyecto_AgregarRol_Invalido_RetornaErrorYNoCambia` (inexistente), `TestAgregarRol_Invalido_RetornaErrorYNoGuarda` (inexistente), `TestRepositorioProyectos_AgregarRol_IntegranteInexistente_RetornaError`, `TestPostRoles_Inexistente_Devuelve404` | `domain.Proyecto.AgregarRol` |
| CA-03.10 Rol vacío | RN-03.4 | No se puede agregar un rol vacío `@error` | `TestProyecto_AgregarRol_Invalido_RetornaErrorYNoCambia` (vacío / solo espacios), `TestAgregarRol_Invalido_RetornaErrorYNoGuarda` (vacío), `TestPostRoles_DatosInvalidos_Devuelve400` (rol vacío) | `domain.Proyecto.AgregarRol` |

Proyecto inexistente (sección 7 de la spec, sin CA propio): `TestAgregarIntegrante_ProyectoInexistente_RetornaError`, `TestQuitarIntegrante_ProyectoInexistente_RetornaError`, `TestAgregarRol_Invalido_RetornaErrorYNoGuarda` (proyecto inexistente), `TestPostIntegrantes_ProyectoInexistente_Devuelve404`, y los casos de proyecto inexistente de `TestDeleteIntegrantes_Inexistente_Devuelve404` y `TestPostRoles_Inexistente_Devuelve404`. Reutiliza `domain.ErrProyectoNoEncontrado`, ya definido en HU-01.

Pedidos simultáneos (sección 5 de la spec; sin CA propio porque no se puede expresar en Gherkin): [integrantes_concurrencia_test.go](../internal/adapters/persistencia/integrantes_concurrencia_test.go) (`TestIntegrantes_DosAltasSimultaneas_QuedanAmbas`, `TestIntegrantes_BajaYAltaSimultaneas_NoReviveElQuitado`, `TestIntegrantes_MismaAltaSimultanea_SegundaEsDuplicado`, `TestIntegrantes_MismaBajaSimultanea_SegundaEsNoEncontrado`, `TestIntegrantes_DosRolesSimultaneos_QuedanAmbos`, `TestIntegrantes_MismoRolSimultaneo_SegundoEsDuplicado`).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/integrante.go](../internal/domain/integrante.go) (+ campo `Integrantes` en [proyecto.go](../internal/domain/proyecto.go)) | [integrante_test.go](../internal/domain/integrante_test.go) |
| Caso de uso | [internal/app/gestionar_integrantes.go](../internal/app/gestionar_integrantes.go) (`AgregarIntegrante`, `QuitarIntegrante`, `AgregarRol`) | [gestionar_integrantes_test.go](../internal/app/gestionar_integrantes_test.go) |
| Persistencia | [internal/adapters/persistencia/proyecto_repo.go](../internal/adapters/persistencia/proyecto_repo.go) (`AgregarIntegrante`, `QuitarIntegrante`, `AgregarRol`, `buscarIntegrantes`), [db.go](../internal/adapters/persistencia/db.go) (tablas `integrantes` e `integrante_roles`) | [integrantes_repo_test.go](../internal/adapters/persistencia/integrantes_repo_test.go), [integrantes_concurrencia_test.go](../internal/adapters/persistencia/integrantes_concurrencia_test.go) |
| HTTP | [internal/adapters/http/integrante_handler.go](../internal/adapters/http/integrante_handler.go) (`POST`/`DELETE /proyectos/{id}/integrantes`, `POST /proyectos/{id}/integrantes/{nombre}/roles`) | [integrante_handler_test.go](../internal/adapters/http/integrante_handler_test.go) |
| BDD | [test/bdd/hu03_sistema_test.go](../test/bdd/hu03_sistema_test.go), [hu03_steps_test.go](../test/bdd/hu03_steps_test.go) | 13 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 3.1 Spec | `5a35d90` docs(HU-03): [3.1] spec de registrar integrantes | — |
| 3.2 Escenarios BDD | `d683854` test(HU-03): [3.2] escenarios BDD | — |
| 3.3 Steps godog | `f834d1c` test(HU-03): [3.3] steps de godog (en rojo) | BDD en rojo |
| 3.4 Entidad Integrante | `94880f1` test(HU-03): [3.4][RED] | **RED** |
| | `dd433e7` feat(HU-03): [3.4][GREEN] | **GREEN** |
| 3.5 Casos de uso | `4181167` test(HU-03): [3.5][RED] | **RED** |
| | `73ebc0e` feat(HU-03): [3.5][GREEN] | **GREEN** |
| 3.6 Persistencia | `e0f8f0d` test(HU-03): [3.6][RED] | **RED** |
| | `e0bc8ed` feat(HU-03): [3.6][GREEN] | **GREEN** |
| 3.7 HTTP | `de67103` test(HU-03): [3.7][RED] | **RED** |
| | `b4df091` feat(HU-03): [3.7][GREEN] adaptador HTTP y BDD en verde | **GREEN** |
| 3.8 Revisión + trazabilidad | `ccf3ba6` | revisión go-reviewer |
| Revisión del PR: spec | `447b0a6` docs(HU-03): issue #14, roles múltiples y concurrencia | — |
| Revisión del PR: concurrencia | `af22850` test(HU-03): [RED] altas y bajas simultáneas se pisan | **RED** |
| | `2b7f78b` fix(HU-03): [GREEN] guardar integrantes fila por fila | **GREEN** |
| | `18a50db` test(HU-03): [RED] la base informa duplicado y no encontrado | **RED** |
| | `7479890` fix(HU-03): [GREEN] traducir UNIQUE y borrado vacío a errores del dominio | **GREEN** |
| Revisión del PR: roles múltiples | `a807b3f` test(HU-03): escenarios BDD y steps de agregar rol | BDD en rojo |
| | `ae0cc1f` test(HU-03): [RED] integrante con varios roles | **RED** |
| | `80612b8` feat(HU-03): [GREEN] integrante con varios roles | **GREEN** |
| | `04bf049` test(HU-03): [RED] persistir los roles | **RED** |
| | `a2b2878` feat(HU-03): [GREEN] persistir los roles (tabla `integrante_roles`) | **GREEN** |
| | `6090836` test(HU-03): [RED] caso de uso agregar rol | **RED** |
| | `f85f2ca` feat(HU-03): [GREEN] caso de uso agregar rol | **GREEN** |
| | `398316e` test(HU-03): [RED] endpoint de agregar rol | **RED** |
| | `5325df6` feat(HU-03): [GREEN] endpoint de agregar rol | **GREEN** |
| | este commit | BDD en verde y trazabilidad |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.

---

## HU-04: Consultar el estado de un proyecto

- **Issue:** #15 · Sprint 1 · 3 SP
- **Spec:** [docs/specs/HU-04-consultar-proyecto.md](specs/HU-04-consultar-proyecto.md)
- **Feature:** [features/HU-04-consultar-proyecto.feature](../features/HU-04-consultar-proyecto.feature) (steps en [test/bdd/hu04_steps_test.go](../test/bdd/hu04_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-04.1 Datos del proyecto (nombre, fechas, estado) | — (sección 3) | Consultar un proyecto existente `@normal` | `TestConsultarProyecto_Existente_DevuelveDatosEIntegrantes`, `TestGetProyecto_Existente_Devuelve200ConElJSONDeLaSpec` | `app.ConsultarProyecto.Ejecutar`, `http.handlerConsultarProyecto` |
| CA-04.2 Integrantes con sus roles, cada uno una sola vez, en orden de alta | RN-04.1, RN-04.2 | La consulta incluye a los integrantes con sus roles, en orden de alta `@normal` · Un integrante con dos roles aparece una sola vez, con los roles en orden de alta `@limite` | `TestConsultarProyecto_Existente_DevuelveDatosEIntegrantes`, `TestGetProyecto_Existente_Devuelve200ConElJSONDeLaSpec` | `app.ConsultarProyecto.Ejecutar`, `http.aRespuestaConsultaProyecto`; el orden, `persistencia.RepositorioProyectos.buscarIntegrantes` (HU-03) |
| CA-04.3 Sin integrantes → lista vacía, `[]` en JSON | RN-04.1 | Un proyecto sin integrantes devuelve una lista vacía `@limite` | `TestGetProyecto_SinIntegrantesNiFechaFin_DevuelveArrayVacioYFechaFinNull` | `http.aRespuestaConsultaProyecto` |
| CA-04.4 Sin fecha de fin → `nil`, `null` en JSON | — (sección 3) | Consultar un proyecto sin fecha de fin `@alternativo` | `TestGetProyecto_SinIntegrantesNiFechaFin_DevuelveArrayVacioYFechaFinNull` | `http.aRespuestaProyecto` (HU-01, embebida en `respuestaConsultaProyecto`) |
| CA-04.5 Id inexistente | — (sección 7) | No se puede consultar un proyecto que no existe `@error` | `TestConsultarProyecto_Inexistente_RetornaErrProyectoNoEncontrado`, `TestGetProyecto_Inexistente_Devuelve404ConElMensajeDelError` | `app.ConsultarProyecto.Ejecutar`, `http.handlerConsultarProyecto` |

En CA-04.3 y CA-04.4, el escenario BDD verifica el dominio (`len(Integrantes) == 0`, `FechaFin == nil`) y el test HTTP verifica el JSON (`[]`, `null`).

Sin CA propio en la spec:
- Id inválido: `abc`, `0`, `-1` y fuera de rango de `int64` (secciones 6 y 7). `TestGetProyecto_IDInvalido_Devuelve400SinLlamarAlCasoDeUso`.
- Falla interna (sección 7). `TestConsultarProyecto_FallaDelRepositorio_SeEnvuelveYNoEsNoEncontrado`, `TestGetProyecto_FallaInterna_Devuelve500SinDetalles`.
- `roles` siempre es un array (sección 3). `TestGetProyecto_IntegranteSinRoles_DevuelveRolesComoArrayVacio`.
- Solo lectura (sección 5). No tiene test: la garantiza el compilador, porque el puerto `app.LectorProyectos` tiene un solo método, `BuscarPorID`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | Sin cambios: usa `domain.Proyecto`, `domain.Integrante` y `domain.ErrProyectoNoEncontrado` | — |
| Caso de uso | [internal/app/consultar_proyecto.go](../internal/app/consultar_proyecto.go) (`LectorProyectos`, `ConsultarProyecto`) | [consultar_proyecto_test.go](../internal/app/consultar_proyecto_test.go) |
| Persistencia | Sin cambios: reutiliza `BuscarPorID` y `buscarIntegrantes` de [proyecto_repo.go](../internal/adapters/persistencia/proyecto_repo.go) | El orden de alta se verifica en los escenarios BDD, contra SQLite |
| HTTP | [internal/adapters/http/consultar_proyecto_handler.go](../internal/adapters/http/consultar_proyecto_handler.go) (`GET /proyectos/{id}`, `RegistrarConsultarProyecto`), registrada en [router.go](../internal/adapters/http/router.go) (`NuevoRouter`) | [consultar_proyecto_handler_test.go](../internal/adapters/http/consultar_proyecto_handler_test.go) |
| BDD | [test/bdd/hu04_sistema_test.go](../test/bdd/hu04_sistema_test.go), [hu04_steps_test.go](../test/bdd/hu04_steps_test.go) | 6 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 4.1 Spec | `0c688e3` docs(HU-04): [4.1] spec de consultar proyecto | — |
| 4.2 Escenarios BDD | `3ab1b36` test(HU-04): [4.2] escenarios BDD | — |
| | `68b9b72` test(HU-04): [4.2] given propios de HU-04 | — |
| 4.3 Steps godog | `66b2b81` test(HU-04): [4.3] steps de godog (en rojo) | BDD en rojo |
| 4.4 Caso de uso | `54d7d19` test(HU-04): [4.4][RED] | **RED** |
| | `6093e55` feat(HU-04): [4.4][GREEN] | **GREEN** |
| | `dbbefcd` feat(HU-04): [4.4] conectar el BDD a ConsultarProyecto | BDD en verde |
| 4.5 HTTP | `abf1c49` test(HU-04): [4.5][RED] | **RED** |
| | `4d289e6` feat(HU-04): [4.5][GREEN] | **GREEN** |
| | `e6a5be2` refactor(HU-04): [4.5][REFACTOR] | **REFACTOR** |
| 4.6 Revisión | `fa53df5` test(HU-04): [4.6][RED] roles siempre array | **RED** |
| | `d59ef41` feat(HU-04): [4.6][GREEN] roles siempre array | **GREEN** |
| | `0206c8a` refactor(HU-04): [4.6] correcciones go-reviewer | **REFACTOR** |
| 4.6 Trazabilidad | este commit | — |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.
- Hay tests sin CA. Los de id inválido y falla interna se trazan a las secciones 6 y 7, y el de `roles` siempre array a la sección 3. Si el equipo quiere CA para estos casos, primero hay que agregarlos a la spec.
- Los commits `66b2b81` (steps en rojo) y `dbbefcd` (BDD en verde) no llevan la etiqueta `[RED]`/`[GREEN]`, así que un `git log --grep="\[RED\]"` no los encuentra.
- RN-04.2 (orden de alta) no está en el código de HU-04: la implementa la persistencia de HU-03 (`ORDER BY i.id, ir.id` en `buscarIntegrantes`). HU-04 solo la conserva y la verifica en los escenarios BDD.

---

## HU-05: Crear un elemento de backlog

- **Issue:** [#16](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/16) · sin sprint asignado · 5 SP
- **Spec:** [docs/specs/HU-05-crear-elemento-backlog.md](specs/HU-05-crear-elemento-backlog.md)
- **Feature:** [features/HU-05-crear-elemento-backlog.feature](../features/HU-05-crear-elemento-backlog.feature) (steps en [test/bdd/hu05_steps_test.go](../test/bdd/hu05_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-05.1 Datos válidos → Pendiente con ID | RN-05.8, RN-05.9 | Crear un elemento con todos los datos válidos `@normal` | `TestNuevoElementoBacklog_DatosValidos_CreaPendiente`, `TestCrearElementoBacklog_DatosValidos_GuardaPendienteConID`, `TestPostBacklog_DatosValidos_Devuelve201ConElemento` | `domain.NuevoElementoBacklog`, `app.CrearElementoBacklog` |
| CA-05.2 Sin story points | RN-05.6 | Crear un elemento sin story points `@alternativo` | `TestNuevoElementoBacklog_CasosLimite_SonValidos`, `TestPostBacklog_SinStoryPoints_DevuelveNull`, `TestRepositorioBacklog_GuardarYBuscar_DevuelveLoMismo` | `domain.NuevoElementoBacklog`, `persistencia.RepositorioBacklog` |
| CA-05.3 Proyecto inexistente | RN-05.1 | No se puede crear un elemento en un proyecto que no existe `@error` | `TestCrearElementoBacklog_ProyectoInexistente_RetornaErrorYNoGuarda`, `TestRepositorioBacklog_ProyectoInexistente`, `TestPostBacklog_ProyectoInexistente_Devuelve404` | `app.CrearElementoBacklog`, `persistencia.RepositorioBacklog` |
| CA-05.4 Título obligatorio, máx. 100 | RN-05.2 | El título es obligatorio y tiene un máximo `@error` (×3) · Un título de exactamente 100 caracteres es válido `@limite` | `TestNuevoElementoBacklog_DatosInvalidos_RetornaError`, `TestNuevoElementoBacklog_CasosLimite_SonValidos`, `TestCrearElementoBacklog_DatosInvalidos_NoConsultaTitulosNiGuarda` | `domain.NuevoElementoBacklog` |
| CA-05.5 Descripción obligatoria | RN-05.4 | La descripción es obligatoria `@error` | `TestNuevoElementoBacklog_DatosInvalidos_RetornaError`, `TestPostBacklog_DatosInvalidos_Devuelve400` | `domain.NuevoElementoBacklog` |
| CA-05.6 Prioridad Alta/Media/Baja | RN-05.5 | La prioridad debe ser Alta, Media o Baja `@error` (×3) | `TestNuevoElementoBacklog_DatosInvalidos_RetornaError`, `TestPrioridadesYEstados_SonLosDeLaSpec`, `TestPostBacklog_DatosInvalidos_Devuelve400` | `domain.NuevoElementoBacklog` |
| CA-05.7 Story points positivos | RN-05.6 | Un story point es el mínimo válido `@limite` · Los story points deben ser positivos `@error` (×2) | `TestNuevoElementoBacklog_DatosInvalidos_RetornaError`, `TestNuevoElementoBacklog_CasosLimite_SonValidos` | `domain.NuevoElementoBacklog` |
| CA-05.8 Al menos un criterio, ninguno vacío | RN-05.7 | Debe haber al menos un criterio `@error` · Ningún criterio puede estar vacío `@error` | `TestNuevoElementoBacklog_DatosInvalidos_RetornaError`, `TestPostBacklog_DatosInvalidos_Devuelve400` | `domain.NuevoElementoBacklog` |
| CA-05.9 Título único por proyecto | RN-05.3 | No se repite un título dentro del mismo proyecto `@error` · El mismo título se puede usar en otro proyecto `@alternativo` | `TestCrearElementoBacklog_TituloDuplicado`, `TestRepositorioBacklog_TituloPorProyecto` | `app.CrearElementoBacklog`, `persistencia.RepositorioBacklog` (UNIQUE `proyecto_id, titulo_clave`) |
| CA-05.10 IDs distintos, criterios en orden | RN-05.9 | Cada elemento recibe un identificador distinto y se recupera con sus criterios `@normal` | `TestCrearElementoBacklog_DatosValidos_GuardaPendienteConID`, `TestRepositorioBacklog_GuardarYBuscar_DevuelveLoMismo` | `persistencia.RepositorioBacklog.GuardarElemento` / `BuscarElementoPorID` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog.go](../internal/domain/backlog.go) | [backlog_test.go](../internal/domain/backlog_test.go) |
| Caso de uso | [internal/app/crear_elemento_backlog.go](../internal/app/crear_elemento_backlog.go) | [crear_elemento_backlog_test.go](../internal/app/crear_elemento_backlog_test.go) |
| Persistencia | [internal/adapters/persistencia/backlog_repo.go](../internal/adapters/persistencia/backlog_repo.go), [db.go](../internal/adapters/persistencia/db.go) (tablas `backlog_elementos` y `backlog_criterios`) | [backlog_repo_test.go](../internal/adapters/persistencia/backlog_repo_test.go) |
| HTTP | [internal/adapters/http/backlog_handler.go](../internal/adapters/http/backlog_handler.go) (`POST /proyectos/{id}/backlog`) | [backlog_handler_test.go](../internal/adapters/http/backlog_handler_test.go) |
| BDD | [test/bdd/hu05_sistema_test.go](../test/bdd/hu05_sistema_test.go) | 19 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 5.1 Spec | `584be96` docs(HU-05): [5.1] spec | — |
| 5.2 Escenarios BDD | `9c015e2` test(HU-05): [5.2] escenarios BDD | — |
| 5.3 Steps godog | `b9ba939` test(HU-05): [5.3] steps (en rojo) | BDD en rojo |
| 5.4 Entidad ElementoBacklog | `97e91bc` test(HU-05): [5.4][RED] | **RED** |
| | `62a8b33` feat(HU-05): [5.4][GREEN] | **GREEN** |
| 5.5 Caso de uso | `62c8686` test(HU-05): [5.5][RED] | **RED** |
| | `4882d1a` feat(HU-05): [5.5][GREEN] | **GREEN** |
| 5.6 Persistencia | `c698466` feat(HU-05): [5.6] repositorio SQLite | adaptador + test de integración |
| 5.7 HTTP | `3617ea4` feat(HU-05): [5.7] POST /proyectos/{id}/backlog y BDD en verde | BDD en verde |
| 5.8 Revisión | `3485268` refactor(HU-05): [5.8] correcciones go-reviewer | **REFACTOR** |
| 5.9 Trazabilidad | este commit | — |

### Huecos

- La spec está en **Borrador** y sus decisiones (prioridad Alta/Media/Baja, story points enteros positivos, estados Pendiente/En progreso/Hecho) se validan en el PR.
- La historia no tiene sprint asignado ni sub-tareas en el tablero: se usó la numeración 5.1 a 5.9 por analogía con HU-01.
- Los adaptadores (5.6 y 5.7) se commitearon con sus tests, sin RED aparte; las reglas de negocio (5.4 y 5.5) sí tienen RED → GREEN.

---

## HU-07: Listar y ordenar el backlog por prioridad

- **Issue:** [#18](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/18) · Sprint 2 · 2 SP
- **Spec:** [docs/specs/HU-07-listar-backlog.md](specs/HU-07-listar-backlog.md)
- **Feature:** [features/HU-07-listar-backlog.feature](../features/HU-07-listar-backlog.feature) (steps en [test/bdd/hu07_steps_test.go](../test/bdd/hu07_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-07.1 Alta primero | RN-07.2 | Por defecto el backlog se ordena por prioridad, de Alta a Baja `@normal` | `TestOrdenarBacklog`, `TestListarBacklog_DevuelveLosDelProyectoOrdenados`, `TestGetBacklog_DevuelveLosElementos`, `TestPantallaBacklog_ListaControlesYFormulario` | `domain.OrdenarBacklog`, `app.ListarBacklog`, `http.handlerListarBacklog` |
| CA-07.2 Orden ascendente | RN-07.2 | En orden ascendente el backlog va de Baja a Alta `@alternativo` | `TestOrdenarBacklog` (ascendente), `TestListarBacklog_DevuelveLosDelProyectoOrdenados`, `TestGetBacklog_PasaOrdenYEstado`, `TestPantallaBacklog_OrdenYFiltro` | `domain.OrdenarBacklog` |
| CA-07.3 Desempate por creación | RN-07.3 | A igual prioridad se respeta el orden de creación `@limite` (2 ejemplos) | `TestOrdenarBacklog`, `TestOrdenarBacklog_DesempataPorCreacionAunqueLleguenDesordenados` | `domain.OrdenarBacklog` |
| CA-07.4 Filtro por estado | RN-07.4 | Se puede filtrar el backlog por estado `@normal` (3 ejemplos) · Sin filtro se listan todos los estados `@normal` | `TestOrdenarBacklog` (filtros), `TestListarBacklog_DevuelveLosDelProyectoOrdenados`, `TestGetBacklog_PasaOrdenYEstado`, `TestPantallaBacklog_OrdenYFiltro` | `domain.OrdenarBacklog` |
| CA-07.5 Lista vacía | RN-07.5 | Un proyecto sin elementos tiene el backlog vacío `@limite` · Un filtro sin coincidencias da un backlog vacío `@limite` | `TestOrdenarBacklog_ListaVacia`, `TestRepositorioBacklog_ListarElementos_SinElementos`, `TestGetBacklog_SinElementos_DevuelveArrayVacio`, `TestPantallaBacklog_SinElementos_MuestraMensaje` | `domain.OrdenarBacklog`, `http.handlerListarBacklog` |
| CA-07.6 Solo su proyecto | RN-07.1 | El backlog solo muestra los elementos de su proyecto `@alternativo` | `TestListarBacklog_DevuelveLosDelProyectoOrdenados`, `TestRepositorioBacklog_ListarElementos_DevuelveLosDelProyectoCompletos` | `persistencia.RepositorioBacklog.ListarElementos` |
| CA-07.7 Proyecto inexistente | RN-07.1 | No se puede listar el backlog de un proyecto que no existe `@error` | `TestListarBacklog_Errores`, `TestGetBacklog_Errores`, `TestPantallaBacklog_ProyectoInexistente_Devuelve404` | `app.ListarBacklog` |
| CA-07.8 Orden o estado inválidos | Sección 7 | El orden tiene que ser asc o desc `@error` · El estado del filtro tiene que ser uno de los tres válidos `@error` | `TestOrdenarBacklog_Invalido_RetornaError`, `TestListarBacklog_Errores`, `TestGetBacklog_Errores`, `TestPantallaBacklog_OrdenInvalido_Devuelve400` | `domain.OrdenarBacklog` |

Pantalla (tarea 7.8, Definition of Done del Sprint 2): `/web/proyectos/{id}/backlog` lista el backlog con el orden y el filtro, y agrega elementos con el caso de uso de HU-05, que no tenía pantalla. Tests del alta: `TestPantallaBacklog_CrearElemento_RedirigeAlBacklog`, `TestPantallaBacklog_CrearSinStoryPoints`, `TestPantallaBacklog_StoryPointsNoNumericos_NoLlamaAlCasoDeUso`, `TestPantallaBacklog_ErrorAlCrear_MuestraElMensajeYConservaElFormulario`, `TestPantallaBacklog_CrearEnProyectoInexistente_Devuelve404` y `TestPantallaDetalle_EnlazaAlBacklog`. Que el alta conserve el orden y el filtro: `TestPantallaBacklog_FormularioConservaOrdenYFiltro`, `TestPantallaBacklog_CrearVuelveConElMismoOrdenYFiltro` y `TestPantallaBacklog_ErrorAlCrearConservaOrdenYFiltro`.

Listar no modifica nada (sección 3): `TestOrdenarBacklog_NoModificaLaListaRecibida`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog_listado.go](../internal/domain/backlog_listado.go) (`OrdenarBacklog`) | [backlog_listado_test.go](../internal/domain/backlog_listado_test.go) |
| Caso de uso | [internal/app/listar_backlog.go](../internal/app/listar_backlog.go) | [listar_backlog_test.go](../internal/app/listar_backlog_test.go) |
| Persistencia | [internal/adapters/persistencia/backlog_repo.go](../internal/adapters/persistencia/backlog_repo.go) (`ListarElementos`) | [listar_backlog_test.go](../internal/adapters/persistencia/listar_backlog_test.go) |
| HTTP | [internal/adapters/http/listar_backlog_handler.go](../internal/adapters/http/listar_backlog_handler.go) (`GET /proyectos/{id}/backlog`, `RegistrarListarBacklog`) | [listar_backlog_handler_test.go](../internal/adapters/http/listar_backlog_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [internal/adapters/http/pantallas_backlog.go](../internal/adapters/http/pantallas_backlog.go), [web/plantillas/backlog.html](../internal/adapters/http/web/plantillas/backlog.html) | [pantallas_backlog_test.go](../internal/adapters/http/pantallas_backlog_test.go) |
| BDD | [test/bdd/hu07_sistema_test.go](../test/bdd/hu07_sistema_test.go), [hu07_steps_test.go](../test/bdd/hu07_steps_test.go) | 14 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 7.1 Spec | `fbc91bb` docs(HU-07): [7.1] spec | — |
| 7.2 Escenarios BDD | `5d29f6b` test(HU-07): [7.2] escenarios BDD | — |
| 7.3 Steps godog | `9ac9ef4` test(HU-07): [7.3] steps de godog (en rojo) | BDD en rojo |
| 7.4 Dominio | `2347e81` test(HU-07): [7.4][RED] | **RED** |
| | `2cdd20f` feat(HU-07): [7.4][GREEN] | **GREEN** |
| 7.5 Caso de uso | `bfff2a8` test(HU-07): [7.5][RED] | **RED** |
| | `f7b57c0` feat(HU-07): [7.5][GREEN] | **GREEN** |
| 7.6 Persistencia | `9ec3df9` test(HU-07): [7.6][RED] | **RED** |
| | `8b37212` feat(HU-07): [7.6][GREEN] | **GREEN** |
| | `7f36d27` test(HU-07): escenarios BDD conectados (en verde) | BDD en verde |
| 7.7 HTTP | `bedaa9e` test(HU-07): [7.7][RED] | **RED** |
| | `fce353f` feat(HU-07): [7.7][GREEN] | **GREEN** |
| | `fd1eb05` refactor(HU-07): [7.7][REFACTOR] HU-05 usa aRespuestaElemento | **REFACTOR** |
| 7.8 Pantalla | `348dc5e` test(HU-07): [7.8][RED] | **RED** |
| | `4549f0f` feat(HU-07): [7.8][GREEN] | **GREEN** |
| | `82ac0e1` test(HU-07): [7.8][RED] el alta conserva el orden y el filtro | **RED** |
| | `7db8183` feat(HU-07): [7.8][GREEN] el alta conserva el orden y el filtro | **GREEN** |
| 7.9 Revisión y trazabilidad | `fb08cbe` docs(HU-07): [7.9] matriz de trazabilidad, y este commit | — |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.

---

## HU-08: Crear un Sprint con Sprint Goal

- **Issue:** [#19](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/19) · Sprint 2 · 3 SP
- **Spec:** [docs/specs/HU-08-crear-sprint.md](specs/HU-08-crear-sprint.md)
- **Feature:** [features/HU-08-crear-sprint.feature](../features/HU-08-crear-sprint.feature) (steps en [test/bdd/hu08_steps_test.go](../test/bdd/hu08_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-08.1 Crear el primer sprint | RN-08.6, RN-08.7, RN-08.8 | Crear el primer sprint de un proyecto `@normal` | `TestNuevoSprint_Valido`, `TestCrearSprint_Valido_GuardaYDevuelveConNumero`, `TestRepositorioSprints_GuardarYListar`, `TestPostSprints_DatosValidos_Devuelve201ConElSprint` | `domain.NuevoSprint`, `app.CrearSprint`, `persistencia.RepositorioSprints.GuardarSprint`, `http.handlerCrearSprint` |
| CA-08.2 Numeración por proyecto | RN-08.7 | El segundo sprint de un proyecto recibe el número siguiente `@normal` · Cada proyecto tiene su propia numeración `@alternativo` | `TestCrearSprint_Valido_GuardaYDevuelveConNumero`, `TestRepositorioSprints_NumeracionPorProyecto`, `TestCrearSprint_DosSimultaneos_RecibenNumerosDistintos` | `persistencia.RepositorioSprints.GuardarSprint` |
| CA-08.3 Proyecto inexistente | RN-08.1 | No se puede crear un sprint en un proyecto que no existe `@error` | `TestCrearSprint_Invalido_RetornaErrorYNoGuarda` (proyecto inexistente), `TestRepositorioSprints_ProyectoInexistente_RetornaError`, `TestPostSprints_ProyectoInexistente_Devuelve404` | `app.CrearSprint` |
| CA-08.4 Proyecto no activo | RN-08.2 | Solo se crean sprints en proyectos activos `@error` (4 ejemplos) | `TestNuevoSprint_Invalido_RetornaError` (4 estados), `TestCrearSprint_Invalido_RetornaErrorYNoGuarda` (proyecto no activo), `TestPostSprints_DatosInvalidos_Devuelve400` | `domain.NuevoSprint` |
| CA-08.5 Fechas | RN-08.3 | Las fechas son obligatorias `@error` (2 ejemplos) · La fecha de fin no puede ser anterior `@error` · Un sprint de un solo día es válido `@limite` | `TestNuevoSprint_Invalido_RetornaError`, `TestNuevoSprint_CasosLimiteValidos` (un día), `TestPostSprints_SinFechas_LlegaComoFechaCero`, `TestPostSprints_DatosInvalidos_Devuelve400` | `domain.NuevoSprint`, `http.pedidoCrearSprint.aDatos` |
| CA-08.6 Sprint Goal | RN-08.6 | El Sprint Goal es obligatorio `@error` · Un Sprint Goal de solo espacios cuenta como vacío `@error` · No puede superar los 200 caracteres `@error` · Uno de exactamente 200 es válido `@limite` | `TestNuevoSprint_Invalido_RetornaError`, `TestNuevoSprint_CasosLimiteValidos` (200 caracteres), `TestCrearSprint_Invalido_RetornaErrorYNoGuarda` (Sprint Goal vacío) | `domain.NuevoSprint` |
| CA-08.7 Dentro del proyecto | RN-08.4 | El sprint tiene que estar dentro de las fechas del proyecto `@error` (2 ejemplos) · Justo en los límites es válido `@limite` (2 ejemplos) · Proyecto sin fecha de fin `@alternativo` | `TestNuevoSprint_Invalido_RetornaError`, `TestNuevoSprint_CasosLimiteValidos` | `domain.NuevoSprint`, `domain.dentroDelProyecto` |
| CA-08.8 Sin superposición | RN-08.5 | Los sprints de un proyecto no se superponen `@error` (5 ejemplos) · Uno que empieza el día siguiente al fin de otro es válido `@limite` | `TestNuevoSprint_Invalido_RetornaError` (5 casos), `TestNuevoSprint_CasosLimiteValidos`, `TestNuevoSprint_ComparaPorDiaCalendario`, `TestRepositorioSprints_SuperpuestoAlGuardar_RetornaError`, `TestCrearSprint_DosSimultaneosSuperpuestos_SegundoEsRechazado` | `domain.seSuperponen`, `persistencia.RepositorioSprints.GuardarSprint` |
| CA-08.9 Listar los sprints | Sección 3 | Los sprints del proyecto se listan ordenados por número `@normal` | `TestPantallaSprints_ListaYFormulario`, `TestPantallaSprints_SinSprints_MuestraMensaje`, `TestListarSprints_DevuelveLosDelProyecto`, `TestListarSprints_ProyectoInexistente_RetornaError`, `TestRepositorioSprints_GuardarYListar`, `TestGetSprints_DevuelveLosSprintsDelProyecto`, `TestGetSprints_SinSprints_DevuelveArrayVacio`, `TestGetSprints_Errores` | `app.ListarSprints`, `persistencia.RepositorioSprints.ListarSprints`, `http.handlerListarSprints` |

Pantalla (tarea 8.8, Definition of Done del Sprint 2): `/web/proyectos/{id}/sprints` lista los sprints y crea uno con los mismos casos de uso que la API. Tests: `TestPantallaSprints_ListaYFormulario`, `TestPantallaSprints_SinSprints_MuestraMensaje`, `TestPantallaSprints_ProyectoInexistente_Devuelve404`, `TestPantallaSprints_Crear_RedirigeALaLista`, `TestPantallaSprints_ErrorDeValidacion_MuestraElMensajeYConservaElFormulario`, `TestPantallaSprints_FechaInvalida_NoLlamaAlCasoDeUso`, `TestPantallaSprints_CrearEnProyectoInexistente_Devuelve404` y `TestPantallaDetalle_EnlazaALosSprints`.

Pedidos simultáneos (sección 5 de la spec; sin CA propio porque no se puede expresar en Gherkin): `TestCrearSprint_DosSimultaneos_RecibenNumerosDistintos` y `TestCrearSprint_DosSimultaneosSuperpuestos_SegundoEsRechazado`. El número y el control de superposición se resuelven en un solo `INSERT`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/sprint.go](../internal/domain/sprint.go) | [sprint_test.go](../internal/domain/sprint_test.go) |
| Casos de uso | [internal/app/gestionar_sprints.go](../internal/app/gestionar_sprints.go) (`CrearSprint`, `ListarSprints`) | [gestionar_sprints_test.go](../internal/app/gestionar_sprints_test.go) |
| Persistencia | [internal/adapters/persistencia/sprint_repo.go](../internal/adapters/persistencia/sprint_repo.go), [db.go](../internal/adapters/persistencia/db.go) (tabla `sprints`) | [sprint_repo_test.go](../internal/adapters/persistencia/sprint_repo_test.go) |
| HTTP | [internal/adapters/http/sprint_handler.go](../internal/adapters/http/sprint_handler.go) (`POST` y `GET /proyectos/{id}/sprints`, `RegistrarSprints`) | [sprint_handler_test.go](../internal/adapters/http/sprint_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [internal/adapters/http/pantallas_sprints.go](../internal/adapters/http/pantallas_sprints.go), [web/plantillas/sprints.html](../internal/adapters/http/web/plantillas/sprints.html) | [pantallas_sprints_test.go](../internal/adapters/http/pantallas_sprints_test.go) |
| BDD | [test/bdd/hu08_sistema_test.go](../test/bdd/hu08_sistema_test.go), [hu08_steps_test.go](../test/bdd/hu08_steps_test.go) | 28 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 8.1 Spec | `4e8aa05` docs(HU-08): [8.1] spec de crear un sprint | — |
| 8.2 Escenarios BDD | `6c91cca` test(HU-08): [8.2] escenarios BDD | — |
| 8.3 Steps godog | `5057f43` test(HU-08): [8.3] steps de godog (en rojo) | BDD en rojo |
| 8.4 Entidad Sprint | `a7b2947` test(HU-08): [8.4][RED] | **RED** |
| | `6679372` feat(HU-08): [8.4][GREEN] | **GREEN** |
| 8.5 Casos de uso | `90786ee` test(HU-08): [8.5][RED] | **RED** |
| | `7b36d87` feat(HU-08): [8.5][GREEN] | **GREEN** |
| 8.6 Persistencia | `9a97f16` test(HU-08): [8.6][RED] | **RED** |
| | `5b371d7` feat(HU-08): [8.6][GREEN] | **GREEN** |
| | `4bafeeb` test(HU-08): escenarios BDD conectados (en verde) | BDD en verde |
| 8.7 HTTP | `0f36e9e` test(HU-08): [8.7][RED] | **RED** |
| | `e603d46` feat(HU-08): [8.7][GREEN] | **GREEN** |
| 8.8 Pantalla | `37d38d3` test(HU-08): [8.8][RED] pantalla de sprints | **RED** |
| | `f00dc5a` feat(HU-08): [8.8][GREEN] pantalla de sprints | **GREEN** |
| | `60e47c8` refactor(HU-08): [8.8][REFACTOR] botones de la cabecera | **REFACTOR** |
| 8.9 Trazabilidad | `6ef97ad` docs(HU-08): [8.9] matriz de trazabilidad, y este commit | — |

### Huecos

- La spec sigue en estado **Borrador**: pasa a Aprobada cuando se apruebe el PR.

---

## HU-18: Registrar esfuerzo en una historia o tarea

- **Issue:** #29 · Sprint 1 · 3 SP · pantalla: #289
- **Spec:** [docs/specs/HU-18-registrar-esfuerzo.md](specs/HU-18-registrar-esfuerzo.md)
- **Feature:** [features/HU-18-registrar-esfuerzo.feature](../features/HU-18-registrar-esfuerzo.feature) (steps en [test/bdd/hu18_steps_test.go](../test/bdd/hu18_steps_test.go))
- **Ruta HTTP:** `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`. **Pendiente de confirmar (P-6):** reemplaza el `POST /esfuerzo` de la tarea 18.7 y falta el OK de Matías.
- **Pantalla:** sección "Registrar esfuerzo" de `GET /web/proyectos/{id}/backlog/{elementoID}` (HU-06), con la acción `POST /web/proyectos/{id}/backlog/{elementoID}/esfuerzo`.

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-18.1 Registro válido con ID y los datos enviados | RN-18.10 | Registrar esfuerzo con todos los datos válidos `@normal` · La actividad se guarda sin los espacios de alrededor `@alternativo` | `TestNuevoRegistroEsfuerzo_Valido`, `TestNuevoRegistroEsfuerzo_RecortaIntegranteYActividad`, `TestRegistrarEsfuerzo_DatosValidos_GuardaConID`, `TestPostEsfuerzo_DatosValidos_Devuelve201ConElRegistro` | `domain.NuevoRegistroEsfuerzo`, `app.RegistrarEsfuerzo.Ejecutar` |
| CA-18.2 Media hora en un elemento sin esfuerzo | RN-18.6, RN-18.7 | Media hora es válida en un elemento sin esfuerzo registrado `@limite` | `TestNuevoRegistroEsfuerzo_HorasLimite_SonValidas` | `domain.validarHoras` |
| CA-18.3 Horas 0 o negativas | RN-18.6 | Las horas deben ser mayores a 0, múltiplo de media hora y como máximo 24 `@error` (ejemplos 0, -2, -1,25) | `TestNuevoRegistroEsfuerzo_HorasNegativas_RetornaError` (incluye NaN y -Inf), `TestRegistrarEsfuerzo_DatosInvalidos_NoLleganAlRepositorio`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `domain.validarHoras` |
| CA-18.4 Actividad vacía o con espacios | RN-18.5 | La actividad es obligatoria `@error` (×2) | `TestNuevoRegistroEsfuerzo_SinActividad_RetornaError`, `TestRegistrarEsfuerzo_DatosInvalidos_NoLleganAlRepositorio`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `domain.NuevoRegistroEsfuerzo` |
| CA-18.5 Integrante vacío | RN-18.2 | El integrante es obligatorio `@error` (×2) | `TestNuevoRegistroEsfuerzo_SinIntegrante_RetornaError`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `domain.NuevoRegistroEsfuerzo` |
| CA-18.6 Fecha obligatoria, no futura; hoy es válida | RN-18.4 | La fecha es obligatoria `@error` · La fecha no puede ser futura `@error` · La fecha de hoy es válida `@limite` | `TestNuevoRegistroEsfuerzo_SinFecha_RetornaError`, `_FechaCeroEnOtraZona_RetornaError`, `_FechaFutura_RetornaError`, `_FechaHoyOPasada_EsValidaYGuardaSoloElDia`, `_FechaYHoyEnZonasDistintas_ComparaElDiaDeCalendario`, `TestRegistrarEsfuerzo_DatosInvalidos_NoLleganAlRepositorio` (reloj), `TestPostEsfuerzo_DatosValidos_PasaLosDatosAlCasoDeUso` (zona local) | `domain.NuevoRegistroEsfuerzo`, `domain.soloDia`, `domain.diaCalendario`, `http.fechaDelPedido` |
| CA-18.7 Proyecto o elemento inexistente, o elemento de otro proyecto | RN-18.1 | No se puede registrar esfuerzo en un proyecto que no existe `@error` · … en un elemento que no existe `@error` · Un elemento de otro proyecto se trata como inexistente `@error` | `TestRegistrarEsfuerzo_ProyectoInexistente_RetornaErrorYNoGuarda`, `_ElementoInexistente_RetornaErrorYNoGuarda`, `_ElementoDeOtroProyecto_RetornaErrorYNoGuarda`, `TestRepositorioEsfuerzo_ElementoInexistente_RetornaErrElementoNoEncontrado`, `TestPostEsfuerzo_ProyectoOElementoInexistente_Devuelve404` | `app.RegistrarEsfuerzo.Ejecutar`, `persistencia.RepositorioEsfuerzo.GuardarRegistro` (clave foránea) |
| CA-18.8 Integrante que no pertenece al proyecto | RN-18.3 | Quien registra tiene que ser integrante del proyecto `@error` · Un integrante de otro proyecto no puede registrar en este `@error` · El nombre del integrante se compara sin espacios ni mayúsculas `@alternativo` · Un integrante quitado del proyecto ya no puede registrar esfuerzo nuevo `@error` | `TestRegistrarEsfuerzo_IntegranteQueNoPertenece_RetornaErrorYNoGuarda`, `_IntegranteConOtrasMayusculas_GuardaElNombreDelProyecto`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `app.RegistrarEsfuerzo.Ejecutar`, `app.nombreEnProyecto` |
| CA-18.9 Horas que no son múltiplo de 0,5 | RN-18.7 | Outline de horas `@error` (ejemplos 1,25 y 25,3) | `TestNuevoRegistroEsfuerzo_HorasNegativas_RetornaError`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `domain.validarHoras` |
| CA-18.10 24 h es válido, 24,5 h no | RN-18.8 | Outline de horas `@error` (ejemplo 24,5) · Un registro de exactamente 24 horas es válido `@limite` | `TestNuevoRegistroEsfuerzo_HorasNegativas_RetornaError`, `_HorasLimite_SonValidas`, `TestPostEsfuerzo_DatosInvalidos_Devuelve400ConElError` | `domain.validarHoras`, `domain.HorasMaximasPorRegistro` |
| CA-18.11 El tope no suma los registros del día | RN-18.8, RN-18.9 | El tope de 24 horas es por registro, no se suman los registros del día `@alternativo` | ✗ sin test unitario propio (ver Huecos) | `domain.validarHoras` (mira un solo registro) |
| CA-18.12 Dos registros iguales, IDs distintos, recuperables | RN-18.9, RN-18.10 | Dos registros con los mismos datos son independientes `@normal` | `TestRegistrarEsfuerzo_DatosValidos_GuardaConID`, `TestRepositorioEsfuerzo_GuardarYBuscar_DevuelveLoMismo`, `_DosRegistrosIguales_TienenIDsDistintos` | `persistencia.RepositorioEsfuerzo.GuardarRegistro` / `BuscarRegistroPorID` |
| CA-18.13 Los registros se conservan al quitar al integrante | RN-18.3 | Los registros se conservan al quitar al integrante del proyecto `@alternativo` · Un integrante quitado del proyecto ya no puede registrar esfuerzo nuevo `@error` | `TestRepositorioEsfuerzo_QuitarAlIntegrante_ConservaSusRegistros` | Tabla `registros_esfuerzo` en `db.go` (sin clave foránea a `integrantes`) |
| CA-18.14 Proyecto no Activo y fecha anterior al inicio | RN-18.4 | Se registra esfuerzo en un proyecto no Activo y con fecha anterior a su inicio `@alternativo` | `TestNuevoRegistroEsfuerzo_FechaHoyOPasada_EsValidaYGuardaSoloElDia` (fecha de hace meses) | `domain.NuevoRegistroEsfuerzo` y `app.RegistrarEsfuerzo` no miran las fechas ni el estado del proyecto |
| CA-18.15 Respuestas HTTP 201 / 400 / 404 | — (secciones 3 y 7) | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y no pasan por HTTP. Se cubre con los tests del handler (tarea 18.7). | `TestPostEsfuerzo_DatosValidos_Devuelve201ConElRegistro`, `_ProyectoOElementoInexistente_Devuelve404`, `_DatosInvalidos_Devuelve400ConElError`, `_PedidoMalFormado_Devuelve400SinLlamarAlCasoDeUso`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerRegistrarEsfuerzo`, `http.RegistrarRegistrarEsfuerzo`, `http.NuevoRouter` |
| CA-18.16 Formulario inicial, sin registros ni totales | Sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_FormularioInicial` | `elemento.html`, `http.pantallas.seccionEsfuerzoInicial`, `http.pantallas.mostrarElemento` |
| CA-18.17 Enlace "Registrar esfuerzo" en cada fila del backlog | Sección 9 | — (test del adaptador) | `TestPantallaBacklog_EnlazaARegistrarEsfuerzoEnCadaElemento` | `backlog.html` |
| CA-18.18 Envío válido: 303, confirmación y sin registrar otro al recargar | Sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_EnvioValido_RedirigeYConfirma`, `TestPantallaEsfuerzo_OtroValorDelParametro_NoConfirma`, `TestPantallaEsfuerzo_ParametroDeConfirmacion_NoInyecta`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.pantallas.registrarPantallasEsfuerzo`, `http.pantallas.registrarEsfuerzoPantalla`, `http.pantallas.elemento`, `elemento.html` |
| CA-18.19 Errores del dominio: 400, mensaje sin envolver y lo tipeado conservado | RN-18.2, RN-18.4 a RN-18.8, sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_ErroresDelDominio_VuelvenAlFormulario`, `TestPantallaEsfuerzo_ConError_ConservaLoTipeado`, `TestPantallaEsfuerzo_LoTipeado_SeEscapa` | `http.mensajeDeValidacionEsfuerzo`, `http.formMensajeDeErrorConocido`, `http.pantallas.mostrarErrorDeEsfuerzo` |
| CA-18.20 Fecha u horas ilegibles, sin llamar al caso de uso; horas vacías | Sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_ErroresDeFormato_NoLlamanAlCasoDeUso`, `TestPantallaEsfuerzo_HorasVacias_MensajeDeHorasInvalidas` | `http.datosRegistroDeForm`, `http.horasDeForm`, `http.fechaDelPedido` |
| CA-18.21 Integrante quitado después de abrir la pantalla | RN-18.3, sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_IntegranteQuitado_400SinElDesplegable` | `http.pantallas.mostrarElemento`, `elemento.html` |
| CA-18.22 Proyecto o elemento inexistente, de otro proyecto o ID inválido | RN-18.1, sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_NoEncontrados_Devuelven404ConLaPaginaDeError`, `TestPantallaEsfuerzo_ErrorDeFormatoConElementoInexistente_Devuelve404`, `TestPantallaEsfuerzo_ErrorDeFormatoConProyectoInexistente_Devuelve404`, `TestPantallaElemento_ProyectoInexistente_Devuelve404`, `TestPantallaEsfuerzo_IDInvalido_Devuelve400` | `http.pantallas.errorDeElemento`, `http.pantallas.mostrarErrorDeEsfuerzo`, `http.pantallas.idsElemento` |
| CA-18.23 Otro sitio: 403; falla interna: 500 sin el detalle | Sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_DesdeOtroSitio_SeRechaza`, `TestPantallaEsfuerzo_FallaInterna_Devuelve500SinDetalle` | `http.pantallas.formulario`, `http.mismoOrigen`, `http.pantallas.fallaInterna` |
| CA-18.24 Proyecto sin integrantes: aviso y enlace al detalle | Sección 9 | — (test del adaptador) | `TestPantallaEsfuerzo_ProyectoSinIntegrantes_AvisaYEnlaza` | `elemento.html` |
| CA-18.25 Los errores de editar y eliminar conservan la sección de esfuerzo | Sección 9 | — (test del adaptador) | `TestPantallaElemento_ErroresDeEditarYEliminar_MuestranLaSeccionDeEsfuerzo` | `http.pantallas.mostrarElemento`, `http.pantallas.editarElemento`, `http.pantallas.eliminarElemento` |

CA-18.16 a CA-18.25 no tienen escenario BDD, como dice la spec: se cubren con tests del adaptador (`httptest`), igual que las pantallas de HU-10 y HU-11. Los tests están en [pantallas_esfuerzo_test.go](../internal/adapters/http/pantallas_esfuerzo_test.go), salvo `TestNuevoRouter_RegistraTodasLasRutas`, que está en [router_test.go](../internal/adapters/http/router_test.go).

**Test de integración entre dos historias, sin criterio propio:** `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` (`067e249`), en el mismo archivo, se traza a la sección 9 de HU-18 y a CA-14.9 de HU-14. Verifica que las dos secciones convivan en la pantalla del elemento, en orden (edición, Planning Poker, esfuerzo, eliminar), sin ids duplicados ni formularios anidados, sin ronda y con una ronda abierta. Ningún CA pide esto, así que no suma a los contadores de la tabla de resumen.

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-18.1 El elemento existe y es del proyecto | CA-18.7 | `app.RegistrarEsfuerzo.Ejecutar` (`e.ProyectoID != p.ID`) |
| RN-18.2 Integrante obligatorio | CA-18.5 | `domain.NuevoRegistroEsfuerzo` |
| RN-18.3 Integrante del proyecto; sus registros se conservan si se lo quita | CA-18.8, CA-18.13 | `app.nombreEnProyecto`; tabla `registros_esfuerzo` sin clave foránea a `integrantes` |
| RN-18.4 Fecha obligatoria, no futura, por día y en la zona del servidor; sin control de fechas ni estado del proyecto | CA-18.6, CA-18.14 | `domain.NuevoRegistroEsfuerzo`, `domain.diaCalendario`, `http.fechaDelPedido`, `persistencia.BuscarRegistroPorID` (lee con `time.Local`) |
| RN-18.5 Actividad obligatoria, texto libre | CA-18.4 | `domain.NuevoRegistroEsfuerzo` |
| RN-18.6 Horas mayores a 0 | CA-18.2, CA-18.3 | `domain.validarHoras` |
| RN-18.7 Múltiplos de media hora | CA-18.2, CA-18.9 | `domain.validarHoras`; en la base, `medias_horas` entero (`persistencia.aMediasHoras`) |
| RN-18.8 Máximo 24 h por registro | CA-18.10, CA-18.11 | `domain.validarHoras`, `domain.HorasMaximasPorRegistro` |
| RN-18.9 Registros independientes | CA-18.11, CA-18.12 | `persistencia.RepositorioEsfuerzo.GuardarRegistro` (sin restricción UNIQUE) |
| RN-18.10 ID único autogenerado | CA-18.1, CA-18.12 | `persistencia.RepositorioEsfuerzo.GuardarRegistro` (`AUTOINCREMENT`) |

Sin CA propio en la spec:
- Orden de los errores (sección 7): `TestNuevoRegistroEsfuerzo_VariosErrores_InformaElPrimero`, `TestRegistrarEsfuerzo_ProyectoInexistente_RetornaErrorYNoGuarda` (con datos inválidos) y `TestRegistrarEsfuerzo_IntegranteAjenoYDatosInvalidos_InformaElErrorDelDominio`. El orden entre las validaciones de horas lo cubren los ejemplos -1,25 y 25,3.
- Cuerpo mal formado, fecha con otro formato o ids inválidos, que dan 400 sin llamar al caso de uso (sección 7): `TestPostEsfuerzo_PedidoMalFormado_Devuelve400SinLlamarAlCasoDeUso`.
- Falla interna (sección 7): `TestRegistrarEsfuerzo_FallaDelRepositorio_SePropaga`, `TestPostEsfuerzo_FallaInterna_Devuelve500SinDetalles`.
- Horas guardadas como medias horas y fecha como `AAAA-MM-DD` (sección 5): `TestRepositorioEsfuerzo_GuardaHorasEnMediasHorasYFechaComoTexto`.
- Búsqueda por ID de un registro inexistente (sección 7, última fila): `TestRepositorioEsfuerzo_BuscarInexistente_RetornaErrRegistroEsfuerzoNoEncontrado`.
- Pantalla, fecha en la zona del servidor (RN-18.4, sección 5): `TestPantallaEsfuerzo_FechaEnLaZonaDelServidor`.
- Pantalla, horas leídas sin los espacios de alrededor (sección 9): `TestPantallaEsfuerzo_HorasConEspacios_SeLeen`.
- Pantalla, falla interna al consultar los integrantes (sección 9): `TestPantallaElemento_FallaAlConsultarElProyecto_Devuelve500SinDetalle`.
- Pantalla, formulario de más de 64 KiB, que da 413 (sección 9): `TestPantallaEsfuerzo_FormularioGrande_Devuelve413`.

### Errores nuevos

Todos están definidos en [internal/domain/registro_esfuerzo.go](../internal/domain/registro_esfuerzo.go). HU-18 también reutiliza `ErrProyectoNoEncontrado` (HU-01), `ErrNombreIntegranteVacio` y `ErrIntegranteNoEncontrado` (HU-03) y `ErrElementoNoEncontrado` (HU-05).

| Error | Mensaje | Línea | Lo devuelve | HTTP |
|---|---|---|---|---|
| `ErrFechaEsfuerzoObligatoria` | "la fecha es obligatoria" | 15 | `domain.NuevoRegistroEsfuerzo` | 400 |
| `ErrFechaEsfuerzoFutura` | "la fecha no puede ser futura" | 16 | `domain.NuevoRegistroEsfuerzo` | 400 |
| `ErrActividadVacia` | "la actividad es obligatoria" | 17 | `domain.NuevoRegistroEsfuerzo` | 400 |
| `ErrHorasInvalidas` | "las horas deben ser mayores a 0" | 18 | `domain.validarHoras` | 400 |
| `ErrHorasNoMultiploMediaHora` | "las horas se registran en múltiplos de media hora" | 19 | `domain.validarHoras` | 400 |
| `ErrHorasExcedenMaximo` | "un registro no puede superar las 24 horas" | 20 | `domain.validarHoras` | 400 |
| `ErrRegistroEsfuerzoNoEncontrado` | "registro de esfuerzo no encontrado" | 23 | `persistencia.RepositorioEsfuerzo.BuscarRegistroPorID` | — (no hay endpoint de lectura en HU-18) |

La pantalla no agrega errores de dominio. El único mensaje nuevo es de formato y vive en el adaptador: "las horas deben ser un número" (`errFormulario` en `http.horasDeForm`).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/registro_esfuerzo.go](../internal/domain/registro_esfuerzo.go) (`RegistroEsfuerzo`, `NuevoRegistroEsfuerzo`, `validarHoras`, `soloDia`, `diaCalendario`) | [registro_esfuerzo_test.go](../internal/domain/registro_esfuerzo_test.go) |
| Caso de uso | [internal/app/registrar_esfuerzo.go](../internal/app/registrar_esfuerzo.go) (`RegistrarEsfuerzo`, puertos `LectorElementos` y `RepositorioEsfuerzo`; reutiliza `LectorProyectos` de HU-04) | [registrar_esfuerzo_test.go](../internal/app/registrar_esfuerzo_test.go) |
| Persistencia | [internal/adapters/persistencia/esfuerzo_repo.go](../internal/adapters/persistencia/esfuerzo_repo.go) (`GuardarRegistro`, `BuscarRegistroPorID`), [db.go](../internal/adapters/persistencia/db.go) (tabla `registros_esfuerzo`) | [esfuerzo_repo_test.go](../internal/adapters/persistencia/esfuerzo_repo_test.go), contra SQLite real |
| HTTP | [internal/adapters/http/esfuerzo_handler.go](../internal/adapters/http/esfuerzo_handler.go) (`RegistrarRegistrarEsfuerzo`), registrada en [router.go](../internal/adapters/http/router.go) y conectada en [cmd/server/main.go](../cmd/server/main.go) con `time.Now` como reloj | [esfuerzo_handler_test.go](../internal/adapters/http/esfuerzo_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla (archivo nuevo) | [internal/adapters/http/pantallas_esfuerzo.go](../internal/adapters/http/pantallas_esfuerzo.go) (`registrarPantallasEsfuerzo`, `registrarEsfuerzoPantalla`, `mostrarErrorDeEsfuerzo`, `seccionEsfuerzoInicial`, `datosRegistroDeForm`, `horasDeForm`, `mensajeDeValidacionEsfuerzo`) | [pantallas_esfuerzo_test.go](../internal/adapters/http/pantallas_esfuerzo_test.go) (nuevo), [router_test.go](../internal/adapters/http/router_test.go) (ruta nueva), [pantallas_test.go](../internal/adapters/http/pantallas_test.go) (`pedirAlMux`) |
| Pantalla (archivos tocados) | [pantallas.go](../internal/adapters/http/pantallas.go) (campo `reloj`, registro de la ruta, `formMensajeDeErrorConocido`), [pantallas_elemento.go](../internal/adapters/http/pantallas_elemento.go) (`datosEdicionElemento`, `mostrarElemento`, `errorDeElemento`), [pantallas_defecto.go](../internal/adapters/http/pantallas_defecto.go) (solo por el refactor), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html), [backlog.html](../internal/adapters/http/web/plantillas/backlog.html), [estilos.css](../internal/adapters/http/web/static/estilos.css) (`.exito`) | Los mismos de la fila anterior |
| BDD | [test/bdd/hu18_sistema_test.go](../test/bdd/hu18_sistema_test.go), [hu18_steps_test.go](../test/bdd/hu18_steps_test.go) (reloj falso) | 28 escenarios en verde |

### Commits

Salen de `git log main..HEAD`, en orden.

| Tarea | Commit | Etapa |
|---|---|---|
| 18.1 Spec | `b02e049` docs(HU-18): [18.1] spec SDD de registrar esfuerzo | — |
| 18.2 Escenarios BDD | `930b2ee` test(HU-18): [18.2] escenarios BDD de registrar esfuerzo | — |
| 18.3 Steps godog | `58f074e` test(HU-18): [18.3][RED] steps de godog | BDD en rojo (28/28) |
| 18.4 Entidad RegistroEsfuerzo | `c2b32b9` test: [18.4][RED] registro de esfuerzo válido | **RED** |
| | `f60b4a3` feat: [18.4][GREEN] entidad RegistroEsfuerzo con constructor | **GREEN** |
| | `bcbcdfd` test: [18.4][RED] integrante y actividad obligatorios y recortados | **RED** |
| | `6a2ab52` feat: [18.4][GREEN] validar y recortar integrante y actividad | **GREEN** |
| | `dedcbcb` test: [18.4][RED] fecha obligatoria, no futura y comparada por día | **RED** |
| | `4422271` feat: [18.4][GREEN] validar la fecha del registro por día | **GREEN** |
| | `dd650d8` test: [18.4][RED] horas mayores a 0, múltiplo de media hora y máximo 24 | **RED** |
| | `034f68e` feat: [18.4][GREEN] validar las horas en el orden de la spec | **GREEN** |
| | `a726057` refactor: [18.4][REFACTOR] extraer validarHoras, soloDia y el tope de 24 h | **REFACTOR** |
| | `0c3dff1` test: [18.4][RED] comparar la fecha por día de calendario aunque las zonas difieran | **RED** |
| | `4ead02e` feat: [18.4][GREEN] comparar fecha y hoy por día de calendario | **GREEN** |
| | `c076181` test: [18.4] steps BDD comparan los errores nuevos con errors.Is | — |
| 18.5 Caso de uso | `1dc245e` test: [18.5][RED] registrar esfuerzo válido y elemento inexistente | **RED** |
| | `089a968` feat: [18.5][GREEN] caso de uso RegistrarEsfuerzo busca el elemento, valida y guarda | **GREEN** |
| | `65c101e` test: [18.5][RED] proyecto inexistente y elemento de otro proyecto | **RED** |
| | `6c66a9f` feat: [18.5][GREEN] verificar el proyecto y que el elemento sea suyo | **GREEN** |
| | `1166a40` test: [18.5][RED] el integrante pertenece al proyecto y se guarda con su nombre | **RED** |
| | `9671ed6` feat: [18.5][GREEN] verificar que el integrante pertenezca al proyecto | **GREEN** |
| | `4c0285e` refactor: [18.5][REFACTOR] extraer nombreEnProyecto y documentar el caso de uso | **REFACTOR** |
| 18.6 Persistencia | `4e30797` test: [18.6][RED] guardar y recuperar registros de esfuerzo en SQLite | **RED** |
| | `929d265` feat: [18.6][GREEN] RepositorioEsfuerzo guarda, busca y cuenta registros en SQLite | **GREEN** |
| | `d826626` test: [18.6][RED] elemento inexistente al guardar el registro | **RED** |
| | `1531f8e` feat: [18.6][GREEN] traducir la clave foránea del elemento a ErrElementoNoEncontrado | **GREEN** |
| | `7d7d40d` refactor: [18.6][REFACTOR] extraer la conversión de horas a medias horas | **REFACTOR** |
| | `7dc1955` test: conectar los escenarios BDD con RegistrarEsfuerzo y RepositorioEsfuerzo | BDD en verde (104/104) |
| — | `cf392b4` Merge branch 'main' into hu-18-registrar-esfuerzo | merge de `main`, no es de HU-18 |
| 18.7 HTTP | `0cd6453` test: [18.7][RED] POST de esfuerzo: 201, 404, 400 y 500 | **RED** |
| | `e96309f` feat: [18.7][GREEN] handler POST /proyectos/{id}/backlog/{elementoID}/esfuerzo | **GREEN** |
| | `0019132` test: [18.7][RED] NuevoRouter registra la ruta de esfuerzo | **RED** |
| | `a686adb` feat: [18.7][GREEN] sumar RegistrarEsfuerzo al router y conectarlo en main con time.Now | **GREEN** |
| | `1c8c7d6` refactor: [18.7][REFACTOR] extraer elementoIDDesdeRuta, fechaDelPedido y la respuesta | **REFACTOR** |
| 18.8 Revisión | `d49fbf2` test: [18.8][RED] fecha cero en otra zona es fecha faltante | **RED** |
| | `3894ece` fix: [18.8][GREEN] preguntar IsZero antes de recortar la fecha al día | **GREEN** |
| | `8cd50f6` test: [18.8][RED] NaN da ErrHorasInvalidas | **RED** |
| | `0462911` fix: [18.8][GREEN] NaN no es mayor a 0 y da ErrHorasInvalidas | **GREEN** |
| | `bb9f3f6` docs: [18.8] separar RN-18.8 y RN-18.9 y documentar ErrRegistroEsfuerzoNoEncontrado | — |
| | `c4864d2` refactor: [18.8] ContarRegistros pasa de producción a helpers de test | refactor (sin etiqueta `[REFACTOR]`) |
| | `c872259` docs: [18.8] explicar por qué la clave foránea se traduce a ErrElementoNoEncontrado | — |
| 18.8 Trazabilidad | este commit | — |

En la tabla se omite el scope `(HU-18)` de los mensajes a partir de 18.4: todos los commits lo llevan.

**Ampliación con la pantalla (issue #289).** Salen de `git log origin/main..HEAD`, en la rama `hu-18-pantalla-esfuerzo`.

| Tarea | Commit | Etapa |
|---|---|---|
| Spec de la pantalla | `df4922e` docs(HU-18): spec de la pantalla de registrar esfuerzo | — |
| Pantalla | `198394e` [RED] → `0ade8e9` [GREEN] → `c6b1694` [REFACTOR] | **RED → GREEN → REFACTOR**. El RED incluye el campo `reloj` en `pantallas.go`: es código de producción mínimo para que los tests compilen y fallen por aserciones, no por compilación. |
| Estilo de la confirmación | `72bb594` style(HU-18): estilo de la confirmación .exito | — |
| Revisión (go-reviewer): error de formato con proyecto inexistente | `33efaf2` [RED] → `790481b` [GREEN] | **RED → GREEN** (sin REFACTOR: el cambio son unas pocas líneas al principio de `mostrarErrorDeEsfuerzo`) |
| Revisión (go-reviewer): escape de HTML | `851ec94` test(HU-18): proteger el escape de HTML en el formulario de esfuerzo | — (tests nuevos sobre código ya en verde; se verificaron con las mutaciones A a D) |
| Revisión (go-reviewer): aserciones menos frágiles | `2eea25c` test(HU-18): aserciones menos frágiles de la pantalla de esfuerzo | — (se verificaron con las mutaciones E a I) |
| Trazabilidad de la pantalla | `071011d` docs(trazabilidad): ampliar HU-18 con la pantalla | — |
| Merge con `main` (HU-14 y HU-15) | `cc1a242` Merge remote-tracking branch 'origin/main' into hu-18-pantalla-esfuerzo | merge de `main`, con dos choques que git no detectó (ver Huecos) |
| Convivencia con Planning Poker | `067e249` test(HU-18): la pantalla del elemento muestra Planning Poker y Registrar esfuerzo juntos | — (sin RED real; se verificó con mutaciones, ver Huecos) |
| Trazabilidad de la convivencia | este commit | — |

### Evidencia de la pantalla: mutaciones manuales

Cada mutación se aplicó sobre el código, se corrieron los tests y se revirtió con `git checkout`. Ninguna quedó en un commit.

| # | Mutación | Resultado |
|---|---|---|
| 1 | `mensajeDeValidacion` en lugar de la lista de HU-18 | Fallan `TestPantallaEsfuerzo_ErroresDelDominio_VuelvenAlFormulario`, `TestPantallaEsfuerzo_ConError_ConservaLoTipeado` y `TestPantallaEsfuerzo_HorasVacias_MensajeDeHorasInvalidas` |
| 2 | `fechaDeForm` (UTC) en lugar de `fechaDelPedido` | Fallan `TestPantallaEsfuerzo_FechaEnLaZonaDelServidor` y `TestPantallaEsfuerzo_ErroresDeFormato_NoLlamanAlCasoDeUso` |
| 3 | `errorDeElemento` sin `ErrProyectoNoEncontrado` | Fallan `TestPantallaEsfuerzo_NoEncontrados_Devuelven404ConLaPaginaDeError` y `TestPantallaElemento_ProyectoInexistente_Devuelve404` |
| 4 | Fecha inicial en UTC en lugar de la zona del reloj | Fallan `TestPantallaEsfuerzo_FormularioInicial`, `TestPantallaEsfuerzo_EnvioValido_RedirigeYConfirma` y `TestPantallaElemento_ErroresDeEditarYEliminar_MuestranLaSeccionDeEsfuerzo` |
| A | `template.HTML` en el atributo `value` | Falla `TestPantallaEsfuerzo_LoTipeado_SeEscapa` |
| B | `template.HTMLAttr` con el valor crudo | Falla `TestPantallaEsfuerzo_LoTipeado_SeEscapa` |
| C | `template.HTML` impreso como texto | Falla `TestPantallaEsfuerzo_LoTipeado_SeEscapa` |
| D | Reflejar `?esfuerzo=` con `template.HTML` | Falla `TestPantallaEsfuerzo_ParametroDeConfirmacion_NoInyecta` |
| E | Tabla de registros en la sección | Falla `TestPantallaEsfuerzo_FormularioInicial` |
| F | Lista de registros en la sección | Falla `TestPantallaEsfuerzo_FormularioInicial` |
| G | Texto "Horas acumuladas: 3,5 h" | Falla `TestPantallaEsfuerzo_FormularioInicial` |
| H | Texto "Llevás 2 horas" | Falla `TestPantallaEsfuerzo_FormularioInicial` |
| I | Texto "Los registros suman al total del elemento." | Pasa, como se buscaba: no es un listado ni un total |

La mutación A falla solo por el `value` esperado: `html/template` quita las etiquetas en ese atributo y no llega a salir un `<script>`. Las mutaciones B, C y D son las que inyectan HTML.

### Huecos

- La spec sigue en **Borrador**. **P-6 (la ruta HTTP) está pendiente de confirmar con Matías**: si cambia, hay que tocar `esfuerzo_handler.go`, `router_test.go`, CA-18.15 y esta matriz.
- **P-14:** como la pantalla no muestra el total, un doble registro no se detecta desde la interfaz, y tampoco se corrige, porque no se pueden editar ni eliminar registros.
- **P-15** sigue abierta en la spec, aunque el código ya hace lo que propone: `errorDeElemento` reconoce `ErrProyectoNoEncontrado`.
- Accesibilidad: `role="status"` y `role="alert"` vienen en el HTML inicial de la página, y esas regiones no siempre se anuncian. `error_form` está en `layout.html`, compartido con HU-30, y no se tocó.
- La pantalla no se probó en un navegador ni con un lector de pantalla: se probó con tests y con curl contra el servidor. Falta ver el salto a `#esfuerzo` y cómo se ven la sección y `.exito`.
- Los tests de zona horaria solo corrieron en UTC-3, porque Go ignora `TZ` en Windows. La independencia real de la zona se verá en el CI.
- `sinListadoNiTotal` no detecta un total sin unidades, como "Total: 3,5", y solo se aplica al formulario inicial (CA-18.16).
- Refactors que tocaron código de otras historias: `formMensajeDeErrorConocido` reemplazó el cuerpo de `mensajeDeValidacion` (HU-30), `mensajeDeValidacionElemento` (HU-06) y `mensajeDeValidacionDefecto` (HU-20), y `mostrarElemento` y `errorDeElemento` (HU-06) cambiaron para la pantalla. Conviene que lo revisen sus autores en el PR.
- `mostrarErrorDeEsfuerzo` consulta el proyecto y `mostrarElemento` lo vuelve a consultar: es una lectura de más, solo cuando el formulario vuelve con un error.
- El merge `cc1a242` con `main` (HU-14 y HU-15) tuvo conflictos en `elemento.html`, `pantallas_elemento.go` y `pantallas.go`, y además dos choques que git no detectó, porque no había líneas en conflicto:
  - `mensajeDeErrorConocido` quedó definida dos veces y con distinto comportamiento: la de HU-18 también reconoce `errFormulario`; la de HU-14 solo mira su lista. La de HU-18 se renombró a `formMensajeDeErrorConocido` y la de HU-14 no cambió. Conviven dos nombres muy parecidos, que se podrían unificar en otro cambio.
  - `iniciarRondaPantalla` (HU-14) llamaba a `mostrarElemento` con la firma anterior, sin el parámetro de esfuerzo. Se le agregó `p.seccionEsfuerzoInicial()`: el `400` de HU-14 sigue igual, y ahora esa pantalla también muestra la sección de esfuerzo.
- `067e249` no tuvo RED real: `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` se escribió con el código ya funcionando. La evidencia son 5 mutaciones manuales sobre `elemento.html`, que no quedaron en ningún commit. Con cada una el test falla, y sin ellas pasa:
  - quitar la sección de Planning Poker: falla porque falta `section titulo-poker` y su `<h2>`;
  - quitar la sección de esfuerzo: falla porque falta `section esfuerzo` y su `<h2>`;
  - invertir el orden (esfuerzo antes que Planning Poker): falla la estructura esperada;
  - duplicar un id (`id="horas"` también en la actividad): falla con `ids duplicados: ["horas"]`;
  - anidar un formulario dentro del de esfuerzo: falla porque aparece `form …/esfuerzo#esfuerzo > form /x`.
- CA-18.15 no tiene escenario BDD a propósito: los escenarios prueban el núcleo y no pasan por HTTP. Lo cubren los tests del handler.
- CA-18.11 no tiene test unitario propio. El dominio no suma registros porque valida uno solo, así que no hay lógica que probar aparte. La regla se verifica en el escenario BDD, que guarda 20 h + 8 h en el mismo día contra SQLite.
- CA-18.14: el test unitario cubre la fecha anterior al inicio del proyecto, pero lo de "proyecto no Activo" solo se verifica en el BDD. El caso de uso no lee el estado, así que no hay rama que probar.
- En 18.4, 18.5 (primer ciclo), 18.6 (primer ciclo) y 18.7 (handler), el RED es de compilación: el test usa código que todavía no existe. Los demás RED fallan al ejecutarse.
- Hay commits sin etiqueta de etapa: `c076181`, `7dc1955`, `bb9f3f6`, `c872259` y `c4864d2`, este último un refactor sin `[REFACTOR]`. Por eso un `git log --grep="\[REFACTOR\]"` no los encuentra.
- El mensaje de `929d265` dice "cuenta registros": `ContarRegistros` después pasó a helpers de test (`c4864d2`), así que el repositorio de producción ya no cuenta.
- Hallazgos del go-reviewer que no se corrigieron todavía y quedan para hablar con el equipo:
  - el cuerpo JSON acepta texto después del objeto, igual que en HU-01 y HU-05;
  - los hooks de BDD de cada historia corren en todos los escenarios;
  - las claves foráneas de `registros_esfuerzo` no tienen `ON DELETE`, y eso va a importar en HU-06;
  - el nombre `RegistrarRegistrarEsfuerzo`.

---

## HU-30: Interfaz web base

- **Issue:** [#103](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/103) · Sprint 2 · 5 SP · ADR-0003
- **Spec:** [docs/specs/HU-30-interfaz-web-base.md](specs/HU-30-interfaz-web-base.md)
- **Feature:** [features/HU-30-interfaz-web-base.feature](../features/HU-30-interfaz-web-base.feature) (steps en [test/bdd/hu30_steps_test.go](../test/bdd/hu30_steps_test.go); los escenarios hacen pedidos HTTP reales a las pantallas)

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-30.1 Layout común y estilos | Restricciones | Las pantallas comparten el layout común `@normal` | `TestPantallas_Lista_TieneLayoutYProyectos`, `TestPantallas_EstilosSeSirven`, `TestNuevoRouter_RegistraTodasLasRutas` | `web/plantillas/layout.html`, `web/static/estilos.css`, `RegistrarPantallas` |
| CA-30.2 Lista ordenada por nombre / vacía | RN-30.1 | La lista muestra los proyectos ordenados por nombre `@normal` · Sin proyectos se muestra la lista vacía `@limite` | `TestListarProyectos_OrdenaPorNombreSinDistinguirMayusculas`, `TestListarProyectos_AIgualNombreOrdenaPorID`, `TestListarProyectos_SinProyectos_DevuelveListaVacia`, `TestRepositorioProyectos_ListarProyectos_*`, `TestPantallas_ListaVacia` | `app.ListarProyectos`, `persistencia.RepositorioProyectos.ListarProyectos`, `web/plantillas/proyectos.html` |
| CA-30.3 Crear desde el formulario | RN-30.2 | Crear un proyecto desde el formulario `@normal` · El formulario muestra el error de dominio `@error` · Fecha con formato inválido `@error` | `TestPantallas_NuevoProyecto` | `pantallas.crear`, `web/plantillas/proyecto_form.html` |
| CA-30.4 Detalle con integrantes / 404 | RN-30.2 | El detalle muestra los datos y sus integrantes `@normal` · Proyecto inexistente es una página 404 `@error` | `TestPantallas_Detalle` | `pantallas.detalle`, `web/plantillas/proyecto_detalle.html` |
| CA-30.5 Editar, quitar fecha de fin | RN-30.2, RN-30.3 | Editar un proyecto y quitarle la fecha de fin `@normal` · El formulario de edición muestra el error `@error` | `TestPantallas_EditarProyecto` | `pantallas.editar`, `cambiosDeForm` |
| CA-30.6 Agregar integrante | RN-30.2 | Agregar un integrante desde el detalle `@normal` · Integrante repetido `@error` | `TestPantallas_Integrantes` | `pantallas.agregarIntegrante` |
| CA-30.7 Quitar integrante | RN-30.2 | Quitar un integrante desde el detalle `@normal` | `TestPantallas_Integrantes`, `TestPantallas_NombreConCaracteresReservados` | `pantallas.quitarIntegrante` |
| CA-30.8 Sumar rol | RN-30.2 | Sumar un rol `@normal` · Rol repetido `@error` | `TestPantallas_Integrantes` | `pantallas.agregarRol` |
| CA-30.9 Contenido escapado | RN-30.4 | El contenido del usuario se muestra escapado `@limite` | `TestPantallas_EscapaElContenido` | `html/template`, función `segmento` |
| CA-30.10 Inicio redirige a la lista | — | La página de inicio lleva a la lista de proyectos `@normal` | `TestPantallas_Inicio_RedirigeALaLista` | `RegistrarPantallas` (`GET /{$}`) |

Además: `TestPantallas_RechazaFormulariosDeOtroSitio` (403, CSRF), `TestPantallas_FormularioDemasiadoGrande` (413), `TestPantallas_CabecerasDeSeguridad` y `TestPantallas_FallaInterna_Devuelve500SinDetalles` (sección 7 de la spec).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Caso de uso | [internal/app/listar_proyectos.go](../internal/app/listar_proyectos.go) | [listar_proyectos_test.go](../internal/app/listar_proyectos_test.go) |
| Persistencia | [internal/adapters/persistencia/proyecto_repo.go](../internal/adapters/persistencia/proyecto_repo.go) (`ListarProyectos`) | [listar_proyectos_test.go](../internal/adapters/persistencia/listar_proyectos_test.go) |
| Pantallas | [internal/adapters/http/pantallas.go](../internal/adapters/http/pantallas.go), [web/plantillas/](../internal/adapters/http/web/plantillas/), [web/static/estilos.css](../internal/adapters/http/web/static/estilos.css) | [pantallas_test.go](../internal/adapters/http/pantallas_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| BDD | [test/bdd/hu30_sistema_test.go](../test/bdd/hu30_sistema_test.go) (servidor HTTP de prueba + SQLite) | 17 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 30.1 Spec | `4089bbd` docs(HU-30): [30.1] spec | — |
| 30.2 Escenarios BDD | `a067092` test(HU-30): [30.2] escenarios BDD | — |
| 30.3 Steps godog | `becbe6d` test(HU-30): [30.3] steps (en rojo) | BDD en rojo |
| 30.4 Caso de uso ListarProyectos | `28cddd3` [RED] → `1e69c44` [GREEN] | **RED → GREEN** |
| 30.5 Persistencia | `74230a8` [RED] → `fe2bd86` [GREEN] | **RED → GREEN** |
| 30.6 Pantallas | `6320169` [RED] → `488bcb4` [GREEN] | **RED → GREEN** |
| 30.7 main y BDD | `ef284be` feat(HU-30): [30.7] BDD en verde | BDD en verde |
| 30.8 Revisión | `29fbdd5` refactor(HU-30): [30.8] correcciones go-reviewer | **REFACTOR** |
| 30.9 Trazabilidad | este commit | — |

### Huecos

- La spec está en **Borrador**: las decisiones (prefijo `/web`, orden por nombre, POST con redirección 303, CSS propio sin JavaScript) se validan en el PR.
- Las pantallas de las demás historias (HU-05 en adelante) no están incluidas: cada historia suma la suya sobre este layout.
- Los tests de los arreglos de la revisión (30.8) se escribieron junto con el arreglo, no en un commit RED aparte.

---

## HU-06: Editar o eliminar un elemento del backlog

- **Issue:** [#17](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/17) · Sprint 2 · 3 SP · tareas #123 a #131
- **Spec:** [docs/specs/HU-06-editar-eliminar-elemento-backlog.md](specs/HU-06-editar-eliminar-elemento-backlog.md)
- **Feature:** [features/HU-06-editar-eliminar-elemento-backlog.feature](../features/HU-06-editar-eliminar-elemento-backlog.feature) (steps en [test/bdd/hu06_steps_test.go](../test/bdd/hu06_steps_test.go); los escenarios de CA-06.10 usan la pantalla real)

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-06.1 Solo la prioridad | RN-06.1 | Cambiar solo la prioridad `@normal` | `TestEditar_SoloPrioridad_ElRestoNoCambia`, `TestEditarElemento_DatosValidos_GuardaSoloLoEnviadoYDevuelveLoGuardado`, `TestPatchElemento_Valido_Devuelve200YSoloEnviaLoRecibido` | `domain.ElementoBacklog.Editar`, `app.EditarElementoBacklog` |
| CA-06.2 Varios campos y estado | RN-06.1, RN-06.4, RN-06.9 | Editar varios campos juntos `@normal` · Un elemento Hecho puede volver a Pendiente `@limite` | `TestEditar_TodosLosCampos`, `TestEditar_EstadoEnCualquierSentido`, `TestRepositorioBacklog_ActualizarElemento_TodosLosCampos` | `domain.ElementoBacklog.Editar`, `persistencia.RepositorioBacklog.ActualizarElemento` |
| CA-06.3 Quitar la estimación | RN-06.3 | Quitar la estimación `@alternativo` | `TestEditar_StoryPointsNil_QuedaSinEstimar`, `TestPatchElemento_StoryPoints`, `TestPantallaElemento_StoryPointsVacioQuitaLaEstimacion` | `domain.ElementoBacklog.Editar`, `pedidoEditarElemento.aCambios` |
| CA-06.4 Datos inválidos | RN-06.2 | Los datos editados cumplen las reglas de HU-05 `@error` (×7) | `TestEditar_DatosInvalidos_RetornaError`, `TestEditarElemento_DatosInvalidos_NoActualiza`, `TestPatchElemento_Errores` | `domain.ElementoBacklog.Editar` (reutiliza `NuevoElementoBacklog`) |
| CA-06.5 Título duplicado | RN-06.2 | No se puede renombrar a un título que ya existe `@error` · Renombrar con otras mayúsculas `@alternativo` | `TestEditarElemento_TituloDuplicado`, `TestRepositorioBacklog_ActualizarElemento_Errores` | `app.EditarElementoBacklog`, UNIQUE `(proyecto_id, titulo_clave)` |
| CA-06.6 Sin cambios | RN-06.5 | Editar sin enviar ningún campo `@error` | `TestEditar_DatosInvalidos_RetornaError` (sin cambios), `TestPatchElemento_Errores` | `domain.ElementoBacklog.Editar` |
| CA-06.7 Otro proyecto / inexistente | RN-06.6 | No se edita ni se elimina un elemento de otro proyecto `@error` · No se elimina un elemento inexistente `@error` | `TestEditarElemento_DeOtroProyectoOInexistente_NoEncontrado`, `TestEliminarElemento_DeOtroProyectoOInexistente_NoEncontrado`, `TestConsultarElemento_DeOtroProyectoOInexistente` | `buscarElementoDelProyecto` |
| CA-06.8 Eliminar sin esfuerzo | RN-06.7 | Eliminar un elemento sin esfuerzo `@normal` | `TestEliminarElemento_SinEsfuerzo_LoElimina`, `TestRepositorioBacklog_EliminarElemento_BorraElElementoYSusCriterios`, `TestDeleteElemento` | `app.EliminarElementoBacklog`, `persistencia.RepositorioBacklog.EliminarElemento` |
| CA-06.9 Con esfuerzo no se elimina | RN-06.7 | No se elimina un elemento con esfuerzo registrado `@error` | `TestEliminarElemento_ConEsfuerzo_NoLoElimina`, `TestRepositorioBacklog_ElementoConEsfuerzo`, `TestDeleteElemento` (409) | `app.EliminarElementoBacklog`, clave foránea de `registros_esfuerzo` |
| CA-06.10 Pantalla | Restricciones | Editar desde la pantalla guarda solo lo que cambió `@normal` · La pantalla muestra el error `@error` · Guardar sin cambios vuelve sin error `@alternativo` · Eliminar desde la pantalla `@normal` | `TestPantallaElemento_*` (9 tests), `TestPantallaBacklog_EnlazaAEditarYEliminarCadaElemento`, `TestConsultarElemento_DelProyecto` | `pantallas_elemento.go`, `web/plantillas/elemento.html`, botones en `web/plantillas/backlog.html` (HU-07), `app.ConsultarElementoBacklog` |
| CA-06.11 Con defectos no se elimina (DEF-02) | RN-06.10 | No se elimina un elemento con un defecto registrado `@error` `@DEF-02` | `TestEliminarElemento_ConDefectos_NoLoElimina`, `TestEliminarElemento_ConEsfuerzoYDefectos_InformaElEsfuerzo`, `TestRepositorioBacklog_ElementoConDefectos`, `TestDeleteElemento` (con defectos), `TestPantallaElemento_Errores` (eliminar con defectos) | `domain.ErrElementoConDefectos`, `app.EliminarElementoBacklog.Ejecutar`, `persistencia.RepositorioBacklog.TieneDefectos`, `persistencia.RepositorioBacklog.errorDeEliminacion` |

RN-06.8 (no eliminar un elemento en un Sprint activo) queda pendiente: la implementa HU-09 en `EliminarElementoBacklog`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog_editar.go](../internal/domain/backlog_editar.go) | [backlog_editar_test.go](../internal/domain/backlog_editar_test.go) |
| Casos de uso | [internal/app/editar_elemento_backlog.go](../internal/app/editar_elemento_backlog.go), [consultar_elemento_backlog.go](../internal/app/consultar_elemento_backlog.go) | [editar_elemento_backlog_test.go](../internal/app/editar_elemento_backlog_test.go), [consultar_elemento_backlog_test.go](../internal/app/consultar_elemento_backlog_test.go) |
| Persistencia | [internal/adapters/persistencia/backlog_editar.go](../internal/adapters/persistencia/backlog_editar.go) | [backlog_editar_test.go](../internal/adapters/persistencia/backlog_editar_test.go) |
| HTTP (JSON) | [internal/adapters/http/backlog_editar_handler.go](../internal/adapters/http/backlog_editar_handler.go) (`PATCH`/`DELETE /proyectos/{id}/backlog/{elementoID}`) | [backlog_editar_handler_test.go](../internal/adapters/http/backlog_editar_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_elemento.go](../internal/adapters/http/pantallas_elemento.go), [web/plantillas/elemento.html](../internal/adapters/http/web/plantillas/elemento.html) (`/web/proyectos/{id}/backlog/{elementoID}`) | [pantallas_elemento_test.go](../internal/adapters/http/pantallas_elemento_test.go) |
| BDD | [test/bdd/hu06_sistema_test.go](../test/bdd/hu06_sistema_test.go) | 22 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 6.1 Spec | `790172d` | — |
| 6.2 Escenarios BDD | `315f3b8` | — |
| 6.3 Steps godog | `2330d68` | BDD en rojo |
| 6.4 Edición en el dominio | `78c6573`, `8183241` [RED] → `0925cd6` [GREEN] | **RED → GREEN** |
| 6.5 Casos de uso | `b4692bf` [RED] → `0555f28` [GREEN] | **RED → GREEN** |
| 6.6 Persistencia | `4321ddb` [RED] → `11283dd` [GREEN] | **RED → GREEN** |
| 6.7 HTTP | `23afd32` [RED] → `f1e086b` [GREEN] | **RED → GREEN** |
| 6.8 Pantalla | `422994e` [RED] → `be127d1` [GREEN] | **RED → GREEN**, BDD en verde |
| 6.9 Revisión y trazabilidad | `a22d553` | **REFACTOR** |
| Guardar sin cambios (alineado con DEF-01) | `b1caa63` [RED] → `d7de2cc` [GREEN] | **RED → GREEN** |
| Editar y eliminar desde la lista del backlog (HU-07) | `01ae24f` [RED] → `72b4724` [GREEN] | **RED → GREEN** |
| Revisión del PR #182: rutas POST de la pantalla en `router_test.go` | `f478777` | test |
| Revisión del PR #182: sin campos `original_*` se envían todos los campos | `9384fca` [RED] → `ba3863c` [GREEN] | **RED → GREEN** |

### Huecos

- La spec está en **Borrador** hasta la aprobación del PR.
- RN-06.8 (Sprint activo) depende de HU-09.
- La lista del backlog (HU-07) enlaza a la pantalla del elemento (Editar) y permite eliminar con confirmación.

---

## HU-09: Asignar historias del backlog a un Sprint

- **Issue:** [#20](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/20) · Sprint 2 · 5 SP
- **Spec:** [docs/specs/HU-09-asignar-historias-sprint.md](specs/HU-09-asignar-historias-sprint.md)
- **Feature:** [features/HU-09-asignar-historias-sprint.feature](../features/HU-09-asignar-historias-sprint.feature) (steps en [test/bdd/hu09_steps_test.go](../test/bdd/hu09_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-09.1 Asignar a un sprint abierto | RN-09.1, RN-09.3 | Asignar un elemento sin sprint a un sprint abierto `@normal` · Reasignar un elemento al mismo sprint que ya tiene no es un error `@limite` | `TestElementoBacklog_AsignarASprint_SinSprintPrevio_Asigna`, `TestElementoBacklog_AsignarASprint_MismoSprintQueYaTiene_NoEsError`, `TestAsignarHistoria_DatosValidos_GuardaElSprintID`, `TestRepositorioBacklog_GuardarAsignacion_Asigna`, `TestPostAsignacion_DatosValidos_Devuelve204`, `TestPantallaSprintDetalle_Asignar_RedirigeAlDetalle` | `domain.ElementoBacklog.AsignarASprint`, `app.AsignarHistoria`, `http.handlerAsignarHistoria`, `http.pantallas.asignarElementoPantalla` |
| CA-09.2 Ya asignado a otro sprint activo | RN-09.2 | No se puede asignar un elemento que ya está en otro sprint abierto `@error` | `TestElementoBacklog_AsignarASprint_YaAsignadoAOtroSprintActivo_RetornaError`, `TestAsignarHistoria_YaAsignadoAOtroSprintActivo_NoGuarda`, `TestPostAsignacion_DatosInvalidos_Devuelve400`, `TestAsignacion_DosAsignacionesSimultaneas_LaSegundaFalla` | `domain.ElementoBacklog.AsignarASprint`, `persistencia.RepositorioBacklog.GuardarAsignacion` (verificado en el propio `UPDATE`) |
| CA-09.3 Reasignar desde un sprint ya cerrado | RN-09.2 | Se puede reasignar un elemento desde un sprint ya cerrado `@alternativo` | `TestElementoBacklog_AsignarASprint_SprintPrevioYaCerrado_Reasigna`, `TestAsignarHistoria_SprintPrevioCerrado_Reasigna`, `TestRepositorioBacklog_GuardarAsignacion_ReasignaDesdeSprintCerrado` | `domain.ElementoBacklog.AsignarASprint` |
| CA-09.4 Desasignar de un sprint abierto | RN-09.4 | Desasignar un elemento de un sprint abierto `@normal` | `TestElementoBacklog_DesasignarDeSprint_SprintActivo_Desasigna`, `TestDesasignarHistoria_SprintActivo_Desasigna`, `TestRepositorioBacklog_GuardarAsignacion_Desasigna`, `TestPostAsignacion` (DELETE), `TestPantallaSprintDetalle_Quitar_RedirigeAlDetalle` | `domain.ElementoBacklog.DesasignarDeSprint`, `app.DesasignarHistoria` |
| CA-09.5 No desasignar de un sprint cerrado | RN-09.4 | No se puede desasignar un elemento de un sprint ya cerrado `@error` | `TestElementoBacklog_DesasignarDeSprint_SprintYaCerrado_RetornaError`, `TestDesasignarHistoria_SprintCerrado_RetornaErrorYNoGuarda`, `TestRepositorioBacklog_GuardarAsignacion_DesasignarConSprintYaCerrado_Falla` | `domain.ElementoBacklog.DesasignarDeSprint`, `persistencia.RepositorioBacklog.GuardarAsignacion` (verificado en el propio `UPDATE`) |
| CA-09.6 No desasignar sin sprint | — (error de dominio) | No se puede desasignar un elemento que no tiene sprint `@error` | `TestElementoBacklog_DesasignarDeSprint_SinSprintAsignado_RetornaError`, `TestDesasignarHistoria_SinSprintAsignado_RetornaError` | `domain.ElementoBacklog.DesasignarDeSprint` |
| CA-09.7 Sprint de otro proyecto | RN-09.5 | No se puede asignar un elemento a un sprint de otro proyecto `@error` | `TestElementoBacklog_AsignarASprint_SprintDeOtroProyecto_RetornaError`, `TestAsignarHistoria_SprintDeOtroProyecto_NoGuarda` | `domain.ElementoBacklog.AsignarASprint` |
| CA-09.8 No asignar a un sprint cerrado | RN-09.6 | No se puede asignar un elemento a un sprint ya cerrado `@error` | `TestElementoBacklog_AsignarASprint_SprintDestinoCerrado_RetornaError`, `TestAsignarHistoria_SprintDestinoCerrado_NoGuarda` | `domain.ElementoBacklog.AsignarASprint` |

RN-09.7 (un elemento de otro proyecto se trata como inexistente, igual que RN-18.1) no tiene CA propio: lo verifican `TestAsignarHistoria_ElementoDeOtroProyecto_RetornaError` y `TestDesasignarHistoria_ElementoDeOtroProyecto_RetornaError`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog_asignacion.go](../internal/domain/backlog_asignacion.go) (+ campo `SprintID` en [backlog.go](../internal/domain/backlog.go)) | [backlog_asignacion_test.go](../internal/domain/backlog_asignacion_test.go) |
| Caso de uso | [internal/app/asignar_historia.go](../internal/app/asignar_historia.go) | [asignar_historia_test.go](../internal/app/asignar_historia_test.go) |
| Persistencia | [internal/adapters/persistencia/asignacion_repo.go](../internal/adapters/persistencia/asignacion_repo.go) (`BuscarSprint`, `GuardarAsignacion` atómico), [backlog_repo.go](../internal/adapters/persistencia/backlog_repo.go) (`BuscarElementoPorID`/`ListarElementos` leen `sprint_id`), [db.go](../internal/adapters/persistencia/db.go) (columna `sprint_id`, agregada con `ALTER TABLE` si falta) | [asignacion_repo_test.go](../internal/adapters/persistencia/asignacion_repo_test.go), [asignacion_concurrencia_test.go](../internal/adapters/persistencia/asignacion_concurrencia_test.go), [db_test.go](../internal/adapters/persistencia/db_test.go) |
| HTTP | [internal/adapters/http/asignacion_handler.go](../internal/adapters/http/asignacion_handler.go) (`POST`/`DELETE /proyectos/{id}/backlog/{elementoID}/sprint`) | [asignacion_handler_test.go](../internal/adapters/http/asignacion_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go) (`/web/proyectos/{id}/sprints/{sprintID}`), [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | [pantallas_asignacion_test.go](../internal/adapters/http/pantallas_asignacion_test.go) |
| BDD | [test/bdd/hu09_sistema_test.go](../test/bdd/hu09_sistema_test.go), [hu09_steps_test.go](../test/bdd/hu09_steps_test.go) | 8 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 9.1 Spec | `b02fcc1` docs(HU-09): [9.1] spec | — |
| 9.2 Escenarios BDD | `cd2050f` test(HU-09): [9.2] escenarios BDD | — |
| 9.3 Steps godog | `f4e7cc3` test(HU-09): [9.3] steps (en rojo) | BDD en rojo |
| 9.4 Reglas de asignación | `ef3020e` [RED] → `22f7fb6` [GREEN] | **RED → GREEN** |
| 9.5 Caso de uso | `4b858de` [RED] → `97a836c` [GREEN] → `d704371` [REFACTOR] | **RED → GREEN → REFACTOR** |
| 9.5 Elemento de otro proyecto (RN-09.7) | `1b754e4` [RED] → `d997df1` [GREEN] | **RED → GREEN** |
| 9.6 Persistencia | `333d0d0` [RED] → `cf64285` [GREEN] | **RED → GREEN** |
| 9.7 HTTP | `86ddb06` [RED] → `616de9b` [GREEN] | **RED → GREEN** |
| 9.8 Pantalla | `7443720` [RED] → `bcce7ea` [GREEN] | **RED → GREEN** |
| 9.8 Fix: ListarElementos no reflejaba el sprint (hallado al verificar la pantalla en el navegador) | `135f945` [RED] → `652c214` [GREEN] | **RED → GREEN** |
| 9.7 Fix: `sprint_id` sin validar en el endpoint JSON (hallado por go-reviewer en 9.9) | `0e974fc` [RED] → `9d3e8c6` [GREEN] | **RED → GREEN** |
| 9.9 Revisión (reviewer humano): migración rompía bases existentes | `fcdb07a` [RED] → `ed0597e` [GREEN] | **RED → GREEN** |
| 9.9 Revisión (reviewer humano): asignaciones simultáneas se pisaban | `e628a16` [RED] → `d974b75` [GREEN] | **RED → GREEN** |
| 9.9 Revisión + trazabilidad | este commit | — |

### Huecos

- La spec está en **Borrador**: pasa a Aprobada cuando se apruebe el PR.
- La pantalla no permite "robarle" a un elemento su sprint activo desde la UI (solo lista elementos sin sprint asignado en el `<select>`); si ese otro sprint ya está Cerrado, sí se puede reasignar por la API aunque la pantalla no lo ofrezca todavía como atajo directo.
- Hallazgo del go-reviewer que no se corrigió todavía y queda para hablar con el equipo: `RepositorioBacklog` (`asignacion_repo.go`) consulta la tabla `sprints` directamente en vez de delegar en `RepositorioSprints` (de HU-08), acoplando el repositorio de backlog a una tabla de otro agregado. Estaba justificado cuando se escribió (HU-08 todavía no estaba mergeada), pero ahora que ya está en `main` convendría refactorizarlo para no duplicar el acceso a `sprints` en dos repositorios.

---

## HU-10: Marcar una historia como completada en el Sprint

- **Issue:** [#21](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/21) · Sprint 2 · 2 SP
- **Spec:** [docs/specs/HU-10-completar-historia.md](specs/HU-10-completar-historia.md) (Aprobada; P-9 tomada por el autor, abierta a objeción en el PR)
- **Feature:** [features/HU-10-completar-historia.feature](../features/HU-10-completar-historia.feature) (steps en [test/bdd/hu10_steps_test.go](../test/bdd/hu10_steps_test.go); los escenarios usan el caso de uso real con SQLite)
- **Rutas HTTP:** `POST /proyectos/{id}/backlog/{elementoID}/completar` (JSON, `204`) y `POST /web/proyectos/{id}/sprints/{sprintID}/completar` (botón "Marcar como hecha").

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-10.1 Pendiente o En progreso en sprint Abierto → Hecho, sigue en el sprint | RN-10.1, RN-10.2 | Una historia de un sprint abierto pasa a Hecho y sigue en el sprint `@normal` (Pendiente) `@alternativo` (En progreso) | `TestElementoBacklog_Completar_Exito` ("desde Pendiente", "desde En progreso"), `TestCompletarHistoria_SprintAbierto_GuardaComoHecho`, `TestRepositorioBacklog_GuardarCompletado_GuardaHecho`, `TestRepositorioBacklog_GuardarCompletado_NoCambiaOtrosCampos`, `TestPantallaSprintDetalle_Completar_RedirigeAlDetalle` | `domain.ElementoBacklog.Completar`, `app.CompletarHistoria.Ejecutar`, `persistencia.RepositorioBacklog.GuardarCompletado` |
| CA-10.2 Ya Hecho en sprint Abierto: sin error, sigue Hecho | RN-10.3 | Completar una historia que ya está Hecho no es un error `@limite` | `TestElementoBacklog_Completar_Exito` ("ya Hecho no cambia nada"), `TestCompletarHistoria_SprintAbierto_GuardaComoHecho` ("ya Hecho, sin error"), `TestRepositorioBacklog_GuardarCompletado_YaHecho_SinError` | `domain.ElementoBacklog.Completar`, `persistencia.RepositorioBacklog.GuardarCompletado` |
| CA-10.3 Sin sprint → `ErrElementoSinSprint`, el estado no cambia | RN-10.2 | No se puede completar una historia que no está en ningún sprint `@error` | `TestElementoBacklog_Completar_Errores_NoCambiaElElemento` ("sin sprint, …"), `TestCompletarHistoria_SinSprint_RetornaErrorYNoGuarda`, `TestRepositorioBacklog_GuardarCompletado_SinSprint_RetornaError`, `TestPostCompletar_ErroresDelCasoDeUso` ("elemento sin sprint", 400), `TestPantallaSprintDetalle_ErrorAlCompletar` ("elemento sin sprint") | `domain.ElementoBacklog.Completar`, `app.sprintDelElementoAbierto`, `persistencia.RepositorioBacklog.errorDeCompletado` |
| CA-10.4 Sprint Cerrado → `ErrSprintCerrado`, el estado no cambia (también si ya está Hecho, P-8) | RN-10.2 | No se puede completar una historia de un sprint cerrado `@error` · Una historia ya Hecho de un sprint cerrado también se rechaza `@limite` | `TestElementoBacklog_Completar_Errores_NoCambiaElElemento` ("sprint cerrado, …", incluido "Hecho (P-8)"), `TestCompletarHistoria_SprintCerrado_RetornaErrorYNoGuarda`, `TestRepositorioBacklog_GuardarCompletado_SprintCerrado_RetornaError`, `TestCompletar_SprintCerradoEntreLecturaYGuardado_Falla`, `TestPostCompletar_ErroresDelCasoDeUso` ("sprint cerrado", 400), `TestPantallaSprintDetalle_ErrorAlCompletar` ("sprint cerrado") | `domain.ElementoBacklog.Completar`, `persistencia.RepositorioBacklog.GuardarCompletado` (verificado en el propio `UPDATE`) |
| CA-10.5 Proyecto o elemento inexistente, o elemento de otro proyecto | RN-10.4 | No se puede completar una historia de un proyecto que no existe `@error` · … una historia que no existe `@error` · Una historia de otro proyecto se trata como inexistente `@error` | `TestCompletarHistoria_ProyectoOElementoInexistente_RetornaErrorYNoGuarda` (incluye "elemento de otro proyecto (RN-10.4)"), `TestRepositorioBacklog_GuardarCompletado_Inexistente_RetornaError`, `TestPostCompletar_ErroresDelCasoDeUso` ("proyecto inexistente", "elemento inexistente", 404), `TestPantallaSprintDetalle_ErrorAlCompletar` ("proyecto inexistente", "elemento inexistente") | `app.CompletarHistoria.Ejecutar` |
| CA-10.6 Respuestas HTTP 204 / 404 / 400 | — (sección 3) | **No aplica, se prueba en el adaptador:** los escenarios prueban el núcleo y no pasan por HTTP. | `TestPostCompletar_PedidoValido_Devuelve204`, `TestPostCompletar_ErroresDelCasoDeUso`, `TestPostCompletar_IDsInvalidos_Devuelve400SinLlamarAlCasoDeUso`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerCompletarHistoria`, `http.RegistrarCompletar`, `http.NuevoRouter`, `cmd/server/main.go` |
| CA-10.7 Pantalla: estado, botón "Marcar como hecha", errores | RN-10.1, RN-10.2, sección 5 | **No aplica, se prueba en el adaptador** (la pantalla no tiene escenario BDD). | `TestPantallaSprintDetalle_BotonCompletar_SprintAbiertoYNoHecho`, `TestPantallaSprintDetalle_SinBotonCompletar`, `TestPantallaSprintDetalle_Completar_RedirigeAlDetalle`, `TestPantallaSprintDetalle_Completar_ElElementoApareceHecho`, `TestPantallaSprintDetalle_ErrorAlCompletar`, `TestMensajeDeValidacion_ErroresDeCompletar`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.pantallas.completarElementoPantalla`, `PuedeCompletarse` en `pantallas_asignacion.go`, botón en `web/plantillas/sprint_detalle.html` |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-10.1 Completar = pasar a `Hecho`, desde Pendiente o En progreso | CA-10.1, CA-10.7 | `domain.ElementoBacklog.Completar`; `GuardarCompletado` cambia solo el estado (`TestRepositorioBacklog_GuardarCompletado_NoCambiaOtrosCampos`) |
| RN-10.2 Solo en un sprint asignado y Abierto | CA-10.1, CA-10.3, CA-10.4, CA-10.7 | `domain.ElementoBacklog.Completar` (recibe `sprintAbierto bool`), `app.sprintDelElementoAbierto` (compartida con `DesasignarHistoria` de HU-09), `persistencia.RepositorioBacklog.GuardarCompletado` (en el propio `UPDATE`) |
| RN-10.3 Ya `Hecho` no es error ni cambia nada | CA-10.2 | `domain.ElementoBacklog.Completar` |
| RN-10.4 Elemento de otro proyecto = inexistente | CA-10.5 | `app.CompletarHistoria.Ejecutar`. Se prueba en el caso de uso y en BDD; el repositorio no distingue el proyecto. |

Sin CA propio en la spec:
- Orden de los errores (sección 7): `TestCompletarHistoria_ProyectoOElementoInexistente_RetornaErrorYNoGuarda` ("proyecto y elemento inexistentes: primero el proyecto") y, para RN-10.2 antes que RN-10.3 (P-8), el subtest "sprint cerrado, Hecho" del dominio y el escenario `@limite` de CA-10.4.
- Pedidos simultáneos (sección 5, P-7): `TestCompletar_SprintCerradoEntreLecturaYGuardado_Falla` (el sprint se cierra en el medio → `ErrSprintCerrado`).
- **P-9**, elemento movido a otro sprint Abierto durante el pedido (sección 6): `TestRepositorioBacklog_GuardarCompletado_MovidoAOtroSprint_RetornaError` y `TestCompletar_MovidoAOtroSprintEntreLecturaYGuardado_Falla` → `ErrElementoSinSprint`, sin completarlo en el sprint nuevo. Lo decide `persistencia.RepositorioBacklog.errorDeCompletado`. No tiene escenario BDD.
- IDs de la ruta inválidos (sección 2): `TestPostCompletar_IDsInvalidos_Devuelve400SinLlamarAlCasoDeUso`; en la pantalla, el subtest "elemento_id inválido" de `TestPantallaSprintDetalle_ErrorAlCompletar`.
- Falla interna: `TestCompletarHistoria_FallaDelRepositorio_SePropagaYNoGuarda`, `TestPostCompletar_FallaInterna_Devuelve500SinDetalles` y el subtest "falla interna" de `TestPantallaSprintDetalle_ErrorAlCompletar`.

HU-10 no define errores nuevos: reutiliza `ErrProyectoNoEncontrado` (HU-01), `ErrElementoNoEncontrado` (HU-05), `ErrElementoSinSprint` y `ErrSprintCerrado` (HU-09).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog_completar.go](../internal/domain/backlog_completar.go) (`ElementoBacklog.Completar`) | [backlog_completar_test.go](../internal/domain/backlog_completar_test.go) |
| Caso de uso | [internal/app/completar_historia.go](../internal/app/completar_historia.go) (`CompletarHistoria`, puerto `RepositorioCompletar`), [sprint_del_elemento.go](../internal/app/sprint_del_elemento.go) (`sprintDelElementoAbierto`, compartida con HU-09) | [completar_historia_test.go](../internal/app/completar_historia_test.go) |
| Persistencia | [internal/adapters/persistencia/completar_repo.go](../internal/adapters/persistencia/completar_repo.go) (`GuardarCompletado` atómico, `errorDeCompletado`) | [completar_repo_test.go](../internal/adapters/persistencia/completar_repo_test.go), [completar_concurrencia_test.go](../internal/adapters/persistencia/completar_concurrencia_test.go), contra SQLite real |
| HTTP | [internal/adapters/http/completar_handler.go](../internal/adapters/http/completar_handler.go) (`RegistrarCompletar`), registrada en [router.go](../internal/adapters/http/router.go) y conectada en [cmd/server/main.go](../cmd/server/main.go) | [completar_handler_test.go](../internal/adapters/http/completar_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [internal/adapters/http/pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go) (`completarElementoPantalla`, `PuedeCompletarse`), ruta en [pantallas.go](../internal/adapters/http/pantallas.go), [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | [pantallas_completar_test.go](../internal/adapters/http/pantallas_completar_test.go) (`pantallas_asignacion_test.go` no tiene tests de HU-10), [router_test.go](../internal/adapters/http/router_test.go) |
| BDD | [test/bdd/hu10_sistema_test.go](../test/bdd/hu10_sistema_test.go) (caso de uso real + SQLite), [hu10_steps_test.go](../test/bdd/hu10_steps_test.go) | 9 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| Spec | `16de21b` borrador → `a3d76dd` aprobada | — |
| Escenarios BDD y steps | `c503bf6` feature y steps (RED) → `219b889` steps usan `domain.ErrSprintNoEncontrado` | BDD en rojo |
| Dominio | `4417453` [RED] → `978155d` [GREEN] → `990744b` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Caso de uso | `2743291` [RED] → `3890a69` [GREEN] → `d1bf254` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Persistencia | `8e4f2b5` [RED] → `1ba612a` [GREEN] | **RED → GREEN** |
| BDD con SQLite real | `ee6283b`, `6b587ea` | BDD en verde |
| Spec: P-9 | `f24a37d` | — |
| HTTP | `7599f50` [RED] → `28abcdb` [GREEN] | **RED → GREEN** |
| Router y `main` | `ec11388` [RED] → `f045837` [GREEN] | **RED → GREEN** |
| Spec: CA-10.7 | `1472da0` | — |
| Pantalla | `484fe58` [RED] → `3a991c0` [GREEN] → `dfdd0ae` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Revisión | `e7863bc` [REFACTOR], `15b1cee` [REFACTOR] | **REFACTOR** |
| Revisión: carrera con el elemento movido (P-9) | `0753a5f` | test (sin etiqueta de etapa) |
| Spec: limitación conocida de la pantalla | `e638e60` | — |
| Trazabilidad | este commit | — |

### Huecos

- **Limitación conocida de la pantalla** (spec, sección 6): `completarElementoPantalla` no verifica que el elemento esté en el sprint de la URL, igual que "Quitar" de HU-09. No se puede provocar desde la interfaz, solo con un POST armado a mano. No tiene test; se propone resolverlo para HU-09 y HU-10 juntas, fuera de esta historia.
- **P-9** la tomó el autor y queda abierta a objeción en el PR. Tiene tests de persistencia, pero no escenario BDD.
- CA-10.6 y CA-10.7 no tienen escenario BDD a propósito: los escenarios prueban el núcleo y no pasan por HTTP. Los cubren los tests del handler y de la pantalla.
- RN-10.4 se prueba en el caso de uso y en BDD; en HTTP y en la pantalla solo con el error simulado. El repositorio no distingue el proyecto.
- Los commits de spec y BDD (`16de21b`, `c503bf6`, `a3d76dd`, `219b889`, `ee6283b`, `6b587ea`, `f24a37d`, `1472da0`, `e638e60`) usan los scopes `(spec)` y `(bdd)` en vez de `(HU-10)`, así que `git log --grep="(HU-10)"` no los encuentra. `0753a5f` es un test sin etiqueta de etapa.
- El sistema BDD cierra el sprint directamente en la base porque HU-11 (cerrar sprint) todavía no está en esta rama.

---

## HU-11: Cerrar un Sprint

- **Issue:** #22 · Sprint 2 · 3 SP
- **Spec:** [docs/specs/HU-11-cerrar-sprint.md](specs/HU-11-cerrar-sprint.md)
- **Feature:** [features/HU-11-cerrar-sprint.feature](../features/HU-11-cerrar-sprint.feature) (steps en [test/bdd/hu11_steps_test.go](../test/bdd/hu11_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-11.1 Cerrar un sprint Abierto, en cualquier momento, sin cambiar sus demás datos | RN-11.1, RN-11.2 | Un sprint abierto se cierra en cualquier momento y conserva sus datos `@normal` `@limite` (3 ejemplos) · Se cierra un sprint sin historias asignadas `@limite` · Se cierra un sprint de un proyecto que no está Activo `@limite` (4 ejemplos) | `TestSprint_Cerrar`, `TestCerrarSprint_SprintAbierto_QuedaCerradoConSusDatos`, `TestRepositorioSprints_CerrarSprint_Abierto_QuedaCerradoConResumenYDevuelveLasSinTerminar` | `domain.Sprint.Cerrar`, `app.CerrarSprint.Ejecutar`, `persistencia.RepositorioSprints.CerrarSprint` |
| CA-11.2 Un sprint ya Cerrado no se vuelve a cerrar | RN-11.1, RN-11.5 | No se puede cerrar un sprint que ya está cerrado `@error` · Un sprint que se acaba de cerrar no se puede volver a cerrar `@alternativo` | `TestSprint_Cerrar` (ya cerrado), `TestSprint_Cerrar_DosVeces_RetornaError`, `TestCerrarSprint_SprintYaCerrado_RetornaError`, `TestRepositorioSprints_CerrarSprint_YaCerrado_RetornaErrorSinCambios`, `TestRepositorioSprints_CerrarSprint_Simultaneos_SoloUnoGana` | `domain.Sprint.Cerrar`, `persistencia.RepositorioSprints.CerrarSprint` (verificado en el propio `UPDATE`), `persistencia.RepositorioSprints.errorDeCierre` |
| CA-11.3 Proyecto o sprint inexistente, o sprint de otro proyecto | RN-11.4 | No se puede cerrar un sprint de un proyecto que no existe `@error` · No se puede cerrar un sprint que no existe `@error` · Un sprint de otro proyecto se trata como inexistente `@error` | `TestCerrarSprint_ProyectoInexistente_RetornaError`, `TestCerrarSprint_SprintInexistente_RetornaError`, `TestCerrarSprint_SprintDeOtroProyecto_RetornaErrorYSigueAbierto`, `TestCerrarSprint_ValidacionFallida_NoLlamaACerrarSprint`, `TestRepositorioSprints_BuscarSprintPorID_Inexistente`, `TestRepositorioSprints_CerrarSprint_Inexistente_RetornaError` | `app.CerrarSprint.Ejecutar` (`s.ProyectoID != d.ProyectoID`), `persistencia.RepositorioSprints.BuscarSprintPorID` |
| CA-11.4 Respuestas HTTP 204 / 404 / 400 | — (sección 3) | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y no pasan por HTTP (comentario del feature). Se cubre con los tests del handler. | `TestPostCerrarSprint_PedidoValido_Devuelve204SinCuerpo`, `TestPostCerrarSprint_NoEncontrado_Devuelve404`, `TestPostCerrarSprint_YaCerrado_Devuelve400`, `TestPostCerrarSprint_YaCerradoEnvuelto_Devuelve400ConMensajeLimpio`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerCerrarSprint`, `http.RegistrarCerrarSprint`, `http.NuevoRouter` |
| CA-11.5 Botón "Cerrar sprint" solo en sprints Abiertos; al usarlo queda Cerrado y sin botón | RN-11.1 (sección 3) | **Sin escenario BDD a propósito:** es la pantalla (comentario del feature). Se cubre con los tests de la pantalla. | `TestSprint_PuedeCerrarse`, `TestPantallaSprints_BotonCerrar_SoloEnSprintsAbiertos`, `TestPantallaSprints_Cerrar_RedirigeALaLista`, `TestPantallaSprints_Cerrar_DePuntaAPunta_QuedaCerradoYSinBoton`, `TestPantallaSprints_Cerrar_Errores`, `TestPantallaSprints_Cerrar_YaCerrado_SigueMostrandoLaLista`, `TestNuevoRouter_RegistraTodasLasRutas` | `domain.Sprint.PuedeCerrarse`, `http.pantallas.cerrarSprintPantalla`, `http.aVistaSprint` (campo `PuedeCerrarse`), [sprints.html](../internal/adapters/http/web/plantillas/sprints.html) |
| CA-11.6 Las historias sin terminar vuelven al Product Backlog; las Hecho quedan en el sprint | RN-11.3 | Las historias sin terminar vuelven al Product Backlog y las terminadas quedan en el sprint `@normal` (sobre SQLite, con historias asignadas por HU-09 y completadas por HU-10) | `TestResumirCierre`, `TestResumirCierre_NoModificaLasHistorias`, `TestRepositorioSprints_CerrarSprint_Abierto_QuedaCerradoConResumenYDevuelveLasSinTerminar`, `TestRepositorioSprints_CerrarSprint_YaCerrado_RetornaErrorSinCambios` | `domain.ResumirCierre` (`DevueltasAlBacklog`), `persistencia.RepositorioSprints.CerrarSprint` (en la misma transacción) |
| CA-11.7 El sprint registra sus story points planificados y completados, y no cambian después | RN-11.6 | Al cerrar, el sprint registra sus story points planificados y completados `@normal` · Los puntos registrados no cambian si después se replanifica una historia `@alternativo` · Un sprint sin historias se cierra con 0 puntos `@limite` | `TestResumirCierre`, `TestRepositorioSprints_CerrarSprint_Abierto_QuedaCerradoConResumenYDevuelveLasSinTerminar`, `TestRepositorioSprints_CerrarSprint_ResumenFijoAunqueSeReasigne`, `TestAbrirDB_BaseExistenteSinPuntosDeSprint_AgregaLasColumnas` | `domain.ResumirCierre`, campos `Sprint.PuntosPlanificados` y `Sprint.PuntosCompletados`, `persistencia.RepositorioSprints.CerrarSprint`, `persistencia.migrarPuntosDeSprint` |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-11.1 Cerrar es pasar de Abierto a Cerrado; cerrar uno ya Cerrado es un error | CA-11.1, CA-11.2, CA-11.5 | `domain.Sprint.Cerrar`, `domain.Sprint.PuedeCerrarse`, `domain.ErrSprintYaCerrado`; en la base, `UPDATE ... WHERE id = ? AND estado = 'Abierto'` |
| RN-11.2 Se cierra en cualquier momento, con o sin historias, y aunque el proyecto no esté Activo (P-10) | CA-11.1 | `app.CerrarSprint.Ejecutar` no mira fechas, historias ni el estado del proyecto |
| RN-11.3 Las historias sin terminar vuelven al Product Backlog; las Hecho quedan en el sprint (P-4, P-13) | CA-11.6 | `domain.ResumirCierre`; `persistencia.RepositorioSprints.CerrarSprint` pone `sprint_id = NULL` a las devueltas |
| RN-11.4 Un sprint de otro proyecto se trata como inexistente | CA-11.3 | `app.CerrarSprint.Ejecutar` (`domain.ErrSprintNoEncontrado`) |
| RN-11.5 Un sprint cerrado no se reabre | CA-11.2 | `domain.Sprint.Cerrar` (nunca vuelve a Abierto); no hay operación de reapertura |
| RN-11.6 Al cerrar se registran los puntos planificados y completados; sin estimar suma 0; quedan fijos (P-11, P-12) | CA-11.7 | `domain.ResumirCierre`; columnas `puntos_planificados` y `puntos_completados` de `sprints` |

Sin CA propio en la spec:
- Pedidos simultáneos (sección 5, P-8): `TestRepositorioSprints_CerrarSprint_Simultaneos_SoloUnoGana` (10 goroutines, exactamente una cierra), `TestRepositorioSprints_CerrarSprint_SimultaneoConCompletados_QuedaCoherente` (un cierre contra 8 completados de HU-10: cada historia queda Hecho en el sprint o Pendiente sin sprint, y los puntos completados coinciden), `TestCerrarSprint_FallaDelRepositorio_SePropaga` (`ErrSprintYaCerrado` desde el repositorio), `TestPostCerrarSprint_YaCerradoEnvuelto_Devuelve400ConMensajeLimpio` y `TestPantallaSprints_Cerrar_Errores` (error envuelto, se muestra el mensaje del dominio).
- Si hay un error no se modifica nada (sección 3): `TestCerrarSprint_ValidacionFallida_NoLlamaACerrarSprint`.
- IDs inválidos en la ruta (400 sin llamar al caso de uso): `TestPostCerrarSprint_IDInvalido_Devuelve400`, `TestPantallaSprints_Cerrar_IDInvalido_Devuelve400`.
- Falla interna (500 sin el detalle): `TestCerrarSprint_FallaDelRepositorio_SePropaga`, `TestPostCerrarSprint_FallaInterna_Devuelve500SinDetalle`, `TestPantallaSprints_Cerrar_FallaInterna_Devuelve500SinDetalle`.
- Formulario enviado desde otro sitio (CSRF, mismo origen): `TestPantallaSprints_Cerrar_DesdeOtroSitio_SeRechaza`.
- Lectura del sprint por ID (puerto `RepositorioCierre`): `TestRepositorioSprints_BuscarSprintPorID_Existente`, `TestRepositorioSprints_BuscarSprintPorID_Inexistente`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/sprint.go](../internal/domain/sprint.go) (`Cerrar`, `PuedeCerrarse`, `ErrSprintYaCerrado`, campos de puntos), [sprint_cierre.go](../internal/domain/sprint_cierre.go) (`ResumirCierre`, `ResumenCierre`) | [sprint_test.go](../internal/domain/sprint_test.go), [sprint_cierre_test.go](../internal/domain/sprint_cierre_test.go) |
| Caso de uso y puerto | [internal/app/cerrar_sprint.go](../internal/app/cerrar_sprint.go) (`CerrarSprint`, puerto `RepositorioCierre`) | [cerrar_sprint_test.go](../internal/app/cerrar_sprint_test.go) |
| Persistencia | [internal/adapters/persistencia/cierre_repo.go](../internal/adapters/persistencia/cierre_repo.go) (`BuscarSprintPorID`, `CerrarSprint` en una transacción, `historiasDelSprint`, `errorDeCierre`), [sprint_repo.go](../internal/adapters/persistencia/sprint_repo.go) (`escanearSprint`, compartido con `ListarSprints` de HU-08), [db.go](../internal/adapters/persistencia/db.go) (columnas de puntos, `migrarPuntosDeSprint`, `agregarColumnas`) | [cierre_repo_test.go](../internal/adapters/persistencia/cierre_repo_test.go), [cierre_concurrencia_test.go](../internal/adapters/persistencia/cierre_concurrencia_test.go), [db_test.go](../internal/adapters/persistencia/db_test.go) |
| HTTP | [internal/adapters/http/cerrar_sprint_handler.go](../internal/adapters/http/cerrar_sprint_handler.go) (`POST /proyectos/{id}/sprints/{sprintID}/cerrar`, `RegistrarCerrarSprint`), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [cerrar_sprint_handler_test.go](../internal/adapters/http/cerrar_sprint_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [internal/adapters/http/pantallas_sprints.go](../internal/adapters/http/pantallas_sprints.go) (`cerrarSprintPantalla`), [pantallas.go](../internal/adapters/http/pantallas.go) (ruta `POST /web/proyectos/{id}/sprints/{sprintID}/cerrar`), [pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go) (`aVistaSprint`), [sprints.html](../internal/adapters/http/web/plantillas/sprints.html) | [pantallas_cerrar_sprint_test.go](../internal/adapters/http/pantallas_cerrar_sprint_test.go) |
| BDD | [test/bdd/hu11_sistema_test.go](../test/bdd/hu11_sistema_test.go) (SQLite real), [hu11_steps_test.go](../test/bdd/hu11_steps_test.go) | 17 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| Spec | `fc46a3d` docs(spec): HU-11 cerrar sprint (borrador) · `6ef2057` docs(spec): HU-11 decisiones tomadas por el autor | — |
| Escenarios BDD y steps | `5019801` test(bdd): HU-11 feature y steps (RED) | BDD en rojo |
| Dominio (`Sprint.Cerrar`) | `6e9ad50` [RED] → `7bf464b` [GREEN] → `ef5542e` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Steps con el error de dominio de HU-09 | `4013785` test(bdd): HU-11 steps usan domain.ErrSprintNoEncontrado | — |
| Caso de uso (`app.CerrarSprint`) | `74f5563` [RED] → `749a531` [GREEN] | **RED → GREEN** (sin REFACTOR: son cuatro pasos lineales con el mismo molde que `CrearSprint` de HU-08; no había nada que mejorar) |
| Persistencia | `8e7f521` [RED] → `3b0948c` [GREEN] → `b48a488` [REFACTOR] | **RED → GREEN → REFACTOR** |
| BDD sobre SQLite real | `bff0b15` test(bdd): HU-11 sistema SQLite real | BDD en verde |
| Endpoint JSON | `b8a6c5b` [RED] → `5fd5a86` [GREEN] | **RED → GREEN** (sin REFACTOR: sigue el patrón de los handlers de HU-09; mover los helpers de IDs, repartidos en archivos de otras historias, tocaba código ajeno a HU-11) |
| Escenario de CA-11.6 | `b9d1c74` test(bdd): HU-11 escenario CA-11.6 | BDD en verde |
| Pantalla | `9a4cb4e` [RED] → `6b66891` [GREEN] → `4aa5329` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Revisión (go-reviewer): `Traza` mal ubicado (B-1) | `3f609e4` test(HU-11): corregir comentarios Traza en sprint_test.go | — |
| Revisión (go-reviewer): CA-11.5 de punta a punta (M-2) | `f9272fd` test(HU-11): CA-11.5 de punta a punta en pantalla de sprints | — (test nuevo sobre código ya en verde; se verificó que falla con `PuedeCerrarse` mutado a `true`) |
| Revisión (go-reviewer): contrato del puerto (B-2) | `100fe21` docs(HU-11): aclarar contrato de CerrarSprint en el puerto | — |
| Revisión (go-reviewer): spec desactualizada (M-1) | `6be321b` docs(HU-11): actualizar spec al estado real del código | — |
| Merge de `main` (HU-10) en la rama | `15bf901` | — (conflictos de líneas agregadas por las dos historias) |
| Revisión del PR: P-4 cambiada y resumen de puntos (spec) | `61bed0a` docs(HU-11): historias sin terminar vuelven al backlog y resumen de puntos al cerrar | — |
| Dominio (`ResumirCierre`) | `7ef1795` [RED] → `ce33313` [GREEN] | **RED → GREEN** (sin REFACTOR: un solo recorrido) |
| Persistencia (transacción, columnas de puntos, migración) | `37fba7b` [RED] → `e1c3b9b` [GREEN] → `8ca607e` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Escenarios de CA-11.6 y CA-11.7 | `c5d3fe7` | BDD en verde (escritos después del código; se verificó que 3 fallan con el cierre anterior) |
| Cierre simultáneo con completados | `adc458d` | — (test nuevo sobre código ya en verde) |
| Trazabilidad | este commit | — |

### Huecos

- **En parte resuelto (2026-10-29):** `TestRepositorios_BaseCerrada_DevuelvenElError` llama a los 31 métodos públicos de los repositorios con la base cerrada y verifica que todos devuelvan el error con su causa, y `TestAbrirDB_RutaInvalida_RetornaError` cubre los errores al abrir la base. La cobertura de persistencia pasó de 80,9% a 86,2%. Siguen sin test los errores que ocurren después de una consulta exitosa (`Scan` de una fila, `RowsAffected`, `Commit`), porque para provocarlos hay que simular la base. Detalle original: Ramas defensivas de persistencia sin test (hallazgo B-3 del go-reviewer): falla de `ExecContext` y de `RowsAffected` en `CerrarSprint`, error genérico de `Scan` en `BuscarSprintPorID`, el "estado guardado inesperado" de `errorDeCierre` y las fechas inválidas de `escanearSprint`. `escanearSprint` también lo usa `ListarSprints` de HU-08, así que ese hueco es compartido.
- `go test -race` no corre en el entorno de desarrollo (gcc de 32 bits, sin cgo de 64 bits): la ausencia de carreras de datos en el test de concurrencia se razonó a mano (cada goroutine escribe solo `errs[i]` y copia el sprint por valor). Conviene correrlo en CI o en Linux.
- El resumen de puntos se guarda pero todavía no se muestra en ninguna pantalla: lo va a mostrar HU-22 (velocidad).
- Si después del cierre se edita con HU-06 una historia `Hecho` que quedó en el sprint cerrado (por ejemplo, se le cambian los puntos), el resumen guardado no cambia: es lo que pide RN-11.6, pero conviene tenerlo en cuenta en HU-22.
- RN-11.2 con un proyecto que no está Activo (P-10) se verifica solo en los escenarios BDD: no hay test unitario del caso de uso con un proyecto en otro estado.
- La ruta JSON `POST /proyectos/{id}/sprints/{sprintID}/cerrar` no verifica el origen del pedido (CSRF): solo las pantallas pasan por `mismoOrigen`. Es deuda técnica de todas las rutas JSON del proyecto, no solo de HU-11.
- `ErrSprintYaCerrado` (HU-11, "el sprint ya fue cerrado") y `ErrSprintCerrado` (HU-09, "el sprint ya está cerrado") describen la misma situación; se mantuvieron separados por el orden de los merges y unificarlos es una decisión del equipo.
- Refactors que tocaron código de otras historias: `escanearSprint` reemplazó el escaneo de `ListarSprints` (HU-08) y `aVistaSprint` (HU-09) ahora completa `PuedeCerrarse` y la usa también `mostrarSprints` (HU-08). Los tests de HU-08 y HU-09 siguen en verde, pero conviene que lo revisen sus autores en el PR.

---

## HU-20: Registrar un defecto

- **Issue:** [#31](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/31) · Sprint 2 · 3 SP
- **Spec:** [docs/specs/HU-20-registrar-defecto.md](specs/HU-20-registrar-defecto.md)
- **Feature:** [features/HU-20-registrar-defecto.feature](../features/HU-20-registrar-defecto.feature) (steps en [test/bdd/hu20_steps_test.go](../test/bdd/hu20_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-20.1 Registrar con Sprint activo | RN-20.1, RN-20.4, RN-20.5, RN-20.6 | Registrar un defecto con todos los datos válidos y un Sprint activo `@normal` · La descripción se guarda sin los espacios de alrededor `@alternativo` | `TestNuevoDefecto_Valido`, `TestSprintVigente_SprintAbiertoConHoyEnRango_DevuelveSuID`, `TestRegistrarDefecto_DatosValidos_GuardaConID`, `TestRegistrarDefecto_ConSprintVigente_AsignaSuID`, `TestRepositorioDefectos_ConSprintDeteccion_LoGuarda`, `TestPostDefecto_DatosValidos_Devuelve201ConElDefecto`, `TestPostDefecto_ConSprintDeteccion_LoInformaEnLaRespuesta`, `TestPantallaDefecto_Crear_ConSprintDeteccion_MuestraConfirmacion` | `domain.NuevoDefecto`, `domain.SprintVigente`, `app.RegistrarDefecto` |
| CA-20.2 Sin Sprint vigente | RN-20.5 | Sin ningún sprint en el proyecto, el defecto queda sin Sprint de detección `@alternativo` · Un sprint que ya terminó no se asigna `@limite` · Un sprint cuyo último día es hoy se asigna `@limite` | `TestSprintVigente_HoyEnLosLimites_DevuelveSuID`, `TestSprintVigente_SinSprints_NoHayNinguno`, `TestSprintVigente_HoyFueraDelRango_NoHayNinguno`, `TestRegistrarDefecto_SinSprintVigente_QuedaSinAsignar`, `TestRepositorioDefectos_SinSprintDeteccion_GuardaNull`, `TestPantallaDefecto_Crear_SinSprintDeteccion_AvisaQueNoSeAsigno` | `domain.SprintVigente` |
| CA-20.3 Descripción obligatoria | RN-20.2 | La descripción es obligatoria `@error` (×2) | `TestNuevoDefecto_SinDescripcion_RetornaError`, `TestRegistrarDefecto_DatosInvalidos_NoLleganAlRepositorio`, `TestPostDefecto_DatosInvalidos_Devuelve400ConElError`, `TestPantallaDefecto_Crear_DatosInvalidos_MuestraElError` | `domain.NuevoDefecto` |
| CA-20.4 Severidad válida | RN-20.3 | La severidad tiene que ser Crítica, Alta, Media o Baja `@error` (×3) · Las cuatro severidades válidas se aceptan `@limite` (×4) | `TestNuevoDefecto_SeveridadInvalida_RetornaError`, `TestNuevoDefecto_SeveridadesValidas` | `domain.NuevoDefecto` |
| CA-20.5 Proyecto o elemento inexistente | RN-20.1 | No se puede registrar un defecto en un proyecto que no existe `@error` · No se puede registrar un defecto sobre un elemento que no existe `@error` · Un elemento de otro proyecto se trata como inexistente `@error` | `TestRegistrarDefecto_ElementoInexistente_RetornaErrorYNoGuarda`, `TestRegistrarDefecto_ProyectoInexistente_RetornaErrorYNoGuarda`, `TestRegistrarDefecto_ElementoDeOtroProyecto_RetornaErrorYNoGuarda`, `TestRepositorioDefectos_ElementoInexistente_RetornaErrElementoNoEncontrado`, `TestPostDefecto_ProyectoOElementoInexistente_Devuelve404` | `app.RegistrarDefecto`, clave foránea de `defectos` |
| CA-20.6 Respuesta HTTP | — | No va acá (como CA-18.15): se cubre con los tests del adaptador | `TestPostDefecto_*` (8 tests) | `internal/adapters/http/defecto_handler.go` |
| CA-20.7 Independencia de IDs | RN-20.6 | Dos defectos sobre el mismo elemento son independientes `@normal` | `TestRepositorioDefectos_Guardar_AsignaIDsDistintos`, `TestPostDefecto_DatosValidos_PasaLosDatosAlCasoDeUso` | `persistencia.RepositorioDefectos.GuardarDefecto` |
| CA-20.8 Pantalla | Restricciones | No va acá (como CA-06.10): se agregó en la tarea 20.8 y se cubre con los tests de la pantalla | `TestPantallaDefecto_*` (6 tests) | `pantallas_defecto.go`, `web/plantillas/defecto_form.html`, link en `backlog.html` (HU-07) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/defecto.go](../internal/domain/defecto.go) (`Defecto`, `Severidad`, `EstadoDefecto`, `NuevoDefecto`, `SprintVigente`) | [defecto_test.go](../internal/domain/defecto_test.go) |
| Caso de uso | [internal/app/registrar_defecto.go](../internal/app/registrar_defecto.go) | [registrar_defecto_test.go](../internal/app/registrar_defecto_test.go) |
| Persistencia | [internal/adapters/persistencia/defecto_repo.go](../internal/adapters/persistencia/defecto_repo.go), [db.go](../internal/adapters/persistencia/db.go) (tabla `defectos`) | [defecto_repo_test.go](../internal/adapters/persistencia/defecto_repo_test.go) |
| HTTP (JSON) | [internal/adapters/http/defecto_handler.go](../internal/adapters/http/defecto_handler.go) (`POST /proyectos/{id}/defectos`) | [defecto_handler_test.go](../internal/adapters/http/defecto_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_defecto.go](../internal/adapters/http/pantallas_defecto.go), [web/plantillas/defecto_form.html](../internal/adapters/http/web/plantillas/defecto_form.html) (`/web/proyectos/{id}/defectos/nuevo`) | [pantallas_defecto_test.go](../internal/adapters/http/pantallas_defecto_test.go) |
| BDD | [test/bdd/hu20_sistema_test.go](../test/bdd/hu20_sistema_test.go), [hu20_steps_test.go](../test/bdd/hu20_steps_test.go) | 18 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 20.1 Spec | `7b8d3ff` docs(HU-20): [20.1] spec SDD | — |
| 20.2 Escenarios BDD | `3e43b7d` test(HU-20): [20.2] escenarios BDD | — |
| 20.3 Steps godog | `6356ec3` test(HU-20): [20.3][RED] steps (en rojo) | BDD en rojo |
| 20.4 Entidad Defecto | `e47f6f9` [RED] → `b76b4aa` [GREEN] | **RED → GREEN** |
| 20.5 SprintVigente | `0a11c53` [RED] → `f0c6af1` [GREEN] | **RED → GREEN** |
| 20.5 Caso de uso RegistrarDefecto | `e55b04c` [RED] → `faf7446` [GREEN] | **RED → GREEN** |
| 20.6 Persistencia | `77996f6` [RED] → `c5060fd` [GREEN] | **RED → GREEN** |
| Conectar los escenarios BDD con el caso de uso real | `383a99d` | BDD en verde |
| 20.7 HTTP | `17e092d` [RED] → `4dbdd36` [GREEN] | **RED → GREEN** |
| 20.8 Pantalla | `d7b29a7` [RED] → `8bd8416` [GREEN] | **RED → GREEN** |
| 20.9 Revisión y trazabilidad | `d76b7c9` | — |

### Huecos

- ~~La spec está en Borrador y P-1 a P-3 pendientes~~: resuelto el 2026-10-29, Matías López confirmó P-1 a P-3 y la spec quedó Aprobada.
- El enlace "Registrar defecto" sale de la lista del backlog (HU-07); no hay todavía un enlace desde el detalle del elemento (HU-06), porque la pantalla de HU-20 elige el elemento de una lista, no parte de uno puntual.
- Fuera de alcance (spec, sección 5): editar, eliminar, listar o cambiar el estado de un defecto; vincularlo manualmente a un Sprint distinto del detectado automáticamente.
- Revisión con el agente `go-reviewer` (tarea 20.9): sin hallazgos bloqueantes. Encontró que el `<select>` de severidad de la pantalla no tenía una opción vacía por defecto ni `required` (a diferencia del de elemento); se corrigió en el mismo commit de esta revisión.

---

## HU-21: Actualizar el estado de un defecto

- **Issue:** a confirmar · Sprint 3 · 2 SP
- **Spec:** [docs/specs/HU-21-actualizar-estado-defecto.md](specs/HU-21-actualizar-estado-defecto.md)
- **Feature:** [features/HU-21-actualizar-estado-defecto.feature](../features/HU-21-actualizar-estado-defecto.feature) (steps en [test/bdd/hu21_steps_test.go](../test/bdd/hu21_steps_test.go), sistema en [hu21_sistema_test.go](../test/bdd/hu21_sistema_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-21.1 Iniciar un defecto Abierto | RN-21.2 | Iniciar un defecto Abierto `@normal` | `TestDefecto_Iniciar_DesdeAbierto_QuedaEnProgreso`, `TestActualizarEstadoDefecto_Iniciar_DesdeAbierto_QuedaEnProgreso` | `domain.Defecto.Iniciar`, `app.ActualizarEstadoDefecto.Iniciar` |
| CA-21.2 Resolver con Sprint activo asigna su ID | RN-21.3, RN-21.6 | Resolver un defecto En progreso con un Sprint activo `@normal` | `TestDefecto_Resolver_DesdeEnProgreso_QuedaResuelto`, `TestActualizarEstadoDefecto_Resolver_ConSprintVigente_AsignaSuID`, `TestRepositorioDefectos_ActualizarEstado_ConSprintResolucion_LoGuarda` | `domain.Defecto.Resolver`, `app.ActualizarEstadoDefecto.Resolver` |
| CA-21.3 Sin Sprint vigente, queda nulo | RN-21.6 | Resolver un defecto En progreso sin ningún Sprint vigente `@alternativo` · Un sprint que ya terminó no se asigna como Sprint de resolución `@limite` | `TestActualizarEstadoDefecto_Resolver_SinSprintVigente_QuedaSinAsignar` | `domain.SprintVigente` (reutilizado de HU-20), `app.ActualizarEstadoDefecto.Resolver` |
| CA-21.4 Cerrar un defecto Resuelto, sin tocar el Sprint de resolución | RN-21.4 | Cerrar un defecto Resuelto `@normal` | `TestDefecto_Cerrar_DesdeResuelto_QuedaCerrado`, `TestActualizarEstadoDefecto_Cerrar_DesdeResuelto_QuedaCerradoConElMismoSprint` | `domain.Defecto.Cerrar`, `app.ActualizarEstadoDefecto.Cerrar` |
| CA-21.5 Iniciar uno que no está Abierto | RN-21.2, RN-21.5 | No se puede iniciar un defecto que ya está En progreso `@error` · No se puede iniciar un defecto que ya está Cerrado `@error` | `TestDefecto_Iniciar_NoAbierto_RetornaError`, `TestActualizarEstadoDefecto_Iniciar_NoAbierto_RetornaErrorYNoGuarda` | `domain.Defecto.Iniciar` (`ErrDefectoNoAbierto`) |
| CA-21.6 Resolver uno que no está En progreso | RN-21.3, RN-21.5 | No se puede resolver un defecto Abierto (saltear En progreso) `@error` · No se puede resolver un defecto que ya está Cerrado `@error` | `TestDefecto_Resolver_NoEnProgreso_RetornaError`, `TestActualizarEstadoDefecto_Resolver_NoEnProgreso_RetornaErrorYNoGuarda` | `domain.Defecto.Resolver` (`ErrDefectoNoEnProgreso`) |
| CA-21.7 Cerrar uno que no está Resuelto | RN-21.4, RN-21.5 | No se puede cerrar un defecto Abierto (saltear En progreso y Resuelto) `@error` · No se puede cerrar un defecto En progreso (saltear Resuelto) `@error` | `TestDefecto_Cerrar_NoResuelto_RetornaError`, `TestActualizarEstadoDefecto_Cerrar_NoResuelto_RetornaErrorYNoGuarda` | `domain.Defecto.Cerrar` (`ErrDefectoNoResuelto`) |
| CA-21.8 Proyecto o defecto inexistente, o de otro proyecto | RN-21.7 | No se puede actualizar un defecto de un proyecto inexistente `@error` · No se puede actualizar un defecto inexistente `@error` · Un defecto de otro proyecto se trata como inexistente `@error` | `TestActualizarEstadoDefecto_ProyectoInexistente_RetornaError`, `TestActualizarEstadoDefecto_DefectoInexistente_RetornaError`, `TestActualizarEstadoDefecto_DefectoDeOtroProyecto_RetornaError`, `TestRepositorioDefectos_BuscarDefectoPorID_Inexistente_RetornaError` | `app.ActualizarEstadoDefecto.buscar` (`ErrDefectoNoEncontrado`) |
| CA-21.9 Respuesta HTTP de iniciar/resolver/cerrar | Sección 3 | No va acá (como CA-20.6): se cubre con los tests del adaptador | `TestPostIniciarDefecto_*`, `TestPostResolverDefecto_*`, `TestPostCerrarDefecto_*`, `TestPostActualizarEstadoDefecto_*` (10 tests) | `internal/adapters/http/actualizar_estado_defecto_handler.go` |
| CA-21.10 Listado JSON de apoyo, ordenado por id | Sección 3 | No va acá: es infraestructura de la pantalla, no una regla de negocio propia | `TestListarDefectos_*` (3), `TestRepositorioDefectos_ListarDefectos_*` (3), `TestGetListarDefectos_*` (3) | `app.ListarDefectos`, `persistencia.RepositorioDefectos.ListarDefectos` (`ORDER BY id`), `internal/adapters/http/listar_defectos_handler.go` |
| CA-21.11 Pantalla con botones según estado | Sección 3 | No va acá (como CA-20.8): es la pantalla | `TestPantallaDefectos_*` (3), `TestPantallaIniciarDefecto_*` (2), `TestPantallaResolverDefecto_*` (2), `TestPantallaCerrarDefecto_*` (2), `TestPantallaActualizarEstadoDefecto_DefectoInexistente_Es404` | `pantallas_defectos_lista.go`, `web/plantillas/defectos.html`, enlaces en `proyecto_detalle.html` y `defecto_form.html` |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-21.1 Cuatro estados: Abierto, En progreso, Resuelto, Cerrado | Todas | `domain.EstadoDefecto`, `domain.DefectoAbierto/EnProgreso/Resuelto/Cerrado` |
| RN-21.2 Iniciar = Abierto→En progreso; error si no está Abierto | CA-21.1, CA-21.5 | `domain.Defecto.Iniciar`, `domain.Defecto.PuedeIniciarse` |
| RN-21.3 Resolver = En progreso→Resuelto; error si no está En progreso | CA-21.2, CA-21.6 | `domain.Defecto.Resolver`, `domain.Defecto.PuedeResolverse` |
| RN-21.4 Cerrar = Resuelto→Cerrado; error si no está Resuelto | CA-21.4, CA-21.7 | `domain.Defecto.Cerrar`, `domain.Defecto.PuedeCerrarse` |
| RN-21.5 Siempre secuencial, sin saltos ni retrocesos | CA-21.5 a CA-21.7 | `domain.Defecto.Iniciar`/`Resolver`/`Cerrar` (no hay forma de volver a un estado anterior ni saltear uno) |
| RN-21.6 Sprint de resolución automático, mismo criterio que detección | CA-21.2, CA-21.3 | `domain.SprintVigente` (HU-20), `app.ActualizarEstadoDefecto.Resolver` |
| RN-21.7 Defecto de otro proyecto es inexistente | CA-21.8 | `app.ActualizarEstadoDefecto.buscar` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/defecto.go](../internal/domain/defecto.go) (`Iniciar`, `Resolver`, `Cerrar`, `PuedeIniciarse`, `PuedeResolverse`, `PuedeCerrarse`, `DefectoEnProgreso/Resuelto/Cerrado`, `ErrDefectoNoEncontrado/NoAbierto/NoEnProgreso/NoResuelto`) | [defecto_test.go](../internal/domain/defecto_test.go) |
| Caso de uso | [internal/app/actualizar_estado_defecto.go](../internal/app/actualizar_estado_defecto.go) (`Iniciar`, `Resolver`, `Cerrar`), [listar_defectos.go](../internal/app/listar_defectos.go) | [actualizar_estado_defecto_test.go](../internal/app/actualizar_estado_defecto_test.go), [listar_defectos_test.go](../internal/app/listar_defectos_test.go) |
| Persistencia | [internal/adapters/persistencia/defecto_estado_repo.go](../internal/adapters/persistencia/defecto_estado_repo.go) (incluye el caso `DefectoEnProgreso` en `errorDeActualizacion` y `ORDER BY id` en `ListarDefectos`), [db.go](../internal/adapters/persistencia/db.go) (columna `sprint_resolucion_id`) | [defecto_estado_repo_test.go](../internal/adapters/persistencia/defecto_estado_repo_test.go) |
| HTTP (JSON) | [internal/adapters/http/actualizar_estado_defecto_handler.go](../internal/adapters/http/actualizar_estado_defecto_handler.go) (`POST .../iniciar`, `.../resolver`, `.../cerrar`), [listar_defectos_handler.go](../internal/adapters/http/listar_defectos_handler.go) (`GET /proyectos/{id}/defectos`) | [actualizar_estado_defecto_handler_test.go](../internal/adapters/http/actualizar_estado_defecto_handler_test.go), [listar_defectos_handler_test.go](../internal/adapters/http/listar_defectos_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_defectos_lista.go](../internal/adapters/http/pantallas_defectos_lista.go), [web/plantillas/defectos.html](../internal/adapters/http/web/plantillas/defectos.html) (`/web/proyectos/{id}/defectos`) | [pantallas_defectos_lista_test.go](../internal/adapters/http/pantallas_defectos_lista_test.go) |
| BDD | [test/bdd/hu21_sistema_test.go](../test/bdd/hu21_sistema_test.go), [hu21_steps_test.go](../test/bdd/hu21_steps_test.go) | 14 escenarios en verde |

### Commits

Rama `hu-21-actualizar-estado-defecto`. El ciclo RED → GREEN completo está en el historial (no se hizo squash).

| Tarea | Commit (mensaje) | Etapa |
|---|---|---|
| 21.1 Spec | `docs(HU-21): [21.1] spec SDD de actualizar el estado de un defecto` | — |
| 21.2 Escenarios BDD | `test(HU-21): [21.2] escenarios BDD de actualizar el estado de un defecto` | — |
| 21.3 Steps godog | `test(HU-21): [21.3][RED] steps de godog` | BDD en rojo |
| 21.4 Entidad Defecto (Resolver/Cerrar) | `test(HU-21): [21.4][RED] Defecto.Resolver y Defecto.Cerrar` → `feat(HU-21): [21.4][GREEN] Defecto.Resolver y Defecto.Cerrar` | **RED → GREEN** |
| 21.5 Caso de uso ActualizarEstadoDefecto | `test(HU-21): [21.5][RED] caso de uso ActualizarEstadoDefecto` → `feat(HU-21): [21.5][GREEN] caso de uso ActualizarEstadoDefecto` | **RED → GREEN** |
| 21.6 Persistencia (RepositorioEstadoDefecto, LectorDefectos) | `test(HU-21): [21.6][RED] adaptador de persistencia` → `feat(HU-21): [21.6][GREEN] RepositorioEstadoDefecto y LectorDefectos sobre SQLite` | **RED → GREEN** |
| Conectar los escenarios BDD con los casos de uso reales | `test(HU-21): conectar los escenarios BDD con ActualizarEstadoDefecto` | BDD en verde |
| 21.7 HTTP (resolver/cerrar) | `test(HU-21): [21.7][RED] adaptador HTTP` → `feat(HU-21): [21.7][GREEN] POST .../resolver y .../cerrar, conexión en el router y main` | **RED → GREEN** |
| 21.7 HTTP (listado, apoyo) | `test(HU-21): [21.7][RED] listado HTTP de defectos` → `feat(HU-21): [21.7][GREEN] GET /proyectos/{id}/defectos` | **RED → GREEN** |
| 21.8 Pantalla | `test(HU-21): [21.8][RED] pantalla de defectos` → `feat(HU-21): [21.8][GREEN] pantalla de defectos con botones según estado` | **RED → GREEN** |
| 21.9 Revisión y trazabilidad | `docs(HU-21): [21.9] revisión con go-reviewer y matriz de trazabilidad` | — |
| 21.10 Corrección de la revisión del PR (estado En progreso, ORDER BY) | spec → BDD → `[RED]` dominio/caso de uso/persistencia/HTTP/pantalla → `[GREEN]` dominio/caso de uso/persistencia/HTTP/pantalla → router/main → trazabilidad | **RED → GREEN** |

### Huecos

- La spec quedó Aprobada el 2026-11-01; el 2026-11-03, en la revisión del PR, Matías López pidió agregar el estado `En progreso` entre `Abierto` y `Resuelto` (lo pedía el enunciado de la cátedra). Ezequiel decidió que siempre hay que pasar por `En progreso` (no se puede saltear directo a `Resuelto`), manteniendo el mismo criterio "sin saltos" que ya tenía la historia. Ver Historial de la spec.
- No hay issue de GitHub todavía asignado (el Sprint 3 ya lo repartió a Eze en la planning, pero el número de issue específico de HU-21 queda a confirmar).
- Igual que HU-20, no hay endpoint para consultar un único defecto por ID: `ActualizarEstadoDefecto` lo busca internamente, y el listado (CA-21.10) alcanza para la pantalla.
- Fuera de alcance (spec, sección 5): reabrir un defecto, editar su descripción/severidad/elemento, eliminarlo, elegir manualmente el Sprint de resolución, filtros o búsqueda en la lista.
- Durante el desarrollo se encontró y corrigió un bug: `RegistrarActualizarEstadoDefecto` armaba un *method value* (`uc.Resolver`) al registrar las rutas, lo que hacía panic con una interfaz nil (pasa en varios fixtures BDD de otras historias que arman el router con un `CasosDeUso` parcial, ej. HU-06). Se corrigió pasando la interfaz completa y llamando al método recién al atender cada pedido, igual que los demás handlers.
- Revisión con el agente `go-reviewer` (tarea 21.9, sobre la versión de 3 estados): Approve, sin hallazgos bloqueantes. Confirmó que el fix del bug de method value quedó aplicado correctamente en todo `internal/adapters/http` (recorrió los ~20 archivos `*_handler.go`/`pantallas*.go`) y no encontró otra instancia del mismo problema. Dos notas de estilo no bloqueantes, ambas ya consistentes con patrones existentes: `app.ListarDefectos.Ejecutar` no envuelve el error de `BuscarPorID` (igual que `ListarBacklog`, HU-07), y `errorDeActualizacion` decide el error a partir de `anterior` en vez de releer el estado (determinístico, documentado en el propio código).
- Revisión de nicoandino y Matías López en el PR (2026-11-03): agregar el estado `En progreso` (resuelto, tarea 21.10), traer `main` resolviendo los conflictos en `main.go`/`router.go`/`router_test.go`/`trazabilidad.md` manteniendo ambos lados (pendiente, orden de merge sugerido #288 → #292 → #294, HU-21 va último), y el `ORDER BY id` en `ListarDefectos` (resuelto, tarea 21.10). El aviso de `router_test.go` sin gofmt ya no aplicaba al revisar (el archivo ya pasaba `gofmt -l` antes de tocarlo).

---

## HU-12: Consultar Sprints anteriores

- **Issue:** #23 · Sprint 3 · 2 SP
- **Spec:** [docs/specs/HU-12-sprints-anteriores.md](specs/HU-12-sprints-anteriores.md)
- **Feature:** [features/HU-12-sprints-anteriores.feature](../features/HU-12-sprints-anteriores.feature) (steps en [test/bdd/hu12_steps_test.go](../test/bdd/hu12_steps_test.go))

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-12.1 Los cerrados, del más reciente al más antiguo, con fechas, Sprint Goal y puntos | RN-12.1 a RN-12.3 | Se listan los sprints cerrados, del más reciente al más antiguo, sin los abiertos `@normal` · Un sprint cerrado sin historias aparece con 0 puntos `@limite` | `TestSprintsAnteriores`, `TestSprintsAnteriores_ConservaLosDatos`, `TestListarSprintsAnteriores_DevuelveLosCerradosDelProyecto` | `domain.SprintsAnteriores`, `app.ListarSprintsAnteriores.Ejecutar` |
| CA-12.2 Sin sprints cerrados: lista vacía, sin error | RN-12.1 | Un proyecto sin sprints no tiene sprints anteriores `@limite` · Los sprints abiertos no aparecen, aunque estén vencidos `@alternativo` | `TestSprintsAnteriores` (sin sprints, abiertos vencidos), `TestListarSprintsAnteriores_SinCerrados_ListaVacia` | `domain.SprintsAnteriores` |
| CA-12.3 El orden sigue la fecha de inicio, no el número | RN-12.3 | El orden sigue las fechas aunque el número no coincida `@alternativo` | `TestSprintsAnteriores` (CA-12.3) | `domain.SprintsAnteriores` (`slices.SortFunc` por `FechaInicio`) |
| CA-12.4 Solo los del proyecto consultado | RN-12.1 | Solo aparecen los sprints del proyecto consultado `@alternativo` | `TestListarSprintsAnteriores_DevuelveLosCerradosDelProyecto` | `persistencia.RepositorioSprints.ListarSprints` (HU-08) |
| CA-12.5 Proyecto inexistente | Sección 7 | No se pueden consultar los sprints anteriores de un proyecto que no existe `@error` | `TestListarSprintsAnteriores_ProyectoInexistente_RetornaError`, `TestListarSprintsAnteriores_FallaDelRepositorio_SePropaga` | `app.ListarSprintsAnteriores.Ejecutar` |
| CA-12.6 HTTP 200 / 404 / 400 | Sección 3 | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo (comentario del feature). | `TestGetSprintsAnteriores_DevuelveLosSprintsConSusPuntos`, `TestGetSprintsAnteriores_SinSprints_DevuelveArrayVacio`, `TestGetSprintsAnteriores_Errores`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerSprintsAnteriores`, `http.RegistrarSprintsAnteriores`, `http.NuevoRouter` |
| CA-12.7 Pantalla "Sprints anteriores", enlazada desde el proyecto y los sprints | Sección 3 | **Sin escenario BDD a propósito:** es la pantalla. | `TestPantallaSprintsAnteriores_MuestraLosSprintsConSusDatos`, `TestPantallaSprintsAnteriores_SinSprints_MuestraMensaje`, `TestPantallaSprintsAnteriores_Errores`, `TestPantallaSprintsAnteriores_EnlazadaDesdeElProyectoYLosSprints` | `http.pantallas.sprintsAnteriores`, [sprints_anteriores.html](../internal/adapters/http/web/plantillas/sprints_anteriores.html) |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-12.1 Un sprint anterior es uno Cerrado | CA-12.1, CA-12.2, CA-12.4 | `domain.SprintsAnteriores` |
| RN-12.2 Se muestran número, fechas, Sprint Goal y los puntos que registró HU-11, sin recalcular | CA-12.1 | `domain.SprintsAnteriores` (devuelve el sprint completo); `PuntosPlanificados` y `PuntosCompletados` de HU-11 |
| RN-12.3 Del más reciente al más antiguo, por fecha de inicio | CA-12.1, CA-12.3 | `domain.SprintsAnteriores` |

Sin CA propio: `TestSprintsAnteriores_NoModificaLaLista` (ordenar no altera la lista que usa la pantalla de HU-08).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/sprint_historial.go](../internal/domain/sprint_historial.go) (`SprintsAnteriores`) | [sprint_historial_test.go](../internal/domain/sprint_historial_test.go) |
| Caso de uso | [internal/app/listar_sprints_anteriores.go](../internal/app/listar_sprints_anteriores.go) (reutiliza `LectorProyectos` y `RepositorioSprints` de HU-08) | [listar_sprints_anteriores_test.go](../internal/app/listar_sprints_anteriores_test.go) |
| Persistencia | Sin cambios: `ListarSprints` de HU-08 ya lee los puntos de HU-11 | Escenarios BDD sobre SQLite real |
| HTTP | [internal/adapters/http/sprints_anteriores_handler.go](../internal/adapters/http/sprints_anteriores_handler.go) (`GET /proyectos/{id}/sprints/anteriores`), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [sprints_anteriores_handler_test.go](../internal/adapters/http/sprints_anteriores_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [pantallas_sprints.go](../internal/adapters/http/pantallas_sprints.go) (`sprintsAnteriores`), [pantallas.go](../internal/adapters/http/pantallas.go) (ruta), [sprints_anteriores.html](../internal/adapters/http/web/plantillas/sprints_anteriores.html); enlaces en [proyecto_detalle.html](../internal/adapters/http/web/plantillas/proyecto_detalle.html) y [sprints.html](../internal/adapters/http/web/plantillas/sprints.html) | [pantallas_sprints_anteriores_test.go](../internal/adapters/http/pantallas_sprints_anteriores_test.go) |
| BDD | [test/bdd/hu12_sistema_test.go](../test/bdd/hu12_sistema_test.go) (SQLite real), [hu12_steps_test.go](../test/bdd/hu12_steps_test.go) | 7 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 12.1 Spec | `e85b947` docs(HU-12): spec de consultar sprints anteriores | — |
| 12.2 Escenarios BDD | `35328d3` | — |
| 12.3 Steps de godog | `81b0806` [RED] (pendientes hasta el caso de uso) | BDD en rojo |
| 12.4 Dominio | `6cfd858` [RED] → `1c65c71` [GREEN] → `eb5de8b` [REFACTOR] | **RED → GREEN → REFACTOR** |
| 12.5 Caso de uso | `1708835` [RED] → `ed73d9f` [GREEN] | **RED → GREEN** (sin REFACTOR: tres pasos lineales con el mismo molde que `ListarSprints` de HU-08) |
| BDD conectado al caso de uso | `e9cef74` | BDD en verde |
| 12.6 Persistencia | — | No hizo falta código nuevo |
| 12.7 Endpoint JSON | `1aa8148` [RED] → `704db13` [GREEN] | **RED → GREEN** |
| 12.8 Pantalla | `f59ce16` [RED] → `902574f` [GREEN] | **RED → GREEN** |
| 12.9 Revisión y trazabilidad | este commit | — |

### Huecos

- No pasó por el agente `go-reviewer`; la revisión queda para el PR.
- Un sprint cerrado antes de que existiera el resumen de HU-11 (base vieja) aparece con 0 y 0 puntos, que es lo que dejó la migración; está documentado como caso límite en la spec, sin test propio.

---

## HU-14: Registrar estimaciones individuales ocultas

- **Issue:** [#25](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/25) · Sprint 3 · 5 SP · tareas #197 a #206
- **Spec:** [docs/specs/HU-14-votos-ocultos.md](specs/HU-14-votos-ocultos.md)
- **Feature:** [features/HU-14-votos-ocultos.feature](../features/HU-14-votos-ocultos.feature) (steps en [test/bdd/hu14_steps_test.go](../test/bdd/hu14_steps_test.go); el resumen se verifica con el endpoint real y los escenarios de CA-14.9 usan la pantalla real)

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-14.1 Iniciar una ronda | RN-14.1, RN-14.2 | Iniciar una ronda sobre un elemento pendiente `@normal` | `TestNuevaRonda_AbiertaConLosIntegrantesComoVotantes`, `TestIniciarRonda_GuardaLaRondaAbiertaConLosIntegrantes`, `TestRepositorioRondas_GuardarYBuscar`, `TestPostRonda_Devuelve201ConElResumen` | `domain.NuevaRonda`, `app.IniciarRonda`, `persistencia.RepositorioRondas.GuardarRonda` |
| CA-14.2 Voto oculto | RN-14.3, RN-14.6 | Un voto no muestra su valor mientras la ronda está abierta `@normal` | `TestVotar_UnVotoNoCierraLaRondaNiMuestraElValor`, `TestVotar_GuardaElVotoYDevuelveElResumenSinValores`, `TestPostVoto_Devuelve200SinElValor` | `domain.Ronda.Votar`, `domain.Ronda.Resumen`, `app.Votar` |
| CA-14.3 Cierre al votar todos | RN-14.6 | La ronda se cierra cuando votaron todos `@normal` · Con un solo integrante, su voto cierra la ronda `@limite` | `TestVotar_CuandoVotanTodosLaRondaSeCierra`, `TestRepositorioRondas_GuardarVoto`, `TestRepositorioRondas_VotosSimultaneos` | `domain.Ronda.Votar`, `persistencia.RepositorioRondas.GuardarVoto` |
| CA-14.4 Cambiar el voto | RN-14.5, RN-14.2 | Cambiar el voto mientras la ronda está abierta `@alternativo` · El nombre del votante no distingue mayúsculas ni espacios `@alternativo` | `TestVotar_CambiarElVotoYNombreSinMayusculas`, `TestRepositorioRondas_GuardarVoto` | `domain.Ronda.Votar` (`ClaveNombre`) |
| CA-14.5 Voto inválido, no votante, ronda cerrada | RN-14.2, RN-14.4, RN-14.6 | Un voto fuera de la escala de Fibonacci `@error` (×4) · Vota alguien que se sumó al proyecto después de abrir la ronda `@error` · No se vota en una ronda cerrada `@error` | `TestVotar_Errores` (dominio y caso de uso), `TestVotar_RondaCerrada`, `TestRepositorioRondas_GuardarVoto_RondaYaCerrada`, `TestVotar_UsaLaEscalaDeLosStoryPoints` | `domain.Ronda.Votar`, `persistencia.RepositorioRondas.GuardarVoto` |
| CA-14.6 No se puede iniciar | RN-14.1, RN-14.2 | No se estima un elemento terminado `@error` · Un proyecto sin integrantes no puede abrir una ronda `@error` · Un elemento tiene una sola ronda abierta a la vez `@error` · Con la ronda anterior cerrada se abre la número 2 `@limite` | `TestNuevaRonda_Errores`, `TestNuevaRonda_ConLaAnteriorCerrada_TomaElNumeroSiguiente`, `TestIniciarRonda_Errores`, `TestIniciarRonda_ConLaAnteriorCerrada_AbreLaSiguiente`, `TestRepositorioRondas_UltimaRonda`, `TestRepositorioRondas_UnaSolaAbiertaPorElemento` | `domain.NuevaRonda`, índice único `rondas_una_abierta_por_elemento` |
| CA-14.7 Otro proyecto o inexistente | RN-14.7 | No se abre una ronda sobre un elemento de otro proyecto `@error` · No se vota en una ronda de otro proyecto `@error` | `TestIniciarRonda_Errores`, `TestVotar_Errores`, `TestConsultarRonda`, `TestConsultarRonda_UltimaDelElemento`, `TestPantallaRonda_Inexistente` | `buscarElementoDelProyecto`, `rondaDelProyecto` |
| CA-14.8 API sin valores | Sección 3, RN-14.3 | El resumen se pide al endpoint real en los escenarios de CA-14.2 y CA-14.3 ("el resumen de la ronda no muestra ningún valor") | `TestPostRonda_Devuelve201ConElResumen`, `TestPostVoto_Devuelve200SinElValor`, `TestGetRonda_Devuelve200ConElResumen`, `TestPlanningPoker_Errores` (17 casos) | `planning_poker_handler.go` |
| CA-14.9 Pantalla | Sección 3 | Iniciar una ronda desde la pantalla del elemento `@normal` · Votar desde la pantalla sin ver los votos `@normal` | `TestPantallaElemento_*` (4 tests de la sección Planning Poker), `TestPantallaIniciarRonda_*`, `TestPantallaRonda_*`, `TestPantallaVotar_*` | `pantallas_planning_poker.go`, `web/plantillas/ronda.html`, sección en `web/plantillas/elemento.html` (HU-06) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/planning_poker.go](../internal/domain/planning_poker.go) (`Ronda`, `NuevaRonda`, `Votar`, `Resumen`; la escala es `EscalaStoryPoints()` de HU-13) | [planning_poker_test.go](../internal/domain/planning_poker_test.go) |
| Casos de uso | [internal/app/planning_poker.go](../internal/app/planning_poker.go) (`IniciarRonda`, `Votar`, `ConsultarRonda`) | [planning_poker_test.go](../internal/app/planning_poker_test.go) |
| Persistencia | [internal/adapters/persistencia/ronda_repo.go](../internal/adapters/persistencia/ronda_repo.go), [db.go](../internal/adapters/persistencia/db.go) (tablas `rondas` y `ronda_votantes`) | [ronda_repo_test.go](../internal/adapters/persistencia/ronda_repo_test.go), métodos de rondas en [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) |
| HTTP (JSON) | [internal/adapters/http/planning_poker_handler.go](../internal/adapters/http/planning_poker_handler.go) (`POST /proyectos/{id}/backlog/{elementoID}/rondas`, `GET /proyectos/{id}/rondas/{rondaID}`, `POST /proyectos/{id}/rondas/{rondaID}/votos`) | [planning_poker_handler_test.go](../internal/adapters/http/planning_poker_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [internal/adapters/http/pantallas_planning_poker.go](../internal/adapters/http/pantallas_planning_poker.go), [web/plantillas/ronda.html](../internal/adapters/http/web/plantillas/ronda.html) (`/web/proyectos/{id}/rondas/{rondaID}`) | [pantallas_planning_poker_test.go](../internal/adapters/http/pantallas_planning_poker_test.go) |
| BDD | [test/bdd/hu14_sistema_test.go](../test/bdd/hu14_sistema_test.go), [hu14_steps_test.go](../test/bdd/hu14_steps_test.go) | 20 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 14.1 Spec | `504cb67` | — |
| 14.2 Escenarios BDD | `beded52` | — |
| 14.3 Steps godog | `f092518` | BDD en rojo (pendientes) |
| 14.4 Ronda en el dominio | `bf87efe` [RED] → `b9b2742` [GREEN] | **RED → GREEN** |
| 14.5 Casos de uso | `61417a0` [RED] → `d84dbdd` [GREEN] | **RED → GREEN** |
| 14.6 Persistencia | `fc5c9ac` [RED] → `cc1b80e` [GREEN] | **RED → GREEN** |
| Conectar los escenarios BDD con los casos de uso reales | `2274c5b` | BDD en verde salvo API y pantalla |
| 14.7 HTTP | `516e1b0` [RED] → `b769bc6` [GREEN] | **RED → GREEN** |
| 14.8 Pantalla | `d4d4bcd` [RED] → `f7c0ca5` [GREEN] | **RED → GREEN**, BDD en verde |
| 14.9 Revisión con go-reviewer | `65613de` | **REFACTOR** |

### Huecos

- La spec está en **Borrador** hasta la aprobación del PR.
- Revisión con el agente `go-reviewer` (tarea 14.9): sin hallazgos críticos ni altos. Confirmó que ningún adaptador muestra el valor de un voto, la concurrencia de `GuardarRonda` y `GuardarVoto` y el `ON DELETE CASCADE`. Se corrigió que `GuardarVoto` informaba "no es votante" cuando la ronda no traía el voto aplicado, y se documentó el choque de `UNIQUE (elemento_id, numero)`. Quedan como en el resto de las historias: las rutas de pantalla suponen que los casos de uso están cargados, y el cuerpo JSON de un voto no rechaza contenido después del primer objeto.
- Mostrar los votos y detectar diferencias es HU-15; reiniciar la ronda, HU-16; y guardar la estimación acordada, HU-17.
- La pantalla no sabe quién está votando (no hay autenticación): el integrante se elige de una lista.

---

## HU-15: Mostrar las estimaciones y detectar diferencias

- **Issue:** [#26](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/26) · Sprint 3 · 3 SP · tareas #207 a #215
- **Spec:** [docs/specs/HU-15-revelar-estimaciones.md](specs/HU-15-revelar-estimaciones.md)
- **Feature:** [features/HU-15-revelar-estimaciones.feature](../features/HU-15-revelar-estimaciones.feature) (steps en [test/bdd/hu15_steps_test.go](../test/bdd/hu15_steps_test.go), sobre el mismo mundo de HU-14: la ronda se arma y se vota con sus pasos)
- **Cambio en HU-14:** RN-14.3, CA-14.8 y CA-14.9 pasan a decir que los valores se ocultan *mientras la ronda está abierta*. Se ajustaron dos escenarios de HU-14 y el test `TestPantallaRonda_Cerrada_SinFormulario`.

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-15.1 Revelar con diferencia | RN-15.1, RN-15.2 | Al cerrarse la ronda se revelan los votos con una diferencia significativa `@normal` | `TestRevelar_RondaCerrada_DevuelveLosVotosElMinimoYElMaximo`, `TestResumen_RondaCerrada_TraeLaRevelacion`, `TestConsultarRonda_RevelaSoloConLaRondaCerrada`, `TestRepositorioRondas_RondaCerradaSeRevelaConTodosLosVotos` | `domain.Ronda.Revelar`, `domain.Ronda.Resumen` |
| CA-15.2 Sin diferencia | RN-15.2, RN-15.3 | Sin diferencia significativa `@alternativo` (×2) | `TestRevelar_DiferenciaSignificativa` | `domain.Ronda.Revelar` |
| CA-15.3 Límites del umbral | RN-15.2, RN-15.3 | Los límites del umbral `@limite` (×3) · Con un solo votante no hay diferencia `@limite` | `TestRevelar_DiferenciaSignificativa` (8 casos) | `domain.Ronda.Revelar` |
| CA-15.4 Ronda abierta sin revelación | RN-15.1 | Con la ronda abierta no se revela nada `@alternativo` | `TestRevelar_RondaAbierta_RetornaError`, `TestConsultarRonda_RevelaSoloConLaRondaCerrada`, `TestGetRonda_Abierta_SinRevelacion`, `TestPantallaRonda_Abierta_SinRevelacion` | `domain.Ronda.Revelar` (`ErrRondaAbierta`) |
| CA-15.5 JSON | Sección 3 | La API muestra la revelación de una ronda cerrada `@normal` | `TestGetRonda_Cerrada_TraeLaRevelacion`, `TestGetRonda_Abierta_SinRevelacion` | `planning_poker_handler.go` (`revelacion`) |
| CA-15.6 Pantalla | Sección 3 | La pantalla de una ronda cerrada muestra los votos y el aviso `@normal` | `TestPantallaRonda_Cerrada_MuestraLosVotosYLaDiferencia`, `TestPantallaRonda_Cerrada_SinDiferencia`, `TestPantallaRonda_Abierta_SinRevelacion` | `web/plantillas/ronda.html` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/planning_poker.go](../internal/domain/planning_poker.go) (`Revelar`, `Revelacion`, `VotoRevelado`; `Resumen` la incluye con la ronda cerrada) | [planning_poker_test.go](../internal/domain/planning_poker_test.go) |
| Caso de uso | Sin cambios: `ConsultarRonda` (HU-14) devuelve el resumen | [planning_poker_test.go](../internal/app/planning_poker_test.go) |
| Persistencia | Sin cambios: `BuscarRonda` (HU-14) ya lee los votos | [ronda_repo_test.go](../internal/adapters/persistencia/ronda_repo_test.go) |
| HTTP (JSON) | [internal/adapters/http/planning_poker_handler.go](../internal/adapters/http/planning_poker_handler.go) (`revelacion` en `GET /proyectos/{id}/rondas/{rondaID}`) | [planning_poker_handler_test.go](../internal/adapters/http/planning_poker_handler_test.go) |
| Pantalla | [web/plantillas/ronda.html](../internal/adapters/http/web/plantillas/ronda.html) | [pantallas_planning_poker_test.go](../internal/adapters/http/pantallas_planning_poker_test.go) |
| BDD | [test/bdd/hu15_steps_test.go](../test/bdd/hu15_steps_test.go) | 10 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 15.1 Spec (y ajuste de HU-14) | `3a75565` | — |
| 15.2 Escenarios BDD (y ajuste de los de HU-14) | `c55fcad` | — |
| 15.3 Steps godog | `72a495c` | BDD en rojo |
| 15.4 Revelar en el dominio | `32bf063` [RED] → `c3075de` [GREEN] | **RED → GREEN** |
| 15.5 Caso de uso | `e46a1bf` | test, sin cambios en el código (ya funcionaba por el dominio) |
| 15.6 Persistencia | `a0f7a6b` | test, sin persistencia nueva |
| 15.7 HTTP | `dbbe2dc` [RED] → `736daa5` [GREEN] | **RED → GREEN** |
| 15.8 Pantalla | `329adbf` [RED] → `ecf26ca` [GREEN] | **RED → GREEN**, BDD en verde |
| 15.9 Revisión y trazabilidad | este commit (go-reviewer sin cambios) | — |

### Huecos

- La spec está en **Borrador** hasta la aprobación del PR.
- Las tareas 15.5 y 15.6 no tienen RED → GREEN: el caso de uso y la persistencia de HU-14 ya alcanzaban, así que sus tests pasaron desde el principio. Se dejaron como evidencia de que la revelación llega hasta la API y la base.
- La ruta sugerida por la tarea 15.7 (`GET /proyectos/{id}/backlog/{elementoID}/rondas/{rondaID}`) no se agregó: la revelación va en el `GET /proyectos/{id}/rondas/{rondaID}` de HU-14, para no tener dos rutas para la misma ronda.
- Revisión con el agente `go-reviewer` (tarea 15.9): sin hallazgos críticos, altos ni medios. Confirmó que con la ronda abierta ningún adaptador revela valores y que el mínimo, el máximo y el umbral son correctos. Tres observaciones bajas sin cambios: el `continue` para un votante sin voto es inalcanzable (la ronda se cierra cuando votaron todos), una ronda sin votantes no se puede crear (RN-14.2) y un comentario se reescribió.

---

## HU-22: Calcular story points y velocidad del equipo

- **Issue:** [#33](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/33) · Sprint 3 · 5 SP
- **Spec:** [docs/specs/HU-22-calcular-velocidad.md](specs/HU-22-calcular-velocidad.md)
- **Feature:** [features/HU-22-calcular-velocidad.feature](../features/HU-22-calcular-velocidad.feature) (steps en [test/bdd/hu22_steps_test.go](../test/bdd/hu22_steps_test.go))
- **Alcance real:** HU-11 ya calcula y guarda los puntos planificados y completados de cada Sprint al cerrarlo (RN-11.6), y HU-12 ya los muestra (pantalla y JSON). HU-22 agrega solo la **velocidad** (el promedio de puntos completados entre los últimos 3 Sprints cerrados); no agrega persistencia nueva.

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-22.1 Planificados vs. completados de un Sprint | RN-22.1 | No va acá: ya lo cubren el endpoint y la pantalla de HU-12 | `TestSprintsAnteriores_ConservaLosDatos` (HU-12) | Ya implementado por HU-11 y HU-12; HU-22 no agrega código |
| CA-22.2 Velocidad = promedio de los cerrados | RN-22.2 | La velocidad es el promedio de los puntos completados de los sprints cerrados `@normal` · Un sprint todavía abierto no cuenta `@alternativo` · Solo cuentan los sprints cerrados del proyecto consultado `@alternativo` | `TestVelocidad`, `TestVelocidad_NoFiltraPorEstado`, `TestConsultarVelocidad_PromedioDeLosCerrados` | `domain.Velocidad`, `app.ConsultarVelocidad` |
| CA-22.3 Sin sprints cerrados, sin velocidad | RN-22.3 | Un proyecto sin ningún sprint no tiene velocidad `@limite` · Un proyecto con sprints, pero ninguno cerrado, no tiene velocidad `@alternativo` | `TestVelocidad` (caso "sin sprints cerrados"), `TestConsultarVelocidad_SinCerrados_NoHayDatos` | `domain.Velocidad` |
| CA-22.4 Con decimales | RN-22.4 | La velocidad se informa con decimales cuando no divide exacto `@limite` | `TestVelocidad` (caso decimales), `TestConsultarVelocidad_ConDecimales` | `domain.Velocidad` |
| CA-22.5 Un sprint en 0 da velocidad 0, no "sin datos" | RN-22.2, RN-22.3 | Un único sprint cerrado con 0 puntos completados da velocidad 0, no "sin datos" `@limite` | `TestVelocidad` (caso 0), `TestConsultarVelocidad_UnSoloSprintEnCero_DaVelocidadCero` | `domain.Velocidad` |
| CA-22.6 Proyecto inexistente | — (error reutilizado) | No se puede consultar la velocidad de un proyecto que no existe `@error` | `TestConsultarVelocidad_ProyectoInexistente_RetornaError` | `app.ConsultarVelocidad` |
| CA-22.7 Respuesta HTTP | — | No va acá (como CA-18.15): se cubre con tests del adaptador | `TestGetVelocidad_*` (5 tests) | `internal/adapters/http/velocidad_handler.go` |
| CA-22.8 Pantalla | Restricciones | No va acá (como CA-06.10): se cubre con tests de la pantalla | `TestPantallaSprintsAnteriores_ConSprintsCerrados_MuestraLaVelocidad`, `TestPantallaSprintsAnteriores_SinSprintsCerrados_AvisaSinVelocidad` | `pantallas_sprints.go` (`sprintsAnteriores`, `formatearVelocidad`), `web/plantillas/sprints_anteriores.html` |
| CA-22.9 Solo los últimos 3 sprints cerrados entran en el promedio | RN-22.2 | Con más de 3 sprints cerrados, solo se promedian los 3 más recientes `@limite` | `TestVelocidad` (caso ventana) | `domain.Velocidad` (`VentanaVelocidad`) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/sprint_historial.go](../internal/domain/sprint_historial.go) (`Velocidad`, `VentanaVelocidad`, junto a `SprintsAnteriores` de HU-12) | [sprint_historial_test.go](../internal/domain/sprint_historial_test.go) |
| Caso de uso | [internal/app/consultar_velocidad.go](../internal/app/consultar_velocidad.go) | [consultar_velocidad_test.go](../internal/app/consultar_velocidad_test.go) |
| HTTP (JSON) | [internal/adapters/http/velocidad_handler.go](../internal/adapters/http/velocidad_handler.go) (`GET /proyectos/{id}/velocidad`) | [velocidad_handler_test.go](../internal/adapters/http/velocidad_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_sprints.go](../internal/adapters/http/pantallas_sprints.go) (reutiliza la pantalla de HU-12), [web/plantillas/sprints_anteriores.html](../internal/adapters/http/web/plantillas/sprints_anteriores.html) | [pantallas_sprints_anteriores_test.go](../internal/adapters/http/pantallas_sprints_anteriores_test.go) |
| BDD | [test/bdd/hu22_sistema_test.go](../test/bdd/hu22_sistema_test.go), [hu22_steps_test.go](../test/bdd/hu22_steps_test.go) | 9 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 22.1 Spec | `c9c046b` docs(HU-22): [22.1] spec SDD | — |
| 22.2 Escenarios BDD | `60c2c24` test(HU-22): [22.2] escenarios BDD | — |
| 22.3 Steps godog | `7f9c2ad` test(HU-22): [22.3][RED] steps (en rojo) | BDD en rojo |
| 22.4 Velocidad (dominio) | `457d6a3` [RED] → `7577e40` [GREEN] | **RED → GREEN** |
| 22.5 Caso de uso ConsultarVelocidad | `d327cbb` [RED] → `28771ad` [GREEN] | **RED → GREEN** |
| Conectar los escenarios BDD con el caso de uso real | `9370946` | BDD en verde |
| 22.7 HTTP | `56dfdb1` [RED] → `9dedf2a` [GREEN] | **RED → GREEN** |
| 22.8 Pantalla (sumada a la de HU-12) | `372945a` [RED] → `203fb56` [GREEN] | **RED → GREEN** |
| 22.9 Revisión y trazabilidad | `83f6053` | — |
| 22.10 Corrección de la revisión del PR (ventana de 3) | `a107acb` (spec) → `b1c51c1` (BDD) → `c4a7558` [RED] → `3514f64` [GREEN] | **RED → GREEN** |

### Huecos

- ~~La spec está en Borrador, P-1 pendiente de confirmar~~: resuelto el 2026-11-02. Matías López (Agile Enabler) pidió en la revisión del PR que la velocidad promedie solo los últimos 3 Sprints cerrados (ventana móvil), no todos; se corrigió la spec (RN-22.2, CA-22.9 nuevo) y el código (tarea 22.10). Spec Aprobada.
- No hay pantalla nueva: la velocidad se agregó a la pantalla ya existente de HU-12 (`/web/proyectos/{id}/sprints/anteriores`), como se acordó al refinar la historia.
- Revisión con el agente `go-reviewer` (tarea 22.9): sin hallazgos bloqueantes. Confirmó la lógica del promedio, la distinción `null` vs. `0` (RN-22.3) y que no hay lógica de negocio en los adaptadores. Notas menores sin cambios: el handler de la pantalla repite la consulta del proyecto y de los sprints por usar dos casos de uso separados (`ListarSprintsAnteriores` y `ConsultarVelocidad`), que es la separación que pide la propia spec; y el ejemplo de la spec (`6.0` en vez de `6`) se corrigió en ese commit.
- `sprints_cerrados` en la respuesta JSON y en `VelocidadProyecto.SprintsCerrados` es cuántos Sprints entraron en el promedio (hasta 3), no el total de Sprints cerrados del proyecto: se aclaró explícitamente en la spec (sección 3) y en el comentario del campo, porque con la ventana podían confundirse los dos números (tarea 22.10).

---

## HU-16: Repetir una ronda de votación

- **Issue:** #27 · Sprint 3 · 2 SP
- **Spec:** [docs/specs/HU-16-nueva-ronda.md](specs/HU-16-nueva-ronda.md)
- **Feature:** [features/HU-16-nueva-ronda.feature](../features/HU-16-nueva-ronda.feature) (steps en [test/bdd/hu16_steps_test.go](../test/bdd/hu16_steps_test.go), sobre el mundo de HU-14)

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-16.1 Con la ronda cerrada se abre la siguiente; la anterior conserva sus votos | RN-16.1, RN-16.2 | Después de una ronda cerrada se abre la siguiente y la anterior conserva sus votos `@normal` | `TestListarRondas_DevuelveElHistorialDelElemento`, `TestRepositorioRondas_ListarRondas` | `domain.NuevaRonda` (HU-14), `app.ListarRondas`, `persistencia.RepositorioRondas.ListarRondas` |
| CA-16.2 Se repite aunque no haya habido diferencias | RN-16.4 | Se puede repetir la ronda aunque no haya habido diferencias `@alternativo` | `TestPuedeRepetirse` (sin diferencia), `TestPantallaRonda_NuevaRonda` (sin diferencia) | `domain.PuedeRepetirse` |
| CA-16.3 Con una ronda abierta no se abre otra | RN-16.1 | No se puede abrir otra ronda mientras hay una abierta `@error` | Los de HU-14 (`ErrRondaYaAbierta`) | `domain.NuevaRonda` (HU-14) |
| CA-16.4 Historial de la más reciente a la más antigua, con votos de las cerradas y sin valores de la abierta | RN-16.3 | El historial muestra todas las rondas, de la más reciente a la más antigua `@normal` | `TestHistorialRondas_OrdenYRevelacion`, `TestHistorialRondas_NoModificaLasRondas`, `TestListarRondas_DevuelveElHistorialDelElemento` | `domain.HistorialRondas` (usa `Ronda.Resumen` de HU-14 y HU-15) |
| CA-16.5 Sin rondas, historial vacío | RN-16.3 | Un elemento sin rondas tiene el historial vacío `@limite` | `TestHistorialRondas_SinRondas`, `TestListarRondas_SinRondas_ListaVacia`, `TestRepositorioRondas_ListarRondas_SinRondas` | `domain.HistorialRondas` |
| CA-16.6 Proyecto o elemento inexistente, o elemento de otro proyecto | RN-16.5 | No se puede consultar el historial desde un proyecto que no existe `@error` · …de un elemento de otro proyecto `@error` | `TestListarRondas_Errores`, `TestListarRondas_FallaDelRepositorio_SePropaga` | `app.ListarRondas.Ejecutar` (`buscarElementoDelProyecto`) |
| CA-16.7 HTTP 200 / 404 / 400 | Sección 3 | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y la pantalla. | `TestGetHistorialRondas_DevuelveLasRondas`, `TestGetHistorialRondas_SinRondas_DevuelveArrayVacio`, `TestGetHistorialRondas_Errores`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerHistorialRondas` (reutiliza `aRespuestaRonda` de HU-14), `http.RegistrarHistorialRondas` |
| CA-16.8 "Nueva ronda" en la última ronda cerrada; historial en el elemento | RN-16.3, RN-16.4 | Desde la pantalla de la ronda cerrada se abre la siguiente `@normal` | `TestPuedeRepetirse`, `TestPantallaRonda_NuevaRonda`, `TestPantallaRonda_CerradaPeroNoEsLaUltima_SinNuevaRonda`, `TestPantallaRonda_ElementoHecho_SinNuevaRonda`, `TestPantallaElemento_HistorialDeRondas`, `TestPantallaElemento_SinRondas_SinHistorial`, `TestPantallaElemento_FallaDelHistorial` | `http.pantallas.mostrarRonda` (`PuedeRepetirse`), `http.pantallas.pokerDelElemento` (`Historial`), [ronda.html](../internal/adapters/http/web/plantillas/ronda.html), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-16.1 Ronda nueva con la última cerrada y el elemento no Hecho; no con una abierta | CA-16.1, CA-16.3 | `domain.NuevaRonda` de HU-14 (HU-16 no la cambia) |
| RN-16.2 Abrir una ronda nueva no modifica las anteriores | CA-16.1 | `domain.HistorialRondas` no modifica lo recibido; la persistencia de HU-14 inserta, no actualiza |
| RN-16.3 Historial de la más reciente a la más antigua, sin valores de la abierta | CA-16.4, CA-16.5, CA-16.8 | `domain.HistorialRondas` |
| RN-16.4 "Nueva ronda" solo en la última ronda, cerrada, de un elemento no Hecho, con o sin diferencia | CA-16.2, CA-16.8 | `domain.PuedeRepetirse` |
| RN-16.5 Un elemento de otro proyecto se trata como inexistente | CA-16.6 | `app.ListarRondas.Ejecutar` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/planning_poker_historial.go](../internal/domain/planning_poker_historial.go) (`HistorialRondas`, `PuedeRepetirse`) | [planning_poker_historial_test.go](../internal/domain/planning_poker_historial_test.go) |
| Caso de uso y puerto | [internal/app/listar_rondas.go](../internal/app/listar_rondas.go) (`ListarRondas`, puerto angosto `LectorHistorialRondas`) | [listar_rondas_test.go](../internal/app/listar_rondas_test.go) (envuelve el fake de HU-14 sin tocarlo) |
| Persistencia | [internal/adapters/persistencia/ronda_historial.go](../internal/adapters/persistencia/ronda_historial.go) (`ListarRondas`, reutiliza `BuscarRonda` de HU-14) | [ronda_historial_test.go](../internal/adapters/persistencia/ronda_historial_test.go), [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) |
| HTTP | [internal/adapters/http/historial_rondas_handler.go](../internal/adapters/http/historial_rondas_handler.go) (`GET /proyectos/{id}/backlog/{elementoID}/rondas`), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [historial_rondas_handler_test.go](../internal/adapters/http/historial_rondas_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [pantallas_planning_poker.go](../internal/adapters/http/pantallas_planning_poker.go) (`PuedeRepetirse` en la ronda, `Historial` en el elemento), [ronda.html](../internal/adapters/http/web/plantillas/ronda.html), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) | [pantallas_historial_rondas_test.go](../internal/adapters/http/pantallas_historial_rondas_test.go) |
| BDD | [test/bdd/hu16_sistema_test.go](../test/bdd/hu16_sistema_test.go) (método nuevo del sistema de HU-14, en archivo propio), [hu16_steps_test.go](../test/bdd/hu16_steps_test.go) | 8 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 16.1 Spec | `5d60f4e` | — |
| 16.2 Escenarios BDD | `1ab6a38` | — |
| 16.3 Steps de godog | `7ea5111` [RED] (pendientes hasta el caso de uso) | BDD en rojo |
| 16.4 Dominio | `5030b7a` [RED] → `75ab7bc` [GREEN] | **RED → GREEN** (sin REFACTOR: dos funciones cortas) |
| 16.5 Caso de uso | `4e2c13f` [RED] → `afe8886` [GREEN] | **RED → GREEN** (sin REFACTOR: mismo molde que HU-14) |
| 16.6 Persistencia | `72ecf17` [RED] → `64cee06` [GREEN] · `efb6090` (base cerrada) | **RED → GREEN** |
| BDD conectado al caso de uso | `35e6f9e` | BDD en verde salvo la pantalla |
| 16.7 Endpoint JSON | `cac55e2` [RED] → `aa49944` [GREEN] | **RED → GREEN** |
| 16.8 Pantalla | `96cf6f9` [RED] → `0485622` [GREEN] | **RED → GREEN**; BDD 8/8 |
| 16.9 Revisión y trazabilidad | `c391460` | — |

### Huecos

- No pasó por el agente `go-reviewer`; la revisión queda para el PR.
- Cambios mínimos en archivos de HU-14 (Nico): registrar los steps de HU-16 en `hu14_steps_test.go`, sumar `ListarRondas` al router del sistema BDD de HU-14 y el campo `Historial` en `datosPokerElemento`. Conviene que lo vea Nico en el PR.
- Si "Nueva ronda" falla (por ejemplo, otro pedido abrió la ronda en el medio), el error se muestra en la pantalla del elemento y no en la de la ronda, porque el botón usa la acción de HU-14 (anotado en la spec).

---

## HU-17: Registrar la estimación acordada

- **Issue:** #28 · Sprint 3 · 2 SP
- **Spec:** [docs/specs/HU-17-estimacion-acordada.md](specs/HU-17-estimacion-acordada.md)
- **Feature:** [features/HU-17-estimacion-acordada.feature](../features/HU-17-estimacion-acordada.feature) (steps en [test/bdd/hu17_steps_test.go](../test/bdd/hu17_steps_test.go), sobre el mundo de HU-14)

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-17.1 La estimación queda en la ronda y como story points del elemento | RN-17.1 | La estimación queda en la ronda y como story points del elemento `@normal` | `TestRegistrarEstimacion_ValoresDeLaEscala`, `TestRegistrarEstimacion_Valida_LaGuarda`, `TestRepositorioRondas_GuardarEstimacion` | `domain.Ronda.RegistrarEstimacion`, `app.RegistrarEstimacion`, `persistencia.RepositorioRondas.GuardarEstimacion` |
| CA-17.2 Reemplaza los story points y no cambia los demás datos | RN-17.1 | La estimación reemplaza los story points y no cambia los demás datos `@alternativo` | `TestRepositorioRondas_GuardarEstimacion` | `persistencia.RepositorioRondas.GuardarEstimacion` (solo actualiza `story_points`) |
| CA-17.3 Cualquier valor de la escala, aunque nadie lo haya votado | RN-17.2 | Se puede acordar un valor que nadie votó `@alternativo` · El valor tiene que ser de la escala `@error` (3 ejemplos) | `TestRegistrarEstimacion_ValoresDeLaEscala`, `TestRegistrarEstimacion_Errores`, `TestRegistrarEstimacion_Errores_NoGuarda` | `domain.Ronda.RegistrarEstimacion` |
| CA-17.4 No en una ronda abierta, que no es la última o con el elemento Hecho | RN-17.3 | No se registra en una ronda abierta `@error` · …que no es la última `@error` · …si el elemento ya está terminado `@error` | `TestRegistrarEstimacion_Errores`, `TestRegistrarEstimacion_Errores_NoGuarda`, `TestRepositorioRondas_GuardarEstimacion_Errores` | `domain.errorParaRegistrar`, `persistencia.RepositorioRondas.errorDeEstimacion` |
| CA-17.5 Una sola estimación por ronda | RN-17.4 | Una ronda tiene una sola estimación `@error` | `TestRegistrarEstimacion_Errores`, `TestRepositorioRondas_GuardarEstimacion_Errores`, `TestRepositorioRondas_GuardarEstimacion_Simultaneos` | `domain.errorParaRegistrar`; en la base, `UPDATE ... AND estimacion_acordada IS NULL` |
| CA-17.6 Una ronda nueva cambia la estimación y la anterior queda registrada | RN-17.4, RN-17.6 | Una ronda nueva puede cambiar la estimación y la anterior queda registrada `@normal` | `TestRepositorioRondas_GuardarEstimacion_RondaNueva` | `persistencia.RepositorioRondas.GuardarEstimacion` |
| CA-17.7 Editar a mano no cambia lo acordado | RN-17.5 | Editar los story points a mano no cambia lo acordado en la ronda `@alternativo` | — (HU-06 no se modificó; lo cubre el escenario sobre SQLite real) | La estimación se guarda en `rondas`, separada de `backlog_elementos.story_points` |
| CA-17.8 Proyecto o ronda inexistentes, o ronda de otro proyecto | Sección 7 | No se registra en una ronda que no existe `@error` · …desde un proyecto que no existe `@error` · Una ronda de otro proyecto se trata como inexistente `@error` | `TestRegistrarEstimacion_Errores_NoGuarda`, `TestRepositorioRondas_GuardarEstimacion_Errores` | `app.rondaDelProyecto` (HU-14) |
| CA-17.9 HTTP 200 / 400 / 404 | Sección 3 | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y la pantalla. | `TestPostEstimacion_Valida_Devuelve200ConLaEstimacion`, `TestRespuestaRonda_SinEstimacion_NoTraeElCampo`, `TestPostEstimacion_ErroresDelCasoDeUso`, `TestPostEstimacion_PedidoInvalido_NoLlamaAlCasoDeUso`, `TestNuevoRouter_RegistraTodasLasRutas` | `http.handlerRegistrarEstimacion` (reutiliza `responderRonda` de HU-14), `http.RegistrarEstimacionAcordada` |
| CA-17.10 Pantalla de la ronda e historial del elemento | Sección 3, RN-17.3 | Desde la pantalla de la ronda se registra la estimación `@normal` · La pantalla de una ronda abierta no ofrece registrar la estimación `@alternativo` | `TestPuedeRegistrarEstimacion`, `TestPantallaRonda_CerradaSinEstimacion_OfreceRegistrar`, `TestPantallaRonda_NoOfreceRegistrar`, `TestPantallaRonda_ConEstimacion_LaMuestra`, `TestPantallaEstimacion_Valida_RedirigeALaRonda`, `TestPantallaEstimacion_ErroresDeValidacion`, `TestPantallaEstimacion_ValorNoNumerico_NoLlamaAlCasoDeUso`, `TestPantallaEstimacion_OtrosErrores`, `TestPantallaEstimacion_DesdeOtroSitio_SeRechaza`, `TestPantallaElemento_HistorialMuestraLaEstimacion` | `domain.PuedeRegistrarEstimacion`, `http.pantallas.registrarEstimacionPantalla`, [ronda.html](../internal/adapters/http/web/plantillas/ronda.html), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-17.1 El valor se guarda en la ronda y se copia al elemento, en una sola operación | CA-17.1, CA-17.2 | `persistencia.RepositorioRondas.GuardarEstimacion` (transacción) |
| RN-17.2 Cualquier valor de la escala | CA-17.3 | `domain.Ronda.RegistrarEstimacion` |
| RN-17.3 Solo en la última ronda, cerrada, con el elemento no Hecho | CA-17.4, CA-17.10 | `domain.errorParaRegistrar`; en la base, las condiciones de los `UPDATE` |
| RN-17.4 Una estimación por ronda; para cambiarla se vota otra | CA-17.5, CA-17.6 | `domain.errorParaRegistrar`; en la base, `estimacion_acordada IS NULL` |
| RN-17.5 La edición manual no se bloquea y la ronda conserva lo acordado | CA-17.7 | HU-06 sin cambios; el dato vive en `rondas.estimacion_acordada` |
| RN-17.6 Registrar no impide abrir otra ronda | CA-17.6 | `domain.PuedeRepetirse` (HU-16) no mira la estimación |

Sin CA propio: `TestRepositorioRondas_GuardarEstimacion_Simultaneos` (10 pedidos a la vez: se registra uno solo, y la ronda y el elemento quedan con el mismo valor), `TestAbrirDB_BaseExistenteSinEstimacionAcordada_AgregaLaColumna` y `TestResumen_TraeLaEstimacionAcordada`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/planning_poker_estimacion.go](../internal/domain/planning_poker_estimacion.go), campo `EstimacionAcordada` en [planning_poker.go](../internal/domain/planning_poker.go) | [planning_poker_estimacion_test.go](../internal/domain/planning_poker_estimacion_test.go) |
| Caso de uso y puerto | [internal/app/registrar_estimacion.go](../internal/app/registrar_estimacion.go) (`RegistrarEstimacion`, puerto angosto `RepositorioEstimacion`) | [registrar_estimacion_test.go](../internal/app/registrar_estimacion_test.go) |
| Persistencia | [internal/adapters/persistencia/ronda_estimacion.go](../internal/adapters/persistencia/ronda_estimacion.go), lectura en [ronda_repo.go](../internal/adapters/persistencia/ronda_repo.go), columna y migración en [db.go](../internal/adapters/persistencia/db.go) | [ronda_estimacion_test.go](../internal/adapters/persistencia/ronda_estimacion_test.go), [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) |
| HTTP | [internal/adapters/http/estimacion_handler.go](../internal/adapters/http/estimacion_handler.go) (`POST /proyectos/{id}/rondas/{rondaID}/estimacion`), [planning_poker_handler.go](../internal/adapters/http/planning_poker_handler.go), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [estimacion_handler_test.go](../internal/adapters/http/estimacion_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla | [pantallas_planning_poker.go](../internal/adapters/http/pantallas_planning_poker.go), [ronda.html](../internal/adapters/http/web/plantillas/ronda.html), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) | [pantallas_estimacion_test.go](../internal/adapters/http/pantallas_estimacion_test.go) |
| BDD | [test/bdd/hu17_sistema_test.go](../test/bdd/hu17_sistema_test.go), [hu17_steps_test.go](../test/bdd/hu17_steps_test.go) | 17 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 17.1 Spec | `8e92530` · `7c8020e` (aprobada) | — |
| 17.2 Escenarios BDD | `e760c80` | — |
| 17.3 Steps de godog | `1ecb594` [RED] | BDD en rojo |
| 17.4 Dominio | `a6d4122` [RED] → `6cc012c` [GREEN] | **RED → GREEN** (sin REFACTOR: dos líneas repetidas que no vale la pena unificar) |
| 17.5 Caso de uso | `bbec012` [RED] → `0569373` [GREEN] | **RED → GREEN** (sin REFACTOR: mismo molde que HU-14) |
| 17.6 Persistencia | `f9edcec` [RED] → `a6b9c13` [GREEN] | **RED → GREEN** |
| BDD conectado al caso de uso | `fdada2a` | BDD en verde salvo la pantalla |
| 17.7 Endpoint JSON | `f4cca58` [RED] → `8b771f3` [GREEN] | **RED → GREEN** |
| 17.8 Pantalla | `050ff0f` [RED] → `a85d395` [GREEN] | **RED → GREEN**; BDD 17/17 |
| 17.9 Revisión y trazabilidad | este commit | — |

### Huecos

- **Escala de HU-13 (#288):** al combinar las dos historias, `RegistrarEstimacion` valida con `escalaStoryPoints` y la pantalla de la ronda usa `EscalaStoryPoints()`, la escala unificada de HU-13 (RN-13.4), en lugar de `valoresVoto` y `ValoresVoto`.
- No pasó por el agente `go-reviewer`; la revisión queda para el PR.
- Las condiciones de la transacción se prueban una a una y con pedidos simultáneos sobre la misma ronda, pero no hay un test de una estimación simultánea con un cambio del elemento a `Hecho`.
- Cambios mínimos en código de HU-14 (Nico): el campo `EstimacionAcordada`, `rondaDelProyecto` con un parámetro más angosto, la lectura de la columna en `BuscarRonda`, y el registro en el sistema BDD. Conviene que lo vea Nico en el PR.

---

## HU-19: Comparar esfuerzo estimado contra real

- **Issue:** #30 · Sprint 3 · 3 SP
- **Spec:** [docs/specs/HU-19-comparar-esfuerzo.md](specs/HU-19-comparar-esfuerzo.md)
- **Feature:** [features/HU-19-comparar-esfuerzo.feature](../features/HU-19-comparar-esfuerzo.feature) (steps en [test/bdd/hu19_steps_test.go](../test/bdd/hu19_steps_test.go), sobre SQLite real en [hu19_sistema_test.go](../test/bdd/hu19_sistema_test.go))
- **Ruta HTTP:** `GET /proyectos/{id}/backlog/{elementoID}/esfuerzo/comparacion`, y el campo `horas_estimadas` en `PATCH /proyectos/{id}/backlog/{elementoID}`.
- **Pantalla:** bloque "Esfuerzo estimado y real" y campo "Horas estimadas" en `GET /web/proyectos/{id}/backlog/{elementoID}` (HU-06).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-19.1 Desviación positiva | RN-19.4, RN-19.5, RN-19.6 | La desviación se calcula en horas y en porcentaje `@normal` (8/10) · Las horas reales suman todos los registros del elemento `@normal` | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_NoCompartePunteroRecibido`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`, `TestRepositorioEsfuerzo_HorasReales_SumaIntegrantesYFechasDistintos`, `TestRepositorioEsfuerzo_HorasReales_ConservaLaMediaHora` | `comparacion_esfuerzo.go` (`domain.CompararEsfuerzo`), `comparar_esfuerzo.go` (`app.CompararEsfuerzo.Ejecutar`), `esfuerzo_repo.go` (`HorasRealesDeElemento`) |
| CA-19.2 Desviación negativa | RN-19.5, RN-19.6 | Outline de desviación `@normal` (8/6) | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento` | `comparacion_esfuerzo.go` |
| CA-19.3 Sin desviación | RN-19.5, RN-19.6 | Outline de desviación `@limite` (8/8) | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento` | `comparacion_esfuerzo.go` |
| CA-19.4 Sin estimar, con esfuerzo | RN-19.1, RN-19.7 | Un elemento sin estimar y con esfuerzo informa solo sus horas reales `@alternativo` | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento` | `comparacion_esfuerzo.go` (`SinEstimar`) |
| CA-19.5 Sin esfuerzo | RN-19.4, RN-19.9 | Outline de desviación `@limite` (8/0) | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`, `TestRepositorioEsfuerzo_HorasReales_SinRegistros_DaCero`, `TestRepositorioEsfuerzo_HorasReales_ElementoInexistente_DaCero` | `comparacion_esfuerzo.go`, `esfuerzo_repo.go` (`COALESCE`) |
| CA-19.6 Sin estimar y sin esfuerzo | RN-19.4, RN-19.7 | Un elemento sin estimar y sin esfuerzo informa 0 horas reales `@limite` | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`, `TestRepositorioEsfuerzo_HorasReales_SinRegistros_DaCero` | `comparacion_esfuerzo.go`, `esfuerzo_repo.go` |
| CA-19.7 Estimado 0 | RN-19.8 | Con estimado 0 hay desviación en horas y no hay porcentaje `@limite` | `TestCompararEsfuerzo` | `comparacion_esfuerzo.go` |
| CA-19.8 Redondeo | RN-19.6 | Outline de desviación `@limite` (3/4, 40/40,5 y 40/39,5) | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_MultiplicaAntesDeDividir`, `TestCompararEsfuerzo_NuncaDevuelveCeroNegativo` | `comparacion_esfuerzo.go` (`redondearUnDecimal`) |
| CA-19.9 Proyecto o elemento inexistente | RN-19.10 | No se puede comparar en un proyecto que no existe `@error` · No se puede comparar un elemento que no existe `@error` | `TestCompararEsfuerzo_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestCompararEsfuerzo_ProyectoInexistente_NoBuscaElElemento`, `TestCompararEsfuerzo_FallaDeInfraestructura_SePropaga`, `TestGetCompararEsfuerzo_NoEncontrado_Devuelve404`, `TestPantallaComparacion_NoEncontrados_Devuelven404` | `comparar_esfuerzo.go`, `comparar_esfuerzo_handler.go`, `pantallas_elemento.go` (`errorDeElemento`) |
| CA-19.10 Elemento de otro proyecto | RN-19.10 | Un elemento de otro proyecto se trata como inexistente `@error` | `TestCompararEsfuerzo_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestGetCompararEsfuerzo_NoEncontrado_Devuelve404`, `TestPantallaComparacion_NoEncontrados_Devuelven404` | `comparar_esfuerzo.go` (`e.ProyectoID != proyectoID`) |
| CA-19.11 Qué suma | RN-19.4 | Suman los registros de un integrante quitado y no los de otro elemento `@alternativo` · La suma de las horas reales no tiene tope `@limite` | `TestCompararEsfuerzo_UsaLasHorasDelElementoPedido`, `TestRepositorioEsfuerzo_HorasReales_LaSumaPuedeSuperar24`, `TestRepositorioEsfuerzo_HorasReales_NoSumaOtroElemento`, `TestRepositorioEsfuerzo_HorasReales_IntegranteQuitado_SigueSumando` | `esfuerzo_repo.go` (`HorasRealesDeElemento`) |
| CA-19.12 Cualquier estado | RN-19.11 | Se compara un elemento Pendiente y sin sprint de un proyecto que no está Activo `@alternativo` (4 ejemplos) | `TestCompararEsfuerzo_CualquierEstadoDeElementoYProyecto` | `comparar_esfuerzo.go` (no mira estados) |
| CA-19.13 Cargar la estimación | RN-19.3 | Se cargan las horas estimadas editando el elemento `@normal` · Editar otro campo no cambia las horas estimadas `@alternativo` · Se quitan las horas estimadas de un elemento `@alternativo` · Se cambia la estimación de un elemento con esfuerzo y la comparación usa el valor nuevo `@alternativo` · Se quita la estimación de un elemento con esfuerzo y conserva sus horas reales `@alternativo` | `TestEditar_ConservaElSprint`, `TestEditar_CargarHorasEstimadas`, `TestEditar_CambiarHorasEstimadas`, `TestEditar_QuitarHorasEstimadas`, `TestEditar_OtrosCampos_ConservanHorasEstimadasYSprint`, `TestEditar_HorasEstimadas_NoCompartenMemoria`, `TestRepositorioBacklog_HorasEstimadas_IdaYVuelta`, `TestRepositorioBacklog_HorasEstimadas_SeGuardanEnMediasHoras`, `TestRepositorioBacklog_HorasEstimadas_CambiarYQuitar`, `TestRepositorioBacklog_EditarOtrosCampos_NoTocaLasHorasEstimadas`, `TestRepositorioBacklog_HorasEstimadas_EdicionesSimultaneasNoSePisan`, `TestRepositorioBacklog_ListarElementos_TraeLasHorasEstimadas`, `TestRepositorioBacklog_HorasEstimadas_SobrevivenAlEditarConSprint`, `TestRepositorioBacklog_HorasEstimadas_ElementoInexistente`, `TestPatchElemento_HorasEstimadas`, `TestGetBacklog_RespuestaConHorasEstimadas`, `TestPantallaElemento_CargarYQuitarHorasEstimadas`, `TestPantallaElemento_HorasSinCambios_NoSeEnvian`, `TestPantallaElemento_SinOriginales_EnviaLasHoras`, `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran` | `backlog_editar.go` del dominio (`Editar`) y de persistencia (`ActualizarElemento`), `backlog_repo.go`, `backlog_editar_handler.go`, `pantallas_elemento.go` |
| CA-19.14 Validar la estimación | RN-19.2 | Se aceptan las horas estimadas en los extremos válidos `@limite` (0,5 y 1000) · Se rechazan las horas estimadas inválidas `@error` (0, -1, 1,25 y 1000,5) | `TestValidarHorasEstimadas_Validas`, `TestValidarHorasEstimadas_Invalidas`, `TestErrHorasEstimadasExcedenMaximo_Mensaje`, `TestValidarHoras_CadaUnaConSuTope`, `TestEditar_HorasEstimadasEnLosExtremos`, `TestEditar_HorasEstimadasInvalidas_RetornaError`, `TestEditar_HorasEstimadasInvalidas_SeInformanAlFinal`, `TestRepositorioBacklog_HorasEstimadas_IdaYVuelta`, `TestPatchElemento_HorasEstimadasQueNoSonNumero_Devuelve400SinLlamarAlCasoDeUso`, `TestPatchElemento_HorasEstimadasInvalidas_Devuelve400ConElMensajeDelDominio`, `TestPatchElemento_ErrorDeHorasEnvuelto_Devuelve400ConMensajeLimpio`, `TestPantallaElemento_HorasInvalidas_VuelvenAlFormularioCon400`, `TestPantallaElemento_HorasQueNoSonNumero_NoLlamanAlCasoDeUso` | `backlog_editar.go` (`validarHorasEstimadas`), `registro_esfuerzo.go` (`validarMediasHoras`), `backlog_editar_handler.go` (`erroresDeValidacionEdicion`) |
| CA-19.15 Independiente de los story points | RN-19.1, RN-19.3 | Un elemento con story points y sin horas estimadas está sin estimar `@alternativo` · Cambiar los story points no cambia las horas estimadas `@alternativo` · Un elemento recién creado nace sin horas estimadas `@normal` | `TestNuevoElementoBacklog_NaceSinHorasEstimadas`, `TestEditar_OtrosCampos_ConservanHorasEstimadasYSprint`, `TestEditar_HorasEstimadas_ConservanStoryPoints`, `TestRepositorioBacklog_EditarOtrosCampos_NoTocaLasHorasEstimadas`, `TestRepositorioBacklog_EditarHorasEstimadas_NoTocaLosDemasCampos`, `TestRepositorioBacklog_GuardarElemento_GuardaLasHorasEstimadas`, `TestPatchElemento_HorasEstimadasYStoryPoints_SonIndependientes`, `TestPostBacklog_RespuestaConHorasEstimadasEnNull`, `TestPostBacklog_ConHorasEstimadas_Devuelve400`, `TestPantallaElemento_HorasYStoryPoints_NoSePisan`, `TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06` | `backlog.go` (`ElementoBacklog.HorasEstimadas`), `backlog_editar.go` (`CambiosElemento`) |
| CA-19.16 Base existente | Sección 5 | **Sin escenario BDD a propósito:** se cubre con tests de persistencia. | `TestAbrirDB_BaseExistenteSinHorasEstimadas_AgregaLaColumna`, `TestAbrirDB_BaseNueva_TieneLaColumnaDeHorasEstimadas`, `TestMigracionHU19_ElementosSinEstimarYConSusDatos`, `TestMigracionHU19_SeCarganHorasYSeComparan`, `TestMigracionHU19_LaSumaDeHorasRealesNoCambia` (los tres últimos, con el esquema exacto de `origin/main`) | `db.go` (`migrarHorasEstimadas`, columna `medias_horas_estimadas`) |
| CA-19.17 Respuestas HTTP | Sección 3 | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y no pasan por HTTP. Se cubre con los tests del handler. | `TestGetCompararEsfuerzo_Devuelve200ConLaComparacion`, `TestGetCompararEsfuerzo_ElPorcentajeEsUnNumero`, `TestGetCompararEsfuerzo_NuncaEscribeCeroNegativo`, `TestGetCompararEsfuerzo_NoEncontrado_Devuelve404`, `TestGetCompararEsfuerzo_IDInvalido_Devuelve400`, `TestGetCompararEsfuerzo_FallaInterna_Devuelve500SinDetalle` | `comparar_esfuerzo_handler.go`, `router.go`, `cmd/server/main.go` |
| CA-19.18 Pantalla | Sección 3 | **Sin escenario BDD a propósito:** es la pantalla. Se cubre con tests del adaptador (`httptest`). | `TestPantallaComparacion_MuestraLasCuatroCifras`, `TestPantallaComparacion_AlVolverConError_MuestraLoGuardado`, `TestPantallaElemento_CampoHorasEstimadas`, `TestPantallaElemento_CargarYQuitarHorasEstimadas`, `TestPantallaElemento_HorasInvalidas_VuelvenAlFormularioCon400`, `TestPantallaElemento_HorasQueNoSonNumero_NoLlamanAlCasoDeUso`, `TestPantallaElemento_HorasTipeadas_SeEscapan`, `TestPantallaElemento_TodosLosCampos_SeEscapan`, `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran`, `TestPantallaComparacion_NoEncontrados_Devuelven404`, `TestPantallaComparacion_FallaInterna_Devuelve500SinDetalle`, `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` (orden de las secciones) | `pantallas_comparacion.go` (`vistaComparacion`), `pantallas_elemento.go` (`mostrarElemento`), `elemento.html` |
| CA-19.19 Campo ausente en el formulario | RN-19.12 | **Sin escenario BDD a propósito:** es la pantalla. Se cubre con tests del adaptador (`httptest`). | `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran`, `TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06` | `pantallas_elemento.go` (`HorasEnviadas`, `cambiosDeFormElemento`) |

Los tests de las pantallas están en [pantallas_comparacion_test.go](../internal/adapters/http/pantallas_comparacion_test.go), salvo `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos`, que está en [pantallas_esfuerzo_test.go](../internal/adapters/http/pantallas_esfuerzo_test.go).

### Reglas de negocio

| RN | CA que la cubren | Tests con la etiqueta (ejemplos) | Código |
|---|---|---|---|
| RN-19.1 Horas estimadas: dato propio y opcional | CA-19.4, CA-19.15 | `TestEditar_HorasEstimadas_ConservanStoryPoints`, `TestPantallaElemento_HorasYStoryPoints_NoSePisan`, `TestMigracionHU19_ElementosSinEstimarYConSusDatos` | `backlog.go` |
| RN-19.2 Mayor a 0, múltiplo de 0,5 y máximo 1000 | CA-19.14 | `TestValidarHorasEstimadas_Invalidas`, `TestEditar_HorasEstimadasInvalidas_RetornaError`, `TestPatchElemento_HorasEstimadasInvalidas_Devuelve400ConElMensajeDelDominio` | `backlog_editar.go` (`validarHorasEstimadas`) |
| RN-19.3 Se cargan, cambian o quitan editando; nacen sin estimar; sin historial | CA-19.13, CA-19.15 | `TestEditar_QuitarHorasEstimadas`, `TestNuevoElementoBacklog_NaceSinHorasEstimadas`, `TestRepositorioBacklog_HorasEstimadas_CambiarYQuitar` | `backlog_editar.go`, `NuevoElementoBacklog` no las recibe |
| RN-19.4 Horas reales: suma de los registros, sin tope, 0 sin registros | CA-19.1, CA-19.5, CA-19.6, CA-19.11 | `TestRepositorioEsfuerzo_HorasReales_LaSumaPuedeSuperar24`, `TestRepositorioEsfuerzo_HorasReales_SinRegistros_DaCero`, `TestMigracionHU19_LaSumaDeHorasRealesNoCambia` | `esfuerzo_repo.go` (`COALESCE(SUM(medias_horas), 0)`) |
| RN-19.5 Desviación absoluta: reales − estimadas | CA-19.1, CA-19.2, CA-19.3 | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_NoCompartePunteroRecibido`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento` | `comparacion_esfuerzo.go` |
| RN-19.6 Desviación porcentual: un decimal, medios alejándose de cero, nunca −0 | CA-19.1, CA-19.2, CA-19.3, CA-19.8 | `TestCompararEsfuerzo_MultiplicaAntesDeDividir`, `TestCompararEsfuerzo_NuncaDevuelveCeroNegativo`, `TestGetCompararEsfuerzo_NuncaEscribeCeroNegativo` | `comparacion_esfuerzo.go` (`redondearUnDecimal`) |
| RN-19.7 Sin estimar no es un error y no hay desviaciones | CA-19.4, CA-19.6 | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`, `TestPantallaComparacion_MuestraLasCuatroCifras` | `comparacion_esfuerzo.go` (`SinEstimar`) |
| RN-19.8 Estimado 0: desviación en horas, sin porcentaje | CA-19.7 | `TestCompararEsfuerzo`, `TestGetCompararEsfuerzo_Devuelve200ConLaComparacion`, `TestPantallaComparacion_MuestraLasCuatroCifras` | `comparacion_esfuerzo.go` |
| RN-19.9 Con estimación y sin esfuerzo: −100 % | CA-19.5 | `TestCompararEsfuerzo`, `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`, `TestPantallaComparacion_MuestraLasCuatroCifras` | `comparacion_esfuerzo.go` |
| RN-19.10 El proyecto existe y el elemento es suyo | CA-19.9, CA-19.10 | `TestCompararEsfuerzo_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestGetCompararEsfuerzo_NoEncontrado_Devuelve404`, `TestPantallaComparacion_NoEncontrados_Devuelven404` | `comparar_esfuerzo.go` |
| RN-19.11 Se compara en cualquier estado | CA-19.12 | `TestCompararEsfuerzo_CualquierEstadoDeElementoYProyecto` | `comparar_esfuerzo.go` (no mira estados) |
| RN-19.12 `horas_estimadas` ausente del formulario no cambia las horas | CA-19.19 | `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran`, `TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06` | `pantallas_elemento.go` (`HorasEnviadas`) |

Sin CA propio en la spec. Algunos llevan en su `Traza` el CA más cercano y la sección:
- Falla interna con 500 genérico (sección 7): `TestCompararEsfuerzo_FallaDeInfraestructura_SePropaga`, `TestGetCompararEsfuerzo_FallaInterna_Devuelve500SinDetalle`, `TestPantallaComparacion_FallaInterna_Devuelve500SinDetalle`. CA-19.17 no menciona el 500.
- Una base cerrada da error y no pánico: `TestRepositorios_BaseCerrada_DevuelvenElError` (suma dos casos de HU-19).
- La ruta está registrada: `TestNuevoRouter_RegistraTodasLasRutas`.
- El JSON esperado del listado de HU-07 suma el campo: `TestGetBacklog_DevuelveLosElementos`.

### Errores nuevos

HU-19 suma un solo error, en [internal/domain/backlog_editar.go](../internal/domain/backlog_editar.go). También reutiliza `ErrProyectoNoEncontrado` (HU-01), `ErrElementoNoEncontrado` (HU-05), `ErrHorasInvalidas` y `ErrHorasNoMultiploMediaHora` (HU-18).

| Error | Mensaje | Lo devuelve | HTTP |
|---|---|---|---|
| `ErrHorasEstimadasExcedenMaximo` | "las horas estimadas no pueden superar las 1000 horas" | `domain.validarHorasEstimadas` | 400 |

Los mensajes nuevos de formato viven en el adaptador: "horas_estimadas debe ser un número o null" (PATCH) y "las horas estimadas deben ser un número" (formulario).

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/comparacion_esfuerzo.go](../internal/domain/comparacion_esfuerzo.go) (`ComparacionEsfuerzo`, `CompararEsfuerzo`, `redondearUnDecimal`), [backlog.go](../internal/domain/backlog.go) (`HorasEstimadas`, `copiarPuntero`), [backlog_editar.go](../internal/domain/backlog_editar.go) (`CambiarHorasEstimadas`, `validarHorasEstimadas`, `Editar`), [registro_esfuerzo.go](../internal/domain/registro_esfuerzo.go) (`validarMediasHoras`, compartida con HU-18) | [comparacion_esfuerzo_test.go](../internal/domain/comparacion_esfuerzo_test.go), [horas_estimadas_test.go](../internal/domain/horas_estimadas_test.go), [backlog_editar_test.go](../internal/domain/backlog_editar_test.go) |
| Caso de uso y puerto | [internal/app/comparar_esfuerzo.go](../internal/app/comparar_esfuerzo.go) (`CompararEsfuerzo`, puerto `LectorHorasReales`; reutiliza `LectorProyectos` de HU-04 y `LectorElementos` de HU-18) | [comparar_esfuerzo_test.go](../internal/app/comparar_esfuerzo_test.go) |
| Persistencia | [internal/adapters/persistencia/esfuerzo_repo.go](../internal/adapters/persistencia/esfuerzo_repo.go) (`HorasRealesDeElemento`), [backlog_repo.go](../internal/adapters/persistencia/backlog_repo.go) (crear, buscar y listar con horas; `aMediasHorasEstimadas`), [backlog_editar.go](../internal/adapters/persistencia/backlog_editar.go) (`ActualizarElemento`), [db.go](../internal/adapters/persistencia/db.go) (`migrarHorasEstimadas`) | [horas_estimadas_repo_test.go](../internal/adapters/persistencia/horas_estimadas_repo_test.go), [esfuerzo_repo_test.go](../internal/adapters/persistencia/esfuerzo_repo_test.go), [migracion_hu19_test.go](../internal/adapters/persistencia/migracion_hu19_test.go), [db_test.go](../internal/adapters/persistencia/db_test.go), [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) |
| Endpoint JSON | [internal/adapters/http/comparar_esfuerzo_handler.go](../internal/adapters/http/comparar_esfuerzo_handler.go) (`RegistrarCompararEsfuerzo`), [backlog_editar_handler.go](../internal/adapters/http/backlog_editar_handler.go) (PATCH), [backlog_handler.go](../internal/adapters/http/backlog_handler.go) y [listar_backlog_handler.go](../internal/adapters/http/listar_backlog_handler.go) (respuestas), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [comparar_esfuerzo_handler_test.go](../internal/adapters/http/comparar_esfuerzo_handler_test.go), [horas_estimadas_handler_test.go](../internal/adapters/http/horas_estimadas_handler_test.go), [listar_backlog_handler_test.go](../internal/adapters/http/listar_backlog_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantalla y plantilla | [internal/adapters/http/pantallas_comparacion.go](../internal/adapters/http/pantallas_comparacion.go) (`vistaComparacion`), [pantallas_elemento.go](../internal/adapters/http/pantallas_elemento.go) (`mostrarElemento`, `leerFormElemento`, `cambiosDeFormElemento`), [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) | [pantallas_comparacion_test.go](../internal/adapters/http/pantallas_comparacion_test.go), [pantallas_esfuerzo_test.go](../internal/adapters/http/pantallas_esfuerzo_test.go), [pantallas_test.go](../internal/adapters/http/pantallas_test.go) |
| BDD | [features/HU-19-comparar-esfuerzo.feature](../features/HU-19-comparar-esfuerzo.feature), [test/bdd/hu19_steps_test.go](../test/bdd/hu19_steps_test.go), [hu19_sistema_test.go](../test/bdd/hu19_sistema_test.go) (SQLite real), [bdd_test.go](../test/bdd/bdd_test.go); cableado de `CompararEsfuerzo` en [hu06_sistema_test.go](../test/bdd/hu06_sistema_test.go) y [hu14_sistema_test.go](../test/bdd/hu14_sistema_test.go) | 34 escenarios en verde |

### Commits

Salen de `git log origin/hu-18-pantalla-esfuerzo..HEAD --no-merges`, en orden. Excluye los commits de HU-18 que trajo el merge.

| Ciclo | Commit | Etapa |
|---|---|---|
| Spec | `714b8fa` docs(HU-19): spec de comparar esfuerzo estimado contra real | — |
| BDD en rojo: feature | `4fae233` test(HU-19): [RED] escenarios BDD del feature | BDD en rojo (34 escenarios y 213 pasos indefinidos) |
| BDD en rojo: steps | `19eb4be` test(bdd): [RED] steps y sistema en memoria de HU-19 | BDD en rojo, por compilación |
| Dominio: cálculo | `3147f82` [RED] → `49ca557` [GREEN] | **RED → GREEN** (sin REFACTOR: no había nada real para mejorar) |
| Dominio: validación | `1a81795` [RED] → `0a68699` [GREEN] → `f4f5dbb` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Dominio: `SprintID` en la edición | `adf99b4` [RED] → `1d701ab` [GREEN] | **RED → GREEN** |
| Dominio: horas en la edición | `3574f61` [RED] → `22f3352` [GREEN] → `ca96df0` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Caso de uso | `147f99f` [RED] → `8c829ae` [GREEN] | **RED → GREEN**; los 34 escenarios pasan con el sistema en memoria |
| Persistencia: horas estimadas | `7fe13b7` [RED] → `2425d80` [GREEN] | **RED → GREEN** |
| Persistencia: horas reales | `5019bea` [RED] → `0b6d3ae` [GREEN] | **RED → GREEN** |
| BDD sobre SQLite | `4d25796` test(bdd): HU-19 sistema SQLite real | BDD en verde |
| Migración | `a271b8e` test(HU-19): migración desde el esquema anterior a HU-19 | — (sin RED real, ver abajo) |
| Comentarios de los steps | `c4b3c8a` test(bdd): actualizar comentarios de los steps de HU-19 | — |
| Endpoint JSON: GET | `42cdbf3` [RED] → `7827bb8` [GREEN] | **RED → GREEN** |
| Endpoint JSON: `horas_estimadas` | `62e9bde` [RED] → `00b5fb9` [GREEN] | **RED → GREEN** (GREEN unificado, ver abajo) |
| Spec: contrato JSON | `2f9d244` docs(HU-19): ajustar la spec al contrato JSON real | — |
| Spec: falla interna | `fe05f12` docs(HU-19): corregir la fila de falla interna de la spec | — |
| Etiquetas | `a3010c1` test(HU-19): etiquetas de trazabilidad literales | — (solo comentarios) |
| Pantalla | `f156b55` [RED] → `0174918` [GREEN] | **RED → GREEN** (GREEN unificado, ver abajo) |
| Spec: pantalla, P-11 y P-13 | `83d39dd` docs(HU-19): cerrar las preguntas P-11 y P-13 y ajustar la pantalla en la spec | — |
| POST sin horas | `374e629` [RED] → `6b2649a` [GREEN] | **RED → GREEN** |
| Escape de HTML | `4292c51` test(HU-19): proteger el escape de HTML en todos los campos de la pantalla del elemento | — (sin RED real, ver abajo) |
| Spec: RN-19.12 y CA-19.19 | `159981b` docs(HU-19): regla del campo horas_estimadas ausente en el formulario | — |
| Etiqueta de CA-19.19 | `9e64052` test(HU-19): etiqueta de trazabilidad de CA-19.19 | — (solo comentarios) |
| Renombre del paso BDD | `a23c87e` (merge): `"X" tiene N story points` → `"X" quedó con N story points` | — (resuelto dentro del merge) |
| Test de escape con el desplegable de HU-13 | `6054837` (merge): `TestPantallaElemento_TodosLosCampos_SeEscapan` exige la `<option>` elegida, con el valor escapado en el atributo y en el texto. Quitar el escape de cualquiera de los dos hace fallar el test. | — (resuelto dentro del merge) |
| Trazabilidad | este commit | — |

**Commits que el filtro `git log --grep="(HU-19)"` no encuentra**, porque usan el prefijo `test(bdd):`: `19eb4be` (RED de los steps), `4d25796` (sistema SQLite) y `c4b3c8a` (comentarios).

**Merges de la rama, que no son de HU-19:** `29991ac` (trae `main`: HU-14 y HU-15), `a106188` (trae `hu-18-pantalla-esfuerzo`), `a23c87e` (trae `main`: HU-16 y HU-17) y `6054837` (trae `main`: HU-13).

**Dos unificaciones de commits GREEN**, hechas con `git reset --soft` antes de subir, porque cada GREEN había dejado tests en rojo. Los hashes de la tabla son los del historial real; los anteriores solo quedan en el reflog local. En las dos, el hash de árbol fue el mismo antes y después.

| Commit final | Unifica | Motivo | Árbol |
|---|---|---|---|
| `00b5fb9` | `bce9940` + `e5a3ebf` | El GREEN del PATCH dejó en rojo `TestGetBacklog_DevuelveLosElementos` de HU-07, que compara el JSON exacto | `06b30303` |
| `0174918` | `edbc4d8` + `c2a979b` | El GREEN de la pantalla dejó en rojo 4 escenarios de HU-06 y HU-14, cuyos sistemas BDD no cableaban `CompararEsfuerzo` | `a9e08442` |

### Evidencia: mutaciones manuales

Ninguna quedó en un commit. Se aplicaron con `go test -overlay` o en una copia descartable del repo.

**Cálculo**

| Mutación | Resultado |
|---|---|
| Sin la protección del `-0` | Falla `TestCompararEsfuerzo_NuncaDevuelveCeroNegativo`: `DesviacionPorcentaje = -0` |
| Dividir antes de multiplicar | Falla `TestCompararEsfuerzo_MultiplicaAntesDeDividir`: 51,5 sobre 40 da 28,7 en vez de 28,8 |
| Compartir el puntero recibido | Falla `TestCompararEsfuerzo_NoCompartePunteroRecibido` |

Sobre todos los empates exactos hasta 1000 h, dividir primero redondea mal 451 de 11.088, y multiplicar primero no falla ninguno.

**BDD sobre SQLite**

| Persistencia rota | Escenarios que fallan (de 34) |
|---|---|
| `ActualizarElemento` no guarda las horas | 24 |
| No se leen las horas del elemento | 24 |
| La suma incluye registros de otros elementos | 1 |

**Migración:** sin `migrarHorasEstimadas`, fallan los tres tests de [migracion_hu19_test.go](../internal/adapters/persistencia/migracion_hu19_test.go) con `no such column: medias_horas_estimadas`.

**Pantalla (nueve mutaciones)**

| # | Mutación | Fallan |
|---|---|---|
| 1 | La plantilla no muestra el bloque | `TestPantallaComparacion_MuestraLasCuatroCifras`, `TestPantallaComparacion_AlVolverConError_MuestraLoGuardado`, `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` |
| 2 | Falta el oculto `original_horas_estimadas` | `TestPantallaElemento_CampoHorasEstimadas`, `TestPantallaElemento_HorasInvalidas_VuelvenAlFormularioCon400` |
| 3 | Sin signo | `TestPantallaComparacion_MuestraLasCuatroCifras`, `TestPantallaComparacion_AlVolverConError_MuestraLoGuardado` |
| 4 | Signo también en el cero | `TestPantallaComparacion_MuestraLasCuatroCifras` |
| 5 | Punto en vez de coma | `TestPantallaComparacion_MuestraLasCuatroCifras`, `TestPantallaComparacion_AlVolverConError_MuestraLoGuardado` |
| 6 | Porcentaje sin decimal fijo | `TestPantallaComparacion_MuestraLasCuatroCifras`, `TestPantallaComparacion_AlVolverConError_MuestraLoGuardado` |
| 7 | El formulario siempre envía el campo | `TestPantallaElemento_HorasSinCambios_NoSeEnvian`, `TestPantallaElemento_HorasYStoryPoints_NoSePisan`, `TestPantallaElemento_SinCambios_VuelveSinError` |
| 8 | Vacío no quita la estimación | `TestPantallaElemento_CargarYQuitarHorasEstimadas` |
| 9 | Se ignoran los errores de la comparación | `TestPantallaComparacion_NoEncontrados_Devuelven404`, `TestPantallaComparacion_FallaInterna_Devuelve500SinDetalle` |

**Orden de las secciones:** mover el bloque antes de "Registrar esfuerzo" hace fallar `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos`.

**POST sin horas**

| Mutación | Fallan |
|---|---|
| Se ignora `HorasEnviadas` (vuelve el bug) | `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran` |
| `HorasEnviadas` siempre falso | `TestPantallaElemento_CargarYQuitarHorasEstimadas` y otros |
| La corrección también se aplica a `story_points` | `TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06` y otros |

**Escape de HTML** (`TestPantallaElemento_TodosLosCampos_SeEscapan`, con `template.HTML` en cada lugar)

| Mutación | Resultado |
|---|---|
| Cada uno de los 7 campos ocultos `original_*` | Detectada, las 7 |
| `input` de título y de horas estimadas (tipeados) | Detectada |
| `<h1>` | Detectada |
| Plantilla sin el oculto de horas | Detectada |
| `<textarea>` de criterios | **Sobrevive: mutante equivalente** |
| `<title>` | **Sobrevive: mutante equivalente** |

Los dos que sobreviven no son un hueco del test: `<textarea>` y `<title>` son contextos RCDATA, y un programa aparte comprobó que `html/template` produce exactamente la misma salida con `string` que con `template.HTML`, ya escapada. En esos dos contextos no hay diferencia que detectar.

**Bug que salió a la luz en el dominio:** `ElementoBacklog.Editar` reconstruye el elemento con `NuevoElementoBacklog` y perdía `SprintID`. En producción no se notaba, porque el repositorio guarda solo los campos enviados y el caso de uso relee el elemento. Se corrigió con su propio ciclo: `adf99b4` [RED], que falla con `SprintID = <nil>, se esperaba 3`, y `1d701ab` [GREEN]. Toca código de HU-06.

### Etapas sin RED real

- **Migración (`a271b8e`):** los tres tests de [migracion_hu19_test.go](../internal/adapters/persistencia/migracion_hu19_test.go) pasaron a la primera, porque la migración ya existía. Su valor es la mutación de arriba.
- **Escape de todos los campos (`4292c51`):** `TestPantallaElemento_TodosLosCampos_SeEscapan` pasó a la primera, porque el código ya escapaba bien. Su valor son las mutaciones.
- **RED por compilación:** `3147f82`, `1a81795`, `3574f61`, `147f99f`, `5019bea` y `42cdbf3`. Para ver fallar las aserciones se usó después un stub temporal, sin commitear. Con ese stub pasaban:
  - `1a81795`: `TestValidarHorasEstimadas_Validas`, `TestErrHorasEstimadasExcedenMaximo_Mensaje`, `TestValidarHoras_CadaUnaConSuTope`
  - `3574f61`: `TestNuevoElementoBacklog_NaceSinHorasEstimadas`, `TestEditar_HorasEstimadasInvalidas_SeInformanAlFinal`
  - `147f99f`: el subcaso "sin estimar y sin esfuerzo" de `TestCompararEsfuerzo_DevuelveLaComparacionDelElemento`
  - `5019bea`: `TestRepositorioEsfuerzo_HorasReales_SinRegistros_DaCero`, `TestRepositorioEsfuerzo_HorasReales_ElementoInexistente_DaCero`
- **RED por comportamiento:** en estos commits, los siguientes tests nuevos ya pasaban, medido volviendo a correr cada commit:
  - `7fe13b7`: `TestRepositorioBacklog_HorasEstimadas_ElementoInexistente` y el subcaso "recién creado, sin estimar" de `TestRepositorioBacklog_GuardarElemento_GuardaLasHorasEstimadas`
  - `62e9bde`: `TestPatchElemento_HorasEstimadasYStoryPoints_SonIndependientes`, `TestPostBacklog_ConHorasEstimadas_Devuelve400`
  - `f156b55`: `TestPantallaElemento_HorasSinCambios_NoSeEnvian`
  - `374e629`: `TestPantallaElemento_StoryPointsAusentes_MantienenElComportamientoDeHU06` y 7 de los 9 subcasos de `TestPantallaElemento_HorasEstimadasAusentes_NoSeBorran`; fijan lo que no debía cambiar

### Huecos

- La spec sigue en **Borrador**.
- **Decisiones confirmadas por el equipo (no abiertas):** P-10 (campo propio `HorasEstimadas`, toca HU-05 y HU-06) y P-12 (tope de 1000 h), confirmadas por Matías el 2026-11-03 en la revisión del PR #297.
- **Carrera de `agregarColumnas` en las migraciones:** #296 (DEF-04). Es preexistente (también `sprint_id` y `puntos_*`, por lectura del código) y compartida. Con 4 `AbrirDB` simultáneos sobre el esquema anterior a HU-19, hubo fallas en 30 de 30 rondas (89 de 120 aperturas, `duplicate column name`), con conexiones en un mismo proceso. Según ADR-0004, en el plan gratuito de Render el disco se borra en cada despliegue, así que hoy esas bases nacen con la columna.
- El mensaje "las horas se registran en múltiplos de media hora" al estimar viene de HU-18 y la sección 7 de la spec lo fija. Se entiende, pero habla de registrar.
- **Asimetría conocida, para decisión del equipo:** un POST del formulario sin campos ocultos conserva las horas estimadas (RN-19.12), pero sigue quitando los story points (comportamiento de HU-06, DEF-01). Comprobado contra el servidor real.
- Si el formulario vuelve con otro error en un pedido sin el campo de horas, las muestra vacías aunque haya guardadas. Reenviarlo no las borra, y desde el navegador el campo siempre viene.
- Si la comparación falla al volver de un error de edición, se responde 500 y se pierde lo tipeado. Es intencional y está probado, pero la spec no lo menciona.
- La pantalla **no se probó en un navegador** ni con un lector de pantalla: se probó con tests y con curl contra el servidor.
- La duración del BDD completo subió a unos 224 s: cada escenario de cualquier historia crea una base más, porque los hooks de todas las historias corren siempre.
- **HU-13 integrada:** HU-13 (que toca `story_points` en los mismos archivos) entró a `main` y se integró en `6054837` (merge de `main` en la rama). Se resolvieron los conflictos de [pantallas_elemento.go](../internal/adapters/http/pantallas_elemento.go) (se conservan `Comparacion` y `OpcionesStoryPoints`) y de este archivo (las secciones de HU-19 y HU-13, que Git había intercalado). `TestPantallaElemento_TodosLosCampos_SeEscapan` se ajustó al desplegable de story points; el cambio quedó dentro del merge. Tras integrar: 355/355 escenarios en el BDD completo (pasos: 1755/1755), y los de HU-19 siguen siendo 34/34.
- **Paso BDD renombrado al integrar con `main`:** `"X" tiene N story points` de HU-19 colisionaba con el mismo paso de HU-17; godog lo rechaza como ambiguo y el compilador no lo detecta. Se renombró en HU-19 a `"X" quedó con N story points` ([features/HU-19-comparar-esfuerzo.feature](../features/HU-19-comparar-esfuerzo.feature) y [test/bdd/hu19_steps_test.go](../test/bdd/hu19_steps_test.go)). La comprobación es la misma. Se tocó HU-19 porque HU-17 ya estaba en `main`. Tras integrar: 336/336 escenarios en el BDD completo (pasos: 1687/1687), y los de HU-19 siguen siendo 34/34. Ese es el total de ese merge (`a23c87e`, que trajo HU-16 y HU-17), no el actual.
- Etiquetas: los tests atados solo a una sección de la spec llevan el CA más cercano. El 500 usa CA-19.17, aunque ese CA no menciona el 500.
- El comentario de cabecera del feature nombra CA-19.16 a CA-19.18, pero no CA-19.19, que se agregó después. `cambiosDeFormElemento` explica la regla, pero no nombra RN-19.12.
- `mostrarElemento` consulta el proyecto y `CompararEsfuerzo` lo vuelve a consultar, igual que el elemento: son lecturas de más, sin consecuencias prácticas.
- Refactors que tocaron código de otras historias: `validarMediasHoras` en `validarHoras` (HU-18), `copiarPuntero` en `NuevoElementoBacklog` (HU-05), `Editar` y `ActualizarElemento` (HU-06), y `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` (HU-18), que suma la sección nueva. La descripción de ese test en HU-18 todavía enumera las secciones anteriores. Conviene que lo revisen sus autores en la PR.

---

## HU-13: Estimar una historia con story points de la escala Fibonacci

- **Issue:** [#24](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/24) · Sprint 3 · 1 SP · tareas #272 a #280
- **Spec:** [docs/specs/HU-13-escala-story-points.md](specs/HU-13-escala-story-points.md)
- **Feature:** [features/HU-13-escala-story-points.feature](../features/HU-13-escala-story-points.feature) (steps en [test/bdd/hu13_steps_test.go](../test/bdd/hu13_steps_test.go), sobre SQLite real; el escenario de CA-13.6 usa el endpoint real)
- **Cambio en HU-05:** RN-05.6, CA-05.7 y el mensaje de `ErrStoryPointsInvalidos` pasan a la escala.

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-13.1 Valores de la escala | RN-13.1 | Los valores de la escala se aceptan `@normal` (×7) | `TestNuevoElementoBacklog_ValoresDeLaEscala` | `domain.NuevoElementoBacklog` |
| CA-13.2 Fuera de la escala al crear | RN-13.1 | Un valor fuera de la escala se rechaza `@error` (×4) | `TestNuevoElementoBacklog_FueraDeLaEscala`, `TestCrearElemento_StoryPointsFueraDeLaEscala_NoSeGuarda` | `domain.NuevoElementoBacklog` |
| CA-13.3 Editar los story points | RN-13.1 | Editar los story points a otro valor de la escala `@normal` · a un valor fuera de la escala `@error` | `TestEditar_StoryPointsFueraDeLaEscala` | `domain.ElementoBacklog.Editar` |
| CA-13.4 Sin estimar | RN-13.2 | Un elemento puede quedar sin estimar `@alternativo` · Quitarle la estimación `@alternativo` | `TestEditar_StoryPointsNil_QuedaSinEstimar` (HU-06), `TestEditar_ValorAnteriorALaEscala` | `domain.NuevoElementoBacklog`, `domain.ElementoBacklog.Editar` |
| CA-13.5 Valores anteriores a la escala | RN-13.3 | Un valor guardado antes de la escala se respeta al editar otro campo `@limite` · no se puede volver a enviar `@error` · se puede corregir `@limite` | `TestEditar_ValorAnteriorALaEscala`, `TestEditarElemento_ValorAnteriorALaEscala`, `TestRepositorioBacklog_ValorAnteriorALaEscala` | `domain.ElementoBacklog.Editar` |
| CA-13.6 API | Sección 3 | La API rechaza un valor fuera de la escala `@error` | `TestPostBacklog_DatosInvalidos_Devuelve400` y `TestPatchElemento_Errores` (casos HU-13) | sin cambios en los handlers |
| CA-13.7 Pantallas | Sección 3, RN-13.3, RN-13.4 | No va acá (como CA-20.8): se cubre con los tests de la pantalla | `TestPantallaBacklog_AltaConLaEscala`, `TestPantallaElemento_EdicionConLaEscala`, `TestPantallaElemento_ValorAnteriorALaEscala`, `TestOpcionesStoryPoints` | `pantallas_escala.go`, `web/plantillas/backlog.html`, `web/plantillas/elemento.html` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/backlog.go](../internal/domain/backlog.go) (`EscalaStoryPoints`, validación en `NuevoElementoBacklog`), [backlog_editar.go](../internal/domain/backlog_editar.go) (`Editar` valida los story points solo si se envían) | [backlog_escala_test.go](../internal/domain/backlog_escala_test.go) |
| Casos de uso | Sin cambios | [escala_story_points_test.go](../internal/app/escala_story_points_test.go) |
| Persistencia | Sin cambios ni migración | [escala_story_points_test.go](../internal/adapters/persistencia/escala_story_points_test.go) |
| HTTP (JSON) | Sin cambios: `ErrStoryPointsInvalidos` ya se respondía con 400 | [backlog_handler_test.go](../internal/adapters/http/backlog_handler_test.go), [backlog_editar_handler_test.go](../internal/adapters/http/backlog_editar_handler_test.go) |
| Pantallas | [internal/adapters/http/pantallas_escala.go](../internal/adapters/http/pantallas_escala.go) (`opcionesStoryPoints`), listas en [backlog.html](../internal/adapters/http/web/plantillas/backlog.html) y [elemento.html](../internal/adapters/http/web/plantillas/elemento.html) | [pantallas_escala_test.go](../internal/adapters/http/pantallas_escala_test.go) |
| BDD | [test/bdd/hu13_steps_test.go](../test/bdd/hu13_steps_test.go) | 19 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 13.1 Spec (y ajuste de HU-05) | `5fa677d` | — |
| 13.2 Escenarios BDD | `711ca5b` | — |
| 13.3 Steps godog | `d7dbeaa` | BDD en rojo |
| 13.4 Escala en el dominio | `a785cf8` [RED] → `f246c79` [GREEN] | **RED → GREEN**, BDD en verde |
| 13.5 Casos de uso | `c679dcf` | test, sin cambios en el código |
| 13.6 Persistencia | `43024a2` | test, sin persistencia nueva |
| 13.7 HTTP | `816bdd7` | test, sin cambios en los handlers |
| 13.8 Pantallas | `5a6432b` [RED] → `ddb9e08` [GREEN] | **RED → GREEN** |
| 13.9 Revisión y trazabilidad | `e5fcbd3` (test de la revisión), `25a3997` | — |
| Revisión del PR #288: escala unificada con HU-14 | `55ba3cf` | **REFACTOR** |

### Huecos

- La spec está en **Borrador** hasta la aprobación del PR.
- Las tareas 13.5 a 13.7 no tienen RED → GREEN: crear y editar ya pasaban por el dominio, la base guarda un entero y los handlers ya respondían `ErrStoryPointsInvalidos` con 400. Se dejaron tests como evidencia, como sugerían las tareas.
- `TestRepositorioSprints_CerrarSprint_SimultaneoConCompletados_QuedaCoherente` (HU-11) creaba historias con 1 a 8 story points; ahora usa valores de la escala.
- **Escala unificada con HU-14 (revisión del PR #288):** después de traer `main`, los votos de HU-14 usan `EscalaStoryPoints()` y se quitó `ValoresVoto`: hay una sola escala en el dominio (RN-13.4). Commit `55ba3cf` [REFACTOR], con `TestVotar_UsaLaEscalaDeLosStoryPoints`.
- Revisión con el agente `go-reviewer` (tarea 13.9): sin hallazgos críticos ni altos; confirmó que `Editar` no comparte el puntero de los story points y que un valor anterior a la escala no se valida si no se envía, pero sí si se envía. Se agregó el test de punta a punta que sugirió (`TestPantallaElemento_GuardarConValorAnteriorALaEscala_NoEnviaLosStoryPoints`, `e5fcbd3`). Dos observaciones bajas sin cambios: en el alta, un POST armado a mano con un valor fuera de la escala se vuelve a mostrar marcado (el dominio igual lo rechaza), y si una edición de 4 a 5 falla por otro error, la lista ya no ofrece el 4.

---

## HU-23: Horas estimadas, reales y desviación de un sprint o proyecto

- **Issue:** #34 · Sprint 3 · 3 SP según el planning y el Product Backlog; la issue dice 5 (P-9)
- **Spec:** [docs/specs/HU-23-horas-desviacion.md](specs/HU-23-horas-desviacion.md)
- **Feature:** [features/HU-23-horas-desviacion.feature](../features/HU-23-horas-desviacion.feature) (steps en [test/bdd/hu23_steps_test.go](../test/bdd/hu23_steps_test.go), sobre SQLite real en [hu23_sistema_test.go](../test/bdd/hu23_sistema_test.go))
- **Rutas HTTP:** `GET /proyectos/{id}/horas` y `GET /proyectos/{id}/sprints/{sprintID}/horas`.
- **Pantallas:** bloque "Horas y desviación" en `GET /web/proyectos/{id}/sprints/{sprintID}` (HU-09) y `GET /web/proyectos/{id}` (HU-30).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-23.1 Horas de un sprint | RN-23.1, RN-23.2, RN-23.3, RN-23.4, RN-23.7 | La desviación de un sprint se calcula sobre la suma de sus elementos `@normal` (8/10 y 4/3) | `TestTotalizarHoras`, `TestTotalizarHoras_ReutilizaCompararEsfuerzo`, `TestConsultarHorasSprint_DevuelveElTotalDeLosElementosDelSprint`, `TestRepositorioEsfuerzo_HorasDeProyecto_SumaIntegrantesYFechas_UnaFilaPorElemento`, `TestRepositorioEsfuerzo_HorasDeSprint_SoloLosElementosDelSprint`, `TestGetHorasSprint_Devuelve200ConElTotal`, `TestGetHorasProyecto_Devuelve200ConElTotal` | `total_horas.go` (`TotalizarHoras`), `consultar_horas.go` (`ConsultarHorasSprint`), `horas_repo.go` (`HorasDeSprint`) |
| CA-23.2 Horas de un proyecto | RN-23.1, RN-23.4, RN-23.10 | Las horas de un proyecto incluyen los elementos con sprint y sin sprint `@normal` | `TestTotalizarHoras`, `TestConsultarHorasProyecto_DevuelveElTotalDeLosElementosDelProyecto`, `TestRepositorioEsfuerzo_HorasDeProyecto_IncluyeConYSinSprint_NoMezclaOtroProyecto`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go`, `consultar_horas.go` (`ConsultarHorasProyecto`), `horas_repo.go` (`HorasDeProyecto`) |
| CA-23.3 Sobre las sumas | RN-23.4 | Outline de la suma `@normal` (10/20 y 90/90) | `TestTotalizarHoras`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go` (una sola llamada a `CompararEsfuerzo`) |
| CA-23.4 Redondeo | RN-23.4 | Outline de la suma `@limite` (33,3, 1,3 y −1,3) | `TestTotalizarHoras`, `TestTotalizarHoras_NuncaDevuelveCeroNegativo`, `TestRepositorioEsfuerzo_HorasDeProyecto_ConservaLasMediasHoras`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal`, `TestGetHoras_NuncaEscribeCeroNegativo` | `comparacion_esfuerzo.go` (HU-19, sin cambios), `horas_repo.go` (`desdeMediasHoras`) |
| CA-23.5 Desvíos que se compensan | RN-23.4 | Outline de la suma `@limite` (8/10 y 8/6) | `TestTotalizarHoras`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go` |
| CA-23.6 Elementos sin estimar | RN-23.2, RN-23.3, RN-23.5 | Un elemento sin estimar suma a las reales y no a la desviación `@alternativo` · Las horas de un elemento sin estimar no compensan la desviación de otro estimado `@limite` | `TestTotalizarHoras`, `TestConsultarHoras_ConElementosSinEstimar_LosPasaAlDominio`, `TestRepositorioEsfuerzo_HorasDeProyecto_EstimadoYSinEstimar`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go` (`HorasRealesSinEstimar`), `horas_repo.go` (`desdeMediasHorasEstimadas`) |
| CA-23.7 Ninguno estimado | RN-23.6 | Un sprint sin ningún elemento estimado informa sus reales y ninguna desviación `@alternativo` | `TestTotalizarHoras`, `TestConsultarHoras_SinElementos_NoEsUnError`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go` (`SinEstimar`, se pasa `nil`) |
| CA-23.8 Sin elementos | RN-23.6 | Un sprint sin elementos informa todo en 0 `@limite` · Un proyecto sin elementos informa todo en 0 `@limite` | `TestTotalizarHoras`, `TestConsultarHoras_SinElementos_NoEsUnError`, `TestRepositorioEsfuerzo_HorasDeProyecto_SinElementos_ListaVacia`, `TestRepositorioEsfuerzo_HorasDeSprint_SinElementos_ListaVacia`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `total_horas.go`, `horas_repo.go` |
| CA-23.9 Sin esfuerzo | RN-23.2, RN-23.3 | Un elemento estimado y sin esfuerzo suma a las estimadas `@limite` | `TestTotalizarHoras`, `TestRepositorioEsfuerzo_HorasDeProyecto_ElementoSinRegistros_ApareceConCero`, `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal` | `horas_repo.go` (`LEFT JOIN`, `COALESCE`) |
| CA-23.10 Cierre del sprint | RN-23.7, RN-23.8, RN-23.10 | Al cerrar un sprint solo cuentan en él los elementos que siguen asignados `@alternativo` · El proyecto sigue contando las horas de un elemento devuelto al backlog `@alternativo` · Las horas de un elemento devuelto al backlog pasan al sprint donde se replanifica `@alternativo` · Un elemento Hecho reasignado deja de contar en el sprint cerrado `@limite` · Un elemento Hecho reasignado cuenta en el sprint nuevo `@limite` · El elemento devuelto y el Hecho reasignado suman juntos en el sprint nuevo `@limite` | `TestRepositorioEsfuerzo_HorasDeSprint_TrasElCierre_SoloLoQueSigueAsignado` | `horas_repo.go` (`e.sprint_id = ?`) |
| CA-23.11 Valores actuales | RN-23.9 | Los totales de un sprint cerrado reflejan los valores actuales `@alternativo` | ✗ Ningún test Go lo cita: solo BDD | `horas_repo.go` (se calcula al consultar) |
| CA-23.12 Qué suma | RN-23.3, RN-23.7, RN-23.9 | Suman los registros de los elementos del sprint y no los de otros elementos `@alternativo` | `TestRepositorioEsfuerzo_HorasDeProyecto_NoMezclaLosRegistrosDeOtroElemento`, `TestRepositorioEsfuerzo_HorasDeProyecto_IncluyeConYSinSprint_NoMezclaOtroProyecto`, `TestRepositorioEsfuerzo_HorasDeSprint_SoloLosElementosDelSprint`, `TestRepositorioEsfuerzo_HorasDeSprint_NoFiltraPorFechaNiPorIntegrante` | `horas_repo.go` (`reg.elemento_id = e.id`, sin filtro de fecha) |
| CA-23.13 Cualquier estado | RN-23.12 | Se consultan las horas de un sprint Abierto en un proyecto que no está Activo `@alternativo` (4 ejemplos) | `TestConsultarHoras_CualquierEstadoDeSprintYProyecto` | `consultar_horas.go` (no mira estados) |
| CA-23.14 Proyecto inexistente | RN-23.11 | No se consultan las horas de un proyecto que no existe `@error` · No se consultan las horas de un sprint desde un proyecto que no existe `@error` | `TestConsultarHorasSprint_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestConsultarHorasProyecto_ProyectoInexistente_RetornaErrorYNoConsultaHoras`, `TestConsultarHorasSprint_ProyectoInexistente_NoBuscaElSprint`, `TestGetHoras_NoEncontrado_Devuelve404`, `TestPantallasHoras_NoEncontrados_Devuelven404`, `TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras` | `consultar_horas.go`, `consultar_horas_handler.go` (`responderErrorHoras`), `pantallas_horas.go` (`errorDeHoras`) |
| CA-23.15 Sprint inexistente o de otro proyecto | RN-23.11 | No se consultan las horas de un sprint que no existe `@error` · Un sprint de otro proyecto se trata como inexistente `@error` | `TestConsultarHorasSprint_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestGetHoras_NoEncontrado_Devuelve404`, `TestPantallasHoras_NoEncontrados_Devuelven404`, `TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras`, `TestPantallaHorasSprint_SprintDeOtroProyecto_Devuelve404` | `consultar_horas.go` (`s.ProyectoID != proyectoID`), `pantallas_asignacion.go` (`buscarSprint`) |
| CA-23.16 Respuestas HTTP | Secciones 3 y 7 | **Sin escenario BDD a propósito:** los escenarios prueban el núcleo y no pasan por HTTP. Se cubre con los tests del handler. | `TestGetHorasProyecto_Devuelve200ConElTotal`, `TestGetHorasSprint_Devuelve200ConElTotal`, `TestGetHoras_NuncaEscribeCeroNegativo`, `TestGetHoras_RespondeJSON`, `TestGetHoras_NoEncontrado_Devuelve404`, `TestGetHoras_IDInvalido_Devuelve400`, `TestGetHoras_FallaInterna_Devuelve500SinDetalle` | `consultar_horas_handler.go`, `router.go`, `cmd/server/main.go` |
| CA-23.17 Pantalla del sprint | Sección 3 | **Sin escenario BDD a propósito:** es la pantalla. Se cubre con tests del adaptador (`httptest`). | `TestPantallasHoras_MuestranLosValores`, `TestPantallasHoras_ElBloqueVaAntesDeLaPrimeraSeccion`, `TestPantallasHoras_NoRompenLoQueYaMostraban`, `TestPantallasHoras_ElBloqueEsDeSoloLectura`, `TestPantallasHoras_NoApareceEnOtrasPantallas`, `TestPantallasHoras_AclaracionDelSprintCerrado`, `TestPantallasHoras_ConsultanConLosIDsDeLaRuta`, `TestPantallasHoras_AlVolverConError_SigueElBloque`, `TestPantallasHoras_NoEncontrados_Devuelven404`, `TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras`, `TestPantallaHorasSprint_SprintDeOtroProyecto_Devuelve404`, `TestPantallasHoras_FallaInterna_Devuelve500SinDetalle`, `TestPantallasHoras_LosValoresSeEscapan` | `pantallas_horas.go` (`vistaHoras`, `avisoSinEstimar`), `pantallas_asignacion.go` (`mostrarSprintDetalle`), `layout.html` (`bloque_horas`), `sprint_detalle.html` |
| CA-23.18 Pantalla del proyecto | Sección 3 | **Sin escenario BDD a propósito:** es la pantalla. Se cubre con tests del adaptador (`httptest`). | `TestPantallasHoras_MuestranLosValores`, `TestPantallasHoras_ElBloqueVaAntesDeLaPrimeraSeccion`, `TestPantallasHoras_NoRompenLoQueYaMostraban`, `TestPantallasHoras_ElBloqueEsDeSoloLectura`, `TestPantallasHoras_NoApareceEnOtrasPantallas`, `TestPantallasHoras_AclaracionDelSprintCerrado`, `TestPantallasHoras_ConsultanConLosIDsDeLaRuta`, `TestPantallasHoras_AlVolverConError_SigueElBloque`, `TestPantallasHoras_NoEncontrados_Devuelven404`, `TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras`, `TestPantallasHoras_FallaInterna_Devuelve500SinDetalle`, `TestPantallasHoras_LosValoresSeEscapan` | `pantallas_horas.go`, `pantallas.go` (`mostrarDetalle`), `layout.html` (`bloque_horas`), `proyecto_detalle.html` |

Los tests se listan según su etiqueta `Traza`. Los escenarios de CA-23.1, CA-23.3, CA-23.4 y CA-23.5 son ejemplos de un mismo `Scenario Outline`, con el tag en cada bloque `Examples`.

### Reglas de negocio

| RN | CA que la cubren | Tests con la etiqueta (ejemplos) | Código |
|---|---|---|---|
| RN-23.1 Se consulta por sprint o por proyecto | CA-23.1, CA-23.2 | `TestConsultarHorasSprint_DevuelveElTotalDeLosElementosDelSprint`, `TestConsultarHorasProyecto_DevuelveElTotalDeLosElementosDelProyecto` | `total_horas.go` (`TotalHoras`), `consultar_horas.go` |
| RN-23.2 Estimadas: suma de los estimados; los story points no cuentan | CA-23.1, CA-23.6, CA-23.9 | `TestTotalizarHoras`, `TestRepositorioEsfuerzo_HorasDeProyecto_EstimadoYSinEstimar`, `TestRepositorioEsfuerzo_Horas_ElementoConStoryPointsYSinHoras_SigueSinEstimar` | `total_horas.go` |
| RN-23.3 Reales: suma de todos los elementos, de cualquier integrante y fecha | CA-23.1, CA-23.6, CA-23.9, CA-23.12 | `TestTotalizarHoras`, `TestRepositorioEsfuerzo_HorasDeProyecto_NoMezclaLosRegistrosDeOtroElemento`, `TestRepositorioEsfuerzo_HorasDeSprint_NoFiltraPorFechaNiPorIntegrante` | `total_horas.go`, `horas_repo.go` |
| RN-23.4 Una sola desviación, sobre las sumas | CA-23.1, CA-23.2, CA-23.3, CA-23.4, CA-23.5 | `TestTotalizarHoras`, `TestTotalizarHoras_ReutilizaCompararEsfuerzo`, `TestPantallasHoras_MuestranLosValores` | `total_horas.go` (`CompararEsfuerzo` de HU-19) |
| RN-23.5 La desviación compara solo a los elementos estimados | CA-23.6 | `TestTotalizarHoras`, `TestTotalizarHoras_ReutilizaCompararEsfuerzo`, `TestConsultarHoras_ConElementosSinEstimar_LosPasaAlDominio`, `TestPantallasHoras_MuestranLosValores` | `total_horas.go`, `pantallas_horas.go` (`avisoSinEstimar`) |
| RN-23.6 Sin estimados no es un error: 0 estimadas y sin desviación (`nil`, no 0) | CA-23.7, CA-23.8 | `TestTotalizarHoras`, `TestConsultarHoras_SinElementos_NoEsUnError`, `TestRepositorioEsfuerzo_HorasDeProyecto_SinElementos_ListaVacia`, `TestPantallasHoras_MuestranLosValores` | `total_horas.go` (`SinEstimar`) |
| RN-23.7 Elementos del sprint: los que hoy tienen ese `SprintID` | CA-23.1, CA-23.10, CA-23.12 | `TestConsultarHorasSprint_DevuelveElTotalDeLosElementosDelSprint`, `TestRepositorioEsfuerzo_HorasDeSprint_SoloLosElementosDelSprint`, `TestRepositorioEsfuerzo_HorasDeSprint_TrasElCierre_SoloLoQueSigueAsignado`, `TestRepositorioEsfuerzo_HorasDeSprint_ElementoQuitado_DejaDeContar` | `horas_repo.go` |
| RN-23.8 Un elemento devuelto al backlog deja de contar en el sprint cerrado (limitación, P-3) | CA-23.10 | `TestRepositorioEsfuerzo_HorasDeSprint_TrasElCierre_SoloLoQueSigueAsignado`, `TestPantallasHoras_AclaracionDelSprintCerrado` | `horas_repo.go`, `pantallas_horas.go` (`textoAclaracionSprintCerrado`) |
| RN-23.9 Valores actuales: sin congelar al cerrar y sin mirar la fecha | CA-23.11, CA-23.12 | `TestRepositorioEsfuerzo_HorasDeSprint_NoFiltraPorFechaNiPorIntegrante` | `horas_repo.go` |
| RN-23.10 Proyecto: todos sus elementos, con sprint o sin él | CA-23.2, CA-23.10 | `TestConsultarHorasProyecto_DevuelveElTotalDeLosElementosDelProyecto`, `TestRepositorioEsfuerzo_HorasDeProyecto_IncluyeConYSinSprint_NoMezclaOtroProyecto`, `TestRepositorioEsfuerzo_HorasDeSprint_TrasElCierre_SoloLoQueSigueAsignado` | `horas_repo.go`, `consultar_horas.go` |
| RN-23.11 El proyecto existe y el sprint es suyo | CA-23.14, CA-23.15 | `TestConsultarHorasSprint_ValidacionFallida_RetornaErrorYNoConsultaHoras`, `TestConsultarHorasProyecto_ProyectoInexistente_RetornaErrorYNoConsultaHoras`, `TestGetHoras_NoEncontrado_Devuelve404`, `TestPantallaHorasSprint_SprintDeOtroProyecto_Devuelve404` | `consultar_horas.go` |
| RN-23.12 Se consulta en cualquier estado | CA-23.13 | `TestConsultarHoras_CualquierEstadoDeSprintYProyecto` | `consultar_horas.go` (no mira estados) |

### Tests sin CA propio

Llevan `Traza: HU-23 (sección N: ...)`, porque verifican lo que la spec dice en una sección y ningún CA lo menciona. `grep "CA-23\."` no los encuentra; `grep "HU-23"`, sí. Se listan a mano:

| Test | Traza | Qué verifica |
|---|---|---|
| `TestTotalizarHoras_NoModificaLosElementos` | Sección 5: cálculo en el dominio | No modifica lo que recibe |
| `TestConsultarHoras_ElLectorRecibeElContextoDelPedido` | Sección 5: contexto | El lector recibe el contexto del pedido |
| `TestGetHoras_ElCasoDeUsoRecibeElContextoDelPedido` | Sección 5: contexto | El caso de uso recibe el contexto del pedido |
| `TestConsultarHoras_FallaDeInfraestructura_SePropaga` | Sección 7: falla interna al consultar | Las búsquedas devuelven el error tal cual; el lector de horas lo envuelve con `%w` |
| `TestRepositorioEsfuerzo_HorasDeProyecto_OrdenaPorElemento` | Sección 5: puerto de lectura | Filas ordenadas por el id del elemento |
| `TestRepositorioEsfuerzo_Horas_ElementoConStoryPointsYSinHoras_SigueSinEstimar` | Sección 6 (RN-23.2) | Los story points no se convierten en horas |
| `TestRepositorioEsfuerzo_HorasDeSprint_ElementoQuitado_DejaDeContar` | Sección 6 (RN-23.7, RN-23.10) | Quitar un elemento de un sprint abierto, con el caso de uso de HU-09 |
| `TestRepositorioEsfuerzo_HorasDeProyecto_ElementoEliminado_DejaDeSumar` | Sección 6 (RN-23.2, RN-23.10) | Eliminar un elemento sin esfuerzo, con el caso de uso de HU-06 |

Sin etiqueta `Traza` de HU-23, porque son tests compartidos que suman casos:
- Una base cerrada da error y no pánico: `TestRepositorios_BaseCerrada_DevuelvenElError` (suma dos casos, marcados `// HU-23`).
- Las rutas están registradas: `TestNuevoRouter_RegistraTodasLasRutas` (suma dos rutas).

### Errores nuevos

Ninguno. HU-23 reutiliza `ErrProyectoNoEncontrado` (HU-01) y `ErrSprintNoEncontrado` (HU-09). No usa `ErrSprintDeOtroProyecto` de HU-09: un sprint de otro proyecto se trata como inexistente (RN-23.11). El mensaje de la falla interna, "error interno al consultar las horas", vive en el adaptador.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/total_horas.go](../internal/domain/total_horas.go) (`HorasElemento`, `TotalHoras`, `SinEstimar`, `TotalizarHoras`; reutiliza `CompararEsfuerzo` de HU-19) | [total_horas_test.go](../internal/domain/total_horas_test.go) |
| Casos de uso y puertos | [internal/app/consultar_horas.go](../internal/app/consultar_horas.go) (`ConsultarHorasProyecto`, `ConsultarHorasSprint`, puertos `LectorHoras` y `LectorSprintPorID`; reutiliza `LectorProyectos` de HU-04) | [consultar_horas_test.go](../internal/app/consultar_horas_test.go) |
| Persistencia | [internal/adapters/persistencia/horas_repo.go](../internal/adapters/persistencia/horas_repo.go) (`HorasDeProyecto`, `HorasDeSprint`, `leerHoras`) | [horas_repo_test.go](../internal/adapters/persistencia/horas_repo_test.go), [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) |
| Endpoints JSON | [internal/adapters/http/consultar_horas_handler.go](../internal/adapters/http/consultar_horas_handler.go) (`RegistrarConsultarHoras`, `responderErrorHoras`), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [consultar_horas_handler_test.go](../internal/adapters/http/consultar_horas_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| Pantallas y plantillas | [internal/adapters/http/pantallas_horas.go](../internal/adapters/http/pantallas_horas.go) (`vistaHoras`, `avisoSinEstimar`, `errorDeHoras`), [pantallas_comparacion.go](../internal/adapters/http/pantallas_comparacion.go) (`textosDesviacion`), [pantallas.go](../internal/adapters/http/pantallas.go) (`mostrarDetalle`), [pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go) (`mostrarSprintDetalle`), [layout.html](../internal/adapters/http/web/plantillas/layout.html), [proyecto_detalle.html](../internal/adapters/http/web/plantillas/proyecto_detalle.html), [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | [pantallas_horas_test.go](../internal/adapters/http/pantallas_horas_test.go), [pantallas_test.go](../internal/adapters/http/pantallas_test.go) |
| BDD | [features/HU-23-horas-desviacion.feature](../features/HU-23-horas-desviacion.feature), [test/bdd/hu23_steps_test.go](../test/bdd/hu23_steps_test.go), [hu23_sistema_test.go](../test/bdd/hu23_sistema_test.go) (SQLite real), [bdd_test.go](../test/bdd/bdd_test.go); cableado de `ConsultarHorasProyecto` en [hu30_sistema_test.go](../test/bdd/hu30_sistema_test.go) | 29 escenarios y 240 pasos en verde |

### Commits

Salen de `git log 4209a72..6be9cf9`, en orden: desde que la rama salió de HU-19 hasta antes de integrar `main`. La última fila es del merge con `main`. HU-19 ya entró a `main` (PR #297), y `hu-19-comparar-esfuerzo..HEAD` ya no sirve para este listado: después de ese merge (`f9d513b`) también incluye commits de HU-13, HU-21 y HU-22.

| Ciclo | Commit | Etapa |
|---|---|---|
| Spec | `62200c5` docs(HU-23): spec SDD de comparar horas y desviación de un sprint o proyecto | — |
| BDD en rojo: feature | `e0b773f` test(HU-23): [RED] escenarios BDD del feature | BDD en rojo (28 escenarios y 227 pasos indefinidos) |
| BDD en rojo: steps | `0519df5` test(bdd): [RED] steps y sistema en memoria de HU-23 | BDD en rojo, por compilación |
| Dominio | `2637cad` [RED] → `e2b9766` [GREEN] → `18e761e` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Casos de uso | `b691bb0` [RED] → `55a8a43` [GREEN] → `cba036b` [REFACTOR] | **RED → GREEN → REFACTOR**; los 28 escenarios de entonces pasan con el sistema en memoria |
| Contexto en el caso de uso | `3813f25` test(HU-23): el lector de horas recibe el contexto del pedido | — (sin RED real, ver abajo) |
| Persistencia | `638c6c8` [RED] → `438bdb4` [GREEN] → `c55a522` [REFACTOR] | **RED → GREEN → REFACTOR** |
| BDD sobre SQLite | `538b2d8` test(bdd): HU-23 sistema SQLite real | BDD en verde |
| Endpoints JSON | `fc9a805` [RED] → `719a47f` [GREEN] → `b927ca0` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Pantallas | `bb380c7` [RED] → `91a31d1` [GREEN] → `cb719c6` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Spec: título del bloque | `0642f62` docs(HU-23): spec ajustada al bloque "Horas y desviación" | — |
| Aviso corto | `72baffb` [RED] → `2eec66a` [GREEN] | **RED → GREEN**; el GREEN incluye el ajuste de la spec |
| Spec: orden, contexto y no modificar | `6e246ec` docs(HU-23): spec con el orden de las filas, el contexto y que totalizar no modifica | — |
| Revisión: CA-23.10 | `d6c808d` test(HU-23): escenario de CA-23.10 con los dos elementos juntos en el sprint nuevo | — (sin RED real, ver abajo) |
| Revisión: sección 6 | `e5d8e46` test(HU-23): casos de la sección 6 con SQLite real | — (sin RED real, ver abajo) |
| Revisión: otro proyecto en la pantalla | `70f5631` test(HU-23): pantalla del sprint con un sprint de otro proyecto | — (sin RED real, ver abajo) |
| Revisión: datos de CA-23.2 en HTTP | `d9a38c1` test(HU-23): las rutas de horas prueban los datos de CA-23.2 con la misma traza | — (sin RED real, ver abajo) |
| Revisión: trazas | `6aba6af` test(HU-23): trazas y comentarios que citan lo que cada test verifica | — (solo comentarios) |
| Trazabilidad | `6be9cf9` docs(trazabilidad): HU-23 | — |
| Revisión: aclaración del sprint cerrado | `cff999d` test(HU-23): [RED] la aclaración del sprint cerrado incluye los elementos reasignados | Test de la pantalla en rojo |
| Revisión: aclaración del sprint cerrado | `c685c06` feat(HU-23): [GREEN] aclaración del sprint cerrado para todo elemento que salió de él | Test de la pantalla y BDD de HU-23 en verde |
| Spec: decisiones confirmadas | `69736ef` docs(HU-23): P-9 a P-12 confirmadas y spec implementada; corrige RN-19.8 de HU-19 | — |
| Integración con `main`: test ajustado | `f9d513b` `TestPantallasHoras_NoApareceEnOtrasPantallas` carga `ConsultarVelocidad`, resuelto dentro del merge | — (sin RED real: una línea de preparación, sin cambiar aserciones) |

**Merges de la rama, que no son de HU-23:** `f9d513b` (trae `origin/main`: HU-19, HU-13, HU-21 y HU-22).

**Commits que el filtro `git log --grep="(HU-23)"` no encuentra**, porque usan el prefijo `test(bdd):`: `0519df5` (RED de los steps) y `538b2d8` (sistema SQLite).

**Commits sin etapa TDD:** `3813f25`, `d6c808d`, `e5d8e46`, `70f5631`, `d9a38c1` y `6aba6af` (tests sin RED real), y `62200c5`, `0642f62` y `6e246ec` (spec).

### Evidencia: mutaciones manuales

Ninguna quedó en un commit. Se aplicaron en una copia descartable del repo. En total, 104 mutantes: 96 murieron y 8 sobrevivieron. De esos 8, M13 de los casos de uso murió después, con `3813f25`; los otros 7 siguen sobreviviendo, por los motivos de la tabla.

| Capa | Mutantes | Murieron | Sobrevivieron y por qué |
|---|---|---|---|
| Dominio (`total_horas.go`) | 14 | 12 | M5 (`SinEstimar` con `<=`) y M14 (desviación como suma por elemento): **equivalentes** |
| Casos de uso (`consultar_horas.go`) | 13 | 12 | M13 (`context.Background()` al lector): sobrevivió en la primera tanda; lo mata `3813f25` |
| Contexto, comprobación de `3813f25` | 2 | 2 | — (M13 en el sprint y en el proyecto) |
| Persistencia (`horas_repo.go`) | 15 | 13 | M14 (sin `ORDER BY`): SQLite ya devuelve las filas ordenadas. M15 (sin `rows.Err()`): no se puede provocar sin un driver simulado |
| Endpoints JSON | 20 | 19 | M20 (sin el log del 500): no cambia lo que ve el cliente |
| Pantallas | 25 | 25 | — |
| Aviso corto | 5 | 5 | — |
| Comprobación de la revisión | 10 | 8 | P1 (la lista de sprints no filtra) y P2 (el caso de uso no valida el proyecto): se cubren entre sí. P3, con las dos protecciones quitadas, muere |

**Lo que detectan los que más importan:**

| Mutación | Fallan |
|---|---|
| Sumar las reales sin estimar a la desviación | `TestTotalizarHoras`, `TestTotalizarHoras_ReutilizaCompararEsfuerzo` |
| Pasar 0 en vez de `nil` | `TestTotalizarHoras` |
| Promediar los porcentajes | `TestTotalizarHoras`, `TestTotalizarHoras_NuncaDevuelveCeroNegativo`, `TestTotalizarHoras_ReutilizaCompararEsfuerzo` |
| No validar que el sprint sea del proyecto | `TestConsultarHorasSprint_ValidacionFallida_RetornaErrorYNoConsultaHoras` |
| `LEFT JOIN` → `JOIN` | `TestRepositorioEsfuerzo_HorasDeProyecto_ElementoSinRegistros_ApareceConCero`, `TestRepositorioEsfuerzo_HorasDeSprint_SoloLosElementosDelSprint` |
| Filtrar por las fechas del sprint | `TestRepositorioEsfuerzo_HorasDeSprint_NoFiltraPorFechaNiPorIntegrante` |
| Sacar `sprint_id` de la respuesta | `TestGetHorasSprint_Devuelve200ConElTotal` |
| Filtrar el texto interno en el 500 | `TestGetHoras_FallaInterna_Devuelve500SinDetalle` |
| Desviación como `template.HTML` | `TestPantallasHoras_LosValoresSeEscapan` |
| Bloque solo en una de las dos pantallas | `TestPantallasHoras_ElBloqueVaAntesDeLaPrimeraSeccion` y otros |
| El sprint devuelve una sola fila, dentro de `@CA-23.10` | Solo el escenario "El elemento devuelto y el Hecho reasignado suman juntos en el sprint nuevo" |

**Cruce de pasos BDD:** un programa temporal, que no está en el repo, cruzó las 367 regex con los pasos de todos los `.feature` de esta rama: 0 pasos con más de una regex y 0 cruces entre HU-23 y otras historias. Plantando un paso duplicado, lo detectaron ese programa y godog.

**Prueba de humo:** `cmd/server` con una base temporal y `curl`, para las dos rutas JSON y las dos pantallas, también después de sus refactors. No se repitió después del aviso corto ni de los tests de la revisión.

### Etapas sin RED real

- **Tests agregados con el código ya correcto** (lo dice cada mensaje de commit):
  - `3813f25`: `TestConsultarHoras_ElLectorRecibeElContextoDelPedido`
  - `d6c808d`: el escenario "El elemento devuelto y el Hecho reasignado suman juntos en el sprint nuevo"
  - `e5d8e46`: los tres tests de la sección 6
  - `70f5631`: `TestPantallaHorasSprint_SprintDeOtroProyecto_Devuelve404`
  - `d9a38c1`: el caso 32 / 28 de `casosHoras`
  - `6aba6af`: solo comentarios
- **RED por compilación:** `0519df5`, `2637cad`, `b691bb0`, `638c6c8` y `fc9a805`. No se pudieron ver fallar las aserciones en esos commits. Para `0519df5` se usó un stub temporal, sin commitear: fallaron los 28 escenarios de HU-23 y pasaron los 336 de las demás historias.
- **RED por comportamiento:** en estos commits, los siguientes tests nuevos ya pasaban:
  - `bb380c7`: `TestPantallasHoras_NoRompenLoQueYaMostraban`, `TestPantallasHoras_NoApareceEnOtrasPantallas`, `TestPantallasHoras_RutaInvalidaOInexistente_NoConsultaLasHoras`; fijan lo que no debía cambiar
  - `72baffb`: los casos que ya existían de `TestPantallasHoras_MuestranLosValores`; fallan solo los seis subcasos nuevos (tres por pantalla)
- **Spec después del código:** el `ORDER BY` y el paso del contexto entraron con `438bdb4` y `719a47f`, y la spec los documentó recién en `6e246ec`. El bloque se implementó como "Horas y desviación" en `91a31d1`, y la spec se ajustó en `0642f62`.

### Huecos

- **CA-23.11 no tiene test Go:** solo el escenario BDD, que corre sobre SQLite real. Cumple con el BDD, pero no aparece si se busca por test.
- **Sin lector de pantalla:** las pantallas se probaron con tests y, en la revisión del PR #298, en un navegador con el servidor real. No se probaron con un lector de pantalla.
- **El contexto en las pantallas no tiene test:** lo pasan, pero solo los handlers JSON y los casos de uso lo verifican. La nota de la sección 5 habla solo de esos dos.
- **Cobertura de `leerHoras`: 84,6 %.** Faltan las ramas de error de `Scan` y `rows.Err()`, que no se pueden provocar sin un driver simulado. El resto del código nuevo de HU-23 está al 100 %.
- **`TestRepositorioEsfuerzo_HorasDeProyecto_OrdenaPorElemento`** detecta el orden descendente, pero no que falte el `ORDER BY` (M14). Los ID siguen siempre el orden de creación, así que no hay datos que lo distingan.
- **La conexión de `cmd/server/main.go` no tiene test automático:** solo el humo. Si faltara, el servidor fallaría al recibir el pedido.
- **BDD completo:** 387 a 403 s en las tres corridas tras integrar `main`, con 407 escenarios (antes, 280 a 293 s con 365). Sube porque cada escenario de cualquier historia crea una base más, ya que los hooks de todas las historias corren siempre.
- **Primer test de `internal/adapters/http` que usa `persistencia`:** `TestPantallaHorasSprint_SprintDeOtroProyecto_Devuelve404`, para usar casos de uso y SQLite reales. Solo afecta al test.
- **Resuelto, cruce de regex al integrar con `main`:** las 407 regex de `test/bdd/*_steps_test.go`, con el `regexp` de Go, contra los 2130 pasos de los 407 escenarios, con los Outline expandidos por el parser de Gherkin de godog: 0 ambiguos, 0 indefinidos y ninguna regex sin usar. El BDD completo, con `Strict: true`, lo confirma.
- **Resuelto, dependencia del PR de HU-19:** HU-19 ya entró a `main` (PR #297), y esta rama trae `CompararEsfuerzo` y `HorasEstimadas` desde ahí con el merge de `origin/main` (`f9d513b`), junto con HU-13, HU-21 y HU-22.
- **Test ajustado al integrar con `main` (`f9d513b`):** HU-22 hizo que la pantalla de sprints anteriores consulte `ConsultarVelocidad`, y `TestPantallasHoras_NoApareceEnOtrasPantallas` quedaba con ese campo en `nil` y entraba en panic. Se le cargó `consultorVelocidadFalso`, sin cambiar sus aserciones.
- **Humo repetido al integrar con `main`:** con el servidor real y una base temporal, por curl. Las rutas de horas, de velocidad y de defectos responden 200, y `id="total-horas"` aparece una vez en cada detalle, sin ids repetidos. Tampoco es una prueba en navegador.
- `mostrarSprintDetalle` y `mostrarDetalle` consultan el proyecto y el sprint, y los casos de uso de horas los vuelven a consultar: son lecturas de más, sin consecuencias prácticas.
- Los puntos planificados de un sprint cerrado (RN-11.6) incluyen lo no terminado, y estas horas no (P-3): no se comparan entre sí.

### Para revisar

Código de otras historias que tocó HU-23. Conviene que lo revisen sus autores en la PR.

| Archivo | Historia | Cambio |
|---|---|---|
| [internal/app/cerrar_sprint.go](../internal/app/cerrar_sprint.go) | HU-11 | Solo el puerto: `RepositorioCierre` incluye `LectorSprintPorID`, con los mismos métodos |
| [pantallas.go](../internal/adapters/http/pantallas.go) | HU-30 | Campo `Horas` y consulta en `mostrarDetalle` |
| [pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go) | HU-09 | Campo `Horas`, consulta y aclaración en `mostrarSprintDetalle` |
| [layout.html](../internal/adapters/http/web/plantillas/layout.html) | HU-30 | Plantilla compartida `bloque_horas` |
| [proyecto_detalle.html](../internal/adapters/http/web/plantillas/proyecto_detalle.html) | HU-30 | Incluye el bloque |
| [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | HU-09 | Incluye el bloque |
| [pantallas_comparacion.go](../internal/adapters/http/pantallas_comparacion.go) | HU-19 | Extracción de `textosDesviacion`, sin cambiar el comportamiento |
| [esfuerzo_repo.go](../internal/adapters/persistencia/esfuerzo_repo.go) | HU-18 y HU-19 | Solo el comentario del tipo |
| [base_cerrada_test.go](../internal/adapters/persistencia/base_cerrada_test.go) | Compartido | Dos casos |
| [router.go](../internal/adapters/http/router.go), [router_test.go](../internal/adapters/http/router_test.go), [cmd/server/main.go](../cmd/server/main.go) | Compartidos | Cableado y dos rutas; el diff se ve mayor por la realineación de `gofmt` |
| [pantallas_test.go](../internal/adapters/http/pantallas_test.go) | HU-30 | Dos dobles en `casosFalsos` |
| [test/bdd/bdd_test.go](../test/bdd/bdd_test.go) | Compartido | Registra `registrarStepsHU23` |
| [test/bdd/hu30_sistema_test.go](../test/bdd/hu30_sistema_test.go) | HU-30 | `ConsultarHorasProyecto`, porque ese sistema abre el detalle del proyecto |

---

## HU-24: Calcular % de story points completados y defectos por Sprint

- **Issue:** [#35](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/35) · Sprint 4 · 3 SP
- **Spec:** [docs/specs/HU-24-porcentaje-defectos-sprint.md](specs/HU-24-porcentaje-defectos-sprint.md)
- **Feature:** [features/HU-24-porcentaje-defectos-sprint.feature](../features/HU-24-porcentaje-defectos-sprint.feature) (steps en [test/bdd/hu24_steps_test.go](../test/bdd/hu24_steps_test.go), sobre SQLite real en [hu24_sistema_test.go](../test/bdd/hu24_sistema_test.go))
- **Alcance real:** HU-11 ya calcula y guarda los puntos planificados y completados de cada Sprint al cerrarlo (RN-11.6), y HU-20/HU-21 ya asignan el Sprint de detección y el de resolución de cada defecto (RN-20.5, RN-21.6). HU-24 agrega solo el **cálculo**: el porcentaje completado de cada Sprint Cerrado y el conteo de sus defectos detectados/resueltos. No agrega persistencia nueva (la tarea 24.6, adaptador de persistencia, no aplica).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-24.1 Porcentaje = completados sobre planificados, redondeado a un decimal | RN-24.1 | El porcentaje completado son los puntos completados sobre los planificados `@normal` · Un sprint cerrado con puntos planificados y 0 completados da 0%, no "sin datos" `@limite` · El porcentaje se redondea a un decimal, con los medios alejándose de cero `@limite` (16 de 1, 6,25% → 6,3%) | `TestPorcentajeCompletado`, `TestConsultarCalidadSprints_SprintsCerrados_InformaElPorcentaje`, `TestPantallaSprintsAnteriores_PorcentajeConDecimales_SeRedondeaAUnDecimal` | `domain.PorcentajeCompletado` (usa `redondearUnDecimal` de HU-19) |
| CA-24.2 Sin puntos planificados, sin porcentaje | RN-24.2 | Un sprint cerrado sin puntos planificados no tiene porcentaje `@limite` | `TestPorcentajeCompletado` (caso "sin puntos planificados"), `TestConsultarCalidadSprints_SinPuntosPlanificados_SinPorcentaje` | `domain.PorcentajeCompletado` |
| CA-24.3 Defectos detectados por Sprint | RN-24.3 | Cuenta los defectos detectados en cada sprint `@normal` | `TestContarDefectosDeSprint`, `TestConsultarCalidadSprints_CuentaLosDefectosPorSprint` | `domain.ContarDefectosDeSprint` |
| CA-24.4 Defectos resueltos por Sprint, aunque luego se cierren | RN-24.4 | Cuenta los defectos resueltos en un sprint aunque ya estén cerrados `@normal` | `TestContarDefectosDeSprint` (caso "aunque ya estén Cerrados"), `TestConsultarCalidadSprints_CuentaLosDefectosPorSprint` | `domain.ContarDefectosDeSprint` |
| CA-24.5 Defecto sin Sprint no cuenta en ninguno | RN-24.3, RN-24.4 | Un defecto sin sprint de detección no cuenta en ningún sprint `@alternativo` | `TestContarDefectosDeSprint` (caso "sin sprint"), `TestConsultarCalidadSprints_CuentaLosDefectosPorSprint` | `domain.ContarDefectosDeSprint` |
| CA-24.6 Solo los Sprints Cerrados | RN-24.5 | Un sprint abierto no aparece en la calidad de los sprints `@alternativo` | `TestConsultarCalidadSprints_SprintsCerrados_InformaElPorcentaje` (incluye un sprint Abierto que no aparece) | `app.ConsultarCalidadSprints` (`domain.SprintsAnteriores`) |
| CA-24.7 Proyecto inexistente | — (error reutilizado) | No se puede consultar la calidad de los sprints de un proyecto que no existe `@error` | `TestConsultarCalidadSprints_ProyectoInexistente_RetornaError` | `app.ConsultarCalidadSprints` |
| CA-24.8 Respuesta HTTP | Sección 3 | No va acá (como CA-22.7): se cubre con tests del adaptador | `TestGetCalidad_*` (3 tests) | `internal/adapters/http/calidad_sprints_handler.go` |
| CA-24.9 Pantalla | Sección 3 | No va acá (como CA-22.8): se cubre con tests de la pantalla | `TestPantallaSprintsAnteriores_MuestraElPorcentajeYLosDefectos`, `TestPantallaSprintsAnteriores_SinPuntosPlanificados_AvisaSinPorcentaje` | `pantallas_sprints.go` (`sprintsAnteriores`, reutiliza `textoPorcentaje` de HU-19/HU-23), `web/plantillas/sprints_anteriores.html` |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-24.1 % = PuntosCompletados / PuntosPlanificados × 100, sin recalcular lo que ya guardó HU-11, redondeado a un decimal | CA-24.1 | `domain.PorcentajeCompletado` |
| RN-24.2 Sin puntos planificados, el porcentaje queda indefinido (no 0%) | CA-24.2 | `domain.PorcentajeCompletado` |
| RN-24.3 Detectados = defectos con ese Sprint como SprintDeteccionID | CA-24.3, CA-24.5 | `domain.ContarDefectosDeSprint` |
| RN-24.4 Resueltos = defectos con ese Sprint como SprintResolucionID, sin importar el estado actual | CA-24.4, CA-24.5 | `domain.ContarDefectosDeSprint` |
| RN-24.5 Solo se calcula para Sprints Cerrados | CA-24.6 | `app.ConsultarCalidadSprints` (reutiliza `domain.SprintsAnteriores` de HU-12) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/sprint_historial.go](../internal/domain/sprint_historial.go) (`PorcentajeCompletado`, `ContarDefectosDeSprint`, junto a `SprintsAnteriores` de HU-12 y `Velocidad` de HU-22) | [sprint_historial_test.go](../internal/domain/sprint_historial_test.go) |
| Caso de uso | [internal/app/consultar_calidad_sprints.go](../internal/app/consultar_calidad_sprints.go) (`ConsultarCalidadSprints`, reutiliza `LectorProyectos`, `RepositorioSprints` de HU-08 y `LectorDefectos` de HU-21) | [consultar_calidad_sprints_test.go](../internal/app/consultar_calidad_sprints_test.go) |
| HTTP (JSON) | [internal/adapters/http/calidad_sprints_handler.go](../internal/adapters/http/calidad_sprints_handler.go) (`GET /proyectos/{id}/calidad`) | [calidad_sprints_handler_test.go](../internal/adapters/http/calidad_sprints_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_sprints.go](../internal/adapters/http/pantallas_sprints.go) (reutiliza la pantalla de HU-12/HU-22), [web/plantillas/sprints_anteriores.html](../internal/adapters/http/web/plantillas/sprints_anteriores.html) | [pantallas_sprints_anteriores_test.go](../internal/adapters/http/pantallas_sprints_anteriores_test.go) |
| BDD | [test/bdd/hu24_sistema_test.go](../test/bdd/hu24_sistema_test.go), [hu24_steps_test.go](../test/bdd/hu24_steps_test.go) | 9 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 24.1 Spec | `01f825e` docs(HU-24): spec SDD | — |
| 24.2 Escenarios BDD | `02ada65` test(HU-24): [24.2] escenarios BDD | — |
| 24.3 Steps de godog | `7447650` test(HU-24): [24.3][RED] steps (en rojo hasta el caso de uso) | BDD en rojo |
| 24.4 Dominio (`PorcentajeCompletado`, `ContarDefectosDeSprint`) | `87d989e` [RED] → `09a168b` [GREEN] | **RED → GREEN** (sin REFACTOR: dos funciones cortas) |
| 24.5 Caso de uso `ConsultarCalidadSprints` | `1c4c1a0` [RED] → `a5bac14` [GREEN] | **RED → GREEN**; BDD en verde (incluyó un fix chico: el step `"X" tiene un sprint abierto` colisionaba con el de HU-22, se distinguió como "tiene un sprint todavía abierto") |
| 24.6 Adaptador de persistencia | No aplica: sin persistencia nueva (reutiliza `RepositorioSprints` de HU-08 y `RepositorioDefectos`/`LectorDefectos` de HU-21) | — |
| 24.7 Adaptador HTTP (`GET .../calidad`) | `8d97166` [RED] → `189ada9` [GREEN] | **RED → GREEN** |
| 24.8 Pantalla (sumada a la de HU-12/HU-22) | `27bbdd8` [RED] → `c4d9253` [GREEN] | **RED → GREEN**; incluyó el fix de un fixture de test (`TestPantallasHoras_NoApareceEnOtrasPantallas`) que armaba un `CasosDeUso` parcial sin `ConsultarCalidadSprints` (mismo patrón de interfaz nil de HU-21) |
| 24.9 Revisión y trazabilidad | `1b44d60` | — |
| 24.10 Corrección de la revisión del PR #366 (ronda 1) | `f681fa5` (spec) → `8721836` [RED] → `86a7137` [GREEN] (dominio) → `f0d11d4` [RED] → `6119967` [GREEN] (pantalla) → `3494b7e` (escenario BDD) → `e5033cf` (trazabilidad) | **RED → GREEN** ×2 |
| Traer `main` (conflicto con HU-29, PR #365, en trazabilidad.md) | `28dfc3e` | — |
| 24.11 Corrección de la revisión del PR #366 (ronda 2, sugerencias opcionales) | `5c05e8b` [RED] → `61bbd61` [GREEN] → `fb91f2f` (spec) | **RED → GREEN** ×1 |

### Huecos

- ~~La spec está en Borrador~~: resuelto el 2026-11-05. Matías López (Agile Enabler) confirmó P-1 y P-2 en la revisión del PR [#366](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/366) y pidió redondear el porcentaje a un decimal (P-3, antes abierta); se corrigió también un caso límite de la sección 6 que en realidad no puede pasar con los datos reales de HU-11 (RN-11.6). Spec Aprobada (tarea 24.10).
- Probado manualmente en el navegador (servidor real, SQLite): proyecto con un Sprint cerrado (8 planificados, 8 completados) y un defecto detectado y resuelto en ese mismo Sprint → la pantalla de sprints anteriores mostró correctamente 100,0%, 1 defecto detectado y 1 resuelto.
- Revisión con el agente `go-reviewer`: **Approve**, sin hallazgos bloqueantes. Confirmó que el bug de HU-21 (method value de una interfaz nil al registrar rutas) no se repite acá (tanto `RegistrarCalidadSprints` como el handler y la pantalla llaman al caso de uso recién dentro del closure, en el momento del pedido), y que no queda ningún otro `CasosDeUso` parcial en los tests existentes sin el campo `ConsultarCalidadSprints` más allá del que ya se corrigió en `pantallas_horas_test.go`. Verificó `go vet`, `gofmt` y los tests de dominio/app/HTTP; no llegó a correr la suite BDD completa por el tiempo que tarda en ese entorno, pero revisó el código BDD a mano (SQL parametrizado, cada escenario aislado en su propia base temporal) y lo encontró consistente con el patrón de HU-22.
- ~~PR #365 (HU-29, en paralelo) también agrega su sección a este archivo~~: resuelto el 2026-11-06 al traer `main` (HU-29 se mergeó primero, PR #365): se dejaron las dos filas del resumen y las dos secciones, como estaba previsto.
- Segunda ronda de revisión del PR #366 (2026-11-06): además de traer `main`, se aplicaron las tres sugerencias opcionales de Matías López: se sacó el escenario BDD "el cálculo no trunca si los completados superan los planificados" (imposible según la spec, sección 5; el test unitario del dominio ya documenta que la función pura no lo valida), el ejemplo de redondeo pasó de 2/3 (no es un medio) a 16/1 = 6,25% → 6,3% (un medio real, prueba "los medios alejándose de cero"), y la pantalla ahora reutiliza `textoPorcentaje` de HU-19/HU-23 en vez de un `formatearPorcentaje` propio, para que todas las pantallas de porcentaje muestren el mismo formato ("75,0 %", con espacio).

---

## HU-29: Documentación técnica del proyecto

- **Issue:** [#40](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/40) · Sprint 4 · 2 SP
- **Tareas:** 29.1 [#357](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/357), 29.2 [#358](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/358), 29.3 [#359](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/359), 29.4 [#360](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/360)
- **Sin spec, BDD ni TDD:** es una historia de documentación y no agrega funcionalidad ni reglas de negocio. SDD, BDD y TDD se aplican a la funcionalidad (`CLAUDE.md`). Los criterios de aceptación se verifican en la revisión del PR.

### Criterios de aceptación

| Criterio (issue #40) | Dónde se cumple | Cómo se verificó |
|---|---|---|
| Describe la arquitectura general | [docs/arquitectura.md](arquitectura.md): capas, puertos y adaptadores, cómo viaja un pedido, tests por capa y dónde va cada cosa al sumar una historia | Se cotejó contra el código: `go list` confirma que el dominio solo usa la biblioteca estándar y que `internal/app` solo depende del dominio, y existen todos los nombres citados |
| Documenta decisiones técnicas relevantes (stack, persistencia, etc.) | Sección "Decisiones técnicas" de [docs/arquitectura.md](arquitectura.md), que enlaza los ADR; [ADR-0002](adr/0002-persistencia.md), [ADR-0003](adr/0003-interfaz-web.md) y [ADR-0004](adr/0004-ci-cd.md) actualizados con notas fechadas | Estado del CI, cobertura y protección de `main` leídos en GitHub; despliegue en Render probado con el pipeline del merge del PR #299 |
| Se mantiene versionada junto al código | Todo está en el repo, en `docs/` y en el [README](../README.md) | El README se siguió paso a paso en un clon nuevo: el servidor arrancó, `/salud` y `/web/proyectos` respondieron y corrió el ejemplo de `BDD_TAGS` |

### Commits

| Etapa | Commit |
|---|---|
| 29.1 Arquitectura general | `46a1181` docs(HU-29): arquitectura general |
| 29.2 Decisiones técnicas y ADR | `fbf95f7` docs(HU-29): decisiones técnicas y ADR al día |
| 29.3 README | `a13c526` docs(HU-29): README y puesta en marcha |
| Render activo | `547fd30` docs(HU-29): despliegue en Render activo |
| 29.4 Trazabilidad | docs(trazabilidad): HU-29 |

### Huecos

- **Los diagramas Mermaid** de `docs/arquitectura.md` no se vieron renderizados antes del PR: GitHub los dibuja al abrirlo.
- **Los ADR-0005 (gráficos) y 0006 (PDF)** están pendientes: los escriben HU-25 y HU-27, y la tabla de decisiones técnicas los nombra. Cuando se mergeen, hay que sumarlos a esa tabla.

---

## HU-25: Ver panel con estado del proyecto y gráficos

- **Issue:** [#36](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/36) · Sprint 4 · 8 SP
- **Spec:** [docs/specs/HU-25-dashboard.md](specs/HU-25-dashboard.md)
- **ADR:** [docs/adr/0005-graficos-dashboard.md](adr/0005-graficos-dashboard.md) (gráficos como SVG armado en Go, sin JavaScript)
- **Feature:** [features/HU-25-dashboard.feature](../features/HU-25-dashboard.feature) (steps en [test/bdd/hu25_steps_test.go](../test/bdd/hu25_steps_test.go), sobre SQLite real en [hu25_sistema_test.go](../test/bdd/hu25_sistema_test.go))
- **Alcance real:** las métricas ya existen (HU-11, HU-22, HU-23 y HU-24). HU-25 cuenta el backlog por estado, resume el sprint en curso con la regla del cierre (RN-11.6), ordena los sprints cerrados para los gráficos y lo muestra todo en una pantalla nueva. No agrega persistencia nueva (la tarea 25.7, adaptador de persistencia, no aplica).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-25.1 Backlog por estado | RN-25.2 | El backlog se cuenta por estado `@normal` | `TestContarBacklog`, `TestConsultarDashboard_JuntaTodasLasMetricas`, `TestPantallaDashboard_MuestraElEstadoDelProyecto` | `domain.ContarBacklog` |
| CA-25.2 Sprint en curso con sus puntos | RN-25.3, RN-25.4 | El sprint en curso suma los puntos de sus historias con la regla del cierre `@normal` · Una historia sin estimar del sprint en curso suma 0 puntos `@limite` | `TestResumirSprintEnCurso`, `TestConsultarDashboard_JuntaTodasLasMetricas`, `TestPantallaDashboard_MuestraElEstadoDelProyecto` | `domain.ResumirSprintEnCurso` (reutiliza `SprintVigente` de HU-20 y `ResumirCierre` de HU-11) |
| CA-25.3 Sin sprint en curso | RN-25.3 (P-1) | Un sprint abierto que ya venció no es el sprint en curso `@alternativo` | `TestResumirSprintEnCurso` (vencido, futuro, cerrado), `TestConsultarDashboard_SprintAbiertoVencido_NoEstaEnCurso`, `TestPantallaDashboard_SinSprintEnCurso` | `domain.ResumirSprintEnCurso` |
| CA-25.4 Sprints cerrados en orden, con sus métricas y la velocidad | RN-25.5 (P-2) | Los sprints cerrados aparecen del más antiguo al más reciente, con sus métricas `@normal` · Un sprint abierto no aparece entre los sprints cerrados `@alternativo` | `TestSprintsCerradosCronologicos` (incluye 7 sprints, sin límite), `TestConsultarDashboard_JuntaTodasLasMetricas` | `domain.SprintsCerradosCronologicos`, `app.ConsultarDashboard` (reutiliza `PorcentajeCompletado`, `ContarDefectosDeSprint` y `Velocidad`) |
| CA-25.5 Horas y desviación del proyecto, con el aviso de elementos sin estimar | RN-23.4 a RN-23.6 (HU-23) | Las horas del proyecto son las estimadas contra las reales, con su desviación `@normal` · Sin elementos estimados en horas, la desviación queda sin datos `@limite` | `TestConsultarDashboard_JuntaTodasLasMetricas`, `TestConsultarDashboard_ProyectoSinDatos`, `TestPantallaDashboard_HorasConElementosSinEstimar_MuestraElAviso`, `TestPantallaDashboard_HorasTodoEstimado_SinAviso` | `app.ConsultarDashboard` (reutiliza `TotalizarHoras` de HU-23), `pantallas_dashboard.go` (reutiliza `avisoSinEstimar` de HU-23) |
| CA-25.6 Datos vigentes | RN-25.1 | El dashboard muestra los datos vigentes al volver a abrirlo `@normal` | `TestConsultarDashboard_LeeLosDatosVigentesEnCadaConsulta` | `app.ConsultarDashboard` (lee todo de nuevo en cada consulta) |
| CA-25.7 Proyecto inexistente | — (error reutilizado) | No se puede consultar el dashboard de un proyecto que no existe `@error` | `TestConsultarDashboard_ProyectoInexistente`, `TestGetDashboard_ProyectoInexistente_Devuelve404`, `TestPantallaDashboard_ProyectoInexistente_Devuelve404` | `app.ConsultarDashboard` |
| CA-25.8 Respuesta HTTP | Sección 3 | No va acá (como CA-24.8): se cubre con tests del adaptador | `TestGetDashboard_*` (6 tests) | `internal/adapters/http/dashboard_handler.go` |
| CA-25.9 Gráficos en la pantalla | RN-25.5, ADR-0005 | No va acá: se cubre con tests de la pantalla y de la geometría | `TestPantallaDashboard_MuestraLosGraficos`, `TestPantallaDashboard_GraficoConSuAnchoReal`, `TestArmarGraficoBarras_*` (7 tests) | `pantallas_dashboard.go`, `grafico_barras.go`, `web/plantillas/dashboard.html` |
| CA-25.10 Enlace desde el proyecto | RN-25.6 | No va acá: se cubre con tests de la pantalla | `TestPantallaProyecto_EnlaceAlDashboard` | `web/plantillas/proyecto_detalle.html` |
| CA-25.11 Proyecto sin datos | Sección 6 | Un proyecto recién creado tiene el dashboard vacío, sin fallar `@limite` | `TestContarBacklog`, `TestSprintsCerradosCronologicos`, `TestConsultarDashboard_ProyectoSinDatos`, `TestGetDashboard_SinDatos_DevuelveNullYArrayVacio`, `TestPantallaDashboard_ProyectoSinDatos`, `TestArmarGraficoBarras_TodoEnCero_SinDividirPorCero` | todas las capas |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-25.1 Datos vigentes en cada pedido, sin cache | CA-25.6 | `app.ConsultarDashboard`, `pantallas.dashboard` (llama al caso de uso al atender el pedido) |
| RN-25.2 Backlog contado por estado actual, esté o no en un sprint | CA-25.1, CA-25.11 | `domain.ContarBacklog` |
| RN-25.3 Sprint en curso = Abierto cuyas fechas incluyen hoy | CA-25.2, CA-25.3 | `domain.ResumirSprintEnCurso` (`domain.SprintVigente` de HU-20) |
| RN-25.4 Puntos del sprint en curso con la regla del cierre | CA-25.2 | `domain.ResumirSprintEnCurso` (`domain.ResumirCierre` de HU-11) |
| RN-25.5 Gráficos con todos los sprints cerrados, del más antiguo al más reciente, con los valores de HU-11/22/24 | CA-25.4, CA-25.9 | `domain.SprintsCerradosCronologicos`, `app.ConsultarDashboard` |
| RN-25.6 Un clic desde la pantalla del proyecto | CA-25.10 | `web/plantillas/proyecto_detalle.html` |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/dashboard.go](../internal/domain/dashboard.go) (`ContarBacklog`, `ResumirSprintEnCurso`, `SprintsCerradosCronologicos`) | [dashboard_test.go](../internal/domain/dashboard_test.go) |
| Caso de uso | [internal/app/consultar_dashboard.go](../internal/app/consultar_dashboard.go) (`ConsultarDashboard`, reutiliza `LectorProyectos`, `RepositorioSprints`, `ListadorElementos`, `LectorDefectos` y la parte de proyecto de `LectorHoras`) | [consultar_dashboard_test.go](../internal/app/consultar_dashboard_test.go) |
| HTTP (JSON) | [internal/adapters/http/dashboard_handler.go](../internal/adapters/http/dashboard_handler.go) (`GET /proyectos/{id}/dashboard`) | [dashboard_handler_test.go](../internal/adapters/http/dashboard_handler_test.go) |
| Pantalla | [internal/adapters/http/pantallas_dashboard.go](../internal/adapters/http/pantallas_dashboard.go), [grafico_barras.go](../internal/adapters/http/grafico_barras.go), [web/plantillas/dashboard.html](../internal/adapters/http/web/plantillas/dashboard.html), [web/static/estilos.css](../internal/adapters/http/web/static/estilos.css) | [pantallas_dashboard_test.go](../internal/adapters/http/pantallas_dashboard_test.go), [grafico_barras_test.go](../internal/adapters/http/grafico_barras_test.go) |
| BDD | [test/bdd/hu25_sistema_test.go](../test/bdd/hu25_sistema_test.go), [hu25_steps_test.go](../test/bdd/hu25_steps_test.go) | 11 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 25.1 ADR de la librería de gráficos | `5fa5fcf` docs(HU-25): [25.1] ADR-0005 | — |
| 25.2 Spec | `0c40b25` docs(HU-25): [25.2] spec SDD → `cd0c207` (P-1 y P-2 confirmadas) | — |
| 25.3 Escenarios BDD | `ac1962d` test(HU-25): [25.3] escenarios BDD | — |
| 25.4 Steps de godog | `4f9bf41` test(HU-25): [25.4][RED] steps (en rojo hasta el caso de uso) | BDD en rojo |
| 25.5 Dominio (`ContarBacklog`, `ResumirSprintEnCurso`, `SprintsCerradosCronologicos`) | `1fd0020` [RED] → `3161194` [GREEN] | **RED → GREEN** |
| 25.6 Caso de uso `ConsultarDashboard` | `2f48728` [RED] → `ca9fb7a` [GREEN] → `62a9fa8` [REFACTOR] → `4954ca2` (BDD en verde) | **RED → GREEN → REFACTOR** (`calidadDe`, compartida con HU-24) |
| 25.7 Adaptador de persistencia | No aplica: sin persistencia nueva | — |
| 25.8 Adaptador HTTP (`GET .../dashboard`) | `80b87be` [RED] → `99b5420` [GREEN] | **RED → GREEN** |
| 25.9 Pantalla con los gráficos | `5f09d89` [RED] → `5a23bb9` [GREEN] (geometría de los SVG) · `90c5645` [RED] → `0e6c59e` [GREEN] (pantalla y enlace) | **RED → GREEN** ×2 |
| 25.10 Revisión y trazabilidad | `958d0c9` [RED] → `08df6ba` [GREEN] (corrección de la revisión: el SVG con su ancho real) · `78c04a0` (trazabilidad) | **RED → GREEN** |
| 25.11 Corrección de la revisión del PR #375 | `5fc949f` [RED] → `3ab9f6e` [GREEN] (aviso de horas sin estimar) · este commit (ADR aceptado, spec aprobada, textos y trazabilidad) | **RED → GREEN** |

### Huecos

- ~~La spec y el ADR-0005 los aprobó solo Nicolás Andino~~: resuelto el 2026-11-07. Matías López (Agile Enabler) aprobó el ADR-0005, P-1 y P-2 en la revisión del PR [#375](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/375). ADR-0005 pasa a *accepted* y la spec a Aprobada (tarea 25.11).
- **RED por compilación:** los RED `1fd0020` (dominio), `2f48728` (caso de uso), `80b87be` (handler JSON) y `5f09d89` (geometría de los SVG) no compilan en lugar de fallar las aserciones, porque la función todavía no existía. En Go es lo normal, pero es evidencia más débil que un RED que falla por la verificación. Los RED de la pantalla (`90c5645`, `958d0c9` y `5fc949f`) sí fallan por las aserciones.
- Revisión del PR #375 (2026-11-07): con elementos sin estimar, el gráfico de horas mostraba, por ejemplo, 10 h contra 18 h y "Desviación: 20,0 %", correcto según RN-23.5 pero sin explicación. Ahora muestra debajo el mismo aviso que el bloque de horas de HU-23 (`avisoSinEstimar`). Se corrigieron también dos textos: el ADR-0005 decía que la política de seguridad era de HU-29 (la agregó HU-30; HU-29 la documentó), y CA-25.5 citaba RN-25.5 en lugar de las reglas de horas de HU-23 (RN-23.4 a RN-23.6).
- Probado manualmente en el navegador (servidor real, SQLite, datos cargados por la API): 3 sprints cerrados (10/8, 13/13 y 3/0 puntos), un sprint en curso con 16 puntos y 5 completados, y un elemento con 10 horas estimadas y 12 reales. El JSON y la pantalla mostraron lo esperado (velocidad 7, 80 %, 100 % y 0 %, desviación del 20,0 %), sin errores en la consola (la CSP no bloqueó nada). En el ancho de un celular (375 px) la página no tiene scroll horizontal.
- Revisión con el agente `go-reviewer`: **Approve**, sin hallazgos CRITICAL ni HIGH. Verificó `go vet`, `gofmt` y los tests de `internal`; que no hay divisiones por cero ni punteros compartidos entre iteraciones; que la pantalla y el handler llaman al caso de uso al atender el pedido (el bug de interfaz nil de HU-21 no se repite), y que la plantilla no tiene scripts ni atributos `style` y escapa todo dentro del SVG. Hallazgos aplicados: (MEDIUM) con muchos sprints el SVG se achicaba al ancho de la pantalla y el texto quedaba ilegible; ahora se dibuja con su ancho real y el recuadro se desplaza en horizontal, sin scroll en la página (verificado a 375 px); (LOW) un fallo interno del dashboard quedaba en el log como "consultar proyecto". Hallazgos descartados: el orden no estable de `SprintsAnteriores` con fechas de inicio iguales no puede pasar, porque RN-08.5 impide sprints superpuestos; buscar el proyecto dos veces (pantalla y caso de uso) sigue el patrón del resto de las pantallas.
- Los defectos de sprints ya cerrados no se pueden cargar a mano desde la interfaz (HU-20 asigna el sprint vigente del día), así que el gráfico de defectos se probó con datos reales solo con ceros; con valores se cubre en los tests y en el BDD.

---

## HU-26: Generar reporte de proyecto o Sprint

- **Issue:** [#37](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/37) · Sprint 4 · 5 SP
- **Spec:** [docs/specs/HU-26-reporte-proyecto-sprint.md](specs/HU-26-reporte-proyecto-sprint.md) (**Aprobada** el 2026-11-08 por Matías López, Agile Enabler, en la revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378))
- **Feature:** [features/HU-26-reporte-proyecto-sprint.feature](../features/HU-26-reporte-proyecto-sprint.feature), 13 escenarios (steps en [hu26_steps_test.go](../test/bdd/hu26_steps_test.go), sobre SQLite real en [hu26_sistema_test.go](../test/bdd/hu26_sistema_test.go))
- **Alcance real:** no agrega ningún cálculo nuevo. Cada una de las cinco secciones del reporte (Resumen, Historias y estimaciones, Esfuerzo, Métricas, Defectos) reutiliza tal cual un caso de uso o una función del dominio que ya existe (HU-04, HU-07, HU-09, HU-21, HU-22, HU-23, HU-24); lo único que agrega HU-26 es juntar esos datos en un `Reporte` y mostrarlos en una pantalla propia, sin persistencia nueva (tarea 26.6 no aplica, mismo caso que la 24.6 de HU-24).
- **Exportar a PDF, de HU-27:** por decisión de Matías en la revisión del PR #378, la conversión a `DocumentoPDF` queda enteramente afuera de HU-26 (antes RN-26.11/CA-26.14; ver "Huecos" para el estado anterior). En su lugar, HU-26 agrega una aclaración de texto en el reporte de un Sprint Cerrado (nueva RN-26.11/CA-26.14, sin relación con HU-27).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-26.1 La sección de historias trae todas, con título/prioridad/estado/story points | RN-26.1 | El reporte de un proyecto trae todas sus historias `@normal` | `TestArmarReporteProyecto_HistoriasYResumen` (dominio), `TestArmarReporteProyecto_TraeTodasLasHistorias` (caso de uso), `TestPantallaReporteProyecto_MuestraLasCincoSecciones` (pantalla) | `domain.ArmarReporteProyecto`, `app.ArmarReporteProyecto.Ejecutar` (`app.ListarBacklog`, HU-07) |
| CA-26.2 El reporte de un Sprint trae solo sus historias asignadas | RN-26.2 | El reporte de un Sprint trae solo sus historias asignadas `@normal` | `TestHistoriasDeSprint` (dominio), `TestArmarReporteSprint_SoloHistoriasAsignadas` (caso de uso), `TestPantallaReporteSprint_Muestra` (pantalla) | `domain.HistoriasDeSprint` (mismo filtro que HU-09) |
| CA-26.3 Resumen: "X de Y historias completadas (Z de W story points)" | RN-26.3 | El resumen cuenta las historias y los story points completados `@normal` | `TestArmarReporteProyecto_HistoriasYResumen`, `TestArmarReporteProyecto_HistoriaSinEstimar_NoSumaStoryPoints` (dominio); `TestPantallaReporteProyecto_MuestraLasCincoSecciones` (pantalla) | `domain.ArmarReporteProyecto`/`ArmarReporteSprint` (`resumenHistorias`) |
| CA-26.4 El esfuerzo coincide con lo que calcula HU-23 | RN-26.4 | El esfuerzo del reporte coincide con las horas de HU-23 `@normal` | `TestArmarReporteProyecto_EsfuerzoYMetricas_SePropaganTalCual` (dominio) | `app.ArmarReporteProyecto`/`ArmarReporteSprint.Ejecutar` (`domain.TotalizarHoras`, HU-23) |
| CA-26.5 Métricas de proyecto: velocidad + una fila por Sprint Cerrado | RN-26.5 | Las métricas de un proyecto traen la velocidad y una fila por sprint cerrado `@normal` | `TestArmarReporteProyecto_EsfuerzoYMetricas_SePropaganTalCual` (dominio), `TestArmarReporteProyecto_Metricas` (caso de uso) | `app.ArmarReporteProyecto.Ejecutar` (`domain.Velocidad` HU-22, `domain.PorcentajeCompletado`/`ContarDefectosDeSprint` HU-24; ver "Huecos" sobre la duplicación con `ConsultarVelocidad`/`ConsultarCalidadSprints`) |
| CA-26.6 Sin Sprints cerrados, sin velocidad y sin filas de métricas | RN-26.5, RN-26.9 | Sin sprints cerrados, el reporte de un proyecto no tiene métricas `@limite` | `TestArmarReporteProyecto_SinMetricas` (dominio), `TestArmarReporteProyecto_SinSprintsCerrados_SinMetricas` (caso de uso) | `domain.ArmarReporteProyecto` / `app.ArmarReporteProyecto.Ejecutar` |
| CA-26.7 Sprint Cerrado con puntos: métrica con su porcentaje real | RN-26.6 | Las métricas de un Sprint cerrado traen su porcentaje real `@normal` | `TestArmarReporteSprint_Metrica_SePropagaTalCual` (dominio), `TestArmarReporteSprint_AbiertoYCerrado` (caso de uso) | `app.ArmarReporteSprint.Ejecutar` (`domain.PorcentajeCompletado` sin exigir `Cerrado`) |
| CA-26.8 Sprint Abierto: "Sin datos" en el porcentaje, pero cuenta sus defectos | RN-26.6 | Un Sprint abierto no tiene datos de porcentaje, pero cuenta sus defectos `@alternativo` | `TestArmarReporteSprint_Metrica_SePropagaTalCual`, `TestArmarReporteSprint_AbiertoYCerrado` | `app.ArmarReporteSprint.Ejecutar` |
| CA-26.9 El reporte de un proyecto trae todos sus defectos | RN-26.7 | El reporte de un proyecto trae todos sus defectos `@normal` | `TestArmarReporteProyecto_Defectos_SePropaganTalCual` (dominio), `TestPantallaReporteProyecto_MuestraLasCincoSecciones` (pantalla) | `app.ArmarReporteProyecto.Ejecutar` (`app.ListarDefectos`, HU-21) |
| CA-26.10 Defecto detectado en un Sprint y resuelto en otro: aparece en los dos reportes | RN-26.8 | Un defecto detectado en un Sprint y resuelto en otro aparece en los dos reportes `@alternativo` | `TestDefectosDeSprint` (dominio), `TestArmarReporteSprint_Defectos` (caso de uso) | `domain.DefectosDeSprint` (mismo criterio que `ContarDefectosDeSprint`, HU-24) |
| CA-26.11 Proyecto inexistente → 404 "proyecto no encontrado" | Sección 7 | No se puede consultar el reporte de un proyecto que no existe `@error` | `TestArmarReporteProyecto_ProyectoInexistente_RetornaError` (caso de uso), `TestPantallaReporteProyecto_Errores` (caso "proyecto inexistente") | `domain.ErrProyectoNoEncontrado` (HU-01, reutilizado) |
| CA-26.12 Sprint inexistente o de otro proyecto → 404 "sprint no encontrado" | RN-26.10, sección 7 | No se puede consultar el reporte de un Sprint que no existe `@error` | `TestArmarReporteSprint_Errores` (caso de uso), `TestPantallaReporteSprint_Errores` (casos "sprint inexistente") | `domain.ErrSprintNoEncontrado` (HU-09, reutilizado) |
| CA-26.13 ID inválido → 400, sin llamar al caso de uso | Sección 7 | No va acá (como CA-24.8): se cubre con tests de adaptador | `TestPantallaReporteProyecto_Errores` (caso "id inválido"), `TestPantallaReporteSprint_Errores` (caso "id de sprint inválido") | `pantallas_reporte.go` (reutiliza `idDeRuta`, `sprintIDDesdeRuta`) |
| CA-26.14 Sprint Cerrado: la sección de Métricas muestra la aclaración; Sprint Abierto, no | RN-26.11 | No va acá, como CA-26.13: es un detalle de la pantalla, no del dominio ni del caso de uso (mismo criterio que la aclaración de HU-23, RN-23.8, que tampoco tiene escenario BDD) | `TestPantallaReporteSprint_Cerrado_MuestraAclaracion`, `TestPantallaReporteSprint_Abierto_NoMuestraAclaracion`, `TestPantallaReporteProyecto_NuncaMuestraAclaracionDeSprint` | `pantallas_reporte.go` (`textoAclaracionReporteSprintCerrado`), `reporte.html` |
| CA-26.15 Proyecto sin datos: cada sección (salvo Resumen) muestra su mensaje, sin romper | RN-26.9 | El reporte de un proyecto sin datos no rompe `@limite` | `TestPantallaReporteProyecto_SinDatos_MuestraMensajes` (pantalla) | `reporte.html` (mensajes de "sin datos") |

### Reglas de negocio

| RN | CA que la cubren | Código |
|---|---|---|
| RN-26.1 El reporte de un proyecto incluye todas las historias de su backlog | CA-26.1 | `app.ArmarReporteProyecto.Ejecutar` (`app.ListarBacklog`, HU-07) |
| RN-26.2 El reporte de un Sprint incluye solo las historias con ese `SprintID` | CA-26.2 | `domain.HistoriasDeSprint` |
| RN-26.3 Completada = `Estado Hecho`; el resumen cuenta historias y story points completados sobre el total | CA-26.3 | `domain.ArmarReporteProyecto`/`ArmarReporteSprint` (`resumenHistorias`) |
| RN-26.4 El esfuerzo se calcula con `domain.TotalizarHoras` sobre las mismas historias del reporte | CA-26.4 | `app.ArmarReporteProyecto`/`ArmarReporteSprint.Ejecutar` |
| RN-26.5 Métricas de proyecto: velocidad del equipo + una fila por Sprint Cerrado | CA-26.5, CA-26.6 | `app.ArmarReporteProyecto.Ejecutar` |
| RN-26.6 Métricas de Sprint: su propio porcentaje y defectos, sin exigir `Cerrado` | CA-26.7, CA-26.8 | `app.ArmarReporteSprint.Ejecutar` |
| RN-26.7 Defectos de proyecto: todos, sin filtrar por Sprint | CA-26.9 | `app.ArmarReporteProyecto.Ejecutar` (`app.ListarDefectos`, HU-21) |
| RN-26.8 Defectos de Sprint: los detectados o resueltos ahí | CA-26.10 | `domain.DefectosDeSprint` |
| RN-26.9 Una sección sin datos muestra un mensaje, no una tabla vacía | CA-26.6, CA-26.15 | `reporte.html` |
| RN-26.10 El reporte se arma para un Sprint en cualquier estado; de otro proyecto es inexistente | CA-26.7, CA-26.8, CA-26.12 | `app.ArmarReporteSprint.Ejecutar` |
| RN-26.11 Aclaración en el reporte de un Sprint Cerrado (por qué el resumen de historias y el % completado pueden parecer contradictorios) | CA-26.14 | `pantallas_reporte.go` (`textoAclaracionReporteSprintCerrado`) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/reporte.go](../internal/domain/reporte.go) (`Reporte`, `FilaMetrica`, `ArmarReporteProyecto`, `ArmarReporteSprint`, `HistoriasDeSprint`, `DefectosDeSprint`) | [reporte_test.go](../internal/domain/reporte_test.go) |
| Caso de uso | [internal/app/armar_reporte.go](../internal/app/armar_reporte.go) (`ArmarReporteProyecto`, `ArmarReporteSprint`; reutilizan `LectorProyectos` HU-04, `ListadorElementos` HU-07, `RepositorioSprints`/`LectorSprintPorID` HU-08/HU-23, `LectorHoras` HU-23, `LectorDefectos` HU-21) | [armar_reporte_test.go](../internal/app/armar_reporte_test.go) |
| HTTP/Pantalla | [internal/adapters/http/pantallas_reporte.go](../internal/adapters/http/pantallas_reporte.go) (`GET /web/proyectos/{id}/reporte`, `GET /web/proyectos/{id}/sprints/{sprintID}/reporte`, sin endpoint JSON — P-1), [web/plantillas/reporte.html](../internal/adapters/http/web/plantillas/reporte.html); enlaces "Ver reporte" en [proyecto_detalle.html](../internal/adapters/http/web/plantillas/proyecto_detalle.html) (HU-04) y [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) (HU-09) | [pantallas_reporte_test.go](../internal/adapters/http/pantallas_reporte_test.go) |
| BDD | [features/HU-26-reporte-proyecto-sprint.feature](../features/HU-26-reporte-proyecto-sprint.feature), [hu26_steps_test.go](../test/bdd/hu26_steps_test.go), [hu26_sistema_test.go](../test/bdd/hu26_sistema_test.go) | 13 escenarios en verde |

### Commits

| Tarea | Commit | Etapa |
|---|---|---|
| 26.1 Spec | `ff2547d` docs(HU-26): spec SDD del reporte de proyecto y de sprint | — |
| 26.2 Escenarios BDD | `9b87d60` test(HU-26): [26.2] escenarios BDD del reporte de proyecto y de sprint | — |
| 26.4 Dominio (`Reporte`, `ArmarReporteProyecto`/`Sprint`, `HistoriasDeSprint`, `DefectosDeSprint`) | `e9e17e5` [RED] → `f1efeae` [GREEN] | **RED → GREEN** |
| 26.5 Casos de uso `ArmarReporteProyecto`/`ArmarReporteSprint` | `3c5a63c` [RED] → `8c6112a` [GREEN] | **RED → GREEN** |
| 26.3 Steps de godog conectados al caso de uso real | `634c5d0` test(HU-26): [26.3] steps de godog, conectados con el caso de uso real (BDD en verde, 13 escenarios) | — |
| 26.6 Adaptador de persistencia | No aplica: sin persistencia nueva (reutiliza los puertos de HU-04/HU-07/HU-08/HU-21/HU-23) | — |
| 26.7 Pantalla, conexión en el router/main y enlaces desde el detalle | `ae8d1ee` [RED] → `317dfdc` [GREEN] | **RED → GREEN** |
| 26.9 Revisión (`go-reviewer`) y trazabilidad | `573ac38` test(HU-26): cobertura del enlace "Ver reporte" en sprint_detalle.html (hallazgo #2 del `go-reviewer`) → `cdefe80` docs(HU-26): trazabilidad | — |
| Revisión de Matías (PR #378): aprobar spec, PDF pasa a HU-27 | `20477eb` docs(HU-26): revisión de Matías en el PR #378 - aprobar la spec | — |
| Revisión de Matías: aclaración del Sprint Cerrado (RN-26.11/CA-26.14 nuevas) | `90d6fd6` [RED] → `496b53c` [GREEN] | **RED → GREEN** |
| Revisión de Matías: trazas de CA-26.11 a CA-26.13 (decían "sección 7") | `d8e5310` test(HU-26): trazas de CA-26.11 y CA-26.12 en los tests del caso de uso (la pantalla se corrigió en el mismo commit que el RED de arriba, `90d6fd6`) | — |

### Huecos

- ~~La spec está en Borrador~~: resuelto el 2026-11-08. Matías López (Agile Enabler) la aprobó en la revisión del PR [#378](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/378): P-1 y P-3 confirmadas tal cual, P-2 se sacó porque ya no aplica. Ver el Historial de la spec para el detalle completo de la revisión.
- ~~CA-26.14 / RN-26.11 (conversión a `DocumentoPDF`) no se implementó en esta rama~~: resuelto el 2026-11-08, pero no implementándolo acá — Matías decidió en la revisión que la conversión a PDF **no es responsabilidad de HU-26**, es enteramente de HU-27. Se sacaron de la spec la RN-26.11/CA-26.14 viejas (PDF), la función `ADocumentoPDF` y los puertos `FuenteDocumento…` de la sección 5. En su lugar, esos mismos números se reusaron para una regla y un criterio nuevos: la aclaración del Sprint Cerrado (ver abajo). Lucas Orellana (`luas01`) sigue la conversión a PDF en su propia rama `hu-27-exportar-reporte-pdf`.
  - **Coordinación con Lucas:** mergeó esta rama (`hu-26-reporte-proyecto-sprint`) en la suya y compila; va a sumar un campo `URLPDF` a `datosReporte` y un botón "Descargar PDF" a `reporte.html` en su rama. Como HU-26 también tocó esos dos archivos para la aclaración del Sprint Cerrado, un futuro merge puede traer un conflicto de líneas agregadas en el mismo lugar (se resuelve dejando las dos, mismo patrón que advirtió Matías para los PR #375/#377). Avisarle a Lucas cuando se mergee el PR #378, como pidió.
- **Nueva RN-26.11 / CA-26.14 (aclaración del Sprint Cerrado):** agregada en la revisión del PR #378 (punto 1), con ciclo RED → GREEN (`90d6fd6` → `496b53c`). Con un Sprint cerrado con una historia Hecha (5 puntos) y otra sin terminar (8 puntos) que volvió al backlog al cerrarse, el reporte decía "1 de 1 historias completadas (5 de 5 puntos de historia)" y, en Métricas, "38,5 %" — Matías señaló que las dos cifras quedan juntas y parecen contradecirse, aunque es la misma regla que ya usa HU-23 (RN-23.8, opción A). La aclaración nueva lo explica junto a la tabla de Métricas, solo cuando el Sprint está `Cerrado`.
- **Numeración de RN en las secciones 2 y 3, corregida** (punto 3 de la revisión de Matías): citaban las reglas con un número de diferencia respecto de la sección 4 (por ejemplo, "cualquier estado" citaba RN-26.9 en vez de RN-26.10). Corregido en el mismo commit que aprueba la spec (`20477eb`).
- **Trazas de CA-26.11 a CA-26.13, corregidas** (punto 4 de la revisión): los tests de errores (`TestPantallaReporteProyecto_Errores`, `TestPantallaReporteSprint_Errores`, `TestArmarReporteProyecto_ProyectoInexistente_RetornaError`, `TestArmarReporteSprint_Errores`) decían `// Traza: HU-26 sección 7`; ahora citan el número de CA que cubre cada subcaso. CA-26.13 ya estaba justificada (sin escenario BDD) en el comentario del feature, que también se actualizó para CA-26.14.
- **Tarea 26.6 no aplica** (sin persistencia nueva), mismo caso que la 24.6 de HU-24.
- **Probado manualmente en el navegador** (servidor real, SQLite, 2026-11-07): proyecto con 2 historias (una Hecha con 5 SP y 10 h estimadas, otra Pendiente con 8 SP), un Sprint cerrado con la primera asignada, y un defecto Alta sobre esa historia (sin Sprint de detección, por quedar detectado después de cerrarse el único Sprint). Las dos pantallas (reporte del proyecto y del Sprint) mostraron las cinco secciones correctamente, con los formatos esperados ("-100,0 %", "38,5 %" con coma decimal) y el mensaje de "sin datos" en Defectos del reporte del Sprint.
- **Revisión con el agente `go-reviewer`:** Approve con cambios sugeridos, ninguno bloqueante. (1) *Medium/High, no resuelto, confirmado por Matías (punto 2 de su revisión):* `ArmarReporteProyecto.Ejecutar` reimplementa la orquestación de `ConsultarVelocidad` (HU-22) y `ConsultarCalidadSprints` (HU-24) en vez de llamarlos, como decía la spec original (sección 5); mismo resultado y bien testeado, pero duplica lógica entre `internal/app/armar_reporte.go` y `consultar_velocidad.go`/`consultar_calidad_sprints.go` — si cambia una regla en HU-22/HU-24 hay que acordarse de replicarla acá. Matías dio dos caminos, ninguno bloqueante: si el PR [#375](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/375) (HU-25, que agrega `calidadDe` justo para este cálculo) entra antes, HU-26 lo reutiliza; si no, queda como está. La spec (sección 5) ya se corrigió para describir lo que hace el código en vez de lo que decía antes (commit `20477eb`). (2) *Medium, resuelto:* faltaba test del enlace "Ver reporte" en `sprint_detalle.html` (sí lo había para `proyecto_detalle.html`); se agregó `TestPantallaSprintDetalle_EnlazaElReporte` (commit `573ac38`).
- **Aviso de conflictos (Matías, revisión del PR #378):** el PR #375 (HU-25) toca `main.go`, `pantallas.go`, `router.go`, `router_test.go`, `bdd_test.go` y la trazabilidad; el PR [#377](https://github.com/matiaslopez237/ingsoft-tpintegrador/pull/377) (HU-31) toca la trazabilidad y `bdd_test.go`. Son conflictos de líneas agregadas en los mismos lugares: el que se fusiona después deja las dos.

---

## HU-27: Exportar el reporte a PDF

- **Issue:** [#38](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/38), sub-issue de la épica 9 (Reportes) · Sprint 4 · 5 SP · depende de HU-26
- **Spec:** [docs/specs/HU-27-exportar-reporte-pdf.md](specs/HU-27-exportar-reporte-pdf.md) (**Borrador**; que la conversión sea de HU-27 lo confirmó Matías López en la issue #38 y en la revisión del PR #378; P-1 a P-11, abiertas a objeción en el PR)
- **ADR:** [docs/adr/0006-libreria-pdf.md](adr/0006-libreria-pdf.md) (`codeberg.org/go-pdf/fpdf` con DejaVu Sans Condensed embebida), todavía `proposed`
- **Feature:** [features/HU-27-exportar-reporte-pdf.feature](../features/HU-27-exportar-reporte-pdf.feature), 14 escenarios y Outlines, 38 escenarios al expandirlos (steps en [hu27_steps_test.go](../test/bdd/hu27_steps_test.go), con fuente y generador en memoria en [hu27_sistema_test.go](../test/bdd/hu27_sistema_test.go))
- **Rutas HTTP:** `GET /web/proyectos/{id}/reporte/pdf` y `GET /web/proyectos/{id}/sprints/{sprintID}/reporte/pdf` (descargas; los errores responden con la página de error HTML).
- **Pantalla:** enlace "Descargar PDF" en la pantalla del reporte de HU-26 (`reporte.html`), del proyecto y del sprint.
- **Integración con HU-26:** la rama `hu-26-reporte-proyecto-sprint` está mergeada en esta (`faea8c3`, y `aa7d69b` con la aclaración del sprint cerrado y la spec aprobada). La conversión del reporte a `DocumentoPDF` es de HU-27 (RN-27.16, CA-27.25) y vive en el adaptador HTTP (P-11).

### Criterios de aceptación

| CA | Regla | Escenario(s) BDD | Tests Go | Código |
|---|---|---|---|---|
| CA-27.1 Documento normal | RN-27.1, RN-27.2, RN-27.3, RN-27.10 | El documento exportado conserva título, subtítulo y secciones, y lleva la fecha de emisión `@normal` | `TestGenerar_DocumentoNormal_TextoEnOrdenEnUnaPagina` | `generador.go` (`Generar`, `dibujarPortada`, `textoEmision`) |
| CA-27.2 Acentos y ñ | RN-27.8, RN-27.9 | Outline "Los textos que el PDF puede mostrar no se reemplazan ni generan aviso", ejemplos de acentos, eñes y signos `@normal` | `TestNormalizarTextoPDF`, `TestGenerar_AcentosYEnie_TextoIdenticoSinAviso`, `TestGenerar_ListaDeRepresentables_SaleIdentica` | `documento_pdf_normalizar.go` (`representableEnPDF`), fuentes DejaVu embebidas en `generador.go` |
| CA-27.3 Tabla larga | RN-27.3, RN-27.4, RN-27.5 | **Sin escenario BDD a propósito:** hay que mirar el binario. Se cubre con tests del adaptador de PDF. | `TestGenerar_TablaLarga_CruzaPaginasSinPerderNiRepetirFilas`, `TestGenerar_EncabezadoYPie`, `TestGenerar_TablaLarga_RepiteLosEncabezadosEnCadaPagina` | `generador.go` (`dibujarTabla`, `encabezado`, `pie`) |
| CA-27.4 Párrafo largo | RN-27.3, RN-27.4 | **Sin escenario BDD a propósito** (adaptador de PDF) | `TestGenerar_ParrafoLargo_ContinuaEnOtrasPaginas` | `generador.go` (`escribir`, salto automático) |
| CA-27.5 Título de sección al pie | RN-27.4 | **Sin escenario BDD a propósito** (adaptador de PDF) | `TestGenerar_TituloDeSeccion_NoQuedaSoloAlPie` | `generador.go` (`dibujarSeccion`, `altoComienzo`) |
| CA-27.6 Largo de una celda | RN-27.6, RN-27.5 | **Sin escenario BDD a propósito** (el tope se ve en el documento normalizado y en el binario) | `TestNormalizarDocumento_TopeDeCaracteresPorCelda`, `TestNormalizarDocumento_TopeSoloEnLasCeldas`, `TestGenerar_TablaEnLosTopes_UnaFilaPorPagina` | `documento_pdf_normalizar.go` (`recortarCelda`), `generador.go` (`dibujarTabla`, `altoFila`) |
| CA-27.7 Vacíos | RN-27.7, RN-27.2, RN-27.11 | Un reporte sin subtítulo ni secciones se exporta igual `@limite` · Un párrafo vacío no impide exportar el reporte `@limite` | `TestDocumentoPDFValidar_ContenidoVacioPermitido_NoEsError`, `TestGenerar_SinSubtituloNiAviso_NoLosMuestra`, `TestGenerar_DocumentoMinimo`, `TestGenerar_SeccionesVacias` | `documento_pdf.go` (`Validar`), `generador.go` (`dibujarSeccion`, `dibujarTabla`: "Sin datos.") |
| CA-27.8 Emojis e ideogramas | RN-27.8, RN-27.9 | Outline "Los caracteres que el PDF no puede mostrar se reemplazan por "?" y se avisa", ejemplos de emojis e ideogramas `@alternativo` · Los caracteres se reemplazan en todos los textos del documento `@alternativo` | `TestNormalizarTextoPDF`, `TestNormalizarDocumento_TodosLosTextos`, `TestTextoAvisoReemplazos`, `TestExportarReporte_NormalizaYCompletaElAviso`, `TestExportarConGeneradorReal_Emojis_PDFConSignoYAviso` | `documento_pdf_normalizar.go` (`NormalizarDocumento`, `TextoAvisoReemplazos`), `exportar_reporte_pdf.go` (`generarPDF`) |
| CA-27.9 Secuencias de emojis | RN-27.8, RN-27.9 | Mismo Outline, ejemplos de secuencias de emojis `@limite` | `TestNormalizarTextoPDF`, `TestTextoAvisoReemplazos`, `TestExportarReporte_NormalizaYCompletaElAviso` | `documento_pdf_normalizar.go` (`invisibleDeFormato`) |
| CA-27.10 Espacios, invisibles y controles | RN-27.8, RN-27.9 | Outline de los textos que se pueden mostrar, ejemplos de espacios e invisibles `@limite` · Outline de los reemplazos, ejemplos de controles y bytes no UTF-8 `@limite` | `TestNormalizarTextoPDF` | `documento_pdf_normalizar.go` (`normalizarTextoPDF`) |
| CA-27.11 Documento inválido | RN-27.11, RN-27.6 (P-10) | Outline "Un documento mal armado no llega al generador", ejemplos de errores de armado `@error` y de títulos que quedan vacíos `@limite` | `TestDocumentoPDFValidar_DocumentoMalArmado_RetornaErrDocumentoInvalido`, `TestDocumentoPDFValidar_ContenidoVacioPermitido_NoEsError`, `TestExportarReporte_DocumentoInvalido_NoLlegaAlGenerador` | `documento_pdf.go` (`Validar`, `MaximoColumnas`, `LargoMaximoEncabezado`) |
| CA-27.12 Determinismo | RN-27.10 | **Sin escenario BDD a propósito:** compara bytes del PDF | `TestGenerar_Determinista`, `TestGenerar_LaFechaDelPDFEsLaEmision` | `generador.go` (`SetCreationDate`, `SetModificationDate`, `SetCatalogSort`) |
| CA-27.13 Fecha de emisión | RN-27.10, RN-27.12 | Outline "La fecha de emisión va en hora de Argentina y da la fecha del nombre del archivo", ejemplos del mismo día `@normal` y del cambio de día `@limite` | `TestHoraDeEmision` | `documento_pdf_emision.go` (`HoraDeEmision`, `zonaEmision`) |
| CA-27.14 Descarga del proyecto | Sección 3, RN-27.12 | Mismo Outline, ejemplos del mismo día (`@CA-27.14`): nombre del archivo | `TestNombreArchivoReporte`, `TestReportePDFProyecto_Descarga`, `TestReportePDF_PasaElContextoDelPedido`, `TestPantallaReporte_BotonDescargarPDF_LlevaALaDescarga` | `documento_pdf_emision.go` (`NombreArchivoReporteProyecto`), `pantallas_reporte_pdf.go` (`descargarReporteProyecto`, `responderPDF`) |
| CA-27.15 Descarga del sprint | Sección 3, RN-27.12 | El reporte de un sprint se exporta con su propio nombre de archivo `@alternativo` | `TestNombreArchivoReporte`, `TestExportarReporteSprint_GeneraElPDFDelDocumento`, `TestReportePDFSprint_Descarga`, `TestReportePDF_PasaElContextoDelPedido`, `TestPantallaReporte_BotonDescargarPDF_LlevaALaDescarga` | `documento_pdf_emision.go` (`NombreArchivoReporteSprint`), `exportar_reporte_pdf.go` (`ExportarReporteSprint`), `pantallas_reporte_pdf.go` (`descargarReporteSprint`) |
| CA-27.16 Nombre seguro | RN-27.12 (P-9) | El nombre del archivo no depende del nombre del proyecto `@limite` | `TestNombreArchivoReporte`, `TestReportePDF_NombreConCaracteresPeligrosos_NoInyectaCabeceras` (traza P-9) | `documento_pdf_emision.go`, `pantallas_reporte_pdf.go` (`nombreArchivoSeguro`) |
| CA-27.17 Proyecto inexistente | RN-27.13, sección 7 | Outline "No se exporta el reporte de un proyecto inexistente" `@error` | `TestExportarReporte_ProyectoInexistente`, `TestReportePDF_NoEncontrados_Devuelven404`, `TestFuenteDocumentoReporte_ErrorDelCasoDeUso_SinCambios` | `exportar_reporte_pdf.go`, `pantallas_reporte_pdf.go` (`errorDeReportePDF`, `erroresReportePDF`) |
| CA-27.18 Sprint inexistente o de otro proyecto | RN-27.13, sección 7 | Outline "No se exporta el reporte de un sprint inexistente o de otro proyecto" `@error` | `TestExportarReporteSprint_SprintInexistenteODeOtroProyecto`, `TestReportePDF_NoEncontrados_Devuelven404`, `TestFuenteDocumentoReporte_ErrorDelCasoDeUso_SinCambios` | `exportar_reporte_pdf.go` (`ProyectoID != proyectoID`), `pantallas_reporte_pdf.go` |
| CA-27.19 IDs inválidos | Sección 7 | **Sin escenario BDD a propósito:** los rechaza el adaptador HTTP sin llamar al caso de uso | `TestReportePDF_IDsInvalidos_Devuelven400SinLlamar` | `pantallas_reporte_pdf.go` (reutiliza `idDeRuta` y `sprintIDDesdeRuta`) |
| CA-27.20 Falla interna | Sección 7 | Outline "Si falla el armado del reporte o el generador, no se obtiene ningún archivo" `@error` (parte del caso de uso; el 500 sin detalle, en tests HTTP) | `TestExportarReporte_FallaDeLaFuenteODelGenerador_SeEnvuelveSinArchivo`, `TestReportePDF_FallaInterna_Devuelve500SinDetalle`, `TestReportePDF_DocumentoInvalido_Devuelve500SinDetalle`, `TestFuenteDocumentoReporte_ErrorDelCasoDeUso_SinCambios` | `exportar_reporte_pdf.go` (`%w`), `pantallas_reporte_pdf.go` (`fallaInterna`) |
| CA-27.21 Botón y revisión visual | Sección 3 | **Sin escenario BDD:** la parte automatizable es la pantalla (tests del adaptador); la revisión visual es manual y no tiene escenario `@manual` (ver "Huecos") | `TestPantallaReporte_BotonDescargarPDF`, `TestPantallaReporte_BotonDescargarPDF_LlevaALaDescarga` | `pantallas_reporte.go` (campo `URLPDF` en `aDatosReporte`), `reporte.html` |
| CA-27.22 Mismas secciones, en el exportador | RN-27.1, RN-27.2, RN-27.7 | El PDF conserva los títulos y el orden de las secciones del reporte, aunque estén vacías `@normal` | `TestGenerar_SeccionesVacias` | `generador.go` (`dibujarSeccion`) |
| CA-27.23 Mismas secciones que la pantalla, integración | RN-27.1, RN-27.2, RN-27.16 | **Sin escenario BDD** (ver "Huecos") | `TestExportarReporte_TitulosDelPDFSonLosDeLaPantalla` (fuente, caso de uso y generador reales; dobles de los casos de uso de HU-26) | `reporte_documento.go`, `exportar_reporte_pdf.go`, `generador.go`, `reporte.html` |
| CA-27.24 PDF válido | RN-27.15 | **Sin escenario BDD a propósito:** mira el binario y las cabeceras | `TestGenerar_EsUnPDFCompleto`, `TestReportePDFProyecto_Descarga`, `TestReportePDFSprint_Descarga` | `generador.go` (`Output` a un buffer), `pantallas_reporte_pdf.go` (`responderPDF`, `Content-Length`) |
| CA-27.25 Conversión del reporte | RN-27.16, RN-27.11 (P-10, P-11) | **Sin escenario BDD a propósito:** compara el documento con la pantalla, en el adaptador HTTP | `TestADocumentoPDF_SeccionesSonLasDeLaPantalla`, `TestADocumentoPDF_TextosSonLosDeLaPantalla`, `TestADocumentoPDF_ReporteDeProyecto`, `TestADocumentoPDF_ReporteDeSprint`, `TestADocumentoPDF_ReporteVacio_MensajesDeLaPantalla`, `TestADocumentoPDF_SiempreValido`, `TestFuenteDocumentoReporte_DocumentoDeProyecto`, `TestFuenteDocumentoReporte_DocumentoDeSprint` | `reporte_documento.go` (`aDocumentoPDF`, `fuenteDocumentoReporte`, `NuevaFuenteDocumentoReporte`) |
| CA-27.26 Aclaración del sprint cerrado | RN-27.17, RN-27.16 | **Sin escenario BDD a propósito:** como CA-27.25, compara el documento con la pantalla, en el adaptador HTTP | `TestADocumentoPDF_ReporteDeSprint` (subtest "sprint cerrado"), `TestADocumentoPDF_TextosSonLosDeLaPantalla` (casos "sprint cerrado" y "sprint cerrado vacío"), `TestADocumentoPDF_SinAclaracion_SprintAbiertoYProyecto` | `reporte_documento.go` (`bloquesMetricas` copia `datosReporte.AclaracionSprintCerrado`, que arma `aDatosReporte` de HU-26) |

Los tests se listan según su etiqueta `Traza`. Los escenarios de CA-27.2, CA-27.8 a CA-27.11, CA-27.13, CA-27.14, CA-27.17, CA-27.18 y CA-27.20 son ejemplos de `Scenario Outline`, con el tag en cada bloque `Examples`.

### Reglas de negocio

| RN | CA que la cubren | Tests con la etiqueta (ejemplos) | Código |
|---|---|---|---|
| RN-27.1 Modelo neutral | CA-27.1, CA-27.22, CA-27.23 | `TestGenerar_DocumentoNormal_TextoEnOrdenEnUnaPagina`, `TestExportarReporte_TitulosDelPDFSonLosDeLaPantalla` | `documento_pdf.go` (`DocumentoPDF`), puertos de `exportar_reporte_pdf.go` |
| RN-27.2 Orden del contenido | CA-27.1, CA-27.7, CA-27.22, CA-27.23 | `TestGenerar_DocumentoNormal_TextoEnOrdenEnUnaPagina`, `TestGenerar_SinSubtituloNiAviso_NoLosMuestra` | `generador.go` (`dibujarPortada`) |
| RN-27.3 Encabezado y pie | CA-27.1, CA-27.3, CA-27.4 | `TestGenerar_EncabezadoYPie` | `generador.go` (`encabezado`, `pie`, `AliasNbPages`) |
| RN-27.4 Saltos de página | CA-27.3, CA-27.4, CA-27.5 | `TestGenerar_TituloDeSeccion_NoQuedaSoloAlPie`, `TestGenerar_TablaDespuesDeUnParrafo_EncabezadosJuntoALaPrimeraFila` | `generador.go` (`dibujarSeccion`, `altoComienzo`) |
| RN-27.5 Tablas | CA-27.3, CA-27.6 | `TestGenerar_CeldaLarga_SeAjustaEnVariasLineas`, `TestGenerar_FilaDeVariasLineas_NoSeParte`, `TestGenerar_TablaLarga_RepiteLosEncabezadosEnCadaPagina`, `TestGenerar_TablaEnLosTopes_UnaFilaPorPagina` | `generador.go` (`dibujarTabla`, `dibujarFila`, `altoFila`) |
| RN-27.6 Largo de una celda; topes de columnas y encabezados (P-10) | CA-27.6, CA-27.11 | `TestNormalizarDocumento_TopeDeCaracteresPorCelda`, `TestDocumentoPDFValidar_DocumentoMalArmado_RetornaErrDocumentoInvalido`, `TestGenerar_TablaEnLosTopes_UnaFilaPorPagina` | `documento_pdf_normalizar.go` (`recortarCelda`), `documento_pdf.go` (`MaximoColumnas`, `LargoMaximoEncabezado`) |
| RN-27.7 Tabla sin filas | CA-27.7, CA-27.22 | `TestGenerar_SeccionesVacias` | `generador.go` (`dibujarTabla`) |
| RN-27.8 Caracteres | CA-27.2, CA-27.8, CA-27.9, CA-27.10 | `TestNormalizarTextoPDF`, `TestGenerar_ListaDeRepresentables_SaleIdentica` | `documento_pdf_normalizar.go` |
| RN-27.9 Aviso de reemplazos | CA-27.8, CA-27.9, CA-27.10 | `TestTextoAvisoReemplazos`, `TestExportarReporte_SinReemplazos_AvisoVacioAunqueLaFuenteTraigaUno` | `documento_pdf_normalizar.go` (`TextoAvisoReemplazos`), `exportar_reporte_pdf.go` |
| RN-27.10 Fecha de emisión | CA-27.1, CA-27.12, CA-27.13 | `TestHoraDeEmision`, `TestGenerar_LaFechaDelPDFEsLaEmision`, `TestExportarReporteProyecto_GeneraElPDFDelDocumento` | `documento_pdf_emision.go` (`HoraDeEmision`), `generador.go` |
| RN-27.11 Documento válido | CA-27.7, CA-27.11, CA-27.25 | `TestDocumentoPDFValidar_DocumentoMalArmado_RetornaErrDocumentoInvalido`, `TestADocumentoPDF_SiempreValido` | `documento_pdf.go` (`Validar`) |
| RN-27.12 Nombre del archivo | CA-27.13 a CA-27.16 | `TestNombreArchivoReporte`, `TestExportarReporteProyecto_GeneraElPDFDelDocumento` | `documento_pdf_emision.go` |
| RN-27.13 Qué se puede descargar | CA-27.17, CA-27.18 | `TestExportarReporteSprint_SprintInexistenteODeOtroProyecto`, `TestReportePDF_NoEncontrados_Devuelven404` | `exportar_reporte_pdf.go` |
| RN-27.14 Hoja A4 vertical | CA-27.1 | `TestGenerar_DocumentoNormal_TextoEnOrdenEnUnaPagina` | `generador.go` (`fpdf.New("P", "mm", "A4", "")`) |
| RN-27.15 Archivo válido | CA-27.24 | `TestGenerar_EsUnPDFCompleto`, `TestReportePDFProyecto_Descarga` | `generador.go`, `pantallas_reporte_pdf.go` |
| RN-27.16 Conversión del reporte | CA-27.23, CA-27.25 | `TestADocumentoPDF_SeccionesSonLasDeLaPantalla`, `TestADocumentoPDF_TextosSonLosDeLaPantalla`, `TestExportarReporte_TitulosDelPDFSonLosDeLaPantalla` | `reporte_documento.go` |
| RN-27.17 Aclaración del sprint cerrado | CA-27.26 | `TestADocumentoPDF_ReporteDeSprint`, `TestADocumentoPDF_SinAclaracion_SprintAbiertoYProyecto` | `reporte_documento.go` (`bloquesMetricas`) |

### Tests sin CA propio

Llevan `Traza: HU-27 (sección N)`, una decisión o una RN sin CA. `grep "CA-27\."` no los encuentra; `grep "HU-27"`, sí:

| Test | Traza | Qué verifica |
|---|---|---|
| `TestNormalizarDocumento_NoModificaElOriginal`, `TestNormalizarDocumento_TopeNoModificaElOriginal` | Sección 5 | La normalización no modifica el documento que recibe |
| `TestNormalizarDocumento_ConservaEmisionYAviso` | Sección 5 | Deja `Emision` y `Aviso` como están |
| `TestExportarReporteProyecto_GeneraElPDFDelDocumento` | RN-27.10, RN-27.12 (sección 5) | Orden de los pasos y nombre del archivo, con un generador falso |
| `TestExportarReporte_SinReemplazos_AvisoVacioAunqueLaFuenteTraigaUno` | RN-27.9 | El aviso que traiga la fuente se pisa |
| `TestGenerar_FallaDeFpdf_DevuelveError` | Sección 7 | El error de `Output` sale tal cual, sin bytes |
| `TestGenerar_CeldaLarga_SeAjustaEnVariasLineas`, `TestGenerar_FilaDeVariasLineas_NoSeParte`, `TestGenerar_TablaDespuesDeUnParrafo_EncabezadosJuntoALaPrimeraFila` | RN-27.4, RN-27.5 | Tablas con celdas de varias líneas y al pie |
| `TestReportePDF_ConviveConLasPantallas` | Sección 5 | Las rutas de descarga se registran junto a las pantallas sin choque de patrones |
| `TestReportePDF_NombreConCaracteresPeligrosos_NoInyectaCabeceras` | P-9 (RN-27.12) | El handler sanea el nombre antes de la cabecera |

Sin etiqueta `Traza` de HU-27, porque es un test compartido que suma casos: `TestNuevoRouter_RegistraTodasLasRutas` (suma las dos rutas `/pdf`).

### Errores nuevos

`domain.ErrDocumentoInvalido` (RN-27.11): es un error de programación de quien arma el documento y se responde con `500`. Se reutilizan `ErrProyectoNoEncontrado` (HU-01) y `ErrSprintNoEncontrado` (HU-09), los dos con `404`.

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/documento_pdf.go](../internal/domain/documento_pdf.go) (`DocumentoPDF`, `SeccionPDF`, `BloquePDF`, `TablaPDF`, `Validar`, `ErrDocumentoInvalido`), [documento_pdf_normalizar.go](../internal/domain/documento_pdf_normalizar.go) (`NormalizarDocumento`, `TextoAvisoReemplazos`), [documento_pdf_emision.go](../internal/domain/documento_pdf_emision.go) (`HoraDeEmision`, `NombreArchivoReporteProyecto`, `NombreArchivoReporteSprint`) | [documento_pdf_test.go](../internal/domain/documento_pdf_test.go), [documento_pdf_normalizar_test.go](../internal/domain/documento_pdf_normalizar_test.go), [documento_pdf_emision_test.go](../internal/domain/documento_pdf_emision_test.go) |
| Casos de uso y puertos | [internal/app/exportar_reporte_pdf.go](../internal/app/exportar_reporte_pdf.go) (`ExportarReporteProyecto`, `ExportarReporteSprint`, puertos `FuenteDocumentoProyecto`, `FuenteDocumentoSprint` y `GeneradorPDF`; reutiliza `LectorProyectos` de HU-04 y `LectorSprintPorID` de HU-23) | [exportar_reporte_pdf_test.go](../internal/app/exportar_reporte_pdf_test.go) |
| Adaptador de PDF | [internal/adapters/pdf/generador.go](../internal/adapters/pdf/generador.go) (`Generador`, `NuevoGenerador`, `Generar`), [fuentes/](../internal/adapters/pdf/fuentes) (DejaVu Sans Condensed 2.37, licencia en [docs/third-party-licenses/](third-party-licenses/dejavu-fonts-LICENSE.txt)) | [generador_test.go](../internal/adapters/pdf/generador_test.go), [integracion_test.go](../internal/adapters/pdf/integracion_test.go) |
| Conversión del reporte (P-11) | [internal/adapters/http/reporte_documento.go](../internal/adapters/http/reporte_documento.go) (`aDocumentoPDF`, `fuenteDocumentoReporte`, `NuevaFuenteDocumentoReporte`) | [reporte_documento_test.go](../internal/adapters/http/reporte_documento_test.go), [reporte_documento_fuente_test.go](../internal/adapters/http/reporte_documento_fuente_test.go), [reporte_pdf_integracion_test.go](../internal/adapters/http/reporte_pdf_integracion_test.go) |
| Descargas HTTP y botón | [internal/adapters/http/pantallas_reporte_pdf.go](../internal/adapters/http/pantallas_reporte_pdf.go) (`RegistrarExportarReportePDF`, `responderPDF`, `nombreArchivoSeguro`), [pantallas_reporte.go](../internal/adapters/http/pantallas_reporte.go) (`URLPDF`), [reporte.html](../internal/adapters/http/web/plantillas/reporte.html), [router.go](../internal/adapters/http/router.go), [cmd/server/main.go](../cmd/server/main.go) | [pantallas_reporte_pdf_test.go](../internal/adapters/http/pantallas_reporte_pdf_test.go), [pantallas_reporte_boton_test.go](../internal/adapters/http/pantallas_reporte_boton_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
| BDD | [features/HU-27-exportar-reporte-pdf.feature](../features/HU-27-exportar-reporte-pdf.feature), [test/bdd/hu27_steps_test.go](../test/bdd/hu27_steps_test.go), [hu27_sistema_test.go](../test/bdd/hu27_sistema_test.go) (fuente y generador en memoria), [bdd_test.go](../test/bdd/bdd_test.go) | 38 escenarios y 252 pasos en verde |

### Commits

Salen de `git log --no-merges origin/main..HEAD ^origin/hu-26-reporte-proyecto-sprint` (72 commits hasta `f68e9a6`, más los de documentación posteriores), en orden. Los merges están al final.

| Ciclo | Commit | Etapa |
|---|---|---|
| ADR | `757de07` docs(HU-27): ADR de la librería de PDF → `8e15310` renumerado a 0006 | — |
| Spec | `ae10399` docs(HU-27): spec SDD de exportar el reporte a PDF → `23999df` ligada a los criterios de la issue #38 | — |
| BDD en rojo: feature | `3849991` test(HU-27): [RED] escenarios BDD de exportar el reporte a PDF | BDD en rojo |
| BDD en rojo: steps | `55d0145` test(bdd): [RED] steps de godog de HU-27 | BDD en rojo, por compilación |
| Spec: `Aviso` y bloque vacío | `90d2002` docs(HU-27) → `ebcba02` test(HU-27): [RED] escenario del bloque vacío | BDD en rojo, por compilación |
| Dominio: validar | `124e49b` [RED] → `2ed5c99` [GREEN] | **RED → GREEN** |
| Dominio: normalizar y avisar | `c6ef36b` [RED] → `9e84213` [GREEN] | **RED → GREEN** |
| Dominio: emisión y nombre | `296c022` [RED] → `8e03f08` [GREEN] | **RED → GREEN** |
| Casos de uso | `d50a616` [RED] → `7a6b1b4` test(bdd): [RED] lectores en el sistema en memoria → `3b3bab5` [GREEN] → `4e319da` [REFACTOR] | **RED → GREEN → REFACTOR**; el BDD de HU-27 en verde desde `3b3bab5` |
| Spec: puertos | `1bf3562` docs(HU-27) | — |
| Tope de 300 por celda | `49d747d` [RED] → `4b1996b` [GREEN] → `9503984` docs (P-8) | **RED → GREEN** |
| Adaptador de PDF: texto en orden | `7bffdcb` [RED] → `8ad30bc` [GREEN] | **RED → GREEN** |
| Adaptador de PDF: páginas, tablas y vacíos | `3926404` [RED] → `f046db4` [GREEN] | **RED → GREEN** |
| Adaptador de PDF: encabezado, pie, títulos al pie, determinismo | `a9b9d2d` [RED] → `685bc65` [GREEN] → `3e52382` [REFACTOR] | **RED → GREEN → REFACTOR** |
| Adaptador de PDF: mutantes sobrevivientes | `9221050` fecha de creación · `f6a58c7` encabezados al pie | — (sin RED real, ver abajo) |
| Descargas HTTP | `fbaf73a` [RED] → `a4abece` [GREEN] → `f973e77` docs (registro propio) → `c9075e5` docs (P-9) | **RED → GREEN** |
| Revisión: topes (P-10) | `311e1d9` docs → `ddd97a7` [RED] → `dd79d27` [GREEN] | **RED → GREEN** |
| Revisión: tests que faltaban | `d3c82ea` representables · `d35b1be` acentos · `3a1d3fe` emojis de punta a punta · `e69d7b9` documento inválido en HTTP · `df58232` docs → `855d1fe` selectores suplementarios · `5f7fde1` trazas | — (sin RED real, ver abajo) |
| Revisión: refactor y preguntas | `8028407` refactor(HU-27): un solo "generar el PDF" · `d241bf1` docs (preguntas 12 y 13) | — (refactor sin cambio de comportamiento) |
| Integración con HU-26 | `faea8c3` merge de `origin/hu-26-reporte-proyecto-sprint` → `7a1ea6f` docs (P-11) | — |
| Conversión del reporte | `0789297` [RED] → `57a1666` [GREEN] | **RED → GREEN** |
| Fuente del documento | `d7855e7` [RED] → `4496a55` [GREEN] | **RED → GREEN** |
| Cambio de alcance (#378) | `4d9d784` revert de la spec de HU-26 · `39b14ab` docs (RN-27.16, CA-27.25) · `c0dbcd1` docs (comentarios) · `720d796` test (fixture propio y trazas) · `ffeae9b` docs (pregunta 9) | — |
| Botón | `165abfa` [RED] → `400fc79` [GREEN] | **RED → GREEN** |
| Rutas en `NuevoRouter` | `a2f5767` [RED] → `3fe8d8f` [GREEN] → `0cb61f0` conexión en `main.go` | **RED → GREEN**; `main.go` sin test automático (humo, ver abajo) |
| Integración CA-27.23 | `bac4b4f` test(HU-27): títulos de sección del PDF iguales a los de la pantalla | — (sin RED real, ver abajo) |
| Spec: implementado | `773d548` docs(HU-27) | — |
| Revisión: comentarios de código | `400426e` docs(HU-27): comentarios de código con Implementa HU-27 y sin referencias viejas. Resuelve dos huecos de la primera versión de esta trazabilidad: el comentario de `DocumentoPDF` que decía que HU-26 completa el documento, y `Implementa HU-27`, que antes estaba solo en los dos handlers de descarga | — (solo comentarios) |
| Aclaración del sprint cerrado (RN-27.17) | `9baf345` [RED] → `f68e9a6` [GREEN] | **RED → GREEN** |

**Merges de la rama, que no son de HU-27:** `f3a0005` y `862daa2` (traen `origin/main`), y `faea8c3` y `aa7d69b` (traen HU-26).

**Commits que el filtro `git log --grep="(HU-27)"` no encuentra**, porque usan el prefijo `test(bdd):`: `55d0145` (RED de los steps) y `7a6b1b4` (lectores en el sistema en memoria).

### Evidencia: mutaciones manuales

Ninguna quedó en un commit: se aplicaron en una copia del repo o en el árbol, y se restauraron. En total, 109 mutantes: 108 murieron y 1 sobrevivió, equivalente para el test que lo corre. Los resultados de las tandas anteriores a la integración con HU-26 (dominio y casos de uso, tope de celda, adaptador de PDF, descargas HTTP y revisión) salen de los registros de las conversaciones donde se hicieron, porque los cuerpos de los commits están vacíos.

| Etapa | Mutantes | Murieron | Sobrevivieron y por qué |
|---|---|---|---|
| Dominio y casos de uso | 10 | 10 | — |
| Tope de 300 por celda | 8 | 8 | — |
| Adaptador de PDF | 13 | 13 | M8 (sin la fecha de emisión en el PDF) y M11 (tabla al pie sin su primera fila) sobrevivieron primero; los matan `9221050` y `f6a58c7`. **M5** (fuente sin glifos para la ñ) muere, pero se aproximó con Helvetica: es **casi equivalente**, porque no se encontró una TTF sin la ñ para probar el caso exacto |
| Descargas HTTP | 14 | 14 | — |
| Revisión: topes y tests que faltaban | 14 | 14 | — (topes: 7; lista de representables, normalizar en el caso de uso, documento inválido como 404: 3; rango de los selectores suplementarios: 4) |
| Conversión y fuente | 25 | 25 | — (14 de `aDocumentoPDF`, 10 de la fuente, repetidos con el fixture propio, y 1 sobre el "—" de `aVistasDefectosReporte`, aplicado solo para comprobar el test) |
| Botón | 10 | 10 | — |
| Rutas en `NuevoRouter` | 3 | 3 | — |
| Integración CA-27.23 | 6 | 5 | "La fuente del sprint usa el caso de uso del proyecto": **equivalente para ese test**, porque los dos dobles devuelven el mismo reporte y los títulos son los mismos. Lo mata `TestFuenteDocumentoReporte_DocumentoDeSprint` |
| Aclaración del sprint cerrado | 6 | 6 | — (quitar la condición, invertirla, otro texto, quitar el párrafo, como tabla en vez de párrafo, y condición por `EsDeSprint` en vez de la aclaración) |

Algunos mutantes no compilaron en el primer intento (un import o una variable sin usar) y se repitieron dejándolos en uso: `Content-Length` sin `strconv`, el encabezado contado en bytes, y la fuente que devuelve un documento vacío. En la aclaración del sprint cerrado, el mutante "como tabla" no compiló la primera vez por un error del script (`\&`) y se repitió.

**Cruce de pasos BDD al integrar HU-26:** un programa temporal, que no está en el repo, cruzó las 475 regex de `test/bdd` (incluidas las armadas con la constante `horas`) con las 1851 líneas de paso de todos los `.feature`, con los Outline expandidos: 0 pasos con más de una regex y 0 sin regex. Se repitió al integrar la aclaración de HU-26 (`aa7d69b`), con el mismo resultado.

**Prueba de humo:** `cmd/server` con una base temporal y `curl`. La pantalla del reporte responde `200` con el botón; la descarga, `200` con `%PDF-1.3` … `%%EOF`, `filename="reporte-proyecto-1-2026-11-07.pdf"` y `nosniff`; un proyecto inexistente, `404`. No es una prueba en navegador.

### Etapas sin RED real

- **Tests agregados con el código ya correcto** (sus mensajes no llevan etiqueta `[RED]`):
  - `9221050` y `f6a58c7`: los agregaron los mutantes M8 y M11 del adaptador de PDF, que sobrevivían; el código ya hacía lo correcto
  - `d3c82ea`, `d35b1be`, `e69d7b9` y `855d1fe`: tests de la revisión con `go-reviewer`, sobre comportamiento que ya existía
  - `3a1d3fe`: emojis de punta a punta, con el caso de uso y el generador reales
  - `bac4b4f`: CA-27.23 con el generador real; prueba código que ya existía
  - `720d796`: fixture propio y trazas; cambia el esperado de la columna Elemento, que con el fixture anterior no se podía probar
  - `5f7fde1`: solo trazas y comentarios
- **RED que ya pasaba en parte:** en `ddd97a7`, `TestGenerar_TablaEnLosTopes_UnaFilaPorPagina` pasó en el RED, porque la regla nueva vive en `Validar` y el dibujo de la tabla no cambia; el RED falló por los dos casos inválidos del dominio.
- **RED por compilación:** `55d0145`, `ebcba02`, `124e49b`, `c6ef36b`, `296c022`, `d50a616`, `7a6b1b4`, `7bffdcb`, `fbaf73a`, `0789297`, `d7855e7` y `a2f5767` (falta la función, el tipo o el campo). `165abfa` falló por aserción (la pantalla todavía no tenía el botón), y `9baf345` también (el documento de un sprint cerrado no tenía la aclaración; `TestADocumentoPDF_SinAclaracion_SprintAbiertoYProyecto` ya pasaba en el RED, porque es la guarda del caso negativo).
- **GREEN que corrigió el test del RED:** `57a1666` cambió el esperado de la columna Elemento ("—" por "Logout"), porque el fixture de HU-26 no tenía IDs. `720d796` volvió a "—" con un fixture propio. `f046db4` y `685bc65` ajustaron tests del RED anterior, como explican sus mensajes.
- **Sin test automático:** `0cb61f0` (`main.go`), cubierto con la prueba de humo.

### Huecos

- **CA-27.23 no tiene escenario BDD.** Hacerlo sobre SQLite real acoplaría los steps de HU-27 a los tipos internos de `hu26_sistema_test.go` (código de Eze) o copiaría unas 160 líneas de esa preparación, y sumaría la lectura del PDF a `test/bdd`. Lo que agregaría es poco: el paso de SQLite al reporte lo cubre el BDD de HU-26 sobre SQLite real, y el del reporte al PDF lo cubre `TestExportarReporte_TitulosDelPDFSonLosDeLaPantalla`. El recorrido completo con el servidor real se probó a mano (prueba de humo).
- **La parte manual de CA-27.21 no tiene escenario `@manual`.** La suite no filtra `@manual` y corre con `Strict`, así que un escenario manual fallaría por pasos indefinidos. Checklist para la revisión visual, con un PDF descargado de la pantalla:
  - [ ] Los acentos, la ñ, "¿", "¡", "€", "–" y "…" se ven bien.
  - [ ] Desde la página 2 está el encabezado con el título del documento; en todas, el pie "Página k de N", con la N correcta.
  - [ ] Ninguna fila de tabla está encimada ni cortada entre dos páginas, y los encabezados se repiten donde la tabla sigue.
  - [ ] Ningún título de sección queda solo al pie.
  - [ ] Con un emoji en un nombre, se ve "?" y el aviso después de la línea de emisión.
- **Mutantes equivalentes o casi:** el de CA-27.23 sobre la fuente del sprint (equivalente para ese test, lo mata el test de la fuente) y M5 del adaptador de PDF (casi equivalente, aproximado con Helvetica).
- **`ledongthuc/pdf` es una dependencia solo de los tests:** `go list -deps ./cmd/server` no la incluye. Está en el bloque `require` de `go.mod`, junto a `fpdf`.
- **La spec está en Borrador**, y el **ADR-0006 sigue en `proposed`** (pregunta abierta 10).
- **`-race` no se corrió localmente:** el único gcc de la máquina es MinGW de 32 bits. Lo corre el CI en Linux (`go test -race ./...`).
- **BDD completo: unos 500 s** (499,7 s en la última corrida). Necesita `-timeout 30m` localmente; en el CI hace falta el timeout de 20 minutos de la rama `ci-timeout-tests` (`5a96c12`), porque `go test` corta cada paquete a los 10 minutos y el CI corre además con `-race`.
- **Decisiones del autor, abiertas a objeción en el PR:** P-8 ("?" que se pierden al recortar una celda, sin test propio), P-9 (saneado del nombre en el handler), P-10 (6 columnas y 100 caracteres por encabezado) y la parte de P-11 que dice que la conversión vive en el adaptador HTTP y no en el dominio. La otra parte de P-11, que la conversión es de HU-27, la confirmó Matías López (issue #38 y revisión del PR #378).
- ~~**RN-27.17 y CA-27.26 (aclaración del sprint cerrado) están pendientes de HU-26.**~~ Resuelto el 2026-11-08: HU-26 agregó `AclaracionSprintCerrado` a `datosReporte` y se integró en esta rama; `aDocumentoPDF` la copia como primer párrafo de Métricas (`9baf345` [RED] → `f68e9a6` [GREEN]). Texto anterior: Matías López le pidió a Eze que el reporte de un sprint cerrado lleve una aclaración, como la de HU-23. Mientras HU-26 no la agregue a `datosReporte`, la pantalla no la muestra y el PDF tampoco, así que no hay nada que probar. Cuando exista: sumarla a `aDocumentoPDF` (como un párrafo en la misma sección que la pantalla), y a `TestADocumentoPDF_ReporteDeSprint` y `TestADocumentoPDF_TextosSonLosDeLaPantalla` con un sprint cerrado, con su RED y su GREEN.
- **Preguntas abiertas:** 12 (recursos de la generación: sin tope de tamaño ni de descargas simultáneas), 13 (texto descompuesto: `e` + `U+0301` sale "e?") y 14 (sprints abiertos y cerrados en el Resumen, que la pantalla no muestra; decisión de Eze). La 15 (aclaración del reporte de un sprint cerrado) quedó resuelta: pasó a RN-27.17 y CA-27.26.
- **Textos fijos repetidos:** los mensajes de "sin datos", "Sin fecha de fin" y "—" están en `reporte.html` y en `reporte_documento.go`. `TestADocumentoPDF_TextosSonLosDeLaPantalla` los compara como subcadenas, así que un recorte del texto ("Sin fecha" en vez de "Sin fecha de fin") lo detecta solo el test del documento esperado completo. La aclaración del sprint cerrado no se repite: se copia del campo de `datosReporte`, y los párrafos de la pantalla se comparan enteros.

### Para revisar

Código de otras historias que tocó HU-27 (`git diff --stat origin/hu-26-reporte-proyecto-sprint...HEAD`, sin los archivos propios de HU-27). Conviene que lo revisen sus autores en la PR.

| Archivo | Historia | Cambio |
|---|---|---|
| [internal/adapters/http/pantallas_reporte.go](../internal/adapters/http/pantallas_reporte.go) | HU-26 (Eze) | **Líneas de Eze tocadas:** campo `URLPDF` en `datosReporte`, y dos líneas en `aDatosReporte` que lo arman (del proyecto, y del sprint con su ID); el import de `fmt` |
| [reporte.html](../internal/adapters/http/web/plantillas/reporte.html) | HU-26 (Eze) | **Línea de Eze tocada:** el enlace "Descargar PDF" dentro de `<div class="botones">`, después de "Volver…" |
| [cmd/server/main.go](../cmd/server/main.go) | Compartido | **Líneas de Eze tocadas:** los dos casos de uso del reporte de HU-26 pasan a variables (`armarReporteProyecto`, `armarReporteSprint`), con la misma construcción, porque los usan la pantalla y la fuente. Se suman el import de `adapters/pdf`, la fuente, el generador y los dos casos de uso de exportación |
| [router.go](../internal/adapters/http/router.go), [router_test.go](../internal/adapters/http/router_test.go) | Compartidos | Dos campos en `CasosDeUso`, el registro antes de `RegistrarPantallas`, dos dobles y dos rutas |
| [test/bdd/bdd_test.go](../test/bdd/bdd_test.go) | Compartido | Registra `registrarStepsHU27`; en el merge con HU-26 se conservaron las dos líneas |
| [go.mod](../go.mod), [go.sum](../go.sum) | Compartidos | `codeberg.org/go-pdf/fpdf` y `github.com/ledongthuc/pdf` (solo tests) |
| [docs/third-party-licenses/README.md](third-party-licenses/README.md) | Compartido | Sección "Fuentes" con la licencia de DejaVu |
| [docs/specs/HU-26-reporte-proyecto-sprint.md](specs/HU-26-reporte-proyecto-sprint.md) | HU-26 (Eze) | Sin cambios: `4d9d784` la dejó igual a la de su rama |

---

## Defectos

| Defecto | Historia | Issue | Escenario BDD | Tests Go | Corrección | Commits |
|---|---|---|---|---|---|---|
| DEF-01 El formulario de edición de proyecto pisa cambios simultáneos | HU-30 (RN-30.3, CA-30.5) | [#179](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/179) | Dos ediciones simultáneas de campos distintos no se pisan `@DEF-01` | `TestPantallas_EditarProyecto_FormularioLlevaLosValoresOriginales`, `TestPantallas_EditarProyecto_SoloEnviaLosCamposCambiados`, `TestPantallas_EditarProyecto_SinCambios_VuelveAlDetalleSinError` | `cambiosDeForm` y `leerOriginalProyecto` en [pantallas.go](../internal/adapters/http/pantallas.go); campos ocultos `original_*` en [proyecto_form.html](../internal/adapters/http/web/plantillas/proyecto_form.html) | `1bbf72a` [RED] → `49c6303` [GREEN] → `7af3f6e` (spec y trazabilidad); sin cambios vuelve al detalle: `4a95342` [RED] → `73439b0` [GREEN] → este commit |
| DEF-02 Eliminar un elemento con defectos informaba "con esfuerzo registrado" | HU-06 (RN-06.10, CA-06.11), al integrarse con HU-20 | [#205](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/205) | No se elimina un elemento con un defecto registrado `@DEF-02` | `TestEliminarElemento_ConDefectos_NoLoElimina`, `TestEliminarElemento_ConEsfuerzoYDefectos_InformaElEsfuerzo`, `TestRepositorioBacklog_ElementoConDefectos`, `TestDeleteElemento`, `TestPantallaElemento_Errores` | `ErrElementoConDefectos` en [backlog_editar.go](../internal/domain/backlog_editar.go); `TieneDefectos` y `errorDeEliminacion` en [backlog_editar.go](../internal/adapters/persistencia/backlog_editar.go); verificación en [editar_elemento_backlog.go](../internal/app/editar_elemento_backlog.go); 409 en el endpoint JSON y en la pantalla | `3d84b17` (spec) → `84bfe20` [RED] → `c2bd622` [GREEN] → `28664a5` (BDD) → este commit |
| DEF-03 La pantalla de un sprint cerrado mostraba las acciones de asignar y quitar | HU-09 (sección 3, RN-09.6, CA-09.5, CA-09.8), al integrarse con HU-11 | [#283](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/283) | Sin escenario BDD a propósito: es la pantalla (los escenarios prueban el núcleo, que ya rechazaba las acciones) | `TestPantallaSprintDetalle_SprintCerrado_SinAcciones`, `TestPantallaSprintDetalle_MuestraAsignadosYDisponibles` (con el sprint abierto siguen las acciones) | campo `Abierto` en `datosSprintDetalle` ([pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go)); aviso y acciones condicionales en [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | `3ffd6d6` (spec) → `a5dee62` [RED] → `1068d74` [GREEN] → este commit |
