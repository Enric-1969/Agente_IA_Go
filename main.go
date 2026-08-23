package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func formatoDocumentos(resultados []ResultadoBusqueda) string {
	var sb strings.Builder
	for i, res := range resultados {
		sb.WriteString(fmt.Sprintf("--- Fragmento %d (Origen: %s | Similitud: %.4f) ---\n", i+1, res.Documento.Chunk.Origen, res.Similitud))
		sb.WriteString(res.Documento.Chunk.Contenido)
		sb.WriteString("\n\n")
	}
	return sb.String()
}

func consultarAgente(pregunta string, gestor *GestorModelos, vs *VectorStore, modeloEmbedding string) {
	ctx := context.Background()

	// 1. Clasificación de intención mediante core_router.go (usa FormatearPromptRouter internamente)
	fmt.Println("🧠 Analizando intención de la consulta...")
	tipoRuta, err := ClasificarConsultaClave(ctx, gestor, pregunta)
	if err != nil {
		fmt.Printf("⚠️ Error al clasificar consulta, asignando RutaGeneral: %v\n", err)
		tipoRuta = RutaGeneral
	}

	fmt.Printf("🎯 Ruta asignada por el router: [%v]\n", tipoRuta)

	// 2. Búsqueda y recuperación vectorial
	var contexto string
	if tipoRuta == RutaRAG || tipoRuta == RutaGeneral {
		fmt.Println("📂 Generando embedding y recuperando contexto relevante...")
		vecConsulta, err := GenerarEmbedding(pregunta, modeloEmbedding)
		if err == nil {
			resultados := vs.BuscarSimilares(vecConsulta, 3)
			if len(resultados) > 0 {
				contexto = formatoDocumentos(resultados)
				fmt.Printf("✅ Se recuperaron %d fragmentos relevantes.\n", len(resultados))
			} else {
				fmt.Println("ℹ️ No se encontraron coincidencias vectoriales suficientes.")
			}
		} else {
			fmt.Printf("⚠️ No se pudo generar embedding de la consulta: %v\n", err)
		}
	}

	// 3. Formateo exacto usando constantes y funciones de core_prompts.go
	var contenidoUsuario string
	if contexto != "" {
		contenidoUsuario = FormatearPromptRAG(contexto, pregunta)
	} else {
		contenidoUsuario = pregunta
	}

	mensajes := []Mensaje{
		{Role: "system", Content: PromptSistemaBase},
		{Role: "user", Content: contenidoUsuario},
	}

	modeloActual := gestor.ModeloActivo
	fmt.Printf("\n🚀 Consultando a Groq [%s]: '%s'\n%s\n", modeloActual, pregunta, strings.Repeat("-", 60))

	// 4. Invocación al conector de Groq
	respuesta, err := EnviarConsultaGroq(gestor, mensajes)
	if err != nil {
		fmt.Printf("❌ Error al comunicarse con Groq: %v\n", err)
		return
	}

	fmt.Println("\n🤖 Respuesta del Agente:")
	fmt.Println(respuesta)
	fmt.Println("\n" + strings.Repeat("=", 60))
}

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

	// Inicialización con los nombres de campos de structs.go
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
			consultarAgente(preguntaUsuario, gestor, vs, modeloEmbedding)
		}
	}
}
