@HU-10
Feature: Marcar una historia como completada en el Sprint
  Como Developer
  quiero marcar una historia como completada dentro de un Sprint
  para reflejar el avance real

  # Si no se indica otro, la historia se marca en el proyecto "Alfa".
  # CA-10.6 (respuesta HTTP) no va acá: se cubre con tests del adaptador cuando HU-09 esté en main.

  Background:
    Given el proyecto "Alfa" tiene en su backlog el elemento "Login"
    And el proyecto "Alfa" tiene un sprint abierto "Sprint 1"

  @CA-10.1
  Scenario Outline: Una historia de un sprint abierto pasa a Hecho y sigue en el sprint
    Given "Login" está en el sprint "Sprint 1" con estado "<estado>"
    When se marca "Login" como completada
    Then la historia se marca como completada
    And "Login" queda en estado "Hecho"
    And "Login" sigue en el sprint "Sprint 1"

    @normal
    Examples:
      | estado    |
      | Pendiente |

    @alternativo
    Examples:
      | estado      |
      | En progreso |

  @CA-10.2 @limite
  Scenario: Completar una historia que ya está Hecho no es un error
    Given "Login" está en el sprint "Sprint 1" con estado "Hecho"
    When se marca "Login" como completada
    Then la historia se marca como completada
    And "Login" queda en estado "Hecho"
    And "Login" sigue en el sprint "Sprint 1"

  @CA-10.3 @error
  Scenario: No se puede completar una historia que no está en ningún sprint
    Given "Login" no está en ningún sprint y tiene estado "Pendiente"
    When se marca "Login" como completada
    Then se rechaza completar la historia con el error "elemento sin sprint"
    And "Login" sigue en estado "Pendiente"

  @CA-10.4 @error
  Scenario: No se puede completar una historia de un sprint cerrado
    Given el proyecto "Alfa" tiene un sprint cerrado "Sprint 0"
    And "Login" está en el sprint "Sprint 0" con estado "En progreso"
    When se marca "Login" como completada
    Then se rechaza completar la historia con el error "sprint cerrado"
    And "Login" sigue en estado "En progreso"

  # P-8: el sprint cerrado (RN-10.2) se valida antes que "ya está Hecho" (RN-10.3).
  @CA-10.4 @limite
  Scenario: Una historia ya Hecho de un sprint cerrado también se rechaza
    Given el proyecto "Alfa" tiene un sprint cerrado "Sprint 0"
    And "Login" está en el sprint "Sprint 0" con estado "Hecho"
    When se marca "Login" como completada
    Then se rechaza completar la historia con el error "sprint cerrado"
    And "Login" sigue en estado "Hecho"

  @CA-10.5 @error
  Scenario: No se puede completar una historia de un proyecto que no existe
    Given "Login" está en el sprint "Sprint 1" con estado "Pendiente"
    When se marca "Login" como completada en un proyecto inexistente
    Then se rechaza completar la historia con el error "proyecto no encontrado"
    And "Login" sigue en estado "Pendiente"

  @CA-10.5 @error
  Scenario: No se puede completar una historia que no existe
    When se marca como completada una historia inexistente
    Then se rechaza completar la historia con el error "elemento no encontrado"

  # Se intenta completar desde la ruta de "Alfa" una historia que es de "Beta".
  @CA-10.5 @error
  Scenario: Una historia de otro proyecto se trata como inexistente
    Given el proyecto "Beta" tiene en su backlog el elemento "Registro"
    And el proyecto "Beta" tiene un sprint abierto "Sprint Beta"
    And "Registro" está en el sprint "Sprint Beta" con estado "Pendiente"
    When se marca "Registro" como completada
    Then se rechaza completar la historia con el error "elemento no encontrado"
    And "Registro" sigue en estado "Pendiente"
