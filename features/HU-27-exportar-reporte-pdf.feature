@HU-27
Feature: Exportar el reporte a PDF
  Como Product Owner
  quiero exportar el reporte a un archivo PDF
  para entregarlo como documento formal

  # Los escenarios corren sobre los casos de uso ExportarReporteProyecto y ExportarReporteSprint. Los pasos
  # llevan "HU-27" para que ninguno coincida con el de otra historia.
  #
  # El documento del reporte lo da una fuente en memoria. SUPUESTO: la fuente real es de HU-26, que
  # convierte su reporte al DocumentoPDF de la spec; las secciones de estos escenarios son de ejemplo, no
  # las de HU-26. El generador de PDF es un doble que guarda el documento que recibe y devuelve unos bytes
  # fijos: los escenarios miran ese documento (ya normalizado) y el archivo que devuelve el caso de uso,
  # no el binario del PDF.
  #
  # Los caracteres que no se pueden escribir en un .feature van entre llaves: {TAB} (tabulación),
  # {LF} (salto de línea), {U+XXXX} (el carácter con ese código) y {0xFF} (ese byte suelto, que no es
  # UTF-8 válido).
  #
  # Van a tests Go del adaptador de PDF o del HTTP, porque hay que mirar el binario o la respuesta:
  # CA-27.3 a CA-27.6 (páginas, saltos, tablas largas y tope por celda), CA-27.12 (bytes idénticos),
  # CA-27.19 (IDs inválidos: los rechaza el adaptador HTTP sin llamar al caso de uso), CA-27.20 en su
  # parte HTTP (500 sin detalle), CA-27.21 (botón y revisión visual @manual) y CA-27.24 (%PDF y %%EOF).
  # CA-27.23 (integración con HU-26) queda pendiente hasta que exista HU-26.

  Background:
    Given el reloj de HU-27 marca el 2026-11-07 a las 17:05 UTC
    And el proyecto 3 de HU-27 tiene un reporte titulado "Reporte del proyecto Alfa"

  # RN-27.1, RN-27.2: el documento conserva lo que armó el reporte, y la emisión la pone el caso de uso.
  @CA-27.1 @normal
  Scenario: El documento exportado conserva título, subtítulo y secciones, y lleva la fecha de emisión
    Given el reporte de HU-27 del proyecto 3 tiene el subtítulo "Del 26/10/2026 al 08/11/2026"
    And el reporte de HU-27 del proyecto 3 tiene la sección "Resumen" con el párrafo "El proyecto está Activo."
    And el reporte de HU-27 del proyecto 3 tiene la sección "Historias" con esta tabla:
      | Historia | Estado      | Puntos |
      | Login    | Hecho       | 5      |
      | Reportes | En progreso | 3      |
      | Pagos    | Pendiente   | 8      |
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene el título "Reporte del proyecto Alfa" y el subtítulo "Del 26/10/2026 al 08/11/2026"
    And el documento de HU-27 se emite el 07/11/2026 a las 14:05 en UTC-3
    And el documento de HU-27 tiene estas secciones, en este orden:
      | sección   |
      | Resumen   |
      | Historias |
    And la sección "Resumen" del documento de HU-27 tiene el párrafo "El proyecto está Activo."
    And la sección "Historias" del documento de HU-27 tiene esta tabla:
      | Historia | Estado      | Puntos |
      | Login    | Hecho       | 5      |
      | Reportes | En progreso | 3      |
      | Pagos    | Pendiente   | 8      |
    And el documento de HU-27 no tiene aviso de reemplazos
    And el archivo PDF de HU-27 tiene los bytes que devolvió el generador

  # Criterio 1 de la issue #38: las mismas secciones, con los mismos títulos y en el mismo orden.
  @CA-27.22 @normal
  Scenario: El PDF conserva los títulos y el orden de las secciones del reporte, aunque estén vacías
    Given el reporte de HU-27 del proyecto 3 tiene la sección "Resumen" con el párrafo "Sprint 2 cerrado."
    And el reporte de HU-27 del proyecto 3 tiene la sección "Historias y estimaciones" con esta tabla:
      | Historia | Puntos |
      | Login    | 5      |
    And el reporte de HU-27 del proyecto 3 tiene la sección "Esfuerzo" sin contenido
    And el reporte de HU-27 del proyecto 3 tiene la sección "Métricas" con el párrafo "Velocidad: 8 puntos."
    And el reporte de HU-27 del proyecto 3 tiene la sección "Defectos" con esta tabla:
      | Defecto | Severidad |
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene estas secciones, en este orden:
      | sección                 |
      | Resumen                 |
      | Historias y estimaciones |
      | Esfuerzo                |
      | Métricas                |
      | Defectos                |
    And la sección "Esfuerzo" del documento de HU-27 no tiene contenido
    And la sección "Defectos" del documento de HU-27 tiene esta tabla:
      | Defecto | Severidad |

  # RN-27.7: un documento sin secciones es válido.
  @CA-27.7 @limite
  Scenario: Un reporte sin subtítulo ni secciones se exporta igual
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene el título "Reporte del proyecto Alfa"
    And el documento de HU-27 no tiene subtítulo
    And el documento de HU-27 no tiene secciones
    And el archivo PDF de HU-27 tiene los bytes que devolvió el generador

  # RN-27.8, pasos 1 a 3: lo representable queda igual, los espacios de Unicode pasan a espacio común y
  # los invisibles de formato se quitan. Nada de eso cuenta para el aviso (RN-27.9).
  Scenario Outline: Los textos que el PDF puede mostrar no se reemplazan ni generan aviso
    Given el título del reporte de HU-27 del proyecto 3 es "<original>"
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene el título "<normalizado>"
    And el documento de HU-27 no tiene aviso de reemplazos

    @CA-27.2 @normal
    Examples: acentos, eñes y signos del español
      | original                                       | normalizado                                    |
      | Proyecto Ñandú: planificación ¿cumplida? ¡Sí! | Proyecto Ñandú: planificación ¿cumplida? ¡Sí! |
      | áéíóúüñÁÉÍÓÚÜÑ€°ªº–—“”…→≤≥%                    | áéíóúüñÁÉÍÓÚÜÑ€°ªº–—“”…→≤≥%                    |

    @CA-27.10 @limite
    Examples: espacios de Unicode e invisibles de formato
      | original              | normalizado  |
      | Sprint{TAB}final      | Sprint final |
      | Sprint{LF}final       | Sprint final |
      | Sprint{U+00A0}final   | Sprint final |
      | Sprint{U+200B}final   | Sprintfinal  |
      | ¿cumplida?            | ¿cumplida?   |

  # RN-27.8, paso 4, y RN-27.9: lo que no se puede mostrar pasa a "?", varios seguidos dan uno solo, y
  # el documento avisa cuántos "?" se insertaron. No es un error: la exportación sigue.
  Scenario Outline: Los caracteres que el PDF no puede mostrar se reemplazan por "?" y se avisa
    Given el título del reporte de HU-27 del proyecto 3 es "<original>"
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene el título "<normalizado>"
    And el documento de HU-27 tiene el aviso "<aviso>"
    And el archivo PDF de HU-27 tiene los bytes que devolvió el generador

    @CA-27.8 @alternativo
    Examples: emojis e ideogramas
      | original           | normalizado    | aviso                                                                                  |
      | Sprint 😀 final    | Sprint ? final | Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?». |
      | Equipo 😀😀 y 中文 | Equipo ? y ?   | Aviso: 2 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?». |

    @CA-27.9 @limite
    Examples: secuencias de emojis
      | original                     | normalizado | aviso                                                                                  |
      | Familia 👨{U+200D}👩{U+200D}👧 | Familia ?   | Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?». |
      | Amor ❤{U+FE0F}               | Amor ?      | Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?». |
      | 😀 😀                        | ? ?         | Aviso: 2 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?». |

    @CA-27.10 @limite
    Examples: caracteres de control y bytes que no son UTF-8 válido
      | original            | normalizado  | aviso                                                                                  |
      | Sprint{U+0007}final | Sprint?final | Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?». |
      | Sprint{0xFF}final   | Sprint?final | Aviso: 1 símbolo que no se puede mostrar en PDF, como los emojis, se reemplazó por «?». |

  # RN-27.8 se aplica a todos los textos: subtítulo, títulos de sección, párrafos, encabezados y celdas.
  @CA-27.8 @alternativo
  Scenario: Los caracteres se reemplazan en todos los textos del documento
    Given el reporte de HU-27 del proyecto 3 tiene el subtítulo "Avance 🚀"
    And el reporte de HU-27 del proyecto 3 tiene la sección "Equipo 😀" con el párrafo "Listo ✅"
    And el reporte de HU-27 del proyecto 3 tiene la sección "Estados" con esta tabla:
      | Estado 🔥 |
      | Hecho 👍  |
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 tiene el subtítulo "Avance ?"
    And la sección "Equipo ?" del documento de HU-27 tiene el párrafo "Listo ?"
    And la sección "Estados" del documento de HU-27 tiene esta tabla:
      | Estado ? |
      | Hecho ?  |
    And el documento de HU-27 tiene el aviso "Aviso: 5 símbolos que no se pueden mostrar en PDF, como los emojis, se reemplazaron por «?»."

  # RN-27.10 y RN-27.12: la emisión es la hora del reloj en UTC-3, y su fecha va en el nombre del archivo.
  Scenario Outline: La fecha de emisión va en hora de Argentina y da la fecha del nombre del archivo
    Given el reloj de HU-27 marca el <fecha_utc> a las <hora_utc> UTC
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el documento de HU-27 se emite el <fecha_local> a las <hora_local> en UTC-3
    And el archivo PDF de HU-27 se llama "<archivo>"

    @CA-27.13 @CA-27.14 @normal
    Examples: emisión en el mismo día
      | fecha_utc  | hora_utc | fecha_local | hora_local | archivo                           |
      | 2026-11-07 | 17:05    | 07/11/2026  | 14:05      | reporte-proyecto-3-2026-11-07.pdf |

    @CA-27.13 @limite
    Examples: el huso horario cambia el día
      | fecha_utc  | hora_utc | fecha_local | hora_local | archivo                           |
      | 2026-11-08 | 02:30    | 07/11/2026  | 23:30      | reporte-proyecto-3-2026-11-07.pdf |
      | 2026-11-08 | 03:00    | 08/11/2026  | 00:00      | reporte-proyecto-3-2026-11-08.pdf |

  # RN-27.12: el archivo del sprint lleva el id del proyecto y el del sprint.
  @CA-27.15 @alternativo
  Scenario: El reporte de un sprint se exporta con su propio nombre de archivo
    Given el proyecto 3 de HU-27 tiene el sprint 12 con un reporte titulado "Reporte del Sprint 2"
    When exporto a PDF el reporte de HU-27 del sprint 12 del proyecto 3
    Then el documento de HU-27 tiene el título "Reporte del Sprint 2"
    And el documento de HU-27 se emite el 07/11/2026 a las 14:05 en UTC-3
    And el archivo PDF de HU-27 se llama "reporte-proyecto-3-sprint-12-2026-11-07.pdf"

  # RN-27.12: el nombre del archivo no incluye el nombre del proyecto, que puede tener cualquier carácter.
  @CA-27.16 @limite
  Scenario: El nombre del archivo no depende del nombre del proyecto
    Given el título del reporte de HU-27 del proyecto 3 es "Ñandú / Q4: «final» 😀"
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then el archivo PDF de HU-27 se llama "reporte-proyecto-3-2026-11-07.pdf"

  # RN-27.11: los errores de quien arma el documento se detectan antes de generar. El título se valida
  # después de normalizarlo (RN-27.8), así que uno hecho solo de invisibles queda vacío.
  Scenario Outline: Un documento mal armado no llega al generador
    Given el reporte de HU-27 del proyecto 3 es inválido por tener <defecto>
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then se rechaza la exportación de HU-27 por documento inválido
    And el generador de PDF de HU-27 no se usó

    @CA-27.11 @error
    Examples: errores de armado
      | defecto                              |
      | un título vacío                      |
      | una sección sin título               |
      | un bloque con párrafo y tabla        |
      | una tabla sin columnas               |
      | una fila con menos celdas que columnas |
      | una fila con más celdas que columnas |

    @CA-27.11 @limite
    Examples: títulos que quedan vacíos al normalizarlos
      | defecto                                  |
      | un título solo de espacios               |
      | un título con solo caracteres invisibles |

  # RN-27.13: el proyecto tiene que existir.
  Scenario Outline: No se exporta el reporte de un proyecto inexistente
    When <exportacion>
    Then se rechaza la exportación de HU-27 con el error "proyecto no encontrado"
    And el generador de PDF de HU-27 no se usó

    @CA-27.17 @error
    Examples:
      | exportacion                                                |
      | exporto a PDF el reporte de HU-27 del proyecto 99          |
      | exporto a PDF el reporte de HU-27 del sprint 12 del proyecto 99 |

  # RN-27.13: un sprint de otro proyecto se trata como inexistente. SUPUESTO: la regla la aplica la fuente
  # de HU-26 (spec, sección 5); este escenario comprueba que el caso de uso devuelve ese error y no genera.
  Scenario Outline: No se exporta el reporte de un sprint inexistente o de otro proyecto
    Given el proyecto 4 de HU-27 tiene un reporte titulado "Reporte del proyecto Beta"
    And el proyecto 4 de HU-27 tiene el sprint 13 con un reporte titulado "Reporte del Sprint 1 de Beta"
    When exporto a PDF el reporte de HU-27 del sprint <sprint> del proyecto 3
    Then se rechaza la exportación de HU-27 con el error "sprint no encontrado"
    And el generador de PDF de HU-27 no se usó

    @CA-27.18 @error
    Examples:
      | sprint |
      | 99     |
      | 13     |

  # Sección 7: una falla al armar el reporte o al generar el PDF no deja un archivo a medias. El 500 sin
  # detalle es del adaptador HTTP y se prueba ahí.
  Scenario Outline: Si falla el armado del reporte o el generador, no se obtiene ningún archivo
    Given <componente> de HU-27 falla
    When exporto a PDF el reporte de HU-27 del proyecto 3
    Then la exportación de HU-27 falla por un error interno

    @CA-27.20 @error
    Examples:
      | componente               |
      | la fuente de documentos  |
      | el generador de PDF      |
