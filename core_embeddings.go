package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OllamaEmbeddingRequest define la estructura para solicitar embeddings a Ollama.
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse representa la respuesta con el vector numérico devuelto por Ollama.
type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerarEmbedding envía un texto a Ollama local y retorna su representación vectorial []float32.
func GenerarEmbedding(texto string, modelo string) ([]float32, error) {
	if modelo == "" {
		modelo = "all-minilm"
	}

	url := "http://localhost:11434/api/embeddings"
	reqBody := OllamaEmbeddingRequest{
		Model:  modelo,
		Prompt: texto,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error al serializar solicitud de embedding: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error al conectar con Ollama en %s (verifica que el servicio esté activo): %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama devolvió un código de estado de error: %d", resp.StatusCode)
	}

	var ollamaResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta del embedding: %w", err)
	}

	if len(ollamaResp.Embedding) == 0 {
		return nil, fmt.Errorf("Ollama devolvió un vector de embeddings vacío")
	}

	return ollamaResp.Embedding, nil
}
