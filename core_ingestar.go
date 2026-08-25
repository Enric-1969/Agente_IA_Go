package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// --- COORDINADOR PRINCIPAL DE INGESTA ---

// ProcesarDominioIngesta recorre el directorio, detecta cambios, genera embeddings y guarda la base vectorial.
func ProcesarDominioIngesta(rutaOrigen, rutaDestinoVS, modeloEmbedding string, vs *VectorStore) error {
	rutaManifest := filepath.Join(rutaDestinoVS, "manifest.json")
	manifest, err := CargarManifest(rutaManifest)
	if err != nil {
		manifest = make(map[string]string)
	}

	nuevoManifest := make(map[string]string)
	for k, v := range manifest {
		nuevoManifest[k] = v
	}

	chunksProcesados := 0

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
				return nil
			}
		}

		// 4. Extraer el texto tras recibir confirmación
		textoExtraido, err := ExtraerTextoDeArchivo(path)
		if err != nil || strings.TrimSpace(textoExtraido) == "" {
			return nil
		}

		// 5. Fragmentar en chunks de texto
		chunks := ConvertirTextoAChunks(relPath, textoExtraido)
		if len(chunks) == 0 {
			return nil
		}

		fmt.Printf("📦 Generando embeddings para [%d chunks]: %s\n", len(chunks), relPath)

		// 6. Generar embeddings con Ollama e insertar en el VectorStore
		huboError := false
		for _, chunk := range chunks {
			vec, err := GenerarEmbedding(chunk.Contenido, modeloEmbedding)
			if err != nil {
				fmt.Printf("  ⚠️ Error al generar embedding para un chunk de %s: %v\n", relPath, err)
				huboError = true
				continue
			}
			vs.AgregarVector(chunk, vec)
			chunksProcesados++
		}

		// Solo marcar el archivo como procesado si no hubo ningun error de embedding
		if !huboError {
			nuevoManifest[relPath] = hashActual
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error durante la lectura del directorio: %w", err)
	}

	if chunksProcesados == 0 {
		fmt.Println("✅ Base vectorial 100% actualizada. No hay archivos nuevos o modificados.")
		return nil
	}

	// 7. Guardar persistencia en disco (.gob y manifest.json)
	if err := vs.GuardarPersistencia(); err != nil {
		return fmt.Errorf("error al guardar el vectorstore: %w", err)
	}

	if err := GuardarManifest(rutaManifest, nuevoManifest); err != nil {
		return fmt.Errorf("error al actualizar manifest.json: %w", err)
	}

	fmt.Printf("✂️ Proceso finalizado. Se integraron %d nuevos chunks vectorizados.\n", chunksProcesados)
	return nil
}
