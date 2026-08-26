package main

import (
	"context"
	"fmt"
	"strings"
)

func formatoDocumentos(resultados []ResultadoHibrido) string {
	var sb strings.Builder
	for i, res := range resultados {
		sb.WriteString(fmt.Sprintf("--- Fragmento %d (Origen: %s | Similitud: %.4f) ---\n", i+1, res.Documento.Chunk.Origen, res.ScoreRRF))
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

	// 4. Búsqueda y recuperación vectorial para RAG / General (Two-Stage Retrieval)
	var contexto string
	if tipoRuta == RutaRAG || tipoRuta == RutaGeneral {
		fmt.Println("📂 Generando embedding y recuperando candidatos RRF (Fase 1)...")
		vecConsulta, err := GenerarEmbedding(pregunta, modeloEmbedding)
		if err == nil {
			// Lectura dinámica desde .env con blindaje contra valores bajos
			limiteK1 := obtenerEnvInt("LIMITE_CANDIDATOS_RECUPERACION", 16, 16)
			limiteK2 := obtenerEnvInt("LIMITE_FRAGMENTOS_FINALES", 6, 6)
			pesoRRF := obtenerEnvFloat("PESO_PUNTUACION_RRF", 0.7)
			pesoDensidad := obtenerEnvFloat("PESO_DENSIDAD_LEXICA", 0.3)

			// Fase 1: Cobertura amplia (K1 candidatos)
			candidatos := vs.BuscarHibrido(pregunta, vecConsulta, limiteK1)

			if len(candidatos) > 0 {
				fmt.Printf("🔍 Fase 1: %d candidatos recuperados. Aplicando Re-Ranking Top-%d (Fase 2)...\n", len(candidatos), limiteK2)

				// Fase 2: Re-Ranking por densidad léxica y recorte al Top-N (K2 final)
				resultadosFinales := EjecutarReordenamiento(pregunta, candidatos, limiteK2, pesoRRF, pesoDensidad)

				contexto = formatoDocumentos(resultadosFinales)
				fmt.Printf("✅ Fase 2: %d fragmentos finales seleccionados para el prompt.\n", len(resultadosFinales))
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

	modeloActual := gestor.ObtenerModeloActivo()
	fmt.Printf("\n🚀 Consultando a Groq [%s]: '%s'\n%s\n", modeloActual, pregunta, strings.Repeat("-", 60))

	// 6. Envío al proveedor LLM
	respuesta, err := EnviarConsultaGroq(gestor, mensajes)
	if err != nil {
		fmt.Printf("❌ Error al comunicarse con Groq: %v\n", err)
		return
	}

	if idx := strings.LastIndex(respuesta, "</think>"); idx != -1 {
		respuesta = strings.TrimSpace(respuesta[idx+len("</think>"):])
	}

	fmt.Println("\n🤖 Respuesta del Agente:")
	fmt.Println(respuesta)
	fmt.Println("\n" + strings.Repeat("=", 60))

	memoria.Agregar("assistant", respuesta)
}
