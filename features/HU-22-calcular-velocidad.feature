@HU-22
Feature: Calcular la velocidad del equipo
  Como Scrum Master
  quiero ver la velocidad del equipo por proyecto
  para medir el rendimiento

  # Los puntos planificados y completados de un Sprint ya los expone HU-12 (CA-22.1): estos
  # escenarios prueban solo la velocidad, que es lo que agrega HU-22.
  # Los sprints se preparan directamente con su estado y puntos completados, igual que en HU-12:
  # el escenario prueba el cálculo de la velocidad, no el cierre de un sprint.
  # CA-22.7 (respuesta HTTP) y CA-22.8 (pantalla) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para la velocidad

  @CA-22.2 @normal
  Scenario: La velocidad es el promedio de los puntos completados de los sprints cerrados
    Given "Alfa" tiene un sprint cerrado que completó 8 puntos
    And "Alfa" tiene un sprint cerrado que completó 6 puntos
    And "Alfa" tiene un sprint cerrado que completó 4 puntos
    When consulto la velocidad de "Alfa"
    Then la velocidad es 6 sobre 3 sprints cerrados

  @CA-22.2 @alternativo
  Scenario: Un sprint todavía abierto no cuenta para la velocidad
    Given "Alfa" tiene un sprint cerrado que completó 10 puntos
    And "Alfa" tiene un sprint abierto
    When consulto la velocidad de "Alfa"
    Then la velocidad es 10 sobre 1 sprints cerrados

  @CA-22.3 @limite
  Scenario: Un proyecto sin ningún sprint no tiene velocidad
    When consulto la velocidad de "Alfa"
    Then no hay velocidad todavía

  @CA-22.3 @alternativo
  Scenario: Un proyecto con sprints, pero ninguno cerrado, no tiene velocidad
    Given "Alfa" tiene un sprint abierto
    When consulto la velocidad de "Alfa"
    Then no hay velocidad todavía

  @CA-22.4 @limite
  Scenario: La velocidad se informa con decimales cuando no divide exacto
    Given "Alfa" tiene un sprint cerrado que completó 5 puntos
    And "Alfa" tiene un sprint cerrado que completó 4 puntos
    When consulto la velocidad de "Alfa"
    Then la velocidad es 4,5 sobre 2 sprints cerrados

  @CA-22.5 @limite
  Scenario: Un único sprint cerrado con 0 puntos completados da velocidad 0, no "sin datos"
    Given "Alfa" tiene un sprint cerrado que completó 0 puntos
    When consulto la velocidad de "Alfa"
    Then la velocidad es 0 sobre 1 sprints cerrados

  @CA-22.2 @alternativo
  Scenario: Solo cuentan los sprints cerrados del proyecto consultado
    Given el proyecto "Beta" para la velocidad
    And "Beta" tiene un sprint cerrado que completó 20 puntos
    And "Alfa" tiene un sprint cerrado que completó 4 puntos
    When consulto la velocidad de "Alfa"
    Then la velocidad es 4 sobre 1 sprints cerrados

  @CA-22.6 @error
  Scenario: No se puede consultar la velocidad de un proyecto que no existe
    When consulto la velocidad de un proyecto inexistente
    Then se rechaza la consulta de velocidad con el error "proyecto no encontrado"
