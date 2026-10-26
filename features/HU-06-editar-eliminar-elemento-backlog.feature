@HU-06
Feature: Editar o eliminar un elemento del backlog
  Como Product Builder
  quiero editar o eliminar un elemento del backlog
  para mantenerlo actualizado

  Background:
    Given el proyecto "Alfa" tiene en su backlog el elemento "Login" con prioridad "Alta" y 5 story points

  @CA-06.1 @normal
  Scenario: Cambiar solo la prioridad
    When edito el elemento "Login" cambiando "prioridad" a "Baja"
    Then el elemento "Login" queda con:
      | campo        | valor               |
      | prioridad    | Baja                |
      | story_points | 5                   |
      | estado       | Pendiente           |
      | descripcion  | Necesidad de prueba |

  @CA-06.2 @normal
  Scenario: Editar varios campos juntos
    When edito el elemento "Login" con:
      | campo        | valor                          |
      | titulo       | Inicio de sesión               |
      | descripcion  | Entrar con email y clave       |
      | story_points | 8                              |
      | criterios    | Con email; Con clave; Con 2FA  |
      | estado       | En progreso                    |
    Then el elemento "Inicio de sesión" queda con:
      | campo        | valor                         |
      | descripcion  | Entrar con email y clave      |
      | story_points | 8                             |
      | criterios    | Con email; Con clave; Con 2FA |
      | estado       | En progreso                   |
      | prioridad    | Alta                          |
    And el elemento "Inicio de sesión" conserva el identificador de "Login"

  @CA-06.3 @alternativo
  Scenario: Quitar la estimación
    When edito el elemento "Login" cambiando "story_points" a "null"
    Then el elemento "Login" queda con:
      | campo        | valor      |
      | story_points | sin estimar |

  @CA-06.2 @limite
  Scenario: Un elemento Hecho puede volver a Pendiente
    Given el elemento "Login" está en estado "Hecho"
    When edito el elemento "Login" cambiando "estado" a "Pendiente"
    Then el elemento "Login" queda con:
      | campo  | valor     |
      | estado | Pendiente |

  @CA-06.4 @error
  Scenario Outline: Los datos editados cumplen las reglas de HU-05
    When edito el elemento "Login" cambiando "<campo>" a "<valor>"
    Then se rechaza la edición del elemento con el error "<error>"
    And el elemento "Login" no cambió

    Examples:
      | campo        | valor   | error                  |
      | titulo       |         | título vacío           |
      | descripcion  |         | descripción vacía      |
      | prioridad    | Urgente | prioridad inválida     |
      | story_points | 0       | story points inválidos |
      | criterios    |         | sin criterios          |
      | estado       | Cerrado | estado inválido        |
      | estado       | hecho   | estado inválido        |

  @CA-06.5 @error
  Scenario: No se puede renombrar a un título que ya existe en el proyecto
    Given el proyecto "Alfa" tiene en su backlog el elemento "Registro" con prioridad "Media" y 3 story points
    When edito el elemento "Login" cambiando "titulo" a "registro"
    Then se rechaza la edición del elemento con el error "título duplicado"
    And el elemento "Login" no cambió

  @CA-06.5 @alternativo
  Scenario: Renombrar al mismo título con otras mayúsculas
    When edito el elemento "Login" cambiando "titulo" a "LOGIN"
    Then el elemento "LOGIN" queda con:
      | campo     | valor |
      | prioridad | Alta  |

  @CA-06.6 @error
  Scenario: Editar sin enviar ningún campo
    When edito el elemento "Login" sin cambiar nada
    Then se rechaza la edición del elemento con el error "sin cambios"

  @CA-06.7 @error
  Scenario: No se edita ni se elimina un elemento de otro proyecto
    Given el proyecto "Beta" tiene en su backlog el elemento "Pagos" con prioridad "Media" y 3 story points
    When edito desde el proyecto "Alfa" el elemento "Pagos" cambiando "prioridad" a "Baja"
    Then se rechaza la edición del elemento con el error "elemento no encontrado"
    When elimino desde el proyecto "Alfa" el elemento "Pagos"
    Then se rechaza la eliminación del elemento con el error "elemento no encontrado"
    And el elemento "Pagos" sigue existiendo

  @CA-06.7 @error
  Scenario: No se elimina un elemento inexistente
    When elimino el elemento "Login"
    And elimino el elemento "Login"
    Then se rechaza la eliminación del elemento con el error "elemento no encontrado"

  @CA-06.8 @normal
  Scenario: Eliminar un elemento sin esfuerzo
    When elimino el elemento "Login"
    Then el elemento "Login" ya no existe
    And el backlog de "Alfa" tiene 0 elementos

  @CA-06.9 @error
  Scenario: No se elimina un elemento con esfuerzo registrado
    Given "Ana" registró 2 horas de esfuerzo en el elemento "Login"
    When elimino el elemento "Login"
    Then se rechaza la eliminación del elemento con el error "elemento con esfuerzo"
    And el elemento "Login" sigue existiendo

  @CA-06.10 @normal
  Scenario: Editar desde la pantalla guarda solo lo que cambió
    When desde la pantalla del elemento "Login" cambio "prioridad" a "Media"
    Then el elemento "Login" queda con:
      | campo        | valor |
      | prioridad    | Media |
      | story_points | 5     |

  @CA-06.10 @error
  Scenario: La pantalla muestra el error de dominio
    When desde la pantalla del elemento "Login" cambio "titulo" a ""
    Then la pantalla del elemento responde 400 y muestra "el título es obligatorio"
    And el elemento "Login" no cambió

  @CA-06.10 @normal
  Scenario: Eliminar desde la pantalla
    When desde la pantalla elimino el elemento "Login"
    Then el elemento "Login" ya no existe
