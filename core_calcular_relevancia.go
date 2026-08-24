package main

import (
	"sort"
	"strings"
)

// CalcularDensidadLexica evalúa la proporción de términos de la consulta presentes en el texto del fragmento.
func CalcularDensidadLexica(query string, contenido string) float64 {
	terminosQuery := Tokenizar(query)
	if len(terminosQuery) == 0 {
		return 0.0
	}

	textoNormalizado := strings.ToLower(contenido)
	coincidencias := 0

	for _, term := range terminosQuery {
		if strings.Contains(textoNormalizado, term) {
			coincidencias++
		}
	}

	return float64(coincidencias) / float64(len(terminosQuery))
}

// PonderarScoreRelevancia combina el score RRF previo con la densidad léxica según pesos configurables.
func PonderarScoreRelevancia(scoreRRF float64, densidadLexica float64, pesoRRF float64, pesoDensidad float64) float64 {
	return (scoreRRF * pesoRRF) + (densidadLexica * pesoDensidad)
}

// OrdenarYFiltrarResultados ordena descendentemente los resultados híbridos según su ScoreRRF recalculado.
func OrdenarYFiltrarResultados(resultados []ResultadoHibrido, limite int) []ResultadoHibrido {
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].ScoreRRF > resultados[j].ScoreRRF
	})

	if limite > 0 && len(resultados) > limite {
		return resultados[:limite]
	}
	return resultados
}
