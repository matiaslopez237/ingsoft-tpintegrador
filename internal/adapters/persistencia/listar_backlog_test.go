package persistencia

import (
	"context"
	"reflect"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// El adaptador tiene que cumplir el puerto del caso de uso (se verifica al compilar).
var _ app.ListadorElementos = (*RepositorioBacklog)(nil)

func guardarElementoDePrueba(t *testing.T, repo *RepositorioBacklog, e domain.ElementoBacklog) domain.ElementoBacklog {
	t.Helper()
	if err := repo.GuardarElemento(context.Background(), &e); err != nil {
		t.Fatalf("guardar elemento %q: %v", e.Titulo, err)
	}
	return e
}

// Traza: HU-07 CA-07.6, RN-07.1: se listan los elementos del proyecto, completos y con sus
// criterios en orden, y no los de otros proyectos.
func TestRepositorioBacklog_ListarElementos_DevuelveLosDelProyectoCompletos(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa", "Beta")
	cinco := 5
	login := guardarElementoDePrueba(t, repo, elementoDePrueba(t, ids[0], "Login", &cinco, "Con clave", "Bloquea tras 3 intentos"))
	guardarElementoDePrueba(t, repo, elementoDePrueba(t, ids[1], "De Beta", nil))
	reportes := guardarElementoDePrueba(t, repo, elementoDePrueba(t, ids[0], "Reportes", nil))

	got, err := repo.ListarElementos(context.Background(), ids[0])

	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	porTitulo := map[string]domain.ElementoBacklog{}
	for _, e := range got {
		porTitulo[e.Titulo] = e
	}
	if len(got) != 2 {
		t.Fatalf("elementos = %+v, se esperaban los 2 de Alfa", got)
	}
	if !reflect.DeepEqual(porTitulo["Login"], login) {
		t.Errorf("Login = %+v, se esperaba %+v", porTitulo["Login"], login)
	}
	if !reflect.DeepEqual(porTitulo["Reportes"], reportes) {
		t.Errorf("Reportes = %+v, se esperaba %+v", porTitulo["Reportes"], reportes)
	}
}

// Traza: HU-07 CA-07.5: un proyecto sin elementos devuelve una lista vacía, sin error.
func TestRepositorioBacklog_ListarElementos_SinElementos(t *testing.T) {
	repo, ids := nuevoBacklogDePrueba(t, "Alfa")

	got, err := repo.ListarElementos(context.Background(), ids[0])

	if err != nil || len(got) != 0 {
		t.Errorf("elementos = %+v, error = %v; se esperaba lista vacía sin error", got, err)
	}
}
