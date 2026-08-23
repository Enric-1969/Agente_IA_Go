package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("=== AGENTE RAG MULTI-DOMINIO ACTIVO (GO) ===")
	fmt.Println("===========================================")
	fmt.Println("Escribe 'salir' o 'exit' para terminar.")

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: GROQ_API_KEY no está configurada en el entorno (.env).")
		return
	}

	modelosEnv := os.Getenv("GROQ_MODELOS_VALIDOS")
	if modelosEnv == "" {
		modelosEnv = "qwen/qwen3.6-27b"
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

	dirOrigen := "./documentos_entrada"
	dirVectorStore := "./mis_vectorstores"
	archivoGOB := filepath.Join(dirVectorStore, "vectorstore.gob")
	modeloEmbedding := "all-minilm"

	_ = os.MkdirAll(dirOrigen, 0755)
	_ = os.MkdirAll(dirVectorStore, 0755)

	vs := NuevoVectorStore(archivoGOB)
	if err := vs.CargarPersistencia(); err != nil {
		fmt.Printf("ℹ️ Estado almacén vectorial: %v\n", err)
	} else {
		fmt.Printf("📚 Base vectorial cargada desde %s (%d vectores activos en RAM)\n", archivoGOB, len(vs.Documentos))
	}

	memoria := NuevoGestorMemoria(6)
	_ = ProcesarDominioIngesta(dirOrigen, dirVectorStore, modeloEmbedding, vs)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n💬 Tu pregunta: ")
		preguntaUsuario, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		preguntaUsuario = strings.TrimSpace(preguntaUsuario)
		cmdLower := strings.ToLower(preguntaUsuario)
		if cmdLower == "salir" || cmdLower == "exit" || cmdLower == "quit" {
			fmt.Println("👋 ¡Hasta luego!")
			break
		}

		if preguntaUsuario != "" {
			ProcesarConsulta(preguntaUsuario, gestor, vs, modeloEmbedding, memoria, dirOrigen)
		}
	}
}
