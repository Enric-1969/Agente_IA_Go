package main

import (
	"strings"
)

// EvaluarIntencionHeuristica analiza patrones de palabras clave locales antes de invocar a la API.
func EvaluarIntencionHeuristica(consulta string) (TipoRuta, bool) {
	q := strings.ToLower(strings.TrimSpace(consulta))

	// Detección directa de tareas de desarrollo y Clean Code
	if strings.Contains(q, "refactoriza") || strings.Contains(q, "clean code") ||
		strings.Contains(q, "crea fichero") || strings.Contains(q, "codigo go") {
		return RutaDesarrollo, true
	}

	if strings.HasPrefix(q, "que es") || strings.HasPrefix(q, "como") || strings.Contains(q, "documentacion") {
		return RutaRAG, true
	}

	if strings.Contains(q, "listar archivos") || strings.Contains(q, "calcular") || strings.Contains(q, "buscar en z:") {
		return RutaHerramienta, true
	}

	return RutaGeneral, false
}
