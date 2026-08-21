// (el despachador central ExtraerTextoDeArchivo)
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ExtraerTextoDeArchivo es el orquestador principal de ingesta.
func ExtraerTextoDeArchivo(ruta string) (string, error) {
	ext := strings.ToLower(filepath.Ext(ruta))

	switch ext {
	case ".txt", ".md", ".go", ".py", ".html", ".css", ".js", ".json", ".sql", ".sh", ".bat":
		return ExtraerTextoPlanoOMarkdown(ruta)
	case ".xlsx", ".xlsm":
		return "", fmt.Errorf("extractor de Excel (.xlsx) pendiente de implementacion")

	case ".docx":
		return "", fmt.Errorf("extractor de Word (.docx) pendiente de implementacion")

	case ".pdf":
		return "", fmt.Errorf("extractor de PDF (.pdf) pendiente de implementacion")

	default:
		return "", fmt.Errorf("formato no soportado: %s", ext)
	}
}
