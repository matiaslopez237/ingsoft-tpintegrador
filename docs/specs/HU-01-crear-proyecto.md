# HU-01: Crear un proyecto

**Historia:** Como Product Builder, quiero crear un proyecto con nombre, fecha de inicio y fin, para poder empezar a planificar el trabajo.
**Issue:** #12  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir que un Product Builder dé de alta un proyecto de software, que es la entidad raíz del sistema: backlog, sprints, esfuerzo, defectos y métricas cuelgan de él. El proyecto se crea con nombre y fechas válidas, queda en estado Activo y con un identificador único para referenciarlo después (HU-02, HU-03, HU-04).

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `nombre` | `string` | Sí | Se recortan espacios al inicio y al final. No vacío. Máximo 100 caracteres (runas). Único entre los proyectos existentes. |
| `fechaInicio` | `time.Time` | Sí | Fecha válida (sin hora; solo día). |
| `fechaFin` | `*time.Time` | No | Si se informa: no puede ser anterior a `fechaInicio`. |

## 3. Salidas esperadas
- Se devuelve el proyecto creado con:
  - `ID`: identificador único generado por el sistema.
  - `Nombre`: ya recortado.
  - `FechaInicio`, `FechaFin` (o sin fecha de fin).
  - `Estado`: `Activo`.
- El proyecto queda persistido y puede recuperarse por su `ID`.
- Si hay un error, no se crea ni persiste nada y se devuelve el error correspondiente.

## 4. Reglas de negocio
- **RN-01.1:** El nombre es obligatorio. Un nombre formado solo por espacios se considera vacío.
- **RN-01.2:** El nombre tiene como máximo 100 caracteres, contados después de recortar los espacios.
- **RN-01.3:** No puede haber dos proyectos con el mismo nombre. La comparación se hace sobre el nombre recortado y sin distinguir mayúsculas de minúsculas ("Alfa" y "alfa" son el mismo nombre).
- **RN-01.4:** La fecha de inicio es obligatoria.
- **RN-01.5:** La fecha de fin es opcional; si se informa, no puede ser anterior a la fecha de inicio (puede ser igual).
- **RN-01.6:** Todo proyecto nuevo queda en estado `Activo`.
- **RN-01.7:** Los estados posibles de un proyecto son: `Activo`, `En pausa`, `Finalizado`, `Postergado` y `Cancelado`. En esta historia solo se usa `Activo`; las transiciones entre estados se definen en otras historias.
- **RN-01.8:** Al crearse, el proyecto recibe un identificador único generado por el sistema: un entero autoincremental (`int64`) asignado por la persistencia.

## 5. Restricciones
- Las validaciones (RN-01.1, .2, .4, .5, .6) viven en el dominio (`internal/domain`), sin dependencias de HTTP ni de base de datos.
- La unicidad del nombre (RN-01.3) y la generación del ID (RN-01.8) se resuelven en el caso de uso `CrearProyecto` (`internal/app`) a través de un puerto de repositorio.
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001).
- No hay autenticación ni roles en esta historia: cualquier usuario de la app puede crear un proyecto.

## 6. Casos límite
- Fecha de fin igual a la de inicio → válido.
- Sin fecha de fin → válido; el proyecto queda sin fecha de cierre.
- Nombre de exactamente 100 caracteres → válido; de 101 → error.
- Nombre con espacios alrededor (`"  Mi proyecto  "`) → se guarda como `"Mi proyecto"`.
- Nombre con caracteres acentuados o ñ (`"Gestión año 2026"`) → se cuentan como 1 carácter cada uno.
- Nombre `"   "` (solo espacios) → error de nombre vacío.
- Nombre ya usado por otro proyecto, aunque cambien mayúsculas o espacios alrededor (`"alfa "` vs `"Alfa"`) → error de duplicado.
- Fecha de fin un día antes de la de inicio → error.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Nombre vacío o solo espacios | `ErrNombreVacio` | "El nombre del proyecto es obligatorio." |
| Nombre de más de 100 caracteres | `ErrNombreDemasiadoLargo` | "El nombre no puede superar los 100 caracteres." |
| Ya existe un proyecto con ese nombre | `ErrNombreDuplicado` | "Ya existe un proyecto con ese nombre." |
| Falta la fecha de inicio | `ErrFechaInicioObligatoria` | "La fecha de inicio es obligatoria." |
| Fecha de fin anterior a la de inicio | `ErrFechaFinAnterior` | "La fecha de fin no puede ser anterior a la de inicio." |

En todos los casos no se persiste nada. Si hay varios errores a la vez, se informa el primero en el orden de la tabla.

## 8. Criterios de aceptación
- **CA-01.1:** Dado un nombre válido, fecha de inicio y fecha de fin posterior, cuando creo el proyecto, entonces se crea en estado `Activo` con un identificador único. (RN-01.6, RN-01.8)
- **CA-01.2:** Dado un nombre vacío o formado solo por espacios, cuando creo el proyecto, entonces obtengo `ErrNombreVacio` y no se crea. (RN-01.1)
- **CA-01.3:** Dado un nombre de más de 100 caracteres, cuando creo el proyecto, entonces obtengo `ErrNombreDemasiadoLargo`; con exactamente 100 se crea. (RN-01.2)
- **CA-01.4:** Dado que ya existe un proyecto "Alfa", cuando creo otro llamado "Alfa" o "alfa", entonces obtengo `ErrNombreDuplicado` y no se crea. (RN-01.3)
- **CA-01.5:** Dado que no informo fecha de inicio, cuando creo el proyecto, entonces obtengo `ErrFechaInicioObligatoria`. (RN-01.4)
- **CA-01.6:** Dada una fecha de fin anterior a la de inicio, cuando creo el proyecto, entonces obtengo `ErrFechaFinAnterior` y no se crea. (RN-01.5)
- **CA-01.7:** Dada una fecha de fin igual a la de inicio, o sin fecha de fin, cuando creo el proyecto, entonces se crea correctamente. (RN-01.5)
- **CA-01.8:** Dado un nombre con espacios alrededor, cuando creo el proyecto, entonces se guarda recortado. (RN-01.1)
- **CA-01.9:** Dados dos proyectos creados, entonces sus identificadores son distintos y cada uno puede recuperarse por su ID. (RN-01.8)

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-22: creada (borrador).
