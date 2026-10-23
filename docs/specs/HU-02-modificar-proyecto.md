# HU-02: Modificar un proyecto

**Historia:** Como Product Builder, quiero modificar los datos de un proyecto existente (incluido su estado), para poder corregir su información o reflejar su situación actual.
**Issue:** #13  ·  **Sprint:** 1  ·  **Story points:** 2
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir editar de forma parcial un proyecto ya creado en HU-01: su nombre, sus fechas y su estado. Solo cambian los campos que se envían; el resto conserva su valor. Las reglas de validación son las mismas que al crear el proyecto (HU-01), y esta historia define las transiciones entre los estados de RN-01.7.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `id` | `int64` | Sí | Debe corresponder a un proyecto existente. |
| `nombre` | `*string` | No | Si se envía: se recorta, no vacío, máximo 100 caracteres (runas), único entre los demás proyectos (RN-01.1 a RN-01.3). |
| `fechaInicio` | `*time.Time` | No | Si se envía: fecha válida (solo día). No se puede quitar (es obligatoria, RN-01.4). |
| `fechaFin` | opcional y anulable | No | Si se envía una fecha: no puede ser anterior a la fecha de inicio resultante. Si se envía explícitamente `null`: el proyecto queda sin fecha de fin. |
| `estado` | `*domain.Estado` | No | Si se envía: debe ser uno de `Activo`, `En pausa`, `Finalizado`, `Postergado`, `Cancelado`. |

Se distingue "no enviado" (el campo no cambia) de "enviado" (el campo se actualiza). Esto importa sobre todo en `fechaFin`, donde enviar `null` significa quitarla.

## 3. Salidas esperadas
- Si la edición es válida: el proyecto queda persistido con los campos enviados actualizados y el resto sin cambios. Se devuelve el proyecto completo (`ID`, `Nombre`, `FechaInicio`, `FechaFin`, `Estado`).
- El `ID` nunca cambia.
- Si hay un error: no se modifica nada y se devuelve el error correspondiente.

## 4. Reglas de negocio
- **RN-02.1:** La edición es parcial: solo se modifican los campos enviados.
- **RN-02.2:** Si se envía `nombre`, se aplican las mismas reglas que al crear: se recorta, no puede quedar vacío (RN-01.1) y tiene como máximo 100 caracteres (RN-01.2).
- **RN-02.3:** El nombre no puede coincidir con el de **otro** proyecto, comparando sin distinguir mayúsculas ni espacios alrededor (RN-01.3). Renombrar un proyecto a su mismo nombre con otras mayúsculas (`"alfa"` → `"Alfa"`) está permitido.
- **RN-02.4:** La fecha de inicio no se puede quitar. Si se envía, reemplaza a la actual.
- **RN-02.5:** La fecha de fin resultante (la nueva, o la vigente si no se envió) no puede ser anterior a la fecha de inicio resultante (RN-01.5). Se puede quitar la fecha de fin enviándola en `null`.
- **RN-02.6:** El estado puede cambiarse a cualquiera de los cinco estados de RN-01.7, desde cualquier estado (por ejemplo, un proyecto `Finalizado` o `Cancelado` puede volver a `Activo`).
- **RN-02.7:** Enviar un campo con el mismo valor que ya tiene no es un error.
- **RN-02.8:** Debe enviarse al menos un campo para modificar.
- **RN-02.9:** El `ID` del proyecto no se modifica.

## 5. Restricciones
- Las validaciones de nombre y fechas reutilizan la lógica de dominio de HU-01 (`internal/domain`); no se duplican.
- La unicidad del nombre (RN-02.3) se resuelve en el caso de uso (`internal/app`) a través del puerto de repositorio, igual que en HU-01, excluyendo al propio proyecto.
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001).
- No hay autenticación ni roles: cualquier usuario de la app puede modificar un proyecto.

## 6. Casos límite
- Enviar solo `nombre`: las fechas y el estado no cambian.
- Enviar solo `fechaFin` anterior a la `fechaInicio` guardada → error.
- Enviar solo `fechaInicio` posterior a la `fechaFin` guardada → error.
- Enviar `fechaFin: null` → el proyecto queda sin fecha de fin.
- Renombrar `"Alfa"` a `"alfa"` (el mismo proyecto) → válido.
- Renombrar a `"Beta"` cuando ya existe otro proyecto `"beta"` → error de duplicado.
- Nombre de exactamente 100 caracteres → válido; de 101 → error.
- Cambiar de `Cancelado` a `Activo` → válido.
- Enviar el mismo estado que ya tiene → válido, sin cambios.
- Enviar un estado que no existe (`"Cerrado"`, `"activo"` en minúscula) → error.
- No enviar ningún campo → error.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe un proyecto con ese ID | `ErrProyectoNoEncontrado` | "Proyecto no encontrado." (HTTP 404) |
| No se envió ningún campo | `ErrSinCambios` | "No se indicó ningún dato para modificar." |
| Nombre vacío o solo espacios | `ErrNombreVacio` | "El nombre del proyecto es obligatorio." |
| Nombre de más de 100 caracteres | `ErrNombreDemasiadoLargo` | "El nombre no puede superar los 100 caracteres." |
| Otro proyecto ya tiene ese nombre | `ErrNombreDuplicado` | "Ya existe un proyecto con ese nombre." |
| Se intenta quitar la fecha de inicio | `ErrFechaInicioObligatoria` | "La fecha de inicio es obligatoria." |
| Fecha de fin resultante anterior a la de inicio | `ErrFechaFinAnterior` | "La fecha de fin no puede ser anterior a la de inicio." |
| Estado que no es uno de los cinco válidos | `ErrEstadoInvalido` | "El estado indicado no es válido." |

`ErrProyectoNoEncontrado` y los errores de nombre y fechas ya existen en `internal/domain` (HU-01). Esta historia agrega `ErrSinCambios` y `ErrEstadoInvalido`. En todos los casos no se persiste nada; si hay varios errores a la vez, se informa el primero en el orden de la tabla.

## 8. Criterios de aceptación
- **CA-02.1:** Dado un proyecto existente, cuando envío solo un nombre válido, entonces se actualiza el nombre y las fechas y el estado no cambian. (RN-02.1, RN-02.2)
- **CA-02.2:** Dado un proyecto existente, cuando envío un nombre vacío o de más de 100 caracteres, entonces obtengo `ErrNombreVacio` o `ErrNombreDemasiadoLargo` y el proyecto no cambia. (RN-02.2)
- **CA-02.3:** Dados los proyectos "Alfa" y "Beta", cuando renombro "Beta" a "alfa", entonces obtengo `ErrNombreDuplicado`; cuando renombro "Alfa" a "alfa", entonces se actualiza. (RN-02.3)
- **CA-02.4:** Dado un proyecto existente, cuando envío una fecha de fin anterior a su fecha de inicio (o una fecha de inicio posterior a su fecha de fin), entonces obtengo `ErrFechaFinAnterior` y el proyecto no cambia. (RN-02.5)
- **CA-02.5:** Dado un proyecto con fecha de fin, cuando envío `fechaFin: null`, entonces queda sin fecha de fin. (RN-02.5)
- **CA-02.6:** Dado un proyecto `Activo`, cuando envío el estado `Finalizado`, entonces queda `Finalizado`; y dado un proyecto `Cancelado`, cuando envío `Activo`, entonces queda `Activo`. (RN-02.6)
- **CA-02.7:** Dado un proyecto existente, cuando envío un estado que no es uno de los cinco válidos, entonces obtengo `ErrEstadoInvalido` y el proyecto no cambia. (RN-02.6)
- **CA-02.8:** Dado un proyecto existente, cuando no envío ningún campo, entonces obtengo `ErrSinCambios`. (RN-02.8)
- **CA-02.9:** Dado un ID que no existe, cuando intento modificarlo, entonces obtengo `ErrProyectoNoEncontrado`. (RN-02.9)
- **CA-02.10:** Dado un proyecto existente, cuando lo modifico, entonces su `ID` sigue siendo el mismo y la modificación queda persistida al recuperarlo. (RN-02.9)
- **CA-02.11:** Dado un proyecto existente, cuando intento quitar su fecha de inicio, entonces obtengo `ErrFechaInicioObligatoria` y el proyecto no cambia. (RN-02.4)

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-22: creada (borrador).
- 2026-10-22: edición parcial y cambio de estado en ambos sentidos, según decisión del equipo.
- 2026-10-23: alineada al dominio implementado en HU-01: cinco estados de RN-01.7 (en lugar de "Activo/Cerrado"), fecha de fin opcional y anulable, reglas de nombre (máximo 100, único sin distinguir mayúsculas) y nombres de errores existentes.
- 2026-10-23: agregado CA-02.11 para cubrir RN-02.4 (no se puede quitar la fecha de inicio).
