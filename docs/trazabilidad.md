# Matriz de trazabilidad

Historia → Spec SDD → Criterios de aceptación → Escenarios BDD → Tests Go → Código Go → Commits TDD.

Cómo se busca cada eslabón: tags `@HU-XX` / `@CA-XX.N` en `features/`, comentarios `// Traza: HU-XX ...` en los tests, `Implementa HU-XX` en el código y `git log --grep="(HU-XX)"`.

## Resumen

| Historia | Spec | CA con BDD | CA con test Go | Código | TDD en git | BDD en verde |
|---|---|---|---|---|---|---|
| HU-01 Crear un proyecto | ✓ (Borrador, pendiente de aprobación en el PR) | 9/9 | 9/9 | ✓ | RED → GREEN ×2, REFACTOR ✓ | ✓ 14/14 escenarios |
| HU-03 Registrar integrantes | ✓ (Borrador, pendiente de aprobación en el PR) | 10/10 | 10/10 | ✓ | RED → GREEN ×10 | ✓ 13/13 escenarios |
| HU-04 Consultar el estado de un proyecto | ✓ (Borrador, pendiente de aprobación en el PR) | 5/5 | 5/5 | ✓ | RED → GREEN ×3, REFACTOR ×2 | ✓ 6/6 escenarios |

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
| HTTP | [internal/adapters/http/consultar_proyecto_handler.go](../internal/adapters/http/consultar_proyecto_handler.go) (`GET /proyectos/{id}`, `RegistrarConsultarProyecto`), registrada en [cmd/server/main.go](../cmd/server/main.go) | [consultar_proyecto_handler_test.go](../internal/adapters/http/consultar_proyecto_handler_test.go) |
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
