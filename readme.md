# 🌱 Agente_IA_Go

**Agente_IA_Go** es 

## Estructura unificada para Agente_IA_Go
```
Agente_IA_Go/
├── structs.go                    # Definición centralizada de las estructuras de datos
├── config_rutas_locales.go       # Escaneo dinámico de la unidad compartida Z:\Files_IA
├── config_validador.go           # Verificación de compatibilidad con Tool Calling
├── config_crea_lista_models.go   # Descarga y filtrado del catálogo activo de Groq
├── config_conector_groq.go       # Cliente HTTP nativo, rotación y reintentos de red
├── core_ingestar.go              # Escaneo MD5, chunking recursivo y almacenamiento
├── core_router.go                # Clasificación de la pregunta por dominio técnico
├── core_prompts.go               # Constantes con las plantillas del sistema
├── main.go                       # Bucle principal de la interfaz
├── mis_vectorstores/             # Bases de datos SQLite independientes por dominio
├── .env
├── .env.example
└── go.mod
'''
