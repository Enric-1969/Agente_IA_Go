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
