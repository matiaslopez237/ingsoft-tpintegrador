@HU-23
Feature: Horas estimadas, reales y desviación de un sprint o proyecto
  Como Scrum Master
  quiero ver las horas estimadas, las horas reales y su desviación
  para evaluar la planificación

  # Las horas se escriben con coma decimal, como en la spec. Los pasos llevan "HU-23" para que ninguno
  # coincida con el de otra historia. En las tablas de elementos, las reales se reparten en registros
  # de hasta 24 h, el tope de HU-18 (RN-18.8), y con 0 reales no hay ningún registro.
  # CA-23.16 (respuesta HTTP), CA-23.17 y CA-23.18 (pantallas) se cubren con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para totalizar horas de HU-23, con la integrante "Ana"
    And el proyecto "Alfa" tiene abierto el sprint "Sprint 1" de HU-23

  # RN-23.4: la desviación se calcula una sola vez sobre las sumas. El porcentaje va redondeado a un
  # decimal, con los medios alejándose de cero (RN-19.6).
  Scenario Outline: La desviación de un sprint se calcula sobre la suma de sus elementos
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas     | reales     |
      | Login    | <estimadas_a> | <reales_a> |
      | Reportes | <estimadas_b> | <reales_b> |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan <estimadas> estimadas, <reales> reales y 0 reales sin estimar
    And la desviación de HU-23 es de <horas> horas y de <porcentaje> %
    And las horas de HU-23 cuentan 2 elementos y 0 sin estimar

    @CA-23.1 @normal
    Examples: se suman las horas de los elementos del sprint
      | estimadas_a | reales_a | estimadas_b | reales_b | estimadas | reales | horas | porcentaje |
      | 8           | 10       | 4           | 3        | 12        | 13     | 1     | 8,3        |

    @CA-23.3 @normal
    Examples: el porcentaje es el de las sumas y no el promedio, que daría 50,0
      | estimadas_a | reales_a | estimadas_b | reales_b | estimadas | reales | horas | porcentaje |
      | 10          | 20       | 90          | 90       | 100       | 110    | 10    | 10,0       |

    @CA-23.4 @limite
    Examples: el porcentaje de la suma se redondea a un decimal
      | estimadas_a | reales_a | estimadas_b | reales_b | estimadas | reales | horas | porcentaje |
      | 1           | 2        | 2           | 2        | 3         | 4      | 1     | 33,3       |
      | 20          | 20,5     | 20          | 20       | 40        | 40,5   | 0,5   | 1,3        |
      | 20          | 19,5     | 20          | 20       | 40        | 39,5   | -0,5  | -1,3       |

    @CA-23.5 @limite
    Examples: los desvíos opuestos se compensan en la suma
      | estimadas_a | reales_a | estimadas_b | reales_b | estimadas | reales | horas | porcentaje |
      | 8           | 10       | 8           | 6        | 16        | 16     | 0     | 0,0        |

  # RN-23.10: el proyecto suma todos sus elementos, estén o no en un sprint.
  @CA-23.2 @normal
  Scenario: Las horas de un proyecto incluyen los elementos con sprint y sin sprint
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
      | Reportes | 4         | 3      |
    And el proyecto "Alfa" de HU-23 tiene estos elementos sin sprint:
      | elemento | estimadas | reales |
      | Pagos    | 20        | 15     |
    When consulto las horas de HU-23 del proyecto "Alfa"
    Then las horas de HU-23 informan 32 estimadas, 28 reales y 0 reales sin estimar
    And la desviación de HU-23 es de -4 horas y de -12,5 %
    And las horas de HU-23 cuentan 3 elementos y 0 sin estimar

  # RN-23.5: las reales de un elemento sin estimar suman al total, pero no entran en la desviación.
  @CA-23.6 @alternativo
  Scenario: Un elemento sin estimar suma a las reales y no a la desviación
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas   | reales |
      | Login    | 8           | 10     |
      | Reportes | sin estimar | 5      |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 15 reales y 5 reales sin estimar
    And la desviación de HU-23 es de 2 horas y de 25,0 %
    And las horas de HU-23 cuentan 2 elementos y 1 sin estimar

  # RN-23.5: las 5 h del elemento sin estimar no compensan la desviación del estimado.
  @CA-23.6 @limite
  Scenario: Las horas de un elemento sin estimar no compensan la desviación de otro estimado
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas   | reales |
      | Login    | 8           | 0      |
      | Reportes | sin estimar | 5      |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 5 reales y 5 reales sin estimar
    And la desviación de HU-23 es de -8 horas y de -100,0 %
    And las horas de HU-23 cuentan 2 elementos y 1 sin estimar

  # RN-23.6: sin ninguna estimación no es un error, y las horas trabajadas se siguen informando.
  @CA-23.7 @alternativo
  Scenario: Un sprint sin ningún elemento estimado informa sus reales y ninguna desviación
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas   | reales |
      | Login    | sin estimar | 5      |
      | Reportes | sin estimar | 3      |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 0 estimadas, 8 reales y 8 reales sin estimar
    And las horas de HU-23 están sin estimar y sin desviaciones
    And las horas de HU-23 cuentan 2 elementos y 2 sin estimar

  @CA-23.8 @limite
  Scenario: Un sprint sin elementos informa todo en 0
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 0 estimadas, 0 reales y 0 reales sin estimar
    And las horas de HU-23 están sin estimar y sin desviaciones
    And las horas de HU-23 cuentan 0 elementos y 0 sin estimar

  @CA-23.8 @limite
  Scenario: Un proyecto sin elementos informa todo en 0
    When consulto las horas de HU-23 del proyecto "Alfa"
    Then las horas de HU-23 informan 0 estimadas, 0 reales y 0 reales sin estimar
    And las horas de HU-23 están sin estimar y sin desviaciones
    And las horas de HU-23 cuentan 0 elementos y 0 sin estimar

  # RN-23.2: un elemento estimado y sin registros suma a las estimadas.
  @CA-23.9 @limite
  Scenario: Un elemento estimado y sin esfuerzo suma a las estimadas
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 0      |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 0 reales y 0 reales sin estimar
    And la desviación de HU-23 es de -8 horas y de -100,0 %
    And las horas de HU-23 cuentan 1 elemento y 0 sin estimar

  # RN-23.8: el cierre deja sin sprint a lo no terminado (RN-11.3), y por eso deja de contar en el sprint.
  @CA-23.10 @alternativo
  Scenario: Al cerrar un sprint solo cuentan en él los elementos que siguen asignados
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
      | Reportes | 4         | 3      |
    And "Login" de HU-23 se completa
    And "Reportes" de HU-23 pasa a En progreso
    And el sprint "Sprint 1" de HU-23 se cierra
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 10 reales y 0 reales sin estimar
    And la desviación de HU-23 es de 2 horas y de 25,0 %
    And las horas de HU-23 cuentan 1 elemento y 0 sin estimar

  # RN-23.10: no se pierden horas: el proyecto sigue contando al elemento devuelto al backlog.
  @CA-23.10 @alternativo
  Scenario: El proyecto sigue contando las horas de un elemento devuelto al backlog
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
      | Reportes | 4         | 3      |
    And "Login" de HU-23 se completa
    And "Reportes" de HU-23 pasa a En progreso
    And el sprint "Sprint 1" de HU-23 se cierra
    When consulto las horas de HU-23 del proyecto "Alfa"
    Then las horas de HU-23 informan 12 estimadas, 13 reales y 0 reales sin estimar
    And la desviación de HU-23 es de 1 horas y de 8,3 %
    And las horas de HU-23 cuentan 2 elementos y 0 sin estimar

  # RN-23.8: al replanificarlo, todas sus horas pasan al sprint nuevo.
  @CA-23.10 @alternativo
  Scenario: Las horas de un elemento devuelto al backlog pasan al sprint donde se replanifica
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
      | Reportes | 4         | 3      |
    And "Login" de HU-23 se completa
    And "Reportes" de HU-23 pasa a En progreso
    And el sprint "Sprint 1" de HU-23 se cierra
    And el proyecto "Alfa" tiene abierto el sprint "Sprint 2" de HU-23
    And "Reportes" de HU-23 se planifica en el sprint "Sprint 2"
    When consulto las horas de HU-23 del sprint "Sprint 2"
    Then las horas de HU-23 informan 4 estimadas, 3 reales y 0 reales sin estimar
    And la desviación de HU-23 es de -1 horas y de -25,0 %
    And las horas de HU-23 cuentan 1 elemento y 0 sin estimar

  # RN-23.8 y RN-09.2: un elemento de un sprint cerrado se puede reasignar, aunque esté Hecho.
  @CA-23.10 @limite
  Scenario: Un elemento Hecho reasignado deja de contar en el sprint cerrado
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
    And "Login" de HU-23 se completa
    And el sprint "Sprint 1" de HU-23 se cierra
    And el proyecto "Alfa" tiene abierto el sprint "Sprint 2" de HU-23
    And "Login" de HU-23 se planifica en el sprint "Sprint 2"
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 0 estimadas, 0 reales y 0 reales sin estimar
    And las horas de HU-23 están sin estimar y sin desviaciones
    And las horas de HU-23 cuentan 0 elementos y 0 sin estimar

  @CA-23.10 @limite
  Scenario: Un elemento Hecho reasignado cuenta en el sprint nuevo
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
    And "Login" de HU-23 se completa
    And el sprint "Sprint 1" de HU-23 se cierra
    And el proyecto "Alfa" tiene abierto el sprint "Sprint 2" de HU-23
    And "Login" de HU-23 se planifica en el sprint "Sprint 2"
    When consulto las horas de HU-23 del sprint "Sprint 2"
    Then las horas de HU-23 informan 8 estimadas, 10 reales y 0 reales sin estimar
    And la desviación de HU-23 es de 2 horas y de 25,0 %
    And las horas de HU-23 cuentan 1 elemento y 0 sin estimar

  # RN-23.9: los totales usan los valores actuales; los puntos que guardó el cierre quedan fijos (RN-11.6).
  @CA-23.11 @alternativo
  Scenario: Los totales de un sprint cerrado reflejan los valores actuales
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 10     |
    And "Login" de HU-23 tiene 5 story points
    And "Login" de HU-23 se completa
    And el sprint "Sprint 1" de HU-23 se cierra
    And "Login" de HU-23 pasa a tener 10 horas estimadas
    And "Ana" registró 2 horas para HU-23 en "Login"
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 10 estimadas, 12 reales y 0 reales sin estimar
    And la desviación de HU-23 es de 2 horas y de 20,0 %
    And el sprint "Sprint 1" de HU-23 conserva 5 puntos planificados y 5 completados

  # RN-23.3, RN-23.7 y RN-23.9: suman los registros del elemento, sin importar quién los cargó ni su
  # fecha; no suman los de otro sprint ni los de otro proyecto.
  @CA-23.12 @alternativo
  Scenario: Suman los registros de los elementos del sprint y no los de otros elementos
    Given "Beto" también integra el proyecto "Alfa" de HU-23
    And el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 0      |
    And "Beto" registró 3 horas para HU-23 en "Login"
    And "Ana" registró 2 horas para HU-23 en "Login" con fecha fuera del sprint
    And "Beto" dejó de integrar el proyecto "Alfa" de HU-23
    And el proyecto "Alfa" tiene abierto el sprint "Sprint 2" de HU-23
    And el sprint "Sprint 2" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Reportes | 4         | 9      |
    And el proyecto "Beta" para totalizar horas de HU-23, con la integrante "Caro"
    And el proyecto "Beta" tiene abierto el sprint "Sprint Beta" de HU-23
    And el sprint "Sprint Beta" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Pagos    | 6         | 7      |
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 5 reales y 0 reales sin estimar
    And la desviación de HU-23 es de -3 horas y de -37,5 %
    And las horas de HU-23 cuentan 1 elemento y 0 sin estimar

  # RN-23.12: se consulta en cualquier estado, y la consulta no modifica nada.
  @CA-23.13 @alternativo
  Scenario Outline: Se consultan las horas de un sprint Abierto en un proyecto que no está Activo
    Given el sprint "Sprint 1" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Login    | 8         | 2      |
    And el proyecto "Alfa" de HU-23 pasa a estar "<estado>"
    When consulto las horas de HU-23 del sprint "Sprint 1"
    Then las horas de HU-23 informan 8 estimadas, 2 reales y 0 reales sin estimar
    And la desviación de HU-23 es de -6 horas y de -75,0 %
    And "Login" de HU-23 sigue Pendiente en el sprint "Sprint 1"
    And el sprint "Sprint 1" de HU-23 sigue Abierto

    Examples:
      | estado     |
      | En pausa   |
      | Finalizado |
      | Postergado |
      | Cancelado  |

  @CA-23.14 @error
  Scenario: No se consultan las horas de un proyecto que no existe
    When consulto las horas de HU-23 de un proyecto inexistente
    Then se rechaza la consulta de horas de HU-23 con el error "proyecto no encontrado"

  # El proyecto se busca antes que el sprint (sección 7).
  @CA-23.14 @error
  Scenario: No se consultan las horas de un sprint desde un proyecto que no existe
    When consulto las horas de HU-23 del sprint "Sprint 1" desde un proyecto inexistente
    Then se rechaza la consulta de horas de HU-23 con el error "proyecto no encontrado"

  @CA-23.15 @error
  Scenario: No se consultan las horas de un sprint que no existe
    When consulto las horas de HU-23 de un sprint inexistente de "Alfa"
    Then se rechaza la consulta de horas de HU-23 con el error "sprint no encontrado"

  # Se pide desde la ruta de "Alfa" un sprint que es de "Beta" (RN-23.11).
  @CA-23.15 @error
  Scenario: Un sprint de otro proyecto se trata como inexistente
    Given el proyecto "Beta" para totalizar horas de HU-23, con la integrante "Caro"
    And el proyecto "Beta" tiene abierto el sprint "Sprint Beta" de HU-23
    And el sprint "Sprint Beta" de HU-23 tiene estos elementos:
      | elemento | estimadas | reales |
      | Pagos    | 6         | 7      |
    When consulto las horas de HU-23 del sprint "Sprint Beta" desde el proyecto "Alfa"
    Then se rechaza la consulta de horas de HU-23 con el error "sprint no encontrado"
