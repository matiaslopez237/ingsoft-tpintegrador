@HU-07
Feature: Listar y ordenar el backlog por prioridad
  Como Product Owner
  quiero listar y ordenar el backlog por prioridad
  para saber qué se debe entregar primero

  Background:
    Given el proyecto de backlog "Alfa"

  @CA-07.1 @normal
  Scenario: Por defecto el backlog se ordena por prioridad, de Alta a Baja
    Given en el backlog de "Alfa" se cargó "Reportes" con prioridad "Baja"
    And en el backlog de "Alfa" se cargó "Login" con prioridad "Alta"
    And en el backlog de "Alfa" se cargó "Perfil" con prioridad "Media"
    When listo el backlog de "Alfa"
    Then el backlog muestra en este orden "Login, Perfil, Reportes"

  @CA-07.2 @alternativo
  Scenario: En orden ascendente el backlog va de Baja a Alta
    Given en el backlog de "Alfa" se cargó "Reportes" con prioridad "Baja"
    And en el backlog de "Alfa" se cargó "Login" con prioridad "Alta"
    And en el backlog de "Alfa" se cargó "Perfil" con prioridad "Media"
    When listo el backlog de "Alfa" en orden "asc"
    Then el backlog muestra en este orden "Reportes, Perfil, Login"

  @CA-07.3 @limite
  Scenario Outline: A igual prioridad se respeta el orden de creación
    Given en el backlog de "Alfa" se cargó "Login" con prioridad "Alta"
    And en el backlog de "Alfa" se cargó "Reportes" con prioridad "Baja"
    And en el backlog de "Alfa" se cargó "Registro" con prioridad "Alta"
    When listo el backlog de "Alfa" en orden "<orden>"
    Then el backlog muestra en este orden "<titulos>"

    Examples:
      | orden | titulos                   |
      | desc  | Login, Registro, Reportes |
      | asc   | Reportes, Login, Registro |

  @CA-07.4 @normal
  Scenario Outline: Se puede filtrar el backlog por estado
    Given en el backlog de "Alfa" se cargó "Login" con prioridad "Alta" y estado "Pendiente"
    And en el backlog de "Alfa" se cargó "Perfil" con prioridad "Media" y estado "En progreso"
    And en el backlog de "Alfa" se cargó "Reportes" con prioridad "Baja" y estado "Hecho"
    When listo el backlog de "Alfa" con el estado "<estado>"
    Then el backlog muestra en este orden "<titulos>"

    Examples:
      | estado      | titulos  |
      | Pendiente   | Login    |
      | En progreso | Perfil   |
      | Hecho       | Reportes |

  @CA-07.4 @normal
  Scenario: Sin filtro se listan los elementos de todos los estados
    Given en el backlog de "Alfa" se cargó "Login" con prioridad "Alta" y estado "Pendiente"
    And en el backlog de "Alfa" se cargó "Reportes" con prioridad "Baja" y estado "Hecho"
    When listo el backlog de "Alfa"
    Then el backlog muestra en este orden "Login, Reportes"

  @CA-07.5 @limite
  Scenario: Un proyecto sin elementos tiene el backlog vacío
    When listo el backlog de "Alfa"
    Then el backlog está vacío

  @CA-07.5 @limite
  Scenario: Un filtro que no coincide con ningún elemento da un backlog vacío
    Given en el backlog de "Alfa" se cargó "Login" con prioridad "Alta"
    When listo el backlog de "Alfa" con el estado "Hecho"
    Then el backlog está vacío

  @CA-07.6 @alternativo
  Scenario: El backlog solo muestra los elementos de su proyecto
    Given el proyecto de backlog "Beta"
    And en el backlog de "Beta" se cargó "Facturación" con prioridad "Alta"
    And en el backlog de "Alfa" se cargó "Login" con prioridad "Media"
    When listo el backlog de "Alfa"
    Then el backlog muestra en este orden "Login"

  @CA-07.7 @error
  Scenario: No se puede listar el backlog de un proyecto que no existe
    When listo el backlog de un proyecto inexistente
    Then se rechaza el listado con el error "proyecto no encontrado"

  @CA-07.8 @error
  Scenario: El orden tiene que ser asc o desc
    When listo el backlog de "Alfa" en orden "ASC"
    Then se rechaza el listado con el error "orden inválido"

  @CA-07.8 @error
  Scenario: El estado del filtro tiene que ser uno de los tres válidos
    When listo el backlog de "Alfa" con el estado "pendiente"
    Then se rechaza el listado con el error "estado inválido"
