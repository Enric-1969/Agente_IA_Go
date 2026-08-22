package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
)

// DocumentoVectorial asocia el chunk de texto con su representación en embedding.
type DocumentoVectorial struct {
	Chunk     ChunkTexto
	Embedding []float32
}

// ResultadoBusqueda almacena la coincidencia encontrada y su puntuación de similitud.
type ResultadoBusqueda struct {
	Documento DocumentoVectorial
	Similitud float32
}

// VectorStore maneja el almacenamiento local en memoria y la persistencia en disco (.gob).
type VectorStore struct {
	mu         sync.RWMutex
	Documentos []DocumentoVectorial
	RutaDisco  string
}

// NuevoVectorStore inicializa la base vectorial local.
func NuevoVectorStore(rutaDisco string) *VectorStore {
	return &VectorStore{
		RutaDisco:  rutaDisco,
		Documentos: []DocumentoVectorial{},
	}
}

// SimilitudCoseno calcula el ángulo de coincidencia entre dos vectores de igual dimensión.
func SimilitudCoseno(v1, v2 []float32) float32 {
	if len(v1) != len(v2) || len(v1) == 0 {
		return 0.0
	}

	var productoEscalar, normaA, normaB float64
	for i := range v1 {
		productoEscalar += float64(v1[i]) * float64(v2[i])
		normaA += float64(v1[i]) * float64(v1[i])
		normaB += float64(v2[i]) * float64(v2[i])
	}

	if normaA == 0 || normaB == 0 {
		return 0.0
	}

	return float32(productoEscalar / (math.Sqrt(normaA) * math.Sqrt(normaB)))
}

// AgregarVector añade un nuevo documento vectorial al almacenamiento en memoria.
func (vs *VectorStore) AgregarVector(chunk ChunkTexto, embedding []float32) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	vs.Documentos = append(vs.Documentos, DocumentoVectorial{
		Chunk:     chunk,
		Embedding: embedding,
	})
}

// GuardarPersistencia exporta los vectores a un archivo binario .gob en disco.
func (vs *VectorStore) GuardarPersistencia() error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	file, err := os.Create(vs.RutaDisco)
	if err != nil {
		return fmt.Errorf("error al crear el archivo de base vectorial: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(vs.Documentos); err != nil {
		return fmt.Errorf("error al serializar vectorstore en .gob: %w", err)
	}

	return nil
}

// CargarPersistencia lee el archivo .gob e hidrata la base de datos vectorial.
func (vs *VectorStore) CargarPersistencia() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	file, err := os.Open(vs.RutaDisco)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Si el archivo aún no existe, iniciamos con la base vacía
		}
		return fmt.Errorf("error al abrir el archivo vectorstore: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&vs.Documentos); err != nil {
		return fmt.Errorf("error al deserializar vectorstore desde .gob: %w", err)
	}

	return nil
}

// BuscarSimilares evalúa el vector de consulta contra toda la base y retorna el Top-K con mayor relevancia.
func (vs *VectorStore) BuscarSimilares(queryVector []float32, topK int) []ResultadoBusqueda {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.Documentos) == 0 {
		return nil
	}

	var resultados []ResultadoBusqueda
	for _, doc := range vs.Documentos {
		sim := SimilitudCoseno(queryVector, doc.Embedding)
		resultados = append(resultados, ResultadoBusqueda{
			Documento: doc,
			Similitud: sim,
		})
	}

	// Ordenamiento descendente por la mejor similitud
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Similitud > resultados[j].Similitud
	})

	if len(resultados) > topK {
		return resultados[:topK]
	}

	return resultados
}
