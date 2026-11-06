package ejemplo

import (
	"context"
	"errors"
	"testing"

	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/domain"
)

type listadorProyectosFalso struct {
	proyectos []domain.Proyecto
	err       error
	llamado   bool
}

func (f *listadorProyectosFalso) Ejecutar(context.Context) ([]domain.Proyecto, error) {
	f.llamado = true
	return f.proyectos, f.err
}

type creadorProyectoFalso struct{ llamado bool }

func (f *creadorProyectoFalso) Ejecutar(context.Context, app.DatosNuevoProyecto) (domain.Proyecto, error) {
	f.llamado = true
	return domain.Proyecto{ID: 1}, nil
}

// Traza: HU-31 CA-31.6, CA-31.7, CA-31.8 (RN-31.1, RN-31.3)
func TestCargarSiCorresponde_SinCargar(t *testing.T) {
	errLectura := errors.New("base caída")
	tests := []struct {
		name        string
		valor       string
		proyectos   []domain.Proyecto
		errListar   error
		want        Resultado
		wantErr     error
		wantListado bool
	}{
		{name: "CA-31.7: variable ausente o vacía", valor: "", want: OpcionApagada},
		{name: "CA-31.7: false", valor: "false", want: OpcionApagada},
		{name: "CA-31.7: 0", valor: "0", want: OpcionApagada},
		{name: "CA-31.8: valor no booleano", valor: "si", want: ValorInvalido},
		{name: "CA-31.6: la base ya tiene un proyecto", valor: "true", proyectos: []domain.Proyecto{{ID: 7, Nombre: "Propio"}}, want: BaseConProyectos, wantListado: true},
		{name: "sección 7: no se pudo leer la lista de proyectos", valor: "true", errListar: errLectura, wantErr: errLectura, wantListado: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listar := &listadorProyectosFalso{proyectos: tt.proyectos, err: tt.errListar}
			crear := &creadorProyectoFalso{}
			cu := CasosDeUso{Reloj: &Reloj{}, ListarProyectos: listar, CrearProyecto: crear}

			got, err := CargarSiCorresponde(context.Background(), cu, tt.valor)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("resultado = %v, se esperaba %v", got, tt.want)
			}
			if listar.llamado != tt.wantListado {
				t.Errorf("se listaron los proyectos: %v, se esperaba %v", listar.llamado, tt.wantListado)
			}
			if crear.llamado {
				t.Error("no se tenía que crear ningún proyecto")
			}
		})
	}
}
