# Sprint 4 — Planning

**Fechas:** del 5 al 8 de noviembre de 2026 (4 días), confirmadas por Matías el 2026-11-04. La fecha de entrega del TP sigue a confirmar.
**Sprint Goal:** ver el avance del proyecto de un vistazo, compartirlo en un reporte en PDF y dejar el sistema documentado para la entrega.
**Compromiso:** 25 story points (velocidad del Sprint 3: 28 puntos, si se confirma en la review que se terminaron todas las historias).

## Sprint Backlog

El orden es la prioridad: si al final del sprint no se llega con todo, quedan afuera las últimas.

| # | Historia | Pts | Depende de | Responsable (propuesta) |
|---|---|---|---|---|
| 1 | HU-24 Porcentaje de historias completadas y defectos | 3 | HU-11 ✓, HU-20 ✓, HU-21 ✓ | Eze |
| 2 | HU-25 Dashboard con el estado del proyecto y gráficos | 8 | HU-22 ✓, HU-23 ✓, HU-24 | Nico |
| 3 | HU-26 Reporte de un proyecto o sprint | 5 | HU-19 ✓, HU-22 ✓, HU-23 ✓, HU-24 | Eze |
| 4 | HU-27 Exportar el reporte a PDF | 5 | HU-26 | Lucas |
| 5 | HU-29 Documentación técnica | 2 | — | Matías |
| 6 | HU-28 Manual breve de usuario | 2 | HU-25, HU-26, HU-27 (para las capturas) | Matías |
| | **Total** | **25** | | |

## Reparto propuesto

Es una propuesta de Matías, siguiendo lo que hizo cada uno en los sprints anteriores. Cada uno la confirma al asignarse sus historias, y se puede cambiar:

| Responsable | Historias | Pts | Por qué |
|---|---|---|---|
| Eze | HU-24, HU-26 | 8 | HU-24 sigue a sus HU-20, HU-21 y HU-22 (defectos y velocidad), y el reporte reúne esas métricas. |
| Nico | HU-25 | 8 | Es la historia más grande y la que tiene más trabajo de pantalla. |
| Lucas | HU-27 | 5 | Toma el reporte de HU-26 y suma la librería de PDF. Arranca por el ADR, que no depende de nadie. |
| Matías | HU-29, HU-28 | 4 | Tiene menos puntos porque también revisa los PR (Agile Enabler). |

**Cada uno se asigna en GitHub sus historias y sus tareas** (Assignees → *assign yourself*), así queda registrado quién las tomó.

## Tareas

Están cargadas como sub-issues en GitHub, 44 en total.

| Historia | Tareas |
|---|---|
| HU-24 (9) | Spec · Escenarios BDD · Steps de godog · TDD del cálculo · TDD del caso de uso · Persistencia · HTTP · Pantalla · Revisión y trazabilidad |
| HU-25 (10) | ADR de gráficos · Spec · Escenarios BDD · Steps de godog · TDD de los datos del panel · TDD del caso de uso · Persistencia · HTTP · Pantalla con los gráficos · Revisión y trazabilidad |
| HU-26 (9) | Spec · Escenarios BDD · Steps de godog · TDD del armado del reporte · TDD del caso de uso · Persistencia · HTTP · Pantalla · Revisión y trazabilidad |
| HU-27 (8) | ADR de PDF · Spec · Escenarios BDD · Steps de godog · TDD del caso de uso · Adaptador de PDF · HTTP y botón de descarga · Revisión y trazabilidad |
| HU-29 (4) | Arquitectura general · Decisiones técnicas y ADR al día · README y puesta en marcha · Revisión |
| HU-28 (4) | Guion de funcionalidades · Pasos de uso · Capturas · Revisión |

## Cómo repartir

- **HU-24 va primero:** HU-25 y HU-26 usan sus métricas. Mientras tanto, HU-25 y HU-26 pueden avanzar con la spec, los escenarios BDD y las métricas que ya están en `main` (HU-22 y HU-23).
- **HU-26 y HU-27 van en orden:** el PDF contiene las mismas secciones que el reporte. Lucas puede avanzar con el ADR y la spec, y conviene que Eze y Lucas acuerden temprano la estructura de datos del reporte.
- **HU-28 va al final,** porque sus capturas necesitan las pantallas terminadas. HU-29 se puede arrancar ya.

## Preguntas que tiene que resolver la spec antes de programar

- **HU-24, porcentaje (decidido por Matías el 2026-11-04):** se calcula **sobre puntos**, con los planificados y completados que guarda el cierre (RN-11.6). El cierre (HU-11) no guarda cuántas historias tenía el sprint, y este cálculo no lo toca y sirve para los sprints ya cerrados. La spec lo aclara, porque el criterio de la issue habla de historias, y define qué se muestra en un sprint abierto y en uno sin puntos planificados.
- **HU-24, defectos:** los defectos ya guardan el sprint de detección y el de resolución (RN-20.5, RN-21.5). Hay que definir qué se muestra si un defecto no tiene sprint.
- **HU-25, métricas:** el criterio pide al menos 3 gráficas. Hay que elegir cuáles (por ejemplo, velocidad, horas estimadas contra reales y defectos) y desde qué pantalla se entra ("sin pasos adicionales desde la vista principal").
- **HU-25, gráficos:** la decisión va en un ADR. Las pantallas tienen `Content-Security-Policy: default-src 'self'`, que bloquea scripts de CDN, scripts en línea y estilos en línea. Propuesta para evaluar: generar SVG desde las plantillas de Go, sin JavaScript ni dependencias. La alternativa es una librería de JavaScript servida desde `web/static`.
- **HU-26, alcance:** qué contiene el reporte de un proyecto y qué el de un sprint. Las horas de un sprint cerrado siguen la opción A de HU-23: solo cuentan los elementos que siguen en él (RN-23.8), y hay que decidir si el reporte lo aclara igual que la pantalla.
- **HU-27, librería:** la decisión va en un ADR. Tiene que funcionar en Render sin dependencias del sistema, como el resto del proyecto, y mostrar bien los acentos. Hay que verificar que la librería elegida siga mantenida.
- **HU-29 y HU-28, ubicación:** el criterio pide que estén en el repositorio. Hay que definir dónde (por ejemplo, `docs/arquitectura.md` y `docs/manual-usuario.md`).

## Reglas para no bloquearnos

Las mismas de los sprints anteriores:
1. Antes de empezar, `git checkout main` y `git pull`; la rama de la historia sale de ahí.
2. Quien depende de otra historia arranca por lo que no depende de nadie: la spec, los escenarios BDD y el dominio.
3. PR chicos y seguido.
4. Cada historia funcional incluye su pantalla (ADR-0003).

## Ceremonias

- **Planning:** este documento. Cada integrante elige sus historias y completa la columna *Responsable*.
- **Daily:** un mensaje en el grupo al terminar el día: qué hice, qué hago y si estoy trabado.
- **Review y retrospectiva:** al cerrar el sprint, unos 30 minutos. Siguen pendientes las actas de los Sprints 1, 2 y 3.
