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
| HU-18 Registrar esfuerzo en una historia o tarea | ✓ (Borrador; P-6, la ruta, pendiente de confirmar) | 14/15 (CA-18.15 se cubre con tests del handler) | 14/15 (CA-18.11 solo en BDD) | ✓ | RED → GREEN ×14, REFACTOR ×4 | ✓ 28/28 escenarios |
| HU-30 Interfaz web base | ✓ (Borrador, pendiente de aprobación en el PR) | 10/10 | 10/10 | ✓ | RED → GREEN ×3, REFACTOR ✓ | ✓ 17/17 escenarios |

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

- **Issue:** #29 · Sprint 1 · 3 SP
- **Spec:** [docs/specs/HU-18-registrar-esfuerzo.md](specs/HU-18-registrar-esfuerzo.md)
- **Feature:** [features/HU-18-registrar-esfuerzo.feature](../features/HU-18-registrar-esfuerzo.feature) (steps en [test/bdd/hu18_steps_test.go](../test/bdd/hu18_steps_test.go))
- **Ruta HTTP:** `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`. **Pendiente de confirmar (P-6):** reemplaza el `POST /esfuerzo` de la tarea 18.7 y falta el OK de Matías.

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

### Archivos de código y tests

| Capa | Código | Tests |
|---|---|---|
| Dominio | [internal/domain/registro_esfuerzo.go](../internal/domain/registro_esfuerzo.go) (`RegistroEsfuerzo`, `NuevoRegistroEsfuerzo`, `validarHoras`, `soloDia`, `diaCalendario`) | [registro_esfuerzo_test.go](../internal/domain/registro_esfuerzo_test.go) |
| Caso de uso | [internal/app/registrar_esfuerzo.go](../internal/app/registrar_esfuerzo.go) (`RegistrarEsfuerzo`, puertos `LectorElementos` y `RepositorioEsfuerzo`; reutiliza `LectorProyectos` de HU-04) | [registrar_esfuerzo_test.go](../internal/app/registrar_esfuerzo_test.go) |
| Persistencia | [internal/adapters/persistencia/esfuerzo_repo.go](../internal/adapters/persistencia/esfuerzo_repo.go) (`GuardarRegistro`, `BuscarRegistroPorID`), [db.go](../internal/adapters/persistencia/db.go) (tabla `registros_esfuerzo`) | [esfuerzo_repo_test.go](../internal/adapters/persistencia/esfuerzo_repo_test.go), contra SQLite real |
| HTTP | [internal/adapters/http/esfuerzo_handler.go](../internal/adapters/http/esfuerzo_handler.go) (`RegistrarRegistrarEsfuerzo`), registrada en [router.go](../internal/adapters/http/router.go) y conectada en [cmd/server/main.go](../cmd/server/main.go) con `time.Now` como reloj | [esfuerzo_handler_test.go](../internal/adapters/http/esfuerzo_handler_test.go), [router_test.go](../internal/adapters/http/router_test.go) |
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

### Huecos

- La spec sigue en **Borrador**. **P-6 (la ruta HTTP) está pendiente de confirmar con Matías**: si cambia, hay que tocar `esfuerzo_handler.go`, `router_test.go`, CA-18.15 y esta matriz.
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

### Huecos

- La spec está en **Borrador** hasta la aprobación del PR.
- RN-06.8 (Sprint activo) depende de HU-09.
- La lista del backlog (HU-07) enlaza a la pantalla del elemento (Editar) y permite eliminar con confirmación.

---

## Defectos

| Defecto | Historia | Issue | Escenario BDD | Tests Go | Corrección | Commits |
|---|---|---|---|---|---|---|
| DEF-01 El formulario de edición de proyecto pisa cambios simultáneos | HU-30 (RN-30.3, CA-30.5) | [#179](https://github.com/matiaslopez237/ingsoft-tpintegrador/issues/179) | Dos ediciones simultáneas de campos distintos no se pisan `@DEF-01` | `TestPantallas_EditarProyecto_FormularioLlevaLosValoresOriginales`, `TestPantallas_EditarProyecto_SoloEnviaLosCamposCambiados`, `TestPantallas_EditarProyecto_SinCambios_VuelveAlDetalleSinError` | `cambiosDeForm` y `leerOriginalProyecto` en [pantallas.go](../internal/adapters/http/pantallas.go); campos ocultos `original_*` en [proyecto_form.html](../internal/adapters/http/web/plantillas/proyecto_form.html) | `1bbf72a` [RED] → `49c6303` [GREEN] → `7af3f6e` (spec y trazabilidad); sin cambios vuelve al detalle: `4a95342` [RED] → `73439b0` [GREEN] → este commit |
