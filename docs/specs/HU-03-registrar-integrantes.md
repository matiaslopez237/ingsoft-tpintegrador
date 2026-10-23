# HU-03: Registrar integrantes de un proyecto

**Historia:** Como Product Builder, quiero registrar los integrantes de un proyecto, para saber quién participa.
**Issue:** #14  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir asociar personas (integrantes) a un proyecto existente, con su nombre y sus roles dentro del proyecto, para que el equipo y los reportes sepan quién participa. Cubre agregar y quitar integrantes, y sumarle roles a un integrante que ya está en el proyecto (en un equipo chico una misma persona suele cumplir varios roles).

## 2. Entradas
**Agregar integrante**

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| id de proyecto | `int64` (`domain.Proyecto.ID`) | Sí | Debe existir un proyecto con ese id |
| nombre del integrante | `string` | Sí | No vacío (ni solo espacios); no puede repetirse dentro del mismo proyecto |
| rol | `string` | Sí | No vacío; texto libre (sin lista cerrada de valores). Es el primer rol del integrante. |

**Agregar rol a un integrante**

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| id de proyecto | `int64` | Sí | Debe existir un proyecto con ese id |
| nombre del integrante | `string` | Sí | Debe ser un integrante del proyecto |
| rol | `string` | Sí | No vacío; el integrante no puede tenerlo ya |

**Quitar integrante:** id de proyecto y nombre del integrante, con las mismas validaciones.

## 3. Salidas esperadas
- **Agregar:** el integrante queda incluido en la lista de integrantes del proyecto, con un rol.
- **Agregar rol:** el integrante suma el rol a los que ya tenía. Sigue figurando una sola vez en la lista.
- **Quitar:** el integrante deja de figurar en la lista de integrantes del proyecto, con todos sus roles.
- No hay una respuesta de "datos" más allá de confirmar el resultado de la operación (204 en el adaptador HTTP); consultar la lista completa es responsabilidad de HU-04.

## 4. Reglas de negocio
- **RN-03.1:** El nombre del integrante es obligatorio; no se puede agregar un integrante sin nombre.
- **RN-03.2:** Un proyecto puede tener múltiples integrantes.
- **RN-03.3:** No se puede duplicar el mismo nombre dos veces en el mismo proyecto. La comparación usa `domain.ClaveNombre` (recorta espacios y no distingue mayúsculas), igual que HU-01 hace para el nombre de proyecto (RN-01.3) — se mantiene la misma convención del dominio.
- **RN-03.4:** El rol es obligatorio, pero no está restringido a una lista cerrada de valores.
- **RN-03.5:** Un integrante puede tener más de un rol en el mismo proyecto. Se agrega con un primer rol, y los demás se le suman con la operación "agregar rol", sin volver a agregar al integrante (que seguiría rechazándose por RN-03.3).
- **RN-03.6:** Un integrante no puede tener dos veces el mismo rol. La comparación usa `domain.ClaveNombre` (recorta espacios y no distingue mayúsculas), la misma convención que RN-03.3.

## 5. Restricciones
- El núcleo (`internal/domain`, `internal/app`) no depende del framework web ni de la capa de persistencia (arquitectura hexagonal).
- La relación "proyecto tiene integrantes" se modela en el dominio (`Proyecto` mantiene su lista de `Integrante`), no en el adaptador HTTP.
- **Pedidos simultáneos:** cada operación cambia solo su propia fila en la base (inserta un integrante o un rol, o borra un integrante); nunca reemplaza la lista completa. Así, dos pedidos sobre el mismo proyecto que llegan al mismo tiempo no se pisan. Las reglas RN-03.3 y RN-03.6 también se garantizan en la base con restricciones `UNIQUE`, por si dos pedidos pasan la validación del dominio a la vez.

## 6. Casos límite
- Proyecto sin ningún integrante todavía (lista vacía): válido, no es un error.
- Quitar al único integrante restante: válido, el proyecto queda con la lista vacía.
- Nombre con espacios al inicio/fin (ej. `" Ana "`): se normaliza con `domain.ClaveNombre` antes de comparar, igual que en HU-01.
- Agregar a "Ana" como Tester cuando ya está como Desarrolladora: se rechaza por duplicado (RN-03.3); el rol se suma con "agregar rol".
- Agregar el rol "desarrolladora" a quien ya tiene "Desarrolladora": se rechaza por rol duplicado (RN-03.6).
- Quitar a un integrante con varios roles: se quitan todos sus roles junto con él.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Nombre vacío al agregar | `ErrNombreIntegranteVacio` | No se agrega el integrante; se informa el motivo |
| Rol vacío al agregar | `ErrRolIntegranteVacio` | No se agrega el integrante; se informa el motivo |
| Nombre duplicado en el mismo proyecto | `ErrIntegranteDuplicado` | No se agrega el integrante; se informa el motivo |
| Quitar un integrante, o agregarle un rol, cuando no pertenece al proyecto | `ErrIntegranteNoEncontrado` | No se modifica la lista; se informa el motivo |
| Rol vacío al agregar un rol | `ErrRolIntegranteVacio` | No se agrega el rol; se informa el motivo |
| El integrante ya tiene ese rol | `ErrRolDuplicado` | No se agrega el rol; se informa el motivo |
| Proyecto inexistente (id no encontrado) | `domain.ErrProyectoNoEncontrado` (ya definido en HU-01, se reutiliza) | No se agrega/quita nada; 404 en el adaptador HTTP |

## 8. Criterios de aceptación
- **CA-03.1:** Dado un proyecto existente, cuando se agrega un integrante con nombre y rol válidos, entonces queda incluido en la lista de integrantes del proyecto.
- **CA-03.2:** Dado un proyecto existente, cuando se intenta agregar un integrante sin nombre, entonces la operación falla con `ErrNombreIntegranteVacio` y la lista no cambia.
- **CA-03.3:** Dado un proyecto con al menos un integrante, cuando se quita un integrante existente, entonces deja de figurar en la lista.
- **CA-03.4:** Dado un proyecto existente, cuando se intenta quitar un integrante que no pertenece al proyecto, entonces la operación falla con `ErrIntegranteNoEncontrado` y la lista no cambia.
- **CA-03.5:** Dado un proyecto con un integrante ya registrado, cuando se intenta agregar otro integrante con el mismo nombre, entonces la operación falla con `ErrIntegranteDuplicado` y la lista no cambia.
- **CA-03.6:** Dado un proyecto existente, cuando se intenta agregar un integrante con nombre válido pero sin rol, entonces la operación falla con `ErrRolIntegranteVacio` y la lista no cambia.
- **CA-03.7:** Dado que "Ana" está en el proyecto como Desarrolladora, cuando se le agrega el rol Tester, entonces Ana tiene los roles Desarrolladora y Tester y figura una sola vez en la lista. (RN-03.5)
- **CA-03.8:** Dado que "Ana" tiene el rol Desarrolladora, cuando se le agrega el rol "desarrolladora", entonces la operación falla con `ErrRolDuplicado` y sus roles no cambian. (RN-03.6)
- **CA-03.9:** Dado un proyecto sin la integrante "Caro", cuando se le agrega un rol a "Caro", entonces la operación falla con `ErrIntegranteNoEncontrado`.
- **CA-03.10:** Dado que "Ana" está en el proyecto, cuando se le agrega un rol vacío, entonces la operación falla con `ErrRolIntegranteVacio` y sus roles no cambian. (RN-03.4)

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-22: creada (borrador).
- 2026-10-23: resueltas las dudas de comparación de nombres duplicados (se usa `domain.ClaveNombre`, igual que HU-01) y el tipo del id de proyecto (`int64`), ya que HU-01 se mergeó a `main` y define `domain.Proyecto` y `domain.ErrProyectoNoEncontrado`.
- 2026-10-23: issue #14. El equipo confirma RN-03.3 (un integrante no se agrega dos veces) y decide que un integrante pueda tener varios roles, sumados con una operación aparte (RN-03.5, RN-03.6, CA-03.7 a CA-03.10). Se resuelve la concurrencia (lost update) guardando fila por fila en vez de reemplazar la lista completa (sección 5).
