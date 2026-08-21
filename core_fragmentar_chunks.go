package main

import "strings"

// FragmentarTextoEnChunks divide un texto en bloques usando runas UTF-8
// para no romper caracteres multibyte (ñ, tildes, símbolos).
func FragmentarTextoEnChunks(texto string, tamanoChunk int, solape int) []string {
	textoLimpio := strings.TrimSpace(texto)
	if textoLimpio == "" {
		return nil
	}

	runas := []rune(textoLimpio)
	totalRunas := len(runas)

	if totalRunas <= tamanoChunk {
		return []string{textoLimpio}
	}

	var chunks []string
	paso := tamanoChunk - solape
	if paso <= 0 {
		paso = 1
	}

	for i := 0; i < totalRunas; i += paso {
		fin := i + tamanoChunk
		if fin > totalRunas {
			fin = totalRunas
		}

		chunk := string(runas[i:fin])
		chunks = append(chunks, chunk)

		if fin == totalRunas {
			break
		}
	}

	return chunks
}

// ConvertirTextoAChunks toma el texto completo de un archivo, lo trocea y
// devuelve la colección reutilizando la estructura ChunkTexto de structs.go.
func ConvertirTextoAChunks(relPath, texto string) []ChunkTexto {
	bloques := FragmentarTextoEnChunks(texto, 1000, 150)
	var resultado []ChunkTexto

	for _, bloque := range bloques {
		if len(strings.TrimSpace(bloque)) > 10 {
			resultado = append(resultado, ChunkTexto{
				Origen:    relPath,
				Contenido: bloque,
				HashMD5:   GenerarHashMD5(bloque), // Aprovechamos la función de core_manifest_hash.go
			})
		}
	}
	return resultado
}
