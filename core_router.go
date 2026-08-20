package main

import (
	"context"
	"fmt"
	"strings"
)

// EvaluarIntencionHeuristica analiza palabras clave antes de llamar a la API.
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

// ClasificarConsultaClave invoca tu función nativa EnviarConsultaGroq.
func ClasificarConsultaClave(ctx context.Context, gestor *GestorModelos, consulta string) (TipoRuta, error) {
	// 1. Intento de resolución local sin costo de tokens
	if ruta, determinada := EvaluarIntencionHeuristica(consulta); determinada {
		return ruta, nil
	}

	// 2. Si es ambigua, preparamos el prompt y el formato []Mensaje
	prompt := FormatearPromptRouter(consulta)
	mensajes := []Mensaje{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// 3. Consulta a la API usando tu conector real
	respuesta, err := EnviarConsultaGroq(gestor, mensajes)
	if err != nil {
		return RutaGeneral, fmt.Errorf("error al clasificar ruta: %w", err)
	}

	respuestaLimpia := strings.TrimSpace(strings.ToLower(respuesta))
	switch {
	case strings.Contains(respuestaLimpia, "rag"):
		return RutaRAG, nil
	case strings.Contains(respuestaLimpia, "herramienta"):
		return RutaHerramienta, nil
	default:
		return RutaGeneral, nil
	}
}
