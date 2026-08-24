package main

import (
	"encoding/xml"
	"sync"
)

// --- ENUMERACIONES Y TIPOS DE DOMINIO ---

type TipoRuta string

const (
	RutaRAG         TipoRuta = "rag"
	RutaHerramienta TipoRuta = "herramienta"
	RutaGeneral     TipoRuta = "general"
)

// --- CONFIGURACIÓN DE LA APLICACIÓN ---

type ConfigApp struct {
	RutaDocsOrigen   string
	RutaVectorstores string
	GroqAPIKey       string
	GroqModel        string
	GroqModelosVal   []string
}

// --- GESTIÓN, ESTADO Y CLIENTE DE MODELOS ---

// GestorModelos es la fuente única de verdad para el estado de la API y las peticiones
type GestorModelos struct {
	mu             sync.Mutex
	APIKey         string
	ModeloActivo   string
	ModelosValidos []string
	recargando     bool
}

type CatalogoModelos struct {
	Data []struct {
		ID     string `json:"id"`
		Active bool   `json:"active"`
	} `json:"data"`
}

// --- ESTRUCTURAS DE COMUNICACIÓN API GROQ ---

type Mensaje struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PeticionGroq struct {
	Model       string    `json:"model"`
	Messages    []Mensaje `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type RespuestaGroq struct {
	Choices []struct {
		Message Mensaje `json:"message"`
	} `json:"choices"`
}

// --- RAG Y ALMACENAMIENTO VECTORIAL ---

type ClasificacionDominio struct {
	Dominio     string `json:"dominio"`
	Explicacion string `json:"explicacion"`
}

type ChunkTexto struct {
	ID        string    `json:"id"`
	Origen    string    `json:"origen"`
	Contenido string    `json:"contenido"`
	Embedding []float32 `json:"embedding"`
	HashMD5   string    `json:"hash_md5"`
}

// Se usa en config_rutas_locales.go para la lectura exclusiva del sistema de archivos
type ConfigRutas struct {
	RutaDocsOrigen   string
	RutaVectorstores string
}

type TextRun struct {
	Value string `xml:",chardata"`
}

type Paragraph struct {
	Text []TextRun `xml:"r>t"`
}

type TableCell struct {
	Paragraphs []Paragraph `xml:"p"`
}

type TableRow struct {
	Cells []TableCell `xml:"tc"`
}

type Table struct {
	Rows []TableRow `xml:"tr"`
}

type DocumentBody struct {
	Paragraphs []Paragraph `xml:"p"`
	Tables     []Table     `xml:"tbl"`
}

type Document struct {
	XMLName xml.Name     `xml:"document"`
	Body    DocumentBody `xml:"body"`
}

// GestorMemoria administra el historial conversacional con ventana deslizante thread-safe.
type GestorMemoria struct {
	mu       sync.RWMutex
	mensajes []Mensaje
	limite   int
}

// ResultadoHerramienta representa la respuesta de la ejecución de una herramienta nativa.
type ResultadoHerramienta struct {
	Nombre    string `json:"nombre"`
	Exito     bool   `json:"exito"`
	Contenido string `json:"contenido"`
	Error     string `json:"error,omitempty"`
}

// ResultadoHibrido combina métricas de similitud vectorial, léxica (BM25) y puntuación RRF.
type ResultadoHibrido struct {
	Documento   DocumentoVectorial
	ScoreVector float64
	ScoreBM25   float64
	ScoreRRF    float64
}

// IndiceBM25 gestiona la recuperación léxica exacta por términos de frecuencia.
type IndiceBM25 struct {
	IDF       map[string]float64
	DocLen    map[int]float64
	AvgDocLen float64
	TotalDocs int
	TermFreqs []map[string]int
	K1        float64
	B         float64
}
