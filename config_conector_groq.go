package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EnviarConsultaGroq transmite los mensajes a la API de Groq e invoca la rotación automática si el modelo en uso falla.
func EnviarConsultaGroq(gestor *GestorModelos, mensajes []Mensaje) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	for {
		modeloActual := gestor.ObtenerModeloActivo()
		if modeloActual == "" {
			return "", fmt.Errorf("no hay modelos disponibles en la cola de reserva")
		}

		peticion := PeticionGroq{
			Model:       modeloActual,
			Messages:    mensajes,
			Temperature: 0.2,
		}

		body, err := json.Marshal(peticion)
		if err != nil {
			return "", fmt.Errorf("error al serializar la petición JSON: %w", err)
		}

		req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
		if err != nil {
			return "", fmt.Errorf("error al preparar la solicitud HTTP: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+gestor.apiKey)

		resp, err := client.Do(req)

		// Si la API falla o la red cae, descartamos el modelo y reintentamos con el siguiente
		if err != nil || resp.StatusCode != http.StatusOK {
			statusCode := 0
			if resp != nil {
				statusCode = resp.StatusCode
				resp.Body.Close()
			}

			fmt.Printf("⚠️ Error al consultar con '%s' (HTTP %d). Conmutando modelo...\n", modeloActual, statusCode)
			siguiente := gestor.ReportarFalloYRotar(modeloActual)
			if siguiente == "" {
				return "", fmt.Errorf("se agotaron los modelos de la reserva y ninguno respondió")
			}
			continue
		}

		var respuesta RespuestaGroq
		err = json.NewDecoder(resp.Body).Decode(&respuesta)
		resp.Body.Close()

		if err != nil {
			return "", fmt.Errorf("error al decodificar la respuesta JSON: %w", err)
		}

		if len(respuesta.Choices) == 0 {
			return "", fmt.Errorf("la respuesta de Groq no contiene elecciones de texto")
		}

		return respuesta.Choices[0].Message.Content, nil
	}
}
