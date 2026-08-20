package main

import "fmt"

// --- PROMPTS BASE DEL AGENTE ---

// PromptSistemaBase establece el comportamiento estricto anti-alucinaciones.
const PromptSistemaBase = `Eres un Agente de IA técnico especializado.
Tu única fuente de verdad son los documentos provistos desde el directorio local.

REGLAS OBLIGATORIAS DE SEGURIDAD:
1. Responde utilizando EXCLUSIVAMENTE la información provista en el contexto recuperado.
2. Si la respuesta no se encuentra en el contexto, di explícitamente: "No dispongo de información suficiente en la documentación local (Z:/Files_IA) para responder a tu consulta."
3. NO asumas, NO deduzcas fuera del texto y NO inventes datos que no estén presentes.`

// PromptRAG inyecta la información extraída de Z:/Files_IA.
const PromptRAG = `Utiliza el siguiente contexto local para responder a la pregunta del usuario de forma precisa y concisa.

--- CONTEXTO RECUPERADO (Z:/Files_IA) ---
%s
-----------------------------------------

Pregunta del usuario: %s`

// PromptClasificador permite al agente elegir el tipo de tarea antes de procesar.
const PromptClasificador = `Analiza la consulta e identifica la intención del usuario.
Responde ÚNICAMENTE con una de estas opciones en minúsculas:
- rag (para consultas sobre documentación de Z:/Files_IA)
- herramienta (para cálculos, búsquedas directas o tareas que requieran herramientas)
- general (para interacciones directas sin necesidad de documentos)

Consulta: %s`

// --- FUNCIONES DE CONSTRUCCIÓN ---

// FormatearPromptRAG une los chunks extraídos con la pregunta del usuario.
func FormatearPromptRAG(contexto string, consulta string) string {
	return fmt.Sprintf(PromptRAG, contexto, consulta)
}

// FormatearPromptRouter inyecta la consulta para clasificar la intención.
func FormatearPromptRouter(consulta string) string {
	return fmt.Sprintf(PromptClasificador, consulta)
}
