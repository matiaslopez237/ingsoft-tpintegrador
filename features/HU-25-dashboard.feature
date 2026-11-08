@HU-25
Feature: Ver panel con estado del proyecto y gráficos
  Como cualquier rol del equipo
  quiero un Dashboard con el estado del proyecto y gráficos de las métricas principales
  para tener una visión rápida del avance

  # Los elementos, los sprints, los defectos y las horas se preparan directamente con sus datos,
  # igual que HU-22 y HU-24: el escenario prueba lo que junta y calcula el Dashboard, no cómo se
  # cargan esos datos (HU-05, HU-08 a HU-11, HU-18, HU-20 y HU-21 ya tienen sus escenarios).
  # CA-25.8 (respuesta HTTP), CA-25.9 (gráficos en la pantalla) y CA-25.10 (enlace desde el
  # proyecto) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para el dashboard
    And hoy es "2026-03-10" para el dashboard

  @CA-25.1 @normal
  Scenario: El backlog se cuenta por estado
    Given "Alfa" tiene 2 elementos "Pendiente" en el backlog
    And "Alfa" tiene 1 elemento "En progreso" en el backlog
    And "Alfa" tiene 3 elementos "Hecho" en el backlog
    When consulto el dashboard de "Alfa"
    Then el backlog del dashboard tiene 2 pendientes, 1 en progreso y 3 hechos, de 6 en total

  @CA-25.2 @normal
  Scenario: El sprint en curso suma los puntos de sus historias con la regla del cierre
    Given "Alfa" tiene un sprint abierto del "2026-03-01" al "2026-03-14"
    And el sprint en curso de "Alfa" tiene una historia "Pendiente" de 3 puntos
    And el sprint en curso de "Alfa" tiene una historia "Hecho" de 5 puntos
    And el sprint en curso de "Alfa" tiene una historia "En progreso" de 2 puntos
    When consulto el dashboard de "Alfa"
    Then el sprint en curso del dashboard es el 1, con 10 puntos planificados y 5 completados

  @CA-25.2 @limite
  Scenario: Una historia sin estimar del sprint en curso suma 0 puntos
    Given "Alfa" tiene un sprint abierto del "2026-03-01" al "2026-03-14"
    And el sprint en curso de "Alfa" tiene una historia "Hecho" de 5 puntos
    And el sprint en curso de "Alfa" tiene una historia "Hecho" sin estimar
    When consulto el dashboard de "Alfa"
    Then el sprint en curso del dashboard es el 1, con 5 puntos planificados y 5 completados

  @CA-25.3 @alternativo
  Scenario: Un sprint abierto que ya venció no es el sprint en curso
    Given "Alfa" tiene un sprint abierto del "2026-02-01" al "2026-02-14"
    When consulto el dashboard de "Alfa"
    Then el dashboard no tiene sprint en curso

  @CA-25.4 @normal
  Scenario: Los sprints cerrados aparecen del más antiguo al más reciente, con sus métricas
    Given "Alfa" tiene un sprint cerrado del "2026-01-01" al "2026-01-14" con 8 puntos planificados y 6 completados
    And "Alfa" tiene un sprint cerrado del "2026-01-15" al "2026-01-28" con 4 puntos planificados y 4 completados
    And se detectaron 2 defectos en el sprint 1 de "Alfa", y 1 se resolvió en el sprint 2
    When consulto el dashboard de "Alfa"
    Then los sprints cerrados del dashboard son el 1 y el 2, en ese orden
    And el sprint 1 del dashboard tiene 8 puntos planificados, 6 completados y 75% completado
    And el sprint 2 del dashboard tiene 4 puntos planificados, 4 completados y 100% completado
    And el sprint 1 del dashboard tiene 2 defectos detectados y 0 resueltos
    And el sprint 2 del dashboard tiene 0 defectos detectados y 1 resueltos
    And la velocidad del dashboard es 5

  @CA-25.4 @alternativo
  Scenario: Un sprint abierto no aparece entre los sprints cerrados
    Given "Alfa" tiene un sprint cerrado del "2026-01-01" al "2026-01-14" con 8 puntos planificados y 6 completados
    And "Alfa" tiene un sprint abierto del "2026-03-01" al "2026-03-14"
    When consulto el dashboard de "Alfa"
    Then los sprints cerrados del dashboard son el 1, en ese orden

  @CA-25.5 @normal
  Scenario: Las horas del proyecto son las estimadas contra las reales, con su desviación
    Given "Alfa" tiene un elemento estimado en 10 horas con 12 horas reales registradas
    When consulto el dashboard de "Alfa"
    Then el dashboard tiene 10 horas estimadas y 12 reales, con una desviación del 20%

  @CA-25.5 @limite
  Scenario: Sin elementos estimados en horas, la desviación queda sin datos
    Given "Alfa" tiene un elemento sin estimar en horas con 3 horas reales registradas
    When consulto el dashboard de "Alfa"
    Then el dashboard no tiene datos de desviación de horas

  @CA-25.6 @normal
  Scenario: El dashboard muestra los datos vigentes al volver a abrirlo
    Given "Alfa" tiene un sprint abierto del "2026-03-01" al "2026-03-14"
    And el sprint en curso de "Alfa" tiene una historia "En progreso" de 3 puntos
    And consulto el dashboard de "Alfa"
    When se completa la historia de 3 puntos del sprint en curso de "Alfa"
    And consulto el dashboard de "Alfa"
    Then el sprint en curso del dashboard es el 1, con 3 puntos planificados y 3 completados

  @CA-25.7 @error
  Scenario: No se puede consultar el dashboard de un proyecto que no existe
    When consulto el dashboard de un proyecto inexistente
    Then se rechaza la consulta del dashboard con el error "proyecto no encontrado"

  @CA-25.11 @limite
  Scenario: Un proyecto recién creado tiene el dashboard vacío, sin fallar
    When consulto el dashboard de "Alfa"
    Then el backlog del dashboard tiene 0 pendientes, 0 en progreso y 0 hechos, de 0 en total
    And el dashboard no tiene sprint en curso
    And el dashboard no tiene sprints cerrados
    And la velocidad del dashboard no tiene datos
    And el dashboard no tiene datos de desviación de horas
