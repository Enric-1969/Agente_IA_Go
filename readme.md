# 🌱 Agente_IA_Go

**Agente_IA_Go** es un sistema RAG (*Retrieval-Augmented Generation*) ultraligero, modular y concurrente escrito en Go, diseñado para procesar documentación técnica local e interactuar con modelos de lenguaje a través de la API de Groq.

## Estructura del Proyecto

```text
Agente_IA_Go/
├── structs.go                  # Definición centralizada de las estructuras de datos
├── config_rutas_locales.go      # Escaneo dinámico de la unidad compartida (ej. Z:\Files_IA)
├── config_validador.go          # Verificación de compatibilidad con Tool Calling
├── config_crea_lista_models.go  # Escaneo concurrente y rotación del catálogo activo de Groq
├── config_conector_groq.go      # Cliente HTTP nativo, reconexión y fallback de modelos
├── core_manifest_hash.go        # Gestión de manifest.json, firmas MD5 y validación de archivos
├── core_fragmentar_chunks.go    # Fragmentación de texto UTF-8 seguro (runas) y solape (overlap)
├── core_ingestar.go             # Coordinador principal del flujo de ingesta incremental
├── core_router.go               # Clasificación de la pregunta por dominio técnico
├── core_prompts.go              # Plantillas de sistema y prompts para el agente
├── extractor_orquestador.go     # Despachador principal de ingesta según extensión del archivo
├── extractor_texto.go           # Extracción para archivos .txt, .md y código fuente
├── extractor_excel.go           # Extracción estructurada multipestaña para .xlsx y .xlsm
├── extractor_word.go            # Extracción de contenido estructurado para .docx
├── extractor_pdf.go             # Extracción y filtrado inteligente de páginas para .pdf
├── main.go                      # Punto de entrada y bucle principal de la aplicación
├── mis_vectorstores/            # Bases de datos vectoriales y manifiestos por dominio
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

# Rutas configurables del sistema
# Windows local/red: "Z:/Files_IA" o "C:/MisDocumentos/IA"
# Linux/macOS: "/var/data/files_ia"
RUTA_DOCS_ORIGEN="./documentos_entrada"
RUTA_VECTORSTORES="./mis_vectorstores"
```

## Requisitos y Dependencias

Asegúrate de contar con Go 1.24+ instalado. Para descargar e instalar las dependencias externas del proyecto, ejecuta:

```bash
### Lectura y procesamiento de archivos Excel (.xlsx, .xlsm)
go get [github.com/xuri/excelize/v2](https://github.com/xuri/excelize/v2)
```