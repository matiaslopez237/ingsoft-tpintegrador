package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: HU-13 CA-13.2 (RN-13.1): crear con un valor fuera de la escala no llega al repositorio.
// No hizo falta cambiar el caso de uso: la regla vive en el dominio.
func TestCrearElemento_StoryPointsFueraDeLaEscala_NoSeGuarda(t *testing.T) {
	repo := nuevoBacklogFalso(1)
	d := datosElemento(1, "Login")
	cuatro := 4
	d.StoryPoints = &cuatro

	_, err := NuevoCrearElementoBacklog(repo).Ejecutar(context.Background(), d)

	if !errors.Is(err, domain.ErrStoryPointsInvalidos) || len(repo.guardados) != 0 {
		t.Errorf("err = %v, guardados = %v", err, repo.guardados)
	}
}

// Traza: HU-13 CA-13.5 (RN-13.3): editar otro campo de un elemento guardado con 4 conserva los 4,
// y enviar 4 se rechaza sin actualizar.
func TestEditarElemento_ValorAnteriorALaEscala(t *testing.T) {
	viejo := elementoDePrueba(t, 7, 1, "Login")
	cuatro := 4
	viejo.StoryPoints = &cuatro
	repo := nuevaEdicionFalsa(viejo)
	uc := NuevoEditarElementoBacklog(repo)

	nuevo := "Ingreso"
	e, err := uc.Ejecutar(context.Background(), DatosEditarElemento{ProyectoID: 1, ElementoID: 7, Cambios: domain.CambiosElemento{Titulo: &nuevo}})
	if err != nil || e.StoryPoints == nil || *e.StoryPoints != 4 {
		t.Errorf("editar el título: %v, %v; se esperaba conservar los 4 story points", e.StoryPoints, err)
	}

	otroCuatro := 4
	_, err = uc.Ejecutar(context.Background(), DatosEditarElemento{ProyectoID: 1, ElementoID: 7,
		Cambios: domain.CambiosElemento{CambiarStoryPoints: true, StoryPoints: &otroCuatro}})
	if !errors.Is(err, domain.ErrStoryPointsInvalidos) || len(repo.actualizados) != 1 {
		t.Errorf("enviar 4: err = %v, actualizaciones = %d; se esperaba el error y ninguna actualización nueva", err, len(repo.actualizados))
	}
}
