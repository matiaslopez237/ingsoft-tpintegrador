# HU-18: Registrar esfuerzo en una historia o tarea

**Historia:** Como Product Builder, quiero registrar integrante, fecha, actividad y horas trabajadas en una historia o tarea, para dejar constancia del esfuerzo real.
**Issue:** #29  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Dejar constancia de las horas que cada integrante trabajó realmente sobre un elemento del backlog: quién, qué día, en qué actividad y cuánto tiempo. Es el dato base del esfuerzo real, que después usan las métricas (estimado vs. real), el dashboard y los reportes. Reutiliza los elementos de backlog de HU-05 y los integrantes de HU-03. En esta historia solo se expone un endpoint JSON; la pantalla para cargar el esfuerzo se suma después con HU-30.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Debe corresponder a un proyecto existente (HU-01). |
| `elementoID` | `int64` | Sí | Debe corresponder a un elemento de backlog (HU-05) del mismo proyecto. |
| `integrante` | `string` | Sí | Se recortan espacios. No vacío. Debe ser integrante del proyecto, comparando con `domain.ClaveNombre`, como en HU-03. |
| `fecha` | `time.Time` | Sí | Solo importa el día (sin hora). No puede ser posterior a hoy. |
| `actividad` | `string` | Sí | Se recortan espacios. No vacía. Texto libre, sin lista cerrada. |
| `horas` | `float64` | Sí | Mayor a 0, múltiplo de 0,5 y como máximo 24, validado en ese orden. |

**En HTTP** (`POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, ruta pendiente de confirmar, ver P-6), el proyecto y el elemento van en la ruta y el resto en un cuerpo JSON, con los nombres en snake_case y la fecha en formato `AAAA-MM-DD`, como en HU-04:

```json
{
  "integrante": "Ana",
  "fecha": "2026-10-25",
  "actividad": "Desarrollo",
  "horas": 1.5
}
```

## 3. Salidas esperadas
- Se devuelve el registro creado con `ID` (único, autogenerado), `ProyectoID`, `ElementoID`, `Integrante` (el nombre tal como figura en el proyecto), `Fecha` (solo el día), `Actividad` recortada y `Horas`.
- El registro queda persistido y se puede recuperar por su `ID`.
- Si hay un error, no se guarda nada.
- **En HTTP:** `201` con el registro creado en JSON (`id`, `proyecto_id`, `elemento_id`, `integrante`, `fecha`, `actividad`, `horas`); `400` con el error si los datos del cuerpo son inválidos, incluido un integrante que no pertenece al proyecto (`ErrIntegranteNoEncontrado`); `404` solo si no existe el proyecto o el elemento, que son los datos que van en la ruta.

## 4. Reglas de negocio
- **RN-18.1:** Todo registro de esfuerzo pertenece a un elemento de backlog existente del proyecto indicado. Las tareas no existen como entidad propia: una tarea es un elemento más del backlog. Un elemento de otro proyecto se trata como inexistente.
- **RN-18.2:** El integrante es obligatorio (un nombre de solo espacios es vacío).
- **RN-18.3:** El integrante tiene que estar registrado en el proyecto (HU-03) al momento de registrar el esfuerzo. Si después se lo quita del proyecto, sus registros de esfuerzo se conservan.
- **RN-18.4:** La fecha es obligatoria y no puede ser futura. Hoy es válido. La comparación es por día, no por hora, y "hoy" es el día actual en la zona horaria del servidor. No se valida que la fecha esté dentro de las fechas del proyecto ni que el proyecto esté Activo.
- **RN-18.5:** La actividad es obligatoria (una de solo espacios es vacía) y es texto libre.
- **RN-18.6:** Las horas son obligatorias y mayores a 0.
- **RN-18.7:** Las horas se registran en múltiplos de media hora (0,5; 1; 1,5; …).
- **RN-18.8:** Un registro tiene como máximo 24 h. El tope es por registro: no se suman los otros registros del integrante en el mismo día.- **RN-18.9:** Se pueden registrar varios esfuerzos sobre el mismo elemento, del mismo integrante y del mismo día; cada uno es un registro independiente.
- **RN-18.10:** El `ID` es único y lo genera el sistema (entero autoincremental), igual que en HU-05.

## 5. Restricciones
- Las validaciones de campos (RN-18.2, 18.4 a 18.8) viven en el dominio (`internal/domain/registro_esfuerzo.go`, constructor `NuevoRegistroEsfuerzo`), sin dependencias externas.
- "Hoy" no se lee del reloj dentro del dominio: el caso de uso recibe un reloj (puerto o función `func() time.Time`) y le pasa la fecha actual al dominio, para que los tests sean deterministas.
- **Zona horaria:** "hoy" es el día actual en la zona horaria local del servidor (`time.Local`). La fecha recibida (`AAAA-MM-DD`) se interpreta como un día de esa misma zona, así ambas se comparan por año, mes y día sin corrimientos. En los tests, el reloj falso fija la zona.
- La existencia del proyecto y del elemento (RN-18.1) y la pertenencia del integrante (RN-18.3) se resuelven en el caso de uso `RegistrarEsfuerzo` (`internal/app/registrar_esfuerzo.go`) a través de puertos de repositorio.
- El registro guarda el nombre del integrante, no una referencia a la fila del integrante. En la base no hay clave foránea hacia la tabla de integrantes, para que quitar a un integrante (HU-03) no borre ni bloquee sus registros (RN-18.3).
- Las horas se reciben como `float64`; para evitar errores de redondeo se recomienda guardarlas internamente como cantidad de medias horas (`int`).
- Persistencia en SQLite (ADR-0002); `net/http` (ADR-0001) con la ruta propuesta `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, en lugar del `POST /esfuerzo` de la tarea 18.7 (pendiente de confirmar, ver P-6).
- Solo endpoint JSON: la vista web se hace en HU-30.
- Fuera de alcance: editar, eliminar o listar registros de esfuerzo, y totales por elemento o por integrante (los usan historias posteriores de métricas). No hay autenticación ni roles.

## 6. Casos límite
- Elemento sin ningún esfuerzo registrado todavía → el primer registro es válido.
- Horas = 0,5 → válido. Horas = 0 o negativas → error.
- Horas = 1,25 → error de múltiplo de media hora.
- Horas = 24 → válido; 24,5 → error.
- Horas = -1,25 → `ErrHorasInvalidas` (se valida primero que sean mayores a 0). Horas = 25,3 → `ErrHorasNoMultiploMediaHora` (el múltiplo se valida antes que el máximo).
- Dos registros del mismo integrante en el mismo día que suman más de 24 h (ej. 20 h + 8 h) → ambos válidos (RN-18.8).
- Fecha = hoy → válido. Fecha = mañana → error.
- Fecha anterior al inicio del proyecto, o posterior a su fin → válido (RN-18.4).
- Proyecto en un estado distinto de Activo → se puede registrar igual (RN-18.4).
- Integrante `" ana "` cuando en el proyecto está `"Ana"` → válido, se guarda como `"Ana"`.
- Integrante que existe en otro proyecto pero no en este → error.
- Elemento que existe pero en otro proyecto → error de elemento no encontrado.
- Integrante con registros que después se quita del proyecto → sus registros siguen existiendo; un registro nuevo suyo falla porque ya no es integrante.
- Dos registros iguales (mismo elemento, integrante, día, actividad y horas) → ambos válidos, con IDs distintos (RN-18.9).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404) |
| No existe el elemento, o es de otro proyecto | `domain.ErrElementoNoEncontrado` (HU-05, se reutiliza) | "elemento de backlog no encontrado" (HTTP 404) |
| Integrante vacío o solo espacios | `domain.ErrNombreIntegranteVacio` (HU-03, se reutiliza) | "el nombre del integrante es obligatorio" (HTTP 400) |
| Fecha no informada (valor cero) | `ErrFechaEsfuerzoObligatoria` | "la fecha es obligatoria" (HTTP 400) |
| Fecha posterior a hoy | `ErrFechaEsfuerzoFutura` | "la fecha no puede ser futura" (HTTP 400) |
| Actividad vacía o solo espacios | `ErrActividadVacia` | "la actividad es obligatoria" (HTTP 400) |
| Horas 0 o negativas | `ErrHorasInvalidas` | "las horas deben ser mayores a 0" (HTTP 400) |
| Horas que no son múltiplo de 0,5 | `ErrHorasNoMultiploMediaHora` | "las horas se registran en múltiplos de media hora" (HTTP 400) |
| Horas mayores a 24 | `ErrHorasExcedenMaximo` | "un registro no puede superar las 24 horas" (HTTP 400) |
| El integrante no pertenece al proyecto | `domain.ErrIntegranteNoEncontrado` (HU-03, se reutiliza) | "integrante no encontrado en el proyecto" (HTTP 400: el integrante viene en el cuerpo, no en la ruta, a diferencia de HU-03) |
| Cuerpo JSON mal formado o fecha con otro formato | — (adaptador HTTP) | HTTP 400; no se llama al caso de uso. |

Los mensajes de los errores reutilizados son los que ya están definidos en `internal/domain`; los nuevos siguen el mismo estilo (minúscula y sin punto final).

Si hay varios errores, se informa el primero en el orden de la tabla: primero se buscan el proyecto y el elemento, después el dominio valida los campos y por último el caso de uso verifica que el integrante pertenezca al proyecto. Dentro de las horas, el orden es: mayores a 0, después múltiplo de 0,5 y al final máximo 24 h.

## 8. Criterios de aceptación
- **CA-18.1:** Dado un proyecto con el elemento "Login" y la integrante "Ana", cuando registro 2 h de "Desarrollo" de Ana con fecha de hoy, entonces se crea un registro con ID asignado y los datos enviados. (RN-18.10)
- **CA-18.2:** Dado un elemento sin esfuerzo registrado, cuando registro 0,5 h, entonces se crea. (RN-18.6, RN-18.7)
- **CA-18.3:** Cuando las horas son 0 o negativas, entonces obtengo `ErrHorasInvalidas`. (RN-18.6)
- **CA-18.4:** Cuando la actividad está vacía o tiene solo espacios, entonces obtengo `ErrActividadVacia`. (RN-18.5)
- **CA-18.5:** Cuando el integrante está vacío, entonces obtengo `ErrNombreIntegranteVacio`. (RN-18.2)
- **CA-18.6:** Cuando la fecha no se informa, entonces obtengo `ErrFechaEsfuerzoObligatoria`; cuando es mañana, `ErrFechaEsfuerzoFutura`; cuando es hoy, se crea. (RN-18.4)
- **CA-18.7:** Dado un ID de proyecto o de elemento que no existe, o un elemento de otro proyecto, cuando registro esfuerzo, entonces obtengo `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado`. (RN-18.1)
- **CA-18.8:** Dado un proyecto sin la integrante "Caro", cuando registro esfuerzo de Caro, entonces obtengo `ErrIntegranteNoEncontrado`. (RN-18.3)
- **CA-18.9:** Cuando las horas son 1,25, entonces obtengo `ErrHorasNoMultiploMediaHora`. (RN-18.7)
- **CA-18.10:** Cuando registro 24 h, entonces se crea; cuando registro 24,5 h, obtengo `ErrHorasExcedenMaximo`. (RN-18.8)
- **CA-18.11:** Dado que Ana tiene un registro de 20 h hoy, cuando registro 8 h más para hoy, entonces se crea, porque el tope no suma los registros del día. (RN-18.8, RN-18.9)
- **CA-18.12:** Dados dos registros con los mismos datos, entonces ambos se crean con IDs distintos y cada uno se recupera por su ID. (RN-18.9, RN-18.10)
- **CA-18.13:** Dado que Ana tiene esfuerzo registrado, cuando se la quita del proyecto (HU-03), entonces sus registros se siguen recuperando por su ID. (RN-18.3)
- **CA-18.14:** Dado un proyecto que no está Activo y una fecha anterior a su fecha de inicio, cuando registro esfuerzo con esa fecha, entonces se crea. (RN-18.4)
- **CA-18.15:** Dado un pedido válido a `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo` (ruta sujeta a P-6), entonces la respuesta es `201` con el registro en JSON; con datos inválidos, `400`; con un integrante que no pertenece al proyecto, `400`; con proyecto o elemento inexistente, `404`.

En todos los casos de error, no se crea ningún registro.

Los escenarios mínimos de la tarea 18.2 (válido, horas negativas, sin actividad) corresponden a CA-18.1, CA-18.3 y CA-18.4.

## Decisiones tomadas
Las acordó el equipo en el chat el 2026-10-25, a partir del borrador de esta spec:
- **P-1:** el esfuerzo se registra solo sobre elementos del backlog (HU-05); las tareas no existen como entidad.
- **P-2:** quien registra tiene que ser integrante del proyecto (HU-03); si después se lo quita, sus registros se conservan.
- **P-3:** la actividad es texto libre y obligatoria.
- **P-4:** se aceptan medias horas, con un tope de 24 h por registro. No se suman los registros del día.
- **P-5:** la fecha solo no puede ser futura; no se valida contra las fechas del proyecto ni contra su estado. Confirmada por Matías.
- HU-18 es solo un endpoint JSON, por el plan del Sprint 1; la pantalla se suma con HU-30. No hubo objeciones.

## Preguntas abiertas
- **P-6 — Ruta HTTP:** se propone `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, en lugar del `POST /esfuerzo` del plan (tarea 18.7), para seguir el estilo de HU-05. Pendiente de confirmar con Matías.

## Historial
- 2026-10-25: creada (borrador), con las decisiones P-1 a P-4 todavía sin confirmar por el equipo.
- 2026-10-25: el equipo resuelve P-1 a P-6. El tope de horas pasa a ser por registro (24 h) en vez de diario: se quitan `ErrTopeHorasDiarias` y la suma por día, y se agrega `ErrHorasExcedenMaximo` en el dominio. Se agregan la conservación de registros al quitar un integrante (RN-18.3, CA-18.13), la falta de validación contra las fechas y el estado del proyecto (RN-18.4, CA-18.14), la ruta y el contrato JSON (CA-18.15), y la aclaración de que la vista web va en HU-30.
- 2026-10-25: ajustes: `ErrIntegranteNoEncontrado` responde HTTP 400 (el integrante va en el cuerpo, no en la ruta) en las secciones 3 y 7 y en CA-18.15; los mensajes de los errores reutilizados se copian de `internal/domain` y los nuevos siguen su estilo (minúscula, sin punto final); "hoy" se define en la zona horaria del servidor (RN-18.4 y Restricciones); se fija el orden de validación de las horas (> 0, múltiplo de 0,5, máximo 24). "Decisiones tomadas" registra que se acordaron en el chat y que Matías confirmó P-5. P-6 (ruta) vuelve a "Preguntas abiertas", pendiente de confirmar con Matías.
