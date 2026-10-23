# HU-03: Registrar integrantes de un proyecto

**Historia:** Como Product Builder, quiero registrar los integrantes de un proyecto, para saber quién participa.
**Issue:** #TODO (completar con el número de GitHub Projects)  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir asociar personas (integrantes) a un proyecto existente, con su nombre y su rol dentro del proyecto, para que el equipo y los reportes sepan quién participa. Cubre tanto agregar como quitar integrantes.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| id de proyecto | `int64` (`domain.Proyecto.ID`) | Sí | Debe existir un proyecto con ese id |
| nombre del integrante | `string` | Sí | No vacío (ni solo espacios); no puede repetirse dentro del mismo proyecto |
| rol | `string` | Sí | No vacío; texto libre (sin lista cerrada de valores) |

## 3. Salidas esperadas
- **Agregar:** el integrante queda incluido en la lista de integrantes del proyecto.
- **Quitar:** el integrante deja de figurar en la lista de integrantes del proyecto.
- No hay una respuesta de "datos" más allá de confirmar el resultado de la operación (200/204 en el adaptador HTTP); consultar la lista completa es responsabilidad de HU-04.

## 4. Reglas de negocio
- **RN-03.1:** El nombre del integrante es obligatorio; no se puede agregar un integrante sin nombre.
- **RN-03.2:** Un proyecto puede tener múltiples integrantes.
- **RN-03.3:** No se puede duplicar el mismo nombre dos veces en el mismo proyecto. La comparación usa `domain.ClaveNombre` (recorta espacios y no distingue mayúsculas), igual que HU-01 hace para el nombre de proyecto (RN-01.3) — se mantiene la misma convención del dominio.
- **RN-03.4:** El rol es obligatorio, pero no está restringido a una lista cerrada de valores.

## 5. Restricciones
- El núcleo (`internal/domain`, `internal/app`) no depende del framework web ni de la capa de persistencia (arquitectura hexagonal).
- La relación "proyecto tiene integrantes" se modela en el dominio (`Proyecto` mantiene su lista de `Integrante`), no en el adaptador HTTP.

## 6. Casos límite
- Proyecto sin ningún integrante todavía (lista vacía): válido, no es un error.
- Quitar al único integrante restante: válido, el proyecto queda con la lista vacía.
- Nombre con espacios al inicio/fin (ej. `" Ana "`): se normaliza con `domain.ClaveNombre` antes de comparar, igual que en HU-01.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Nombre vacío al agregar | `ErrNombreIntegranteVacio` | No se agrega el integrante; se informa el motivo |
| Rol vacío al agregar | `ErrRolIntegranteVacio` | No se agrega el integrante; se informa el motivo |
| Nombre duplicado en el mismo proyecto | `ErrIntegranteDuplicado` | No se agrega el integrante; se informa el motivo |
| Quitar un integrante que no pertenece al proyecto | `ErrIntegranteNoEncontrado` | No se modifica la lista; se informa el motivo |
| Proyecto inexistente (id no encontrado) | `domain.ErrProyectoNoEncontrado` (ya definido en HU-01, se reutiliza) | No se agrega/quita nada; 404 en el adaptador HTTP |

## 8. Criterios de aceptación
- **CA-03.1:** Dado un proyecto existente, cuando se agrega un integrante con nombre y rol válidos, entonces queda incluido en la lista de integrantes del proyecto.
- **CA-03.2:** Dado un proyecto existente, cuando se intenta agregar un integrante sin nombre, entonces la operación falla con `ErrNombreIntegranteVacio` y la lista no cambia.
- **CA-03.3:** Dado un proyecto con al menos un integrante, cuando se quita un integrante existente, entonces deja de figurar en la lista.
- **CA-03.4:** Dado un proyecto existente, cuando se intenta quitar un integrante que no pertenece al proyecto, entonces la operación falla con `ErrIntegranteNoEncontrado` y la lista no cambia.
- **CA-03.5:** Dado un proyecto con un integrante ya registrado, cuando se intenta agregar otro integrante con el mismo nombre, entonces la operación falla con `ErrIntegranteDuplicado` y la lista no cambia.
- **CA-03.6:** Dado un proyecto existente, cuando se intenta agregar un integrante con nombre válido pero sin rol, entonces la operación falla con `ErrRolIntegranteVacio` y la lista no cambia.

## Preguntas abiertas
- ¿Puede un mismo nombre de integrante tener roles distintos en el mismo proyecto (ej. agregarlo de nuevo con otro rol), o el nombre solo puede estar una vez sin importar el rol? Por ahora se asume que el nombre es único en el proyecto sin importar el rol (RN-03.3), por eso no se permite re-agregarlo con otro rol.

## Limitaciones conocidas
- **Concurrencia (lost update):** `AgregarIntegrante`/`QuitarIntegrante` leen el proyecto completo, modifican la lista de integrantes en memoria y la reemplazan entera, en dos pasos separados sin una transacción que cubra ambos. Si dos pedidos llegan casi al mismo tiempo para el mismo proyecto, el segundo puede pisar el cambio del primero (se pierde un integrante agregado/quitado). Detectado en la revisión de la Tarea 3.8 (go-reviewer). No se corrige en este sprint porque el TP no tiene carga concurrente real; si se vuelve relevante, la solución es hacer que `GuardarIntegrantes` opere fila por fila (INSERT/DELETE atómico) en vez de reemplazar toda la lista.

## Historial
- 2026-10-22: creada (borrador).
- 2026-10-23: resueltas las dudas de comparación de nombres duplicados (se usa `domain.ClaveNombre`, igual que HU-01) y el tipo del id de proyecto (`int64`), ya que HU-01 se mergeó a `main` y define `domain.Proyecto` y `domain.ErrProyectoNoEncontrado`.
