// Package bdd ejecuta los escenarios Given-When-Then de features/ con godog.
package bdd

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

// TestFeatures corre todos los escenarios. Para correr solo los de una historia, por ejemplo:
//
//	BDD_TAGS=@HU-08 go test ./test/bdd/
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../features"},
			Strict:   true,
			Tags:     os.Getenv("BDD_TAGS"), // vacío: todos los escenarios
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron escenarios BDD")
	}
}

// InitializeScenario registra los steps de cada historia de usuario.
func InitializeScenario(ctx *godog.ScenarioContext) {
	registrarStepsHU01(ctx)
	registrarStepsHU02(ctx)
	registrarStepsHU03(ctx)
	registrarStepsHU04(ctx)
	registrarStepsHU05(ctx)
	registrarStepsHU07(ctx)
	registrarStepsHU08(ctx)
	registrarStepsHU11(ctx)
	registrarStepsHU09(ctx)
	registrarStepsHU18(ctx)
	registrarStepsHU06(ctx)
	registrarStepsHU30(ctx)
}
