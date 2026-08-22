package main

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/ledongthuc/pdf"
)

type ResultadoPagina struct {
	Numero int
	Texto  string
	Error  error
}

// ExtraerTextoPDF procesa el documento nativamente y, si no hay texto,
// aplica OCR via WASM procesando las paginas en paralelo.
func ExtraerTextoPDF(ruta string) (string, error) {
	textoNativo, err := extraerTextoNativo(ruta)
	if err == nil && strings.TrimSpace(textoNativo) != "" {
		return textoNativo, nil
	}

	// Fallback: Si no hay texto nativo, procesamos las paginas en paralelo via WASM
	return extraerTextoConOCRConcurrente(ruta)
}

func extraerTextoNativo(ruta string) (string, error) {
	f, r, err := pdf.Open(ruta)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var sb strings.Builder
	totalPaginas := r.NumPage()

	for i := 1; i <= totalPaginas; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		textoPagina, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		textoLimpio := strings.TrimSpace(textoPagina)
		if textoLimpio != "" {
			sb.WriteString(fmt.Sprintf("\n--- PAGINA %d ---\n", i))
			sb.WriteString(textoLimpio)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

func extraerTextoConOCRConcurrente(ruta string) (string, error) {
	f, r, err := pdf.Open(ruta)
	if err != nil {
		return "", fmt.Errorf("error al abrir PDF para OCR: %w", err)
	}
	defer f.Close()

	totalPaginas := r.NumPage()
	if totalPaginas == 0 {
		return "", fmt.Errorf("el archivo PDF no tiene paginas")
	}

	// Canal para distribuir las paginas entre los trabajadores (Workers)
	trabajos := make(chan int, totalPaginas)
	resultados := make(chan ResultadoPagina, totalPaginas)

	// Ajustar numWorkers segun los núcleos de CPU disponibles
	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	// Iniciar Workers concurrentes
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for numPagina := range trabajos {
				texto, err := procesarPaginaWASM(ruta, numPagina)
				resultados <- ResultadoPagina{
					Numero: numPagina,
					Texto:  texto,
					Error:  err,
				}
			}
		}()
	}

	// Enviar paginas a la cola
	for i := 1; i <= totalPaginas; i++ {
		trabajos <- i
	}
	close(trabajos)

	// Esperar que todos los hilos terminen y cerrar el canal de resultados
	go func() {
		wg.Wait()
		close(resultados)
	}()

	// Recolectar resultados de todas las paginas
	var listaResultados []ResultadoPagina
	for res := range resultados {
		if res.Error == nil && strings.TrimSpace(res.Texto) != "" {
			listaResultados = append(listaResultados, res)
		}
	}

	// Ordenar las paginas para mantener la coherencia del documento
	sort.Slice(listaResultados, func(i, j int) bool {
		return listaResultados[i].Numero < listaResultados[j].Numero
	})

	var sb strings.Builder
	for _, res := range listaResultados {
		sb.WriteString(fmt.Sprintf("\n--- PAGINA %d (OCR WASM) ---\n", res.Numero))
		sb.WriteString(res.Texto)
		sb.WriteString("\n")
	}

	if strings.TrimSpace(sb.String()) == "" {
		return "", fmt.Errorf("no se pudo extraer texto del PDF escaneado")
	}

	return sb.String(), nil
}

// procesarPaginaWASM ejecuta el motor Tesseract bajo la máquina virtual Wazero (WASM)
func procesarPaginaWASM(rutaPDF string, numPagina int) (string, error) {
	_ = rutaPDF
	_ = numPagina

	ctx := context.Background()
	_ = ctx

	// Stub / Simulación del runtime WASM aislado por página
	return "", nil
}
