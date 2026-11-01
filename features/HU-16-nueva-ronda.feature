@HU-16
Feature: Repetir una ronda de votación
  Como equipo
  quiero iniciar una nueva ronda de votación sobre la misma historia
  para converger en una estimación

  # Las rondas se abren y se votan con los pasos de HU-14; la revelación es la de HU-15.
  # CA-16.7 (respuesta HTTP) se cubre con tests del adaptador.

  Background:
    Given el proyecto de Planning Poker "Alfa" con los integrantes "Ana, Beto"
    And en "Alfa" hay un elemento a estimar "Login" en estado "Pendiente"

  @CA-16.1 @normal
  Scenario: Después de una ronda cerrada se abre la siguiente y la anterior conserva sus votos
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    When inicio una ronda de estimación sobre "Login"
    Then la ronda queda "Abierta" con el número 2
    And en la ronda falta votar "Ana, Beto"
    And en el historial de "Login", la ronda 1 revela los votos "Ana: 3, Beto: 13" con mínimo 3 y máximo 13

  # RN-16.4: repetir es decisión del equipo, aunque no haya habido diferencias (HU-15).
  @CA-16.2 @alternativo
  Scenario: Se puede repetir la ronda aunque no haya habido diferencias
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 5"
    When inicio una ronda de estimación sobre "Login"
    Then la ronda queda "Abierta" con el número 2

  @CA-16.3 @error
  Scenario: No se puede abrir otra ronda mientras hay una abierta
    Given hay una ronda abierta sobre "Login"
    When inicio una ronda de estimación sobre "Login"
    Then se rechaza la ronda con el error "ronda ya abierta"
    And el historial de rondas de "Login" es "1: Abierta"

  @CA-16.4 @normal
  Scenario: El historial muestra todas las rondas, de la más reciente a la más antigua
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    And "Login" ya tuvo una ronda cerrada con los votos "Ana: 5, Beto: 8"
    And hay una ronda abierta sobre "Login"
    And "Ana" ya votó 5 en la ronda
    When consulto el historial de rondas de "Login"
    Then el historial de rondas es "3: Abierta, 2: Cerrada, 1: Cerrada"
    And en el historial, la ronda 2 revela los votos "Ana: 5, Beto: 8" con mínimo 5 y máximo 8
    And en el historial, la ronda 1 revela los votos "Ana: 3, Beto: 13" con mínimo 3 y máximo 13
    And en el historial, la ronda 3 no revela ningún voto

  @CA-16.5 @limite
  Scenario: Un elemento sin rondas tiene el historial vacío
    When consulto el historial de rondas de "Login"
    Then el historial de rondas está vacío

  @CA-16.6 @error
  Scenario: No se puede consultar el historial desde un proyecto que no existe
    When consulto el historial de rondas de "Login" desde un proyecto inexistente
    Then se rechaza el historial con el error "proyecto no encontrado"

  # RN-16.5: un elemento de otro proyecto se trata como inexistente.
  @CA-16.6 @error
  Scenario: No se puede consultar el historial de un elemento de otro proyecto
    Given el proyecto de Planning Poker "Beta" con los integrantes "Caro"
    When consulto el historial de rondas de "Login" desde el proyecto "Beta"
    Then se rechaza el historial con el error "elemento no encontrado"

  @CA-16.8 @normal
  Scenario: Desde la pantalla de la ronda cerrada se abre la siguiente
    Given "Login" ya tuvo una ronda cerrada con los votos "Ana: 3, Beto: 13"
    When desde la pantalla de la ronda 1 de "Login" uso "Nueva ronda"
    Then llego a la pantalla de la ronda número 2 de "Login"
    And la pantalla de la ronda 1 de "Login" no ofrece "Nueva ronda"
    And la pantalla de la ronda 2 de "Login" no ofrece "Nueva ronda"
    And la pantalla del elemento "Login" lista las rondas "2, 1"
