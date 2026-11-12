# HU-18: Registrar esfuerzo en una historia o tarea

**Historia:** Como Product Builder, quiero registrar integrante, fecha, actividad y horas trabajadas en una historia o tarea, para dejar constancia del esfuerzo real.
**Issue:** #29  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Implementada

## 1. Objetivo
Dejar constancia de las horas que cada integrante trabajó realmente sobre un elemento del backlog: quién, qué día, en qué actividad y cuánto tiempo. Es el dato base del esfuerzo real, que después usan las métricas (estimado vs. real), el dashboard y los reportes. Reutiliza los elementos de backlog de HU-05 y los integrantes de HU-03. Se expone con un endpoint JSON y con una pantalla (sección 9).

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
- **En la pantalla** del elemento: ver la sección 9.

## 4. Reglas de negocio
- **RN-18.1:** Todo registro de esfuerzo pertenece a un elemento de backlog existente del proyecto indicado. Las tareas no existen como entidad propia: una tarea es un elemento más del backlog. Un elemento de otro proyecto se trata como inexistente.
- **RN-18.2:** El integrante es obligatorio (un nombre de solo espacios es vacío).
- **RN-18.3:** El integrante tiene que estar registrado en el proyecto (HU-03) al momento de registrar el esfuerzo. Si después se lo quita del proyecto, sus registros de esfuerzo se conservan.
- **RN-18.4:** La fecha es obligatoria y no puede ser futura. Hoy es válido. La comparación es por día, no por hora, y "hoy" es el día actual en la zona horaria del servidor. No se valida que la fecha esté dentro de las fechas del proyecto ni que el proyecto esté Activo.
- **RN-18.5:** La actividad es obligatoria (una de solo espacios es vacía) y es texto libre.
- **RN-18.6:** Las horas son obligatorias y mayores a 0.
- **RN-18.7:** Las horas se registran en múltiplos de media hora (0,5; 1; 1,5; …).
- **RN-18.8:** Un registro tiene como máximo 24 h. El tope es por registro: no se suman los otros registros del integrante en el mismo día.
- **RN-18.9:** Se pueden registrar varios esfuerzos sobre el mismo elemento, del mismo integrante y del mismo día; cada uno es un registro independiente.
- **RN-18.10:** El `ID` es único y lo genera el sistema (entero autoincremental), igual que en HU-05.

## 5. Restricciones
- Las validaciones de campos (RN-18.2, 18.4 a 18.8) viven en el dominio (`internal/domain/registro_esfuerzo.go`, constructor `NuevoRegistroEsfuerzo`), sin dependencias externas.
- "Hoy" no se lee del reloj dentro del dominio: el caso de uso recibe un reloj (puerto o función `func() time.Time`) y le pasa la fecha actual al dominio, para que los tests sean deterministas.
- **Zona horaria:** "hoy" es el día actual en la zona horaria local del servidor (`time.Local`). La fecha recibida (`AAAA-MM-DD`) se interpreta como un día de esa misma zona, así ambas se comparan por año, mes y día sin corrimientos. En los tests, el reloj falso fija la zona.
- La existencia del proyecto y del elemento (RN-18.1) y la pertenencia del integrante (RN-18.3) se resuelven en el caso de uso `RegistrarEsfuerzo` (`internal/app/registrar_esfuerzo.go`) a través de puertos de repositorio.
- El registro guarda el nombre del integrante, no una referencia a la fila del integrante. En la base no hay clave foránea hacia la tabla de integrantes, para que quitar a un integrante (HU-03) no borre ni bloquee sus registros (RN-18.3).
- Las horas se reciben como `float64`; para evitar errores de redondeo se recomienda guardarlas internamente como cantidad de medias horas (`int`).
- Persistencia en SQLite (ADR-0002); `net/http` (ADR-0001) con la ruta propuesta `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, en lugar del `POST /esfuerzo` de la tarea 18.7 (pendiente de confirmar, ver P-6).
- La pantalla se especifica en la sección 9 (ADR-0003) y usa el mismo caso de uso que el endpoint JSON.
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
| Se busca por ID un registro que no existe | `domain.ErrRegistroEsfuerzoNoEncontrado` | "registro de esfuerzo no encontrado". Lo devuelve la búsqueda por ID de la persistencia (sección 3: el registro se recupera por su ID). No es parte de registrar y no tiene endpoint en HU-18. |

Los mensajes de los errores reutilizados son los que ya están definidos en `internal/domain`; los nuevos siguen el mismo estilo (minúscula y sin punto final).

Si hay varios errores al registrar, se informa el primero en el orden de la tabla (la última fila es de la búsqueda, no de registrar): primero se buscan el proyecto y el elemento, después el dominio valida los campos y por último el caso de uso verifica que el integrante pertenezca al proyecto. Dentro de las horas, el orden es: mayores a 0, después múltiplo de 0,5 y al final máximo 24 h.

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
- **CA-18.16:** Dada la pantalla de un elemento de un proyecto con las integrantes "Ana" y "Caro", entonces la sección "Registrar esfuerzo" muestra un desplegable con esas dos integrantes y sin ninguna elegida, la fecha de hoy, y la actividad y las horas vacías. No muestra registros ni totales de esfuerzo. (Sección 9)
- **CA-18.17:** Dada la lista del backlog, entonces cada elemento tiene el enlace "Registrar esfuerzo", que lleva a la sección del formulario en la pantalla de ese elemento. (Sección 9)
- **CA-18.18:** Dado el formulario con datos válidos, cuando lo envío, entonces el caso de uso recibe esos datos una sola vez y la respuesta es `303` a `/web/proyectos/{id}/backlog/{elementoID}?esfuerzo=registrado#esfuerzo`. Esa pantalla muestra "Esfuerzo registrado." y el formulario con sus valores iniciales. Si la recargo, no se crea otro registro. (Sección 9)
- **CA-18.19:** Dado el formulario con un dato que el dominio rechaza (integrante vacío, fecha vacía o futura, actividad vacía, horas 0, 1,25 o 24,5), cuando lo envío, entonces la respuesta es `400`, la misma pantalla muestra en esa sección el mensaje de la sección 7, sin el texto envuelto por el caso de uso, y conserva lo que tipeé. El formulario de edición muestra los valores guardados del elemento. (RN-18.2, RN-18.4 a RN-18.8, sección 9)
- **CA-18.20:** Dada una fecha con otro formato o unas horas que no son un número, cuando envío el formulario, entonces la respuesta es `400` con el mensaje de formato de la sección 9, conserva lo que tipeé y no se llama al caso de uso. Con las horas vacías, el mensaje es el de `ErrHorasInvalidas`, igual que en JSON. (Sección 9)
- **CA-18.21:** Dado que a "Caro" la quitaron del proyecto después de abrir la pantalla, cuando envío el formulario con Caro, entonces la respuesta es `400` con "integrante no encontrado en el proyecto", los demás campos conservan lo que tipeé y Caro ya no aparece en el desplegable. (RN-18.3, sección 9)
- **CA-18.22:** Dado un proyecto o un elemento inexistente, o un elemento de otro proyecto, cuando envío el formulario, entonces la respuesta es `404` y la página de error muestra "proyecto no encontrado" o "elemento de backlog no encontrado". Con un ID que no es un entero positivo, la respuesta es `400` con la página de error. (RN-18.1, sección 9)
- **CA-18.23:** Dado un envío desde otro sitio, entonces la respuesta es `403` y no se llama al caso de uso. Dada una falla interna, entonces la respuesta es `500` con el mensaje genérico, sin el detalle. (Sección 9)
- **CA-18.24:** Dado un proyecto sin integrantes, entonces la sección avisa que no hay integrantes y enlaza al detalle del proyecto. (Sección 9)
- **CA-18.25:** Dado un error al editar (`400`) o al eliminar (`409`) el elemento, entonces la pantalla sigue mostrando ese mensaje arriba, como en HU-06, y también la sección "Registrar esfuerzo" con sus valores iniciales. (Sección 9)

En todos los casos de error, no se crea ningún registro.

Los escenarios mínimos de la tarea 18.2 (válido, horas negativas, sin actividad) corresponden a CA-18.1, CA-18.3 y CA-18.4.

CA-18.16 a CA-18.25 se cubren con tests del adaptador (`httptest`), como las pantallas de HU-10 y HU-11: las reglas ya están cubiertas por los escenarios BDD de CA-18.1 a CA-18.14.

## 9. Pantalla
Agregada por el issue #289. HU-18 se entregó solo con el endpoint JSON y HU-30 no incluyó la pantalla, que la DoD del proyecto exige (ADR-0003).

**Dónde vive.** En la pantalla del elemento de HU-06 (`GET /web/proyectos/{id}/backlog/{elementoID}`, plantilla `elemento.html`), en una sección nueva "Registrar esfuerzo" con `id="esfuerzo"`, entre el formulario de edición y "Eliminar elemento". No hay plantilla nueva. La lista del backlog (HU-07) suma en cada fila el enlace "Registrar esfuerzo" a `/web/proyectos/{id}/backlog/{elementoID}#esfuerzo`, porque "Editar" no sugiere que ahí se cargan horas.

**Qué muestra.**

| Campo (`name`) | Control | Valor inicial |
|---|---|---|
| `integrante` | `<select>` obligatorio con los integrantes del proyecto (HU-03), más una primera opción vacía "Elegí un integrante" | Sin elegir |
| `fecha` | `<input type="date">` obligatorio | Hoy, en la zona del servidor (`AAAA-MM-DD`) |
| `actividad` | `<input type="text">` obligatorio | Vacío |
| `horas` | `<input type="number" step="0.5">` obligatorio | Vacío |

- `step="0.5"` hace falta porque, sin él, el navegador no acepta decimales. No se agregan `min`, `max` ni tope de fecha: los límites siguen en un solo lugar, el dominio. El navegador envía el número con punto aunque lo muestre con coma.
- No se listan registros ni se muestran totales: siguen fuera de alcance (sección 5). Las horas reales del elemento se suman con HU-19.
- Sin integrantes en el proyecto, el formulario se muestra igual, con un aviso y un enlace a `/web/proyectos/{id}` para agregarlos. No se puede elegir a nadie, así que un envío termina en `ErrNombreIntegranteVacio`.
- El formulario se muestra con cualquier estado del proyecto (RN-18.4) y del elemento: ninguna regla lo limita y el caso de uso no los mira.

**Acción.** `POST /web/proyectos/{id}/backlog/{elementoID}/esfuerzo`, con el `action` terminado en `#esfuerzo` para que, si vuelve con un error, el navegador muestre la sección y no el principio de la página.

- Se registra con `p.formulario`, como las demás acciones: chequeo `mismoOrigen` (`403`) y cuerpo de hasta 64 KiB (`413`).
- Llama a `CasosDeUso.RegistrarEsfuerzo`, el mismo caso de uso del endpoint JSON.
- Si sale bien: POST-redirect-GET con `303` a `/web/proyectos/{id}/backlog/{elementoID}?esfuerzo=registrado#esfuerzo`. El `GET` muestra "Esfuerzo registrado." (`class="exito" role="status"`, como en HU-20) cuando recibe ese parámetro. No hay sesiones para guardar un aviso. Escribir el parámetro a mano solo muestra el texto, sin ningún efecto sobre los datos.

**Lectura del formulario** (adaptador, antes de llamar al caso de uso, como la fila "cuerpo mal formado" de la sección 7):

| Dato | Lectura | Si no se puede leer |
|---|---|---|
| `fecha` | `fechaDelPedido`, que la interpreta en `time.Local` (sección 5). No se usa `fechaDeForm` de las pantallas de proyecto, que la interpreta en UTC. Vacía es la fecha cero: responde `ErrFechaEsfuerzoObligatoria`. | "la fecha debe tener el formato AAAA-MM-DD" (`400`) |
| `horas` | `strconv.ParseFloat`, sin los espacios de alrededor. Vacías valen 0: responden `ErrHorasInvalidas`, igual que en JSON cuando faltan. | "las horas deben ser un número" (`400`) |

Si fallan las dos lecturas, se informa primero la fecha. `errFormulario` lleva estos mensajes, como "story_points debe ser un número entero" en HU-06. `NaN` e `Inf` se pueden leer y los rechaza el dominio.

**Errores.** Se muestra siempre el mensaje del error del dominio de la sección 7, nunca el texto envuelto por el caso de uso. Se comparan con `errors.Is` contra las listas que ya usa el handler JSON (`erroresNoEncontradoEsfuerzo` y `erroresDeValidacionEsfuerzo`): `mensajeDeValidacion` no incluye los errores propios de HU-18 y respondería `500`.

| Situación | Código | Qué se ve |
|---|---|---|
| Error de validación de la sección 7, incluido `ErrIntegranteNoEncontrado`, o de lectura del formulario | `400` | La pantalla del elemento con el mensaje en la sección "Registrar esfuerzo" (`error_form`) |
| `ErrProyectoNoEncontrado` o `ErrElementoNoEncontrado` | `404` | `error.html` con el mensaje |
| ID de la ruta inválido | `400` | `error.html` (`idsElemento`) |
| Envío desde otro sitio | `403` | `error.html` |
| Formulario demasiado grande | `413` | `error.html` |
| Cualquier otro error | `500` | `error.html` con el mensaje genérico (`fallaInterna`), sin el detalle |

**Valores al volver con un error.** `actividad`, `horas` y `fecha` conservan exactamente lo tipeado, aunque la fecha haya quedado vacía. `integrante` conserva la opción elegida si sigue en el proyecto; si no, queda sin elegir y el mensaje lo explica. El formulario de edición se rearma con los valores guardados del elemento, porque cada formulario se envía por separado: lo que se haya tipeado ahí sin guardar se pierde, igual que al usar "Eliminar".

**Cambios en el adaptador.**
- Archivo nuevo `pantallas_esfuerzo.go`, con `registrarPantallasEsfuerzo(mux)`, sumado a `RegistrarPantallas` y con la ruta sumada a `router_test.go`.
- `datosEdicionElemento` suma los integrantes, el formulario de esfuerzo, su error propio (separado de `Error`, que sigue siendo el de editar y eliminar) y el aviso de registrado.
- `mostrarElemento` también consulta el proyecto con `ConsultarProyecto`, porque lo usan el `GET` y los errores de editar y eliminar (CA-18.25). `errorDeElemento` suma `ErrProyectoNoEncontrado`, que hoy terminaría en `500`.
- `pantallas` suma el campo `reloj func() time.Time`, que `nuevasPantallas` inicializa con `time.Now`, solo para el valor inicial de la fecha. Los tests lo fijan y registran la ruta con `registrarPantallasElemento`. La validación de "hoy" sigue en el caso de uso, con su propio reloj.

**Qué no cambia.** Las reglas RN-18.1 a RN-18.10, el dominio, el caso de uso `RegistrarEsfuerzo`, la persistencia y el endpoint JSON `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo` con sus códigos.

## Decisiones tomadas
Las acordó el equipo en el chat el 2026-10-25, a partir del borrador de esta spec:
- **P-1:** el esfuerzo se registra solo sobre elementos del backlog (HU-05); las tareas no existen como entidad.
- **P-2:** quien registra tiene que ser integrante del proyecto (HU-03); si después se lo quita, sus registros se conservan.
- **P-3:** la actividad es texto libre y obligatoria.
- **P-4:** se aceptan medias horas, con un tope de 24 h por registro. No se suman los registros del día.
- **P-5:** la fecha solo no puede ser futura; no se valida contra las fechas del proyecto ni contra su estado. Confirmada por Matías.
- HU-18 es solo un endpoint JSON, por el plan del Sprint 1; la pantalla se suma con HU-30. No hubo objeciones. (Reemplazada por P-7 a P-13: HU-30 no incluyó la pantalla.)

P-7 a P-13 las tomó el autor (Lucas Orellana) el 2026-11-01; quedan abiertas a objeción en el PR.
- **P-7:** el formulario vive en la pantalla del elemento y no en una propia. El esfuerzo es de un elemento, esa pantalla ya resuelve los IDs y los `404`, y es donde aparece "no se puede eliminar un elemento con esfuerzo registrado". Sumar a una pantalla existente es lo que hicieron HU-10 y HU-11, y HU-19 propone mostrar ahí las horas reales (su P-11). Una pantalla propia, como la de HU-20, evitaría tocar `mostrarElemento`, pero sumaría una plantilla y una ruta `GET` más y separaría la carga de lo demás del elemento.
- **P-8:** no se muestra ni el listado ni el total, y se respeta la sección 5, que deja fuera a los dos; se muestra solo la confirmación "Esfuerzo registrado." Es lo más chico que le dice a quien carga que el registro quedó. El total necesitaría un puerto, un caso de uso y una consulta SQL, que la spec de HU-19 (en su rama, en borrador) ya define como `LectorHorasReales`, para mostrarlo en esta misma pantalla: hacerlo acá duplicaría ese trabajo.
- **P-9:** la acción es `POST /web/proyectos/{id}/backlog/{elementoID}/esfuerzo`, con `303` al elemento. La confirmación viaja en el parámetro `?esfuerzo=registrado` porque no hay sesiones, y `#esfuerzo` lleva a la sección, que queda debajo del formulario de edición.
- **P-10:** los códigos son los de las otras pantallas: `400` en el mismo formulario, `404` con `error.html`, `403` por otro origen, `413` por tamaño y `500` genérico. `ErrIntegranteNoEncontrado` es `400`, como en JSON.
- **P-11:** la lectura del formulario sigue al endpoint JSON: la fecha se lee con `fechaDelPedido` y las horas vacías valen 0. "las horas deben ser un número" es el único mensaje nuevo y vive en el adaptador, porque es de formato, no una regla de negocio.
- **P-12:** en el navegador solo se valida lo que hace usable al formulario (`required` y `step="0.5"`). Sin `min`, `max` ni tope de fecha, los límites de RN-18.4, RN-18.6 y RN-18.8 quedan solo en el dominio y sus mensajes se ven en la pantalla.
- **P-13:** el enlace "Registrar esfuerzo" va en la lista del backlog. Queda afuera el detalle del sprint (HU-09), para no agrandar el cambio. Sin integrantes, el formulario se muestra igual, con un aviso y un enlace para agregarlos.

## Preguntas abiertas
- **P-6 — Ruta HTTP. Resuelta (2026-11-11, por lo que hace el código):** quedó `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, como se proponía. *Propuesta original:* se propone `POST /proyectos/{id}/backlog/{elementoID}/esfuerzo`, en lugar del `POST /esfuerzo` del plan (tarea 18.7), para seguir el estilo de HU-05. Pendiente de confirmar con Matías.
- **P-14. Resuelta (2026-11-11):** HU-19 entró en el Sprint 3 (PR #297), y la pantalla del elemento muestra las horas estimadas y reales. *Pregunta original:* ¿alcanza la confirmación si HU-19 no llega a la entrega? Sin el total, la pantalla no muestra lo ya cargado, y como no se pueden editar ni eliminar registros, un doble registro no se detecta ni se corrige desde la interfaz. **Propuesta:** aceptarlo y dejar el total para HU-19. Si HU-19 se cae del Sprint 3, se suma solo el total a esta pantalla, con el mismo puerto `LectorHorasReales`, como cambio de alcance de la sección 5.
- **P-15. Resuelta (2026-11-11, por lo que hace el código):** `errorDeElemento` reconoce `ErrProyectoNoEncontrado` (`pantallas_elemento.go`). *Pregunta original:* quién suma `ErrProyectoNoEncontrado` a `errorDeElemento`. HU-18 y HU-19 piden el mismo cambio, y las dos tocan `mostrarElemento`. **Propuesta:** que lo haga esta pantalla, que se implementa antes (lo propone P-13 de HU-19), y que HU-19 lo reciba al actualizarse con `main`.

## Historial
- 2026-10-25: creada (borrador), con las decisiones P-1 a P-4 todavía sin confirmar por el equipo.
- 2026-10-25: el equipo resuelve P-1 a P-6. El tope de horas pasa a ser por registro (24 h) en vez de diario: se quitan `ErrTopeHorasDiarias` y la suma por día, y se agrega `ErrHorasExcedenMaximo` en el dominio. Se agregan la conservación de registros al quitar un integrante (RN-18.3, CA-18.13), la falta de validación contra las fechas y el estado del proyecto (RN-18.4, CA-18.14), la ruta y el contrato JSON (CA-18.15), y la aclaración de que la vista web va en HU-30.
- 2026-10-25: ajustes: `ErrIntegranteNoEncontrado` responde HTTP 400 (el integrante va en el cuerpo, no en la ruta) en las secciones 3 y 7 y en CA-18.15; los mensajes de los errores reutilizados se copian de `internal/domain` y los nuevos siguen su estilo (minúscula, sin punto final); "hoy" se define en la zona horaria del servidor (RN-18.4 y Restricciones); se fija el orden de validación de las horas (> 0, múltiplo de 0,5, máximo 24). "Decisiones tomadas" registra que se acordaron en el chat y que Matías confirmó P-5. P-6 (ruta) vuelve a "Preguntas abiertas", pendiente de confirmar con Matías.
- 2026-10-25: correcciones de la revisión con go-reviewer (tarea 18.8): se separan RN-18.8 y RN-18.9, que habían quedado en la misma línea; se agrega a la sección 7 `ErrRegistroEsfuerzoNoEncontrado`, que devuelve la búsqueda por ID de la persistencia (sin endpoint en HU-18). En el código, una fecha cero en cualquier zona horaria y unas horas NaN se rechazan con `ErrFechaEsfuerzoObligatoria` y `ErrHorasInvalidas`, sin cambiar las reglas.
- 2026-11-01: se especifica la pantalla para registrar esfuerzo (issue #289): sección 9, CA-18.16 a CA-18.25, decisiones P-7 a P-13 del autor, abiertas a objeción en el PR, y preguntas P-14 y P-15. Se actualizan las secciones 1, 3 y 5, que decían que la pantalla iba en HU-30. No cambian las reglas RN-18.x ni el endpoint JSON. La spec sigue en borrador.
- 2026-11-11: implementada en el PR #104, ya mergeado, que aprobó Nicolás Andino. Pasa de Borrador a Implementada al ordenar la documentación para la entrega; las reglas y los criterios no cambian. Las preguntas abiertas quedan resueltas por lo que hace el código, a confirmar por Matías en el PR.
