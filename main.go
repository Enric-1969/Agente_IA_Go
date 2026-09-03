package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("=== AGENTE RAG MULTI-DOMINIO ACTIVO (GO) ===")
	fmt.Println("===========================================")
	fmt.Println("Escribe 'salir', 'exit' o 'quit' para terminar.")

	// Comprobar y gestionar el ciclo de vida del servicio Ollama
	var cleanupOllama func()
	var err error

	cleanupOllama, err = AsegurarOllamaActivo()
	if err != nil {
		fmt.Printf("⚠️ Advertencia servicio Ollama: %v\n", err)
	} else if cleanupOllama != nil {
		defer cleanupOllama()
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: GROQ_API_KEY no está configurada en el entorno (.env).")
		return
	}

	modelosEnv := os.Getenv("GROQ_MODELOS_VALIDOS")
	if modelosEnv == "" {
		modelosEnv = "qwen/qwen3.8-27b"
	}

	var modelos []string
	for _, m := range strings.Split(modelosEnv, ",") {
		mTr := strings.TrimSpace(m)
		if mTr != "" {
			modelos = append(modelos, mTr)
		}
	}

	gestor := &GestorModelos{
		APIKey:         apiKey,
		ModelosValidos: modelos,
		ModeloActivo:   modelos[0],
	}

	// Carga de rutas desde el entorno con valores por defecto
	dirOrigen := os.Getenv("RUTA_DOCS_ORIGEN")
	if dirOrigen == "" {
		dirOrigen = "./documentos_entrada"
	}

	dirVectorStore := os.Getenv("RUTA_VECTORSTORES")
	if dirVectorStore == "" {
		dirVectorStore = "./mis_vectorstores"
	}

	dirHistoricoConversaciones := os.Getenv("RUTA_HISTORICO_CONVERSACIONES")
	if dirHistoricoConversaciones == "" {
		dirHistoricoConversaciones = "./Historico_Conversaciones"
	}

	archivoGOB := filepath.Join(dirVectorStore, "vectorstore.gob")

	modeloEmbedding := os.Getenv("MODELO_EMBEDDING")
	if modeloEmbedding == "" {
		modeloEmbedding = "nomic-embed-text"
	}

	_ = os.MkdirAll(dirOrigen, 0755)
	_ = os.MkdirAll(dirVectorStore, 0755)
	_ = os.MkdirAll(dirHistoricoConversaciones, 0755)

	vs := NuevoVectorStore(archivoGOB)
	if err := vs.CargarPersistencia(); err != nil {
		fmt.Printf("ℹ️ Estado almacén vectorial: %v\n", err)
	} else {
		fmt.Printf("📚 Base vectorial cargada desde %s (%d vectores activos en RAM)\n", archivoGOB, len(vs.Documentos))
	}

	memoria := NuevoGestorMemoria(6)
	_ = ProcesarDominioIngesta(dirOrigen, dirVectorStore, modeloEmbedding, vs)

	// Helper para sanitizar el asunto y generar la ruta
	obtenerRutaExportacion := func(asunto string) string {
		timestamp := time.Now().Format("2006-01-02_15-04-05")

		asuntoLimpio := strings.TrimSpace(asunto)
		asuntoLimpio = strings.ReplaceAll(asuntoLimpio, " ", "_")
		for _, char := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"} {
			asuntoLimpio = strings.ReplaceAll(asuntoLimpio, char, "")
		}

		var nombreFichero string
		if asuntoLimpio != "" {
			nombreFichero = fmt.Sprintf("sesion_%s_%s.md", timestamp, asuntoLimpio)
		} else {
			nombreFichero = fmt.Sprintf("sesion_%s.md", timestamp)
		}

		return filepath.Join(dirHistoricoConversaciones, nombreFichero)
	}

	// --- CAPTURA ELEGANTE DE CTRL+C (SIGINT / SIGTERM) ---
	canalSenales := make(chan os.Signal, 1)
	signal.Notify(canalSenales, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-canalSenales
		fmt.Println("\n\n⚠️ Interrupción detectada (Ctrl+C). Guardando estado y cerrando de forma segura...")

		rutaMD := obtenerRutaExportacion("interrupcion")
		_ = ExportarHistorialMarkdown(memoria, rutaMD)
		fmt.Printf("💾 Sesión guardada en: %s\n", rutaMD)

		if cleanupOllama != nil {
			cleanupOllama()
		}
		fmt.Println("👋 ¡Hasta luego!")
		os.Exit(0)
	}()
	// -----------------------------------------------------

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n💬 Tu pregunta: ")
		preguntaUsuario, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		preguntaUsuario = strings.TrimSpace(preguntaUsuario)

		// Si presiona Enter directamente sin escribir texto
		if preguntaUsuario == "" {
			fmt.Println("ℹ️ Escribe una pregunta o escribe 'salir' / 'exit' para finalizar la sesión.")
			continue
		}

		cmdLower := strings.ToLower(preguntaUsuario)
		if cmdLower == "salir" || cmdLower == "exit" || cmdLower == "quit" {
			fmt.Print("📝 Introduce el asunto o tema de la sesión (Enter para omitir): ")
			asuntoInput, _ := reader.ReadString('\n')

			rutaMD := obtenerRutaExportacion(asuntoInput)
			_ = ExportarHistorialMarkdown(memoria, rutaMD)
			fmt.Printf("💾 Sesión guardada automáticamente en: %s\n", rutaMD)
			fmt.Println("👋 ¡Hasta luego!")
			break
		}

		ProcesarConsulta(preguntaUsuario, gestor, vs, modeloEmbedding, memoria, dirOrigen)
	}
}
