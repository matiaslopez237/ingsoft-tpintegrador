@HU-11
Feature: Cerrar un Sprint
  Como Scrum Master
  quiero cerrar un Sprint
  para dar por finalizado el período de trabajo

  # Si no se indica otro, el sprint se cierra desde la ruta de su propio proyecto.
  # Si no se indica otro, los sprints no tienen historias asignadas; las que se asignan quedan
  # sin terminar (Pendiente).
  # CA-11.4 (respuesta HTTP) y CA-11.5 (pantalla) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" con el sprint 1 abierto del "2026-10-12" al "2026-10-25"

  # RN-11.2: el cierre no mira la fecha de hoy. Un sprint que termina en 2099 todavía no llegó a
  # su fin, y uno que empieza en 2099 todavía no empezó.
  @CA-11.1
  Scenario Outline: Un sprint abierto se cierra en cualquier momento y conserva sus datos
    Given el proyecto "Alfa" con el sprint 2 abierto del "<inicio>" al "<fin>"
    When cierro el sprint 2 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 2 de "Alfa" está Cerrado
    And el sprint 2 de "Alfa" conserva sus datos

    @normal
    Examples: el sprint ya terminó
      | inicio     | fin        |
      | 2020-01-01 | 2020-01-14 |

    @limite
    Examples: la fecha de fin todavía no llegó, o el sprint todavía no empezó
      | inicio     | fin        |
      | 2026-10-28 | 2099-12-31 |
      | 2099-01-01 | 2099-01-14 |

  @CA-11.1 @limite
  Scenario: Se cierra un sprint sin historias asignadas
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" está Cerrado

  # RN-11.2 y P-10 de la spec: cerrar no depende del estado del proyecto.
  @CA-11.1 @limite
  Scenario Outline: Se cierra un sprint de un proyecto que no está Activo
    Given el proyecto "Alfa" está "<estado>"
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" está Cerrado

    Examples:
      | estado     |
      | En pausa   |
      | Finalizado |
      | Postergado |
      | Cancelado  |

  # RN-11.3: cerrar no toca las historias; se pueden reasignar después por RN-09.2 (HU-09).
  @CA-11.6 @normal
  Scenario: Las historias sin terminar siguen asignadas al sprint cerrado
    Given la historia "Login de usuarios" de "Alfa" está asignada al sprint 1
    And la historia "Alta de proyectos" de "Alfa" está asignada al sprint 1
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" está Cerrado
    And la historia "Login de usuarios" de "Alfa" sigue asignada al sprint 1
    And la historia "Alta de proyectos" de "Alfa" sigue asignada al sprint 1

  @CA-11.2 @error
  Scenario: No se puede cerrar un sprint que ya está cerrado
    Given el proyecto "Alfa" con el sprint 2 cerrado del "2026-10-26" al "2026-10-27"
    When cierro el sprint 2 de "Alfa"
    Then se rechaza el cierre del sprint con el error "sprint ya cerrado"
    And el sprint 2 de "Alfa" está Cerrado
    And el sprint 2 de "Alfa" conserva sus datos

  @CA-11.2 @alternativo
  Scenario: Un sprint que se acaba de cerrar no se puede volver a cerrar
    Given cerré el sprint 1 de "Alfa"
    When cierro el sprint 1 de "Alfa"
    Then se rechaza el cierre del sprint con el error "sprint ya cerrado"
    And el sprint 1 de "Alfa" está Cerrado

  @CA-11.3 @error
  Scenario: No se puede cerrar un sprint de un proyecto que no existe
    When cierro el sprint 1 de "Alfa" desde un proyecto inexistente
    Then se rechaza el cierre del sprint con el error "proyecto no encontrado"
    And el sprint 1 de "Alfa" está Abierto

  @CA-11.3 @error
  Scenario: No se puede cerrar un sprint que no existe
    When cierro un sprint inexistente de "Alfa"
    Then se rechaza el cierre del sprint con el error "sprint no encontrado"

  # Se intenta cerrar desde la ruta de "Alfa" un sprint que es de "Beta" (RN-11.4).
  @CA-11.3 @error
  Scenario: Un sprint de otro proyecto se trata como inexistente
    Given el proyecto "Beta" con el sprint 1 abierto del "2026-10-12" al "2026-10-25"
    When cierro el sprint 1 de "Beta" desde el proyecto "Alfa"
    Then se rechaza el cierre del sprint con el error "sprint no encontrado"
    And el sprint 1 de "Beta" está Abierto
