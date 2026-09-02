package main

import "fmt"

// --- FUNCIONES AUXILIARES Y FORMATEADORES DE PROMPTS ---

// FormatearPromptRAG une los chunks extraídos con la pregunta del usuario.
func FormatearPromptRAG(contexto string, consulta string) string {
	return fmt.Sprintf(PromptRAG, contexto, consulta)
}

// FormatearPromptRouter inyecta la consulta para clasificar la intención.
func FormatearPromptRouter(consulta string) string {
	return fmt.Sprintf(PromptClasificador, consulta)
}

// SeleccionarPromptSistema devuelve el prompt de sistema adecuado según la intención clasificada.
func SeleccionarPromptSistema(ruta TipoRuta) string {
	if ruta == RutaDesarrollo {
		return PromptSistemaDesarrollador
	}
	return PromptSistemaBase
}

// FormatearPromptDesarrolloGo inyecta la consulta de programación.
func FormatearPromptDesarrolloGo(contexto string, consulta string) string {
	return fmt.Sprintf(PromptRAGDesarrolloGo, contexto, consulta)
}
