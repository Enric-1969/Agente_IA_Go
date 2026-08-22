package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// --- COORDINADOR PRINCIPAL DE INGESTA ---

func ProcesarDominioIngesta(rutaOrigen, rutaDestinoVS string) ([]ChunkTexto, error) {
	rutaManifest := filepath.Join(rutaDestinoVS, "manifest.json")
	manifest, err := CargarManifest(rutaManifest)
	if err != nil {
		return nil, fmt.Errorf("error al cargar manifest: %w", err)
	}
	nuevoManifest := make(map[string]string)
	for k, v := range manifest {
		nuevoManifest[k] = v
	}

	var todosLosChunks []ChunkTexto
	err = filepath.Walk(rutaOrigen, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !esExtensionSoportada(ext) {
			return nil
		}

		relPath, _ := filepath.Rel(rutaOrigen, path)
		hashActual, _ := ObtenerHashArchivo(path)

		// 1. Verificar si es nuevo o si cambió
		if !archivoHaCambiado(relPath, hashActual, manifest) {
			return nil
		}

		// 2. Anunciar la detección al usuario
		fmt.Printf("⚡ Detectado cambio o nuevo archivo: %s\n", relPath)

		// 3. Evaluar tamaño en MB antes de cargar en RAM
		pesoMB := float64(info.Size()) / (1024 * 1024)
		if pesoMB > 50.0 {
			limiteAjustado := int(math.Ceil(pesoMB/10.0)) * 10
			fmt.Printf("⚠️  [ALERTA MEMORIA] El archivo pesa %.1f MB (supera el estándar de 50 MB).\n", pesoMB)
			fmt.Printf("¿Deseas permitir su análisis ampliando el límite a %d MB solo para este archivo? (s/n): ", limiteAjustado)

			var respuesta string
			fmt.Scanln(&respuesta)

			if strings.ToLower(strings.TrimSpace(respuesta)) != "s" {
				fmt.Println("⏭️ Archivo ignorado por decisión del usuario. Pasando al siguiente...")
				return nil // Salta al siguiente fichero sin actualizar manifest ni extraer texto
			}
		}

		// 4. Extraer el texto solo tras recibir confirmación
		textoExtraido, err := ExtraerTextoDeArchivo(path)
		if err != nil || strings.TrimSpace(textoExtraido) == "" {
			return nil
		}

		chunks := ConvertirTextoAChunks(relPath, textoExtraido)
		todosLosChunks = append(todosLosChunks, chunks...)

		nuevoManifest[relPath] = hashActual
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(todosLosChunks) == 0 {
		fmt.Println("✅ Base vectorial 100% actualizada. No hay archivos nuevos o modificados.")
		return nil, nil
	}

	_ = GuardarManifest(rutaManifest, nuevoManifest)
	fmt.Printf("✂️ Generados %d chunks limpios de texto.\n", len(todosLosChunks))
	return todosLosChunks, nil
}
