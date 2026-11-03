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
| HU-14 Registrar estimaciones individuales ocultas | ✓ (Borrador, pendiente de aprobación en el PR) | 9/9 | 9/9 | ✓ | RED → GREEN ×5, REFACTOR ✓ | ✓ 20/20 escenarios |
| HU-15 Mostrar las estimaciones y detectar diferencias | ✓ (Borrador, pendiente de aprobación en el PR) | 6/6 | 6/6 | ✓ | RED → GREEN ×3 | ✓ 10/10 escenarios |
| HU-16 Repetir una ronda de votación | ✓ (Implementada) | 7/8 (CA-16.7 se cubre con tests del handler) | 8/8 | ✓ | RED → GREEN ×5 | ✓ 8/8 escenarios |
| HU-17 Registrar la estimación acordada | ✓ (Implementada) | 9/10 (CA-17.9 se cubre con tests del handler) | 10/10 | ✓ | RED → GREEN ×5 | ✓ 17/17 escenarios |
| HU-19 Comparar esfuerzo estimado contra real | ✓ (Borrador; P-10 y P-12 abiertas a objeción del equipo en la PR) | 15/19 (CA-19.16 se cubre con tests de persistencia; CA-19.17, del handler; CA-19.18 y CA-19.19, de la pantalla) | 19/19 | ✓ | RED → GREEN ×11, REFACTOR ×2 | ✓ 34/34 escenarios (336/336 en el BDD completo; pasos: 1687/1687) |

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
| CA-14.5 Voto inválido, no votante, ronda cerrada | RN-14.2, RN-14.4, RN-14.6 | Un voto fuera de la escala de Fibonacci `@error` (×4) · Vota alguien que se sumó al proyecto después de abrir la ronda `@error` · No se vota en una ronda cerrada `@error` | `TestVotar_Errores` (dominio y caso de uso), `TestVotar_RondaCerrada`, `TestRepositorioRondas_GuardarVoto_RondaYaCerrada`, `TestValoresVoto` | `domain.Ronda.Votar`, `persistencia.RepositorioRondas.GuardarVoto` |
| CA-14.6 No se puede iniciar | RN-14.1, RN-14.2 | No se estima un elemento terminado `@error` · Un proyecto sin integrantes no puede abrir una ronda `@error` · Un elemento tiene una sola ronda abierta a la vez `@error` · Con la ronda anterior cerrada se abre la número 2 `@limite` | `TestNuevaRonda_Errores`, `TestNuevaRonda_ConLaAnteriorCerrada_TomaElNumeroSiguiente`, `TestIniciarRonda_Errores`, `TestIniciarRonda_ConLaAnteriorCerrada_AbreLaSiguiente`, `TestRepositorioRondas_UltimaRonda`, `TestRepositorioRondas_UnaSolaAbiertaPorElemento` | `domain.NuevaRonda`, índice único `rondas_una_abierta_por_elemento` |
| CA-14.7 Otro proyecto o inexistente | RN-14.7 | No se abre una ronda sobre un elemento de otro proyecto `@error` · No se vota en una ronda de otro proyecto `@error` | `TestIniciarRonda_Errores`, `TestVotar_Errores`, `TestConsultarRonda`, `TestConsultarRonda_UltimaDelElemento`, `TestPantallaRonda_Inexistente` | `buscarElementoDelProyecto`, `rondaDelProyecto` |
| CA-14.8 API sin valores | Sección 3, RN-14.3 | El resumen se pide al endpoint real en los escenarios de CA-14.2 y CA-14.3 ("el resumen de la ronda no muestra ningún valor") | `TestPostRonda_Devuelve201ConElResumen`, `TestPostVoto_Devuelve200SinElValor`, `TestGetRonda_Devuelve200ConElResumen`, `TestPlanningPoker_Errores` (17 casos) | `planning_poker_handler.go` |
| CA-14.9 Pantalla | Sección 3 | Iniciar una ronda desde la pantalla del elemento `@normal` · Votar desde la pantalla sin ver los votos `@normal` | `TestPantallaElemento_*` (4 tests de la sección Planning Poker), `TestPantallaIniciarRonda_*`, `TestPantallaRonda_*`, `TestPantallaVotar_*` | `pantallas_planning_poker.go`, `web/plantillas/ronda.html`, sección en `web/plantillas/elemento.html` (HU-06) |

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/planning_poker.go](../internal/domain/planning_poker.go) (`Ronda`, `NuevaRonda`, `Votar`, `Resumen`, `ValoresVoto`) | [planning_poker_test.go](../internal/domain/planning_poker_test.go) |
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
| 16.9 Revisión y trazabilidad | este commit | — |

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

- **Depende de HU-13 (#288):** HU-17 usa `valoresVoto` y `ValoresVoto`, que esa rama reemplaza por la escala unificada (RN-13.4). Al combinarlas hay que cambiar esas referencias.
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
| Trazabilidad | este commit | — |

**Commits que el filtro `git log --grep="(HU-19)"` no encuentra**, porque usan el prefijo `test(bdd):`: `19eb4be` (RED de los steps), `4d25796` (sistema SQLite) y `c4b3c8a` (comentarios).

**Merges de la rama, que no son de HU-19:** `29991ac` (trae `main`: HU-14 y HU-15), `a106188` (trae `hu-18-pantalla-esfuerzo`) y `a23c87e` (trae `main`: HU-16 y HU-17).

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
- **P-10 (campo nuevo de horas estimadas, toca HU-05 y HU-06)** y **P-12 (tope de 1000 h, arbitrario)** están implementadas según la propuesta y abiertas a objeción del equipo en la PR.
- **Carrera de `agregarColumnas` en las migraciones:** issue a crear. Es preexistente (también `sprint_id` y `puntos_*`, por lectura del código) y compartida. Con 4 `AbrirDB` simultáneos sobre el esquema anterior a HU-19, hubo fallas en 30 de 30 rondas (89 de 120 aperturas, `duplicate column name`), con conexiones en un mismo proceso. Según ADR-0004, en el plan gratuito de Render el disco se borra en cada despliegue, así que hoy esas bases nacen con la columna.
- El mensaje "las horas se registran en múltiplos de media hora" al estimar viene de HU-18 y la sección 7 de la spec lo fija. Se entiende, pero habla de registrar.
- **Asimetría conocida, para decisión del equipo:** un POST del formulario sin campos ocultos conserva las horas estimadas (RN-19.12), pero sigue quitando los story points (comportamiento de HU-06, DEF-01). Comprobado contra el servidor real.
- Si el formulario vuelve con otro error en un pedido sin el campo de horas, las muestra vacías aunque haya guardadas. Reenviarlo no las borra, y desde el navegador el campo siempre viene.
- Si la comparación falla al volver de un error de edición, se responde 500 y se pierde lo tipeado. Es intencional y está probado, pero la spec no lo menciona.
- La pantalla **no se probó en un navegador** ni con un lector de pantalla: se probó con tests y con curl contra el servidor.
- La duración del BDD completo subió a unos 224 s: cada escenario de cualquier historia crea una base más, porque los hooks de todas las historias corren siempre.
- Si HU-13 (que toca `story_points` en los mismos archivos) entra antes, habrá que traer `main` de nuevo.
- **Paso BDD renombrado al integrar con `main`:** `"X" tiene N story points` de HU-19 colisionaba con el mismo paso de HU-17; godog lo rechaza como ambiguo y el compilador no lo detecta. Se renombró en HU-19 a `"X" quedó con N story points` ([features/HU-19-comparar-esfuerzo.feature](../features/HU-19-comparar-esfuerzo.feature) y [test/bdd/hu19_steps_test.go](../test/bdd/hu19_steps_test.go)). La comprobación es la misma. Se tocó HU-19 porque HU-17 ya estaba en `main`. Tras integrar: 336/336 escenarios en el BDD completo (pasos: 1687/1687), y los de HU-19 siguen siendo 34/34.
- Etiquetas: los tests atados solo a una sección de la spec llevan el CA más cercano. El 500 usa CA-19.17, aunque ese CA no menciona el 500.
- El comentario de cabecera del feature nombra CA-19.16 a CA-19.18, pero no CA-19.19, que se agregó después. `cambiosDeFormElemento` explica la regla, pero no nombra RN-19.12.
- `mostrarElemento` consulta el proyecto y `CompararEsfuerzo` lo vuelve a consultar, igual que el elemento: son lecturas de más, sin consecuencias prácticas.
- Refactors que tocaron código de otras historias: `validarMediasHoras` en `validarHoras` (HU-18), `copiarPuntero` en `NuevoElementoBacklog` (HU-05), `Editar` y `ActualizarElemento` (HU-06), y `TestPantallaElemento_PlanningPokerYEsfuerzo_Juntos` (HU-18), que suma la sección nueva. La descripción de ese test en HU-18 todavía enumera las secciones anteriores. Conviene que lo revisen sus autores en la PR.

---

## Defectos

| Defecto | Historia | Issue | Escenario BDD | Tests Go | Corrección | Commits |
|---|---|---|---|---|---|---|
| DEF-01 El formulario de edición de proyecto pisa cambios simultáneos | HU-30 (RN-30.3, CA-30.5) | [#179](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/179) | Dos ediciones simultáneas de campos distintos no se pisan `@DEF-01` | `TestPantallas_EditarProyecto_FormularioLlevaLosValoresOriginales`, `TestPantallas_EditarProyecto_SoloEnviaLosCamposCambiados`, `TestPantallas_EditarProyecto_SinCambios_VuelveAlDetalleSinError` | `cambiosDeForm` y `leerOriginalProyecto` en [pantallas.go](../internal/adapters/http/pantallas.go); campos ocultos `original_*` en [proyecto_form.html](../internal/adapters/http/web/plantillas/proyecto_form.html) | `1bbf72a` [RED] → `49c6303` [GREEN] → `7af3f6e` (spec y trazabilidad); sin cambios vuelve al detalle: `4a95342` [RED] → `73439b0` [GREEN] → este commit |
| DEF-02 Eliminar un elemento con defectos informaba "con esfuerzo registrado" | HU-06 (RN-06.10, CA-06.11), al integrarse con HU-20 | [#205](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/205) | No se elimina un elemento con un defecto registrado `@DEF-02` | `TestEliminarElemento_ConDefectos_NoLoElimina`, `TestEliminarElemento_ConEsfuerzoYDefectos_InformaElEsfuerzo`, `TestRepositorioBacklog_ElementoConDefectos`, `TestDeleteElemento`, `TestPantallaElemento_Errores` | `ErrElementoConDefectos` en [backlog_editar.go](../internal/domain/backlog_editar.go); `TieneDefectos` y `errorDeEliminacion` en [backlog_editar.go](../internal/adapters/persistencia/backlog_editar.go); verificación en [editar_elemento_backlog.go](../internal/app/editar_elemento_backlog.go); 409 en el endpoint JSON y en la pantalla | `3d84b17` (spec) → `84bfe20` [RED] → `c2bd622` [GREEN] → `28664a5` (BDD) → este commit |
| DEF-03 La pantalla de un sprint cerrado mostraba las acciones de asignar y quitar | HU-09 (sección 3, RN-09.6, CA-09.5, CA-09.8), al integrarse con HU-11 | [#283](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/283) | Sin escenario BDD a propósito: es la pantalla (los escenarios prueban el núcleo, que ya rechazaba las acciones) | `TestPantallaSprintDetalle_SprintCerrado_SinAcciones`, `TestPantallaSprintDetalle_MuestraAsignadosYDisponibles` (con el sprint abierto siguen las acciones) | campo `Abierto` en `datosSprintDetalle` ([pantallas_asignacion.go](../internal/adapters/http/pantallas_asignacion.go)); aviso y acciones condicionales en [sprint_detalle.html](../internal/adapters/http/web/plantillas/sprint_detalle.html) | `3ffd6d6` (spec) → `a5dee62` [RED] → `1068d74` [GREEN] → este commit |
