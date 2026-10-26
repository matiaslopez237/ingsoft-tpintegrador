package domain

import (
	"errors"
	"slices"
	"testing"
)

// elementosDePrueba son elementos de un proyecto en orden de creación (ID creciente).
func elementosDePrueba() []ElementoBacklog {
	return []ElementoBacklog{
		{ID: 1, Titulo: "Reportes", Prioridad: PrioridadBaja, Estado: EstadoHecho},
		{ID: 2, Titulo: "Login", Prioridad: PrioridadAlta, Estado: EstadoPendiente},
		{ID: 3, Titulo: "Perfil", Prioridad: PrioridadMedia, Estado: EstadoEnProgreso},
		{ID: 4, Titulo: "Registro", Prioridad: PrioridadAlta, Estado: EstadoPendiente},
	}
}

func titulos(elementos []ElementoBacklog) []string {
	var t []string
	for _, e := range elementos {
		t = append(t, e.Titulo)
	}
	return t
}

// Traza: HU-07 CA-07.1 a CA-07.5, RN-07.2 a RN-07.5
func TestOrdenarBacklog(t *testing.T) {
	tests := []struct {
		nombre string
		orden  string
		estado string
		want   []string
	}{
		{"por defecto, de Alta a Baja y a igual prioridad por creación", "", "", []string{"Login", "Registro", "Perfil", "Reportes"}},
		{"descendente explícito", "desc", "", []string{"Login", "Registro", "Perfil", "Reportes"}},
		{"ascendente, de Baja a Alta y a igual prioridad por creación", "asc", "", []string{"Reportes", "Perfil", "Login", "Registro"}},
		{"filtro Pendiente", "", "Pendiente", []string{"Login", "Registro"}},
		{"filtro En progreso", "", "En progreso", []string{"Perfil"}},
		{"filtro Hecho", "", "Hecho", []string{"Reportes"}},
		{"filtro y orden ascendente juntos", "asc", "Pendiente", []string{"Login", "Registro"}},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			got, err := OrdenarBacklog(elementosDePrueba(), tt.orden, tt.estado)

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !slices.Equal(titulos(got), tt.want) {
				t.Errorf("orden = %v, se esperaba %v", titulos(got), tt.want)
			}
		})
	}
}

// Traza: HU-07 RN-07.3: el desempate es por orden de creación aunque la lista llegue desordenada.
func TestOrdenarBacklog_DesempataPorCreacionAunqueLleguenDesordenados(t *testing.T) {
	desordenados := []ElementoBacklog{
		{ID: 9, Titulo: "Tercero", Prioridad: PrioridadAlta},
		{ID: 2, Titulo: "Primero", Prioridad: PrioridadAlta},
		{ID: 5, Titulo: "Segundo", Prioridad: PrioridadAlta},
	}

	got, err := OrdenarBacklog(desordenados, "", "")

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"Primero", "Segundo", "Tercero"}; !slices.Equal(titulos(got), want) {
		t.Errorf("orden = %v, se esperaba %v", titulos(got), want)
	}
}

// Traza: HU-07 CA-07.5, RN-07.5: sin elementos, o sin coincidencias con el filtro, lista vacía y sin error.
func TestOrdenarBacklog_ListaVacia(t *testing.T) {
	if got, err := OrdenarBacklog(nil, "", ""); err != nil || len(got) != 0 {
		t.Errorf("sin elementos: %v, %v; se esperaba lista vacía sin error", got, err)
	}
	soloPendientes := []ElementoBacklog{{ID: 1, Titulo: "Login", Prioridad: PrioridadAlta, Estado: EstadoPendiente}}
	if got, err := OrdenarBacklog(soloPendientes, "", "Hecho"); err != nil || len(got) != 0 {
		t.Errorf("filtro sin coincidencias: %v, %v; se esperaba lista vacía sin error", got, err)
	}
}

// Traza: HU-07 CA-07.8, sección 7
func TestOrdenarBacklog_Invalido_RetornaError(t *testing.T) {
	tests := []struct {
		nombre, orden, estado string
		wantErr               error
	}{
		{"orden en mayúsculas", "ASC", "", ErrOrdenInvalido},
		{"orden desconocido", "prioridad", "", ErrOrdenInvalido},
		{"estado en minúsculas", "", "pendiente", ErrEstadoElementoInvalido},
		{"estado desconocido", "", "Cerrado", ErrEstadoElementoInvalido},
		{"orden y estado inválidos: se informa el orden", "ASC", "Cerrado", ErrOrdenInvalido},
	}
	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			_, err := OrdenarBacklog(elementosDePrueba(), tt.orden, tt.estado)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, se esperaba %v", err, tt.wantErr)
			}
		})
	}
}

// Traza: HU-07 sección 3 (listar no modifica nada): la lista recibida no cambia de orden.
func TestOrdenarBacklog_NoModificaLaListaRecibida(t *testing.T) {
	recibidos := elementosDePrueba()

	if _, err := OrdenarBacklog(recibidos, "asc", ""); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if want := titulos(elementosDePrueba()); !slices.Equal(titulos(recibidos), want) {
		t.Errorf("la lista recibida quedó %v, se esperaba %v", titulos(recibidos), want)
	}
}
