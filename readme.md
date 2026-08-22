# 🌱 Agente_IA_Go

**Agente_IA_Go** es un sistema RAG (*Retrieval-Augmented Generation*) ultraligero, modular y concurrente escrito en Go, diseñado para procesar documentación técnica local e interactuar con modelos de lenguaje a través de la API de Groq.

## Estructura del Proyecto

```text
Agente_IA_Go/
├── structs.go                 # Definición centralizada de las estructuras de datos
├── config_rutas_locales.go     # Escaneo dinámico de carpetas de entrada y almacenamiento
├── config_validador.go         # Verificación de compatibilidad con Tool Calling
├── config_crea_lista_models.go # Escaneo concurrente y rotación del catálogo activo de Groq
├── config_conector_groq.go     # Cliente HTTP nativo, reconexión y fallback de modelos
├── core_manifest_hash.go       # Gestión de manifest.json, firmas MD5 y validación de cambios
├── core_fragmentar_chunks.go   # Fragmentación UTF-8 segura (runas) y solape de texto
├── core_embeddings.go         # Cliente HTTP nativo para generación de vectores con Ollama
├── core_vectorstore.go        # VectorStore nativo en memoria, persistencia .gob y similitud coseno
├── core_ingestar.go            # Coordinador principal de ingesta con control de memoria (>50MB)
├── core_router.go              # Clasificación de preguntas por dominio técnico
├── core_prompts.go             # Plantillas de sistema y prompts para el agente
├── extractor_orquestador.go    # Despachador principal de ingesta según extensión
├── extractor_texto.go          # Extracción para archivos .txt, .md y código fuente
├── extractor_excel.go          # Extracción estructurada multipestaña para .xlsx y .xlsm
├── extractor_word.go           # Extracción de contenido estructurado para .docx
├── extractor_pdf.go            # Extracción nativa de PDF con fallback OCR WASM concurrente
├── main.go                     # CLI interactiva y bucle de orquestación RAG
├── storage/                    # Archivos de persistencia (.gob y manifest.json)
├── .env
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

## Configuración del Entorno (.env)

Antes de ejecutar la aplicación, crea un archivo .env en la raíz del proyecto (puedes duplicar .env.example) y configura las siguientes variables:
```
# Configuración de API Groq
# Groq no es un modelo, es el motor donde se ejecutan los modelos.
# 1. Ve a https://console.groq.com/home y crea una cuenta gratuita.
# 2. Genera una API key en https://console.groq.com/keys.
GROQ_API_KEY="gsk_tu_clave_aqui"
GROQ_MODEL='qwen/qwen3.6-27b'
GROQ_MODELOS_VALIDOS='qwen/qwen3.6-27b'

# Configuración de Embeddings Locales (Ollama)
OLLAMA_EMBEDDING_MODEL="all-minilm"
OLLAMA_HOST="http://localhost:11434"

# Rutas configurables del sistema
# Windows local/red: "Z:/Files_IA" o "C:/MisDocumentos/IA"
# Linux/macOS: "/var/data/files_ia"
RUTA_DOCS_ORIGEN="./documentos_entrada"
RUTA_VECTORSTORES="./mis_vectorstores"
```

## Requisitos y Dependencias

Asegúrate de contar con Go 1.24+ instalado. Para descargar e instalar las dependencias externas del proyecto, ejecuta:

```bash
# Lectura y procesamiento de archivos Excel (.xlsx, .xlsm)
go get github.com/xuri/excelize/v2
# Lectura y procesamiento de archivos PDF (.pdf)
go get github.com/ledongthuc/pdf
# Iniciar el servicio y descargar el modelo de embeddings ligero
ollama pull all-minilm
```
