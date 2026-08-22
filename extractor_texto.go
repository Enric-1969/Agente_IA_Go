// (lectura de .txt y .md)
package main

import (
	"fmt"
	"os"
)

// ExtraerTextoPlanoOMarkdown lee el contenido de archivos .txt y .md.
func ExtraerTextoPlanoOMarkdown(ruta string) (string, error) {
	contenido, err := os.ReadFile(ruta)
	if err != nil {
		return "", fmt.Errorf("error al leer archivo de texto: %w", err)
	}

	return string(contenido), nil
}
