// (lectura de .txt y .md)
package main

import (
	"fmt"
	"os"
)

// Limite maximo de tamaño para evitar desbordamiento de RAM (50 MB)
const MaxTamanoTextoBytes = 50 * 1024 * 1024

// ExtraerTextoPlanoOMarkdown lee el contenido de archivos .txt y .md.
func ExtraerTextoPlanoOMarkdown(ruta string) (string, error) {
	info, err := os.Stat(ruta)
	if err != nil {
		return "", fmt.Errorf("error al obtener info del archivo: %w", err)
	}

	if info.Size() > MaxTamanoTextoBytes {
		return "", fmt.Errorf("el archivo supera el limite maximo permitido de %d MB", MaxTamanoTextoBytes/(1024*1024))
	}

	contenido, err := os.ReadFile(ruta)
	if err != nil {
		return "", fmt.Errorf("error al leer archivo de texto: %w", err)
	}

	return string(contenido), nil
}
