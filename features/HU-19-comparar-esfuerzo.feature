@HU-19
Feature: Comparar esfuerzo estimado contra real
  Como Scrum Master
  quiero comparar el esfuerzo estimado contra el real de una historia
  para evaluar la precisión de las estimaciones

  # Las horas se escriben con coma decimal, como en la spec. Los elementos nacen Pendiente, sin sprint
  # y sin horas estimadas (RN-19.3); las horas estimadas se cargan editando el elemento (HU-06).
  # CA-19.16 (base creada antes de HU-19) se cubre con tests de persistencia, y CA-19.17 (respuesta
  # HTTP) y CA-19.18 (pantalla), con tests del adaptador.

  Background:
    Given el proyecto "Alfa" para comparar esfuerzo, con la integrante "Ana"
    And el proyecto "Alfa" tiene el elemento "Login" sin horas estimadas

  # "acumula" reparte las horas en registros de hasta 24 h, el tope de HU-18 (RN-18.8); con 0 horas
  # no hay ningún registro. El porcentaje va redondeado a un decimal (RN-19.6).
  Scenario Outline: La desviación se calcula en horas y en porcentaje
    Given "Login" tiene cargadas <estimadas> horas estimadas
    And "Login" acumula <reales> horas de esfuerzo
    When comparo el esfuerzo de "Login"
    Then la comparación informa <estimadas> horas estimadas, <reales> horas reales y una desviación de <horas> horas y de <porcentaje> %

    @CA-19.1 @normal
    Examples: se trabajó más de lo estimado
      | estimadas | reales | horas | porcentaje |
      | 8         | 10     | 2     | 25,0       |

    @CA-19.2 @normal
    Examples: se trabajó menos de lo estimado
      | estimadas | reales | horas | porcentaje |
      | 8         | 6      | -2    | -25,0      |

    @CA-19.3 @limite
    Examples: se trabajó exactamente lo estimado
      | estimadas | reales | horas | porcentaje |
      | 8         | 8      | 0     | 0,0        |

    @CA-19.5 @limite
    Examples: con estimación y sin esfuerzo registrado
      | estimadas | reales | horas | porcentaje |
      | 8         | 0      | -8    | -100,0     |

    @CA-19.8 @limite
    Examples: el porcentaje se redondea a un decimal, con los medios alejándose de cero
      | estimadas | reales | horas | porcentaje |
      | 3         | 4      | 1     | 33,3       |
      | 40        | 40,5   | 0,5   | 1,3        |
      | 40        | 39,5   | -0,5  | -1,3       |

  # RN-19.4: las horas reales son la suma de todos los registros del elemento.
  @CA-19.1 @normal
  Scenario: Las horas reales suman todos los registros del elemento
    Given "Login" tiene cargadas 8 horas estimadas
    And "Ana" registró 6 horas de esfuerzo en "Login"
    And "Ana" registró 4 horas de esfuerzo en "Login"
    When comparo el esfuerzo de "Login"
    Then la comparación informa 8 horas estimadas, 10 horas reales y una desviación de 2 horas y de 25,0 %

  # RN-19.7: sin estimar no es un error. Tomar la estimación como 0 mostraría todo como desviación.
  @CA-19.4 @alternativo
  Scenario: Un elemento sin estimar y con esfuerzo informa solo sus horas reales
    Given "Login" acumula 10 horas de esfuerzo
    When comparo el esfuerzo de "Login"
    Then la comparación está sin estimar, con 10 horas reales y sin desviaciones

  @CA-19.6 @limite
  Scenario: Un elemento sin estimar y sin esfuerzo informa 0 horas reales
    When comparo el esfuerzo de "Login"
    Then la comparación está sin estimar, con 0 horas reales y sin desviaciones

  # RN-19.8: un elemento no puede tener estimado 0 (RN-19.2), pero HU-23 le va a pasar sumas a la
  # función del dominio, que sí pueden dar 0. Por eso este escenario la llama directamente.
  @CA-19.7 @limite
  Scenario: Con estimado 0 hay desviación en horas y no hay porcentaje
    When calculo la comparación en el dominio con 0 horas estimadas y 10 horas reales
    Then la comparación no está sin estimar, con 10 horas reales, una desviación de 10 horas y sin porcentaje

  @CA-19.9 @error
  Scenario: No se puede comparar en un proyecto que no existe
    When comparo el esfuerzo de "Login" desde un proyecto inexistente
    Then se rechaza la comparación con el error "proyecto no encontrado"

  @CA-19.9 @error
  Scenario: No se puede comparar un elemento que no existe
    When comparo el esfuerzo de un elemento inexistente de "Alfa"
    Then se rechaza la comparación con el error "elemento no encontrado"

  # Se pide desde la ruta de "Alfa" un elemento que es de "Beta" (RN-19.10).
  @CA-19.10 @error
  Scenario: Un elemento de otro proyecto se trata como inexistente
    Given el proyecto "Beta" para comparar esfuerzo, con la integrante "Caro"
    And el proyecto "Beta" tiene el elemento "Pagos" con 8 horas estimadas
    And "Caro" registró 2 horas de esfuerzo en "Pagos"
    When comparo el esfuerzo de "Pagos" desde el proyecto "Alfa"
    Then se rechaza la comparación con el error "elemento no encontrado"

  # RN-19.4 y RN-18.3: las horas de quien ya no está en el proyecto se conservan; las de otro
  # elemento del mismo proyecto no suman.
  @CA-19.11 @alternativo
  Scenario: Suman los registros de un integrante quitado y no los de otro elemento
    Given "Beto" también integra el proyecto "Alfa"
    And el proyecto "Alfa" tiene el elemento "Reportes" sin horas estimadas
    And "Login" tiene cargadas 8 horas estimadas
    And "Beto" registró 3 horas de esfuerzo en "Login"
    And "Ana" registró 2 horas de esfuerzo en "Login"
    And "Ana" registró 5 horas de esfuerzo en "Reportes"
    And "Beto" dejó de integrar el proyecto "Alfa"
    When comparo el esfuerzo de "Login"
    Then la comparación informa 8 horas estimadas, 5 horas reales y una desviación de -3 horas y de -37,5 %

  # RN-19.4: el tope de 24 h es por registro, no para la suma.
  @CA-19.11 @limite
  Scenario: La suma de las horas reales no tiene tope
    Given "Login" tiene cargadas 8 horas estimadas
    And "Ana" registró 24 horas de esfuerzo en "Login"
    And "Ana" registró 24 horas de esfuerzo en "Login"
    And "Ana" registró 8 horas de esfuerzo en "Login"
    When comparo el esfuerzo de "Login"
    Then la comparación informa 8 horas estimadas, 56 horas reales y una desviación de 48 horas y de 600,0 %

  # RN-19.11: se compara con el elemento en cualquier estado y con el proyecto en cualquier estado.
  @CA-19.12 @alternativo
  Scenario Outline: Se compara un elemento Pendiente y sin sprint de un proyecto que no está Activo
    Given "Login" tiene cargadas 8 horas estimadas
    And "Login" acumula 2 horas de esfuerzo
    And "Login" está Pendiente y sin sprint
    And el proyecto "Alfa" que se compara pasa a estar "<estado>"
    When comparo el esfuerzo de "Login"
    Then la comparación informa 8 horas estimadas, 2 horas reales y una desviación de -6 horas y de -75,0 %

    Examples:
      | estado     |
      | En pausa   |
      | Finalizado |
      | Postergado |
      | Cancelado  |

  @CA-19.13 @normal
  Scenario: Se cargan las horas estimadas editando el elemento
    When edito las horas estimadas de "Login" a 8
    Then la edición de las horas estimadas termina sin errores
    And "Login" tiene 8 horas estimadas

  # RN-19.3: un campo que no se envía en la edición conserva su valor.
  @CA-19.13 @alternativo
  Scenario: Editar otro campo no cambia las horas estimadas
    Given "Login" tiene cargadas 8 horas estimadas
    When edito los story points de "Login" a 5
    Then la edición de las horas estimadas termina sin errores
    And "Login" tiene 5 story points
    And "Login" tiene 8 horas estimadas

  @CA-19.13 @alternativo
  Scenario: Se quitan las horas estimadas de un elemento
    Given "Login" tiene cargadas 8 horas estimadas
    When quito las horas estimadas de "Login"
    Then la edición de las horas estimadas termina sin errores
    And "Login" está sin horas estimadas

  # RN-19.3: se puede cambiar aunque tenga esfuerzo, y la comparación usa el valor actual.
  @CA-19.13 @alternativo
  Scenario: Se cambia la estimación de un elemento con esfuerzo y la comparación usa el valor nuevo
    Given "Login" tiene cargadas 8 horas estimadas
    And "Login" acumula 10 horas de esfuerzo
    When edito las horas estimadas de "Login" a 10
    And comparo el esfuerzo de "Login"
    Then la edición de las horas estimadas termina sin errores
    And la comparación informa 10 horas estimadas, 10 horas reales y una desviación de 0 horas y de 0,0 %

  @CA-19.13 @alternativo
  Scenario: Se quita la estimación de un elemento con esfuerzo y conserva sus horas reales
    Given "Login" tiene cargadas 8 horas estimadas
    And "Login" acumula 10 horas de esfuerzo
    When quito las horas estimadas de "Login"
    And comparo el esfuerzo de "Login"
    Then la edición de las horas estimadas termina sin errores
    And la comparación está sin estimar, con 10 horas reales y sin desviaciones

  @CA-19.14 @limite
  Scenario Outline: Se aceptan las horas estimadas en los extremos válidos
    When edito las horas estimadas de "Login" a <horas>
    Then la edición de las horas estimadas termina sin errores
    And "Login" tiene <horas> horas estimadas

    Examples:
      | horas |
      | 0,5   |
      | 1000  |

  # RN-19.2: se valida en este orden: mayor a 0, múltiplo de 0,5 y máximo 1000. Si hay error, el
  # elemento conserva la estimación que tenía.
  @CA-19.14 @error
  Scenario Outline: Se rechazan las horas estimadas inválidas
    Given "Login" tiene cargadas 8 horas estimadas
    When edito las horas estimadas de "Login" a <horas>
    Then se rechaza la edición de las horas estimadas con el error "<error>"
    And "Login" tiene 8 horas estimadas

    Examples:
      | horas  | error                             |
      | 0      | horas inválidas                   |
      | -1     | horas inválidas                   |
      | 1,25   | horas no múltiplo de media hora   |
      | 1000,5 | horas estimadas exceden el máximo |

  # RN-19.1: las horas estimadas no se derivan de los story points.
  @CA-19.15 @alternativo
  Scenario: Un elemento con story points y sin horas estimadas está sin estimar
    Given el proyecto "Alfa" tiene el elemento "Pagos" con 5 story points y sin horas estimadas
    When comparo el esfuerzo de "Pagos"
    Then la comparación está sin estimar, con 0 horas reales y sin desviaciones

  @CA-19.15 @alternativo
  Scenario: Cambiar los story points no cambia las horas estimadas
    Given el proyecto "Alfa" tiene el elemento "Pagos" con 5 story points y sin horas estimadas
    And "Pagos" tiene cargadas 8 horas estimadas
    When edito los story points de "Pagos" a 13
    Then la edición de las horas estimadas termina sin errores
    And "Pagos" tiene 13 story points
    And "Pagos" tiene 8 horas estimadas

  @CA-19.15 @normal
  Scenario: Un elemento recién creado nace sin horas estimadas
    When creo el elemento "Ayuda" en el proyecto "Alfa" que se compara
    Then "Ayuda" está sin horas estimadas
