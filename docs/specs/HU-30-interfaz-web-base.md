# HU-30: Interfaz web base

**Historia:** Como usuario del sistema, quiero una interfaz web con una base común (diseño, navegación y pantallas de proyectos e integrantes), para operar el sistema desde el navegador.
**Issue:** #103  ·  **Sprint:** 2  ·  **Story points:** 5
**Estado de la spec:** Borrador

## 1. Objetivo
Dar a la aplicación pantallas HTML además de la API JSON (ADR-0003). HU-30 arma la base que reutilizan todas las pantallas (layout, navegación, estilos, forma de mostrar errores) y las pantallas de lo que ya existe en el núcleo: proyectos (HU-01, HU-02, HU-04) e integrantes (HU-03). Suma un caso de uso nuevo, `ListarProyectos`, porque no había forma de listar proyectos.

## 2. Entradas
Las pantallas reciben formularios HTML (`application/x-www-form-urlencoded`) y rutas:

| Pantalla | Ruta | Datos del formulario |
|---|---|---|
| Inicio | `GET /` | — (redirige a la lista de proyectos) |
| Lista de proyectos | `GET /web/proyectos` | — |
| Nuevo proyecto | `GET /web/proyectos/nuevo`, `POST /web/proyectos` | `nombre`, `fecha_inicio`, `fecha_fin` (AAAA-MM-DD; fin opcional) |
| Detalle del proyecto | `GET /web/proyectos/{id}` | — |
| Editar proyecto | `GET /web/proyectos/{id}/editar`, `POST /web/proyectos/{id}/editar` | `nombre`, `fecha_inicio`, `fecha_fin` (vacía = sin fecha de fin), `estado` |
| Agregar integrante | `POST /web/proyectos/{id}/integrantes` | `nombre`, `rol` |
| Quitar integrante | `POST /web/proyectos/{id}/integrantes/{nombre}/quitar` | — |
| Sumar rol | `POST /web/proyectos/{id}/integrantes/{nombre}/roles` | `rol` |
| Estilos | `GET /static/estilos.css` | — |

## 3. Salidas esperadas
- Páginas HTML con el layout común: encabezado con el nombre de la aplicación, navegación y la hoja de estilos `/static/estilos.css`.
- Después de guardar con éxito, redirección `303 See Other` a la pantalla que corresponde (patrón POST-redirección-GET), así recargar la página no repite el envío.
- Ante un error de validación, la misma pantalla con el mensaje del error de dominio y los datos que escribió el usuario, con estado HTTP 400.
- Proyecto inexistente: página de error con estado 404.

## 4. Reglas de negocio
- **RN-30.1:** La lista de proyectos muestra todos los proyectos ordenados por nombre, sin distinguir mayúsculas (`domain.ClaveNombre`); a igual clave, por ID. El orden lo decide el caso de uso `ListarProyectos`, no la plantilla.
- **RN-30.2:** Las pantallas no tienen reglas de negocio propias: cada acción llama al mismo caso de uso que su endpoint JSON (HU-01 a HU-04) y muestra el error que este devuelve.
- **RN-30.3:** El formulario de edición envía siempre todos los campos. Una fecha de fin vacía quita la fecha de fin (RN-02.5); una fecha de inicio vacía da el error de HU-02 (RN-02.4).
- **RN-30.4:** Todo texto que viene del usuario se muestra escapado (sin interpretar HTML).

## 5. Restricciones
- `html/template` de la librería estándar, plantillas y CSS embebidos en el binario con `embed` (ADR-0003). Sin JavaScript ni dependencias nuevas.
- Las pantallas viven en `internal/adapters/http`, bajo el prefijo `/web`, para no chocar con las rutas JSON (`/proyectos/...`). Se registran con `RegistrarPantallas` desde `NuevoRouter`.
- Layout reutilizable: cada historia nueva agrega su plantilla definiendo los bloques `titulo` y `contenido` y la renderiza con el mismo helper.
- `ListarProyectos` está en `internal/app` y usa un puerto de solo lectura; el orden (RN-30.1) se testea con TDD.

## 6. Casos límite
- Sin proyectos: la lista muestra un mensaje que invita a crear el primero.
- Nombres con mayúsculas distintas ("beta", "Alfa", "Gamma") se listan "Alfa", "beta", "Gamma".
- Un nombre con `<script>` se muestra como texto.
- Fecha con formato inválido en el formulario: se muestra el mensaje de formato en la misma pantalla.
- Editar quitando la fecha de fin (campo vacío): el proyecto queda sin fecha de fin.
- Id de proyecto no numérico en la ruta: 400 con página de error.

## 7. Condiciones de error
| Situación | Error (Go) | Comportamiento en pantalla |
|---|---|---|
| Proyecto inexistente | `domain.ErrProyectoNoEncontrado` | Página de error, HTTP 404 |
| Id de proyecto inválido en la ruta | — (adaptador) | Página de error, HTTP 400 |
| Fecha con formato distinto de AAAA-MM-DD | — (adaptador) | Mismo formulario con "fecha_inicio debe tener el formato AAAA-MM-DD" (o fecha_fin), HTTP 400 |
| Cualquier error de validación de HU-01, HU-02 o HU-03 | los de `internal/domain` | Misma pantalla con el mensaje del error, HTTP 400 |
| Falla interna | — | Página de error genérica, HTTP 500, sin mostrar el detalle |

## 8. Criterios de aceptación
- **CA-30.1:** Todas las pantallas comparten el layout: encabezado, navegación con enlace a Proyectos y la hoja de estilos, que se sirve en `/static/estilos.css`. (Restricciones)
- **CA-30.2:** Dados proyectos "beta", "Alfa" y "Gamma", cuando abro la lista, entonces aparecen en el orden "Alfa", "beta", "Gamma"; sin proyectos, aparece el mensaje de lista vacía. (RN-30.1)
- **CA-30.3:** Cuando creo un proyecto válido desde el formulario, entonces se guarda y me lleva a su detalle; con el nombre vacío, veo "el nombre del proyecto es obligatorio" en el mismo formulario y no se crea. (RN-30.2)
- **CA-30.4:** Dado un proyecto con integrantes, cuando abro su detalle, entonces veo sus datos, su estado y sus integrantes con sus roles; si no existe, veo una página de error 404. (RN-30.2)
- **CA-30.5:** Cuando edito un proyecto desde el formulario, entonces se guardan los cambios y vuelvo al detalle; si dejo la fecha de fin vacía, queda sin fecha de fin; si la fecha de fin es anterior a la de inicio, veo el mensaje del error en el formulario. (RN-30.2, RN-30.3)
- **CA-30.6:** Cuando agrego un integrante con rol desde el detalle, entonces aparece en la lista; si ya existe, veo el mensaje del error en el detalle. (RN-30.2)
- **CA-30.7:** Cuando quito un integrante desde el detalle, entonces deja de aparecer. (RN-30.2)
- **CA-30.8:** Cuando le sumo un rol a un integrante, entonces aparece con ese rol; si ya lo tiene, veo el mensaje del error. (RN-30.2)
- **CA-30.9:** Dado un proyecto llamado `<script>alert(1)</script>`, cuando abro la lista, entonces el nombre se ve como texto y no como HTML. (RN-30.4)
- **CA-30.10:** Cuando abro `/`, entonces me redirige a la lista de proyectos.

## Decisiones tomadas
Definidas por Nicolás Andino al implementar la historia (2026-10-26), para validar con el equipo en el PR:
- Prefijo `/web` para las pantallas y `/static` para los archivos estáticos.
- Orden de la lista de proyectos por nombre (sin distinguir mayúsculas).
- Formularios con POST y redirección 303; errores en la misma pantalla con HTTP 400.
- CSS propio mínimo, sin frameworks; sin JavaScript.

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-10-26: creada (borrador).
