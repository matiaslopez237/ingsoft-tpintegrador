@HU-03
Feature: Registrar integrantes de un proyecto
  Como Product Builder
  quiero registrar los integrantes de un proyecto
  para saber quién participa

  Background:
    Given un proyecto "TP Integrador" sin integrantes

  @CA-03.1 @normal
  Scenario: Agregar un integrante válido
    When se agrega el integrante "Ana" con rol "Desarrolladora"
    Then el proyecto tiene como integrantes a "Ana"

  @CA-03.1 @alternativo
  Scenario: Agregar varios integrantes distintos al mismo proyecto
    When se agrega el integrante "Ana" con rol "Desarrolladora"
    And se agrega el integrante "Beto" con rol "Tester"
    Then el proyecto tiene como integrantes a "Ana" y "Beto"

  @CA-03.2 @error
  Scenario: No se puede agregar un integrante sin nombre
    When se agrega el integrante "" con rol "Desarrolladora"
    Then se rechaza la operación con el error "nombre vacío"
    And el proyecto no tiene integrantes

  @CA-03.5 @error
  Scenario: No se puede agregar un integrante con nombre duplicado
    Given el integrante "Ana" con rol "Desarrolladora" ya está en el proyecto
    When se agrega el integrante "Ana" con rol "Tester"
    Then se rechaza la operación con el error "integrante duplicado"
    And el proyecto tiene como integrantes a "Ana"

  @CA-03.6 @error
  Scenario: No se puede agregar un integrante sin rol
    When se agrega el integrante "Ana" con rol ""
    Then se rechaza la operación con el error "rol vacío"
    And el proyecto no tiene integrantes

  @CA-03.3 @normal
  Scenario: Quitar un integrante existente quedando otros en el proyecto
    Given los integrantes "Ana" y "Beto" ya están en el proyecto
    When se quita al integrante "Ana"
    Then el proyecto tiene como integrantes a "Beto"

  @CA-03.3 @limite
  Scenario: Quitar al único integrante restante deja la lista vacía
    Given el integrante "Ana" con rol "Desarrolladora" ya está en el proyecto
    When se quita al integrante "Ana"
    Then el proyecto no tiene integrantes

  @CA-03.4 @error
  Scenario: No se puede quitar un integrante que no pertenece al proyecto
    Given el integrante "Ana" con rol "Desarrolladora" ya está en el proyecto
    When se quita al integrante "Caro"
    Then se rechaza la operación con el error "integrante no encontrado"
    And el proyecto tiene como integrantes a "Ana"
