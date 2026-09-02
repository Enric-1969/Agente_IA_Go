package main

import (
	"strings"
)

// EvaluarIntencionHeuristica analiza patrones de palabras clave locales (multilenguaje).
func EvaluarIntencionHeuristica(consulta string) (TipoRuta, bool) {
	q := strings.ToLower(strings.TrimSpace(consulta))

	// Detección genérica de tareas de desarrollo
	if strings.Contains(q, "refactoriza") || strings.Contains(q, "clean code") ||
		strings.Contains(q, "crea fichero") || strings.Contains(q, "codigo") ||
		strings.Contains(q, "funcion") || strings.Contains(q, "script") {
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
