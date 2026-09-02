package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"strings"
)

// NuevoVectorStore inicializa la base vectorial local.
func NuevoVectorStore(rutaDisco string) *VectorStore {
	return &VectorStore{
		RutaDisco:  rutaDisco,
		Documentos: []DocumentoVectorial{},
	}
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

// buscarSimilaresSinLock evalúa similitud sin adquirir mutex (asume que el llamador ya tiene el bloqueo RLock).
func (vs *VectorStore) buscarSimilaresSinLock(queryVector []float32, topK int) []ResultadoBusqueda {
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

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Similitud > resultados[j].Similitud
	})

	if len(resultados) > topK {
		return resultados[:topK]
	}

	return resultados
}

// BuscarSimilares evalúa el vector de consulta contra toda la base y retorna el Top-K con mayor relevancia.
func (vs *VectorStore) BuscarSimilares(queryVector []float32, topK int) []ResultadoBusqueda {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.buscarSimilaresSinLock(queryVector, topK)
}

// BuscarSimilaresConFiltro evalúa solo los vectores cuyo Origen contenga el filtro especificado (ej: "Go\", "Python\").
func (vs *VectorStore) BuscarSimilaresConFiltro(queryVector []float32, topK int, filtroOrigen string) []ResultadoBusqueda {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.Documentos) == 0 {
		return nil
	}

	var resultados []ResultadoBusqueda
	for _, doc := range vs.Documentos {
		if filtroOrigen != "" && !strings.Contains(strings.ToLower(doc.Chunk.Origen), strings.ToLower(filtroOrigen)) {
			continue
		}

		sim := SimilitudCoseno(queryVector, doc.Embedding)
		resultados = append(resultados, ResultadoBusqueda{
			Documento: doc,
			Similitud: sim,
		})
	}

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Similitud > resultados[j].Similitud
	})

	if len(resultados) > topK {
		return resultados[:topK]
	}

	return resultados
}

// BuscarHibrido ejecuta la búsqueda por Similitud Coseno y BM25, fusionando sus rankings mediante RRF.
func (vs *VectorStore) BuscarHibrido(pregunta string, queryVector []float32, topK int) []ResultadoHibrido {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.Documentos) == 0 {
		return nil
	}

	// 1. Recuperación vectorial amplia (usando método interno seguro)
	resVector := vs.buscarSimilaresSinLock(queryVector, topK*2)

	// 2. Recuperación léxica por palabras clave (BM25)
	indiceBM25 := NuevoIndiceBM25(vs.Documentos)
	resBM25 := indiceBM25.BuscarBM25(pregunta, vs.Documentos, topK*2)

	// 3. Fusión matemática de rankings (RRF)
	return FusionarResultadosRRF(resVector, resBM25, topK)
}
