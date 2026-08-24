# 🌱 Agente_IA_Go

**Agente_IA_Go** es un sistema RAG (*Retrieval-Augmented Generation*) ultraligero, modular y concurrente escrito en Go, diseñado para procesar documentación técnica local e interactuar con modelos de lenguaje a través de la API de Groq.

## Estructura del Proyecto

```text
Agente_IA_Go/
├── structs.go                  # Definición centralizada de tipos y estructuras de datos
├── config_rutas_locales.go     # Escaneo dinámico de carpetas de entrada y almacenamiento
├── config_validador.go         # Verificación de compatibilidad con Tool Calling
├── config_crea_lista_models.go # Escaneo concurrente y rotación del catálogo activo de Groq
├── config_conector_groq.go     # Cliente HTTP nativo, reconexión y fallback de modelos
├── core_agente.go              # Orquestador del ciclo RAG, gestión de memoria y flujo de consultas
├── core_bm25.go                # Indización y puntuación léxica exactas por palabras clave (BM25)
├── core_rrf.go                 # Algoritmo de fusión matemática de rankings (Reciprocal Rank Fusion)
├── core_memoria.go             # Gestor de memoria conversacional multi-turno (sliding window)
├── core_herramientas.go        # Catálogo y despachador de herramientas nativas en Go (Tool Calling)
├── core_manifest_hash.go       # Gestión de manifest.json, firmas MD5 y validación de cambios
├── core_fragmentar_chunks.go   # Fragmentación UTF-8 segura (runas) y solape de texto
├── core_embeddings.go          # Cliente HTTP nativo para generación de vectores con Ollama
├── core_vectorstore.go         # VectorStore nativo en memoria, persistencia .gob y similitud coseno
├── core_ingestar.go            # Coordinador principal de ingesta con control de memoria (>50MB)
├── core_router.go              # Enrutador híbrido (heurístico + LLM) de intenciones (RAG/Herramienta/General)
├── core_prompts.go             # Plantillas de sistema anti-alucinaciones y formateadores RAG
├── extractor_orquestador.go    # Despachador principal de ingesta según extensión
├── extractor_texto.go          # Extracción para archivos .txt, .md y código fuente
├── extractor_excel.go          # Extracción estructurada multipestaña para .xlsx y .xlsm
├── extractor_word.go           # Extracción de contenido estructurado para .docx
├── extractor_pdf.go            # Extracción nativa de PDF con fallback OCR WASM concurrente
├── main.go                     # Punto de entrada CLI/REPL e inicialización del entorno
├── documentos_entrada/         # Carpeta de entrada para documentos a procesar
├── mis_vectorstores/           # Almacenamiento persistente (.gob y manifest.json)
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
## 📐 Decisiones de Arquitectura y Diseño

Este proyecto fue diseñado priorizando la **portabilidad, privacidad y velocidad de ejecución nativa**, evitando la sobreingeniería y las dependencias pesadas habituales en el ecosistema de IA.

### 1. Ingesta y Embeddings 100% Locales (Soberanía de Datos)
* **Decisión:** Comunicación HTTP nativa con **Ollama** (`core_embeddings.go`) en lugar de servicios en la nube.
* **Por qué:** Garantiza privacidad absoluta (ningún dato sale de la máquina local), elimina costes de API de terceros y suprime la latencia de red externa.

### 2. VectorStore Nativo en Go Pure + Persistencia Binaria `.gob`
* **Decisión:** Base de datos vectorial en memoria (`core_vectorstore.go`) con serialización en formato binario nativo de Go (`.gob`).
* **Por qué:** Elimina la necesidad de desplegar contenedores Docker o bases de datos externas. El formato `.gob` ofrece una tasa de compresión y velocidad de hidratación superior a JSON. La búsqueda por similitud de coseno en RAM responde en menos de 5 ms sin consumo de tokens de la API.

### 3. Enrutamiento Híbrido de Cero-Coste
* **Decisión:** Clasificación en dos fases (`core_router.go`): reglas heurísticas en local primero, evaluación por LLM en caso de ambigüedad.

* **Por qué:** Resuelve consultas directas en microsegundos sin consumir tokens de API ni incurrir en latencia de red.

### 4. Motor Aislado OCR en WASM
* **Decisión:** Extracción de texto PDF dual (`extractor_pdf.go`) con motor de respaldo OCR bajo WebAssembly.
* **Por qué:** Mantiene el principio de ejecutable único y portable sin requerir la instalación de binarios pesados en el sistema operativo del cliente (C++ / Tesseract tradicional).

### 5. Ingesta Incremental Inteligente por Hashing MD5
* **Decisión:** Registro de estado mediante firmas MD5 por archivo (`core_manifest_hash.go`) y control interactivo de memoria RAM (`core_ingestar.go`).
* **Por qué:** Evita la re-vectorización innecesaria de documentos no modificados (ahorrando tiempo de CPU) y protege el sistema contra problemas de falta de memoria (*Out Of Memory*) al procesar archivos de gran tamaño (>50 MB).

### 6. Fragmentación UTF-8 Segura
* **Decisión:** Procesamiento de texto por vectores de runas (`core_fragmentar_chunks.go`) con solape (*overlap*) de 150 caracteres.
* **Por qué:** Evita la corrupción de caracteres multibyte (tildes, `ñ`, símbolos) al cortar palabras y preserva el contexto semántico entre bloques contiguos.

### 7. Cliente HTTP Nativo con Resiliencia y Fallback de Modelos
* **Decisión:** Implementación de un cliente HTTP sin SDKs de terceros (`config_conector_groq.go`) con rotación dinámica y automática entre ModelosValidos (`GestorModelos`).
* **Por qué:** Elimina dependencias pesadas de SDKs externos y evita interrupciones en el servicio si un modelo específico de Groq alcanza límites de frecuencia (rate limits) o sufre una caída temporal.

### 8. Pipeline Polimórfico de Extracción Multiformato
* **Decisión:** Despachador modular por extensiones (`extractor_orquestador.go`) con extractores nativos dedicados para .docx, .xlsx/.xlsm, .pdf y código/texto plano.
* **Por qué:** Preserva la estructura original de tablas y documentos jerárquicos sin depender de scripts externos de Python o microservicios de parseo.

### 9. Memoria Conversacional Multi-turno con Ventana Deslizante
* **Decisión:** Gestor de historial en memoria (`core_memoria.go`) basado en ventana deslizante (*sliding window*) thread-safe mediante `sync.RWMutex`.
* **Por qué:** Mantiene la coherencia contextual en conversaciones extensas sin saturar la ventana de tokens del modelo ni incrementar la latencia por procesamiento de historial acumulado.

### 10. Despacho Nativo de Herramientas (Tool Calling)
* **Decisión:** Módulo de herramientas Go nativas (`core_herramientas.go`) activado dinámicamente tras la clasificación por `RutaHerramienta`.
* **Por qué:** Permite ejecutar operaciones del sistema, inspección de archivos y cálculos a velocidad nativa en local, reduciendo el consumo de tokens y evitando la necesidad de entornos de ejecución secundarios.

### 11. Búsqueda Híbrida (BM25 + Similitud Coseno) mediante Reciprocal Rank Fusion (RRF)
* **Decisión:** Combinación de búsqueda vectorial semántica con el algoritmo de frecuencia de términos BM25 (`core_bm25.go`), fusionados dinámicamente mediante RRF (`core_rrf.go`).
* **Por qué:** Resuelve las limitaciones del RAG puramente vectorial en consultas que contienen nombres de funciones, códigos de error o identificadores exactos, ofreciendo máxima precisión contextual antes de invocar al modelo.