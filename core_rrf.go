package main

import (
	"fmt"
	"sort"
)

func generarIDChunk(doc DocumentoVectorial) string {
	return fmt.Sprintf("%s_%s", doc.Chunk.Origen, doc.Chunk.Contenido)
}

// FusionarResultadosRRF combina ranking vectorial y BM25 mediante Reciprocal Rank Fusion.
func FusionarResultadosRRF(resVector []ResultadoBusqueda, resBM25 []ResultadoHibrido, topK int) []ResultadoHibrido {
	kRRF := 60.0
	scoresMap := make(map[string]*ResultadoHibrido)

	for rank, v := range resVector {
		id := generarIDChunk(v.Documento)
		if _, ok := scoresMap[id]; !ok {
			scoresMap[id] = &ResultadoHibrido{
				Documento:   v.Documento,
				ScoreVector: float64(v.Similitud),
			}
		}
		scoresMap[id].ScoreRRF += 1.0 / (kRRF + float64(rank+1))
	}

	for rank, b := range resBM25 {
		id := generarIDChunk(b.Documento)
		if _, ok := scoresMap[id]; !ok {
			scoresMap[id] = &ResultadoHibrido{
				Documento: b.Documento,
				ScoreBM25: b.ScoreBM25,
			}
		} else {
			scoresMap[id].ScoreBM25 = b.ScoreBM25
		}
		scoresMap[id].ScoreRRF += 1.0 / (kRRF + float64(rank+1))
	}

	lista := make([]ResultadoHibrido, 0, len(scoresMap))
	for _, item := range scoresMap {
		lista = append(lista, *item)
	}

	sort.Slice(lista, func(i, j int) bool {
		return lista[i].ScoreRRF > lista[j].ScoreRRF
	})

	if topK > len(lista) {
		topK = len(lista)
	}
	return lista[:topK]
}
