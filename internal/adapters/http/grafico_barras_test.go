package http

import (
	"math"
	"testing"
)

func valor(v float64) *float64 { return &v }

func cerca(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// alturaUtil es el alto que pueden ocupar las barras: el alto del gráfico sin los márgenes.
const alturaUtil = altoGrafico - margenArriba - margenAbajo

func especDosSeries(grupos ...grupoGrafico) especGrafico {
	return especGrafico{
		ID: "puntos", Titulo: "Puntos por sprint",
		Series: []serieGrafico{{Nombre: "Planificados", Clase: "serie-1"}, {Nombre: "Completados", Clase: "serie-2"}},
		Grupos: grupos,
	}
}

// Traza: HU-25 CA-25.9 (ADR-0005): el alto de cada barra es proporcional a su valor, con el mayor
// valor ocupando todo el alto útil, y las barras crecen desde el eje.
func TestArmarGraficoBarras_AltoProporcional(t *testing.T) {
	g := armarGraficoBarras(especDosSeries(grupoGrafico{Etiqueta: "S1", Valores: []*float64{valor(8), valor(6)}}))

	b := g.Grupos[0].Barras
	if !cerca(b[0].Alto, alturaUtil) || !cerca(b[1].Alto, alturaUtil*0.75) {
		t.Errorf("altos = %v y %v, se esperaba %v y %v", b[0].Alto, b[1].Alto, float64(alturaUtil), alturaUtil*0.75)
	}
	for _, barra := range b {
		if !cerca(barra.Y+barra.Alto, g.BaseY) {
			t.Errorf("la barra %+v no apoya en el eje (y=%v)", barra, g.BaseY)
		}
	}
	if b[0].Clase != "serie-1" || b[1].Clase != "serie-2" {
		t.Errorf("clases = %q y %q, se esperaba serie-1 y serie-2", b[0].Clase, b[1].Clase)
	}
}

// Traza: HU-25, sección 6: con todos los valores en 0 no se divide por cero; las barras quedan en 0.
func TestArmarGraficoBarras_TodoEnCero_SinDividirPorCero(t *testing.T) {
	g := armarGraficoBarras(especDosSeries(grupoGrafico{Etiqueta: "S1", Valores: []*float64{valor(0), valor(0)}}))

	for _, b := range g.Grupos[0].Barras {
		if b.Alto != 0 || math.IsNaN(b.Y) || math.IsInf(b.Y, 0) {
			t.Errorf("barra = %+v, se esperaba alto 0 y una posición válida", b)
		}
	}
}

// Traza: HU-25 CA-25.9 (RN-24.2): un valor sin datos no tiene barra y lo dice con texto.
func TestArmarGraficoBarras_ValorSinDatos(t *testing.T) {
	spec := especGrafico{
		ID: "porcentaje", Titulo: "% completado", Sufijo: "%", MinimoEscala: 100,
		Series: []serieGrafico{{Nombre: "% completado", Clase: "serie-1"}},
		Grupos: []grupoGrafico{{Etiqueta: "S1", Valores: []*float64{nil}}},
	}

	b := armarGraficoBarras(spec).Grupos[0].Barras[0]

	if !b.SinDatos || b.Alto != 0 || b.Texto != "sin datos" {
		t.Errorf("barra = %+v, se esperaba sin datos, alto 0 y el texto \"sin datos\"", b)
	}
}

// Traza: HU-25 CA-25.9: el porcentaje se dibuja sobre 100 aunque ningún sprint llegue a 100.
func TestArmarGraficoBarras_MinimoDeEscala(t *testing.T) {
	spec := especGrafico{
		ID: "porcentaje", Titulo: "% completado", Sufijo: "%", MinimoEscala: 100,
		Series: []serieGrafico{{Nombre: "% completado", Clase: "serie-1"}},
		Grupos: []grupoGrafico{{Etiqueta: "S1", Valores: []*float64{valor(50)}}},
	}

	b := armarGraficoBarras(spec).Grupos[0].Barras[0]

	if !cerca(b.Alto, alturaUtil/2) {
		t.Errorf("alto = %v, se esperaba la mitad del alto útil (%v)", b.Alto, float64(alturaUtil)/2)
	}
}

// Traza: HU-25 CA-25.9: el valor va escrito sobre cada barra, con coma decimal y sin ",0" de más.
func TestArmarGraficoBarras_TextoDeLosValores(t *testing.T) {
	spec := especGrafico{
		ID: "porcentaje", Titulo: "% completado", Sufijo: "%",
		Series: []serieGrafico{{Nombre: "% completado", Clase: "serie-1"}},
		Grupos: []grupoGrafico{
			{Etiqueta: "S1", Valores: []*float64{valor(66.7)}},
			{Etiqueta: "S2", Valores: []*float64{valor(100)}},
		},
	}

	g := armarGraficoBarras(spec)

	if got := g.Grupos[0].Barras[0].Texto; got != "66,7%" {
		t.Errorf("texto = %q, se esperaba \"66,7%%\"", got)
	}
	if got := g.Grupos[1].Barras[0].Texto; got != "100%" {
		t.Errorf("texto = %q, se esperaba \"100%%\"", got)
	}
}

// Traza: HU-25 RN-25.5 (P-2): los grupos van de izquierda a derecha en el orden recibido, sin
// superponerse, y el ancho del gráfico alcanza para todos, por muchos que sean.
func TestArmarGraficoBarras_GruposEnOrdenYSinSuperponerse(t *testing.T) {
	var grupos []grupoGrafico
	for i := range 12 {
		grupos = append(grupos, grupoGrafico{Etiqueta: "S" + string(rune('A'+i)), Valores: []*float64{valor(float64(i)), valor(1)}})
	}

	g := armarGraficoBarras(especDosSeries(grupos...))

	finAnterior := 0.0
	for i, grupo := range g.Grupos {
		if grupo.Etiqueta != grupos[i].Etiqueta {
			t.Fatalf("grupo %d = %q, se esperaba %q", i, grupo.Etiqueta, grupos[i].Etiqueta)
		}
		for _, b := range grupo.Barras {
			if b.X < finAnterior {
				t.Errorf("la barra %+v del grupo %q se superpone con la anterior (que termina en %v)", b, grupo.Etiqueta, finAnterior)
			}
			finAnterior = b.X + b.Ancho
		}
	}
	if finAnterior > float64(g.Ancho) {
		t.Errorf("la última barra termina en %v, fuera del ancho del gráfico (%d)", finAnterior, g.Ancho)
	}
}

// Traza: HU-25 CA-25.4, CA-25.9: la velocidad se dibuja como una línea a su altura.
func TestArmarGraficoBarras_LineaDeReferencia(t *testing.T) {
	spec := especDosSeries(grupoGrafico{Etiqueta: "S1", Valores: []*float64{valor(8), valor(4)}})
	spec.Referencia = valor(4)
	spec.EtiquetaReferencia = "Velocidad"

	g := armarGraficoBarras(spec)

	if g.Referencia == nil {
		t.Fatal("se esperaba la línea de referencia")
	}
	if want := g.BaseY - alturaUtil/2; !cerca(g.Referencia.Y, want) {
		t.Errorf("y de la línea = %v, se esperaba %v (la mitad del alto útil)", g.Referencia.Y, want)
	}
	if g.Referencia.Etiqueta != "Velocidad: 4" {
		t.Errorf("etiqueta = %q, se esperaba \"Velocidad: 4\"", g.Referencia.Etiqueta)
	}
}
