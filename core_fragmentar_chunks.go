package main

import (
	"fmt"
	"strings"
)

// --- FRAGMENTACIÓN Y CREACIÓN DE CHUNKS ---

func DividirTextoEnChunks(texto string, tamanio, solape int) []string {
	var chunks []string
	runas := []rune(texto)
	total := len(runas)
	if total == 0 {
		return chunks
	}

	for inicio := 0; inicio < total; {
		fin := inicio + tamanio
		if fin > total {
			fin = total
		}
		chunks = append(chunks, string(runas[inicio:fin]))
		if fin == total {
			break
		}
		inicio += (tamanio - solape)
	}
	return chunks
}

func ConvertirTextoAChunks(relPath, texto string) []ChunkTexto {
	var resultado []ChunkTexto
	fragmentos := DividirTextoEnChunks(texto, 1000, 200)

	for i, frag := range fragmentos {
		if len(strings.TrimSpace(frag)) <= 10 {
			continue
		}
		resultado = append(resultado, ChunkTexto{
			ID:        fmt.Sprintf("%s_chunk_%d", relPath, i),
			Origen:    relPath,
			Contenido: frag,
			HashMD5:   GenerarHashMD5(frag),
			Embedding: []float32{},
		})
	}
	return resultado
}
