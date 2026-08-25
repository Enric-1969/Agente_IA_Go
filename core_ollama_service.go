package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// AsegurarOllamaActivo verifica si Ollama responde en 127.0.0.1:11434.
// Si no responde, comprueba si el ejecutable está instalado en el sistema antes de iniciar 'ollama serve'.
// Devuelve una función de limpieza (cleanup) para cerrar el subproceso solo si fue abierto por el agente.
func AsegurarOllamaActivo() (func(), error) {
	urlPing := "http://127.0.0.1:11434/api/version"
	client := http.Client{Timeout: 2 * time.Second}

	// 1. Comprobar si Ollama ya está activo
	resp, err := client.Get(urlPing)
	if err == nil {
		resp.Body.Close()
		fmt.Println("ℹ️  Servicio Ollama detectado en ejecución previa (127.0.0.1:11434).")
		return func() {}, nil // No se requiere apagar nada al terminar
	}

	// 2. Comprobar si el ejecutable existe en el PATH del sistema operativo
	_, errPath := exec.LookPath("ollama")
	if errPath != nil {
		return nil, fmt.Errorf("❌ Ollama no está instalado en el sistema operativo.\n👉 Por favor, descárgalo e instálalo desde: https://ollama.com/download")
	}

	// 3. Arrebatar el servicio iniciando el subproceso 'ollama serve'
	fmt.Println("🔄 Ollama inactivo. Levantando servicio local ('ollama serve')...")
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error al iniciar el subproceso de Ollama: %w", err)
	}

	// Esperar brevemente a que el socket levante
	time.Sleep(3 * time.Second)

	// Función de limpieza ejecutada mediante defer en main.go
	cleanup := func() {
		if cmd.Process != nil {
			fmt.Println("🧹 Cerrando subproceso de Ollama instanciado por el agente...")
			_ = cmd.Process.Kill()
		}
	}

	return cleanup, nil
}
