package main

// EjecutarReordenamiento orquesta la re-evaluación y filtrado Top-N de candidatos devueltos por RRF.
func EjecutarReordenamiento(query string, candidatos []ResultadoHibrido, limiteTopN int, pesoRRF float64, pesoDensidad float64) []ResultadoHibrido {
	if len(candidatos) == 0 {
		return candidatos
	}

	resultadosReordenados := make([]ResultadoHibrido, len(candidatos))

	for i, candidato := range candidatos {
		densidad := CalcularDensidadLexica(query, candidato.Documento.Chunk.Contenido)
		scoreFinal := PonderarScoreRelevancia(candidato.ScoreRRF, densidad, pesoRRF, pesoDensidad)

		res := candidato
		res.ScoreRRF = scoreFinal
		resultadosReordenados[i] = res
	}

	return OrdenarYFiltrarResultados(resultadosReordenados, limiteTopN)
}
