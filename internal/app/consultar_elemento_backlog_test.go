package app

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

// Traza: HU-06 CA-06.10 (la pantalla muestra el elemento del proyecto)
func TestConsultarElemento_DelProyecto(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 1, "Login"))

	e, err := NuevoConsultarElementoBacklog(repo).Ejecutar(context.Background(), 1, 7)

	if err != nil || e.ID != 7 || e.Titulo != "Login" {
		t.Errorf("se obtuvo %+v, %v", e, err)
	}
}

// Traza: HU-06 CA-06.7 (RN-06.6)
func TestConsultarElemento_DeOtroProyectoOInexistente(t *testing.T) {
	repo := nuevaEdicionFalsa(elementoDePrueba(t, 7, 2, "Pagos"))
	uc := NuevoConsultarElementoBacklog(repo)

	for _, ids := range [][2]int64{{1, 7}, {2, 99}} {
		if _, err := uc.Ejecutar(context.Background(), ids[0], ids[1]); !errors.Is(err, domain.ErrElementoNoEncontrado) {
			t.Errorf("proyecto %d, elemento %d: error = %v", ids[0], ids[1], err)
		}
	}
}
