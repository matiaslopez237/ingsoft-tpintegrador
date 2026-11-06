@HU-24
Feature: Calcular el porcentaje de story points completados y los defectos por Sprint
  Como Scrum Master
  quiero ver el porcentaje de story points completados y la cantidad de defectos
  detectados y resueltos de cada Sprint
  para evaluar la calidad del Sprint

  # Los sprints y los defectos se preparan directamente con su estado y sus datos, igual que HU-22
  # con los sprints: el escenario prueba el cálculo de HU-24, no el cierre de un sprint (HU-11) ni
  # el registro o la resolución de un defecto (HU-20, HU-21).
  # CA-24.8 (respuesta HTTP) y CA-24.9 (pantalla) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para la calidad del sprint

  @CA-24.1 @normal
  Scenario: El porcentaje completado son los puntos completados sobre los planificados
    Given "Alfa" tiene un sprint cerrado con 8 puntos planificados y 6 completados
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 75% completado

  @CA-24.1 @limite
  Scenario: Un sprint cerrado con puntos planificados y 0 completados da 0%, no "sin datos"
    Given "Alfa" tiene un sprint cerrado con 5 puntos planificados y 0 completados
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 0% completado

  @CA-24.1 @limite
  Scenario: El porcentaje se redondea a un decimal, con los medios alejándose de cero
    Given "Alfa" tiene un sprint cerrado con 16 puntos planificados y 1 completados
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 6,3% completado

  @CA-24.2 @limite
  Scenario: Un sprint cerrado sin puntos planificados no tiene porcentaje
    Given "Alfa" tiene un sprint cerrado con 0 puntos planificados y 0 completados
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 no tiene datos de porcentaje

  @CA-24.3 @normal
  Scenario: Cuenta los defectos detectados en cada sprint
    Given "Alfa" tiene un sprint cerrado con 5 puntos planificados y 5 completados
    And "Alfa" tiene otro sprint cerrado con 5 puntos planificados y 5 completados
    And se detectaron 2 defectos en el sprint 1 de "Alfa"
    And se detectó 1 defecto en el sprint 2 de "Alfa"
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 2 defectos detectados
    And el sprint 2 tiene 1 defectos detectados

  @CA-24.4 @normal
  Scenario: Cuenta los defectos resueltos en un sprint aunque ya estén cerrados
    Given "Alfa" tiene un sprint cerrado con 5 puntos planificados y 5 completados
    And se resolvió un defecto en el sprint 1 de "Alfa", y después se cerró
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 1 defectos resueltos

  @CA-24.5 @alternativo
  Scenario: Un defecto sin sprint de detección no cuenta en ningún sprint
    Given "Alfa" tiene un sprint cerrado con 5 puntos planificados y 5 completados
    And se detectó un defecto en "Alfa" sin ningún sprint
    When consulto la calidad de los sprints de "Alfa"
    Then el sprint 1 tiene 0 defectos detectados

  @CA-24.6 @alternativo
  Scenario: Un sprint abierto no aparece en la calidad de los sprints
    Given "Alfa" tiene un sprint cerrado con 5 puntos planificados y 5 completados
    And "Alfa" tiene un sprint todavía abierto
    When consulto la calidad de los sprints de "Alfa"
    Then solo aparece 1 sprint en la calidad

  @CA-24.7 @error
  Scenario: No se puede consultar la calidad de los sprints de un proyecto que no existe
    When consulto la calidad de los sprints de un proyecto inexistente
    Then se rechaza la consulta de calidad con el error "proyecto no encontrado"
