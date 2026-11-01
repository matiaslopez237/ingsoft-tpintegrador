@HU-15
Feature: Mostrar las estimaciones y detectar diferencias
  Como equipo
  quiero ver todas las estimaciones reveladas y si hay diferencias significativas
  para saber si hace falta discutir

  # Usa los pasos de HU-14 para armar la ronda y votar.
  Background:
    Given el proyecto de Planning Poker "Alfa" con los integrantes "Ana, Beto"
    And en "Alfa" hay un elemento a estimar "Login" en estado "Pendiente"
    And hay una ronda abierta sobre "Login"

  @CA-15.1 @normal
  Scenario: Al cerrarse la ronda se revelan los votos con una diferencia significativa
    Given "Ana" ya votó 3 en la ronda
    When "Beto" vota 8 en la ronda
    Then la ronda revela los votos "Ana: 3, Beto: 8"
    And el voto mínimo es 3 y el máximo es 8
    And hay una diferencia significativa

  @CA-15.2 @alternativo
  Scenario Outline: Sin diferencia significativa
    Given "Ana" ya votó <ana> en la ronda
    When "Beto" vota <beto> en la ronda
    Then la ronda revela los votos "Ana: <ana>, Beto: <beto>"
    And no hay una diferencia significativa

    Examples:
      | ana | beto |
      | 5   | 8    |
      | 5   | 5    |

  @CA-15.3 @limite
  Scenario Outline: Los límites del umbral (el mayor más del doble del menor)
    Given "Ana" ya votó <ana> en la ronda
    When "Beto" vota <beto> en la ronda
    Then <resultado> una diferencia significativa

    Examples:
      | ana | beto | resultado |
      | 1   | 2    | no hay    |
      | 1   | 3    | hay       |
      | 21  | 8    | hay       |

  @CA-15.3 @limite
  Scenario: Con un solo votante no hay diferencia
    Given el proyecto de Planning Poker "Solo" con los integrantes "Carla"
    And en "Solo" hay un elemento a estimar "Pagos" en estado "Pendiente"
    And hay una ronda abierta sobre "Pagos"
    When "Carla" vota 13 en la ronda
    Then la ronda revela los votos "Carla: 13"
    And el voto mínimo es 13 y el máximo es 13
    And no hay una diferencia significativa

  @CA-15.4 @alternativo
  Scenario: Con la ronda abierta no se revela nada
    When "Ana" vota 5 en la ronda
    Then la ronda no revela ningún voto
    And el resumen de la ronda no muestra ningún valor

  @CA-15.5 @normal
  Scenario: La API muestra la revelación de una ronda cerrada
    Given "Ana" ya votó 3 en la ronda
    When "Beto" vota 8 en la ronda
    Then la API revela los votos "Ana: 3, Beto: 8" con mínimo 3, máximo 8 y diferencia significativa

  @CA-15.6 @normal
  Scenario: La pantalla de una ronda cerrada muestra los votos y el aviso
    Given "Ana" ya votó 3 en la ronda
    When desde la pantalla de la ronda "Beto" vota 8
    Then la pantalla de la ronda dice que está cerrada
    And la pantalla de la ronda muestra los votos "Ana: 3, Beto: 8"
    And la pantalla de la ronda avisa que hay diferencias significativas entre 3 y 8
