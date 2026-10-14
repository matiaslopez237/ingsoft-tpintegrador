---
name: bdd-go
description: Escribe escenarios BDD Given-When-Then a partir de una spec SDD y los automatiza en Go con godog. Usar cuando el usuario pide escenarios, BDD, Gherkin, archivos .feature, step definitions, o después de aprobar una spec en docs/specs/.
---

# BDD en Go con godog

Los escenarios se derivan de los **criterios de aceptación** de la spec (`docs/specs/HU-XX-*.md`). La cátedra exige que cubran: **casos normales, alternativos, límite y errores**, y que se automaticen siempre que sea posible.

## Convenciones

- Un archivo por historia: `features/HU-XX-nombre-corto.feature`.
- Palabras clave en inglés (`Feature`, `Scenario`, `Given`, `When`, `Then`, `And`, `But`) y el texto en español. Así respeta literalmente el "Given – When – Then" del enunciado.
- **Tags obligatorios** para la trazabilidad:
  - `@HU-XX` en la `Feature`.
  - `@CA-XX.N` en cada `Scenario` (uno o más CAs).
  - Tipo de caso: `@normal`, `@alternativo`, `@limite` o `@error`.
- Antes de terminar, verificá que **todos los CA de la spec** tienen al menos un escenario y que hay al menos un escenario de cada tipo de caso. Si un CA no se puede automatizar, dejalo con el tag `@manual` y explicá por qué en un comentario.

## Ejemplo de .feature

```gherkin
@HU-04
Feature: Planning Poker
  Como integrante del equipo
  quiero votar la estimación de una historia sin ver los votos de los demás
  para evitar sesgos de anclaje

  Background:
    Given un proyecto "TP Integrador" con los integrantes "Ana", "Beto" y "Caro"
    And la historia "HU-07" en el backlog sin estimar

  @CA-04.1 @normal
  Scenario: Los votos permanecen ocultos hasta cerrar la votación
    Given una ronda de votación abierta para "HU-07"
    When "Ana" vota 5 puntos
    Then los votos visibles de la ronda son 0

  @CA-04.3 @alternativo
  Scenario: Se detecta diferencia y se abre una nueva ronda
    Given una ronda con los votos "Ana"=3, "Beto"=13, "Caro"=5
    When se cierra la votación
    Then se informa diferencia entre estimaciones
    And se puede iniciar la ronda 2

  @CA-04.5 @error
  Scenario: Un integrante no puede votar dos veces en la misma ronda
    Given una ronda de votación abierta para "HU-07"
    And "Ana" ya votó 5 puntos
    When "Ana" vota 8 puntos
    Then se rechaza el voto con el error "voto duplicado"
```

## Automatización con godog

Dependencia: `github.com/cucumber/godog` (agregala con `go get github.com/cucumber/godog@latest` si no está en `go.mod`).

Los steps llaman al **núcleo de dominio en Go** (la capa de casos de uso), no a la interfaz web: así los escenarios prueban las reglas de negocio y corren rápido con `go test`.

Estructura sugerida:

```
features/HU-04-planning-poker.feature
test/bdd/bdd_test.go          # runner: una función TestFeatures
test/bdd/hu04_steps_test.go   # steps de una historia
```

Runner (`test/bdd/bdd_test.go`):

```go
package bdd

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../features"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron escenarios BDD")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	registrarStepsHU04(ctx)
	// registrarStepsHU05(ctx) ...
}
```

Steps (`test/bdd/hu04_steps_test.go`): el estado vive en un struct que se crea limpio en cada escenario.

```go
package bdd

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

type pokerWorld struct {
	// referencias al servicio de dominio, la ronda actual, último error, etc.
	ultimoErr error
}

func registrarStepsHU04(ctx *godog.ScenarioContext) {
	w := &pokerWorld{}
	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		*w = pokerWorld{} // estado limpio por escenario
		return c, nil
	})
	ctx.Step(`^"([^"]*)" vota (\d+) puntos$`, w.vota)
	ctx.Step(`^los votos visibles de la ronda son (\d+)$`, w.votosVisibles)
	// ...
}

func (w *pokerWorld) vota(nombre string, puntos int) error {
	// w.ultimoErr = w.ronda.Votar(nombre, puntos)
	return nil
}

func (w *pokerWorld) votosVisibles(esperado int) error {
	// if got := len(w.ronda.VotosVisibles()); got != esperado { return fmt.Errorf(...) }
	return fmt.Errorf("pendiente: implementar")
}
```

Reglas:
- Un step que todavía no está implementado devuelve `godog.ErrPending`; no lo dejes pasando en falso.
- Los steps de error guardan el error en el world y un `Then` lo verifica con `errors.Is` contra el error definido en la spec.
- Correr todo: `go test ./...` (incluye BDD). Solo BDD: `go test ./test/bdd/ -v`.
- El escenario se escribe **antes** de implementar (falla en rojo) y pasa en verde al terminar el ciclo TDD de la historia.

## Al terminar

Mostrá la tabla: CA → escenario(s) → tipo → automatizado sí/no. Si falta cubrir algo, decilo.
