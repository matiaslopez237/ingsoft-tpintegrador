@HU-18
Feature: Registrar esfuerzo en una historia o tarea
  Como Product Builder
  quiero registrar integrante, fecha, actividad y horas trabajadas en una historia o tarea
  para dejar constancia del esfuerzo real

  # "Hoy" lo fija un reloj falso, para que los escenarios no dependan de la fecha en que corren.
  # Si no se indica otro, el elemento se busca en el proyecto "Alfa".
  # Las horas se escriben con coma decimal, como en la spec (0,5 = media hora).
  # CA-18.15 (respuesta HTTP) no va acá: se cubre con tests del adaptador en la tarea 18.7.

  Background:
    Given hoy es "2026-10-25"
    And el proyecto "Alfa" con inicio "2026-03-01" tiene la integrante "Ana" y el elemento "Login"

  @CA-18.1 @normal
  Scenario: Registrar esfuerzo con todos los datos válidos
    When "Ana" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And el registro tiene un identificador asignado
    And el registro tiene la integrante "Ana", la fecha "2026-10-25", la actividad "Desarrollo" y 2 horas

  @CA-18.1 @alternativo
  Scenario: La actividad se guarda sin los espacios de alrededor
    When "Ana" registra 2 horas de "  Desarrollo  " en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And el registro tiene la integrante "Ana", la fecha "2026-10-25", la actividad "Desarrollo" y 2 horas

  @CA-18.2 @limite
  Scenario: Media hora es válida en un elemento sin esfuerzo registrado
    Given "Login" no tiene esfuerzo registrado
    When "Ana" registra 0,5 horas de "Testing" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And "Login" tiene 1 registros de esfuerzo

  @CA-18.3 @CA-18.9 @CA-18.10 @error
  Scenario Outline: Las horas deben ser mayores a 0, múltiplo de media hora y como máximo 24
    When "Ana" registra <horas> horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "<error>"
    And "Login" tiene 0 registros de esfuerzo

    # -1,25 y 25,3 verifican el orden de validación: primero > 0, después múltiplo de 0,5, al final máximo 24.
    Examples:
      | horas | error                           |
      | 0     | horas inválidas                 |
      | -2    | horas inválidas                 |
      | -1,25 | horas inválidas                 |
      | 1,25  | horas no múltiplo de media hora |
      | 25,3  | horas no múltiplo de media hora |
      | 24,5  | horas exceden el máximo         |

  @CA-18.10 @limite
  Scenario: Un registro de exactamente 24 horas es válido
    When "Ana" registra 24 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente

  @CA-18.4 @error
  Scenario Outline: La actividad es obligatoria
    When "Ana" registra 2 horas de <actividad> en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "actividad vacía"
    And "Login" tiene 0 registros de esfuerzo

    # Las comillas van dentro de la celda porque Gherkin recorta los espacios de las celdas.
    Examples:
      | actividad |
      | ""        |
      | "   "     |

  @CA-18.5 @error
  Scenario Outline: El integrante es obligatorio
    When <integrante> registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "integrante vacío"
    And "Login" tiene 0 registros de esfuerzo

    # Las comillas van dentro de la celda porque Gherkin recorta los espacios de las celdas.
    Examples:
      | integrante |
      | ""         |
      | "   "      |

  @CA-18.6 @error
  Scenario: La fecha es obligatoria
    When "Ana" registra 2 horas de "Desarrollo" en "Login" sin fecha
    Then se rechaza el esfuerzo con el error "fecha obligatoria"
    And "Login" tiene 0 registros de esfuerzo

  @CA-18.6 @error
  Scenario: La fecha no puede ser futura
    When "Ana" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-26"
    Then se rechaza el esfuerzo con el error "fecha futura"
    And "Login" tiene 0 registros de esfuerzo

  @CA-18.6 @limite
  Scenario: La fecha de hoy es válida
    When "Ana" registra 1 horas de "Reunión" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente

  @CA-18.7 @error
  Scenario: No se puede registrar esfuerzo en un proyecto que no existe
    When "Ana" registra 2 horas de "Desarrollo" en "Login" de un proyecto inexistente
    Then se rechaza el esfuerzo con el error "proyecto no encontrado"

  @CA-18.7 @error
  Scenario: No se puede registrar esfuerzo en un elemento que no existe
    When "Ana" registra 2 horas de "Desarrollo" en un elemento inexistente con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "elemento no encontrado"

  @CA-18.7 @error
  Scenario: Un elemento de otro proyecto se trata como inexistente
    Given el proyecto "Beta" con inicio "2026-03-01" tiene la integrante "Beto" y el elemento "Registro"
    When "Ana" registra 2 horas de "Desarrollo" en "Registro" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "elemento no encontrado"
    And "Registro" tiene 0 registros de esfuerzo

  @CA-18.8 @error
  Scenario: Quien registra tiene que ser integrante del proyecto
    When "Caro" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "integrante no encontrado"
    And "Login" tiene 0 registros de esfuerzo

  @CA-18.8 @error
  Scenario: Un integrante de otro proyecto no puede registrar en este
    Given el proyecto "Beta" con inicio "2026-03-01" tiene la integrante "Beto" y el elemento "Registro"
    When "Beto" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "integrante no encontrado"

  @CA-18.8 @alternativo
  Scenario: El nombre del integrante se compara sin espacios ni mayúsculas
    When " ana " registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And el registro tiene la integrante "Ana", la fecha "2026-10-25", la actividad "Desarrollo" y 2 horas

  @CA-18.11 @alternativo
  Scenario: El tope de 24 horas es por registro, no se suman los registros del día
    Given "Ana" ya registró 20 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    When "Ana" registra 8 horas de "Testing" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And "Login" tiene 2 registros de esfuerzo

  @CA-18.12 @normal
  Scenario: Dos registros con los mismos datos son independientes
    Given "Ana" ya registró 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    When "Ana" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then el esfuerzo se registra correctamente
    And los dos registros tienen identificadores distintos
    And cada registro se recupera por su identificador con sus datos

  @CA-18.13 @alternativo
  Scenario: Los registros se conservan al quitar al integrante del proyecto
    Given "Ana" ya registró 3 horas de "Desarrollo" en "Login" con fecha "2026-10-24"
    When se quita a "Ana" del proyecto "Alfa"
    Then el registro de "Ana" se sigue recuperando por su identificador

  @CA-18.13 @CA-18.8 @error
  Scenario: Un integrante quitado del proyecto ya no puede registrar esfuerzo nuevo
    Given "Ana" ya registró 3 horas de "Desarrollo" en "Login" con fecha "2026-10-24"
    And se quita a "Ana" del proyecto "Alfa"
    When "Ana" registra 1 horas de "Desarrollo" en "Login" con fecha "2026-10-25"
    Then se rechaza el esfuerzo con el error "integrante no encontrado"
    And "Login" tiene 1 registros de esfuerzo

  @CA-18.14 @alternativo
  Scenario: Se registra esfuerzo en un proyecto no Activo y con fecha anterior a su inicio
    Given el proyecto "Alfa" pasó al estado "Finalizado"
    When "Ana" registra 2 horas de "Desarrollo" en "Login" con fecha "2026-02-15"
    Then el esfuerzo se registra correctamente
