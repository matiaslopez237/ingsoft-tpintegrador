@HU-31
Feature: Proyecto de ejemplo precargado
  Como integrante del equipo
  quiero que la aplicación arranque con un proyecto de ejemplo ya cargado
  para mostrar cómo se ve un proyecto completo sin tener que cargarlo a mano

  # Los escenarios corren la carga contra una base SQLite nueva y verifican el resultado con los
  # casos de uso de consulta de las otras historias, como lo vería la aplicación.
  # "arranca la carga del ejemplo con CARGAR_EJEMPLO ..." le pasa a la carga el valor de la
  # variable de entorno, igual que main.go al arrancar.
  # CA-31.10 (un caso de uso que falla a mitad de la carga) se cubre con un test Go del adaptador,
  # con un caso de uso falso que devuelve un error: no se puede provocar con los casos de uso reales.
  # CA-31.11 se cubre con un test Go que lee render.yaml, y se verifica a mano en la URL de Render
  # después de desplegar.

  Background:
    Given una base de datos vacía para el ejemplo

  @CA-31.1 @normal
  Scenario: Con la base vacía y la opción activada se carga el proyecto de ejemplo con su equipo
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then la carga del ejemplo informa que se cargó
    And existe el proyecto de ejemplo "Software Metrics & Estimation" en estado "Activo"
    And el proyecto de ejemplo tiene estos integrantes:
      | nombre           | rol             |
      | Matías López     | Agile Enabler   |
      | Ezequiel Assenza | Product Builder |
      | Nicolás Andino   | Product Builder |
      | Lucas Orellana   | Product Builder |

  @CA-31.2 @normal
  Scenario: El ejemplo tiene el Product Backlog y los sprints del TP
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then el backlog del ejemplo tiene 31 elementos, de "HU-01" a "HU-31", que suman 101 story points
    And los sprints del ejemplo son:
      | número | estado  | inicio     | fin        | planificados | completados |
      | 1      | Cerrado | 2026-10-22 | 2026-10-25 | 19           | 19          |
      | 2      | Cerrado | 2026-10-26 | 2026-10-27 | 26           | 26          |
      | 3      | Cerrado | 2026-10-29 | 2026-10-31 | 28           | 28          |
      | 4      | Cerrado | 2026-11-05 | 2026-11-13 | 28           | 28          |
    And en el sprint 4 del ejemplo están hechas las historias "HU-24", "HU-25", "HU-26", "HU-27", "HU-28", "HU-29" y "HU-31"

  @CA-31.3 @normal
  Scenario: La velocidad del ejemplo sale de los tres sprints cerrados
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then la velocidad del ejemplo es el promedio de 26, 28 y 28

  @CA-31.4 @normal
  Scenario: Los defectos del ejemplo tienen los sprints que corresponden a sus fechas
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then los defectos del ejemplo son:
      | defecto | estado  | sprint de detección | sprint de resolución |
      | DEF-01  | Cerrado | 2                   | 2                    |
      | DEF-02  | Cerrado | 3                   | 3                    |
      | DEF-03  | Cerrado | 3                   | ninguno              |
      | DEF-04  | Abierto | ninguno             | ninguno              |

  @CA-31.5 @normal
  Scenario: El ejemplo tiene esfuerzo registrado y una historia que se votó dos veces
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then las horas del proyecto de ejemplo son 151,5 estimadas y 166 reales, con una desviación de 14,5 h y 9,6 %
    And el historial de rondas de "HU-14" en el ejemplo tiene 2 rondas, la última con estimación acordada 5

  @CA-31.1 @limite
  Scenario Outline: Otros valores booleanos verdaderos también cargan el ejemplo
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "<valor>"
    Then la carga del ejemplo informa que se cargó
    And existe el proyecto de ejemplo "Software Metrics & Estimation" en estado "Activo"

    Examples:
      | valor |
      | 1     |
      | TRUE  |

  @CA-31.6 @alternativo
  Scenario: Si la base ya tiene un proyecto, no se carga el ejemplo
    Given la base del ejemplo ya tiene el proyecto "Propio"
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then la carga del ejemplo informa que no se cargó porque la base ya tiene proyectos
    And la base del ejemplo tiene solo el proyecto "Propio"

  @CA-31.7 @alternativo
  Scenario Outline: Con la opción apagada no se carga nada
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "<valor>"
    Then la carga del ejemplo informa que no se cargó porque la opción está apagada
    And la base del ejemplo no tiene proyectos

    Examples:
      | valor |
      |       |
      | false |
      | 0     |

  @CA-31.8 @error
  Scenario: Un valor que no es booleano no carga el ejemplo, pero no impide arrancar
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "si"
    Then la carga del ejemplo informa que el valor no es booleano
    And la base del ejemplo no tiene proyectos

  @CA-31.9 @limite
  Scenario: Volver a arrancar con el ejemplo cargado no lo duplica
    Given arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    When arranca la carga del ejemplo con CARGAR_EJEMPLO "true"
    Then la carga del ejemplo informa que no se cargó porque la base ya tiene proyectos
    And la base del ejemplo tiene solo el proyecto "Software Metrics & Estimation"
    And el backlog del ejemplo tiene 31 elementos, de "HU-01" a "HU-31", que suman 101 story points
