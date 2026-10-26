@HU-08
Feature: Crear un Sprint con Sprint Goal
  Como Scrum Master
  quiero crear un Sprint con fecha de inicio, fin y un Sprint Goal
  para organizar el trabajo del equipo

  Background:
    Given un proyecto "Alfa" que va del "2026-10-01" al "2026-12-31"

  @CA-08.1 @normal
  Scenario: Crear el primer sprint de un proyecto
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con el Sprint Goal "  Gestionar el backlog desde la web  "
    Then el sprint se crea con el número 1
    And el sprint queda en estado "Abierto"
    And el sprint va del "2026-10-26" al "2026-10-27"
    And el sprint tiene el Sprint Goal "Gestionar el backlog desde la web"

  @CA-08.2 @normal
  Scenario: El segundo sprint de un proyecto recibe el número siguiente
    Given "Alfa" ya tiene un sprint del "2026-10-26" al "2026-10-27"
    When creo en "Alfa" un sprint del "2026-10-28" al "2026-10-29" con el Sprint Goal "Estimar con Planning Poker"
    Then el sprint se crea con el número 2

  @CA-08.2 @alternativo
  Scenario: Cada proyecto tiene su propia numeración
    Given "Alfa" ya tiene un sprint del "2026-10-26" al "2026-10-27"
    And un proyecto "Beta" que va del "2026-10-01" al "2026-12-31"
    When creo en "Beta" un sprint del "2026-10-26" al "2026-10-27" con el Sprint Goal "Primer sprint de Beta"
    Then el sprint se crea con el número 1

  @CA-08.3 @error
  Scenario: No se puede crear un sprint en un proyecto que no existe
    When creo un sprint en un proyecto inexistente
    Then se rechaza el sprint con el error "proyecto no encontrado"

  @CA-08.4 @error
  Scenario Outline: Solo se crean sprints en proyectos activos
    Given el proyecto "Alfa" pasa al estado "<estado>"
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con el Sprint Goal "Objetivo"
    Then se rechaza el sprint con el error "proyecto no activo"
    And "Alfa" tiene 0 sprints

    Examples:
      | estado     |
      | En pausa   |
      | Finalizado |
      | Postergado |
      | Cancelado  |

  @CA-08.5 @error
  Scenario Outline: Las fechas son obligatorias y la de fin no puede ser anterior
    When creo en "Alfa" un sprint sin fecha de <fecha>
    Then se rechaza el sprint con el error "fecha de <fecha> obligatoria"
    And "Alfa" tiene 0 sprints

    Examples:
      | fecha  |
      | inicio |
      | fin    |

  @CA-08.5 @error
  Scenario: La fecha de fin no puede ser anterior a la de inicio
    When creo en "Alfa" un sprint del "2026-10-28" al "2026-10-27" con el Sprint Goal "Objetivo"
    Then se rechaza el sprint con el error "fecha de fin anterior"
    And "Alfa" tiene 0 sprints

  @CA-08.5 @limite
  Scenario: Un sprint de un solo día es válido
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-26" con el Sprint Goal "Objetivo"
    Then el sprint se crea con el número 1

  @CA-08.6 @error
  Scenario: El Sprint Goal es obligatorio
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con el Sprint Goal ""
    Then se rechaza el sprint con el error "Sprint Goal vacío"
    And "Alfa" tiene 0 sprints

  @CA-08.6 @error
  Scenario: Un Sprint Goal de solo espacios cuenta como vacío
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con el Sprint Goal "   "
    Then se rechaza el sprint con el error "Sprint Goal vacío"
    And "Alfa" tiene 0 sprints

  @CA-08.6 @error
  Scenario: El Sprint Goal no puede superar los 200 caracteres
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con un Sprint Goal de 201 caracteres
    Then se rechaza el sprint con el error "Sprint Goal demasiado largo"
    And "Alfa" tiene 0 sprints

  @CA-08.6 @limite
  Scenario: Un Sprint Goal de exactamente 200 caracteres es válido
    When creo en "Alfa" un sprint del "2026-10-26" al "2026-10-27" con un Sprint Goal de 200 caracteres
    Then el sprint se crea con el número 1

  @CA-08.7 @error
  Scenario Outline: El sprint tiene que estar dentro de las fechas del proyecto
    When creo en "Alfa" un sprint del "<inicio>" al "<fin>" con el Sprint Goal "Objetivo"
    Then se rechaza el sprint con el error "sprint fuera del proyecto"
    And "Alfa" tiene 0 sprints

    Examples:
      | inicio     | fin        |
      | 2026-09-30 | 2026-10-02 |
      | 2026-12-30 | 2027-01-01 |

  @CA-08.7 @limite
  Scenario Outline: Un sprint justo en los límites del proyecto es válido
    When creo en "Alfa" un sprint del "<inicio>" al "<fin>" con el Sprint Goal "Objetivo"
    Then el sprint se crea con el número 1

    Examples:
      | inicio     | fin        |
      | 2026-10-01 | 2026-10-02 |
      | 2026-12-30 | 2026-12-31 |

  @CA-08.7 @alternativo
  Scenario: En un proyecto sin fecha de fin el sprint puede terminar en cualquier fecha posterior
    Given un proyecto "Gamma" que empieza el "2026-10-01" y no tiene fecha de fin
    When creo en "Gamma" un sprint del "2027-05-01" al "2027-05-14" con el Sprint Goal "Objetivo"
    Then el sprint se crea con el número 1

  @CA-08.8 @error
  Scenario Outline: Los sprints de un proyecto no se superponen
    Given "Alfa" ya tiene un sprint del "2026-10-26" al "2026-10-27"
    When creo en "Alfa" un sprint del "<inicio>" al "<fin>" con el Sprint Goal "Objetivo"
    Then se rechaza el sprint con el error "sprint superpuesto"
    And "Alfa" tiene 1 sprint

    Examples:
      | inicio     | fin        |
      | 2026-10-27 | 2026-10-28 |
      | 2026-10-25 | 2026-10-26 |
      | 2026-10-20 | 2026-10-30 |
      | 2026-10-26 | 2026-10-26 |
      | 2026-10-26 | 2026-10-27 |

  @CA-08.8 @limite
  Scenario: Un sprint que empieza el día siguiente al fin de otro es válido
    Given "Alfa" ya tiene un sprint del "2026-10-26" al "2026-10-27"
    When creo en "Alfa" un sprint del "2026-10-28" al "2026-10-29" con el Sprint Goal "Objetivo"
    Then el sprint se crea con el número 2

  @CA-08.9 @normal
  Scenario: Los sprints del proyecto se listan ordenados por número
    Given "Alfa" ya tiene un sprint del "2026-10-28" al "2026-10-29"
    And "Alfa" ya tiene un sprint del "2026-10-26" al "2026-10-27"
    Then los sprints de "Alfa" son los números "1, 2"
    And el sprint 1 de "Alfa" va del "2026-10-28" al "2026-10-29"
