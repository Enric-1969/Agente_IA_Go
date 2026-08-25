package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// AsegurarOllamaActivo verifica el estado del puerto local 11434.
// Si el servicio no responde, inicia 'ollama serve' en segundo plano
// y retorna una función cleanup para cerrar el proceso al finalizar el agente.
func AsegurarOllamaActivo() (func(), error) {
	urlPing := "http://localhost:11434/api/version"
	client := http.Client{Timeout: 1 * time.Second}

	// 1. Verificar si Ollama ya está activo en el sistema
	_, err := client.Get(urlPing)
	if err == nil {
		fmt.Println("🟢 Servicio Ollama detectado y activo.")
		return func() {}, nil
	}

	// 2. Levantar el servicio si no está corriendo
	fmt.Println("⚡ Ollama no responde. Arrancando servicio local 'ollama serve'...")
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("fallo al ejecutar 'ollama serve': %w", err)
	}

	// 3. Reintentar conexión hasta confirmar disponibilidad
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		if _, err := client.Get(urlPing); err == nil {
			fmt.Println("✅ Servicio Ollama iniciado con éxito.")
			cleanup := func() {
				fmt.Println("🛑 Apagando el proceso de Ollama iniciado por el agente...")
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
			}
			return cleanup, nil
		}
	}

	return nil, fmt.Errorf("Ollama arrancó pero no aceptó conexiones en http://localhost:11434")
}
