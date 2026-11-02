@HU-17
Feature: Registrar la estimación acordada
  Como Product Owner
  quiero registrar la estimación final acordada por el equipo
  para dejarla asociada a la historia

  # Las rondas se abren y se votan con los pasos de HU-14 y HU-16.
  # CA-17.9 (respuesta HTTP) se cubre con tests del adaptador.

  Background:
    Given el proyecto de Planning Poker "Alfa" con los integrantes "Ana, Beto"
    And en "Alfa" hay un elemento a estimar "Login" en estado "Pendiente"

  @CA-17.1 @normal
  Scenario: La estimación queda en la ronda y como story points del elemento
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then la estimación se registra sin errores
    And la ronda 1 de "Login" tiene la estimación acordada 8
    And "Login" tiene 8 story points

  @CA-17.2 @alternativo
  Scenario: La estimación reemplaza los story points y no cambia los demás datos
    Given "Login" tiene 5 story points cargados a mano
    And "Login" ya tuvo una ronda cerrada con los votos "Ana: 8, Beto: 8"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then "Login" tiene 8 story points
    And los demás datos de "Login" no cambiaron

  # RN-17.2: después de discutir, el equipo puede acordar un valor que nadie votó.
  @CA-17.3 @alternativo
  Scenario: Se puede acordar un valor que nadie votó
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then la estimación se registra sin errores
    And "Login" tiene 8 story points

  @CA-17.3 @error
  Scenario Outline: El valor tiene que ser de la escala
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    When registro la estimación <valor> en la ronda 1 de "Login"
    Then se rechaza la estimación con el error "estimación inválida"
    And la ronda 1 de "Login" no tiene estimación acordada
    And "Login" sigue sin estimar

    Examples:
      | valor |
      | 4     |
      | 0     |
      | 34    |

  @CA-17.4 @error
  Scenario: No se registra en una ronda abierta
    Given hay una ronda abierta sobre "Login"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then se rechaza la estimación con el error "ronda abierta"
    And "Login" sigue sin estimar

  @CA-17.4 @error
  Scenario: No se registra en una ronda que no es la última
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    And "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then se rechaza la estimación con el error "no es la última"
    And "Login" sigue sin estimar

  # P-5: un elemento Hecho no se estima (RN-14.1).
  @CA-17.4 @error
  Scenario: No se registra si el elemento ya está terminado
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    And "Login" pasó a estado "Hecho"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then se rechaza la estimación con el error "elemento terminado"
    And "Login" sigue sin estimar

  @CA-17.5 @error
  Scenario: Una ronda tiene una sola estimación
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    And registré la estimación 5 en la ronda 1 de "Login"
    When registro la estimación 8 en la ronda 1 de "Login"
    Then se rechaza la estimación con el error "estimación ya registrada"
    And la ronda 1 de "Login" tiene la estimación acordada 5
    And "Login" tiene 5 story points

  # RN-17.4 y RN-17.6: para cambiarla se vota otra ronda; la anterior queda en el historial.
  @CA-17.6 @normal
  Scenario: Una ronda nueva puede cambiar la estimación y la anterior queda registrada
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 5"
    And registré la estimación 5 en la ronda 1 de "Login"
    And "Login" ya tuvo una ronda cerrada con los votos "Ana: 8, Beto: 8"
    When registro la estimación 8 en la ronda 2 de "Login"
    Then "Login" tiene 8 story points
    And la ronda 1 de "Login" tiene la estimación acordada 5
    And la ronda 2 de "Login" tiene la estimación acordada 8

  # RN-17.5: la ronda es un registro de lo que se acordó.
  @CA-17.7 @alternativo
  Scenario: Editar los story points a mano no cambia lo acordado en la ronda
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 8, Beto: 8"
    And registré la estimación 8 en la ronda 1 de "Login"
    When edito a mano los story points de "Login" a 13
    Then "Login" tiene 13 story points
    And la ronda 1 de "Login" tiene la estimación acordada 8

  @CA-17.8 @error
  Scenario: No se registra en una ronda que no existe
    When registro la estimación 8 en una ronda inexistente de "Alfa"
    Then se rechaza la estimación con el error "ronda no encontrada"

  @CA-17.8 @error
  Scenario: No se registra desde un proyecto que no existe
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    When registro la estimación 8 en la ronda 1 de "Login" desde un proyecto inexistente
    Then se rechaza la estimación con el error "proyecto no encontrado"
    And "Login" sigue sin estimar

  @CA-17.8 @error
  Scenario: Una ronda de otro proyecto se trata como inexistente
    Given el proyecto de Planning Poker "Beta" con los integrantes "Caro"
    And "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    When registro la estimación 8 en la ronda 1 de "Login" desde el proyecto "Beta"
    Then se rechaza la estimación con el error "ronda no encontrada"
    And "Login" sigue sin estimar

  @CA-17.10 @normal
  Scenario: Desde la pantalla de la ronda se registra la estimación
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    When desde la pantalla de la ronda 1 de "Login" registro la estimación 8
    Then "Login" tiene 8 story points
    And la pantalla de la ronda 1 de "Login" muestra la estimación acordada 8
    And la pantalla de la ronda 1 de "Login" no ofrece registrar la estimación
    And el historial de la pantalla de "Login" muestra acordada 8

  @CA-17.10 @alternativo
  Scenario: La pantalla de una ronda abierta no ofrece registrar la estimación
    Given hay una ronda abierta sobre "Login"
    Then la pantalla de la ronda 1 de "Login" no ofrece registrar la estimación
