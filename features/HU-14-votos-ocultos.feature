@HU-14
Feature: Registrar estimaciones individuales ocultas
  Como Product Builder
  quiero registrar mi estimación de forma oculta hasta que todos voten
  para evitar sesgos por ver el voto de otros

  Background:
    Given el proyecto de Planning Poker "Alfa" con los integrantes "Ana, Beto"
    And en "Alfa" hay un elemento a estimar "Login" en estado "Pendiente"

  @CA-14.1 @normal
  Scenario: Iniciar una ronda sobre un elemento pendiente
    When inicio una ronda de estimación sobre "Login"
    Then la ronda queda "Abierta" con el número 1
    And los votantes de la ronda son "Ana, Beto"
    And en la ronda falta votar "Ana, Beto"

  @CA-14.2 @normal
  Scenario: Un voto no muestra su valor mientras la ronda está abierta
    Given hay una ronda abierta sobre "Login"
    When "Ana" vota 5 en la ronda
    Then la ronda queda "Abierta" con el número 1
    And en la ronda ya votó "Ana"
    And en la ronda falta votar "Beto"
    And el resumen de la ronda no muestra ningún valor

  @CA-14.3 @normal
  Scenario: La ronda se cierra cuando votaron todos
    Given hay una ronda abierta sobre "Login"
    And "Ana" ya votó 5 en la ronda
    When "Beto" vota 8 en la ronda
    Then la ronda queda "Cerrada" con el número 1
    And en la ronda ya votó "Ana, Beto"
    And el resumen de la ronda no muestra ningún valor

  @CA-14.3 @limite
  Scenario: Con un solo integrante, su voto cierra la ronda
    Given el proyecto de Planning Poker "Solo" con los integrantes "Carla"
    And en "Solo" hay un elemento a estimar "Pagos" en estado "En progreso"
    And hay una ronda abierta sobre "Pagos"
    When "Carla" vota 3 en la ronda
    Then la ronda queda "Cerrada" con el número 1

  @CA-14.4 @alternativo
  Scenario: Cambiar el voto mientras la ronda está abierta
    Given hay una ronda abierta sobre "Login"
    And "Ana" ya votó 5 en la ronda
    When "Ana" vota 8 en la ronda
    Then la ronda queda "Abierta" con el número 1
    And en la ronda ya votó "Ana"
    And en la ronda falta votar "Beto"

  @CA-14.4 @alternativo
  Scenario: El nombre del votante no distingue mayúsculas ni espacios
    Given hay una ronda abierta sobre "Login"
    When " ana " vota 5 en la ronda
    Then en la ronda ya votó "Ana"

  @CA-14.5 @error
  Scenario Outline: Un voto fuera de la escala de Fibonacci
    Given hay una ronda abierta sobre "Login"
    When "Ana" vota <valor> en la ronda
    Then se rechaza el voto con el error "voto inválido"
    And en la ronda falta votar "Ana, Beto"

    Examples:
      | valor |
      | 0     |
      | 4     |
      | 34    |
      | -1    |

  @CA-14.5 @error
  Scenario: Vota alguien que se sumó al proyecto después de abrir la ronda
    Given hay una ronda abierta sobre "Login"
    And "Dani" se sumó a "Alfa" después de abrir la ronda
    When "Dani" vota 5 en la ronda
    Then se rechaza el voto con el error "no es votante"
    And los votantes de la ronda son "Ana, Beto"

  @CA-14.5 @error
  Scenario: No se vota en una ronda cerrada
    Given hay una ronda abierta sobre "Login"
    And "Ana" ya votó 5 en la ronda
    And "Beto" ya votó 8 en la ronda
    When "Ana" vota 3 en la ronda
    Then se rechaza el voto con el error "ronda cerrada"

  @CA-14.6 @error
  Scenario: No se estima un elemento terminado
    Given en "Alfa" hay un elemento a estimar "Registro" en estado "Hecho"
    When inicio una ronda de estimación sobre "Registro"
    Then se rechaza la ronda con el error "elemento terminado"

  @CA-14.6 @error
  Scenario: Un proyecto sin integrantes no puede abrir una ronda
    Given el proyecto de Planning Poker "Vacío" sin integrantes
    And en "Vacío" hay un elemento a estimar "Pagos" en estado "Pendiente"
    When inicio una ronda de estimación sobre "Pagos"
    Then se rechaza la ronda con el error "sin votantes"

  @CA-14.6 @error
  Scenario: Un elemento tiene una sola ronda abierta a la vez
    Given hay una ronda abierta sobre "Login"
    When inicio una ronda de estimación sobre "Login"
    Then se rechaza la ronda con el error "ronda ya abierta"

  @CA-14.6 @limite
  Scenario: Con la ronda anterior cerrada se abre la número 2
    Given hay una ronda abierta sobre "Login"
    And "Ana" ya votó 5 en la ronda
    And "Beto" ya votó 8 en la ronda
    When inicio una ronda de estimación sobre "Login"
    Then la ronda queda "Abierta" con el número 2
    And en la ronda falta votar "Ana, Beto"

  @CA-14.7 @error
  Scenario: No se abre una ronda sobre un elemento de otro proyecto
    Given el proyecto de Planning Poker "Beta" con los integrantes "Eva"
    When inicio desde el proyecto "Beta" una ronda de estimación sobre "Login"
    Then se rechaza la ronda con el error "elemento no encontrado"

  @CA-14.7 @error
  Scenario: No se vota en una ronda de otro proyecto
    Given el proyecto de Planning Poker "Beta" con los integrantes "Eva"
    And hay una ronda abierta sobre "Login"
    When "Ana" vota 5 en la ronda desde el proyecto "Beta"
    Then se rechaza el voto con el error "ronda no encontrada"

  @CA-14.9 @normal
  Scenario: Iniciar una ronda desde la pantalla del elemento
    When desde la pantalla del elemento "Login" inicio una ronda de estimación
    Then llego a la pantalla de la ronda número 1 de "Login"
    And la pantalla de la ronda muestra que falta votar "Ana, Beto"

  @CA-14.9 @normal
  Scenario: Votar desde la pantalla sin ver los votos
    Given hay una ronda abierta sobre "Login"
    And "Ana" ya votó 13 en la ronda
    When desde la pantalla de la ronda "Beto" vota 8
    Then la pantalla de la ronda muestra que votaron "Ana, Beto"
    And la pantalla de la ronda dice que está cerrada
    And la pantalla de la ronda no muestra los valores "13, 8"
