@HU-05
Feature: Crear un elemento de backlog
  Como Product Builder
  quiero crear un elemento de backlog con título, descripción, prioridad, story points y criterios de aceptación
  para registrar una necesidad del producto

  Background:
    Given existe el proyecto "Alfa"

  @CA-05.1 @normal
  Scenario: Crear un elemento con todos los datos válidos
    When creo en "Alfa" el elemento "Login" con prioridad "Alta", 5 story points y los criterios:
      | criterio                               |
      | El usuario ingresa con email y clave   |
      | Se muestra un error si la clave es mal |
    Then el elemento se crea correctamente
    And el elemento queda en estado "Pendiente"
    And el elemento tiene un identificador asignado
    And el elemento tiene 5 story points
    And el elemento tiene 2 criterios de aceptación

  @CA-05.2 @alternativo
  Scenario: Crear un elemento sin story points
    When creo en "Alfa" el elemento "Login" con prioridad "Media" y sin story points
    Then el elemento se crea correctamente
    And el elemento queda sin estimar

  @CA-05.3 @error
  Scenario: No se puede crear un elemento en un proyecto que no existe
    When creo en un proyecto inexistente el elemento "Login" con prioridad "Alta" y sin story points
    Then se rechaza el elemento con el error "proyecto no encontrado"

  @CA-05.4 @error
  Scenario Outline: El título es obligatorio y tiene un máximo
    When creo en "Alfa" un elemento con <titulo> y prioridad "Alta"
    Then se rechaza el elemento con el error "<error>"
    And el proyecto "Alfa" tiene 0 elementos en el backlog

    Examples:
      | titulo                  | error                  |
      | un título vacío         | título vacío           |
      | un título de espacios   | título vacío           |
      | un título de 101 letras | título demasiado largo |

  @CA-05.4 @limite
  Scenario: Un título de exactamente 100 caracteres es válido
    When creo en "Alfa" un elemento con un título de 100 letras y prioridad "Alta"
    Then el elemento se crea correctamente

  @CA-05.5 @error
  Scenario: La descripción es obligatoria
    When creo en "Alfa" el elemento "Login" con la descripción "   "
    Then se rechaza el elemento con el error "descripción vacía"

  @CA-05.6 @error
  Scenario Outline: La prioridad debe ser Alta, Media o Baja
    When creo en "Alfa" el elemento "Login" con prioridad "<prioridad>" y sin story points
    Then se rechaza el elemento con el error "prioridad inválida"

    Examples:
      | prioridad |
      | alta      |
      | Urgente   |
      |           |

  @CA-05.7 @limite
  Scenario: Un story point es el mínimo válido
    When creo en "Alfa" el elemento "Login" con prioridad "Baja" y 1 story points
    Then el elemento tiene 1 story points

  @CA-05.7 @error
  Scenario Outline: Los story points deben ser positivos
    When creo en "Alfa" el elemento "Login" con prioridad "Baja" y <puntos> story points
    Then se rechaza el elemento con el error "story points inválidos"

    Examples:
      | puntos |
      | 0      |
      | -3     |

  @CA-05.8 @error
  Scenario: Debe haber al menos un criterio de aceptación
    When creo en "Alfa" el elemento "Login" sin criterios de aceptación
    Then se rechaza el elemento con el error "sin criterios"

  @CA-05.8 @error
  Scenario: Ningún criterio de aceptación puede estar vacío
    When creo en "Alfa" el elemento "Login" con prioridad "Alta", 3 story points y los criterios:
      | criterio                 |
      | El usuario puede entrar  |
      | <vacío>                  |
    Then se rechaza el elemento con el error "criterio vacío"

  @CA-05.9 @error
  Scenario: No se repite un título dentro del mismo proyecto
    Given en "Alfa" existe el elemento "Login"
    When creo en "Alfa" el elemento "  login " con prioridad "Alta" y sin story points
    Then se rechaza el elemento con el error "título duplicado"
    And el proyecto "Alfa" tiene 1 elementos en el backlog

  @CA-05.9 @alternativo
  Scenario: El mismo título se puede usar en otro proyecto
    Given existe el proyecto "Beta"
    And en "Alfa" existe el elemento "Login"
    When creo en "Beta" el elemento "Login" con prioridad "Alta" y sin story points
    Then el elemento se crea correctamente

  @CA-05.10 @normal
  Scenario: Cada elemento recibe un identificador distinto y se recupera con sus criterios
    Given en "Alfa" existe el elemento "Login"
    When creo en "Alfa" el elemento "Registro" con prioridad "Alta", 8 story points y los criterios:
      | criterio        |
      | Primer criterio |
      | Segundo         |
      | Tercero         |
    Then los elementos "Login" y "Registro" tienen identificadores distintos
    And cada elemento se recupera por su identificador con sus criterios en orden
