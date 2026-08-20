package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ValidarEntornoCompleto verifica que las variables de entorno principales existan.
func ValidarEntornoCompleto() error {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" || apiKey == "gsk_tu_clave_aqui" {
		return fmt.Errorf("la variable GROQ_API_KEY no está configurada")
	}

	modelo := os.Getenv("GROQ_MODEL")
	if modelo == "" {
		return fmt.Errorf("la variable GROQ_MODEL no está definida en el archivo .env")
	}

	return nil
}

// ProbarSoporteToolCalling comprueba mediante una petición HTTP si un modelo soporta herramientas.
func ProbarSoporteToolCalling(apiKey string, modelID string) bool {
	url := "https://api.groq.com/openai/v1/chat/completions"

	payload := map[string]interface{}{
		"model": modelID,
		"messages": []map[string]string{
			{"role": "user", "content": "Responde True"},
		},
		"tools": []map[string]interface{}{
			{
				"type": "function",
				"function": map[string]interface{}{
					"name":        "prueba_tool",
					"description": "Confirmacion de prueba",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"confirmacion": map[string]string{"type": "boolean"},
						},
						"required": []string{"confirmacion"},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
