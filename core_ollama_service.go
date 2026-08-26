package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func AsegurarOllamaActivo() (func(), error) {
	urlPing := "http://127.0.0.1:11434/api/version"
	// Reducimos el timeout de 2s a 500ms para no demorar el chequeo inicial
	client := http.Client{Timeout: 500 * time.Millisecond}

	// 1. Comprobar si Ollama ya está activo
	resp, err := client.Get(urlPing)
	if err == nil {
		resp.Body.Close()
		fmt.Println("ℹ️  Servicio Ollama detectado en ejecución previa (127.0.0.1:11434).")
		return func() {}, nil
	}

	// 2. Comprobar ejecutable en PATH
	_, errPath := exec.LookPath("ollama")
	if errPath != nil {
		return nil, fmt.Errorf("❌ Ollama no está instalado en el sistema operativo.\n👉 Descárgalo en: https://ollama.com/download")
	}

	// 3. Levantar servicio local
	fmt.Println("🔄 Ollama inactivo. Levantando servicio local ('ollama serve')...")
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error al iniciar el subproceso de Ollama: %w", err)
	}

	// Reemplazo de time.Sleep(3s): Polling activo cada 200ms durante máximo 3 segundos
	ollamaListo := false
	for i := 0; i < 15; i++ {
		time.Sleep(200 * time.Millisecond)
		r, errPing := client.Get(urlPing)
		if errPing == nil {
			r.Body.Close()
			ollamaListo = true
			break
		}
	}

	if !ollamaListo {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("el servicio Ollama tardó demasiado en responder en el puerto 11434")
	}

	cleanup := func() {
		if cmd.Process != nil {
			fmt.Println("🧹 Cerrando subproceso de Ollama instanciado por el agente...")
			_ = cmd.Process.Kill()
		}
	}

	return cleanup, nil
}
