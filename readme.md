# 🌱 Agente_IA_Go

**Agente_IA_Go** es un sistema RAG (*Retrieval-Augmented Generation*) ultraligero, modular y concurrente escrito en Go, diseñado para procesar documentación técnica local e interactuar con modelos de lenguaje a través de la API de Groq.

## Estructura del Proyecto

```text
Agente_IA_Go/
├── structs.go                 # Definición centralizada de tipos y estructuras de datos
├── config_rutas_locales.go    # Escaneo dinámico de carpetas de entrada y almacenamiento
├── config_validador.go        # Verificación de compatibilidad con Tool Calling
├── config_crea_lista_models.go# Escaneo concurrente y rotación del catálogo activo de Groq
├── config_conector_groq.go    # Cliente HTTP nativo, reconexión y fallback de modelos
├── core_agente.go             # Orquestador del ciclo RAG, gestión de memoria y flujo de consultas
├── core_agente_helpers.go     # Formateo y serialización de contexto RAG para el orquestador
├── core_bm25.go               # Indización y puntuación léxica exactas por palabras clave (BM25)
├── core_bm25_helpers.go       # Tokenización, limpieza y normalización de texto aislada para BM25
├── core_calcular_relevancia.go   # Cálculo de densidad y re-puntuación
├── core_embeddings.go         # Cliente HTTP nativo para generación de vectores con Ollama
├── core_fragmentar_chunks.go  # Fragmentación UTF-8 segura (runas) y solape de texto

├── core_herramientas.go       # Catálogo y despachador de herramientas nativas en Go (Tool Calling)
├── core_ingestar.go           # Coordinador principal de ingesta con control de memoria (>50MB)
├── core_manifest_hash.go      # Gestión de manifest.json, firmas MD5 y validación de cambios
├── core_memoria.go            # Gestor de memoria conversacional multi-turno (sliding window)
├── core_ollama_service.go    # Autogestión del ciclo de vida y arranque/parada del servicio local Ollama
├── core_orquestador_ranking.go  # Pipeline de filtrado Top-N
├── core_prompts.go            # Plantillas de sistema, roles y constantes del agente
├── core_prompts_helpers.go    # Formateadores y constructores de cadenas de prompts
├── core_router.go             # Enrutador híbrido (heurístico + LLM) de intenciones (RAG/Herramienta/General)
├── core_router_helpers.go     # Reglas heurísticas y filtrado léxico local del enrutador
├── core_rrf.go                # Algoritmo de fusión matemática de rankings (Reciprocal Rank Fusion)
├── core_vectorstore.go        # VectorStore en memoria, persistencia .gob y orquestación de búsqueda híbrida
├── core_vectorstore_helpers.go# Funciones auxiliares matemáticas aisladas (Similitud Coseno)
├── extractor_orquestador.go   # Despachador principal de ingesta según extensión
├── extractor_texto.go         # Extracción para archivos .txt, .md y código fuente
├── extractor_excel.go         # Extracción estructurada multipestaña para .xlsx y .xlsm
├── extractor_word.go          # Extracción de contenido estructurado para .docx
├── extractor_pdf.go           # Extracción nativa de PDF con fallback OCR WASM concurrente
├── main.go                    # Punto de entrada CLI/REPL e inicialización del entorno
├── documentos_entrada/        # Carpeta de entrada para documentos a procesar
├── mis_vectorstores/          # Almacenamiento persistente (.gob y manifest.json)
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
MODELO_EMBEDDING="nomic-embed-text"
OLLAMA_HOST="http://localhost:11434"

# Rutas configurables del sistema
# Windows local/red: "Z:/Files_IA" o "C:/MisDocumentos/IA"
# Linux/macOS: "/var/data/files_ia"
RUTA_DOCS_ORIGEN="./documentos_entrada"
RUTA_VECTORSTORES="./mis_vectorstores"

# Parámetros de Recuperación y Re-Ranking (Two-Stage Retrieval)
LIMITE_CANDIDATOS_RECUPERACION=10
LIMITE_FRAGMENTOS_FINALES=3
PESO_PUNTUACION_RRF=0.7
PESO_DENSIDAD_LEXICA=0.3
```

## Requisitos y Dependencias

Asegúrate de contar con Go 1.24+ instalado. Para descargar e instalar las dependencias externas del proyecto, ejecuta:

```bash
# Lectura y procesamiento de archivos Excel (.xlsx, .xlsm)
go get github.com/xuri/excelize/v2
# Lectura y procesamiento de archivos PDF (.pdf)
go get github.com/ledongthuc/pdf
# Iniciar el servicio y descargar el modelo de embeddings local
ollama pull nomic-embed-text
```
## 📐 Decisiones de Arquitectura y Diseño

Este proyecto fue diseñado priorizando la **portabilidad, privacidad y velocidad de ejecución nativa**, evitando la sobreingeniería y las dependencias pesadas habituales en el ecosistema de IA.

### 1. Ingesta y Embeddings 100% Locales (Soberanía de Datos)
* **Decisión:** Comunicación HTTP nativa con **Ollama** (`core_embeddings.go`) utilizando el modelo de incrustación `nomic-embed-text` en lugar de servicios en la nube o modelos obsoletos.
* **Por qué:** Garantiza privacidad absoluta (ningún dato sale de la máquina local), elimina costes de API de terceros, ofrece mayor resolución semántica y suprime la latencia de red externa.

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
* **Decisión:** Combinación de búsqueda vectorial semántica con el algoritmo de frecuencia de términos BM25 (`core_bm25.go` y `core_bm25_helpers.go`), fusionados dinámicamente mediante RRF (`core_rrf.go`).
* **Por qué:** Resuelve las limitaciones del RAG puramente vectorial en consultas que contienen nombres de funciones, códigos de error o identificadores exactos, ofreciendo máxima precisión contextual antes de invocar al modelo.

### 12. Selección Dinámica de Prompts de Sistema por Intención
* **Decisión:** Inyección condicional del prompt de sistema (`PromptSistemaBase` vs `PromptSistemaDesarrollador`) mediante `SeleccionarPromptSistema` en `core_prompts_helpers.go`, activada tras la clasificación de `RutaDesarrollo`.
* **Por qué:** Aplica las directrices estrictas de *Clean Code*, modularidad y separación de funciones auxiliares (*SRP*) únicamente en consultas de desarrollo o refactorización, evitando sobrecargar el contexto de las consultas RAG convencionales.

### 13. Re-Ranking Híbrido y Ponderación por Densidad Léxica (Top-N)
* **Decisión:** Pipeline dedicado de re-evaluación (`core_orquestador_ranking.go`) y cálculo de densidad léxica (`core_calcular_relevancia.go`) sobre los candidatos devueltos por RRF.
* **Por qué:** Filtra candidatos ruidosos recalculando la relevancia real de cada fragmento antes de inyectar el contexto al prompt del modelo, optimizando el uso de la ventana de contexto.

### 14. Control Determinista de Salida del LLM (`max_tokens`)
* **Decisión:** Inclusión explícita del parámetro `MaxTokens` en la estructura nativa `PeticionGroq` mediante la etiqueta `json:"max_tokens,omitempty"`.
* **Por qué:** Previene respuestas sobre-extendidas o truncamientos inesperados, acotando el consumo de tokens y la latencia de generación del LLM en la API de Groq sin dependencias externas.

### 15. Estimación Pre-Vuelo de Tokens Local de Cero-Latencia
* **Decisión:** Cálculo heurístico en memoria RAM (`EstimarTokensLocal`) del tamaño del prompt antes de iniciar la llamada HTTP, con capacidad de ajuste dinámico del Top-$K_2$.
* **Por qué:** Elimina la necesidad de realizar peticiones de red adicionales a APIs remotas de conteo de tokens (`CountTokens`), conservando la velocidad de respuesta nativa en una única solicitud I/O.

### 16. Gestión Automática del Ciclo de Vida del Servicio Ollama
* **Decisión:** Comprobación proactiva del socket local `http://127.0.0.1:11434` (`core_ollama_service.go`) al arrancar en `main.go`. Si está inactivo, valida la presencia de `ollama` mediante `exec.LookPath`; si existe, inicia `ollama serve` y programa su cierre automático vía `defer`. Si no existe, detiene la ejecución informando el enlace oficial de descarga.
* **Por qué:** Evita fallos de red (`dial tcp`), previene comportamientos bloqueantes por parte de antivirus al no autodescargar binarios sin permisos, ahorra memoria RAM al salir si inició el servicio y respeta las instancias que el usuario ya tenía abiertas.
* **Regla de cortesía con el entorno:** Si el agente inicia Ollama por necesidad, se encarga de destruirlo al terminar para evitar procesos huérfanos en RAM. Si Ollama ya estaba abierto previamente por el usuario, el agente lo reutiliza sin alterar su estado al cerrar.
  