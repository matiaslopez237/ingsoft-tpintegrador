# HU-04: Consultar el estado de un proyecto

**Historia:** Como Product Builder, quiero consultar el estado general de un proyecto, para ver su progreso.
**Issue:** #15  ·  **Sprint:** 1  ·  **Story points:** 3
**Estado de la spec:** Borrador

## 1. Objetivo
Permitir que un Product Builder consulte, a partir del identificador de un proyecto, su información general: nombre, fechas, estado actual y los integrantes registrados con sus roles. Es una operación de solo lectura que reúne lo definido en HU-01 (datos y estado) y HU-03 (integrantes) en una única respuesta.

## 2. Entradas
| Dato | Tipo (Go) | Obligatorio | Validación |
|---|---|---|---|
| id de proyecto | `int64` (`domain.Proyecto.ID`) | Sí | Debe existir un proyecto con ese id. En HTTP llega como `{id}` en la ruta y tiene que ser un entero positivo que entre en `int64`; si no, el adaptador responde 400 (igual que HU-03). |

## 3. Salidas esperadas
Se devuelve el proyecto con:
- `ID`, `Nombre`.
- `FechaInicio` y `FechaFin` (esta última puede no existir).
- `Estado`: el estado guardado del proyecto. Esta historia no lo valida al leerlo; hoy solo HU-01 lo escribe, siempre como `Activo` (RN-01.6, RN-01.7).
- `Integrantes`: lista de `Integrante` (`Nombre` y `Roles []string`). Cada integrante figura una sola vez. Los integrantes van en el orden en que se agregaron al proyecto, y los roles de cada uno, en el orden en que se le asignaron. Lista vacía si no hay integrantes.

No se modifica ningún dato. Si el proyecto no existe, se devuelve un error y no se devuelve ningún dato.

**En HTTP** (`GET /proyectos/{id}`), si la consulta funciona se responde `200` con este JSON:

```json
{
  "id": 1,
  "nombre": "TP Integrador",
  "fecha_inicio": "2026-03-01",
  "fecha_fin": "2026-11-30",
  "estado": "Activo",
  "integrantes": [
    { "nombre": "Ana", "roles": ["Tester", "Desarrolladora"] }
  ]
}
```

- `fecha_fin` es `null` si el proyecto no tiene fecha de fin.
- `integrantes` es **siempre un array**: `[]` si no hay ninguno, nunca `null`.
- `roles` es **siempre un array**, aunque el integrante tenga un solo rol: `["Tester"]`, nunca `"Tester"`.

## 4. Reglas de negocio
Esta historia es de solo lectura y no tiene reglas de cálculo. Los datos que devuelve están en la sección 3 y el caso de proyecto inexistente, en la sección 7. Las reglas propias son:

- **RN-04.1:** La consulta incluye los integrantes registrados en el proyecto (HU-03). Un proyecto sin integrantes devuelve una lista vacía, no un error.
- **RN-04.2:** Los integrantes se devuelven en el orden en que se agregaron al proyecto, y los roles de cada integrante en el orden en que se le asignaron. Quitar un integrante no cambia el orden relativo de los demás. (Decisión del equipo, 2026-10-23.)

## 5. Restricciones
- La consulta no tiene reglas de cálculo: se resuelve en el caso de uso `ConsultarProyecto` (`internal/app`), sin lógica de negocio en el adaptador HTTP.
- `ConsultarProyecto` define su propio puerto, con un solo método: `LectorProyectos { BuscarPorID(ctx context.Context, id int64) (domain.Proyecto, error) }`. No reutiliza `RepositorioIntegrantes` (HU-03), porque ese puerto incluye métodos de escritura (`AgregarIntegrante`, `QuitarIntegrante`, `AgregarRol`) que esta historia no necesita. `persistencia.RepositorioProyectos` ya lo cumple sin cambios.
- El núcleo (`internal/domain`, `internal/app`) no depende de HTTP ni de la base de datos (arquitectura hexagonal).
- **Solo lectura:** consultar un proyecto no modifica ningún dato. Se garantiza por diseño: el puerto que usa `ConsultarProyecto` no tiene métodos de escritura, así que el caso de uso no puede modificar nada y eso lo controla el compilador. En HTTP se expone con `GET`, que es un método *safe* (RFC 9110).
- Persistencia en SQLite (ADR-0002); interfaz web con `net/http` (ADR-0001).
- No hay autenticación ni roles en esta historia: cualquier usuario de la app puede consultar cualquier proyecto (igual que HU-01).
- El id de la ruta se valida con el helper `proyectoIDDesdeRuta` de HU-03 (`internal/adapters/http`), sin duplicar el parseo.
- **"Estado general"** = datos del proyecto + estado del ciclo de vida (RN-01.7) + integrantes con sus roles. El progreso medido (story points, % de historias completadas, defectos) queda para HU-22 a HU-24 (decisión del equipo, Matías, 2026-10-23). Su visualización corresponde al Dashboard (HU-25).
- El orden de RN-04.2 sale de la persistencia: cada integrante y cada rol se insertan como filas nuevas con `id` AUTOINCREMENT, y `buscarIntegrantes` ordena por `i.id, ir.id`. Se verifica con un test de integración contra SQLite, porque un repositorio falso conserva el orden aunque la base no lo haga.

## 6. Casos límite
- Proyecto sin integrantes → lista vacía.
- Proyecto sin fecha de fin → `FechaFin == nil`; en JSON, `"fecha_fin": null`.
- Proyecto con muchos integrantes → se devuelven todos, en el orden en que se agregaron.
- Integrante con varios roles → aparece **una sola vez**, con todos sus roles en orden de alta.
- Integrante con un solo rol → `roles` es un array de un elemento.
- Id positivo que no existe (por ejemplo, 9999) → error de no encontrado (404).
- Id `0` o negativo (`-1`) → 400; no llega al caso de uso.
- Id no numérico (`abc`) → 400.
- Id más grande que `math.MaxInt64` (por ejemplo, `99999999999999999999`) → 400 (`strconv.ParseInt` falla por rango).

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| El id no corresponde a ningún proyecto | `domain.ErrProyectoNoEncontrado` (ya definido en HU-01) | "proyecto no encontrado". En HTTP: 404. |
| Id inválido en la ruta: no numérico, `0`, negativo o fuera del rango de `int64` | No hay error de dominio: lo rechaza `proyectoIDDesdeRuta` en el adaptador HTTP | 400, "el id del proyecto debe ser un número entero positivo" (mismo mensaje que HU-03). No llega al caso de uso. |
| Falla interna (error de base de datos, fecha guardada inválida, etc.) | Cualquier error que no cumpla `errors.Is(err, domain.ErrProyectoNoEncontrado)` | En HTTP: 500, "error interno al consultar el proyecto" (mismo patrón que HU-01). El detalle se registra con `log.Printf` y no se le muestra al cliente. |

Si el caso de uso envuelve errores, usa `fmt.Errorf("...: %w", err)` y no `%v`, para que `errors.Is` siga encontrando `ErrProyectoNoEncontrado` y el adaptador responda 404 en lugar de 500.

## 8. Criterios de aceptación
- **CA-04.1:** Dado un proyecto existente, cuando lo consulto por su id, entonces obtengo su nombre, fechas y estado. (Sección 3)
- **CA-04.2:** Dado un proyecto con integrantes, cuando lo consulto, entonces obtengo cada integrante una sola vez, con su nombre y su lista de roles; los integrantes, y los roles de cada uno, en el orden en que se agregaron. (RN-04.1, RN-04.2)
- **CA-04.3:** Dado un proyecto sin integrantes, cuando lo consulto, entonces obtengo una lista de integrantes vacía y no un error; en la respuesta HTTP, `integrantes` es `[]` y no `null`. (RN-04.1)
- **CA-04.4:** Dado un proyecto sin fecha de fin, cuando lo consulto, entonces el proyecto devuelto tiene `FechaFin == nil` y, en la respuesta HTTP, `"fecha_fin": null`. (Sección 3)
- **CA-04.5:** Dado un id que no corresponde a ningún proyecto, cuando lo consulto, entonces obtengo `ErrProyectoNoEncontrado`. (Sección 7)

## Preguntas abiertas
- Ninguna.

## Limitaciones conocidas
- **Lectura no atómica (*torn read*):** `BuscarPorID` hace dos consultas separadas: primero el proyecto, después sus integrantes con sus roles (estos dos, en una sola consulta con `JOIN`, así que entre sí son consistentes). Si otra operación cambia algo entre esas dos consultas (por ejemplo, HU-02 modifica el proyecto o HU-03 agrega un integrante), la respuesta puede mezclar datos de antes y de después del cambio. No se corrige en este sprint porque el TP no tiene carga concurrente real. Si se vuelve relevante, la solución es hacer las dos lecturas dentro de una misma transacción (`db.BeginTx`), para que ambas vean la misma foto de la base.

## Historial
- 2026-10-23: creada (borrador).
- 2026-10-23: correcciones de la revisión: contrato JSON de salida (`integrantes` siempre array, `fecha_fin` null); "solo lectura" pasa de RN a restricción y se quita el CA-04.6 de idempotencia; RN reducidas a reglas reales (se agrega el orden de integrantes, pendiente de confirmar); validación del id igual a HU-03; puerto propio `LectorProyectos`; error 500 definido con `errors.Is`; el estado no se promete validado; sin autenticación; limitación de lectura no atómica.
- 2026-10-23: decisión del equipo (Matías): "estado general" = datos del proyecto + estado + integrantes; el progreso medido queda para HU-22 a HU-24; integrantes y roles en orden de alta. Quedan resueltas las preguntas abiertas. Se adapta a HU-03 mergeada (PR #87): `Integrante.Roles []string`.
