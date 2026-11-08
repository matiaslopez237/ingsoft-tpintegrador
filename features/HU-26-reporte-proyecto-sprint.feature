@HU-26
Feature: Generar reporte de proyecto o Sprint
  Como Product Owner
  quiero generar un reporte de un proyecto o Sprint con historias, estimaciones,
  esfuerzo, métricas y defectos
  para compartirlo con los interesados

  # Las historias, los sprints y los defectos se preparan directamente con sus datos, igual que
  # HU-22/HU-24 con los sprints: el escenario prueba que HU-26 junta bien los datos, no el registro
  # ni el ciclo de vida de cada uno (HU-05, HU-06, HU-08, HU-09, HU-10, HU-18, HU-20, HU-21).
  # CA-26.13 (ID inválido) y CA-26.14 (conversión a DocumentoPDF) se cubren con tests de adaptador.

  Background:
    Given el proyecto "Alfa" para el reporte

  @CA-26.1 @normal
  Scenario: El reporte de un proyecto trae todas sus historias
    Given "Alfa" tiene una historia "Login" Pendiente con 5 story points
    And "Alfa" tiene una historia "Logout" Hecha con 3 story points
    When consulto el reporte de "Alfa"
    Then el reporte tiene 2 historias

  @CA-26.2 @normal
  Scenario: El reporte de un Sprint trae solo sus historias asignadas
    Given "Alfa" tiene para el reporte un sprint cerrado con 8 puntos planificados y 8 completados
    And "Alfa" tiene una historia "Login" Hecha con 5 story points asignada al sprint 1
    And "Alfa" tiene una historia "Sin sprint" Pendiente con 3 story points
    When consulto el reporte del sprint 1 de "Alfa"
    Then el reporte tiene 1 historias

  @CA-26.3 @normal
  Scenario: El resumen cuenta las historias y los story points completados
    Given "Alfa" tiene una historia "Login" Hecha con 5 story points
    And "Alfa" tiene una historia "Logout" Hecha con 8 story points
    And "Alfa" tiene una historia "Perfil" Pendiente con 2 story points
    When consulto el reporte de "Alfa"
    Then el reporte tiene 2 de 3 historias completadas y 13 de 15 story points

  @CA-26.4 @normal
  Scenario: El esfuerzo del reporte coincide con las horas de HU-23
    Given "Alfa" tiene una historia "Login" Pendiente con 2 horas estimadas
    And se registraron 5 horas reales en "Login" de "Alfa"
    When consulto el reporte de "Alfa"
    Then el esfuerzo del reporte es 2 horas estimadas y 5 horas reales

  @CA-26.5 @normal
  Scenario: Las métricas de un proyecto traen la velocidad y una fila por sprint cerrado
    Given "Alfa" tiene para el reporte un sprint cerrado con 8 puntos planificados y 6 completados
    And "Alfa" tiene para el reporte otro sprint cerrado con 8 puntos planificados y 8 completados
    When consulto el reporte de "Alfa"
    Then el reporte tiene velocidad 7 sobre 2 sprints cerrados
    And el sprint 1 del reporte tiene 75% completado
    And el sprint 2 del reporte tiene 100% completado

  @CA-26.6 @limite
  Scenario: Sin sprints cerrados, el reporte de un proyecto no tiene métricas
    When consulto el reporte de "Alfa"
    Then el reporte no tiene velocidad
    And el reporte no tiene métricas de sprints

  @CA-26.7 @normal
  Scenario: Las métricas de un Sprint cerrado traen su porcentaje real
    Given "Alfa" tiene para el reporte un sprint cerrado con 4 puntos planificados y 3 completados
    When consulto el reporte del sprint 1 de "Alfa"
    Then el sprint 1 del reporte tiene 75% completado

  @CA-26.8 @alternativo
  Scenario: Un Sprint abierto no tiene datos de porcentaje, pero cuenta sus defectos
    Given "Alfa" tiene para el reporte un sprint todavía abierto
    And se detectó un defecto en el sprint 1 de "Alfa"
    When consulto el reporte del sprint 1 de "Alfa"
    Then el sprint 1 del reporte no tiene datos de porcentaje
    And el sprint 1 del reporte tiene 1 defectos detectados

  @CA-26.9 @normal
  Scenario: El reporte de un proyecto trae todos sus defectos
    Given "Alfa" tiene para el reporte un sprint cerrado con 5 puntos planificados y 5 completados
    And "Alfa" tiene para el reporte otro sprint cerrado con 5 puntos planificados y 5 completados
    And se detectó un defecto en el sprint 1 de "Alfa"
    And se detectó un defecto en el sprint 2 de "Alfa"
    When consulto el reporte de "Alfa"
    Then el reporte tiene 2 defectos

  @CA-26.10 @alternativo
  Scenario: Un defecto detectado en un Sprint y resuelto en otro aparece en los dos reportes
    Given "Alfa" tiene para el reporte un sprint cerrado con 5 puntos planificados y 5 completados
    And "Alfa" tiene para el reporte otro sprint cerrado con 5 puntos planificados y 5 completados
    And se detectó un defecto en el sprint 1 de "Alfa" y se resolvió en el sprint 2
    When consulto el reporte del sprint 1 de "Alfa"
    Then el reporte tiene 1 defectos
    When consulto el reporte del sprint 2 de "Alfa"
    Then el reporte tiene 1 defectos

  @CA-26.11 @error
  Scenario: No se puede consultar el reporte de un proyecto que no existe
    When consulto el reporte de un proyecto inexistente
    Then se rechaza la consulta del reporte con el error "proyecto no encontrado"

  @CA-26.12 @error
  Scenario: No se puede consultar el reporte de un Sprint que no existe
    When consulto el reporte de un sprint inexistente de "Alfa"
    Then se rechaza la consulta del reporte con el error "sprint no encontrado"

  @CA-26.15 @limite
  Scenario: El reporte de un proyecto sin datos no rompe
    When consulto el reporte de "Alfa"
    Then el reporte tiene 0 historias
    And el reporte tiene 0 defectos
