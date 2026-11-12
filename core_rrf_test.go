package main

import (
	"testing"
)

func TestFusionarResultadosRRF(t *testing.T) {
	// 1. Preparar (Arrange): Creamos 3 documentos simulados
	docA := DocumentoVectorial{Chunk: ChunkTexto{Origen: "manual.go", Contenido: "func main() {}"}}
	docB := DocumentoVectorial{Chunk: ChunkTexto{Origen: "readme.md", Contenido: "Instalacion en Windows"}}
	docC := DocumentoVectorial{Chunk: ChunkTexto{Origen: "notas.txt", Contenido: "Texto secundario"}}

	// docB sale 2º en Vector y 1º en BM25 (al aparecer en ambos, debe ganar el ranking RRF)
	resVector := []ResultadoBusqueda{
		{Documento: docA, Similitud: 0.95},
		{Documento: docB, Similitud: 0.80},
	}

	resBM25 := []ResultadoHibrido{
		{Documento: docB, ScoreBM25: 12.5},
		{Documento: docC, ScoreBM25: 8.0},
	}

	// 2. Actuar (Act): Pedimos el Top 2 fusionado
	top2 := FusionarResultadosRRF(resVector, resBM25, 2)

	// 3. Comprobar (Assert)
	if len(top2) != 2 {
		t.Fatalf("Se esperaban 2 resultados por topK=2, se obtuvieron %d", len(top2))
	}

	// El ganador debe ser docB ("readme.md") por sumar puntos en ambos rankings
	if top2[0].Documento.Chunk.Origen != "readme.md" {
		t.Errorf("Se esperaba 'readme.md' en primer lugar, pero salió '%s'", top2[0].Documento.Chunk.Origen)
	}

	// Debe conservar tanto su puntuación vectorial como su puntuación BM25
	if top2[0].ScoreVector == 0 || top2[0].ScoreBM25 == 0 {
		t.Errorf("Se esperaba que docB conservara ambos scores (Vector: %f, BM25: %f)", top2[0].ScoreVector, top2[0].ScoreBM25)
	}

	// Comprobación extra: si topK es mayor que el total de elementos (3), no debe dar panic
	todos := FusionarResultadosRRF(resVector, resBM25, 10)
	if len(todos) != 3 {
		t.Errorf("Se esperaban 3 elementos únicos totales, se obtuvieron %d", len(todos))
	}
}
