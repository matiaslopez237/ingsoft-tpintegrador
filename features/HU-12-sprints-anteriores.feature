@HU-12
Feature: Consultar Sprints anteriores
  Como Scrum Master
  quiero consultar Sprints anteriores de un proyecto
  para revisar el historial

  # Los sprints se preparan directamente con su estado y, si están cerrados, con los puntos que
  # registró el cierre de HU-11 (RN-11.6): el escenario prueba el historial, no el cierre.
  # CA-12.6 (respuesta HTTP) y CA-12.7 (pantalla) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para el historial

  @CA-12.1 @normal
  Scenario: Se listan los sprints cerrados, del más reciente al más antiguo, sin los abiertos
    Given "Alfa" tiene el sprint 1 cerrado del "2026-10-01" al "2026-10-12" con el Sprint Goal "Armar el MVP", 13 puntos planificados y 10 completados
    And "Alfa" tiene el sprint 2 cerrado del "2026-10-13" al "2026-10-25" con el Sprint Goal "Interfaz web", 8 puntos planificados y 5 completados
    And "Alfa" tiene el sprint 3 abierto del "2026-10-26" al "2026-11-08"
    When consulto los sprints anteriores de "Alfa"
    Then los sprints anteriores son los números "2, 1"
    And el sprint anterior 2 va del "2026-10-13" al "2026-10-25" con el Sprint Goal "Interfaz web"
    And el sprint anterior 2 tiene 8 puntos planificados y 5 completados
    And el sprint anterior 1 va del "2026-10-01" al "2026-10-12" con el Sprint Goal "Armar el MVP"
    And el sprint anterior 1 tiene 13 puntos planificados y 10 completados

  # Un sprint cerrado sin historias registró 0 y 0 al cerrarse (RN-11.6).
  @CA-12.1 @limite
  Scenario: Un sprint cerrado sin historias aparece con 0 puntos
    Given "Alfa" tiene el sprint 1 cerrado del "2026-10-01" al "2026-10-12" con el Sprint Goal "Sin historias", 0 puntos planificados y 0 completados
    When consulto los sprints anteriores de "Alfa"
    Then los sprints anteriores son los números "1"
    And el sprint anterior 1 tiene 0 puntos planificados y 0 completados

  @CA-12.2 @limite
  Scenario: Un proyecto sin sprints no tiene sprints anteriores
    When consulto los sprints anteriores de "Alfa"
    Then no hay sprints anteriores

  # RN-12.1: un sprint abierto no aparece aunque su fecha de fin ya haya pasado.
  @CA-12.2 @alternativo
  Scenario: Los sprints abiertos no aparecen, aunque estén vencidos
    Given "Alfa" tiene el sprint 1 abierto del "2020-01-01" al "2020-01-14"
    And "Alfa" tiene el sprint 2 abierto del "2026-10-26" al "2026-11-08"
    When consulto los sprints anteriores de "Alfa"
    Then no hay sprints anteriores

  # RN-12.3: el orden es por fecha de inicio, no por número. Acá el sprint 2 se creó después pero
  # transcurrió antes.
  @CA-12.3 @alternativo
  Scenario: El orden sigue las fechas aunque el número no coincida
    Given "Alfa" tiene el sprint 1 cerrado del "2026-10-13" al "2026-10-25" con el Sprint Goal "Segundo en el tiempo", 5 puntos planificados y 5 completados
    And "Alfa" tiene el sprint 2 cerrado del "2026-10-01" al "2026-10-12" con el Sprint Goal "Primero en el tiempo", 3 puntos planificados y 3 completados
    When consulto los sprints anteriores de "Alfa"
    Then los sprints anteriores son los números "1, 2"

  @CA-12.4 @alternativo
  Scenario: Solo aparecen los sprints del proyecto consultado
    Given el proyecto "Beta" para el historial
    And "Beta" tiene el sprint 1 cerrado del "2026-10-01" al "2026-10-12" con el Sprint Goal "Objetivo de Beta", 8 puntos planificados y 8 completados
    And "Alfa" tiene el sprint 1 cerrado del "2026-10-13" al "2026-10-25" con el Sprint Goal "Objetivo de Alfa", 5 puntos planificados y 3 completados
    When consulto los sprints anteriores de "Alfa"
    Then los sprints anteriores son los números "1"
    And el sprint anterior 1 va del "2026-10-13" al "2026-10-25" con el Sprint Goal "Objetivo de Alfa"

  @CA-12.5 @error
  Scenario: No se pueden consultar los sprints anteriores de un proyecto que no existe
    When consulto los sprints anteriores de un proyecto inexistente
    Then se rechaza la consulta de sprints anteriores con el error "proyecto no encontrado"
