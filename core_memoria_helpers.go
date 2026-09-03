package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ExportarHistorialMarkdown convierte los mensajes del GestorMemoria en un archivo Markdown.
func ExportarHistorialMarkdown(memoria *GestorMemoria, rutaDestino string) error {
	historial := memoria.ObtenerHistorial()
	if len(historial) == 0 {
		return fmt.Errorf("no hay mensajes en memoria para exportar")
	}

	var sb strings.Builder
	sb.WriteString("# 📝 Registro de Sesión del Agente IA\n")
	sb.WriteString(fmt.Sprintf("**Fecha:** %s\n\n---\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for _, msg := range historial {
		switch msg.Role {
		case "user":
			sb.WriteString(fmt.Sprintf("### 👤 Usuario\n%s\n\n", msg.Content))
		case "assistant":
			sb.WriteString(fmt.Sprintf("### 🤖 Agente\n%s\n\n---\n\n", msg.Content))
		}
	}

	return os.WriteFile(rutaDestino, []byte(sb.String()), 0644)
}
