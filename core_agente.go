package main

import (
	"context"
	"fmt"
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

// ProcesarConsulta coordina la intención, memoria, herramientas y respuesta del LLM.
func ProcesarConsulta(pregunta string, gestor *GestorModelos, vs *VectorStore, modeloEmbedding string, memoria *GestorMemoria, dirOrigen string) {
	ctx := context.Background()

	// 1. Registrar entrada en la memoria conversacional
	memoria.Agregar("user", pregunta)

	// 2. Clasificar intención
	fmt.Println("🧠 Analizando intención de la consulta...")
	tipoRuta, err := ClasificarConsultaClave(ctx, gestor, pregunta)
	if err != nil {
		fmt.Printf("⚠️ Error al clasificar consulta, asignando RutaGeneral: %v\n", err)
		tipoRuta = RutaGeneral
	}

	fmt.Printf("🎯 Ruta asignada por el router: [%v]\n", tipoRuta)

	// 3. Ejecutar herramientas Go nativas si aplica
	if tipoRuta == RutaHerramienta {
		comando := "info_sistema"
		pLower := strings.ToLower(pregunta)
		if strings.Contains(pLower, "archivo") || strings.Contains(pLower, "documento") || strings.Contains(pLower, "lista") || strings.Contains(pLower, "carpet") {
			comando = "listar_documentos"
		}

		fmt.Printf("⚙️ [Ejecutando herramienta nativa: %s]\n", comando)
		res := DespacharHerramienta(comando, dirOrigen)

		var respTexto string
		if res.Exito {
			respTexto = res.Contenido
		} else {
			respTexto = fmt.Sprintf("Error ejecutando herramienta: %s", res.Error)
		}

		fmt.Println("\n🤖 Respuesta de Herramienta Nativa:")
		fmt.Println(respTexto)
		fmt.Println("\n" + strings.Repeat("=", 60))

		memoria.Agregar("assistant", respTexto)
		return
	}

	// 4. Búsqueda y recuperación vectorial para RAG / General
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

	// 5. Inyección de contexto y construcción del historial
	var contenidoUsuario string
	if contexto != "" {
		contenidoUsuario = FormatearPromptRAG(contexto, pregunta)
	} else {
		contenidoUsuario = pregunta
	}

	historial := memoria.ObtenerHistorial()
	if len(historial) > 0 {
		historial[len(historial)-1].Content = contenidoUsuario
	}

	mensajes := make([]Mensaje, 0, len(historial)+1)
	mensajes = append(mensajes, Mensaje{Role: "system", Content: PromptSistemaBase})
	mensajes = append(mensajes, historial...)

	modeloActual := gestor.ModeloActivo
	fmt.Printf("\n🚀 Consultando a Groq [%s]: '%s'\n%s\n", modeloActual, pregunta, strings.Repeat("-", 60))

	// 6. Envío al proveedor LLM
	respuesta, err := EnviarConsultaGroq(gestor, mensajes)
	if err != nil {
		fmt.Printf("❌ Error al comunicarse con Groq: %v\n", err)
		return
	}

	fmt.Println("\n🤖 Respuesta del Agente:")
	fmt.Println(respuesta)
	fmt.Println("\n" + strings.Repeat("=", 60))

	memoria.Agregar("assistant", respuesta)
}
