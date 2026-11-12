package main

import (
	"strings"
)

// EvaluarIntencionHeuristica analiza patrones de palabras clave locales (multilenguaje).
func EvaluarIntencionHeuristica(consulta string) (TipoRuta, bool) {
	q := strings.ToLower(strings.TrimSpace(consulta))

	// 1. Primero detectamos preguntas directas sobre documentación (ej. "Como funciona...")
	if strings.HasPrefix(q, "que es") || strings.HasPrefix(q, "como") || strings.Contains(q, "documentacion") {
		return RutaRAG, true
	}

	// 2. Detección genérica de tareas de desarrollo (añadimos espacio a "funcion" para evitar "funciona")
	if strings.Contains(q, "refactoriza") || strings.Contains(q, "clean code") ||
		strings.Contains(q, "crea fichero") || strings.Contains(q, "codigo") ||
		strings.Contains(q, "funcion ") || strings.HasSuffix(q, "funcion") ||
		strings.Contains(q, "script") {
		return RutaDesarrollo, true
	}

	// 3. Herramientas locales del sistema
	if strings.Contains(q, "listar archivos") || strings.Contains(q, "calcular") || strings.Contains(q, "buscar en z:") {
		return RutaHerramienta, true
	}

	return RutaGeneral, false
}
