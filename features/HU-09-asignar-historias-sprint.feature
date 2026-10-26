@HU-09
Feature: Asignar historias del backlog a un Sprint
  Como equipo
  quiero asignar historias del Product Backlog a un Sprint
  para definir el Sprint Backlog

  Background:
    Given un proyecto "TP Integrador" con un elemento de backlog "Login de usuarios"

  @CA-09.1 @normal
  Scenario: Asignar un elemento sin sprint a un sprint abierto
    Given un sprint abierto "Sprint 1" del proyecto
    When se asigna "Login de usuarios" al sprint "Sprint 1"
    Then "Login de usuarios" queda asignado al sprint "Sprint 1"

  @CA-09.1 @limite
  Scenario: Reasignar un elemento al mismo sprint que ya tiene no es un error
    Given un sprint abierto "Sprint 1" del proyecto
    And "Login de usuarios" ya está asignado al sprint "Sprint 1"
    When se asigna "Login de usuarios" al sprint "Sprint 1"
    Then "Login de usuarios" queda asignado al sprint "Sprint 1"

  @CA-09.2 @error
  Scenario: No se puede asignar un elemento que ya está en otro sprint abierto
    Given un sprint abierto "Sprint 1" del proyecto
    And un sprint abierto "Sprint 2" del proyecto
    And "Login de usuarios" ya está asignado al sprint "Sprint 1"
    When se asigna "Login de usuarios" al sprint "Sprint 2"
    Then se rechaza la operación con el error "elemento ya asignado"
    And "Login de usuarios" sigue asignado al sprint "Sprint 1"

  @CA-09.3 @alternativo
  Scenario: Se puede reasignar un elemento desde un sprint ya cerrado
    Given un sprint cerrado "Sprint 1" del proyecto
    And un sprint abierto "Sprint 2" del proyecto
    And "Login de usuarios" ya está asignado al sprint "Sprint 1"
    When se asigna "Login de usuarios" al sprint "Sprint 2"
    Then "Login de usuarios" queda asignado al sprint "Sprint 2"

  @CA-09.4 @normal
  Scenario: Desasignar un elemento de un sprint abierto
    Given un sprint abierto "Sprint 1" del proyecto
    And "Login de usuarios" ya está asignado al sprint "Sprint 1"
    When se desasigna "Login de usuarios"
    Then "Login de usuarios" no tiene sprint asignado

  @CA-09.5 @error
  Scenario: No se puede desasignar un elemento de un sprint ya cerrado
    Given un sprint cerrado "Sprint 1" del proyecto
    And "Login de usuarios" ya está asignado al sprint "Sprint 1"
    When se desasigna "Login de usuarios"
    Then se rechaza la operación con el error "sprint cerrado"
    And "Login de usuarios" sigue asignado al sprint "Sprint 1"

  @CA-09.6 @error
  Scenario: No se puede desasignar un elemento que no tiene sprint
    When se desasigna "Login de usuarios"
    Then se rechaza la operación con el error "elemento sin sprint"

  @CA-09.7 @error
  Scenario: No se puede asignar un elemento a un sprint de otro proyecto
    Given un sprint abierto "Sprint de otro proyecto" de otro proyecto
    When se asigna "Login de usuarios" al sprint "Sprint de otro proyecto"
    Then se rechaza la operación con el error "sprint de otro proyecto"
    And "Login de usuarios" no tiene sprint asignado

  @CA-09.8 @error
  Scenario: No se puede asignar un elemento a un sprint ya cerrado
    Given un sprint cerrado "Sprint 1" del proyecto
    When se asigna "Login de usuarios" al sprint "Sprint 1"
    Then se rechaza la operación con el error "sprint cerrado"
    And "Login de usuarios" no tiene sprint asignado
