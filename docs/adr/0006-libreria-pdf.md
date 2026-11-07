# ADR-0006: Librería para exportar el reporte a PDF

**Date**: 2026-11-05
**Status**: proposed
**Deciders**: Propuesto por Lucas Orellana (HU-27); falta la confirmación del equipo (Matías López, Ezequiel Assenza, Nicolás Andino)

## Context

HU-27 ("Exportar el reporte a PDF", 5 pts) exporta el reporte que arma HU-26. La planificación del Sprint 4 pide decidir la librería en un ADR: tiene que funcionar en Render sin dependencias del sistema, mostrar bien los acentos y seguir mantenida.

Cómo se construye y despliega hoy el proyecto, según el repositorio:
- `go.mod` declara `go 1.27.0`. El equipo y el CI usan esa línea de Go (`go-version-file: go.mod`).
- No hay `Dockerfile`. `render.yaml` usa el runtime nativo de Go, con `go build -o bin/server ./cmd/server` y el plan `free`. ADR-0004 descartó Docker porque el servicio es un solo binario sin dependencias de sistema.
- ADR-0002 eligió un driver de SQLite sin CGO para que `go build` funcione en cualquier máquina.
- Las pantallas usan `html/template` (ADR-0003) y envían `Content-Security-Policy: default-src 'self'`. Esa cabecera se aplica en `internal/adapters/http/pantallas.go`, solo a las respuestas HTML; ADR-0003 no la menciona. Un PDF descargado no depende de ella.

### Restricciones

1. **Sin dependencias del sistema:** nada de navegador sin cabeza, wkhtmltopdf, LaTeX ni binarios externos. Compila con `CGO_ENABLED=0`. La documentación de Render indica que una herramienta que no viene en el runtime nativo requiere desplegar con Docker, y eso contradice ADR-0004.
2. **Español correcto:** UTF-8 con una fuente TrueType embebida.
3. **Documentos largos:** tablas y texto con saltos de página, y comprobables en un test sin abrir el PDF a mano.
4. **Licencia compatible con un TP académico:** se descartan las que exigen licencia comercial o clave, y la AGPL. La AGPL obligaría a ofrecer bajo esa licencia el código de toda la aplicación a quienes la usen por la red, una condición que el equipo no eligió para su código.
5. **Mantenida:** sin archivar y con versiones recientes.

## Decision

Usar **`codeberg.org/go-pdf/fpdf` v0.12.0** (MIT) con la fuente **DejaVu Sans Condensed 2.37**, regular y negrita, embebida con `go:embed`.

Los motivos, en orden de peso:
- **No suma módulos al binario:** el paquete `fpdf` solo usa la biblioteca estándar de Go. `maroto`, la otra finalista, suma 13 módulos externos.
- **Falla de forma controlable:** ante un emoji, `fpdf` devuelve un `error`; `maroto` entra en pánico (`panic`). También entra en pánico con UTF-8 inválido, y `fpdf` no.
- **Cumple todas las restricciones en el spike**, incluidos los saltos de página con el encabezado de la tabla repetido.
- **Es un adaptador fácil de testear:** con las fechas fijas, generó un PDF idéntico byte a byte en Windows y en Linux.

**La ruta del módulo importa.** `github.com/go-pdf/fpdf` está archivado y su README remite a Codeberg. `go get` de la ruta de GitHub instala v0.9.0, de 2023.

El costo de esta elección es que `fpdf` no trae un componente de tabla: el adaptador necesita una función propia que calcule la altura de cada fila y decida el salto de página (unas 30 líneas en el spike).

Reglas para el adaptador:
- Solo `internal/adapters/pdf` importa `fpdf`. El dominio y los casos de uso no lo conocen (arquitectura hexagonal).
- Se fijan la fecha con `SetCreationDate` y `SetModificationDate` del documento, y `SetCatalogSort(true)`. No se usan las funciones `SetDefault*`, que cambian estado global del paquete y afectarían a pedidos simultáneos.
- Se devuelve siempre el error de `Output`: `fpdf` acumula el primer error interno y lo informa ahí.

### Fuente

DejaVu Sans Condensed 2.37 viene en el módulo de `fpdf` (`font/`); la versión se leyó de la propia fuente. Se copian al repo `DejaVuSansCondensed.ttf` y `DejaVuSansCondensed-Bold.ttf`.

- **Cobertura verificada:** tiene todos los glifos de `áéíóúüñÁÉÍÓÚÜÑ¿¡€°ªº–—“”…→≤≥%`. La variante condensada aprovecha mejor el ancho en tablas.
- **Licencia:** la de Bitstream Vera y Arev, publicada en dejavu-fonts.github.io (consultada el 2026-11-05); los cambios de DejaVu son de dominio público. Permite usar y redistribuir la fuente si se conservan los avisos. Si se la modifica, hay que renombrarla, y no se puede vender sola. No es OFL.

## Alternatives Considered

Evidencia obtenida el **2026-11-05**. Versiones y fechas: `go list -m -json <módulo>@latest`. Estado del repositorio: API de GitHub y de Codeberg. Licencia: archivo `LICENSE` del módulo descargado. Módulos externos: `go list -deps` del paquete principal, sin contar el propio. Compilación: `CGO_ENABLED=0 go build` en un módulo temporal.

| Candidata | Última versión (fecha) | Repositorio | Licencia | Módulos externos | Compila sin CGO | Sin dep. del sistema | Tablas y saltos de página | Resultado |
|---|---|---|---|---:|---|---|---|---|
| **`codeberg.org/go-pdf/fpdf`** | v0.12.0 (2026-05-18) | Activo, sin archivar | MIT | 0 | Sí | Sí | Salto automático; tabla con función propia. Probado | **Elegida** |
| `johnfercher/maroto/v2` | v2.4.4 (2026-10-22) | Activo, sin archivar | MIT | 13 | Sí | Sí | Filas y paginado propios. Probado | Finalista: más dependencias y pánico con emoji |
| `signintech/gopdf` | v0.38.4 (2026-10-31) | Activo, sin archivar | MIT | 2 | Sí | Sí | **No salta de página.** Probado: 200 filas en 1 hoja | Descartada |
| `github.com/go-pdf/fpdf` | v0.9.0 (2023-09-13) | **Archivado** (último push: 2025-03-04) | MIT | 0 | Sí | Sí | No probado | Descartada: la reemplaza la de Codeberg |
| `jung-kurt/gofpdf` | v2.17.3 (2021-11-13) | **Archivado** | MIT | 0 (v1.16.2) | Sí (v1.16.2) | Sí | No probado | Descartada: archivada |
| `phpdave11/gofpdf` (agregada) | v1.4.3 (2025-05-01) | Sin archivar (último push: 2026-04-16) | MIT | 0 | Sí | Sí | No probado; `maroto` la usa por dentro | Descartada: sin ventaja sobre la de Codeberg y con la última versión más antigua |
| `unidoc/unipdf/v5` | v5.2.0 (2026-10-20) | Activo | **Comercial** (EULA) | 17 | Sí | Sí | No probado | Descartada: "requires a license code to operate" |
| `chromedp/chromedp` (HTML→PDF) | v0.22.1 (2026-11-04) | Activo | MIT | 1 | Sí | **No: necesita Chrome** | No probado | Descartada: exige Docker en Render y el plan free tiene 512 MB de RAM |
| `wkhtmltopdf` (HTML→PDF) | No es un módulo de Go | **Archivado** (último push: 2022-11-22) | LGPL-3.0 | No aplica | No aplica | **No: binario externo** | No probado | Descartada |
| `pdfcpu/pdfcpu` (agregada) | v0.16.4 (2026-10-27) | Activo | Apache-2.0 | 8 | Sí | Sí | No probado | Descartada: se presenta como herramienta para PDF ("PDF tooling"), no como generador de reportes, y suma 8 módulos |
| `carlos7ags/folio` (agregada) | v0.11.0 (2026-10-23) | Activo; creado el 2026-03-14 | Apache-2.0 | 3 | Sí | Sí | Según su README, sí. No probado | No finalista: el proyecto tiene 8 meses y no se probó. Es la primera alternativa si `fpdf` no alcanza |

Ninguna candidata usa hoy AGPL. En `unipdf` v3 el último tag es de 2025-04-27, pero la línea vigente es v5. Con la misma búsqueda aparecieron `tdewolff/canvas`, `VantageDataChat/GoPDF2` y `benoitkugler/pdf`; no se evaluaron a fondo porque las dos primeras no publican versiones con tag y la tercera no se revisó más allá de su versión.

### Spike de verificación

Se hizo en un directorio temporal fuera del repo, ya borrado. Generó el texto `Acentuación: áéíóúñ ¿?¡! Proyecto Ñandú` y una tabla de 3 columnas con Go 1.27.3, en versiones de 3 y de 200 filas.

| Comprobación | `fpdf` | `maroto` | `gopdf` |
|---|---|---|---|
| `go test` con `CGO_ENABLED=0` compila | Sí | Sí | Sí |
| Texto extraído con `ledongthuc/pdf` (3 filas) | Correcto | Correcto | Correcto |
| Texto extraído con `pdftotext` (3 filas) | Correcto | Correcto | No se comparó |
| Páginas con 200 filas | 6 | 4 | **1** |
| Filas leídas por `pdftotext` (200 filas) | 200 | 200 | **38** |
| Revisión visual de una imagen de la página | Correcta | Correcta | **Filas encimadas y tabla cortada** |
| Tamaño con 200 filas | 41 KB | 252 KB | 35 KB |
| Tiempo y memoria por documento con 200 filas | 21,5 ms; 20,3 MB | 21,7 ms; 18,7 MB | No se midió |
| Emoji (`😀`) | **`error`** | **`panic`** | No se probó |
| UTF-8 inválido | Muestra `�` | **`panic`** | No se probó |
| Texto vacío, salto de línea, tabulación, carácter de control y palabra de 400 letras | Sin errores | Sin errores | No se probó |
| Ideogramas chinos o japoneses | Se extraen, pero se dibujan como cuadrados | Se extraen | No se probó |

Además, con `fpdf`:
- 32 generaciones simultáneas dieron el mismo resultado.
- Un binario Linux compilado con `CGO_ENABLED=0` (estático según `file` y `ldd`) corrió en WSL2. Generó un PDF idéntico byte a byte al de Windows.

**Lo que el spike no cubrió:** no se ejecutó en Render, y `-race` no corrió (el gcc de MinGW instalado es de 32 bits y WSL no tiene Go). Lo cubrirá el CI cuando exista el adaptador.

## Consequences

### Positive
- **`go.mod`:** dos líneas en el bloque de dependencias directas, comprobadas sobre una copia del repo: `codeberg.org/go-pdf/fpdf v0.12.0` y `github.com/ledongthuc/pdf` (BSD-3-Clause, solo para tests). `go.sum` suma 4 líneas. La versión de Go no cambia: `fpdf` pide 1.25.
- **Sin cambios en el despliegue:** `render.yaml` y el CI siguen iguales.
- **Binario:** para Linux amd64 y con el comando de Render, pasa de 21,1 MB a 23,6 MB (+2,5 MB, +12 %), de los que 1,35 MB son las dos fuentes. `go version -m` confirma que `ledongthuc/pdf` no entra al binario.
- **Tests sin intervención manual:** el adaptador genera el PDF, lo lee con `ledongthuc/pdf` y comprueba el texto y la cantidad de páginas. También sirve para los steps de godog.

### Negative
- Hay que escribir y testear la función de tabla del adaptador.
- `maroto` tiene versiones más recientes y `fpdf` una comunidad más chica en Codeberg (31 estrellas; el repo archivado de GitHub tiene 642).
- Se agregan al repo dos archivos binarios de fuente, de 1,35 MB en total.
- Hay que sumar la licencia de DejaVu a `docs/third-party-licenses/`, junto con el adaptador. La licencia exige conservar los avisos.
- `ledongthuc/pdf` no publica versiones con tag: se fija una pseudo-versión.

### Cómo se testea el PDF
- **El texto extraído no prueba lo que se ve.** Sale del mapa de caracteres del PDF, no del dibujo: `ledongthuc/pdf` leyó las 200 filas de `gopdf`, aunque `pdftotext` solo encontró 38. Por eso se comprueba también la cantidad de páginas, y la revisión visual queda como escenario `@manual`.
- **El extractor corta el texto** donde el PDF cambia de línea: hay que normalizar los espacios o buscar fragmentos cortos.
- **Para repetir bytes exactos**, el adaptador recibe la fecha por parámetro; así es posible una comparación con un archivo de referencia si hiciera falta.

### Qué queda fuera
- Convertir HTML a PDF: el PDF se arma desde los datos, no desde las plantillas.
- PDF/A, firma digital, formularios, marcadores e imágenes.
- Soporte de escrituras no latinas: los dos ideogramas probados no están en DejaVu y las demás escrituras no se probaron.
- Gráficos, salvo que HU-26 los incluya en el reporte (ver preguntas abiertas).

### Risks
- **Un emoji puede impedir la exportación.** El dominio solo limita el largo de nombres y títulos, así que acepta emojis, y `fpdf` rechaza los caracteres fuera del plano básico de Unicode aunque la fuente los tenga. Mitigación: tratarlos antes de generar y testearlo; la regla la define la spec (ver preguntas abiertas).
- **Mantenimiento más lento:** la última versión de `fpdf` es de mayo de 2026. Si se abandonara, el código ya descargado sigue compilando, y el acceso queda aislado en el adaptador, así que el cambio de librería no toca el núcleo.
- **Dependencia de Codeberg:** el spike descargó el módulo por `proxy.golang.org`, el proxy por defecto de Go, y el CI no lo cambia. No se probó qué pasa si Codeberg no responde, ni qué proxy usa Render.
- **Memoria en el plan free:** cada documento reserva unos 20 MB, en buena parte por procesar las fuentes, y Render ofrece 512 MB. No se ve un problema con el uso previsto (ADR-0002), pero no se midió en Render.
- **Numeración de este ADR:** se escribió como 0005, pero la rama de HU-25 ya usa ese número para el ADR de gráficos (`docs/adr/0005-graficos-dashboard.md`). Por eso pasó a ser el 0006, el número que ya le asignaban `docs/arquitectura.md` y `docs/trazabilidad.md`.

## Preguntas abiertas

Dependen de HU-26 (estructura de datos del reporte):
1. **Estructura y puerto:** qué secciones y tipos tiene el reporte de proyecto y el de sprint, y por lo tanto qué firma tiene el puerto de exportación de HU-27. Conviene acordarlo temprano con Ezequiel, como pide la planificación.
2. **Gráficos:** si el reporte incluye gráficos (HU-25), hay que probar cómo llevarlos al PDF. El spike no lo verificó.
3. **Formato de las tablas:** la cantidad de columnas y el largo de los textos definen la orientación de la hoja y los anchos.
4. **Fecha de emisión:** si el reporte la muestra, hay que definir el huso horario y pasarla por parámetro.

Para la spec de HU-27, porque son reglas que este ADR no debe inventar:
5. **Caracteres no representables:** qué se hace con emojis, ideogramas y caracteres de control (reemplazarlos, quitarlos u otra opción) y si se le avisa al usuario.
6. **Aclaración del sprint cerrado:** si HU-26 la incluye (RN-23.8), el PDF lleva el mismo texto.

Para el equipo:
7. **Aceptación:** confirmar este ADR para pasarlo de `proposed` a `accepted`.

## Notes

Primera tarea de HU-27 en la planificación del Sprint 4. Esta etapa solo agrega el ADR: `go.mod` y el código de producción no cambian hasta el adaptador.
