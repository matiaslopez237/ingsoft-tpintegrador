@HU-30
Feature: Interfaz web base
  Como usuario del sistema
  quiero una interfaz web con una base común y pantallas de proyectos e integrantes
  para operar el sistema desde el navegador

  @CA-30.10 @normal
  Scenario: La página de inicio lleva a la lista de proyectos
    When abro la página "/"
    Then me redirige a "/web/proyectos"

  @CA-30.1 @normal
  Scenario: Las pantallas comparten el layout común
    When abro la página "/web/proyectos"
    Then la respuesta tiene estado 200
    And la página tiene el layout común
    And la hoja de estilos se sirve en "/static/estilos.css"

  @CA-30.2 @normal
  Scenario: La lista muestra los proyectos ordenados por nombre
    Given están guardados los proyectos:
      | nombre |
      | beta   |
      | Alfa   |
      | Gamma  |
    When abro la página "/web/proyectos"
    Then veo los proyectos en este orden:
      | nombre |
      | Alfa   |
      | beta   |
      | Gamma  |

  @CA-30.2 @limite
  Scenario: Sin proyectos se muestra la lista vacía
    When abro la página "/web/proyectos"
    Then veo el texto "Todavía no hay proyectos"

  @CA-30.3 @normal
  Scenario: Crear un proyecto desde el formulario
    When envío el formulario de nuevo proyecto con nombre "Alfa", inicio "2026-03-01" y fin "2026-11-30"
    Then me redirige al detalle del proyecto "Alfa"
    And hay 1 proyectos guardados

  @CA-30.3 @error
  Scenario: El formulario de nuevo proyecto muestra el error de dominio
    When envío el formulario de nuevo proyecto con nombre "", inicio "2026-03-01" y fin ""
    Then la respuesta tiene estado 400
    And veo el texto "el nombre del proyecto es obligatorio"
    And hay 0 proyectos guardados

  @CA-30.3 @error
  Scenario: El formulario de nuevo proyecto avisa una fecha con formato inválido
    When envío el formulario de nuevo proyecto con nombre "Alfa", inicio "01/03/2026" y fin ""
    Then la respuesta tiene estado 400
    And veo el texto "fecha_inicio debe tener el formato AAAA-MM-DD"

  @CA-30.4 @normal
  Scenario: El detalle muestra los datos del proyecto y sus integrantes
    Given está guardado el proyecto "Alfa" con el integrante "Ana" con rol "Tester"
    When abro el detalle del proyecto "Alfa"
    Then la respuesta tiene estado 200
    And veo el texto "Alfa"
    And veo el texto "Activo"
    And el detalle de "Alfa" muestra al integrante "Ana" con rol "Tester"

  @CA-30.4 @error
  Scenario: El detalle de un proyecto inexistente es una página 404
    When abro la página "/web/proyectos/999"
    Then la respuesta tiene estado 404
    And veo el texto "proyecto no encontrado"

  @CA-30.5 @normal
  Scenario: Editar un proyecto y quitarle la fecha de fin
    Given está guardado el proyecto "Alfa"
    When envío el formulario de edición de "Alfa" con nombre "Alfa 2", inicio "2026-03-01", fin "" y estado "En pausa"
    Then me redirige al detalle del proyecto "Alfa 2"
    And el proyecto "Alfa 2" queda en estado "En pausa" y sin fecha de fin

  @CA-30.5 @error
  Scenario: El formulario de edición muestra el error de dominio
    Given está guardado el proyecto "Alfa"
    When envío el formulario de edición de "Alfa" con nombre "Alfa", inicio "2026-03-01", fin "2026-02-01" y estado "Activo"
    Then la respuesta tiene estado 400
    And veo el texto "la fecha de fin no puede ser anterior a la de inicio"

  @CA-30.6 @normal
  Scenario: Agregar un integrante desde el detalle
    Given está guardado el proyecto "Alfa"
    When agrego desde la web al integrante "Ana" con rol "Tester" en "Alfa"
    Then me redirige al detalle del proyecto "Alfa"
    And el detalle de "Alfa" muestra al integrante "Ana" con rol "Tester"

  @CA-30.6 @error
  Scenario: Agregar un integrante repetido muestra el error en el detalle
    Given está guardado el proyecto "Alfa" con el integrante "Ana" con rol "Tester"
    When agrego desde la web al integrante "ana" con rol "Dev" en "Alfa"
    Then la respuesta tiene estado 400
    And veo el texto "ya hay un integrante con ese nombre en el proyecto"

  @CA-30.7 @normal
  Scenario: Quitar un integrante desde el detalle
    Given está guardado el proyecto "Alfa" con el integrante "Ana" con rol "Tester"
    When quito desde la web al integrante "Ana" de "Alfa"
    Then me redirige al detalle del proyecto "Alfa"
    And el detalle de "Alfa" no muestra al integrante "Ana"

  @CA-30.8 @normal
  Scenario: Sumar un rol a un integrante desde el detalle
    Given está guardado el proyecto "Alfa" con el integrante "Ana" con rol "Tester"
    When le sumo desde la web el rol "Desarrolladora" a "Ana" en "Alfa"
    Then me redirige al detalle del proyecto "Alfa"
    And el detalle de "Alfa" muestra al integrante "Ana" con rol "Desarrolladora"

  @CA-30.8 @error
  Scenario: Sumar un rol repetido muestra el error en el detalle
    Given está guardado el proyecto "Alfa" con el integrante "Ana" con rol "Tester"
    When le sumo desde la web el rol "tester" a "Ana" en "Alfa"
    Then la respuesta tiene estado 400
    And veo el texto "el integrante ya tiene ese rol"

  @CA-30.9 @limite
  Scenario: El contenido del usuario se muestra escapado
    Given están guardados los proyectos:
      | nombre                        |
      | <script>alert(1)</script>     |
    When abro la página "/web/proyectos"
    Then la página no contiene "<script>alert(1)</script>"
    And la página contiene "&lt;script&gt;alert(1)&lt;/script&gt;"
