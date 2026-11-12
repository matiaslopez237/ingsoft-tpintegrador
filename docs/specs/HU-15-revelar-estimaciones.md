# HU-15: Mostrar las estimaciones y detectar diferencias

**Historia:** Como equipo, quiero ver todas las estimaciones reveladas y si hay diferencias significativas, para saber si hace falta discutir.
**Issue:** #26  ·  **Sprint:** 3  ·  **Story points:** 3
**Estado de la spec:** Implementada

## 1. Objetivo
Cerrar el ciclo de una ronda de Planning Poker (HU-14): cuando votaron todos y la ronda se cerró, mostrar el voto de cada integrante, el menor y el mayor, y avisar si la diferencia entre ellos es lo bastante grande como para discutir antes de acordar una estimación (HU-16 repite la ronda; HU-17 registra la estimación acordada).

## 2. Entradas
No hay entradas nuevas: se usa la consulta de la ronda de HU-14 (`GET /proyectos/{id}/rondas/{rondaID}` y la pantalla `/web/proyectos/{id}/rondas/{rondaID}`), con `proyectoID` y `rondaID`.

## 3. Salidas esperadas
- El resumen de la ronda de HU-14 suma la **revelación**, solo si la ronda está `Cerrada`:
  - el voto de cada votante, en el orden de los votantes;
  - el voto **mínimo** y el **máximo**;
  - si hay una **diferencia significativa** (RN-15.2).
- Con la ronda `Abierta`, el resumen es el de HU-14, sin revelación ni valores (RN-14.3).
- **En HTTP:** el JSON de la ronda suma `revelacion` cuando está cerrada; mientras está abierta, no aparece.

```json
"revelacion": {
  "votos": [ { "nombre": "Ana", "valor": 3 }, { "nombre": "Beto", "valor": 8 } ],
  "minimo": 3, "maximo": 8, "diferencia_significativa": true
}
```

- **En la pantalla de la ronda:** cerrada, la tabla de votantes muestra el valor de cada uno, el mínimo y el máximo, y un aviso visible: "Hay diferencias significativas…" con los dos valores, o "No hay diferencias significativas".

## 4. Reglas de negocio
- **RN-15.1:** Los votos se revelan solo cuando la ronda está `Cerrada`, es decir, cuando votaron todos (RN-14.6). Mientras está `Abierta`, siguen ocultos (RN-14.3).
- **RN-15.2:** Hay una diferencia significativa cuando el voto mayor es **más del doble** del menor (`máximo > 2 × mínimo`). Ejemplos: 3 y 8 sí (8 > 6); 3 y 5 no; 1 y 2 no (2 no es más que 2); 1 y 3 sí.
- **RN-15.3:** Si todos votaron lo mismo, o la ronda tiene un solo votante, no hay diferencia significativa (el mínimo y el máximo son iguales).
- **RN-15.4:** La revelación se calcula con los votos guardados al cerrarse la ronda; no cambia después (una ronda cerrada no acepta votos, RN-14.6).

## 5. Restricciones
- RN-15.1 a RN-15.3 viven en el dominio (`internal/domain/planning_poker.go`): `Ronda.Revelar` devuelve la revelación, o `ErrRondaAbierta` si la ronda sigue abierta, y `Ronda.Resumen` la incluye solo si está cerrada. Los adaptadores muestran lo que trae el resumen, así no repiten la regla de cuándo revelar.
- El caso de uso `ConsultarRonda` de HU-14 no cambia: devuelve el resumen, que ya trae la revelación cuando corresponde. No hace falta persistencia nueva: `BuscarRonda` (HU-14) ya lee los votos.
- **Cambio en HU-14:** RN-14.3 y CA-14.8 decían que el resumen "nunca" trae los valores; con HU-15, nunca los trae **mientras la ronda está abierta**. El escenario de pantalla de CA-14.9 pasa a verificar que no se vean los valores con la ronda abierta.
- Fuera de alcance: repetir la ronda (HU-16), acordar la estimación (HU-17), promedios u otras estadísticas.

## 6. Casos límite
- Votos 3 y 8 → mínimo 3, máximo 8, diferencia significativa.
- Votos 5 y 8 → sin diferencia (8 ≤ 10).
- Votos 1 y 2 → sin diferencia (justo el doble no alcanza).
- Votos 1 y 3 → diferencia.
- Votos 5, 5 y 5 → mínimo y máximo 5, sin diferencia.
- Un solo votante → sin diferencia.
- Ronda abierta con algunos votos → sin revelación; pedirla al dominio da `ErrRondaAbierta`.

## 7. Condiciones de error
| Situación | Error (Go) | Mensaje / comportamiento |
|---|---|---|
| Se pide revelar una ronda abierta (en el dominio) | `domain.ErrRondaAbierta` | "los votos se ven cuando votaron todos". La API y la pantalla no lo muestran como error: con la ronda abierta, el resumen simplemente no trae la revelación. |
| Proyecto o ronda inexistentes | los de HU-14 (`ErrProyectoNoEncontrado`, `ErrRondaNoEncontrada`) | HTTP 404 |

## 8. Criterios de aceptación
- **CA-15.1:** Dada una ronda en la que Ana votó 3 y Beto 8, cuando se cierra, entonces el resumen muestra el voto de cada uno, el mínimo 3, el máximo 8 y una diferencia significativa. (RN-15.1, RN-15.2)
- **CA-15.2:** Dada una ronda cerrada con votos 5 y 8, o con todos los votos iguales, entonces no hay diferencia significativa. (RN-15.2, RN-15.3)
- **CA-15.3:** Dados los límites del umbral, 1 y 2 no es diferencia significativa y 1 y 3 sí. (RN-15.2)
- **CA-15.4:** Dada una ronda abierta en la que alguien ya votó, cuando la consulto, entonces no hay revelación ni valores. (RN-15.1)
- **CA-15.5:** El JSON de una ronda cerrada trae `revelacion` con los votos, el mínimo, el máximo y `diferencia_significativa`; el de una ronda abierta no la trae. (Sección 3)
- **CA-15.6:** En la pantalla de una ronda cerrada veo el voto de cada integrante, el mínimo, el máximo y el aviso de si hay diferencias significativas. (Sección 3)

## Decisiones tomadas
Definidas por Nicolás Andino al refinar la historia (2026-11-01):
- Umbral: el voto mayor es más del doble del menor (RN-15.2), como sugiere la tarea 15.4.
- Se muestran el voto de cada integrante, el mínimo y el máximo; sin promedio (en Planning Poker la estimación se acuerda, no se promedia).

## Preguntas abiertas
- Ninguna.

## Historial
- 2026-11-01: creada (borrador).
- 2026-11-11: implementada en el PR #287, ya mergeado y aprobado por otro integrante. Pasa de Borrador a Implementada al ordenar la documentación para la entrega; las reglas y los criterios no cambian.
