@HU-01
Feature: Crear un proyecto
  Como Product Builder
  quiero crear un proyecto con nombre, fecha de inicio y fin
  para poder empezar a planificar el trabajo

  @CA-01.1 @normal
  Scenario: Crear un proyecto con datos válidos
    When creo un proyecto "TP Integrador" con inicio "2026-03-01" y fin "2026-11-30"
    Then el proyecto se crea correctamente
    And el proyecto queda en estado "Activo"
    And el proyecto tiene un identificador asignado

  @CA-01.7 @alternativo
  Scenario: Crear un proyecto sin fecha de fin
    When creo un proyecto "TP Integrador" con inicio "2026-03-01" y sin fecha de fin
    Then el proyecto se crea correctamente
    And el proyecto no tiene fecha de fin

  @CA-01.8 @alternativo
  Scenario: El nombre se guarda sin los espacios de alrededor
    When creo un proyecto "   TP Integrador   " con inicio "2026-03-01" y fin "2026-11-30"
    Then el proyecto se crea correctamente
    And el nombre del proyecto es "TP Integrador"

  @CA-01.7 @limite
  Scenario: La fecha de fin puede ser igual a la de inicio
    When creo un proyecto "Proyecto de un día" con inicio "2026-03-01" y fin "2026-03-01"
    Then el proyecto se crea correctamente

  @CA-01.3 @limite
  Scenario: Un nombre de exactamente 100 caracteres es válido
    When creo un proyecto con un nombre de 100 caracteres
    Then el proyecto se crea correctamente

  @CA-01.9 @normal
  Scenario: Cada proyecto recibe un identificador distinto
    Given existe un proyecto "Alfa"
    When creo un proyecto "Beta" con inicio "2026-03-01" y fin "2026-11-30"
    Then los proyectos "Alfa" y "Beta" tienen identificadores distintos
    And cada proyecto puede recuperarse por su identificador

  @CA-01.2 @error
  Scenario: No se puede crear un proyecto sin nombre
    When creo un proyecto "" con inicio "2026-03-01" y fin "2026-11-30"
    Then se rechaza la creación con el error "nombre vacío"
    And no se creó ningún proyecto

  @CA-01.2 @error
  Scenario: Un nombre de solo espacios cuenta como vacío
    When creo un proyecto "     " con inicio "2026-03-01" y fin "2026-11-30"
    Then se rechaza la creación con el error "nombre vacío"
    And no se creó ningún proyecto

  @CA-01.3 @error
  Scenario: Un nombre de más de 100 caracteres es rechazado
    When creo un proyecto con un nombre de 101 caracteres
    Then se rechaza la creación con el error "nombre demasiado largo"
    And no se creó ningún proyecto

  @CA-01.4 @error
  Scenario Outline: No se puede repetir el nombre de un proyecto
    Given existe un proyecto "Alfa"
    When creo un proyecto "<nombre>" con inicio "2026-03-01" y fin "2026-11-30"
    Then se rechaza la creación con el error "nombre duplicado"
    And solo existe 1 proyecto

    Examples:
      | nombre |
      | Alfa   |
      | alfa   |
      | ALFA   |

  @CA-01.5 @error
  Scenario: La fecha de inicio es obligatoria
    When creo un proyecto "TP Integrador" sin fecha de inicio
    Then se rechaza la creación con el error "fecha de inicio obligatoria"
    And no se creó ningún proyecto

  @CA-01.6 @error
  Scenario: La fecha de fin no puede ser anterior a la de inicio
    When creo un proyecto "TP Integrador" con inicio "2026-03-02" y fin "2026-03-01"
    Then se rechaza la creación con el error "fecha de fin anterior"
    And no se creó ningún proyecto
