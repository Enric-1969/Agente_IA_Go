package main

import (
	"math"
	"sort"
)

// NuevoIndiceBM25 construye el índice léxico calculando frecuencias e IDF sobre el VectorStore.
func NuevoIndiceBM25(docs []DocumentoVectorial) *IndiceBM25 {
	k1 := 1.2
	b := 0.75
	totalDocs := len(docs)
	if totalDocs == 0 {
		return &IndiceBM25{K1: k1, B: b}
	}

	docLen := make(map[int]float64)
	termFreqs := make([]map[string]int, totalDocs)
	docFreq := make(map[string]int)
	totalLen := 0.0

	for i, doc := range docs {
		tokens := Tokenizar(doc.Chunk.Contenido)
		tf := make(map[string]int)
		for _, t := range tokens {
			tf[t]++
		}
		termFreqs[i] = tf
		docLen[i] = float64(len(tokens))
		totalLen += float64(len(tokens))

		for t := range tf {
			docFreq[t]++
		}
	}

	avgDocLen := totalLen / float64(totalDocs)
	idf := make(map[string]float64)

	for term, df := range docFreq {
		val := math.Log(1.0 + (float64(totalDocs)-float64(df)+0.5)/(float64(df)+0.5))
		idf[term] = math.Max(0.0, val)
	}

	return &IndiceBM25{
		IDF:       idf,
		DocLen:    docLen,
		AvgDocLen: avgDocLen,
		TotalDocs: totalDocs,
		TermFreqs: termFreqs,
		K1:        k1,
		B:         b,
	}
}

// BuscarBM25 evalúa y ordena los fragmentos por relevancia de palabras clave.
func (bm *IndiceBM25) BuscarBM25(query string, docs []DocumentoVectorial, topK int) []ResultadoHibrido {
	tokensQuery := Tokenizar(query)
	if len(tokensQuery) == 0 || bm.TotalDocs == 0 {
		return nil
	}

	resultados := make([]ResultadoHibrido, bm.TotalDocs)

	for i, doc := range docs {
		score := 0.0
		tfMap := bm.TermFreqs[i]
		dLen := bm.DocLen[i]

		for _, qTerm := range tokensQuery {
			tf, ok := tfMap[qTerm]
			if !ok {
				continue
			}
			idfVal := bm.IDF[qTerm]
			num := float64(tf) * (bm.K1 + 1.0)
			denom := float64(tf) + bm.K1*(1.0-bm.B+bm.B*(dLen/bm.AvgDocLen))
			score += idfVal * (num / denom)
		}

		resultados[i] = ResultadoHibrido{
			Documento: doc,
			ScoreBM25: score,
		}
	}

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].ScoreBM25 > resultados[j].ScoreBM25
	})

	if topK > len(resultados) {
		topK = len(resultados)
	}
	return resultados[:topK]
}
