package main

import (
	"fmt"
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

		if !archivoHaCambiado(relPath, hashActual, manifest) {
			return nil
		}

		textoExtraido, err := ExtraerTextoDeArchivo(path)
		if err != nil || strings.TrimSpace(textoExtraido) == "" {
			return nil // Si falla la lectura o está vacío, saltamos al siguiente archivo
		}

		fmt.Printf("⚡ Detectado cambio o nuevo archivo: %s\n", relPath)
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
