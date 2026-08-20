package main

import "sync"

type ConfigApp struct {
	RutaDocsOrigen   string
	RutaVectorstores string
	GroqAPIKey       string
	GroqModel        string
	GroqModelosVal   []string // Convertido a slice para manejo directo
}

// GestorModelos administra el estado dinámico, el bloqueo por hilos y la rotación
type GestorModelos struct {
	mu             sync.Mutex
	ModelosValidos []string
	apiKey         string
	recargando     bool
}

// Mensaje representa el formato estándar de conversación
type Mensaje struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PeticionGroq define el cuerpo enviando a la API de Groq
type PeticionGroq struct {
	Model       string    `json:"model"`
	Messages    []Mensaje `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// RespuestaGroq procesa la devolución de Groq
type RespuestaGroq struct {
	Choices []struct {
		Message Mensaje `json:"message"`
	} `json:"choices"`
}

// ClasificacionDominio captura la decisión tomada por el router RAG
type ClasificacionDominio struct {
	Dominio     string `json:"dominio"`
	Explicacion string `json:"explicacion"`
}

// ChunkTexto representa el fragmento de texto procesado con su vector
type ChunkTexto struct {
	ID        string    `json:"id"`
	Origen    string    `json:"origen"`
	Contenido string    `json:"contenido"`
	Embedding []float32 `json:"embedding"`
	HashMD5   string    `json:"hash_md5"`
}

// CatalogoModelos mapea la respuesta del endpoint de modelos de Groq
type CatalogoModelos struct {
	Data []struct {
		ID     string `json:"id"`
		Active bool   `json:"active"`
	} `json:"data"`
}

type ConfigRutas struct {
	RutaDocsOrigen   string
	RutaVectorstores string
}
