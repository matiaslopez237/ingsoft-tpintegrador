@HU-11
Feature: Cerrar un Sprint
  Como Scrum Master
  quiero cerrar un Sprint
  para dar por finalizado el período de trabajo

  # Si no se indica otro, el sprint se cierra desde la ruta de su propio proyecto.
  # Si no se indica otro, los sprints no tienen historias asignadas; las que se asignan quedan
  # sin terminar (Pendiente) y sin estimar, hasta que se completan con HU-10.
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

  # RN-11.3 (P-4, P-13): lo no terminado vuelve al Product Backlog para planificarlo de nuevo con
  # HU-09; lo Hecho queda en el sprint, que es donde se hizo.
  @CA-11.6 @normal
  Scenario: Las historias sin terminar vuelven al Product Backlog y las terminadas quedan en el sprint
    Given la historia "Login de usuarios" de "Alfa" está asignada al sprint 1
    And la historia "Alta de proyectos" de "Alfa" está asignada al sprint 1
    And la historia "Alta de proyectos" de "Alfa" está completada
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" está Cerrado
    And la historia "Login de usuarios" de "Alfa" volvió al backlog sin sprint y sigue "Pendiente"
    And la historia "Alta de proyectos" de "Alfa" sigue asignada al sprint 1

  # RN-11.6 (P-11, P-12): planificados suma todas las historias asignadas al cerrar; completados,
  # solo las Hecho; una historia sin estimar suma 0.
  @CA-11.7 @normal
  Scenario: Al cerrar, el sprint registra sus story points planificados y completados
    Given la historia "Login" de "Alfa", de 5 puntos, está asignada al sprint 1
    And la historia "Login" de "Alfa" está completada
    And la historia "Reportes" de "Alfa", de 3 puntos, está asignada al sprint 1
    And la historia "Ayuda" de "Alfa", sin estimar, está asignada al sprint 1
    And la historia "Ayuda" de "Alfa" está completada
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" registra 8 puntos planificados y 5 completados
    And la historia "Reportes" de "Alfa" volvió al backlog sin sprint y sigue "Pendiente"

  # RN-11.6: el resumen queda fijo aunque la historia devuelta se planifique en el sprint siguiente.
  @CA-11.7 @alternativo
  Scenario: Los puntos registrados no cambian si después se replanifica una historia
    Given la historia "Reportes" de "Alfa", de 3 puntos, está asignada al sprint 1
    And cerré el sprint 1 de "Alfa"
    And el proyecto "Alfa" con el sprint 2 abierto del "2026-10-26" al "2026-11-08"
    When asigno la historia "Reportes" de "Alfa" al sprint 2
    Then el sprint 1 de "Alfa" registra 3 puntos planificados y 0 completados
    And la historia "Reportes" de "Alfa" sigue asignada al sprint 2

  @CA-11.7 @limite
  Scenario: Un sprint sin historias se cierra con 0 puntos
    When cierro el sprint 1 de "Alfa"
    Then el cierre del sprint termina sin errores
    And el sprint 1 de "Alfa" registra 0 puntos planificados y 0 completados

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
