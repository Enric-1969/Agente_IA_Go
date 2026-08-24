package main

import (
	"context"
	"fmt"
	"strings"
)

// ClasificarConsultaClave coordina la clasificación de intenciones (heurística local + fallback LLM).
func ClasificarConsultaClave(ctx context.Context, gestor *GestorModelos, consulta string) (TipoRuta, error) {
	// 1. Intento de resolución local sin costo de tokens
	if ruta, determinada := EvaluarIntencionHeuristica(consulta); determinada {
		return ruta, nil
	}

	// 2. Si es ambigua, se construye el prompt y la consulta al LLM
	prompt := FormatearPromptRouter(consulta)
	mensajes := []Mensaje{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// 3. Consulta a la API mediante el conector nativo de Groq
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
	case strings.Contains(respuestaLimpia, "desarrollo") || strings.Contains(respuestaLimpia, "codigo"):
		return RutaDesarrollo, nil
	default:
		return RutaGeneral, nil
	}
}
