@HU-21
Feature: Actualizar el estado de un defecto
  Como Product Builder
  quiero actualizar el estado de un defecto y registrar su Sprint de resolución
  para reflejar su progreso

  # La fecha de referencia la fija un reloj falso, para que los escenarios no dependan de la fecha
  # en que corren. CA-21.9 (respuesta HTTP) y CA-21.10 (listado JSON) no van acá: se cubren con
  # tests del adaptador, como en HU-18/HU-20. CA-21.11 (pantalla) se agrega en la tarea 21.8, como
  # en HU-20. La secuencia es siempre Abierto → En progreso → Resuelto → Cerrado (RN-21.5).

  Background:
    Given la fecha de referencia es "2026-10-25"
    And el proyecto "Alfa" iniciado el "2026-03-01" ya tiene el elemento "Login"
    And hay un defecto "No carga la pantalla de login" con severidad "Alta" en "Login"

  @CA-21.1 @normal
  Scenario: Iniciar un defecto Abierto
    When inicio el defecto
    Then el defecto se actualiza correctamente
    And el defecto queda en estado "En progreso"

  @CA-21.2 @normal
  Scenario: Resolver un defecto En progreso con un Sprint activo
    Given el proyecto "Alfa" cuenta con un sprint vigente entre "2026-10-20" y "2026-10-30"
    And inicié el defecto
    When resuelvo el defecto
    Then el defecto se actualiza correctamente
    And el defecto queda en estado "Resuelto" con el Sprint de resolución asignado

  @CA-21.3 @alternativo
  Scenario: Resolver un defecto En progreso sin ningún Sprint vigente
    Given inicié el defecto
    When resuelvo el defecto
    Then el defecto se actualiza correctamente
    And el defecto queda en estado "Resuelto" sin Sprint de resolución

  @CA-21.3 @limite
  Scenario: Un sprint que ya terminó no se asigna como Sprint de resolución
    Given el proyecto "Alfa" cuenta con un sprint vigente entre "2026-09-01" y "2026-09-15"
    And inicié el defecto
    When resuelvo el defecto
    Then el defecto se actualiza correctamente
    And el defecto queda en estado "Resuelto" sin Sprint de resolución

  @CA-21.4 @normal
  Scenario: Cerrar un defecto Resuelto
    Given el proyecto "Alfa" cuenta con un sprint vigente entre "2026-10-20" y "2026-10-30"
    And inicié el defecto
    And resolví el defecto
    When cierro el defecto
    Then el defecto se actualiza correctamente
    And el defecto queda en estado "Cerrado" con el Sprint de resolución asignado

  @CA-21.5 @error
  Scenario: No se puede iniciar un defecto que ya está En progreso
    Given inicié el defecto
    When inicio el defecto
    Then se rechaza la actualización con el error "el defecto no está Abierto"

  @CA-21.5 @error
  Scenario: No se puede iniciar un defecto que ya está Cerrado
    Given inicié el defecto
    And resolví el defecto
    And cerré el defecto
    When inicio el defecto
    Then se rechaza la actualización con el error "el defecto no está Abierto"

  @CA-21.6 @error
  Scenario: No se puede resolver un defecto Abierto (saltear En progreso)
    When resuelvo el defecto
    Then se rechaza la actualización con el error "el defecto no está En progreso"

  @CA-21.6 @error
  Scenario: No se puede resolver un defecto que ya está Cerrado
    Given inicié el defecto
    And resolví el defecto
    And cerré el defecto
    When resuelvo el defecto
    Then se rechaza la actualización con el error "el defecto no está En progreso"

  @CA-21.7 @error
  Scenario: No se puede cerrar un defecto Abierto (saltear En progreso y Resuelto)
    When cierro el defecto
    Then se rechaza la actualización con el error "el defecto no está Resuelto"

  @CA-21.7 @error
  Scenario: No se puede cerrar un defecto En progreso (saltear Resuelto)
    Given inicié el defecto
    When cierro el defecto
    Then se rechaza la actualización con el error "el defecto no está Resuelto"

  @CA-21.8 @error
  Scenario: No se puede actualizar un defecto de un proyecto inexistente
    When resuelvo el defecto en un proyecto inexistente
    Then se rechaza la actualización con el error "proyecto no encontrado"

  @CA-21.8 @error
  Scenario: No se puede actualizar un defecto inexistente
    When resuelvo un defecto inexistente
    Then se rechaza la actualización con el error "defecto no encontrado"

  @CA-21.8 @error
  Scenario: Un defecto de otro proyecto se trata como inexistente
    Given el proyecto "Beta" iniciado el "2026-03-01" ya tiene el elemento "Registro"
    When resuelvo el defecto de "Alfa" indicando el proyecto "Beta"
    Then se rechaza la actualización con el error "defecto no encontrado"
