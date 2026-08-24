package main

import (
	"strings"
)

// EvaluarIntencionHeuristica analiza patrones de palabras clave locales antes de invocar a la API.
func EvaluarIntencionHeuristica(consulta string) (TipoRuta, bool) {
	q := strings.ToLower(strings.TrimSpace(consulta))

	if strings.HasPrefix(q, "que es") || strings.HasPrefix(q, "como") || strings.Contains(q, "documentacion") {
		return RutaRAG, true
	}

	if strings.Contains(q, "listar archivos") || strings.Contains(q, "calcular") || strings.Contains(q, "buscar en z:") {
		return RutaHerramienta, true
	}

	return RutaGeneral, false
}
