package main

import (
	"testing"
)

func TestFragmentarTextoEnChunks(t *testing.T) {
	// 1. Caso texto vacío o solo espacios
	if res := FragmentarTextoEnChunks("   \n\t  ", 10, 2); res != nil {
		t.Errorf("Se esperaba nil para texto vacío, pero se obtuvo: %v", res)
	}

	// 2. Caso texto más corto que el tamaño del chunk
	textoCorto := "Hola mundo"
	resCorto := FragmentarTextoEnChunks(textoCorto, 50, 5)
	if len(resCorto) != 1 || resCorto[0] != textoCorto {
		t.Errorf("Se esperaba 1 chunk con '%s', pero se obtuvo: %v", textoCorto, resCorto)
	}

	// 3. Caso texto con caracteres UTF-8 multibyte (ñ, tildes, emojis) y solape
	textoUTF8 := "Cañón_áéíóú_🚀" // 13 runas exactas
	chunks := FragmentarTextoEnChunks(textoUTF8, 6, 2)

	// Paso = 6 - 2 = 4.
	// Chunk 0: runas[0:6]  -> "Cañón_"
	// Chunk 1: runas[4:10] -> "n_áéíó"
	// Chunk 2: runas[8:13] -> "íóú_🚀"
	esperados := []string{"Cañón_", "n_áéíó", "íóú_🚀"}

	if len(chunks) != len(esperados) {
		t.Fatalf("Se esperaban %d chunks, pero se obtuvieron %d: %v", len(esperados), len(chunks), chunks)
	}

	for i, esperado := range esperados {
		if chunks[i] != esperado {
			t.Errorf("En chunk %d se esperaba '%s', pero se obtuvo '%s'", i, esperado, chunks[i])
		}
	}
}

func TestConvertirTextoAChunks(t *testing.T) {
	// Texto menor o igual a 10 caracteres debe filtrarse
	resCorto := ConvertirTextoAChunks("doc.txt", "Corto")
	if len(resCorto) != 0 {
		t.Errorf("Se esperaba 0 chunks para texto <= 10 caracteres, se obtuvo: %d", len(resCorto))
	}

	// Texto válido debe generar Origen, Contenido y HashMD5
	textoValido := "Este es un documento técnico válido para probar el RAG."
	resValido := ConvertirTextoAChunks("manual_go.md", textoValido)

	if len(resValido) != 1 {
		t.Fatalf("Se esperaba 1 ChunkTexto, se obtuvo: %d", len(resValido))
	}

	if resValido[0].Origen != "manual_go.md" {
		t.Errorf("Origen incorrecto: se esperaba 'manual_go.md', se obtuvo '%s'", resValido[0].Origen)
	}

	if resValido[0].HashMD5 == "" {
		t.Error("El campo HashMD5 no debería estar vacío")
	}
}
