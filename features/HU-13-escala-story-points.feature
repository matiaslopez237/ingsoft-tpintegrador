@HU-13
Feature: Estimar una historia con story points de la escala Fibonacci
  Como Product Builder
  quiero asignar story points a una historia
  para reflejar su tamaño relativo

  Background:
    Given el proyecto "Alfa" para estimar con story points

  @CA-13.1 @normal
  Scenario Outline: Los valores de la escala se aceptan
    When creo en "Alfa" el elemento "Login" estimado en <puntos> story points
    Then el elemento estimado "Login" tiene <puntos> story points

    Examples:
      | puntos |
      | 1      |
      | 2      |
      | 3      |
      | 5      |
      | 8      |
      | 13     |
      | 21     |

  @CA-13.2 @error
  Scenario Outline: Un valor fuera de la escala se rechaza
    When creo en "Alfa" el elemento "Login" estimado en <puntos> story points
    Then se rechaza la estimación con el error de la escala
    And el backlog de "Alfa" no tiene el elemento "Login"

    Examples:
      | puntos |
      | 4      |
      | 10     |
      | 34     |
      | 0      |

  @CA-13.3 @normal
  Scenario: Editar los story points a otro valor de la escala
    Given en "Alfa" existe el elemento "Login" estimado en 5 story points
    When cambio los story points de "Login" a 13
    Then el elemento estimado "Login" tiene 13 story points

  @CA-13.3 @error
  Scenario: Editar los story points a un valor fuera de la escala
    Given en "Alfa" existe el elemento "Login" estimado en 5 story points
    When cambio los story points de "Login" a 10
    Then se rechaza la estimación con el error de la escala
    And el elemento estimado "Login" tiene 5 story points

  @CA-13.4 @alternativo
  Scenario: Un elemento puede quedar sin estimar
    When creo en "Alfa" el elemento "Login" sin estimar
    Then el elemento estimado "Login" queda sin estimar

  @CA-13.4 @alternativo
  Scenario: Quitarle la estimación a un elemento
    Given en "Alfa" existe el elemento "Login" estimado en 8 story points
    When le quito la estimación a "Login"
    Then el elemento estimado "Login" queda sin estimar

  @CA-13.5 @limite
  Scenario: Un valor guardado antes de la escala se respeta al editar otro campo
    Given en "Alfa" existe el elemento "Login" guardado antes de la escala con 4 story points
    When cambio el título de "Login" a "Ingreso"
    Then el elemento estimado "Ingreso" tiene 4 story points

  @CA-13.5 @error
  Scenario: Un valor guardado antes de la escala no se puede volver a enviar
    Given en "Alfa" existe el elemento "Login" guardado antes de la escala con 4 story points
    When cambio los story points de "Login" a 4
    Then se rechaza la estimación con el error de la escala

  @CA-13.5 @limite
  Scenario: Un valor guardado antes de la escala se puede corregir
    Given en "Alfa" existe el elemento "Login" guardado antes de la escala con 4 story points
    When cambio los story points de "Login" a 5
    Then el elemento estimado "Login" tiene 5 story points

  @CA-13.6 @error
  Scenario: La API rechaza un valor fuera de la escala
    When por la API creo en "Alfa" el elemento "Login" con 4 story points
    Then la API responde 400 con el mensaje de la escala
