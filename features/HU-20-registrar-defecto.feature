@HU-20
Feature: Registrar un defecto
  Como Product Builder
  quiero registrar un defecto con descripción, severidad, estado, historia relacionada y Sprint de detección
  para darle seguimiento

  # La fecha de hoy la fija un reloj falso, para que los escenarios no dependan de la fecha en que corren.
  # Si no se indica otro proyecto, el elemento se busca en el proyecto "Alfa".
  # CA-20.6 (respuesta HTTP) no va acá: se cubre con tests del adaptador en la tarea 20.7, como en HU-18.
  # CA-20.8 (pantalla) se agrega en la tarea 20.8, como en HU-06.

  Background:
    Given la fecha de hoy es "2026-10-25"
    And el proyecto "Alfa" con inicio "2026-03-01" tiene el elemento "Login"

  @CA-20.1 @normal
  Scenario: Registrar un defecto con todos los datos válidos y un Sprint activo
    Given el proyecto "Alfa" tiene un sprint con fechas "2026-10-20" a "2026-10-30"
    When registro el defecto "No carga la pantalla de login en Firefox" con severidad "Alta" en "Login"
    Then el defecto se registra correctamente
    And el defecto tiene un identificador asignado
    And el defecto tiene la descripción "No carga la pantalla de login en Firefox", la severidad "Alta", el estado "Abierto" y el Sprint de detección asignado

  @CA-20.1 @alternativo
  Scenario: La descripción se guarda sin los espacios de alrededor
    When registro el defecto "  No carga la pantalla  " con severidad "Media" en "Login"
    Then el defecto se registra correctamente
    And el defecto tiene la descripción "No carga la pantalla", la severidad "Media", el estado "Abierto" y sin Sprint de detección

  @CA-20.2 @alternativo
  Scenario: Sin ningún sprint en el proyecto, el defecto queda sin Sprint de detección
    When registro el defecto "Falla el login" con severidad "Baja" en "Login"
    Then el defecto se registra correctamente
    And el defecto tiene la descripción "Falla el login", la severidad "Baja", el estado "Abierto" y sin Sprint de detección

  @CA-20.2 @limite
  Scenario: Un sprint que ya terminó no se asigna como Sprint de detección
    Given el proyecto "Alfa" tiene un sprint con fechas "2026-09-01" a "2026-09-15"
    When registro el defecto "Falla el login" con severidad "Baja" en "Login"
    Then el defecto se registra correctamente
    And el defecto tiene la descripción "Falla el login", la severidad "Baja", el estado "Abierto" y sin Sprint de detección

  @CA-20.2 @limite
  Scenario: Un sprint cuyo último día es hoy se asigna como Sprint de detección
    Given el proyecto "Alfa" tiene un sprint con fechas "2026-10-10" a "2026-10-25"
    When registro el defecto "Falla el login" con severidad "Baja" en "Login"
    Then el defecto se registra correctamente
    And el defecto tiene la descripción "Falla el login", la severidad "Baja", el estado "Abierto" y el Sprint de detección asignado

  @CA-20.3 @error
  Scenario Outline: La descripción es obligatoria
    When registro el defecto <descripcion> con severidad "Alta" en "Login"
    Then se rechaza el defecto con el error "descripción vacía"

    # Las comillas van dentro de la celda porque Gherkin recorta los espacios de las celdas.
    Examples:
      | descripcion |
      | ""          |
      | "   "       |

  @CA-20.4 @error
  Scenario Outline: La severidad tiene que ser Crítica, Alta, Media o Baja
    When registro el defecto "Falla el login" con severidad <severidad> en "Login"
    Then se rechaza el defecto con el error "severidad inválida"

    Examples:
      | severidad |
      | ""        |
      | "urgente" |
      | "alta"    |

  @CA-20.4 @limite
  Scenario Outline: Las cuatro severidades válidas se aceptan
    When registro el defecto "Falla el login" con severidad <severidad> en "Login"
    Then el defecto se registra correctamente

    Examples:
      | severidad |
      | "Crítica" |
      | "Alta"    |
      | "Media"   |
      | "Baja"    |

  @CA-20.5 @error
  Scenario: No se puede registrar un defecto en un proyecto que no existe
    When registro el defecto "Falla el login" con severidad "Alta" en un proyecto inexistente
    Then se rechaza el defecto con el error "proyecto no encontrado"

  @CA-20.5 @error
  Scenario: No se puede registrar un defecto sobre un elemento que no existe
    When registro el defecto "Falla el login" con severidad "Alta" en un elemento inexistente
    Then se rechaza el defecto con el error "elemento no encontrado"

  @CA-20.5 @error
  Scenario: Un elemento de otro proyecto se trata como inexistente
    Given el proyecto "Beta" con inicio "2026-03-01" tiene el elemento "Registro"
    When registro el defecto "Falla el login" con severidad "Alta" en "Registro"
    Then se rechaza el defecto con el error "elemento no encontrado"

  @CA-20.7 @normal
  Scenario: Dos defectos sobre el mismo elemento son independientes
    When registro el defecto "Falla el login" con severidad "Alta" en "Login"
    And registro el defecto "Falla el login" con severidad "Alta" en "Login"
    Then los dos defectos tienen identificadores distintos
