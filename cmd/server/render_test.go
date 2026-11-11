package main

import (
	"os"
	"strings"
	"testing"
)

// valorEnRender busca en render.yaml el valor de una variable de entorno (envVars: "- key: X" seguido
// de "value: Y"), sin una librería de YAML: el archivo es chico y lo escribimos nosotros.
func valorEnRender(t *testing.T, clave string) (string, bool) {
	t.Helper()
	b, err := os.ReadFile("../../render.yaml")
	if err != nil {
		t.Fatalf("leer render.yaml: %v", err)
	}
	lineas := strings.Split(string(b), "\n")
	for i, l := range lineas {
		if strings.TrimSpace(l) != "- key: "+clave || i+1 == len(lineas) {
			continue
		}
		valor, ok := strings.CutPrefix(strings.TrimSpace(lineas[i+1]), "value: ")
		return strings.Trim(valor, `"'`), ok
	}
	return "", false
}

// Traza: HU-31 CA-31.11 (sección 5): el despliegue en Render arranca con el proyecto de ejemplo.
func TestRender_CargaElProyectoDeEjemplo(t *testing.T) {
	valor, ok := valorEnRender(t, "CARGAR_EJEMPLO")
	if !ok || valor != "true" {
		t.Errorf("render.yaml: CARGAR_EJEMPLO = %q (definida: %v), se esperaba \"true\"", valor, ok)
	}
}
