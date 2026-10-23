@HU-04
Feature: Consultar el estado de un proyecto
  Como Product Builder
  quiero consultar el estado general de un proyecto
  para ver su progreso

  @CA-04.1 @normal
  Scenario: Consultar un proyecto existente
    Given existe un proyecto "TP Integrador" con inicio "2026-03-01" y fin "2026-11-30"
    When consulto el proyecto "TP Integrador"
    Then veo el nombre "TP Integrador"
    And veo el estado "Activo"
    And veo la fecha de inicio "2026-03-01" y la fecha de fin "2026-11-30"

  @CA-04.2 @normal
  Scenario: La consulta incluye a los integrantes con sus roles, en orden de alta
    Given existe un proyecto "TP Integrador" sin integrantes
    And "Beto" se sumó al proyecto con el rol "Tester"
    And "Ana" se sumó al proyecto con el rol "Desarrolladora"
    When consulto el proyecto "TP Integrador"
    Then veo los integrantes "Beto" y "Ana" en ese orden
    And veo que "Beto" tiene solo el rol "Tester"
    And veo que "Ana" tiene solo el rol "Desarrolladora"

  @CA-04.2 @limite
  Scenario: Un integrante con dos roles aparece una sola vez, con los roles en orden de alta
    Given existe un proyecto "TP Integrador" sin integrantes
    And "Ana" se sumó al proyecto con el rol "Tester"
    And "Ana" también sumó el rol "Desarrolladora"
    When consulto el proyecto "TP Integrador"
    Then veo al integrante "Ana" una sola vez
    And veo que "Ana" tiene los roles "Tester" y "Desarrolladora" en ese orden

  @CA-04.3 @limite
  Scenario: Un proyecto sin integrantes devuelve una lista vacía
    Given existe un proyecto "TP Integrador" sin integrantes
    When consulto el proyecto "TP Integrador"
    Then veo que el proyecto no tiene integrantes

  @CA-04.4 @alternativo
  Scenario: Consultar un proyecto sin fecha de fin
    Given existe un proyecto "TP Integrador" con inicio "2026-03-01" y sin fecha de fin
    When consulto el proyecto "TP Integrador"
    Then veo la fecha de inicio "2026-03-01" y ninguna fecha de fin

  @CA-04.5 @error
  Scenario: No se puede consultar un proyecto que no existe
    When consulto el proyecto con identificador 9999
    Then se rechaza la consulta con el error "proyecto no encontrado"
