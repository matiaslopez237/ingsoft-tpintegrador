# HU-27: Exportar el reporte a PDF

**Historia:** Como Product Owner, quiero exportar el reporte a un archivo PDF, para entregarlo como documento formal.
**Issue:** #38, sub-issue de la épica 9 (Reportes)  ·  **Sprint:** 4  ·  **Story points:** 5  ·  **Depende de:** HU-26
**Tareas:** #349 ADR · #350 Spec · #351 Escenarios BDD · #352 Steps de godog · #353 TDD del caso de uso · #354 Adaptador de PDF · #355 HTTP y botón de descarga · #356 Revisión y trazabilidad
**Estado de la spec:** Borrador

**Criterios de aceptación de la issue #38** (literales):
1. El PDF contiene las mismas secciones que el reporte en pantalla.
2. El archivo se puede descargar/guardar correctamente.

Los RN y CA de esta spec detallan esos dos criterios; la sección 9 dice cuáles cubren cada uno. Cada RN y CA está marcado como **[Esencial]** (hace falta para cumplir los criterios de la issue) o **[Deseable]** (calidad adicional). Si falta tiempo, los deseables son lo primero que se recorta (sección 9).

## 1. Objetivo
Permitir descargar como archivo PDF el reporte de un proyecto o de un sprint que arma HU-26, para entregarlo fuera del sistema. El PDF tiene que mostrar bien el español, paginar sin cortar mal las tablas y no fallar por un emoji en un nombre.

**Alcance real de esta historia:** HU-27 no decide qué contiene el reporte; eso es de HU-26, que todavía no tiene spec ni rama. Para no esperarla, esta spec propone desacoplar las dos historias con un **modelo neutral de documento** (`DocumentoPDF`): HU-27 genera un PDF a partir de cualquier documento con título, subtítulo, fecha de emisión y secciones con párrafos y tablas de texto, y HU-26 convierte su reporte a ese modelo. Así HU-27 se puede programar y probar completa con documentos de prueba. Todo lo que depende de HU-26 está marcado como **SUPUESTO** y repetido en "Preguntas abiertas".

La librería es `codeberg.org/go-pdf/fpdf` con la fuente DejaVu Sans Condensed embebida ([ADR-0006](../adr/0006-libreria-pdf.md), todavía `proposed`).

## 2. Entradas
Son dos descargas sin cuerpo.

**Reporte de un proyecto.** **En HTTP:** `GET /web/proyectos/{id}/reporte/pdf` (**SUPUESTO**: la ruta cuelga de la pantalla del reporte de HU-26, ver sección 3).

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |

**Reporte de un sprint.** **En HTTP:** `GET /web/proyectos/{id}/sprints/{sprintID}/reporte/pdf` (**SUPUESTO**, igual que la anterior).

| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| `proyectoID` | `int64` | Sí | Entero positivo. Debe corresponder a un proyecto existente (HU-01). |
| `sprintID` | `int64` | Sí | Entero positivo. Debe corresponder a un sprint de ese proyecto (RN-27.13). |

**Entrada interna: el documento.** El generador de PDF recibe un `domain.DocumentoPDF`, que arma HU-26 a partir de su reporte:

```go
type DocumentoPDF struct {
	Titulo    string     // obligatorio
	Subtitulo string     // opcional
	Emision   time.Time  // fecha y hora de emisión, con su huso horario (RN-27.10)
	Aviso     string     // aviso de reemplazos (RN-27.9): lo completa el caso de uso; HU-26 nunca lo completa ni lo lee
	Secciones []SeccionPDF
}

type SeccionPDF struct {
	Titulo  string       // obligatorio
	Bloques []BloquePDF  // en el orden en que se muestran
}

// Un bloque es un párrafo o una tabla, no los dos (RN-27.11). Sin ninguno de los dos es un párrafo
// vacío: no es un error y no se muestra.
type BloquePDF struct {
	Parrafo string
	Tabla   *TablaPDF
}

type TablaPDF struct {
	Columnas []string    // encabezados; al menos uno
	Filas    [][]string  // cada fila con tantas celdas como columnas
}
```

HU-26 completa todos los campos menos `Emision` y `Aviso`, que pone el caso de uso de HU-27 (RN-27.9, RN-27.10). El campo `Aviso` se agregó al escribir los steps de HU-27: RN-27.2 ubica el aviso en el documento, pero el modelo no tenía dónde guardarlo.

Todo es texto: los números, porcentajes y horas llegan ya formateados por HU-26, para que el PDF diga lo mismo que la pantalla (por ejemplo, `+25,0 %` como en HU-19).

## 3. Salidas esperadas

**El PDF.** Hoja A4 vertical (RN-27.14). En la primera página, en este orden (RN-27.2):
1. El título y, si hay, el subtítulo.
2. La línea de emisión: "Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3)." (RN-27.10).
3. Si se reemplazaron caracteres, el aviso de RN-27.9.
4. Las secciones, cada una con su título y sus bloques en orden.

Desde la segunda página, un encabezado con el título del documento. En todas, el pie "Página 2 de 5" (RN-27.3).

**La descarga.**
- `200` con el PDF completo en el cuerpo y estas cabeceras:
  - `Content-Type: application/pdf`
  - `Content-Disposition: attachment; filename="reporte-proyecto-3-2026-11-07.pdf"` (o `reporte-proyecto-3-sprint-12-2026-11-07.pdf` para un sprint; RN-27.12)
  - `Content-Length` con el tamaño del PDF
  - `X-Content-Type-Options: nosniff`
- Los errores (sección 7) responden con la página de error HTML que ya usan las pantallas (`error.html`, con su CSP), porque quien llega a esta ruta es el navegador al seguir el botón de descarga.
- La descarga no modifica nada.

No lleva `Content-Security-Policy`: esa cabecera se aplica a las respuestas HTML (`renderizar` en `pantallas.go`) y un PDF descargado no la usa (hallazgo del ADR-0006).

### Pantalla
Solo lectura, sin formularios y sin JavaScript (ADR-0003).

- **Dónde (SUPUESTO):** en la pantalla del reporte de HU-26. Se supone que esa pantalla vive en `/web/proyectos/{id}/reporte` para el proyecto y en `/web/proyectos/{id}/sprints/{sprintID}/reporte` para el sprint, y que la ruta del PDF es la misma con `/pdf` al final. Si HU-26 elige otras rutas, se ajustan las de esta spec sin cambiar las reglas.
- **Qué se agrega:** un enlace con estilo de botón (`<a class="boton" href="…/reporte/pdf">Descargar PDF</a>`), como los enlaces "Registrar defecto" y "Editar" de `backlog.html`. Es un enlace `GET` y no un formulario, porque la descarga no cambia nada (no necesita la comprobación de origen de los formularios).
- **Si la descarga falla:** el navegador muestra la página de error, con el mensaje de la sección 7 y el enlace "Volver a la lista de proyectos". Un `500` muestra "Ocurrió un error interno. Probá de nuevo más tarde.", sin el detalle.
- **Si se reemplazaron caracteres:** el aviso va dentro del PDF (RN-27.9). La pantalla no lo anticipa, porque tendría que generar el documento para saberlo.

## 4. Reglas de negocio
- **RN-27.1 (modelo neutral):** [Esencial] el PDF se genera solo a partir de un `DocumentoPDF`. HU-27 no conoce la estructura del reporte de HU-26: qué secciones hay y qué dice cada una lo decide HU-26 al armar el documento.
- **RN-27.2 (orden del contenido):** [Esencial] título, subtítulo (si no está vacío), línea de emisión, aviso de reemplazos (`DocumentoPDF.Aviso`, si no está vacío) y secciones, en ese orden. Dentro de una sección, los bloques salen en el orden del documento.
- **RN-27.3 (encabezado y pie):** [Deseable] desde la segunda página, un encabezado con el título del documento. En todas las páginas, el pie "Página k de N", donde N es la cantidad total de páginas del PDF.
- **RN-27.4 (saltos de página):** [Esencial; la última oración, Deseable] cuando el contenido no entra, sigue en la página siguiente. Los párrafos largos continúan en la página siguiente. El título de una sección no queda solo al pie: si no entra junto con la primera línea de su primer bloque (o con el encabezado y la primera fila, si es una tabla), pasa entero a la página siguiente.
- **RN-27.5 (tablas):** [Esencial; la repetición de encabezados, Deseable] el texto de cada celda se ajusta al ancho de su columna en varias líneas, y la fila toma la altura de su celda más alta. Una fila nunca se parte entre dos páginas: si no entra, pasa entera a la siguiente. En cada página donde la tabla continúa se repite la fila de encabezados.
- **RN-27.6 (largo de una celda):** [Deseable] una celda de tabla tiene como máximo 300 caracteres, contados después de RN-27.8. Si tiene más, se cortan los primeros 299 y se agrega "…". Con ese tope, toda fila entra en una página (sección 5), así que RN-27.5 siempre se puede cumplir. Los párrafos no tienen tope.
- **RN-27.7 (tabla sin filas):** [Esencial] se muestra la fila de encabezados y debajo el texto "Sin datos.". Una sección sin bloques muestra solo su título. Un documento sin secciones es válido: muestra título, subtítulo y emisión.
- **RN-27.8 (caracteres):** [Esencial] antes de generar, cada texto del documento (título, subtítulo, títulos de sección, párrafos, encabezados y celdas) se normaliza carácter por carácter:
  1. Los espacios en blanco de Unicode (`unicode.IsSpace`: tabulación, salto de línea, retorno de carro, espacio duro, etc.) pasan a ser un espacio común.
  2. Los caracteres invisibles de formato (categoría Unicode `Cf`, como el espacio de ancho cero, el guion opcional o el unidor de emojis, y los selectores de variante `U+FE00`–`U+FE0F`) se quitan.
  3. Los caracteres representables se dejan como están. Son: `U+0021`–`U+007E` (ASCII imprimible), `U+00A1`–`U+017F` (Latin-1 y Latín extendido A, salvo `U+00AD`, que es `Cf`), `– — ‘ ’ “ ” • …` (`U+2013`, `U+2014`, `U+2018`, `U+2019`, `U+201C`, `U+201D`, `U+2022`, `U+2026`), `€` (`U+20AC`), `← ↑ → ↓` (`U+2190`–`U+2193`) y `≤ ≥` (`U+2264`, `U+2265`).
  4. Cualquier otro (emojis, ideogramas, otras escrituras, caracteres de control y bytes que no son UTF-8 válido) se reemplaza por `?`. Varios seguidos, aunque entre ellos hubiera caracteres quitados en el paso 2, dan **un solo** `?`.
- **RN-27.9 (aviso de reemplazos):** [Deseable] si RN-27.8 insertó al menos un `?`, el caso de uso escribe el aviso en `DocumentoPDF.Aviso` y el PDF lo muestra después de la línea de emisión; si no hubo reemplazos, el campo queda vacío. HU-26 nunca completa ni lee ese campo. El texto es: "Aviso: 3 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?»." En singular: "Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?»." El número es la cantidad de `?` insertados. Los espacios convertidos (paso 1) y los invisibles quitados (paso 2) no cuentan, porque no cambian lo que se ve.
- **RN-27.10 (fecha de emisión):** [Esencial] la fija el caso de uso con su reloj al pedir la descarga, en la hora de Argentina (UTC-3, sin horario de verano), y la pasa en `DocumentoPDF.Emision`. Se muestra como `dd/mm/aaaa` y `HH:MM` en 24 horas. Es la única fecha que usa el generador: también es la de creación y modificación en los metadatos del PDF, de modo que el mismo documento con la misma emisión genera siempre los mismos bytes.
- **RN-27.11 (documento válido):** [Esencial] el generador rechaza con `domain.ErrDocumentoInvalido` un documento con título vacío (después de RN-27.8 y sin contar espacios), una sección con título vacío, un bloque con párrafo y tabla a la vez, una tabla sin columnas, o una fila con distinta cantidad de celdas que columnas. Son errores de programación de quien arma el documento (HU-26), no del usuario: se responden como falla interna (sección 7). Un bloque sin párrafo ni tabla **no** es un error: es un párrafo vacío y no se muestra nada (sección 6).
- **RN-27.12 (nombre del archivo):** [Esencial] `reporte-proyecto-{proyectoID}-{aaaa-mm-dd}.pdf` para un proyecto y `reporte-proyecto-{proyectoID}-sprint-{sprintID}-{aaaa-mm-dd}.pdf` para un sprint, con los IDs de la ruta y la fecha de emisión de RN-27.10. No incluye el nombre del proyecto ni del sprint: el nombre tiene solo minúsculas sin acentos, dígitos, `-` y `.`.
- **RN-27.13 (qué se puede descargar):** [Esencial] el proyecto tiene que existir y el sprint tiene que ser de ese proyecto; un sprint de otro proyecto se trata como inexistente (igual que RN-23.11). **SUPUESTO:** se puede descargar en cualquier estado del proyecto y del sprint, igual que HU-23 (RN-23.12); si HU-26 restringe algún caso, HU-27 responde lo mismo que la pantalla del reporte.
- **RN-27.14 (hoja):** [Esencial] A4 en orientación vertical.
- **RN-27.15 (archivo válido):** [Esencial] el cuerpo de una descarga exitosa es un PDF completo: empieza con `%PDF-`, termina con `%%EOF` (como mucho seguido de un salto de línea), se puede abrir y extraer su texto, y `Content-Length` coincide con su tamaño. Nunca se envía un PDF a medias (sección 5: se genera entero en memoria antes de responder).

## 5. Restricciones
- **Dominio** (`internal/domain`, sin dependencias externas). Son reglas de negocio y se desarrollan con TDD:
  - `documento_pdf.go`: los tipos `DocumentoPDF`, `SeccionPDF`, `BloquePDF` y `TablaPDF` de la sección 2, `ErrDocumentoInvalido` y `(DocumentoPDF) Validar() error` (RN-27.11).
  - `documento_pdf_normalizar.go`:
    - `NormalizarDocumento(d DocumentoPDF) (DocumentoPDF, int)`: aplica RN-27.8 a todos los textos y RN-27.6 a las celdas, y devuelve cuántos `?` insertó. Es una función pura: no modifica el documento que recibe, y deja `Emision` y `Aviso` como están.
    - `TextoAvisoReemplazos(n int) string`: el texto de RN-27.9, vacío con `n` en 0.
  - `documento_pdf_emision.go`: `HoraDeEmision(t time.Time) time.Time` (RN-27.10), y `NombreArchivoReporteProyecto` y `NombreArchivoReporteSprint` (RN-27.12).
- **Casos de uso** (`internal/app/exportar_reporte_pdf.go`): `ExportarReporteProyecto.Ejecutar(ctx, proyectoID)` y `ExportarReporteSprint.Ejecutar(ctx, proyectoID, sprintID)`. Devuelven `ArchivoPDF{Nombre string; Contenido []byte}`. **La existencia del proyecto y del sprint la verifica el caso de uso, no la fuente** (RN-27.13). Cada uno, en este orden, que es el de la tabla de la sección 7:
  1. busca el proyecto; si no existe, devuelve `ErrProyectoNoEncontrado`;
  2. solo el del sprint: busca el sprint y devuelve `ErrSprintNoEncontrado` si no existe o si es de otro proyecto;
  3. pide el documento a la fuente;
  4. le pone la emisión con su reloj (inyectado, igual que `RegistrarDefecto` de HU-20);
  5. lo normaliza, completa `Aviso` (vacío si no hubo reemplazos, aunque la fuente trajera uno; RN-27.9) y lo valida;
  6. llama al generador y arma el nombre del archivo.

  Con un error no sigue con los pasos siguientes y no devuelve ningún archivo. Los errores de los lectores salen tal cual, como en HU-23. Las fallas de la fuente y del generador se envuelven con `%w`. La fuente y el generador reciben el contexto del pedido.
- **Puertos** (definidos en `internal/app`). Cada caso de uso recibe solo los que usa:

  | Caso de uso | Puertos |
  |---|---|
  | `ExportarReporteProyecto` | `LectorProyectos`, `FuenteDocumentoProyecto`, `GeneradorPDF` y el reloj |
  | `ExportarReporteSprint` | `LectorProyectos`, `LectorSprintPorID`, `FuenteDocumentoSprint`, `GeneradorPDF` y el reloj |

  ```go
  // Se reutilizan: LectorProyectos (HU-04) y LectorSprintPorID (HU-23).

  // FuenteDocumentoProyecto y FuenteDocumentoSprint las implementa HU-26 (SUPUESTO): convierten su
  // reporte al modelo neutral. Tienen un método cada una.
  type FuenteDocumentoProyecto interface {
  	DocumentoDeProyecto(ctx context.Context, proyectoID int64) (domain.DocumentoPDF, error)
  }
  type FuenteDocumentoSprint interface {
  	DocumentoDeSprint(ctx context.Context, proyectoID, sprintID int64) (domain.DocumentoPDF, error)
  }

  // GeneradorPDF lo implementa el adaptador internal/adapters/pdf con fpdf (ADR-0006).
  type GeneradorPDF interface {
  	Generar(ctx context.Context, doc domain.DocumentoPDF) ([]byte, error)
  }
  ```

  Hasta que exista HU-26, los tests y los steps de godog usan una fuente de prueba. En los steps, la fuente en memoria también hace de lector de proyectos y de sprints, con los mismos datos.
- **Huso horario:** se usa `time.FixedZone("UTC-3", -3*60*60)` y no `time.LoadLocation("America/Argentina/Buenos_Aires")`. Esta última necesita la base de husos horarios del sistema o embeber `time/tzdata`, y Render corre en UTC. Argentina no usa horario de verano desde 2009.
- **Adaptador de PDF** (`internal/adapters/pdf`): es el único paquete que importa `fpdf`. Las reglas de uso de la librería son las del ADR-0006: fuente embebida con `go:embed`, fechas fijas con `SetCreationDate` y `SetModificationDate`, `SetCatalogSort(true)`, sin las funciones `SetDefault*`, y devolver siempre el error de `Output`. La tabla (RN-27.5) la arma una función propia del adaptador; el ancho de cada columna lo decide el adaptador y no es regla de negocio. El adaptador recibe el documento ya normalizado, pero si `fpdf` igual devuelve un error, lo devuelve (se responde como falla interna).
- **Por qué alcanza el tope de RN-27.6:** en A4 vertical quedan unos 260 mm de alto útil. Una celda de 300 caracteres en una columna angosta, de unos 20 mm (unos 13 caracteres por línea), ocupa unas 23 líneas de unos 4,5 mm, es decir, unos 105 mm. Se comprueba en el test de la celda más larga con la columna más angosta que permita el adaptador.
- **HTTP:** el handler (`pantallas_reporte_pdf.go`) genera el PDF completo en memoria antes de escribir las cabeceras, igual que `renderizar`, para poder responder un `500` limpio si algo falla. Las dos rutas se registran con su propia `RegistrarExportarReportePDF(mux, proyecto, sprint)`, en el mismo archivo, como las demás historias; esa función se suma a `NuevoRouter` y los dos casos de uso a `CasosDeUso` en `router.go`, y las rutas a `router_test.go`. Reutiliza `proyectoIDDesdeRuta`, `sprintIDDesdeRuta`, `renderizar` y `fallaInterna`.
- **Concurrencia:** cada pedido genera su propio documento sin estado global compartido (ADR-0006: 32 generaciones simultáneas dieron el mismo resultado).
- **Dependencias nuevas:** `codeberg.org/go-pdf/fpdf` y, solo para los tests, `github.com/ledongthuc/pdf`. Se suman a `go.mod` con el adaptador, no con esta spec. La licencia de DejaVu va a `docs/third-party-licenses/` en el mismo cambio.
- **Sin persistencia nueva:** no hay tablas, columnas ni migraciones.
- **Fuera de alcance:**
  - **Para HU-26:** el contenido del reporte, su pantalla y la conversión a `DocumentoPDF`.
  - **Para HU-25 o una historia posterior:** gráficos en el PDF (ver "Preguntas abiertas").
  - **Sin historia:** HTML a PDF, imágenes, marcadores, PDF/A, firma digital, otras escrituras además de las de RN-27.8, elegir la orientación desde la pantalla, una ruta JSON de descarga y controlar quién descarga (no hay usuarios ni sesiones).

### Cómo se prueba
- **Dominio y casos de uso:** tests de tabla para RN-27.6, RN-27.8, RN-27.9, RN-27.10, RN-27.11 y RN-27.12, con un reloj fijo y dobles de la fuente y del generador.
- **Adaptador de PDF:** genera el PDF y lo lee con `ledongthuc/pdf`. Comprueba **las dos cosas**, porque el texto extraído no prueba lo que se ve (ADR-0006: ese extractor leyó 200 filas de un PDF que en pantalla mostraba 38):
  - el **texto de cada página**: el encabezado de la tabla en cada página donde continúa, el pie "Página k de N" y que las celdas de una misma fila estén en la misma página;
  - la **cantidad de páginas**, que tiene que coincidir con la N del pie.
- **Lista de representables:** un test genera un PDF con todos los caracteres de RN-27.8, paso 3, y comprueba que `fpdf` no devuelve error y que el texto extraído es idéntico. El ADR-0006 verificó los glifos de una parte de la lista; la revisión `@manual` cubre el resto.
- **Normalizar el texto extraído:** el extractor corta donde el PDF cambia de línea, así que se normalizan los espacios o se buscan fragmentos cortos.
- **Determinismo:** con la misma emisión fija, dos generaciones dan los mismos bytes. No se usa `time.Now` en ningún test.
- **HTTP:** `httptest` con dobles de los casos de uso: cabeceras, cuerpo que empieza con `%PDF-`, códigos de error y que el `500` no muestre el detalle.
- **Revisión visual:** un escenario `@manual` (CA-27.21), porque ningún test automático mira el dibujo.

## 6. Casos límite
- Celda de exactamente 300 caracteres → queda igual. Con 301 → 299 caracteres y "…" (RN-27.6).
- Celda de 300 caracteres en la columna más angosta → la fila entra en una página (sección 5).
- Fila que no entra al final de una página → pasa entera a la siguiente, con el encabezado repetido (RN-27.5).
- Tabla que termina justo al final de una página → no aparece una página siguiente con solo el encabezado.
- Título de sección en la última línea libre de una página → pasa a la siguiente junto con su primer bloque (RN-27.4).
- Párrafo de miles de caracteres → continúa en las páginas siguientes y se extrae completo.
- Tabla sin filas → encabezados y "Sin datos." (RN-27.7).
- Sección sin bloques → solo su título. Documento sin secciones → una página con título, subtítulo y emisión.
- Subtítulo vacío → no se muestra ni deja una línea en blanco.
- Párrafo vacío (un bloque sin párrafo ni tabla) → no es un error, no se muestra nada y el documento se exporta (RN-27.11).
- Título hecho solo de espacios o de caracteres invisibles → documento inválido (RN-27.11).
- `😀😀` → `?`, cuenta 1. `👨‍👩‍👧` (tres emojis unidos con `U+200D`) → `?`, cuenta 1. `❤️` (`U+2764` + `U+FE0F`) → `?`, cuenta 1. `😀 😀` → `? ?`, cuenta 2 (RN-27.8).
- `中文` → `?`. Una palabra en griego o en cirílico → `?`, aunque la fuente tenga esos glifos (RN-27.8, paso 3: solo lo que se probó).
- Tabulación o salto de línea dentro de una celda → espacio, sin aviso. Espacio de ancho cero → se quita, sin aviso.
- `U+0007` (campana) u otro carácter de control → `?`, con aviso.
- Bytes que no son UTF-8 válido → `?`, con aviso.
- Texto con `?` propios (por ejemplo, "¿cumplida?") → quedan igual y no cuentan en el aviso.
- Emisión el 2026-11-08 a las 02:30 UTC → "07/11/2026 a las 23:30" y el archivo lleva `2026-11-07` (cambio de día por el huso horario).
- Proyecto con nombre "Ñandú / Q4: \"final\"" → el nombre del archivo no cambia, porque no lo incluye (RN-27.12).
- ID `0`, negativo, `abc` o `anteriores` → `400`, sin llamar al caso de uso.
- Sprint de otro proyecto → `404`, sprint no encontrado (RN-27.13).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| ID de proyecto que no es un entero positivo | — (adaptador HTTP) | "el id del proyecto debe ser un número entero positivo" (HTTP 400, página de error); no se llama al caso de uso. |
| ID de sprint que no es un entero positivo | — (adaptador HTTP) | "el id del sprint debe ser un número entero positivo" (HTTP 400, página de error); no se llama al caso de uso. |
| No existe el proyecto | `domain.ErrProyectoNoEncontrado` (HU-01, se reutiliza) | "proyecto no encontrado" (HTTP 404, página de error) |
| No existe el sprint, o es de otro proyecto | `domain.ErrSprintNoEncontrado` (HU-09, se reutiliza) | "sprint no encontrado" (HTTP 404, página de error) |
| Documento inválido | `domain.ErrDocumentoInvalido` (nuevo) | HTTP 500, "Ocurrió un error interno. Probá de nuevo más tarde."; el detalle va solo al log. |
| Falla al armar el reporte (fuente) o al generar el PDF | — (se envuelven con `%w`) | HTTP 500, el mismo mensaje genérico; el detalle va solo al log. No se envía un PDF parcial. |

Si hay varios errores, se informa el primero en el orden de la tabla. **SUPUESTO:** si HU-26 define errores propios (por ejemplo, un reporte que no se puede armar para un sprint abierto), HU-27 los muestra igual que la pantalla del reporte; se agregan a esta tabla cuando HU-26 tenga spec.

Un emoji, un ideograma o un carácter de control **no son errores**: se reemplazan y se avisa (RN-27.8, RN-27.9).

## 8. Criterios de aceptación
- **CA-27.1 (documento normal):** [Esencial] Dado un documento con título, subtítulo, emisión 2026-11-07 14:05 UTC-3 y dos secciones (un párrafo y una tabla de 3 columnas y 3 filas), cuando se genera el PDF, entonces tiene 1 página, el texto extraído incluye en orden el título, el subtítulo, "Emitido el 07/11/2026 a las 14:05 (hora de Argentina, UTC-3).", los títulos de sección, el párrafo, los encabezados y las celdas, y el pie dice "Página 1 de 1". (RN-27.1, RN-27.2, RN-27.3, RN-27.10)
- **CA-27.2 (acentos y ñ):** [Esencial] Dado un título "Proyecto Ñandú: planificación ¿cumplida? ¡Sí!" y celdas con `áéíóúüñÁÉÍÓÚÜÑ€°ªº–—“”…→≤≥%`, cuando se genera, entonces el texto extraído es idéntico y no hay aviso. (RN-27.8, RN-27.9)
- **CA-27.3 (tabla larga):** [Esencial] Dada una tabla de 200 filas, cuando se genera, entonces el PDF tiene más de una página, la cantidad de páginas coincide con la N de "Página k de N", cada página donde sigue la tabla tiene la fila de encabezados, las 200 filas se extraen y las celdas de cada fila están en la misma página. (RN-27.3, RN-27.4, RN-27.5)
- **CA-27.4 (párrafo largo):** [Esencial] Dado un párrafo de 10 000 caracteres, cuando se genera, entonces ocupa más de una página, todas las páginas tienen el pie con la numeración y el texto se extrae completo. (RN-27.3, RN-27.4)
- **CA-27.5 (título de sección al pie):** [Deseable] Dado un documento en el que el título de una sección cae en la última línea libre de una página, cuando se genera, entonces ese título aparece en la página siguiente junto con su primer bloque. (RN-27.4)
- **CA-27.6 (largo de una celda):** [Deseable] Dada una celda de 301 caracteres, entonces queda en 299 más "…"; dada una de 300, queda igual; dada una de 300 en la columna más angosta, la fila entra en una página. (RN-27.6, RN-27.5)
- **CA-27.7 (vacíos):** [Esencial] Dada una tabla sin filas, entonces se ven sus encabezados y "Sin datos."; dada una sección sin bloques, solo su título; dado un documento sin secciones, una página con título, subtítulo y emisión; dado un subtítulo vacío, no se muestra; dado un bloque sin párrafo ni tabla (párrafo vacío), no es un error, no se muestra nada y el documento se exporta. (RN-27.7, RN-27.2, RN-27.11)
- **CA-27.8 (emojis e ideogramas):** [Esencial] Dado el título "Sprint 😀 final", cuando se genera, entonces no hay error, el título dice "Sprint ? final" y aparece "Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?».". Dado "Equipo 😀😀 y 中文", queda "Equipo ? y ?" y el aviso dice 2 símbolos. (RN-27.8, RN-27.9)
- **CA-27.9 (secuencias de emojis):** [Deseable] Dados `👨‍👩‍👧`, `❤️` y `😀 😀`, cuando se normalizan, entonces quedan `?`, `?` y `? ?`, y cuentan 1, 1 y 2. (RN-27.8, RN-27.9)
- **CA-27.10 (espacios, invisibles y controles):** [Esencial] Dada una celda con una tabulación, un salto de línea y un espacio de ancho cero, entonces los dos primeros pasan a espacio, el tercero se quita y no hay aviso; dada una con `U+0007` o con bytes que no son UTF-8 válido, cada uno pasa a `?` con aviso; dado "¿cumplida?", queda igual y no cuenta. (RN-27.8, RN-27.9)
- **CA-27.11 (documento inválido):** [Esencial] Dado un documento con título vacío o solo con espacios, una sección sin título, un bloque con párrafo y tabla, una tabla sin columnas o una fila con distinta cantidad de celdas que columnas, cuando se valida, entonces se obtiene `ErrDocumentoInvalido` y no se llama al generador. (RN-27.11)
- **CA-27.12 (determinismo):** [Deseable] Dado el mismo documento con la misma emisión, cuando se genera dos veces, entonces los bytes son idénticos. (RN-27.10)
- **CA-27.13 (fecha de emisión):** [Esencial] Dado un reloj en 2026-11-07 17:05 UTC, entonces la emisión dice "07/11/2026 a las 14:05"; dado uno en 2026-11-08 02:30 UTC, dice "07/11/2026 a las 23:30" y el archivo lleva `2026-11-07`. (RN-27.10, RN-27.12)
- **CA-27.14 (descarga del proyecto):** [Esencial] Dado el proyecto 3 y la emisión del 2026-11-07, cuando se pide `GET /web/proyectos/3/reporte/pdf`, entonces la respuesta es `200` con `Content-Type: application/pdf`, `Content-Disposition: attachment; filename="reporte-proyecto-3-2026-11-07.pdf"`, `X-Content-Type-Options: nosniff` y un cuerpo que empieza con `%PDF-`. (Sección 3, RN-27.12)
- **CA-27.15 (descarga del sprint):** [Esencial] Dado el sprint 12 del proyecto 3, cuando se pide `GET /web/proyectos/3/sprints/12/reporte/pdf`, entonces la respuesta es `200` con `filename="reporte-proyecto-3-sprint-12-2026-11-07.pdf"` y las mismas cabeceras. (Sección 3, RN-27.12)
- **CA-27.16 (nombre seguro):** [Esencial] Dado un proyecto llamado "Ñandú / Q4: \"final\"", cuando se descarga su reporte, entonces el nombre del archivo solo tiene minúsculas sin acentos, dígitos, `-` y `.`, y no incluye el nombre del proyecto. (RN-27.12)
- **CA-27.17 (proyecto inexistente):** [Esencial] Dado un ID de proyecto que no existe, cuando se descarga el reporte del proyecto o de un sprint, entonces la respuesta es `404` con la página de error "proyecto no encontrado". (RN-27.13, sección 7)
- **CA-27.18 (sprint inexistente o de otro proyecto):** [Esencial] Dado un ID de sprint que no existe, o un sprint de otro proyecto, cuando se descarga, entonces la respuesta es `404` con la página de error "sprint no encontrado". (RN-27.13, sección 7)
- **CA-27.19 (IDs inválidos):** [Esencial] Dado un ID de proyecto o de sprint `0`, negativo o `abc`, cuando se descarga, entonces la respuesta es `400` con la página de error y el caso de uso no se llama. (Sección 7)
- **CA-27.20 (falla interna):** [Esencial] Dada una falla de la fuente, del generador o un documento inválido, cuando se descarga, entonces la respuesta es `500` con `Content-Type: text/html`, la página "Ocurrió un error interno. Probá de nuevo más tarde." sin el detalle, y el detalle queda en el log. (Sección 7)
- **CA-27.21 (botón y revisión visual, `@manual` en parte):** [Esencial] Dada la pantalla del reporte de un proyecto o de un sprint, entonces tiene un enlace "Descargar PDF" a la ruta de la sección 2, sin formulario ni JavaScript (automatizable). Al abrir el PDF descargado, se ven los acentos, el encabezado y el pie en cada página, las tablas sin filas encimadas ni cortadas y el aviso si hubo reemplazos (`@manual`, porque ningún test mira el dibujo). (Sección 3)
- **CA-27.22 (mismas secciones, en el exportador):** [Esencial] Dado un documento de ejemplo con cinco secciones, "Resumen", "Historias y estimaciones", "Esfuerzo", "Métricas" y "Defectos" (una sin bloques y una con una tabla sin filas), cuando se genera el PDF, entonces el texto extraído tiene los cinco títulos, exactos, una vez cada uno y en ese orden. Los nombres son de ejemplo, no las secciones de HU-26. (RN-27.1, RN-27.2, RN-27.7)
- **CA-27.23 (mismas secciones que la pantalla, integración):** [Esencial] **PENDIENTE DE HU-26.** Dado un proyecto con historias, esfuerzo y defectos, cuando se arma su reporte con HU-26, se convierte a `DocumentoPDF` con la fuente de HU-26 y se genera el PDF, entonces los títulos de sección del PDF son los mismos y en el mismo orden que los títulos de sección de la pantalla del reporte. Lo mismo para el reporte de un sprint. Se automatiza cuando HU-26 tenga pantalla y conversión; HU-27 no se cierra sin este CA en verde. (RN-27.1, RN-27.2; criterio 1 de la issue)
- **CA-27.24 (PDF válido):** [Esencial] Dada una descarga exitosa (CA-27.14 o CA-27.15), entonces el cuerpo empieza con `%PDF-`, termina con `%%EOF` (como mucho seguido de un salto de línea), `ledongthuc/pdf` lo abre, informa al menos una página y extrae el título del documento, y `Content-Length` coincide con la cantidad de bytes del cuerpo. (RN-27.15; criterio 2 de la issue)

Los escenarios mínimos corresponden a: feliz, CA-27.1 y CA-27.14; alternativo, CA-27.15; límite, CA-27.3, CA-27.6 y CA-27.13; caracteres, CA-27.8; errores, CA-27.17 a CA-27.20; criterios de la issue, CA-27.22, CA-27.23 y CA-27.24.

## 9. Trazabilidad con los criterios de la issue

| Criterio de la issue #38 | Cómo se verifica | RN | CA |
|---|---|---|---|
| 1. "El PDF contiene las mismas secciones que el reporte en pantalla." | **(a) En el exportador:** el PDF respeta el orden y los títulos de las secciones del `DocumentoPDF`, y no pierde contenido al paginar ni por los caracteres. Automatizable hoy, con documentos de prueba. **(b) Integración:** el reporte de HU-26, convertido al `DocumentoPDF`, da en el PDF los mismos títulos de sección que su pantalla. **PENDIENTE DE HU-26.** | RN-27.1, RN-27.2, RN-27.4, RN-27.5, RN-27.7, RN-27.8 | (a) CA-27.1, CA-27.2, CA-27.3, CA-27.4, CA-27.7, CA-27.8, CA-27.10, CA-27.22 · (b) CA-27.23 |
| 2. "El archivo se puede descargar/guardar correctamente." | Las descargas HTTP: tipo, nombre seguro y errores 404, 400 y 500. Los bytes son un PDF válido que se puede abrir. El botón de la pantalla y la revisión visual. | RN-27.8, RN-27.10, RN-27.11, RN-27.12, RN-27.13, RN-27.14, RN-27.15 | CA-27.11, CA-27.13, CA-27.14, CA-27.15, CA-27.16, CA-27.17, CA-27.18, CA-27.19, CA-27.20, CA-27.21, CA-27.24 |

RN-27.8 está en los dos criterios: sin esa regla, los acentos no saldrían bien (criterio 1) y un emoji haría fallar la descarga (criterio 2).

**Por qué (b) queda pendiente:** HU-26 no tiene rama, spec ni código; tampoco hay en el repositorio un reporte de ejemplo. Probarlo con otra pantalla existente (por ejemplo, el detalle del proyecto) no verificaría el criterio, porque esa pantalla no es el reporte.

### Prioridad
- **Esenciales:** 12 RN (RN-27.1, 27.2, 27.4, 27.5, 27.7, 27.8, 27.10 a 27.15) y 20 CA (CA-27.1 a 27.4, 27.7, 27.8, 27.10, 27.11, 27.13 a 27.24). RN-27.4 y RN-27.5 tienen una parte deseable, indicada en su marca.
- **Deseables:** 3 RN (RN-27.3, 27.6, 27.9) y 4 CA (CA-27.5, 27.6, 27.9, 27.12).

**Si falta tiempo y se recortan deseables,** estos son los ajustes que arrastra cada recorte:
- **RN-27.3 (encabezado y pie):** CA-27.1, CA-27.3 y CA-27.4 dejan de comprobar el pie; la cantidad de páginas de CA-27.3 se comprueba solo con la que informa el lector.
- **Encabezados repetidos (RN-27.5):** CA-27.3 deja de comprobarlos en cada página.
- **Título de sección al pie (RN-27.4, última oración):** sale CA-27.5.
- **RN-27.6 (tope por celda):** sale CA-27.6. Conviene recortarlo último: sin el tope, una fila más alta que una página no se puede cumplir con RN-27.5 y habría que decidir otra regla para ese caso.
- **RN-27.9 (aviso):** CA-27.8 deja de comprobar el aviso, sale CA-27.9 y RN-27.2 queda sin el paso del aviso. Los reemplazos de RN-27.8 se mantienen.
- **CA-27.12 (determinismo):** sale el test de bytes idénticos. RN-27.10 sigue fijando las fechas.

## Decisiones tomadas por el autor
Tomadas por el autor (Lucas Orellana) el 2026-11-07, abiertas a objeción en el PR:

- **P-1 — Modelo neutral de documento:** el puerto de exportación recibe un `DocumentoPDF` y devuelve los bytes; HU-26 convierte su reporte a ese modelo (RN-27.1, sección 5).

  | Opción | Ventajas | Desventajas |
  |---|---|---|
  | **A. Modelo neutral** (elegida) | HU-27 se programa y se prueba sin esperar a HU-26. El adaptador de PDF no cambia si el reporte suma secciones. | HU-26 tiene que escribir la conversión, y las dos historias tienen que acordar el modelo. |
  | B. El generador recibe el reporte de HU-26 | Sin capa intermedia. | Bloquea HU-27 hasta que HU-26 tenga sus tipos, y cada cambio del reporte toca el adaptador de PDF. |
  | C. Convertir a PDF el HTML de la pantalla | El PDF y la pantalla son idénticos. | Necesita un navegador sin cabeza, que el ADR-0006 descartó. |

- **P-2 — Caracteres no representables:** reemplazarlos por `?`, juntar los seguidos en uno solo y avisar dentro del PDF (RN-27.8, RN-27.9). Motivos:
  - un emoji en un nombre no impide la exportación (riesgo del ADR-0006: `fpdf` devuelve error);
  - el usuario sabe que algo cambió y no se pierde en silencio;
  - la lista de representables es corta y se probó contra la fuente.

  Se descartó quitarlos sin aviso, porque "Sprint 😀 final" pasaría a "Sprint  final" sin que nadie se entere. Se descartó rechazar la exportación, porque el dominio acepta esos caracteres desde HU-01.
- **P-3 — Nombre del archivo:** con IDs y fecha, sin el nombre del proyecto (RN-27.12). Así no hace falta transliterar ni escapar el nombre en la cabecera, y la fecha distingue descargas de distintos días. Se usa el ID del sprint y no su número (RN-08.7) porque es el que está en la ruta y no requiere otra consulta. Objeción posible: el número dice más al usuario ("sprint-2"); cambiarlo no toca otras reglas.
- **P-4 — Orientación:** A4 vertical (RN-27.14). A4 es el tamaño habitual en Argentina para un documento formal. La vertical alcanza para tablas de pocas columnas. Si las tablas de HU-26 tienen muchas columnas, se revisa (pregunta abierta 4).
- **P-5 — Fecha de emisión:** la fija el caso de uso con un reloj inyectado, en UTC-3 fijo y con formato `dd/mm/aaaa HH:MM` (RN-27.10, sección 5). Así el PDF es determinista en los tests, la hora no depende de que Render corra en UTC, y no hace falta la base de husos horarios.
- **P-6 — Tope de 300 caracteres por celda:** garantiza que toda fila entre en una página (RN-27.6). La alternativa, partir una fila entre páginas, complica la función de tabla del adaptador y deja filas difíciles de leer. El texto largo va en párrafos, que no tienen tope.
- **P-7 — Ruta bajo `/web/` y errores como página HTML:** quien pide el PDF es el navegador al seguir el botón, así que un error tiene que mostrarse como página y no como JSON. No se agrega una ruta JSON.
- **P-8 — "?" que se pierden al recortar una celda** (2026-11-07, abierta a objeción en el PR): el aviso de RN-27.9 cuenta los `?` al normalizar, antes de recortar las celdas (RN-27.6). Un `?` que después queda fuera del recorte sigue contando. Motivos:
  - el aviso informa lo que se le hizo al texto;
  - contar solo lo visible obligaría a volver a procesar cada celda recortada;
  - es un caso raro: hace falta una celda de más de 300 caracteres con un emoji después del carácter 299.

  No tiene test propio.

## Preguntas abiertas
**Criterios de la historia:**
1. **Resuelta (2026-11-07, fuente: issue #38). Criterios de aceptación originales de HU-27:** el texto de los criterios no estaba en el repositorio. La issue #38 tiene dos: "El PDF contiene las mismas secciones que el reporte en pantalla." y "El archivo se puede descargar/guardar correctamente.". Se citan al principio de la spec y la sección 9 los liga con los RN y CA.

**Dependen de HU-26** (conviene acordarlas con Ezequiel antes de aprobar la spec):

2. **Secciones del reporte (SUPUESTO):** esta spec no supone ninguna sección concreta. HU-26 define qué secciones lleva y cómo se reparte el contenido entre párrafos y tablas. De esto depende CA-27.23, que verifica el criterio 1 de la issue de punta a punta.
3. **Proyecto y sprint (SUPUESTO):** se supone que los dos reportes se pueden expresar con el mismo modelo, aunque tengan secciones distintas. Si alguno necesita algo que el modelo no tiene, se amplía el modelo.
4. **Columnas de las tablas:** la cantidad de columnas y el largo de los textos confirman o cambian la orientación vertical (P-4) y el tope de RN-27.6.
5. **Gráficos de HU-25 (SUPUESTO: no lleva):** si el reporte incluye gráficos, el modelo necesita un tipo de bloque nuevo y hay que probar cómo dibujarlo con `fpdf`. El ADR-0006 no lo verificó.
6. **Pantalla y rutas (SUPUESTO):** se supone que la pantalla del reporte vive en `/web/proyectos/{id}/reporte` y `/web/proyectos/{id}/sprints/{sprintID}/reporte`.
7. **Dónde vive la conversión:** los formatos de la pantalla (`textoHoras`, `textoPorcentaje`, `conSigno`) están en `internal/adapters/http`. Si la conversión a `DocumentoPDF` los usa, tiene que estar en ese paquete o hay que moverlos. Lo decide HU-26.
8. **Restricciones del reporte (SUPUESTO: ninguna):** si HU-26 no permite el reporte de un sprint abierto, o define otros errores, HU-27 responde igual (RN-27.13, sección 7).
9. **Aclaración del sprint cerrado:** si el reporte incluye la aclaración de RN-23.8, el PDF la lleva como un párrafo más, sin regla propia de HU-27.

**Para el equipo:**

10. **ADR-0006:** confirmar la librería y la fuente para pasarlo de `proposed` a `accepted`. Esta spec depende de esa decisión: si cambia la librería, cambian la lista de RN-27.8 y la sección 5.
11. **Decisiones P-1 a P-8:** sin objeciones hasta ahora; se confirman en el PR.

## Historial
- 2026-11-07: creada (borrador). Decisiones P-1 a P-7 tomadas por el autor, abiertas a objeción en el PR. Preguntas 2 a 9 pendientes de la spec de HU-26.
- 2026-11-07: ligada a los criterios de la issue #38, citados literalmente, y pregunta 1 resuelta. Se agregan RN-27.15 (archivo válido), CA-27.22 (mismas secciones, en el exportador), CA-27.23 (integración con HU-26, pendiente) y CA-27.24 (PDF válido), y la sección 9, con la trazabilidad a los criterios y la prioridad esencial/deseable. Las reglas, criterios y decisiones anteriores no cambian. La spec sigue en borrador.
- 2026-11-07: dos huecos corregidos al escribir los steps de HU-27, con el OK de Lucas Orellana:
  - `DocumentoPDF` suma el campo `Aviso`, que completa el caso de uso y HU-26 nunca completa ni lee (sección 2, RN-27.2, RN-27.9, sección 5).
  - Un bloque sin párrafo ni tabla deja de ser un error, porque contradecía el caso límite "Párrafo vacío": se saca esa condición de RN-27.11 y de CA-27.11, y CA-27.7 suma el caso.

  No cambia ninguna otra regla ni decisión. La spec sigue en borrador.
- 2026-11-07: la sección 5 se ajusta a los puertos que quedaron en el TDD del caso de uso:
  - el caso de uso verifica la existencia del proyecto y del sprint con `LectorProyectos` y `LectorSprintPorID`, en vez de esperar esos errores de la fuente;
  - `FuenteDocumentoReporte` se divide en `FuenteDocumentoProyecto` y `FuenteDocumentoSprint`, de un método cada una;
  - se detalla el orden de los pasos, y las funciones del dominio con sus archivos.

  No cambia ninguna regla, criterio ni decisión. La spec sigue en borrador.
- 2026-11-07: decisión P-8 del autor, abierta a objeción en el PR: los `?` se cuentan al normalizar, antes de recortar las celdas. No cambia ninguna regla ni criterio. La spec sigue en borrador.
- 2026-11-07: la parte HTTP de la sección 5 se ajusta a la convención de rutas de CLAUDE.md: las descargas se registran con `RegistrarExportarReportePDF` en `pantallas_reporte_pdf.go`, y no dentro de `RegistrarPantallas`, y usan su propia tabla de errores 404 en vez de `errorDeConsulta`, que no conoce el sprint. No cambia ninguna regla, criterio ni decisión. La spec sigue en borrador.
