@HU-02
Feature: Modificar un proyecto
  Como Product Builder
  quiero modificar los datos de un proyecto existente (incluido su estado)
  para poder corregir su información o reflejar su situación actual

  Background:
    Given hay un proyecto "Alfa" con inicio "2026-03-01" y fin "2026-11-30"

  @CA-02.1 @normal
  Scenario: Cambiar solo el nombre no toca el resto de los datos
    When modifico el proyecto "Alfa" con el nombre "Alfa Renovado"
    Then la modificación se realiza correctamente
    And el proyecto modificado se llama "Alfa Renovado"
    And el proyecto modificado tiene inicio "2026-03-01" y fin "2026-11-30"
    And el proyecto modificado está en estado "Activo"

  @CA-02.1 @alternativo
  Scenario: El nuevo nombre se guarda sin los espacios de alrededor
    When modifico el proyecto "Alfa" con el nombre "   Alfa Renovado   "
    Then la modificación se realiza correctamente
    And el proyecto modificado se llama "Alfa Renovado"

  @CA-02.2 @error
  Scenario: Un nombre vacío es rechazado
    When modifico el proyecto "Alfa" con el nombre ""
    Then se rechaza la modificación con el error "nombre vacío"
    And el proyecto "Alfa" no cambió

  @CA-02.2 @error
  Scenario: Un nombre de solo espacios cuenta como vacío
    When modifico el proyecto "Alfa" con el nombre "     "
    Then se rechaza la modificación con el error "nombre vacío"
    And el proyecto "Alfa" no cambió

  @CA-02.2 @limite
  Scenario: Un nombre de exactamente 100 caracteres es válido
    When modifico el proyecto "Alfa" con un nombre de 100 caracteres
    Then la modificación se realiza correctamente

  @CA-02.2 @error
  Scenario: Un nombre de más de 100 caracteres es rechazado
    When modifico el proyecto "Alfa" con un nombre de 101 caracteres
    Then se rechaza la modificación con el error "nombre demasiado largo"
    And el proyecto "Alfa" no cambió

  @CA-02.3 @error
  Scenario Outline: No se puede usar el nombre de otro proyecto
    Given hay un proyecto "Beta" con inicio "2026-04-01" y fin "2026-10-31"
    When modifico el proyecto "Beta" con el nombre "<nombre>"
    Then se rechaza la modificación con el error "nombre duplicado"
    And el proyecto "Beta" no cambió

    Examples:
      | nombre |
      | Alfa   |
      | alfa   |
      | ALFA   |

  @CA-02.3 @alternativo
  Scenario: Un proyecto puede cambiar solo las mayúsculas de su propio nombre
    When modifico el proyecto "Alfa" con el nombre "ALFA"
    Then la modificación se realiza correctamente
    And el proyecto modificado se llama "ALFA"

  @CA-02.4 @error
  Scenario: La nueva fecha de fin no puede ser anterior a la fecha de inicio guardada
    When modifico el proyecto "Alfa" con la fecha de fin "2026-02-28"
    Then se rechaza la modificación con el error "fecha de fin anterior"
    And el proyecto "Alfa" no cambió

  @CA-02.4 @error
  Scenario: La nueva fecha de inicio no puede ser posterior a la fecha de fin guardada
    When modifico el proyecto "Alfa" con la fecha de inicio "2026-12-01"
    Then se rechaza la modificación con el error "fecha de fin anterior"
    And el proyecto "Alfa" no cambió

  @CA-02.4 @limite
  Scenario: La nueva fecha de fin puede ser igual a la fecha de inicio
    When modifico el proyecto "Alfa" con la fecha de fin "2026-03-01"
    Then la modificación se realiza correctamente
    And el proyecto modificado tiene inicio "2026-03-01" y fin "2026-03-01"

  @CA-02.5 @alternativo
  Scenario: Se puede quitar la fecha de fin
    When modifico el proyecto "Alfa" quitando la fecha de fin
    Then la modificación se realiza correctamente
    And el proyecto modificado no tiene fecha de fin

  @CA-02.6 @normal
  Scenario: Finalizar un proyecto activo
    When modifico el proyecto "Alfa" con el estado "Finalizado"
    Then la modificación se realiza correctamente
    And el proyecto modificado está en estado "Finalizado"

  @CA-02.6 @alternativo
  Scenario: Reactivar un proyecto cancelado
    Given el proyecto "Alfa" está en estado "Cancelado"
    When modifico el proyecto "Alfa" con el estado "Activo"
    Then la modificación se realiza correctamente
    And el proyecto modificado está en estado "Activo"

  @CA-02.6 @limite
  Scenario: Enviar el mismo estado que ya tiene no es un error
    When modifico el proyecto "Alfa" con el estado "Activo"
    Then la modificación se realiza correctamente
    And el proyecto modificado está en estado "Activo"

  @CA-02.7 @error
  Scenario Outline: Un estado que no existe es rechazado
    When modifico el proyecto "Alfa" con el estado "<estado>"
    Then se rechaza la modificación con el error "estado inválido"
    And el proyecto "Alfa" no cambió

    Examples:
      | estado  |
      | Cerrado |
      | activo  |

  @CA-02.8 @error
  Scenario: Hay que indicar al menos un dato para modificar
    When modifico el proyecto "Alfa" sin indicar ningún dato
    Then se rechaza la modificación con el error "sin cambios"
    And el proyecto "Alfa" no cambió

  @CA-02.9 @error
  Scenario: No se puede modificar un proyecto que no existe
    When modifico el proyecto con identificador 999 con el nombre "Gamma"
    Then se rechaza la modificación con el error "proyecto no encontrado"

  @CA-02.10 @normal
  Scenario: La modificación queda guardada y el identificador no cambia
    When modifico el proyecto "Alfa" con el nombre "Alfa Renovado"
    Then el proyecto modificado conserva su identificador
    And al recuperarlo por su identificador se llama "Alfa Renovado"

  @CA-02.11 @error
  Scenario: No se puede quitar la fecha de inicio
    When modifico el proyecto "Alfa" quitando la fecha de inicio
    Then se rechaza la modificación con el error "fecha de inicio obligatoria"
    And el proyecto "Alfa" no cambió

  @CA-02.12 @normal
  Scenario: Cambiar el nombre conserva los integrantes y sus roles
    Given el proyecto "Alfa" tiene a "Ana" con los roles "Desarrolladora" y "Tester"
    And el proyecto "Alfa" tiene a "Beto" con el rol "Scrum Master"
    When modifico el proyecto "Alfa" con el nombre "Alfa Renovado"
    Then la modificación se realiza correctamente
    And al recuperarlo, "Ana" tiene los roles "Desarrolladora" y "Tester"
    And al recuperarlo, "Beto" tiene solo el rol "Scrum Master"
